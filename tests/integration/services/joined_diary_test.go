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

func TestJoinedDiariesExclusion(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// Init services
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	config := ctx.GetConfig().GetAppConfig()
    
    // Explicitly set config for test stability
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

	// Create users
	me := createTestUser(t, ctx, userService, "me")
	other := createTestUser(t, ctx, userService, "other")

	t.Run("Joined Diaries Should Exclude Self-Created", func(t *testing.T) {
		// 1. Me creates 2 diaries
		diary1 := &diary.Diary{
			Title:      "My Diary 1",
			Remark:     "Remark 1",
			AddUserID:  me.UserId,
			AddUTCTime: time.Now().Unix(),
			Status:     1,
		}
		_, err := diaryService.AddDiary(diary1)
		require.NoError(t, err)

		diary2 := &diary.Diary{
			Title:      "My Diary 2",
			Remark:     "Remark 2",
			AddUserID:  me.UserId,
			AddUTCTime: time.Now().Unix(),
			Status:     1,
		}
		_, err = diaryService.AddDiary(diary2)
		require.NoError(t, err)

		// 2. Other creates a diary
		diaryOther := &diary.Diary{
			Title:      "Other's Diary",
			Remark:     "Remark Other",
			AddUserID:  other.UserId,
			AddUTCTime: time.Now().Unix(),
			Status:     1,
		}
		idOther, err := diaryService.AddDiary(diaryOther)
		require.NoError(t, err)

		// 3. Me joins Other's diary
		// First, Other invites Me
		err = diaryService.InviteAuthor(idOther, other.UserId, me.Email, "en")
		require.NoError(t, err)
		// Me accepts
		err = diaryService.AcceptInvitation(idOther, me.UserId)
		require.NoError(t, err)

		// 4. Verify ListJoinedDiariesByUserID for Me
		// Should ONLY contain Other's Diary
		joined, err := diaryService.ListJoinedDiariesByUserID(me.UserId)
		require.NoError(t, err)
		
		assert.Equal(t, 1, len(joined), "Should have exactly 1 joined diary")
		if len(joined) > 0 {
			assert.Equal(t, idOther, joined[0].DiaryID, "Should be Other's diary")
		}

		// 5. Verify Paged
		joinedPaged, count, err := diaryService.ListJoinedDiariesByUserIDPaged(me.UserId, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, 1, len(joinedPaged))
		assert.Equal(t, 1, count)
		if len(joinedPaged) > 0 {
			assert.Equal(t, idOther, joinedPaged[0].DiaryID)
		}
	})
}
