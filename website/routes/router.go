package routes

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/modules/diary"
	"github.com/biuren/hebe/internal/modules/post"
	"github.com/biuren/hebe/internal/modules/user"
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
	// 创建模块处理器
	userHandler := user.NewHandler()
	diaryHandler := diary.NewHandler()
	postHandler := post.NewHandler()

	// 基础页面路由
	server.GET("/index", controllers.RenderIndex)

	// 用户相关路由 - 使用新的User模块
	userHandler.RegisterRoutes(server)

	// 日记本相关路由 - 使用新的Diary模块（包含所有diary相关路由）
	diaryHandler.RegisterRoutes(server)

	// 日记相关路由 - 使用新的Post模块（包含所有post相关路由）
	postHandler.RegisterRoutes(server)

}
