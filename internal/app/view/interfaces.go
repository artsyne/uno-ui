package view

import (
	"github.com/biuren/hebe/internal/shared/types"
	"github.com/gin-gonic/gin"
)

// ViewServiceInterface 视图服务接口
type ViewServiceInterface interface {
	// ===================== 页面数据组装 =====================

	// PageData 创建页面数据
	PageData(c *gin.Context, titleKey string) *types.PageData[map[string]interface{}]

	// BuildPageData 创建页面数据
	BuildPageData(c *gin.Context, titleKey string, data interface{}) *types.PageData[interface{}]

	// ===================== 响应处理门面 =====================

	// JSON响应
	JSON(c *gin.Context, code int, data interface{})
	JSONSuccess(c *gin.Context, data interface{})
	JSONError(c *gin.Context, err error, data ...interface{})

	// HTML响应
	HTML(c *gin.Context, name string, data interface{})
	HTMLError(c *gin.Context, title, message string)

	// 其他响应
	Redirect(c *gin.Context, url string)
}
