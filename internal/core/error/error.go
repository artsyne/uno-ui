package errors

import (
	"fmt"
	"net/http"

    "github.com/gin-gonic/gin"
)

// 定义所有的错误码
var (
    InputParamErr            = NewError(100000, "input_param_error")
    RegisterUserNameExists   = NewError(100001, "register_user_name_exists")
    RegisterUserEmailExists  = NewError(100002, "register_user_email_exists")
    ParseInputParamErr       = NewError(100003, "parse_input_param_error")
    UploadFileTypeErr        = NewError(100004, "upload_file_type_error")
    UploadFileTooLarge       = NewError(100005, "upload_file_too_large")
    UploadFileErr            = NewError(100006, "upload_file_error")
    UploadFileEmpty          = NewError(100007, "upload_file_empty")

    BizErr          = NewError(200000, "biz_error")
    UserNotExists   = NewError(200001, "user_not_exists")
    UserPwdErr      = NewError(200002, "user_pwd_error")
    ParseJsonErr    = NewError(200003, "parse_json_error")
    MaxDaybookLimitErr = NewError(200004, "max_daybook_limit_error")

    DbOperationErr  = NewError(300000, "db_operation_error")
    ExternalErr     = NewError(400000, "external_error")
    InternalErr     = NewError(500000, "internal_error")
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