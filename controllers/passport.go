package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"fmt"
	"strconv"

	gobalerr "github.com/biuren/hebe/config"
	"github.com/biuren/hebe/helpers"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/system"
)

// RenderFindPwd : 渲染查找密码界面
func RenderFindPwd(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页

	h := helpers.DefaultH(c)
	h["Title"] = "找回密码"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "passport/findpwd", h)
}

// RenderLogin : 渲染登录界面
func RenderLogin(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页

	h := helpers.DefaultH(c)
	h["Title"] = "登录"
	c.HTML(http.StatusOK, "passport/login", h)
}

// RenderRegister : 渲染注册界面
func RenderRegister(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页
	h := helpers.DefaultH(c)
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
	if err == nil && len(sign) == system.VALID_CODE_LEN {
		if result, _ := models.CheckValidCode(userID, sign, system.VALID_CODE_RESET_PWD_EMAIL); result == true {
			validCodeOK = true
		}
	}

	//设置渲染返回对象
	h := helpers.DefaultH(c)
	h["Title"] = "重置密码"
	h["UID"] = uid
	h["Sign"] = sign
	h["IsRightCode"] = validCodeOK
	c.HTML(http.StatusOK, "passport/resetpwd", h)
}

// RenderActivateUser : 渲染激活页面
func RenderActivateUser(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "用户激活"
	uid, _ := c.GetQuery("uid")
	sign, _ := c.GetQuery("sign")
	h["IsActivated"] = ActiviteUser(uid, sign)
	c.HTML(http.StatusOK, "passport/activateuser", h)
}

// RenderPassport : 渲染登录页面
func RenderPassport(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Passport"
	c.HTML(http.StatusOK, "passport/index", h)
}

// PostLoginOut : 退出登录
func PostLoginOut(c *gin.Context) {
	//如果登录删除cookie
	if IsLogin(c) {
		helpers.DelCookie(c.Writer, "ud")
	}

	h := helpers.ResponseH(c)
	h["Status"] = 0
	c.JSON(http.StatusOK, h)
}

// PostResetPwd : 重置密码
func PostResetPwd(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	sign, _ := c.GetPostForm("sign")
	newPwd, _ := c.GetPostForm("new_pwd")
	userID, err := strconv.ParseInt(uid, 10, 64) //字符串转换
	//检查确认码是否正确
	if err != nil || len(newPwd) < system.PWD_MIN_LEN || len(newPwd) > system.PWD_MAX_LEN ||
		len(sign) != system.VALID_CODE_LEN || userID < 0 {
		gobalerr.ResponseErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}
	//检查确认码是否正确
	result, validID := models.CheckValidCode(userID, sign, system.VALID_CODE_RESET_PWD_EMAIL)
	if !result {
		gobalerr.ResponseErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}
	//修改密码
	if rows, err := models.UpdateUserPwd(userID, helpers.Sha256(newPwd)); err != nil || rows != 1 {
		gobalerr.ResponseErr(c, gobalerr.DB_OPERATION_ERR)
		return
	}
	//删除验证码
	if rows, err := models.DeleteValidCode(validID); err != nil || rows != 1 {
		gobalerr.ResponseErr(c, gobalerr.DB_OPERATION_ERR)
		return
	}

	c.JSON(http.StatusOK, helpers.ResponseH(c))
}

// PostLogin : 处理登录POST请求
func PostLogin(c *gin.Context) {
	login := &models.Login{}
	err := c.Bind(login)
	if err != nil {
		gobalerr.ResponseErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}

	//判断参数是否到位
	if len(login.LoginName) < 6 || len(login.Pwd) < 6 {
		gobalerr.ResponseErr(c, gobalerr.INPUT_PARAM_ERR)
		return
	}
	//判断用户名与密码是否正确，login name 可以是用户名或者邮件地址
	login.LoginName = strings.ToLower(login.LoginName) //登录名改为小写
	result, user, _ := CheckNameAndPwd(login.LoginName, login.Pwd)
	if result != 0 || user == nil {
		gobalerr.ResponseErr(c, result)
		return
	}

	//创建用户Login数据
	login.UserId = user.UserId
	login.AccessUrl = user.AccessUrl
	login.FaceSmallUrl = user.FaceSmallUrl
	login.LoginUTCTime = helpers.GetNowTimeUnixMilli()
	login.UserGroup = user.UserGroup
	login.UserName = user.UserName
	login.Pwd = user.Pwd
	login.TimeZone = user.TimeZone
	login.Lang = user.Lang
	login.Status = user.Status

	//序列化login对象
	ud, err := json.Marshal(login)
	if err != nil {
		gobalerr.ResponseErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}
	//加密cookie
	cookieUD, _ := helpers.TripleDesEncrypt(string(ud), system.GetConfig().Cookie.SecretKey)

	//写入cookie
	helpers.WriteCookie(c.Writer, "ud", cookieUD)

	//返回结果
	h := helpers.DefaultH(c)
	//设置返回值
	h["IsLogin"] = true
	h["UserId"] = user.UserId
	h["AccessUrl"] = user.AccessUrl

	c.JSON(http.StatusOK, h)
}

// PostRegister : 处理用户注册请求
func PostRegister(c *gin.Context) {
	IfLoginGotoIndex(c) //如果登录跳转首页

	// 序列化请求为user对象
	user := &models.User{}
	err := c.Bind(user)
	if err != nil {
		gobalerr.ResponseErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}

	//判断参数是否到位
	if len(user.UserName) < 6 || len(user.Pwd) < 6 {
		gobalerr.ResponseErr(c, gobalerr.INPUT_PARAM_ERR)
		return
	}

	//转换用户名，邮件地址为小写并去除空格
	user.UserName = strings.TrimSpace(strings.ToLower(user.UserName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	//判断用户时区与语言是否正确
	if !helpers.IsTimeZone(user.TimeZone) || !helpers.IsLang(user.Lang) {
		gobalerr.ResponseErr(c, gobalerr.INPUT_PARAM_ERR)
		return
	}

	//判断用户名是否存在
	if isExistName, err := models.IsExistUserName(user.UserName); err == nil && isExistName {
		gobalerr.ResponseErr(c, gobalerr.REGISTER_USER_NAME_EXISTS)
		return
	} else if err != nil {
		gobalerr.ResponseErrWithDesc(c, gobalerr.DB_OPERATION_ERR,
			gobalerr.EnumErrDesc[gobalerr.DB_OPERATION_ERR]+","+err.Error())
		return
	}

	//判断邮件地址是否存在
	if isExistEmail, err := models.IsExistsEmail(user.Email); err == nil && isExistEmail {
		gobalerr.ResponseErr(c, gobalerr.REGISTER_USER_EMAIL_EXISTS)
		return
	} else if err != nil {
		gobalerr.ResponseErrWithDesc(c, gobalerr.DB_OPERATION_ERR,
			gobalerr.EnumErrDesc[gobalerr.DB_OPERATION_ERR]+","+err.Error())
		return
	}

	//初始化用户信息
	//随机生成salt，保存密码为pwd+salt的sha256
	user.PwdSalt = string(helpers.Rand(64, helpers.RAND_KIND_ALL))
	user.Pwd = helpers.Sha256(user.Pwd + user.PwdSalt)
	user.UserGroup = system.DEFAULT_USER_GROUP //用户分组为1
	nowTime := helpers.GetNowTimeUnixMilli()
	user.RegUTCTime = nowTime
	user.RegIp = c.ClientIP()
	user.LastLoginUTCTime = nowTime
	user.LastIp = c.ClientIP()
	user.Status = system.USER_STATUS_WAIT_ACTIVE

	//插入数据库
	uid, err := user.Insert()
	if err != nil {
		gobalerr.ResponseErrWithDesc(c, gobalerr.DB_OPERATION_ERR,
			gobalerr.EnumErrDesc[gobalerr.DB_OPERATION_ERR]+","+err.Error())
		return
	}

	//更新用户访问地址
	accessUrl:= fmt.Sprintf("/u/%d", uid)
	rows, err := models.UpdateAccessUrlByUserId(uid, accessUrl)
	if err != nil || rows == 0 {
		gobalerr.ResponseErrWithDesc(c, gobalerr.DB_OPERATION_ERR,
			gobalerr.EnumErrDesc[gobalerr.DB_OPERATION_ERR]+","+err.Error())
		return
	}

	//发送激活邮件
	AddActivateUserEmailTask(user)

	h := helpers.DefaultH(c)
	h["UserId"] = uid
	h["AccessUrl"] = accessUrl
	c.JSON(http.StatusOK, h)
}

// PostFindPwd : 处理找回密码链接
func PostFindPwd(c *gin.Context) {
	email, _ := c.GetPostForm("email")
	//创建正则表达式
	regexCheck := &helpers.RegexCheck{}

	//判断参数是否正确
	if !regexCheck.IsEmail(email) {
		gobalerr.ResponseErr(c, gobalerr.INPUT_PARAM_ERR)
		return
	}
	//根据邮件地址获取用户信息
	user, err := models.GetUserByEmail(email)
	if user == nil || err != nil {
		gobalerr.ResponseErr(c, gobalerr.USER_NOT_EXISTS)
		return
	}

	//创建修改密码凭证
	//插入激活码
	validCode := &models.ValidCode{}
	validCode.ValidCode = string(helpers.Rand(system.VALID_CODE_LEN, helpers.RAND_KIND_ALL))
	validCode.UserId = user.UserId
	validCode.ValidTypeId = system.VALID_CODE_RESET_PWD_EMAIL
	nowTime := helpers.GetNowTime()
	validCode.CreateUTCTime = nowTime.UnixMilli()
	validCode.ExpireUTCTime = nowTime.AddDate(0, 0, system.FIND_PWD_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli()
	if _, err := validCode.Insert(); err != nil {
		gobalerr.ResponseErr(c, gobalerr.INPUT_PARAM_ERR)
		return
	}

	//插入发送邮件队列
	config := system.GetConfig()
	params := system.GetParams()
	//激活链接
	url := fmt.Sprintf("%s/resetpwd?uid=%d&sign=%s", config.Domain, user.UserId, validCode.ValidCode)

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
		gobalerr.ResponseErr(c, gobalerr.PARSE_JSON_ERR)
		return
	}
	jobQueue.Body = string(jsonBytes) //插入邮件任务json
	jobQueue.JobType = system.JOB_TYPE_EMAIL
	jobQueue.PlanUTCTime = nowTime.UnixMilli()
	jobQueue.TryNum = 0
	jobQueue.Result = 1
	_, err = jobQueue.Insert()
	if err != nil {
		gobalerr.ResponseErr(c, gobalerr.DB_OPERATION_ERR)
		return
	}
	h := helpers.ResponseH(c)
	c.JSON(http.StatusOK, h)
}

// AddActivateUserEmailTask : 增加发送激活用户邮件任务
func AddActivateUserEmailTask(user *models.User) error {
	//插入激活码
	validCode := &models.ValidCode{}
	validCode.ValidCode = string(helpers.Rand(system.VALID_CODE_LEN, helpers.RAND_KIND_ALL))
	validCode.UserId = user.UserId
	validCode.ValidTypeId = system.VALID_CODE_ACTIVE_USER_EMAIL
	nowTime := helpers.GetNowTime()
	validCode.CreateUTCTime = nowTime.UnixMilli()
	validCode.ExpireUTCTime = nowTime.AddDate(0, 0, system.REG_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli() //增加一天过期
	if _, err := validCode.Insert(); err != nil {
		return err
	}

	//插入发送邮件队列
	config := system.GetConfig()
	params := system.GetParams()
	//激活链接
	url := fmt.Sprintf("%s/activateuser?uid=%d&sign=%s", config.Domain, user.UserId, validCode.ValidCode)

	//邮件task对象
	emailTask := &models.EmailTask{}
	emailTask.To = user.Email
	emailTask.Subject = params["reg_email_subject"]
	emailTask.Body = fmt.Sprintf(params["reg_email_body"], user.UserName, url, url)

	//插入邮件队列
	jobQueue := &models.JobQueue{}
	jobQueue.AddUTCTime = nowTime.UnixMilli()
	jsonBytes, err := json.Marshal(emailTask)
	if err != nil {
		return err
	}
	jobQueue.Body = string(jsonBytes) //插入邮件任务json
	jobQueue.JobType = system.JOB_TYPE_EMAIL
	jobQueue.PlanUTCTime = nowTime.UnixMilli()
	jobQueue.TryNum = 0
	jobQueue.Result = 1
	_, err = jobQueue.Insert()

	return err
}

// CheckNameAndPwd : 检查用户名与密码是否正确
func CheckNameAndPwd(name string, pwd string) (int, *models.User, error) {
	//根据用户名获取用户信息
	user, err := models.GetUserByUserName(name)

	if err != nil {
		return gobalerr.USER_NOT_EXISTS, nil, err
	}

	//如果用户不存在，用邮件地址查找用户信息
	if user == nil {
		user, err = models.GetUserByEmail(name)
		if err != nil || user == nil {
			return gobalerr.USER_NOT_EXISTS, nil, err
		}
	}

	//判断密码是否正确,密码为pwd+salt的sha256
	if user.Pwd != helpers.Sha256(pwd+user.PwdSalt) {
		return gobalerr.USER_PWD_ERR, nil, nil
	}

	return 0, user, nil
}

// ActiviteUser : 激活用户
func ActiviteUser(uid string, sign string) bool {
	result := false                              //返回结果
	userID, err := strconv.ParseInt(uid, 10, 64) //字符串转换
	if err != nil {
		return false
	}
	//获取验证码
	if result, validID := models.CheckValidCode(userID, sign, system.VALID_CODE_ACTIVE_USER_EMAIL); result == true {
		//更新用户状态
		if rows, err := models.UpdateUserStatus(userID, system.USER_STATUS_OK); err != nil || rows != 1 {
			return false
		}
		//删除验证码
		if rows, err := models.DeleteValidCode(validID); err != nil || rows != 1 {
			return false
		}
		result = true //通过所有校验
	}

	return result
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
