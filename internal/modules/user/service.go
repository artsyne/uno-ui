package user

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/pkg"
)

// Service 用户业务逻辑层
type Service struct {
	repo *Repository
}

// NewService 创建用户业务逻辑层实例
func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// RegisterUser 注册新用户
func (s *Service) RegisterUser(user *User) error {
	// 转换用户名和邮件地址为小写并去除空格
	s.normalizeUserInput(user)

	// 校验用户输入
	if err := s.validateUserInput(user); err != nil {
		return err
	}

	// 初始化用户信息
	s.initUserRegInfo(user)

	// 插入数据库
	if err := s.insertUserIntoDB(user); err != nil {
		return errors.DbOperationErr
	}

	// 发送激活邮件
	s.addActivateUserEmailTask(user)

	return nil
}

// validateUserInput 验证用户输入参数
func (s *Service) validateUserInput(user *User) error {
	// 判断参数是否到位
	if len(user.UserName) < 6 || len(user.Pwd) < 6 {
		return errors.InputParamErr
	}

	// 判断用户时区与语言是否正确
	if !pkg.IsTimeZone(user.TimeZone) {
		return errors.InputParamErr
	}

	// 判断用户名是否存在
	if err := s.isExistUserName(user.UserName); err != nil {
		return err
	}

	// 判断邮件地址是否存在
	if err := s.isExistEmail(user.Email); err != nil {
		return err
	}

	return nil
}

// isExistEmail 判断邮件地址是否存在
func (s *Service) isExistEmail(email string) error {
	// 判断邮件地址是否存在
	if isExistEmail, err := s.repo.IsExistsEmail(email); err == nil && isExistEmail {
		return errors.RegisterUserEmailExists
	} else if err != nil {
		return errors.DbOperationErr
	}
	return nil
}

// IsExistEmail 对外暴露的检查邮箱是否存在的方法
func (s *Service) IsExistEmail(email string) error {
	return s.isExistEmail(email)
}

// isExistUserName 判断用户名是否存在
func (s *Service) isExistUserName(userName string) error {
	// 判断用户名是否存在
	if isExistName, err := s.repo.IsExistUserName(userName); err == nil && isExistName {
		return errors.RegisterUserNameExists
	} else if err != nil {
		return errors.DbOperationErr
	}
	return nil
}

// IsExistUserName 对外暴露的检查用户名是否存在的方法
func (s *Service) IsExistUserName(userName string) error {
	return s.isExistUserName(userName)
}

// normalizeUserInput 转换用户名和邮件地址为小写并去除空格
func (s *Service) normalizeUserInput(user *User) {
	user.UserName = strings.TrimSpace(strings.ToLower(user.UserName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
}

// initUserRegInfo 初始化用户注册信息
func (s *Service) initUserRegInfo(user *User) {
	// 随机生成salt，保存密码为pwd+salt的sha256
	user.PwdSalt = string(pkg.Rand(64, pkg.RAND_KIND_ALL))
	user.Pwd = pkg.Sha256(user.Pwd + user.PwdSalt)
	// 用户默认分组为1
	user.UserGroup = configs.DEFAULT_USER_GROUP
	// 获取当前UTC时间
	nowTime := pkg.GetNowTimeUnixMilli()
	user.CreatedUTCTime = nowTime
	user.LastLoginUTCTime = nowTime
	// 设置默认语言
	user.Lang = s.getSupportedLang(user.Lang)
	user.Status = configs.USER_STATUS_WAIT_ACTIVE
}

// insertUserIntoDB 将用户插入数据库
func (s *Service) insertUserIntoDB(user *User) error {
	uid, err := s.repo.InsertUser(user)
	if err != nil {
		return err
	}
	user.UserId = uid
	return nil
}

// CheckNameAndPassword 检查用户名与密码是否正确
func (s *Service) CheckNameAndPassword(name string, pwd string) (*User, error) {
	// 根据用户名获取用户信息
	user, err := s.repo.GetUserByUserName(name)

	if err != nil {
		return nil, errors.UserNotExists
	}

	// 如果用户不存在，用邮件地址查找用户信息
	if user == nil {
		user, err = s.repo.GetUserByEmail(name)
		if err != nil || user == nil {
			return nil, errors.UserNotExists
		}
	}

	// 判断密码是否正确,密码为pwd+salt的sha256
	if user.Pwd != pkg.Sha256(pwd+user.PwdSalt) {
		return nil, errors.UserPwdErr
	}

	return user, nil
}

// ActivateUser 激活用户
func (s *Service) ActivateUser(uid string, sign string) bool {
	// 将字符串转换为整数类型的用户ID
	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return false
	}

	// 检查验证码是否有效
	result, validID := s.repo.CheckValidCode(userID, sign, configs.VALID_CODE_ACTIVE_USER_EMAIL)
	if !result {
		return false
	}

	// 更新用户状态
	if rows, err := s.repo.UpdateUserStatus(userID, configs.USER_STATUS_OK); err != nil || rows != 1 {
		return false
	}

	// 删除验证码
	if rows, err := s.repo.DeleteValidCode(validID); err != nil || rows != 1 {
		return false
	}

	return true
}

// addActivateUserEmailTask 增加发送激活用户邮件任务
func (s *Service) addActivateUserEmailTask(user *User) error {
	// 插入激活码
	validCode, err := s.insertValidCode(user.UserId)
	if err != nil {
		return err
	}

	// 插入发送邮件队列
	return s.insertEmailTask(user, validCode)
}

// insertValidCode 插入有效验证码
func (s *Service) insertValidCode(userID int64) (*ValidCode, error) {
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(configs.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   configs.VALID_CODE_ACTIVE_USER_EMAIL,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().AddDate(0, 0, configs.REG_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli(),
	}

	if _, err := s.repo.InsertValidCode(validCode); err != nil {
		return nil, err
	}

	return validCode, nil
}

// insertEmailTask 插入邮件任务
func (s *Service) insertEmailTask(user *User, validCode *ValidCode) error {
	// 生成激活链接
	activateUrl := fmt.Sprintf("http://%s/passport/activateuser?uid=%d&sign=%s",
		configs.GetConfig().App.Domain, user.UserId, validCode.ValidCode)

	// 创建邮件任务
	emailTask := &EmailTask{
		To:      user.Email,
		Subject: "账户激活",
		Body:    fmt.Sprintf("请点击以下链接激活账户：%s", activateUrl),
	}

	// 序列化邮件任务
	body, err := json.Marshal(emailTask)
	if err != nil {
		return err
	}

	// 创建作业队列
	jobQueue := &JobQueue{
		JobType:     "email",
		Body:        string(body),
		AddUTCTime:  pkg.GetNowTime().UnixMilli(),
		PlanUTCTime: pkg.GetNowTime().UnixMilli(),
		Result:      0,
		TryNum:      0,
	}

	// 插入作业队列
	_, err = s.repo.InsertJobQueue(jobQueue)
	return err
}

// getSupportedLang 获取支持的语言
func (s *Service) getSupportedLang(language string) string {
	supportedLangs := []string{"zh-cn", "zh-tw", "en"}
	for _, lang := range supportedLangs {
		if lang == language {
			return language
		}
	}
	return "zh-cn" // 默认中文
}

// CreateResetPasswordCode 创建重置密码验证码
func (s *Service) CreateResetPasswordCode(userID int64) error {
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(configs.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   configs.VALID_CODE_TYPE_RESET_PWD,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().AddDate(0, 0, 1).UnixMilli(), // 1天有效
	}

	_, err := s.repo.InsertValidCode(validCode)
	return err
}

// SendResetPasswordEmail 发送重置密码邮件
func (s *Service) SendResetPasswordEmail(user *User) error {
	// 创建重置密码验证码
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(configs.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        user.UserId,
		ValidTypeId:   configs.VALID_CODE_TYPE_RESET_PWD,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().AddDate(0, 0, 1).UnixMilli(), // 1天有效
	}

	// 插入验证码
	_, err := s.repo.InsertValidCode(validCode)
	if err != nil {
		return err
	}

	// 生成重置链接
	resetUrl := fmt.Sprintf("http://%s/passport/resetpwd?uid=%d&sign=%s",
		configs.GetConfig().App.Domain, user.UserId, validCode.ValidCode)

	// 创建邮件任务
	emailTask := &EmailTask{
		To:      user.Email,
		Subject: "重置密码",
		Body:    fmt.Sprintf("请点击以下链接重置密码：%s", resetUrl),
	}

	// 序列化邮件任务
	body, err := json.Marshal(emailTask)
	if err != nil {
		return err
	}

	// 创建作业队列
	jobQueue := &JobQueue{
		JobType:     "email",
		Body:        string(body),
		AddUTCTime:  pkg.GetNowTime().UnixMilli(),
		PlanUTCTime: pkg.GetNowTime().UnixMilli(),
		Result:      0,
		TryNum:      0,
	}

	// 插入作业队列
	_, err = s.repo.InsertJobQueue(jobQueue)
	return err
}

// CheckValidCode 检查验证码
func (s *Service) CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	return s.repo.CheckValidCode(userID, sign, validType)
}

// UpdateUserStatus 更新用户状态
func (s *Service) UpdateUserStatus(userId int64, status int8) (int64, error) {
	return s.repo.UpdateUserStatus(userId, status)
}

// UpdateUserPwd 更新用户密码
func (s *Service) UpdateUserPwd(userId int64, pwd string) (int64, error) {
	return s.repo.UpdateUserPwd(userId, pwd)
}

// DeleteValidCode 删除验证码
func (s *Service) DeleteValidCode(id int64) (int64, error) {
	return s.repo.DeleteValidCode(id)
}

// GetUserByEmail 根据邮箱获取用户信息
func (s *Service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetUserByEmail(email)
}

// GetUserByUserId 根据用户ID获取用户信息
func (s *Service) GetUserByUserId(userId int64) (*User, error) {
	return s.repo.GetUserByUserId(userId)
}

// ParseUserID 解析用户ID
func (s *Service) ParseUserID(idStr string) (int64, error) {
	// 尝试直接解析数字
	if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
		return id, nil
	}

	return 0, fmt.Errorf("无效的用户ID格式")
}
