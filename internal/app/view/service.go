package view

import (
	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/shared/i18n"
	"github.com/biuren/hebe/internal/shared/http"
	"github.com/biuren/hebe/internal/shared/contracts"
	"github.com/gin-gonic/gin"
)

// ViewService 视图服务实现 - 组合各个专门的服务
type ViewService struct {
	config *config.Config
	
	// 组合的服务
	pageDataService PageDataService
	responseService ResponseService
	i18nService     I18nService
}

// NewViewService 创建视图服务
func NewViewService(cfg *config.Config) *ViewService {
	// 初始化国际化管理器
	i18nManager := i18n.NewLocalizationManager(cfg)

	// 初始化响应处理器
	responseHandler := response.NewHandler()

	// 创建各个专门的服务
	pageDataService := NewPageDataService(i18nManager)
	responseService := NewResponseService(responseHandler)
	i18nService := NewI18nService(i18nManager)

	return &ViewService{
		config:          cfg,
		pageDataService: pageDataService,
		responseService: responseService,
		i18nService:     i18nService,
	}
}

// ===================== 页面数据组装 =====================
// 委托给 pageDataService

// PageData 创建页面数据（保持向后兼容）
func (s *ViewService) PageData(c *gin.Context, titleKey string) *types.PageData[map[string]interface{}] {
	return s.pageDataService.PageData(c, titleKey)
}

// BuildPageData 创建页面数据（保持向后兼容）
func (s *ViewService) BuildPageData(c *gin.Context, titleKey string, data interface{}) *types.PageData[interface{}] {
	return s.pageDataService.BuildPageData(c, titleKey, data)
}

// BuildRenderContext 构建渲染上下文
func (s *ViewService) BuildRenderContext(c *gin.Context) *types.RenderContext {
	return s.pageDataService.BuildRenderContext(c)
}

// ===================== 响应处理门面 =====================
// 委托给 responseService

// JSON 返回JSON响应
func (s *ViewService) JSON(c *gin.Context, code int, data interface{}) {
	s.responseService.JSON(c, code, data)
}

// JSONSuccess 返回成功JSON响应
func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
	s.responseService.JSONSuccess(c, data)
}

// JSONError 返回错误JSON响应
func (s *ViewService) JSONError(c *gin.Context, err error, data ...interface{}) {
	s.responseService.JSONError(c, err, data...)
}

// HTML 渲染HTML模板
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
	s.responseService.HTML(c, name, data)
}

// HTMLError 渲染错误页面
func (s *ViewService) HTMLError(c *gin.Context, title, message string) {
	s.responseService.HTMLError(c, title, message)
}

// Redirect 重定向响应
func (s *ViewService) Redirect(c *gin.Context, url string) {
	s.responseService.Redirect(c, url)
}

// ===================== 国际化服务 =====================
// 委托给 i18nService

// GetMessage 获取国际化消息
func (s *ViewService) GetMessage(key, lang string) string {
	return s.i18nService.GetMessage(key, lang)
}

// GetLocaleFromContext 从上下文获取语言设置
func (s *ViewService) GetLocaleFromContext(c *gin.Context) string {
	return s.i18nService.GetLocaleFromContext(c)
}

// ===================== 配置访问 =====================

// GetConfig 获取配置（用于向后兼容）
func (s *ViewService) GetConfig() *config.Config {
	return s.config
}