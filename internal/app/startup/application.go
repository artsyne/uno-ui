package bootstrap

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/app/route"
	"github.com/biuren/hebe/internal/core/initializer"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	utilsMiddleware "github.com/biuren/hebe/internal/utils/middleware"
	"github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
)

// Application 应用程序主体
type Application struct {
	config    *config.Config
	server    *gin.Engine
	container *di.Container
	coreInit  *bootstrap.Initializer
}

// NewApplication 创建应用程序实例
func NewApplication(config *config.Config) *Application {
	return &Application{
		config:    config,
		coreInit: bootstrap.NewInitializer(config),
	}
}

// Bootstrap 启动应用程序
// 初始化核心基础设施，然后初始化应用层组件
func (a *Application) Bootstrap() error {
	// 1. 初始化核心基础设施
	if err := a.coreInit.InitializeCore(); err != nil {
		return err // 已经是统一错误格式
	}

	// 2. 获取核心组件
	a.container = a.coreInit.GetContainer()

	// 验证核心组件是否正确初始化
	if a.config == nil {
		return errors.ConfigError.WithDetails("配置获取失败")
	}
	if a.container == nil {
		return errors.DIServiceNotRegistered.WithDetails("DI容器获取失败")
	}

	// 3. 初始化应用层组件
	if err := a.initApplicationLayer(); err != nil {
		return err // 已经是统一错误格式
	}

	return nil
}

// Run 运行应用程序
func (a *Application) Run() error {
	// 启动服务器
	addr := ":" + a.config.App.Port
	slog.Info("Starting server", "addr", addr)

	if err := a.server.Run(addr); err != nil {
		slog.Error("服务器启动失败", "error", err, "addr", addr)
		os.Exit(1)
	}

	return nil
}

// initApplicationLayer 初始化应用层组件
func (a *Application) initApplicationLayer() error {
	// 1. 注册应用层服务到DI容器
	if err := a.registerAppServices(); err != nil {
		return err
	}

	// 2. 设置 Cookie 默认值
	if err := a.initCookie(); err != nil {
		return err
	}

	// 3. 启动定时任务
	if err := a.initScheduler(); err != nil {
		return err
	}

	// 4. 初始化 Gin 服务器
	if err := a.initGinServer(); err != nil {
		return err
	}

	return nil
}

// registerAppServices 注册应用层服务到DI容器
func (a *Application) registerAppServices() error {
	// 创建 I18nManager 实例
	i18nManager := i18n.NewLocalizationManager(a.config)
	di.Value(a.container, i18nManager)
	
	// 创建 ViewService 实例
	viewService := view.NewViewService(a.config, i18nManager)
	di.Value(a.container, viewService)
	
	// 创建仓储实例
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	postRepo := repository.NewPostRepository(database.GetDB())
	
	// 创建服务实例
	taskService := task.NewService(a.config, taskRepo)
	userService := user.NewService(a.config, userRepo, taskService)
	diaryService := diary.NewService(a.config, diaryRepo)
	postService := post.NewService(a.config, userRepo, taskService, diaryRepo, postRepo)
	
	// 注册为预创建实例，避免循环依赖
	di.Value(a.container, userRepo)
	di.Value(a.container, taskRepo)
	di.Value(a.container, diaryRepo)
	di.Value(a.container, postRepo)
	di.Value(a.container, taskService)
	di.Value(a.container, userService)
	di.Value(a.container, diaryService)
	di.Value(a.container, postService)

	// 注册 AuthHelper 为单例服务
	di.Singleton(a.container, func() *user.AuthHelper {
		return user.NewAuthHelper(a.config)
	})

	slog.Info("成功注册应用层服务")
	return nil
}

// initCookie 设置 Cookie 默认值
func (a *Application) initCookie() error {
	pkg.SetCookieDefVal(
		a.config.Cookie.Domain,
		"/",
		a.config.Cookie.ExpireDays,
		a.config.Cookie.HttpOnly,
	)
	slog.Info("成功设置 Cookie 默认值")
	return nil
}

// initScheduler 启动定时任务
func (a *Application) initScheduler() error {
	// 从 DI 容器安全获取 TaskService 实例并启动调度器
	taskService, err := di.TryGet[*task.Service](a.container)
	if err != nil {
		return fmt.Errorf("获取TaskService失败: %w", err)
	}
	taskService.StartScheduler()
	slog.Info("成功启动定时任务")
	return nil
}

// initGinServer 初始化 Gin 服务器
func (a *Application) initGinServer() error {
	// 设置 Gin 模式
	gin.SetMode(a.config.App.GinMode)

	// 初始化默认gin server
	server := gin.New()

	// 初始化中间件
	server.Use(utilsMiddleware.PanicRecovery())   // panic recovery中间件
	server.Use(utilsMiddleware.AppErrorHandler()) // 应用错误处理中间件
	server.Use(middleware.SharedData())            // 用户共享数据中间件
	slog.Info("成功初始化中间件")

	// 初始化路由
	route.SetWithContainer(server, a.container)
	slog.Info("成功初始化网站静态文件目录及路由")

	a.server = server
	return nil
}
