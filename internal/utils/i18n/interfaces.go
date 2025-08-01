package i18n

import (
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
)

// I18nManager 国际化管理器接口
// 提供完整的国际化功能，包括消息获取、语言检测和管理
type I18nManager interface {
	// GetMessage 获取指定键和语言的本地化消息
	GetMessage(key, locale string) string

	// GetMessageWithParams 获取带参数的翻译消息，支持 {{参数名}} 格式的占位符替换
	GetMessageWithParams(key, locale string, params map[string]string) string

	// GetSupportedLanguages 获取支持的语言列表
	GetSupportedLanguages() []config.SupportLang

	// GetSupportedLanguageCodes 获取支持的语言代码列表
	GetSupportedLanguageCodes() []string

	// GetDefaultLanguage 获取默认语言
	GetDefaultLanguage() string

	// IsLanguageSupported 检查语言是否支持
	IsLanguageSupported(locale string) bool

	// GetLanguageInfo 获取语言信息
	GetLanguageInfo(locale string) (config.SupportLang, bool)

	// GetMessageCount 获取已加载的消息数量（用于调试）
	GetMessageCount() int

	// GetMessagesByLanguage 获取指定语言的所有消息（用于调试）
	GetMessagesByLanguage(locale string) map[string]string

	// GetLocaleFromContext 从 Gin 上下文中智能检测语言
	// 优先级：用户设置 > Cookie > 浏览器偏好 > 系统默认
	GetLocaleFromContext(c *gin.Context) string

	// ReloadMessages 重新加载消息（开发环境使用）
	ReloadMessages() error
}