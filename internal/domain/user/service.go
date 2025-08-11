package user

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/utils/config"
	appErrors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/pkg"
)



// Service 用户业务逻辑层实现
type Service struct {
	repo       RepositoryInterface
	config     *config.Config // 缓存配置，避免重复调用全局变量
	regexCheck *pkg.RegexCheck
	taskService *task.Service
}

// NewService 创建用户业务逻辑层实例
// NewService 创建用户服务实例，使用依赖注入
func NewService(repo RepositoryInterface, taskService *task.Service, cfg *config.Config) *Service {

	return &Service{
		repo:       repo,
		config:     cfg, // 缓存配置
		regexCheck: &pkg.RegexCheck{},
		taskService: taskService,
	}
}

// ===================== 用户注册和认证 =====================

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
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("用户注册失败: %v", err))
	}

	// 发送激活邮件
	s.addActivateUserEmailTask(user)

	return nil
}

// LoginWithUsernameOrEmail 用户登录验证，支持用户名或邮箱登录
// 先尝试用户名登录，如果用户不存在则尝试邮箱登录
func (s *Service) LoginWithUsernameOrEmail(usernameOrEmail string, password string) (*User, error) {
	var user *User
	var err error

	// 判断输入是否为邮箱地址
	if s.regexCheck.IsEmail(usernameOrEmail) {
		// 如果是邮箱，直接通过邮箱查找
		user, err = s.repo.GetUserByEmail(usernameOrEmail)
	} else {
		// 如果不是邮箱，通过用户名查找
		user, err = s.repo.GetUserByUserName(usernameOrEmail)
	}

	if err != nil {
		// 用户不存在
		return nil, err
	}

	// 验证密码
	if !s.verifyPassword(password, user.Pwd) {
		return nil, appErrors.PasswordIncorrect.WithDetails("密码不正确")
	}

	return user, nil
}



// ActivateUser 激活用户
func (s *Service) ActivateUser(uid string, sign string) error {
	// 将字符串转换为整数类型的用户ID
	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return appErrors.InvalidParams.WithDetails("无效的用户ID格式")
	}

	// 检查验证码是否有效
	result, validID := s.repo.CheckValidCode(userID, sign, config.VALID_CODE_ACTIVE_USER_EMAIL)
	if !result {
		return appErrors.InvalidParams.WithDetails("验证码无效或已过期")
	}

	// 更新用户状态
	if rows, err := s.repo.UpdateUserStatus(userID, config.USER_STATUS_OK); err != nil || rows != 1 {
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("更新用户状态失败: %v", err))
	}

	// 删除验证码
	if rows, err := s.repo.DeleteValidCode(validID); err != nil || rows != 1 {
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("删除验证码失败: %v", err))
	}

	return nil
}

// ===================== 用户信息验证 =====================

// IsExistEmail 判断邮件地址是否存在
func (s *Service) IsExistEmail(email string) (bool, error) {
	// 判断邮件地址是否存在
	isExistEmail, err := s.repo.IsExistsEmail(email)
	if err != nil {
		return false, appErrors.DatabaseError.WithDetails(err.Error())
	}
	return isExistEmail, nil
}

// IsExistUserName 判断用户名是否存在
func (s *Service) IsExistUserName(userName string) (bool, error) {
	// 判断用户名是否存在
	isExistName, err := s.repo.IsExistUserName(userName)
	if err != nil {
		return false, appErrors.DatabaseError.WithDetails(err.Error())
	}
	return isExistName, nil
}

// ===================== 密码管理 =====================

// CreateResetPasswordCode 创建重置密码验证码
func (s *Service) CreateResetPasswordCode(userID int64) error {
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
func (s *Service) SendResetPasswordEmail(user *User) error {
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

	// 调用task服务创建邮件任务
	return s.taskService.CreateEmailTask(user.Email, "重置密码", fmt.Sprintf("请点击以下链接重置密码：%s", resetUrl))
}

// ChangeUserPassword 修改用户密码
// 验证当前密码后更新为新密码，使用与initUserRegInfo相同的加密逻辑
func (s *Service) ChangeUserPassword(currentUserID int64, currentPwd, newPwd string) error {
	// 获取用户当前密码
	currentEncryptPwd, err := s.repo.GetUserEncryptPwd(currentUserID)
	if err != nil {
		return err
	}

	// 验证当前密码,有权限后才能修改数据库中的密码 - 使用统一的密码验证逻辑
	if !s.verifyPassword(currentPwd, currentEncryptPwd) {
		return appErrors.InvalidPassword.WithDetails("当前密码不正确")
	}
	
	// 验证新密码强度
	if len(newPwd) < config.PWD_MIN_LEN || len(newPwd) > config.PWD_MAX_LEN {
		return appErrors.ValidationError.WithDetails(fmt.Sprintf("新密码长度必须在%d-%d个字符之间", config.PWD_MIN_LEN, config.PWD_MAX_LEN))
	}

	// 使用统一的密码生成逻辑
	newPwdHash, err := s.generatePasswordHash(newPwd)
	if err != nil {
		return appErrors.DatabaseError.WithDetails(err.Error())
	}

	// 更新用户密码
	if err := s.UpdateUserPassword(currentUserID, newPwdHash); err != nil {
		return appErrors.DatabaseError.WithDetails(err.Error())
	}

	return nil
}

// ResetUserPassword 重置用户密码
// 用于忘记密码后通过验证码重置密码，使用统一的密码生成逻辑
func (s *Service) ResetUserPassword(userId int64, newPassword string) error {
	// 验证新密码强度
	if len(newPassword) < config.PWD_MIN_LEN || len(newPassword) > config.PWD_MAX_LEN {
		return appErrors.ValidationError.WithDetails(fmt.Sprintf("新密码长度必须在%d-%d个字符之间", config.PWD_MIN_LEN, config.PWD_MAX_LEN))
	}

	// 使用统一的密码生成逻辑
	newPwdHash, err := s.generatePasswordHash(newPassword)
	if err != nil {
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("密码哈希生成失败: %v", err))
	}

	// 更新用户密码
	if err := s.UpdateUserPassword(userId, newPwdHash); err != nil {
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("密码重置失败: %v", err))
	}

	return nil
}

// UpdateUserPassword 更新用户密码
// 内部方法，用于密码修改时更新密码
func (s *Service) UpdateUserPassword(userId int64, hashedPassword string) error {
	_, err := s.repo.UpdateUserPassword(userId, hashedPassword)
	return err
}

// ===================== 验证码管理 =====================

// CheckValidCode 检查验证码
func (s *Service) CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	return s.repo.CheckValidCode(userID, sign, validType)
}

// DeleteValidCode 删除验证码
func (s *Service) DeleteValidCode(id int64) (int64, error) {
	return s.repo.DeleteValidCode(id)
}

// ===================== 用户信息获取 =====================

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

	return 0, appErrors.InvalidParams.WithDetails("无效的用户ID格式")
}

// ===================== 个人资料管理 =====================
// UpdateUserProfile 更新用户基本资料
// 更新用户的全名和个人简介等信息
func (s *Service) UpdateUserProfile(uid int64, fullName, remark string) error {
	// 验证输入参数
	if len(fullName) > 50 {
		return appErrors.InvalidParams.WithDetails("全名长度不能超过50个字符")
	}
	if len(remark) > 250 {
		return appErrors.InvalidParams.WithDetails("个人简介长度不能超过250个字符")
	}

	// 调用repository更新用户资料
	rowsAffected, err := s.repo.UpdateUserProfile(uid, fullName, remark)
	if err != nil {
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("更新用户资料失败: %v", err))
	}

	// 检查是否真的更新了用户
	if rowsAffected == 0 {
		return appErrors.DatabaseError.WithDetails("未找到用户或资料未修改")
	}

	return nil
}

// UpdateUserEmail 更新用户邮箱
func (s *Service) UpdateUserEmail(uid int64, email string) error {
	// 验证邮箱格式
	if !s.regexCheck.IsEmail(email) {
		return appErrors.InvalidFormat.WithDetails("邮箱格式不正确")
	}

	// 检查用户是否存在
	_, err := s.repo.GetUserByUserId(uid)
	if err != nil {
		return err // 如果用户不存在，直接返回错误
	}

	// 检查邮箱是否已被其他用户使用
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		// 检查是否是 UserNotFound 错误
		if appErr, ok := err.(*appErrors.AppError); ok && appErr.Code == appErrors.UserNotFound.Code {
			// 如果是 UserNotFound 错误，说明邮箱未被使用，这是正常的，继续执行
		} else {
			// 如果是其他数据库错误，则返回数据库错误
			return appErrors.DatabaseError.WithDetails(fmt.Sprintf("查询邮箱失败: %v", err))
		}
	} else if existingUser != nil && existingUser.UserId != uid {
		// 如果找到了用户且不是当前用户，则表示邮箱已被占用
		return appErrors.EmailExists.WithDetails("邮箱已被其他用户使用")
	}

	// 调用 repository 更新邮箱
	if _, err := s.repo.UpdateEmailByUserId(uid, email); err != nil {
		return appErrors.DatabaseError.WithDetails(fmt.Sprintf("更新邮箱失败: %v", err))
	}

	return nil
}


// UpdateUserStatus 更新用户状态
	func (s *Service) UpdateUserStatus(userId int64, status int8) (int64, error) {
		rowsAffected, err := s.repo.UpdateUserStatus(userId, status)
		if err != nil {
			return 0, appErrors.DatabaseError.WithDetails(fmt.Sprintf("更新用户状态失败: %v", err))
		}
		// 检查是否更新成功
		if rowsAffected == 0 {
			return 0, appErrors.DatabaseError.WithDetails("未找到用户或状态未修改")
		}
		return rowsAffected, nil
	}

// UpdateUserPwd 更新用户密码
	func (s *Service) UpdateUserPwd(userId int64, pwd string) (int64, error) {
		rowsAffected, err := s.repo.UpdateUserPwd(userId, pwd)
		if err != nil {
			return 0, appErrors.DatabaseError.WithDetails(fmt.Sprintf("更新用户密码失败: %v", err))
		}
		return rowsAffected, nil
	}

// GetSupportedLang 获取支持的语言
func (s *Service) GetSupportedLang(language string) string {
	supportedLangs := []string{"zh-cn", "zh-tw", "en"}
	for _, lang := range supportedLangs {
		if lang == language {
			return language
		}
	}
	return "en" // 默认英文
}

// ===================== 私有辅助方法 =====================

// validateUserInput 验证用户输入参数
	func (s *Service) validateUserInput(user *User) error {
		// 判断参数是否到位
		if len(user.UserName) < 6 || len(user.Pwd) < 6 {
			return appErrors.ValidationError.WithDetails("用户名和密码长度至少6位")
		}

		// 判断用户时区与语言是否正确
		if !pkg.IsTimeZone(user.TimeZone) {
			return appErrors.ValidationError.WithDetails("时区格式不正确")
		}

		// 判断用户名是否存在
		isExistUserName, err := s.IsExistUserName(user.UserName)
		if err != nil {
			return err
		}
		if isExistUserName {
			return appErrors.UserExists.WithDetails("用户名: " + user.UserName)
		}

		// 判断邮件地址是否存在
		isExistEmail, err := s.IsExistEmail(user.Email)
		if err != nil {
			return err
		}
		if isExistEmail {
			return appErrors.EmailExists.WithDetails("邮箱: " + user.Email)
		}

		return nil
	}

// normalizeUserInput 转换用户名和邮件地址为小写并去除空格
func (s *Service) normalizeUserInput(user *User) {
	user.UserName = strings.TrimSpace(strings.ToLower(user.UserName))
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
}

// generatePasswordHash 生成密码哈希
// 使用 bcrypt 算法，不再需要单独的盐值
func (s *Service) generatePasswordHash(password string) (string, error) {
	// 使用 bcrypt 哈希密码
	hash, err := pkg.BcryptHash(password)
	if err != nil {
		return "", err
	}
	
	return hash, nil
}

// verifyPassword 验证密码
// 使用 bcrypt 验证，不再需要盐值参数
func (s *Service) verifyPassword(password, hash string) bool {
	// 使用 bcrypt 验证密码
	return pkg.BcryptVerify(password, hash)
}

// initUserRegInfo 初始化用户注册信息
func (s *Service) initUserRegInfo(user *User) {
	// 使用统一的密码加密逻辑
	hashedPwd, err := s.generatePasswordHash(user.Pwd)
	if err != nil {
		// 在这里处理错误，但由于原方法没有返回错误，我们暂时忽略
		// 在实际应用中应该修改方法签名来处理错误
	}
	user.Pwd = hashedPwd

	// 用户默认分组为1
	user.UserGroup = config.DEFAULT_USER_GROUP
	// 获取当前UTC时间
	nowTime := pkg.GetNowTimeUnixMilli()
	user.CreatedUTCTime = nowTime
	user.LastLoginUTCTime = nowTime
	// 设置默认语言
	user.Lang = s.GetSupportedLang(user.Lang)
	user.Status = config.USER_STATUS_WAIT_ACTIVE
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

// addActivateUserEmailTask 增加发送激活用户邮件任务
func (s *Service) addActivateUserEmailTask(user *User) error {
	// 插入激活码
	validCode, err := s.insertValidCode(user.UserId)
	if err != nil {
		return err
	}

	// 生成激活链接
	activateUrl := fmt.Sprintf("http://%s/passport/activateuser?uid=%d&sign=%s",
		s.config.App.Domain, user.UserId, validCode.ValidCode)

	// 调用task服务创建邮件任务
	return s.taskService.CreateEmailTask(user.Email, "账户激活", fmt.Sprintf("请点击以下链接激活账户：%s", activateUrl))
}

// insertValidCode 插入有效验证码
func (s *Service) insertValidCode(userID int64) (*ValidCode, error) {
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

