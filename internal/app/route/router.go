package route

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/biuren/hebe/internal/app/controller"
	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/module/diary"
	"github.com/biuren/hebe/internal/module/post"
	"github.com/biuren/hebe/internal/module/user"
	"github.com/gin-gonic/gin"
)

// SetWithContainer 使用DI容器设置路由系统
func SetWithContainer(server *gin.Engine, container *di.Container) {
	setTemplate(server, config.GetViewsPath(), container)
	setStaticFs(server, config.GetStaticPath())
	setPath(server, container)
}

// setPath 设置站点路由，使用DI容器注入服务
func setPath(server *gin.Engine, container *di.Container) {
	// 通过DI容器创建Controller
	controllerInstance := createController(container)

	// 基础页面路由
	server.GET("/index", controllerInstance.RenderIndex)

	// 通过DI容器创建handlers，注册模块路由
	createUserHandler(container).RegisterRoutes(server)
	createDiaryHandler(container).RegisterRoutes(server)
	createPostHandler(container).RegisterRoutes(server)
}

// createController 通过DI容器创建Controller
func createController(container *di.Container) *controller.Controller {
	renderService := di.Get[render.RenderServiceInterface](container)
	return controller.NewController(renderService)
}

// createUserHandler 通过DI容器创建UserHandler
func createUserHandler(container *di.Container) user.HandlerInterface {
	renderService := di.Get[render.RenderServiceInterface](container)
	return user.NewHandlerWithServices(renderService, renderService)
}

// createDiaryHandler 通过DI容器创建DiaryHandler
func createDiaryHandler(container *di.Container) diary.HandlerInterface {
	renderService := di.Get[render.RenderServiceInterface](container)
	return diary.NewHandlerWithServices(renderService, renderService)
}

// createPostHandler 通过DI容器创建PostHandler
func createPostHandler(container *di.Container) post.HandlerInterface {
	renderService := di.Get[render.RenderServiceInterface](container)
	return post.NewHandlerWithServices(renderService, renderService)
}

// #region 加载模版文件，并注入全局变量
// 通过原生 Go 能力加载 "views" 目录中的模版文件
func setTemplate(server *gin.Engine, basePath string, container *di.Container) {
	// 从DI容器获取RenderService
	renderService := di.Get[render.RenderServiceInterface](container)

	// 注册所有模版函数
	tmpl := template.New("").Funcs(renderService.RegisterFuncs())

	// 加载所有 .html 文件
	err := loadTemplates(tmpl, basePath)
	if err != nil {
		panic(err)
	}

	// 设置 Gin 原生的 HTML 渲染器
	server.SetHTMLTemplate(tmpl)
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
