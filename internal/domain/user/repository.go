package user

// RepositoryInterface 用户数据访问层接口
// 定义所有与用户相关的数据库操作，包括用户管理和验证码管理
type RepositoryInterface interface {
	// ===================== User 相关方法 =====================

	// InsertUser 插入用户数据
	InsertUser(user *User) (int64, error)

	// UpdateUser 更新用户数据
	UpdateUser(user *User) (int64, error)

	// DeleteByUserId 删除用户数据（假删除，设置status = 0）
	DeleteByUserId(id int64) (int64, error)

	// IsExistHandle 判断Handle是否存在
	IsExistHandle(handle string) (bool, error)

	// IsExistsEmail 判断邮箱是否存在
	IsExistsEmail(email string) (bool, error)

	// UpdateEmailByUserId 根据用户ID更新邮箱
	UpdateEmailByUserId(userId int64, email string) (int64, error)

	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(userId int64, status int8) (int64, error)

	// UpdateUserPwd 更新用户密码
	UpdateUserPwd(userId int64, pwd string) (int64, error)

	// DeleteUser 彻底删除用户
	DeleteUser(userId int64) (int64, error)

	// GetUserByUserId 根据用户ID获取用户信息
	GetUserByUserId(userId int64) (*User, error)

	// GetUserByHandle 根据Handle获取用户信息
	GetUserByHandle(handle string) (*User, error)

	// GetUserByEmail 根据邮箱获取用户信息
	GetUserByEmail(email string) (*User, error)

	// GetUserEncryptPwd 获取用户加密密码
	GetUserEncryptPwd(userId int64) (string, error)

	// GetUserTimeZone 获取用户时区
	GetUserTimeZone(userId int64) (string, error)

	// UpdateUserProfile 更新用户基本资料
	UpdateUserProfile(userId int64, displayName, remark, timeZone, lang string) (int64, error)

	// UpdateUserPassword 更新用户密码
	UpdateUserPassword(userId int64, pwd string) (int64, error)

	// UpdateUserHandle 更新用户Handle
	UpdateUserHandle(userId int64, handle string, modifiedUTCTime int64) (int64, error)

	// ===================== ValidCode 相关方法 =====================

	// InsertValidCode 插入验证码
	InsertValidCode(validCode *ValidCode) (int64, error)

	// DeleteValidCode 删除验证码
	DeleteValidCode(id int64) (int64, error)

	// GetValidCode 获取验证码
	GetValidCode(userId int64, code string) (*ValidCode, error)

	// CheckValidCode 检查确认码
	CheckValidCode(userID int64, sign string, validType int8) (bool, int64)
}
