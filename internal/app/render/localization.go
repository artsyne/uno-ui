package render

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// LocalizationManager 管理用于国际化的、特定于地域的消息。
type LocalizationManager struct {
	cfg      *config.Config
	messages map[string]map[string]string
}

// NewLocalizationManager 创建并初始化本地化管理器，预加载所有支持的地域语言文件。
func NewLocalizationManager(cfg *config.Config) (*LocalizationManager, error) {
	manager := &LocalizationManager{
		cfg:      cfg,
		messages: make(map[string]map[string]string),
	}

	var loadedCount int
	for _, supportedLocale := range cfg.SupportLangs {
		if err := manager.loadLocale(supportedLocale.Code); err != nil {
			slog.Error("加载语言文件失败", "locale", supportedLocale.Code, "error", err)
			continue
		}
		loadedCount++
	}

	// 确保至少加载了一个语言文件
	if loadedCount == 0 {
		slog.Error("初始化本地化管理器失败", "error", "没有成功加载任何语言文件")
		return nil, fmt.Errorf("没有成功加载任何语言文件，请检查语言文件是否存在")
	}

	// 确保默认语言已加载
	if !manager.isValidLocale(cfg.App.DefaultLang) {
		slog.Warn("默认语言未成功加载", "default_lang", cfg.App.DefaultLang, "loaded_locales", manager.GetSupportedLocales())
	}

	slog.Info("本地化管理器初始化完成", "loaded_count", loadedCount, "total_configured", len(cfg.SupportLangs))
	return manager, nil
}

// loadLocale 从磁盘加载单个地域语言文件。
func (m *LocalizationManager) loadLocale(locale string) error {
	localeFile := filepath.Join(m.cfg.App.LocalesDir, locale+".json")

	data, err := os.ReadFile(localeFile)
	if err != nil {
		slog.Error("无法读取地域文件", "file", localeFile, "locale", locale, "error", err)
		return fmt.Errorf("无法读取地域文件 %s: %w", localeFile, err)
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		slog.Error("无法解析地域文件", "file", localeFile, "locale", locale, "error", err)
		return fmt.Errorf("无法解析地域文件 %s: %w", localeFile, err)
	}

	m.messages[locale] = messages
	slog.Info("成功加载地域文件", "locale", locale, "messages_count", len(messages))
	return nil
}

// GetMessage 获取指定键和地域的消息，如果找不到则回退到默认地域
func (m *LocalizationManager) GetMessage(key, locale string) string {
	if messages, exists := m.messages[locale]; exists {
		if message, found := messages[key]; found {
			return message
		}
	}

	// 回退到默认语言
	defaultLocale := m.cfg.App.DefaultLang
	if locale != defaultLocale {
		if messages, exists := m.messages[defaultLocale]; exists {
			if message, found := messages[key]; found {
				return message
			}
		}
	}

	// 如果都找不到，返回键本身
	return key
}

// GetMessageWithParams 获取带参数的翻译消息，支持 {{参数名}} 格式的占位符替换
func (m *LocalizationManager) GetMessageWithParams(key string, locale string, params map[string]string) string {
	message := m.GetMessage(key, locale)

	for paramKey, paramValue := range params {
		placeholder := "{{" + paramKey + "}}"
		message = strings.ReplaceAll(message, placeholder, paramValue)
	}

	return message
}

// GetSupportedLocales 返回所有已加载的地域代码列表
func (m *LocalizationManager) GetSupportedLocales() []string {
	locales := make([]string, 0, len(m.messages))
	for locale := range m.messages {
		locales = append(locales, locale)
	}
	return locales
}

// isValidLocale 检查语言是否被支持
func (m *LocalizationManager) isValidLocale(locale string) bool {
	_, exists := m.messages[locale]
	return exists
}

// GetUserLocale 获取用户的首选地域，如果未设置则回退到默认值
func (m *LocalizationManager) GetUserLocale(user RenderUser) string {
	if user != nil {
		if locale := user.GetLang(); locale != "" {
			return locale
		}
	}
	return m.cfg.App.DefaultLang
}

// GetLocaleFromContext 从 Gin 上下文中智能检测语言
// 优先级：用户设置 > Cookie > 浏览器偏好 > 系统默认
func (m *LocalizationManager) GetLocaleFromContext(c *gin.Context) string {
	// 优先级1：用户设置 (如果已登录)
	if isLogin, exists := c.Get("IsLogin"); exists && isLogin == true {
		if authUserValue, exists := c.Get("AuthUser"); exists && authUserValue != nil {
			if user, ok := authUserValue.(RenderUser); ok {
				if locale := user.GetLang(); locale != "" && m.isValidLocale(locale) {
					return locale
				}
			}
		}
	}

	// 优先级2：Cookie 设置
	if cookieLocale, err := c.Cookie("lang"); err == nil && cookieLocale != "" {
		if m.isValidLocale(cookieLocale) {
			return cookieLocale
		}
	}

	// 优先级3：浏览器偏好 (Accept-Language)
	if browserLocale := m.getBrowserPreferredLocale(c); browserLocale != "" {
		return browserLocale
	}

	// 优先级4：系统默认
	return m.cfg.App.DefaultLang
}

// getBrowserPreferredLocale 从 Accept-Language 头中提取最佳匹配的语言
func (m *LocalizationManager) getBrowserPreferredLocale(c *gin.Context) string {
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang == "" {
		return ""
	}

	languages := m.parseAcceptLanguage(acceptLang)

	for _, lang := range languages {
		// 直接匹配
		if m.isValidLocale(lang) {
			return lang
		}

		// 尝试标准化匹配 (如 en-US -> en, zh-CN -> zh-cn)
		if normalized := m.normalizeLocale(lang); normalized != "" && m.isValidLocale(normalized) {
			return normalized
		}
	}

	return ""
}

// normalizeLocale 标准化语言代码格式
func (m *LocalizationManager) normalizeLocale(locale string) string {
	parts := strings.Split(strings.ToLower(locale), "-")
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) >= 2 {
		// 处理特殊情况：zh-CN -> zh-cn
		return parts[0] + "-" + parts[1]
	}
	return ""
}

// parseAcceptLanguage 解析 Accept-Language 头部内容
func (m *LocalizationManager) parseAcceptLanguage(acceptLang string) []string {
	var languages []string
	parts := strings.Split(acceptLang, ",")

	for _, part := range parts {
		lang := strings.TrimSpace(part)
		// 移除质量因子 (如 ;q=0.9)
		if idx := strings.Index(lang, ";"); idx != -1 {
			lang = lang[:idx]
		}
		lang = strings.TrimSpace(lang)
		if lang != "" {
			languages = append(languages, lang)
		}
	}

	return languages
}
