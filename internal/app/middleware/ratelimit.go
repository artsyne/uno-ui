package middleware

import (
	"log/slog"
	"net/http"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
	"github.com/gin-gonic/gin"
)

// headerCapturer 用于捕获 tollbooth 设置的响应头
type headerCapturer struct {
	headers http.Header
}

func (w *headerCapturer) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

func (w *headerCapturer) Write([]byte) (int, error) { return 0, nil }
func (w *headerCapturer) WriteHeader(int)           {}

// RateLimit 限流中间件
func RateLimit(limit float64, burst int) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(limit, &limiter.ExpirableOptions{
		DefaultExpirationTTL: 0, // 默认不自动过期，依赖 tollbooth 内部管理
	})

	// 设置 IP 识别规则
	lmt.SetIPLookup(limiter.IPLookup{
		Name:           "X-Real-IP",
		IndexFromRight: 0,
	})

	return func(c *gin.Context) {
		// Pass headerCapturer to prevent tollbooth from writing response directly
		// but allow it to set headers (which we can copy)
		capturer := &headerCapturer{}
		httpErr := tollbooth.LimitByRequest(lmt, capturer, c.Request)
		if httpErr != nil {
			// Copy captured headers (like X-Rate-Limit-*)
			for k, v := range capturer.headers {
				for _, s := range v {
					c.Writer.Header().Add(k, s)
				}
			}

			c.Header("Retry-After", "60") // 建议 60 秒后重试

			slog.Warn("rate limit exceeded",
				"ip", c.ClientIP(),
				"path", c.Request.URL.Path,
				"request_id", getRequestID(c),
			)

			// 使用统一的错误响应
			resp := errors.NewAPIErrorResponse(
				errors.RateLimitExceeded.Code(),
				errors.RateLimitExceeded.Key(),
				"Rate limit exceeded. Please try again later.",
				getRequestID(c),
			)
			c.AbortWithStatusJSON(errors.RateLimitExceeded.HttpCode(), resp)
			return
		}
		c.Next()
	}
}

// AuthRateLimit 针对认证接口的严格限流
func AuthRateLimit(cfg *config.Config) gin.HandlerFunc {
	return RateLimit(cfg.RateLimit.AuthLimit, cfg.RateLimit.AuthBurst)
}

// GeneralRateLimit 全局基础限流
func GeneralRateLimit(cfg *config.Config) gin.HandlerFunc {
	return RateLimit(cfg.RateLimit.GeneralLimit, cfg.RateLimit.GeneralBurst)
}

// Helper to avoid circular dependency or duplication if getRequestID is not exported.
// Assuming getRequestID is available in the package or I should duplicate it/export it.
// In security.go it was using `getRequestID(c)`.
// `getRequestID` is likely in `common.go` or similar in the same package?
// Wait, I saw `getRequestID` usage in `security.go`.
// Let's check where `getRequestID` is defined.
