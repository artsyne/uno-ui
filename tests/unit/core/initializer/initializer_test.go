package initializer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biuren/hebe/internal/utils/config"
	initializer "github.com/biuren/hebe/internal/core/initializer"
)

func TestNewInitializer(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "test",
			User:     "root",
			Password: "password",
			Conn:     "mysql",
		},
	}

	init := initializer.NewInitializer(cfg)
	assert.NotNil(t, init)
}

func TestInitializeCore_WithoutDatabase(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "nonexistent",
			User:     "root",
			Password: "wrongpassword",
			Conn:     "mysql",
		},
	}

	init := initializer.NewInitializer(cfg)
	assert.NotNil(t, init)

	// 测试初始化核心组件
	err := init.InitializeCore()
	// 由于涉及真实数据库连接，这里主要测试方法是否存在和基本逻辑
	// 在实际项目中，应该使用 mock 数据库
	assert.NotNil(t, err) // 预期会有错误，因为没有真实数据库连接
}

func TestGetContainer(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "test",
			User:     "root",
			Password: "password",
			Conn:     "mysql",
		},
	}

	init := initializer.NewInitializer(cfg)
	assert.NotNil(t, init)

	// 测试获取容器
	container := init.GetContainer()
	// 初始时容器应该为 nil，因为还没有调用 InitializeCore
	assert.Nil(t, container)
}

func TestInitializer_Fields(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "test",
			User:     "root",
			Password: "password",
			Conn:     "mysql",
		},
	}

	init := initializer.NewInitializer(cfg)
	assert.NotNil(t, init)

	// 验证字段初始化
	// 注意：由于字段是私有的，我们只能通过公共方法来验证
	assert.NotNil(t, init)
}

// TestInitializer_ConfigValidation 测试配置验证
func TestInitializer_ConfigValidation(t *testing.T) {
	tests := []struct {
		name string
		config *config.Config
		expectPanic bool
	}{
		{
			name: "valid database config",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Host: "localhost",
					Port: 3306,
					User: "test",
					Password: "test",
					Name: "test",
					Conn: "mysql",
				},
			},
			expectPanic: false,
		},
		{
			name: "empty database config",
			config: &config.Config{
				Database: config.DatabaseConfig{},
			},
			expectPanic: false, // 初始化器应该能处理空配置
		},
		{
			name: "nil config",
			config: nil,
			expectPanic: false, // 初始化器应该能处理 nil 配置
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					initializer.NewInitializer(tt.config)
				})
			} else {
				assert.NotPanics(t, func() {
					init := initializer.NewInitializer(tt.config)
					assert.NotNil(t, init)
				})
			}
		})
	}
}

// TestInitializer_ContainerOperations 测试容器操作
func TestInitializer_ContainerOperations(t *testing.T) {
	config := &config.Config{
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "test",
			Password: "test",
			Name: "test",
			Conn: "mysql",
		},
	}

	init := initializer.NewInitializer(config)
	assert.NotNil(t, init)

	// 初始时容器应该为 nil
	container := init.GetContainer()
	assert.Nil(t, container)

	// 调用 InitializeCore 会因为数据库连接失败而报错
	// 由于数据库初始化失败，容器初始化也不会执行
	err := init.InitializeCore()
	assert.NotNil(t, err) // 预期数据库连接失败

	// 由于数据库连接失败，容器不会被初始化
	container = init.GetContainer()
	assert.Nil(t, container)
}

// TestInitializer_MultipleInstances 测试多个实例
func TestInitializer_MultipleInstances(t *testing.T) {
	config1 := &config.Config{
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "test1",
			Password: "test1",
			Name: "test1",
			Conn: "mysql",
		},
	}

	config2 := &config.Config{
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "test2",
			Password: "test2",
			Name: "test2",
			Conn: "mysql",
		},
	}

	init1 := initializer.NewInitializer(config1)
	init2 := initializer.NewInitializer(config2)

	assert.NotNil(t, init1)
	assert.NotNil(t, init2)
	assert.NotEqual(t, init1, init2)

	// 初始时两个实例的容器都应该为 nil
	assert.Nil(t, init1.GetContainer())
	assert.Nil(t, init2.GetContainer())

	// 验证它们是不同的实例
	assert.NotSame(t, init1, init2)
}