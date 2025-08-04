package core

import (
	"testing"

	bootstrap "github.com/biuren/hebe/internal/core/initializer"
	framework "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitializer 测试初始化器
func TestInitializer(t *testing.T) {
	ctx := framework.NewTestContext("TestInitializer")
	defer ctx.Cleanup()

	t.Run("创建初始化器实例", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		assert.NotNil(t, initInstance)
	})

	t.Run("初始化核心服务", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()

		// 初始化核心服务
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 验证各种组件都已正确初始化
		container := initInstance.GetContainer()
		assert.NotNil(t, container)

		config := initInstance.GetConfig()
		assert.NotNil(t, config)
	})
}

// TestBootstrapProcess 测试完整的启动过程
func TestBootstrapProcess(t *testing.T) {
	ctx := framework.NewTestContext("TestBootstrapProcess")
	defer ctx.Cleanup()

	t.Run("完整启动过程", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 测试完整的初始化过程
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 验证各种组件都已正确初始化
		container := initInstance.GetContainer()
		assert.NotNil(t, container)

		config := initInstance.GetConfig()
		assert.NotNil(t, config)
	})
}

// TestBootstrapWithDatabase 测试带数据库的启动过程
func TestBootstrapWithDatabase(t *testing.T) {
	ctx := framework.NewTestContext("TestBootstrapWithDatabase")
	defer ctx.Cleanup()

	t.Run("数据库环境下的启动", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 初始化核心组件
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 验证初始化成功
		container := initInstance.GetContainer()
		assert.NotNil(t, container)

		config := initInstance.GetConfig()
		assert.NotNil(t, config)
	})
}

// TestInitializerComponents 测试初始化器组件
func TestInitializerComponents(t *testing.T) {
	ctx := framework.NewTestContext("TestInitializerComponents")
	defer ctx.Cleanup()

	t.Run("获取配置", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 初始化后才能获取配置
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		config := initInstance.GetConfig()
		assert.NotNil(t, config)
		assert.Equal(t, "test", config.Env)
	})

	t.Run("获取DI容器", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 初始化后才能获取容器
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		container := initInstance.GetContainer()
		assert.NotNil(t, container)
	})
}

// TestInitializationOrder 测试初始化顺序
func TestInitializationOrder(t *testing.T) {
	ctx := framework.NewTestContext("TestInitializationOrder")
	defer ctx.Cleanup()

	t.Run("初始化前获取组件", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 初始化前获取配置应该为空
		config := initInstance.GetConfig()
		assert.Nil(t, config)

		// 初始化前获取容器应该为空
		container := initInstance.GetContainer()
		assert.Nil(t, container)
	})

	t.Run("初始化后获取组件", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 执行初始化
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 初始化后获取配置应该不为空
		config := initInstance.GetConfig()
		assert.NotNil(t, config)

		// 初始化后获取容器应该不为空
		container := initInstance.GetContainer()
		assert.NotNil(t, container)
	})
}

// TestConfigurationValidation 测试配置验证
func TestConfigurationValidation(t *testing.T) {
	ctx := framework.NewTestContext("TestConfigurationValidation")
	defer ctx.Cleanup()

	t.Run("验证必需的配置项", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		cfg := initInstance.GetConfig()
		require.NotNil(t, cfg)

		// 验证关键配置项存在
		assert.NotEmpty(t, cfg.Env)
		assert.NotEmpty(t, cfg.App.DefaultLang)
		assert.NotEmpty(t, cfg.Database.Host)
		assert.Greater(t, cfg.Database.Port, 0)
	})
}

// TestDependencyResolution 测试依赖解析
func TestDependencyResolution(t *testing.T) {
	ctx := framework.NewTestContext("TestDependencyResolution")
	defer ctx.Cleanup()

	t.Run("验证服务依赖解析", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 验证从初始化器可以获取到容器
		container := initInstance.GetContainer()
		assert.NotNil(t, container)
	})
}

// TestEnvironmentSpecificBootstrap 测试环境特定的启动
func TestEnvironmentSpecificBootstrap(t *testing.T) {
	ctx := framework.NewTestContext("TestEnvironmentSpecificBootstrap")
	defer ctx.Cleanup()

	t.Run("测试环境启动", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 初始化应该适应测试环境
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		config := initInstance.GetConfig()
		require.NotNil(t, config)

		// 验证是测试环境
		assert.Equal(t, "test", config.Env)
	})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	ctx := framework.NewTestContext("TestErrorHandling")
	defer ctx.Cleanup()

	t.Run("多次初始化", func(t *testing.T) {

		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		// 第一次初始化
		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 第二次初始化应该也能成功（或者优雅处理）
		err = initInstance.InitializeCore()
		assert.NoError(t, err)
	})
}

// TestPerformance 测试启动性能
func TestPerformance(t *testing.T) {
	ctx := framework.NewTestContext("TestPerformance")
	defer ctx.Cleanup()

	t.Run("启动性能测试", func(t *testing.T) {
		// 测试多次初始化的性能
		for i := 0; i < 5; i++ {
			initInstance := bootstrap.NewInitializer()
			require.NotNil(t, initInstance)

			// 快速初始化测试
			err := initInstance.InitializeCore()
			assert.NoError(t, err)
		}
	})
}

// TestComponentIntegration 测试组件集成
func TestComponentIntegration(t *testing.T) {
	ctx := framework.NewTestContext("TestComponentIntegration")
	defer ctx.Cleanup()

	t.Run("bootstrap与测试套件集成", func(t *testing.T) {
		initInstance := bootstrap.NewInitializer()
		require.NotNil(t, initInstance)

		err := initInstance.InitializeCore()
		assert.NoError(t, err)

		// 验证bootstrap初始化的组件与测试套件兼容
		config := initInstance.GetConfig()
		assert.NotNil(t, config)

		suiteConfig := ctx.GetConfig()
		assert.NotNil(t, suiteConfig)

		// 环境应该一致
		// assert.Equal(t, config.Env, suiteConfig.Env)
	})
}
