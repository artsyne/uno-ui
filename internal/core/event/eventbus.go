// Package event 提供事件驱动架构的核心功能
// 事件总线实现
package event

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// subscriptionIDCounter 全局订阅ID计数器
var subscriptionIDCounter uint64

// NewEventBus 创建新的泛型事件总线
func NewEventBus[T Event](logger *slog.Logger) EventBus[T] {
	if logger == nil {
		logger = slog.Default()
	}
	return &eventBus[T]{
		subscriptions:  make([]subscription[T], 0),
		logger:         logger,
		handlerTimeout: 30 * time.Second, // 默认30秒超时
	}
}

// NewEventBusWithTimeout 创建带自定义超时的事件总线
func NewEventBusWithTimeout[T Event](logger *slog.Logger, timeout time.Duration) EventBus[T] {
	if logger == nil {
		logger = slog.Default()
	}
	return &eventBus[T]{
		subscriptions:  make([]subscription[T], 0),
		logger:         logger,
		handlerTimeout: timeout,
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
	eb.logger.Debug("Publishing event", "event_type", eventType)

	if len(subscriptionsCopy) == 0 {
		eb.logger.Debug("No handlers found for event type", "event_type", eventType)
		return nil
	}

	// 并发执行所有处理器
	var wg sync.WaitGroup
	errorChan := make(chan error, len(subscriptionsCopy))

	for _, sub := range subscriptionsCopy {
		wg.Add(1)
		go func(s subscription[T]) {
			defer wg.Done()
			
			// 创建带超时的上下文
			handlerCtx, cancel := context.WithTimeout(ctx, eb.handlerTimeout)
			defer cancel()

			// 在单个Goroutine中处理，避免嵌套Goroutine
			var err error
			var completed bool
			
			// 使用内部函数处理panic恢复
			func() {
				defer func() {
					if r := recover(); r != nil {
						eb.logger.Error("Panic in event handler", "handler_id", s.id, "error", r)
						err = fmt.Errorf("panic in event handler (ID: %d): %v", s.id, r)
					}
					completed = true
				}()
				err = s.handler(handlerCtx, event)
			}()

			// 检查是否超时
			select {
			case <-handlerCtx.Done():
				if !completed {
					// 处理器还没完成就超时了
					eb.logger.Error("Event handler timeout", "handler_id", s.id)
					errorChan <- fmt.Errorf("handler (ID: %d) timeout", s.id)
				} else if err != nil {
					// 处理器完成了但有错误，同时也超时了
					eb.logger.Error("Event handler error and timeout", "handler_id", s.id, "error", err)
					errorChan <- fmt.Errorf("handler (ID: %d) error with timeout: %w", s.id, err)
				}
				// 如果completed=true且err=nil，说明处理成功但刚好到达超时边界
			default:
				// 没有超时
				if err != nil {
					eb.logger.Error("Error in event handler", "handler_id", s.id, "error", err)
					errorChan <- fmt.Errorf("handler (ID: %d) error: %w", s.id, err)
				}
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
		eb.logger.Error("Event handling errors", "error_count", len(errors), "errors", errors)
		return fmt.Errorf("event handling failed with %d errors: %v", len(errors), errors)
	}

	eb.logger.Debug("Event published successfully", "event_type", eventType)
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

	eb.logger.Debug("Handler subscribed successfully", "subscription_id", id)
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
			eb.logger.Debug("Handler unsubscribed successfully", "subscription_id", id)
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