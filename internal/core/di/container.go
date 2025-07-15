package di

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	errors "github.com/biuren/hebe/internal/core/error"
)

// Container 轻量级依赖注入容器
//
// 提供基本的依赖注入功能，包括：
// - 瞬态服务注册（每次解析创建新实例）
// - 单例服务注册（容器生命周期内只创建一个实例）
// - 实例注册（直接注册现有实例）
// - 类型安全的服务解析
// - 容器生命周期管理
type Container struct {
	mu          sync.RWMutex                 // 读写锁，保护并发访问
	providers   map[reflect.Type]*provider   // 服务提供者映射
	singletons  map[reflect.Type]interface{} // 单例实例缓存
	initialized bool                         // 初始化状态标志
}

// NewContainer 创建新的DI容器
//
// 返回一个初始化完成的空容器，可以立即开始注册服务
func NewContainer() *Container {
	return &Container{
		providers:  make(map[reflect.Type]*provider),
		singletons: make(map[reflect.Type]interface{}),
	}
}

// Register 注册瞬态依赖
//
// 每次解析时都会调用factory函数创建新的实例
// 适用于有状态的服务或需要每次都创建新实例的场景
//
// 参数：
//
//	serviceType: 服务类型的反射信息
//	factory: 工厂函数，接受容器参数并返回服务实例
//
// 返回：
//
//	error: 注册失败时的错误信息
func (c *Container) Register(serviceType reflect.Type, factory interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return errors.DIContainerNotInitialized.Err().WithDetails("container is already initialized")
	}

	if _, exists := c.providers[serviceType]; exists {
		return newAlreadyRegisteredError(serviceType)
	}

	// 包装工厂函数
	wrappedFactory, err := WrapFactory(factory)
	if err != nil {
		return err
	}

	provider := newProvider(serviceType, wrappedFactory, false)
	c.providers[serviceType] = provider

	return nil
}

// RegisterSingleton 注册单例依赖
//
// 在容器生命周期内只会创建一个实例，后续解析都返回同一个实例
// 适用于无状态服务或需要共享状态的服务
//
// 参数：
//
//	serviceType: 服务类型的反射信息
//	factory: 工厂函数，接受容器参数并返回服务实例
//
// 返回：
//
//	error: 注册失败时的错误信息
func (c *Container) RegisterSingleton(serviceType reflect.Type, factory interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return errors.DIContainerNotInitialized.Err().WithDetails("container is already initialized")
	}

	if _, exists := c.providers[serviceType]; exists {
		return newAlreadyRegisteredError(serviceType)
	}

	// 包装工厂函数
	wrappedFactory, err := WrapFactory(factory)
	if err != nil {
		return err
	}

	provider := newProvider(serviceType, wrappedFactory, true)
	c.providers[serviceType] = provider

	return nil
}

// RegisterInstance 注册已创建的实例
//
// 直接注册现有实例，等同于单例模式
// 适用于已经创建好的实例或外部创建的对象
//
// 参数：
//
//	serviceType: 服务类型的反射信息
//	instance: 要注册的实例
//
// 返回：
//
//	error: 注册失败时的错误信息
func (c *Container) RegisterInstance(serviceType reflect.Type, instance interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return errors.DIContainerNotInitialized.Err().WithDetails("container is already initialized")
	}

	if _, exists := c.providers[serviceType]; exists {
		return newAlreadyRegisteredError(serviceType)
	}

	// 创建单例提供者
	provider := newProvider(serviceType, func(c ContainerInterface) (interface{}, error) {
		return instance, nil
	}, true)

	c.providers[serviceType] = provider
	c.singletons[serviceType] = instance

	return nil
}

// Resolve 解析依赖
//
// 根据服务类型返回对应的实例
// 对于瞬态服务，每次调用都会创建新实例
// 对于单例服务，会返回缓存的实例
//
// 参数：
//
//	serviceType: 要解析的服务类型
//
// 返回：
//
//	interface{}: 解析得到的服务实例
//	error: 解析失败时的错误信息
func (c *Container) Resolve(serviceType reflect.Type) (interface{}, error) {
	c.mu.RLock()
	provider, exists := c.providers[serviceType]
	c.mu.RUnlock()

	if !exists {
		return nil, newNotRegisteredError(serviceType)
	}

	return provider.provide(c)
}

// IsRegistered 检查类型是否已注册
//
// 检查指定类型的服务是否已在容器中注册
//
// 参数：
//
//	serviceType: 要检查的服务类型
//
// 返回：
//
//	bool: 如果已注册返回true，否则返回false
func (c *Container) IsRegistered(serviceType reflect.Type) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.providers[serviceType]
	return exists
}

// Initialize 初始化容器
//
// 预加载所有单例服务，确保容器状态正确
// 这一步会创建所有已注册的单例实例，如果有任何单例创建失败，初始化会失败
//
// 参数：
//
//	ctx: 上下文，可用于超时控制
//
// 返回：
//
//	error: 初始化失败时的错误信息
func (c *Container) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return errors.DIContainerNotInitialized.Err().WithDetails("container is already initialized")
	}

	// 预加载所有单例
	for serviceType, provider := range c.providers {
		if provider.isSingleton() {
			instance, err := provider.provide(c)
			if err != nil {
				return errors.DIContainerNotInitialized.Err().WithDetails(
					fmt.Sprintf("container initialization failed: %v", err),
				)
			}
			c.singletons[serviceType] = instance
		}
	}

	c.initialized = true
	return nil
}

// Shutdown 优雅关闭容器
//
// 清理资源，重置容器状态
// 调用后容器将回到未初始化状态，可以重新初始化
//
// 参数：
//
//	ctx: 上下文，可用于超时控制
//
// 返回：
//
//	error: 关闭失败时的错误信息
func (c *Container) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return errors.DIContainerNotInitialized.Err().WithDetails("container is not initialized")
	}

	// 清理单例实例
	c.singletons = make(map[reflect.Type]interface{})
	c.initialized = false

	return nil
}

// newAlreadyRegisteredError 创建已注册错误
func newAlreadyRegisteredError(serviceType reflect.Type) error {
	return errors.DIServiceAlreadyRegistered.WithDetails(
		fmt.Sprintf("service type %s is already registered", serviceType.String()),
	)
}

// newNotRegisteredError 创建未注册错误
func newNotRegisteredError(serviceType reflect.Type) error {
	return errors.DIServiceNotRegistered.WithDetails(
		fmt.Sprintf("service type %s is not registered", serviceType.String()),
	)
}
