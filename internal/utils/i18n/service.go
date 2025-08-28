package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
)

// 全局消息缓存 - 启动时一次性加载，运行时只读
// 使用全局变量避免过度封装，提升性能
var (
	globalMessages = make(map[string]string)
	globalMessagesMutex sync.RWMutex
)

// BuildCacheKey 构建缓存键
// 统一将locale和key都转换为小写，确保大小写一致性
func BuildCacheKey(locale, key string) string {
	return strings.ToLower(locale + ":" + key)
}

// SetMessage 设置消息到全局缓存（仅在初始化时使用）
func SetMessage(cacheKey, message string) {
	globalMessagesMutex.Lock()
	defer globalMessagesMutex.Unlock()
	globalMessages[cacheKey] = message
}

// GetMessage 获取消息（高性能版本）
func GetMessage(cacheKey string) (string, bool) {
	globalMessagesMutex.RLock()
	defer globalMessagesMutex.RUnlock()
	message, exists := globalMessages[cacheKey]
	return message, exists
}

// ClearMessages 清空消息缓存（仅在重新加载时使用）
func ClearMessages() {
	globalMessagesMutex.Lock()
	defer globalMessagesMutex.Unlock()
	globalMessages = make(map[string]string)
}

// GetMessageCount 获取消息数量
func GetMessageCount() int {
	globalMessagesMutex.RLock()
	defer globalMessagesMutex.RUnlock()
	return len(globalMessages)
}

// GetMessagesByPrefix 根据前缀获取消息（用于调试）
func GetMessagesByPrefix(prefix string) map[string]string {
	globalMessagesMutex.RLock()
	defer globalMessagesMutex.RUnlock()
	result := make(map[string]string)
	for key, value := range globalMessages {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			result[key] = value
		}
	}
	return result
}

// LocalizationManager 本地化管理器 - 基于配置的预编译版本
type LocalizationManager struct {
	cfg *config.Config
}

// NewLocalizationManager 创建本地化管理器
func NewLocalizationManager(cfg *config.Config) I18nManager {
	manager := &LocalizationManager{
		cfg: cfg,
	}

	// 预加载所有消息
	if err := manager.preloadAllMessages(); err != nil {
		slog.Error("预加载消息失败", "error", err)
		// 降级处理：创建空管理器但不返回错误
	}

	slog.Info("本地化管理器初始化完成",
		"默认语言", manager.GetDefaultLanguage(),
		"支持语言数", len(manager.GetSupportedLanguageCodes()),
		"消息数量", GetMessageCount(),
		"全局消息缓存大小", len(globalMessages))

	return manager
}

// preloadAllMessages 预加载所有语言的所有消息
func (m *LocalizationManager) preloadAllMessages() error {
	return m.loadAllMessages()
}

// loadAllMessages 加载所有语言的消息文件
func (m *LocalizationManager) loadAllMessages() error {
	localesDir := m.cfg.App.LocalesDir
	if localesDir == "" {
		return fmt.Errorf("locales directory not configured")
	}

	// 检查目录是否存在
	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		return fmt.Errorf("locales directory does not exist: %s", localesDir)
	}

	// 遍历所有JSON文件
	return filepath.WalkDir(localesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理JSON文件
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}

		// 从文件名提取语言代码
		langCode := strings.TrimSuffix(d.Name(), ".json")
		langCode = m.normalizeLocale(langCode)

		// 加载文件内容
		if err := m.loadLanguageFile(path, langCode); err != nil {
			slog.Error("加载语言文件失败", "file", path, "error", err)
			// 继续处理其他文件
		}

		return nil
	})
}

// loadLanguageFile 加载单个语言文件
func (m *LocalizationManager) loadLanguageFile(filePath, langCode string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	var messages map[string]interface{}
	if err := json.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	// 扁平化并存储消息
	flattened := m.flattenMessages(messages, "")
	for key, value := range flattened {
		cacheKey := BuildCacheKey(langCode, key)
		SetMessage(cacheKey, value)
	}

	slog.Debug("语言文件加载完成", "file", filePath, "language", langCode, "messages", len(flattened))
	return nil
}

// flattenMessages 扁平化嵌套的消息结构
func (m *LocalizationManager) flattenMessages(messages map[string]interface{}, prefix string) map[string]string {
	result := make(map[string]string)

	for key, value := range messages {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			// 递归处理嵌套对象
			nested := m.flattenMessages(v, fullKey)
			for nestedKey, nestedValue := range nested {
				result[nestedKey] = nestedValue
			}
		default:
			// 如果不是字符串也不是嵌套对象，则忽略或按需处理
			// 目前，我们只关心字符串值作为最终消息
		}
	}

	return result
}


// normalizeLocale 规范化语言代码
func (m *LocalizationManager) normalizeLocale(locale string) string {
	if locale == "" {
		return m.GetDefaultLanguage()
	}

	// 转换为小写
	locale = strings.ToLower(locale)

	// 处理常见的语言代码映射
	switch locale {
	case "zh", "zh-hans", "zh_hans", "zh-cn", "zh_cn":
		return "zh-cn"
	case "zh-hant", "zh_hant", "zh-tw", "zh_tw", "zh-hk", "zh_hk":
		return "zh-tw"
	case "en", "en-us", "en_us":
		return "en"
	case "ja", "ja-jp", "ja_jp":
		return "ja-jp"
	case "ko", "ko-kr", "ko_kr":
		return "ko-kr"
	case "fr", "fr-fr", "fr_fr":
		return "fr-fr"
	case "de", "de-de", "de_de":
		return "de-de"
	case "es", "es-es", "es_es":
		return "es-es"
	case "pt", "pt-br", "pt_br":
		return "pt-br"
	case "ru", "ru-ru", "ru_ru":
		return "ru-ru"
	case "ar", "ar-sa", "ar_sa":
		return "ar-sa"
	case "hi", "hi-in", "hi_in":
		return "hi-in"
	default:
		return locale
	}
}

// GetMessage 获取本地化消息 - 超高性能版本
func (m *LocalizationManager) GetMessage(key, locale string) string {
	// 规范化语言代码
	normalizedLocale := m.normalizeLocale(locale)
	cacheKey := BuildCacheKey(normalizedLocale, key)

	// 高效的缓存查找
	if message, exists := GetMessage(cacheKey); exists {
		return message
	}

	// 尝试默认语言
	defaultLang := m.GetDefaultLanguage()
	if normalizedLocale != defaultLang {
		defaultCacheKey := BuildCacheKey(defaultLang, key)
		if defaultMessage, exists := GetMessage(defaultCacheKey); exists {
			slog.Debug("使用默认语言消息", "key", key, "requested_locale", locale, "default_locale", defaultLang)
			return defaultMessage
		}
	}

	// 最终降级：返回原始key
	slog.Warn("消息未找到", "key", key, "locale", locale, "normalized_locale", normalizedLocale, "cache_key", cacheKey)
	return key
}

// GetMessageWithParams 获取带参数的翻译消息，支持 {{参数名}} 格式的占位符替换
func (m *LocalizationManager) GetMessageWithParams(key, locale string, params map[string]string) string {
	message := m.GetMessage(key, locale)

	for paramKey, paramValue := range params {
		placeholder := "{{" + paramKey + "}}"
		message = strings.ReplaceAll(message, placeholder, paramValue)
	}

	return message
}

// GetSupportedLanguages 获取支持的语言列表
func (m *LocalizationManager) GetSupportedLanguages() []config.SupportLang {
	return m.cfg.SupportLangs
}

// IsLanguageSupported 判断输入语言是否支持
func (m *LocalizationManager) IsLanguageSupported(locale string) bool {
	normalizedLocale := m.normalizeLocale(locale)
	for _, lang := range m.cfg.SupportLangs {
		if lang.Code == normalizedLocale {
			return true
		}
	}
	return false
}

// GetSupportedLanguageCodes 获取支持的语言代码列表
func (m *LocalizationManager) GetSupportedLanguageCodes() []string {
	codes := make([]string, len(m.cfg.SupportLangs))
	for i, lang := range m.cfg.SupportLangs {
		codes[i] = lang.Code
	}
	return codes
}

// GetDefaultLanguage 获取默认语言
func (m *LocalizationManager) GetDefaultLanguage() string {
	return m.cfg.App.DefaultLang
}

// GetLanguageInfo 获取语言信息
func (m *LocalizationManager) GetLanguageInfo(locale string) (config.SupportLang, bool) {
	normalizedLocale := m.normalizeLocale(locale)
	for _, lang := range m.cfg.SupportLangs {
		if lang.Code == normalizedLocale {
			return lang, true
		}
	}
	return config.SupportLang{}, false
}

// GetMessageCount 获取已加载的消息数量（用于调试）
func (m *LocalizationManager) GetMessageCount() int {
	return GetMessageCount()
}

// GetMessagesByLanguage 获取指定语言的所有消息（用于调试）
func (m *LocalizationManager) GetMessagesByLanguage(locale string) map[string]string {
	normalized := m.normalizeLocale(locale)
	prefix := BuildCacheKey(normalized, "")
	return GetMessagesByPrefix(prefix)
}

// GetLocaleFromContext 从 Gin 上下文中智能检测语言
// 优先级：用户设置 > Cookie > 浏览器偏好 > 系统默认
func (m *LocalizationManager) GetLocaleFromContext(c *gin.Context) string {
	// 优先级1：用户设置 (如果已登录)
	if locale := m.getUserLocaleFromContext(c); locale != "" {
		return locale
	}

	// 优先级2：Cookie 设置
	if cookieLocale, err := c.Cookie("lang"); err == nil && cookieLocale != "" {
		if m.IsLanguageSupported(cookieLocale) {
			slog.Debug("从 Cookie 获取语言", "locale", cookieLocale)
			return cookieLocale
		}
	}

	// 优先级3：浏览器偏好
	if browserLocale := m.getBrowserPreferredLocale(c); browserLocale != "" {
		slog.Debug("从浏览器偏好获取语言", "locale", browserLocale)
		return browserLocale
	}

	// 优先级4：系统默认
	defaultLocale := m.GetDefaultLanguage()
	slog.Debug("使用系统默认语言", "locale", defaultLocale)
	return defaultLocale
}

// getUserLocaleFromContext 从上下文中获取已认证用户的语言设置
func (m *LocalizationManager) getUserLocaleFromContext(c *gin.Context) string {
	isLogin, exists := c.Get("IsLogin")
	if !exists || isLogin != true {
		return ""
	}

	authUserRaw, exists := c.Get("AuthUser")
	if !exists || authUserRaw == nil {
		return ""
	}

	authUser, ok := authUserRaw.(*types.AuthUser)
	if !ok || authUser.Lang == "" || !m.IsLanguageSupported(authUser.Lang) {
		return ""
	}

	slog.Debug("从用户设置获取语言", "locale", authUser.Lang)
	return authUser.Lang
}

// getBrowserPreferredLocale 从 Accept-Language 头中提取最佳匹配的语言
func (m *LocalizationManager) getBrowserPreferredLocale(c *gin.Context) string {
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang == "" {
		return ""
	}

	languages := m.parseAcceptLanguage(acceptLang)

	for _, lang := range languages {
		// 尝试规范化匹配
		if normalized := m.normalizeLocale(lang); normalized != "" && m.IsLanguageSupported(normalized) {
			return normalized
		}
	}

	return ""
}

// parseAcceptLanguage 解析 Accept-Language 头
func (m *LocalizationManager) parseAcceptLanguage(acceptLang string) []string {
	// 简化实现，实际项目中可能需要更复杂的解析
	var languages []string
	parts := strings.Split(acceptLang, ",")

	for _, part := range parts {
		// 移除质量值 (如 ";q=0.9")
		lang := strings.TrimSpace(strings.Split(part, ";")[0])
		if lang != "" {
			languages = append(languages, lang)
		}
	}

	return languages
}

// ReloadMessages 重新加载消息（开发环境使用）
func (m *LocalizationManager) ReloadMessages() error {
	ClearMessages()
	return m.preloadAllMessages()
}