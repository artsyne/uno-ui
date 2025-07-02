package middleware

import (
	"github.com/gin-gonic/gin"
)

// ResponseHeaders 响应头处理中间件
// 设置通用的响应头
func ResponseHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置CORS头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		// 设置安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// 设置缓存控制（对于API请求）
		if isAPIRequest(c) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isAPIRequest 判断是否为API请求
func isAPIRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	// 简单判断：以/api/开头的路径为API请求
	return len(path) >= 5 && path[:5] == "/api/"
}
