package framework

import (
	"context"
	"database/sql"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

// TestOption 测试选项函数类型（先声明，具体定义在core包中）
type TestOption interface{}

// TestConfig 统一的测试配置接口
type TestConfigInterface interface {
	GetDatabaseConfig() DatabaseConfig
	GetCleanupConfig() CleanupConfig
	GetIsolationConfig() IsolationConfig
	GetPerformanceConfig() PerformanceConfig
}

// DatabaseManagerInterface 数据库管理器接口
type DatabaseManagerInterface interface {
	Initialize() error
	GetDB() *sqlx.DB
	GetRawDB() *sql.DB
	Close() error
	IsHealthy() bool
}

// DataIsolatorInterface 数据隔离器接口
type DataIsolatorInterface interface {
	IsolateString(value string) string
	IsolateEmail(email string) string
	IsolateUsername(username string) string
	IsolateTitle(title string) string
	TrackCreated(table string, id interface{})
	GetCreatedData() map[string][]interface{}
	GetPrefix() string
	GetTestID() string
	GetDataCount() int
	GetTableCount() int
}

// CleanupManagerInterface 清理管理器接口
type CleanupManagerInterface interface {
	Cleanup(strategy string) error
	CleanupByPrefix() error
	CleanupTrackedData() error
	ForceCleanupDatabase() error
	AddCleanupFunc(fn func() error)
}

// DataGeneratorInterface 测试数据生成器接口
type DataGeneratorInterface interface {
	CreateUser(options ...UserOption) *user.User
	CreateDiary(userID int64, options ...DiaryOption) *diary.Diary
	CreatePost(userID, diaryID int64, options ...PostOption) *post.Post
	GenerateUniqueString(prefix string) string
	GenerateUniqueEmail(prefix string) string
}

// ServiceRegistryInterface 服务注册表接口
type ServiceRegistryInterface interface {
	GetUserService() user.ServiceInterface
	GetDiaryService() diary.ServiceInterface
	GetPostService() post.ServiceInterface
	InitializeServices(db *sqlx.DB) error
}

// HTTPTestHelperInterface HTTP测试助手接口
type HTTPTestHelperInterface interface {
	SetupRouter() error
	GET(path string, headers map[string]string) *HTTPResponse
	POST(path string, body interface{}, headers map[string]string) *HTTPResponse
	PUT(path string, body interface{}, headers map[string]string) *HTTPResponse
	DELETE(path string, headers map[string]string) *HTTPResponse
}

// TestContextInterface 测试上下文接口
type TestContextInterface interface {
	GetTestID() string
	GetConfig() TestConfigInterface
	GetDatabase() DatabaseManagerInterface
	GetDataGenerator() DataGeneratorInterface
	GetIsolator() DataIsolatorInterface
	GetCleanup() CleanupManagerInterface
	GetServices() ServiceRegistryInterface
	GetHTTP() HTTPTestHelperInterface
	Context() context.Context
	Cleanup() error
	ForceCleanupDatabase() error
}
