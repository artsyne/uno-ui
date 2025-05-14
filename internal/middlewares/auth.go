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

// 需要登录访问地址
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if controllers.IsLogin(c) {
			c.Next()
		} else {
			logrus.Warnf("User not authorized to visit %s", c.Request.RequestURI)
			c.HTML(http.StatusForbidden, "error/403", nil)
			c.Abort()
		}
	}
}

// http 中间件
// 在http上下中注入共享数据，如用户信息、css代码等
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取cookie
		cookieEnUserData := pkg.ReadCookie(c.Request, "ud")
		logrus.Debugf("Cookie value: %s", cookieEnUserData)

		//获取是否有cookie
		if len(cookieEnUserData) < 1 {
			logrus.Debug("No cookie found")
			return
		}

		//解密cookie
		byteUserData, err := pkg.AESDecrypt(cookieEnUserData, configs.GetConfig().Cookie.Secret)
		if err != nil || byteUserData == nil {
			logrus.Errorf("Failed to decrypt cookie: %v", err)
			return
		}
		logrus.Debugf("Decrypted cookie data: %s", string(byteUserData))

		//反序列化为对象
		var userData models.Login
		err = json.Unmarshal(byteUserData, &userData)
		if err != nil {
			logrus.Errorf("Failed to unmarshal user data: %v", err)
			return
		}
		logrus.Debugf("User data: %+v", userData)

		//判断认证信息

		//写入上下文信息
		c.Set("UserData", userData)
		c.Set("IsLogin", true)
		c.Next()
	}
}
