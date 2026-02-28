package view

import (
	stdErrors "errors"
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

// APIErrorResponse 标准化错误响应结构 (Deprecated: use errors.APIErrorResponse instead)
type APIErrorResponse = errors.APIErrorResponse

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
	if responseData == nil && appErr.Details() != nil {
		if str, ok := appErr.Details().(string); ok && str != "" {
			responseData = gin.H{"error_details": str}
		} else if appErr.Details() != nil {
			responseData = appErr.Details()
		}
	}

	errorResponse := errors.NewAPIErrorResponse(appErr.Code(), appErr.Key(), responseData, getRequestID(c))
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

// HTML 渲染HTML模板
func (s *ViewService) HTML(c *gin.Context, name string, data interface{}) {
	c.HTML(http.StatusOK, name, data)
}

// HTMLError 渲染错误页面
func (s *ViewService) HTMLError(c *gin.Context, title, message string) {
	pageData := s.CreatePageData(c, "error.title")
	pageData.Data["Title"] = title
	pageData.Data["Message"] = message
	c.HTML(http.StatusInternalServerError, "error", pageData)
}

// RenderErrorPage 渲染统一的错误页面
// 当发生错误时显示用户友好的错误信息，统一错误页面的展示格式
//
// 参数:
//   - c: Gin上下文对象
//   - err: 错误对象（优先处理 AppError）
func (s *ViewService) RenderErrorPage(c *gin.Context, err error) {
	lang := s.GetLocaleFromContext(c)

	var appErr errors.AppError
	if !stdErrors.As(err, &appErr) {
		// 如果不是 AppError，回退到 InternalError
		slog.Warn("Error is not an AppError, falling back to InternalError", "error", err)
		appErr = errors.InternalError
	}

	title := s.GetMessage("error.title", lang)
	messageKey := "error." + appErr.Key()
	message := s.GetMessage(messageKey, lang)

	// Fallback if translation missing
	if message == messageKey {
		// 尝试只用 Key 查找，有些 Key 可能没有 error. 前缀
		message = s.GetMessage(appErr.Key(), lang)
		if message == appErr.Key() {
			slog.Warn("Translation missing, using default message", "key", messageKey)
			message = s.GetMessage("error.message", lang)
		}
	}

	// Use "error.title" to avoid warning about missing "error" key
	data := s.CreatePageData(c, "error.title")
	// Template uses .Data.Title and .Data.Message (Capitalized keys)
	data.Data["Title"] = title
	data.Data["Message"] = message

	c.HTML(appErr.HttpCode(), "error", data)
}

// HTMLErrorDefault 渲染默认错误页面
func (s *ViewService) HTMLErrorDefault(c *gin.Context) {
	pageData := s.CreatePageData(c, "error.title")
	c.HTML(http.StatusInternalServerError, "error", pageData)
}

// Redirect 重定向响应
func (s *ViewService) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}
