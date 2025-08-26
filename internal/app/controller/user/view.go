package user

import (
	"github.com/gin-gonic/gin"
)

// ShowLoginPage 显示登录页面
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
	
	ctrl.renderService.HTML(c, "member/settings", data)
}