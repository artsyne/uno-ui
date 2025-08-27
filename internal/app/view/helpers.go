package view

import (
	"fmt"

	"github.com/gin-gonic/gin"

	contracts "github.com/biuren/hebe/internal/contracts"
)

// ===================== 内部辅助方法 =====================

// getLoginStatus 获取登录状态
// 从Gin上下文中获取用户的登录状态
//
// 参数:
//   - c: Gin上下文对象
//
// 返回:
//   - bool: true表示已登录，false表示未登录
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
// 从Gin上下文中获取当前认证的用户信息
//
// 参数:
//   - c: Gin上下文对象
//
// 返回:
//   - *contracts.AuthUser: 认证用户对象，如果未登录或获取失败则返回nil
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
// 根据标题键和当前语言环境生成完整的页面标题
//
// 参数:
//   - c: Gin上下文对象，用于获取语言环境
//   - titleKey: 标题的国际化键名
//
// 返回:
//   - string: 格式为"页面标题 - 网站名称"的完整标题
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
// 根据标题键和当前语言环境生成页面描述信息
//
// 参数:
//   - c: Gin上下文对象，用于获取语言环境
//   - titleKey: 标题的国际化键名，描述键为titleKey + ".desc"
//
// 返回:
//   - string: 页面描述文本，如果找不到对应描述则返回默认描述
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