package core

import (
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	testhelpers "github.com/biuren/hebe/tests/helpers"
)

// TestDIContainer 测试DI容器基本功能
func TestDIContainer(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		// 测试容器创建
		container := suite.Container
		if container == nil {
			t.Fatal("DI容器不应该为nil")
		}

		// 测试瞬态服务注册和获取
		err := di.Bind(container, func() string {
			return "test-service"
		})
		if err != nil {
			t.Fatalf("注册瞬态服务失败: %v", err)
		}

		// 验证服务是否已注册
		if !di.Has[string](container) {
			t.Error("服务应该已注册")
		}

		// 获取服务
		service, err := di.TryGet[string](container)
		if err != nil {
			t.Fatalf("获取服务失败: %v", err)
		}

		if service != "test-service" {
			t.Errorf("期望服务值为 'test-service'，实际得到: %s", service)
		}

		// 测试重复注册检查
		err = di.Bind(container, func() string {
			return "duplicate-service"
		})
		if err == nil {
			t.Error("重复注册应该返回错误")
		}
	})
}

// TestDISingleton 测试DI容器单例功能
func TestDISingleton(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		container := suite.Container

		// 计数器用于验证单例只创建一次
		var createCount int

		// 注册单例服务
		err := di.Singleton(container, func() *int {
			createCount++
			return &createCount
		})
		if err != nil {
			t.Fatalf("注册单例服务失败: %v", err)
		}

		// 第一次获取
		instance1, err := di.TryGet[*int](container)
		if err != nil {
			t.Fatalf("获取单例服务失败: %v", err)
		}

		// 第二次获取
		instance2, err := di.TryGet[*int](container)
		if err != nil {
			t.Fatalf("第二次获取单例服务失败: %v", err)
		}

		// 验证是同一个实例
		if instance1 != instance2 {
			t.Error("单例服务应该返回相同的实例")
		}

		// 验证只创建了一次
		if *instance1 != 1 {
			t.Errorf("单例服务应该只创建一次，实际创建了 %d 次", *instance1)
		}
	})
}

// TestDIValue 测试DI容器实例注册
func TestDIValue(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		container := suite.Container

		// 创建一个实例
		testValue := "pre-created-instance"

		// 注册实例
		err := di.Value(container, testValue)
		if err != nil {
			t.Fatalf("注册实例失败: %v", err)
		}

		// 获取实例
		retrievedValue, err := di.TryGet[string](container)
		if err != nil {
			t.Fatalf("获取实例失败: %v", err)
		}

		if retrievedValue != testValue {
			t.Errorf("期望得到 '%s'，实际得到: '%s'", testValue, retrievedValue)
		}
	})
}

// TestDIContainerClear 测试DI容器清空功能
func TestDIContainerClear(t *testing.T) {
	container := di.New()

	// 注册多种类型的服务
	di.Bind(container, func() string { return "transient" })
	di.Singleton(container, func() int { return 42 })
	di.Value(container, true)

	// 验证服务已注册
	if !di.Has[string](container) || !di.Has[int](container) || !di.Has[bool](container) {
		t.Error("服务应该已注册")
	}

	// 清空容器
	container.Clear()

	// 验证服务已清空
	if di.Has[string](container) || di.Has[int](container) || di.Has[bool](container) {
		t.Error("清空后容器不应该有任何服务")
	}
}
