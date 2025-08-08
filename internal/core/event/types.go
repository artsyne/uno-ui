// Package event 提供事件驱动架构的核心功能
package event

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Event 定义事件接口，所有事件都必须实现此接口
type Event interface {
	// EventType 返回事件类型标识符
	EventType() string
	// Timestamp 返回事件发生时间
	Timestamp() time.Time
}

// EventHandler 定义事件处理器函数类型
// T 必须实现 Event 接口
type EventHandler[T Event] func(ctx context.Context, event T) error

// SubscriptionID 订阅唯一标识符
type SubscriptionID uint64

// EventBus 泛型事件总线接口
type EventBus[T Event] interface {
	// Publish 发布事件到所有订阅的处理器
	Publish(ctx context.Context, event T) error
	// Subscribe 订阅事件处理器，返回订阅ID
	Subscribe(handler EventHandler[T]) (SubscriptionID, error)
	// Unsubscribe 根据订阅ID取消订阅
	Unsubscribe(id SubscriptionID) error
	// Close 关闭事件总线
	Close() error
	// GetSubscriptionCount 获取当前订阅数量
	GetSubscriptionCount() int
}

// EventBusManager 事件总线管理器接口
type EventBusManager[T Event] interface {
	// CreateBus 创建并存储命名的事件总线
	CreateBus(name string, logger *slog.Logger) error
	// GetBus 获取命名的事件总线
	GetBus(name string) (EventBus[T], error)
	// RemoveBus 移除并关闭命名的事件总线
	RemoveBus(name string) error
	// ListBuses 列出所有事件总线名称
	ListBuses() []string
	// CloseAll 关闭所有事件总线
	CloseAll()
	// IsClosed 检查管理器是否已关闭
	IsClosed() bool
}

// subscription 内部订阅结构
type subscription[T Event] struct {
	id      SubscriptionID
	handler EventHandler[T]
}

// eventBus 事件总线实现
type eventBus[T Event] struct {
	// subscriptions 订阅列表
	subscriptions []subscription[T]
	// mu 保护并发操作的互斥锁
	mu sync.RWMutex
	// closed 标识总线是否已关闭
	closed int32
	// logger 日志记录器
	logger *slog.Logger
	// handlerTimeout 处理器超时时间
	handlerTimeout time.Duration
}

// GetSubscriptionCount 返回当前订阅数量
func (eb *eventBus[T]) GetSubscriptionCount() int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.subscriptions)
}

// eventBusManager 事件总线管理器实现
type eventBusManager[T Event] struct {
	// buses 存储命名的事件总线
	buses map[string]EventBus[T]
	// mu 保护并发操作的互斥锁
	mu sync.RWMutex
	// closed 标识管理器是否已关闭
	closed bool
	// logger 日志记录器
	logger *slog.Logger
}

// TypedEventBusManager 类型化事件总线管理器
type TypedEventBusManager[T Event] struct {
	// buses 存储命名的事件总线
	buses map[string]EventBus[T]
	// mu 保护并发操作的互斥锁
	mu sync.RWMutex
	// closed 标识管理器是否已关闭
	closed int32
	// logger 日志记录器
	logger *slog.Logger
}

// GlobalEventBusRegistry 全局事件总线注册表
type GlobalEventBusRegistry struct {
	managers map[string]interface{}
	mu       sync.RWMutex
	closed   int32
}

// EventConfig 事件总线配置
type EventConfig struct {
	// HandlerTimeout 处理器超时时间，默认30秒
	HandlerTimeout time.Duration
	// Logger 日志记录器，如果为nil则使用默认slog
	Logger *slog.Logger
	// EnableMetrics 是否启用指标收集
	EnableMetrics bool
}

// DefaultEventConfig 返回默认配置
func DefaultEventConfig() *EventConfig {
	return &EventConfig{
		HandlerTimeout: 30 * time.Second,
		Logger:         slog.Default(),
		EnableMetrics:  false,
	}
}

// EventMetrics 事件指标
type EventMetrics struct {
	// PublishedEvents 已发布事件数量
	PublishedEvents int64
	// ProcessedEvents 已处理事件数量
	ProcessedEvents int64
	// FailedEvents 失败事件数量
	FailedEvents int64
	// ActiveSubscriptions 活跃订阅数量
	ActiveSubscriptions int64
}

// BaseEvent 基础事件实现
type BaseEvent struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

// EventType 返回事件类型
func (e *BaseEvent) EventType() string {
	return e.Type
}

// Timestamp 返回事件时间戳
func (e *BaseEvent) Timestamp() time.Time {
	return e.CreatedAt
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventType string, data any) *BaseEvent {
	return &BaseEvent{
		Type:      eventType,
		CreatedAt: time.Now(),
		Data:      data,
	}
}