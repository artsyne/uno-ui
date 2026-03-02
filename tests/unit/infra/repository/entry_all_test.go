package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/stretchr/testify/assert"
)

func TestEntryRepository_ListAllByUserID(t *testing.T) {
	repo, mock := setupTestRepo(t)

	userID := int64(456)
	pageIndex := 1
	pageSize := 10
	status := entry.StatusDraft

	// Mock rows for Select - limit+1 模式，查询 pageSize+1 = 11 条
	rows := sqlmock.NewRows([]string{
		"entry_id", "diary_id", "is_top", "title", "summary", "body", 
		"view_num", "like_num", "comment_num", "add_user_id", "add_user_ip", 
		"add_utc_time", "update_utc_time", "status", "info_type",
		"handle", "display_name",
	}).AddRow(
		1, 100, 0, "Draft Title", "Summary", "Body", 
		0, 0, 0, userID, "127.0.0.1", 
		1234567890000, 1234567890000, status, 1,
		"userhandle", "User Name",
	)

	// limit+1 模式：参数为 pageSize+1 = 11
	mock.ExpectQuery(`SELECT e.*, COALESCE\(u.handle, ''\) as handle, COALESCE\(u.display_name, ''\) as display_name`).WithArgs(
		userID, status, pageSize+1, 0, // limit = pageSize + 1
	).WillReturnRows(rows)

	entries, hasMore, err := repo.ListAllByUserID(userID, status, pageIndex, pageSize)

	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, false, hasMore) // 只有1条记录，不超过 pageSize，所以 hasMore = false
	assert.Equal(t, "Draft Title", entries[0].Title)
}