package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
	gobalerr "github.com/biuren/hebe/utils"
)

// 根据用户的语言偏好和系统支持的语言列表，返回匹配的系统支持语言。
// 该函数首先将用户的语言偏好转换为小写，然后进行以下判断：
// 1. 如果用户偏好的语言是简体中文，则返回 "zh-cn"；
// 2. 如果用户偏好的语言是香港或台湾的本地语言，则返回 "zh-tw"；
// 3. 如果用户偏好的语言在系统支持的语言列表中，则返回该语言；
// 4. 如果用户偏好的语言不在系统支持的语言列表中，则返回默认语言 "en"。
// 参数：
// language: 用户的语言偏好，通常来自浏览器的 Accept-Language 头部；
// supportedLanguages: 系统支持的语言列表，来自配置文件。
func GetSupportedLang(language string, supportedLanguages []settings.SupportLangs) string {
	language = strings.ToLower(language)

	switch {
	case strings.Contains(language, "zh-cn"), strings.Contains(language, "zh-hans"):
		return "zh-cn"
	case strings.Contains(language, "zh-hk"), strings.Contains(language, "zh-tw"), strings.Contains(language, "zh-hant"):
		return "zh-tw"
	default:
		language = language[0:2]
		for _, supportedLanguage := range supportedLanguages {
			if supportedLanguage.Code == language {
				return language
			}
		}
		return "en"
	}
}

// AddActivateUserEmailTask : 增加发送激活用户邮件任务
func AddActivateUserEmailTask(user *models.User) error {
	//插入激活码
	validCode := &models.ValidCode{}
	validCode.ValidCode = string(utils.Rand(settings.VALID_CODE_LEN, utils.RAND_KIND_ALL))
	validCode.UserId = user.UserId
	validCode.ValidTypeId = settings.VALID_CODE_ACTIVE_USER_EMAIL
	nowTime := utils.GetNowTime()
	validCode.CreateUTCTime = nowTime.UnixMilli()
	validCode.ExpireUTCTime = nowTime.AddDate(0, 0, settings.REG_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli() //增加一天过期
	if _, err := validCode.Insert(); err != nil {
		return err
	}

	//插入发送邮件队列
	config := settings.GetConfig()
	params := settings.GetParams()
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
	jobQueue.JobType = settings.JOB_TYPE_EMAIL
	jobQueue.PlanUTCTime = nowTime.UnixMilli()
	jobQueue.TryNum = 0
	jobQueue.Result = 1
	_, err = jobQueue.Insert()

	return err
}

// 用户注册
// 返回值：0 表示成功，非0为错误码
func RegUser(user *models.User) int {

	//转换用户名，邮件地址为小写并去除空格
	user.UserName = strings.TrimSpace(strings.ToLower(user.UserName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	//判断参数是否到位
	if len(user.UserName) < 6 || len(user.Pwd) < 6 {
		return gobalerr.INPUT_PARAM_ERR
	}

	//判断用户时区与语言是否正确
	if !utils.IsTimeZone(user.TimeZone) {
		return gobalerr.INPUT_PARAM_ERR
	}

	//判断用户名是否存在
	if isExistName, err := models.IsExistUserName(user.UserName); err == nil && isExistName {
		return gobalerr.REGISTER_USER_NAME_EXISTS
	} else if err != nil {
		// db error 需要写log
		// TODO
		return gobalerr.DB_OPERATION_ERR
	}

	//判断邮件地址是否存在
	if isExistEmail, err := models.IsExistsEmail(user.Email); err == nil && isExistEmail {
		return gobalerr.REGISTER_USER_EMAIL_EXISTS
	} else if err != nil {
		// db error 需要写log
		// TODO
		return gobalerr.DB_OPERATION_ERR
	}

	//初始化用户信息
	//随机生成salt，保存密码为pwd+salt的sha256
	user.PwdSalt = string(utils.Rand(64, utils.RAND_KIND_ALL))
	user.Pwd = utils.Sha256(user.Pwd + user.PwdSalt)
	user.UserGroup = settings.DEFAULT_USER_GROUP //用户分组为1
	nowTime := utils.GetNowTimeUnixMilli()
	user.RegUTCTime = nowTime
	user.LastLoginUTCTime = nowTime
	user.Lang = GetSupportedLang(user.Lang, settings.GetConfig().SupportLangs)
	user.Status = settings.USER_STATUS_WAIT_ACTIVE

	//插入数据库
	uid, err := user.Insert()
	if err != nil {
		return gobalerr.DB_OPERATION_ERR
	}

	//写入用户ID
	user.UserId = uid

	//更新用户访问地址
	user.AccessUrl = fmt.Sprintf("/u/%d", user.UserId)
	rows, err := models.UpdateAccessUrlByUserId(user.UserId, user.AccessUrl)
	if err != nil || rows == 0 {
		return gobalerr.DB_OPERATION_ERR
	}

	//发送激活邮件
	AddActivateUserEmailTask(user)

	return 0
}

// CheckNameAndPwd : 检查用户名与密码是否正确
func CheckNameAndPwd(name string, pwd string) (*models.User, error) {
	//根据用户名获取用户信息
	user, err := models.GetUserByUserName(name)

	if err != nil {
		errors.UserNotExists.Data = err
		errors.UserNotExists.Source = "models.GetUserByUserName"
		return nil, errors.UserNotExists
	}

	//如果用户不存在，用邮件地址查找用户信息
	if user == nil {
		user, err = models.GetUserByEmail(name)
		if err != nil || user == nil {
			errors.UserNotExists.Data = err
			errors.UserNotExists.Source = "models.GetUserByEmail"
			return nil, errors.UserNotExists
		}
	}

	//判断密码是否正确,密码为pwd+salt的sha256
	if user.Pwd != utils.Sha256(pwd+user.PwdSalt) {
		errors.UserPwdErr.Data = err
		errors.UserNotExists.Source = "models.CheckNameAndPwd"
		return nil, errors.UserPwdErr
	}

	return user, nil
}

// ActiviteUser : 激活用户
func ActiviteUser(uid string, sign string) bool {
	result := false                              //返回结果
	userID, err := strconv.ParseInt(uid, 10, 64) //字符串转换
	if err != nil {
		return false
	}
	//获取验证码
	if result, validID := models.CheckValidCode(userID, sign, settings.VALID_CODE_ACTIVE_USER_EMAIL); result == true {
		//更新用户状态
		if rows, err := models.UpdateUserStatus(userID, settings.USER_STATUS_OK); err != nil || rows != 1 {
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
