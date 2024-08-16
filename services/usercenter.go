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
)

// reg start
// RegUser 注册新用户
// user: 用户信息
// 返回值：用户ID和错误
// 函数会检查输入参数，验证用户名和邮件地址的唯一性，初始化用户信息，插入数据库，并发送激活邮件
func RegUser(user *models.User) error {

	//转换用户名，邮件地址为小写并去除空格
	normalizeUserInput(user)

	//校验用户输入
	if err := ValidateUserInput(user); err != nil {
		return err
	}

	//初始化用户信息
	initUserRegInfo(user)

	//插入数据库
	if err := insertUserIntoDB(user); err != nil {
		return errors.DbOperationErr
	}

	//更新用户访问地址
	user.AccessUrl = fmt.Sprintf("/u/%d", user.UserId)
	if err := UpdateUserAccessUrl(user, user.AccessUrl); err != nil {
		return err
	}

	//发送激活邮件
	AddActivateUserEmailTask(user)

	return nil
}

// ValidateUserInput 验证用户输入参数
func ValidateUserInput(user *models.User) error {
	//判断参数是否到位
	if len(user.UserName) < 6 || len(user.Pwd) < 6 {
		return errors.InputParamErr
	}

	//判断用户时区与语言是否正确，TODO
	if !utils.IsTimeZone(user.TimeZone) {
		return errors.InputParamErr
	}

	//判断用户名是否存在
	if err := IsExistUserName(user.UserName); err != nil {
		return err
	}

	//判断邮件地址是否存在
	if err := IsExistEmail(user.Email); err != nil {
		return err
	}

	return nil
}

// IsExistEmail 判断邮件地址是否存在
func IsExistEmail(email string) error {
	//判断邮件地址是否存在
	if isExistEmail, err := models.IsExistsEmail(email); err == nil && isExistEmail {
		return errors.RegisterUserEmailExists
	} else if err != nil {
		return errors.DbOperationErr
	}
	return nil
}

// IsExistUserName 判断用户名是否存在
func IsExistUserName(userName string) error {
	//判断用户名是否存在
	if isExistName, err := models.IsExistUserName(userName); err == nil && isExistName {
		return errors.RegisterUserNameExists
	} else if err != nil {
		// todo:db error 需要写log
		return errors.DbOperationErr
	}
	return nil
}

// normalizeUserInput 转换用户名和邮件地址为小写并去除空格
func normalizeUserInput(user *models.User) {
	user.UserName = strings.TrimSpace(strings.ToLower(user.UserName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
}

// initUserRegInfo 初始化用户信息
func initUserRegInfo(user *models.User) {
	//随机生成salt，保存密码为pwd+salt的sha256
	user.PwdSalt = string(utils.Rand(64, utils.RAND_KIND_ALL))
	user.Pwd = utils.Sha256(user.Pwd + user.PwdSalt)
	//用户默认分组为1
	user.UserGroup = settings.DEFAULT_USER_GROUP
	//获取当前UTC时间
	nowTime := utils.GetNowTimeUnixMilli()
	user.RegUTCTime = nowTime
	user.LastLoginUTCTime = nowTime
	//设置默认语言
	user.Lang = GetSupportedLang(user.Lang)
	user.Status = settings.USER_STATUS_WAIT_ACTIVE
}

// insertUserIntoDB 将用户插入数据库
func insertUserIntoDB(user *models.User) error {
	uid, err := user.Insert()
	if err != nil {
		return err
	}
	user.UserId = uid
	return nil
}

// UpdateUserAccessUrl 更新用户访问地址
func UpdateUserAccessUrl(user *models.User, url string) error {
	rows, err := models.UpdateAccessUrlByUserId(user.UserId, url)
	if err != nil || rows == 0 {
		return err
	}
	return nil
}

//reg end

// CheckNameAndPwd : 检查用户名与密码是否正确
func CheckNameAndPwd(name string, pwd string) (*models.User, error) {
	//根据用户名获取用户信息
	user, err := models.GetUserByUserName(name)

	if err != nil {
		return nil, errors.UserNotExists
	}

	//如果用户不存在，用邮件地址查找用户信息
	if user == nil {
		user, err = models.GetUserByEmail(name)
		if err != nil || user == nil {
			return nil, errors.UserNotExists
		}
	}

	//判断密码是否正确,密码为pwd+salt的sha256
	if user.Pwd != utils.Sha256(pwd+user.PwdSalt) {
		return nil, errors.UserPwdErr
	}

	return user, nil
}

// ActivateUser : 激活用户
func ActivateUser(uid string, sign string) bool {
	// 将字符串转换为整数类型的用户ID
	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return false
	}

	// 检查验证码是否有效
	result, validID := models.CheckValidCode(userID, sign, settings.VALID_CODE_ACTIVE_USER_EMAIL)
	if !result {
		return false
	}

	// 更新用户状态
	if rows, err := models.UpdateUserStatus(userID, settings.USER_STATUS_OK); err != nil || rows != 1 {
		return false
	}

	// 删除验证码
	if rows, err := models.DeleteValidCode(validID); err != nil || rows != 1 {
		return false
	}

	return true
}

// region
// AddActivateUserEmailTask : 增加发送激活用户邮件任务
func AddActivateUserEmailTask(user *models.User) error {
    // 插入激活码
    validCode, err := insertValidCode(user.UserId)
    if err != nil {
        return err
    }

    // 插入发送邮件队列
    return insertEmailTask(user, validCode)
}

// insertValidCode 插入有效验证码
func insertValidCode(userID int64) (*models.ValidCode, error) {
    validCode := &models.ValidCode{
        ValidCode:     string(utils.Rand(settings.VALID_CODE_LEN, utils.RAND_KIND_ALL)),
        UserId:        userID,
        ValidTypeId:   settings.VALID_CODE_ACTIVE_USER_EMAIL,
        CreateUTCTime: utils.GetNowTime().UnixMilli(),
        ExpireUTCTime: utils.GetNowTime().AddDate(0, 0, settings.REG_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli(),
    }

    if _, err := validCode.Insert(); err != nil {
        return nil, err
    }

    return validCode, nil
}

// insertEmailTask 插入邮件任务
func insertEmailTask(user *models.User, validCode *models.ValidCode) error {
    config := settings.GetConfig()
    params := settings.GetParams()
    nowTime := utils.GetNowTime()

    // 激活链接
    url := fmt.Sprintf("%s/activateuser?uid=%d&sign=%s", config.Domain, user.UserId, validCode.ValidCode)

    // 邮件task对象
    emailTask := &models.EmailTask{
        To:      user.Email,
        Subject: params["reg_email_subject"],
        Body:    fmt.Sprintf(params["reg_email_body"], user.UserName, url, url),
    }

    // 插入邮件队列
    jobQueue := &models.JobQueue{
        AddUTCTime: nowTime.UnixMilli(),
        JobType:    settings.JOB_TYPE_EMAIL,
        PlanUTCTime: nowTime.UnixMilli(),
        TryNum:     0,
        Result:     1,
    }

    jsonBytes, err := json.Marshal(emailTask)
    if err != nil {
        return err
    }
    jobQueue.Body = string(jsonBytes)

    _, err = jobQueue.Insert()
    return err
}
// endregion

// 根据用户的语言偏好和系统支持的语言列表，返回匹配的系统支持语言。
// 该函数首先将用户的语言偏好转换为小写，然后进行以下判断：
// 1. 如果用户偏好的语言是简体中文，则返回 "zh-cn"；
// 2. 如果用户偏好的语言是香港或台湾的本地语言，则返回 "zh-tw"；
// 3. 如果用户偏好的语言在系统支持的语言列表中，则返回该语言；
// 4. 如果用户偏好的语言不在系统支持的语言列表中，则返回默认语言 "en"。
// 参数：
// language: 用户的语言偏好，通常来自浏览器的 Accept-Language 头部；
// supportedLanguages: 系统支持的语言列表，来自配置文件。
func GetSupportedLang(language string) string {
	//转换为小写
	language = strings.ToLower(language)

	//获取系统支持的语言列表
	supportedLangs := settings.GetConfig().SupportLangs

	switch {
	case strings.Contains(language, "zh-cn"), strings.Contains(language, "zh-hans"):
		return "zh-cn"
	case strings.Contains(language, "zh-hk"), strings.Contains(language, "zh-tw"), strings.Contains(language, "zh-hant"):
		return "zh-tw"
	default:
		language = language[0:2]
		for _, supportedLanguage := range supportedLangs {
			if supportedLanguage.Code == language {
				return language
			}
		}
		return "en"
	}
}
