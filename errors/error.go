package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError 定义了一个应用错误的结构体
type AppError struct {
    // Code 是错误的唯一标识符
    Code int
    // MsgKey 是错误信息的键，可以用于支持多语言
    MsgKey  string
}


// NewError 创建一个新的 AppError。
func NewError(code int, msgKey string) *AppError {
    return  &AppError{
        Code: code,
        MsgKey: msgKey,
    }
}

// 实现error接口
// 返回错误信息
// todo：未来增加多语言支持，通过MsgKey获取对应的错误信息
func (e *AppError) Error() string {
    // 返回一个包含所有错误信息的字符串
    return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.MsgKey)
}

// 显示错误Json并返回
// todo：未来增加 log 记录
func ResponseErr(c *gin.Context, err error) {
    if appErr, ok := err.(*AppError); ok {
        // Handle the error as before.
        c.JSON(http.StatusBadRequest, gin.H{
            "Status":  appErr.Code,
            "ErrDesc": appErr.Error(),
        })
    } else {
        // Handle other types of errors.
        c.JSON(http.StatusBadRequest, gin.H{
            "Status":  http.StatusInternalServerError,
            "ErrDesc": err.Error(),
        })
    }
}