package testhelpers

import (
	"os"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/tests/testdata/mock"
)

// TestSuite 测试套件
// 提供完整的测试环境，包括配置、DI容器、数据库等
type TestSuite struct {
	Config    *ConfigTestHelper
	Container *di.Container
	Database  *DatabaseTestHelper
	HTTP      *HTTPTestHelper

	// 真实服务接口（用于集成测试）
	UserService  user.ServiceInterface
	DiaryService diary.ServiceInterface
	PostService  post.ServiceInterface
	TaskService  task.ServiceInterface

	// Mock服务（用于单元测试）
	UserServiceMock  *mock.UserServiceMock
	DiaryServiceMock *mock.DiaryServiceMock
	PostServiceMock  *mock.PostServiceMock
	TaskServiceMock  *mock.TaskServiceMock

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

	// 根据环境变量或测试名称决定是否使用真实数据库
	useRealDB := os.Getenv("USE_REAL_DB") == "true" || strings.Contains(t.Name(), "Integration")
	suite.UseRealDB = useRealDB

	if useRealDB {
		// 使用真实数据库进行集成测试
		suite.Database = SetupTestDatabase(t)
		suite.setupRealServices(t)
	} else {
		// 使用Mock数据库进行单元测试
		suite.setupMockServices(t)
	}

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
	s.DiaryService = di.Get[diary.ServiceInterface](s.Container)
	s.PostService = di.Get[post.ServiceInterface](s.Container)
	s.TaskService = di.Get[task.ServiceInterface](s.Container)
}

// setupMockServices 设置Mock服务（用于单元测试）
func (s *TestSuite) setupMockServices(t *testing.T) {
	t.Helper()

	s.UserServiceMock = mock.NewUserServiceMock()
	s.DiaryServiceMock = mock.NewDiaryServiceMock()
	s.PostServiceMock = mock.NewPostServiceMock()
	s.TaskServiceMock = mock.NewTaskServiceMock()

	// 设置Mock服务的容器引用
	s.UserServiceMock.Container = s.Container
	s.DiaryServiceMock.Container = s.Container
	s.PostServiceMock.Container = s.Container
	s.TaskServiceMock.Container = s.Container

	// 将Mock服务注册到DI容器
	di.Value[user.ServiceInterface](s.Container, s.UserServiceMock)
	di.Value[diary.ServiceInterface](s.Container, s.DiaryServiceMock)
	di.Value[post.ServiceInterface](s.Container, s.PostServiceMock)
	di.Value[task.ServiceInterface](s.Container, s.TaskServiceMock)
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
