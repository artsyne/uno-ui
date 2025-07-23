package view

import (
	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/shared/i18n"
	"github.com/biuren/hebe/internal/shared/response"
	"github.com/biuren/hebe/internal/shared/types"
	"github.com/gin-gonic/gin"
)

// ViewService 视图服务实现
type ViewService struct {
	config          *config.Config
	responseHandler *response.Handler
	i18nManager     i18n.I18nManager
}

// NewViewService 创建视图服务
func NewViewService(cfg *config.Config) *ViewService {
	// 初始化国际化管理器
	i18nManager := i18n.NewLocalizationManager(cfg)

	// 初始化响应处理器
	responseHandler := response.NewHandler()

	return &ViewService{
		config:          cfg,
		responseHandler: responseHandler,
		i18nManager:     i18nManager,
	}
}

// ===================== 页面数据组装 ====================// PageData 创建页面数据（保持原有接口）
func (s *ViewService) PageData(c *gin.Context, titleKey string) *types.PageData[map[string]interface{}] {
	// 构建渲染上下文
	ctx := s.buildRenderContext(c)
	
	// 获取翻译后的标题
	title := s.i18nManager.GetMessage(titleKey, ctx.Lang)
	
	return &types.PageData[map[string]interface{}]{
		Title:   title,
		Context: ctx,
		Data:    make(map[string]interface{}),
	}
}

// BuildPageData 创建页面数据
func (s *ViewService) BuildPageData(c *gin.Context, titleKey string, data interface{}) *types.PageData[interface{}] {
	// 构建渲染上下文
	ctx := s.buildRenderContext(c)
	
	// 获取翻译后的标题
	title := s.i18nManager.GetMessage(titleKey, ctx.Lang)
	
	return &types.PageData[interface{}]{
		Title:   title,
		Context: ctx,
		Data:    data,
	}
}

// buildRenderContext 构建渲染上下文
func (s *ViewService) buildRenderContext(c *gin.Context) *types.RenderContext {
	// 检测语言
	lang := s.i18nManager.GetLocaleFromContext(c)
	
	// 检查登录状态
	isLogin := false
	var authUser *types.AuthUser
	if isLoginValue, exists := c.Get("IsLogin"); exists {
		isLogin = isLoginValue.(bool)
	}
	if authUserValue, exists := c.Get("AuthUser"); exists && authUserValue != nil {
		if user, ok := authUserValue.(*types.AuthUser); ok {
			authUser = user
		}
	}
	
	return &types.RenderContext{
		IsLogin:  isLogin,
		AuthUser: authUser,
		Lang:     lang,
	}
}

// ===================== 响应处理门面 =====================

// JSON响应
func (s *ViewService) JSON(c *gin.Context, code int, data interface{}) {
	s.responseHandler.JSON(c, code, data)
}

func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
	s.responseHandler.JSONSuccess(c, data)
}

func (s *ViewService) JSONError(c *gin.Context, err error, data ...interface{}) {
	s.responseHandler.JSONError(c, err, data...)
}

// HTML响应
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
	s.responseHandler.HTML(c, name, data)
}

func (s *ViewService) HTMLError(c *gin.Context, title, message string) {
	s.responseHandler.HTMLError(c, title, message)
}

// 其他响应
func (s *ViewService) Redirect(c *gin.Context, url string) {
	s.responseHandler.Redirect(c, url)
}