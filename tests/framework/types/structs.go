package framework

import (
	"time"

	"github.com/biuren/hebe/internal/utils/config"
)

// TestConfig 统一的测试配置
type TestConfig struct {
	// 应用配置
	App config.AppConfig `json:"app"`

	// 数据库配置
	Database DatabaseConfig `json:"database"`

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
	Strategy     string `json:"strategy" default:"prefix"`                           // prefix|namespace|schema
	PrefixFormat string `json:"prefix_format" default:"test_{date}_{time}_{random}"` // 前缀格式
	MaxPrefixLen int    `json:"max_prefix_len" default:"20"`                         // 最大前缀长度
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	Parallel      bool          `json:"parallel" default:"true"`         // 允许并发测试
	Timeout       time.Duration `json:"timeout" default:"30s"`           // 测试超时
	MaxRetries    int           `json:"max_retries" default:"3"`         // 最大重试次数
	SkipSlowTests bool          `json:"skip_slow_tests" default:"false"` // 跳过慢速测试
}

// UserOptions 用户创建选项
type UserOptions struct {
	Username string
	Email    string
	Password string
	FullName string
	TimeZone string
	Lang     string
}

// DiaryOptions 日记创建选项
type DiaryOptions struct {
	Title         string
	Description   string
	ViewSetting   int
	SocialSetting int
	ThemeID       int
}

// PostOptions 文章创建选项
type PostOptions struct {
	Title    string
	Body     string
	Summary  string
	InfoType int
}

// UserOption 用户选项函数类型
type UserOption func(*UserOptions)

// DiaryOption 日记选项函数类型
type DiaryOption func(*DiaryOptions)

// PostOption 文章选项函数类型
type PostOption func(*PostOptions)

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
