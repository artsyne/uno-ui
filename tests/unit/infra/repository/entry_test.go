package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo 创建测试用的Repo和Mock
func setupTestRepo(t *testing.T) (*repository.EntryRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewEntryRepositoryWithDB(sqlxDB)
	return repo, mock
}

func TestEntryRepository_Insert(t *testing.T) {
	repo, mock := setupTestRepo(t)
	// defer repo.Close() // sqlmock db close if needed

	now := time.Now().Unix()
	testEntry := &entry.Entry{
		DiaryID:       123,
		IsTop:         0,
		Title:         "测试日记",
		Summary:       "摘要",
		Body:          "这是一篇测试日记的内容",
		ViewNum:       0,
		LikeNum:       0,
		CommentNum:    0,
		AddUserID:     456,
		AddUserIP:     "127.0.0.1",
		AddUTCTime:    now,
		UpdateUTCTime: now,
		Status:        entry.StatusActive,
		InfoType:      entry.InfoTypeText,
	}

	// 模拟成功插入
	// 注意：参数顺序必须与 repository.Insert 中的 SQL 参数顺序一致
	mock.ExpectExec(`INSERT INTO br_entry`).WithArgs(
		testEntry.DiaryID,
		testEntry.IsTop,
		testEntry.Title,
		testEntry.Summary,
		testEntry.Body,
		testEntry.ViewNum,
		testEntry.LikeNum,
		testEntry.CommentNum,
		testEntry.AddUserID,
		testEntry.AddUserIP,
		testEntry.AddUTCTime,
		testEntry.UpdateUTCTime,
		testEntry.Status,
		testEntry.InfoType,
	).WillReturnResult(sqlmock.NewResult(789, 1))

	id, err := repo.Insert(testEntry)

	assert.NoError(t, err)
	assert.Equal(t, int64(789), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_Update(t *testing.T) {
	repo, mock := setupTestRepo(t)

	now := time.Now().Unix()
	testEntry := &entry.Entry{
		EntryID:       789,
		DiaryID:       123,
		IsTop:         1,
		Title:         "更新的日记",
		Summary:       "更新摘要",
		Body:          "这是一篇更新的日记内容",
		ViewNum:       10,
		LikeNum:       5,
		CommentNum:    2,
		UpdateUTCTime: now,
		Status:        entry.StatusActive,
		InfoType:      entry.InfoTypeText,
	}

	// 模拟成功更新
	mock.ExpectExec(`UPDATE br_entry`).WithArgs(
		testEntry.DiaryID,
		testEntry.IsTop,
		testEntry.Title,
		testEntry.Summary,
		testEntry.Body,
		testEntry.ViewNum,
		testEntry.LikeNum,
		testEntry.CommentNum,
		testEntry.UpdateUTCTime,
		testEntry.Status,
		testEntry.InfoType,
		testEntry.EntryID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(testEntry)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_ListByUserID(t *testing.T) {
	repo, mock := setupTestRepo(t)

	userID := int64(456)
	diaryID := int64(123)
	pageIndex := 1
	pageSize := 10

	// 模拟查询列表
	rows := sqlmock.NewRows([]string{"entry_id", "title", "handle", "display_name", "diary_title", "social_setting"}).
		AddRow(1, "Title 1", "user1", "User One", "Diary Title", 1).
		AddRow(2, "Title 2", "user1", "User One", "Diary Title", 1)

	// limit+1 模式：pageSize=10, 查询 10+1=11 条
	mock.ExpectQuery(`SELECT e\.\*, COALESCE\(u\.handle, ''\) as handle, COALESCE\(u\.display_name, ''\) as display_name, d\.title as diary_title, d\.social_setting FROM br_entry e LEFT JOIN br_user u ON e\.add_user_id = u\.user_id JOIN br_diary d ON e\.diary_id = d\.diary_id WHERE e\.diary_id = \? AND e\.add_user_id = \? and e\.status = \? ORDER BY e\.entry_id DESC LIMIT \? offset \?`).WithArgs(
		diaryID, userID, entry.StatusActive, 11, 0, // pageSize+1=11
	).WillReturnRows(rows)

	// 不再需要 COUNT 查询
	entries, hasMore, err := repo.ListByUserID(userID, diaryID, entry.StatusActive, pageIndex, pageSize)

	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.False(t, hasMore) // 只有 2 条，返回 2 条，len=2 <= pageSize=10，所以 hasMore=false
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_ListByUserID_NoRows(t *testing.T) {
	repo, mock := setupTestRepo(t)

	// Mock empty result set (not error)
	rows := sqlmock.NewRows([]string{"entry_id", "title", "handle", "display_name"})
	mock.ExpectQuery(`SELECT e.*`).WillReturnRows(rows)

	entries, hasMore, err := repo.ListByUserID(1, 1, entry.StatusActive, 1, 10)

	assert.NoError(t, err)
	assert.Empty(t, entries)
	assert.False(t, hasMore)
}
