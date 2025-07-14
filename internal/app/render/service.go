package render

import (
	"html/template"
	"log/slog"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// RenderService 渲染服务实现（合并response、view、template功能）
type RenderService struct {
	config              *config.Config
	contextBuilder      *ContextBuilder
	pageBuilder         *PageBuilder
	responseHandler     *ResponseHandler
	templateFuncs       *TemplateFuncs
	localizationManager *LocalizationManager
	utils               *RenderUtils
}

// NewRenderService 创建RenderService实例
func NewRenderService(cfg *config.Config) RenderServiceInterface {
	// 初始化本地化管理器
	localizationManager, err := NewLocalizationManager(cfg)
	if err != nil {
		slog.Error("初始化本地化管理器失败", "error", err)
		// 如果初始化失败，使用一个空的管理器避免系统崩溃
		localizationManager = &LocalizationManager{
			cfg:      cfg,
			messages: make(map[string]map[string]string),
		}
	}

	service := &RenderService{
		config:              cfg,
		contextBuilder:      NewContextBuilder(localizationManager),
		localizationManager: localizationManager,
	}

	// 初始化各个组件
	service.pageBuilder = NewPageBuilder(cfg, service.contextBuilder, service.localizationManager)
	service.responseHandler = NewResponseHandler(service)
	service.templateFuncs = NewTemplateFuncs(service)
	service.utils = NewRenderUtils(service)

	return service
}

// ===================== 页面数据组装（原view功能） =====================

// PageData 创建页面数据
func (s *RenderService) PageData(c *gin.Context, titleKey string) *PageData {
	return s.pageBuilder.BuildPageData(c, titleKey)
}

// Message 获取国际化消息
func (s *RenderService) Message(key string, args ...interface{}) string {
	// 从gin context中获取当前语言
	lang := "en"
	if c, ok := args[0].(*gin.Context); ok {
		if user := s.contextBuilder.BuildContext(c).AuthUser; user != nil {
			lang = user.GetLang()
		}
	}

	return s.localizationManager.GetMessage(key, lang)
}

// ===================== HTTP响应处理（原response功能） =====================

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

// ===================== 模板函数（原template功能） =====================

// RegisterFuncs 注册模板函数
func (s *RenderService) RegisterFuncs() template.FuncMap {
	return s.templateFuncs.RegisterFuncs()
}

// RenderBreadcrumb 渲染面包屑导航
func (s *RenderService) RenderBreadcrumb(breadcrumb []BreadcrumbItem) template.HTML {
	return s.templateFuncs.RenderBreadcrumb(breadcrumb)
}

// GetUserAvatar 获取用户头像URL
func (s *RenderService) GetUserAvatar(user interface{}, size string) string {
	return s.templateFuncs.GetUserAvatar(user, size)
}

// IsUserLoggedIn 检查用户是否已登录
func (s *RenderService) IsUserLoggedIn(user interface{}) bool {
	return s.templateFuncs.IsUserLoggedIn(user)
}

// GetCurrentLanguage 获取当前有效语言
func (s *RenderService) GetCurrentLanguage(c *gin.Context, user interface{}) string {
	if renderUser, ok := user.(RenderUser); ok {
		return s.localizationManager.GetUserLocale(renderUser)
	}
	return "en"
}

// FormatTime 格式化时间
func (s *RenderService) FormatTime(timestamp int64, format string, timezone ...string) string {
	return s.templateFuncs.FormatTime(timestamp, format, timezone...)
}
