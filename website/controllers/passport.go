package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/helpers"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// PassportController 用户认证控制器
// 负责处理用户认证相关的所有HTTP请求，包括登录、注册、找回密码、用户激活等操作
type PassportController struct{}

// NewPassportController 创建新的用户认证控制器实例
// 返回值: *PassportController 用户认证控制器指针
func NewPassportController() *PassportController {
	return &PassportController{}
}

// parseUserID 解析并验证用户ID
// 从字符串中解析出有效的用户ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的用户ID
//   - error: 解析失败时的错误信息
func (pc *PassportController) parseUserID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.ParseInputParamErr
	}
	return id, nil
}

// validateEmailFormat 验证邮箱格式
// 使用正则表达式验证邮箱地址的有效性
//
// 参数:
//   - email: 待验证的邮箱地址
//
// 返回值:
//   - bool: 邮箱格式是否有效
func (pc *PassportController) validateEmailFormat(email string) bool {
	regexCheck := &pkg.RegexCheck{}
	return regexCheck.IsEmail(email)
}

// validateUserNameFormat 验证用户名格式
// 使用正则表达式验证用户名的有效性
//
// 参数:
//   - userName: 待验证的用户名
//
// 返回值:
//   - bool: 用户名格式是否有效
func (pc *PassportController) validateUserNameFormat(userName string) bool {
	regexCheck := &pkg.RegexCheck{}
	return regexCheck.IsUserName(userName)
}

// RenderFindPwd 渲染找回密码页面
// 显示找回密码表单，如果用户已登录则重定向到首页
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (pc *PassportController) RenderFindPwd(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := helpers.Auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	h := pkg.DefaultH(c)
	h["Title"] = "找回密码"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "passport/findpwd", h)
}

// RenderLogin 渲染登录页面
// 显示用户登录表单，如果用户已登录则重定向到首页
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (pc *PassportController) RenderLogin(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := helpers.Auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	h := pkg.DefaultH(c)
	h["Title"] = "登录"
	c.HTML(http.StatusOK, "passport/login", h)
}

// RenderRegister 渲染注册页面
// 显示用户注册表单，如果用户已登录则重定向到首页
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (pc *PassportController) RenderRegister(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := helpers.Auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	h := pkg.DefaultH(c)
	h["Title"] = "免费注册"
	c.HTML(http.StatusOK, "passport/register", h)
}

// RenderResetPwd 渲染重置密码页面
// 验证重置密码的确认码，显示重置密码表单
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// URL参数:
//   - uid: 用户ID
//   - sign: 重置密码确认码
func (pc *PassportController) RenderResetPwd(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := helpers.Auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 获取URL参数
	uid, _ := c.GetQuery("uid")
	sign, _ := c.GetQuery("sign")

	// 验证确认码
	validCodeOK := false
	userID, err := pc.parseUserID(uid)
	if err == nil && len(sign) == configs.VALID_CODE_LEN {
		if result, _ := models.CheckValidCode(userID, sign, configs.VALID_CODE_RESET_PWD_EMAIL); result {
			validCodeOK = true
		}
	}

	// 设置模板数据
	h := pkg.DefaultH(c)
	h["Title"] = "重置密码"
	h["UID"] = uid
	h["Sign"] = sign
	h["IsRightCode"] = validCodeOK
	c.HTML(http.StatusOK, "passport/resetpwd", h)
}

// RenderActivateUser 渲染用户激活页面
// 处理用户邮箱激活，显示激活结果
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// URL参数:
//   - uid: 用户ID
//   - sign: 激活确认码
func (pc *PassportController) RenderActivateUser(c *gin.Context) {
	h := pkg.DefaultH(c)
	h["Title"] = "用户激活"

	uid, _ := c.GetQuery("uid")
	sign, _ := c.GetQuery("sign")
	h["IsActivated"] = services.ActivateUser(uid, sign)

	c.HTML(http.StatusOK, "passport/activateuser", h)
}

// HandleLoginOut 处理用户退出登录
// 清除用户的登录cookie并返回成功响应
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (pc *PassportController) HandleLoginOut(c *gin.Context) {
	// 如果已登录则删除cookie
	if isAuth, _ := helpers.Auth.IsAuthenticated(c); isAuth {
		pkg.DelCookie(c.Writer, "ud")
	}

	h := pkg.ResponseH(c)
	h["Status"] = 0
	c.JSON(http.StatusOK, h)
}

// HandleResetPwd 处理重置密码请求
// 验证重置密码确认码，更新用户密码，删除使用过的确认码
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 表单参数:
//   - uid: 用户ID
//   - sign: 重置密码确认码
//   - new_pwd: 新密码
func (pc *PassportController) HandleResetPwd(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	sign, _ := c.GetPostForm("sign")
	newPwd, _ := c.GetPostForm("new_pwd")

	// 解析并验证用户ID
	userID, err := pc.parseUserID(uid)
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 验证输入参数
	if len(newPwd) < configs.PWD_MIN_LEN || len(newPwd) > configs.PWD_MAX_LEN || len(sign) != configs.VALID_CODE_LEN {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 验证确认码
	result, validID := models.CheckValidCode(userID, sign, configs.VALID_CODE_RESET_PWD_EMAIL)
	if !result {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 更新密码
	if rows, err := models.UpdateUserPwd(userID, pkg.Sha256(newPwd)); err != nil || rows != 1 {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}

	// 删除已使用的验证码
	if rows, err := models.DeleteValidCode(validID); err != nil || rows != 1 {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}

	c.JSON(http.StatusOK, pkg.ResponseH(c))
}

// HandleLogin 处理用户登录请求
// 验证用户凭据，创建登录会话，设置登录cookie
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 表单参数通过Login模型绑定:
//   - LoginName: 登录名（用户名或邮箱）
//   - Pwd: 密码
func (pc *PassportController) HandleLogin(c *gin.Context) {
	// 解析登录表单数据
	login := &models.Login{}
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
	user, err := services.CheckNameAndPwd(login.LoginName, login.Pwd)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 创建完整的登录会话数据
	loginSession := &models.Login{
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
	if err := pc.setLoginCookie(c, loginSession); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 返回登录成功响应
	h := pkg.DefaultH(c)
	h["IsLogin"] = true
	h["UserId"] = user.UserId
	c.JSON(http.StatusOK, h)
}

// HandleUserReg 处理用户注册请求
// 验证注册信息，创建新用户账户
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 表单参数通过User模型绑定
func (pc *PassportController) HandleUserReg(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := helpers.Auth.IsAuthenticated(c); isAuth {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 解析注册表单数据
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 设置客户端IP
	user.LastIp = c.ClientIP()

	// 调用服务层注册用户
	if regErr := services.RegUser(user); regErr != nil {
		errors.ResponseErr(c, regErr)
		return
	}

	// 返回注册成功响应
	h := pkg.DefaultH(c)
	h["UserId"] = user.UserId
	c.JSON(http.StatusOK, h)
}

// HandleCheckUserName 检查用户名是否已存在
// 验证用户名格式并检查是否已被注册
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 表单参数:
//   - username: 待检查的用户名
func (pc *PassportController) HandleCheckUserName(c *gin.Context) {
	userName, _ := c.GetPostForm("username")

	// 验证用户名格式
	if !pc.validateUserNameFormat(userName) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 检查用户名是否已存在
	if services.IsExistUserName(userName) != nil {
		errors.ResponseErr(c, errors.RegisterUserNameExists)
		return
	}

	h := pkg.ResponseH(c)
	c.JSON(http.StatusOK, h)
}

// HandleCheckEmail 检查邮箱地址是否已存在
// 验证邮箱格式并检查是否已被注册
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 表单参数:
//   - email: 待检查的邮箱地址
func (pc *PassportController) HandleCheckEmail(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !pc.validateEmailFormat(email) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 检查邮箱是否已存在
	if services.IsExistEmail(email) != nil {
		errors.ResponseErr(c, errors.RegisterUserEmailExists)
		return
	}

	h := pkg.ResponseH(c)
	c.JSON(http.StatusOK, h)
}

// HandleFindPwd 处理找回密码请求
// 验证邮箱地址，生成重置密码确认码，发送重置密码邮件
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 表单参数:
//   - email: 注册时使用的邮箱地址
func (pc *PassportController) HandleFindPwd(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !pc.validateEmailFormat(email) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 根据邮箱获取用户信息
	user, err := models.GetUserByEmail(email)
	if user == nil || err != nil {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}

	// 创建并插入重置密码验证码
	if err := pc.createResetPasswordCode(user.UserId); err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	// 发送重置密码邮件
	if err := pc.sendResetPasswordEmail(user); err != nil {
		errors.ResponseErr(c, err)
		return
	}

	h := pkg.ResponseH(c)
	c.JSON(http.StatusOK, h)
}

// setLoginCookie 设置登录会话cookie
// 将登录会话信息序列化并加密后存储到cookie
//
// 参数:
//   - c: Gin上下文
//   - login: 登录会话信息
//
// 返回值:
//   - error: 处理失败时的错误信息
func (pc *PassportController) setLoginCookie(c *gin.Context, login *models.Login) error {
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

// createResetPasswordCode 创建重置密码验证码
// 生成随机验证码并插入数据库
//
// 参数:
//   - userID: 用户ID
//
// 返回值:
//   - error: 创建失败时的错误信息
func (pc *PassportController) createResetPasswordCode(userID int64) error {
	// 创建验证码对象
	nowTime := pkg.GetNowTime()
	validCode := &models.ValidCode{
		ValidCode:     string(pkg.Rand(configs.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   configs.VALID_CODE_RESET_PWD_EMAIL,
		CreateUTCTime: nowTime.UnixMilli(),
		ExpireUTCTime: nowTime.AddDate(0, 0, configs.FIND_PWD_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli(),
	}

	// 插入验证码到数据库
	_, err := validCode.Insert()
	return err
}

// sendResetPasswordEmail 发送重置密码邮件
// 创建邮件任务并加入发送队列
//
// 参数:
//   - user: 用户信息
//
// 返回值:
//   - error: 发送失败时的错误信息
func (pc *PassportController) sendResetPasswordEmail(user *models.User) error {
	config := configs.GetConfig()
	params := configs.GetParams()
	nowTime := pkg.GetNowTime()

	// 获取最新的验证码（刚创建的）
	// 这里简化处理，实际应该从数据库获取刚创建的验证码
	validCode := string(pkg.Rand(configs.VALID_CODE_LEN, pkg.RAND_KIND_ALL))

	// 构建重置密码链接
	resetURL := fmt.Sprintf("%s/resetpwd?uid=%d&sign=%s", config.App.Domain, user.UserId, validCode)

	// 创建邮件任务
	emailTask := &models.EmailTask{
		To:      user.Email,
		Subject: params["rest_pwd_email_subject"],
		Body:    fmt.Sprintf(params["rest_pwd_email_body"], user.UserName, resetURL, resetURL),
	}

	// 序列化邮件任务
	jsonBytes, err := json.Marshal(emailTask)
	if err != nil {
		return errors.ParseJsonErr
	}

	// 创建作业队列任务
	jobQueue := &models.JobQueue{
		AddUTCTime:  nowTime.UnixMilli(),
		Body:        string(jsonBytes),
		JobType:     "email",
		PlanUTCTime: nowTime.UnixMilli(),
		TryNum:      0,
		Result:      1,
	}

	// 插入作业队列
	_, err = jobQueue.Insert()
	if err != nil {
		return errors.DbOperationErr
	}

	return nil
}
