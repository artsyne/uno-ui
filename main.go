package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/biuren/hebe/middlewares"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/routes"
	"github.com/biuren/hebe/services"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
)

func main() {
	//设置gin模式∏
	gin.SetMode(gin.DebugMode)

	//设置日志
	utils.SetLogger()

	//初始化配置信息
	settings.LoadConfig()
	logrus.Info("config,params 加载成功！")

	//链接数据库，若异常将 panic
	models.ConnectToDB()
	logrus.Info("成功数据库连接")

	//设置cookie默认值
	config := settings.GetConfig()
	utils.SetCookieDefVal(config.Cookie.Domain, "/", config.Cookie.ExpireDays, config.Cookie.HttpOnly)
	logrus.Info("成功设置Cookie默认值")

	//初始化Gin server
	startServer(config.SitePort)
}

// 初始化服务，设置路由、加载中间件等
func startServer(sitePort string) {
	// 初始化默认gin server
	server := gin.Default()

	// 初始化中间件
	// TODO: gin-csrf,gin-limit，error
	AttachMiddleware(server)
	logrus.Info("成功初始化中间件")

	// 初始化Router，创建gin 默认router
	// 设置模板文件，静态文件Path, 加载上下文共享数据
	routes.Set(server)
	logrus.Info("成功初始化网站静态文件目录及路由")

	// 启动定时任务
	services.RunTask()
	logrus.Info("成功启动定时任务")

	// Listen and server
	server.Run(":" + sitePort)
}

// 在http请求上下文中注入中间件
func AttachMiddleware(server *gin.Engine){
	// 加载上下文共享数据，如登录数据
	server.Use(middlewares.ErrorHandler(),middlewares.SharedData())
}
