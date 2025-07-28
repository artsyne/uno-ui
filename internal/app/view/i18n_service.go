package view

import (
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/shared/i18n"
)

// i18nService 国际化服务实现
type i18nService struct {
	manager i18n.I18nManager
}

// NewI18nService 创建国际化服务
func NewI18nService(manager i18n.I18nManager) I18nService {
	return &i18nService{
		manager: manager,
	}
}

// GetMessage 获取国际化消息
func (s *i18nService) GetMessage(key, lang string) string {
	return s.manager.GetMessage(key, lang)
}

// GetLocaleFromContext 从上下文获取语言设置
func (s *i18nService) GetLocaleFromContext(c *gin.Context) string {
	return s.manager.GetLocaleFromContext(c)
}