package view

import (
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
	contracts "github.com/biuren/hebe/internal/contracts"
)

// ViewServiceInterface 视图服务接口，提供完整的视图服务功能
type ViewServiceInterface interface {
	// ===================== 页面数据组装 =====================
	
	// CreatePageData 创建页面数据
	CreatePageData(c *gin.Context, titleKey string) *contracts.PageData[map[string]interface{}]
	
	// CreatePageDataWithData 创建带数据的页面数据
	CreatePageDataWithData(c *gin.Context, titleKey string, data interface{}) *contracts.PageData[interface{}]
	
	// BuildRenderContext 构建渲染上下文
	BuildRenderContext(c *gin.Context) *contracts.RenderContext
	
	// ===================== 响应处理 =====================
	
	// JSON 返回JSON响应
	JSON(c *gin.Context, code int, data interface{})
	
	// JSONSuccess 返回成功JSON响应
	JSONSuccess(c *gin.Context, data interface{})
	
	// JSONError 返回错误JSON响应
	JSONError(c *gin.Context, err error, data ...interface{})
	
	// HTML 渲染HTML模板
	HTML(c *gin.Context, name string, data interface{})
	
	// HTMLError 渲染错误页面
	HTMLError(c *gin.Context, title, message string)
	
	// Redirect 重定向响应
	Redirect(c *gin.Context, url string)
	
	// ===================== 国际化服务 =====================
	
	// GetMessage 获取国际化消息
	GetMessage(key, lang string) string
	
	// GetLocaleFromContext 获取当前请求的语言环境
	GetLocaleFromContext(c *gin.Context) string

	// GetSupportLangs 获取支持的语言列表
	GetSupportLangs() []config.SupportLang

	// GetTimeZones 获取时区列表
	GetTimeZones() []string
}
