package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/utils/config"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := framework.NewTestContext("TestSecurityMiddleware")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("RateLimit", func(t *testing.T) {
		// 测试限流: 1 req/sec
		cfg := &config.Config{
			RateLimit: config.RateLimitConfig{
				GeneralLimit: 1,
				GeneralBurst: 1,
			},
		}

		router := gin.New()
		router.Use(middleware.GeneralRateLimit(cfg))
		router.GET("/", func(c *gin.Context) {
			c.String(200, "OK")
		})

		// 第一次请求 - 应该成功
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/", nil)
		req1.Header.Set("X-Real-IP", "127.0.0.1")
		router.ServeHTTP(w1, req1)
		assert.Equal(t, 200, w1.Code)

		// 第二次请求 (立即) - 应该被限流
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/", nil)
		req2.Header.Set("X-Real-IP", "127.0.0.1")
		router.ServeHTTP(w2, req2)
		assert.Equal(t, 429, w2.Code)
		assert.Contains(t, w2.Body.String(), "rate_limit_exceeded")
		assert.Equal(t, "60", w2.Header().Get("Retry-After"))
	})

	t.Run("CSRF", func(t *testing.T) {
		cfg := &config.Config{
			Env: "production", // Try production to enable Secure=true
			App: config.AppConfig{
				Domain: "example.com",
			},
			Cookie: config.CookieConfig{
				Secret: "12345678901234567890123456789012", // 32 bytes
			},
		}

		router := gin.New()
		router.Use(middleware.CSRF(cfg))
		
		router.GET("/", func(c *gin.Context) {
			token := c.GetString("CsrfToken")
			c.String(200, token)
		})
		
		router.POST("/", func(c *gin.Context) {
			c.String(200, "Posted")
		})

		// 1. GET 请求应该生成 Token
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w1, req1)
		
		assert.Equal(t, 200, w1.Code)
		token := w1.Body.String()
		assert.NotEmpty(t, token)
		
		// 获取 Cookie
		cookies := w1.Result().Cookies()
		var csrfCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "idiary_csrf" { // 自定义 cookie 名
				csrfCookie = c
				break
			}
		}
		
		assert.NotNil(t, csrfCookie, "CSRF cookie should be present. Cookies: %v", cookies)
        if csrfCookie == nil {
            return
        }

		// 2. POST 请求如果不带 Token 应该失败
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/", nil)
		// 带上 cookie 但不带 header
		req2.AddCookie(csrfCookie)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, 403, w2.Code)

		// 3. POST 请求带正确 Token 应该成功
		w3 := httptest.NewRecorder()
		req3 := httptest.NewRequest("POST", "https://example.com/", nil)
		req3.AddCookie(csrfCookie)
		req3.Header.Set("X-CSRF-Token", token)
		req3.Header.Set("Origin", "https://example.com")
		req3.Header.Set("Referer", "https://example.com/")
		router.ServeHTTP(w3, req3)
		assert.Equal(t, 200, w3.Code, "Response body: %s. Token: %s, Cookie: %v", w3.Body.String(), token, csrfCookie)
	})
}
