package di_test

import (
	"testing"

	"github.com/biuren/hebe/internal/core/di"
)

// 测试接口
type UserService interface {
	GetName() string
}

type DatabaseService interface {
	Connect() bool
}

// 实现
type userServiceImpl struct {
	name string
	db   DatabaseService
}

func (u *userServiceImpl) GetName() string {
	return u.name
}

type databaseServiceImpl struct {
	connected bool
}

func (d *databaseServiceImpl) Connect() bool {
	return d.connected
}

// 功能测试
func TestNewDI(t *testing.T) {
	container := di.New()

	// 测试瞬态服务
	di.Bind[UserService](container, func() UserService {
		return &userServiceImpl{name: "alice"}
	})

	user1 := di.Get[UserService](container)
	user2 := di.Get[UserService](container)

	if user1 == user2 {
		t.Error("瞬态服务应该返回不同实例")
	}

	if user1.GetName() != "alice" {
		t.Error("服务功能不正确")
	}

	// 测试单例服务
	di.Singleton[DatabaseService](container, func() DatabaseService {
		return &databaseServiceImpl{connected: true}
	})

	db1 := di.Get[DatabaseService](container)
	db2 := di.Get[DatabaseService](container)

	if db1 != db2 {
		t.Error("单例服务应该返回相同实例")
	}

	if !db1.Connect() {
		t.Error("数据库服务功能不正确")
	}

	// 测试实例注册
	config := &userServiceImpl{name: "config"}
	di.Value[*userServiceImpl](container, config)

	retrieved := di.Get[*userServiceImpl](container)
	if retrieved != config {
		t.Error("实例注册失败")
	}

	// 测试Has函数
	if !di.Has[UserService](container) {
		t.Error("Has函数检测失败")
	}

	if di.Has[string](container) {
		t.Error("Has函数误报")
	}
}

// 依赖注入测试
func TestDependencyInjection(t *testing.T) {
	container := di.New()

	// 注册依赖
	di.Singleton[DatabaseService](container, func() DatabaseService {
		return &databaseServiceImpl{connected: true}
	})

	// 注册依赖于DatabaseService的服务
	di.Bind[UserService](container, func() UserService {
		db := di.Get[DatabaseService](container)
		return &userServiceImpl{name: "bob", db: db}
	})

	user := di.Get[UserService](container)
	if user.GetName() != "bob" {
		t.Error("依赖注入失败")
	}
}

// 基准测试
func BenchmarkBind(b *testing.B) {
	container := di.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		di.Bind[UserService](container, func() UserService {
			return &userServiceImpl{name: "test"}
		})
		container.Clear()
		container = di.New()
	}
}

func BenchmarkGet(b *testing.B) {
	container := di.New()
	di.Bind[UserService](container, func() UserService {
		return &userServiceImpl{name: "test"}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = di.Get[UserService](container)
	}
}

func BenchmarkGetSingleton(b *testing.B) {
	container := di.New()
	di.Singleton[UserService](container, func() UserService {
		return &userServiceImpl{name: "singleton"}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = di.Get[UserService](container)
	}
}

func BenchmarkHas(b *testing.B) {
	container := di.New()
	di.Bind[UserService](container, func() UserService {
		return &userServiceImpl{name: "test"}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = di.Has[UserService](container)
	}
}
