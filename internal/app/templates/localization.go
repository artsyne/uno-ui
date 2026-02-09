package templates

import (
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// 全局本地化管理器实例（用于需要manager方法的函数）
var globalManager i18n.I18nManager

// SetGlobalManager 设置全局管理器实例
// 设置全局的国际化管理器实例，用于需要访问管理器方法的模板函数
//
// 参数:
//   - manager: 国际化管理器实例
func SetGlobalManager(manager i18n.I18nManager) {
	globalManager = manager
}

// GetLocalizationFuncs 获取本地化相关的模板函数
// 返回包含所有国际化和本地化相关的模板函数映射
//
// 返回:
//   - template.FuncMap: 本地化相关的模板函数映射
//
// 包含的函数:
//   - getMessage: 获取本地化消息（核心函数）
//   - getMessageWithDefault: 带默认值的消息获取
//   - getMessages: 批量获取消息
//   - getSupportedLanguages: 获取支持的语言列表
//   - isLanguageSupported: 检查语言是否支持
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
// 根据消息键和语言环境获取对应的本地化消息，支持参数替换
//
// 参数:
//   - key: 消息键，用于标识特定的消息
//   - locale: 语言环境代码（如"en", "zh-CN"等）
//   - args: 可选参数，通常是一个map[string]interface{}（通过dict函数生成），用于替换消息中的占位符 {{key}}
//
// 返回:
//   - string: 本地化消息内容，如果未找到则返回原始键作为降级处理
func getMessage(key, locale string, args ...interface{}) string {
	cacheKey := i18n.BuildCacheKey(locale, key)
	var message string
	var exists bool

	if message, exists = i18n.GetMessage(cacheKey); !exists {
		slog.Warn("getMessage: 消息未找到", "key", key, "locale", locale, "cacheKey", cacheKey)
		message = key // 降级处理：返回原始key
	}

	// 处理参数替换
	if len(args) > 0 {
		if params, ok := args[0].(map[string]interface{}); ok {
			for k, v := range params {
				val := fmt.Sprint(v)
				message = strings.ReplaceAll(message, "{{"+k+"}}", val)
			}
		} else if params, ok := args[0].(map[string]string); ok {
			for k, v := range params {
				message = strings.ReplaceAll(message, "{{"+k+"}}", v)
			}
		}
	}

	return message
}

// getMessageWithDefault 带默认值的消息获取
// 获取本地化消息，如果未找到则使用提供的默认消息
//
// 参数:
//   - key: 消息键，用于标识特定的消息
//   - locale: 语言环境代码
//   - defaultMsg: 默认消息内容，当找不到本地化消息时使用
//
// 返回:
//   - string: 本地化消息、默认消息或原始键（按优先级返回）
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
// 批量获取多个消息键对应的本地化消息
//
// 参数:
//   - keys: 消息键列表
//   - locale: 语言环境代码
//
// 返回:
//   - map[string]string: 键值对映射，键为消息键，值为对应的本地化消息
func getMessages(keys []string, locale string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		result[key] = getMessage(key, locale)
	}
	return result
}

// getSupportedLanguages 获取支持的语言列表
// 获取系统支持的所有语言配置列表
//
// 返回:
//   - []config.SupportLang: 支持的语言配置列表，如果全局管理器未设置则返回nil
func getSupportedLanguages() []config.SupportLang {
	if globalManager != nil {
		return globalManager.GetSupportedLanguages()
	}
	return nil
}

// isLanguageSupported 检查语言是否支持
// 检查指定的语言环境是否在系统支持的语言列表中
//
// 参数:
//   - locale: 要检查的语言环境代码
//
// 返回:
//   - bool: true表示支持该语言，false表示不支持或全局管理器未设置
func isLanguageSupported(locale string) bool {
	if globalManager != nil {
		return globalManager.IsLanguageSupported(locale)
	}
	return false
}
