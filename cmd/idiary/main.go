// Package main 是应用程序的入口点
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/biuren/hebe/internal/app/route"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/core/middleware"
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
		// 使用logrus记录错误，Fatal会自动调用os.Exit(1)
		logrus.WithError(err).Fatal("初始化配置失败")
	}

	// 使用应用配置初始化日志
	if err := setupLogger(); err != nil {
		// 保持一致的错误处理风格
		logrus.WithError(err).Fatal("初始化日志系统失败")
	}
	logrus.Info("成功初始化日志系统")

	// 创建并启动应用
	app, err := NewApp()
	if err != nil {
		logrus.WithError(err).Fatal("初始化应用程序失败")
	}

	app.Start()
}

// App 结构体包含应用程序的核心组件
type App struct {
	config *config.Config // 应用配置
	server *gin.Engine    // Gin HTTP 服务器实例
}

// NewApp 创建并初始化一个新的应用实例
// 包括配置加载、日志设置、数据库连接、Cookie 设置和 HTTP 服务器初始化
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
	logrus.Info("数据库连接成功")

	// 设置 cookie 默认值
	pkg.SetCookieDefVal(config.Cookie.Domain, "/", config.Cookie.ExpireDays, config.Cookie.HttpOnly)
	logrus.Info("成功设置 Cookie 默认值")

	// 启动定时任务
	task.NewService().StartScheduler()
	logrus.Info("成功启动定时任务")

	// 初始化 Gin server
	server := setupGinServer(config.App.GinMode)

	return &App{
		config: config,
		server: server,
	}, nil
}

// Start 启动应用程序
// 包括HTTP服务器启动
func (app *App) Start() {
	// 启动服务器
	addr := ":" + app.config.App.Port
	logrus.Infof("Starting server on %s", addr)

	if err := app.server.Run(addr); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"addr":  addr,
		}).Fatal("服务器启动失败")
	}
}

// Shutdown 关闭应用程序
// 清理资源并优雅退出
func (app *App) Shutdown() {
	logrus.Info("开始关闭应用程序...")
	// 此处可以添加资源清理代码
	logrus.Info("应用程序已关闭")
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
			// 尝试使用 logrus，但如果 logrus 还没初始化，则使用标准输出
			stack := string(debug.Stack())

			// 检查 logrus 是否已初始化
			if logrus.StandardLogger().Out == os.Stderr {
				// logrus 还没初始化，使用标准输出
				fmt.Printf("启动过程中发生 panic: %v\n堆栈: %s\n", err, stack)
				os.Exit(1)
			} else {
				// logrus 已初始化，使用 logrus 记录
				logrus.WithFields(logrus.Fields{
					"error": err,
					"stack": stack,
				}).Fatal("启动过程中发生 panic") // 使用Fatal代替Error+os.Exit
			}
		}
	}()
}

// setupGinServer 初始化并配置 Gin HTTP 服务器
// 设置运行模式、中间件和路由
func setupGinServer(mode string) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(mode)

	// 初始化默认gin server
	server := gin.New()

	// 初始化中间件
	server.Use(middleware.PanicHandler(), middleware.ErrorHandler(), user.SharedData())
	logrus.Info("成功初始化中间件")

	// 初始化Router，创建 gin 默认 router
	// 设置模板文件，静态文件Path, 加载上下文共享数据
	route.Set(server)
	logrus.Info("成功初始化网站静态文件目录及路由")

	return server
}
