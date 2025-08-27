package middleware

import (
	"log/slog"
	"net/http"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// AuthRequired 需要登录访问的中间件
// 检查用户是否已登录，未登录用户将被重定向到错误页面
//
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
//
// 功能:
//   - 验证用户登录状态
//   - 已登录用户继续执行后续处理
//   - 未登录用户返回403错误页面并中止请求
//
// 使用场景:
//   - 需要用户登录才能访问的路由
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
// 解析用户登录信息并注入到请求上下文中，供后续处理使用
//
// 返回:
//   - gin.HandlerFunc: Gin中间件处理函数
//
// 功能:
//   - 从Cookie中解析用户信息
//   - 设置登录状态标识(IsLogin)
//   - 注入用户数据到上下文(AuthUser)
//   - 为模板渲染提供用户状态数据
//
// 上下文变量:
//   - IsLogin: bool - 用户是否已登录
//   - AuthUser: *user.User - 当前登录用户信息(仅登录时)
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