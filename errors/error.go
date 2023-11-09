package errors

import (
	"fmt"
)

// AppError 定义了一个应用错误的结构体
type AppError struct {
    // Code 是错误的唯一标识符
    Code int
    // MsgKey 是错误信息的键，可以用于支持多语言
    MsgKey  string
    // Data 是与错误相关的额外数据
    Data interface{}
    // Source 是引发错误的源头，例如函数或方法的名称
    Source string
}

// ErrorOption 定义了一个修改 AppError 的函数的类型
type ErrorOption func(*AppError)

// WithData 返回一个设置 Data 字段的 ErrorOption
func WithData(data interface{}) ErrorOption {
    return func(e *AppError) {
        e.Data = data
    }
}

// WithSource 返回一个设置 Source 字段的 ErrorOption
func WithSource(source string) ErrorOption {
    return func(e *AppError) {
        e.Source = source
    }
}


// NewError 创建一个新的 AppError。你可以使用 ErrorOption 来设置可选的字段。
//
// 例如，创建一个没有任何可选字段的错误：
//     err := NewError(100000, "input_param_error")
//
// 创建一个带有 Data 字段的错误：
//     err := NewError(100000, "input_param_error", WithData(someData))
//
func NewError(code int, msgKey string, opts ...ErrorOption) *AppError {
    e := &AppError{
        Code: code,
        MsgKey: msgKey,
    }
    for _, opt := range opts {
        opt(e)
    }
    return e
}

// 实现error接口
// 返回错误信息
// todo：未来增加多语言支持，通过MsgKey获取对应的错误信息
func (e *AppError) Error() string {
    // 返回一个包含所有错误信息的字符串
    return fmt.Sprintf("Code: %d, Message: %s,  Source: %s", e.Code, e.MsgKey, e.Source)
}
