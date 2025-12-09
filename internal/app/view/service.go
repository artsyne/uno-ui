package view

import (
	"github.com/gin-gonic/gin"

	contracts "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// ViewService 视图服务实现 - 直接实现核心逻辑
// 提供页面渲染、数据组装、响应处理等视图层相关功能
type ViewService struct {
	config      *config.Config   // 应用配置
	i18nManager i18n.I18nManager // 国际化管理器
}

// NewViewService 创建视图服务
// 初始化视图服务实例，注入必要的依赖
//
// 参数:
//   - cfg: 应用配置对象
//   - i18nManager: 国际化管理器
//
// 返回:
//   - *ViewService: 视图服务实例
func NewViewService(cfg *config.Config, i18nManager i18n.I18nManager) *ViewService {
	return &ViewService{
		config:      cfg,
		i18nManager: i18nManager,
	}
}

// ===================== 页面数据组装 =====================

// CreatePageData 创建页面数据（统一命名）
// 创建标准的页面数据结构，包含标题、上下文和空的数据映射
//
// 参数:
//   - c: Gin上下文对象
//   - titleKey: 页面标题的国际化键名
//
// 返回:
//   - *contracts.PageData[map[string]interface{}]: 包含空数据映射的页面数据
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
// 创建包含指定数据的页面数据结构
//
// 参数:
//   - c: Gin上下文对象
//   - titleKey: 页面标题的国际化键名
//   - data: 页面数据，可以是任意类型
//
// 返回:
//   - *contracts.PageData[interface{}]: 包含指定数据的页面数据
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
// 从Gin上下文中提取并构建页面渲染所需的上下文信息
//
// 参数:
//   - c: Gin上下文对象
//
// 返回:
//   - *contracts.RenderContext: 包含登录状态、用户信息、语言设置的渲染上下文
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
