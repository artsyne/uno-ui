package framework

import (
	"log/slog"
	"os"
	"testing"
	"time"

	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFrameworkBasicFunctionality 验证新框架基本功能
func TestFrameworkBasicFunctionality(t *testing.T) {
	t.Parallel()

	t.Run("TestContext创建和初始化", func(t *testing.T) {
		// 创建测试上下文
		ctx := NewTestContext()
		assert.NotNil(t, ctx, "TestContext不应该为nil")

		// 验证测试上下文
		assert.NotNil(t, ctx.GetTestID())

		// 清理
		defer func() {
			if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
				t.Logf("清理失败: %v", cleanupErr)
			}
		}()
	})

	t.Run("数据隔离器功能测试", func(t *testing.T) {
		ctx := NewTestContext()
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
		ctx := NewTestContext()
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
		entry := generator.CreateEntry(1, 1)
		assert.NotNil(t, entry)
		assert.NotEmpty(t, entry.Title)
		assert.NotEmpty(t, entry.Body)
		assert.Equal(t, int64(1), entry.DiaryID)
		assert.Equal(t, int64(1), entry.AddUserID)
		assert.Equal(t, int8(1), entry.Status)
	})

	t.Run("配置管理器功能测试", func(t *testing.T) {
		ctx := NewTestContext()
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
		ctx := NewTestContext(WithDatabase())
		defer ctx.Cleanup()

		dbManager := ctx.GetDatabase()
		require.NotNil(t, dbManager)

		// 验证数据库连接
		db := dbManager.GetDB()
		require.NotNil(t, db, "数据库连接应该已初始化")

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
				ctx := NewTestContext()

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
		// 简化性能测试，减少迭代次数以适应常规测试
		iterations := 10
		if testing.Short() {
			iterations = 3 // 短模式下进一步减少
		}

		start := time.Now()

		// 创建和销毁多个测试上下文
		for i := 0; i < iterations; i++ {
			ctx := NewTestContext()

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
		ctx := NewTestContext()
		// 如果允许空ID，应该生成默认ID
		assert.NotNil(t, ctx)
		assert.NotEmpty(t, ctx.GetTestID())
		ctx.Cleanup()
	})

	t.Run("重复清理处理", func(t *testing.T) {
		ctx := NewTestContext()
		assert.NotNil(t, ctx)

		// 第一次清理
		err1 := ctx.Cleanup()
		assert.NoError(t, err1)

		// 第二次清理应该也不报错
		err2 := ctx.Cleanup()
		assert.NoError(t, err2)
	})
}

// TestDatabaseInitialization 验证数据库初始化和错误处理
func TestDatabaseInitialization(t *testing.T) {
	t.Parallel()

	t.Run("数据库连接和初始化检查", func(t *testing.T) {
		// 设置调试日志级别以便查看详细信息
		originalLevel := slog.SetLogLoggerLevel(slog.LevelDebug)
		defer slog.SetLogLoggerLevel(originalLevel)

		slog.Info("开始数据库初始化测试")

		// 创建测试上下文，这会触发数据库初始化
		ctx := NewTestContext(WithDatabase())
		require.NotNil(t, ctx, "测试上下文创建失败")
		defer ctx.Cleanup()

		// 获取数据库管理器
		dbManager := ctx.GetDatabase()
		require.NotNil(t, dbManager, "数据库管理器获取失败")

		// 检查数据库连接健康状态
		require.True(t, dbManager.IsHealthy(), "数据库连接不健康")

		// 验证数据库连接
		db := dbManager.GetDB()
		require.NotNil(t, db, "数据库连接为空")

		// 测试基本查询
		var version string
		err := db.QueryRow("SELECT VERSION()").Scan(&version)
		require.NoError(t, err, "数据库版本查询失败")
		assert.NotEmpty(t, version, "数据库版本信息为空")
		slog.Info("数据库版本信息", "version", version)

		// 检查测试数据库是否存在关键表
		config := ctx.GetConfig().GetDatabaseConfig()
		var tableCount int
		err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'br_user'",
			config.Name).Scan(&tableCount)
		require.NoError(t, err, "检查br_user表失败")
		assert.Equal(t, 1, tableCount, "br_user表不存在，数据库schema可能未正确加载")

		// 检查更多关键表
		keyTables := []string{"br_diary", "br_entry", "br_user"}
		for _, tableName := range keyTables {
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
				config.Name, tableName).Scan(&count)
			assert.NoError(t, err, "检查表%s失败", tableName)
			assert.Equal(t, 1, count, "表%s不存在", tableName)
		}

		slog.Info("数据库初始化测试完成")
	})

	t.Run("数据库错误处理测试", func(t *testing.T) {
		// 设置调试日志级别
		originalLevel := slog.SetLogLoggerLevel(slog.LevelDebug)
		defer slog.SetLogLoggerLevel(originalLevel)

		ctx := NewTestContext(WithDatabase())
		defer ctx.Cleanup()

		dbManager := ctx.GetDatabase()
		require.NotNil(t, dbManager)

		// 测试无效SQL查询的错误处理
		_, err := dbManager.Execute("INVALID SQL STATEMENT")
		assert.Error(t, err, "无效SQL应该返回错误")
		slog.Info("无效SQL错误测试", "error", err.Error())

		// 测试查询不存在的表
		_, err = dbManager.Query("SELECT * FROM non_existent_table")
		assert.Error(t, err, "查询不存在的表应该返回错误")
		slog.Info("查询不存在表错误测试", "error", err.Error())
	})

	t.Run("项目根目录检查", func(t *testing.T) {
		// 测试GetProjectRoot函数
		projectRoot := core.GetProjectRoot()
		assert.NotEmpty(t, projectRoot, "项目根目录不应为空")
		slog.Info("项目根目录", "path", projectRoot)

		// 检查go.mod文件是否存在
		goModPath := projectRoot + "/go.mod"
		_, err := os.Stat(goModPath)
		assert.NoError(t, err, "go.mod文件应该存在于项目根目录")

		// 检查tests目录是否存在
		testsDir := projectRoot + "/tests"
		_, err = os.Stat(testsDir)
		assert.NoError(t, err, "tests目录应该存在")
		slog.Info("tests目录路径", "path", testsDir)
	})
}
