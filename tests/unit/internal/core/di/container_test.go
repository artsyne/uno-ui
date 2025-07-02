package di_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
)

// 测试用的接口和实现
type TestService interface {
	GetName() string
}

type TestServiceImpl struct {
	name string
}

func (t *TestServiceImpl) GetName() string {
	return t.name
}

// 测试用的工厂函数
func createTestService(c di.ContainerInterface) (TestService, error) {
	return &TestServiceImpl{name: "test"}, nil
}

func TestContainer_BasicOperations(t *testing.T) {
	container := di.NewContainer()
	ctx := context.Background()

	// 测试注册
	err := container.RegisterSingleton(reflect.TypeOf((*TestService)(nil)).Elem(), createTestService)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	// 测试初始化
	err = container.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize container: %v", err)
	}

	// 测试解析
	instance, err := container.Resolve(reflect.TypeOf((*TestService)(nil)).Elem())
	if err != nil {
		t.Fatalf("Failed to resolve service: %v", err)
	}

	service, ok := instance.(TestService)
	if !ok {
		t.Fatalf("Failed to cast to TestService")
	}

	if service.GetName() != "test" {
		t.Errorf("Expected 'test', got '%s'", service.GetName())
	}

	// 测试类型安全的解析
	service2, err := di.Resolve[TestService](container)
	if err != nil {
		t.Fatalf("Failed to resolve service with type-safe function: %v", err)
	}

	if service2.GetName() != "test" {
		t.Errorf("Expected 'test', got '%s'", service2.GetName())
	}
}

func TestContainer_TypeSafeFunctions(t *testing.T) {
	container := di.NewContainer()

	// 测试类型安全的注册
	err := di.RegisterSingleton[TestService](container, createTestService)
	if err != nil {
		t.Fatalf("Failed to register service with type-safe function: %v", err)
	}

	// 测试类型安全的注册检查
	if !di.IsRegistered[TestService](container) {
		t.Error("Service should be registered")
	}

	// 测试类型安全的解析
	service, err := di.Resolve[TestService](container)
	if err != nil {
		t.Fatalf("Failed to resolve service: %v", err)
	}

	if service.GetName() != "test" {
		t.Errorf("Expected 'test', got '%s'", service.GetName())
	}
}

func TestContainer_ErrorHandling(t *testing.T) {
	container := di.NewContainer()

	// 测试解析未注册的服务
	_, err := container.Resolve(reflect.TypeOf((*TestService)(nil)).Elem())
	if err == nil {
		t.Error("Expected error when resolving unregistered service")
	}

	// 测试重复注册
	err = container.RegisterSingleton(reflect.TypeOf((*TestService)(nil)).Elem(), createTestService)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	err = container.RegisterSingleton(reflect.TypeOf((*TestService)(nil)).Elem(), createTestService)
	if err == nil {
		t.Error("Expected error when registering duplicate service")
	}
}

// 性能基准测试
func BenchmarkContainer_Resolve(b *testing.B) {
	container := di.NewContainer()

	// 注册服务
	err := di.RegisterSingleton[TestService](container, createTestService)
	if err != nil {
		b.Fatalf("Failed to register service: %v", err)
	}

	// 初始化容器
	ctx := context.Background()
	err = container.Initialize(ctx)
	if err != nil {
		b.Fatalf("Failed to initialize container: %v", err)
	}

	b.ResetTimer()

	// 基准测试解析性能
	for i := 0; i < b.N; i++ {
		_, err := di.Resolve[TestService](container)
		if err != nil {
			b.Fatalf("Failed to resolve service: %v", err)
		}
	}
}

func BenchmarkContainer_Register(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		container := di.NewContainer()

		err := di.RegisterSingleton[TestService](container, createTestService)
		if err != nil {
			b.Fatalf("Failed to register service: %v", err)
		}
	}
}
