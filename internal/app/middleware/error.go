package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/biuren/hebe/internal/app/render"
	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 应用级错误处理中间件
// 依赖RenderService处理错误响应
func ErrorHandler(renderService render.RenderServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 处理panic恢复
				stack := string(debug.Stack())

				// 记录panic信息
				slog.Error("panic recovered",
					"panic", r,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", stack,
				)

				// 创建内部错误响应
				err := errors.InternalErr()
				renderService.JSONError(c, err)
			}
		}()

		c.Next()
	}
}
