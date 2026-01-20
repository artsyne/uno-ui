package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDiaryRepository 模拟日记本仓储实现
type MockDiaryRepository struct {
	db *sqlx.DB
}

// NewMockDiaryRepository 创建模拟日记本仓储实例
func NewMockDiaryRepository(t *testing.T) (*MockDiaryRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &MockDiaryRepository{db: sqlxDB}, mock
}

// 实现 diary.RepositoryInterface 接口的所有方法

func (r *MockDiaryRepository) Insert(diary *diary.Diary) (int64, error) {
	sqlQuery := `
		INSERT INTO br_diary (
			add_user_id,
			title,
			remark,
			info_num,
			add_utc_time,
			status
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(sqlQuery,
		diary.AddUserID,
		diary.Title,
		diary.Remark,
		diary.InfoNum,
		diary.AddUTCTime,
		diary.Status)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *MockDiaryRepository) Update(diary *diary.Diary) (int64, error) {
	sqlQuery := `
		UPDATE br_diary
		SET
			add_user_id = ?,
			title = ?,
			remark = ?,
			info_num = ?,
			add_utc_time = ?,
			status = ?
		WHERE diary_id = ?
	`

	res, err := r.db.Exec(sqlQuery,
		diary.AddUserID,
		diary.Title,
		diary.Remark,
		diary.InfoNum,
		diary.AddUTCTime,
		diary.Status,
		diary.DiaryID)

	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockDiaryRepository) Delete(diaryId int64) (int64, error) {
	sqlQuery := "UPDATE br_diary SET status = 0 WHERE diary_id = ?"
	res, err := r.db.Exec(sqlQuery, diaryId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockDiaryRepository) GetById(diaryId int64) (*diary.Diary, error) {
	sqlQuery := "SELECT * FROM br_diary WHERE diary_id = ?"
	diary := &diary.Diary{}
	err := r.db.Get(diary, sqlQuery, diaryId)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return diary, err
}

func (r *MockDiaryRepository) GetListByUserId(userId int64, offset, limit int) ([]*diary.Diary, error) {
	sqlQuery := "SELECT * FROM br_diary WHERE add_user_id = ? AND status = 1 ORDER BY add_utc_time DESC LIMIT ? OFFSET ?"
	var diaries []*diary.Diary
	err := r.db.Select(&diaries, sqlQuery, userId, limit, offset)
	return diaries, err
}

func (r *MockDiaryRepository) GetCountByUserId(userId int64) (int64, error) {
	var count int64
	sqlQuery := "SELECT COUNT(*) FROM br_diary WHERE add_user_id = ? AND status = 1"
	err := r.db.Get(&count, sqlQuery, userId)
	return count, err
}

func (r *MockDiaryRepository) IncrementEntryNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
	sqlQuery := `
		UPDATE br_diary SET info_num = info_num + 1,
		last_action_utc_time = ?,
		last_action_user_id=?, 
		last_action_user_name = ? 
		WHERE diary_id = ?
	`
	_, err := r.db.Exec(sqlQuery, actionUTCTime, userID, userName, diaryID)
	return err
}

func (r *MockDiaryRepository) DecrementEntryNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error {
	sqlQuery := `
		UPDATE br_diary SET info_num = info_num - 1,
		last_action_utc_time = ?,
		last_action_user_id=?, 
		last_action_user_name = ? 
		WHERE diary_id = ?
	`
	_, err := r.db.Exec(sqlQuery, actionUTCTime, userID, userName, diaryID)
	return err
}

// 测试用例

func TestDiaryRepository_Insert(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	testDiary := &diary.Diary{
		AddUserID:      123,
		Title:          "测试日记本",
		Remark:         "这是一个测试日记本",
		InfoNum:        0,
		AddUTCTime:     1640995200,
		Status:         1,
	}

	// 模拟成功插入
	mock.ExpectExec(`INSERT INTO br_diary`).WithArgs(
		testDiary.AddUserID,
		testDiary.Title,
		testDiary.Remark,
		testDiary.InfoNum,
		testDiary.AddUTCTime,
		testDiary.Status,
	).WillReturnResult(sqlmock.NewResult(456, 1))

	id, err := repo.Insert(testDiary)

	assert.NoError(t, err)
	assert.Equal(t, int64(456), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_Update(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	testDiary := &diary.Diary{
		DiaryID:        456,
		AddUserID:      123,
		Title:          "更新的日记本",
		Remark:         "这是一个更新的日记本",
		InfoNum:        5,
		AddUTCTime:     1640995200,
		Status:         1,
	}

	// 模拟成功更新
	mock.ExpectExec(`UPDATE br_diary`).WithArgs(
		testDiary.AddUserID,
		testDiary.Title,
		testDiary.Remark,
		testDiary.InfoNum,
		testDiary.AddUTCTime,
		testDiary.Status,
		testDiary.DiaryID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.Update(testDiary)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_GetById(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(456)

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"diary_id", "add_user_id", "title", "remark", "info_num", "add_utc_time", "status",
	}).AddRow(
		456, 123, "测试日记本", "这是一个测试日记本", 0, 1640995200, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_diary WHERE diary_id = \?`).WithArgs(diaryID).WillReturnRows(rows)

	result, err := repo.GetById(diaryID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, diaryID, result.DiaryID)
	assert.Equal(t, "测试日记本", result.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_GetById_NotFound(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(999)

	// 模拟没有找到日记本
	mock.ExpectQuery(`SELECT \* FROM br_diary WHERE diary_id = \?`).WithArgs(diaryID).WillReturnError(sql.ErrNoRows)

	result, err := repo.GetById(diaryID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_GetListByUserId(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	userID := int64(123)
	offset := 0
	limit := 10

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"diary_id", "add_user_id", "title", "remark", "info_num", "add_utc_time", "status",
	}).AddRow(
		456, 123, "日记本1", "描述1", 5, 1640995200, 1,
	).AddRow(
		457, 123, "日记本2", "描述2", 3, 1640995100, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_diary WHERE add_user_id = \? AND status = 1 ORDER BY add_utc_time DESC LIMIT \? OFFSET \?`).WithArgs(userID, limit, offset).WillReturnRows(rows)

	result, err := repo.GetListByUserId(userID, offset, limit)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "日记本1", result[0].Title)
	assert.Equal(t, "日记本2", result[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_GetCountByUserId(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	userID := int64(123)

	// 模拟计数结果
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_diary WHERE add_user_id = \? AND status = 1`).WithArgs(userID).WillReturnRows(rows)

	count, err := repo.GetCountByUserId(userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_IncrementEntryNum(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(456)
	userID := int64(123)
	userName := "testuser"
	actionUTCTime := int64(1640995200)

	// 模拟成功增加
	mock.ExpectExec(`UPDATE br_diary SET info_num = info_num \+ 1,\s*last_action_utc_time = \?,\s*last_action_user_id=\?,\s*last_action_user_name = \?\s*WHERE diary_id = \?`).WithArgs(actionUTCTime, userID, userName, diaryID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.IncrementEntryNum(diaryID, userID, userName, actionUTCTime)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_DecrementEntryNum(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(456)
	userID := int64(123)
	userName := "testuser"
	actionUTCTime := int64(1640995200)

	// 模拟成功减少
	mock.ExpectExec(`UPDATE br_diary SET info_num = info_num - 1,\s*last_action_utc_time = \?,\s*last_action_user_id=\?,\s*last_action_user_name = \?\s*WHERE diary_id = \?`).WithArgs(actionUTCTime, userID, userName, diaryID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DecrementEntryNum(diaryID, userID, userName, actionUTCTime)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_IncrementEntryNum_DatabaseError(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(456)
	userID := int64(123)
	userName := "testuser"
	actionUTCTime := int64(1640995200)

	// 模拟数据库错误
	mock.ExpectExec(`UPDATE br_diary SET info_num = info_num \+ 1,\s*last_action_utc_time = \?,\s*last_action_user_id=\?,\s*last_action_user_name = \?\s*WHERE diary_id = \?`).WithArgs(actionUTCTime, userID, userName, diaryID).WillReturnError(sql.ErrConnDone)

	err := repo.IncrementEntryNum(diaryID, userID, userName, actionUTCTime)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_DecrementEntryNum_DatabaseError(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(456)
	userID := int64(123)
	userName := "testuser"
	actionUTCTime := int64(1640995200)

	// 模拟数据库错误
	mock.ExpectExec(`UPDATE br_diary SET info_num = info_num - 1,\s*last_action_utc_time = \?,\s*last_action_user_id=\?,\s*last_action_user_name = \?\s*WHERE diary_id = \?`).WithArgs(actionUTCTime, userID, userName, diaryID).WillReturnError(sql.ErrConnDone)

	err := repo.DecrementEntryNum(diaryID, userID, userName, actionUTCTime)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_Delete(t *testing.T) {
	repo, mock := NewMockDiaryRepository(t)
	defer repo.db.Close()

	diaryID := int64(456)

	// 模拟成功删除（软删除，设置status=0）
	mock.ExpectExec(`UPDATE br_diary SET status = 0 WHERE diary_id = \?`).WithArgs(diaryID).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.Delete(diaryID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}