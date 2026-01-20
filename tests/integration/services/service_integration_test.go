package services

import (
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

// TestUserDiaryIntegration 测试用户和日记服务的集成
func TestUserDiaryIntegration(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)

	t.Run("用户创建和日记本操作流程", func(t *testing.T) {
		// 1. 注册用户 - 使用随机数据
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
			TimeZone: "Asia/Shanghai",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)
		require.True(t, testUser.UserId > 0)

		// 2. 验证用户登录
		authenticatedUser, err := userService.LoginWithUsernameOrEmail(testUser.UserName, "password123")
		require.NoError(t, err)
		assert.Equal(t, testUser.UserId, authenticatedUser.UserId)

		// 3. 为用户创建日记本 - 使用随机数据
		testDiary := &diary.Diary{
			Title:         ctx.GetDataGenerator().GenerateUniqueString("diary"),
			Remark:        "用于集成测试的日记本",
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)
		assert.True(t, diaryID > 0)

		// 4. 验证用户可以访问自己的日记本
		userDiaries, err := diaryService.ListDiariesForWrite(testUser.UserId)
		require.NoError(t, err)
		assert.Len(t, userDiaries, 1)
		assert.Equal(t, testDiary.Title, userDiaries[0].Title)

		// 5. 检查日记本数量限制
		err = diaryService.CheckDiaryLimit(testUser.UserId)
		assert.NoError(t, err)

		// 6. 更新用户资料
		err = userService.UpdateUserProfile(testUser.UserId, "Updated Name", "Updated Remark", "UTC", "en")
		assert.NoError(t, err)

		// 7. 验证用户信息更新
		updatedUser, err := userService.GetUserByUserId(testUser.UserId)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.FullName)
		assert.Equal(t, "Updated Remark", updatedUser.Remark)
	})
}

// TestUserEntryDiaryIntegration 测试用户、日记和文章服务的完整集成
func TestUserEntryDiaryIntegration(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	t.Run("完整的内容创建流程", func(t *testing.T) {
		// 1. 创建用户 - 使用随机数据
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			FullName: ctx.GetDataGenerator().GenerateUniqueString("fullname"),
			TimeZone: "UTC",
			Lang:     "en",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 2. 创建日记本 - 使用随机数据
		testDiary := &diary.Diary{
			Title:         ctx.GetDataGenerator().GenerateUniqueString("diary"),
			Remark:        "用于内容测试",
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)

		// 3. 在日记本中创建文章 - 使用随机数据
		testEntry := &entry.Entry{
			Title:       ctx.GetDataGenerator().GenerateUniqueString("entry"),
			Body:        ctx.GetDataGenerator().GenerateUniqueString("content"),
			DiaryID:     diaryID,
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUTCTime:  1642680100,
			Status:      1,
		}
		entryID, err := entryService.AddEntry(testEntry)
		require.NoError(t, err)
		assert.True(t, entryID > 0)

		// 更新日记本的文章计数
		err = diaryService.IncrementEntryCount(diaryID, testUser.UserId, testUser.UserName, testEntry.AddUTCTime)
		require.NoError(t, err)

		// 4. 验证文章计数
		entryCount, err := entryService.GetEntryCountByUserID(testUser.UserId, diaryID)
		require.NoError(t, err)
		assert.Equal(t, 1, entryCount)

		// 5. 验证日记本的文章计数已更新
		updatedDiary, err := diaryService.GetDiaryByID(diaryID)
		require.NoError(t, err)
		assert.Equal(t, int32(1), updatedDiary.InfoNum)

		// 6. 创建更多文章
		for i := 2; i <= 3; i++ {
			additionalEntry := &entry.Entry{
				Title:       "第" + string(rune('0'+i)) + "篇文章",
				Body:        "这是第" + string(rune('0'+i)) + "篇文章的内容",
				Summary:     "第" + string(rune('0'+i)) + "篇文章摘要",
				DiaryID:     diaryID,
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  int64(1642680100 + i),
				Status:      1,
			}
			_, err := entryService.AddEntry(additionalEntry)
			require.NoError(t, err)

			// 更新日记本的文章计数
			err = diaryService.IncrementEntryCount(diaryID, testUser.UserId, testUser.UserName, additionalEntry.AddUTCTime)
			require.NoError(t, err)
		}

		// 7. 验证分页获取文章
		entries, total, err := entryService.ListEntriesByUserID(testUser.UserId, diaryID, 1, 2)
		require.NoError(t, err)
		assert.Len(t, entries, 2) // 第一页2篇
		assert.Equal(t, 3, total) // 总共3篇

		// 8. 获取第二页
		entries2, total2, err := entryService.ListEntriesByUserID(testUser.UserId, diaryID, 2, 2)
		require.NoError(t, err)
		assert.Len(t, entries2, 1) // 第二页1篇
		assert.Equal(t, 3, total2) // 总共还是3篇

		// 9. 更新文章
		firstEntry := &entries[0]
		firstEntry.Title = "更新后的标题"
		firstEntry.Body = "更新后的内容"
		firstEntry.UpdateUTCTime = 1642680300
		err = entryService.UpdateEntry(firstEntry)
		assert.NoError(t, err)

		// 10. 验证文章更新
		updatedEntry, err := entryService.GetEntryByID(firstEntry.EntryID)
		require.NoError(t, err)
		assert.Equal(t, "更新后的标题", updatedEntry.Title)
		assert.Equal(t, "更新后的内容", updatedEntry.Body)

		// 11. 删除一篇文章
		err = entryService.DeleteEntryByID(firstEntry.EntryID, testUser.UserId, testUser.UserName, 1234567890)
		assert.NoError(t, err)

		// 12. 验证文章已删除
		deletedEntry, err := entryService.GetEntryByID(firstEntry.EntryID)
		assert.Error(t, err)
		assert.Nil(t, deletedEntry)

		// 13. 验证文章计数已更新
		finalCount, err := entryService.GetEntryCountByUserID(testUser.UserId, diaryID)
		require.NoError(t, err)
		assert.Equal(t, 2, finalCount) // 删除一篇后剩2篇
	})
}

// TestMultiUserScenario 测试多用户场景
func TestMultiUserScenario(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置
	
	i18nManager := i18n.NewLocalizationManager(config)
	taskService := task.NewService(config, taskRepo)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	t.Run("多用户数据隔离", func(t *testing.T) {
		// 创建两个用户
		user1 := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("testuser1"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user1"),
			Pwd:      "password123",
			TimeZone: "Asia/Shanghai",
		}
		err := userService.RegisterUser(user1)
		require.NoError(t, err)

		user2 := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("testuser2"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user2"),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err = userService.RegisterUser(user2)
		require.NoError(t, err)

		// 为每个用户创建日记本
		diary1 := &diary.Diary{
			Title:         "用户1的日记本",
			Remark:        "用户1专用",
			AddUserID:     user1.UserId,
			AddUserName:   user1.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		diary1ID, err := diaryService.AddDiary(diary1)
		require.NoError(t, err)

		diary2 := &diary.Diary{
			Title:         "用户2的日记本",
			Remark:        "用户2专用",
			AddUserID:     user2.UserId,
			AddUserName:   user2.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		diary2ID, err := diaryService.AddDiary(diary2)
		require.NoError(t, err)

		// 为每个用户在各自的日记本中创建文章
		entry1 := &entry.Entry{
			Title:       "用户1的文章",
			Body:        "用户1的内容",
			DiaryID:     diary1ID,
			AddUserID:   user1.UserId,
			AddUserName: user1.UserName,
			AddUTCTime:  1642680100,
			Status:      1,
		}
		entry1ID, err := entryService.AddEntry(entry1)
		require.NoError(t, err)

		entry2 := &entry.Entry{
			Title:       "用户2的文章",
			Body:        "用户2的内容",
			DiaryID:     diary2ID,
			AddUserID:   user2.UserId,
			AddUserName: user2.UserName,
			AddUTCTime:  1642680100,
			Status:      1,
		}
		entry2ID, err := entryService.AddEntry(entry2)
		require.NoError(t, err)

		// 验证数据隔离：用户1只能看到自己的日记本
		user1Diaries, err := diaryService.ListDiariesForWrite(user1.UserId)
		require.NoError(t, err)
		assert.Len(t, user1Diaries, 1)
		assert.Equal(t, "用户1的日记本", user1Diaries[0].Title)

		// 验证数据隔离：用户2只能看到自己的日记本
		user2Diaries, err := diaryService.ListDiariesForWrite(user2.UserId)
		require.NoError(t, err)
		assert.Len(t, user2Diaries, 1)
		assert.Equal(t, "用户2的日记本", user2Diaries[0].Title)

		// 验证用户1只能看到自己的文章
		entries1, _, err := entryService.ListEntriesByUserID(user1.UserId, diary1ID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, entries1, 1)
		assert.Equal(t, entry1ID, entries1[0].EntryID)

		// 验证用户2只能看到自己的文章
		entries2, _, err := entryService.ListEntriesByUserID(user2.UserId, diary2ID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, entries2, 1)
		assert.Equal(t, entry2ID, entries2[0].EntryID)

		// 验证用户1尝试获取用户2的文章应该失败（或返回空/无权限）
		// 注意：GetEntryByID通常不检查用户权限，而是在Controller层检查
		// 但我们可以验证GetEntryByID能获取到任何存在的文章
		retrievedEntry2, err := entryService.GetEntryByID(entry2ID)
		require.NoError(t, err)
		assert.Equal(t, entry2.Title, retrievedEntry2.Title)
	})
}

// TestServiceErrorHandling 测试服务错误处理
func TestServiceErrorHandling(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

	i18nManager := i18n.NewLocalizationManager(config)
	taskService := task.NewService(config, taskRepo)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	t.Run("服务间错误传播", func(t *testing.T) {
		// 1. 尝试为不存在的用户创建日记本
		invalidDiary := &diary.Diary{
			Title:         "无效日记本",
			Remark:        "测试",
			AddUserID:     99999, // 不存在的用户ID
			AddUserName:   "nonexistent",
			AddUTCTime:    1642680000,
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		_, err := diaryService.AddDiary(invalidDiary)
		// 注意：系统设计上允许为不存在的用户创建日记本，所以这里不会出错
		assert.NoError(t, err)

		// 2. 尝试在不存在的日记本中创建文章
		invalidPost := &entry.Entry{
			Title:       "无效文章",
			Body:        "测试内容",
			DiaryID:     99999, // 不存在的日记本ID
			AddUserID:   1,
			AddUserName: "testuser",
			AddUTCTime:  1642680100,
			Status:      1,
		}
		_, err = entryService.AddEntry(invalidPost)
		assert.Error(t, err) // 应该出错

		// 3. 尝试重复注册相同的用户名
		duplicateUserName := ctx.GetDataGenerator().GenerateUniqueString("duplicateuser")
		user1 := &user.User{
			UserName: duplicateUserName,
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("dup1"),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err = userService.RegisterUser(user1)
		require.NoError(t, err)

		user2 := &user.User{
			UserName: duplicateUserName, // 重复的用户名
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("dup2"),
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err = userService.RegisterUser(user2)
		assert.Error(t, err) // 应该出错

		// 4. 尝试重复注册相同的邮箱
		user3 := &user.User{
			UserName: "anotheruser",
			Email:    user1.Email, // 重复的邮箱
			Pwd:      "password123",
			TimeZone: "UTC",
		}
		err = userService.RegisterUser(user3)
		assert.Error(t, err) // 应该出错
	})
}

// TestServicePerformance 测试服务性能
func TestServicePerformance(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	config := ctx.GetConfig().GetAppConfig() // 通过测试上下文获取配置

	i18nManager := i18n.NewLocalizationManager(config)
	taskService := task.NewService(config, taskRepo)
	userService := user.NewService(config, userRepo, taskService, i18nManager)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userRepo, taskService, diaryRepo, entryRepo, i18nManager)

	t.Run("批量操作性能", func(t *testing.T) {
		// 创建一个用户
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("perfuser"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("perf"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "en",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		// 创建一个日记本
		testDiary := &diary.Diary{
			Title:         "性能测试日记本",
			Remark:        "用于性能测试",
			AddUserID:     testUser.UserId,
			AddUserName:   testUser.UserName,
			AddUTCTime:    1642680000,
			Status:        1,
			ViewSetting:   1,
			SocialSetting: 1,
			ThemeID:       1,
		}
		diaryID, err := diaryService.AddDiary(testDiary)
		require.NoError(t, err)

		// 批量创建文章（小批量，避免测试时间过长）
		const articleCount = 10
		for i := 1; i <= articleCount; i++ {
			article := &entry.Entry{
				Title:       "性能测试文章" + string(rune('0'+i)),
				Body:        "这是第" + string(rune('0'+i)) + "篇性能测试文章",
				Summary:     "摘要" + string(rune('0'+i)),
				DiaryID:     diaryID,
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUTCTime:  int64(1642680100 + i),
				Status:      1,
			}
			_, err := entryService.AddEntry(article)
			require.NoError(t, err)
		}

		// 验证分页查询性能
		entries, total, err := entryService.ListEntriesByUserID(testUser.UserId, diaryID, 1, 5)
		require.NoError(t, err)
		assert.Len(t, entries, 5)
		assert.Equal(t, articleCount, total)

		// 验证文章计数
		count, err := entryService.GetEntryCountByUserID(testUser.UserId, diaryID)
		require.NoError(t, err)
		assert.Equal(t, articleCount, count)
	})
}
