package testhelpers

import (
	"os"
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

	// 确保必要的目录存在
	ensureTestDirectories(t)

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
			{Code: "zh-cn", Label: "简体中文", LocalText: "简体中文"},
			{Code: "en", Label: "English", LocalText: "English"},
		},
		Timezones: []string{"Asia/Shanghai", "UTC", "America/New_York"},
		Params: []config.Param{
			{Key: "site_name", Val: "iDiary"},
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

// ensureTestDirectories 确保测试所需的目录存在
func ensureTestDirectories(t *testing.T) {
	t.Helper()

	directories := []string{
		"website/locales",
		"website/static",
		"test_uploads",
		"test_logs",
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// 创建基本的本地化文件
	createBasicLocalizationFiles(t)
}

// createBasicLocalizationFiles 创建基本的本地化文件
func createBasicLocalizationFiles(t *testing.T) {
	t.Helper()

	// 创建中文本地化文件
	zhCNContent := `{
		"nav": {
			"home": "首页"
		},
		"common": {
			"default_page_desc": "默认页面描述"
		}
	}`

	zhCNFile := "website/locales/zh-cn.json"
	if err := os.WriteFile(zhCNFile, []byte(zhCNContent), 0644); err != nil {
		t.Logf("Warning: Failed to create %s: %v", zhCNFile, err)
	}

	// 创建英文本地化文件
	enContent := `{
		"nav": {
			"home": "Home"
		},
		"common": {
			"default_page_desc": "Default page description"
		}
	}`

	enFile := "website/locales/en.json"
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Logf("Warning: Failed to create %s: %v", enFile, err)
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
