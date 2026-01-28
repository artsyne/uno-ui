package middleware

import (
	"log/slog"
	"net/http"

	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

// AppErrorHandler 应用错误处理中间件
// 处理通过 c.Error() 添加的业务错误和应用错误
// 使用统一的错误响应格式
func AppErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) > 0 {
				lastError := c.Errors.Last()

				// 处理 AppError 类型的业务错误
				if appErr, ok := lastError.Err.(*errors.AppError); ok {
					logAppError(c, appErr)
					c.JSON(appErr.HttpCode(), errors.NewAPIErrorResponse(appErr.Code(), appErr.Key(), appErr.Details(), getRequestID(c)))
					return
				}

				// 处理其他类型的错误
				slog.Error("未知应用错误",
					"error", lastError.Err.Error(),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"user_agent", c.Request.UserAgent(),
					"client_ip", c.ClientIP())

				internalErr := errors.InternalError
				c.JSON(http.StatusInternalServerError, errors.NewAPIErrorResponse(internalErr.Code(), internalErr.Key(), nil, getRequestID(c)))
			}
		}()

		c.Next()
	}
}

// logAppError 记录应用错误的详细信息
// 实现应用特定的错误日志格式
func logAppError(c *gin.Context, err *errors.AppError) {
	slog.Error("应用业务错误",
		"error_code", err.Code(),
		"error_key", err.Key(),
		"http_code", err.HttpCode(),
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"user_agent", c.Request.UserAgent(),
		"client_ip", c.ClientIP(),
		"details", err.Details())
}
