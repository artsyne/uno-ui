package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if recovered != nil {
			// 记录 panic 详情
			log.Printf("系统发生panic: %v, 路径: %s, 方法: %s, 客户端IP: %s",
				recovered, c.Request.URL.Path, c.Request.Method, c.ClientIP())
			log.Printf("堆栈跟踪: %s", string(debug.Stack()))

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

// logAppError 记录应用错误的详细信息
func logAppError(c *gin.Context, err *errors.AppError) {
	log.Printf("API错误: 错误码=%d, 错误键=%s, HTTP状态=%d, 方法=%s, 路径=%s, 客户端IP=%s, 详情=%s",
		err.Code, err.MsgKey, err.HttpCode, c.Request.Method, c.Request.URL.Path,
		c.ClientIP(), err.Details)
}

// GlobalErrorHandlerMiddleware 全局错误处理中间件
func GlobalErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
			if len(c.Errors) > 0 {
				// 获取最后一个错误
				lastError := c.Errors.Last()

				// 处理 AppError
				if appErr, ok := lastError.Err.(*errors.AppError); ok {
					logAppError(c, appErr)
					c.JSON(appErr.HttpCode, gin.H{
						"code":    appErr.Code,
						"message": appErr.MsgKey,
						"data":    nil,
					})
					return
				}

				// 处理其他错误
				log.Printf("未知错误: %s, 路径: %s, 方法: %s, 客户端IP: %s",
					lastError.Err.Error(), c.Request.URL.Path, c.Request.Method, c.ClientIP())

				// 返回错误响应
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
