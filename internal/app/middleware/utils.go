package middleware

import "github.com/gin-gonic/gin"

// getRequestID 辅助函数：获取请求ID
func getRequestID(c *gin.Context) string {
	reqID := c.GetHeader("X-Request-ID")
	if reqID == "" {
		if v, exists := c.Get("RequestId"); exists {
			reqID, _ = v.(string)
		}
	}
	return reqID
}
