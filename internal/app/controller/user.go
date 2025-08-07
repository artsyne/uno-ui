package controller

import (
	"net/http"
	"strings"

	"github.com/biuren/hebe/internal/app/view"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/config"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// UserController 用户HTTP处理层接口
// 定义所有与用户相关的HTTP处理操作，包括认证、注册、设置等
type UserController interface {
	// RegisterRoutes 注册所有用户相关路由
	RegisterRoutes(router *gin.Engine)

	// ===================== 页面渲染方法 =====================

	// ShowLoginPage 显示登录页面
	ShowLoginPage(c *gin.Context)

	// ShowRegisterPage 显示注册页面
	ShowRegisterPage(c *gin.Context)

	// ShowPasswordResetPage 显示找回密码页面
	ShowPasswordResetPage(c *gin.Context)

	// ShowPasswordResetForm 显示重置密码表单页面
	ShowPasswordResetForm(c *gin.Context)

	// ShowUserActivationPage 显示激活用户页面
	ShowUserActivationPage(c *gin.Context)

	// ShowUserSettingsPage 显示用户设置页面
	ShowUserSettingsPage(c *gin.Context)

	// ===================== 认证相关方法 =====================

	// Login 处理用户登录请求
	Login(c *gin.Context)

	// Register 处理用户注册请求
	Register(c *gin.Context)

	// Logout 处理用户退出登录请求
	Logout(c *gin.Context)

	// ===================== 验证方法 =====================

	// ValidateUsername 验证用户名是否可用
	ValidateUsername(c *gin.Context)

	// ValidateEmail 验证邮箱是否可用
	ValidateEmail(c *gin.Context)

	// ===================== 密码管理方法 =====================

	// RequestPasswordReset 处理找回密码请求
	RequestPasswordReset(c *gin.Context)

	// ResetPassword 处理重置密码请求
	ResetPassword(c *gin.Context)

	// ===================== 用户设置相关方法 =====================

	// UpdateProfile 更新用户基本资料
	UpdateProfile(c *gin.Context)

	// UpdateEmail 更新用户邮箱
	UpdateEmail(c *gin.Context)

	// UpdatePassword 修改用户密码
	UpdatePassword(c *gin.Context)
}

// userController 用户HTTP处理层实现
type userController struct {
	service       user.ServiceInterface
	auth          user.AuthHelperInterface
	config        *config.Config // 缓存配置，避免重复调用全局变量
	renderService view.ViewServiceInterface
}

// NewUserController 创建用户HTTP处理层实例
func NewUserController() UserController {
	config := config.GetConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(taskRepo)
	return &userController{
		service:       user.NewService(userRepo, taskService),
		auth:          user.NewAuthHelper(config),
		config:        config,
		renderService: view.NewViewService(config),
	}
}

// RegisterRoutes 注册所有用户相关路由
// 保持与原系统完全一致的路由路径
func (h *userController) RegisterRoutes(router *gin.Engine) {
	// 用户认证相关路由（无需鉴权）
	router.GET("/register", h.ShowRegisterPage)
	router.POST("/register", h.Register)
	router.POST("/check-username", h.ValidateUsername)
	router.POST("/check-email", h.ValidateEmail)

	router.GET("/login", h.ShowLoginPage)
	router.POST("/login", h.Login)

	router.POST("/loginout", h.Logout)

	router.GET("/findpwd", h.ShowPasswordResetPage)
	router.POST("/findpwd", h.RequestPasswordReset)

	router.GET("/resetpwd", h.ShowPasswordResetForm)
	router.POST("/resetpwd", h.ResetPassword)

	router.GET("/activateuser", h.ShowUserActivationPage)

	// 用户设置相关路由（需要鉴权）
	authGroup := router.Group("/u")
	authGroup.Use(h.auth.RequireAuth())
	{
		authGroup.GET("/settings", h.ShowUserSettingsPage)
		authGroup.POST("/settings/profile", h.UpdateProfile)
		authGroup.POST("/settings/email", h.UpdateEmail)
		authGroup.POST("/settings/password", h.UpdatePassword)
	}
}

// ===================== 页面渲染方法 =====================

// ShowLoginPage 显示登录页面
func (h *userController) ShowLoginPage(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	data := h.renderService.CreatePageData(c, "user.login.title")
	h.renderService.HTML(c, "passport/login", data)
}

// ShowRegisterPage 显示注册页面
func (h *userController) ShowRegisterPage(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	data := h.renderService.CreatePageData(c, "user.register.title")
	data.Data["SupportLangs"] = config.GetSupportLangs()
	data.Data["TimeZones"] = config.GetTimeZone()
	h.renderService.HTML(c, "passport/register", data)
}

// ShowPasswordResetPage 显示找回密码页面
func (h *userController) ShowPasswordResetPage(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	data := h.renderService.CreatePageData(c, "user.findpwd.title")
	h.renderService.HTML(c, "passport/findpwd", data)
}

// ShowPasswordResetForm 显示重置密码表单页面
func (h *userController) ShowPasswordResetForm(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	if uid == "" || sign == "" {
		data := h.renderService.CreatePageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "缺少必要参数"
		h.renderService.HTML(c, "error", data)
		return
	}

	// 解析用户ID
	userID, err := h.service.ParseUserID(uid)
	if err != nil {
		data := h.renderService.CreatePageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "用户ID格式错误"
		h.renderService.HTML(c, "error", data)
		return
	}

	// 检查验证码
	if valid, _ := h.service.CheckValidCode(userID, sign, config.VALID_CODE_RESET_PWD_EMAIL); !valid {
		data := h.renderService.CreatePageData(c, "error")
		data.Data["title"] = "验证码无效"
		data.Data["message"] = "验证码已过期或无效"
		h.renderService.HTML(c, "error.html", data)
		return
	}

	data := h.renderService.CreatePageData(c, "user.resetpwd.title")
	data.Data["uid"] = uid
	data.Data["sign"] = sign
	h.renderService.HTML(c, "passport/resetpwd", data)
}

// ShowUserActivationPage 显示激活用户页面
func (h *userController) ShowUserActivationPage(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	if uid == "" || sign == "" {
		data := h.renderService.CreatePageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "缺少必要参数"
		h.renderService.HTML(c, "error", data)
		return
	}

	// 激活用户
	if err := h.service.ActivateUser(uid, sign); err != nil {
		data := h.renderService.CreatePageData(c, "user.activate.failed.title")
		data.Data["message"] = "激活链接无效或已过期"
		h.renderService.HTML(c, "passport/activate_failed", data)
	} else {
		data := h.renderService.CreatePageData(c, "user.activate.success.title")
		data.Data["message"] = "账户激活成功，请登录"
		h.renderService.HTML(c, "passport/activate_success", data)
	}
}

// ShowUserSettingsPage 显示用户设置页面
func (h *userController) ShowUserSettingsPage(c *gin.Context) {
	// 读取用户认证信息
	authUser := h.auth.MustGetCurrentUser(c)

	// 获取用户设置信息
	userInfo, err := h.service.GetUserByUserId(authUser.UserId)
	if err != nil {
		data := h.renderService.CreatePageData(c, "error")
		data.Data["title"] = "错误"
		data.Data["message"] = "获取用户信息失败"
		h.renderService.HTML(c, "error", data)
		return
	}

	data := h.renderService.CreatePageData(c, "user.settings.title")
	data.Data["User"] = userInfo
	h.renderService.HTML(c, "member/settings", data)
}

// ===================== 认证相关方法 =====================

// Login 处理用户登录请求
func (h *userController) Login(c *gin.Context) {
	// 解析登录表单数据，创建 LoginRequest 对象
	req := &user.LoginRequest{}

	if err := c.Bind(req); err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 验证输入参数长度
	if len(req.LoginName) < 6 || len(req.Pwd) < 6 {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("用户名和密码长度至少6位"))
		return
	}

	// 验证用户名和密码，登录名统一转为小写
	req.LoginName = strings.ToLower(req.LoginName)
	user, err := h.service.LoginWithUsernameOrEmail(req.LoginName, req.Pwd)
	if err != nil {
		h.renderService.JSONError(c, errors.PasswordIncorrect.WithDetails(err.Error()))
		return
	}

	// 创建登录会话并设置Cookie
	authUser := &types.AuthUser{
		UserId:   user.UserId,
		UserName: user.UserName,
		FullName: user.FullName,
		Email:    user.Email,
		Lang:     user.Lang,
		TimeZone: user.TimeZone,
	}

	// 设置登录Cookie
	h.auth.SetLoginCookie(c, authUser)

	// 返回登录成功响应
	responseData := map[string]interface{}{
		"IsLogin": true,
		"UserId":  user.UserId,
	}
	h.renderService.JSONSuccess(c, responseData)
}

// Register 处理用户注册请求
func (h *userController) Register(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	// 解析注册表单数据
	user := &user.User{}
	if err := c.Bind(user); err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 设置客户端IP
	user.LastIp = c.ClientIP()

	// 调用服务层注册用户
	if regErr := h.service.RegisterUser(user); regErr != nil {
		h.renderService.JSONError(c, regErr)
		return
	}

	// 返回注册成功响应
	responseData := map[string]interface{}{
		"UserId": user.UserId,
	}
	h.renderService.JSONSuccess(c, responseData)
}

// Logout 处理登出请求
func (h *userController) Logout(c *gin.Context) {
	// 清除登录cookie
	pkg.DelCookie(c.Writer, "ud")
	c.Redirect(http.StatusFound, "/passport/login")
}

// ===================== 验证方法 =====================

// ValidateUsername 验证用户名是否可用
func (h *userController) ValidateUsername(c *gin.Context) {
	userName, _ := c.GetPostForm("username")

	// 验证用户名格式
	if !h.validateUserNameFormat(userName) {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("用户名格式不正确"))
		return
	}

	// 检查用户名是否已存在
	isExistUserName, err := h.service.IsExistUserName(userName)
	if err != nil {
		h.renderService.JSONError(c, err)
		return
	}
	if isExistUserName {
		h.renderService.JSONError(c, errors.UserExists.WithDetails("用户名: "+userName))
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// ValidateEmail 验证邮箱是否可用
func (h *userController) ValidateEmail(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !h.validateEmailFormat(email) {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("邮箱格式不正确"))
		return
	}

	// 检查邮箱是否已存在
	isExistEmail, err := h.service.IsExistEmail(email)
	if err != nil {
		h.renderService.JSONError(c, err)
		return
	}
	if isExistEmail {
		h.renderService.JSONError(c, errors.EmailExists.WithDetails("邮箱: "+email))
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// ===================== 密码管理方法 =====================

// RequestPasswordReset 处理找回密码请求
func (h *userController) RequestPasswordReset(c *gin.Context) {
	// 解析请求数据
	req := &user.FindPasswordRequest{}
	if err := c.Bind(req); err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 验证邮箱格式
	if !h.validateEmailFormat(req.Email) {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("邮箱格式不正确"))
		return
	}

	// 获取用户信息
	userInfo, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		h.renderService.JSONError(c, errors.UserNotFound.WithDetails("邮箱: "+req.Email))
		return
	}

	// 创建重置密码验证码
	if err := h.service.CreateResetPasswordCode(userInfo.UserId); err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	// 发送重置密码邮件
	if err := h.service.SendResetPasswordEmail(userInfo); err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// ResetPassword 处理重置密码请求
func (h *userController) ResetPassword(c *gin.Context) {
	// 解析重置密码请求
	req := &user.ResetPasswordRequest{}
	if err := c.Bind(req); err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 验证新密码格式
	if !h.validatePasswordFormat(req.NewPwd) {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("密码格式不正确"))
		return
	}

	// 解析用户ID
	userID, err := h.service.ParseUserID(req.UID)
	if err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("用户ID格式错误"))
		return
	}

	// 检查验证码
	valid, validCodeID := h.service.CheckValidCode(userID, req.Sign, config.VALID_CODE_RESET_PWD_EMAIL)
	if !valid {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("验证码无效或已过期"))
		return
	}

	// 重置密码
	if err := h.service.ResetUserPassword(userID, req.NewPwd); err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	// 删除验证码
	h.service.DeleteValidCode(validCodeID)

	h.renderService.JSONSuccess(c, nil)
}

// ===================== 用户设置相关方法 =====================

// UpdateProfile 更新用户基本资料
func (h *userController) UpdateProfile(c *gin.Context) {
	// 获取用户登录信息
	authUser := h.auth.MustGetCurrentUser(c)

	// 解析更新资料请求
	req := &user.UpdateProfileRequest{}
	if err := c.Bind(req); err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 设置用户ID,只能修改登录用户 profile
	req.UserId = authUser.UserId

	// 验证全名格式
	if req.FullName != "" && !h.validateFullNameFormat(req.FullName) {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("全名格式不正确"))
		return
	}

	// 调用service更新资料
	if err := h.service.UpdateUserProfile(req.UserId, req.FullName, req.Remark); err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	// 设置新的cookie
	authUser.FullName = req.FullName
	h.auth.SetLoginCookie(c, authUser)

	h.renderService.JSONSuccess(c, nil)
}

// UpdateEmail 更新用户邮箱
func (h *userController) UpdateEmail(c *gin.Context) {
	// 获取用户登录信息
	authUser := h.auth.MustGetCurrentUser(c)

	req := &user.UpdateEmailRequest{}

	req.UserId = authUser.UserId

	// 获取表单数据
	req.Email = c.PostForm("email")

	// 调用service更新邮箱
	if err := h.service.UpdateUserEmail(req.UserId, req.Email); err != nil {
		h.renderService.JSONError(c, errors.DatabaseError.WithDetails(err.Error()))
		return
	}

	h.renderService.JSONSuccess(c, gin.H{"message": "邮箱更新成功"})
}

// UpdatePassword 修改用户密码
func (h *userController) UpdatePassword(c *gin.Context) {
	// 读取用户认证信息
	authUser := h.auth.MustGetCurrentUser(c)

	// 绑定请求数据
	req := &user.ChangePasswordRequest{}
	if err := c.Bind(req); err != nil {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 判断输入密码是否正确
	if req.NewPwd != req.ConfirmPwd {
		h.renderService.JSONError(c, errors.ValidationError.WithDetails("两次输入密码不一致"))
		return
	}
	// 调用service更改密码
	err := h.service.ChangeUserPassword(authUser.UserId, req.CurrentPwd, req.NewPwd)
	if err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	// 更新 cookie
	h.auth.SetLoginCookie(c, authUser)

	h.renderService.JSONSuccess(c, gin.H{"message": "密码修改成功"})
}

// ===================== 辅助方法 =====================
// validateUserNameFormat 验证用户名格式
func (h *userController) validateUserNameFormat(userName string) bool {
	return pkg.IsValidUserName(userName)
}

// parseUserID 解析用户ID
func (h *userController) parseUserID(idStr string) (int64, error) {
	return h.service.ParseUserID(idStr)
}

// validateEmailFormat 验证邮箱格式
func (h *userController) validateEmailFormat(email string) bool {
	regexCheck := &pkg.RegexCheck{}
	return regexCheck.IsEmail(email)
}

// validatePasswordFormat 验证密码格式
func (h *userController) validatePasswordFormat(password string) bool {
	return pkg.IsValidPassword(password)
}

// validateFullNameFormat 验证全名格式
func (h *userController) validateFullNameFormat(fullName string) bool {
	return len(strings.TrimSpace(fullName)) > 0 && len(fullName) <= 50
}
