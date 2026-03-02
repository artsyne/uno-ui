package entry_test

import (
	"testing"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryService_DraftsAndTrash(t *testing.T) {
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
	entryRepo := repository.NewEntryRepository()

	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
	entryService := entry.NewService(config, userService, taskService, diaryRepo, entryRepo, i18nManager)

	testUser, testDiary := setupTestUserAndDiary(t, ctx.(*framework.TestContext), userService, diaryService)

	t.Run("Test Recover Entry", func(t *testing.T) {
		// Create a deleted entry
		e := &entry.Entry{
			Title:     "Deleted Entry",
			Body:      "Content",
			DiaryID:   testDiary.DiaryID,
			AddUserID: testUser.UserId,
			Status:    entry.StatusDeleted,
		}
		id, err := entryService.AddEntry(e)
		require.NoError(t, err)

		saved, err := entryService.GetEntryByIDForUser(id, testUser.UserId)
		require.NoError(t, err)
		// If AddEntry sets status to 1 by default if 0?
		// Repo Insert uses the struct value.
		// But let's force update to be sure.
		saved.Status = entry.StatusDeleted
		saved.UpdateUTCTime++ // Ensure row is updated
		err = entryService.UpdateEntry(saved, testUser.UserId)
		require.NoError(t, err)

		// Get Diary Count
		d, _ := diaryService.GetDiaryByID(testDiary.DiaryID)
		initialCount := d.InfoNum

		// Recover
		err = entryService.RecoverEntryByID(id, testUser.UserId)
		require.NoError(t, err)

		// Verify Status
		recovered, err := entryService.GetEntryByIDForUser(id, testUser.UserId)
		require.NoError(t, err)
		assert.Equal(t, entry.StatusActive, recovered.Status)

		// Verify Count Incremented
		d, _ = diaryService.GetDiaryByID(testDiary.DiaryID)
		assert.Equal(t, initialCount+1, d.InfoNum)
	})

	t.Run("Test Hard Delete Entry", func(t *testing.T) {
		// Create an active entry
		e := &entry.Entry{
			Title:     "Active Entry",
			Body:      "Content",
			DiaryID:   testDiary.DiaryID,
			AddUserID: testUser.UserId,
			Status:    entry.StatusActive,
		}
		id, err := entryService.AddEntry(e)
		require.NoError(t, err)

		// Increment count manually if AddEntry didn't?
		// Controller calls Increment. Service AddEntry does NOT increment count.
		// So we must increment it manually to simulate controller behavior
		if err := diaryService.IncrementEntryCount(testDiary.DiaryID); err != nil {
			t.Fatal(err)
		}

		// Get Diary Count
		d, _ := diaryService.GetDiaryByID(testDiary.DiaryID)
		initialCount := d.InfoNum

		// Hard Delete
		err = entryService.HardDeleteEntryByID(id, testUser.UserId)
		require.NoError(t, err)

		// Verify Gone
		gone, err := entryService.GetEntryByIDForUser(id, testUser.UserId)
		assert.Error(t, err)
		assert.Nil(t, gone)

		// Verify Count Decremented
		d, _ = diaryService.GetDiaryByID(testDiary.DiaryID)
		assert.Equal(t, initialCount-1, d.InfoNum)
	})

	t.Run("Test List Drafts", func(t *testing.T) {
		// Create 2 drafts
		for i := 0; i < 2; i++ {
			e := &entry.Entry{
				Title:     fmt.Sprintf("Draft %d", i),
				Body:      "Content",
				DiaryID:   testDiary.DiaryID,
				AddUserID: testUser.UserId,
				Status:    entry.StatusDraft,
			}
			_, err := entryService.AddEntry(e)
			require.NoError(t, err)
		}

		// List Drafts
		drafts, hasMore, err := entryService.ListAllEntriesByUserID(testUser.UserId, entry.StatusDraft, 1, 10)
		require.NoError(t, err)
		// hasMore 表示是否有更多数据（返回记录数 > pageSize）
		// 只有当实际数据超过 pageSize 时才为 true
		_ = hasMore // 忽略此值，因为取决于测试数据量
		for _, d := range drafts {
			assert.Equal(t, entry.StatusDraft, d.Status)
		}
	})
}
