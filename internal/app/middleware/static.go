package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// StaticCache 返回静态资源缓存中间件
// 根据文件类型设置不同的缓存时间
//
// 缓存策略:
//   - CSS/JS: 缓存 1 年 (31536000 秒)
//   - 图片(PNG/JPG/GIF/SVG/ICO): 缓存 1 周 (604800 秒)
func StaticCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js") {
			c.Header("Cache-Control", "public, max-age=31536000")
		} else if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") ||
			strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".gif") ||
			strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".ico") {
			c.Header("Cache-Control", "public, max-age=604800")
		}

		c.Next()
	}
}
