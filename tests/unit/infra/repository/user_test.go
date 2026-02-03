package repository_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockUserRepository 模拟用户仓储实现
type MockUserRepository struct {
	db *sqlx.DB
}

// NewMockUserRepository 创建模拟用户仓储实例
func NewMockUserRepository(t *testing.T) (*MockUserRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &MockUserRepository{db: sqlxDB}, mock
}

// 实现 user.RepositoryInterface 接口的所有方法

func (r *MockUserRepository) InsertUser(user *user.User) (int64, error) {
	sql := `
                insert into br_user  (
                    handle,
                    display_name,
                    user_group,
                    pwd,
                    profile_url,
                    face_url,
                    email,
                    email_verified,
                    reg_ip,
                    created_utc_time,
                    remark,
                    last_ip,
                    last_login_utc_time,
                    time_zone,
                    lang,
                    status
                )
                values (
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?
                )
                `

	res, err := r.db.Exec(sql,
		user.Handle,
		user.DisplayName,
		user.UserGroup,
		user.Pwd,
		user.ProfileUrl,
		user.FaceUrl,
		user.Email,
		user.EmailVerified,
		user.RegIp,
		user.CreatedUTCTime,
		user.Remark,
		user.LastIp,
		user.LastLoginUTCTime,
		user.TimeZone,
		user.Lang,
		user.Status)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *MockUserRepository) UpdateUser(user *user.User) (int64, error) {
	sql := `
                update br_user
                 set
                        user_id = ?,
                        handle = ?,
                        display_name = ?,
                        user_group = ?,
                        pwd = ?,
                        profile_url = ?,
                        face_url = ?,
                        email = ?,
                        email_verified = ?,
                        reg_ip = ?,
                        created_utc_time = ?,
                        remark = ?,
                        last_ip = ?,
                        last_login_utc_time = ?,
                        time_zone = ?,
                        lang = ?,
                        status = ?
                where user_id = ?
           `
	res, err := r.db.Exec(sql,
		user.UserId,
		user.Handle,
		user.DisplayName,
		user.UserGroup,
		user.Pwd,
		user.ProfileUrl,
		user.FaceUrl,
		user.Email,
		user.EmailVerified,
		user.RegIp,
		user.CreatedUTCTime,
		user.Remark,
		user.LastIp,
		user.LastLoginUTCTime,
		user.TimeZone,
		user.Lang,
		user.Status,
		user.UserId)

	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) DeleteByUserId(id int64) (int64, error) {
	sql := "UPDATE br_user SET status = 0 WHERE user_id = ?"
	res, err := r.db.Exec(sql, id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) IsExistHandle(handle string) (bool, error) {
	var count int
	sql := "SELECT count(*) FROM br_user WHERE handle = ?"
	err := r.db.Get(&count, sql, handle)
	return count > 0, err
}

func (r *MockUserRepository) IsExistsEmail(email string) (bool, error) {
	var count int
	sql := "SELECT count(*) FROM br_user WHERE email = ?"
	err := r.db.Get(&count, sql, email)
	return count > 0, err
}

func (r *MockUserRepository) UpdateEmailByUserId(userId int64, email string) (int64, error) {
	sql := "UPDATE br_user SET email = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, email, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) UpdateUserStatus(userId int64, status int8) (int64, error) {
	sql := "UPDATE br_user SET status = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, status, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) UpdateUserPwd(userId int64, pwd string) (int64, error) {
	sql := "UPDATE br_user SET pwd = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, pwd, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) DeleteUser(userId int64) (int64, error) {
	sql := "UPDATE br_user SET status = 0 WHERE user_id = ?"
	res, err := r.db.Exec(sql, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) GetUserByUserId(userId int64) (*user.User, error) {
	sqlQuery := "SELECT * FROM br_user WHERE user_id = ?"
	user := &user.User{}
	err := r.db.Get(user, sqlQuery, userId)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return user, err
}

func (r *MockUserRepository) GetUserByHandle(handle string) (*user.User, error) {
	sqlQuery := "SELECT * FROM br_user WHERE handle = ?"
	user := &user.User{}
	err := r.db.Get(user, sqlQuery, handle)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return user, err
}

func (r *MockUserRepository) GetUserByEmail(email string) (*user.User, error) {
	sqlQuery := "SELECT * FROM br_user WHERE email = ?"
	user := &user.User{}
	err := r.db.Get(user, sqlQuery, email)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return user, err
}

func (r *MockUserRepository) GetUserEncryptPwd(userId int64) (string, error) {
	var pwd string
	sql := "SELECT pwd FROM br_user WHERE user_id = ?"
	err := r.db.Get(&pwd, sql, userId)
	return pwd, err
}

func (r *MockUserRepository) UpdateUserProfile(userId int64, displayName, remark, timeZone, lang string) (int64, error) {
	sql := "UPDATE br_user SET display_name = ?, remark = ?, time_zone = ?, lang = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, displayName, remark, timeZone, lang, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) UpdateUserPassword(userId int64, pwd string) (int64, error) {
	sql := "UPDATE br_user SET pwd = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, pwd, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) InsertValidCode(validCode *user.ValidCode) (int64, error) {
	sqlQuery := `INSERT INTO br_valid_code (user_id, valid_code, valid_type_id, create_utc_time, expire_utc_time) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.Exec(sqlQuery, validCode.UserId, validCode.ValidCode, validCode.ValidTypeId, validCode.CreateUTCTime, validCode.ExpireUTCTime)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *MockUserRepository) DeleteValidCode(id int64) (int64, error) {
	sql := "DELETE FROM br_valid_code WHERE id = ?"
	res, err := r.db.Exec(sql, id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

func (r *MockUserRepository) GetValidCode(userId int64, code string) (*user.ValidCode, error) {
	sqlQuery := "SELECT * FROM br_valid_code WHERE user_id = ? AND valid_code = ?"
	validCode := &user.ValidCode{}
	err := r.db.Get(validCode, sqlQuery, userId, code)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return validCode, err
}

func (r *MockUserRepository) CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	var count int
	var id int64
	sql := "SELECT COUNT(*), COALESCE(MAX(id), 0) FROM br_valid_code WHERE user_id = ? AND code = ? AND valid_type = ?"
	err := r.db.QueryRow(sql, userID, sign, validType).Scan(&count, &id)
	if err != nil {
		return false, 0
	}
	return count > 0, id
}

// 测试用例

func TestUserRepository_InsertUser(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	testUser := &user.User{
		Handle:             "testuser",
		DisplayName:        "Test User",
		UserGroup:          1,
		Pwd:                "hashedpassword",
		ProfileUrl:         "/profile/test",
		FaceUrl:            "/avatar/test.jpg",
		Email:              "test@example.com",
		EmailVerified:      1,
		RegIp:              "192.168.1.1",
		CreatedUTCTime:     1640995200,
		Remark:             "Test user",
		LastIp:             "192.168.1.1",
		LastLoginUTCTime:   1640995200,
		TimeZone:           "UTC",
		Lang:               "en",
		Status:             1,
	}

	// 模拟成功插入
	mock.ExpectExec(`insert into br_user`).WithArgs(
		testUser.Handle, testUser.DisplayName, testUser.UserGroup,
		testUser.Pwd, testUser.ProfileUrl, testUser.FaceUrl,
		testUser.Email, testUser.EmailVerified, testUser.RegIp,
		testUser.CreatedUTCTime, testUser.Remark, testUser.LastIp,
		testUser.LastLoginUTCTime, testUser.TimeZone, testUser.Lang,
		testUser.Status,
	).WillReturnResult(sqlmock.NewResult(123, 1))

	id, err := repo.InsertUser(testUser)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	testUser := &user.User{
		UserId:             123,
		Handle:             "testuser",
		DisplayName:        "Updated User",
		UserGroup:          1,
		Pwd:                "newhashedpassword",
		ProfileUrl:         "/profile/updated",
		FaceUrl:            "/avatar/updated.jpg",
		Email:              "updated@example.com",
		EmailVerified:      1,
		RegIp:              "192.168.1.1",
		CreatedUTCTime:     1640995200,
		Remark:             "Updated user",
		LastIp:             "192.168.1.2",
		LastLoginUTCTime:   1640995300,
		TimeZone:           "UTC",
		Lang:               "en",
		Status:             1,
	}

	// 模拟成功更新
	mock.ExpectExec(`update br_user`).WithArgs(
		testUser.UserId, testUser.Handle, testUser.DisplayName,
		testUser.UserGroup, testUser.Pwd, testUser.ProfileUrl,
		testUser.FaceUrl, testUser.Email, testUser.EmailVerified,
		testUser.RegIp, testUser.CreatedUTCTime, testUser.Remark,
		testUser.LastIp, testUser.LastLoginUTCTime, testUser.TimeZone,
		testUser.Lang, testUser.Status, testUser.UserId,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.UpdateUser(testUser)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByUserId(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	userID := int64(123)
	expectedUser := &user.User{
		UserId:      123,
		Handle:      "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Status:      1,
	}

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"user_id", "handle", "display_name", "user_group", "pwd",
		"profile_url", "face_url", "email", "email_verified", "reg_ip",
		"created_utc_time", "remark", "last_ip", "last_login_utc_time",
		"time_zone", "lang", "status",
	}).AddRow(
		123, "testuser", "Test User", 1, "hashedpassword",
		"/profile/test", "/avatar/test.jpg", "test@example.com", 1, "192.168.1.1",
		1640995200, "Test user", "192.168.1.1", 1640995200,
		"UTC", "en", 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_user WHERE user_id = \?`).WithArgs(userID).WillReturnRows(rows)

	result, err := repo.GetUserByUserId(userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.UserId, result.UserId)
	assert.Equal(t, expectedUser.Handle, result.Handle)
	assert.Equal(t, expectedUser.DisplayName, result.DisplayName)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByHandle(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	handle := "testuser"

	// 模拟查询结果
	rows := sqlmock.NewRows([]string{
		"user_id", "handle", "display_name", "user_group", "pwd",
		"profile_url", "face_url", "email", "email_verified", "reg_ip",
		"created_utc_time", "remark", "last_ip", "last_login_utc_time",
		"time_zone", "lang", "status",
	}).AddRow(
		123, "testuser", "Test User", 1, "hashedpassword",
		"/profile/test", "/avatar/test.jpg", "test@example.com", 1, "192.168.1.1",
		1640995200, "Test user", "192.168.1.1", 1640995200,
		"UTC", "en", 1,
	)

	mock.ExpectQuery(`SELECT \* FROM br_user WHERE handle = \?`).WithArgs(handle).WillReturnRows(rows)

	result, err := repo.GetUserByHandle(handle)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, handle, result.Handle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByHandle_NotFound(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	handle := "nonexistent"

	// 模拟没有找到用户
	mock.ExpectQuery(`SELECT \* FROM br_user WHERE handle = \?`).WithArgs(handle).WillReturnError(sql.ErrNoRows)

	result, err := repo.GetUserByHandle(handle)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_IsExistHandle(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	t.Run("Handle存在", func(t *testing.T) {
		handle := "existinguser"
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM br_user WHERE handle = \?`).WithArgs(handle).WillReturnRows(rows)

		exists, err := repo.IsExistHandle(handle)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Handle不存在", func(t *testing.T) {
		handle := "nonexistentuser"
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM br_user WHERE handle = \?`).WithArgs(handle).WillReturnRows(rows)

		exists, err := repo.IsExistHandle(handle)

		assert.NoError(t, err)
		assert.False(t, exists)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_IsExistsEmail(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	t.Run("邮箱存在", func(t *testing.T) {
		email := "existing@example.com"
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM br_user WHERE email = \?`).WithArgs(email).WillReturnRows(rows)

		exists, err := repo.IsExistsEmail(email)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("邮箱不存在", func(t *testing.T) {
		email := "nonexistent@example.com"
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM br_user WHERE email = \?`).WithArgs(email).WillReturnRows(rows)

		exists, err := repo.IsExistsEmail(email)

		assert.NoError(t, err)
		assert.False(t, exists)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUserStatus(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	userID := int64(123)
	status := int8(0) // 禁用用户

	// 模拟成功更新
	mock.ExpectExec(`UPDATE br_user SET status = \? WHERE user_id = \?`).WithArgs(status, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.UpdateUserStatus(userID, status)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUserPwd(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	userID := int64(123)
	newPassword := "newhashedpassword"

	// 模拟成功更新
	mock.ExpectExec(`UPDATE br_user SET pwd = \? WHERE user_id = \?`).WithArgs(newPassword, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.UpdateUserPwd(userID, newPassword)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser(t *testing.T) {
	repo, mock := NewMockUserRepository(t)
	defer repo.db.Close()

	userID := int64(123)

	// 模拟成功删除（软删除，设置status=0）
	mock.ExpectExec(`UPDATE br_user SET status = 0 WHERE user_id = \?`).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.DeleteUser(userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
