package render

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/biuren/hebe/internal/core/config"

	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// RenderService 渲染服务实现（合并response、view、template功能）
type RenderService struct {
	config     *config.Config
	messages   map[string]map[string]string
}

// NewRenderService 创建RenderService实例
func NewRenderService(cfg *config.Config) RenderServiceInterface {
	service := &RenderService{
		config:     cfg,
		messages:   make(map[string]map[string]string),
	}
	service.initMessages()
	return service
}

// ===================== 页面数据组装（原view功能） =====================
// getBoolFromContext 安全地从上下文获取 bool 值
func (s *RenderService) getBoolFromContext(c *gin.Context, key string) bool {
    if value, exists := c.Get(key); exists {
        if boolValue, ok := value.(bool); ok {
            return boolValue
        }
    }
    return false
}

// PageData 创建页面数据（不再处理用户信息）
func (s *RenderService) PageData(c *gin.Context, titleKey string) *PageData {
	pageData := &PageData{
		Title:       s.generatePageTitle(titleKey),
		Description: s.generateDescription(titleKey),
		Data:        make(map[string]interface{}),
		Breadcrumb:  s.generateBreadcrumb(c),
		IsLogin: s.getBoolFromContext(c, "IsLogin"),
	}

	return pageData
}

// generatePageTitle 生成页面标题
func (s *RenderService) generatePageTitle(titleKey string) string {
	title := s.Message(titleKey)
	siteName := "iDiary"

	if s.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			siteName = name
		}
	}

	if title == titleKey {
		// 如果没有找到对应的标题，使用key本身
		return fmt.Sprintf("%s - %s", titleKey, siteName)
	}

	return fmt.Sprintf("%s - %s", title, siteName)
}

// generateDescription 生成页面描述
func (s *RenderService) generateDescription(titleKey string) string {
	descriptions := map[string]string{
		"diary.title":        "记录生活点滴，分享精彩瞬间",
		"diary.detail.title": "查看日记详细内容",
		"user.settings":      "管理个人账户和偏好设置",
		"user.profile":       "编辑个人资料信息",
	}

	if desc, exists := descriptions[titleKey]; exists {
		return desc
	}

	return "iDiary - 记录生活，分享美好"
}

// generateBreadcrumb 生成面包屑导航
func (s *RenderService) generateBreadcrumb(c *gin.Context) []BreadcrumbItem {
	path := c.Request.URL.Path
	breadcrumb := []BreadcrumbItem{
		{Name: "首页", URL: "/"},
	}

	// 根据路径模式智能生成面包屑
	switch {
	case strings.HasPrefix(path, "/u/"):
		// 用户页面：首页 > 用户名 > 具体页面
		if authUser, exists := c.Get("AuthUser"); exists {
			if user, ok := authUser.(map[string]interface{}); ok {
				if userName, exists := user["user_name"]; exists {
					if userID, exists := user["user_id"]; exists {
						breadcrumb = append(breadcrumb, BreadcrumbItem{
							Name: fmt.Sprintf("%v", userName),
							URL:  fmt.Sprintf("/u/%v", userID),
						})
					}
				}
			}
		}

		if strings.Contains(path, "/settings") {
			breadcrumb = append(breadcrumb, BreadcrumbItem{
				Name: "设置",
			})
		}

	case strings.HasPrefix(path, "/d/"):
		// 日记页面：首页 > 用户名 > 日记本名称
		// 这里需要从数据中获取具体信息

	case strings.HasPrefix(path, "/p/"):
		// 文章页面：首页 > 用户名 > 日记本名称 > 文章标题
		// 这里需要从数据中获取具体信息
	}

	return breadcrumb
}

// Message 获取国际化消息
func (s *RenderService) Message(key string, args ...interface{}) string {
	// 基础消息映射
	messages := map[string]string{
		"success":             "成功",
		"error":               "错误",
		"internal_error":      "系统内部错误",
		"validation_error":    "参数验证错误",
		"user_exists":         "用户名已存在",
		"email_exists":        "邮箱已存在",
		"user_pwd_error":      "用户名或密码错误",
		"user_not_exists":     "用户不存在",
		"login_required":      "请先登录",
		"permission_denied":   "权限不足",
		"not_found":           "页面未找到",
		"diary.title":         "我的日记",
		"diary.detail.title":  "日记详情",
		"user.settings":       "个人设置",
		"user.profile":        "个人资料",
		"user.login.title":    "用户登录",
		"user.register.title": "用户注册",
		"user.findpwd.title":  "找回密码",
		"user.resetpwd.title": "重置密码",
		"user.activate.title": "激活账户",
	}

	if msg, exists := messages[key]; exists {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}

	// 如果找不到对应的消息，返回key本身
	return key
}

// ===================== HTTP响应处理（原response功能） =====================

// JSON 返回JSON响应
func (s *RenderService) JSON(c *gin.Context, code int, data interface{}) {
	response := JSONResponse{
		Code:    code,
		Message: "success",
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// JSONSuccess 返回成功响应
func (s *RenderService) JSONSuccess(c *gin.Context, data interface{}) {
	s.JSON(c, 0, data)
}

// JSONError 返回错误响应
func (s *RenderService) JSONError(c *gin.Context, err error, data ...interface{}) {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}

	// 处理AppError
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.HttpCode, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.MsgKey,
			Data:    errorData,
		})

		// 记录错误的详细信息
		s.logAppError(c, appErr)
		return
	}

	// 处理其他类型的错误，转换为内部错误
	internalErr := errors.InternalErr()
	response := ErrorResponse{
		Code:    internalErr.Code,
		Message: internalErr.MsgKey,
		Data:    errorData,
	}
	c.JSON(internalErr.HttpCode, response)
}

// HTML 返回HTML页面
func (s *RenderService) HTML(c *gin.Context, name string, data interface{}) {
	c.HTML(http.StatusOK, name, data)
}

// HTMLError 渲染错误页面
func (s *RenderService) HTMLError(c *gin.Context, title, message string) {
	pageData := &PageData{
		Title: title,
		Data: map[string]interface{}{
			"error_title":   title,
			"error_message": message,
		},
	}
	c.HTML(http.StatusOK, "error.html", pageData)
}

// Redirect 重定向
func (s *RenderService) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}

// logAppError 记录应用错误的详细信息
func (s *RenderService) logAppError(c *gin.Context, err *errors.AppError) {
	// 使用slog记录错误
	slog.Error("API错误",
		"error_code", err.Code,
		"error_key", err.MsgKey,
		"http_code", err.HttpCode,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"user_agent", c.Request.UserAgent(),
		"client_ip", c.ClientIP(),
		"details", err.Details,
	)
}

// ===================== 模板函数（原template功能） =====================

// RegisterFuncs 注册模板函数（无依赖注入）
func (s *RenderService) RegisterFuncs() template.FuncMap {
	return template.FuncMap{
		"renderBreadcrumb": s.RenderBreadcrumb,
		"getMessage":       s.Message,
		"formatTime":       s.FormatTime,
		}
}

// RenderBreadcrumb 渲染面包屑导航
func (s *RenderService) RenderBreadcrumb(breadcrumb []BreadcrumbItem) template.HTML {
	if len(breadcrumb) == 0 {
		return ""
	}

	var parts []string
	for i, item := range breadcrumb {
		if i == len(breadcrumb)-1 {
			// 最后一个项目不需要链接
			parts = append(parts, fmt.Sprintf(`<span>%s</span>`, item.Name))
		} else if item.URL != "" {
			parts = append(parts, fmt.Sprintf(`<a href="%s">%s</a>`, item.URL, item.Name))
		} else {
			parts = append(parts, fmt.Sprintf(`<span>%s</span>`, item.Name))
		}
	}

	html := fmt.Sprintf(`<nav class="breadcrumb">%s</nav>`, strings.Join(parts, ` <span class="separator">/</span> `))
	return template.HTML(html)
}

// GetUserAvatar 获取用户头像URL
func (s *RenderService) GetUserAvatar(user interface{}, size string) string {
	if user == nil {
		return "/public/img/default-avatar.jpg"
	}

	// 尝试从用户对象中获取头像URL
	if userMap, ok := user.(map[string]interface{}); ok {
		if faceSmallUrl, exists := userMap["face_small_url"]; exists {
			if url, ok := faceSmallUrl.(string); ok && url != "" {
				return url
			}
		}
	}

	// 默认头像
	return "/public/img/face.jpg"
}

// IsUserLoggedIn 检查用户是否已登录
func (s *RenderService) IsUserLoggedIn(user interface{}) bool {
	return user != nil
}

// GetCurrentLanguage 获取当前有效语言
// 优先级：用户设置 > 请求头 > 默认语言
func (s *RenderService) GetCurrentLanguage(c *gin.Context, user interface{}) string {
	// 1. 如果用户已登录且设置了语言偏好，优先使用用户设置
	if user != nil {
		if userMap, ok := user.(map[string]interface{}); ok {
			if lang, exists := userMap["lang"]; exists {
				if langStr, ok := lang.(string); ok && langStr != "" {
					return langStr
				}
			}
		}
	}

	// 2. 从请求头获取语言
	return s.detectLanguageFromHeader(c)
}

// FormatTime 格式化时间
func (s *RenderService) FormatTime(timestamp int64, format string, timezone ...string) string {
	t := time.Unix(timestamp/1000, (timestamp%1000)*1000000)

	// 如果指定了时区，进行转换
    if len(timezone) > 0 && timezone[0] != "" {
        if loc, err := time.LoadLocation(timezone[0]); err == nil {
            t = t.In(loc)
        }
    }

	return t.Format(format)
}

// ===================== 私有辅助方法 =====================

// initMessages 初始化消息映射
func (s *RenderService) initMessages() {
	// 这里可以从配置文件或数据库加载消息
	// 目前使用硬编码的方式
}

// detectLanguageFromHeader 从请求头检测语言
func (s *RenderService) detectLanguageFromHeader(c *gin.Context) string {
	// 从上下文获取语言设置
	if lang, exists := c.Get("language"); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}

	// 从请求头获取Accept-Language
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang != "" {
		// 简单解析Accept-Language，取第一个
		langs := strings.Split(acceptLang, ",")
		if len(langs) > 0 {
			lang := strings.TrimSpace(strings.Split(langs[0], ";")[0])
			// 转换为支持的语言格式
			switch {
			case strings.HasPrefix(lang, "zh-CN"), strings.HasPrefix(lang, "zh-cn"):
				return "zh-cn"
			case strings.HasPrefix(lang, "zh-TW"), strings.HasPrefix(lang, "zh-tw"):
				return "zh-tw"
			case strings.HasPrefix(lang, "en"):
				return "en"
			}
		}
	}

	return "zh-cn" // 默认中文
}