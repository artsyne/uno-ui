package user

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/core/config"
	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/biuren/hebe/internal/module/task"
	"github.com/biuren/hebe/pkg"
)

// ServiceInterface 用户业务逻辑层接口
// 定义所有与用户相关的业务逻辑操作，包括注册、登录、激活、密码管理、个人资料管理等
type ServiceInterface interface {
	// ===================== 用户注册和认证 =====================

	// RegisterUser 注册新用户
	RegisterUser(user *User) error

	// CheckNameAndPassword 检查用户名与密码是否正确
	CheckNameAndPassword(name string, pwd string) (*User, error)

	// ActivateUser 激活用户
	ActivateUser(uid string, sign string) bool

	// ===================== 用户信息验证 =====================

	// IsExistEmail 对外暴露的检查邮箱是否存在的方法
	IsExistEmail(email string) error

	// IsExistUserName 对外暴露的检查用户名是否存在的方法
	IsExistUserName(userName string) error

	// ===================== 密码管理 =====================

	// CreateResetPasswordCode 创建重置密码验证码
	CreateResetPasswordCode(userID int64) error

	// SendResetPasswordEmail 发送重置密码邮件
	SendResetPasswordEmail(user *User) error

	// ChangeUserPassword 修改用户密码
	ChangeUserPassword(uid int64, currentUserID int64, currentPwd, newPwd string) error

	// ResetUserPassword 重置用户密码
	ResetUserPassword(userId int64, newPassword string) error

	// ===================== 验证码管理 =====================

	// CheckValidCode 检查验证码
	CheckValidCode(userID int64, sign string, validType int8) (bool, int64)

	// DeleteValidCode 删除验证码
	DeleteValidCode(id int64) (int64, error)

	// ===================== 用户信息获取 =====================

	// GetUserByEmail 根据邮箱获取用户信息
	GetUserByEmail(email string) (*User, error)

	// GetUserByUserId 根据用户ID获取用户信息
	GetUserByUserId(userId int64) (*User, error)

	// ParseUserID 解析用户ID
	ParseUserID(idStr string) (int64, error)

	// ===================== 个人资料管理 =====================

	// GetUserProfile 获取用户个人档案信息
	GetUserProfile(uid int64, currentUserID int64) (*User, bool, error)

	// GetUserSettings 获取用户设置信息
	GetUserSettings(uid int64, currentUserID int64) (*User, error)

	// UpdateUserProfile 更新用户基本资料
	UpdateUserProfile(uid int64, currentUserID int64, fullName, remark string) error

	// UpdateUserEmail 更新用户邮箱
	UpdateUserEmail(uid int64, currentUserID int64, email string) error

	// ===================== 系统管理相关 =====================

	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(userId int64, status int8) (int64, error)

	// UpdateUserPwd 更新用户密码
	UpdateUserPwd(userId int64, pwd string) (int64, error)

	// GetSupportedLang 获取支持的语言（公开方法，供测试使用）
	GetSupportedLang(language string) string
}

// service 用户业务逻辑层实现
type service struct {
	repo   RepositoryInterface
	config *config.Config // 缓存配置，避免重复调用全局变量
}

// NewService 创建用户业务逻辑层实例
func NewService() ServiceInterface {
	config := config.GetConfig()
	if config == nil {
		// 在初始化时就检查配置，如果配置未正确加载则panic
		// 这是合理的，因为没有配置应用无法正常运行
		panic("Failed to get config during user service initialization - config is nil")
	}

	return &service{
		repo:   NewRepository(),
		config: config, // 缓存配置
	}
}

// RegisterUser 注册新用户
func (s *service) RegisterUser(user *User) error {
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
		return errors.DatabaseErr().WithDetails(err.Error())
	}

	// 发送激活邮件
	s.addActivateUserEmailTask(user)

	return nil
}

// validateUserInput 验证用户输入参数
func (s *service) validateUserInput(user *User) error {
	// 判断参数是否到位
	if len(user.UserName) < 6 || len(user.Pwd) < 6 {
		return errors.ValidationErr().WithDetails("用户名和密码长度至少6位")
	}

	// 判断用户时区与语言是否正确
	if !pkg.IsTimeZone(user.TimeZone) {
		return errors.ValidationErr().WithDetails("时区格式不正确")
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
func (s *service) isExistEmail(email string) error {
	// 判断邮件地址是否存在
	if isExistEmail, err := s.repo.IsExistsEmail(email); err == nil && isExistEmail {
		return errors.EmailExistsErr().WithDetails("邮箱: " + email)
	} else if err != nil {
		return errors.DatabaseErr().WithDetails(err.Error())
	}
	return nil
}

// IsExistEmail 对外暴露的检查邮箱是否存在的方法
func (s *service) IsExistEmail(email string) error {
	return s.isExistEmail(email)
}

// isExistUserName 判断用户名是否存在
func (s *service) isExistUserName(userName string) error {
	// 判断用户名是否存在
	if isExistName, err := s.repo.IsExistUserName(userName); err == nil && isExistName {
		return errors.UserExistsErr().WithDetails("用户名: " + userName)
	} else if err != nil {
		return errors.DatabaseErr().WithDetails(err.Error())
	}
	return nil
}

// IsExistUserName 对外暴露的检查用户名是否存在的方法
func (s *service) IsExistUserName(userName string) error {
	return s.isExistUserName(userName)
}

// normalizeUserInput 转换用户名和邮件地址为小写并去除空格
func (s *service) normalizeUserInput(user *User) {
	user.UserName = strings.TrimSpace(strings.ToLower(user.UserName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
}

// generatePasswordHash 生成密码哈希和盐值
// 统一的密码加密逻辑，用于用户注册和密码修改
// 返回: (加密后的密码, 盐值)
func (s *service) generatePasswordHash(plainPassword string) (string, string) {
	// 随机生成64位盐值
	salt := string(pkg.Rand(64, pkg.RAND_KIND_ALL))
	// 使用SHA256加密 密码+盐值
	hashedPassword := pkg.Sha256(plainPassword + salt)
	return hashedPassword, salt
}

// verifyPassword 验证密码是否正确
// 用于登录验证和修改密码时的当前密码验证
func (s *service) verifyPassword(plainPassword, hashedPassword, salt string) bool {
	return hashedPassword == pkg.Sha256(plainPassword+salt)
}

// updateUserPasswordAndSalt 更新用户密码和盐值
// 内部方法，用于密码修改时同时更新密码和盐值
func (s *service) updateUserPasswordAndSalt(userId int64, hashedPassword, salt string) error {
	_, err := s.repo.UpdateUserPasswordAndSalt(userId, hashedPassword, salt)
	return err
}

// initUserRegInfo 初始化用户注册信息
func (s *service) initUserRegInfo(user *User) {
	// 使用统一的密码加密逻辑
	hashedPwd, salt := s.generatePasswordHash(user.Pwd)
	user.Pwd = hashedPwd
	user.PwdSalt = salt

	// 用户默认分组为1
	user.UserGroup = config.DEFAULT_USER_GROUP
	// 获取当前UTC时间
	nowTime := pkg.GetNowTimeUnixMilli()
	user.CreatedUTCTime = nowTime
	user.LastLoginUTCTime = nowTime
	// 设置默认语言
	user.Lang = s.getSupportedLang(user.Lang)
	user.Status = config.USER_STATUS_WAIT_ACTIVE
}

// insertUserIntoDB 将用户插入数据库
func (s *service) insertUserIntoDB(user *User) error {
	uid, err := s.repo.InsertUser(user)
	if err != nil {
		return err
	}
	user.UserId = uid
	return nil
}

// CheckNameAndPassword 检查用户名与密码是否正确
func (s *service) CheckNameAndPassword(name string, pwd string) (*User, error) {
	// 根据用户名获取用户信息
	user, err := s.repo.GetUserByUserName(name)

	if err != nil {
		return nil, errors.UserNotFoundErr().WithDetails(err.Error())
	}

	// 如果用户不存在，用邮件地址查找用户信息
	if user == nil {
		user, err = s.repo.GetUserByEmail(name)
		if err != nil || user == nil {
			return nil, errors.UserNotFoundErr().WithDetails(err.Error())
		}
	}

	// 判断密码是否正确,密码为pwd+salt的sha256
	if !s.verifyPassword(pwd, user.Pwd, user.PwdSalt) {
		return nil, errors.PasswordIncorrectErr().WithDetails("密码不正确")
	}

	return user, nil
}

// ActivateUser 激活用户
func (s *service) ActivateUser(uid string, sign string) bool {
	// 将字符串转换为整数类型的用户ID
	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return false
	}

	// 检查验证码是否有效
	result, validID := s.repo.CheckValidCode(userID, sign, config.VALID_CODE_ACTIVE_USER_EMAIL)
	if !result {
		return false
	}

	// 更新用户状态
	if rows, err := s.repo.UpdateUserStatus(userID, config.USER_STATUS_OK); err != nil || rows != 1 {
		return false
	}

	// 删除验证码
	if rows, err := s.repo.DeleteValidCode(validID); err != nil || rows != 1 {
		return false
	}

	return true
}

// addActivateUserEmailTask 增加发送激活用户邮件任务
func (s *service) addActivateUserEmailTask(user *User) error {
	// 插入激活码
	validCode, err := s.insertValidCode(user.UserId)
	if err != nil {
		return err
	}

	// 生成激活链接
	activateUrl := fmt.Sprintf("http://%s/passport/activateuser?uid=%d&sign=%s",
		s.config.App.Domain, user.UserId, validCode.ValidCode)

	// 调用task模块创建邮件任务
	taskService := task.NewService()
	return taskService.CreateEmailTask(
		user.Email,
		"账户激活",
		fmt.Sprintf("请点击以下链接激活账户：%s", activateUrl),
	)
}

// insertValidCode 插入有效验证码
func (s *service) insertValidCode(userID int64) (*ValidCode, error) {
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(config.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   config.VALID_CODE_ACTIVE_USER_EMAIL,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().AddDate(0, 0, config.REG_EMAIL_VAILID_CODE_EXPIRE_DAY).UnixMilli(),
	}

	if _, err := s.repo.InsertValidCode(validCode); err != nil {
		return nil, err
	}

	return validCode, nil
}

// getSupportedLang 获取支持的语言
func (s *service) getSupportedLang(language string) string {
	supportedLangs := []string{"zh-cn", "zh-tw", "en"}
	for _, lang := range supportedLangs {
		if lang == language {
			return language
		}
	}
	return "zh-cn" // 默认中文
}

// GetSupportedLang 获取支持的语言（公开方法，供测试使用）
func (s *service) GetSupportedLang(language string) string {
	return s.getSupportedLang(language)
}

// CreateResetPasswordCode 创建重置密码验证码
func (s *service) CreateResetPasswordCode(userID int64) error {
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(config.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   config.VALID_CODE_TYPE_RESET_PWD,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().AddDate(0, 0, 1).UnixMilli(), // 1天有效
	}

	_, err := s.repo.InsertValidCode(validCode)
	return err
}

// SendResetPasswordEmail 发送重置密码邮件
func (s *service) SendResetPasswordEmail(user *User) error {
	// 创建重置密码验证码
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(config.VALID_CODE_LEN, pkg.RAND_KIND_ALL)),
		UserId:        user.UserId,
		ValidTypeId:   config.VALID_CODE_TYPE_RESET_PWD,
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
		s.config.App.Domain, user.UserId, validCode.ValidCode)

	// 调用task模块创建邮件任务
	taskService := task.NewService()
	return taskService.CreateEmailTask(
		user.Email,
		"重置密码",
		fmt.Sprintf("请点击以下链接重置密码：%s", resetUrl),
	)
}

// CheckValidCode 检查验证码
func (s *service) CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	return s.repo.CheckValidCode(userID, sign, validType)
}

// UpdateUserStatus 更新用户状态
func (s *service) UpdateUserStatus(userId int64, status int8) (int64, error) {
	return s.repo.UpdateUserStatus(userId, status)
}

// UpdateUserPwd 更新用户密码
func (s *service) UpdateUserPwd(userId int64, pwd string) (int64, error) {
	return s.repo.UpdateUserPwd(userId, pwd)
}

// DeleteValidCode 删除验证码
func (s *service) DeleteValidCode(id int64) (int64, error) {
	return s.repo.DeleteValidCode(id)
}

// GetUserByEmail 根据邮箱获取用户信息
func (s *service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetUserByEmail(email)
}

// GetUserByUserId 根据用户ID获取用户信息
func (s *service) GetUserByUserId(userId int64) (*User, error) {
	return s.repo.GetUserByUserId(userId)
}

// ParseUserID 解析用户ID
func (s *service) ParseUserID(idStr string) (int64, error) {
	// 尝试直接解析数字
	if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
		return id, nil
	}

	return 0, fmt.Errorf("无效的用户ID格式")
}

// ===================== Member 相关业务逻辑 =====================

// GetUserProfile 获取用户个人档案信息
// 用于渲染用户个人主页，包含基本信息验证
func (s *service) GetUserProfile(uid int64, currentUserID int64) (*User, bool, error) {
	// 获取目标用户信息
	user, err := s.repo.GetUserByUserId(uid)
	if err != nil {
		return nil, false, err
	}

	// 判断是否为用户本人
	isOwner := (currentUserID == uid)

	return user, isOwner, nil
}

// GetUserSettings 获取用户设置信息
// 仅允许用户本人访问，返回完整的用户信息用于设置页面
func (s *service) GetUserSettings(uid int64, currentUserID int64) (*User, error) {
	// 验证权限：只能访问自己的设置页面
	if currentUserID != uid {
		return nil, fmt.Errorf("无权访问他人的设置页面")
	}

	// 获取用户信息
	user, err := s.repo.GetUserByUserId(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserProfile 更新用户基本资料
// 更新用户的全名和个人简介等信息
func (s *service) UpdateUserProfile(uid int64, currentUserID int64, fullName, remark string) error {
	// 验证权限：只能修改自己的资料
	if currentUserID != uid {
		return fmt.Errorf("无权修改他人的资料")
	}

	// 验证输入参数
	if len(fullName) > 50 {
		return fmt.Errorf("全名长度不能超过50个字符")
	}
	if len(remark) > 250 {
		return fmt.Errorf("个人简介长度不能超过250个字符")
	}

	// 调用repository更新用户资料
	_, err := s.repo.UpdateUserProfile(uid, fullName, remark)
	return err
}

// UpdateUserEmail 更新用户邮箱
// 更新用户的邮箱地址，需要验证邮箱格式
func (s *service) UpdateUserEmail(uid int64, currentUserID int64, email string) error {
	// 验证权限：只能修改自己的邮箱
	if currentUserID != uid {
		return fmt.Errorf("无权修改他人的邮箱")
	}

	// 验证邮箱格式 - 使用与注册时相同的验证逻辑
	regexCheck := &pkg.RegexCheck{}
	if !regexCheck.IsEmail(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	// 检查邮箱是否已被其他用户使用
	existingUser, err := s.repo.GetUserByEmail(email)
	if err == nil && existingUser != nil && existingUser.UserId != uid {
		return fmt.Errorf("该邮箱已被其他用户使用")
	}

	// 调用repository更新邮箱
	_, err = s.repo.UpdateEmailByUserId(uid, email)
	return err
}

// ChangeUserPassword 修改用户密码
// 验证当前密码后更新为新密码，使用与initUserRegInfo相同的加密逻辑
func (s *service) ChangeUserPassword(uid int64, currentUserID int64, currentPwd, newPwd string) error {
	// 验证权限：只能修改自己的密码
	if currentUserID != uid {
		return fmt.Errorf("无权修改他人的密码")
	}

	// 获取用户当前信息
	user, err := s.repo.GetUserByUserId(uid)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证当前密码 - 使用统一的密码验证逻辑
	if !s.verifyPassword(currentPwd, user.Pwd, user.PwdSalt) {
		return fmt.Errorf("当前密码不正确")
	}

	// 验证新密码强度
	if len(newPwd) < config.PWD_MIN_LEN || len(newPwd) > config.PWD_MAX_LEN {
		return fmt.Errorf("新密码长度必须在%d-%d个字符之间", config.PWD_MIN_LEN, config.PWD_MAX_LEN)
	}

	// 使用统一的密码生成逻辑
	newPwdHash, newSalt := s.generatePasswordHash(newPwd)

	// 更新用户密码和盐值 - 需要先更新repository方法支持盐值
	if err := s.updateUserPasswordAndSalt(uid, newPwdHash, newSalt); err != nil {
		return fmt.Errorf("密码更新失败: %v", err)
	}

	return nil
}

// ResetUserPassword 重置用户密码
// 用于忘记密码后通过验证码重置密码，使用统一的密码生成逻辑
func (s *service) ResetUserPassword(userId int64, newPassword string) error {
	// 验证新密码强度
	if len(newPassword) < config.PWD_MIN_LEN || len(newPassword) > config.PWD_MAX_LEN {
		return fmt.Errorf("新密码长度必须在%d-%d个字符之间", config.PWD_MIN_LEN, config.PWD_MAX_LEN)
	}

	// 使用统一的密码生成逻辑
	newPwdHash, newSalt := s.generatePasswordHash(newPassword)

	// 更新用户密码和盐值
	if err := s.updateUserPasswordAndSalt(userId, newPwdHash, newSalt); err != nil {
		return fmt.Errorf("密码重置失败: %v", err)
	}

	return nil
}
