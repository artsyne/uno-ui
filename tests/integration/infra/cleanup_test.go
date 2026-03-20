package infra_test

import (
	"database/sql"
	"fmt"
	"testing"

	"log/slog"
	"os"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/tests/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataCleanup(t *testing.T) {
	// Enable debug logging
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	// Initialize framework with DB (Cleanup should be enabled by default now)
	ctx := framework.NewTestContext(framework.WithDatabase())
	
	// We defer cleanup to ensure it runs at end of test as well (idempotent check)
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// Get database connection via context
	db := ctx.GetDatabase().GetDB()

	// 1. Create and Insert Data
	// User
	testUser := ctx.GetDataGenerator().CreateUser()
	userId, err := insertUser(db.DB, testUser)
	require.NoError(t, err)
	testUser.UserId = userId
	// Track the created user manually since we inserted it manually (Framework's DataGenerator tracks it if using service, but here we use raw insert)
	// Actually, DataGenerator.CreateUser() calls isolator.TrackCreated? 
	// Let's check DataGenerator implementation.
	// If it just creates struct, it doesn't track ID because ID is 0.
	// So we need to manually track it if we want 'tracked' strategy to work.
	// But we are testing 'prefix' strategy primarily in Cleanup(strategy).
	// The CleanupManager.Cleanup("prefix") relies on prefix, so manual tracking is not strictly needed for prefix strategy.
	
	// Diary
	testDiary := ctx.GetDataGenerator().CreateDiary(userId)
	diaryId, err := insertDiary(db.DB, testDiary)
	require.NoError(t, err)
	testDiary.DiaryID = diaryId

	// Entry
	testEntry := ctx.GetDataGenerator().CreateEntry(userId, diaryId)
	entryId, err := insertEntry(db.DB, testEntry)
	require.NoError(t, err)
	testEntry.EntryID = entryId

	// Verify data exists
	assert.True(t, recordExists(db.DB, "br_user", "user_id", userId))
	assert.True(t, recordExists(db.DB, "br_diary", "diary_id", diaryId))
	assert.True(t, recordExists(db.DB, "br_entry", "entry_id", entryId))

	// 2. Perform Cleanup
	// Manually invoke cleanup to verify it deletes the data
	err = ctx.GetCleanup().Cleanup("prefix")
	require.NoError(t, err)

	// 3. Verify data is gone
	assert.False(t, recordExists(db.DB, "br_user", "user_id", userId), "User should be deleted")
	assert.False(t, recordExists(db.DB, "br_diary", "diary_id", diaryId), "Diary should be deleted")
	assert.False(t, recordExists(db.DB, "br_entry", "entry_id", entryId), "Entry should be deleted")
}

func recordExists(db *sql.DB, table, pk string, id int64) bool {
	// 使用白名单验证表名和主键名，防止 SQL 注入
	validTables := map[string]string{
		"br_user":    "user_id",
		"br_diary":   "diary_id",
		"br_entry":   "entry_id",
		"br_comment": "comment_id",
		"br_like":    "like_id",
	}

	expectedPK, ok := validTables[table]
	if !ok || expectedPK != pk {
		return false
	}

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, pk)
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func insertUser(db *sql.DB, u *user.User) (int64, error) {
	query := `INSERT INTO br_user (handle, email, pwd, display_name, time_zone, lang, status, created_utc_time, last_login_utc_time) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, u.Handle, u.Email, u.Pwd, u.DisplayName, u.TimeZone, u.Lang, u.Status, u.CreatedUTCTime, u.LastLoginUTCTime)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertDiary(db *sql.DB, d *diary.Diary) (int64, error) {
	query := `INSERT INTO br_diary (title, remark, add_user_id, add_utc_time, status, view_setting, social_setting, theme_id) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, d.Title, d.Remark, d.AddUserID, d.AddUTCTime, d.Status, d.ViewSetting, d.SocialSetting, d.ThemeID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func insertEntry(db *sql.DB, e *entry.Entry) (int64, error) {
	query := `INSERT INTO br_entry (diary_id, title, body, summary, info_type, add_user_id, add_utc_time, status) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, e.DiaryID, e.Title, e.Body, e.Summary, e.InfoType, e.AddUserID, e.AddUTCTime, e.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
