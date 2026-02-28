package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiaryAuthorPagination(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// Initialize Services
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()

	// Ensure config settings
	config.App.UserJoinedDiariesLimit = 50
	config.App.DefaultLang = "zh-cn"

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
	redisClient.FlushDB()

	userService := user.NewService(config, userRepo, taskService, i18nManager, redisClient)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// Setup Data: 1 Owner, 5 Active Authors, 3 Pending Authors
	// Use short prefixes to avoid handle truncation collision (max 16 chars)
	owner := createTestUser(t, ctx, userService, "own")
	
	d := &diary.Diary{
		Title:      "Pagination Test Diary",
		AddUserID:  owner.UserId,
		AddUTCTime: time.Now().Unix(),
		Status:     1,
	}
	diaryID, err := diaryService.AddDiary(d)
	require.NoError(t, err)

	// Add 5 Active Authors
	for i := 0; i < 5; i++ {
		u := createTestUser(t, ctx, userService, fmt.Sprintf("a%d", i))
		require.NoError(t, diaryService.InviteAuthor(diaryID, owner.UserId, u.Email, "en"))
		require.NoError(t, diaryService.AcceptInvitation(diaryID, u.UserId))
	}

	// Add 3 Pending Authors
	for i := 0; i < 3; i++ {
		u := createTestUser(t, ctx, userService, fmt.Sprintf("p%d", i))
		require.NoError(t, diaryService.InviteAuthor(diaryID, owner.UserId, u.Email, "en"))
		// Do not accept invitation -> Pending status
	}

	// Total Authors: 1 Owner (Active) + 5 Active + 3 Pending = 9 Total (6 Active, 3 Pending)

	t.Run("Owner should see all authors with pagination", func(t *testing.T) {
		// Page 1, Size 5 -> Should get 5 authors
		authors, total, err := diaryService.GetDiaryAuthors(diaryID, false, 1, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(9), total)
		assert.Len(t, authors, 5)

		// Page 2, Size 5 -> Should get remaining 4 authors
		authors, total, err = diaryService.GetDiaryAuthors(diaryID, false, 2, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(9), total)
		assert.Len(t, authors, 4)
	})

	t.Run("Non-Owner should see only active authors with pagination", func(t *testing.T) {
		// Page 1, Size 5 -> Should get 5 active authors (Owner + 4 others)
		authors, total, err := diaryService.GetDiaryAuthors(diaryID, true, 1, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(6), total) // 1 Owner + 5 Active
		assert.Len(t, authors, 5)

		// Page 2, Size 5 -> Should get remaining 1 active author
		authors, total, err = diaryService.GetDiaryAuthors(diaryID, true, 2, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(6), total)
		assert.Len(t, authors, 1)

		// Verify that none of the authors are pending
		for _, a := range authors {
			assert.Equal(t, int8(1), a.Status) // 1 is Active
		}
	})

	t.Run("Invalid Page Index should default to 1", func(t *testing.T) {
		authors, total, err := diaryService.GetDiaryAuthors(diaryID, true, -1, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(6), total)
		assert.Len(t, authors, 5)
	})
}
