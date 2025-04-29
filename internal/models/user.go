package models

import (
	"database/sql"
)

type Login struct {
	LoginName    string `form:"login_name" json:"login_name"`
	Pwd          string `form:"pwd" json:"pwd" `
	AccessUrl    string `form:"access_url" json:"access_url"`
	UserId       int64  `form:"user_id" json:"user_id"`
	UserName     string `form:"user_name" json:"user_name"`
	FullName     string `form:"full_name" json:"full_name"`
	UserGroup    int8   `form:"user_group" json:"user_group"`
	FaceSmallUrl string `form:"face_small_url" json:"face_small_url"`
	LoginUTCTime int64  `form:"login_utc_time" json:"login_utc_time"`
	TimeZone     string `form:"time_zone" json:"time_zone"`
	Lang         string `form:"lang" json:"lang"`
	Status       int8   `form:"status" json:"status"`
}

type User struct {
	UserId           int64  `form:"user_id" json:"user_id" db:"user_id"`
	UserName         string `form:"user_name" json:"user_name" db:"user_name"`
	FullName         string `form:"full_name" json:"full_name" db:"full_name"`
	UserGroup        int8   `form:"user_group" json:"user_group" db:"user_group"`
	Pwd              string `form:"pwd" json:"pwd" db:"pwd"`
	PwdSalt          string `form:"pwd_salt" json:"pwd_salt" db:"pwd_salt"`
	ProfileUrl       string `form:"profile_url" json:"profile_url" db:"profile_url"`
	FaceUrl          string `form:"face_url" json:"face_url" db:"face_url"`
	Email            string `form:"email" json:"email" db:"email"`
	EmailVerified    int8   `form:"email_verified" json:"email_verified" db:"email_verified"`
	RegIp            string `form:"reg_ip" json:"reg_ip" db:"reg_ip"`
	CreatedUTCTime   int64  `form:"created_utc_time" json:"created_utc_time" db:"created_utc_time"`
	Remark           string `form:"remark" json:"remark" db:"remark"`
	LastIp           string `form:"last_ip" json:"last_ip" db:"last_ip"`
	LastLoginUTCTime int64  `form:"last_login_utc_time" json:"last_login_utc_time" db:"last_login_utc_time"`
	TimeZone         string `form:"time_zone" json:"time_zone" db:"time_zone"`
	Lang             string `form:"lang" json:"lang" db:"lang"`
	Status           int8   `form:"status" json:"status" db:"status"`
}

// 插入数据
func (user *User) Insert() (int64, error) {
	sql := `
                insert into br_user  (
                    user_id,
                    user_name,
                    full_name,
                    user_group,
                    pwd,
                    pwd_salt,
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
                    ?,
                    ?,
                    ?
                )
                `

	res, err := db.Exec(sql,
		user.UserId,
		user.UserName,
		user.FullName,
		user.UserGroup,
		user.Pwd,
		user.PwdSalt,
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

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

// 更新数据，方法应该废弃，拆分多个方法，不能一次修改多个关键信息
func (user *User) Update() (int64, error) {
	sql := `
                update br_user
                 set
                        user_id = ?,
                        user_name = ?,
                        full_name = ?,
                        user_group = ?,
                        pwd = ?,
                        pwd_salt = ?,
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
	res, err := db.Exec(sql,
		user.UserId,
		user.UserName,
		user.FullName,
		user.UserGroup,
		user.Pwd,
		user.PwdSalt,
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

// 根据uid更新用户空间地址
// 此方法已废弃，因为 access_url 字段已被移除
// func UpdateAccessUrlByUserId(userId int64, accessUrl string) (int64, error) {
//     sql := "UPDATE br_user SET access_url = ? WHERE user_id = ?"
//     res, err := db.Exec(sql, accessUrl, userId)
//     if err != nil {
//         return 0, err
//     }
//
//     rows, err := res.RowsAffected()
//     return rows, err
// }

// 删除数据，假删除，status = 0
func DeleteByUserId(id int64) (int64, error) {
	sql := `
		UPDATE br_user SET status = 0 where user_id = ?
	`
	res, err := db.Exec(sql, id)

	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// 判断用户名是否存在
func IsExistUserName(userName string) (bool, error) {
	var total int
	err := db.Get(&total, "SELECT count(*) as num FROM br_user WHERE user_name = ?", userName)
	if err != nil {
		return true, err
	}

	result := true
	if total == 0 {
		result = false
	}
	return result, err
}

// 判断邮箱是否存在
func IsExistsEmail(email string) (bool, error) {
	var total int
	err := db.Get(&total, "SELECT COUNT(*) AS num FROM br_user WHERE email =?", email)
	if err != nil {
		return true, err
	}

	result := true
	if total == 0 {
		result = false
	}
	return result, err
}

// 更新邮件地址，通过用户Id
func UpdateEmailByUserId(userId int64, email string) (int64, error) {
	sql := "UPDATE br_user SET email = ? WHERE user_id = ?"
	res, err := db.Exec(sql, userId, email)

	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// 激活用户
func UpdateUserStatus(userId int64, status int8) (int64, error) {
	sql := "UPDATE br_user SET status = ? where user_id = ?"
	res, err := db.Exec(sql, status, userId)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateUserPwd 更新用户密码
func UpdateUserPwd(userId int64, pwd string) (int64, error) {
	sql := "UPDATE br_user SET pwd = ? where user_id = ?"
	res, err := db.Exec(sql, pwd, userId)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// 删除用户，假删除
func DeleteUser(userId int64) (int64, error) {
	sql := "UPDATE br_user SET status = 0 where user_id = ?"
	res, err := db.Exec(sql, userId)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// 获取用户信息
func GetUserByUserId(userId int64) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM br_user WHERE user_id = ?", userId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// 根据用户名获取用户信息
func GetUserByUserName(userName string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM br_user WHERE user_name = ?", userName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// 根据电子邮件获取用户信息
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM br_user WHERE email = ?", email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
