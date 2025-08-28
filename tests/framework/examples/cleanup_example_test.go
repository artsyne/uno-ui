package examples

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/tests/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseCleanupExample 演示数据库清理功能的使用
// 支持命令行参数控制:
//   go test ./tests/framework/examples/ -v -cleardb  # 强制清理数据库
//   go test ./tests/framework/examples/ -v           # 使用默认配置
// 包括自动清理、强制清理、手动控制清理时机等
func TestDatabaseCleanupExample(t *testing.T) {
	t.Parallel()

	t.Run("正常清理流程", func(t *testing.T) {
		// 创建测试上下文
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup() // 正常清理

		// 创建服务
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		diaryRepo := repository.NewDiaryRepository()
		postRepo := repository.NewPostRepository(ctx.GetDatabase().GetDB())
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager)
		diaryService := diary.NewService(config, diaryRepo)
		postService := post.NewService(config, userRepo, taskService, diaryRepo, postRepo, i18nManager)

		// 创建一些测试数据
		generator := ctx.GetDataGenerator()
		isolator := ctx.GetIsolator()

		// 创建用户并保存到数据库
		testUser := generator.CreateUser()
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)
		isolator.TrackCreated("br_user", testUser.UserId)

		// 创建日记并保存到数据库
		testDiary := generator.CreateDiary(testUser.UserId)
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		testDiary.DiaryID = diaryID
		isolator.TrackCreated("br_diary", diaryID)

		// 创建文章并保存到数据库
		testPost := generator.CreatePost(testUser.UserId, testDiary.DiaryID)
		postID, err := postService.AddPost(testPost)
		require.NoError(t, err)
		testPost.PostID = postID
		isolator.TrackCreated("br_post", postID)

		// 验证数据已创建
		assert.True(t, isolator.GetDataCount() > 0)
		assert.True(t, isolator.GetTableCount() > 0)

		t.Logf("创建了 %d 条数据，涉及 %d 个表", isolator.GetDataCount(), isolator.GetTableCount())
	})

	t.Run("强制清理数据库", func(t *testing.T) {
		// 创建测试上下文，禁用自动清理
		ctx := framework.NewTestContext(framework.WithDatabase(), framework.WithCleanupEnabled(false))

		// 创建服务
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		diaryRepo := repository.NewDiaryRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager)
		diaryService := diary.NewService(config, diaryRepo)

		// 创建一些测试数据
		generator := ctx.GetDataGenerator()
		isolator := ctx.GetIsolator()

		// 创建用户并保存到数据库
		testUser := generator.CreateUser()
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)
		isolator.TrackCreated("br_user", testUser.UserId)

		// 创建日记并保存到数据库
		testDiary := generator.CreateDiary(testUser.UserId)
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		isolator.TrackCreated("br_diary", diaryID)

		// 验证数据已创建
		assert.True(t, isolator.GetDataCount() > 0)
		t.Logf("创建了 %d 条数据，涉及 %d 个表", isolator.GetDataCount(), isolator.GetTableCount())

		// 强制清理数据库（无论配置如何）
		err = ctx.ForceCleanupDatabase()
		assert.NoError(t, err, "强制清理数据库应该成功")

		t.Log("强制清理数据库完成")

		// 注意：由于禁用了自动清理，defer ctx.Cleanup() 不会清理数据
		// 但我们已经通过 ForceCleanupDatabase() 手动清理了
	})

	t.Run("手动控制清理时机", func(t *testing.T) {
		// 创建测试上下文
		ctx := framework.NewTestContext(framework.WithDatabase())

		// 创建服务
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager)

		// 创建第一批测试数据
		generator := ctx.GetDataGenerator()
		isolator := ctx.GetIsolator()

		// 创建第一个用户并保存到数据库
		user1 := generator.CreateUser()
		err := userService.RegisterUser(user1)
		require.NoError(t, err)
		isolator.TrackCreated("br_user", user1.UserId)

		// 在测试中间清理数据
		err = ctx.ForceCleanupDatabase()
		assert.NoError(t, err)
		t.Log("中间清理完成")

		// 创建第二批测试数据（使用新的测试上下文以获得新的前缀）
		ctx2 := framework.NewTestContext(framework.WithDatabase())
		defer ctx2.Cleanup()

		generator2 := ctx2.GetDataGenerator()
		isolator2 := ctx2.GetIsolator()

		// 创建第二个用户并保存到数据库
		user2 := generator2.CreateUser()
		err = userService.RegisterUser(user2)
		require.NoError(t, err)
		isolator2.TrackCreated("br_user", user2.UserId)

		// 验证新数据的用户名不同（因为前缀会更新）
		assert.NotEqual(t, user1.UserName, user2.UserName)

		// 最终清理
		defer ctx.Cleanup()
	})

	t.Run("清理管理器直接使用", func(t *testing.T) {
			// 创建测试上下文
			ctx := framework.NewTestContext(framework.WithDatabase())
			defer ctx.Cleanup()

		// 创建服务
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager)

		// 创建测试数据
		generator := ctx.GetDataGenerator()
		isolator := ctx.GetIsolator()

		// 创建用户并保存到数据库
		testUser := generator.CreateUser()
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)
		isolator.TrackCreated("br_user", testUser.UserId)

		// 直接使用清理管理器
		cleanupManager := ctx.GetCleanup()
		assert.NotNil(t, cleanupManager)

		// 强制清理
		err = cleanupManager.ForceCleanupDatabase()
		assert.NoError(t, err)
		t.Log("通过清理管理器强制清理完成")
	})
}

// TestCommandLineFlagsExample 演示命令行参数控制清理功能
// 运行方式:
//   go test ./tests/framework/examples/ -v -cleardb -run TestCommandLineFlagsExample
func TestCommandLineFlagsExample(t *testing.T) {
	t.Parallel()

	// 检查命令行参数状态
	flags := framework.GetTestFlags()
	t.Logf("命令行参数状态: cleardb=%v, verbose=%v, parallel=%v", 
		flags.ClearDB, flags.Verbose, flags.Parallel)

	// 检查清理配置
	shouldCleanup := framework.ShouldCleanupDatabase()
	t.Logf("是否清理数据库: %v", shouldCleanup)

	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建测试数据
	generator := ctx.GetDataGenerator()
	user := generator.CreateUser()
	assert.NotNil(t, user)
	t.Logf("创建用户: %s (ID: %d)", user.UserName, user.UserId)

	diary := generator.CreateDiary(user.UserId)
	assert.NotNil(t, diary)
	t.Logf("创建日记: %s (ID: %d)", diary.Title, diary.DiaryID)

	// 如果指定了 -cleardb 参数，演示强制清理
	if flags.ClearDB {
		t.Log("检测到 -cleardb 参数，执行强制清理")
		err := ctx.ForceCleanupDatabase()
		assert.NoError(t, err)
		t.Log("强制清理完成")
	} else {
		t.Log("未指定 -cleardb 参数，使用默认清理策略")
	}
}

// TestCleanupStrategies 演示不同的清理策略
func TestCleanupStrategies(t *testing.T) {
	t.Parallel()

	t.Run("按前缀清理", func(t *testing.T) {
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		// 创建测试数据
		generator := ctx.GetDataGenerator()
		user := generator.CreateUser()
		assert.NotNil(t, user)

		// 使用清理管理器按前缀清理
		cleanupManager := ctx.GetCleanup()
		err := cleanupManager.CleanupByPrefix()
		assert.NoError(t, err)
	})

	t.Run("按跟踪数据清理", func(t *testing.T) {
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		// 创建测试数据
		generator := ctx.GetDataGenerator()
		user := generator.CreateUser()
		assert.NotNil(t, user)

		// 使用清理管理器按跟踪数据清理
		cleanupManager := ctx.GetCleanup()
		err := cleanupManager.CleanupTrackedData()
		assert.NoError(t, err)
	})
}