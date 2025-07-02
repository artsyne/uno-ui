package view

import (
	"fmt"
	"strings"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// ViewService 视图服务实现
type ViewService struct {
	config *config.Config
}

// NewViewService 创建ViewService实例
func NewViewService(cfg *config.Config) ViewServiceInterface {
	return &ViewService{
		config: cfg,
	}
}

// PageData 创建页面数据
func (s *ViewService) PageData(c *gin.Context, titleKey string) *PageData {
	pageData := &PageData{
		Title:       s.generatePageTitle(titleKey),
		Description: s.generateDescription(titleKey),
		Data:        make(map[string]interface{}),
		Lang:        s.getCurrentLanguage(c),
	}

	// 从上下文获取用户信息（如果存在）
	if user, exists := c.Get("AuthUser"); exists {
		pageData.User = user
	}

	// 生成导航数据
	pageData.Navigation = s.generateNavigation(c)

	return pageData
}

// PageTitle 生成页面标题
func (s *ViewService) PageTitle(content, author, pageType string) string {
	siteName := "iDiary"
	if s.config != nil {
		// 从配置获取站点名称
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			siteName = name
		}
	}

	var parts []string

	if content != "" {
		parts = append(parts, content)
	}

	if author != "" {
		parts = append(parts, author)
	}

	if pageType != "" {
		parts = append(parts, pageType)
	}

	parts = append(parts, siteName)

	return strings.Join(parts, " - ")
}

// Message 获取国际化消息
func (s *ViewService) Message(key string, args ...interface{}) string {
	// 基础消息映射
	messages := map[string]string{
		"success":            "成功",
		"error":              "错误",
		"internal_error":     "系统内部错误",
		"validation_error":   "参数验证错误",
		"user_exists":        "用户名已存在",
		"email_exists":       "邮箱已存在",
		"user_pwd_error":     "用户名或密码错误",
		"user_not_exists":    "用户不存在",
		"login_required":     "请先登录",
		"permission_denied":  "权限不足",
		"not_found":          "页面未找到",
		"diary.title":        "我的日记",
		"diary.detail.title": "日记详情",
		"user.settings":      "个人设置",
		"user.profile":       "个人资料",
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

// generatePageTitle 生成页面标题
func (s *ViewService) generatePageTitle(titleKey string) string {
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
func (s *ViewService) generateDescription(titleKey string) string {
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

// getCurrentLanguage 获取当前语言
func (s *ViewService) getCurrentLanguage(c *gin.Context) string {
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

// generateNavigation 生成导航数据
func (s *ViewService) generateNavigation(c *gin.Context) *NavigationData {
	currentPath := c.Request.URL.Path

	items := []NavigationItem{
		{
			Name:   "首页",
			URL:    "/",
			Active: currentPath == "/",
		},
		{
			Name:   "日记",
			URL:    "/diary",
			Active: strings.HasPrefix(currentPath, "/diary"),
		},
		{
			Name:   "文章",
			URL:    "/post",
			Active: strings.HasPrefix(currentPath, "/post"),
		},
	}

	// 如果用户已登录，添加个人中心导航
	if _, exists := c.Get("AuthUser"); exists {
		items = append(items, NavigationItem{
			Name:   "个人中心",
			URL:    "/member",
			Active: strings.HasPrefix(currentPath, "/member"),
		})
	}

	return &NavigationData{
		Current: currentPath,
		Items:   items,
	}
}
