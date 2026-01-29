package types

import (
	"time"

	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/utils/config"
)

// TestConfig 统一的测试配置
type TestConfig struct {
	// 应用配置
	App config.AppConfig `json:"app"`

	// 用户配置
	User config.UserConfig `json:"user"`

	// 数据库配置
	Database DatabaseConfig `json:"database"`

	// Redis配置
	Redis config.RedisConfig `json:"redis"`

	// 清理配置
	Cleanup CleanupConfig `json:"cleanup"`

	// 隔离配置
	Isolation IsolationConfig `json:"isolation"`

	// 性能配置
	Performance PerformanceConfig `json:"performance"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host" default:"127.0.0.1"`
	Port     int    `json:"port" default:"3306"`
	Name     string `json:"name" default:"hebe_test"`
	User     string `json:"user" default:"root"`
	Password string `json:"password" default:"123123"`
	Charset  string `json:"charset" default:"utf8mb4"`
	MaxConns int    `json:"max_conns" default:"10"`
	MaxIdle  int    `json:"max_idle" default:"5"`
}

// CleanupConfig 清理配置
type CleanupConfig struct {
	Enabled     bool          `json:"enabled" default:"false"`   // 默认不清理，便于调试
	Strategy    string        `json:"strategy" default:"prefix"` // prefix|manual|immediate
	Retention   time.Duration `json:"retention" default:"24h"`   // 数据保留时间
	BatchSize   int           `json:"batch_size" default:"100"`  // 批量清理大小
	TablesOrder []string      `json:"tables_order"`              // 清理表的顺序（考虑外键约束）
}

// IsolationConfig 隔离配置
type IsolationConfig struct {
	Strategy     string `json:"strategy" default:"prefix"`          // prefix|namespace|schema
	PrefixFormat string `json:"prefix_format" default:"t_{random}"` // 前缀格式 - 缩短以适应用户名长度限制
	MaxPrefixLen int    `json:"max_prefix_len" default:"10"`        // 最大前缀长度 - 缩短以适应用户名长度限制(16)
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	Parallel      bool          `json:"parallel" default:"true"`         // 允许并发测试
	Timeout       time.Duration `json:"timeout" default:"30s"`           // 测试超时
	MaxRetries    int           `json:"max_retries" default:"3"`         // 最大重试次数
	SkipSlowTests bool          `json:"skip_slow_tests" default:"false"` // 跳过慢速测试
}

// TestFlags 测试标志
type TestFlags struct {
	ClearDB  bool `json:"cleardb"`  // 是否清理数据库
	Verbose  bool `json:"verbose"`  // 详细模式
	Parallel bool `json:"parallel"` // 并行测试
	Timeout  int  `json:"timeout"`  // 超时时间（秒）
}

// UserOptions 用户创建选项
type UserOptions struct {
	Handle      string
	Email       string
	Password    string
	DisplayName string
	TimeZone    string
	Lang        string
}

// DiaryOptions 日记创建选项
type DiaryOptions struct {
	Title         string
	Description   string
	ViewSetting   int
	SocialSetting int
	ThemeID       int
}

// EntryOptions 日记创建选项
type EntryOptions struct {
	Title    string
	Body     string
	Summary  string
	InfoType entry.InfoType
}

// UserOption 用户选项函数类型
type UserOption func(*UserOptions)

// DiaryOption 日记选项函数类型
type DiaryOption func(*DiaryOptions)

// EntryOption 日记选项函数类型
type EntryOption func(*EntryOptions)

// HTTPResponse HTTP响应结构
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// TestStats 测试统计信息
type TestStats struct {
	TestID       string        `json:"test_id"`
	TestName     string        `json:"test_name"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	DataCreated  int           `json:"data_created"`
	DataCleaned  int           `json:"data_cleaned"`
	DatabaseOps  int           `json:"database_ops"`
	HTTPRequests int           `json:"http_requests"`
	Success      bool          `json:"success"`
	ErrorMessage string        `json:"error_message,omitempty"`
}

// TestConfig接口实现
func (tc *TestConfig) GetDatabaseConfig() DatabaseConfig {
	return tc.Database
}

func (tc *TestConfig) GetCleanupConfig() CleanupConfig {
	return tc.Cleanup
}

func (tc *TestConfig) GetIsolationConfig() IsolationConfig {
	return tc.Isolation
}

func (tc *TestConfig) GetPerformanceConfig() PerformanceConfig {
	return tc.Performance
}

func (tc *TestConfig) GetAppConfig() *config.Config {
	// 返回完整的config.Config结构，包含所有必要的配置字段
	// 这确保了所有服务都能获得完整的配置信息
	return &config.Config{
		Env:  "test",
		App:  tc.App,
		User: tc.User,
		Database: config.DatabaseConfig{
			Host:     tc.Database.Host,
			Port:     tc.Database.Port,
			Name:     tc.Database.Name,
			User:     tc.Database.User,
			Password: tc.Database.Password,
			Conn:     "%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_general_ci",
		},
		Redis: tc.Redis,
		RateLimit: config.RateLimitConfig{
			GeneralLimit: 10.0,
			GeneralBurst: 20,
			AuthLimit:    5.0,
			AuthBurst:    10,
		},
		BruteForce: config.BruteForceConfig{
			MaxFailures:    5,
			LockoutMinutes: 15,
		},
		ValidCode: config.ValidCodeConfig{
			Length:         6,
			TypeActiveUser: 1,
			TypeResetPwd:   2,
			ExpireMinutes:  30,
		},
		Email: config.EmailConfig{
			SMTPServer:   "smtp.test.com",
			SMTPUser:     "test@test.com",
			SMTPPassword: "test_password",
			SMTPPort:     25,
			SMTPFrom:     "test@test.com",
			MaxRetries:   3,
			BatchSize:    10,
		},
		Cookie: config.CookieConfig{
			Domain:     "127.0.0.1",
			Secret:     "test_secret_key_12345678901234567890",
			ExpireDays: 7,
			HttpOnly:   true,
		},
		Upload: config.UploadConfig{
			Path:    "/tmp/test_uploads",
			Types:   "jpg,jpeg,png,gif,pdf,txt",
			MaxSize: 10485760,
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
			Level:      "debug",
			Format:     "text",
			Path:       "/tmp/test_logs",
			Filename:   "test.log",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 5,
			Compress:   true,
			LocalTime:  true,
		},
	}
}
