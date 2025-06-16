package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/pkg"
	"github.com/biuren/hebe/website/controllers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/**
 * 认证中间件模块
 *
 * 负责认证相关的中间件功能：
 * - 用户会话数据注入
 * - 访问权限控制
 * - 与 helpers/auth.go 协作提供完整的认证体系
 *
 * @author hebe
 * @version 1.0.0
 */

// AuthRequired 需要登录访问的中间件
// 检查用户是否已登录，如果未登录则返回403错误页面
// 适用于API接口或需要强制登录验证的页面
//
// 注意：这个中间件依赖于 SharedData() 中间件注入的认证信息
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if controllers.IsLogin(c) {
			c.Next()
		} else {
			logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "error", nil)
			c.Abort()
		}
	}
}

// SharedData 共享数据中间件
// 从Cookie中解析用户信息并注入到Gin上下文中
// 这是认证体系的基础，必须在其他认证相关中间件之前使用
//
// 注入的上下文数据：
//   - IsLogin: bool - 用户是否已登录
//   - UserData: models.Login - 用户登录信息（如果已登录）
//   - 其他系统配置信息
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析用户登录信息
		userData := parseUserFromCookie(c)

		// 注入登录状态和用户数据到上下文
		if userData != nil {
			c.Set("IsLogin", true)
			c.Set("UserData", *userData)
		} else {
			c.Set("IsLogin", false)
		}

		// 注入其他共享数据
		injectSystemData(c)

		c.Next()
	}
}

// parseUserFromCookie 从Cookie中解析用户登录信息
// 这是内部函数，负责具体的用户认证逻辑
func parseUserFromCookie(c *gin.Context) *models.Login {
	// 获取cookie
	cookieEnUserData := pkg.ReadCookie(c.Request, "ud")
	logrus.Debugf("Cookie value: %s", cookieEnUserData)

	// 检查是否有cookie
	if len(cookieEnUserData) < 1 {
		logrus.Debug("No cookie found")
		return nil
	}

	// 解密cookie
	byteUserData, err := pkg.AESDecrypt(cookieEnUserData, configs.GetConfig().Cookie.Secret)
	if err != nil || byteUserData == nil {
		logrus.Errorf("Failed to decrypt cookie: %v", err)
		return nil
	}
	logrus.Debugf("Decrypted cookie data: %s", string(byteUserData))

	// 反序列化为对象
	var userData models.Login
	err = json.Unmarshal(byteUserData, &userData)
	if err != nil {
		logrus.Errorf("Failed to unmarshal user data: %v", err)
		return nil
	}
	logrus.Debugf("User data: %+v", userData)

	return &userData
}

// injectSystemData 注入系统级共享数据
// 可以在这里添加其他需要在模板中使用的全局数据
func injectSystemData(c *gin.Context) {
	// 可以在这里添加其他系统数据
	// 例如：站点配置、主题设置等
	// c.Set("SiteName", configs.GetSiteName())
	// c.Set("Version", configs.GetVersion())
}
