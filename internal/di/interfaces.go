package di

import (
	"context"
	"reflect"
)

// ContainerInterface 定义依赖注入容器的核心接口
//
// 这是一个轻量级的DI容器接口，提供基本的依赖注入功能：
// - 服务注册（普通服务、单例服务、实例服务）
// - 服务解析（支持类型安全）
// - 容器生命周期管理
type ContainerInterface interface {
	// Register 注册瞬态依赖到容器中
	// 每次解析时都会创建新的实例
	Register(serviceType reflect.Type, factory interface{}) error

	// RegisterSingleton 注册单例依赖
	// 在容器生命周期内只会创建一个实例
	RegisterSingleton(serviceType reflect.Type, factory interface{}) error

	// RegisterInstance 注册已创建的实例
	// 直接注册现有实例，等同于单例模式
	RegisterInstance(serviceType reflect.Type, instance interface{}) error

	// Resolve 解析依赖
	// 根据服务类型返回对应的实例
	Resolve(serviceType reflect.Type) (interface{}, error)

	// IsRegistered 检查类型是否已注册
	IsRegistered(serviceType reflect.Type) bool

	// Initialize 初始化容器
	// 预加载所有单例服务，确保容器状态正确
	Initialize(ctx context.Context) error

	// Shutdown 优雅关闭容器
	// 清理资源，重置容器状态
	Shutdown(ctx context.Context) error
}

// --- 类型安全的包级别泛型函数 ---

// Register 类型安全的瞬态服务注册函数
func Register[T any](container ContainerInterface, factory func(ContainerInterface) (T, error)) error {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()

	wrappedFactory := func(c ContainerInterface) (interface{}, error) {
		return factory(c)
	}

	return container.Register(serviceType, wrappedFactory)
}

// RegisterSingleton 类型安全的单例服务注册函数
func RegisterSingleton[T any](container ContainerInterface, factory func(ContainerInterface) (T, error)) error {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()

	wrappedFactory := func(c ContainerInterface) (interface{}, error) {
		return factory(c)
	}

	return container.RegisterSingleton(serviceType, wrappedFactory)
}

// RegisterInstance 类型安全的实例注册函数
func RegisterInstance[T any](container ContainerInterface, instance T) error {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()
	return container.RegisterInstance(serviceType, instance)
}

// Resolve 类型安全的依赖解析函数
func Resolve[T any](container ContainerInterface) (T, error) {
	var zero T
	serviceType := reflect.TypeOf((*T)(nil)).Elem()

	instance, err := container.Resolve(serviceType)
	if err != nil {
		return zero, err
	}

	result, ok := instance.(T)
	if !ok {
		return zero, NewTypeError(serviceType, reflect.TypeOf(instance))
	}

	return result, nil
}

// IsRegistered 类型安全的注册检查函数
func IsRegistered[T any](container ContainerInterface) bool {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()
	return container.IsRegistered(serviceType)
}
