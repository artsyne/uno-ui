package render

import (
	"html/template"
	"log/slog"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// RenderService 渲染服务实现（合并response、view、template功能）
type RenderService struct {
	config          *config.Config
	contextBuilder  *ContextBuilder
	pageBuilder     *PageBuilder
	responseHandler *ResponseHandler
	templateFuncs   *TemplateFuncs
	i18nManager     *I18nManager
	utils           *RenderUtils
}

// NewRenderService 创建RenderService实例
func NewRenderService(cfg *config.Config) RenderServiceInterface {
	service := &RenderService{
		config:         cfg,
		contextBuilder: NewContextBuilder(),
	}

	// 初始化各个组件
	service.i18nManager = NewI18nManager(service)
	service.pageBuilder = NewPageBuilder(cfg, service.contextBuilder, service.i18nManager)
	service.responseHandler = NewResponseHandler(service)
	service.templateFuncs = NewTemplateFuncs(service)
	service.utils = NewRenderUtils(service)

	// 初始化国际化消息
	service.initMessages()

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

	return s.i18nManager.GetMessage(key, lang)
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
		return s.i18nManager.GetUserLanguage(renderUser)
	}
	return "en"
}

// FormatTime 格式化时间
func (s *RenderService) FormatTime(timestamp int64, format string, timezone ...string) string {
	return s.templateFuncs.FormatTime(timestamp, format, timezone...)
}

// ===================== 私有方法 =====================

// initMessages 初始化国际化消息
func (s *RenderService) initMessages() {
	// 预加载常用语言
	supportedLanguages := []string{"en", "zh-cn", "zh-tw"}
	for _, lang := range supportedLanguages {
		if err := s.i18nManager.LoadMessages(lang); err != nil {
			slog.Error("Failed to load messages", "lang", lang, "error", err)
		}
	}

	// 设置默认消息
	s.setDefaultMessages()
}

// setDefaultMessages 设置默认消息
func (s *RenderService) setDefaultMessages() {
	defaultMessages := map[string]string{
		"success":             "Success",
		"error":               "Error",
		"internal_error":      "Internal Server Error",
		"validation_error":    "Validation Error",
		"user_exists":         "User already exists",
		"email_exists":        "Email already exists",
		"user_pwd_error":      "Invalid username or password",
		"user_not_exists":     "User does not exist",
		"login_required":      "Please log in first",
		"permission_denied":   "Permission denied",
		"not_found":           "Page not found",
		"diary.title":         "My Diary",
		"diary.detail.title":  "Diary Details",
		"user.settings":       "User Settings",
		"user.profile":        "User Profile",
		"user.login.title":    "User Login",
		"user.register.title": "User Registration",
		"user.findpwd.title":  "Find Password",
		"user.resetpwd.title": "Reset Password",
		"user.activate.title": "Activate Account",
	}

	for key, message := range defaultMessages {
		s.i18nManager.SetMessage(key, message)
	}
}

// ===================== MessageService 接口实现 =====================

// GetMessage 获取国际化消息
func (s *RenderService) GetMessage(key string, args ...interface{}) string {
	return s.Message(key, args...)
}

// GetMessageWithLang 获取指定语言的消息
func (s *RenderService) GetMessageWithLang(key, lang string, args ...interface{}) string {
	return s.i18nManager.GetMessage(key, lang)
}

// HasMessage 检查是否存在消息
func (s *RenderService) HasMessage(key string) bool {
	return s.i18nManager.GetMessage(key, "en") != key
}

// SetMessage 设置消息
func (s *RenderService) SetMessage(key, message string) {
	s.i18nManager.SetMessage(key, message)
}

// LoadMessages 加载消息
func (s *RenderService) LoadMessages(lang string) error {
	return s.i18nManager.LoadMessages(lang)
}
