package templates

import (
	"html/template"

	"github.com/biuren/hebe/internal/contracts"
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
func isLogin(data *types.PageData[map[string]interface{}]) bool {
	return data.IsAuthenticated()
}

// authUser 获取认证用户对象
func authUser(data *types.PageData[map[string]interface{}]) *types.AuthUser {
	return data.GetAuthUser()
}

// getLoginName 获取登录名
func getLoginName(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.LoginName
	}
	return ""
}

// getUserName 获取用户名
func getUserName(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.UserName
	}
	return ""
}

// getEmail 获取邮箱
func getEmail(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.Email
	}
	return ""
}

// getFaceUrl 获取头像URL
func getFaceUrl(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.FaceUrl
	}
	return ""
}

// getUserAvatar 获取用户头像（兼容性函数）
func getUserAvatar(data *types.PageData[map[string]interface{}]) string {
	if url := getFaceUrl(data); url != "" {
		return url
	}
	return "/static/img/face.jpg"
}

// getUserID 获取用户ID
func getUserID(data *types.PageData[map[string]interface{}]) int64 {
	if u := authUser(data); u != nil {
		return u.UserId
	}
	return 0
}

// getFullName 获取全名
func getFullName(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.FullName
	}
	return ""
}

// getUserGroup 获取用户组
func getUserGroup(data *types.PageData[map[string]interface{}]) int8 {
	if u := authUser(data); u != nil {
		return u.UserGroup
	}
	return 0
}

// getTimeZone 获取时区
func getTimeZone(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.TimeZone
	}
	return ""
}

// getLang 获取语言
func getLang(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.Lang
	}
	return ""
}

// getStatus 获取状态
func getStatus(data *types.PageData[map[string]interface{}]) int8 {
	if u := authUser(data); u != nil {
		return u.Status
	}
	return 0
}
