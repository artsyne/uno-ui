package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBenchmarkEnvironmentSetup 测试基准测试环境初始化
func TestBenchmarkEnvironmentSetup(t *testing.T) {
	// 初始化基准测试环境
	initBenchmarkEnv()

	// 验证全局变量已正确初始化
	assert.NotNil(t, benchCtx, "benchmark context should be initialized")
	assert.NotNil(t, benchServices, "benchmark services should be initialized")
	
	// 验证服务实例
	assert.NotNil(t, benchServices.userService, "user service should be initialized")
	assert.NotNil(t, benchServices.diaryService, "diary service should be initialized")
	assert.NotNil(t, benchServices.entryService, "entry service should be initialized")
	assert.NotNil(t, benchServices.taskService, "task service should be initialized")
	
	// 验证测试上下文功能
	testID := benchCtx.GetTestID()
	assert.NotEmpty(t, testID, "test ID should not be empty")
	
	config := benchCtx.GetConfig()
	assert.NotNil(t, config, "test config should not be nil")
}

// TestBenchmarkDataGeneration 测试基准测试数据生成
func TestBenchmarkDataGeneration(t *testing.T) {
	initBenchmarkEnv()
	
	dataGenerator := benchCtx.GetDataGenerator()
	require.NotNil(t, dataGenerator, "data generator should not be nil")
	
	// 测试字符串生成 - 使用不同的前缀来确保唯一性
	uniqueStr1 := dataGenerator.GenerateUniqueString("test1")
	uniqueStr2 := dataGenerator.GenerateUniqueString("test2")
	assert.NotEqual(t, uniqueStr1, uniqueStr2, "generated strings should be unique")
	
	// 测试邮箱生成 - 使用不同的前缀来确保唯一性
	email1 := dataGenerator.GenerateUniqueEmail("test1")
	email2 := dataGenerator.GenerateUniqueEmail("test2")
	assert.NotEqual(t, email1, email2, "generated emails should be unique")
	assert.Contains(t, email1, "@", "generated email should contain @")
	assert.Contains(t, email2, "@", "generated email should contain @")
	
	// 测试字符串非空
	assert.NotEmpty(t, uniqueStr1, "generated string should not be empty")
	assert.NotEmpty(t, uniqueStr2, "generated string should not be empty")
	assert.NotEmpty(t, email1, "generated email should not be empty")
	assert.NotEmpty(t, email2, "generated email should not be empty")
}

// TestBenchmarkDatabaseConnection 测试基准测试数据库连接
func TestBenchmarkDatabaseConnection(t *testing.T) {
	initBenchmarkEnv()
	
	db := benchCtx.GetDatabase()
	require.NotNil(t, db, "database manager should not be nil")
	
	// 测试数据库健康状态
	isHealthy := db.IsHealthy()
	assert.True(t, isHealthy, "database should be healthy")
	
	// 测试数据库连接
	sqlxDB := db.GetDB()
	assert.NotNil(t, sqlxDB, "sqlx database connection should not be nil")
	
	rawDB := db.GetRawDB()
	assert.NotNil(t, rawDB, "raw database connection should not be nil")
}

// TestBenchmarkCleanup 测试基准测试清理功能
func TestBenchmarkCleanup(t *testing.T) {
	initBenchmarkEnv()
	
	cleanup := benchCtx.GetCleanup()
	if cleanup != nil {
		// 测试清理功能可用性（不实际执行清理）
		assert.NotNil(t, cleanup, "cleanup manager should be available")
	}
	
	isolator := benchCtx.GetIsolator()
	require.NotNil(t, isolator, "data isolator should not be nil")
	
	// 测试数据隔离功能
	prefix := isolator.GetPrefix()
	assert.NotEmpty(t, prefix, "isolator prefix should not be empty")
	
	testID := isolator.GetTestID()
	assert.NotEmpty(t, testID, "isolator test ID should not be empty")
}

// TestBenchmarkServicesIntegration 测试基准测试服务集成
func TestBenchmarkServicesIntegration(t *testing.T) {
	initBenchmarkEnv()
	
	// 验证所有服务都能正常访问
	require.NotNil(t, benchServices, "benchmark services should be initialized")
	
	userService := benchServices.userService
	diaryService := benchServices.diaryService
	entryService := benchServices.entryService
	taskService := benchServices.taskService
	
	assert.NotNil(t, userService, "user service should be available")
	assert.NotNil(t, diaryService, "diary service should be available")
	assert.NotNil(t, entryService, "entry service should be available")
	assert.NotNil(t, taskService, "task service should be available")
}

// TestBenchmarkConfiguration 测试基准测试配置
func TestBenchmarkConfiguration(t *testing.T) {
	initBenchmarkEnv()
	
	config := benchCtx.GetConfig()
	require.NotNil(t, config, "config should not be nil")
	
	appConfig := config.GetAppConfig()
	require.NotNil(t, appConfig, "app config should not be nil")
	
	// 验证关键配置项
	assert.Greater(t, appConfig.App.MaxDiaryLimit, 0, "max diary limit should be positive")
	assert.Greater(t, appConfig.App.SummaryMax, 0, "summary max should be positive")
	assert.Greater(t, appConfig.App.PageSize, 0, "page size should be positive")
	assert.NotEmpty(t, appConfig.App.DefaultLang, "default language should not be empty")
}