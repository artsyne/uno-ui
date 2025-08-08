package view

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/config"
	contracts "github.com/biuren/hebe/internal/contracts"
	response "github.com/biuren/hebe/internal/utils/http"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// ViewService 视图服务实现 - 直接实现核心逻辑
type ViewService struct {
	config          *config.Config
	i18nManager     i18n.I18nManager
	responseHandler *response.Handler
}

// NewViewService 创建视图服务
func NewViewService(cfg *config.Config) *ViewService {
	// 初始化国际化管理器
	i18nManager := i18n.NewLocalizationManager(cfg)

	// 初始化响应处理器
	responseHandler := &response.Handler{}

	return &ViewService{
		config:          cfg,
		i18nManager:     i18nManager,
		responseHandler: responseHandler,
	}
}

// ===================== 页面数据组装 =====================

// CreatePageData 创建页面数据（统一命名）
func (s *ViewService) CreatePageData(c *gin.Context, titleKey string) *contracts.PageData[map[string]interface{}] {
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

// CreatePageDataWithData 创建带数据的页面数据
func (s *ViewService) CreatePageDataWithData(c *gin.Context, titleKey string, data interface{}) *contracts.PageData[interface{}] {
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
func (s *ViewService) BuildRenderContext(c *gin.Context) *contracts.RenderContext {
	// 获取登录状态
	isLogin := s.getLoginStatus(c)
	
	// 获取认证用户
	authUser := s.getAuthUser(c)
	
	// 获取语言设置
	lang := s.GetLocaleFromContext(c)
	
	return &contracts.RenderContext{
		IsLogin:  isLogin,
		AuthUser: authUser,
		Lang:     lang,
	}
}

// ===================== 响应处理 =====================

// JSON 返回JSON响应
func (s *ViewService) JSON(c *gin.Context, code int, data interface{}) {
	s.responseHandler.JSON(c, code, data)
}

// JSONSuccess 返回成功JSON响应
func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
	s.responseHandler.JSONSuccess(c, data)
}

// JSONError 返回错误JSON响应
func (s *ViewService) JSONError(c *gin.Context, err error, data ...interface{}) {
	s.responseHandler.JSONError(c, err, data...)
}

// HTML 渲染HTML模板
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
	s.responseHandler.HTML(c, name, data)
}

// HTMLError 渲染错误页面
func (s *ViewService) HTMLError(c *gin.Context, title, message string) {
	s.responseHandler.HTMLError(c, title, message)
}

// Redirect 重定向响应
func (s *ViewService) Redirect(c *gin.Context, url string) {
	s.responseHandler.Redirect(c, url)
}

// ===================== 国际化服务 =====================

// GetMessage 获取国际化消息
func (s *ViewService) GetMessage(key, lang string) string {
	return s.i18nManager.GetMessage(key, lang)
}

// GetLocaleFromContext 从上下文获取语言设置
func (s *ViewService) GetLocaleFromContext(c *gin.Context) string {
	return s.i18nManager.GetLocaleFromContext(c)
}

// ===================== 配置服务 =====================

// GetSupportLangs 获取支持的语言列表
func (s *ViewService) GetSupportLangs() []config.SupportLang {
	return s.config.GetSupportLangs()
}

// GetTimeZones 获取时区列表
func (s *ViewService) GetTimeZones() []string {
	return s.config.GetTimeZone()
}

// ===================== 内部辅助方法 =====================

// getLoginStatus 获取登录状态
func (s *ViewService) getLoginStatus(c *gin.Context) bool {
	// 从上下文中获取登录状态
	if isLogin, exists := c.Get("IsLogin"); exists {
		if loginStatus, ok := isLogin.(bool); ok {
			return loginStatus
		}
	}
	return false
}

// getAuthUser 获取认证用户
func (s *ViewService) getAuthUser(c *gin.Context) *contracts.AuthUser {
	// 从上下文中获取认证用户
	if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
		if authUser, ok := authUserRaw.(*contracts.AuthUser); ok {
			return authUser
		}
	}
	return nil
}

// generateTitle 生成页面标题
func (s *ViewService) generateTitle(c *gin.Context, titleKey string) string {
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
func (s *ViewService) getSiteName() string {
	if s.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			return name
		}
	}
	return "iDiary"
}

// generateDescription 生成页面描述
func (s *ViewService) generateDescription(c *gin.Context, titleKey string) string {
	locale := s.i18nManager.GetLocaleFromContext(c)
	descKey := titleKey + ".desc"
	desc := s.i18nManager.GetMessage(descKey, locale)

	// 如果没有找到对应的描述，使用通用默认描述
	if desc == descKey {
		return s.i18nManager.GetMessage("common.default_page_desc", locale)
	}

	return desc
}