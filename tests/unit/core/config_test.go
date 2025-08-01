package core

import (
	"os"
	"testing"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 设置测试环境变量
func setupTestEnv() {
	os.Setenv("IDIARY_ENV", "test")
	os.Setenv("IDIARY_GIN_MODE", "debug")
	os.Setenv("IDIARY_DOMAIN", "http://localhost")
	os.Setenv("IDIARY_PORT", "8080")
	os.Setenv("IDIARY_SUMMARY_MAX", "200")
	os.Setenv("IDIARY_SUFFIX", ".test")
	os.Setenv("IDIARY_PAGE_SIZE", "20")
	os.Setenv("IDIARY_MAX_DAYBOOK_LIMIT", "100")
	os.Setenv("IDIARY_HASH_SALT", "test_salt_12345")
	os.Setenv("IDIARY_STATIC_DIR", "static")
	os.Setenv("IDIARY_LOCALES_DIR", "locales")
	os.Setenv("IDIARY_DEFAULT_LANG", "zh-CN")

	// Database配置
	os.Setenv("IDIARY_DB_HOST", "localhost")
	os.Setenv("IDIARY_DB_PORT", "3306")
	os.Setenv("IDIARY_DB_NAME", "test_db")
	os.Setenv("IDIARY_DB_USER", "test_user")
	os.Setenv("IDIARY_DB_PASSWORD", "test_password")
	os.Setenv("IDIARY_DB_CONN", "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local")

	// Redis配置
	os.Setenv("IDIARY_REDIS_HOST", "localhost")
	os.Setenv("IDIARY_REDIS_PORT", "6379")
	os.Setenv("IDIARY_REDIS_PASSWORD", "")

	// Email配置
	os.Setenv("IDIARY_SMTP_SERVER", "smtp.test.com")
	os.Setenv("IDIARY_SMTP_USER", "test@test.com")
	os.Setenv("IDIARY_SMTP_PASSWORD", "test_password")
	os.Setenv("IDIARY_SMTP_PORT", "587")
	os.Setenv("IDIARY_SMTP_FROM", "test@test.com")
	os.Setenv("IDIARY_EMAIL_MAX_RETRIES", "3")
	os.Setenv("IDIARY_EMAIL_BATCH_SIZE", "10")

	// Cookie配置
	os.Setenv("IDIARY_COOKIE_DOMAIN", "localhost")
	os.Setenv("IDIARY_COOKIE_SECRET", "test_cookie_secret_12345")
	os.Setenv("IDIARY_COOKIE_EXPIRE_DAYS", "7")
	os.Setenv("IDIARY_COOKIE_HTTP_ONLY", "true")

	// Upload配置
	os.Setenv("IDIARY_UPLOAD_PATH", "uploads")
	os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")

	// Log配置
	os.Setenv("IDIARY_LOG_LEVEL", "debug")
	os.Setenv("IDIARY_LOG_FORMAT", "json")
	os.Setenv("IDIARY_LOG_PATH", "logs")
	os.Setenv("IDIARY_LOG_FILENAME", "app.log")
	os.Setenv("IDIARY_LOG_MAX_SIZE", "100")
	os.Setenv("IDIARY_LOG_MAX_AGE", "30")
	os.Setenv("IDIARY_LOG_MAX_BACKUPS", "10")
	os.Setenv("IDIARY_LOG_COMPRESS", "true")
	os.Setenv("IDIARY_LOG_LOCAL_TIME", "true")

	// JSON配置
	os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-CN","label":"简体中文","local-text":"简体中文"},{"code":"en-US","label":"English","local-text":"English"}]`)
	os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("IDIARY_PARAMS", `[{"key":"test_param","val":"test_value"}]`)
}

// 清理测试环境变量
func cleanupTestEnv() {
	envVars := []string{
		"IDIARY_ENV", "IDIARY_GIN_MODE", "IDIARY_DOMAIN", "IDIARY_PORT", "IDIARY_SUMMARY_MAX", "IDIARY_SUFFIX", "IDIARY_PAGE_SIZE",
		"IDIARY_MAX_DAYBOOK_LIMIT", "IDIARY_HASH_SALT", "IDIARY_STATIC_DIR", "IDIARY_LOCALES_DIR", "IDIARY_DEFAULT_LANG",
		"IDIARY_DB_HOST", "IDIARY_DB_PORT", "IDIARY_DB_NAME", "IDIARY_DB_USER", "IDIARY_DB_PASSWORD", "IDIARY_DB_CONN",
		"IDIARY_REDIS_HOST", "IDIARY_REDIS_PORT", "IDIARY_REDIS_PASSWORD",
		"IDIARY_SMTP_SERVER", "IDIARY_SMTP_USER", "IDIARY_SMTP_PASSWORD", "IDIARY_SMTP_PORT", "IDIARY_SMTP_FROM",
		"IDIARY_EMAIL_MAX_RETRIES", "IDIARY_EMAIL_BATCH_SIZE",
		"IDIARY_COOKIE_DOMAIN", "IDIARY_COOKIE_SECRET", "IDIARY_COOKIE_EXPIRE_DAYS", "IDIARY_COOKIE_HTTP_ONLY",
		"IDIARY_UPLOAD_PATH", "IDIARY_UPLOAD_TYPES", "IDIARY_UPLOAD_MAX_SIZE",
		"IDIARY_LOG_LEVEL", "IDIARY_LOG_FORMAT", "IDIARY_LOG_PATH", "IDIARY_LOG_FILENAME", "IDIARY_LOG_MAX_SIZE",
		"IDIARY_LOG_MAX_AGE", "IDIARY_LOG_MAX_BACKUPS", "IDIARY_LOG_COMPRESS", "IDIARY_LOG_LOCAL_TIME",
		"IDIARY_SUPPORT_LANGS", "IDIARY_TIMEZONES", "IDIARY_PARAMS",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

func TestConfig_InitConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 测试配置初始化
	err := config.InitConfig()
	require.NoError(t, err)

	// 获取配置
	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证基本配置
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "debug", cfg.App.GinMode)
	assert.Equal(t, "http://localhost", cfg.App.Domain)
	assert.Equal(t, "8080", cfg.App.Port)
}

func TestConfig_AppConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证App配置
	assert.Equal(t, "debug", cfg.App.GinMode)
	assert.Equal(t, "http://localhost", cfg.App.Domain)
	assert.Equal(t, "8080", cfg.App.Port)
	assert.Equal(t, 200, cfg.App.SummaryMax)
	assert.Equal(t, ".test", cfg.App.Suffix)
	assert.Equal(t, 20, cfg.App.PageSize)
	assert.Equal(t, 100, cfg.App.MaxDaybookLimit)
	assert.Equal(t, "static", cfg.App.StaticDir)
	assert.Equal(t, "locales", cfg.App.LocalesDir)
	assert.Equal(t, "zh-CN", cfg.App.DefaultLang)
}

func TestConfig_DatabaseConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Database配置
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 3306, cfg.Database.Port)
	assert.Equal(t, "test_db", cfg.Database.Name)
	assert.Equal(t, "test_user", cfg.Database.User)
	assert.Equal(t, "test_password", cfg.Database.Password)
	assert.NotEmpty(t, cfg.Database.Conn)
}

func TestConfig_RedisConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Redis配置
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, "", cfg.Redis.Password)
}

func TestConfig_EmailConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Email配置
	assert.Equal(t, "smtp.test.com", cfg.Email.SMTPServer)
	assert.Equal(t, "test@test.com", cfg.Email.SMTPUser)
	assert.Equal(t, "test_password", cfg.Email.SMTPPassword)
	assert.Equal(t, 587, cfg.Email.SMTPPort)
	assert.Equal(t, "test@test.com", cfg.Email.SMTPFrom)
	assert.Equal(t, 3, cfg.Email.MaxRetries)
	assert.Equal(t, 10, cfg.Email.BatchSize)
}

func TestConfig_CookieConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Cookie配置
	assert.Equal(t, "localhost", cfg.Cookie.Domain)
	assert.Equal(t, "test_cookie_secret_12345", cfg.Cookie.Secret)
	assert.Equal(t, 7, cfg.Cookie.ExpireDays)
	assert.True(t, cfg.Cookie.HttpOnly)
}

func TestConfig_UploadConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Upload配置
	assert.Equal(t, "uploads", cfg.Upload.Path)
	assert.Equal(t, "jpg,jpeg,png,gif", cfg.Upload.Types)
	assert.Equal(t, 10485760, cfg.Upload.MaxSize)
}

func TestConfig_LogConfig(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Log配置
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "logs", cfg.Log.Path)
	assert.Equal(t, "app.log", cfg.Log.Filename)
	assert.Equal(t, 100, cfg.Log.MaxSize)
	assert.Equal(t, 30, cfg.Log.MaxAge)
	assert.Equal(t, 10, cfg.Log.MaxBackups)
	assert.True(t, cfg.Log.Compress)
	assert.True(t, cfg.Log.LocalTime)
}

func TestConfig_SupportLangs(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证支持的语言
	require.Len(t, cfg.SupportLangs, 2)
	assert.Equal(t, "zh-CN", cfg.SupportLangs[0].Code)
	assert.Equal(t, "简体中文", cfg.SupportLangs[0].Label)
	assert.Equal(t, "简体中文", cfg.SupportLangs[0].LocalText)
	assert.Equal(t, "en-US", cfg.SupportLangs[1].Code)
	assert.Equal(t, "English", cfg.SupportLangs[1].Label)
	assert.Equal(t, "English", cfg.SupportLangs[1].LocalText)
}

func TestConfig_Timezones(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证时区配置
	require.Len(t, cfg.Timezones, 3)
	assert.Equal(t, "Asia/Shanghai", cfg.Timezones[0])
	assert.Equal(t, "UTC", cfg.Timezones[1])
	assert.Equal(t, "America/New_York", cfg.Timezones[2])
}

func TestConfig_Params(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证参数配置
	require.Len(t, cfg.Params, 1)
	assert.Equal(t, "test_param", cfg.Params[0].Key)
	assert.Equal(t, "test_value", cfg.Params[0].Val)
}

func TestConfig_GetParams(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	// 测试GetParams函数
	params := config.GetParams()
	require.NotNil(t, params)
	assert.Equal(t, "test_value", params["test_param"])
}

func TestConfig_GetTimeZone(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	// 测试GetTimeZone函数
	timezones := config.GetTimeZone()
	require.NotNil(t, timezones)
	require.Len(t, timezones, 3)
	assert.Equal(t, "Asia/Shanghai", timezones[0])
}

func TestConfig_GetSupportLangs(t *testing.T) {
	// 设置测试环境
	setupTestEnv()
	defer cleanupTestEnv()

	// 初始化配置
	err := config.InitConfig()
	require.NoError(t, err)

	// 测试GetSupportLangs函数
	supportLangs := config.GetSupportLangs()
	require.NotNil(t, supportLangs)
	require.Len(t, supportLangs, 2)
	assert.Equal(t, "zh-CN", supportLangs[0].Code)
	assert.Equal(t, "en-US", supportLangs[1].Code)
}

func TestConfig_MissingRequiredFields(t *testing.T) {
	// 清理环境变量
	cleanupTestEnv()

	// 测试缺少必需字段时的错误
	err := config.InitConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "配置不能为空")
}

func TestConfig_InvalidJSON(t *testing.T) {
	// 设置测试环境，但使用无效的JSON
	setupTestEnv()
	defer cleanupTestEnv()

	// 设置无效的JSON
	os.Setenv("IDIARY_SUPPORT_LANGS", "invalid json")

	// 测试无效JSON时的错误
	err := config.InitConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "解析SUPPORT_LANGS失败")
}