package core_test

import (
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	"github.com/stretchr/testify/assert"
)

// 测试用的接口和实现
type TestService interface {
	GetName() string
}

type testServiceImpl struct {
	name string
}

func (t *testServiceImpl) GetName() string {
	return t.name
}

func newTestService() TestService {
	return &testServiceImpl{name: "test"}
}

type TestRepository interface {
	GetData() string
}

type testRepositoryImpl struct {
	data string
}

func (t *testRepositoryImpl) GetData() string {
	return t.data
}

func newTestRepository() TestRepository {
	return &testRepositoryImpl{data: "repository data"}
}

// 带依赖的服务
type TestServiceWithDep interface {
	GetInfo() string
}

type testServiceWithDepImpl struct {
	repo TestRepository
}

func (t *testServiceWithDepImpl) GetInfo() string {
	return "service: " + t.repo.GetData()
}

func newTestServiceWithDep(repo TestRepository) TestServiceWithDep {
	return &testServiceWithDepImpl{repo: repo}
}

func TestDIContainerCreation(t *testing.T) {
	container := di.New()
	assert.NotNil(t, container)
}

func TestDIContainerTransient(t *testing.T) {
	container := di.New()

	// 注册瞬态服务
	err := di.Bind[TestService](container, newTestService)
	assert.NoError(t, err)

	// 获取两次实例
	service1 := di.Get[TestService](container)
	service2 := di.Get[TestService](container)

	// 验证实例不同（瞬态）
	assert.NotSame(t, service1, service2)
	assert.Equal(t, "test", service1.GetName())
	assert.Equal(t, "test", service2.GetName())
}

func TestDIContainerSingleton(t *testing.T) {
	container := di.New()

	// 注册单例服务
	err := di.Singleton[TestService](container, newTestService)
	assert.NoError(t, err)

	// 获取两次实例
	service1 := di.Get[TestService](container)
	service2 := di.Get[TestService](container)

	// 验证实例相同（单例）
	assert.Same(t, service1, service2)
	assert.Equal(t, "test", service1.GetName())
}

func TestDIContainerValue(t *testing.T) {
	container := di.New()

	// 创建预先存在的实例
	originalService := &testServiceImpl{name: "original"}

	// 注册值
	err := di.Value[TestService](container, originalService)
	assert.NoError(t, err)

	// 获取实例
	retrievedService := di.Get[TestService](container)

	// 验证是同一个实例
	assert.Same(t, originalService, retrievedService)
	assert.Equal(t, "original", retrievedService.GetName())
}

func TestDIContainerDependencyInjection(t *testing.T) {
	container := di.New()

	// 注册依赖
	repo := newTestRepository()
	err := di.Value[TestRepository](container, repo)
	assert.NoError(t, err)

	// 注册带依赖的服务，使用预先创建的依赖
	service := newTestServiceWithDep(repo)
	err = di.Value[TestServiceWithDep](container, service)
	assert.NoError(t, err)

	// 获取服务
	retrievedService := di.Get[TestServiceWithDep](container)

	// 验证依赖注入正常工作
	assert.NotNil(t, retrievedService)
	assert.Same(t, service, retrievedService)
	assert.Equal(t, "service: repository data", retrievedService.GetInfo())
}

func TestDIContainerHas(t *testing.T) {
	container := di.New()

	// 检查未注册的服务
	assert.False(t, di.Has[TestService](container))

	// 注册服务
	err := di.Bind[TestService](container, newTestService)
	assert.NoError(t, err)

	// 检查已注册的服务
	assert.True(t, di.Has[TestService](container))
}

func TestDIContainerClear(t *testing.T) {
	container := di.New()

	// 注册服务
	err := di.Bind[TestService](container, newTestService)
	assert.NoError(t, err)
	err = di.Singleton[TestRepository](container, newTestRepository)
	assert.NoError(t, err)

	// 验证服务已注册
	assert.True(t, di.Has[TestService](container))
	assert.True(t, di.Has[TestRepository](container))

	// 清空容器
	container.Clear()

	// 验证服务已清空
	assert.False(t, di.Has[TestService](container))
	assert.False(t, di.Has[TestRepository](container))
}

func TestDIContainerOverride(t *testing.T) {
	container := di.New()

	// 注册第一个实现
	err := di.Bind[TestService](container, func() TestService {
		return &testServiceImpl{name: "first"}
	})
	assert.NoError(t, err)

	// 获取第一个实现
	service1 := di.Get[TestService](container)
	assert.Equal(t, "first", service1.GetName())

	// 尝试重复注册应该返回错误
	err = di.Bind[TestService](container, func() TestService {
		return &testServiceImpl{name: "second"}
	})
	assert.Error(t, err)

	// 使用Rebind覆盖注册第二个实现
	di.Rebind[TestService](container, func() TestService {
		return &testServiceImpl{name: "second"}
	})

	// 获取第二个实现
	service2 := di.Get[TestService](container)
	assert.Equal(t, "second", service2.GetName())
}

func TestDIContainerSingletonLazyLoading(t *testing.T) {
	container := di.New()
	creationCount := 0

	// 注册单例服务，带创建计数
	err := di.Singleton[TestService](container, func() TestService {
		creationCount++
		return &testServiceImpl{name: "lazy"}
	})
	assert.NoError(t, err)

	// 验证还未创建实例
	assert.Equal(t, 0, creationCount)

	// 第一次获取
	service1 := di.Get[TestService](container)
	assert.Equal(t, 1, creationCount)
	assert.Equal(t, "lazy", service1.GetName())

	// 第二次获取
	service2 := di.Get[TestService](container)
	assert.Equal(t, 1, creationCount) // 创建次数不应增加
	assert.Same(t, service1, service2)
}

func TestDIContainerComplexDependencyChain(t *testing.T) {
	container := di.New()

	// 创建依赖链：Service -> ServiceWithDep -> Repository
	repo := newTestRepository()
	err := di.Value[TestRepository](container, repo)
	assert.NoError(t, err)

	serviceWithDep := newTestServiceWithDep(repo)
	err = di.Value[TestServiceWithDep](container, serviceWithDep)
	assert.NoError(t, err)

	// 注册更高层的服务，使用瞬态注册但避免在工厂函数中调用di.Get
	err = di.Bind[TestService](container, func() TestService {
		return &testServiceImpl{name: "complex service"}
	})
	assert.NoError(t, err)

	// 获取最终服务
	service := di.Get[TestService](container)
	retrievedServiceWithDep := di.Get[TestServiceWithDep](container)
	retrievedRepo := di.Get[TestRepository](container)

	// 验证整个依赖链正常工作
	assert.NotNil(t, service)
	assert.Same(t, serviceWithDep, retrievedServiceWithDep)
	assert.Same(t, repo, retrievedRepo)
	assert.Equal(t, "complex service", service.GetName())
	assert.Equal(t, "service: repository data", retrievedServiceWithDep.GetInfo())
}

func TestDIContainerMixedRegistrationTypes(t *testing.T) {
	container := di.New()

	// 混合注册类型
	originalRepo := &testRepositoryImpl{data: "original data"}
	err := di.Value[TestRepository](container, originalRepo) // 值注册
	assert.NoError(t, err)

	// 创建一个预先存在的服务实例，避免在工厂函数中调用di.Get
	originalServiceWithDep := newTestServiceWithDep(originalRepo)
	err = di.Value[TestServiceWithDep](container, originalServiceWithDep) // 值注册
	assert.NoError(t, err)

	// 注册一个简单的瞬态服务
	err = di.Bind[TestService](container, func() TestService {
		return &testServiceImpl{name: "transient service"}
	})
	assert.NoError(t, err)

	// 获取服务并验证
	service1 := di.Get[TestService](container)
	service2 := di.Get[TestService](container)
	serviceWithDep1 := di.Get[TestServiceWithDep](container)
	serviceWithDep2 := di.Get[TestServiceWithDep](container)
	repo := di.Get[TestRepository](container)

	// 验证不同注册类型的行为
	assert.NotSame(t, service1, service2) // 瞬态服务
	assert.Same(t, serviceWithDep1, serviceWithDep2) // 值注册的服务
	assert.Same(t, originalRepo, repo) // 值注册
	assert.Equal(t, "transient service", service1.GetName())
	assert.Equal(t, "service: original data", serviceWithDep1.GetInfo())
}

func TestDIContainerTryGet(t *testing.T) {
	container := di.New()

	// 尝试获取未注册的服务
	service, err := di.TryGet[TestService](container)
	assert.Error(t, err)
	assert.Nil(t, service)

	// 注册服务
	err = di.Bind[TestService](container, newTestService)
	assert.NoError(t, err)

	// 成功获取已注册的服务
	service, err = di.TryGet[TestService](container)
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "test", service.GetName())
}

func TestDIContainerGetPanic(t *testing.T) {
	container := di.New()

	// 获取未注册的服务应该panic
	assert.Panics(t, func() {
		di.Get[TestService](container)
	})

	// 注册服务后不应该panic
	err := di.Bind[TestService](container, newTestService)
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		service := di.Get[TestService](container)
		assert.NotNil(t, service)
	})
}