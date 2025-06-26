package routes

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/website/controllers"
	"github.com/gin-gonic/gin"
)

// 初始化Router，加载静态文件，上下文共享数据，中间件
func Set(server *gin.Engine) {
	// 初始化 Views 模板
	setTemplate(server, configs.GetViewsPath())

	// 初始化静态文件目录
	setStaticFs(server, configs.GetStaticPath())

	// 初始化站点动态目录
	setPath(server)
}

// #region 加载模版文件，并注入全局变量
// 通过原生 Go 能力加载 "views" 目录中的模版文件
func setTemplate(server *gin.Engine, basePath string) {
	// 注册模版函数
	tmpl := template.New("").Funcs(template.FuncMap{
		"formatTime": formatTime,
	})

	// 加载所有 .html 文件
	err := loadTemplates(tmpl, basePath)
	if err != nil {
		panic(err)
	}

	// 设置 Gin 原生的 HTML 渲染器
	server.SetHTMLTemplate(tmpl)
}

// formatTime 格式化时间的辅助函数，utc 时间需精度到毫秒
func formatTime(utc int64, layout string) string {
	return time.Unix(utc/1000, 0).Format(layout)
}

// loadTemplates 递归遍历目录，加载所有 .html 文件
func loadTemplates(tmpl *template.Template, basePath string) error {
	return filepath.Walk(basePath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
			// 使用相对路径作为模板名称，避免命名冲突
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}

			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// 使用相对路径作为模板名称解析
			_, err = tmpl.New(relPath).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// #endregion

// 设置静态文件目录，存放js、css，以及上传文件目录upload
func setStaticFs(server *gin.Engine, staticPath string) {
	//初始化静态文件
	//better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	server.StaticFS("/public", http.Dir(staticPath))
	server.StaticFS("/static", http.Dir(staticPath))
	server.StaticFile("/favicon.ico", "./public/favicon.ico")
}

// setPath 设置站点路由
// 创建控制器实例并配置所有HTTP路由
func setPath(server *gin.Engine) {
	// 创建控制器实例
	diaryController := controllers.NewDiaryController()
	postController := controllers.NewPostController()

	// 基础页面路由
	server.GET("/index", controllers.RenderIndex)

	// 日记相关路由组
	diary := server.Group("/diary")
	{
		diary.GET("/", diaryController.RenderDiary)
	}

	// 用户认证相关路由
	server.GET("/register", controllers.RenderRegister)
	server.POST("/register", controllers.HandleUserReg)
	server.POST("/check-username", controllers.HandleCheckUserName)
	server.POST("/check-email", controllers.HandleCheckEmail)

	server.GET("/login", controllers.RenderLogin)
	server.POST("/login", controllers.HandleLogin)

	server.POST("/loginout", controllers.HandleLoginOut)

	server.GET("/findpwd", controllers.RenderFindPwd)
	server.POST("/findpwd", controllers.HandleFindPwd)

	server.GET("/resetpwd", controllers.RenderResetPwd)
	server.POST("/resetpwd", controllers.HandleResetPwd)

	server.GET("/activateuser", controllers.RenderActivateUser)

	// 用户主页路由
	server.GET("/u/:uid", controllers.RenderMemberIndex)
	server.GET("u/:uid/settings", controllers.RenderMemberSettings)
	server.GET("/u/:uid/:pageindex", controllers.RenderMemberIndex)

	// 日记本页面路由组
	daybook := server.Group("/d")
	{
		daybook.GET("/:diaryid", diaryController.RenderDiaryIndex)
		daybook.GET("/:diaryid/:pageindex", diaryController.RenderDiaryIndex)
	}

	// 表单渲染路由组
	form := server.Group("/form")
	{
		form.GET("/diary", diaryController.RenderDiaryForm)
		form.GET("/post", postController.RenderPostForm)
	}

	// 数据操作路由
	// 日记本相关操作
	server.POST("/diary", diaryController.HandleDiary)
	// 日记相关操作
	server.POST("/post", postController.HandlePost)
}
