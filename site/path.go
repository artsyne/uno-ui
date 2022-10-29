package site

import (
	"net/http"

	"chinadiary.com/website/controllers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//设置站点router
func SetPath(router *gin.Engine) {

	router.GET("/index", controllers.RenderIndex)

	// daybook
	daybook := router.Group("/daybook")
	{
		daybook.GET("/:spaceid", controllers.RenderDaybookIndex)
	}

	//表单渲染
	form := router.Group("/form")
	{
		form.GET("/daybook", controllers.RenderDaybookForm)
	}

	// diary
	diary := router.Group("/diary")
	{
		diary.GET("/", controllers.RenderDiary)
		diary.GET("/write", controllers.RenderDiaryWrite)
	}
	router.POST("/info", controllers.PostInfo)

	// my
	router.GET("/my", controllers.RenderMyIndex)

	// passport
	router.GET("/passport/", controllers.RenderPassport)
	router.GET("/register", controllers.RenderRegister)
	router.POST("/register", controllers.PostRegister)
	router.GET("/login", controllers.RenderLogin)
	router.POST("/login", controllers.PostLogin)
	router.POST("/loginout", controllers.PostLoginOut)
	router.GET("/findpwd", controllers.RenderFindPwd)
	router.POST("/findpwd", controllers.PostFindPwd)
	router.GET("/resetpwd", controllers.RenderResetPwd)
	router.POST("/resetpwd", controllers.PostResetPwd)
	router.GET("/activateuser", controllers.RenderActivateUser)

	router.GET("/u/:uid", controllers.RenderMemberIndex)

	member := router.Group("/member")
	{
		member.GET("/", controllers.RenderMemberIndex)
	}

	// 日记本操作
	router.POST("/daybook", controllers.DoSpace)

	logrus.Info("router初始化成功！")
}

//需要登录访问地址
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
