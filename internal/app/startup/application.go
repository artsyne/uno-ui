package bootstrap

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/biuren/hebe/internal/app/controller/validation"
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/app/route"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/core/di"
	bootstrap "github.com/biuren/hebe/internal/core/initializer"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/cache"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/pkg"
)

// Application 应用程序主体
// 负责整个应用程序的启动、初始化和运行管理
//
// 字段:
//   - config: 应用程序配置
//   - server: Gin HTTP服务器实例
//   - container: DI依赖注入容器
//   - coreInit: 核心组件初始化器
type Application struct {
	config    *config.Config
	server    *gin.Engine
	container *di.Container
	coreInit  *bootstrap.Initializer
}

// NewApplication 创建应用程序实例
// 使用提供的配置创建新的应用程序实例
//
// 参数:
//   - config: 应用程序配置对象
//
// 返回:
//   - *Application: 应用程序实例
func NewApplication(config *config.Config) *Application {
	return &Application{
		config:    config,
		coreInit: bootstrap.NewInitializer(config),
	}
}

// Bootstrap 启动应用程序
// 按顺序初始化核心基础设施和应用层组件，为应用程序运行做好准备
//
// 返回:
//   - error: 初始化过程中的错误信息
//
// 初始化步骤:
//   1. 初始化核心基础设施（数据库、缓存等）
//   2. 获取并验证核心组件
//   3. 初始化应用层组件（服务、中间件、路由等）
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
// 启动HTTP服务器，开始监听和处理客户端请求
//
// 返回:
//   - error: 服务器启动过程中的错误信息
//
// 注意:
//   - 此方法会阻塞当前goroutine直到服务器停止
//   - 如果启动失败会调用os.Exit(1)退出程序
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
// 按顺序初始化应用层的各个组件，包括服务注册、配置设置等
//
// 返回:
//   - error: 初始化过程中的错误信息
//
// 初始化步骤:
//   1. 注册应用层服务到DI容器
//   2. 设置Cookie默认值
//   3. 启动定时任务调度器
//   4. 初始化Gin HTTP服务器
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
// 创建并注册所有应用层服务实例，建立服务间的依赖关系
//
// 返回:
//   - error: 服务注册过程中的错误信息
//
// 注册的服务:
//   - I18nManager: 国际化管理器
//   - ViewService: 视图渲染服务
//   - Repository层: 用户、任务、日记、文章仓储
//   - Service层: 用户、任务、日记、文章业务服务
//   - AuthHelper: 认证辅助服务
func (a *Application) registerAppServices() error {
	// 创建 I18nManager 实例
	i18nManager := i18n.NewLocalizationManager(a.config)
	if err := di.Value(a.container, i18nManager); err != nil {
		return err
	}

	// 初始化 Redis 客户端
	redisClient, err := cache.NewRedisClient(a.config.Redis)
	if err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}
	
	// 创建 ViewService 实例
	viewService := view.NewViewService(a.config, i18nManager)
	if err := di.Value(a.container, viewService); err != nil {
		return err
	}
	
	// 创建仓储实例
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	commentRepo := repository.NewCommentRepository()
	likeRepo := repository.NewLikeRepository()
	
	// 创建服务实例
	taskService := task.NewService(a.config, taskRepo)
	userService := user.NewService(a.config, userRepo, taskService, i18nManager, redisClient)
	diaryService := diary.NewService(a.config, diaryRepo, userService, taskService, i18nManager)
	entryService := entry.NewService(a.config, userService, taskService, diaryRepo, entryRepo, i18nManager)
	commentService := comment.NewService(a.config, commentRepo, userService, entryService, diaryService)
	likeService := like.NewLikeService(likeRepo)
	
	// 注册为预创建实例，避免循环依赖
	if err := di.Value(a.container, userRepo); err != nil {
		return err
	}
	if err := di.Value(a.container, taskRepo); err != nil {
		return err
	}
	if err := di.Value(a.container, diaryRepo); err != nil {
		return err
	}
	if err := di.Value(a.container, entryRepo); err != nil {
		return err
	}
	if err := di.Value(a.container, commentRepo); err != nil {
		return err
	}
	if err := di.Value(a.container, likeRepo); err != nil {
		return err
	}
	if err := di.Value(a.container, taskService); err != nil {
		return err
	}
	if err := di.Value(a.container, userService); err != nil {
		return err
	}
	if err := di.Value(a.container, diaryService); err != nil {
		return err
	}
	if err := di.Value(a.container, entryService); err != nil {
		return err
	}
	if err := di.Value(a.container, commentService); err != nil {
		return err
	}
	if err := di.Value(a.container, likeService); err != nil {
		return err
	}

	// 注册 AuthHelper 为单例服务
	if err := di.Singleton(a.container, func() *user.AuthHelper {
		return user.NewAuthHelper(a.config)
	}); err != nil {
		return err
	}

	slog.Info("成功注册应用层服务")
	return nil
}

// initCookie 设置Cookie默认值
// 根据配置文件设置全局Cookie的默认参数
//
// 返回:
//   - error: 设置过程中的错误信息
//
// 设置的参数:
//   - Domain: Cookie域名
//   - Path: Cookie路径（默认为"/"）
//   - ExpireDays: Cookie过期天数
//   - HttpOnly: 是否仅HTTP访问
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

// initScheduler 启动定时任务调度器
// 从DI容器获取TaskService并启动后台定时任务调度器
//
// 返回:
//   - error: 启动过程中的错误信息
//
// 功能:
//   - 获取TaskService实例
//   - 启动后台任务调度器
//   - 处理定时任务的执行
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

// initGinServer 初始化Gin HTTP服务器
// 配置Gin服务器的运行模式、中间件和路由系统
//
// 返回:
//   - error: 初始化过程中的错误信息
//
// 配置步骤:
//   1. 设置Gin运行模式（开发/生产）
//   2. 初始化自定义数据校验器
//   3. 创建Gin引擎实例
//   4. 注册全局中间件
//   5. 设置路由和静态文件服务
func (a *Application) initGinServer() error {
	// 设置 Gin 模式
	gin.SetMode(a.config.App.GinMode)

	// 设置自定义 validator 绑定到 Gin
	if err := a.initCustomValidator(); err != nil {
		return err
	}

	// 初始化默认gin server
	server := gin.New()

	// 初始化中间件
	server.Use(middleware.PanicRecovery())   // panic recovery中间件
	server.Use(middleware.AppErrorHandler()) // 应用错误处理中间件
	server.Use(middleware.CORS(a.config))         // CORS中间件
	server.Use(middleware.SharedData())           // 用户共享数据中间件
	server.Use(middleware.CSRF(a.config))         // CSRF防护中间件
	server.Use(middleware.GeneralRateLimit(a.config)) // 全局限流中间件
	slog.Info("成功初始化中间件")

	// 初始化路由
	route.SetWithContainer(server, a.container)
	slog.Info("成功初始化网站静态文件目录及路由")

	a.server = server
	return nil
}

// customValidator 自定义数据校验器
// 实现Gin框架的binding.StructValidator接口，提供结构体数据校验功能
//
// 字段:
//   - validate: go-playground/validator校验器实例
type customValidator struct {
	validate *validator.Validate
}

// ValidateStruct 实现结构体校验
// 对传入的结构体对象进行数据校验
//
// 参数:
//   - obj: 需要校验的结构体对象
//
// 返回:
//   - error: 校验失败时的错误信息
func (v *customValidator) ValidateStruct(obj interface{}) error {
	return v.validate.Struct(obj)
}

// Engine 返回底层的校验器引擎
// 提供对底层validator实例的访问
//
// 返回:
//   - interface{}: 底层的validator.Validate实例
func (v *customValidator) Engine() interface{} {
	return v.validate
}

// initCustomValidator 初始化自定义校验器并绑定到Gin
// 将项目的自定义校验器设置为Gin框架的默认校验器
//
// 返回:
//   - error: 初始化过程中的错误信息
//
// 功能:
//   - 使用项目预配置的validation.Validator
//   - 替换Gin默认的数据校验器
//   - 支持自定义校验规则和错误消息
func (a *Application) initCustomValidator() error {
	// 设置 Gin 使用项目的自定义校验器
	binding.Validator = &customValidator{
		validate: validation.Validator,
	}
	
	slog.Info("成功绑定自定义校验器到 Gin 框架")
	return nil
}
