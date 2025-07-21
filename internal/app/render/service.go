package render

import (
	"github.com/biuren/hebe/internal/app/render/localization"
	"github.com/biuren/hebe/internal/app/render/page"
	"github.com/biuren/hebe/internal/app/render/response"
	"github.com/biuren/hebe/internal/app/render/types"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// PageDataBuilderImpl 页面数据构建器实现
type PageDataBuilderImpl struct{}

// BuildPageData 构建页面数据
func (p *PageDataBuilderImpl) BuildPageData(title, description, keywords string, data map[string]interface{}, breadcrumb []types.BreadcrumbItem, context *types.RenderContext, localizationManager localization.ManagerInterface) interface{} {
	return NewPageData(title, description, keywords, data, breadcrumb, context, localizationManager)
}

// RenderService 渲染服务实现
type RenderService struct {
	config              *config.Config
	contextBuilder      *page.ContextBuilder
	pageBuilder         *page.PageBuilder
	responseHandler     *response.Handler
	localizationManager localization.ManagerInterface
}

// NewRenderService 创建RenderService实例
func NewRenderService(cfg *config.Config) RenderServiceInterface {
	// 初始化本地化管理器
	localizationManager := localization.NewManager(cfg)

	service := &RenderService{
		config:              cfg,
		contextBuilder:      page.NewContextBuilder(localizationManager),
		localizationManager: localizationManager,
	}

	// 初始化各个组件
	pageDataBuilder := &PageDataBuilderImpl{}
	service.pageBuilder = page.NewPageBuilder(cfg, service.contextBuilder, service.localizationManager, pageDataBuilder)
	service.responseHandler = response.NewHandler()

	return service
}

// ===================== 页面数据组装（原view功能） =====================

// PageData 创建页面数据
func (s *RenderService) PageData(c *gin.Context, titleKey string) *PageData {
	result := s.pageBuilder.BuildPageData(c, titleKey)
	if pageData, ok := result.(*PageData); ok {
		return pageData
	}
	return nil
}

// ===================== 响应处理门面 =====================

// JSON响应
func (s *RenderService) JSON(c *gin.Context, code int, data interface{}) {
	s.responseHandler.JSON(c, code, data)
}

func (s *RenderService) JSONSuccess(c *gin.Context, data interface{}) {
	s.responseHandler.JSONSuccess(c, data)
}

func (s *RenderService) JSONError(c *gin.Context, err error, data ...interface{}) {
	s.responseHandler.JSONError(c, err, data...)
}

// HTML响应
func (s *RenderService) HTML(c *gin.Context, name string, data interface{}) {
	s.responseHandler.HTML(c, name, data)
}

func (s *RenderService) HTMLError(c *gin.Context, title, message string) {
	s.responseHandler.HTMLError(c, title, message)
}

// 其他响应
func (s *RenderService) Redirect(c *gin.Context, url string) {
	s.responseHandler.Redirect(c, url)
}
