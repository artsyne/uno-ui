package user

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/pkg"
	"github.com/go-redis/redis"
)

// Service 用户业务逻辑层实现
type Service struct {
	config      *config.Config
	repo        RepositoryInterface
	taskService *task.Service
	i18nService i18n.I18nManager
	redis       *redis.Client
}

// NewService 创建用户业务逻辑层实例
// NewService 创建用户服务实例，使用依赖注入
func NewService(config *config.Config, repo RepositoryInterface, taskService *task.Service, i18nService i18n.I18nManager, redis *redis.Client) *Service {

	return &Service{
		config:      config,
		repo:        repo,
		taskService: taskService,
		i18nService: i18nService,
		redis:       redis,
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
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "registration_failed"})
	}

	// 发送激活邮件
	if err := s.addActivateUserEmailTask(user); err != nil {
		slog.Error("发送激活邮件任务失败", "error", err, "user_id", user.UserId)
	}

	return nil
}

// LoginWithUsernameOrEmail 用户登录验证，支持用户名或邮箱登录
// 先尝试用户名登录，如果用户不存在则尝试邮箱登录
func (s *Service) LoginWithUsernameOrEmail(loginName string, password string) (*User, error) {
	var user *User
	var err error

	// 判断输入是否为邮箱地址
	if pkg.IsEmail(loginName) {
		// 如果是邮箱，直接通过邮箱查找
		user, err = s.repo.GetUserByEmail(loginName)
	} else {
		// 如果不是邮箱，通过Handle查找
		user, err = s.repo.GetUserByHandle(loginName)
	}

	if err != nil {
		// 用户不存在
		return nil, err
	}

	// 检查账户是否被锁定
	if err := s.checkBruteForce(user.UserId); err != nil {
		return nil, err
	}

	// 验证密码
	if !s.verifyPassword(password, user.Pwd) {
		// 记录登录失败
		s.recordLoginFailure(user.UserId)
		return nil, errors.PasswordIncorrect
	}

	// 登录成功，清除失败记录
	s.clearLoginFailure(user.UserId)

	return user, nil
}

// checkBruteForce 检查暴力破解限制
func (s *Service) checkBruteForce(userID int64) error {
	if s.redis == nil {
		return nil
	}

	key := fmt.Sprintf("brute_force:uid:%d", userID)
	count, err := s.redis.Get(key).Int()
	if err != nil && err != redis.Nil {
		// Redis错误，降级处理，允许通过
		return nil
	}

	if count >= s.config.BruteForce.MaxFailures {
		return errors.AccountLocked.WithDetails(map[string]interface{}{
			"lockout_minutes": s.config.BruteForce.LockoutMinutes,
		})
	}

	return nil
}

// recordLoginFailure 记录登录失败
func (s *Service) recordLoginFailure(userID int64) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("brute_force:uid:%d", userID)

	// 使用Pipeline原子操作
	pipe := s.redis.Pipeline()
	incr := pipe.Incr(key)
	if _, err := pipe.Exec(); err != nil {
		slog.Error("登录失败记录写入失败", "error", err, "user_id", userID)
		return
	}

	// 如果是第一次失败（计数为1），设置过期时间
	if incr.Val() == 1 {
		s.redis.Expire(key, time.Duration(s.config.BruteForce.LockoutMinutes)*time.Minute)
	}
}

// clearLoginFailure 清除登录失败记录
func (s *Service) clearLoginFailure(userID int64) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("brute_force:uid:%d", userID)
	s.redis.Del(key)
}

// ActivateUser 激活用户
func (s *Service) ActivateUser(uid string, sign string) error {
	// 将字符串转换为整数类型的用户ID
	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_user_id_format"})
	}

	// 检查验证码是否有效
	result, validID := s.repo.CheckValidCode(userID, sign, s.config.ValidCode.TypeActiveUser)
	if !result {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_or_expired_code"})
	}

	// 更新用户状态
	if rows, err := s.repo.UpdateUserStatus(userID, s.config.User.StatusOK); err != nil || rows != 1 {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()})
	}

	// 删除验证码
	if rows, err := s.repo.DeleteValidCode(validID); err != nil || rows != 1 {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()})
	}

	return nil
}

// ===================== 用户信息验证 =====================

// IsExistEmail 判断邮件地址是否存在
func (s *Service) IsExistEmail(email string) (bool, error) {
	// 判断邮件地址是否存在
	isExistEmail, err := s.repo.IsExistsEmail(email)
	if err != nil {
		return false, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "email_check_failed"})
	}
	return isExistEmail, nil
}

// IsExistHandle 判断Handle是否存在
func (s *Service) IsExistHandle(handle string) (bool, error) {
	// 判断Handle是否存在
	isExistHandle, err := s.repo.IsExistHandle(handle)
	if err != nil {
		return false, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "handle_check_failed"})
	}
	return isExistHandle, nil
}

// ===================== 密码管理 =====================

// CreateResetPasswordCode 创建重置密码验证码
func (s *Service) CreateResetPasswordCode(userID int64) error {
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(s.config.ValidCode.Length, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   s.config.ValidCode.TypeResetPwd,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().Add(time.Duration(s.config.ValidCode.ExpireMinutes) * time.Minute).UnixMilli(), // 使用配置中的分钟数
	}

	_, err := s.repo.InsertValidCode(validCode)
	return err
}

// SendResetPasswordEmail 发送重置密码邮件
func (s *Service) SendResetPasswordEmail(user *User) error {
	// 创建重置密码验证码
	validCode := &ValidCode{
		ValidCode:     string(pkg.Rand(s.config.ValidCode.Length, pkg.RAND_KIND_ALL)),
		UserId:        user.UserId,
		ValidTypeId:   s.config.ValidCode.TypeResetPwd,
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
		return errors.InvalidPassword.WithDetails(map[string]interface{}{"reason": "incorrect_current_password"})
	}

	// 验证新密码强度
	if len(newPwd) < s.config.User.PwdMinLen || len(newPwd) > s.config.User.PwdMaxLen {
		return errors.ValidationError.WithDetails(map[string]interface{}{
			"reason": "password_length_invalid",
			"min":    s.config.User.PwdMinLen,
			"max":    s.config.User.PwdMaxLen,
		})
	}

	// 使用统一的密码生成逻辑
	newPwdHash, err := s.generatePasswordHash(newPwd)
	if err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "password_hash_failed"})
	}

	// 更新用户密码
	if err := s.UpdateUserPassword(currentUserID, newPwdHash); err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "password_update_failed"})
	}

	return nil
}

// ResetUserPassword 重置用户密码
// 用于忘记密码后通过验证码重置密码，使用统一的密码生成逻辑
func (s *Service) ResetUserPassword(userId int64, newPassword string) error {
	// 验证新密码强度
	if len(newPassword) < s.config.User.PwdMinLen || len(newPassword) > s.config.User.PwdMaxLen {
		return errors.ValidationError.WithDetails(map[string]interface{}{
			"reason": "password_length_invalid",
			"min":    s.config.User.PwdMinLen,
			"max":    s.config.User.PwdMaxLen,
		})
	}

	// 使用统一的密码生成逻辑
	newPwdHash, err := s.generatePasswordHash(newPassword)
	if err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "password_hash_failed"})
	}

	// 更新用户密码
	if err := s.UpdateUserPassword(userId, newPwdHash); err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "password_reset_failed"})
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

// GetUserTimeZone 获取用户时区
func (s *Service) GetUserTimeZone(userId int64) (string, error) {
	return s.repo.GetUserTimeZone(userId)
}

// ParseUserID 解析用户ID（已弃用，建议使用validation.ParseID）
// 为了保持向后兼容性暂时保留，内部调用validation.ParseID
func (s *Service) ParseUserID(idStr string) (int64, error) {
	// 使用统一的ID解析方法
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_user_id_format"})
	}
	return id, nil
}

// ===================== 个人资料管理 =====================
// UpdateUserProfile 更新用户基本资料
// 更新用户的全名、个人简介、时区和语言等信息
func (s *Service) UpdateUserProfile(uid int64, displayName, remark, timeZone, lang string) error {

	// 验证语言格式
	if lang != "" {
		supportedLang := s.i18nService.IsLanguageSupported(lang)
		if !supportedLang {
			return &errors.InvalidParams
		}
	}

	// 调用repository更新用户资料
	rowsAffected, err := s.repo.UpdateUserProfile(uid, displayName, remark, timeZone, lang)
	if err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "user_profile_update_failed"})
	}

	// 检查是否更新成功
	if rowsAffected == 0 {
		// 如果影响行数为0，可能是因为数据没有变化，也可能是用户不存在
		// 我们需要通过查询用户是否存在来区分
		_, err := s.repo.GetUserByUserId(uid)
		if err != nil {
			// 如果查询用户报错（包括用户不存在），则返回错误
			return err
		}
		// 用户存在但没有更新行数，说明数据未变更，视为成功
	}

	return nil
}

// UpdateUserHandle 更新用户Handle
func (s *Service) UpdateUserHandle(userId int64, newHandle string) error {
	// 1. 验证 Handle 格式
	if !pkg.IsHandle(newHandle) {
		return errors.ValidationError.WithDetails(map[string]interface{}{"reason": "invalid_handle_format"})
	}

	// 2. 获取当前用户信息
	user, err := s.repo.GetUserByUserId(userId)
	if err != nil {
		return err
	}

	// 3. 检查冷却时间
	now := pkg.GetNowTimeUnixMilli()
	cooldownMs := int64(s.config.User.HandleCooldownDays) * 24 * 3600 * 1000
	if now-user.ModifiedUTCTime < cooldownMs {
		remainingMs := cooldownMs - (now - user.ModifiedUTCTime)
		daysRemaining := int64(math.Ceil(float64(remainingMs) / float64(24*3600*1000)))
		return errors.ValidationError.WithDetails(map[string]interface{}{
			"reason":         "handle_cooldown_active",
			"days_remaining": daysRemaining,
		})
	}

	// 4. 检查 Handle 是否已存在
	isExist, err := s.IsExistHandle(newHandle)
	if err != nil {
		return err
	}
	if isExist {
		return errors.UserExists.WithDetails(map[string]interface{}{"handle": newHandle, "reason": "handle_exists"})
	}

	// 5. 更新 Handle 和 ModifiedUTCTime
	_, err = s.repo.UpdateUserHandle(userId, newHandle, now)
	if err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "handle_update_failed"})
	}

	return nil
}

// UpdateUserEmail 更新用户邮箱
func (s *Service) UpdateUserEmail(uid int64, email string) error {
	// 验证邮箱格式
	if !pkg.IsEmail(email) {
		return errors.InvalidFormat.WithDetails(map[string]interface{}{"reason": "invalid_email_format"})
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
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code() == errors.UserNotFound.Code() {
			// 如果是 UserNotFound 错误，说明邮箱未被使用，这是正常的，继续执行
		} else {
			// 如果是其他数据库错误，则返回数据库错误
			return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "email_query_failed"})
		}
	} else if existingUser != nil && existingUser.UserId != uid {
		// 如果找到了用户且不是当前用户，则表示邮箱已被占用
		return errors.EmailExists.WithDetails(map[string]interface{}{"reason": "email_already_used"})
	}

	// 调用 repository 更新邮箱
	if _, err := s.repo.UpdateEmailByUserId(uid, email); err != nil {
		return errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "email_update_failed"})
	}

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *Service) UpdateUserStatus(userId int64, status int8) (int64, error) {
	rowsAffected, err := s.repo.UpdateUserStatus(userId, status)
	if err != nil {
		return 0, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "user_status_update_failed"})
	}
	// 检查是否更新成功
	if rowsAffected == 0 {
		return 0, errors.DatabaseError.WithDetails(map[string]interface{}{"reason": "user_not_found_or_status_unchanged"})
	}
	return rowsAffected, nil
}

// UpdateUserPwd 更新用户密码
func (s *Service) UpdateUserPwd(userId int64, pwd string) (int64, error) {
	rowsAffected, err := s.repo.UpdateUserPwd(userId, pwd)
	if err != nil {
		return 0, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error(), "reason": "password_update_failed"})
	}
	return rowsAffected, nil
}

// ===================== 私有辅助方法 =====================

// validateUserInput 验证用户输入参数
func (s *Service) validateUserInput(user *User) error {
	// 判断参数是否到位
	if !pkg.IsHandle(user.Handle) {
		return errors.ValidationError.WithDetails(map[string]interface{}{"reason": "invalid_handle_format"})
	}

	// 判断密码格式
	if !pkg.IsPwd(user.Pwd) {
		return errors.ValidationError.WithDetails(map[string]interface{}{"reason": "invalid_password_format"})
	}

	// 判断用户时区与语言是否正确
	if !pkg.IsTimeZone(user.TimeZone) {
		return errors.ValidationError.WithDetails(map[string]interface{}{"reason": "invalid_timezone_format"})
	}

	// 判断Handle是否存在
	isExistHandle, err := s.IsExistHandle(user.Handle)
	if err != nil {
		return err
	}
	if isExistHandle {
		return errors.UserExists.WithDetails(map[string]interface{}{"handle": user.Handle, "reason": "handle_exists"})
	}

	// 判断邮件地址是否存在
	isExistEmail, err := s.IsExistEmail(user.Email)
	if err != nil {
		return err
	}
	if isExistEmail {
		return errors.EmailExists.WithDetails(map[string]interface{}{"email": user.Email, "reason": "email_exists"})
	}

	return nil
}

// normalizeUserInput 转换用户名和邮件地址为小写并去除空格
func (s *Service) normalizeUserInput(user *User) {
	user.Handle = strings.TrimSpace(strings.ToLower(user.Handle))
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
		slog.Error("密码哈希生成失败", "error", err)
	}
	user.Pwd = hashedPwd

	// 用户默认分组为1
	user.UserGroup = s.config.User.DefaultGroup
	// 获取当前UTC时间
	nowTime := pkg.GetNowTimeUnixMilli()
	user.CreatedUTCTime = nowTime
	user.ModifiedUTCTime = nowTime
	user.LastLoginUTCTime = nowTime
	// 设置默认语言
	user.Lang = s.config.App.DefaultLang
	user.Status = s.config.User.StatusWaitActive
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
		ValidCode:     string(pkg.Rand(s.config.ValidCode.Length, pkg.RAND_KIND_ALL)),
		UserId:        userID,
		ValidTypeId:   s.config.ValidCode.TypeActiveUser,
		CreateUTCTime: pkg.GetNowTime().UnixMilli(),
		ExpireUTCTime: pkg.GetNowTime().Add(time.Duration(s.config.ValidCode.ExpireMinutes) * time.Minute).UnixMilli(),
	}

	if _, err := s.repo.InsertValidCode(validCode); err != nil {
		return nil, err
	}

	return validCode, nil
}
