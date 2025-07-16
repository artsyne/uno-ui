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

	// ===================== HTTP响应处理 =====================

	// JSON响应
	JSON(c *gin.Context, code int, data interface{})
	JSONSuccess(c *gin.Context, data interface{})
	JSONError(c *gin.Context, err error, data ...interface{})

	// HTML响应
	HTML(c *gin.Context, name string, data interface{})
	HTMLError(c *gin.Context, title, message string)

	// 其他响应
	Redirect(c *gin.Context, url string)

	// ===================== 模板函数=====================

	// RegisterFuncs 注册模板函数
	RegisterFuncs() template.FuncMap
}
