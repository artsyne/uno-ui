package templates

import (
	"html/template"

	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/app/render/types"
)

// GetAuthFuncs 获取认证相关的模板函数
func GetAuthFuncs() template.FuncMap {
	return template.FuncMap{
		"isLogin":       isLogin,
		"authUser":      authUser,
		"getLoginName":  getLoginName,
		"getUserName":   getUserName,
		"getEmail":      getEmail,
		"getFaceUrl":    getFaceUrl,
		"getUserAvatar": getUserAvatar,
		"getUserID":     getUserID,
		"getFullName":   getFullName,
		"getUserGroup":  getUserGroup,
		"getTimeZone":   getTimeZone,
		"getLang":       getLang,
		"getStatus":     getStatus,
	}
}

// isLogin 检查用户是否已登录
func isLogin(data *render.PageData) bool {
	return data != nil && data.IsAuthenticated()
}

// authUser 获取认证用户对象
func authUser(data *render.PageData) types.RenderUser {
	if data == nil {
		return nil
	}
	return data.GetAuthUser()
}

// getLoginName 获取登录名
func getLoginName(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetLoginName()
	}
	return ""
}

// getUserName 获取用户名
func getUserName(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetUserName()
	}
	return ""
}

// getEmail 获取邮箱
func getEmail(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetEmail()
	}
	return ""
}

// getFaceUrl 获取头像URL
func getFaceUrl(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetFaceUrl()
	}
	return ""
}

// getUserAvatar 获取用户头像（兼容性函数）
func getUserAvatar(data *render.PageData) string {
	if url := getFaceUrl(data); url != "" {
		return url
	}
	return "/static/img/face.jpg"
}

// getUserID 获取用户ID
func getUserID(data *render.PageData) int64 {
	if u := authUser(data); u != nil {
		return u.GetUserID()
	}
	return 0
}

// getFullName 获取全名
func getFullName(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetFullName()
	}
	return ""
}

// getUserGroup 获取用户组
func getUserGroup(data *render.PageData) int8 {
	if u := authUser(data); u != nil {
		return u.GetUserGroup()
	}
	return 0
}

// getTimeZone 获取时区
func getTimeZone(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetTimeZone()
	}
	return ""
}

// getLang 获取语言
func getLang(data *render.PageData) string {
	if u := authUser(data); u != nil {
		return u.GetLang()
	}
	return ""
}

// getStatus 获取状态
func getStatus(data *render.PageData) int8 {
	if u := authUser(data); u != nil {
		return u.GetStatus()
	}
	return 0
}
