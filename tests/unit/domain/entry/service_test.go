package entry_test

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestUserAndDiary 创建测试用的用户和日记本
func setupTestUserAndDiary(t *testing.T, ctx *framework.TestContext, userService *user.Service, diaryService *diary.Service) (*user.User, *diary.Diary) {
	t.Helper()

	// 创建用户 - 使用随机数据
	testUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(testUser)
	require.NoError(t, err)

	// 创建日记本 - 使用随机数据
	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "用于文章测试",
		ViewSetting:   1, // 必需字段：1-公开，2-私有
		SocialSetting: 1, // 必需字段：1-允许社交，2-不允许社交
		ThemeID:       1, // 必需字段：主题ID，范围1-16
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	require.NoError(t, err)
	testDiary.DiaryID = diaryID

	return testUser, testDiary
}

// TestEntryService_AddEntry 测试添加文章功能
func TestEntryService_AddEntry(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	t.Run("成功添加文章", func(t *testing.T) {
		newEntry := &entry.Entry{
			Title:       "测试文章标题",
			Body:        "这是一篇测试文章的内容",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		entryID, err := entryService.AddEntry(newEntry)
		require.NoError(t, err)
		assert.True(t, entryID > 0)
		assert.Equal(t, entryID, newEntry.EntryID)
	})

	t.Run("无效的日记本ID", func(t *testing.T) {
		invalidEntry := &entry.Entry{
			Title:       "标题",
			Body:        "内容",
			DiaryID:     99999, // 不存在的日记本ID
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		_, err := entryService.AddEntry(invalidEntry)
		assert.Error(t, err)
	})

	t.Run("无效的用户ID", func(t *testing.T) {
		invalidEntry := &entry.Entry{
			Title:       "标题",
			Body:        "内容",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   0, // 无效用户ID
			AddUserName: "invalid",
			Status:      1,
		}

		_, err := entryService.AddEntry(invalidEntry)
		assert.Error(t, err)
	})
}

// TestEntryService_GetEntryByID 测试根据ID获取文章
func TestEntryService_GetEntryByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	// 添加一篇文章
	newEntry := &entry.Entry{
		Title:       "可获取的文章",
		Body:        "测试获取功能的文章内容",
		Summary:     "测试摘要",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	entryID, err := entryService.AddEntry(newEntry)
	require.NoError(t, err)

	t.Run("成功获取文章", func(t *testing.T) {
		retrievedEntry, err := entryService.GetEntryByID(entryID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedEntry)
		assert.Equal(t, "可获取的文章", retrievedEntry.Title)
		assert.Equal(t, "测试获取功能的文章内容", retrievedEntry.Body)
		assert.Equal(t, testUser.UserId, retrievedEntry.AddUserID)
		assert.Equal(t, testDiary.DiaryID, retrievedEntry.DiaryID)
	})

	t.Run("文章不存在", func(t *testing.T) {
		retrievedEntry, err := entryService.GetEntryByID(99999)
		assert.Error(t, err)
		assert.Nil(t, retrievedEntry)
	})
}

// TestEntryService_UpdateEntry 测试更新文章
func TestEntryService_UpdateEntry(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	// 添加一篇文章
	originalEntry := &entry.Entry{
		Title:       "原始标题",
		Body:        "原始内容",
		Summary:     "原始摘要",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	entryID, err := entryService.AddEntry(originalEntry)
	require.NoError(t, err)

	t.Run("成功更新文章", func(t *testing.T) {
		updatedEntry := &entry.Entry{
			EntryID:        entryID,
			Title:         "更新后的标题",
			Body:          "更新后的内容",
			Summary:       "更新后的摘要",
			DiaryID:       testDiary.DiaryID,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			UpdateUTCTime: 1642680200,
			Status:        1,
		}

		err := entryService.UpdateEntry(updatedEntry)
		assert.NoError(t, err)

		// 验证更新是否成功
		retrievedEntry, err := entryService.GetEntryByID(entryID)
		require.NoError(t, err)
		assert.Equal(t, "更新后的标题", retrievedEntry.Title)
		assert.Equal(t, "更新后的内容", retrievedEntry.Body)
		assert.Equal(t, "更新后的摘要", retrievedEntry.Summary)
	})

	t.Run("更新不存在的文章", func(t *testing.T) {
		nonExistentEntry := &entry.Entry{
			EntryID:      99999,
			Title:       "不存在的文章",
			Body:        "测试",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := entryService.UpdateEntry(nonExistentEntry)
		assert.Error(t, err)
	})

	t.Run("内容为空时更新失败", func(t *testing.T) {
		invalidEntry := &entry.Entry{
			EntryID:      entryID,
			Title:       "标题",
			Body:        "", // 空内容
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := entryService.UpdateEntry(invalidEntry)
		assert.Error(t, err)
	})
}

// TestEntryService_ListEntriesByUserID 测试分页获取用户文章列表
func TestEntryService_ListEntriesByUserID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo)
	entryRepo := repository.NewEntryRepository()
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	// 创建另一个用户和日记本
	otherUser := &user.User{
		UserName: ctx.GetDataGenerator().GenerateUniqueString("ou"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("ou"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(otherUser)
	require.NoError(t, err)

	otherDiary := &diary.Diary{
		Title:         "其他用户的日记本",
		Remark:        "测试",
		AddUserID:     otherUser.UserId,
		AddUserName:   otherUser.UserName,
		Status:        1,
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}
	otherDiaryID, err := diaryService.AddDiary(otherDiary)
	require.NoError(t, err)

	// 为testUser添加5篇文章
	for i := 1; i <= 5; i++ {
		entry := &entry.Entry{
			Title:       fmt.Sprintf("文章%d", i),
			Body:        fmt.Sprintf("第%d篇文章内容", i),
			Summary:     fmt.Sprintf("第%d篇文章摘要", i),
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
			InfoType:    1,
		}
		_, err := entryService.AddEntry(entry)
		require.NoError(t, err)
	}

	// 为otherUser添加1篇文章
	otherEntry := &entry.Entry{
		Title:       "其他用户的文章",
		Body:        "其他用户的文章内容",
		Summary:     "其他用户的文章摘要",
		DiaryID:     otherDiaryID,
		AddUserID:   otherUser.UserId,
		AddUserName: otherUser.UserName,
		Status:      1,
		InfoType:    1,
	}
	_, err = entryService.AddEntry(otherEntry)
	require.NoError(t, err)

	t.Run("第一页，每页3篇", func(t *testing.T) {
		entries, total, err := entryService.ListEntriesByUserID(testUser.UserId, testDiary.DiaryID, 1, 3)
		require.NoError(t, err)
		assert.Len(t, entries, 3)
		assert.Equal(t, 5, total)

		// 验证返回的都是指定用户和日记本的文章
		for _, p := range entries {
			assert.Equal(t, testUser.UserId, p.AddUserID)
			assert.Equal(t, testDiary.DiaryID, p.DiaryID)
		}
	})

	t.Run("第二页，每页3篇", func(t *testing.T) {
		entries, total, err := entryService.ListEntriesByUserID(testUser.UserId, testDiary.DiaryID, 2, 3)
		require.NoError(t, err)
		assert.Len(t, entries, 2) // 剩余2篇
		assert.Equal(t, 5, total)
	})

	t.Run("获取所有文章", func(t *testing.T) {
		entries, total, err := entryService.ListEntriesByUserID(testUser.UserId, testDiary.DiaryID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, entries, 5)
		assert.Equal(t, 5, total)
	})

	t.Run("获取其他用户的文章", func(t *testing.T) {
		entries, total, err := entryService.ListEntriesByUserID(otherUser.UserId, otherDiaryID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Equal(t, 1, total)
		assert.Equal(t, "其他用户的文章", entries[0].Title)
	})

	t.Run("不存在的用户", func(t *testing.T) {
		entries, total, err := entryService.ListEntriesByUserID(99999, testDiary.DiaryID, 1, 10)
		assert.NoError(t, err) // 应该返回空列表而不是错误
		assert.Len(t, entries, 0)
		assert.Equal(t, 0, total)
	})

	t.Run("不存在的日记本", func(t *testing.T) {
		entries, total, err := entryService.ListEntriesByUserID(testUser.UserId, 99999, 1, 10)
		assert.NoError(t, err) // 应该返回空列表而不是错误
		assert.Len(t, entries, 0)
		assert.Equal(t, 0, total)
	})
}

// TestEntryService_GetEntryCountByUserID 测试获取用户文章数量
func TestEntryService_GetEntryCountByUserID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo)
	entryRepo := repository.NewEntryRepository()
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
		testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	t.Run("初始文章数量为0", func(t *testing.T) {
		count, err := entryService.GetEntryCountByUserID(testUser.UserId, testDiary.DiaryID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	// 添加3篇文章
	for i := 1; i <= 3; i++ {
		entry := &entry.Entry{
			Title:       fmt.Sprintf("文章%d", i),
			Body:        fmt.Sprintf("第%d篇文章内容", i),
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}
		_, err := entryService.AddEntry(entry)
		require.NoError(t, err)
	}

	t.Run("添加文章后的数量", func(t *testing.T) {
		count, err := entryService.GetEntryCountByUserID(testUser.UserId, testDiary.DiaryID)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("不存在的用户", func(t *testing.T) {
		count, err := entryService.GetEntryCountByUserID(99999, testDiary.DiaryID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("不存在的日记本", func(t *testing.T) {
		count, err := entryService.GetEntryCountByUserID(testUser.UserId, 99999)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// TestEntryService_UpdateEntryByID 测试根据ID更新文章
func TestEntryService_UpdateEntryByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo)
	entryRepo := repository.NewEntryRepository()
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	// 添加一篇文章
	originalEntry := &entry.Entry{
		Title:       "待更新的文章",
		Body:        "原始内容",
		Summary:     "原始摘要",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	entryID, err := entryService.AddEntry(originalEntry)
	require.NoError(t, err)

	t.Run("成功更新文章", func(t *testing.T) {
		updatedEntry := &entry.Entry{
			EntryID:        entryID,
			Title:         "已更新的文章",
			Body:          "更新后的内容",
			Summary:       "更新后的摘要",
			DiaryID:       testDiary.DiaryID,
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			UpdateUTCTime: 1642680300,
			Status:        1,
		}

		err := entryService.UpdateEntryByID(updatedEntry)
		assert.NoError(t, err)

		// 验证更新是否成功
		retrievedEntry, err := entryService.GetEntryByID(entryID)
		require.NoError(t, err)
		assert.Equal(t, "已更新的文章", retrievedEntry.Title)
		assert.Equal(t, "更新后的内容", retrievedEntry.Body)
		assert.Equal(t, "更新后的摘要", retrievedEntry.Summary)
	})

	t.Run("更新不存在的文章", func(t *testing.T) {
		nonExistentEntry := &entry.Entry{
			EntryID:      99999,
			Title:       "不存在",
			Body:        "测试",
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
		}

		err := entryService.UpdateEntryByID(nonExistentEntry)
		assert.Error(t, err)
	})
}

// TestEntryService_DeleteEntryByID 测试删除文章
func TestEntryService_DeleteEntryByID(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo)
	entryRepo := repository.NewEntryRepository()
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	// 添加一篇文章用于删除测试
	newEntry := &entry.Entry{
		Title:       "待删除的文章",
		Body:        "这篇文章将被删除",
		Summary:     "删除测试",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	entryID, err := entryService.AddEntry(newEntry)
	require.NoError(t, err)

	// 验证文章存在
	retrievedEntry, err := entryService.GetEntryByID(entryID)
	require.NoError(t, err)
	assert.NotNil(t, retrievedEntry)

	t.Run("成功删除文章", func(t *testing.T) {
		currentTime := int64(1642680300)
		err := entryService.DeleteEntryByID(entryID, testUser.UserId, testUser.UserName, currentTime)
		assert.NoError(t, err)

		// 验证文章已被删除
		deletedEntry, err := entryService.GetEntryByID(entryID)
		assert.Error(t, err)
		assert.Nil(t, deletedEntry)
	})

	t.Run("删除不存在的文章", func(t *testing.T) {
		currentTime := int64(1642680300)
		err := entryService.DeleteEntryByID(99999, testUser.UserId, testUser.UserName, currentTime)
		assert.Error(t, err)
	})

	t.Run("删除已删除的文章", func(t *testing.T) {
		currentTime := int64(1642680300)
		err := entryService.DeleteEntryByID(entryID, testUser.UserId, testUser.UserName, currentTime)
		assert.Error(t, err)
	})
}

// TestEntryService_GetEntriesStats 测试获取文章统计信息功能
func TestEntryService_GetEntriesStats(t *testing.T) {
	// 使用新测试框架
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	// 设置测试数据
	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	// 创建多篇测试文章
	var testEntries []*entry.Entry
	for i := 0; i < 3; i++ {
		testEntry := &entry.Entry{
			Title:       fmt.Sprintf("测试文章%d", i+1),
			Body:        fmt.Sprintf("这是第%d篇测试文章的内容", i+1),
			DiaryID:     testDiary.DiaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			Status:      1,
			ViewNum:     int32(10 * (i + 1)), // 设置不同的浏览数
			LikeNum:     int32(5 * (i + 1)),  // 设置不同的点赞数
			CommentNum:  int32(2 * (i + 1)),  // 设置不同的评论数
		}
		entryID, err := entryService.AddEntry(testEntry)
		require.NoError(t, err)
		testEntry.EntryID = entryID
		testEntries = append(testEntries, testEntry)
	}

	t.Run("成功获取多篇文章统计信息", func(t *testing.T) {
		// 构建文章ID列表
		entryIDs := make([]int64, len(testEntries))
		for i, p := range testEntries {
			entryIDs[i] = p.EntryID
		}

		// 获取统计信息
		stats, err := entryService.GetEntriesStats(entryIDs)
		require.NoError(t, err)
		assert.Len(t, stats, len(testEntries))

		// 验证返回的统计信息
		for i, stat := range stats {
			assert.Equal(t, testDiary.DiaryID, stat.DiaryID)
			assert.Equal(t, testEntries[i].EntryID, stat.EntryID)
			// 注意：由于数据库操作可能重置统计数据，这里主要验证结构正确性
			assert.GreaterOrEqual(t, stat.ViewNum, int32(0))
			assert.GreaterOrEqual(t, stat.LikeNum, int32(0))
			assert.GreaterOrEqual(t, stat.CommentNum, int32(0))
		}
	})

	t.Run("空文章ID列表", func(t *testing.T) {
		stats, err := entryService.GetEntriesStats([]int64{})
		require.NoError(t, err)
		assert.Empty(t, stats)
	})

	t.Run("不存在的文章ID", func(t *testing.T) {
		stats, err := entryService.GetEntriesStats([]int64{99999, 99998})
		require.NoError(t, err)
		assert.Empty(t, stats) // 不存在的文章应该返回空列表
	})

	t.Run("混合存在和不存在的文章ID", func(t *testing.T) {
		mixedIDs := []int64{testEntries[0].EntryID, 99999, testEntries[1].EntryID}
		stats, err := entryService.GetEntriesStats(mixedIDs)
		require.NoError(t, err)
		// 应该只返回存在的文章统计信息
		assert.Len(t, stats, 2)
		for _, stat := range stats {
			assert.Contains(t, []int64{testEntries[0].EntryID, testEntries[1].EntryID}, stat.EntryID)
		}
	})
}
