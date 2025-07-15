package render

import (
	"log/slog"
	"net/http"

	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// ResponseHandler 响应处理器
type ResponseHandler struct {
	service *RenderService
}

// NewResponseHandler 创建响应处理器
func NewResponseHandler(service *RenderService) *ResponseHandler {
	return &ResponseHandler{
		service: service,
	}
}

// JSON 返回JSON响应
func (rh *ResponseHandler) JSON(c *gin.Context, code int, data interface{}) {
	response := JSONResponse{
		Code:    code,
		Message: "success",
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// JSONSuccess 返回成功响应
func (rh *ResponseHandler) JSONSuccess(c *gin.Context, data interface{}) {
	rh.JSON(c, 0, data)
}

// JSONError 返回错误响应
func (rh *ResponseHandler) JSONError(c *gin.Context, err error, data ...interface{}) {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}

	// 处理AppError
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.HttpCode, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Key,
			Data:    errorData,
		})

		// 记录错误的详细信息
		rh.logAppError(c, appErr)
		return
	}

	// 处理其他类型的错误，转换为内部错误
	internalErr := errors.InternalError.Err()
	response := ErrorResponse{
		Code:    internalErr.Code,
		Message: internalErr.Key,
		Data:    errorData,
	}
	c.JSON(internalErr.HttpCode, response)
}

// HTML 返回HTML页面
func (rh *ResponseHandler) HTML(c *gin.Context, name string, data interface{}) {
	c.HTML(http.StatusOK, name, data)
}

// HTMLError 渲染错误页面
func (rh *ResponseHandler) HTMLError(c *gin.Context, title, message string) {
	pageData := &PageData{
		Title: title,
		Data: map[string]interface{}{
			"error_title":   title,
			"error_message": message,
		},
	}
	c.HTML(http.StatusOK, "error.html", pageData)
}

// Redirect 重定向
func (rh *ResponseHandler) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}

// logAppError 记录应用错误的详细信息
func (rh *ResponseHandler) logAppError(c *gin.Context, err *errors.AppError) {
	// 使用slog记录错误
	slog.Error("API错误",
		"error_code", err.Code,
		"error_key", err.Key,
		"http_code", err.HttpCode,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"user_agent", c.Request.UserAgent(),
		"client_ip", c.ClientIP(),
		"details", err.Details,
	)
}
