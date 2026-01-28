package middleware

import (
	"net/http"
	"strings"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
)

// CORS Cross-Origin Resource Sharing middleware
// 处理跨域请求，允许特定的源访问API
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// Check if origin is allowed
			allowed := false
			
			// Always allow the configured domain
			if origin == "http://"+cfg.App.Domain || origin == "https://"+cfg.App.Domain {
				allowed = true
			}

			// In non-production environments, allow localhost and 127.0.0.1 with any port
			if cfg.Env != "production" {
				if strings.HasPrefix(origin, "http://localhost") || 
				   strings.HasPrefix(origin, "https://localhost") || 
				   strings.HasPrefix(origin, "http://127.0.0.1") || 
				   strings.HasPrefix(origin, "https://127.0.0.1") {
					allowed = true
				}
			}

			if allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, HX-Request, HX-Target, HX-Current-URL")
				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
