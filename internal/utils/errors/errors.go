package errors

import (
	"fmt"
)

// AppError 统一的错误结构体，包含所有必要信息
type AppError struct {
	code     int    // 错误数字码（私有）
	key      string // 错误标识符（私有）
	httpCode int    // HTTP状态码（私有）
	details  string // 错误详情（私有，可选）
}

// NewAppError 创建新的错误定义（构造函数）
func NewAppError(code int, key string, httpCode int) AppError {
	return AppError{
		code:     code,
		key:      key,
		httpCode: httpCode,
		details:  "", // 默认无详情
	}
}

// Code 获取错误码（只读访问）
func (def AppError) Code() int {
	return def.code
}

// Key 获取错误键（只读访问）
func (def AppError) Key() string {
	return def.key
}

// HttpCode 获取HTTP状态码（只读访问）
func (def AppError) HttpCode() int {
	return def.httpCode
}

// Details 获取错误详情（只读访问）
func (def AppError) Details() string {
	return def.details
}

// Error 实现 error 接口
func (def AppError) Error() string {
	if def.details != "" {
		return fmt.Sprintf("[%d] %s: %s", def.code, def.key, def.details)
	}
	return fmt.Sprintf("[%d] %s", def.code, def.key)
}

// WithDetails 返回带详细信息的新错误实例
func (def AppError) WithDetails(details string) AppError {
	return AppError{
		code:     def.code,
		key:      def.key,
		httpCode: def.httpCode,
		details:  details,
	}
}

// 统一的错误定义 - 所有错误信息在一个地方维护
// 使用构造函数确保数据完整性和不可变性
var (
	// 1000-1099: 参数验证错误 (400)
	ValidationError      = NewAppError(1001, "validation_error", 400)
	InvalidParams        = NewAppError(1002, "invalid_params", 400)
	MissingRequiredField = NewAppError(1003, "missing_required_field", 400)
	InvalidFormat        = NewAppError(1004, "invalid_format", 400)
	ValueOutOfRange      = NewAppError(1005, "value_out_of_range", 400)
	InvalidPassword      = NewAppError(1006, "invalid_password", 400)
	InvalidEmail         = NewAppError(1007, "invalid_email", 400)
	InvalidUsername      = NewAppError(1008, "invalid_username", 400)
	PasswordTooShort     = NewAppError(1009, "password_too_short", 400)
	PasswordTooLong      = NewAppError(1010, "password_too_long", 400)
	FieldTooShort        = NewAppError(1011, "field_too_short", 400)
	FieldTooLong         = NewAppError(1012, "field_too_long", 400)
	FieldMismatch        = NewAppError(1013, "field_mismatch", 400)

	// 1100-1199: 认证错误 (401)
	Unauthorized       = NewAppError(1101, "unauthorized", 401)
	InvalidToken       = NewAppError(1102, "invalid_token", 401)
	TokenExpired       = NewAppError(1103, "token_expired", 401)
	LoginRequired      = NewAppError(1104, "login_required", 401)
	InvalidCredentials = NewAppError(1105, "invalid_credentials", 401)

	// 1200-1299: 权限错误 (403)
	PermissionDenied       = NewAppError(1201, "permission_denied", 403)
	AccessForbidden        = NewAppError(1202, "access_forbidden", 403)
	InsufficientPrivileges = NewAppError(1203, "insufficient_privileges", 403)
	OperationNotAllowed    = NewAppError(1204, "operation_not_allowed", 403)

	// 1300-1399: 资源不存在 (404)
	UserNotFound     = NewAppError(1301, "user_not_found", 404)
	ResourceNotFound = NewAppError(1302, "resource_not_found", 404)
	DiaryNotFound    = NewAppError(1303, "diary_not_found", 404)
	PostNotFound     = NewAppError(1304, "post_not_found", 404)
	PageNotFound     = NewAppError(1305, "page_not_found", 404)

	// 2000-2999: 用户相关业务错误 (400)
	UserExists        = NewAppError(2001, "user_exists", 400)
	EmailExists       = NewAppError(2002, "email_exists", 400)
	UsernameExists    = NewAppError(2003, "username_exists", 400)
	DuplicateResource = NewAppError(2004, "duplicate_resource", 400)

	// 2100-2199: 用户状态错误 (400)
	UserInactive     = NewAppError(2101, "user_inactive", 400)
	UserBanned       = NewAppError(2102, "user_banned", 400)
	AccountLocked    = NewAppError(2103, "account_locked", 400)
	EmailNotVerified = NewAppError(2104, "email_not_verified", 400)

	// 2200-2299: 密码相关错误 (400)
	PasswordIncorrect    = NewAppError(2201, "password_incorrect", 400)
	PasswordTooWeak      = NewAppError(2202, "password_too_weak", 400)
	PasswordExpired      = NewAppError(2203, "password_expired", 400)
	OldPasswordIncorrect = NewAppError(2204, "old_password_incorrect", 400)

	// 2300-2399: 限制相关错误 (400)
	DiaryLimitExceeded = NewAppError(2301, "diary_limit_exceeded", 400)
	ContentTooLong     = NewAppError(2302, "content_too_long", 400)
	FileTooLarge       = NewAppError(2303, "file_too_large", 400)
	QuotaExceeded      = NewAppError(2304, "quota_exceeded", 400)

	// 2400-2499: 通信相关错误 (400)
	EmailSendFailed    = NewAppError(2401, "email_send_failed", 400)
	SMSSendFailed      = NewAppError(2402, "sms_send_failed", 400)
	NotificationFailed = NewAppError(2403, "notification_failed", 400)

	// 2500-2599: 文件处理错误 (400)
	UploadFailed          = NewAppError(2501, "upload_failed", 400)
	ImageProcessFailed    = NewAppError(2502, "image_process_failed", 400)
	FileFormatUnsupported = NewAppError(2503, "file_format_unsupported", 400)

	// 6000-6099: 数据库错误 (500)
	DatabaseError     = NewAppError(6001, "database_error", 500)
	ConnectionFailed  = NewAppError(6002, "connection_failed", 500)
	QueryTimeout      = NewAppError(6003, "query_timeout", 500)
	TransactionFailed = NewAppError(6004, "transaction_failed", 500)

	// 6100-6199: 缓存错误 (500)
	RedisError   = NewAppError(6101, "redis_error", 500)
	CacheError   = NewAppError(6102, "cache_error", 500)
	SessionError = NewAppError(6103, "session_error", 500)

	// 6200-6299: 服务错误 (500)
	ConfigError          = NewAppError(6201, "config_error", 500)
	ServiceUnavailable   = NewAppError(6202, "service_unavailable", 500)
	ExternalServiceError = NewAppError(6203, "external_service_error", 500)

	// 6250-6269: DI容器错误 (500)
	DIServiceAlreadyRegistered = NewAppError(6251, "di_service_already_registered", 500)
	DIServiceNotRegistered     = NewAppError(6252, "di_service_not_registered", 500)

	// 6300-6399: 系统负载错误 (500)
	RateLimitExceeded = NewAppError(6301, "rate_limit_exceeded", 500)
	ServerOverloaded  = NewAppError(6302, "server_overloaded", 500)
	MaintenanceMode   = NewAppError(6303, "maintenance_mode", 500)

	// 9999: 通用内部错误 (500)
	InternalError = NewAppError(9999, "internal_error", 500)
)