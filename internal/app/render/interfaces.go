package render

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

// RenderServiceInterface 渲染服务接口（合并response、view、template功能）
type RenderServiceInterface interface {
	// ===================== 页面数据组装（原view功能） =====================

	// PageData 创建页面数据
	PageData(c *gin.Context, titleKey string) *PageData

	// Message 获取国际化消息
	Message(key string, args ...interface{}) string

	// ===================== HTTP响应处理（原response功能） =====================

	// JSON响应
	JSON(c *gin.Context, code int, data interface{})
	JSONSuccess(c *gin.Context, data interface{})
	JSONError(c *gin.Context, err error, data ...interface{})

	// HTML响应
	HTML(c *gin.Context, name string, data interface{})
	HTMLError(c *gin.Context, title, message string)

	// 其他响应
	Redirect(c *gin.Context, url string)

	// ===================== 模板函数（原template功能） =====================

	// RegisterFuncs 注册模板函数
	RegisterFuncs() template.FuncMap

	// RenderBreadcrumb 渲染面包屑导航
	RenderBreadcrumb(breadcrumb []BreadcrumbItem) template.HTML

	// GetUserAvatar 获取用户头像URL
	GetUserAvatar(user interface{}, size string) string

	// IsUserLoggedIn 检查用户是否已登录
	IsUserLoggedIn(user interface{}) bool

	// GetCurrentLanguage 获取当前有效语言
	GetCurrentLanguage(c *gin.Context, user interface{}) string

	// FormatTime 格式化时间
	FormatTime(timestamp int64, format string, timezone ...string) string
}
