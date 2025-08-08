package view

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "github.com/biuren/hebe/internal/utils/error"
)

// JSONResponse JSON响应结构
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON 返回JSON响应
func (s *ViewService) JSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// JSONSuccess 返回成功JSON响应
func (s *ViewService) JSONSuccess(c *gin.Context, data interface{}) {
	response := JSONResponse{
		Success: true,
		Message: "success",
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// JSONError 返回错误JSON响应
func (s *ViewService) JSONError(c *gin.Context, err error, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	if appErr, ok := err.(*appErrors.AppError); ok {
		s.logAppError(appErr)
		errorResponse := ErrorResponse{
			Success: false,
			Message: appErr.Key,
			Code:    fmt.Sprintf("%d", appErr.Code),
			Data:    responseData,
		}
		c.JSON(appErr.HttpCode, errorResponse)
	} else {
		slog.Error("Unexpected error", "error", err.Error())
		errorResponse := ErrorResponse{
			Success: false,
			Message: "Internal server error",
			Code:    "INTERNAL_ERROR",
			Data:    responseData,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
	}
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
	c.HTML(http.StatusInternalServerError, "error.html", errorData)
}

// Redirect 重定向响应
func (s *ViewService) Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}

// logAppError 记录应用错误日志
func (s *ViewService) logAppError(appErr *appErrors.AppError) {
	slog.Error("Application error",
		"code", appErr.Code,
		"key", appErr.Key,
		"httpCode", appErr.HttpCode,
		"details", appErr.Details,
	)
}