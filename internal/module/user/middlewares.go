package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// AuthRequired 需要登录访问的中间件
func AuthRequired() gin.HandlerFunc {
	auth := NewAuthHelper()
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
		// 解析用户登录信息
		authUser := parseUserFromCookie(c)

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


// parseUserFromCookie 从Cookie中解析用户登录信息
func parseUserFromCookie(c *gin.Context) *AuthUser {
	cookieEnUserData := pkg.ReadCookie(c.Request, "ud")
	if len(cookieEnUserData) < 1 {
		return nil
	}

	config := config.GetConfig()
	if config == nil {
		return nil
	}

	byteUserData, err := pkg.AESDecrypt(cookieEnUserData, config.Cookie.Secret)
	if err != nil || byteUserData == nil {
		return nil
	}

	var userData AuthUser
	err = json.Unmarshal(byteUserData, &userData)
	if err != nil {
		return nil
	}

	return &userData
}
