package like_test

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	ftypes "github.com/biuren/hebe/tests/framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLikeService_ToggleLike 测试点赞切换功能
func TestLikeService_ToggleLike(t *testing.T) {
	t.Run("成功点赞日记本", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer func() {
			if err := ctx.Cleanup(); err != nil {
				t.Log(err)
			}
		}()

		// 获取配置和服务
		config := ctx.GetConfig().GetAppConfig()
		likeRepo := repository.NewLikeRepository()
		likeService := like.NewLikeService(likeRepo)

		// 创建测试用户
		testUser := ctx.GetDataGenerator().CreateUser()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建测试日记本
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserHandle: testUser.Handle,
			Status:        1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		require.True(t, diaryID > 0)

		// 测试点赞日记本
		err = likeService.ToggleLike(testUser.UserId, diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)

		// 验证点赞状态
		status, err := likeService.GetLikeStatus(testUser.UserId, diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.True(t, status.IsLiked)

		// 验证点赞数量
		count, err := likeService.GetLikeCount(diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("成功取消点赞日记本", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer func() {
			if err := ctx.Cleanup(); err != nil {
				t.Log(err)
			}
		}()

		// 获取配置和服务
		config := ctx.GetConfig().GetAppConfig()
		likeRepo := repository.NewLikeRepository()
		likeService := like.NewLikeService(likeRepo)

		// 创建测试用户
		testUser := ctx.GetDataGenerator().CreateUser()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建测试日记本
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserHandle: testUser.Handle,
			Status:        1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		require.True(t, diaryID > 0)

		// 先点赞
		err = likeService.ToggleLike(testUser.UserId, diaryID, like.LikeTypeDiary)
		require.NoError(t, err)

		// 再取消点赞
		err = likeService.ToggleLike(testUser.UserId, diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)

		// 验证点赞状态
		status, err := likeService.GetLikeStatus(testUser.UserId, diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.False(t, status.IsLiked)

		// 验证点赞数量
		count, err := likeService.GetLikeCount(diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// TestLikeService_GetLikeStatus 测试获取点赞状态功能
func TestLikeService_GetLikeStatus(t *testing.T) {
	t.Run("获取未点赞状态", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer func() {
			if err := ctx.Cleanup(); err != nil {
				t.Log(err)
			}
		}()

		// 获取配置和服务
		config := ctx.GetConfig().GetAppConfig()
		likeRepo := repository.NewLikeRepository()
		likeService := like.NewLikeService(likeRepo)

		// 创建测试用户
		testUser := ctx.GetDataGenerator().CreateUser()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建测试日记本
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserHandle: testUser.Handle,
			Status:        1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		require.True(t, diaryID > 0)

		// 获取点赞状态（未点赞）
		status, err := likeService.GetLikeStatus(testUser.UserId, diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.False(t, status.IsLiked)
		assert.Equal(t, 0, status.LikeCount)
	})
}

// TestLikeService_GetLikeCount 测试获取点赞数量功能
func TestLikeService_GetLikeCount(t *testing.T) {
	t.Run("获取零点赞数量", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer func() {
			if err := ctx.Cleanup(); err != nil {
				t.Log(err)
			}
		}()

		// 获取配置和服务
		config := ctx.GetConfig().GetAppConfig()
		likeRepo := repository.NewLikeRepository()
		likeService := like.NewLikeService(likeRepo)

		// 创建测试用户
		testUser := ctx.GetDataGenerator().CreateUser()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建测试日记本
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserHandle: testUser.Handle,
			Status:        1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		require.True(t, diaryID > 0)

		// 获取点赞数量（应该为0）
		count, err := likeService.GetLikeCount(diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("多用户点赞统计", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer func() {
			if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
				t.Logf("Cleanup failed: %v", cleanupErr)
			}
		}()

		// 获取配置和服务
		config := ctx.GetConfig().GetAppConfig()
		likeRepo := repository.NewLikeRepository()
		likeService := like.NewLikeService(likeRepo)

		// 创建测试用户1
		testUser1 := ctx.GetDataGenerator().CreateUser()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		err := userService.RegisterUser(testUser1)
		require.NoError(t, err)

		// 创建测试用户2（确保用户名和邮箱不重复）
		testUser2 := ctx.GetDataGenerator().CreateUser(
			func(opts *ftypes.UserOptions) {
				opts.Handle = ctx.GetIsolator().IsolateHandle("u2")
				opts.Email = ctx.GetIsolator().IsolateEmail("u2@test.com")
			},
		)
		err = userService.RegisterUser(testUser2)
		require.NoError(t, err)

		// 创建测试日记本
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser1.UserId,
			AddUserHandle: testUser1.Handle,
			Status:        1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		require.True(t, diaryID > 0)

		// 用户1点赞
		err = likeService.ToggleLike(testUser1.UserId, diaryID, like.LikeTypeDiary)
		require.NoError(t, err)

		// 用户2点赞
		err = likeService.ToggleLike(testUser2.UserId, diaryID, like.LikeTypeDiary)
		require.NoError(t, err)

		// 验证点赞数量
		count, err := likeService.GetLikeCount(diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)

		// 用户1取消点赞
		err = likeService.ToggleLike(testUser1.UserId, diaryID, like.LikeTypeDiary)
		require.NoError(t, err)

		// 验证点赞数量减少
		count, err = likeService.GetLikeCount(diaryID, like.LikeTypeDiary)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
