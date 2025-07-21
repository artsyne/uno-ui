package localization

import (
	"log/slog"
	"strings"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// 全局消息缓存 - 启动时一次性加载，运行时只读
// 使用全局变量避免过度封装，提升性能
var globalMessages = make(map[string]string)

// BuildCacheKey 构建缓存键
func BuildCacheKey(locale, key string) string {
	return locale + ":" + key
}

// SetMessage 设置消息到全局缓存（仅在初始化时使用）
func SetMessage(cacheKey, message string) {
	globalMessages[cacheKey] = message
}

// GetMessage 获取消息（高性能版本）
func GetMessage(cacheKey string) (string, bool) {
	message, exists := globalMessages[cacheKey]
	return message, exists
}

// ClearMessages 清空消息缓存（仅在重新加载时使用）
func ClearMessages() {
	globalMessages = make(map[string]string)
}

// GetMessageCount 获取消息数量
func GetMessageCount() int {
	return len(globalMessages)
}

// GetMessagesByPrefix 根据前缀获取消息（用于调试）
func GetMessagesByPrefix(prefix string) map[string]string {
	result := make(map[string]string)
	for key, value := range globalMessages {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			result[key] = value
		}
	}
	return result
}

// ManagerInterface 本地化管理器接口
type ManagerInterface interface {
	GetMessage(key, locale string) string
	GetMessageWithParams(key, locale string, params map[string]string) string
	GetSupportedLanguages() []config.SupportLang
	GetSupportedLanguageCodes() []string
	GetDefaultLanguage() string
	IsLanguageSupported(locale string) bool
	GetLanguageInfo(locale string) (config.SupportLang, bool)
	GetMessageCount() int
	GetMessagesByLanguage(locale string) map[string]string
	GetLocaleFromContext(c *gin.Context) string
	ReloadMessages() error
}

// Manager 本地化管理器 - 基于配置的预编译版本
type Manager struct {
	cfg        *config.Config
	loader     *Loader
	normalizer *Normalizer
}

// NewManager 创建本地化管理器
func NewManager(cfg *config.Config) ManagerInterface {
	normalizer := NewNormalizer(cfg)
	loader := NewLoader(cfg, normalizer)

	manager := &Manager{
		cfg:        cfg,
		loader:     loader,
		normalizer: normalizer,
	}

	// 预加载所有消息
	if err := manager.preloadAllMessages(); err != nil {
		slog.Error("预加载消息失败", "error", err)
		// 降级处理：创建空管理器但不返回错误
	}

	slog.Info("本地化管理器初始化完成",
		"默认语言", manager.normalizer.GetDefaultLanguage(),
		"支持语言数", len(manager.normalizer.GetSupportedLanguageCodes()),
		"消息数量", GetMessageCount())

	return manager
} // preloadAllMessages 预加载所有语言的所有消息
func (m *Manager) preloadAllMessages() error {
	return m.loader.LoadAllMessages()
}

// GetMessage 获取本地化消息 - 超高性能版本
func (m *Manager) GetMessage(key, locale string) string {
	// 规范化语言代码
	normalizedLocale := m.normalizer.NormalizeLocale(locale)
	cacheKey := BuildCacheKey(normalizedLocale, key)

	// 高效的缓存查找
	if message, exists := GetMessage(cacheKey); exists {
		return message
	}

	// 尝试默认语言
	defaultLang := m.normalizer.GetDefaultLanguage()
	if normalizedLocale != defaultLang {
		defaultCacheKey := BuildCacheKey(defaultLang, key)
		if defaultMessage, exists := GetMessage(defaultCacheKey); exists {
			slog.Debug("使用默认语言消息", "key", key, "requested_locale", locale, "default_locale", defaultLang)
			return defaultMessage
		}
	}

	// 最终降级：返回原始key
	slog.Warn("消息未找到", "key", key, "locale", locale)
	return key
}

// GetMessageWithParams 获取带参数的翻译消息，支持 {{参数名}} 格式的占位符替换
func (m *Manager) GetMessageWithParams(key, locale string, params map[string]string) string {
	message := m.GetMessage(key, locale)

	for paramKey, paramValue := range params {
		placeholder := "{{" + paramKey + "}}"
		message = strings.ReplaceAll(message, placeholder, paramValue)
	}

	return message
}

// GetSupportedLanguages 获取支持的语言列表
func (m *Manager) GetSupportedLanguages() []config.SupportLang {
	return m.normalizer.GetSupportedLanguages()
}

// GetSupportedLanguageCodes 获取支持的语言代码列表
func (m *Manager) GetSupportedLanguageCodes() []string {
	return m.normalizer.GetSupportedLanguageCodes()
}

// GetDefaultLanguage 获取默认语言
func (m *Manager) GetDefaultLanguage() string {
	return m.normalizer.GetDefaultLanguage()
}

// IsLanguageSupported 检查语言是否支持
func (m *Manager) IsLanguageSupported(locale string) bool {
	return m.normalizer.IsLanguageSupported(locale)
}

// GetLanguageInfo 获取语言信息
func (m *Manager) GetLanguageInfo(locale string) (config.SupportLang, bool) {
	return m.normalizer.GetLanguageInfo(locale)
}

// GetMessageCount 获取已加载的消息数量（用于调试）
func (m *Manager) GetMessageCount() int {
	return GetMessageCount()
}

// GetMessagesByLanguage 获取指定语言的所有消息（用于调试）
func (m *Manager) GetMessagesByLanguage(locale string) map[string]string {
	normalized := m.normalizer.NormalizeLocale(locale)
	prefix := BuildCacheKey(normalized, "")
	return GetMessagesByPrefix(prefix)
}

// GetLocaleFromContext 从 Gin 上下文中智能检测语言
// 优先级：用户设置 > Cookie > 浏览器偏好 > 系统默认
func (m *Manager) GetLocaleFromContext(c *gin.Context) string {
	// 优先级1：用户设置 (如果已登录)
	if isLogin, exists := c.Get("IsLogin"); exists && isLogin == true {
		if authUserValue, exists := c.Get("AuthUser"); exists && authUserValue != nil {
			// 这里使用接口避免循环依赖
			if user, ok := authUserValue.(interface{ GetLang() string }); ok {
				if locale := user.GetLang(); locale != "" && m.IsLanguageSupported(locale) {
					return locale
				}
			}
		}
	}

	// 优先级2：Cookie 设置
	if cookieLocale, err := c.Cookie("lang"); err == nil && cookieLocale != "" {
		if m.IsLanguageSupported(cookieLocale) {
			return cookieLocale
		}
	}

	// 优先级3：浏览器偏好 (Accept-Language)
	if browserLocale := m.getBrowserPreferredLocale(c); browserLocale != "" {
		return browserLocale
	}

	// 优先级4：系统默认
	return m.GetDefaultLanguage()
}

// getBrowserPreferredLocale 从 Accept-Language 头中提取最佳匹配的语言
func (m *Manager) getBrowserPreferredLocale(c *gin.Context) string {
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang == "" {
		return ""
	}

	languages := m.parseAcceptLanguage(acceptLang)

	for _, lang := range languages {
		// 直接匹配
		if m.IsLanguageSupported(lang) {
			return lang
		}

		// 尝试规范化匹配
		if normalized := m.normalizer.NormalizeLocale(lang); normalized != "" && m.IsLanguageSupported(normalized) {
			return normalized
		}
	}

	return ""
}

// parseAcceptLanguage 解析 Accept-Language 头
func (m *Manager) parseAcceptLanguage(acceptLang string) []string {
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
func (m *Manager) ReloadMessages() error {
	ClearMessages()
	return m.preloadAllMessages()
}
