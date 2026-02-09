package templates

import (
	"html/template"
	"reflect"
	"strings"

	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/pkg"
)

// GetAuthFuncs 获取认证相关的模板函数
// 返回包含所有用户认证和信息获取相关的模板函数映射
//
// 返回:
//   - template.FuncMap: 认证相关的模板函数映射
//
// 包含的函数:
//   - isLogin: 检查用户是否已登录
//   - authUser: 获取认证用户对象
//   - getLoginName: 获取登录名
//   - getHandle: 获取Handle(原用户名)
//   - getEmail: 获取邮箱地址
//   - getFaceUrl: 获取头像URL
//   - getUserAvatar: 获取用户头像（带默认值）
//   - getUserID: 获取用户ID
//   - getDisplayName: 获取显示名称(原全名)
//   - getUserGroup: 获取用户组
//   - getTimeZone: 获取时区
//   - getLang: 获取语言
//   - getStatus: 获取状态
//   - getFirstCharOfName: 获取姓名首字符
//   - getRandomVersion: 获取随机版本号
func GetAuthFuncs() template.FuncMap {
	return template.FuncMap{
		"isLogin":             isLogin,
		"authUser":            authUser,
		"getLoginName":        getLoginName,
		"getHandle":           getHandle,
		"getEmail":            getEmail,
		"getFaceUrl":          getFaceUrl,
		"getUserAvatar":       getUserAvatar,
		"getUserID":           getUserID,
		"getDisplayName":      getDisplayName,
		"getUserGroup":        getUserGroup,
		"getTimeZone":         getTimeZone,
		"getLang":             getLang,
		"getStatus":           getStatus,
		"getFirstCharOfName":  getFirstCharOfName,
		"getRandomVersion":    getRandomVersion,
		"getAutoSaveEnabled":  getAutoSaveEnabled,
		"getAutoSaveInterval": getAutoSaveInterval,
		"showUserName":        showUserName,
	}
}

// showUserName 获取用户显示名称
// 优先显示 DisplayName，如果为空则显示 Handle
// 支持 User, AuthUser, Entry (AddUser...) 等结构体
func showUserName(v interface{}) string {
	if v == nil {
		return ""
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ""
	}

	// Helper to get string field value
	getString := func(name string) string {
		f := val.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
		return ""
	}

	// 1. Try DisplayName / Handle
	dn := getString("DisplayName")
	h := getString("Handle")
	if dn != "" {
		return dn
	}
	if h != "" {
		return h
	}

	// 2. Try AddUserDisplayName / AddUserHandle (e.g. Entry)
	adn := getString("AddUserDisplayName")
	ah := getString("AddUserHandle")
	if adn != "" {
		return adn
	}
	if ah != "" {
		return ah
	}

	// 3. Try Context.AuthUser (PageData)
	ctxField := val.FieldByName("Context")
	if ctxField.IsValid() {
		ctxVal := ctxField
		if ctxVal.Kind() == reflect.Ptr {
			if ctxVal.IsNil() {
				return ""
			}
			ctxVal = ctxVal.Elem()
		}
		if ctxVal.Kind() == reflect.Struct {
			authUserField := ctxVal.FieldByName("AuthUser")
			if authUserField.IsValid() && !authUserField.IsNil() {
				return showUserName(authUserField.Interface())
			}
		}
	}

	return ""
}

func getAutoSaveEnabled(data *types.PageData[map[string]interface{}]) bool {
	cfg := config.GetConfig()
	if cfg == nil {
		return true
	}
	return cfg.AutoSave.Enabled
}

func getAutoSaveInterval(data *types.PageData[map[string]interface{}]) int {
	cfg := config.GetConfig()
	if cfg == nil || cfg.AutoSave.Interval <= 0 {
		return 60000
	}
	return cfg.AutoSave.Interval
}

// isLogin 检查用户是否已登录
// 通过页面数据判断当前用户的登录状态
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - bool: true表示已登录，false表示未登录
func isLogin(data *types.PageData[map[string]interface{}]) bool {
	return data.IsAuthenticated()
}

// authUser 获取认证用户对象
// 从页面数据中提取当前认证用户的完整信息
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - *types.AuthUser: 认证用户对象，未登录时返回nil
func authUser(data *types.PageData[map[string]interface{}]) *types.AuthUser {
	return data.GetAuthUser()
}

// getLoginName 获取登录名
// 获取当前认证用户的登录名称
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户登录名，未登录时返回空字符串
func getLoginName(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.LoginName
	}
	return ""
}

// getHandle 获取Handle
// 获取当前认证用户的Handle
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: Handle，未登录时返回空字符串
func getHandle(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.Handle
	}
	return ""
}

// getEmail 获取邮箱
// 获取当前认证用户的邮箱地址
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户邮箱地址，未登录时返回空字符串
func getEmail(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.Email
	}
	return ""
}

// getFaceUrl 获取头像URL
// 获取当前认证用户的头像URL地址
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户头像URL，未登录时返回空字符串
func getFaceUrl(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.FaceUrl
	}
	return ""
}

// getUserAvatar 获取用户头像（兼容性函数）
// 获取用户头像URL，如果用户未设置头像则返回默认头像
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户头像URL，未设置时返回默认头像路径
func getUserAvatar(data *types.PageData[map[string]interface{}]) string {
	if url := getFaceUrl(data); url != "" {
		return url
	}
	return "/static/img/face.jpg"
}

// getUserID 获取用户ID
// 获取当前认证用户的唯一标识ID
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - int64: 用户ID，未登录时返回0
func getUserID(data *types.PageData[map[string]interface{}]) int64 {
	if u := authUser(data); u != nil {
		return u.UserId
	}
	return 0
}

// getDisplayName 获取显示名称
// 获取当前认证用户的显示名称
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户显示名称，未登录时返回空字符串
func getDisplayName(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.DisplayName
	}
	return ""
}

// getUserGroup 获取用户组
// 获取当前认证用户所属的用户组标识
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - int8: 用户组标识，未登录时返回0
func getUserGroup(data *types.PageData[map[string]interface{}]) int8 {
	if u := authUser(data); u != nil {
		return u.UserGroup
	}
	return 0
}

// getTimeZone 获取时区
// 获取当前认证用户设置的时区信息
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户时区设置，未登录时返回空字符串
func getTimeZone(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.TimeZone
	}
	return ""
}

// getLang 获取语言
// 获取当前认证用户设置的语言偏好
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户语言设置，未登录时返回空字符串
func getLang(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		return u.Lang
	}
	return ""
}

// getStatus 获取状态
// 获取当前认证用户的账户状态
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - int8: 用户账户状态，未登录时返回0
func getStatus(data *types.PageData[map[string]interface{}]) int8 {
	if u := authUser(data); u != nil {
		return u.Status
	}
	return 0
}

// getFirstCharOfName 获取用户姓名首字符
// 获取用户显示名称或Handle的首字符，并返回大写形式，用于头像显示等场景
//
// 参数:
//   - data: 页面数据对象，包含认证信息
//
// 返回:
//   - string: 用户姓名首字符的大写形式，优先使用显示名称，其次Handle，默认返回"I"
//
// 逻辑:
//  1. 优先使用用户显示名称的首字符
//  2. 如果显示名称为空，使用Handle的首字符
//  3. 如果都为空或未登录，返回默认值"I"
func getFirstCharOfName(data *types.PageData[map[string]interface{}]) string {
	if u := authUser(data); u != nil {
		name := ""
		if u.DisplayName != "" {
			name = u.DisplayName
		} else if u.Handle != "" {
			name = u.Handle
		}

		if name != "" {
			for _, r := range name {
				return strings.ToUpper(string(r))
			}
		}
	}
	return "I" // 如果没有显示名称或Handle，返回默认值 "I"
}

// getRandomVersion 获取随机版本号
// 生成一个32位的随机字符串，通常用于静态资源的版本控制或缓存破坏
//
// 参数:
//   - data: 页面数据对象（此函数不使用该参数，保持接口一致性）
//
// 返回:
//   - string: 32位小写字母组成的随机字符串
func getRandomVersion(data *types.PageData[map[string]interface{}]) string {
	return pkg.Rand(32, pkg.RAND_KIND_LOWER)
}
