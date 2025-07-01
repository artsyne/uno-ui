package user

import (
	"encoding/json"
	"net/http"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Middleware 用户中间件结构体，用于缓存配置
type Middleware struct {
	config *configs.Config
	auth   *AuthHelper
}

// NewMiddleware 创建用户中间件实例
func NewMiddleware() *Middleware {
	config := configs.GetConfig()
	if config == nil {
		panic("Failed to get config during user middleware initialization - config is nil")
	}

	return &Middleware{
		config: config,
		auth:   NewAuthHelper(),
	}
}

// AuthRequired 需要登录访问的中间件
func AuthRequired() gin.HandlerFunc {
	middleware := NewMiddleware()
	return func(c *gin.Context) {
		if isAuth, _ := middleware.auth.IsAuthenticated(c); isAuth {
			c.Next()
		} else {
			logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "error", nil)
			c.Abort()
		}
	}
}

// SharedData 共享数据中间件
func SharedData() gin.HandlerFunc {
	middleware := NewMiddleware()
	return func(c *gin.Context) {
		// 解析用户登录信息
		userData := middleware.parseUserFromCookie(c)

		// 注入登录状态和用户数据到上下文
		if userData != nil {
			c.Set("IsLogin", true)
			c.Set("UserData", *userData)
		} else {
			c.Set("IsLogin", false)
		}

		c.Next()
	}
}

// parseUserFromCookie 从Cookie中解析用户登录信息
func (m *Middleware) parseUserFromCookie(c *gin.Context) *Login {
	cookieEnUserData := pkg.ReadCookie(c.Request, "ud")
	if len(cookieEnUserData) < 1 {
		return nil
	}

	byteUserData, err := pkg.AESDecrypt(cookieEnUserData, m.config.Cookie.Secret)
	if err != nil || byteUserData == nil {
		return nil
	}

	var userData Login
	err = json.Unmarshal(byteUserData, &userData)
	if err != nil {
		return nil
	}

	return &userData
}
