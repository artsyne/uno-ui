package services_test

import (
	"fmt"
	"testing"

	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommentServiceIntegration 评论服务集成测试
// 测试评论服务与数据库、其他服务的完整集成
func TestCommentServiceIntegration(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	commentRepo := repository.NewCommentRepository()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
	
	userService := user.NewService(config, userRepo, taskService, i18nManager, redisClient)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userService, taskService, diaryRepo, entryRepo, i18nManager)
	commentService := comment.NewService(config, commentRepo, userService, entryService, diaryService)

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

	// 创建测试日记本
	testDiary := &diary.Diary{
		Title:         "集成测试日记本",
		Remark:        "用于评论集成测试",
		ViewSetting:   1, // 公开
		SocialSetting: 1, // 允许社交
		ThemeID:       1,
		AddUserID:     testUser.UserId,
		AddUserName:   testUser.UserName,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	require.NoError(t, err)
	testDiary.DiaryID = diaryID

	// 创建测试文章
	testEntry := &entry.Entry{
		Title:       "集成测试文章",
		Body:        "这是一篇用于评论集成测试的文章",
		DiaryID:     testDiary.DiaryID,
		AddUserID:   testUser.UserId,
		AddUserName: testUser.UserName,
		Status:      1,
	}
	entryID, err := entryService.AddEntry(testEntry)
	require.NoError(t, err)
	testEntry.EntryID = entryID

	t.Run("完整的评论生命周期测试", func(t *testing.T) {
		// 1. 添加评论
		newComment := &comment.Comment{
			EntryID:     testEntry.EntryID,
			Content:     "这是一条集成测试评论",
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUserIP:   "127.0.0.1",
		}
		commentID, err := commentService.AddComment(newComment)
		assert.NoError(t, err)
		assert.Greater(t, commentID, int64(0))

		// 2. 获取评论
		savedComment, err := commentService.GetCommentByID(commentID)
		assert.NoError(t, err)
		assert.Equal(t, newComment.EntryID, savedComment.EntryID)
		assert.Equal(t, newComment.Content, savedComment.Content)
		assert.Equal(t, newComment.AddUserID, savedComment.AddUserID)
		assert.NotEmpty(t, savedComment.AddUtcTime)

		// 3. 验证评论数量
		count, err := commentService.GetCommentCountByEntryID(testEntry.EntryID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		// 4. 获取评论列表
		commentList, total, err := commentService.ListCommentsByEntryID(testEntry.EntryID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, commentList, 1)
		assert.Equal(t, commentID, commentList[0].CommentID)

		// 5. 删除评论
		err = commentService.DeleteComment(commentID, testUser.UserId)
		assert.NoError(t, err)

		// 6. 验证评论已被删除
		_, err = commentService.GetCommentByID(commentID)
		assert.Error(t, err)

		// 7. 验证评论数量已减少
		count, err = commentService.GetCommentCountByEntryID(testEntry.EntryID)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("多用户评论交互测试", func(t *testing.T) {
		// 创建第二个用户
		testUser2 := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("u2"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("u2"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(testUser2)
		require.NoError(t, err)

		// 两个用户分别添加评论
		comment1 := &comment.Comment{
			EntryID:     testEntry.EntryID,
			Content:     "用户1的评论",
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUserIP:   "127.0.0.1",
		}
		commentID1, err := commentService.AddComment(comment1)
		assert.NoError(t, err)

		comment2 := &comment.Comment{
			EntryID:     testEntry.EntryID,
			Content:     "用户2的评论",
			AddUserID:   testUser2.UserId,
			AddUserName: testUser2.UserName,
			AddUserIP:   "127.0.0.2",
		}
		commentID2, err := commentService.AddComment(comment2)
		assert.NoError(t, err)

		// 验证评论总数
		count, err := commentService.GetCommentCountByEntryID(testEntry.EntryID)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)

		// 获取评论列表
		commentList, total, err := commentService.ListCommentsByEntryID(testEntry.EntryID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, commentList, 2)

		// 用户1删除自己的评论
		err = commentService.DeleteComment(commentID1, testUser.UserId)
		assert.NoError(t, err)

		// 验证只剩下用户2的评论
		count, err = commentService.GetCommentCountByEntryID(testEntry.EntryID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		// 清理剩余评论
		err = commentService.DeleteComment(commentID2, testUser2.UserId)
		assert.NoError(t, err)
	})

	t.Run("分页功能测试", func(t *testing.T) {
		// 创建多条评论
		var commentIDs []int64
		for i := 1; i <= 5; i++ {
			testComment := &comment.Comment{
				EntryID:     testEntry.EntryID,
				Content:     "分页测试评论" + string(rune(i)),
				AddUserID:   testUser.UserId,
				AddUserName: testUser.UserName,
				AddUserIP:   "127.0.0.1",
			}
			commentID, err := commentService.AddComment(testComment)
			assert.NoError(t, err)
			commentIDs = append(commentIDs, commentID)
		}

		// 测试第一页（每页2条）
		commentList, total, err := commentService.ListCommentsByEntryID(testEntry.EntryID, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, commentList, 2)

		// 测试第二页
		commentList, total, err = commentService.ListCommentsByEntryID(testEntry.EntryID, 2, 2)
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, commentList, 2)

		// 测试第三页（最后一页，只有1条）
		commentList, total, err = commentService.ListCommentsByEntryID(testEntry.EntryID, 3, 2)
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, commentList, 1)

		// 清理测试数据
		for _, commentID := range commentIDs {
			err = commentService.DeleteComment(commentID, testUser.UserId)
			assert.NoError(t, err)
		}
	})
}

// TestCommentServiceErrorHandling 评论服务错误处理集成测试
func TestCommentServiceErrorHandling(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 创建真实的服务
	config := ctx.GetConfig().GetAppConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	commentRepo := repository.NewCommentRepository()
	
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
	
	userService := user.NewService(config, userRepo, taskService, i18nManager, redisClient)
	diaryService := diary.NewService(config, diaryRepo)
	entryService := entry.NewService(config, userService, taskService, diaryRepo, entryRepo, i18nManager)
	commentService := comment.NewService(config, commentRepo, userService, entryService, diaryService)

	t.Run("无效用户ID错误处理", func(t *testing.T) {
		invalidComment := &comment.Comment{
			EntryID:     1,
			Content:     "测试评论",
			AddUserID:   99999, // 不存在的用户ID
			AddUserName: "不存在的用户",
			AddUserIP:   "127.0.0.1",
		}
		_, err := commentService.AddComment(invalidComment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_not_found")
	})

	t.Run("无效文章ID错误处理", func(t *testing.T) {
		// 先创建一个有效用户
		testUser := &user.User{
			UserName: ctx.GetDataGenerator().GenerateUniqueString("user"),
			Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user"),
			Pwd:      "password123",
			TimeZone: "UTC",
			Lang:     "zh-cn",
		}
		err := userService.RegisterUser(testUser)
		require.NoError(t, err)

		invalidComment := &comment.Comment{
			EntryID:     99999, // 不存在的文章ID
			Content:     "测试评论",
			AddUserID:   testUser.UserId,
			AddUserName: testUser.UserName,
			AddUserIP:   "127.0.0.1",
		}
		_, err = commentService.AddComment(invalidComment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "entry_not_found")
	})

	t.Run("获取不存在的评论", func(t *testing.T) {
		_, err := commentService.GetCommentByID(99999)
		assert.Error(t, err)
	})

	t.Run("删除不存在的评论", func(t *testing.T) {
		err := commentService.DeleteComment(99999, 1)
		assert.Error(t, err)
	})
}