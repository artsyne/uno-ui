package user

import (
	"encoding/json"
	"net/http"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MiddlewareInterface 用户中间件接口
// 定义了用户认证中间件相关的核心方法
type MiddlewareInterface interface {
	// ParseUserFromCookie 从Cookie中解析用户登录信息
	ParseUserFromCookie(c *gin.Context) *Login

	// GetAuthHelper 获取认证助手实例
	GetAuthHelper() AuthHelperInterface
}

// middleware 用户中间件结构体，用于缓存配置
type middleware struct {
	config *config.Config
	auth   AuthHelperInterface
}

// NewMiddleware 创建用户中间件实例
func NewMiddleware() MiddlewareInterface {
	config := config.GetConfig()
	if config == nil {
		panic("Failed to get config during user middleware initialization - config is nil")
	}

	return &middleware{
		config: config,
		auth:   NewAuthHelper(),
	}
}

// ParseUserFromCookie 从Cookie中解析用户登录信息
func (m *middleware) ParseUserFromCookie(c *gin.Context) *Login {
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

// GetAuthHelper 获取认证助手实例
func (m *middleware) GetAuthHelper() AuthHelperInterface {
	return m.auth
}

// AuthRequired 需要登录访问的中间件
// 包级别函数，提供全局可用的认证中间件
func AuthRequired() gin.HandlerFunc {
	middleware := NewMiddleware()
	return func(c *gin.Context) {
		if isAuth, _ := middleware.GetAuthHelper().IsAuthenticated(c); isAuth {
			c.Next()
		} else {
			logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "error", nil)
			c.Abort()
		}
	}
}

// SharedData 共享数据中间件
// 包级别函数，提供全局可用的共享数据中间件
func SharedData() gin.HandlerFunc {
	middleware := NewMiddleware()
	return func(c *gin.Context) {
		// 解析用户登录信息
		userData := middleware.ParseUserFromCookie(c)

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
