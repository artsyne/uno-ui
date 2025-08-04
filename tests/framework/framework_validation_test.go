package framework

import (
	"testing"
	"time"

	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
)

// TestFrameworkBasicFunctionality 验证新框架基本功能
func TestFrameworkBasicFunctionality(t *testing.T) {
	t.Parallel()

	t.Run("TestContext创建和初始化", func(t *testing.T) {
		// 创建测试上下文
		ctx := core.NewTestContext("framework_test")
		assert.NotNil(t, ctx, "TestContext不应该为nil")

		// 验证测试ID
		assert.Equal(t, "framework_test", ctx.GetTestID())

		// 清理
		defer func() {
			if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
				t.Logf("清理失败: %v", cleanupErr)
			}
		}()
	})

	t.Run("数据隔离器功能测试", func(t *testing.T) {
		ctx := core.NewTestContext("isolation_test")
		defer ctx.Cleanup()

		isolator := ctx.GetIsolator()
		assert.NotNil(t, isolator)

		// 测试字符串隔离
		isolated := isolator.IsolateString("test")
		assert.Contains(t, isolated, "test")
		assert.NotEqual(t, "test", isolated)

		// 测试邮箱隔离
		email := isolator.IsolateEmail("user@test.com")
		assert.Contains(t, email, "@test.com")
		assert.NotEqual(t, "user@test.com", email)

		// 测试用户名隔离
		username := isolator.IsolateUsername("testuser")
		assert.Contains(t, username, "testuser")
		assert.True(t, len(username) <= 90) // 验证长度限制

		// 测试标题隔离
		title := isolator.IsolateTitle("test title")
		assert.Contains(t, title, "test title")
		assert.True(t, len(title) <= 180) // 验证长度限制
	})

	t.Run("数据生成器功能测试", func(t *testing.T) {
		ctx := core.NewTestContext("generator_test")
		defer ctx.Cleanup()

		generator := ctx.GetDataGenerator()
		assert.NotNil(t, generator)

		// 测试用户生成
		user := generator.CreateUser()
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.UserName)
		assert.NotEmpty(t, user.Email)
		assert.Contains(t, user.Email, "@")
		assert.Equal(t, int8(3), user.Status) // 正常状态

		// 测试日记生成
		diary := generator.CreateDiary(1)
		assert.NotNil(t, diary)
		assert.NotEmpty(t, diary.Title)
		assert.Equal(t, int64(1), diary.AddUserID)
		assert.Equal(t, int8(1), diary.Status)

		// 测试文章生成
		post := generator.CreatePost(1, 1)
		assert.NotNil(t, post)
		assert.NotEmpty(t, post.Title)
		assert.NotEmpty(t, post.Body)
		assert.Equal(t, int64(1), post.DiaryID)
		assert.Equal(t, int64(1), post.AddUserID)
		assert.Equal(t, int8(1), post.Status)
	})

	t.Run("配置管理器功能测试", func(t *testing.T) {
		ctx := core.NewTestContext("config_test")
		defer ctx.Cleanup()

		config := ctx.GetConfig()
		assert.NotNil(t, config)

		// 验证数据库配置
		dbConfig := config.GetDatabaseConfig()
		assert.NotEmpty(t, dbConfig.Host)
		assert.NotEmpty(t, dbConfig.User)
		assert.NotEmpty(t, dbConfig.Password)
		assert.NotEmpty(t, dbConfig.Name)

		// 验证隔离配置
		isolationConfig := config.GetIsolationConfig()
		assert.NotEmpty(t, isolationConfig.PrefixFormat)
		assert.True(t, isolationConfig.MaxPrefixLen > 0)

		// 验证清理配置
		cleanupConfig := config.GetCleanupConfig()
		assert.NotEmpty(t, cleanupConfig.Strategy)

		// 验证性能配置
		perfConfig := config.GetPerformanceConfig()
		assert.True(t, perfConfig.Parallel)
	})

	t.Run("数据库管理器功能测试", func(t *testing.T) {
		ctx := core.NewTestContext("database_test")
		defer ctx.Cleanup()

		dbManager := ctx.GetDatabase()
		assert.NotNil(t, dbManager)

		// 验证数据库连接
		db := dbManager.GetDB()
		assert.NotNil(t, db)

		rawDB := dbManager.GetRawDB()
		assert.NotNil(t, rawDB)

		// 验证健康检查
		assert.True(t, dbManager.IsHealthy())

		// 简单的数据库查询测试
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?",
			ctx.GetConfig().GetDatabaseConfig().Name).Scan(&count)
		assert.NoError(t, err)
		assert.True(t, count > 0, "测试数据库应该包含表")
	})

	t.Run("并发安全性测试", func(t *testing.T) {
		const goroutines = 10
		results := make(chan error, goroutines)

		// 并发创建多个TestContext
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				ctx := core.NewTestContext("concurrent_test")

				// 执行一些操作
			isolator := ctx.GetIsolator()
				generator := ctx.GetDataGenerator()

				user := generator.CreateUser()
				isolated := isolator.IsolateString("test")

				// 验证结果
				if user == nil || isolated == "" {
					results <- assert.AnError
					return
				}

				// 清理
				ctx.Cleanup()
				results <- nil
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < goroutines; i++ {
			select {
			case err := <-results:
				assert.NoError(t, err, "并发测试应该成功")
			case <-time.After(30 * time.Second):
				t.Fatal("并发测试超时")
			}
		}
	})

	t.Run("性能基准测试", func(t *testing.T) {
		if testing.Short() {
			t.Skip("跳过性能测试")
		}

		start := time.Now()

		// 创建和销毁100个测试上下文
		for i := 0; i < 100; i++ {
			ctx := core.NewTestContext("perf_test")

			// 执行一些基本操作
			generator := ctx.GetDataGenerator()
			user := generator.CreateUser()
			assert.NotNil(t, user)

			ctx.Cleanup()
		}

		duration := time.Since(start)
		t.Logf("创建100个TestContext耗时: %v (平均: %v/个)", duration, duration/100)

		// 验证性能要求 (每个TestContext创建应该在100ms内)
		assert.True(t, duration < 10*time.Second, "性能要求: 100个TestContext应该在10秒内完成")
	})
}

// TestFrameworkErrorHandling 验证框架错误处理
func TestFrameworkErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("无效测试ID处理", func(t *testing.T) {
		ctx := core.NewTestContext("")
		// 如果允许空ID，应该生成默认ID
		assert.NotNil(t, ctx)
		assert.NotEmpty(t, ctx.GetTestID())
		ctx.Cleanup()
	})

	t.Run("重复清理处理", func(t *testing.T) {
		ctx := core.NewTestContext("cleanup_test")
		assert.NotNil(t, ctx)

		// 第一次清理
		err1 := ctx.Cleanup()
		assert.NoError(t, err1)

		// 第二次清理应该也不报错
		err2 := ctx.Cleanup()
		assert.NoError(t, err2)
	})
}
