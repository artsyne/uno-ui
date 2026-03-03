package view

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ===================== 内部辅助方法 =====================

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
