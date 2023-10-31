package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/middlewares"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/routes"
	"github.com/biuren/hebe/services"
	"github.com/biuren/hebe/utils"
)

func main() {
	//设置gin模式∏
	gin.SetMode(gin.DebugMode)

	//设置日志
	utils.SetLogger()

	//初始化配置信息
	settings.LoadConfig()

	//链接数据库
	models.ConnectToDB()
	logrus.Info("成功数据库连接")

	//设置cookie默认值
	config := settings.GetConfig()
	utils.SetCookieDefVal(config.Cookie.Domain, "/", config.Cookie.ExpireDays, config.Cookie.HttpOnly)
	logrus.Info("成功设置Cookie默认值")

	//初始化默认gin router
	router := gin.Default()

	// 初始化中间件
	middlewares.Attach(router)
	logrus.Info("成功初始化中间件")

	//初始化Router，创建gin 默认router
	//设置模板文件，静态文件Path, 加载上下文共享数据
	routes.Set(router)
	logrus.Info("成功初始化网站静态文件目录及路由")

	//启动定时任务
	services.RunTask()
	logrus.Info("成功启动定时任务")

	// Listen and server
	router.Run(":" + config.SitePort)
}
