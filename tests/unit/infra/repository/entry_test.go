package repository

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
	rows := sqlmock.NewRows([]string{"entry_id", "title", "handle", "display_name"}).
		AddRow(1, "Title 1", "user1", "User One").
		AddRow(2, "Title 2", "user1", "User One")

	// 注意：ListByUserID 中有两个查询，一个是 Select，一个是 GetCountByUserID
	// 这里只模拟 Select，因为 GetCountByUserID 在 Select 之后调用
	// 但是 repo.ListByUserID 内部会调用 GetCountByUserID，所以也需要模拟它
	// 不过如果 Select 返回 empty，它可能不会调用 Count? 
	// 看代码：Select 成功后调用 GetCountByUserID.

	mock.ExpectQuery(`SELECT e.*, COALESCE\(u.handle, ''\) as handle, COALESCE\(u.display_name, ''\) as display_name FROM br_entry e`).WithArgs(
		diaryID, userID, entry.StatusActive, pageSize, 0,
	).WillReturnRows(rows)

	// 模拟 Count 查询
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_entry WHERE diary_id = \? AND add_user_id = \? and status = \?`).WithArgs(
		diaryID, userID, entry.StatusActive,
	).WillReturnRows(countRows)

	entries, total, err := repo.ListByUserID(userID, diaryID, entry.StatusActive, pageIndex, pageSize)

	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, 2, total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_ListByUserID_NoRows(t *testing.T) {
	repo, mock := setupTestRepo(t)

	// Mock empty result set (not error)
	rows := sqlmock.NewRows([]string{"entry_id", "title", "handle", "display_name"})
	mock.ExpectQuery(`SELECT e.*`).WillReturnRows(rows)

	entries, total, err := repo.ListByUserID(1, 1, entry.StatusActive, 1, 10)

	assert.NoError(t, err)
	assert.Empty(t, entries)
	assert.Equal(t, 0, total)
}
