// Package main 是应用程序的入口点
package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/app/route"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/module/task"
	"github.com/biuren/hebe/internal/module/user"
	"github.com/biuren/hebe/pkg"
)

// main 是应用程序的入口函数
// 负责初始化日志、panic 恢复，以及启动应用程序
func main() {
	setupPanicRecovery()

	// 设置基本日志，用于应用启动阶段
	pkg.SetupBasicLogger()

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		// 使用slog记录错误并退出
		slog.Error("初始化配置失败", "error", err)
		os.Exit(1)
	}

	// 使用应用配置初始化日志
	if err := setupLogger(); err != nil {
		// 保持一致的错误处理风格
		slog.Error("初始化日志系统失败", "error", err)
		os.Exit(1)
	}
	slog.Info("成功初始化日志系统")

	// 创建并启动应用
	app, err := NewApp()
	if err != nil {
		slog.Error("初始化应用程序失败", "error", err)
		os.Exit(1)
	}

	app.Start()
}

// App 结构体包含应用程序的核心组件
type App struct {
	config    *config.Config        // 应用配置
	server    *gin.Engine           // Gin HTTP 服务器实例
	container di.ContainerInterface // DI容器
}

// NewApp 创建并初始化一个新的应用实例
// 包括配置加载、日志设置、数据库连接、DI容器初始化和HTTP服务器初始化
// 返回初始化后的应用实例和可能的错误
func NewApp() (*App, error) {
	config := config.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("配置未初始化")
	}

	// 初始化数据库连接
	if err := initializeDatabase(); err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	slog.Info("数据库连接成功")

	// 创建DI容器并注册服务
	container := di.NewContainer()
	if err := registerServices(container); err != nil {
		return nil, fmt.Errorf("注册服务失败: %w", err)
	}
	slog.Info("成功注册DI服务")

	// 设置 cookie 默认值
	pkg.SetCookieDefVal(config.Cookie.Domain, "/", config.Cookie.ExpireDays, config.Cookie.HttpOnly)
	slog.Info("成功设置 Cookie 默认值")

	// 启动定时任务
	task.NewService().StartScheduler()
	slog.Info("成功启动定时任务")

	// 初始化 Gin server
	server, err := setupGinServer(config.App.GinMode, container)
	if err != nil {
		return nil, fmt.Errorf("初始化Gin服务器失败: %w", err)
	}

	return &App{
		config:    config,
		server:    server,
		container: container,
	}, nil
}

// registerServices 注册所有DI服务
func registerServices(container di.ContainerInterface) error {
	// 注册App层服务
	if err := di.RegisterAppServices(container); err != nil {
		return fmt.Errorf("注册App层服务失败: %w", err)
	}

	return nil
}

// setupGinServer 初始化并配置 Gin HTTP 服务器
// 设置运行模式、中间件和路由，使用DI容器注入服务
func setupGinServer(mode string, container di.ContainerInterface) (*gin.Engine, error) {
	// 设置 Gin 模式
	gin.SetMode(mode)

	// 初始化默认gin server
	server := gin.New()

	// 从DI容器获取RenderService
	renderService, err := di.Resolve[render.RenderServiceInterface](container)
	if err != nil {
		return nil, fmt.Errorf("获取RenderService失败: %w", err)
	}

	// 初始化中间件 - 使用App层的ErrorHandler
	server.Use(gin.Recovery())                         // 使用gin内置的panic recovery
	server.Use(middleware.ErrorHandler(renderService)) // 使用App层错误处理中间件
	server.Use(user.SharedData())                      // 用户共享数据中间件
	slog.Info("成功初始化中间件")

	// 初始化Router - 传递DI容器以供路由创建handlers
	route.SetWithContainer(server, container)
	slog.Info("成功初始化网站静态文件目录及路由")

	return server, nil
}

// Start 启动应用程序
// 包括HTTP服务器启动
func (app *App) Start() {
	// 启动服务器
	addr := ":" + app.config.App.Port
	slog.Info("Starting server", "addr", addr)

	if err := app.server.Run(addr); err != nil {
		slog.Error("服务器启动失败", "error", err, "addr", addr)
		os.Exit(1)
	}
}

// Shutdown 关闭应用程序
// 清理资源并优雅退出
func (app *App) Shutdown() {
	slog.Info("开始关闭应用程序...")
	// 此处可以添加资源清理代码
	slog.Info("应用程序已关闭")
}

// initializeDatabase 初始化数据库连接
func initializeDatabase() error {
	// 使用新的数据库连接方法
	if err := database.Connect(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	return nil
}

// setupLogger 初始化日志配置
// 设置日志级别和输出格式
func setupLogger() error {
	config := config.GetConfig()
	if config == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 初始化日志系统
	if err := pkg.SetLogger(&config.Log); err != nil {
		return fmt.Errorf("初始化日志系统失败: %w", err)
	}

	return nil
}

// setupPanicRecovery 设置全局 panic 恢复机制
// 捕获并记录启动过程中的 panic
func setupPanicRecovery() {
	defer func() {
		if err := recover(); err != nil {
			// 使用slog记录panic
			stack := string(debug.Stack())
			slog.Error("启动过程中发生 panic", "error", err, "stack", stack)
			os.Exit(1)
		}
	}()
}
