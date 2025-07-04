package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// PanicRecoveryMiddleware panic 恢复中间件
// 捕获并处理运行时 panic，记录详细的错误信息和堆栈跟踪
func PanicRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if recovered != nil {
			// 记录 panic 详细信息
			slog.Error("panic recovered",
				"panic", recovered,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"stack", string(debug.Stack()))

			response := gin.H{
				"code":    errors.InternalErr().Code,
				"message": errors.InternalErr().MsgKey,
				"data":    nil,
			}

			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
		}
	})
}

// AppErrorHandlerMiddleware 应用错误处理中间件
// 处理通过 c.Error() 添加的业务错误和应用错误
func AppErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) > 0 {
				lastError := c.Errors.Last()

				// 处理 AppError 类型的业务错误
				if appErr, ok := lastError.Err.(*errors.AppError); ok {
					logAppError(c, appErr)
					c.JSON(appErr.HttpCode, gin.H{
						"code":    appErr.Code,
						"message": appErr.MsgKey,
						"data":    nil,
					})
					return
				}

				// 处理其他类型的错误
				slog.Error("未知应用错误",
					"error", lastError.Err.Error(),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"user_agent", c.Request.UserAgent(),
					"client_ip", c.ClientIP())

				internalErr := errors.InternalErr()
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    internalErr.Code,
					"message": internalErr.MsgKey,
					"data":    nil,
				})
			}
		}()

		c.Next()
	}
}

// logAppError 记录应用错误的详细信息
func logAppError(c *gin.Context, err *errors.AppError) {
	slog.Error("应用业务错误",
		"error_code", err.Code,
		"error_key", err.MsgKey,
		"http_code", err.HttpCode,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"user_agent", c.Request.UserAgent(),
		"client_ip", c.ClientIP(),
		"details", err.Details)
}
