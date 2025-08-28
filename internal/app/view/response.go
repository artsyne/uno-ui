package view

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/errors"
)

// JSONResponse JSON响应结构
// 用于标准化API成功响应的数据结构
type JSONResponse struct {
	Success bool        `json:"success"` // 请求是否成功
	Message string      `json:"message"` // 响应消息
	Code    int         `json:"code"`    // 业务状态码
	Data    interface{} `json:"data"`    // 响应数据
}

// ErrorResponse 错误响应结构
// 用于标准化API错误响应的数据结构
type ErrorResponse struct {
	Success bool        `json:"success"`          // 请求是否成功，错误时为false
	Message string      `json:"message"`          // 错误消息
	Code    int         `json:"code"`             // 错误代码
	Data    interface{} `json:"data,omitempty"`   // 可选的错误详细数据
}

// JSON 返回JSON响应
// 直接返回指定HTTP状态码和数据的JSON响应
//
// 参数:
//   - c: Gin上下文对象
//   - code: HTTP状态码
//   - data: 响应数据
func (s *ViewService) JSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// JSONSuccess 返回成功JSON响应
// 返回标准化的成功响应，HTTP状态码为200
//
// 参数:
//   - c: Gin上下文对象
//   - data: 成功响应的数据
func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
	response := JSONResponse{
		Success: true,
		Message: "success",
		Code:    0,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// JSONError 返回错误JSON响应
// 处理应用错误并返回标准化的错误响应
//
// 参数:
//   - c: Gin上下文对象
//   - err: 错误对象，应为*appErrors.AppError类型
//   - data: 可选的错误详细数据
func (s *ViewService) JSONError(c *gin.Context, err error, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	if appErr, ok := err.(*errors.AppError); ok {
		s.logAppError(appErr)
		errorResponse := ErrorResponse{
			Success: false,
			Message: appErr.Key(),
			Code:    appErr.Code(),
			Data:    responseData,
		}
		c.JSON(appErr.HttpCode(), errorResponse)
	}
}

// HTML 渲染HTML模板
// 使用指定的模板名称和数据渲染HTML页面
//
// 参数:
//   - c: Gin上下文对象
//   - name: 模板文件名
//   - data: 传递给模板的数据
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
	c.HTML(http.StatusOK, name, data)
}

// HTMLError 渲染错误页面
// 渲染通用错误页面，HTTP状态码为500
//
// 参数:
//   - c: Gin上下文对象
//   - title: 错误页面标题
//   - message: 错误消息内容
func (s *ViewService) HTMLError(c *gin.Context, title, message string) {
	errorData := gin.H{
		"Title":   title,
		"Message": message,
	}
	c.HTML(http.StatusInternalServerError, "error.html", errorData)
}

// Redirect 重定向响应
// 执行HTTP重定向到指定URL，状态码为302
//
// 参数:
//   - c: Gin上下文对象
//   - url: 重定向目标URL
func (s *ViewService) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}

// logAppError 记录应用错误日志
// 使用结构化日志记录应用错误的详细信息
//
// 参数:
//   - appErr: 应用错误对象，包含错误代码、键名、HTTP状态码等信息
func (s *ViewService) logAppError(appErr *errors.AppError) {
	slog.Error("Application error",
		"code", appErr.Code,
		"key", appErr.Key,
		"httpCode", appErr.HttpCode,
		"details", appErr.Details,
	)
}