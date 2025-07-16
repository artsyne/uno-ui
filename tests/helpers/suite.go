package testhelpers

import (
	"os"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/tests/testdata/mock"
)

// TestSuite 测试套件，提供完整的测试环境
type TestSuite struct {
	Config    *ConfigTestHelper
	Container *di.Container
	Database  *DatabaseTestHelper
	User      *mock.UserServiceMock
	Diary     *mock.DiaryServiceMock
	Post      *mock.PostServiceMock
	Task      *mock.TaskServiceMock
}

// SetupTestSuite 设置完整的测试套件
func SetupTestSuite(t *testing.T) *TestSuite {
	t.Helper()

	// 创建独立的测试DI容器
	container := di.New()

	// 设置配置（如果需要真实数据库，先设置数据库）
	var config *ConfigTestHelper
	var database *DatabaseTestHelper

	// 检查是否需要真实数据库（通过环境变量控制）
	if os.Getenv("USE_REAL_DB") == "true" || strings.Contains(t.Name(), "Database") {
		// 设置真实数据库
		database = SetupTestDatabase(t)
		// 数据库设置会修改环境变量，所以配置在数据库之后初始化
		config = &ConfigTestHelper{} // 配置已经在数据库设置中初始化
	} else {
		// 设置配置
		config = SetupTestConfig(t)
		// 设置模拟数据库
		mockDB := mock.SetupMockDatabase(t, container)
		_ = mockDB // 使用模拟数据库的场景暂时保留原有逻辑
	}

	// 设置模拟服务
	user := mock.SetupMockUserService(t, container)
	diary := mock.SetupMockDiaryService(t, container)
	post := mock.SetupMockPostService(t, container)
	task := mock.SetupMockTaskService(t, container)

	return &TestSuite{
		Config:    config,
		Container: container,
		Database:  database,
		User:      user,
		Diary:     diary,
		Post:      post,
		Task:      task,
	}
}

// Cleanup 清理测试套件
func (s *TestSuite) Cleanup(t *testing.T) {
	t.Helper()

	// 清理数据库（真实数据库或模拟数据库）
	if s.Database != nil {
		s.Database.Cleanup(t)
	}

	// 清理配置
	if s.Config != nil {
		s.Config.Cleanup(t)
	}

	// 清空DI容器
	if s.Container != nil {
		s.Container.Clear()
	}
}

// WithCleanup 自动清理的测试套件
func WithCleanup(t *testing.T, testFunc func(*TestSuite)) {
	t.Helper()

	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)

	testFunc(suite)
}
