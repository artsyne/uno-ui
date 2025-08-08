package route

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/biuren/hebe/internal/app/controller"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/internal/app/templates"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/di"
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

	// 通过DI容器创建controllers，注册模块路由
	createUserController(container).RegisterRoutes(server)
	createDiaryController(container).RegisterRoutes(server)
	createPostController(container).RegisterRoutes(server)
}

// createController 通过DI容器创建Controller
func createController(container *di.Container) *controller.Controller {
	renderService := di.Get[*view.ViewService](container)
	return controller.NewController(renderService)
}

// createUserController 通过DI容器创建UserController
func createUserController(container *di.Container) *controller.UserController {
	userService := di.Get[user.ServiceInterface](container)
	renderService := di.Get[*view.ViewService](container)
	authHelper := di.Get[user.AuthHelperInterface](container)
	return controller.NewUserController(userService, renderService, authHelper)
}

// createDiaryController 通过DI容器创建DiaryController
func createDiaryController(container *di.Container) *controller.DiaryController {
	return controller.NewDiaryController()
}

// createPostController 通过DI容器创建PostController
func createPostController(container *di.Container) *controller.PostController {
	return controller.NewPostController()
}

// #region 加载模版文件，并注入全局变量
// 通过原生 Go 能力加载 "views" 目录中的模版文件
func setTemplate(server *gin.Engine, basePath string, container *di.Container) {
	// 1. 获取配置
	cfg := di.Get[*config.Config](container)

	// 2. 初始化本地化管理器（消息会自动加载到全局缓存）
	i18nManager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(i18nManager)

	// 3. 注册所有模版函数（直接使用templates包）
	tmpl := template.New("").Funcs(templates.RegisterAllFuncs())

	// 6. 加载所有 .html 文件
	if err := loadTemplates(tmpl, basePath); err != nil {
		panic(err)
	}

	// 7. 设置 Gin 原生的 HTML 渲染器
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
