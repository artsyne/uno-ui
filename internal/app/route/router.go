package route

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/biuren/hebe/internal/app/controller"
	commentController "github.com/biuren/hebe/internal/app/controller/comment"
	diaryController "github.com/biuren/hebe/internal/app/controller/diary"
	entryController "github.com/biuren/hebe/internal/app/controller/entry"
	userController "github.com/biuren/hebe/internal/app/controller/user"
	"github.com/biuren/hebe/internal/app/templates"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/core/di"
	commentDomain "github.com/biuren/hebe/internal/domain/comment"
	diaryDomain "github.com/biuren/hebe/internal/domain/diary"
	entryDomain "github.com/biuren/hebe/internal/domain/entry"
	likeDomain "github.com/biuren/hebe/internal/domain/like"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
)

// SetWithContainer 使用DI容器设置路由系统
// 初始化完整的Web服务器路由配置，包括模板、静态文件和业务路由
//
// 参数:
//   - server: Gin引擎实例
//   - container: DI依赖注入容器
//
// 功能:
//   - 设置HTML模板系统
//   - 配置静态文件服务
//   - 注册所有业务路由
func SetWithContainer(server *gin.Engine, container *di.Container) {
	setTemplate(server, config.GetViewsPath(), container)
	setStaticFs(server, config.GetStaticPath())
	setPath(server, container)
}

// setPath 设置站点路由，使用DI容器注入服务
// 注册所有业务模块的路由，包括用户、日记、文章等功能模块
//
// 参数:
//   - server: Gin引擎实例
//   - container: DI依赖注入容器
//
// 功能:
//   - 创建并注册IndexController路由
//   - 创建并注册UserController路由
//   - 创建并注册DiaryController路由
//   - 创建并注册PostController路由
func setPath(server *gin.Engine, container *di.Container) {
	// 通过DI容器创建IndexController
	controllerInstance := createIndexController(container)

	// 基础页面路由
	server.GET("/index", controllerInstance.RenderIndex)

	// 通过DI容器创建controllers，注册模块路由
	createUserController(container).RegisterRoutes(server)
	createDiaryController(container).RegisterRoutes(server)
	createEntryController(container).RegisterRoutes(server)

	// 注册评论路由
	commentCtrl := createCommentController(container)
	apiGroup := server.Group("/api")
	commentController.RegisterRoutes(apiGroup, commentCtrl)
}

// controllerDependencies 控制器通用依赖结构
type controllerDependencies struct {
	userService    *userDomain.Service
	diaryService   *diaryDomain.Service
	entryService   *entryDomain.Service
	commentService *commentDomain.Service
	likeService    likeDomain.Service
	authHelper     *userDomain.AuthHelper
	config         *config.Config
	renderService  *view.ViewService
}

// getControllerDependencies 从DI容器获取控制器通用依赖
// 统一获取所有控制器需要的公共服务依赖，避免重复代码
//
// 参数:
//   - container: DI依赖注入容器
//
// 返回:
//   - *controllerDependencies: 包含所有公共依赖的结构体
//
// 包含的依赖:
//   - userService: 用户业务服务
//   - diaryService: 日记业务服务
//   - entryService: 日记业务服务
//   - likeService: 点赞业务服务
//   - authHelper: 认证辅助服务
//   - config: 应用配置
//   - renderService: 视图渲染服务
func getControllerDependencies(container *di.Container) *controllerDependencies {
	return &controllerDependencies{
		userService:    di.Get[*userDomain.Service](container),
		diaryService:   di.Get[*diaryDomain.Service](container),
		entryService:   di.Get[*entryDomain.Service](container),
		commentService: di.Get[*commentDomain.Service](container),
		likeService:    di.Get[likeDomain.Service](container),
		authHelper:     di.Get[*userDomain.AuthHelper](container),
		config:         di.Get[*config.Config](container),
		renderService:  di.Get[*view.ViewService](container),
	}
}

// createIndexController 通过DI容器创建IndexController
// 创建首页控制器实例，处理网站首页相关请求
//
// 参数:
//   - container: DI依赖注入容器
//
// 返回:
//   - *controller.IndexController: 首页控制器实例
func createIndexController(container *di.Container) *controller.IndexController {
	deps := getControllerDependencies(container)
	return controller.NewIndexController(deps.config, deps.renderService)
}

// createUserController 通过DI容器创建UserController
// 创建用户控制器实例，处理用户注册、登录、个人资料等功能
//
// 参数:
//   - container: DI依赖注入容器
//
// 返回:
//   - *userController.UserController: 用户控制器实例
func createUserController(container *di.Container) *userController.UserController {
	deps := getControllerDependencies(container)
	return userController.NewUserController(deps.config, deps.userService, deps.renderService, deps.authHelper)
}

// createDiaryController 通过DI容器创建DiaryController
// 创建日记控制器实例，处理日记本的创建、编辑、删除等功能
//
// 参数:
//   - container: DI依赖注入容器
//
// 返回:
//   - *diaryController.DiaryController: 日记控制器实例
func createDiaryController(container *di.Container) *diaryController.DiaryController {
	deps := getControllerDependencies(container)
	return diaryController.NewDiaryController(deps.config, deps.diaryService, deps.userService, deps.likeService, deps.authHelper, deps.renderService)
}

// createEntryController 通过DI容器创建EntryController
// 创建文章控制器实例，处理文章的创建、编辑、删除、查看等功能
//
// 参数:
//   - container: DI依赖注入容器
//
// 返回:
//   - *entryController.EntryController: 文章控制器实例
func createEntryController(container *di.Container) *entryController.EntryController {
	deps := getControllerDependencies(container)
	return entryController.NewEntryController(deps.config, deps.entryService, deps.diaryService, deps.authHelper, deps.renderService, deps.userService, deps.likeService)
}

// createCommentController 通过DI容器创建CommentController
// 创建评论控制器实例，处理评论的创建、删除、查看等功能
//
// 参数:
//   - container: DI依赖注入容器
//
// 返回:
//   - *commentController.CommentController: 评论控制器实例
func createCommentController(container *di.Container) *commentController.CommentController {
	deps := getControllerDependencies(container)
	return commentController.NewCommentController(deps.config, deps.commentService, deps.authHelper, deps.renderService)
}

// setTemplate 设置模板引擎
// 配置Gin的HTML模板系统，加载所有模板文件并设置自定义函数
//
// 参数:
//   - server: Gin引擎实例
//   - basePath: 模板文件路径
//   - container: DI依赖注入容器
//
// 功能:
//   - 加载指定路径下的所有HTML模板
//   - 注册模板自定义函数
//   - 设置模板渲染引擎
func setTemplate(server *gin.Engine, basePath string, container *di.Container) {
	// 1. 从DI容器获取ViewService（已包含i18nManager）
	viewService := di.Get[*view.ViewService](container)

	// 2. 使用ViewService中的i18nManager，避免重复初始化
	templates.SetGlobalManager(viewService.GetI18nManager())

	// 3. 注册所有模版函数（直接使用templates包）
	tmpl := template.New("").Funcs(templates.RegisterAllFuncs())

	// 6. 加载所有 .html 文件
	if err := loadTemplates(tmpl, basePath); err != nil {
		panic(err)
	}

	// 7. 设置 Gin 原生的 HTML 渲染器
	server.SetHTMLTemplate(tmpl)
}

// loadTemplates 加载模板文件
// 递归扫描模板目录，加载所有HTML模板文件
//
// 参数:
//   - tmpl: 模板对象
//   - basePath: 模板文件根目录路径
//
// 返回:
//   - error: 加载过程中的错误信息
//
// 功能:
//   - 递归遍历模板目录
//   - 解析所有.html文件
//   - 将模板添加到模板对象中
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

// setStaticFs 设置静态文件服务
// 配置静态资源文件的HTTP服务，包括CSS、JS、图片等
//
// 参数:
//   - server: Gin引擎实例
//   - staticPath: 静态文件路径
//
// 功能:
//   - 设置/public和/static路径的静态文件服务
//   - 配置favicon.ico文件访问
//   - 支持静态资源的HTTP访问
func setStaticFs(server *gin.Engine, staticPath string) {
	//初始化静态文件
	//better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	server.StaticFS("/public", http.Dir(staticPath))
	server.StaticFS("/static", http.Dir(staticPath))
	server.StaticFile("/favicon.ico", "./public/favicon.ico")
}
