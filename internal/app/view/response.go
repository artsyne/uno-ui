package view

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/errors"
)

// APIResponse 标准化API响应结构
type APIResponse[T any] struct {
	Success   bool   `json:"success" example:"true"`
	Message   string `json:"message,omitempty" example:"操作成功"`
	Data      T      `json:"data,omitempty"`
	RequestID string `json:"request_id" example:"req_123456"`
}

// APIErrorResponse 标准化错误响应结构
type APIErrorResponse struct {
	Success   bool        `json:"success" example:"false"`
	Code      int         `json:"code" example:"4001001"`
	Message   string      `json:"message" example:"验证失败"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id"`
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	// 尝试从 header 获取
	reqID := c.GetHeader("X-Request-ID")
	if reqID == "" {
		// 尝试从 context 获取
		if v, exists := c.Get("RequestId"); exists {
			reqID, _ = v.(string)
		}
	}
	return reqID
}

// JSONSuccess 返回标准化的成功响应
func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
	s.JSONSuccessWithMsg(c, data, "success")
}

// JSONSuccessWithMsg 返回带消息的成功响应
func (s *ViewService) JSONSuccessWithMsg(c *gin.Context, data interface{}, msg string) {
	// 处理 HTMX 触发器
	if trigger := c.GetString("HX-Trigger"); trigger != "" {
		c.Header("HX-Trigger", trigger)
	}

	response := APIResponse[interface{}]{
		Success:   true,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
	}
	c.JSON(http.StatusOK, response)
}

// JSONError 返回标准化的错误响应
func (s *ViewService) JSONError(c *gin.Context, err error, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else if e, ok := err.(errors.AppError); ok {
		appErr = &e
	} else {
		// 如果是普通错误，转换为 InternalError，但保留原始错误信息作为详情
		e := errors.InternalError.WithDetails(err.Error())
		appErr = &e
	}

	s.logAppError(c, appErr)

	// 如果没有提供详细数据，且错误包含详情，则将详情放入Data中
	if responseData == nil && appErr.Details() != "" {
		responseData = gin.H{"error_details": appErr.Details()}
	}

	errorResponse := APIErrorResponse{
		Success:   false,
		Code:      appErr.Code(),
		Message:   appErr.Key(),
		Details:   responseData,
		RequestID: getRequestID(c),
	}
	c.JSON(appErr.HttpCode(), errorResponse)
}

// logAppError 记录应用错误日志
func (s *ViewService) logAppError(c *gin.Context, appErr *errors.AppError) {
	slog.Error("Application error",
		"request_id", getRequestID(c),
		"code", appErr.Code(),
		"key", appErr.Key(),
		"http_code", appErr.HttpCode(),
		"details", appErr.Details(),
	)
}

// JSON 直接返回指定HTTP状态码和数据的JSON响应 (保留以兼容旧代码，建议逐步迁移)
func (s *ViewService) JSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// HTML 渲染HTML模板
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
	c.HTML(http.StatusOK, name, data)
}

// HTMLError 渲染错误页面
func (s *ViewService) HTMLError(c *gin.Context, title, message string) {
	errorData := gin.H{
		"Title":   title,
		"Message": message,
	}
	c.HTML(http.StatusInternalServerError, "error", errorData)
}

// HTMLErrorDefault 渲染默认错误页面
func (s *ViewService) HTMLErrorDefault(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "error", nil)
}

// Redirect 重定向响应
func (s *ViewService) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}
