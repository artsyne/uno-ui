package render

import (
	"fmt"
	"strings"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// PageBuilder 页面构建器
type PageBuilder struct {
	config         *config.Config
	contextBuilder *ContextBuilder
	messageService MessageService
}

// NewPageBuilder 创建页面构建器
func NewPageBuilder(cfg *config.Config, contextBuilder *ContextBuilder, messageService MessageService) *PageBuilder {
	return &PageBuilder{
		config:         cfg,
		contextBuilder: contextBuilder,
		messageService: messageService,
	}
}

// BuildPageData 构建页面数据
func (pb *PageBuilder) BuildPageData(c *gin.Context, titleKey string) *PageData {
	return &PageData{
		Title:       pb.generatePageTitle(titleKey),
		Description: pb.generateDescription(titleKey),
		Data:        make(map[string]interface{}),
		Breadcrumb:  pb.generateBreadcrumb(c),
		Context:     pb.contextBuilder.BuildContext(c),
	}
}

// generatePageTitle 生成页面标题
func (pb *PageBuilder) generatePageTitle(titleKey string) string {
	title := pb.messageService.GetMessage(titleKey, "en")
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
func (pb *PageBuilder) generateDescription(titleKey string) string {
	descriptions := map[string]string{
		"diary.title":        "记录生活点滴，分享精彩瞬间",
		"diary.detail.title": "查看日记详细内容",
		"user.settings":      "管理个人账户和偏好设置",
		"user.profile":       "编辑个人资料信息",
	}

	if desc, exists := descriptions[titleKey]; exists {
		return desc
	}

	return "iDiary - 记录生活，分享美好"
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
