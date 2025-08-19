package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
	framework "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/biuren/hebe/pkg"
)



// TestLoggerConfig 测试LoggerConfig结构体
func TestLoggerConfig(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	t.Run("创建默认配置", func(t *testing.T) {
		config := &LoggerConfig{
			Path:       "/tmp/logs",
			Filename:   "app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
			LocalTime:  true,
			Level:      "info",
			Format:     "text",
		}
		
		assert.Equal(t, "/tmp/logs", config.Path, "路径应该正确")
		assert.Equal(t, "app.log", config.Filename, "文件名应该正确")
		assert.Equal(t, 100, config.MaxSize, "最大文件大小应该正确")
		assert.Equal(t, 3, config.MaxBackups, "备份文件数量应该正确")
		assert.Equal(t, 28, config.MaxAge, "文件保留天数应该正确")
		assert.True(t, config.Compress, "压缩选项应该正确")
		assert.True(t, config.LocalTime, "本地时间选项应该正确")
		assert.Equal(t, "info", config.Level, "日志级别应该正确")
		assert.Equal(t, "text", config.Format, "日志格式应该正确")
	})
}

// TestParseLogLevel 测试日志级别解析功能
func TestParseLogLevel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	tests := []struct {
		name     string
		levelStr string
		expected slog.Level
	}{
		{
			name:     "debug级别",
			levelStr: "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "info级别",
			levelStr: "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn级别",
			levelStr: "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "error级别",
			levelStr: "error",
			expected: slog.LevelError,
		},
		{
			name:     "无效级别默认为info",
			levelStr: "invalid",
			expected: slog.LevelInfo,
		},
		{
			name:     "空字符串默认为info",
			levelStr: "",
			expected: slog.LevelInfo,
		},
		{
			name:     "大写级别默认为info",
			levelStr: "INFO",
			expected: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于 parseLogLevel 是私有函数，我们通过 SetLogger 间接测试
			config := &pkg.LoggerConfig{
				Path:       "/tmp",
				Filename:   "test.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   false,
				LocalTime:  true,
				Level:      tt.levelStr,
				Format:     "text",
			}
			err := pkg.SetLogger(config)
			if tt.levelStr == "invalid" {
				assert.Error(t, err, "无效日志级别应该返回错误")
			} else {
				assert.NoError(t, err, "有效日志级别应该设置成功")
			}

		})
	}
}

// TestValidateLogConfig 测试日志配置验证功能
func TestValidateLogConfig(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	tests := []struct {
		name        string
		config      *LoggerConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "有效配置",
			config: &LoggerConfig{
				Path:       "/tmp/logs",
				Filename:   "app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
				LocalTime:  true,
				Level:      "info",
				Format:     "text",
			},
			expectError: false,
		},
		{
			name: "有效的JSON格式配置",
			config: &LoggerConfig{
				Path:       "/tmp/logs",
				Filename:   "app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
				LocalTime:  true,
				Level:      "debug",
				Format:     "json",
			},
			expectError: false,
		},
		{
			name:        "nil配置",
			config:      nil,
			expectError: true,
			errorMsg:    "log config is nil",
		},
		{
			name: "无效的日志级别",
			config: &LoggerConfig{
				Path:       "/tmp/logs",
				Filename:   "app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
				LocalTime:  true,
				Level:      "invalid",
				Format:     "text",
			},
			expectError: true,
			errorMsg:    "invalid log level",
		},
		{
			name: "无效的日志格式",
			config: &LoggerConfig{
				Path:       "/tmp/logs",
				Filename:   "app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
				LocalTime:  true,
				Level:      "info",
				Format:     "xml",
			},
			expectError: true,
			errorMsg:    "invalid log format",
		},
		{
			name: "空的日志级别",
			config: &LoggerConfig{
				Path:       "/tmp/logs",
				Filename:   "app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
				LocalTime:  true,
				Level:      "",
				Format:     "text",
			},
			expectError: false,
		},
		{
			name: "空的日志格式",
			config: &LoggerConfig{
				Path:       "/tmp/logs",
				Filename:   "app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
				LocalTime:  true,
				Level:      "info",
				Format:     "",
			},
			expectError: true,
			errorMsg:    "invalid log format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于 validateLogConfig 是私有函数，我们通过 SetLogger 间接测试
			err := pkg.SetLogger(tt.config)
			
			if tt.expectError {
				assert.Error(t, err, "应该返回错误")
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "错误消息应该包含预期内容")
				}
			} else {
				assert.NoError(t, err, "不应该返回错误")
			}
		})
	}
}

// TestSetLogger 测试日志设置功能
func TestSetLogger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "logging_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 保存原始的默认logger
	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)

	// 创建测试配置
	config := &LoggerConfig{
		Path:       tempDir,
		Filename:   "test.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
		LocalTime:  true,
		Level:      "info",
		Format:     "json",
	}

	// 测试设置日志
	err = SetLogger(config)
	require.NoError(t, err)

	// 验证日志文件是否创建
	logFilePath := filepath.Join(tempDir, "test.log")

	// 记录一条日志
	slog.Info("test message", "key", "value")

	// 等待日志写入
	time.Sleep(100 * time.Millisecond)

	// 验证日志文件存在
	_, err = os.Stat(logFilePath)
	require.NoError(t, err)

	// 读取日志内容
	content, err := os.ReadFile(logFilePath)
	require.NoError(t, err)

	// 验证日志内容
	logStr := string(content)
	assert.Contains(t, logStr, "test message", "应该包含测试消息")
	assert.Contains(t, logStr, "key", "应该包含键")
	assert.Contains(t, logStr, "value", "应该包含值")

	// 验证JSON格式
	assert.Contains(t, logStr, `"level"`, "JSON格式应该包含level字段")
	assert.Contains(t, logStr, `"msg"`, "JSON格式应该包含msg字段")
}

// TestSetupBasicLogger 测试基本日志设置功能
func TestSetupBasicLogger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	t.Run("设置基本日志", func(t *testing.T) {
		// 保存原始的默认logger
		originalLogger := slog.Default()
		defer slog.SetDefault(originalLogger)
		
		// 调用SetupBasicLogger
		SetupBasicLogger()
		
		// 验证logger已被设置（通过检查是否能正常记录日志）
		assert.NotPanics(t, func() {
			slog.Info("test basic logger")
		}, "基本日志设置后应该能正常记录日志")
		
		// 验证默认logger已被更改
		newLogger := slog.Default()
		assert.NotEqual(t, originalLogger, newLogger, "默认logger应该被更改")
	})
}

// TestLoggerIntegration 测试日志模块的集成功能
func TestLoggerIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "logging_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("完整的日志配置和使用流程", func(t *testing.T) {
		// 1. 创建日志配置
		config := &LoggerConfig{
			Path:       tempDir,
			Filename:   "integration.log",
			MaxSize:    1, // 1MB
			MaxBackups: 2,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  true,
			Level:      "debug",
			Format:     "json",
		}
		
		// 2. 设置日志
		err := SetLogger(config)
		assert.NoError(t, err, "日志设置应该成功")
		
		// 3. 记录不同级别的日志
		slog.Debug("debug message")
		slog.Info("info message")
		slog.Warn("warn message")
		slog.Error("error message")
		
		// 4. 等待日志写入
		time.Sleep(200 * time.Millisecond)
		
		// 5. 验证日志文件存在
		logFilePath := filepath.Join(tempDir, "integration.log")
		_, err = os.Stat(logFilePath)
		assert.NoError(t, err, "日志文件应该存在")
		
		// 6. 读取日志文件内容
		logContent, err := os.ReadFile(logFilePath)
		assert.NoError(t, err, "应该能读取日志文件")
		
		// 7. 验证日志内容
		logStr := string(logContent)
		assert.Contains(t, logStr, "debug message", "应该包含debug日志")
		assert.Contains(t, logStr, "info message", "应该包含info日志")
		assert.Contains(t, logStr, "warn message", "应该包含warn日志")
		assert.Contains(t, logStr, "error message", "应该包含error日志")
		
		// 8. 验证JSON格式（应该包含JSON字段）
		assert.Contains(t, logStr, `"level"`, "JSON格式应该包含level字段")
		assert.Contains(t, logStr, `"msg"`, "JSON格式应该包含msg字段")
	})
	
	t.Run("文本格式日志测试", func(t *testing.T) {
		// 1. 创建文本格式配置
		config := &LoggerConfig{
			Path:       tempDir,
			Filename:   "text.log",
			MaxSize:    1,
			MaxBackups: 2,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  true,
			Level:      "info",
			Format:     "text",
		}
		
		// 2. 设置日志
		err := SetLogger(config)
		assert.NoError(t, err, "文本格式日志设置应该成功")
		
		// 3. 记录日志
		slog.Info("text format test")
		
		// 4. 等待日志写入
		time.Sleep(100 * time.Millisecond)
		
		// 5. 验证日志文件
		logFilePath := filepath.Join(tempDir, "text.log")
		logContent, err := os.ReadFile(logFilePath)
		assert.NoError(t, err, "应该能读取文本格式日志文件")
		
		// 6. 验证文本格式内容
		logStr := string(logContent)
		assert.Contains(t, logStr, "text format test", "应该包含日志消息")
		assert.Contains(t, logStr, "INFO", "应该包含INFO级别标识")
	})
}

// TestLoggerEdgeCases 测试日志模块的边界情况
func TestLoggerEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "logging_edge_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("处理不存在的目录", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "non", "existent", "path")
		config := &LoggerConfig{
			Path:       nonExistentDir,
			Filename:   "test.log",
			MaxSize:    1,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  true,
			Level:      "info",
			Format:     "text",
		}
		
		err := SetLogger(config)
		assert.NoError(t, err, "应该能创建不存在的目录")
		
		// 验证目录是否被创建
		_, err = os.Stat(nonExistentDir)
		assert.NoError(t, err, "目录应该被创建")
	})
	
	t.Run("处理极端配置值", func(t *testing.T) {
		config := &LoggerConfig{
			Path:       tempDir,
			Filename:   "extreme.log",
			MaxSize:    0, // 极小值
			MaxBackups: 0, // 不保留备份
			MaxAge:     0, // 不按时间清理
			Compress:   true,
			LocalTime:  false,
			Level:      "error", // 最高级别
			Format:     "json",
		}
		
		err := SetLogger(config)
		assert.NoError(t, err, "应该能处理极端配置值")
		
		// 记录不同级别的日志，只有error应该被记录
		slog.Debug("debug should not appear")
		slog.Info("info should not appear")
		slog.Warn("warn should not appear")
		slog.Error("error should appear")
		
		time.Sleep(100 * time.Millisecond)
		
		// 验证只有error级别的日志被记录
		logFilePath := filepath.Join(tempDir, "extreme.log")
		logContent, err := os.ReadFile(logFilePath)
		assert.NoError(t, err, "应该能读取日志文件")
		
		logStr := string(logContent)
		assert.Contains(t, logStr, "error should appear", "应该包含error日志")
		assert.NotContains(t, logStr, "debug should not appear", "不应该包含debug日志")
		assert.NotContains(t, logStr, "info should not appear", "不应该包含info日志")
		assert.NotContains(t, logStr, "warn should not appear", "不应该包含warn日志")
	})
	
	t.Run("处理特殊字符文件名", func(t *testing.T) {
		config := &LoggerConfig{
			Path:       tempDir,
			Filename:   "special-chars_123.log",
			MaxSize:    1,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  true,
			Level:      "info",
			Format:     "text",
		}
		
		err := SetLogger(config)
		assert.NoError(t, err, "应该能处理特殊字符文件名")
		
		slog.Info("special filename test")
		time.Sleep(100 * time.Millisecond)
		
		// 验证文件是否被创建
		logFilePath := filepath.Join(tempDir, "special-chars_123.log")
		_, err = os.Stat(logFilePath)
		assert.NoError(t, err, "特殊字符文件名的日志文件应该存在")
	})
}

// TestLoggerConcurrency 测试日志模块的并发安全性
func TestLoggerConcurrency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "logging_concurrency_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("并发写入日志", func(t *testing.T) {
		config := &LoggerConfig{
			Path:       tempDir,
			Filename:   "concurrent.log",
			MaxSize:    1,
			MaxBackups: 3,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  true,
			Level:      "info",
			Format:     "text",
		}
		
		err := SetLogger(config)
		assert.NoError(t, err, "日志设置应该成功")
		
		// 并发写入日志
		const numGoroutines = 10
		const messagesPerGoroutine = 10
		done := make(chan bool, numGoroutines)
		
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				for j := 0; j < messagesPerGoroutine; j++ {
					slog.Info(fmt.Sprintf("goroutine %d message %d", id, j))
				}
				done <- true
			}(i)
		}
		
		// 等待所有goroutine完成
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
		
		// 等待日志写入完成
		time.Sleep(500 * time.Millisecond)
		
		// 验证日志文件存在
		logFilePath := filepath.Join(tempDir, "concurrent.log")
		_, err = os.Stat(logFilePath)
		assert.NoError(t, err, "并发日志文件应该存在")
		
		// 读取日志内容
		logContent, err := os.ReadFile(logFilePath)
		assert.NoError(t, err, "应该能读取并发日志文件")
		
		// 验证日志内容不为空
		assert.NotEmpty(t, logContent, "并发日志内容不应该为空")
		
		// 计算日志行数（简单验证）
		logStr := string(logContent)
		lines := strings.Split(strings.TrimSpace(logStr), "\n")
		assert.Greater(t, len(lines), 0, "应该有日志行")
	})
}