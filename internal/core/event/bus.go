// Package event 提供全局事件总线管理
package event

import (
	"sync"
	"sync/atomic"
)

// GlobalEventBusRegistry 全局事件总线注册表
type GlobalEventBusRegistry struct {
	// managers 存储不同类型的事件总线管理器
	managers map[string]interface{}
	// mu 保护并发操作的互斥锁
	mu sync.RWMutex
	// closed 标识注册表是否已关闭
	closed int32
}

var (
	// globalRegistry 全局事件总线注册表实例
	globalRegistry = &GlobalEventBusRegistry{
		managers: make(map[string]interface{}),
	}
)

// GetTypedEventBusManager 获取特定类型的事件总线管理器
// 如果不存在则创建一个新的
func GetTypedEventBusManager[T Event](eventTypeName string, logger Logger) *TypedEventBusManager[T] {
	if atomic.LoadInt32(&globalRegistry.closed) == 1 {
		return nil
	}

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// 再次检查关闭状态（双重检查）
	if atomic.LoadInt32(&globalRegistry.closed) == 1 {
		return nil
	}

	if manager, exists := globalRegistry.managers[eventTypeName]; exists {
		if typedManager, ok := manager.(*TypedEventBusManager[T]); ok {
			return typedManager
		}
	}

	// 创建新的类型化管理器
	manager := NewTypedEventBusManager[T](logger)
	globalRegistry.managers[eventTypeName] = manager
	return manager
}

// CloseAllManagers 关闭所有事件总线管理器
func CloseAllManagers() {
	// 原子性设置关闭状态
	if !atomic.CompareAndSwapInt32(&globalRegistry.closed, 0, 1) {
		return // 已经关闭
	}

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	for _, manager := range globalRegistry.managers {
		if mgr, ok := manager.(interface{ CloseAll() }); ok {
			mgr.CloseAll()
		}
	}

	// 清空注册表
	globalRegistry.managers = make(map[string]interface{})
}