package testhelpers

import (
	"os"
	"sync"
	"testing"

	"github.com/biuren/hebe/internal/utils/config"
)

// ConfigManager 配置管理器单例
type ConfigManager struct {
	config *config.Config
	mu     sync.RWMutex
}

var (
	configManager     *ConfigManager
	configManagerOnce sync.Once
)

// GetConfigManager 获取配置管理器单例
func GetConfigManager() *ConfigManager {
	configManagerOnce.Do(func() {
		configManager = &ConfigManager{}
	})
	return configManager
}

// InitializeConfig 初始化配置（线程安全）
func (cm *ConfigManager) InitializeConfig(t *testing.T) *config.Config {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 如果配置已经初始化，直接返回
	if cm.config != nil {
		return cm.config
	}

	// 强制设置测试环境变量（绕过sync.Once限制）
	forceSetupTestEnvironment()

	// 验证关键环境变量是否已设置
	if supportLangs := os.Getenv("IDIARY_SUPPORT_LANGS"); supportLangs == "" {
		if t != nil {
			t.Fatalf("IDIARY_SUPPORT_LANGS环境变量未设置，当前值: %q", supportLangs)
		}
		panic("IDIARY_SUPPORT_LANGS环境变量未设置")
	}

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		if t != nil {
			t.Fatalf("配置初始化失败: %v", err)
		}
		panic("配置初始化失败: " + err.Error())
	}

	// 获取并缓存配置
	cm.config = config.GetConfig()
	if cm.config != nil {
		// 确保使用测试数据库
		cm.config.Database.Name = "hebe_test"
	}

	return cm.config
}

// GetConfig 获取配置（线程安全）
func (cm *ConfigManager) GetConfig() *config.Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// ResetConfig 重置配置（主要用于测试）
func (cm *ConfigManager) ResetConfig() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = nil
}

// SetupTestConfigSafe 安全的测试配置设置
func SetupTestConfigSafe(t *testing.T) *config.Config {
	t.Helper()
	return GetConfigManager().InitializeConfig(t)
}

// forceSetupTestEnvironment 强制设置测试环境变量（绕过sync.Once限制）
func forceSetupTestEnvironment() {
	// 直接设置环境变量，不使用sync.Once
	os.Setenv("IDIARY_ENV", "test")
	os.Setenv("IDIARY_GIN_MODE", "debug")
	os.Setenv("IDIARY_DOMAIN", "http://localhost:8080")
	os.Setenv("IDIARY_PORT", "8080")
	os.Setenv("IDIARY_SUMMARY_MAX", "200")
	os.Setenv("IDIARY_SUFFIX", "...")
	os.Setenv("IDIARY_PAGE_SIZE", "10")
	os.Setenv("IDIARY_MAX_DAYBOOK_LIMIT", "100")
	os.Setenv("IDIARY_HASH_SALT", "test_salt_12345")
	os.Setenv("IDIARY_STATIC_DIR", "website/static")
	os.Setenv("IDIARY_LOCALES_DIR", "website/locales")
	os.Setenv("IDIARY_DEFAULT_LANG", "zh-cn")
	os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-cn","label":"中文","local-text":"中文"},{"code":"en","label":"English","local-text":"English"}]`)
	os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("IDIARY_PARAMS", `[{"key":"test_key","val":"test_value"}]`)
	
	// 数据库配置
	os.Setenv("IDIARY_DB_HOST", "127.0.0.1")
	os.Setenv("IDIARY_DB_PORT", "3306")
	os.Setenv("IDIARY_DB_NAME", "hebe_test")
	os.Setenv("IDIARY_DB_USER", "root")
	os.Setenv("IDIARY_DB_PASSWORD", "123123")
	os.Setenv("IDIARY_DB_CONN", "%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_general_ci")
	
	// Cookie配置
	os.Setenv("IDIARY_COOKIE_DOMAIN", "localhost")
	os.Setenv("IDIARY_COOKIE_SECRET", "test_secret_key_12345678901234567890")
	os.Setenv("IDIARY_COOKIE_EXPIRE_DAYS", "30")
	os.Setenv("IDIARY_COOKIE_HTTP_ONLY", "true")
	
	// 上传配置
	os.Setenv("IDIARY_UPLOAD_PATH", "tests/test_uploads")
	os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")
	
	// Redis配置
	os.Setenv("IDIARY_REDIS_HOST", "127.0.0.1")
	os.Setenv("IDIARY_REDIS_PORT", "6379")
	os.Setenv("IDIARY_REDIS_PASSWORD", "woaibeijingtiananmeng")
	
	// 邮件配置
	os.Setenv("IDIARY_SMTP_SERVER", "smtp.test.com")
	os.Setenv("IDIARY_SMTP_USER", "test@test.com")
	os.Setenv("IDIARY_SMTP_PASSWORD", "test_password")
	os.Setenv("IDIARY_SMTP_PORT", "587")
	os.Setenv("IDIARY_SMTP_FROM", "test@test.com")
	os.Setenv("IDIARY_EMAIL_MAX_RETRIES", "3")
	os.Setenv("IDIARY_EMAIL_BATCH_SIZE", "10")
	
	// 日志配置
	os.Setenv("IDIARY_LOG_LEVEL", "info")
	os.Setenv("IDIARY_LOG_FORMAT", "json")
	os.Setenv("IDIARY_LOG_PATH", "tests/test_logs")
	os.Setenv("IDIARY_LOG_FILENAME", "test.log")
	os.Setenv("IDIARY_LOG_MAX_SIZE", "100")
	os.Setenv("IDIARY_LOG_MAX_BACKUPS", "3")
	os.Setenv("IDIARY_LOG_MAX_AGE", "7")
	os.Setenv("IDIARY_LOG_COMPRESS", "true")
	os.Setenv("IDIARY_LOG_LOCAL_TIME", "true")
}

// GetTestConfigSafe 安全的获取测试配置
func GetTestConfigSafe() *config.Config {
	cm := GetConfigManager()
	cfg := cm.GetConfig()
	if cfg == nil {
		// 如果配置未初始化，尝试初始化（不传入t，使用panic处理错误）
		cfg = cm.InitializeConfig(nil)
	}
	return cfg
}