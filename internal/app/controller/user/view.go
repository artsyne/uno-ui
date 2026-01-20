package user

import (
	"github.com/gin-gonic/gin"
)

// ShowLoginPage 显示登录页面
// 如果用户已登录则重定向到首页，否则显示登录表单
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 检查用户登录状态
//   - 已登录用户重定向到首页
//   - 未登录用户显示登录页面
func (ctrl *UserController) ShowLoginPage(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := ctrl.auth.IsAuthenticated(c); isAuth {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.login.title")
	ctrl.renderService.HTML(c, "passport/login", data)
}

// ShowRegisterPage 显示注册页面
// 如果用户已登录则重定向到首页，否则显示注册表单及相关配置选项
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 检查用户登录状态
//   - 已登录用户重定向到首页
//   - 未登录用户显示注册页面
//   - 提供支持的语言列表和时区选项
func (ctrl *UserController) ShowRegisterPage(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := ctrl.auth.IsAuthenticated(c); isAuth {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.register.title")
	supportLangs := ctrl.renderService.GetSupportLangs()
	langMaps := make([]map[string]string, len(supportLangs))
	for i, lang := range supportLangs {
		langMaps[i] = map[string]string{
			"Code": lang.Code,
			"Name": lang.Label,
		}
	}
	data.Data["SupportLangs"] = langMaps
	data.Data["TimeZones"] = ctrl.renderService.GetTimeZones()
	ctrl.renderService.HTML(c, "passport/register", data)
}

// ShowPasswordResetPage 显示找回密码页面
// 如果用户已登录则重定向到首页，否则显示密码找回表单
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 检查用户登录状态
//   - 已登录用户重定向到首页
//   - 未登录用户显示密码找回页面
func (ctrl *UserController) ShowPasswordResetPage(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := ctrl.auth.IsAuthenticated(c); isAuth {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.findpwd.title")
	ctrl.renderService.HTML(c, "passport/findpwd", data)
}

// ShowPasswordResetForm 显示重置密码表单页面
// 验证重置密码链接的有效性，如果有效则显示密码重置表单
//
// URL参数:
//   - uid: 用户ID
//   - sign: 验证签名
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 验证URL参数完整性
//   - 解析和验证用户ID
//   - 检查验证码有效性
//   - 显示密码重置表单或错误页面
func (ctrl *UserController) ShowPasswordResetForm(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	if uid == "" || sign == "" {
		data := ctrl.renderService.CreatePageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "缺少必要参数"
		ctrl.renderService.HTML(c, "error", data)
		return
	}

	// 解析用户ID
	userID, err := ctrl.service.ParseUserID(uid)
	if err != nil {
		data := ctrl.renderService.CreatePageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "用户ID格式错误"
		ctrl.renderService.HTML(c, "error", data)
		return
	}

	// 检查验证码
	if valid, _ := ctrl.service.CheckValidCode(userID, sign, ctrl.config.ValidCode.TypeResetPwd); !valid {
		data := ctrl.renderService.CreatePageData(c, "error")
		data.Data["title"] = "验证码无效"
		data.Data["message"] = "验证码已过期或无效"
		ctrl.renderService.HTML(c, "error.html", data)
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.resetpwd.title")
	data.Data["uid"] = uid
	data.Data["sign"] = sign
	ctrl.renderService.HTML(c, "passport/resetpwd", data)
}

// ShowUserActivationPage 显示激活用户页面
// 处理用户账户激活请求，验证激活链接并执行激活操作
//
// URL参数:
//   - uid: 用户ID
//   - sign: 激活验证签名
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 验证URL参数完整性
//   - 执行用户激活操作
//   - 显示激活成功或失败页面
func (ctrl *UserController) ShowUserActivationPage(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	if uid == "" || sign == "" {
		data := ctrl.renderService.CreatePageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "缺少必要参数"
		ctrl.renderService.HTML(c, "error", data)
		return
	}

	// 激活用户
	if err := ctrl.service.ActivateUser(uid, sign); err != nil {
		data := ctrl.renderService.CreatePageData(c, "user.activate.failed.title")
		data.Data["message"] = "激活链接无效或已过期"
		ctrl.renderService.HTML(c, "passport/activate_failed", data)
	} else {
		data := ctrl.renderService.CreatePageData(c, "user.activate.success.title")
		data.Data["message"] = "账户激活成功，请登录"
		ctrl.renderService.HTML(c, "passport/activate_success", data)
	}
}

// ShowUserSettingsPage 显示用户设置页面
// 显示当前登录用户的个人设置页面，包括个人信息、语言和时区设置
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 获取当前登录用户信息
//   - 获取用户详细设置信息
//   - 提供支持的语言列表和时区选项
//   - 渲染用户设置页面
//
// 注意:
//   - 此方法需要用户已登录（通过中间件保证）
func (ctrl *UserController) ShowUserSettingsPage(c *gin.Context) {
	// 读取用户认证信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取用户设置信息
	userInfo, err := ctrl.service.GetUserByUserId(authUser.UserId)
	if err != nil {
		data := ctrl.renderService.CreatePageData(c, "error")
		data.Data["title"] = "错误"
		data.Data["message"] = "获取用户信息失败"
		ctrl.renderService.HTML(c, "error", data)
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.settings.title")
	data.Data["User"] = userInfo
	
	// 添加支持的语言和时区列表
	supportLangs := ctrl.renderService.GetSupportLangs()
	data.Data["SupportLangs"] = supportLangs
	data.Data["TimeZones"] = ctrl.renderService.GetTimeZones()
	
	ctrl.renderService.HTML(c, "user/settings", data)
}