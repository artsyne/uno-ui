package diary_test

import (
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	framework.RunTestsWithFlags(m)
}

// TestDiaryService_AddDiary 测试添加日记本
func TestDiaryService_AddDiary(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if err := ctx.Cleanup(); err != nil {
			t.Log(err)
		}
	}()

	// 准备服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 先注册一个用户
	newUser := &user.User{
		Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
	}
	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	t.Run("成功添加日记本", func(t *testing.T) {
		newDiary := &diary.Diary{
			Title:         "My Diary",
			AddUserID:     newUser.UserId,
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
			AddUTCTime:    time.Now().Unix(),
		}

		id, err := diaryService.AddDiary(newDiary)
		require.NoError(t, err)
		assert.True(t, id > 0)
		assert.Equal(t, id, newDiary.DiaryID)
	})

	t.Run("标题为空", func(t *testing.T) {
		invalidDiary := &diary.Diary{
			Title:         "",
			AddUserID:     newUser.UserId,
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
		}

		id, err := diaryService.AddDiary(invalidDiary)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("用户ID无效", func(t *testing.T) {
		invalidDiary := &diary.Diary{
			Title:         "My Diary",
			AddUserID:     0,
			ViewSetting:   diary.ViewSettingPublic,
			SocialSetting: 1,
			ThemeID:       1,
		}

		id, err := diaryService.AddDiary(invalidDiary)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}

// TestDiaryService_UpdateDiary 测试更新日记本
func TestDiaryService_UpdateDiary(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if err := ctx.Cleanup(); err != nil {
			t.Log(err)
		}
	}()

	// 准备服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 先注册一个用户
	newUser := &user.User{
		Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
	}
	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	// 先添加一个日记本
	newDiary := &diary.Diary{
		Title:         "Original Title",
		AddUserID:     newUser.UserId,
		ViewSetting:   diary.ViewSettingPublic,
		SocialSetting: 1,
		ThemeID:       1,
		AddUTCTime:    time.Now().Unix(),
	}
	id, err := diaryService.AddDiary(newDiary)
	require.NoError(t, err)

	t.Run("成功更新日记本", func(t *testing.T) {
		newDiary.Title = "Updated Title"
		err := diaryService.UpdateDiary(newDiary)
		require.NoError(t, err)

		// 验证更新
		updatedDiary, err := diaryService.GetDiaryByID(id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updatedDiary.Title)
	})
}

// TestDiaryService_DeleteDiary 测试删除日记本
func TestDiaryService_DeleteDiary(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if err := ctx.Cleanup(); err != nil {
			t.Log(err)
		}
	}()

	// 准备服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 先注册一个用户
	newUser := &user.User{
		Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
	}
	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	// 先添加一个日记本
	newDiary := &diary.Diary{
		Title:         "To Delete",
		AddUserID:     newUser.UserId,
		ViewSetting:   diary.ViewSettingPublic,
		SocialSetting: 1,
		ThemeID:       1,
		AddUTCTime:    time.Now().Unix(),
	}
	id, err := diaryService.AddDiary(newDiary)
	require.NoError(t, err)

	t.Run("成功删除日记本", func(t *testing.T) {
		err := diaryService.DeleteDiary(id, newUser.UserId)
		require.NoError(t, err)

		// 验证已删除
		_, err = diaryService.GetDiaryByID(id)
		assert.Error(t, err)
	})

	t.Run("非拥有者无法删除", func(t *testing.T) {
		// 再添加一个日记本
		anotherDiary := &diary.Diary{
			Title:         "Another Diary",
		AddUserID:     newUser.UserId,
		ViewSetting:   diary.ViewSettingPublic,
		SocialSetting: 1,
			ThemeID:       1,
			AddUTCTime:    time.Now().Unix(),
		}
		id, err := diaryService.AddDiary(anotherDiary)
		require.NoError(t, err)

		// 尝试用另一个用户ID删除
		err = diaryService.DeleteDiary(id, newUser.UserId+1)
		assert.Error(t, err)
	})
}

// TestDiaryService_IncrementEntryCount 测试增加日记数量
func TestDiaryService_IncrementEntryCount(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if err := ctx.Cleanup(); err != nil {
			t.Log(err)
		}
	}()

	// 准备服务
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
	
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// 先注册一个用户
	newUser := &user.User{
		Handle: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
	}
	err := userService.RegisterUser(newUser)
	require.NoError(t, err)

	// 先添加一个日记本
	newDiary := &diary.Diary{
		Title:         "Diary For Counts",
		AddUserID:     newUser.UserId,
		ViewSetting:   diary.ViewSettingPublic,
		SocialSetting: 1,
		ThemeID:       1,
		AddUTCTime:    time.Now().Unix(),
		InfoNum:       0,
	}
	id, err := diaryService.AddDiary(newDiary)
	require.NoError(t, err)

	t.Run("成功增加计数", func(t *testing.T) {
		// 移除userName参数
		err := diaryService.IncrementEntryCount(id)
		require.NoError(t, err)

		// 验证计数增加
		updatedDiary, err := diaryService.GetDiaryByID(id)
		require.NoError(t, err)
		assert.Equal(t, int32(1), updatedDiary.InfoNum)
	})
}
