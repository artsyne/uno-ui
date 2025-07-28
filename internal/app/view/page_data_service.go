package view

import (
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/shared/contracts"
	"github.com/biuren/hebe/internal/shared/i18n"
)

// pageDataService 页面数据服务实现
type pageDataService struct {
	i18nManager i18n.I18nManager
}

// NewPageDataService 创建页面数据服务
func NewPageDataService(i18nManager i18n.I18nManager) PageDataService {
	return &pageDataService{
		i18nManager: i18nManager,
	}
}

// PageData 创建页面数据（保持向后兼容）
func (s *pageDataService) PageData(c *gin.Context, titleKey string) *types.PageData[map[string]interface{}] {
	// 构建渲染上下文
	ctx := s.BuildRenderContext(c)
	
	// 获取翻译后的标题
	title := s.i18nManager.GetMessage(titleKey, ctx.Lang)
	
	return &types.PageData[map[string]interface{}]{
		Title:   title,
		Context: ctx,
		Data:    make(map[string]interface{}),
	}
}

// BuildPageData 创建页面数据（保持向后兼容）
func (s *pageDataService) BuildPageData(c *gin.Context, titleKey string, data interface{}) *types.PageData[interface{}] {
	// 构建渲染上下文
	ctx := s.BuildRenderContext(c)
	
	// 获取翻译后的标题
	title := s.i18nManager.GetMessage(titleKey, ctx.Lang)
	
	return &types.PageData[interface{}]{
		Title:   title,
		Context: ctx,
		Data:    data,
	}
}

// BuildRenderContext 构建渲染上下文
func (s *pageDataService) BuildRenderContext(c *gin.Context) *types.RenderContext {
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