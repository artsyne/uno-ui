package di

import (
	"fmt"
	"sync"

	errors "github.com/biuren/hebe/internal/utils/error"
)

/*
Package di 提供了一个基于Go 1.24泛型的高性能依赖注入容器

为什么依赖注入是最佳实践:
1. 控制反转(IoC): 对象不再负责创建其依赖，而是由外部容器注入
2. 松耦合: 组件间通过接口依赖，而非具体实现
3. 易测试: 可以轻松注入mock对象进行单元测试
4. 配置集中: 所有依赖关系在一处配置，便于管理
5. 生命周期管理: 统一管理对象的创建和销毁

设计原则:
- 零反射: 使用泛型替代反射，提供编译时类型安全
- 高性能: 基于map的快速查找，单例使用sync.Once优化
- 简洁API: 最少的概念，最直观的使用方式
- 并发安全: 使用读写锁保证线程安全
- 错误明确: 重复注册时明确报错，避免意外覆盖

使用示例:

    container := di.New()

    // 注册服务（会检查重复）
    if err := di.Bind(container, func() UserService {
        return &UserServiceImpl{db: GetDB()}
    }); err != nil {
        log.Fatal(err)
    }

    // 获取服务
    userService := di.Get[UserService](container)
*/

// Container 高性能DI容器
//
// 基于泛型设计的依赖注入容器，提供三种服务生命周期管理:
// - factories: 瞬态服务，每次调用都创建新实例
// - singletons: 单例服务，全局唯一实例，使用sync.Once确保线程安全的延迟初始化
// - instances: 预创建实例，直接返回注册的对象
// - cache: 单例实例缓存
// - once: 单例初始化控制器
//
// 所有操作都是并发安全的，使用读写锁优化性能
type Container struct {
	mu sync.RWMutex // 读写锁，优化并发性能

	factories  map[string]func() any // 瞬态服务工厂函数
	singletons map[string]func() any // 单例服务工厂函数
	instances  map[string]any        // 预创建实例存储
	cache      map[string]any        // 单例实例缓存
	once       map[string]*sync.Once // 单例初始化控制器
}

// New 创建新的DI容器
//
// 创建一个空的依赖注入容器，所有map都已初始化
// 容器是并发安全的，可以在多个goroutine中同时使用
//
// 返回:
//
//	*Container: 新的DI容器实例
func New() *Container {
	return &Container{
		factories:  make(map[string]func() any),
		singletons: make(map[string]func() any),
		instances:  make(map[string]any),
		cache:      make(map[string]any),
		once:       make(map[string]*sync.Once),
	}
}

// key 生成类型的唯一标识符
//
// 使用Go泛型在编译时生成类型的字符串表示作为map的key
// 这种方式避免了反射的性能开销，同时保证了类型安全
//
// 类型参数:
//
//	T: 要生成key的类型
//
// 返回:
//
//	string: 类型的唯一字符串标识
func key[T any]() string {
	return fmt.Sprintf("%T", (*T)(nil))
}

// checkConflict 检查服务是否已经注册
//
// 在注册新服务前检查是否存在冲突，防止意外覆盖
// 这是一个内部函数，假设调用时已经持有写锁
//
// 类型参数:
//
//	T: 要检查的服务类型
//
// 参数:
//
//	c: DI容器
//	k: 类型键
//
// 返回:
//
//	error: 如果存在冲突则返回详细错误信息
func checkConflict[T any](c *Container, k string) error {
	// 检查各种注册类型是否存在冲突
	if _, exists := c.factories[k]; exists {
		return errors.DIServiceAlreadyRegistered.WithDetails(
			fmt.Sprintf("transient service %T already registered", (*T)(nil)),
		)
	}
	if _, exists := c.singletons[k]; exists {
		return errors.DIServiceAlreadyRegistered.WithDetails(
			fmt.Sprintf("singleton service %T already registered", (*T)(nil)),
		)
	}
	if _, exists := c.instances[k]; exists {
		return errors.DIServiceAlreadyRegistered.WithDetails(
			fmt.Sprintf("instance service %T already registered", (*T)(nil)),
		)
	}
	return nil
}

// Bind 注册瞬态服务
//
// 瞬态服务每次获取时都会调用factory函数创建新实例
// 适用于无状态服务或需要独立实例的场景
//
// 注意: 如果该类型的服务已经注册（任何生命周期），将返回错误
// 这样设计是为了防止意外覆盖，确保服务注册的明确性
//
// 示例:
//
//	err := di.Bind(container, func() UserService {
//	    return &UserServiceImpl{id: uuid.New()}
//	})
//	if err != nil {
//	    log.Fatal("注册失败:", err)
//	}
//
// 类型参数:
//
//	T: 服务类型
//
// 参数:
//
//	c: DI容器
//	factory: 服务创建函数，每次调用都应返回新实例
//
// 返回:
//
//	error: 如果服务已注册或其他错误则返回详细错误信息
func Bind[T any](c *Container, factory func() T) error {
	k := key[T]()
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否存在注册冲突
	if err := checkConflict[T](c, k); err != nil {
		return err
	}

	// 注册瞬态服务
	c.factories[k] = func() any { return factory() }
	return nil
}

// Singleton 注册单例服务
//
// 单例服务在整个应用生命周期中只创建一次实例
// 使用sync.Once确保在并发环境下的线程安全初始化
// 适用于数据库连接、配置对象等共享资源
//
// 注意: 如果该类型的服务已经注册（任何生命周期），将返回错误
//
// 示例:
//
//	err := di.Singleton(container, func() DatabaseService {
//	    return &PostgreSQLService{dsn: os.Getenv("DB_DSN")}
//	})
//	if err != nil {
//	    log.Fatal("注册失败:", err)
//	}
//
// 类型参数:
//
//	T: 服务类型
//
// 参数:
//
//	c: DI容器
//	factory: 服务创建函数，只会被调用一次
//
// 返回:
//
//	error: 如果服务已注册或其他错误则返回详细错误信息
func Singleton[T any](c *Container, factory func() T) error {
	k := key[T]()
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否存在注册冲突
	if err := checkConflict[T](c, k); err != nil {
		return err
	}

	// 注册单例服务
	c.singletons[k] = func() any { return factory() }
	c.once[k] = &sync.Once{}
	return nil
}

// Value 注册已创建的实例
//
// 直接注册一个已经创建好的对象实例
// 适用于配置对象、常量等预先创建的实例
//
// 注意: 如果该类型的服务已经注册（任何生命周期），将返回错误
//
// 示例:
//
//	config := Config{Env: "production", Debug: false}
//	err := di.Value(container, config)
//	if err != nil {
//	    log.Fatal("注册失败:", err)
//	}
//
// 类型参数:
//
//	T: 服务类型
//
// 参数:
//
//	c: DI容器
//	instance: 要注册的实例对象
//
// 返回:
//
//	error: 如果服务已注册或其他错误则返回详细错误信息
func Value[T any](c *Container, instance T) error {
	k := key[T]()
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否存在注册冲突
	if err := checkConflict[T](c, k); err != nil {
		return err
	}

	// 注册实例
	c.instances[k] = instance
	return nil
}

// Rebind 强制重新注册瞬态服务
//
// 允许覆盖已存在的服务注册，用于特殊场景如测试或动态配置
// 使用时需要明确意图，避免在生产代码中滥用
//
// 示例:
//
//	// 在测试中覆盖生产服务
//	di.Rebind(container, func() UserService {
//	    return &MockUserService{}
//	})
//
// 类型参数:
//
//	T: 服务类型
//
// 参数:
//
//	c: DI容器
//	factory: 服务创建函数
func Rebind[T any](c *Container, factory func() T) {
	k := key[T]()
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清理可能存在的旧注册
	delete(c.singletons, k)
	delete(c.instances, k)
	delete(c.cache, k)
	delete(c.once, k)

	// 注册新的瞬态服务
	c.factories[k] = func() any { return factory() }
}

// Get 获取服务实例
//
// 获取指定类型的服务实例，如果服务未注册则panic
// 这是最常用的服务获取方式，适用于确定服务已注册的场景
//
// 查找优先级: 实例 > 单例 > 瞬态
//
// 示例:
//
//	userService := di.Get[UserService](container)
//	config := di.Get[Config](container)
//
// 类型参数:
//
//	T: 服务类型
//
// 参数:
//
//	c: DI容器
//
// 返回:
//
//	T: 指定类型的服务实例
//
// Panic:
//
//	如果服务未注册，会panic并提供详细错误信息
func Get[T any](c *Container) T {
	result, err := TryGet[T](c)
	if err != nil {
		panic(fmt.Sprintf("service %T not registered", (*T)(nil)))
	}
	return result
}

// TryGet 安全获取服务实例
//
// 获取指定类型的服务实例，如果服务未注册则返回错误而不是panic
// 适用于不确定服务是否已注册，需要错误处理的场景
//
// 查找优先级: 实例 > 单例 > 瞬态
//
// 对于单例服务，使用sync.Once确保只初始化一次，即使在高并发环境下也是线程安全的
//
// 示例:
//
//	userService, err := di.TryGet[UserService](container)
//	if err != nil {
//	    log.Printf("UserService not available: %v", err)
//	    return
//	}
//
// 类型参数:
//
//	T: 服务类型
//
// 参数:
//
//	c: DI容器
//
// 返回:
//
//	T: 指定类型的服务实例（如果找到）
//	error: 如果服务未注册则返回详细错误信息
func TryGet[T any](c *Container) (T, error) {
	var zero T
	k := key[T]()

	c.mu.RLock()

	// 优先级1: 检查预注册实例（最快）
	if v, ok := c.instances[k]; ok {
		c.mu.RUnlock()
		return v.(T), nil
	}

	// 优先级2: 检查单例服务
	if factory, ok := c.singletons[k]; ok {
		once := c.once[k]
		c.mu.RUnlock()

		// 使用sync.Once确保单例只初始化一次
		// 即使在高并发环境下也是线程安全的
		once.Do(func() {
			c.mu.Lock()
			c.cache[k] = factory()
			c.mu.Unlock()
		})

		c.mu.RLock()
		result := c.cache[k]
		c.mu.RUnlock()
		return result.(T), nil
	}

	// 优先级3: 检查瞬态服务（每次都创建新实例）
	if factory, ok := c.factories[k]; ok {
		c.mu.RUnlock()
		return factory().(T), nil
	}

	c.mu.RUnlock()
	return zero, errors.DIServiceNotRegistered.WithDetails(
		fmt.Sprintf("service %T not registered", (*T)(nil)),
	)
}

// Has 检查服务是否已注册
//
// 检查指定类型的服务是否在容器中注册
// 这是一个只读操作，不会触发服务的创建
//
// 示例:
//
//	if di.Has[UserService](container) {
//	    userService := di.Get[UserService](container)
//	    // 使用userService
//	}
//
// 类型参数:
//
//	T: 要检查的服务类型
//
// 参数:
//
//	c: DI容器
//
// 返回:
//
//	bool: 如果服务已注册返回true，否则返回false
func Has[T any](c *Container) bool {
	k := key[T]()
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 检查所有可能的注册类型
	_, inFactories := c.factories[k]
	_, inSingletons := c.singletons[k]
	_, inInstances := c.instances[k]

	return inFactories || inSingletons || inInstances
}

// Clear 清空容器中的所有服务
//
// 移除容器中注册的所有服务，包括瞬态、单例和实例
// 这个操作会清空所有内部map，释放内存
//
// 注意: 这个操作不会调用任何清理方法，如果服务需要特殊的清理逻辑，
// 应该在调用Clear之前手动处理
//
// 使用Go 1.24的clear函数，高效清理map内容
//
// 示例:
//
//	container.Clear()  // 清空所有注册的服务
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 使用Go 1.24的clear函数高效清理所有map
	clear(c.factories)
	clear(c.singletons)
	clear(c.instances)
	clear(c.cache)
	clear(c.once)
}

/*
设计说明:

1. 重复注册检查:
   - Bind、Singleton、Value方法都会检查服务是否已注册
   - 如果已注册则返回明确的错误信息，防止意外覆盖
   - 提供Rebind方法用于明确的覆盖场景

2. 错误处理:
   - 使用项目统一的错误处理机制
   - 提供详细的错误信息，包含类型和冲突原因
   - TryGet返回错误，Get在失败时panic

3. 并发安全:
   - 使用sync.RWMutex优化读多写少场景
   - 单例使用sync.Once确保线程安全初始化
   - 所有公共方法都是并发安全的

4. 性能优化:
   - 零反射设计，使用编译时类型信息
   - 基于字符串键的快速map查找
   - defer确保锁的正确释放

5. API设计:
   - 简洁的泛型API，编译时类型安全
   - 明确的生命周期管理
   - 直观的错误处理模式
*/
