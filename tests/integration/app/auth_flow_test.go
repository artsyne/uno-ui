package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginFlow_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Config
	cfg := &config.Config{
		Env: "development",
		App: config.AppConfig{
			Domain: "example.com",
		},
		Cookie: config.CookieConfig{
			Secret: "12345678901234567890123456789012", // 32 bytes
		},
	}

	// Setup Router
	router := gin.New()
	// Apply Middleware in order
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.CSRF(cfg))

	// Mock Login Handler
	router.GET("/login", func(c *gin.Context) {
		// Just to get the token
		token := c.GetString("CsrfToken")
		c.String(200, token)
	})

	router.POST("/login", func(c *gin.Context) {
		c.String(200, "Login Success")
	})

	t.Run("Full Login Flow (Same Origin)", func(t *testing.T) {
		// 1. GET /login to establish session and get CSRF token
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/login", nil)
		router.ServeHTTP(w1, req1)

		assert.Equal(t, 200, w1.Code)
		token := w1.Body.String()
		assert.NotEmpty(t, token)

		// Get the CSRF cookie
		var csrfCookie *http.Cookie
		for _, c := range w1.Result().Cookies() {
			if c.Name == "idiary_csrf" {
				csrfCookie = c
				break
			}
		}
		assert.NotNil(t, csrfCookie, "CSRF cookie not found")

		// 2. POST /login with token
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/login", nil)
		// Set Host and Origin to simulate same-origin
		req2.Host = "example.com"
		req2.URL.Scheme = "http"
		req2.URL.Host = "example.com"
		req2.Header.Set("Origin", "http://example.com")
		
		req2.Header.Set("X-CSRF-Token", token)
		req2.AddCookie(csrfCookie)
		
		router.ServeHTTP(w2, req2)
		assert.Equal(t, 200, w2.Code)
		assert.Equal(t, "Login Success", w2.Body.String())
	})

	t.Run("Cross Origin Login Flow (Allowed Origin)", func(t *testing.T) {
		// 1. GET /login (Simulate visiting the page)
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/login", nil)
		router.ServeHTTP(w1, req1)
		
		token := w1.Body.String()
		var csrfCookie *http.Cookie
		for _, c := range w1.Result().Cookies() {
			if c.Name == "idiary_csrf" {
				csrfCookie = c
				break
			}
		}

		// 2. POST /login from Localhost (Cross Origin)
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/login", nil)
		// Set Host to target domain, Origin to client domain
		req2.Host = "example.com"
		req2.URL.Scheme = "http"
		req2.URL.Host = "example.com"
		req2.Header.Set("Origin", "http://localhost:9090")
		
		req2.Header.Set("X-CSRF-Token", token)
		req2.AddCookie(csrfCookie)

		router.ServeHTTP(w2, req2)
		
		// Should succeed because CORS allows localhost in dev, and TrustedOrigins includes localhost
		assert.Equal(t, 200, w2.Code)
		assert.Equal(t, "Login Success", w2.Body.String())
		assert.Equal(t, "http://localhost:9090", w2.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Cross Origin Login Flow (Disallowed Origin)", func(t *testing.T) {
		// 1. GET /login
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/login", nil)
		router.ServeHTTP(w1, req1)
		
		token := w1.Body.String()
		var csrfCookie *http.Cookie
		for _, c := range w1.Result().Cookies() {
			if c.Name == "idiary_csrf" {
				csrfCookie = c
				break
			}
		}

		// 2. POST /login from Evil Site
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/login", nil)
		req2.Header.Set("X-CSRF-Token", token)
		req2.Header.Set("Origin", "http://evil.com")
		req2.AddCookie(csrfCookie)

		router.ServeHTTP(w2, req2)
		
		// CORS might NOT block the request execution in Gin unless using a specific middleware that aborts.
		// My CORS middleware aborts OPTIONS, but for POST it just sets headers.
		// However, CSRF middleware CHECKS the Origin/Referer.
		// "http://evil.com" is NOT in TrustedOrigins.
		// So CSRF middleware should block it with 403.
		
		assert.Equal(t, 403, w2.Code)
		// Access-Control-Allow-Origin should NOT be set for evil.com
		assert.Empty(t, w2.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Preflight Request (OPTIONS)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/login", nil)
		req.Header.Set("Origin", "http://localhost:9090")
		req.Header.Set("Access-Control-Request-Method", "POST")
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, 204, w.Code)
		assert.Equal(t, "http://localhost:9090", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})
}
