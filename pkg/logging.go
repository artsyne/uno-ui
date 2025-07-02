package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetLogger 配置日志
// 设置日志输出位置、格式和轮转策略
func SetLogger(config *config.LogConfig) error {
	if config == nil {
		return fmt.Errorf("log configuration is nil")
	}

	// 验证日志配置
	if err := validateLogConfig(config); err != nil {
		return fmt.Errorf("invalid log configuration: %w", err)
	}

	// 确保日志目录存在
	logPath := filepath.Join(config.Path)
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
    }

	// 配置日志轮转
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, config.Filename),
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // 文件个数
		MaxAge:     config.MaxAge,     // 天数
		Compress:   config.Compress,   // 压缩
		LocalTime:  config.LocalTime,  // 使用本地时间
	}

	// 同时输出到文件和控制台
	writers := []io.Writer{logFile, os.Stdout}
	mw := io.MultiWriter(writers...)
	logrus.SetOutput(mw)

	// 设置日志格式
	if config.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}
	logrus.SetLevel(level)

	return nil
}

// validateLogConfig 验证日志配置的有效性
func validateLogConfig(config *config.LogConfig) error {
	if config == nil {
		return fmt.Errorf("log config is nil")
	}

	// 验证日志级别
	if _, err := logrus.ParseLevel(config.Level); err != nil {
		return fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

	// 验证日志格式
	if config.Format != "text" && config.Format != "json" {
		return fmt.Errorf("invalid log format %q: must be 'text' or 'json'", config.Format)
	}

	return nil
}

// SetupBasicLogger 设置基本日志，用于配置加载前
func SetupBasicLogger() {
	// 基本控制台日志设置
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
}
