package core

import (
	"testing"

	"github.com/biuren/hebe/internal/core/bootstrap"
	testhelpers "github.com/biuren/hebe/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitializer 测试初始化器
func TestInitializer(t *testing.T) {
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	t.Run("创建初始化器实例", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		assert.NotNil(t, initializer)
	})

	t.Run("初始化核心服务", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()

		// 初始化核心服务
		err := initializer.InitializeCore()
		assert.NoError(t, err)

		// 验证DI容器不为空
		container := initializer.GetContainer()
		assert.NotNil(t, container)

		// 验证配置不为空
		config := initializer.GetConfig()
		assert.NotNil(t, config)
	})
}

// TestBootstrapProcess 测试完整的启动过程
func TestBootstrapProcess(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		t.Run("完整启动过程", func(t *testing.T) {
			initializer := bootstrap.NewInitializer()
			require.NotNil(t, initializer)

			// 测试完整的初始化过程
			err := initializer.InitializeCore()
			assert.NoError(t, err)

			// 验证各种组件都已正确初始化
			container := initializer.GetContainer()
			assert.NotNil(t, container)

			config := initializer.GetConfig()
			assert.NotNil(t, config)
		})

		t.Run("验证服务可用性", func(t *testing.T) {
			// 验证用户服务可用
			supportedLang := suite.UserService.GetSupportedLang("zh-cn")
			assert.Equal(t, "zh-cn", supportedLang)

			// 验证各种服务的基本功能
			assert.NotNil(t, suite.UserService)
			assert.NotNil(t, suite.DiaryService)
			assert.NotNil(t, suite.PostService)
			assert.NotNil(t, suite.TaskService)
		})
	})
}

// TestBootstrapWithDatabase 测试带数据库的启动过程
func TestBootstrapWithDatabase(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		configHelper := testhelpers.SetupTestConfig(t)
		defer configHelper.Cleanup(t)

		t.Run("数据库环境下的启动", func(t *testing.T) {
			initializer := bootstrap.NewInitializer()
			require.NotNil(t, initializer)

			// 初始化核心组件
			err := initializer.InitializeCore()
			assert.NoError(t, err)

			// 验证初始化成功
			container := initializer.GetContainer()
			assert.NotNil(t, container)

			config := initializer.GetConfig()
			assert.NotNil(t, config)
		})
	})
}

// TestInitializerComponents 测试初始化器组件
func TestInitializerComponents(t *testing.T) {
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	t.Run("获取配置", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		// 初始化后才能获取配置
		err := initializer.InitializeCore()
		assert.NoError(t, err)

		config := initializer.GetConfig()
		assert.NotNil(t, config)
		assert.Equal(t, "test", config.Env)
	})

	t.Run("获取DI容器", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		// 初始化后才能获取容器
		err := initializer.InitializeCore()
		assert.NoError(t, err)

		container := initializer.GetContainer()
		assert.NotNil(t, container)
	})
}

// TestInitializationOrder 测试初始化顺序
func TestInitializationOrder(t *testing.T) {
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	t.Run("初始化前获取组件", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		// 初始化前获取配置应该为空
		config := initializer.GetConfig()
		assert.Nil(t, config)

		// 初始化前获取容器应该为空
		container := initializer.GetContainer()
		assert.Nil(t, container)
	})

	t.Run("初始化后获取组件", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		// 执行初始化
		err := initializer.InitializeCore()
		assert.NoError(t, err)

		// 初始化后获取配置应该不为空
		config := initializer.GetConfig()
		assert.NotNil(t, config)

		// 初始化后获取容器应该不为空
		container := initializer.GetContainer()
		assert.NotNil(t, container)
	})
}

// TestConfigurationValidation 测试配置验证
func TestConfigurationValidation(t *testing.T) {
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	t.Run("验证必需的配置项", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		err := initializer.InitializeCore()
		assert.NoError(t, err)

		cfg := initializer.GetConfig()
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
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		t.Run("验证服务依赖解析", func(t *testing.T) {
			initializer := bootstrap.NewInitializer()
			require.NotNil(t, initializer)

			err := initializer.InitializeCore()
			assert.NoError(t, err)

			// 验证从初始化器可以获取到容器
			container := initializer.GetContainer()
			assert.NotNil(t, container)
		})

		t.Run("验证服务间依赖", func(t *testing.T) {
			// 验证服务之间的依赖关系正确
			assert.NotNil(t, suite.UserService)
			assert.NotNil(t, suite.DiaryService)
			assert.NotNil(t, suite.PostService)
			assert.NotNil(t, suite.TaskService)
		})
	})
}

// TestEnvironmentSpecificBootstrap 测试环境特定的启动
func TestEnvironmentSpecificBootstrap(t *testing.T) {
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	t.Run("测试环境启动", func(t *testing.T) {
		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		// 初始化应该适应测试环境
		err := initializer.InitializeCore()
		assert.NoError(t, err)

		config := initializer.GetConfig()
		require.NotNil(t, config)

		// 验证是测试环境
		assert.Equal(t, "test", config.Env)
	})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	t.Run("多次初始化", func(t *testing.T) {
		// 设置临时的测试配置
		configHelper := testhelpers.SetupTestConfig(t)
		defer configHelper.Cleanup(t)

		initializer := bootstrap.NewInitializer()
		require.NotNil(t, initializer)

		// 第一次初始化
		err := initializer.InitializeCore()
		assert.NoError(t, err)

		// 第二次初始化应该也能成功（或者优雅处理）
		err = initializer.InitializeCore()
		assert.NoError(t, err)
	})
}

// TestPerformance 测试启动性能
func TestPerformance(t *testing.T) {
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	t.Run("启动性能测试", func(t *testing.T) {
		// 测试多次初始化的性能
		for i := 0; i < 5; i++ {
			initializer := bootstrap.NewInitializer()
			require.NotNil(t, initializer)

			// 快速初始化测试
			err := initializer.InitializeCore()
			assert.NoError(t, err)
		}
	})
}

// TestComponentIntegration 测试组件集成
func TestComponentIntegration(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		t.Run("bootstrap与测试套件集成", func(t *testing.T) {
			initializer := bootstrap.NewInitializer()
			require.NotNil(t, initializer)

			err := initializer.InitializeCore()
			assert.NoError(t, err)

			// 验证bootstrap初始化的组件与测试套件兼容
			config := initializer.GetConfig()
			assert.NotNil(t, config)

			suiteConfig := suite.Config.GetTestConfig()
			assert.NotNil(t, suiteConfig)

			// 环境应该一致
			assert.Equal(t, config.Env, suiteConfig.Env)
		})
	})
}
