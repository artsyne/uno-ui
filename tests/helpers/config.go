package testhelpers

import (
	"testing"

	"github.com/biuren/hebe/internal/core/config"
)

// ConfigTestHelper 配置测试辅助工具
type ConfigTestHelper struct {
	originalConfig *config.Config
}

// SetupTestConfig 使用预定义配置初始化测试环境
// 基于.env.local的值，无需环境变量设置
func SetupTestConfig(t *testing.T) *ConfigTestHelper {
	t.Helper()

	// 创建测试配置（基于实际配置结构）
	testConfig := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode:         "debug",
			Domain:          "http://localhost:8080",
			Port:            "8080",
			SummaryMax:      200,
			Suffix:          "...",
			PageSize:        10,
			MaxDaybookLimit: 100,
			HashSalt:        "test_salt_12345678",
			StaticDir:       "website/static",
			LocalesDir:      "website/locales",
			DefaultLang:     "zh-cn",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "hebe_test",
			User:     "test_user",
			Password: "test_password",
			Conn:     "mysql",
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
		},
		Cookie: config.CookieConfig{
			Domain:     "localhost",
			Secret:     "test_secret_key_12345678901234567890",
			ExpireDays: 30,
			HttpOnly:   true,
		},
		Upload: config.UploadConfig{
			Path:    "test_uploads",
			Types:   "jpg,jpeg,png,gif",
			MaxSize: 10485760, // 10MB
		},
		Email: config.EmailConfig{
			SMTPServer:   "smtp.test.com",
			SMTPUser:     "test@test.com",
			SMTPPassword: "test_password",
			SMTPPort:     587,
			SMTPFrom:     "test@test.com",
			MaxRetries:   3,
			BatchSize:    10,
		},
		Log: config.LogConfig{
			Level:      "info",
			Format:     "json",
			Path:       "test_logs",
			Filename:   "test.log",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 5,
			Compress:   true,
			LocalTime:  true,
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "中文", LocalText: "中文"},
			{Code: "en", Label: "English", LocalText: "English"},
		},
		Timezones: []string{"Asia/Shanghai", "UTC", "America/New_York"},
		Params: []config.Param{
			{Key: "test_key", Val: "test_value"},
		},
	}

	// 备份原始配置
	originalConfig := config.GetConfig()

	// 设置测试配置
	config.SetConfig(testConfig)

	return &ConfigTestHelper{
		originalConfig: originalConfig,
	}
}

// Cleanup 清理测试配置，恢复原始状态
func (h *ConfigTestHelper) Cleanup(t *testing.T) {
	t.Helper()
	if h.originalConfig != nil {
		config.SetConfig(h.originalConfig)
	}
}

// GetTestConfig 获取当前测试配置
func (h *ConfigTestHelper) GetTestConfig() *config.Config {
	return config.GetConfig()
}
