package services_test

import (
	"fmt"
	"testing"
	"time"

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

func TestCommentOnSoftDeletedEntry(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// 初始化服务
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

	// 创建主用户（日记作者）
	user1 := &user.User{
		Handle:   ctx.GetDataGenerator().GenerateUniqueString("u1"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user1"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err := userService.RegisterUser(user1)
	require.NoError(t, err)

	// 创建第二个用户（旁观者）
	user2 := &user.User{
		Handle:   ctx.GetDataGenerator().GenerateUniqueString("u2"),
		Email:    ctx.GetDataGenerator().GenerateUniqueEmail("user2"),
		Pwd:      "password123",
		TimeZone: "UTC",
		Lang:     "zh-cn",
	}
	err = userService.RegisterUser(user2)
	require.NoError(t, err)

	// 创建日记本
	testDiary := &diary.Diary{
		Title:         "Soft Delete Test Diary",
		Remark:        "For testing comments on soft deleted entries",
		ViewSetting:   1, // 公开
		SocialSetting: 1, // 允许社交
		ThemeID:       1,
		AddUserID:     user1.UserId,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(testDiary)
	require.NoError(t, err)
	testDiary.DiaryID = diaryID

	// 1. 测试草稿状态 (StatusDraft = 2)
	t.Run("Draft Entry Comments", func(t *testing.T) {
		// 创建草稿
		draftEntry := &entry.Entry{
			DiaryID:       testDiary.DiaryID,
			Title:         "Draft Entry",
			Body:          "This is a draft",
			Status:        2, // StatusDraft
			InfoType:      1, // InfoTypeText
			AddUserID:     user1.UserId,
		}
		draftID, err := entryService.AddEntry(draftEntry)
		require.NoError(t, err)

		// 作者评论草稿 - 应该成功
		comment1 := &comment.Comment{
			EntryID:   draftID,
			Content:   "Comment by owner on draft",
			AddUserID: user1.UserId,
			AddUserIP: "127.0.0.1",
		}
		_, err = commentService.AddComment(comment1)
		assert.NoError(t, err, "Owner should be able to comment on draft")

		// 其他用户评论草稿 - 应该失败
		comment2 := &comment.Comment{
			EntryID:   draftID,
			Content:   "Comment by other on draft",
			AddUserID: user2.UserId,
			AddUserIP: "127.0.0.1",
		}
		_, err = commentService.AddComment(comment2)
		assert.Error(t, err, "Other user should not be able to comment on draft")
		// Check for EntryNotFound because GetEntryByID (active only) fails and GetEntryByIDForUser (owner only) fails
		// Using string check since errors.Is might not work with AppError details
		assert.Contains(t, err.Error(), "entry_not_found", "Should return entry_not_found error")
	})

	// 2. 测试软删除状态 (StatusDeleted = 0 or specific status)
	// Assuming logic allows commenting on soft deleted entries if owner
	t.Run("Soft Deleted Entry Comments", func(t *testing.T) {
		// 创建并删除日记
		trashEntry := &entry.Entry{
			DiaryID:       testDiary.DiaryID,
			Title:         "Trash Entry",
			Body:          "This is trash",
			Status:        1, // Start Active
			InfoType:      1,
			AddUserID:     user1.UserId,
		}
		trashID, err := entryService.AddEntry(trashEntry)
		require.NoError(t, err)

		// 软删除
		err = entryService.DeleteEntryByID(trashID, user1.UserId, time.Now().Unix())
		require.NoError(t, err)

		// Verify it is deleted (status != Active)
		// GetEntryByID should fail
		_, err = entryService.GetEntryByID(trashID)
		assert.Error(t, err)

		// 作者评论软删除日记 - 应该成功 (This is the fix verification)
		comment1 := &comment.Comment{
			EntryID:   trashID,
			Content:   "Comment by owner on trash",
			AddUserID: user1.UserId,
			AddUserIP: "127.0.0.1",
		}
		_, err = commentService.AddComment(comment1)
		assert.NoError(t, err, "Owner should be able to comment on soft deleted entry")

		// 列表评论也应该成功
		comments, _, err := commentService.ListCommentsByEntryID(trashID, 1, 10)
		assert.NoError(t, err, "Owner should be able to list comments of soft deleted entry")
		assert.Len(t, comments, 1)

		// 其他用户评论软删除日记 - 应该失败
		comment2 := &comment.Comment{
			EntryID:   trashID,
			Content:   "Comment by other on trash",
			AddUserID: user2.UserId,
			AddUserIP: "127.0.0.1",
		}
		_, err = commentService.AddComment(comment2)
		assert.Error(t, err, "Other user should not be able to comment on soft deleted entry")
	})
}
