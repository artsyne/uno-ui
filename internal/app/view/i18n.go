package view

import (
	"github.com/gin-gonic/gin"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// ===================== 国际化服务 =====================

// GetMessage 获取国际化消息
// 根据指定的键和语言获取对应的国际化消息文本
//
// 参数:
//   - key: 消息键名
//   - lang: 语言代码，如"en"、"zh-CN"等
//
// 返回:
//   - string: 对应的国际化消息文本
func (s *ViewService) GetMessage(key, lang string) string {
	return s.i18nManager.GetMessage(key, lang)
}

// GetLocaleFromContext 从上下文获取语言设置
// 从Gin上下文中提取当前请求的语言环境设置
//
// 参数:
//   - c: Gin上下文对象
//
// 返回:
//   - string: 当前请求的语言代码
func (s *ViewService) GetLocaleFromContext(c *gin.Context) string {
	return s.i18nManager.GetLocaleFromContext(c)
}

// GetI18nManager 获取国际化管理器实例
// 用于在其他组件中重用已创建的i18nManager，避免重复初始化
//
// 返回:
//   - i18n.I18nManager: 国际化管理器实例
func (s *ViewService) GetI18nManager() i18n.I18nManager {
	return s.i18nManager
}