package pkg

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerConfig 独立的日志配置结构
type LoggerConfig struct {
	Path       string // 日志文件路径
	Filename   string // 日志文件名
	MaxSize    int    // 单个日志文件最大尺寸，单位MB
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
	LocalTime  bool   // 是否使用本地时间
	Level      string // 日志级别: debug, info, warn, error
	Format     string // 日志格式: text, json
}

// SetLogger 配置slog日志
// 设置日志输出位置、格式和轮转策略
func SetLogger(config *LoggerConfig) error {
	if config == nil {
		return fmt.Errorf("log config is nil")
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

	// 解析日志级别
	level := parseLogLevel(config.Level)

	// 创建handler
	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(mw, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(mw, &slog.HandlerOptions{
			Level: level,
		})
	}

	// 设置默认logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// normalizeLogLevel 规范化日志级别字符串
func normalizeLogLevel(levelStr string) string {
	// 只处理大写情况，其他保持原样
	switch levelStr {
	case "DEBUG":
		return "debug"
	case "INFO":
		return "info"
	case "WARN", "WARNING":
		return "warn"
	case "ERROR":
		return "error"
	default:
		// 对于其他情况（包括小写的有效级别和无效级别），保持原样
		return levelStr
	}
}

// parseLogLevel 将字符串日志级别转换为slog.Level
func parseLogLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // 默认级别
	}
}

// validateLogConfig 验证日志配置的有效性
func validateLogConfig(config *LoggerConfig) error {
	if config == nil {
		return fmt.Errorf("log config is nil")
	}

	// 处理空字符串级别，默认为 info
	if config.Level == "" {
		config.Level = "info"
	}

	// 规范化日志级别（处理大写）
	normalizedLevel := normalizeLogLevel(config.Level)
	config.Level = normalizedLevel

	// 验证日志级别
	validLevels := []string{"debug", "info", "warn", "error"}
	isValidLevel := false
	for _, level := range validLevels {
		if config.Level == level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		return fmt.Errorf("invalid log level %q: must be one of %v", config.Level, validLevels)
	}

	// 验证日志格式
	if config.Format != "text" && config.Format != "json" {
		return fmt.Errorf("invalid log format %q: must be 'text' or 'json'", config.Format)
	}

	return nil
}

// SetupBasicLogger 设置基本slog日志，用于配置加载前
func SetupBasicLogger() {
	// 基本控制台日志设置
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
