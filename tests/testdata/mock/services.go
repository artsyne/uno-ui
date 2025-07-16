package mock

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
)

// MockUser 模拟用户结构
type MockUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

// UserServiceMock 用户服务模拟器
type UserServiceMock struct {
	users     map[int]*MockUser
	Container *di.Container
}

// NewUserServiceMock 创建用户服务模拟器
func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{
		users: make(map[int]*MockUser),
	}
}

// AddUser 添加模拟用户
func (m *UserServiceMock) AddUser(u *MockUser) {
	m.users[u.ID] = u
}

// GetUser 获取模拟用户
func (m *UserServiceMock) GetUser(id int) (*MockUser, error) {
	if u, exists := m.users[id]; exists {
		return u, nil
	}
	return nil, fmt.Errorf("用户不存在")
}

// SetupMockUserService 设置模拟用户服务
func SetupMockUserService(t *testing.T, container *di.Container) *UserServiceMock {
	t.Helper()

	mock := NewUserServiceMock()

	// 添加测试用户
	mock.AddUser(&MockUser{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   1,
	})

	mock.Container = container

	// 这里可以根据实际的用户服务接口进行绑定
	// err := di.Bind(container, func() UserServiceInterface {
	//     return mock
	// })
	// if err != nil {
	//     t.Fatalf("绑定用户服务到DI容器失败: %v", err)
	// }

	return mock
}

// MockDiary 模拟日记结构
type MockDiary struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// DiaryServiceMock 日记服务模拟器
type DiaryServiceMock struct {
	diaries   map[int]*MockDiary
	Container *di.Container
}

// NewDiaryServiceMock 创建日记服务模拟器
func NewDiaryServiceMock() *DiaryServiceMock {
	return &DiaryServiceMock{
		diaries: make(map[int]*MockDiary),
	}
}

// SetupMockDiaryService 设置模拟日记服务
func SetupMockDiaryService(t *testing.T, container *di.Container) *DiaryServiceMock {
	t.Helper()

	mock := NewDiaryServiceMock()
	mock.Container = container

	return mock
}

// MockPost 模拟文章结构
type MockPost struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PostServiceMock 文章服务模拟器
type PostServiceMock struct {
	posts     map[int]*MockPost
	Container *di.Container
}

// NewPostServiceMock 创建文章服务模拟器
func NewPostServiceMock() *PostServiceMock {
	return &PostServiceMock{
		posts: make(map[int]*MockPost),
	}
}

// SetupMockPostService 设置模拟文章服务
func SetupMockPostService(t *testing.T, container *di.Container) *PostServiceMock {
	t.Helper()

	mock := NewPostServiceMock()
	mock.Container = container

	return mock
}

// MockTask 模拟任务结构
type MockTask struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// TaskServiceMock 任务服务模拟器
type TaskServiceMock struct {
	tasks     map[int]*MockTask
	Container *di.Container
}

// NewTaskServiceMock 创建任务服务模拟器
func NewTaskServiceMock() *TaskServiceMock {
	return &TaskServiceMock{
		tasks: make(map[int]*MockTask),
	}
}

// SetupMockTaskService 设置模拟任务服务
func SetupMockTaskService(t *testing.T, container *di.Container) *TaskServiceMock {
	t.Helper()

	mock := NewTaskServiceMock()
	mock.Container = container

	return mock
}
