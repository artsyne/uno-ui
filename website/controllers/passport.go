package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"fmt"
	"strconv"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
)

// RenderFindPwd : 渲染查找密码界面
func RenderFindPwd(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页

	h := pkg.DefaultH(c)
	h["Title"] = "找回密码"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "passport/findpwd", h)
}

// RenderLogin : 渲染登录界面
func RenderLogin(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页

	h := pkg.DefaultH(c)
	h["Title"] = "登录"
	c.HTML(http.StatusOK, "passport/login", h)
}

// RenderRegister : 渲染注册界面
func RenderRegister(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页
	h := pkg.DefaultH(c)
	h["Title"] = "免费注册"
	c.HTML(http.StatusOK, "passport/register", h)
}

// RenderResetPwd: 渲染重设密码界面
func RenderResetPwd(c *gin.Context) {
	IfLoginGotoIndex(c) // 如果登录跳转首页
	//获取参数
	uid, _ := c.GetQuery("uid")
	sign, _ := c.GetQuery("sign")

	validCodeOK := false
	userID, err := strconv.ParseInt(uid, 10, 64) //字符串转换
	//验证确认码
	if err == nil && len(sign) == configs.VALID_CODE_LEN {
		if result, _ := models.CheckValidCode(userID, sign, configs.VALID_CODE_RESET_PWD_EMAIL); result == true {
			validCodeOK = true
		}
	}

	//设置渲染返回对象
	h := pkg.DefaultH(c)
	h["Title"] = "重置密码"
	h["UID"] = uid
	h["Sign"] = sign
	h["IsRightCode"] = validCodeOK
	c.HTML(http.StatusOK, "passport/resetpwd", h)
}

// RenderActivateUser : 渲染激活页面
func RenderActivateUser(c *gin.Context) {
	h := pkg.DefaultH(c)
	h["Title"] = "用户激活"
	uid, _ := c.GetQuery("uid")
	sign, _ := c.GetQuery("sign")
	h["IsActivated"] = services.ActivateUser(uid, sign)
	c.HTML(http.StatusOK, "passport/activateuser", h)
}

// HandleLoginOut : 退出登录
func HandleLoginOut(c *gin.Context) {
	//如果登录删除cookie
	if IsLogin(c) {
		pkg.DelCookie(c.Writer, "ud")
	}

	h := pkg.ResponseH(c)
	h["Status"] = 0
	c.JSON(http.StatusOK, h)
}

// HandleResetPwd : 重置密码
func HandleResetPwd(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	sign, _ := c.GetPostForm("sign")
	newPwd, _ := c.GetPostForm("new_pwd")
	userID, err := strconv.ParseInt(uid, 10, 64) //字符串转换
	//检查确认码是否正确
	if err != nil || len(newPwd) < configs.PWD_MIN_LEN || len(newPwd) > configs.PWD_MAX_LEN ||
		len(sign) != configs.VALID_CODE_LEN || userID < 0 {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}
	//检查确认码是否正确
	result, validID := models.CheckValidCode(userID, sign, configs.VALID_CODE_RESET_PWD_EMAIL)
	if !result {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}
	//修改密码
	if rows, err := models.UpdateUserPwd(userID, pkg.Sha256(newPwd)); err != nil || rows != 1 {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}
	//删除验证码
	if rows, err := models.DeleteValidCode(validID); err != nil || rows != 1 {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}

	c.JSON(http.StatusOK, pkg.ResponseH(c))
}

// HandleLogin : 处理登录POST请求
func HandleLogin(c *gin.Context) {
	login := &models.Login{}
	err := c.Bind(login)
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	//判断参数是否到位
	if len(login.LoginName) < 6 || len(login.Pwd) < 6 {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}
	//判断用户名与密码是否正确，login name 可以是用户名或者邮件地址
	login.LoginName = strings.ToLower(login.LoginName) //登录名改为小写
	user, err := services.CheckNameAndPwd(login.LoginName, login.Pwd)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//创建用户Login数据
	login.UserId = user.UserId
	login.FullName = user.FullName
	login.LoginUTCTime = pkg.GetNowTimeUnixMilli()
	login.UserGroup = user.UserGroup
	login.UserName = user.UserName
	login.Pwd = user.Pwd
	login.TimeZone = user.TimeZone
	login.Lang = user.Lang
	login.Status = user.Status

	//序列化login对象
	ud, err := json.Marshal(login)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//加密cookie
	cookieUD, _ := pkg.AESEncrypt(string(ud), configs.GetConfig().Cookie.Secret)

	//写入cookie
	pkg.WriteCookie(c.Writer, "ud", cookieUD)

	//返回结果
	h := pkg.DefaultH(c)
	//设置返回值
	h["IsLogin"] = true
	h["UserId"] = user.UserId

	c.JSON(http.StatusOK, h)
}

// 处理用户注册请求
func HandleUserReg(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页

	// 序列化请求为user对象
	user := &models.User{}
	err := c.Bind(user)
	if err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	//获取客户端IP
	user.LastIp = c.ClientIP()

	// 注册用户
	regErr := services.RegUser(user)

	//判断注册是否成功,todo: 写regerr到日志
	if regErr != nil {
		errors.ResponseErr(c, regErr)
		return
	}

	h := pkg.DefaultH(c)
	h["UserId"] = user.UserId
	c.JSON(http.StatusOK, h)
}

// HandleFindPwd : 处理找回密码链接
func HandleFindPwd(c *gin.Context) {
	email, _ := c.GetPostForm("email")
	//创建正则表达式
	regexCheck := &pkg.RegexCheck{}

	//判断参数是否正确
	if !regexCheck.IsEmail(email) {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}
	//根据邮件地址获取用户信息
	user, err := models.GetUserByEmail(email)
	if user == nil || err != nil {
		errors.ResponseErr(c, errors.UserNotExists)
		return
	}

	//创建修改密码凭证
	//插入激活码
	validCode := &models.ValidCode{}
	validCode.ValidCode = string(pkg.Rand(configs.VALID_CODE_LEN, pkg.RAND_KIND_ALL))
	validCode.UserId = user.UserId
	validCode.ValidTypeId = configs.VALID_CODE_RESET_PWD_EMAIL
	nowTime := pkg.GetNowTime()
	validCode.CreateUTCTime = nowTime.UnixMilli()
	validCode.ExpireUTCTime = nowTime.AddDate(0, 0, configs.FIND_PWD_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli()
	if _, err := validCode.Insert(); err != nil {
		errors.ResponseErr(c, errors.InputParamErr)
		return
	}

	//插入发送邮件队列
	config := configs.GetConfig()
	params := configs.GetParams()
	//激活链接
	url := fmt.Sprintf("%s/resetpwd?uid=%d&sign=%s", config.App.Domain, user.UserId, validCode.ValidCode)

	//邮件task对象
	emailTask := &models.EmailTask{}
	emailTask.To = user.Email
	emailTask.Subject = params["rest_pwd_email_subject"]
	emailTask.Body = fmt.Sprintf(params["rest_pwd_email_body"], user.UserName, url, url)

	//插入邮件队列
	jobQueue := &models.JobQueue{}
	jobQueue.AddUTCTime = nowTime.UnixMilli()
	jsonBytes, err := json.Marshal(emailTask)
	if err != nil {
		errors.ResponseErr(c, errors.ParseJsonErr)
		return
	}
	jobQueue.Body = string(jsonBytes) //插入邮件任务json
	jobQueue.JobType = "email"
	jobQueue.PlanUTCTime = nowTime.UnixMilli()
	jobQueue.TryNum = 0
	jobQueue.Result = 1
	_, err = jobQueue.Insert()
	if err != nil {
		errors.ResponseErr(c, errors.DbOperationErr)
		return
	}
	h := pkg.ResponseH(c)
	c.JSON(http.StatusOK, h)
}

// 判断用户是否登录
func IsLogin(c *gin.Context) bool {
	isLogin, _ := c.Get("IsLogin")
	return isLogin == true
}

// 如果登录跳转到首页
func IfLoginGotoIndex(c *gin.Context) {
	if IsLogin(c) {
		c.Redirect(http.StatusFound, "/")
	}
}
