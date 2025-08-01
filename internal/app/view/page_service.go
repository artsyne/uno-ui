package view

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/config"
	contracts "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// pageDataService 页面数据服务实现
type pageDataService struct {
	i18nManager i18n.I18nManager
	config      *config.Config
}

// NewPageDataService 创建页面数据服务
func NewPageDataService(i18nManager i18n.I18nManager) *pageDataService {
	return &pageDataService{
		i18nManager: i18nManager,
		config:      config.GetConfig(),
	}
}

// PageData 创建页面数据（保持向后兼容）
func (s *pageDataService) PageData(c *gin.Context, titleKey string) *contracts.PageData[map[string]interface{}] {
	// 构建渲染上下文
	ctx := s.BuildRenderContext(c)
	
	// 生成完整标题
	title := s.generateTitle(c, titleKey)
	
	return &contracts.PageData[map[string]interface{}]{
		Title:   title,
		Context: ctx,
		Data:    make(map[string]interface{}),
	}
}

// BuildPageData 创建页面数据（保持向后兼容）
func (s *pageDataService) BuildPageData(c *gin.Context, titleKey string, data interface{}) *contracts.PageData[interface{}] {
	// 构建渲染上下文
	ctx := s.BuildRenderContext(c)
	
	// 生成完整标题
	title := s.generateTitle(c, titleKey)
	
	return &contracts.PageData[interface{}]{
		Title:   title,
		Context: ctx,
		Data:    data,
	}
}

// BuildRenderContext 构建渲染上下文
func (s *pageDataService) BuildRenderContext(c *gin.Context) *contracts.RenderContext {
	// 获取登录状态
	isLogin := s.getLoginStatus(c)
	
	// 获取认证用户
	authUser := s.getAuthUser(c)
	
	// 获取语言设置
	lang := s.i18nManager.GetLocaleFromContext(c)
	
	return &contracts.RenderContext{
		IsLogin:  isLogin,
		AuthUser: authUser,
		Lang:     lang,
	}
}

// getLoginStatus 提取登录状态
func (s *pageDataService) getLoginStatus(c *gin.Context) bool {
	if isLogin, exists := c.Get("IsLogin"); exists {
		if login, ok := isLogin.(bool); ok {
			return login
		}
	}
	return false
}

// getAuthUser 提取认证用户信息
func (s *pageDataService) getAuthUser(c *gin.Context) *contracts.AuthUser {
	if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
		if authUser, ok := authUserRaw.(*contracts.AuthUser); ok {
			return authUser
		}
		slog.Warn("AuthUser 类型断言失败", "value", authUserRaw)
	}
	return nil
}

// generateTitle 生成页面标题
func (s *pageDataService) generateTitle(c *gin.Context, titleKey string) string {
	locale := s.i18nManager.GetLocaleFromContext(c)
	title := s.i18nManager.GetMessage(titleKey, locale)
	siteName := s.getSiteName()

	if title == titleKey {
		// 如果没有找到对应的标题，使用key本身
		return fmt.Sprintf("%s - %s", titleKey, siteName)
	}

	return fmt.Sprintf("%s - %s", title, siteName)
}

// getSiteName 获取网站名称
func (s *pageDataService) getSiteName() string {
	if s.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			return name
		}
	}
	return "iDiary"
}

// generateDescription 生成页面描述
func (s *pageDataService) generateDescription(c *gin.Context, titleKey string) string {
	locale := s.i18nManager.GetLocaleFromContext(c)
	descKey := titleKey + ".desc"
	desc := s.i18nManager.GetMessage(descKey, locale)

	// 如果没有找到对应的描述，使用通用默认描述
	if desc == descKey {
		return s.i18nManager.GetMessage("common.default_page_desc", locale)
	}

	return desc
}