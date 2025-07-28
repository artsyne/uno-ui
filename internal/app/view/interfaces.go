package view

import (
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/shared/contracts"
)

// PageDataService 负责页面数据组装
type PageDataService interface {
	// PageData 创建页面数据（保持向后兼容）
	PageData(c *gin.Context, titleKey string) *types.PageData[map[string]interface{}]
	
	// BuildPageData 创建页面数据（保持向后兼容）
	BuildPageData(c *gin.Context, titleKey string, data interface{}) *types.PageData[interface{}]
	
	// BuildRenderContext 构建渲染上下文
	BuildRenderContext(c *gin.Context) *types.RenderContext
}

// ResponseService 负责HTTP响应处理
type ResponseService interface {
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
}

// I18nService 负责国际化处理
type I18nService interface {
	// GetMessage 获取国际化消息
	GetMessage(key, lang string) string
	
	// GetLocaleFromContext 从上下文获取语言设置
	GetLocaleFromContext(c *gin.Context) string
}

// ViewServiceInterface 组合接口，提供完整的视图服务功能
// 保持向后兼容，现有代码无需修改
type ViewServiceInterface interface {
	PageDataService
	ResponseService
	I18nService
	
	// GetConfig 获取配置（用于向后兼容）
	GetConfig() *config.Config
}
