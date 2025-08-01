package testhelpers

import (
	"flag"
	"os"
	"path/filepath"
	"sync"
)

var (
	// 命令行参数
	cleanupDatabase = flag.Bool("cleanup-db", false, "清理测试数据库数据（默认保留数据以便观察）")
	useRealDB       = flag.Bool("real-db", false, "强制使用真实数据库进行所有测试")

	// 确保环境变量只设置一次
	envSetupOnce sync.Once
	isEnvSetup   bool
)

// SetupTestEnvironment 统一设置测试环境变量
// 这个函数确保环境变量只设置一次，避免重复设置
func SetupTestEnvironment() {
	envSetupOnce.Do(func() {
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
		os.Setenv("IDIARY_UPLOAD_PATH", "tests/test_uploads")
		os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
		os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")

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

		// 设置必需的JSON格式环境变量（使用IDIARY_前缀）
		os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-cn","label":"中文","local-text":"中文"},{"code":"en","label":"English","local-text":"English"}]`)
		os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
		os.Setenv("IDIARY_PARAMS", `[{"key":"test_key","val":"test_value"}]`)

		isEnvSetup = true
	})
}

// IsEnvSetup 检查环境变量是否已设置
func IsEnvSetup() bool {
	return isEnvSetup
}

// ShouldCleanupDatabase 检查是否应该清理数据库
// 通过命令行参数 -cleanup-db 控制，默认为false（不清理）
func ShouldCleanupDatabase() bool {
	// 确保flag已解析
	if !flag.Parsed() {
		flag.Parse()
	}
	return *cleanupDatabase
}

// ShouldUseRealDB 检查是否应该使用真实数据库
func ShouldUseRealDB() bool {
	// 确保flag已解析
	if !flag.Parsed() {
		flag.Parse()
	}

	// 命令行参数优先
	if *useRealDB {
		return true
	}

	// 环境变量次之
	if os.Getenv("USE_REAL_DB") == "true" {
		return true
	}

	return false
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 从当前文件位置开始，向上查找包含go.mod的目录
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达根目录，回退到默认路径
			return "/Users/yasong.cy/go/src/idiary.ai/hebe"
		}
		dir = parent
	}
}
