package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinedDiariesVisibilityRepro(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	config := ctx.GetConfig().GetAppConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
	redisClient.FlushDB(context.Background())

	userService := user.NewService(config, userRepo, taskService, i18nManager, redisClient)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)

	// Helper to create user
	createUser := func(email, handle string) *user.User {
		u := &user.User{
			Email:          email,
			Handle:         handle,
			DisplayName:    handle,
			Pwd:            "Test@123456",
			Status:         1,
			TimeZone:       "UTC",
			Lang:           "en",
			CreatedUTCTime: time.Now().Unix(),
		}
		err := userService.RegisterUser(u)
		require.NoError(t, err)
		return u
	}

	// 1. Create Users
	ts := time.Now().Unix() % 10000
	alice := createUser(fmt.Sprintf("alice_%d@example.com", ts), fmt.Sprintf("a_%d", ts))         // Profile Owner
	charlie := createUser(fmt.Sprintf("charlie_%d@example.com", ts), fmt.Sprintf("c_%d", ts)) // Diary Creator

	// 2. Charlie creates Public Diary
	publicDiary := &diary.Diary{
		Title:         "Charlie Public Diary",
		AddUserID:     charlie.UserId,
		ViewSetting:   diary.ViewSettingPublic, // 1
		Status:        diary.StatusActive,      // 1
		AddUTCTime:    time.Now().Unix(),
		SocialSetting: diary.SocialSettingEnabled,
	}
	_, err := diaryService.AddDiary(publicDiary)
	require.NoError(t, err)

	// 3. Charlie creates Private Diary
	privateDiary := &diary.Diary{
		Title:         "Charlie Private Diary",
		AddUserID:     charlie.UserId,
		ViewSetting:   diary.ViewSettingPrivate, // 2
		Status:        diary.StatusActive,       // 1
		AddUTCTime:    time.Now().Unix(),
		SocialSetting: diary.SocialSettingEnabled,
	}
	_, err = diaryService.AddDiary(privateDiary)
	require.NoError(t, err)

	// 4. Alice joins Public Diary
	err = diaryService.InviteAuthor(publicDiary.DiaryID, charlie.UserId, alice.Email, "en")
	require.NoError(t, err)
	err = diaryService.AcceptInvitation(publicDiary.DiaryID, alice.UserId)
	require.NoError(t, err)

	// 5. Alice joins Private Diary
	err = diaryService.InviteAuthor(privateDiary.DiaryID, charlie.UserId, alice.Email, "en")
	require.NoError(t, err)
	err = diaryService.AcceptInvitation(privateDiary.DiaryID, alice.UserId)
	require.NoError(t, err)

	// 6. Verify Alice (Owner) sees BOTH
	// includePrivate = true
	listOwner, err := diaryService.ListJoinedDiariesByUserIDWithLimit(alice.UserId, true, 10)
	require.NoError(t, err)
	assert.Len(t, listOwner, 2, "Owner should see both joined diaries")

	// 7. Verify Bob (Visitor) sees ONLY Public
	// includePrivate = false
	listVisitor, err := diaryService.ListJoinedDiariesByUserIDWithLimit(alice.UserId, false, 10)
	require.NoError(t, err)
	
	if assert.Len(t, listVisitor, 1, "Visitor should see exactly 1 joined diary") {
		assert.Equal(t, publicDiary.DiaryID, listVisitor[0].DiaryID, "Visitor should see the Public diary")
		assert.Equal(t, diary.ViewSettingPublic, listVisitor[0].ViewSetting)
	} else {
		// Debug info
		for _, d := range listVisitor {
			t.Logf("Visible Diary: ID=%d, ViewSetting=%d, Title=%s", d.DiaryID, d.ViewSetting, d.Title)
		}
	}

	// 8. Verify Pagination (ListJoinedDiariesByUserIDPaged)
	// Visitor Paged
	pagedList, hasMore, err := diaryService.ListJoinedDiariesByUserIDPaged(alice.UserId, false, 1, 10)
	require.NoError(t, err)
	assert.False(t, hasMore, "Paged should have no more")
	assert.Len(t, pagedList, 1, "Paged List should be 1")
	if len(pagedList) > 0 {
		assert.Equal(t, publicDiary.DiaryID, pagedList[0].DiaryID)
	}

	// Owner Paged
	pagedListOwner, hasMoreOwner, err := diaryService.ListJoinedDiariesByUserIDPaged(alice.UserId, true, 1, 10)
	require.NoError(t, err)
	assert.False(t, hasMoreOwner, "Owner Paged should have no more")
	assert.Len(t, pagedListOwner, 2, "Owner Paged List should be 2")
}
