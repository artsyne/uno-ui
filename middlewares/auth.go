package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/biuren/hebe/controllers"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
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
		cookieEnUserData := utils.ReadCookie(c.Request, "ud")
		//获取是否有cookie
		if len(cookieEnUserData) < 1 {
			return
		}

		//解密cookie
		byteUserData, err := utils.AESDecrypt(cookieEnUserData, settings.GetConfig().Cookie.SecretKey)
		if err != nil || byteUserData == nil {
			return
		}

		//反序列化为对象
		var userData models.Login
		err = json.Unmarshal(byteUserData, &userData)
		if err != nil {
			return
		}

		//判断认证信息

		//写入上下文信息
		c.Set("UserData", userData)
		c.Set("IsLogin", true)
		//css code 注入
		c.Set("CSS", settings.GetCSS())
		c.Next()
	}
}