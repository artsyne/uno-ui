// Package event 提供事件驱动架构的核心功能
// 支持类型安全的事件发布和订阅机制，使用 Go 1.24 泛型特性
package event

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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

// subscription 内部订阅结构
type subscription[T Event] struct {
	id      SubscriptionID
	handler EventHandler[T]
}

// subscriptionIDCounter 全局订阅ID计数器
var subscriptionIDCounter uint64

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
}

// eventBus 泛型事件总线实现
type eventBus[T Event] struct {
	// subscriptions 存储订阅信息
	subscriptions []subscription[T]
	// mu 保护并发操作的互斥锁
	mu sync.RWMutex
	// closed 标识事件总线是否已关闭
	closed int32 // 使用atomic操作
	// logger 用于调试和错误记录
	logger Logger
}

// TypedEventBusManager 类型化事件总线管理器
type TypedEventBusManager[T Event] struct {
	// buses 存储命名的事件总线
	buses map[string]EventBus[T]
	// mu 保护并发操作的互斥锁
	mu sync.RWMutex
	// closed 标识管理器是否已关闭
	closed int32 // 使用atomic操作
	// logger 用于调试和错误记录
	logger Logger
}



// Logger 定义日志接口，用于调试和错误记录
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// defaultLogger 默认日志实现
type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+msg+"\n", args...)
}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func (l *defaultLogger) Warn(msg string, args ...interface{}) {
	fmt.Printf("[WARN] "+msg+"\n", args...)
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}

// NewEventBus 创建新的泛型事件总线
func NewEventBus[T Event](logger Logger) EventBus[T] {
	if logger == nil {
		logger = &defaultLogger{}
	}
	return &eventBus[T]{
		subscriptions: make([]subscription[T], 0),
		logger:        logger,
	}
}

// NewTypedEventBusManager 创建新的类型化事件总线管理器
func NewTypedEventBusManager[T Event](logger Logger) *TypedEventBusManager[T] {
	if logger == nil {
		logger = &defaultLogger{}
	}
	return &TypedEventBusManager[T]{
		buses:  make(map[string]EventBus[T]),
		closed: 0,
		logger: logger,
	}
}

// Publish 发布事件到所有订阅的处理器
func (eb *eventBus[T]) Publish(ctx context.Context, event T) error {
	// 快速检查是否已关闭（无锁）
	if atomic.LoadInt32(&eb.closed) == 1 {
		return fmt.Errorf("event bus is closed")
	}

	// 获取当前订阅的快照
	eb.mu.RLock()
	subscriptionsCopy := make([]subscription[T], len(eb.subscriptions))
	copy(subscriptionsCopy, eb.subscriptions)
	eb.mu.RUnlock()

	eventType := event.EventType()
	eb.logger.Debug("Publishing event: %s", eventType)

	if len(subscriptionsCopy) == 0 {
		eb.logger.Debug("No handlers found for event type: %s", eventType)
		return nil
	}

	// 并发执行所有处理器
	var wg sync.WaitGroup
	errorChan := make(chan error, len(subscriptionsCopy))

	for _, sub := range subscriptionsCopy {
		wg.Add(1)
		go func(s subscription[T]) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					eb.logger.Error("Panic in event handler (ID: %d): %v", s.id, r)
					errorChan <- fmt.Errorf("panic in event handler (ID: %d): %v", s.id, r)
				}
			}()
			if err := s.handler(ctx, event); err != nil {
				eb.logger.Error("Error in event handler (ID: %d): %v", s.id, err)
				errorChan <- fmt.Errorf("handler (ID: %d) error: %w", s.id, err)
			}
		}(sub)
	}

	wg.Wait()
	close(errorChan)

	// 收集所有错误
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		eb.logger.Error("Event handling errors: %v", errors)
		return fmt.Errorf("event handling failed with %d errors: %v", len(errors), errors)
	}

	eb.logger.Debug("Event published successfully: %s", eventType)
	return nil
}

// Subscribe 订阅事件处理器，返回订阅ID
func (eb *eventBus[T]) Subscribe(handler EventHandler[T]) (SubscriptionID, error) {
	if atomic.LoadInt32(&eb.closed) == 1 {
		return 0, fmt.Errorf("event bus is closed")
	}

	// 生成唯一订阅ID
	id := SubscriptionID(atomic.AddUint64(&subscriptionIDCounter, 1))

	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 再次检查关闭状态（双重检查）
	if atomic.LoadInt32(&eb.closed) == 1 {
		return 0, fmt.Errorf("event bus is closed")
	}

	eb.subscriptions = append(eb.subscriptions, subscription[T]{
		id:      id,
		handler: handler,
	})

	eb.logger.Debug("Handler subscribed successfully with ID: %d", id)
	return id, nil
}

// Unsubscribe 根据订阅ID取消订阅
func (eb *eventBus[T]) Unsubscribe(id SubscriptionID) error {
	if atomic.LoadInt32(&eb.closed) == 1 {
		return fmt.Errorf("event bus is closed")
	}

	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 再次检查关闭状态（双重检查）
	if atomic.LoadInt32(&eb.closed) == 1 {
		return fmt.Errorf("event bus is closed")
	}

	// 查找并移除指定ID的订阅
	for i, sub := range eb.subscriptions {
		if sub.id == id {
			// 使用切片技巧高效移除元素
			eb.subscriptions[i] = eb.subscriptions[len(eb.subscriptions)-1]
			eb.subscriptions = eb.subscriptions[:len(eb.subscriptions)-1]
			eb.logger.Debug("Handler unsubscribed successfully with ID: %d", id)
			return nil
		}
	}

	return fmt.Errorf("subscription with ID %d not found", id)
}

// Close 关闭事件总线
func (eb *eventBus[T]) Close() error {
	// 原子性设置关闭状态
	if !atomic.CompareAndSwapInt32(&eb.closed, 0, 1) {
		return fmt.Errorf("event bus is already closed")
	}

	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 清空订阅
	eb.subscriptions = nil
	eb.logger.Info("Event bus closed")
	return nil
}

// CreateBus 创建并存储命名的事件总线
func (manager *TypedEventBusManager[T]) CreateBus(name string, logger Logger) error {
	if atomic.LoadInt32(&manager.closed) == 1 {
		return fmt.Errorf("manager is closed")
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	// 再次检查关闭状态（双重检查）
	if atomic.LoadInt32(&manager.closed) == 1 {
		return fmt.Errorf("manager is closed")
	}

	if _, exists := manager.buses[name]; exists {
		return fmt.Errorf("event bus '%s' already exists", name)
	}

	bus := NewEventBus[T](logger)
	manager.buses[name] = bus
	manager.logger.Debug("Created event bus: %s", name)
	return nil
}

// GetBus 获取命名的事件总线
func (manager *TypedEventBusManager[T]) GetBus(name string) (EventBus[T], error) {
	if atomic.LoadInt32(&manager.closed) == 1 {
		return nil, fmt.Errorf("manager is closed")
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	bus, exists := manager.buses[name]
	if !exists {
		return nil, fmt.Errorf("event bus '%s' not found", name)
	}

	return bus, nil
}

// RemoveBus 移除命名的事件总线
func (manager *TypedEventBusManager[T]) RemoveBus(name string) error {
	if atomic.LoadInt32(&manager.closed) == 1 {
		return fmt.Errorf("manager is closed")
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	// 再次检查关闭状态（双重检查）
	if atomic.LoadInt32(&manager.closed) == 1 {
		return fmt.Errorf("manager is closed")
	}

	bus, exists := manager.buses[name]
	if !exists {
		return fmt.Errorf("event bus '%s' not found", name)
	}

	// 关闭事件总线
	bus.Close()
	delete(manager.buses, name)
	manager.logger.Debug("Removed event bus: %s", name)
	return nil
}

// ListBuses 列出所有事件总线名称
func (manager *TypedEventBusManager[T]) ListBuses() []string {
	if atomic.LoadInt32(&manager.closed) == 1 {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	names := make([]string, 0, len(manager.buses))
	for name := range manager.buses {
		names = append(names, name)
	}
	return names
}

// CloseAll 关闭所有事件总线
func (manager *TypedEventBusManager[T]) CloseAll() {
	// 原子性设置关闭状态
	if !atomic.CompareAndSwapInt32(&manager.closed, 0, 1) {
		return // 已经关闭
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	for name, bus := range manager.buses {
		bus.Close()
		manager.logger.Debug("Closed event bus: %s", name)
	}
	manager.buses = make(map[string]EventBus[T])
	manager.logger.Info("All event buses closed")
}