package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCommentTestRepo(t *testing.T) (*repository.CommentRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewCommentRepositoryWithDB(sqlxDB)
	return repo, mock
}

func TestCommentRepository_Insert(t *testing.T) {
	repo, mock := setupCommentTestRepo(t)

	now := time.Now().Unix()
	c := &comment.Comment{
		EntryID:    1,
		Content:    "Test Comment",
		AddUtcTime: now,
		AddUserID:  100,
		AddUserIP:  "127.0.0.1",
	}

	mock.ExpectExec(`INSERT INTO br_comment`).WithArgs(
		c.EntryID, c.Content, c.AddUtcTime, c.AddUserID, c.AddUserIP,
	).WillReturnResult(sqlmock.NewResult(10, 1))

	id, err := repo.Insert(c)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_Delete(t *testing.T) {
	repo, mock := setupCommentTestRepo(t)
	commentID := int64(10)

	// Case 1: Success
	mock.ExpectExec(`DELETE FROM br_comment WHERE comment_id = \?`).WithArgs(commentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(commentID)
	assert.NoError(t, err)

	// Case 2: Not Found
	mock.ExpectExec(`DELETE FROM br_comment WHERE comment_id = \?`).WithArgs(commentID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(commentID)
	assert.Error(t, err)
	// assert.IsType(t, &errors.AppError{}, err) // Can check specific error type if needed

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_GetByID(t *testing.T) {
	repo, mock := setupCommentTestRepo(t)
	commentID := int64(10)

	rows := sqlmock.NewRows([]string{"comment_id", "entry_id", "content", "add_utc_time", "add_user_id", "add_user_ip", "handle", "display_name"}).
		AddRow(commentID, 1, "Content", 1234567890, 100, "127.0.0.1", "user1", "User One")

	mock.ExpectQuery(`SELECT .* FROM br_comment c .* WHERE c.comment_id = \?`).WithArgs(commentID).
		WillReturnRows(rows)

	c, err := repo.GetByID(commentID)
	assert.NoError(t, err)
	assert.Equal(t, commentID, c.CommentID)
	assert.Equal(t, "user1", c.AddUserHandle)

	// Case: Not Found
	mock.ExpectQuery(`SELECT .* FROM br_comment c .* WHERE c.comment_id = \?`).WithArgs(commentID).
		WillReturnError(sql.ErrNoRows)

	c, err = repo.GetByID(commentID)
	assert.Error(t, err)
	assert.Nil(t, c)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentRepository_ListByEntryID(t *testing.T) {
	repo, mock := setupCommentTestRepo(t)
	entryID := int64(1)

	rows := sqlmock.NewRows([]string{"comment_id", "entry_id", "content", "add_utc_time", "add_user_id", "add_user_ip", "handle", "display_name"}).
		AddRow(10, entryID, "Content 1", 1234567890, 100, "127.0.0.1", "user1", "User One").
		AddRow(11, entryID, "Content 2", 1234567891, 101, "127.0.0.1", "user2", "User Two")

	mock.ExpectQuery(`SELECT .* FROM br_comment c .* WHERE c.entry_id = \? ORDER BY c.comment_id DESC LIMIT \? OFFSET \?`).
		WithArgs(entryID, 10, 0).
		WillReturnRows(rows)

	// Expect Count Query
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_comment WHERE entry_id = \?`).
		WithArgs(entryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	list, count, err := repo.ListByEntryID(entryID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, 2, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}
