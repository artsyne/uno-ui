package render

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// TemplateFuncs 模板函数管理器
type TemplateFuncs struct {
	service *RenderService
}

// NewTemplateFuncs 创建模板函数管理器
func NewTemplateFuncs(service *RenderService) *TemplateFuncs {
	return &TemplateFuncs{
		service: service,
	}
}

// RegisterFuncs 注册模板函数
func (tf *TemplateFuncs) RegisterFuncs() template.FuncMap {
	return template.FuncMap{
		"renderBreadcrumb": tf.RenderBreadcrumb,
		"getMessage":       tf.service.Message,
		"formatTime":       tf.FormatTime,

		// 用户认证相关函数
		"isLogin": func(data *PageData) bool {
			return data.IsAuthenticated()
		},

		// 登录用户信息获取函数
		"getLoginName": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetLoginName()
			}
			return ""
		},
		"getAccessUrl": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetAccessUrl()
			}
			return ""
		},
		"getUserID": func(data *PageData) int64 {
			if user := data.GetAuthUser(); user != nil {
				return user.GetUserID()
			}
			return 0
		},
		"getUserName": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetUserName()
			}
			return ""
		},
		"getFullName": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetFullName()
			}
			return ""
		},
		"getEmail": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetEmail()
			}
			return ""
		},
		"getUserGroup": func(data *PageData) int8 {
			if user := data.GetAuthUser(); user != nil {
				return user.GetUserGroup()
			}
			return 0
		},
		"getFaceUrl": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetFaceUrl()
			}
			return ""
		},
		"getUserAvatar": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				if avatar := user.GetFaceUrl(); avatar != "" {
					return avatar
				}
			}
			return "/public/img/default-avatar.jpg"
		},
		"getLoginUTCTime": func(data *PageData) int64 {
			if user := data.GetAuthUser(); user != nil {
				return user.GetLoginUTCTime()
			}
			return 0
		},
		"getTimeZone": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetTimeZone()
			}
			return ""
		},
		"getLang": func(data *PageData) string {
			if user := data.GetAuthUser(); user != nil {
				return user.GetLang()
			}
			return ""
		},
		"getStatus": func(data *PageData) int8 {
			if user := data.GetAuthUser(); user != nil {
				return user.GetStatus()
			}
			return 0
		},
	}
}

// RenderBreadcrumb 渲染面包屑导航
func (tf *TemplateFuncs) RenderBreadcrumb(breadcrumb []BreadcrumbItem) template.HTML {
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

// FormatTime 格式化时间
func (tf *TemplateFuncs) FormatTime(timestamp int64, format string, timezone ...string) string {
	t := time.Unix(timestamp/1000, (timestamp%1000)*1000000)

	// 如果指定了时区，进行转换
	if len(timezone) > 0 && timezone[0] != "" {
		if loc, err := time.LoadLocation(timezone[0]); err == nil {
			t = t.In(loc)
		}
	}

	return t.Format(format)
}

// GetUserAvatar 获取用户头像URL
func (tf *TemplateFuncs) GetUserAvatar(user interface{}, size string) string {
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
func (tf *TemplateFuncs) IsUserLoggedIn(user interface{}) bool {
	return user != nil
}
