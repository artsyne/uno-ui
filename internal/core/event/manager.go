// Package event 提供事件驱动架构的核心功能
// 事件总线管理器实现
package event

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
)

// NewTypedEventBusManager 创建新的类型化事件总线管理器
func NewTypedEventBusManager[T Event](logger *slog.Logger) *TypedEventBusManager[T] {
	if logger == nil {
		logger = slog.Default()
	}
	return &TypedEventBusManager[T]{
		buses:  make(map[string]EventBus[T]),
		mu:     sync.RWMutex{},
		closed: 0,
		logger: logger,
	}
}

// CreateBus 创建并存储命名的事件总线
func (manager *TypedEventBusManager[T]) CreateBus(name string, logger *slog.Logger) error {
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
	manager.logger.Debug("Created event bus", "bus_name", name)
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
	manager.logger.Debug("Removed event bus", "bus_name", name)
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
		manager.logger.Debug("Closed event bus", "bus_name", name)
	}
	manager.buses = make(map[string]EventBus[T])
	manager.logger.Info("All event buses closed")
}

// IsClosed 检查管理器是否已关闭
func (manager *TypedEventBusManager[T]) IsClosed() bool {
	return atomic.LoadInt32(&manager.closed) == 1
}