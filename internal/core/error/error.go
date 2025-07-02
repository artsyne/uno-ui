package errors

import (
	"fmt"
)

// 错误码常量定义 - 按功能模块和HTTP状态码范围分组
const (
	// 1000-1099: 参数验证错误 (400)
	ValidationError      = 1001 // 参数验证失败
	InvalidParams        = 1002 // 参数无效
	MissingRequiredField = 1003 // 缺少必填字段
	InvalidFormat        = 1004 // 格式无效
	ValueOutOfRange      = 1005 // 值超出范围

	// 1100-1199: 认证错误 (401)
	Unauthorized       = 1101 // 未授权访问
	InvalidToken       = 1102 // 无效令牌
	TokenExpired       = 1103 // 令牌已过期
	LoginRequired      = 1104 // 需要登录
	InvalidCredentials = 1105 // 凭据无效

	// 1200-1299: 权限错误 (403)
	PermissionDenied       = 1201 // 权限被拒绝
	AccessForbidden        = 1202 // 访问被禁止
	InsufficientPrivileges = 1203 // 权限不足
	OperationNotAllowed    = 1204 // 操作不被允许

	// 1300-1399: 资源不存在 (404)
	UserNotFound     = 1301 // 用户不存在
	ResourceNotFound = 1302 // 资源不存在
	DiaryNotFound    = 1303 // 日记不存在
	PostNotFound     = 1304 // 文章不存在
	PageNotFound     = 1305 // 页面不存在

	// 2000-2999: 用户相关业务错误 (400)
	UserExists        = 2001 // 用户已存在
	EmailExists       = 2002 // 邮箱已存在
	UsernameExists    = 2003 // 用户名已存在
	DuplicateResource = 2004 // 资源重复

	// 2100-2199: 用户状态错误 (400)
	UserInactive     = 2101 // 用户未激活
	UserBanned       = 2102 // 用户被封禁
	AccountLocked    = 2103 // 账户被锁定
	EmailNotVerified = 2104 // 邮箱未验证

	// 2200-2299: 密码相关错误 (400)
	PasswordIncorrect    = 2201 // 密码错误
	PasswordTooWeak      = 2202 // 密码太弱
	PasswordExpired      = 2203 // 密码已过期
	OldPasswordIncorrect = 2204 // 旧密码错误

	// 2300-2399: 限制相关错误 (400)
	DiaryLimitExceeded = 2301 // 日记数量超限
	ContentTooLong     = 2302 // 内容太长
	FileTooLarge       = 2303 // 文件太大
	QuotaExceeded      = 2304 // 配额超限

	// 2400-2499: 通信相关错误 (400)
	EmailSendFailed    = 2401 // 邮件发送失败
	SMSSendFailed      = 2402 // 短信发送失败
	NotificationFailed = 2403 // 通知发送失败

	// 2500-2599: 文件处理错误 (400)
	UploadFailed          = 2501 // 上传失败
	ImageProcessFailed    = 2502 // 图片处理失败
	FileFormatUnsupported = 2503 // 文件格式不支持

	// 6000-6099: 数据库错误 (500)
	DatabaseError     = 6001 // 数据库错误
	ConnectionFailed  = 6002 // 连接失败
	QueryTimeout      = 6003 // 查询超时
	TransactionFailed = 6004 // 事务失败

	// 6100-6199: 缓存错误 (500)
	RedisError   = 6101 // Redis错误
	CacheError   = 6102 // 缓存错误
	SessionError = 6103 // 会话错误

	// 6200-6299: 服务错误 (500)
	ConfigError          = 6201 // 配置错误
	ServiceUnavailable   = 6202 // 服务不可用
	ExternalServiceError = 6203 // 外部服务错误

	// 6300-6399: 系统负载错误 (500)
	RateLimitExceeded = 6301 // 频率限制超出
	ServerOverloaded  = 6302 // 服务器过载
	MaintenanceMode   = 6303 // 维护模式

	// 9999: 通用内部错误 (500)
	InternalError = 9999 // 内部错误
)

// 错误码到错误键的映射
var errorKeyMap = map[int]string{
	// 参数验证错误
	ValidationError:      "validation_error",
	InvalidParams:        "invalid_params",
	MissingRequiredField: "missing_required_field",
	InvalidFormat:        "invalid_format",
	ValueOutOfRange:      "value_out_of_range",

	// 认证错误
	Unauthorized:       "unauthorized",
	InvalidToken:       "invalid_token",
	TokenExpired:       "token_expired",
	LoginRequired:      "login_required",
	InvalidCredentials: "invalid_credentials",

	// 权限错误
	PermissionDenied:       "permission_denied",
	AccessForbidden:        "access_forbidden",
	InsufficientPrivileges: "insufficient_privileges",
	OperationNotAllowed:    "operation_not_allowed",

	// 资源不存在
	UserNotFound:     "user_not_found",
	ResourceNotFound: "resource_not_found",
	DiaryNotFound:    "diary_not_found",
	PostNotFound:     "post_not_found",
	PageNotFound:     "page_not_found",

	// 用户相关业务错误
	UserExists:        "user_exists",
	EmailExists:       "email_exists",
	UsernameExists:    "username_exists",
	DuplicateResource: "duplicate_resource",

	// 用户状态错误
	UserInactive:     "user_inactive",
	UserBanned:       "user_banned",
	AccountLocked:    "account_locked",
	EmailNotVerified: "email_not_verified",

	// 密码相关错误
	PasswordIncorrect:    "password_incorrect",
	PasswordTooWeak:      "password_too_weak",
	PasswordExpired:      "password_expired",
	OldPasswordIncorrect: "old_password_incorrect",

	// 限制相关错误
	DiaryLimitExceeded: "diary_limit_exceeded",
	ContentTooLong:     "content_too_long",
	FileTooLarge:       "file_too_large",
	QuotaExceeded:      "quota_exceeded",

	// 通信相关错误
	EmailSendFailed:    "email_send_failed",
	SMSSendFailed:      "sms_send_failed",
	NotificationFailed: "notification_failed",

	// 文件处理错误
	UploadFailed:          "upload_failed",
	ImageProcessFailed:    "image_process_failed",
	FileFormatUnsupported: "file_format_unsupported",

	// 数据库错误
	DatabaseError:     "database_error",
	ConnectionFailed:  "connection_failed",
	QueryTimeout:      "query_timeout",
	TransactionFailed: "transaction_failed",

	// 缓存错误
	RedisError:   "redis_error",
	CacheError:   "cache_error",
	SessionError: "session_error",

	// 服务错误
	ConfigError:          "config_error",
	ServiceUnavailable:   "service_unavailable",
	ExternalServiceError: "external_service_error",

	// 系统负载错误
	RateLimitExceeded: "rate_limit_exceeded",
	ServerOverloaded:  "server_overloaded",
	MaintenanceMode:   "maintenance_mode",

	// 通用内部错误
	InternalError: "internal_error",
}

// AppError 极简应用错误结构体
type AppError struct {
	Code     int    `json:"code"`    // 错误数字码
	MsgKey   string `json:"msg_key"` // 错误标识符
	HttpCode int    `json:"-"`       // HTTP状态码（自动推断）
	Details  string `json:"-"`       // 错误详情（包含所有额外信息）
}

// getHttpStatusByErrorNum 根据错误码范围自动确定HTTP状态码
func getHttpStatusByErrorNum(errorNum int) int {
	switch {
	case errorNum >= 1000 && errorNum <= 1099:
		return 400 // 参数验证错误
	case errorNum >= 1100 && errorNum <= 1199:
		return 401 // 认证错误
	case errorNum >= 1200 && errorNum <= 1299:
		return 403 // 权限错误
	case errorNum >= 1300 && errorNum <= 1399:
		return 404 // 资源不存在
	case errorNum >= 2000 && errorNum <= 5999:
		return 400 // 业务逻辑错误
	case errorNum >= 6000 && errorNum <= 9999:
		return 500 // 系统内部错误
	default:
		return 500 // 未知错误
	}
}

// getErrorKey 根据错误码获取错误键
func getErrorKey(errorCode int) string {
	if key, exists := errorKeyMap[errorCode]; exists {
		return key
	}
	return "unknown_error"
}

// NewAppError 创建应用错误
func NewAppError(errorCode int) *AppError {
	return &AppError{
		Code:     errorCode,
		MsgKey:   getErrorKey(errorCode),
		HttpCode: getHttpStatusByErrorNum(errorCode),
	}
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return e.Details
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.MsgKey)
}

// WithDetails 添加详细描述
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// ==================== 快捷错误创建函数 ====================
// 为常用错误提供简洁的创建方式

// 参数验证错误
func ValidationErr() *AppError    { return NewAppError(ValidationError) }
func InvalidParamsErr() *AppError { return NewAppError(InvalidParams) }
func MissingFieldErr() *AppError  { return NewAppError(MissingRequiredField) }
func InvalidFormatErr() *AppError { return NewAppError(InvalidFormat) }

// 认证错误
func UnauthorizedErr() *AppError  { return NewAppError(Unauthorized) }
func InvalidTokenErr() *AppError  { return NewAppError(InvalidToken) }
func TokenExpiredErr() *AppError  { return NewAppError(TokenExpired) }
func LoginRequiredErr() *AppError { return NewAppError(LoginRequired) }

// 权限错误
func PermissionErr() *AppError      { return NewAppError(PermissionDenied) }
func AccessForbiddenErr() *AppError { return NewAppError(AccessForbidden) }

// 资源不存在错误
func UserNotFoundErr() *AppError     { return NewAppError(UserNotFound) }
func ResourceNotFoundErr() *AppError { return NewAppError(ResourceNotFound) }
func DiaryNotFoundErr() *AppError    { return NewAppError(DiaryNotFound) }
func PostNotFoundErr() *AppError     { return NewAppError(PostNotFound) }

// 业务错误
func UserExistsErr() *AppError         { return NewAppError(UserExists) }
func EmailExistsErr() *AppError        { return NewAppError(EmailExists) }
func UsernameExistsErr() *AppError     { return NewAppError(UsernameExists) }
func UserInactiveErr() *AppError       { return NewAppError(UserInactive) }
func PasswordIncorrectErr() *AppError  { return NewAppError(PasswordIncorrect) }
func DiaryLimitExceededErr() *AppError { return NewAppError(DiaryLimitExceeded) }

// 文件处理错误
func UploadFailedErr() *AppError          { return NewAppError(UploadFailed) }
func FileTooLargeErr() *AppError          { return NewAppError(FileTooLarge) }
func FileFormatUnsupportedErr() *AppError { return NewAppError(FileFormatUnsupported) }

// 系统错误
func DatabaseErr() *AppError { return NewAppError(DatabaseError) }
func RedisErr() *AppError    { return NewAppError(RedisError) }
func InternalErr() *AppError { return NewAppError(InternalError) }
