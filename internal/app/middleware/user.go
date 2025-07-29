package middleware

import (
	"log/slog"
	"net/http"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// AuthRequired 需要登录访问的中间件
func AuthRequired() gin.HandlerFunc {
	auth := user.NewAuthHelper(config.GetConfig())
	return func(c *gin.Context) {
		if isAuth, _ := auth.IsAuthenticated(c); isAuth {
			c.Next()
		} else {
			slog.Warn("User not authorized", "uri", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "error", nil)
			c.Abort()
		}
	}
}

// SharedData 共享数据中间件
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := user.NewAuthHelper(config.GetConfig())
		// 解析用户登录信息
		authUser := auth.ParseUserFromCookie(c)

		// 注入登录状态和用户数据到上下文
		if authUser != nil {
			c.Set("IsLogin", true)
			c.Set("AuthUser", authUser)
		} else {
			c.Set("IsLogin", false)
		}

		c.Next()
	}
}