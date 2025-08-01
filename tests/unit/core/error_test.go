package core

import (
	"testing"

	errorsPkg "github.com/biuren/hebe/internal/utils/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorDef_Err(t *testing.T) {
	// 测试ErrorDef.Err()方法
	errorDef := errorsPkg.ValidationError
	appError := errorDef.Err()
	
	require.NotNil(t, appError)
	assert.Equal(t, errorDef.Code, appError.Code)
	assert.Equal(t, errorDef.Key, appError.Key)
	assert.Equal(t, errorDef.HttpCode, appError.HttpCode)
	assert.Empty(t, appError.Details)
}

func TestErrorDef_WithDetails(t *testing.T) {
	// 测试ErrorDef.WithDetails()方法
	errorDef := errorsPkg.ValidationError
	details := "用户名不能为空"
	appError := errorDef.WithDetails(details)
	
	require.NotNil(t, appError)
	assert.Equal(t, errorDef.Code, appError.Code)
	assert.Equal(t, errorDef.Key, appError.Key)
	assert.Equal(t, errorDef.HttpCode, appError.HttpCode)
	assert.Equal(t, details, appError.Details)
}

func TestAppError_Error_WithoutDetails(t *testing.T) {
	// 测试AppError.Error()方法 - 无详细信息
	appError := errorsPkg.ValidationError.Err()
	errorMsg := appError.Error()
	
	expected := "[1001] validation_error"
	assert.Equal(t, expected, errorMsg)
}

func TestAppError_Error_WithDetails(t *testing.T) {
	// 测试AppError.Error()方法 - 有详细信息
	details := "用户名格式不正确"
	appError := errorsPkg.ValidationError.WithDetails(details)
	errorMsg := appError.Error()
	
	expected := "[1001] validation_error: 用户名格式不正确"
	assert.Equal(t, expected, errorMsg)
}

// 测试参数验证错误 (1000-1099)
func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"ValidationError", errorsPkg.ValidationError, 1001, "validation_error", 400},
		{"InvalidParams", errorsPkg.InvalidParams, 1002, "invalid_params", 400},
		{"MissingRequiredField", errorsPkg.MissingRequiredField, 1003, "missing_required_field", 400},
		{"InvalidFormat", errorsPkg.InvalidFormat, 1004, "invalid_format", 400},
		{"ValueOutOfRange", errorsPkg.ValueOutOfRange, 1005, "value_out_of_range", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试认证错误 (1100-1199)
func TestAuthenticationErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"Unauthorized", errorsPkg.Unauthorized, 1101, "unauthorized", 401},
		{"InvalidToken", errorsPkg.InvalidToken, 1102, "invalid_token", 401},
		{"TokenExpired", errorsPkg.TokenExpired, 1103, "token_expired", 401},
		{"LoginRequired", errorsPkg.LoginRequired, 1104, "login_required", 401},
		{"InvalidCredentials", errorsPkg.InvalidCredentials, 1105, "invalid_credentials", 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试权限错误 (1200-1299)
func TestPermissionErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"PermissionDenied", errorsPkg.PermissionDenied, 1201, "permission_denied", 403},
		{"AccessForbidden", errorsPkg.AccessForbidden, 1202, "access_forbidden", 403},
		{"InsufficientPrivileges", errorsPkg.InsufficientPrivileges, 1203, "insufficient_privileges", 403},
		{"OperationNotAllowed", errorsPkg.OperationNotAllowed, 1204, "operation_not_allowed", 403},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试资源不存在错误 (1300-1399)
func TestResourceNotFoundErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"UserNotFound", errorsPkg.UserNotFound, 1301, "user_not_found", 404},
		{"ResourceNotFound", errorsPkg.ResourceNotFound, 1302, "resource_not_found", 404},
		{"DiaryNotFound", errorsPkg.DiaryNotFound, 1303, "diary_not_found", 404},
		{"PostNotFound", errorsPkg.PostNotFound, 1304, "post_not_found", 404},
		{"PageNotFound", errorsPkg.PageNotFound, 1305, "page_not_found", 404},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试用户相关业务错误 (2000-2999)
func TestUserBusinessErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"UserExists", errorsPkg.UserExists, 2001, "user_exists", 400},
		{"EmailExists", errorsPkg.EmailExists, 2002, "email_exists", 400},
		{"UsernameExists", errorsPkg.UsernameExists, 2003, "username_exists", 400},
		{"DuplicateResource", errorsPkg.DuplicateResource, 2004, "duplicate_resource", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试用户状态错误 (2100-2199)
func TestUserStatusErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"UserInactive", errorsPkg.UserInactive, 2101, "user_inactive", 400},
		{"UserBanned", errorsPkg.UserBanned, 2102, "user_banned", 400},
		{"AccountLocked", errorsPkg.AccountLocked, 2103, "account_locked", 400},
		{"EmailNotVerified", errorsPkg.EmailNotVerified, 2104, "email_not_verified", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试密码相关错误 (2200-2299)
func TestPasswordErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"PasswordIncorrect", errorsPkg.PasswordIncorrect, 2201, "password_incorrect", 400},
		{"PasswordTooWeak", errorsPkg.PasswordTooWeak, 2202, "password_too_weak", 400},
		{"PasswordExpired", errorsPkg.PasswordExpired, 2203, "password_expired", 400},
		{"OldPasswordIncorrect", errorsPkg.OldPasswordIncorrect, 2204, "old_password_incorrect", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试数据库错误 (6000-6099)
func TestDatabaseErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"DatabaseError", errorsPkg.DatabaseError, 6001, "database_error", 500},
		{"ConnectionFailed", errorsPkg.ConnectionFailed, 6002, "connection_failed", 500},
		{"QueryTimeout", errorsPkg.QueryTimeout, 6003, "query_timeout", 500},
		{"TransactionFailed", errorsPkg.TransactionFailed, 6004, "transaction_failed", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试DI容器错误 (6250-6269)
func TestDIErrors(t *testing.T) {
	tests := []struct {
		name     string
		errorDef errorsPkg.ErrorDef
		code     int
		key      string
		httpCode int
	}{
		{"DIServiceAlreadyRegistered", errorsPkg.DIServiceAlreadyRegistered, 6251, "di_service_already_registered", 500},
		{"DIServiceNotRegistered", errorsPkg.DIServiceNotRegistered, 6252, "di_service_not_registered", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.errorDef.Code)
			assert.Equal(t, tt.key, tt.errorDef.Key)
			assert.Equal(t, tt.httpCode, tt.errorDef.HttpCode)
		})
	}
}

// 测试内部错误 (9999)
func TestInternalError(t *testing.T) {
	errorDef := errorsPkg.InternalError
	assert.Equal(t, 9999, errorDef.Code)
	assert.Equal(t, "internal_error", errorDef.Key)
	assert.Equal(t, 500, errorDef.HttpCode)
}

// 测试错误链式调用
func TestErrorChaining(t *testing.T) {
	// 测试从ErrorDef到AppError的链式调用
	appError := errorsPkg.ValidationError.WithDetails("测试详细信息")
	
	assert.Equal(t, 1001, appError.Code)
	assert.Equal(t, "validation_error", appError.Key)
	assert.Equal(t, 400, appError.HttpCode)
	assert.Equal(t, "测试详细信息", appError.Details)
	
	// 测试Error()方法
	errorMsg := appError.Error()
	assert.Equal(t, "[1001] validation_error: 测试详细信息", errorMsg)
}

// 测试错误作为Go error接口的使用
func TestErrorInterface(t *testing.T) {
	appError := errorsPkg.UserNotFound.WithDetails("用户ID: 123")
	
	// 测试作为error接口使用
	var err error = appError
	assert.NotNil(t, err)
	assert.Equal(t, "[1301] user_not_found: 用户ID: 123", err.Error())
	
	// 测试类型断言
	if appErr, ok := err.(*errorsPkg.AppError); ok {
		assert.Equal(t, 1301, appErr.Code)
		assert.Equal(t, "user_not_found", appErr.Key)
		assert.Equal(t, 404, appErr.HttpCode)
		assert.Equal(t, "用户ID: 123", appErr.Details)
	} else {
		t.Fatal("类型断言失败")
	}
}

// 测试错误码的唯一性（抽样测试）
func TestErrorCodeUniqueness(t *testing.T) {
	// 收集一些错误码进行唯一性测试
	errorCodes := map[int]string{
		errorsPkg.ValidationError.Code:              "ValidationError",
		errorsPkg.Unauthorized.Code:                 "Unauthorized",
		errorsPkg.PermissionDenied.Code:             "PermissionDenied",
		errorsPkg.UserNotFound.Code:                 "UserNotFound",
		errorsPkg.UserExists.Code:                   "UserExists",
		errorsPkg.DatabaseError.Code:                "DatabaseError",
		errorsPkg.DIServiceNotRegistered.Code:       "DIServiceNotRegistered",
		errorsPkg.InternalError.Code:                "InternalError",
	}
	
	// 验证错误码数量与映射数量一致（确保没有重复）
	assert.Equal(t, 8, len(errorCodes), "错误码应该是唯一的")
	
	// 验证特定的错误码
	assert.Equal(t, "ValidationError", errorCodes[1001])
	assert.Equal(t, "Unauthorized", errorCodes[1101])
	assert.Equal(t, "UserNotFound", errorCodes[1301])
	assert.Equal(t, "InternalError", errorCodes[9999])
}