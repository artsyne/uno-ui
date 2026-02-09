package repository_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDiaryRepository(t *testing.T) (diary.RepositoryInterface, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewDiaryRepositoryWithDB(sqlxDB)
	return repo, mock
}

func TestDiaryRepository_Insert(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	testDiary := &diary.Diary{
		Title:         "测试日记本",
		Remark:        "这是一个测试日记本",
		AccessURL:     "test-url",
		CoverURL:      "cover.jpg",
		CoverMidURL:   "cover_mid.jpg",
		LikeNum:       0,
		MemberNum:     1,
		ViewNum:       0,
		ViewSetting:   1,
		AccessPwd:     "",
		SocialSetting: 1,
		AddUTCTime:    1640995200,
		AddUserID:     123,
		Status:        1,
		InfoNum:       0,
		ThemeID:       1,
	}

	mock.ExpectExec(`INSERT INTO br_diary`).WithArgs(
		testDiary.Title, testDiary.Remark, testDiary.AccessURL,
		testDiary.CoverURL, testDiary.CoverMidURL, testDiary.LikeNum,
		testDiary.MemberNum, testDiary.ViewNum,
		testDiary.ViewSetting, testDiary.AccessPwd, testDiary.SocialSetting,
		testDiary.AddUTCTime,
		testDiary.AddUserID, testDiary.Status, testDiary.InfoNum,
		testDiary.ThemeID,
	).WillReturnResult(sqlmock.NewResult(456, 1))

	id, err := repo.Insert(testDiary)

	assert.NoError(t, err)
	assert.Equal(t, int64(456), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_UpdateAuthorStatus(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(1001)
	userID := int64(2001)
	status := int8(2)

	mock.ExpectExec(`UPDATE br_diary_author SET status = \? WHERE diary_id = \? AND user_id = \?`).
		WithArgs(status, diaryID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateAuthorStatus(diaryID, userID, status)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_UpdateAuthorAdmin(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(1001)
	userID := int64(2001)
	isAdmin := int8(1)

	mock.ExpectExec(`UPDATE br_diary_author SET is_admin = \? WHERE diary_id = \? AND user_id = \?`).
		WithArgs(isAdmin, diaryID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateAuthorAdmin(diaryID, userID, isAdmin)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_DeleteAuthor(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(1001)
	userID := int64(2001)

	mock.ExpectExec(`DELETE FROM br_diary_author WHERE diary_id = \? AND user_id = \?`).
		WithArgs(diaryID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteAuthor(diaryID, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_GetByDiaryID(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(456)
	rows := sqlmock.NewRows([]string{
		"diary_id", "title", "remark", "add_user_id", "status", "handle", "display_name",
	}).AddRow(
		diaryID, "测试日记本", "备注", 123, 1, "testuser", "Test User",
	)

	// SQL select includes join with br_user
	mock.ExpectQuery(`SELECT d\.\*, u1\.handle, u1\.display_name FROM br_diary d`).
		WithArgs(diaryID, 1). // 1 is StatusActive
		WillReturnRows(rows)

	result, err := repo.GetByDiaryID(diaryID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, diaryID, result.DiaryID)
	assert.Equal(t, "测试日记本", result.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_InsertAuthor(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	author := &diary.Author{
		DiaryID:     1001,
		UserID:      2001,
		IsAdmin:     1,
		JoinUTCTime: 1640995200,
		PostNum:     0,
		Status:      1,
	}

	mock.ExpectExec(`INSERT INTO br_diary_author`).WithArgs(
		author.DiaryID, author.UserID, author.IsAdmin,
		author.JoinUTCTime, author.PostNum, author.Status,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertAuthor(author)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_GetDiaryAuthor(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(1001)
	userID := int64(2001)

	rows := sqlmock.NewRows([]string{
		"diary_id", "user_id", "is_admin", "status",
	}).AddRow(
		diaryID, userID, 1, 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_diary_author WHERE diary_id = \? AND user_id = \?`).
		WithArgs(diaryID, userID).
		WillReturnRows(rows)

	author, err := repo.GetDiaryAuthor(diaryID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, author)
	assert.Equal(t, diaryID, author.DiaryID)
	assert.Equal(t, userID, author.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_ListAuthorsByDiaryID(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(1001)

	rows := sqlmock.NewRows([]string{
		"diary_id", "user_id", "is_admin", "status", "handle", "display_name", "face_url",
	}).AddRow(
		diaryID, 2001, 1, 1, "owner", "Owner", "face1.jpg",
	).AddRow(
		diaryID, 2002, 0, 1, "member", "Member", "face2.jpg",
	)

	mock.ExpectQuery(`SELECT da\.\*, u\.handle, u\.display_name, u\.face_url FROM br_diary_author da`).
		WithArgs(diaryID).
		WillReturnRows(rows)

	authors, err := repo.ListAuthorsByDiaryID(diaryID)

	assert.NoError(t, err)
	assert.Len(t, authors, 2)
	assert.Equal(t, int64(2001), authors[0].UserID)
	assert.Equal(t, int64(2002), authors[1].UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDiaryRepository_IsDiaryOwner(t *testing.T) {
	repo, mock := newTestDiaryRepository(t)

	diaryID := int64(1001)
	userID := int64(2001)

	// Case 1: Is Owner
	rows := sqlmock.NewRows([]string{"add_user_id"}).AddRow(userID)
	mock.ExpectQuery(`SELECT add_user_id FROM br_diary`).
		WithArgs(diaryID, 1). // 1 is StatusActive
		WillReturnRows(rows)

	isOwner, err := repo.IsDiaryOwner(diaryID, userID)
	assert.NoError(t, err)
	assert.True(t, isOwner)

	// Case 2: Not Owner
	rows2 := sqlmock.NewRows([]string{"add_user_id"}).AddRow(userID + 1)
	mock.ExpectQuery(`SELECT add_user_id FROM br_diary`).
		WithArgs(diaryID, 1).
		WillReturnRows(rows2)

	isOwner, err = repo.IsDiaryOwner(diaryID, userID)
	assert.NoError(t, err)
	assert.False(t, isOwner)

	// Case 3: Diary Not Found (or deleted)
	mock.ExpectQuery(`SELECT add_user_id FROM br_diary`).
		WithArgs(diaryID, 1).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.IsDiaryOwner(diaryID, userID)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
