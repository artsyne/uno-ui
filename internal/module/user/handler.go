package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/core/config"
	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// HandlerInterface 用户HTTP处理层接口
// 定义所有与用户相关的HTTP处理操作，包括认证、注册、设置等
type HandlerInterface interface {
	// RegisterRoutes 注册所有用户相关路由
	RegisterRoutes(router *gin.Engine)

	// ===================== 页面渲染方法 =====================

	// RenderLogin 渲染登录页面
	RenderLogin(c *gin.Context)

	// RenderRegister 渲染注册页面
	RenderRegister(c *gin.Context)

	// RenderFindPwd 渲染找回密码页面
	RenderFindPwd(c *gin.Context)

	// RenderResetPwd 渲染重置密码页面
	RenderResetPwd(c *gin.Context)

	// RenderActivateUser 渲染激活用户页面
	RenderActivateUser(c *gin.Context)

	// RenderUserSettings 渲染用户设置页面
	RenderUserSettings(c *gin.Context)

	// ===================== 请求处理方法 =====================

	// HandleLogin 处理用户登录请求
	HandleLogin(c *gin.Context)

	// HandleUserReg 处理用户注册请求
	HandleUserReg(c *gin.Context)

	// HandleCheckUserName 检查用户名是否存在
	HandleCheckUserName(c *gin.Context)

	// HandleCheckEmail 检查邮箱是否存在
	HandleCheckEmail(c *gin.Context)

	// HandleFindPwd 处理找回密码请求
	HandleFindPwd(c *gin.Context)

	// HandleResetPwd 处理重置密码请求
	HandleResetPwd(c *gin.Context)

	// HandleLoginOut 处理用户退出登录请求
	HandleLoginOut(c *gin.Context)

	// ===================== 用户设置相关方法 =====================

	// UpdateProfile 更新用户基本资料
	UpdateProfile(c *gin.Context)

	// UpdateEmail 更新用户邮箱
	UpdateEmail(c *gin.Context)

	// ChangePassword 修改用户密码
	ChangePassword(c *gin.Context)
}

// handler 用户HTTP处理层实现
type handler struct {
	service       ServiceInterface
	auth          AuthHelperInterface
	config        *config.Config // 缓存配置，避免重复调用全局变量
	renderService render.RenderServiceInterface
}

// NewHandlerWithServices 通过DI容器创建用户HTTP处理层实例
func NewHandlerWithServices(renderService render.RenderServiceInterface, _ render.RenderServiceInterface) HandlerInterface {
	return &handler{
		service:       NewService(),
		auth:          NewAuthHelper(),
		config:        config.GetConfig(),
		renderService: renderService,
	}
}

// RegisterRoutes 注册所有用户相关路由
// 保持与原系统完全一致的路由路径
func (h *handler) RegisterRoutes(router *gin.Engine) {
	// 用户认证相关路由（无需鉴权）
	router.GET("/register", h.RenderRegister)
	router.POST("/register", h.HandleUserReg)
	router.POST("/check-username", h.HandleCheckUserName)
	router.POST("/check-email", h.HandleCheckEmail)

	router.GET("/login", h.RenderLogin)
	router.POST("/login", h.HandleLogin)

	router.POST("/loginout", h.HandleLoginOut)

	router.GET("/findpwd", h.RenderFindPwd)
	router.POST("/findpwd", h.HandleFindPwd)

	router.GET("/resetpwd", h.RenderResetPwd)
	router.POST("/resetpwd", h.HandleResetPwd)

	router.GET("/activateuser", h.RenderActivateUser)

	// 用户设置相关路由（需要鉴权）
	authGroup := router.Group("/u")
	authGroup.Use(AuthRequired())
	{
		authGroup.GET("/:uid/settings", h.RenderUserSettings)
		authGroup.POST("/:uid/settings/profile", h.UpdateProfile)
		authGroup.POST("/:uid/settings/email", h.UpdateEmail)
		authGroup.POST("/:uid/settings/password", h.ChangePassword)
	}
}

// ===================== 页面渲染方法 =====================

// RenderLogin 渲染登录页面
func (h *handler) RenderLogin(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	data := h.renderService.PageData(c, "user.login.title")
	h.renderService.HTML(c, "passport/login", data)
}

// RenderRegister 渲染注册页面
func (h *handler) RenderRegister(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	data := h.renderService.PageData(c, "user.register.title")
	data.Data["SupportLangs"] = config.GetSupportLangs()
	data.Data["TimeZones"] = config.GetTimeZone()
	h.renderService.HTML(c, "passport/register", data)
}

// RenderFindPwd 渲染找回密码页面
func (h *handler) RenderFindPwd(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	data := h.renderService.PageData(c, "user.findpwd.title")
	h.renderService.HTML(c, "passport/findpwd", data)
}

// RenderResetPwd 渲染重置密码页面
func (h *handler) RenderResetPwd(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	if uid == "" || sign == "" {
		data := h.renderService.PageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "缺少必要参数"
		h.renderService.HTML(c, "error.html", data)
		return
	}

	// 解析用户ID
	userID, err := h.parseUserID(uid)
	if err != nil {
		data := h.renderService.PageData(c, "error")
		data.Data["title"] = "参数错误"
		data.Data["message"] = "用户ID格式错误"
		h.renderService.HTML(c, "error.html", data)
		return
	}

	// 检查验证码
	if valid, _ := h.service.CheckValidCode(userID, sign, config.VALID_CODE_RESET_PWD_EMAIL); !valid {
		data := h.renderService.PageData(c, "error")
		data.Data["title"] = "验证码无效"
		data.Data["message"] = "验证码已过期或无效"
		h.renderService.HTML(c, "error.html", data)
		return
	}

	data := h.renderService.PageData(c, "user.resetpwd.title")
	data.Data["uid"] = uid
	data.Data["sign"] = sign
	h.renderService.HTML(c, "passport/resetpwd.html", data)
}

// RenderActivateUser 渲染激活用户页面
func (h *handler) RenderActivateUser(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	data := h.renderService.PageData(c, "user.activate.title")

	if uid == "" || sign == "" {
		data.Data["message"] = "缺少必要参数"
		h.renderService.HTML(c, "passport/activateuser.html", data)
		return
	}

	// 激活用户
	if h.service.ActivateUser(uid, sign) {
		data.Data["message"] = "账户激活成功，请登录"
		h.renderService.HTML(c, "passport/activateuser.html", data)
	} else {
		data.Data["message"] = "激活链接已过期或无效"
		h.renderService.HTML(c, "passport/activateuser.html", data)
	}
}

// ===================== 请求处理方法 =====================

// HandleLogin 处理用户登录请求
func (h *handler) HandleLogin(c *gin.Context) {
	// 解析登录表单数据
	login := &AuthUser{}
	if err := c.Bind(login); err != nil {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails(err.Error()))
		return
	}

	// 验证输入参数长度
	if len(login.LoginName) < 6 || len(login.Pwd) < 6 {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails("用户名和密码长度至少6位"))
		return
	}

	// 验证用户名和密码，登录名统一转为小写
	login.LoginName = strings.ToLower(login.LoginName)
	user, err := h.service.CheckNameAndPassword(login.LoginName, login.Pwd)
	if err != nil {
		h.renderService.JSONError(c, errors.PasswordIncorrectErr().WithDetails(err.Error()))
		return
	}

	// 创建完整的登录会话数据
	loginSession := &AuthUser{
		UserId:       user.UserId,
		FullName:     user.FullName,
		LoginUTCTime: pkg.GetNowTimeUnixMilli(),
		UserGroup:    user.UserGroup,
		UserName:     user.UserName,
		Pwd:          user.Pwd,
		TimeZone:     user.TimeZone,
		Lang:         user.Lang,
		Status:       user.Status,
		LoginName:    login.LoginName,
		FaceUrl: user.FaceUrl,
	}

	// 序列化登录会话并加密存储到cookie
	if err := h.setLoginCookie(c, loginSession); err != nil {
		h.renderService.JSONError(c, errors.InternalErr().WithDetails(err.Error()))
		return
	}

	// 返回登录成功响应
	responseData := map[string]interface{}{
		"IsLogin": true,
		"UserId":  user.UserId,
	}
	h.renderService.JSONSuccess(c, responseData)
}

// HandleUserReg 处理用户注册请求
func (h *handler) HandleUserReg(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		h.renderService.Redirect(c, "/")
		return
	}

	// 解析注册表单数据
	user := &User{}
	if err := c.Bind(user); err != nil {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails(err.Error()))
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

// HandleCheckUserName 检查用户名是否已存在
func (h *handler) HandleCheckUserName(c *gin.Context) {
	userName, _ := c.GetPostForm("username")

	// 验证用户名格式
	if !h.validateUserNameFormat(userName) {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails("用户名格式不正确"))
		return
	}

	// 检查用户名是否已存在
	if h.service.IsExistUserName(userName) != nil {
		h.renderService.JSONError(c, errors.UserExistsErr().WithDetails("用户名: "+userName))
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// HandleCheckEmail 检查邮箱地址是否已存在
func (h *handler) HandleCheckEmail(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !h.validateEmailFormat(email) {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails("邮箱格式不正确"))
		return
	}

	// 检查邮箱是否已存在
	if h.service.IsExistEmail(email) != nil {
		h.renderService.JSONError(c, errors.EmailExistsErr().WithDetails("邮箱: "+email))
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// HandleFindPwd 处理找回密码请求
func (h *handler) HandleFindPwd(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !h.validateEmailFormat(email) {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails("邮箱格式不正确"))
		return
	}

	// 根据邮箱获取用户信息
	user, err := h.service.GetUserByEmail(email)
	if user == nil || err != nil {
		h.renderService.JSONError(c, errors.UserNotFoundErr().WithDetails(err.Error()))
		return
	}

	// 创建并插入重置密码验证码
	if err := h.service.CreateResetPasswordCode(user.UserId); err != nil {
		h.renderService.JSONError(c, errors.InternalErr().WithDetails(err.Error()))
		return
	}

	// 发送重置密码邮件
	if err := h.service.SendResetPasswordEmail(user); err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// HandleResetPwd 处理重置密码请求
func (h *handler) HandleResetPwd(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	sign, _ := c.GetPostForm("sign")
	newPwd, _ := c.GetPostForm("new_pwd")

	// 验证参数
	if uid == "" || sign == "" || len(newPwd) < 6 {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails("参数不完整或密码长度不足6位"))
		return
	}

	// 解析用户ID
	userID, err := h.parseUserID(uid)
	if err != nil {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails(err.Error()))
		return
	}

	// 验证确认码
	result, validID := h.service.CheckValidCode(userID, sign, config.VALID_CODE_RESET_PWD_EMAIL)
	if !result {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails("验证码无效"))
		return
	}

	// 获取用户信息以获取盐值
	user, err := h.service.GetUserByUserId(userID)
	if err != nil || user == nil {
		h.renderService.JSONError(c, errors.UserNotFoundErr().WithDetails(err.Error()))
		return
	}

	// 更新密码
	newPwdHash := pkg.Sha256(newPwd + user.PwdSalt)
	if rows, err := h.service.UpdateUserPwd(userID, newPwdHash); err != nil || rows != 1 {
		h.renderService.JSONError(c, errors.DatabaseErr().WithDetails(err.Error()))
		return
	}

	// 删除已使用的验证码
	if rows, err := h.service.DeleteValidCode(validID); err != nil || rows != 1 {
		h.renderService.JSONError(c, errors.DatabaseErr().WithDetails(err.Error()))
		return
	}

	h.renderService.JSONSuccess(c, nil)
}

// HandleLoginOut 处理登出请求
func (h *handler) HandleLoginOut(c *gin.Context) {
	// 清除登录cookie
	pkg.DelCookie(c.Writer, "ud")
	c.Redirect(http.StatusFound, "/passport/login")
}

// ===================== 私有辅助方法 =====================

// parseUserID 解析用户ID
func (h *handler) parseUserID(idStr string) (int64, error) {
	return h.service.ParseUserID(idStr)
}

// validateEmailFormat 验证邮箱格式
func (h *handler) validateEmailFormat(email string) bool {
	regexCheck := &pkg.RegexCheck{}
	return regexCheck.IsEmail(email)
}

// validateUserNameFormat 验证用户名格式
func (h *handler) validateUserNameFormat(userName string) bool {
	return len(userName) >= 6 && len(userName) <= 50
}

// setLoginCookie 设置登录会话cookie
func (h *handler) setLoginCookie(c *gin.Context, login *AuthUser) error {
	// 序列化登录对象
	ud, err := json.Marshal(login)
	if err != nil {
		return err
	}

	// 加密cookie内容
	cookieUD, err := pkg.AESEncrypt(string(ud), h.config.Cookie.Secret)
	if err != nil {
		return err
	}

	// 写入cookie
	pkg.WriteCookie(c.Writer, "ud", cookieUD)
	return nil
}

// ===================== Member 相关处理方法 =====================

// RenderUserSettings 渲染用户设置页面
// 从website/controllers/member.go中迁移的RenderMemberSettings方法
func (h *handler) RenderUserSettings(c *gin.Context) {
	// 确保用户已登录
	isAuth, authUser := h.auth.IsAuthenticated(c)
	if !isAuth {
		h.renderService.Redirect(c, "/passport/login")
		return
	}

	// 获取用户详细信息
	user, err := h.service.GetUserByUserId(authUser.UserId)
	if err != nil || user == nil {
		h.renderService.JSONError(c, errors.UserNotFoundErr().WithDetails(err.Error()))
		return
	}

	// 构建页面数据
	data := h.renderService.PageData(c, "user.settings")
	data.Data["User"] = user
	data.Data["SupportLangs"] = config.GetSupportLangs()
	data.Data["TimeZones"] = config.GetTimeZone()

	h.renderService.HTML(c, "member/settings", data)
}

// UpdateProfile 更新用户基本资料
func (h *handler) UpdateProfile(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails(err.Error()))
		return
	}

	// 读取用户认证信息
	userData, exists := c.Get("UserData")
	if !exists {
		h.renderService.JSONError(c, errors.UnauthorizedErr().WithDetails("用户认证信息不存在"))
		return
	}
	userLogin := userData.(AuthUser)

	// 获取表单数据
	fullName := c.PostForm("full_name")
	remark := c.PostForm("remark")

	// 调用service更新资料
	err = h.service.UpdateUserProfile(uid, userLogin.UserId, fullName, remark)
	if err != nil {
		h.renderService.JSONError(c, errors.DatabaseErr().WithDetails(err.Error()))
		return
	}

	h.renderService.JSONSuccess(c, gin.H{"message": "资料更新成功"})
}

// UpdateEmail 更新用户邮箱
func (h *handler) UpdateEmail(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails(err.Error()))
		return
	}

	// 读取用户认证信息
	userData, exists := c.Get("UserData")
	if !exists {
		h.renderService.JSONError(c, errors.UnauthorizedErr().WithDetails("用户认证信息不存在"))
		return
	}
	userLogin := userData.(AuthUser)

	// 获取表单数据
	email := c.PostForm("email")

	// 调用service更新邮箱
	err = h.service.UpdateUserEmail(uid, userLogin.UserId, email)
	if err != nil {
		h.renderService.JSONError(c, errors.DatabaseErr().WithDetails(err.Error()))
		return
	}

	h.renderService.JSONSuccess(c, gin.H{"message": "邮箱更新成功"})
}

// ChangePassword 修改用户密码
func (h *handler) ChangePassword(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		h.renderService.JSONError(c, errors.ValidationErr().WithDetails(err.Error()))
		return
	}

	// 读取用户认证信息
	userData, exists := c.Get("UserData")
	if !exists {
		h.renderService.JSONError(c, errors.UnauthorizedErr().WithDetails("用户认证信息不存在"))
		return
	}
	userLogin := userData.(AuthUser)

	// 获取表单数据
	currentPwd := c.PostForm("current_pwd")
	newPwd := c.PostForm("new_pwd")

	// 调用service更改密码
	err = h.service.ChangeUserPassword(uid, userLogin.UserId, currentPwd, newPwd)
	if err != nil {
		h.renderService.JSONError(c, errors.DatabaseErr().WithDetails(err.Error()))
		return
	}

	h.renderService.JSONSuccess(c, gin.H{"message": "密码修改成功"})
}
