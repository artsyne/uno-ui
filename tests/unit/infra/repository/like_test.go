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

func TestLikeRepository_BatchIsLiked(t *testing.T) {
	repo, mock := setupLikeTestRepo(t)
	userID := int64(1)
	ownerIDs := []int64{101, 102, 103}
	likeType := 1

	// 设置期望的查询 - 返回 101 和 103 被点赞
	rows := sqlmock.NewRows([]string{"owner_id"}).
		AddRow(101).
		AddRow(103)

	mock.ExpectQuery("SELECT owner_id FROM br_like WHERE").
		WithArgs(userID, likeType, ownerIDs[0], ownerIDs[1], ownerIDs[2]).
		WillReturnRows(rows)

	// 执行测试
	result, err := repo.BatchIsLiked(userID, ownerIDs, likeType)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.True(t, result[101], "Owner 101 should be liked")
	assert.False(t, result[102], "Owner 102 should not be liked")
	assert.True(t, result[103], "Owner 103 should be liked")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeRepository_BatchIsLiked_Empty(t *testing.T) {
	repo, _ := setupLikeTestRepo(t)

	// 测试空数组
	result, err := repo.BatchIsLiked(1, []int64{}, 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestLikeRepository_BatchIsLiked_NoLikes(t *testing.T) {
	repo, mock := setupLikeTestRepo(t)
	userID := int64(1)
	ownerIDs := []int64{101, 102, 103}
	likeType := 1

	// 设置期望的查询 - 没有任何点赞
	rows := sqlmock.NewRows([]string{"owner_id"})

	mock.ExpectQuery("SELECT owner_id FROM br_like WHERE").
		WithArgs(userID, likeType, ownerIDs[0], ownerIDs[1], ownerIDs[2]).
		WillReturnRows(rows)

	// 执行测试
	result, err := repo.BatchIsLiked(userID, ownerIDs, likeType)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.False(t, result[101])
	assert.False(t, result[102])
	assert.False(t, result[103])

	assert.NoError(t, mock.ExpectationsWereMet())
}
