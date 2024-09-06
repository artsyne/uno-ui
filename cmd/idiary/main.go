package main

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/middlewares"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
	"github.com/biuren/hebe/website/routes"
)

func main() {
	// 设置 gin 模式
	gin.SetMode(gin.DebugMode)

	// 设置日志
	pkg.SetLogger()

	// 捕获启动过程中的 panic 错误
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"stack": string(debug.Stack()),
			}).Error("启动函数 panic 错误")
		}
	}()

	// 初始化 site
	initialize()

	// 初始化 Gin server
	startServer(configs.GetConfig().SitePort)
}

// 初始化 site
func initialize() {

	// 初始化配置信息
	configs.LoadConfig()
	logrus.Info("config、params 加载成功！")

	// 链接数据库，若异常将 panic
	models.ConnectToDB()
	logrus.Info("数据库连接成功")

	// 设置 cookie 默认值
	config := configs.GetConfig()
	pkg.SetCookieDefVal(config.Cookie.Domain, "/", config.Cookie.ExpireDays, config.Cookie.HttpOnly)
	logrus.Info("成功设置 Cookie 默认值")
}

// 初始化服务，设置路由、加载中间件等
func startServer(sitePort string) {
	// 初始化默认gin server
	server := gin.Default()

	// 初始化中间件
	// TODO: gin-csrf、gin-limit、error
	attachMiddleware(server)
	logrus.Info("成功初始化中间件")

	// 初始化Router，创建 gin 默认 router
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
func attachMiddleware(server *gin.Engine) {
	// 加载上下文共享数据，如登录数据
	server.Use(middlewares.PanicHandler(), middlewares.ErrorHandler(), middlewares.SharedData())
}
