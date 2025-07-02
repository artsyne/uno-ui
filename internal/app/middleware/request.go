package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestID 请求ID中间件
// 为每个请求生成唯一ID，便于日志追踪
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头是否已有RequestID
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			// 生成新的RequestID
			requestID = generateRequestID()
		}

		// 设置到上下文
		c.Set("RequestID", requestID)

		// 设置响应头
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		// 如果生成失败，使用时间戳作为备选
		return "fallback"
	}
	return hex.EncodeToString(bytes)
}
