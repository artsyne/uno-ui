package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/biuren/hebe/internal/utils/config"
	. "github.com/biuren/hebe/tests/framework/types"
)

var (
	globalTestConfig   *TestConfig
	globalConfig       *config.Config
	testConfigMutex    sync.RWMutex
	globalConfigMutex  sync.RWMutex
)

// ResetGlobalConfig 重置全局配置单例状态（仅用于测试）
func ResetGlobalConfig() {
	globalConfigMutex.Lock()
	defer globalConfigMutex.Unlock()
	globalConfig = nil
}

// ResetTestConfig 重置测试配置单例状态（仅用于测试）
func ResetTestConfig() {
	testConfigMutex.Lock()
	defer testConfigMutex.Unlock()
	globalTestConfig = nil
}

// GetTestConfig 获取测试配置
// 单例模式，确保配置只初始化一次，替代当前190次重复的环境变量设置
func GetTestConfig() *TestConfig {
	testConfigMutex.RLock()
	if globalTestConfig != nil {
		defer testConfigMutex.RUnlock()
		return globalTestConfig
	}
	testConfigMutex.RUnlock()

	testConfigMutex.Lock()
	defer testConfigMutex.Unlock()
	if globalTestConfig == nil {
		globalTestConfig = loadTestConfig()
	}
	return globalTestConfig
}

// GetGlobalConfig 获取全局配置对象
// 使用 config 包中定义的标准配置结构体
func GetGlobalConfig() *config.Config {
	globalConfigMutex.RLock()
	if globalConfig != nil {
		defer globalConfigMutex.RUnlock()
		return globalConfig
	}
	globalConfigMutex.RUnlock()

	globalConfigMutex.Lock()
	defer globalConfigMutex.Unlock()
	if globalConfig == nil {
		globalConfig = loadGlobalConfig()
	}
	return globalConfig
}

// loadGlobalConfig 加载全局配置对象
// 使用 config 包中定义的标准配置结构体
func loadGlobalConfig() *config.Config {
	// 首先设置测试环境变量
	setupTestEnvironment()

	return &config.Config{
		Env: getEnvString("IDIARY_ENV", "test"),
		App: config.AppConfig{
			GinMode:         getEnvString("IDIARY_GIN_MODE", "debug"),
			Domain:          getEnvString("IDIARY_DOMAIN", "http://127.0.0.1:8080"),
			Port:            getEnvString("IDIARY_PORT", "8080"),
			SummaryMax:      getEnvInt("IDIARY_SUMMARY_MAX", 555),
			Suffix:          getEnvString("IDIARY_SUFFIX", "..."),
			PageSize:        getEnvInt("IDIARY_PAGE_SIZE", 20),
			MaxDaybookLimit: getEnvInt("IDIARY_MAX_DAYBOOK_LIMIT", 500),
			StaticDir:       getEnvString("IDIARY_STATIC_DIR", filepath.Join(GetProjectRoot(), "website/static")),
			LocalesDir:      getEnvString("IDIARY_LOCALES_DIR", filepath.Join(GetProjectRoot(), "website/locales")),
			DefaultLang:     getEnvString("IDIARY_DEFAULT_LANG", "zh-cn"),
		},
		Database: config.DatabaseConfig{
			Host:     getEnvString("IDIARY_DB_HOST", "127.0.0.1"),
			Port:     getEnvInt("IDIARY_DB_PORT", 3306),
			Name:     getEnvString("IDIARY_DB_NAME", "hebe_test"),
			User:     getEnvString("IDIARY_DB_USER", "root"),
			Password: getEnvString("IDIARY_DB_PASSWORD", "123123"),
			Conn:     getEnvString("IDIARY_DB_CONN", "%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_general_ci"),
	},
	User: config.UserConfig{
		PwdMinLen:        getEnvInt("IDIARY_PWD_MIN_LEN", 6),
		PwdMaxLen:        getEnvInt("IDIARY_PWD_MAX_LEN", 50),
		DefaultGroup:     int8(getEnvInt("IDIARY_USER_DEFAULT_GROUP", 1)),
		StatusOK:         int8(getEnvInt("IDIARY_USER_STATUS_OK", 1)),
		StatusWaitActive: int8(getEnvInt("IDIARY_USER_STATUS_WAIT_ACTIVE", 0)),
		StatusDeleted:    int8(getEnvInt("IDIARY_USER_STATUS_DEL", 2)),
	},
	Redis: config.RedisConfig{
			Host:     getEnvString("IDIARY_REDIS_HOST", "127.0.0.1"),
			Port:     getEnvInt("IDIARY_REDIS_PORT", 6379),
			Password: getEnvString("IDIARY_REDIS_PASSWORD", ""),
		},
		Email: config.EmailConfig{
			SMTPServer:   getEnvString("IDIARY_SMTP_SERVER", "smtp.test.com"),
			SMTPUser:     getEnvString("IDIARY_SMTP_USER", "test@test.com"),
			SMTPPassword: getEnvString("IDIARY_SMTP_PASSWORD", "test_password"),
			SMTPPort:     getEnvInt("IDIARY_SMTP_PORT", 25),
			SMTPFrom:     getEnvString("IDIARY_SMTP_FROM", "test@test.com"),
			MaxRetries:   getEnvInt("IDIARY_EMAIL_MAX_RETRIES", 3),
			BatchSize:    getEnvInt("IDIARY_EMAIL_BATCH_SIZE", 10),
		},
		Cookie: config.CookieConfig{
			Domain:     getEnvString("IDIARY_COOKIE_DOMAIN", "127.0.0.1"),
			Secret:     getEnvString("IDIARY_COOKIE_SECRET", "test_secret_key_12345678901234567890"),
			ExpireDays: getEnvInt("IDIARY_COOKIE_EXPIRE_DAYS", 7),
			HttpOnly:   getEnvBool("IDIARY_COOKIE_HTTP_ONLY", true),
		},
		Upload: config.UploadConfig{
			Path:    getEnvString("IDIARY_UPLOAD_PATH", "/tmp/test_uploads"),
			Types:   getEnvString("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif,pdf,txt"),
			MaxSize: getEnvInt("IDIARY_UPLOAD_MAX_SIZE", 10485760),
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "简体中文", LocalText: "简体中文"},
			{Code: "en", Label: "English", LocalText: "English"},
		},
		Timezones: []string{"Asia/Shanghai", "UTC"},
		Params: []config.Param{
			{Key: "test_param", Val: "test_value"},
		},
		Log: config.LogConfig{
			Level:      getEnvString("IDIARY_LOG_LEVEL", "debug"),
			Format:     getEnvString("IDIARY_LOG_FORMAT", "text"),
			Path:       getEnvString("IDIARY_LOG_PATH", "/tmp/test_logs"),
			Filename:   getEnvString("IDIARY_LOG_FILENAME", "test.log"),
			MaxSize:    getEnvInt("IDIARY_LOG_MAX_SIZE", 100),
			MaxAge:     getEnvInt("IDIARY_LOG_MAX_AGE", 30),
			MaxBackups: getEnvInt("IDIARY_LOG_MAX_BACKUPS", 5),
			Compress:   getEnvBool("IDIARY_LOG_COMPRESS", true),
			LocalTime:  getEnvBool("IDIARY_LOG_LOCAL_TIME", true),
		},
	}
}

// loadTestConfig 加载测试配置
// 统一从环境变量加载，消除分散在4个不同文件中的配置设置
func loadTestConfig() *TestConfig {
	// 首先设置测试环境变量
	setupTestEnvironment()

	config := &TestConfig{
		App: config.AppConfig{
			GinMode:         getEnvString("IDIARY_GIN_MODE", "debug"),
			Domain:          getEnvString("IDIARY_DOMAIN", "http://127.0.0.1:8080"),
			Port:            getEnvString("IDIARY_PORT", "8080"),
			SummaryMax:      getEnvInt("IDIARY_SUMMARY_MAX", 555),
			Suffix:          getEnvString("IDIARY_SUFFIX", "..."),
			PageSize:        getEnvInt("IDIARY_PAGE_SIZE", 20),
			MaxDaybookLimit: getEnvInt("IDIARY_MAX_DAYBOOK_LIMIT", 500),
			StaticDir:       getEnvString("IDIARY_STATIC_DIR", "website/static"),
			LocalesDir:      getEnvString("IDIARY_LOCALES_DIR", "website/locales"),
			DefaultLang:     getEnvString("IDIARY_DEFAULT_LANG", "zh-cn"),
		},
		Database: DatabaseConfig{
			Host:     getEnvString("IDIARY_DB_HOST", "127.0.0.1"),
			Port:     getEnvInt("IDIARY_DB_PORT", 3306),
			Name:     getEnvString("IDIARY_DB_NAME", "hebe_test"),
			User:     getEnvString("IDIARY_DB_USER", "root"),
			Password: getEnvString("IDIARY_DB_PASSWORD", "123123"),
			Charset:  getEnvString("IDIARY_DB_CHARSET", "utf8mb4"),
			MaxConns: getEnvInt("IDIARY_DB_MAX_CONNS", 10),
			MaxIdle:  getEnvInt("IDIARY_DB_MAX_IDLE", 5),
		},
		Cleanup: CleanupConfig{
			Enabled:     ShouldCleanupDatabase(), // 使用新的flag解析逻辑
			Strategy:    getEnvString("IDIARY_TEST_CLEANUP_STRATEGY", "prefix"),
			Retention:   getEnvDuration("IDIARY_TEST_CLEANUP_RETENTION", 24*time.Hour),
			BatchSize:   getEnvInt("IDIARY_TEST_CLEANUP_BATCH_SIZE", 100),
			TablesOrder: getDefaultTableOrder(),
		},
		Isolation: IsolationConfig{
			Strategy:     getEnvString("IDIARY_TEST_ISOLATION_STRATEGY", "prefix"),
			PrefixFormat: getEnvString("IDIARY_TEST_PREFIX_FORMAT", "test_{date}_{time}_{random}"),
			MaxPrefixLen: getEnvInt("IDIARY_TEST_MAX_PREFIX_LEN", 20),
		},
		Performance: PerformanceConfig{
			Parallel:      ShouldRunParallel(), // 使用新的flag解析逻辑
			Timeout:       getEnvDuration("IDIARY_TEST_TIMEOUT", 30*time.Second),
			MaxRetries:    getEnvInt("IDIARY_TEST_MAX_RETRIES", 3),
			SkipSlowTests: getEnvBool("IDIARY_TEST_SKIP_SLOW", false),
		},
	}

	// 验证配置
	if err := validateTestConfig(config); err != nil {
		panic(fmt.Sprintf("无效的测试配置: %v", err))
	}

	return config
}

// setupTestEnvironment 统一设置测试环境变量
// 替代当前分散在多个文件中的环境变量设置，基于实际的.env.local配置
func setupTestEnvironment() {
	// 基本环境配置
	setEnvIfNotExists("IDIARY_ENV", "test")
	setEnvIfNotExists("IDIARY_GIN_MODE", "debug")

	// 应用配置
	setEnvIfNotExists("IDIARY_DOMAIN", "http://127.0.0.1:8080")
	setEnvIfNotExists("IDIARY_PORT", "8080")
	setEnvIfNotExists("IDIARY_SUMMARY_MAX", "555")
	setEnvIfNotExists("IDIARY_SUFFIX", "...")
	setEnvIfNotExists("IDIARY_PAGE_SIZE", "20")
	setEnvIfNotExists("IDIARY_MAX_DAYBOOK_LIMIT", "500")
	setEnvIfNotExists("IDIARY_STATIC_DIR", filepath.Join(GetProjectRoot(), "website", "static"))
	setEnvIfNotExists("IDIARY_LOCALES_DIR", filepath.Join(GetProjectRoot(), "website", "locales"))
	setEnvIfNotExists("IDIARY_DEFAULT_LANG", "zh-cn")

	// 用户配置
	setEnvIfNotExists("IDIARY_USER_DEFAULT_GROUP", "1")
	setEnvIfNotExists("IDIARY_PWD_MIN_LEN", "6")
	setEnvIfNotExists("IDIARY_PWD_MAX_LEN", "20")
	setEnvIfNotExists("IDIARY_USER_STATUS_DEL", "0")
	setEnvIfNotExists("IDIARY_USER_STATUS_OK", "1")
	setEnvIfNotExists("IDIARY_USER_STATUS_WAIT_ACTIVE", "2")

	// 验证码配置
	setEnvIfNotExists("IDIARY_VALID_CODE_LEN", "6")
	setEnvIfNotExists("IDIARY_VALID_CODE_ACTIVE_USER", "1")
	setEnvIfNotExists("IDIARY_VALID_CODE_RESET_PWD", "2")
	setEnvIfNotExists("IDIARY_VALID_CODE_EXPIRE_MINUTES", "30")

	// 必需的JSON配置
	setEnvIfNotExists("IDIARY_SUPPORT_LANGS", `[{"code":"zh-cn","label":"简体中文","local-text":"简体中文"},{"code":"en","label":"English","local-text":"English"}]`)
	setEnvIfNotExists("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC"]`)
	setEnvIfNotExists("IDIARY_PARAMS", `[{"key":"test_param","val":"test_value"}]`)

	// 数据库配置 - 使用测试数据库
	setEnvIfNotExists("IDIARY_DB_HOST", "127.0.0.1")
	setEnvIfNotExists("IDIARY_DB_PORT", "3306")
	setEnvIfNotExists("IDIARY_DB_NAME", "hebe_test") // 重要：使用测试数据库
	setEnvIfNotExists("IDIARY_DB_USER", "root")
	setEnvIfNotExists("IDIARY_DB_PASSWORD", "123123")
	setEnvIfNotExists("IDIARY_DB_CONN", "%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_general_ci")

	// Redis配置（测试用）
	setEnvIfNotExists("IDIARY_REDIS_HOST", "127.0.0.1")
	setEnvIfNotExists("IDIARY_REDIS_PORT", "6379")
	setEnvIfNotExists("IDIARY_REDIS_PASSWORD", "woaibeijingtiananmeng")

	// Cookie配置
	setEnvIfNotExists("IDIARY_COOKIE_DOMAIN", "127.0.0.1")
	setEnvIfNotExists("IDIARY_COOKIE_SECRET", "test_secret_key_12345678901234567890")
	setEnvIfNotExists("IDIARY_COOKIE_EXPIRE_DAYS", "7")
	setEnvIfNotExists("IDIARY_COOKIE_HTTP_ONLY", "true")

	// 邮件配置（测试用）
	setEnvIfNotExists("IDIARY_SMTP_SERVER", "smtp.test.com")
	setEnvIfNotExists("IDIARY_SMTP_USER", "test@test.com")
	setEnvIfNotExists("IDIARY_SMTP_PASSWORD", "test_password")
	setEnvIfNotExists("IDIARY_SMTP_PORT", "25")
	setEnvIfNotExists("IDIARY_SMTP_FROM", "test@test.com")
	setEnvIfNotExists("IDIARY_EMAIL_MAX_RETRIES", "3")
	setEnvIfNotExists("IDIARY_EMAIL_BATCH_SIZE", "10")

	// 文件上传配置 - 使用跨平台路径
	testUploadsPath := filepath.Join(GetProjectRoot(), "tests", "test_uploads")
	setEnvIfNotExists("IDIARY_UPLOAD_PATH", testUploadsPath)
	setEnvIfNotExists("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif,pdf,txt")
	setEnvIfNotExists("IDIARY_UPLOAD_MAX_SIZE", "10485760") // 10MB

	// 日志配置 - 使用跨平台路径
	testLogsPath := filepath.Join(GetProjectRoot(), "tests", "test_logs")
	setEnvIfNotExists("IDIARY_LOG_LEVEL", "debug")
	setEnvIfNotExists("IDIARY_LOG_FORMAT", "text")
	setEnvIfNotExists("IDIARY_LOG_PATH", testLogsPath)
	setEnvIfNotExists("IDIARY_LOG_FILENAME", "test.log")
	setEnvIfNotExists("IDIARY_LOG_MAX_SIZE", "100")
	setEnvIfNotExists("IDIARY_LOG_MAX_AGE", "30")
	setEnvIfNotExists("IDIARY_LOG_MAX_BACKUPS", "5")
	setEnvIfNotExists("IDIARY_LOG_COMPRESS", "true")
	setEnvIfNotExists("IDIARY_LOG_LOCAL_TIME", "true")

	// JSON格式配置
	setSupportedLanguages()
	setTimezones()
	setParams()
}

// setSupportedLanguages 设置支持的语言
func setSupportedLanguages() {
	langs := `[{"code":"zh-cn","label":"Chinese Simplified","local-text":"简体中文"},{"code":"en","label":"English","local-text":"English"}]`
	setEnvIfNotExists("IDIARY_SUPPORT_LANGS", langs)
}

// setTimezones 设置时区
func setTimezones() {
	timezones := `["Asia/Shanghai","UTC","America/New_York"]`
	setEnvIfNotExists("IDIARY_TIMEZONES", timezones)
}

// setParams 设置参数
func setParams() {
	params := `[{"key":"test_key","val":"test_value"}]`
	setEnvIfNotExists("IDIARY_PARAMS", params)
}

// getDefaultTableOrder 获取默认的表清理顺序
// 基于sql/schema.sql中的实际表结构，考虑外键约束
func getDefaultTableOrder() []string {
	return []string{
		"br_comment",      // 评论表（子表）
		"br_like",         // 点赞表（关联表）
		"br_follow",       // 关注表（关联表）
		"br_job_history",  // 任务历史表
		"br_job_queue",    // 任务队列表
		"br_diary_member", // 日记成员表（依赖日记和用户）
		"br_post",         // 文章表（依赖日记和用户）
		"br_diary",        // 日记表（依赖用户）
		"br_user_stat",    // 用户统计表（依赖用户）
		"br_valid_code",   // 验证码表（依赖用户）
		"br_user",         // 用户表（主表）
	}
}

// 辅助函数来获取配置（替代方法）
func GetDatabaseConfig(config *TestConfig) DatabaseConfig {
	return config.Database
}

func GetCleanupConfig(config *TestConfig) CleanupConfig {
	return config.Cleanup
}

func GetIsolationConfig(config *TestConfig) IsolationConfig {
	return config.Isolation
}

func GetPerformanceConfig(config *TestConfig) PerformanceConfig {
	return config.Performance
}

// 工具函数
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func setEnvIfNotExists(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// validateTestConfig 验证测试配置
func validateTestConfig(config *TestConfig) error {
	// 验证数据库配置
	if config.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("数据库端口必须在1-65535之间")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if config.Database.User == "" {
		return fmt.Errorf("数据库用户不能为空")
	}

	// 验证清理配置
	validStrategies := []string{"prefix", "manual", "immediate"}
	if !contains(validStrategies, config.Cleanup.Strategy) {
		return fmt.Errorf("无效的清理策略: %s", config.Cleanup.Strategy)
	}

	// 验证隔离配置
	if config.Isolation.MaxPrefixLen <= 0 {
		return fmt.Errorf("前缀最大长度必须大于0")
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetProjectRoot 获取项目根目录
func GetProjectRoot() string {
	// 从当前文件位置开始，向上查找包含go.mod的目录
	dir, err := os.Getwd()
	if err != nil {
		// 如果无法获取当前工作目录，尝试从可执行文件路径开始
		if execPath, execErr := os.Executable(); execErr == nil {
			dir = filepath.Dir(execPath)
		} else {
			panic(fmt.Sprintf("无法获取当前工作目录和可执行文件路径: %v, %v", err, execErr))
		}
	}
	
	// 向上查找包含go.mod的目录
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达根目录，无法找到go.mod文件
			panic(fmt.Sprintf("无法找到项目根目录，请确保在Go项目目录中运行测试。当前目录: %s", dir))
		}
		dir = parent
	}
}
