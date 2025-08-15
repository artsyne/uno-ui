package templates

import (
	"html/template"
	"log/slog"

	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/internal/utils/config"
)

// 全局本地化管理器实例（用于需要manager方法的函数）
var globalManager i18n.I18nManager

// SetGlobalManager 设置全局管理器实例
func SetGlobalManager(manager i18n.I18nManager) {
	globalManager = manager
}

// GetLocalizationFuncs 获取本地化相关的模板函数
func GetLocalizationFuncs() template.FuncMap {
	return template.FuncMap{
		// 核心消息获取函数
		"getMessage": getMessage,

		// 扩展本地化函数
		"getMessageWithDefault": getMessageWithDefault,
		"getMessages":           getMessages,
		"getSupportedLanguages": getSupportedLanguages,
		"isLanguageSupported":   isLanguageSupported,
	}
}

// getMessage 获取本地化消息（核心函数）
func getMessage(key, locale string) string {
	cacheKey := i18n.BuildCacheKey(locale, key)
	if message, exists := i18n.GetMessage(cacheKey); exists {
		slog.Debug("getMessage: 消息已找到", "key", key, "locale", locale, "message", message)
		return message
	}
	slog.Warn("getMessage: 消息未找到", "key", key, "locale", locale, "cacheKey", cacheKey)
	return key // 降级处理：返回原始key
}

// getMessageWithDefault 带默认值的消息获取
func getMessageWithDefault(key, locale, defaultMsg string) string {
	cacheKey := i18n.BuildCacheKey(locale, key)
	if message, exists := i18n.GetMessage(cacheKey); exists {
		return message
	}
	if defaultMsg != "" {
		return defaultMsg
	}
	return key
}

// getMessages 批量获取消息
func getMessages(keys []string, locale string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		result[key] = getMessage(key, locale)
	}
	return result
}

// getSupportedLanguages 获取支持的语言列表
func getSupportedLanguages() []config.SupportLang {
	if globalManager != nil {
		return globalManager.GetSupportedLanguages()
	}
	return nil
}

// isLanguageSupported 检查语言是否支持
func isLanguageSupported(locale string) bool {
	if globalManager != nil {
		return globalManager.IsLanguageSupported(locale)
	}
	return false
}
