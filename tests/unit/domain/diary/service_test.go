package diary_test

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelperImport 测试 framework 包导入
func TestHelperImport(t *testing.T) {
	// 简单测试 framework 包是否可用
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	if ctx == nil {
		t.Fatal("framework.NewTestContext returned nil")
	}
}

// TestDiaryService_AddDiary 测试添加日记本功能
func TestDiaryService_AddDiary(t *testing.T) {
	t.Run("成功添加日记本", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		// 获取服务
		diaryRepo := repository.NewDiaryRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		diaryService := diary.NewService(config, diaryRepo)

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
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			Status:        1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)

		// 验证日记本创建成功
		assert.True(t, diaryID > 0)
		assert.Equal(t, testUser.UserId, testDiary.AddUserID)
		assert.Equal(t, testUser.UserName, testDiary.AddUserName)
	})

	t.Run("用户ID无效", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		// 获取服务
		diaryRepo := repository.NewDiaryRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		diaryService := diary.NewService(config, diaryRepo)

		// 尝试创建无效用户ID的日记本
		invalidDiary := &diary.Diary{
			Title:       "测试日记本",
			Remark:      "测试备注",
			AddUserID:   0, // 无效用户ID
			AddUserName: "invalid_user",
			Status:      1,
		}

		_, err := diaryService.AddDiary(invalidDiary)
		assert.Error(t, err)
	})

	t.Run("检查日记本数量限制", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		// 获取服务
		diaryRepo := repository.NewDiaryRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		diaryService := diary.NewService(config, diaryRepo)

		// 创建测试用户
		testUser := ctx.GetDataGenerator().CreateUser()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 首先验证限制检查功能
		err = diaryService.CheckDiaryLimit(testUser.UserId)
		assert.NoError(t, err) // 当前没有日记本，应该不会超限

		// 创建一个日记本
		testDiary := &diary.Diary{
			Title:         "测试日记本",
			Remark:        "测试备注",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			Status:        1,
		}
		_, err = diaryService.AddDiary(testDiary)
		require.NoError(t, err)

		// 再次检查限制
		err = diaryService.CheckDiaryLimit(testUser.UserId)
		assert.NoError(t, err) // 应该还没超限
	})
}

// TestDiaryService_AddDiary_EmptyTitle 测试添加空标题日记本
func TestDiaryService_AddDiary_EmptyTitle(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 获取服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	diaryService := diary.NewService(config, diaryRepo)

	// 创建测试用户
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		Lang:     "zh-cn",
	}
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	config = ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 尝试创建空标题的日记本
	invalidDiary := &diary.Diary{
		Title:       "", // 空标题
		Remark:      "测试备注",
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}

	_, err = diaryService.AddDiary(invalidDiary)
	assert.Error(t, err)
}

// TestDiaryService_GetDiaryByID 测试根据ID获取日记本
func TestDiaryService_GetDiaryByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo)

	// 先创建一个用户
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 添加一个日记本
	newDiary := &diary.Diary{
		Title:         "测试日记本标题",
		Remark:        "用于测试获取功能",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(newDiary)
	require.NoError(t, err)

	t.Run("成功获取日记本", func(t *testing.T) {
		retrievedDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedDiary)
		assert.Equal(t, newDiary.Title, retrievedDiary.Title)
		assert.Equal(t, "用于测试获取功能", retrievedDiary.Remark)
		assert.Equal(t, testUser.UserId, retrievedDiary.AddUserID)
	})

	t.Run("日记本不存在", func(t *testing.T) {
		retrievedDiary, err := diaryService.GetDiaryByID(99999)
		assert.Error(t, err)
		assert.Nil(t, retrievedDiary)
	})
}

// TestDiaryService_UpdateDiary 测试更新日记本
func TestDiaryService_UpdateDiary(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo)

	// 先创建一个用户
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 添加一个日记本
	originalDiary := &diary.Diary{
		Title:         "原始日记本标题",
		Remark:        "原始备注",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(originalDiary)
	require.NoError(t, err)

	t.Run("成功更新日记本", func(t *testing.T) {
		updatedTitle := "测试日记本标题"
		updatedDiary := &diary.Diary{
			DiaryID:       diaryID,
			Title:         updatedTitle,
			Remark:        "更新后的备注",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
		}

		err := diaryService.UpdateDiary(updatedDiary)
		assert.NoError(t, err)

		// 验证更新是否成功
		retrievedDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		assert.Equal(t, updatedTitle, retrievedDiary.Title)
		assert.Equal(t, "更新后的备注", retrievedDiary.Remark)
	})

	t.Run("更新不存在的日记本", func(t *testing.T) {
		nonExistentDiary := &diary.Diary{
			DiaryID:       99999,
			Title:         "测试日记本标题",
			Remark:        "测试",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
		}

		err := diaryService.UpdateDiary(nonExistentDiary)
		assert.Error(t, err)
	})

	t.Run("标题为空时更新失败", func(t *testing.T) {
		invalidDiary := &diary.Diary{
			DiaryID:       diaryID,
			Title:         "", // 空标题
			Remark:        "备注",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
		}

		err := diaryService.UpdateDiary(invalidDiary)
		assert.Error(t, err)
	})
}

// TestDiaryService_ListDiariesForWrite 测试获取可写入的日记本列表
func TestDiaryService_ListDiariesForWrite(t *testing.T) {
	t.Run("获取用户1的日记本列表", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo)

		user1 := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(user1)
		require.NoError(t, err)

		diary1 := &diary.Diary{
			Title:         ctx.GetDataGenerator().GenerateUniqueString("fullname"),
			Remark:        "测试",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     user1.UserId,
			AddUserName:   user1.UserName,
			Status:        1,
		}
		_, err = diaryService.AddDiary(diary1)
		require.NoError(t, err)

		diary2 := &diary.Diary{
			Title:         ctx.GetDataGenerator().GenerateUniqueString("fullname"),
			Remark:        "测试",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     user1.UserId,
			AddUserName:   user1.UserName,
			Status:        1,
		}
		_, err = diaryService.AddDiary(diary2)
		require.NoError(t, err)

		diaries, err := diaryService.ListDiariesForWrite(user1.UserId)
		require.NoError(t, err)
		assert.Len(t, diaries, 2)

		for _, d := range diaries {
			assert.Equal(t, user1.UserId, d.AddUserID)
		}
	})

	t.Run("获取用户2的日记本列表", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		// 创建真实的服务
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		i18nManager := i18n.NewLocalizationManager(config)
		taskService := task.NewService(config, taskRepo)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
		diaryRepo := repository.NewDiaryRepository()
		diaryService := diary.NewService(config, diaryRepo)

		user2 := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(user2)
		require.NoError(t, err)

		diary3 := &diary.Diary{
			Title:         ctx.GetDataGenerator().GenerateUniqueString("fullname"),
			Remark:        "测试",
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     user2.UserId,
			AddUserName:   user2.UserName,
			Status:        1,
		}
		_, err = diaryService.AddDiary(diary3)
		require.NoError(t, err)

		diaries, err := diaryService.ListDiariesForWrite(user2.UserId)
		require.NoError(t, err)
		assert.Len(t, diaries, 1)
		assert.Equal(t, diary3.Title, diaries[0].Title)
	})

	t.Run("获取不存在用户的日记本列表", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		diaryRepo := repository.NewDiaryRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		diaryService := diary.NewService(config, diaryRepo)
		// 使用一个极大的用户ID，确保不存在
		diaries, err := diaryService.ListDiariesForWrite(999999999)
		assert.NoError(t, err)
		assert.Len(t, diaries, 0)
	})
}

// TestDiaryService_ListDiariesByUser 测试分页获取用户日记本列表
func TestDiaryService_ListDiariesByUser(t *testing.T) {
	t.Run("第一页，每页3个", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		diaryRepo := repository.NewDiaryRepository()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

		diaryService := diary.NewService(config, diaryRepo)
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)

		// 创建测试用户
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建5个日记本
		for i := 1; i <= 5; i++ {
			testDiary := &diary.Diary{
				Title:         fmt.Sprintf("diary%d", i),
				Remark:        "测试",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     testUser.UserId,
				AddUserName:   testUser.UserName,
				Status:        1,
			}
			_, err = diaryService.AddDiary(testDiary)
			require.NoError(t, err)
		}

		diaries, total, err := diaryService.ListDiariesByUser(testUser.UserId, 1, 3)
		require.NoError(t, err)
		assert.Len(t, diaries, 3)
		assert.Equal(t, 5, total)
	})

	t.Run("第二页，每页3个", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		diaryRepo := repository.NewDiaryRepository()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

		diaryService := diary.NewService(config, diaryRepo)
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)

		// 创建测试用户
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建5个日记本
		for i := 1; i <= 5; i++ {
			testDiary := &diary.Diary{
				Title:         fmt.Sprintf("diary%d", i),
				Remark:        "测试",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     testUser.UserId,
				AddUserName:   testUser.UserName,
				Status:        1,
			}
			_, err = diaryService.AddDiary(testDiary)
			require.NoError(t, err)
		}

		diaries, total, err := diaryService.ListDiariesByUser(testUser.UserId, 2, 3)
		require.NoError(t, err)
		assert.Len(t, diaries, 2)
		assert.Equal(t, 5, total)
	})

	t.Run("获取所有", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		diaryRepo := repository.NewDiaryRepository()
		userRepo := repository.NewUserRepository()
		taskRepo := repository.NewTaskRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

		diaryService := diary.NewService(config, diaryRepo)
		taskService := task.NewService(config, taskRepo)
		i18nManager := i18n.NewLocalizationManager(config)
		userService := user.NewService(config, userRepo, taskService, i18nManager, nil)

		// 创建测试用户
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建5个日记本
		for i := 1; i <= 5; i++ {
			testDiary := &diary.Diary{
				Title:         fmt.Sprintf("diary%d", i),
				Remark:        "测试",
				ViewSetting:   1,
				SocialSetting: 1,
				ThemeID:       1,
				AddUserID:     testUser.UserId,
				AddUserName:   testUser.UserName,
				Status:        1,
			}
			_, err = diaryService.AddDiary(testDiary)
			require.NoError(t, err)
		}

		diaries, total, err := diaryService.ListDiariesByUser(testUser.UserId, 1, 10)
		require.NoError(t, err)
		assert.Len(t, diaries, 5)
		assert.Equal(t, 5, total)
	})

	t.Run("用户不存在", func(t *testing.T) {
		// 使用新测试框架
		ctx := framework.NewTestContext(framework.WithDatabase())
		defer ctx.Cleanup()

		diaryRepo := repository.NewDiaryRepository()
		config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
		diaryService := diary.NewService(config, diaryRepo)
		// 使用一个极大的用户ID，确保不存在
		diaries, total, err := diaryService.ListDiariesByUser(999999999, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, diaries, 0)
		assert.Equal(t, 0, total)
	})
}

// TestDiaryService_EntryCountOperations 测试日记数量操作
func TestDiaryService_EntryCountOperations(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo)

	// 创建一个用户
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 添加一个日记本
	testDiary := &diary.Diary{
		Title:         ctx.GetDataGenerator().GenerateUniqueString("fullname"),
		Remark:        "测试日记数量操作",
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		Status:        1,
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	require.NoError(t, err)

	// 获取初始日记本状态
	initialDiary, err := diaryService.GetDiaryByID(diaryID)
	require.NoError(t, err)
	initialCount := initialDiary.InfoNum

	t.Run("增加日记计数", func(t *testing.T) {
		actionTime := int64(1642680100)
		err := diaryService.IncrementEntryCount(diaryID, testUser.UserId, testUser.UserName, actionTime)
		assert.NoError(t, err)

		// 验证计数是否增加
		updatedDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		assert.Equal(t, initialCount+1, updatedDiary.InfoNum)
		assert.Equal(t, testUser.UserId, updatedDiary.LastActionUserID)
		assert.Equal(t, testUser.UserName, updatedDiary.LastActionUserName)
		assert.Equal(t, actionTime, updatedDiary.LastActionUTCTime)
	})

	t.Run("减少日记计数", func(t *testing.T) {
		actionTime := int64(1642680200)

		// 获取当前计数
		currentDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		currentCount := currentDiary.InfoNum

		err = diaryService.DecrementEntryCount(diaryID, testUser.UserId, testUser.UserName, actionTime)
		assert.NoError(t, err)

		// 验证计数是否减少
		updatedDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		assert.Equal(t, currentCount-1, updatedDiary.InfoNum)
		assert.Equal(t, actionTime, updatedDiary.LastActionUTCTime)
	})

	t.Run("对不存在的日记本操作", func(t *testing.T) {
		actionTime := int64(1642680300)

		// 对不存在的日记本操作不会返回错误，只是不会影响任何行
		err := diaryService.IncrementEntryCount(99999, testUser.UserId, testUser.UserName, actionTime)
		assert.NoError(t, err)

		err = diaryService.DecrementEntryCount(99999, testUser.UserId, testUser.UserName, actionTime)
		assert.NoError(t, err)
	})

	t.Run("计数不会变为负数", func(t *testing.T) {
		// 先重置日记本计数为0（通过减少足够多次）
		currentDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)

		// 尝试减少计数，即使当前计数已经很低
		actionTime := int64(1642680400)
		for i := 0; i < int(currentDiary.InfoNum)+5; i++ { // 减少超过当前计数的次数
			err = diaryService.DecrementEntryCount(diaryID, testUser.UserId, testUser.UserName, actionTime)
			// 可能有些减少操作会成功，有些可能失败，这取决于具体实现
		}

		// 验证计数没有变为负数
		finalDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		assert.True(t, finalDiary.InfoNum >= 0, "日记计数不应该为负数")
	})
}
