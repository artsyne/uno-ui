package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPostRepository 模拟日记仓储实现
type MockPostRepository struct {
	db *sqlx.DB
}

// NewMockPostRepository 创建模拟日记仓储实例
func NewMockPostRepository(t *testing.T) (*MockPostRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &MockPostRepository{db: sqlxDB}, mock
}

// 实现 post.RepositoryInterface 接口的所有方法

func (r *MockPostRepository) Insert(post *post.Post) (int64, error) {
	sqlQuery := `
		INSERT INTO br_post (
			diary_id,
			title,
			body,
			add_user_id,
			add_utc_time,
			status
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(sqlQuery,
		post.DiaryID,
		post.Title,
		post.Body,
		post.AddUserID,
		post.AddUTCTime,
		post.Status)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *MockPostRepository) Update(post *post.Post) (int64, error) {
	sqlQuery := `
		UPDATE br_post
		SET
			diary_id = ?,
			title = ?,
			body = ?,
			add_user_id = ?,
			add_utc_time = ?,
			status = ?
		WHERE post_id = ?
	`

	res, err := r.db.Exec(sqlQuery,
		post.DiaryID,
		post.Title,
		post.Body,
		post.AddUserID,
		post.AddUTCTime,
		post.Status,
		post.PostID)

	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockPostRepository) Delete(postId int64) (int64, error) {
	sqlQuery := "UPDATE br_post SET status = 0 WHERE post_id = ?"
	res, err := r.db.Exec(sqlQuery, postId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockPostRepository) GetByID(postId int64) (*post.Post, error) {
	sqlQuery := "SELECT * FROM br_post WHERE post_id = ?"
	post := &post.Post{}
	err := r.db.Get(post, sqlQuery, postId)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return post, err
}

func (r *MockPostRepository) ListByUserID(userId int64, offset, limit int) ([]*post.Post, error) {
	sqlQuery := "SELECT * FROM br_post WHERE add_user_id = ? AND status = 1 ORDER BY add_utc_time DESC LIMIT ? OFFSET ?"
	var posts []*post.Post
	err := r.db.Select(&posts, sqlQuery, userId, limit, offset)
	return posts, err
}

func (r *MockPostRepository) ListByDiaryID(diaryId int64, offset, limit int) ([]*post.Post, error) {
	sqlQuery := "SELECT * FROM br_post WHERE diary_id = ? AND status = 1 ORDER BY add_utc_time DESC LIMIT ? OFFSET ?"
	var posts []*post.Post
	err := r.db.Select(&posts, sqlQuery, diaryId, limit, offset)
	return posts, err
}

func (r *MockPostRepository) GetCountByUserID(userId int64) (int64, error) {
	var count int64
	sqlQuery := "SELECT COUNT(*) FROM br_post WHERE add_user_id = ? AND status = 1"
	err := r.db.Get(&count, sqlQuery, userId)
	return count, err
}

func (r *MockPostRepository) GetCountByDiaryID(diaryId int64) (int64, error) {
	var count int64
	sqlQuery := "SELECT COUNT(*) FROM br_post WHERE diary_id = ? AND status = 1"
	err := r.db.Get(&count, sqlQuery, diaryId)
	return count, err
}

// 测试用例

func TestPostRepository_Insert(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	testPost := &post.Post{
		DiaryID:    123,
		Title:      "测试日记",
		Body:       "这是一篇测试日记的内容",
		AddUserID:  456,
		AddUTCTime: 1640995200,
		Status:     1,
	}

	// 模拟成功插入
	mock.ExpectExec(`INSERT INTO br_post`).WithArgs(
		testPost.DiaryID,
		testPost.Title,
		testPost.Body,
		testPost.AddUserID,
		testPost.AddUTCTime,
		testPost.Status,
	).WillReturnResult(sqlmock.NewResult(789, 1))

	id, err := repo.Insert(testPost)

	assert.NoError(t, err)
	assert.Equal(t, int64(789), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Update(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	testPost := &post.Post{
		PostID:     789,
		DiaryID:    123,
		Title:      "更新的日记",
		Body:       "这是一篇更新的日记内容",
		AddUserID:  456,
		AddUTCTime: 1640995200,
		Status:     1,
	}

	// 模拟成功更新
	mock.ExpectExec(`UPDATE br_post`).WithArgs(
		testPost.DiaryID,
		testPost.Title,
		testPost.Body,
		testPost.AddUserID,
		testPost.AddUTCTime,
		testPost.Status,
		testPost.PostID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.Update(testPost)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByID(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	postID := int64(789)

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"post_id", "diary_id", "title", "body", "add_user_id", "add_utc_time", "status",
	}).AddRow(
		789, 123, "测试日记", "这是一篇测试日记的内容", 456, 1640995200, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_post WHERE post_id = \?`).WithArgs(postID).WillReturnRows(rows)

	result, err := repo.GetByID(postID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.PostID)
	assert.Equal(t, "测试日记", result.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetByID_NotFound(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	postID := int64(999)

	// 模拟没有找到日记
	mock.ExpectQuery(`SELECT \* FROM br_post WHERE post_id = \?`).WithArgs(postID).WillReturnError(sql.ErrNoRows)

	result, err := repo.GetByID(postID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_ListByUserID(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	userID := int64(456)
	offset := 0
	limit := 10

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"post_id", "diary_id", "title", "body", "add_user_id", "add_utc_time", "status",
	}).AddRow(
		789, 123, "日记1", "内容1", 456, 1640995200, 1,
	).AddRow(
		790, 124, "日记2", "内容2", 456, 1640995100, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_post WHERE add_user_id = \? AND status = 1 ORDER BY add_utc_time DESC LIMIT \? OFFSET \?`).WithArgs(userID, limit, offset).WillReturnRows(rows)

	result, err := repo.ListByUserID(userID, offset, limit)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "日记1", result[0].Title)
	assert.Equal(t, "日记2", result[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_ListByDiaryID(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	diaryID := int64(123)
	offset := 0
	limit := 10

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"post_id", "diary_id", "title", "body", "add_user_id", "add_utc_time", "status",
	}).AddRow(
		789, 123, "日记1", "内容1", 456, 1640995200, 1,
	).AddRow(
		791, 123, "日记3", "内容3", 457, 1640995000, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_post WHERE diary_id = \? AND status = 1 ORDER BY add_utc_time DESC LIMIT \? OFFSET \?`).WithArgs(diaryID, limit, offset).WillReturnRows(rows)

	result, err := repo.ListByDiaryID(diaryID, offset, limit)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "日记1", result[0].Title)
	assert.Equal(t, "日记3", result[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetCountByUserID(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	userID := int64(456)

	// 模拟计数结果
	rows := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_post WHERE add_user_id = \? AND status = 1`).WithArgs(userID).WillReturnRows(rows)

	count, err := repo.GetCountByUserID(userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetCountByDiaryID(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	diaryID := int64(123)

	// 模拟计数结果
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_post WHERE diary_id = \? AND status = 1`).WithArgs(diaryID).WillReturnRows(rows)

	count, err := repo.GetCountByDiaryID(diaryID)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_Delete(t *testing.T) {
	repo, mock := NewMockPostRepository(t)
	defer repo.db.Close()

	postID := int64(789)

	// 模拟成功删除（软删除，设置status=0）
	mock.ExpectExec(`UPDATE br_post SET status = 0 WHERE post_id = \?`).WithArgs(postID).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.Delete(postID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}