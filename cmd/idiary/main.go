// Package main 是应用程序的入口点
package main

import (
	"log/slog"
	"os"

	"github.com/biuren/hebe/internal/app/startup"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/pkg"
)

// main 是应用程序的入口函数
// 负责初始化日志、panic 恢复，以及启动应用程序
func main() {
	// 1. 设置基本日志，用于应用启动阶段
	pkg.SetupBasicLogger()

	// 2. 初始化配置
	if err := config.InitConfig(); err != nil {
		slog.Error("初始化配置失败", "error", err)
		os.Exit(1)
	}

	// 3. 配置完整的日志系统（包含文件输出）
	logConfig := &pkg.LoggerConfig{
		Path:       config.GetConfig().Log.Path,
		Filename:   config.GetConfig().Log.Filename,
		MaxSize:    config.GetConfig().Log.MaxSize,
		MaxBackups: config.GetConfig().Log.MaxBackups,
		MaxAge:     config.GetConfig().Log.MaxAge,
		Compress:   config.GetConfig().Log.Compress,
		LocalTime:  config.GetConfig().Log.LocalTime,
		Level:      config.GetConfig().Log.Level,
		Format:     config.GetConfig().Log.Format,
	}
	if err := pkg.SetLogger(logConfig); err != nil {
		slog.Error("日志系统初始化失败", "error", err)
		os.Exit(1)
	}

	slog.Info("应用启动", "version", "1.0.0")

	// 4. 创建并启动应用
	app := bootstrap.NewApplication()
	if err := app.Bootstrap(); err != nil {
		slog.Error("应用启动失败", "error", err)
		os.Exit(1)
	}

	// 5. 运行应用
	if err := app.Run(); err != nil {
		slog.Error("应用运行失败", "error", err)
		os.Exit(1)
	}
}
