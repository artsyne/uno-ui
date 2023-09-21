package main

import (
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/site"
	"github.com/biuren/hebe/system"
)

func main() {
	//设置gin模式∏
	gin.SetMode(gin.DebugMode)

	//设置日志
	system.SetLogger()

	//初始化配置信息
	system.LoadConfig()

	//初始化站点
	//加载配置文件，链接DB，设置cookie默认值
	//启动定时任务
	router := site.InitSite()

	// 初始化Router，创建gin 默认router
	//设置模板文件，静态文件Path, 加载上下文共享数据，加载中间件
	site.SetRouter(router)

	//设置路由Path
	site.SetPath(router)

	// Listen and server
	router.Run(":" + system.GetConfig().SitePort)
}