package models

import (
	"database/sql"

	"github.com/biuren/hebe/helpers"
)

type ValidCode struct {
	ValidCodeId   int64  `form:"valid_code_id" json:"valid_code_id" db:"valid_code_id"`
	UserId        int64  `form:"user_id" json:"user_id" db:"user_id"`
	ValidCode     string `form:"valid_code" json:"valid_code" db:"valid_code"`
	ValidTypeId   int8   `form:"valid_type_id" json:"valid_type_id" db:"valid_type_id"`
	ExpireUTCTime int64  `form:"expire_utc_time" json:"expire_utc_time" db:"expire_utc_time"`
	CreateUTCTime int64  `form:"create_utc_time" json:"gmt_create" db:"create_utc_time"`
}

// 插入数据
func (validCode *ValidCode) Insert() (int64, error) {
	//删除已存在确认码
	_, err := db.Exec("DELETE FROM br_valid_code WHERE user_id = ? and valid_type_id=?", validCode.UserId, validCode.ValidTypeId)
	if err != nil {
		return 0, err
	}
	//增加确认码
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
	res, err := db.Exec(sql,
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

// 删除数据
func DeleteValidCode(id int64) (int64, error) {
	sql := "DELETE FROM br_valid_code WHERE valid_code_id = ?"
	res, err := db.Exec(sql, id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// 获取验证码
// 获取当前用户制定类型最新一条数据
func GetValidCode(userId int64, code string) (*ValidCode, error) {
	sqlstr := "SELECT * FROM br_valid_code WHERE user_id = ? AND valid_code = ? order by valid_code_id desc limit 1"
	validCode := &ValidCode{}
	err := db.Get(validCode, sqlstr, userId, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return validCode, err
}

// CheckValidCode : 检查确认码
// userId 用户ID sign 验证码 validType 在configuration.go定义
func CheckValidCode(userID int64, sign string, validType int8) (bool, int64) {
	//默认返回值
	result := false
	var validId int64 = 0
	//获取验证码
	if validCode, err := GetValidCode(userID, sign); err == nil && validCode != nil {
		//判断验证码类型, 验证码类型不对
		if validCode.ValidTypeId != validType {
			return false, 0
		}
		//判断验证码有效时间
		if validCode.ExpireUTCTime < helpers.GetNowTimeUnixMilli() {
			return false, 0
		}
		result = true
		validId = validCode.ValidCodeId
	}

	return result, validId
}
