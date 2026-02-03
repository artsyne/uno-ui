package repository_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLikeTestRepo(t *testing.T) (*repository.LikeRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewLikeRepositoryWithDB(sqlxDB)
	return repo, mock
}

func TestLikeRepository_ToggleLike(t *testing.T) {
	repo, mock := setupLikeTestRepo(t)
	userID := int64(1)
	ownerID := int64(100)
	likeType := 1

	// Case 1: Like does not exist -> Insert
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_like`).WithArgs(userID, ownerID, likeType).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(`INSERT INTO br_like`).WithArgs(userID, ownerID, likeType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.ToggleLike(userID, ownerID, likeType)
	assert.NoError(t, err)

	// Case 2: Like exists -> Delete
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_like`).WithArgs(userID, ownerID, likeType).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectExec(`DELETE FROM br_like`).WithArgs(userID, ownerID, likeType).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.ToggleLike(userID, ownerID, likeType)
	assert.NoError(t, err)

	// Case 3: Error checking existence
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_like`).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.ToggleLike(userID, ownerID, likeType)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeRepository_IsLiked(t *testing.T) {
	repo, mock := setupLikeTestRepo(t)
	userID := int64(1)
	ownerID := int64(100)
	likeType := 1

	// Case 1: Liked
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_like`).WithArgs(userID, ownerID, likeType).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	liked, err := repo.IsLiked(userID, ownerID, likeType)
	assert.NoError(t, err)
	assert.True(t, liked)

	// Case 2: Not Liked
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_like`).WithArgs(userID, ownerID, likeType).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	liked, err = repo.IsLiked(userID, ownerID, likeType)
	assert.NoError(t, err)
	assert.False(t, liked)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeRepository_GetLikeCount(t *testing.T) {
	repo, mock := setupLikeTestRepo(t)
	ownerID := int64(100)
	likeType := 1

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM br_like`).WithArgs(ownerID, likeType).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	count, err := repo.GetLikeCount(ownerID, likeType)
	assert.NoError(t, err)
	assert.Equal(t, 10, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}
