package view

import (
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/config"
	contracts "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// ViewService 视图服务实现 - 直接实现核心逻辑
type ViewService struct {
	config      *config.Config
	i18nManager i18n.I18nManager
}

// NewViewService 创建视图服务
func NewViewService(cfg *config.Config) *ViewService {
	// 初始化国际化管理器
	i18nManager := i18n.NewLocalizationManager(cfg)

	return &ViewService{
		config:      cfg,
		i18nManager: i18nManager,
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