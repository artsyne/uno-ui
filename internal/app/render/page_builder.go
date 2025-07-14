package render

import (
	"fmt"
	"strings"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// PageBuilder 页面构建器
type PageBuilder struct {
	config              *config.Config
	contextBuilder      *ContextBuilder
	localizationManager *LocalizationManager
}

// NewPageBuilder 创建页面构建器
func NewPageBuilder(cfg *config.Config, contextBuilder *ContextBuilder, localizationManager *LocalizationManager) *PageBuilder {
	return &PageBuilder{
		config:              cfg,
		contextBuilder:      contextBuilder,
		localizationManager: localizationManager,
	}
}

// BuildPageData 构建页面数据
func (pb *PageBuilder) BuildPageData(c *gin.Context, titleKey string) *PageData {
	return &PageData{
		Title:               pb.generatePageTitle(c, titleKey),
		Description:         pb.generateDescription(c, titleKey),
		Data:                make(map[string]interface{}),
		Breadcrumb:          pb.generateBreadcrumb(c),
		Context:             pb.contextBuilder.BuildContext(c),
		localizationManager: pb.localizationManager, // 注入本地化管理器
	}
}

// generatePageTitle 生成页面标题
func (pb *PageBuilder) generatePageTitle(c *gin.Context, titleKey string) string {
	// 使用智能语言检测
	locale := pb.localizationManager.GetLocaleFromContext(c)

	// 从 localization 获取标题
	title := pb.localizationManager.GetMessage(titleKey, locale)
	siteName := "iDiary"

	if pb.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			siteName = name
		}
	}

	if title == titleKey {
		// 如果没有找到对应的标题，使用key本身
		return fmt.Sprintf("%s - %s", titleKey, siteName)
	}

	return fmt.Sprintf("%s - %s", title, siteName)
}

// generateDescription 生成页面描述
func (pb *PageBuilder) generateDescription(c *gin.Context, titleKey string) string {
	// 使用智能语言检测
	locale := pb.localizationManager.GetLocaleFromContext(c)

	// 从 localization 获取描述（使用 .desc 后缀）
	descKey := titleKey + ".desc"
	desc := pb.localizationManager.GetMessage(descKey, locale)

	// 如果没有找到对应的描述，使用通用默认描述
	if desc == descKey {
		return pb.localizationManager.GetMessage("common.default_page_desc", locale)
	}

	return desc
}

// generateBreadcrumb 生成面包屑导航
func (pb *PageBuilder) generateBreadcrumb(c *gin.Context) []BreadcrumbItem {
	path := c.Request.URL.Path
	breadcrumb := []BreadcrumbItem{
		{Name: "首页", URL: "/"},
	}

	// 根据路径模式智能生成面包屑
	switch {
	case strings.HasPrefix(path, "/u/"):
		// 用户页面：首页 > 用户名 > 具体页面
		if authUser, exists := c.Get("AuthUser"); exists {
			if user, ok := authUser.(map[string]interface{}); ok {
				if userName, exists := user["user_name"]; exists {
					if userID, exists := user["user_id"]; exists {
						breadcrumb = append(breadcrumb, BreadcrumbItem{
							Name: fmt.Sprintf("%v", userName),
							URL:  fmt.Sprintf("/u/%v", userID),
						})
					}
				}
			}
		}

		if strings.Contains(path, "/settings") {
			breadcrumb = append(breadcrumb, BreadcrumbItem{
				Name: "设置",
			})
		}

	case strings.HasPrefix(path, "/d/"):
		// 日记页面：首页 > 用户名 > 日记本名称
		// 这里需要从数据中获取具体信息
		// TODO: 实现日记页面面包屑

	case strings.HasPrefix(path, "/p/"):
		// 文章页面：首页 > 用户名 > 日记本名称 > 文章标题
		// 这里需要从数据中获取具体信息
		// TODO: 实现文章页面面包屑
	}

	return breadcrumb
}
