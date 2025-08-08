package services

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/core/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UserEvent 用户事件
type UserEvent struct {
	UserID    string
	Action    string
	Metadata  map[string]interface{}
	timestamp time.Time
}

// EventType 实现 Event 接口
func (e *UserEvent) EventType() string {
	return "user.event"
}

// Timestamp 实现 Event 接口
func (e *UserEvent) Timestamp() time.Time {
	return e.timestamp
}

// TaskEvent 任务事件
type TaskEvent struct {
	TaskID    string
	Status    string
	Metadata  map[string]interface{}
	timestamp time.Time
}

// EventType 实现 Event 接口
func (e *TaskEvent) EventType() string {
	return "task.event"
}

// Timestamp 实现 Event 接口
func (e *TaskEvent) Timestamp() time.Time {
	return e.timestamp
}

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	UserID    string
	Username  string
	Email     string
	timestamp time.Time
}

// EventType 实现 Event 接口
func (e *UserCreatedEvent) EventType() string {
	return "user.created"
}

// Timestamp 实现 Event 接口
func (e *UserCreatedEvent) Timestamp() time.Time {
	return e.timestamp
}

// TaskCreatedEvent 任务创建事件
type TaskCreatedEvent struct {
	TaskID      string
	Title       string
	Description string
	UserID      string
	timestamp   time.Time
}

// EventType 实现 Event 接口
func (e *TaskCreatedEvent) EventType() string {
	return "task.created"
}

// Timestamp 实现 Event 接口
func (e *TaskCreatedEvent) Timestamp() time.Time {
	return e.timestamp
}

// User 用户模型
type User struct {
	ID       string
	Username string
	Email    string
}

// Task 任务模型
type Task struct {
	ID          string
	Title       string
	Description string
	UserID      string
}

// MockUserService 模拟用户服务
type MockUserService struct {
	eventBus event.EventBus[*UserCreatedEvent]
	users    map[string]*User
	mu       sync.RWMutex
}

// NewMockUserService 创建模拟用户服务
func NewMockUserService(eventBus event.EventBus[*UserCreatedEvent]) *MockUserService {
	return &MockUserService{
		eventBus: eventBus,
		users:    make(map[string]*User),
	}
}

// CreateUser 创建用户
func (s *MockUserService) CreateUser(ctx context.Context, id, username, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		ID:       id,
		Username: username,
		Email:    email,
	}
	s.users[id] = user

	// 发布用户创建事件
	userEvent := &UserCreatedEvent{
		UserID:    id,
		Username:  username,
		Email:     email,
		timestamp: time.Now(),
	}
	return s.eventBus.Publish(ctx, userEvent)
}

// GetUser 获取用户
func (s *MockUserService) GetUser(id string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

// MockTaskService 模拟任务服务
type MockTaskService struct {
	eventBus event.EventBus[*TaskCreatedEvent]
	tasks    map[string]*Task
	mu       sync.RWMutex
}

// NewMockTaskService 创建模拟任务服务
func NewMockTaskService(eventBus event.EventBus[*TaskCreatedEvent]) *MockTaskService {
	return &MockTaskService{
		eventBus: eventBus,
		tasks:    make(map[string]*Task),
	}
}

// CreateTask 创建任务
func (s *MockTaskService) CreateTask(ctx context.Context, id, title, description, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		UserID:      userID,
	}
	s.tasks[id] = task

	// 发布任务创建事件
	taskEvent := &TaskCreatedEvent{
		TaskID:      id,
		Title:       title,
		Description: description,
		UserID:      userID,
		timestamp:   time.Now(),
	}
	return s.eventBus.Publish(ctx, taskEvent)
}

// GetTask 获取任务
func (s *MockTaskService) GetTask(id string) *Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id]
}

// MockNotificationService 模拟通知服务
type MockNotificationService struct {
	notifications []string
	mu            sync.Mutex
}

// NewMockNotificationService 创建模拟通知服务
func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{
		notifications: make([]string, 0),
	}
}

// HandleUserCreated 处理用户创建事件
func (s *MockNotificationService) HandleUserCreated(ctx context.Context, event *UserCreatedEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	message := "Welcome " + event.Username + "! Your account has been created."
	s.notifications = append(s.notifications, message)
	return nil
}

// HandleTaskCreated 处理任务创建事件
func (s *MockNotificationService) HandleTaskCreated(ctx context.Context, event *TaskCreatedEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	message := "Task '" + event.Title + "' has been created."
	s.notifications = append(s.notifications, message)
	return nil
}

// GetNotifications 获取通知列表
func (s *MockNotificationService) GetNotifications() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string{}, s.notifications...)
}

// TestEventDrivenArchitecture 测试事件驱动架构集成
func TestEventDrivenArchitecture(t *testing.T) {
	// 创建用户事件总线
	userEventBus := event.NewEventBus[*UserCreatedEvent](nil)
	defer userEventBus.Close()

	// 创建任务事件总线
	taskEventBus := event.NewEventBus[*TaskCreatedEvent](nil)
	defer taskEventBus.Close()

	// 创建服务
	userService := NewMockUserService(userEventBus)
	taskService := NewMockTaskService(taskEventBus)
	notificationService := NewMockNotificationService()

	// 订阅事件
	_, err := userEventBus.Subscribe(notificationService.HandleUserCreated)
	if err != nil {
		t.Fatalf("Failed to subscribe user created handler: %v", err)
	}

	_, err = taskEventBus.Subscribe(notificationService.HandleTaskCreated)
	if err != nil {
		t.Fatalf("Failed to subscribe task created handler: %v", err)
	}

	ctx := context.Background()

	// 创建用户
	err = userService.CreateUser(ctx, "user1", "john_doe", "john@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 创建任务
	err = taskService.CreateTask(ctx, "task1", "Complete project", "Finish the project by deadline", "user1")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// 等待事件处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证用户和任务已创建
	user := userService.GetUser("user1")
	if user == nil {
		t.Fatal("User not found")
	}
	if user.Username != "john_doe" {
		t.Errorf("Expected username 'john_doe', got '%s'", user.Username)
	}

	task := taskService.GetTask("task1")
	if task == nil {
		t.Fatal("Task not found")
	}
	if task.Title != "Complete project" {
		t.Errorf("Expected task title 'Complete project', got '%s'", task.Title)
	}

	// 验证通知已发送
	notifications := notificationService.GetNotifications()
	if len(notifications) != 2 {
		t.Fatalf("Expected 2 notifications, got %d", len(notifications))
	}

	// 检查通知内容
	expectedNotifications := []string{
		"Welcome john_doe! Your account has been created.",
		"Task 'Complete project' has been created.",
	}

	for i, expected := range expectedNotifications {
		if notifications[i] != expected {
			t.Errorf("Expected notification '%s', got '%s'", expected, notifications[i])
		}
	}
}

// TestEventIntegration 测试事件系统集成
func TestEventIntegration(t *testing.T) {
	// 获取用户事件总线管理器
	userManager := event.GetTypedEventBusManager[*UserEvent]("user", nil)

	// 获取任务事件总线管理器
	taskManager := event.GetTypedEventBusManager[*TaskEvent]("task", nil)

	// 创建用户事件总线
	err := userManager.CreateBus("user_bus", nil)
	require.NoError(t, err)

	// 创建任务事件总线
	err = taskManager.CreateBus("task_bus", nil)
	require.NoError(t, err)

	// 获取事件总线
	userBus, err := userManager.GetBus("user_bus")
	require.NoError(t, err)

	taskBus, err := taskManager.GetBus("task_bus")
	require.NoError(t, err)

	// 用户事件处理器
	userEventReceived := false
	userHandler := func(ctx context.Context, e *UserEvent) error {
		assert.Equal(t, "user123", e.UserID)
		assert.Equal(t, "login", e.Action)
		userEventReceived = true
		return nil
	}

	// 任务事件处理器
	taskEventReceived := false
	taskHandler := func(ctx context.Context, e *TaskEvent) error {
		assert.Equal(t, "task456", e.TaskID)
		assert.Equal(t, "completed", e.Status)
		taskEventReceived = true
		return nil
	}

	// 订阅事件
	_, err = userBus.Subscribe(userHandler)
	require.NoError(t, err)

	_, err = taskBus.Subscribe(taskHandler)
	require.NoError(t, err)

	// 发布用户事件
	userEvent := &UserEvent{
		UserID: "user123",
		Action: "login",
		Metadata: map[string]interface{}{
			"ip": "192.168.1.1",
		},
		timestamp: time.Now(),
	}
	ctx := context.Background()
	err = userBus.Publish(ctx, userEvent)
	require.NoError(t, err)

	// 发布任务事件
	taskEvent := &TaskEvent{
		TaskID: "task456",
		Status: "completed",
		Metadata: map[string]interface{}{
			"duration": "5m",
		},
		timestamp: time.Now(),
	}
	err = taskBus.Publish(ctx, taskEvent)
	require.NoError(t, err)

	// 等待事件处理
	time.Sleep(100 * time.Millisecond)

	// 验证事件是否被正确处理
	assert.True(t, userEventReceived, "User event should be received")
	assert.True(t, taskEventReceived, "Task event should be received")

	// 清理
	userManager.CloseAll()
	taskManager.CloseAll()
}

// TestTypedEventBusManager 测试类型化事件总线管理器
func TestTypedEventBusManager(t *testing.T) {
	manager := event.GetTypedEventBusManager[*UserEvent]("test", nil)
	defer manager.CloseAll()

	// 创建命名的事件总线
	err := manager.CreateBus("user-events", nil)
	require.NoError(t, err)

	err = manager.CreateBus("user-actions", nil)
	require.NoError(t, err)

	// 获取事件总线
	userBus, err := manager.GetBus("user-events")
	require.NoError(t, err)

	actionBus, err := manager.GetBus("user-actions")
	require.NoError(t, err)

	// 测试事件处理
	var receivedEvents []*UserEvent
	var mu sync.Mutex

	handler := func(ctx context.Context, event *UserEvent) error {
		mu.Lock()
		defer mu.Unlock()
		receivedEvents = append(receivedEvents, event)
		return nil
	}

	// 订阅事件
	_, err = userBus.Subscribe(handler)
	require.NoError(t, err)

	_, err = actionBus.Subscribe(handler)
	require.NoError(t, err)

	// 发布事件
	event1 := &UserEvent{UserID: "user1", Action: "login", timestamp: time.Now()}
	event2 := &UserEvent{UserID: "user2", Action: "logout", timestamp: time.Now()}

	ctx := context.Background()
	err = userBus.Publish(ctx, event1)
	require.NoError(t, err)

	err = actionBus.Publish(ctx, event2)
	require.NoError(t, err)

	// 等待事件处理
	time.Sleep(100 * time.Millisecond)

	// 验证事件接收
	mu.Lock()
	assert.Len(t, receivedEvents, 2)
	assert.Equal(t, "user1", receivedEvents[0].UserID)
	assert.Equal(t, "login", receivedEvents[0].Action)
	assert.Equal(t, "user2", receivedEvents[1].UserID)
	assert.Equal(t, "logout", receivedEvents[1].Action)
	mu.Unlock()

	// 测试总线列表
	buses := manager.ListBuses()
	assert.Len(t, buses, 2)
	assert.Contains(t, buses, "user-events")
	assert.Contains(t, buses, "user-actions")

	// 测试移除总线
	err = manager.RemoveBus("user-actions")
	require.NoError(t, err)

	buses = manager.ListBuses()
	assert.Len(t, buses, 1)
	assert.Contains(t, buses, "user-events")
}