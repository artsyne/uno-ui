package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// Handler 用户HTTP处理层
type Handler struct {
	service *Service
	auth    *AuthHelper
}

// NewHandler 创建用户HTTP处理层实例
func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
		auth:    NewAuthHelper(),
	}
}

// RegisterRoutes 注册所有用户相关路由
// 保持与原系统完全一致的路由路径
func (h *Handler) RegisterRoutes(router *gin.Engine) {
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
func (h *Handler) RenderLogin(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data := pkg.DefaultH(c)
	data["title"] = "用户登录"
	c.HTML(http.StatusOK, "passport/login", data)
}

// RenderRegister 渲染注册页面
func (h *Handler) RenderRegister(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data := pkg.DefaultH(c)
	data["title"] = "用户注册"
	data["SupportLangs"] = configs.GetSupportLangs()
	data["TimeZones"] = configs.GetTimeZone()
	c.HTML(http.StatusOK, "passport/register", data)
}

// RenderFindPwd 渲染找回密码页面
func (h *Handler) RenderFindPwd(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data := pkg.DefaultH(c)
	data["title"] = "找回密码"
	c.HTML(http.StatusOK, "passport/findpwd", data)
}

// RenderResetPwd 渲染重置密码页面
func (h *Handler) RenderResetPwd(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	if uid == "" || sign == "" {
		data := pkg.DefaultH(c)
		data["title"] = "参数错误"
		data["message"] = "缺少必要参数"
		c.HTML(http.StatusBadRequest, "error", data)
		return
	}

	// 解析用户ID
	userID, err := h.parseUserID(uid)
	if err != nil {
		data := pkg.DefaultH(c)
		data["title"] = "参数错误"
		data["message"] = "用户ID格式错误"
		c.HTML(http.StatusBadRequest, "error", data)
		return
	}

	// 检查验证码
	if valid, _ := h.service.CheckValidCode(userID, sign, configs.VALID_CODE_RESET_PWD_EMAIL); !valid {
		data := pkg.DefaultH(c)
		data["title"] = "验证码无效"
		data["message"] = "验证码已过期或无效"
		c.HTML(http.StatusBadRequest, "error", data)
		return
	}

	data := pkg.DefaultH(c)
	data["title"] = "重置密码"
	data["uid"] = uid
	data["sign"] = sign
	c.HTML(http.StatusOK, "passport/resetpwd", data)
}

// RenderActivateUser 渲染激活用户页面
func (h *Handler) RenderActivateUser(c *gin.Context) {
	uid := c.Query("uid")
	sign := c.Query("sign")

	data := pkg.DefaultH(c)
	data["title"] = "账户激活"

	if uid == "" || sign == "" {
		data["message"] = "缺少必要参数"
		c.HTML(http.StatusBadRequest, "passport/activateuser", data)
		return
	}

	// 激活用户
	if h.service.ActivateUser(uid, sign) {
		data["message"] = "账户激活成功，请登录"
		c.HTML(http.StatusOK, "passport/activateuser", data)
	} else {
		data["message"] = "激活链接已过期或无效"
		c.HTML(http.StatusBadRequest, "passport/activateuser", data)
	}
}

// ===================== 请求处理方法 =====================

// HandleLogin 处理用户登录请求
func (h *Handler) HandleLogin(c *gin.Context) {
	// 解析登录表单数据
	login := &Login{}
	if err := c.Bind(login); err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 验证输入参数长度
	if len(login.LoginName) < 6 || len(login.Pwd) < 6 {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 验证用户名和密码，登录名统一转为小写
	login.LoginName = strings.ToLower(login.LoginName)
	user, err := h.service.CheckNameAndPassword(login.LoginName, login.Pwd)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 创建完整的登录会话数据
	loginSession := &Login{
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
	}

	// 序列化登录会话并加密存储到cookie
	if err := h.setLoginCookie(c, loginSession); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 返回登录成功响应
	data := pkg.DefaultH(c)
	data["IsLogin"] = true
	data["UserId"] = user.UserId
	c.JSON(http.StatusOK, data)
}

// HandleUserReg 处理用户注册请求
func (h *Handler) HandleUserReg(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := h.auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 解析注册表单数据
	user := &User{}
	if err := c.Bind(user); err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 设置客户端IP
	user.LastIp = c.ClientIP()

	// 调用服务层注册用户
	if regErr := h.service.RegisterUser(user); regErr != nil {
		errors.ResponseErr(c, regErr)
		return
	}

	// 返回注册成功响应
	data := pkg.DefaultH(c)
	data["UserId"] = user.UserId
	c.JSON(http.StatusOK, data)
}

// HandleCheckUserName 检查用户名是否已存在
func (h *Handler) HandleCheckUserName(c *gin.Context) {
	userName, _ := c.GetPostForm("username")

	// 验证用户名格式
	if !h.validateUserNameFormat(userName) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 检查用户名是否已存在
	if h.service.IsExistUserName(userName) != nil {
		errors.ResponseErr(c, errors.RegisterUserNameExists)
		return
	}

	data := pkg.ResponseH(c)
	c.JSON(http.StatusOK, data)
}

// HandleCheckEmail 检查邮箱地址是否已存在
func (h *Handler) HandleCheckEmail(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !h.validateEmailFormat(email) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 检查邮箱是否已存在
	if h.service.IsExistEmail(email) != nil {
		errors.ResponseErr(c, errors.RegisterUserEmailExists)
		return
	}

	data := pkg.ResponseH(c)
	c.JSON(http.StatusOK, data)
}

// HandleFindPwd 处理找回密码请求
func (h *Handler) HandleFindPwd(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !h.validateEmailFormat(email) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 根据邮箱获取用户信息
	user, err := h.service.GetUserByEmail(email)
	if user == nil || err != nil {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}

	// 创建并插入重置密码验证码
	if err := h.service.CreateResetPasswordCode(user.UserId); err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 发送重置密码邮件
	if err := h.service.SendResetPasswordEmail(user); err != nil {
		errors.ResponseErr(c, err)
		return
	}

	data := pkg.ResponseH(c)
	c.JSON(http.StatusOK, data)
}

// HandleResetPwd 处理重置密码请求
func (h *Handler) HandleResetPwd(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	sign, _ := c.GetPostForm("sign")
	newPwd, _ := c.GetPostForm("new_pwd")

	// 验证参数
	if uid == "" || sign == "" || len(newPwd) < 6 {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 解析用户ID
	userID, err := h.parseUserID(uid)
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 验证确认码
	result, validID := h.service.CheckValidCode(userID, sign, configs.VALID_CODE_RESET_PWD_EMAIL)
	if !result {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 获取用户信息以获取盐值
	user, err := h.service.GetUserByUserId(userID)
	if err != nil || user == nil {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}

	// 更新密码
	newPwdHash := pkg.Sha256(newPwd + user.PwdSalt)
	if rows, err := h.service.UpdateUserPwd(userID, newPwdHash); err != nil || rows != 1 {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}

	// 删除已使用的验证码
	if rows, err := h.service.DeleteValidCode(validID); err != nil || rows != 1 {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}

	c.JSON(http.StatusOK, pkg.ResponseH(c))
}

// HandleLoginOut 处理登出请求
func (h *Handler) HandleLoginOut(c *gin.Context) {
	// 清除登录cookie
	pkg.DelCookie(c.Writer, "ud")
	c.Redirect(http.StatusFound, "/passport/login")
}

// ===================== 私有辅助方法 =====================

// parseUserID 解析用户ID
func (h *Handler) parseUserID(idStr string) (int64, error) {
	return h.service.ParseUserID(idStr)
}

// validateEmailFormat 验证邮箱格式
func (h *Handler) validateEmailFormat(email string) bool {
	regexCheck := &pkg.RegexCheck{}
	return regexCheck.IsEmail(email)
}

// validateUserNameFormat 验证用户名格式
func (h *Handler) validateUserNameFormat(userName string) bool {
	return len(userName) >= 6 && len(userName) <= 50
}

// setLoginCookie 设置登录会话cookie
func (h *Handler) setLoginCookie(c *gin.Context, login *Login) error {
	// 序列化登录对象
	ud, err := json.Marshal(login)
	if err != nil {
		return err
	}

	// 加密cookie内容
	cookieUD, err := pkg.AESEncrypt(string(ud), configs.GetConfig().Cookie.Secret)
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
// 现在在路由级别已确保用户已登录
func (h *Handler) RenderUserSettings(c *gin.Context) {
	// 获取当前登录用户信息（由中间件保证已登录）
	loginUser := h.auth.MustGetCurrentUser(c)

	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	data := pkg.DefaultH(c)
	data["Title"] = "设置"
	data["Active"] = "setting"

	// 判断是否是用户本人
	isOwner := (loginUser.UserId == uid)

	// 如果是用户本人，获取详细的用户信息
	var userInfo *User
	if isOwner {
		userInfo, err = h.service.GetUserSettings(uid, loginUser.UserId)
		if err != nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	data["IsOwner"] = isOwner
	data["UserData"] = *loginUser
	data["UserInfo"] = userInfo

	c.HTML(http.StatusOK, "member/settings", data)
}

// UpdateProfile 更新用户基本资料
func (h *Handler) UpdateProfile(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 读取用户认证信息
	userData, exists := c.Get("UserData")
	if !exists {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}
	userLogin := userData.(Login)

	// 获取表单数据
	fullName := c.PostForm("full_name")
	remark := c.PostForm("remark")

	// 调用service更新资料
	err = h.service.UpdateUserProfile(uid, userLogin.UserId, fullName, remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "资料更新成功"})
}

// UpdateEmail 更新用户邮箱
func (h *Handler) UpdateEmail(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 读取用户认证信息
	userData, exists := c.Get("UserData")
	if !exists {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}
	userLogin := userData.(Login)

	// 获取表单数据
	email := c.PostForm("email")

	// 调用service更新邮箱
	err = h.service.UpdateUserEmail(uid, userLogin.UserId, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邮箱更新成功"})
}

// ChangePassword 修改用户密码
func (h *Handler) ChangePassword(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseUserID(c.Param("uid"))
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 读取用户认证信息
	userData, exists := c.Get("UserData")
	if !exists {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}
	userLogin := userData.(Login)

	// 获取表单数据
	currentPwd := c.PostForm("current_pwd")
	newPwd := c.PostForm("new_pwd")

	// 调用service更改密码
	err = h.service.ChangeUserPassword(uid, userLogin.UserId, currentPwd, newPwd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
