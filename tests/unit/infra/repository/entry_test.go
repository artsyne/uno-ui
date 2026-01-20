package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockEntryRepository 模拟日记仓储实现
type MockEntryRepository struct {
	db *sqlx.DB
}

// NewMockEntryRepository 创建模拟日记仓储实例
func NewMockEntryRepository(t *testing.T) (*MockEntryRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &MockEntryRepository{db: sqlxDB}, mock
}

// 实现 entry.RepositoryInterface 接口的所有方法

func (r *MockEntryRepository) Insert(entry *entry.Entry) (int64, error) {
	sqlQuery := `
		INSERT INTO br_entry (
			diary_id,
			title,
			body,
			add_user_id,
			add_utc_time,
			status
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(sqlQuery,
		entry.DiaryID,
		entry.Title,
		entry.Body,
		entry.AddUserID,
		entry.AddUTCTime,
		entry.Status)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *MockEntryRepository) Update(entry *entry.Entry) error {
	sqlQuery := `
		UPDATE br_entry
		SET
			diary_id = ?,
			title = ?,
			body = ?,
			add_user_id = ?,
			add_utc_time = ?,
			status = ?
		WHERE entry_id = ?
	`

	_, err := r.db.Exec(sqlQuery,
		entry.DiaryID,
		entry.Title,
		entry.Body,
		entry.AddUserID,
		entry.AddUTCTime,
		entry.Status,
		entry.EntryID)

	if err != nil {
		return err
	}

	return nil
}

func (r *MockEntryRepository) Delete(entryID int64) error {
	sqlQuery := "UPDATE br_entry SET status = 0 WHERE entry_id = ?"
	_, err := r.db.Exec(sqlQuery, entryID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MockEntryRepository) GetByID(entryID int64) (*entry.Entry, error) {
	sqlQuery := "SELECT * FROM br_entry WHERE entry_id = ?"
	entry := &entry.Entry{}
	err := r.db.Get(entry, sqlQuery, entryID)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return entry, err
}

func (r *MockEntryRepository) ListByUserID(userId int64, diaryID int64, pageIndex, pageSize int) ([]entry.Entry, int, error) {
	sqlQuery := "SELECT * FROM br_entry WHERE add_user_id = ? AND status = 1 ORDER BY add_utc_time DESC LIMIT ? OFFSET ?"
	offset := (pageIndex - 1) * pageSize
	var entries []entry.Entry
	err := r.db.Select(&entries, sqlQuery, userId, pageSize, offset)
	return entries, 0, err
}

func (r *MockEntryRepository) GetCountByUserID(userId int64, diaryID int64) (int, error) {
	var count int
	sqlQuery := "SELECT COUNT(*) FROM br_entry WHERE add_user_id = ? AND status = 1"
	err := r.db.Get(&count, sqlQuery, userId)
	return count, err
}

func (r *MockEntryRepository) UpdateByID(entry *entry.Entry) error {
    return r.Update(entry)
}

func (r *MockEntryRepository) DeleteByID(entryID int64) error {
    return r.Delete(entryID)
}

func (r *MockEntryRepository) GetPreviousEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
    return 0, nil
}

func (r *MockEntryRepository) GetNextEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
    return 0, nil
}

func (r *MockEntryRepository) GetStatsByIDs(entryIDs []int64) ([]entry.Entry, error) {
    return nil, nil
}

func (r *MockEntryRepository) IncrementLikeNum(entryID int64) error {
    return nil
}

func (r *MockEntryRepository) DecrementLikeNum(entryID int64) error {
    return nil
}

// 测试用例

func TestEntryRepository_Insert(t *testing.T) {
	repo, mock := NewMockEntryRepository(t)
	defer repo.db.Close()

	testEntry := &entry.Entry{
		DiaryID:    123,
		Title:      "测试日记",
		Body:       "这是一篇测试日记的内容",
		AddUserID:  456,
		AddUTCTime: 1640995200,
		Status:     1,
	}

	// 模拟成功插入
	mock.ExpectExec(`INSERT INTO br_entry`).WithArgs(
		testEntry.DiaryID,
		testEntry.Title,
		testEntry.Body,
		testEntry.AddUserID,
		testEntry.AddUTCTime,
		testEntry.Status,
	).WillReturnResult(sqlmock.NewResult(789, 1))

	id, err := repo.Insert(testEntry)

	assert.NoError(t, err)
	assert.Equal(t, int64(789), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_Update(t *testing.T) {
	repo, mock := NewMockEntryRepository(t)
	defer repo.db.Close()

	testEntry := &entry.Entry{
		EntryID:    789,
		DiaryID:    123,
		Title:      "更新的日记",
		Body:       "这是一篇更新的日记内容",
		AddUserID:  456,
		AddUTCTime: 1640995200,
		Status:     1,
	}

	// 模拟成功更新
	mock.ExpectExec(`UPDATE br_entry`).WithArgs(
		testEntry.DiaryID,
		testEntry.Title,
		testEntry.Body,
		testEntry.AddUserID,
		testEntry.AddUTCTime,
		testEntry.Status,
		testEntry.EntryID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(testEntry)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_GetByID(t *testing.T) {
	repo, mock := NewMockEntryRepository(t)
	defer repo.db.Close()

	entryID := int64(789)

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"entry_id", "diary_id", "title", "body", "add_user_id", "add_utc_time", "status",
	}).AddRow(
		789, 123, "测试日记", "这是一篇测试日记的内容", 456, 1640995200, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_entry WHERE entry_id = \?`).WithArgs(entryID).WillReturnRows(rows)

	result, err := repo.GetByID(entryID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entryID, result.EntryID)
	assert.Equal(t, "测试日记", result.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_GetByID_NotFound(t *testing.T) {
	repo, mock := NewMockEntryRepository(t)
	defer repo.db.Close()

	entryID := int64(999)

	// 模拟没有找到日记
	mock.ExpectQuery(`SELECT \* FROM br_entry WHERE entry_id = \?`).WithArgs(entryID).WillReturnError(sql.ErrNoRows)

	result, err := repo.GetByID(entryID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEntryRepository_Delete(t *testing.T) {
	repo, mock := NewMockEntryRepository(t)
	defer repo.db.Close()

	entryID := int64(789)

	// 模拟成功删除（软删除，设置status=0）
	mock.ExpectExec(`UPDATE br_entry SET status = 0 WHERE entry_id = \?`).WithArgs(entryID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(entryID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}