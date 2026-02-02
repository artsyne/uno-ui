package repository

import (
	"database/sql"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/jmoiron/sqlx"
)

// NewUserRepository 创建用户仓储实例
func NewUserRepository() user.RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &UserRepository{db: db}
}

// UserRepository 用户数据访问层实现
type UserRepository struct {
	db *sqlx.DB
}

// ===================== User 相关方法 =====================

// InsertUser 插入用户数据
func (r *UserRepository) InsertUser(user *user.User) (int64, error) {
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
                    modified_utc_time,
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
		user.ModifiedUTCTime,
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

// UpdateUser 更新用户数据
func (r *UserRepository) UpdateUser(user *user.User) (int64, error) {
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
                        modified_utc_time = ?,
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
		user.ModifiedUTCTime,
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

// DeleteByUserId 删除用户数据（假删除，设置status = 0）
func (r *UserRepository) DeleteByUserId(id int64) (int64, error) {
	sql := `UPDATE br_user SET status = ? where user_id = ?`
	res, err := r.db.Exec(sql, user.StatusDeleted, id)

	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// IsExistHandle 判断用户名是否存在
func (r *UserRepository) IsExistHandle(handle string) (bool, error) {
	var total int
	err := r.db.Get(&total, "SELECT count(*) as num FROM br_user WHERE handle = ?", handle)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

// IsExistsEmail 判断邮箱是否存在
func (r *UserRepository) IsExistsEmail(email string) (bool, error) {
	var total int
	err := r.db.Get(&total, "SELECT count(*) as num FROM br_user WHERE email = ?", email)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

// UpdateEmailByUserId 根据用户ID更新邮箱
func (r *UserRepository) UpdateEmailByUserId(userId int64, email string) (int64, error) {
	sql := "UPDATE br_user SET email = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, email, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateUserStatus 更新用户状态
func (r *UserRepository) UpdateUserStatus(userId int64, status int8) (int64, error) {
	sql := "UPDATE br_user SET status = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, status, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateUserPwd 更新用户密码
func (r *UserRepository) UpdateUserPwd(userId int64, pwd string) (int64, error) {
	sql := "UPDATE br_user SET pwd = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, pwd, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateUserHandle 更新用户Handle
func (r *UserRepository) UpdateUserHandle(userId int64, handle string, modifiedUTCTime int64) (int64, error) {
	sql := "UPDATE br_user SET handle = ?, modified_utc_time = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, handle, modifiedUTCTime, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// DeleteUser 彻底删除用户
func (r *UserRepository) DeleteUser(userId int64) (int64, error) {
	sql := "DELETE FROM br_user WHERE user_id = ?"
	res, err := r.db.Exec(sql, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// GetUserByUserId 根据用户ID获取用户信息
func (r *UserRepository) GetUserByUserId(userId int64) (*user.User, error) {
	sqlstr := "SELECT * FROM br_user WHERE user_id = ?"
	userEntity := &user.User{}
	err := r.db.Get(userEntity, sqlstr, userId)
	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound.WithDetails(map[string]interface{}{"reason": "user_not_found"})
	}
	return userEntity, err
}

// GetUserByHandle 根据用户名获取用户信息
func (r *UserRepository) GetUserByHandle(handle string) (*user.User, error) {
	sqlstr := "SELECT * FROM br_user WHERE handle = ?"
	userEntity := &user.User{}
	err := r.db.Get(userEntity, sqlstr, handle)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFound.WithDetails(map[string]interface{}{"reason": "user_not_found"})
		}
		return nil, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()})
	}
	return userEntity, nil
}

// GetUserByEmail 根据邮箱获取用户信息
func (r *UserRepository) GetUserByEmail(email string) (*user.User, error) {
	sqlstr := "SELECT * FROM br_user WHERE email = ?"
	userEntity := &user.User{}
	err := r.db.Get(userEntity, sqlstr, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFound.WithDetails(map[string]interface{}{"reason": "user_not_found"})
		}
		return nil, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()})
	}
	return userEntity, nil
}

// GetUserTimeZone 获取用户时区
func (r *UserRepository) GetUserTimeZone(userId int64) (string, error) {
	sql := "SELECT time_zone FROM br_user WHERE user_id = ?"
	var timeZone string
	err := r.db.Get(&timeZone, sql, userId)
	if err != nil {
		return "", err
	}
	return timeZone, nil
}

// GetUserEncryptPwd 获取用户加密密码
func (r *UserRepository) GetUserEncryptPwd(userId int64) (string, error) {
	sql := "SELECT pwd FROM br_user WHERE user_id = ?"
	var pwd string
	err := r.db.Get(&pwd, sql, userId)
	if err != nil {
		return "", err
	}
	return pwd, nil
}

// UpdateUserProfile 更新用户基本资料
func (r *UserRepository) UpdateUserProfile(userId int64, displayName, remark, timeZone, lang string) (int64, error) {
	sql := "UPDATE br_user SET display_name = ?, remark = ?, time_zone = ?, lang = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, displayName, remark, timeZone, lang, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateUserPassword 更新用户密码
func (r *UserRepository) UpdateUserPassword(userId int64, pwd string) (int64, error) {
	sql := "UPDATE br_user SET pwd = ? WHERE user_id = ?"
	res, err := r.db.Exec(sql, pwd, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// ===================== ValidCode 相关方法 =====================

// InsertValidCode 插入验证码
func (r *UserRepository) InsertValidCode(validCode *user.ValidCode) (int64, error) {
	// 删除已存在确认码
	_, err := r.db.Exec("DELETE FROM br_valid_code WHERE user_id = ? and valid_type_id=?", validCode.UserId, validCode.ValidTypeId)
	if err != nil {
		return 0, err
	}
	// 增加确认码
	sql := `
                insert into br_valid_code  (
                    user_id  ,
                    valid_code  ,
                    valid_type_id  ,
                    expire_utc_time  ,
                    create_utc_time
                )
                values (
                    ?  ,
                    ?  ,
                    ?  ,
                    ?  ,
                    ? 
                )
                `
	res, err := r.db.Exec(sql,
		validCode.UserId,
		validCode.ValidCode,
		validCode.ValidTypeId,
		validCode.ExpireUTCTime,
		validCode.CreateUTCTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// DeleteValidCode 删除验证码
func (r *UserRepository) DeleteValidCode(id int64) (int64, error) {
	sql := "DELETE FROM br_valid_code WHERE valid_code_id = ?"
	res, err := r.db.Exec(sql, id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// GetValidCode 获取验证码
func (r *UserRepository) GetValidCode(userId int64, code string) (*user.ValidCode, error) {
	sqlstr := "SELECT * FROM br_valid_code WHERE user_id = ? AND valid_code = ? order by valid_code_id desc limit 1"
	validCode := &user.ValidCode{}
	err := r.db.Get(validCode, sqlstr, userId, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return validCode, err
}

// CheckValidCode 检查确认码
func (r *UserRepository) CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	// 默认返回值
	result := false
	var validId int64 = 0
	// 获取验证码
	if validCode, err := r.GetValidCode(userID, sign); err == nil && validCode != nil {
		// 判断验证码类型, 验证码类型不对
		if validCode.ValidTypeId != validType {
			return false, 0
		}
		// 判断验证码有效时间
		if validCode.ExpireUTCTime < pkg.GetNowTimeUnixMilli() {
			return false, 0
		}
		result = true
		validId = validCode.ValidCodeId
	}

	return result, validId
}
