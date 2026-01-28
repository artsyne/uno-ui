package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		config         *config.Config
		origin         string
		expectedStatus int
		expectHeaders  bool
	}{
		{
			name: "Allowed Origin (Localhost Dev)",
			config: &config.Config{
				Env: "development",
				App: config.AppConfig{Domain: "example.com"},
			},
			origin:         "http://localhost:9090",
			expectedStatus: 204,
			expectHeaders:  true,
		},
		{
			name: "Allowed Origin (127.0.0.1 Dev)",
			config: &config.Config{
				Env: "development",
				App: config.AppConfig{Domain: "example.com"},
			},
			origin:         "http://127.0.0.1:8080",
			expectedStatus: 204,
			expectHeaders:  true,
		},
		{
			name: "Allowed Origin (Domain)",
			config: &config.Config{
				Env: "production",
				App: config.AppConfig{Domain: "example.com"},
			},
			origin:         "https://example.com",
			expectedStatus: 204,
			expectHeaders:  true,
		},
		{
			name: "Disallowed Origin (Production Localhost)",
			config: &config.Config{
				Env: "production",
				App: config.AppConfig{Domain: "example.com"},
			},
			origin:         "http://localhost:9090",
			expectedStatus: 204, // Still returns 204 for OPTIONS, but NO headers
			expectHeaders:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.CORS(tt.config))
			router.GET("/", func(c *gin.Context) {
				c.String(200, "OK")
			})

			// Test OPTIONS request (Preflight)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("OPTIONS", "/", nil)
			req.Header.Set("Origin", tt.origin)
			req.Header.Set("Access-Control-Request-Method", "POST")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectHeaders {
				assert.Equal(t, tt.origin, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}
