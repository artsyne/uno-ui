package error_test

import (
	"net/http"
	"testing"

	errors "github.com/biuren/hebe/internal/core/error"
)

// TestAppError_BasicFunctionality 测试AppError基础功能
func TestAppError_BasicFunctionality(t *testing.T) {
	// 测试现代化的创建方式
	err := errors.NewAppError(errors.CodeValidationError, "validation_error")
	if err.Code != errors.CodeValidationError {
		t.Errorf("Expected code %d, got %d", errors.CodeValidationError, err.Code)
	}
	if err.MsgKey != "validation_error" {
		t.Errorf("Expected MsgKey 'validation_error', got %s", err.MsgKey)
	}
	if err.HttpCode != http.StatusBadRequest {
		t.Errorf("Expected HttpCode %d, got %d", http.StatusBadRequest, err.HttpCode)
	}
}

// TestAppError_WithHttpCode 测试指定HTTP状态码的创建方式
func TestAppError_WithHttpCode(t *testing.T) {
	// 测试指定HTTP状态码的创建方式
	err := errors.NewAppErrorWithHttpCode(errors.CodeUserExists, "user_exists", http.StatusConflict)
	if err.Code != errors.CodeUserExists {
		t.Errorf("Expected code %d, got %d", errors.CodeUserExists, err.Code)
	}
	if err.MsgKey != "user_exists" {
		t.Errorf("Expected MsgKey 'user_exists', got %s", err.MsgKey)
	}
	if err.HttpCode != http.StatusConflict {
		t.Errorf("Expected HttpCode %d, got %d", http.StatusConflict, err.HttpCode)
	}
}

// TestAppError_ChainMethods 测试链式方法调用
func TestAppError_ChainMethods(t *testing.T) {
	err := errors.NewAppError(errors.CodeUserExists, "user_exists").
		WithMeta(map[string]interface{}{"username": "testuser"}).
		WithRequestID("req-123").
		WithUserID(456).
		WithContext("user", "register").
		WithDetails("Username already exists in database")

	// 验证链式调用结果
	if err.RequestID != "req-123" {
		t.Errorf("Expected RequestID 'req-123', got %s", err.RequestID)
	}
	if err.UserID != 456 {
		t.Errorf("Expected UserID 456, got %v", err.UserID)
	}
	if err.Module != "user" {
		t.Errorf("Expected Module 'user', got %s", err.Module)
	}
	if err.Action != "register" {
		t.Errorf("Expected Action 'register', got %s", err.Action)
	}
	if err.Details != "Username already exists in database" {
		t.Errorf("Expected Details match, got %s", err.Details)
	}
	if username, exists := err.Metadata["username"]; !exists || username != "testuser" {
		t.Errorf("Expected Metadata username 'testuser', got %v", err.Metadata["username"])
	}
}

// TestAppError_ErrorTypes 测试错误类型分类
func TestAppError_ErrorTypes(t *testing.T) {
	tests := []struct {
		code     int
		expected errors.ErrorType
	}{
		{errors.CodeValidationError, errors.ErrorTypeValidation},
		{errors.CodeUserExists, errors.ErrorTypeBusiness},
		{errors.CodeDatabaseError, errors.ErrorTypeDatabase},
		{errors.CodeNotFound, errors.ErrorTypeNotFound},
		{errors.CodeUnauthorized, errors.ErrorTypePermission},
	}

	for _, test := range tests {
		err := errors.NewAppError(test.code, "test_error")
		if err.GetType() != test.expected {
			t.Errorf("For code %d, expected type %d, got %d", test.code, test.expected, err.GetType())
		}
	}
}

// TestAppError_ErrorSeverity 测试错误严重程度
func TestAppError_ErrorSeverity(t *testing.T) {
	tests := []struct {
		code     int
		expected errors.ErrorSeverity
	}{
		{errors.CodeValidationError, errors.SeverityLow},
		{errors.CodeUserExists, errors.SeverityMedium},
		{errors.CodeDatabaseError, errors.SeverityHigh},
		{errors.CodeInternalError, errors.SeverityHigh},
	}

	for _, test := range tests {
		err := errors.NewAppError(test.code, "test_error")
		if err.GetSeverity() != test.expected {
			t.Errorf("For code %d, expected severity %d, got %d", test.code, test.expected, err.GetSeverity())
		}
	}
}

// TestAppError_HttpStatusMapping 测试HTTP状态码映射
func TestAppError_HttpStatusMapping(t *testing.T) {
	tests := []struct {
		code     int
		expected int
	}{
		{errors.CodeValidationError, http.StatusBadRequest},
		{errors.CodeUnauthorized, http.StatusUnauthorized},
		{errors.CodePermissionDenied, http.StatusForbidden},
		{errors.CodeNotFound, http.StatusNotFound},
		{errors.CodeUserExists, http.StatusConflict},
		{errors.CodeDatabaseError, http.StatusInternalServerError},
	}

	for _, test := range tests {
		err := errors.NewAppError(test.code, "test_error")
		if err.HttpCode != test.expected {
			t.Errorf("For code %d, expected HTTP status %d, got %d", test.code, test.expected, err.HttpCode)
		}
	}
}

// TestAppError_ErrorInterface 测试error接口实现
func TestAppError_ErrorInterface(t *testing.T) {
	err := errors.NewAppError(errors.CodeUserExists, "user_exists")

	// 测试默认Error()输出
	errorStr := err.Error()
	expected := "Code: 2001, Message: user_exists"
	if errorStr != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, errorStr)
	}

	// 测试带Details的Error()输出
	err.WithDetails("Custom error details")
	errorStr = err.Error()
	if errorStr != "Custom error details" {
		t.Errorf("Expected error string 'Custom error details', got '%s'", errorStr)
	}
}

// TestAppError_PredefinedErrors 测试预定义错误
func TestAppError_PredefinedErrors(t *testing.T) {
	// 测试预定义的错误变量
	if errors.ErrValidation.Code != errors.CodeValidationError {
		t.Errorf("Expected ErrValidation code %d, got %d", errors.CodeValidationError, errors.ErrValidation.Code)
	}
	if errors.ErrUserNotFound.Code != errors.CodeUserNotFound {
		t.Errorf("Expected ErrUserNotFound code %d, got %d", errors.CodeUserNotFound, errors.ErrUserNotFound.Code)
	}
	if errors.ErrDatabase.Code != errors.CodeDatabaseError {
		t.Errorf("Expected ErrDatabase code %d, got %d", errors.CodeDatabaseError, errors.ErrDatabase.Code)
	}
}

// TestAppError_ModernErrorCodes 测试现代化的4位错误码
func TestAppError_ModernErrorCodes(t *testing.T) {
	// 验证错误码在正确的范围内
	tests := []struct {
		name     string
		code     int
		minRange int
		maxRange int
	}{
		{"通用错误", errors.CodeValidationError, 1000, 1999},
		{"用户模块", errors.CodeUserExists, 2000, 2999},
		{"日记模块", errors.CodeDiaryNotFound, 3000, 3999},
		{"文章模块", errors.CodePostNotFound, 4000, 4999},
		{"任务模块", errors.CodeTaskNotFound, 5000, 5999},
		{"系统配置", errors.CodeConfigError, 6000, 6999},
	}

	for _, test := range tests {
		if test.code < test.minRange || test.code > test.maxRange {
			t.Errorf("%s错误码 %d 不在范围 %d-%d 内", test.name, test.code, test.minRange, test.maxRange)
		}
	}
}

// TestAppError_PrebuiltInstances 测试预构建的错误实例
func TestAppError_PrebuiltInstances(t *testing.T) {
	// 测试预定义错误实例的基本属性
	tests := []struct {
		err      *errors.AppError
		code     int
		msgKey   string
		httpCode int
	}{
		{errors.ErrValidation, errors.CodeValidationError, "validation_error", http.StatusBadRequest},
		{errors.ErrUserExists, errors.CodeUserExists, "user_exists", http.StatusConflict},
		{errors.ErrNotFound, errors.CodeNotFound, "not_found", http.StatusNotFound},
		{errors.ErrUnauthorized, errors.CodeUnauthorized, "unauthorized", http.StatusUnauthorized},
	}

	for _, test := range tests {
		if test.err.Code != test.code {
			t.Errorf("Expected code %d, got %d", test.code, test.err.Code)
		}
		if test.err.MsgKey != test.msgKey {
			t.Errorf("Expected msgKey %s, got %s", test.msgKey, test.err.MsgKey)
		}
		if test.err.HttpCode != test.httpCode {
			t.Errorf("Expected httpCode %d, got %d", test.httpCode, test.err.HttpCode)
		}
	}
}
