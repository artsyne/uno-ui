package view

import (
	"github.com/gin-gonic/gin"
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