package testhelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/biuren/hebe/internal/core/config"
)

// ConfigTestHelper 配置测试辅助工具
type ConfigTestHelper struct {
	originalConfig *config.Config
}

// SetupTestConfig 使用.env.local文件初始化测试环境
// 直接使用项目根目录的.env.local配置，但修改数据库名为测试数据库
func SetupTestConfig(t *testing.T) *ConfigTestHelper {
	t.Helper()

	// 确保必要的目录存在
	ensureTestDirectories(t)

	// 设置环境变量指向.env.local文件
	projectRoot := getProjectRoot()
	envFile := filepath.Join(projectRoot, ".env.local")
	os.Setenv("ENV_FILE", envFile)

	// 设置必需的环境变量
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
	
	// 数据库配置
	os.Setenv("IDIARY_DB_HOST", "127.0.0.1")
	os.Setenv("IDIARY_DB_PORT", "3306")
	os.Setenv("IDIARY_DB_NAME", "hebe_test")
	os.Setenv("IDIARY_DB_USER", "root")
	os.Setenv("IDIARY_DB_PASSWORD", "123123")
	os.Setenv("IDIARY_DB_CONN", "%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_general_ci")
	
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
	
	// Cookie配置
	os.Setenv("IDIARY_COOKIE_DOMAIN", "localhost")
	os.Setenv("IDIARY_COOKIE_SECRET", "test_secret_key_12345678901234567890")
	os.Setenv("IDIARY_COOKIE_EXPIRE_DAYS", "30")
	os.Setenv("IDIARY_COOKIE_HTTP_ONLY", "true")
	
	// 上传配置
	os.Setenv("IDIARY_UPLOAD_PATH", "uploads")
	os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")
	
	// 日志配置
	os.Setenv("IDIARY_LOG_LEVEL", "info")
	os.Setenv("IDIARY_LOG_FORMAT", "json")
	os.Setenv("IDIARY_LOG_PATH", "logs")
	os.Setenv("IDIARY_LOG_FILENAME", "test.log")
	os.Setenv("IDIARY_LOG_MAX_SIZE", "100")
	os.Setenv("IDIARY_LOG_MAX_BACKUPS", "3")
	os.Setenv("IDIARY_LOG_MAX_AGE", "7")
	os.Setenv("IDIARY_LOG_COMPRESS", "true")
	os.Setenv("IDIARY_LOG_LOCAL_TIME", "true")
	
	// 设置必需的JSON格式环境变量（使用IDIARY_前缀）
	os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-cn","label":"中文","local-text":"中文"},{"code":"en","label":"English","local-text":"English"}]`)
	os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("IDIARY_PARAMS", `[{"key":"test_key","val":"test_value"}]`)

	// 备份原始配置
	originalConfig := config.GetConfig()

	// 初始化配置（会自动读取.env.local）
	if err := config.InitConfig(); err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 获取配置并修改数据库名为测试数据库
	cfg := config.GetConfig()
	if cfg != nil {
		// 修改数据库名为测试数据库
		cfg.Database.Name = "hebe_test"
	}

	return &ConfigTestHelper{
		originalConfig: originalConfig,
	}
}

// ensureTestDirectories 确保测试所需的目录存在
func ensureTestDirectories(t *testing.T) {
	t.Helper()

	projectRoot := getProjectRoot()
	directories := []string{
		filepath.Join(projectRoot, "website/locales"),
		filepath.Join(projectRoot, "website/static"),
		filepath.Join(projectRoot, "test_uploads"),
		filepath.Join(projectRoot, "test_logs"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Logf("Warning: Failed to create test directory %s: %v", dir, err)
		}
	}

	// 创建基本的本地化文件（如果不存在）
	createBasicLocalizationFiles(t)
}

// createBasicLocalizationFiles 创建基本的本地化文件（如果不存在）
func createBasicLocalizationFiles(t *testing.T) {
	t.Helper()

	projectRoot := getProjectRoot()

	// 检查并创建中文本地化文件
	zhCNFile := filepath.Join(projectRoot, "website/locales/zh-cn.json")
	if _, err := os.Stat(zhCNFile); os.IsNotExist(err) {
		zhCNContent := `{
	"nav": {
		"home": "首页"
	},
	"common": {
		"default_page_desc": "默认页面描述"
	}
}`
		if err := os.WriteFile(zhCNFile, []byte(zhCNContent), 0644); err != nil {
			t.Logf("Warning: Failed to create %s: %v", zhCNFile, err)
		}
	}

	// 检查并创建英文本地化文件
	enFile := filepath.Join(projectRoot, "website/locales/en.json")
	if _, err := os.Stat(enFile); os.IsNotExist(err) {
		enContent := `{
	"nav": {
		"home": "Home"
	},
	"common": {
		"default_page_desc": "Default page description"
	}
}`
		if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
			t.Logf("Warning: Failed to create %s: %v", enFile, err)
		}
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
