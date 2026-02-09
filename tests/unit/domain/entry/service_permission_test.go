package entry_test

import (
	std_errors "errors"
	"fmt"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryService_Permissions(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	config := ctx.GetConfig().GetAppConfig()
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	taskService := task.NewService(config, taskRepo)
	i18nManager := i18n.NewLocalizationManager(config)
	userService := user.NewService(config, userRepo, taskService, i18nManager, nil)
	diaryRepo := repository.NewDiaryRepository()
	diaryService := diary.NewService(config, diaryRepo, userService, taskService, i18nManager)
	entryRepo := repository.NewEntryRepository()
	entryService := entry.NewService(config, userService, taskService, diaryRepo, entryRepo, i18nManager)

	// Helper to create user
	createUser := func(prefix string) *user.User {
		u := &user.User{
			Handle:      fmt.Sprintf("%s%d", prefix, time.Now().UnixNano()%1000000),
			Email:       ctx.GetDataGenerator().GenerateUniqueEmail(prefix),
			Pwd:         "password123",
			DisplayName: prefix,
			Status:      user.StatusActive,
		}
		err := userService.RegisterUser(u)
		require.NoError(t, err)
		return u
	}

	owner := createUser("owner")
	admin := createUser("admin")
	author := createUser("author")
	outsider := createUser("outsider")

	// Create Private Diary
	privateDiary := &diary.Diary{
		Title:         "Private Diary",
		ViewSetting:   2, // Private
		SocialSetting: 1,
		ThemeID:       1,
		AddUserID:     owner.UserId,
		Status:        1,
	}
	diaryID, err := diaryService.AddDiary(privateDiary)
	require.NoError(t, err)
	privateDiary.DiaryID = diaryID

	// Add members
	err = diaryService.InviteAuthor(diaryID, owner.UserId, admin.Email, "en")
	require.NoError(t, err)
	err = diaryService.AcceptInvitation(diaryID, admin.UserId)
	require.NoError(t, err)
	err = diaryService.UpdateAuthorAdmin(diaryID, owner.UserId, admin.UserId, diary.AuthorIsAdminYes)
	require.NoError(t, err)

	err = diaryService.InviteAuthor(diaryID, owner.UserId, author.Email, "en")
	require.NoError(t, err)
	err = diaryService.AcceptInvitation(diaryID, author.UserId)
	require.NoError(t, err)

	// Add Entries
	// Owner's entry
	ownerEntry := &entry.Entry{
		Title:     "Owner Entry",
		Body:      "Content",
		DiaryID:   diaryID,
		AddUserID: owner.UserId,
		Status:    entry.StatusActive,
		InfoType:  entry.InfoTypeText,
	}
	ownerEntryID, err := entryService.AddEntry(ownerEntry)
	require.NoError(t, err)
	_ = diaryService.IncrementEntryCount(diaryID)

	// Author's entry
	authorEntry := &entry.Entry{
		Title:     "Author Entry",
		Body:      "Content",
		DiaryID:   diaryID,
		AddUserID: author.UserId,
		Status:    entry.StatusActive,
		InfoType:  entry.InfoTypeText,
	}
	authorEntryID, err := entryService.AddEntry(authorEntry)
	require.NoError(t, err)
	_ = diaryService.IncrementEntryCount(diaryID)

	t.Run("CanUpdateEntry", func(t *testing.T) {
		// 1. Owner can update own entry
		can, err := entryService.CanUpdateEntry(ownerEntryID, owner.UserId)
		assert.NoError(t, err)
		assert.True(t, can, "Owner should update own entry")

		// 2. Owner can update author's entry (as owner of diary)
		can, err = entryService.CanUpdateEntry(authorEntryID, owner.UserId)
		assert.NoError(t, err)
		assert.True(t, can, "Owner should update author's entry")

		// 3. Author can update own entry
		can, err = entryService.CanUpdateEntry(authorEntryID, author.UserId)
		assert.NoError(t, err)
		assert.True(t, can, "Author should update own entry")

		// 4. Author cannot update owner's entry
		can, err = entryService.CanUpdateEntry(ownerEntryID, author.UserId)
		assert.NoError(t, err)
		assert.False(t, can, "Author should NOT update owner's entry")

		// 5. Admin can update author's entry
		can, err = entryService.CanUpdateEntry(authorEntryID, admin.UserId)
		assert.NoError(t, err)
		assert.True(t, can, "Admin should update author's entry")

		// 6. Admin can update owner's entry (as admin of diary)
		can, err = entryService.CanUpdateEntry(ownerEntryID, admin.UserId)
		assert.NoError(t, err)
		assert.True(t, can, "Admin should update owner's entry")

		// 7. Outsider cannot update
		can, err = entryService.CanUpdateEntry(ownerEntryID, outsider.UserId)
		assert.NoError(t, err)
		assert.False(t, can, "Outsider should NOT update")
	})

	t.Run("ListEntriesByDiaryID_Private", func(t *testing.T) {
		// 1. Owner can list
		entries, _, err := entryService.ListEntriesByDiaryID(diaryID, owner.UserId, entry.StatusActive, 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, entries)

		// 2. Author can list
		entries, _, err = entryService.ListEntriesByDiaryID(diaryID, author.UserId, entry.StatusActive, 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, entries)

		// 3. Admin can list
		entries, _, err = entryService.ListEntriesByDiaryID(diaryID, admin.UserId, entry.StatusActive, 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, entries)

		// 4. Outsider cannot list
		entries, _, err = entryService.ListEntriesByDiaryID(diaryID, outsider.UserId, entry.StatusActive, 1, 10)
		assert.Error(t, err)
		var appErr errors.AppError
		if std_errors.As(err, &appErr) {
			assert.Equal(t, errors.PermissionDenied.Code(), appErr.Code())
		}
		assert.Empty(t, entries)
	})

	t.Run("ListEntriesByDiaryID_Public", func(t *testing.T) {
		publicDiary := &diary.Diary{
			Title:         "Public Diary",
			ViewSetting:   1, // Public
			SocialSetting: 1,
			ThemeID:       1,
			AddUserID:     owner.UserId,
			Status:        1,
		}
		pubID, err := diaryService.AddDiary(publicDiary)
		require.NoError(t, err)

		// Add entry
		e := &entry.Entry{
			Title:     "Public Entry",
			Body:      "Content",
			DiaryID:   pubID,
			AddUserID: owner.UserId,
			Status:    entry.StatusActive,
			InfoType:  entry.InfoTypeText,
		}
		_, err = entryService.AddEntry(e)
		require.NoError(t, err)

		// Outsider can list
		entries, _, err := entryService.ListEntriesByDiaryID(pubID, outsider.UserId, entry.StatusActive, 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, entries)
	})

	t.Run("UpdateEntryByID_Permissions", func(t *testing.T) {
		// Helper to get full entry
		getFullEntry := func(id int64) *entry.Entry {
			e, err := entryService.GetEntryByID(id)
			require.NoError(t, err)
			return e
		}

		// 1. Author can update own entry
		toUpdate := getFullEntry(authorEntryID)
		toUpdate.Body = "Updated by Author"
		err := entryService.UpdateEntryByID(toUpdate, author.UserId)
		assert.NoError(t, err)

		// 2. Owner can update author's entry
		toUpdate = getFullEntry(authorEntryID)
		toUpdate.Body = "Updated by Owner"
		err = entryService.UpdateEntryByID(toUpdate, owner.UserId)
		assert.NoError(t, err)

		// 3. Admin can update author's entry
		toUpdate = getFullEntry(authorEntryID)
		toUpdate.Body = "Updated by Admin"
		err = entryService.UpdateEntryByID(toUpdate, admin.UserId)
		assert.NoError(t, err)

		// 4. Outsider cannot update
		toUpdate = getFullEntry(authorEntryID)
		toUpdate.Body = "Updated by Outsider"
		err = entryService.UpdateEntryByID(toUpdate, outsider.UserId)
		assert.Error(t, err)
		var appErr errors.AppError
		if std_errors.As(err, &appErr) {
			assert.Equal(t, errors.PermissionDenied.Code(), appErr.Code())
		}
	})

	t.Run("DeleteAndRecover_Permissions", func(t *testing.T) {
		// Create a temporary entry for testing delete/recover
		tempEntry := &entry.Entry{
			Title:     "Temp Entry",
			Body:      "Content",
			DiaryID:   diaryID,
			AddUserID: author.UserId,
			Status:    entry.StatusActive,
			InfoType:  entry.InfoTypeText,
		}
		tempID, err := entryService.AddEntry(tempEntry)
		require.NoError(t, err)
		
		// 1. Outsider cannot delete
		err = entryService.DeleteEntryByID(tempID, outsider.UserId, time.Now().UnixMilli())
		assert.Error(t, err)
		assert.True(t, std_errors.Is(err, errors.PermissionDenied))

		// 2. Author can delete own entry
		err = entryService.DeleteEntryByID(tempID, author.UserId, time.Now().UnixMilli())
		assert.NoError(t, err)

		// Verify deleted
		deleted, err := entryService.GetEntryByIDForUser(tempID, author.UserId)
		assert.NoError(t, err)
		assert.Equal(t, entry.StatusDeleted, deleted.Status)

		// 3. Outsider cannot recover
		err = entryService.RecoverEntryByID(tempID, outsider.UserId)
		assert.Error(t, err)
		assert.True(t, std_errors.Is(err, errors.PermissionDenied))

		// 4. Owner can recover (Author's entry)
		err = entryService.RecoverEntryByID(tempID, owner.UserId)
		assert.NoError(t, err)

		// Verify recovered
		recovered, err := entryService.GetEntryByIDForUser(tempID, author.UserId)
		assert.NoError(t, err)
		assert.Equal(t, entry.StatusActive, recovered.Status)

		// 5. Admin can delete
		err = entryService.DeleteEntryByID(tempID, admin.UserId, time.Now().UnixMilli())
		assert.NoError(t, err)

		// 6. Admin can recover
		err = entryService.RecoverEntryByID(tempID, admin.UserId)
		assert.NoError(t, err)
	})

	t.Run("GetEntryByIDForUser_Permissions", func(t *testing.T) {
		// Create a draft entry
		draftEntry := &entry.Entry{
			Title:     "Draft Entry",
			Body:      "Content",
			DiaryID:   diaryID,
			AddUserID: author.UserId,
			Status:    entry.StatusDraft,
			InfoType:  entry.InfoTypeText,
		}
		draftID, err := entryService.AddEntry(draftEntry)
		require.NoError(t, err)
		// No increment for draft

		// 1. Author can get draft
		got, err := entryService.GetEntryByIDForUser(draftID, author.UserId)
		assert.NoError(t, err)
		assert.Equal(t, draftID, got.EntryID)

		// 2. Outsider cannot get draft
		_, err = entryService.GetEntryByIDForUser(draftID, outsider.UserId)
		assert.Error(t, err)
		assert.True(t, std_errors.Is(err, errors.PermissionDenied), "Expected PermissionDenied, got %v", err)

		// 3. Owner can get draft (New Logic)
		got, err = entryService.GetEntryByIDForUser(draftID, owner.UserId)
		assert.NoError(t, err)
		assert.Equal(t, draftID, got.EntryID)

		// 4. Admin can get draft (New Logic)
		got, err = entryService.GetEntryByIDForUser(draftID, admin.UserId)
		assert.NoError(t, err)
		assert.Equal(t, draftID, got.EntryID)
	})

	t.Run("HardDelete_Permissions", func(t *testing.T) {
		// Create a temporary entry for testing hard delete
		tempEntry := &entry.Entry{
			Title:     "Temp Entry for Hard Delete",
			Body:      "Content",
			DiaryID:   diaryID,
			AddUserID: author.UserId,
			Status:    entry.StatusActive,
			InfoType:  entry.InfoTypeText,
		}
		tempID, err := entryService.AddEntry(tempEntry)
		require.NoError(t, err)
		_ = diaryService.IncrementEntryCount(diaryID)

		// 1. Outsider cannot hard delete
		err = entryService.HardDeleteEntryByID(tempID, outsider.UserId)
		assert.Error(t, err)
		assert.True(t, std_errors.Is(err, errors.PermissionDenied))

		// 2. Author can hard delete
		err = entryService.HardDeleteEntryByID(tempID, author.UserId)
		assert.NoError(t, err)

		// Verify physically deleted
		_, err = entryService.GetEntryByID(tempID)
		assert.Error(t, err) // Should be not found or error
		
		// Create another one for Owner/Admin test
		tempEntry2 := &entry.Entry{
			Title:     "Temp Entry 2 for Hard Delete",
			Body:      "Content",
			DiaryID:   diaryID,
			AddUserID: author.UserId,
			Status:    entry.StatusActive,
			InfoType:  entry.InfoTypeText,
		}
		tempID2, err := entryService.AddEntry(tempEntry2)
		require.NoError(t, err)
		_ = diaryService.IncrementEntryCount(diaryID)

		// 3. Owner can hard delete
		err = entryService.HardDeleteEntryByID(tempID2, owner.UserId)
		assert.NoError(t, err)
	})
}
