package di_test

import (
	"reflect"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
)

// 反射版本的服务注册和解析（模拟老版本）
type ReflectionContainer struct {
	services map[reflect.Type]any
}

func NewReflectionContainer() *ReflectionContainer {
	return &ReflectionContainer{
		services: make(map[reflect.Type]any),
	}
}

func (c *ReflectionContainer) RegisterSingleton(instance any) {
	t := reflect.TypeOf(instance)
	c.services[t] = instance
}

func (c *ReflectionContainer) Resolve(serviceType reflect.Type) any {
	return c.services[serviceType]
}

// 测试服务结构体
type PerformanceService struct {
	Name string
	ID   int
}

// 性能对比基准测试 - 新的泛型版本
func BenchmarkGenericContainer_Register(b *testing.B) {
	container := di.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := di.Singleton(container, func() *PerformanceService {
			return &PerformanceService{Name: "test", ID: i}
		}); err != nil {
			b.Fatal(err)
		}
		container.Clear()
		container = di.New()
	}
}

func BenchmarkGenericContainer_Resolve(b *testing.B) {
	container := di.New()
	if err := di.Singleton(container, func() *PerformanceService {
		return &PerformanceService{Name: "test", ID: 123}
	}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = di.Get[*PerformanceService](container)
	}
}

// 性能对比基准测试 - 反射版本
func BenchmarkReflectionContainer_Register(b *testing.B) {
	container := NewReflectionContainer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service := &PerformanceService{Name: "test", ID: i}
		container.RegisterSingleton(service)
		container = NewReflectionContainer()
	}
}

func BenchmarkReflectionContainer_Resolve(b *testing.B) {
	container := NewReflectionContainer()
	service := &PerformanceService{Name: "test", ID: 123}
	container.RegisterSingleton(service)
	serviceType := reflect.TypeOf((*PerformanceService)(nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container.Resolve(serviceType).(*PerformanceService)
	}
}

// 内存分配对比测试
func BenchmarkGenericContainer_MemoryAllocation(b *testing.B) {
	container := di.New()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := di.Singleton(container, func() *PerformanceService {
			return &PerformanceService{Name: "test", ID: i}
		}); err != nil {
			b.Fatal(err)
		}
		_ = di.Get[*PerformanceService](container)
		container.Clear()
		container = di.New()
	}
}

func BenchmarkReflectionContainer_MemoryAllocation(b *testing.B) {
	serviceType := reflect.TypeOf((*PerformanceService)(nil))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		container := NewReflectionContainer()
		service := &PerformanceService{Name: "test", ID: i}
		container.RegisterSingleton(service)
		_ = container.Resolve(serviceType).(*PerformanceService)
	}
}

// 并发性能测试
func BenchmarkGenericContainer_Concurrent(b *testing.B) {
	container := di.New()
	if err := di.Singleton(container, func() *PerformanceService {
		return &PerformanceService{Name: "concurrent", ID: 999}
	}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = di.Get[*PerformanceService](container)
		}
	})
}

func BenchmarkReflectionContainer_Concurrent(b *testing.B) {
	container := NewReflectionContainer()
	service := &PerformanceService{Name: "concurrent", ID: 999}
	container.RegisterSingleton(service)
	serviceType := reflect.TypeOf((*PerformanceService)(nil))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = container.Resolve(serviceType).(*PerformanceService)
		}
	})
}
