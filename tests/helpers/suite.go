package testhelpers

import (
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
)

// TestSuite 测试套件
// 提供完整的测试环境，包括配置、DI容器、数据库等
// 统一使用真实服务进行测试，确保测试的真实性
type TestSuite struct {
	Config    *ConfigTestHelper
	Container *di.Container
	Database  *DatabaseTestHelper
	HTTP      *HTTPTestHelper

	// 真实服务接口（所有测试都使用真实服务）
	UserService  user.ServiceInterface
	DiaryService diary.ServiceInterface
	PostService  post.ServiceInterface
	TaskService  task.ServiceInterface

	// 是否使用真实数据库
	UseRealDB bool
}

// SetupTestSuite 设置完整的测试套件
// 根据环境变量或测试名称决定是否使用真实数据库
func SetupTestSuite(t *testing.T) *TestSuite {
	t.Helper()

	suite := &TestSuite{}

	// 设置配置
	suite.Config = SetupTestConfig(t)

	// 设置DI容器
	suite.Container = di.New()

	// 检查是否使用真实数据库
	useRealDB := ShouldUseRealDB() ||
		strings.Contains(t.Name(), "Integration") ||
		strings.Contains(t.Name(), "Database")
	suite.UseRealDB = useRealDB

	// 设置数据库（如果需要）
	if useRealDB {
		suite.Database = SetupTestDatabase(t)
	}

	// 始终设置真实服务
	suite.setupRealServices(t)

	// 设置HTTP测试助手（需要在服务设置之后）
	suite.HTTP = SetupHTTPTest(suite)

	return suite
}

// setupRealServices 设置真实服务（用于集成测试）
func (s *TestSuite) setupRealServices(t *testing.T) {
	t.Helper()

	// 注册真实的域服务到DI容器
	di.Singleton(s.Container, func() user.ServiceInterface {
		return user.NewService()
	})

	di.Singleton(s.Container, func() diary.ServiceInterface {
		return diary.NewService()
	})

	di.Singleton(s.Container, func() post.ServiceInterface {
		return post.NewService()
	})

	di.Singleton(s.Container, func() task.ServiceInterface {
		return task.NewService()
	})

	// 获取注册的真实服务
	s.UserService = di.Get[user.ServiceInterface](s.Container)
	if s.UserService == nil {
		t.Fatal("UserService 未能从DI容器中获取")
	}
	s.DiaryService = di.Get[diary.ServiceInterface](s.Container)
	if s.DiaryService == nil {
		t.Fatal("DiaryService 未能从DI容器中获取")
	}
	s.PostService = di.Get[post.ServiceInterface](s.Container)
	if s.PostService == nil {
		t.Fatal("PostService 未能从DI容器中获取")
	}
	s.TaskService = di.Get[task.ServiceInterface](s.Container)
	if s.TaskService == nil {
		t.Fatal("TaskService 未能从DI容器中获取")
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
