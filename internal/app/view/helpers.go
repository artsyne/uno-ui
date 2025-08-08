package view

import (
	"fmt"

	"github.com/gin-gonic/gin"

	contracts "github.com/biuren/hebe/internal/contracts"
)

// ===================== 内部辅助方法 =====================

// getLoginStatus 获取登录状态
func (s *ViewService) getLoginStatus(c *gin.Context) bool {
	// 从上下文中获取登录状态
	if isLogin, exists := c.Get("IsLogin"); exists {
		if loginStatus, ok := isLogin.(bool); ok {
			return loginStatus
		}
	}
	return false
}

// getAuthUser 获取认证用户
func (s *ViewService) getAuthUser(c *gin.Context) *contracts.AuthUser {
	// 从上下文中获取认证用户
	if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
		if authUser, ok := authUserRaw.(*contracts.AuthUser); ok {
			return authUser
		}
	}
	return nil
}

// generateTitle 生成页面标题
func (s *ViewService) generateTitle(c *gin.Context, titleKey string) string {
	locale := s.i18nManager.GetLocaleFromContext(c)
	title := s.i18nManager.GetMessage(titleKey, locale)
	siteName := s.getSiteName()

	if title == titleKey {
		// 如果没有找到对应的标题，使用key本身
		return fmt.Sprintf("%s - %s", titleKey, siteName)
	}

	return fmt.Sprintf("%s - %s", title, siteName)
}

// generateDescription 生成页面描述
func (s *ViewService) generateDescription(c *gin.Context, titleKey string) string {
	locale := s.i18nManager.GetLocaleFromContext(c)
	descKey := titleKey + ".desc"
	desc := s.i18nManager.GetMessage(descKey, locale)

	// 如果没有找到对应的描述，使用通用默认描述
	if desc == descKey {
		return s.i18nManager.GetMessage("common.default_page_desc", locale)
	}

	return desc
}