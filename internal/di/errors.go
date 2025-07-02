package di

import (
	"fmt"
	"reflect"
)

// DIError DI容器相关错误的基础接口
//
// 扩展了标准error接口，提供错误代码功能
// 便于错误分类和处理
type DIError interface {
	error
	ErrorCode() string
}

// baseDIError DI错误的基础实现
//
// 实现了DIError接口，提供统一的错误格式
type baseDIError struct {
	code    string // 错误代码，用于分类
	message string // 错误描述信息
	cause   error  // 原始错误（可选）
}

func (e *baseDIError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

func (e *baseDIError) ErrorCode() string {
	return e.code
}

func (e *baseDIError) Unwrap() error {
	return e.cause
}

// 简化的错误代码
const (
	ErrCodeRegistration = "DI001" // 注册相关错误
	ErrCodeResolution   = "DI002" // 解析相关错误
	ErrCodeLifecycle    = "DI003" // 生命周期相关错误
)

// NewAlreadyRegisteredError 创建已注册错误
func NewAlreadyRegisteredError(serviceType reflect.Type) DIError {
	return &baseDIError{
		code:    ErrCodeRegistration,
		message: fmt.Sprintf("service type %s is already registered", serviceType.String()),
	}
}

// NewInvalidFactoryError 创建无效工厂函数错误
func NewInvalidFactoryError(factoryType reflect.Type) DIError {
	return &baseDIError{
		code:    ErrCodeRegistration,
		message: fmt.Sprintf("invalid factory function type: %s", factoryType.String()),
	}
}

// NewNotRegisteredError 创建未注册错误
func NewNotRegisteredError(serviceType reflect.Type) DIError {
	return &baseDIError{
		code:    ErrCodeResolution,
		message: fmt.Sprintf("service type %s is not registered", serviceType.String()),
	}
}

// NewInitializationFailedError 创建初始化失败错误
func NewInitializationFailedError(cause error) DIError {
	return &baseDIError{
		code:    ErrCodeLifecycle,
		message: "container initialization failed",
		cause:   cause,
	}
}

// NewAlreadyInitializedError 创建已初始化错误
func NewAlreadyInitializedError() DIError {
	return &baseDIError{
		code:    ErrCodeLifecycle,
		message: "container is already initialized",
	}
}

// NewNotInitializedError 创建未初始化错误
func NewNotInitializedError() DIError {
	return &baseDIError{
		code:    ErrCodeLifecycle,
		message: "container is not initialized",
	}
}

// NewTypeError 创建类型转换错误
func NewTypeError(expectedType, actualType reflect.Type) DIError {
	return &baseDIError{
		code:    ErrCodeResolution,
		message: fmt.Sprintf("type conversion failed: expected %s, got %s", expectedType.String(), actualType.String()),
	}
}
