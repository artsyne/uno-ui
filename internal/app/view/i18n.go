package view

import (
	"github.com/gin-gonic/gin"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// ===================== 国际化服务 =====================

// GetMessage 获取国际化消息
func (s *ViewService) GetMessage(key, lang string) string {
	return s.i18nManager.GetMessage(key, lang)
}

// GetLocaleFromContext 从上下文获取语言设置
func (s *ViewService) GetLocaleFromContext(c *gin.Context) string {
	return s.i18nManager.GetLocaleFromContext(c)
}

// GetI18nManager 获取国际化管理器实例
// 用于在其他组件中重用已创建的i18nManager，避免重复初始化
func (s *ViewService) GetI18nManager() i18n.I18nManager {
	return s.i18nManager
}