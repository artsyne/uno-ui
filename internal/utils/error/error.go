package errors

import (
	"fmt"
)

// ErrorDef 错误定义结构体，统一管理错误码、错误键和HTTP状态码
type ErrorDef struct {
	Code     int    `json:"code"`      // 错误数字码
	Key      string `json:"key"`       // 错误标识符
	HttpCode int    `json:"http_code"` // HTTP状态码
}

// Error 实现 error 接口，返回 Key 作为错误字符串
func (def ErrorDef) Error() string {
	return def.Key
}

// Err 创建应用错误
func (def ErrorDef) Err() *AppError {
	return &AppError{
		ErrorDef: def,
	}
}

// WithDetails 创建带详细信息的应用错误
func (def ErrorDef) WithDetails(details string) *AppError {
	return &AppError{
		ErrorDef: def,
		Details:  details,
	}
}

// 统一的错误定义 - 所有错误信息在一个地方维护
var (
	// 1000-1099: 参数验证错误 (400)
	ValidationError      = ErrorDef{1001, "validation_error", 400}
	InvalidParams        = ErrorDef{1002, "invalid_params", 400}
	MissingRequiredField = ErrorDef{1003, "missing_required_field", 400}
	InvalidFormat        = ErrorDef{1004, "invalid_format", 400}
	ValueOutOfRange      = ErrorDef{1005, "value_out_of_range", 400}
	InvalidPassword      = ErrorDef{1006, "invalid_password", 400}
	InvalidEmail         = ErrorDef{1007, "invalid_email", 400}
	InvalidUsername      = ErrorDef{1008, "invalid_username", 400}
	PasswordTooShort     = ErrorDef{1009, "password_too_short", 400}
	PasswordTooLong      = ErrorDef{1010, "password_too_long", 400}
	FieldTooShort        = ErrorDef{1011, "field_too_short", 400}
	FieldTooLong         = ErrorDef{1012, "field_too_long", 400}
	FieldMismatch        = ErrorDef{1013, "field_mismatch", 400}

	// 1100-1199: 认证错误 (401)
	Unauthorized       = ErrorDef{1101, "unauthorized", 401}
	InvalidToken       = ErrorDef{1102, "invalid_token", 401}
	TokenExpired       = ErrorDef{1103, "token_expired", 401}
	LoginRequired      = ErrorDef{1104, "login_required", 401}
	InvalidCredentials = ErrorDef{1105, "invalid_credentials", 401}

	// 1200-1299: 权限错误 (403)
	PermissionDenied       = ErrorDef{1201, "permission_denied", 403}
	AccessForbidden        = ErrorDef{1202, "access_forbidden", 403}
	InsufficientPrivileges = ErrorDef{1203, "insufficient_privileges", 403}
	OperationNotAllowed    = ErrorDef{1204, "operation_not_allowed", 403}

	// 1300-1399: 资源不存在 (404)
	UserNotFound     = ErrorDef{1301, "user_not_found", 404}
	ResourceNotFound = ErrorDef{1302, "resource_not_found", 404}
	DiaryNotFound    = ErrorDef{1303, "diary_not_found", 404}
	PostNotFound     = ErrorDef{1304, "post_not_found", 404}
	PageNotFound     = ErrorDef{1305, "page_not_found", 404}

	// 2000-2999: 用户相关业务错误 (400)
	UserExists        = ErrorDef{2001, "user_exists", 400}
	EmailExists       = ErrorDef{2002, "email_exists", 400}
	UsernameExists    = ErrorDef{2003, "username_exists", 400}
	DuplicateResource = ErrorDef{2004, "duplicate_resource", 400}

	// 2100-2199: 用户状态错误 (400)
	UserInactive     = ErrorDef{2101, "user_inactive", 400}
	UserBanned       = ErrorDef{2102, "user_banned", 400}
	AccountLocked    = ErrorDef{2103, "account_locked", 400}
	EmailNotVerified = ErrorDef{2104, "email_not_verified", 400}

	// 2200-2299: 密码相关错误 (400)
	PasswordIncorrect    = ErrorDef{2201, "password_incorrect", 400}
	PasswordTooWeak      = ErrorDef{2202, "password_too_weak", 400}
	PasswordExpired      = ErrorDef{2203, "password_expired", 400}
	OldPasswordIncorrect = ErrorDef{2204, "old_password_incorrect", 400}

	// 2300-2399: 限制相关错误 (400)
	DiaryLimitExceeded = ErrorDef{2301, "diary_limit_exceeded", 400}
	ContentTooLong     = ErrorDef{2302, "content_too_long", 400}
	FileTooLarge       = ErrorDef{2303, "file_too_large", 400}
	QuotaExceeded      = ErrorDef{2304, "quota_exceeded", 400}

	// 2400-2499: 通信相关错误 (400)
	EmailSendFailed    = ErrorDef{2401, "email_send_failed", 400}
	SMSSendFailed      = ErrorDef{2402, "sms_send_failed", 400}
	NotificationFailed = ErrorDef{2403, "notification_failed", 400}

	// 2500-2599: 文件处理错误 (400)
	UploadFailed          = ErrorDef{2501, "upload_failed", 400}
	ImageProcessFailed    = ErrorDef{2502, "image_process_failed", 400}
	FileFormatUnsupported = ErrorDef{2503, "file_format_unsupported", 400}

	// 6000-6099: 数据库错误 (500)
	DatabaseError     = ErrorDef{6001, "database_error", 500}
	ConnectionFailed  = ErrorDef{6002, "connection_failed", 500}
	QueryTimeout      = ErrorDef{6003, "query_timeout", 500}
	TransactionFailed = ErrorDef{6004, "transaction_failed", 500}

	// 6100-6199: 缓存错误 (500)
	RedisError   = ErrorDef{6101, "redis_error", 500}
	CacheError   = ErrorDef{6102, "cache_error", 500}
	SessionError = ErrorDef{6103, "session_error", 500}

	// 6200-6299: 服务错误 (500)
	ConfigError          = ErrorDef{6201, "config_error", 500}
	ServiceUnavailable   = ErrorDef{6202, "service_unavailable", 500}
	ExternalServiceError = ErrorDef{6203, "external_service_error", 500}

	// 6250-6269: DI容器错误 (500)
	DIServiceAlreadyRegistered = ErrorDef{6251, "di_service_already_registered", 500}
	DIServiceNotRegistered     = ErrorDef{6252, "di_service_not_registered", 500}

	// 6300-6399: 系统负载错误 (500)
	RateLimitExceeded = ErrorDef{6301, "rate_limit_exceeded", 500}
	ServerOverloaded  = ErrorDef{6302, "server_overloaded", 500}
	MaintenanceMode   = ErrorDef{6303, "maintenance_mode", 500}

	// 9999: 通用内部错误 (500)
	InternalError = ErrorDef{9999, "internal_error", 500}
)

// AppError 应用错误结构体
type AppError struct {
	ErrorDef        // 嵌入错误定义
	Details  string `json:"-"` // 错误详情
}

// Error 实现error接口
func (e *AppError) Error() string {
    if e.Details != "" {
        // 如果有详细信息，则输出：[Code] Key: Details
        return fmt.Sprintf("[%d] %s: %s", e.Code, e.Key, e.Details)
    }
    // 否则只输出：[Code] Key
    return fmt.Sprintf("[%d] %s", e.Code, e.Key)
}