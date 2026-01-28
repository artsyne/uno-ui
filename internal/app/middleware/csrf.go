package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRF CSRF中间件
func CSRF(cfg *config.Config) gin.HandlerFunc {
	// 使用 Cookie 中的 Secret 作为 CSRF 的加密密钥
	// gorilla/csrf 要求密钥长度为 32 位
	secret := []byte(cfg.Cookie.Secret)
	if len(secret) < 32 {
		// 如果密钥长度不足 32 位，则通过填充使其达到 32 位
		newSecret := make([]byte, 32)
		copy(newSecret, secret)
		secret = newSecret
	}

	// 处理域名：移除 http:// 或 https:// 前缀
	domain := cfg.App.Domain
	if strings.HasPrefix(domain, "http://") {
		domain = strings.TrimPrefix(domain, "http://")
	} else if strings.HasPrefix(domain, "https://") {
		domain = strings.TrimPrefix(domain, "https://")
	}

	trustedOrigins := []string{
		domain,
		"localhost",
		"localhost:9090",
		"127.0.0.1",
		"127.0.0.1:9090",
	}
	slog.Info("CSRF Config", "trusted_origins", trustedOrigins, "secure", cfg.Env == "production")

	// 创建 CSRF 处理器
	csrfHandler := csrf.Protect(
		secret,
		csrf.Secure(cfg.Env == "production"), // 生产环境开启 Secure 属性
		csrf.HttpOnly(true),                  // 保护令牌 Cookie 不被 JS 访问
		csrf.SameSite(csrf.SameSiteLaxMode),  // 防护 CSRF 的标准做法
		csrf.Path("/"),                       // Cookie 作用域
		csrf.FieldName("csrf_token"),         // 表单中的字段名
		csrf.CookieName("idiary_csrf"),       // 自定义 Cookie 名称，避免使用默认的 _gorilla_csrf
		csrf.RequestHeader("X-CSRF-Token"),   // 自定义请求头
		csrf.TrustedOrigins(trustedOrigins),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CSRF 验证失败时的处理
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)

			reason := csrf.FailureReason(r)
			reqID := r.Header.Get("X-Request-ID")

			slog.Warn("CSRF check failed",
				"reason", reason,
				"path", r.URL.Path,
				"method", r.Method,
				"request_id", reqID,
				"remote_addr", r.RemoteAddr,
				"origin", r.Header.Get("Origin"),
				"referer", r.Header.Get("Referer"),
			)

			response := errors.NewAPIErrorResponse(
				errors.ErrOriginInvalid.Code(),
				errors.ErrOriginInvalid.Key(),
				reason,
				reqID,
			)
			json.NewEncoder(w).Encode(response)
		})),
	)

	return func(c *gin.Context) {
		// 标记是否调用了 next handler
		nextCalled := false

		// 使用适配器将 gorilla/csrf 集成到 Gin
		csrfHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 将生成的 Token 存入上下文，供后续页面渲染使用
			c.Set("CsrfToken", csrf.Token(r))
			// 更新请求对象（包含 CSRF context）
			c.Request = r
			nextCalled = true
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)

		// 如果 next 未被调用（说明 CSRF 验证失败），则终止后续处理
		if !nextCalled {
			c.Abort()
		}
	}
}
