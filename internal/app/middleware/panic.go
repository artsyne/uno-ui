package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// PanicRecovery panic 恢复中间件
// 捕获并处理运行时 panic，记录详细的错误信息和堆栈跟踪
// 使用统一的错误响应格式
func PanicRecovery() gin.HandlerFunc {
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

			// 使用统一的错误响应格式
			internalErr := errors.InternalError
			response := gin.H{
				"code":    internalErr.Code,
				"message": internalErr.Key,
				"data":    nil,
			}

			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
		}
	})
}
