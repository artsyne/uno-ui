package models

import (
	"database/sql"
)

// info 表单对象
type InfoForm struct {
	InfoID       int64  `form:"info_id" json:"info_id" `
	SpaceID      int64  `form:"space_id" json:"space_id"`
	Title        string `form:"info_title" json:"title"`
	InfoType	 int8   `form:"info_type" json:"info_type"`
	Body         string `form:"info_body" json:"body"`
}

//info 数据库对象
type Info struct {
	InfoID      int64     `form:"info_id" json:"info_id" db:"info_id"`
	SpaceID     int64     `form:"space_id" json:"space_id" db:"space_id"`
	IsTop       int8      `form:"is_top" json:"is_top" db:"is_top"`
	Title       string    `form:"title" json:"title" db:"title"`
	Summary     string    `form:"summary" json:"summary" db:"summary"`
	Body        string    `form:"body" json:"body" db:"body"`
	ViewNum     int32     `form:"view_num" json:"view_num" db:"view_num"`
	LikeNum     int32     `form:"like_num" json:"like_num" db:"like_num"`
	CommentNum  int32     `form:"comment_num" json:"comment_num" db:"comment_num"`
	AddUserID   int64     `form:"add_user_id" json:"add_user_id" db:"add_user_id"`
	AddUserName string    `form:"add_user_name" json:"add_user_name" db:"add_user_name"`
	AddUserIP   string    `form:"add_user_ip" json:"add_user_ip" db:"add_user_ip"`
	AddUTCTime     int64 `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	UpdateUTCTime  int64 `form:"update_utc_time" json:"update_utc_time" db:"update_utc_time"`
	// 用户状态：删除：0，邮箱待激活：1，锁定：2，正常：3
	Status      int8      `form:"status" json:"status" db:"status"`
	// 信息类型参数：Text：1，Image：2，Video：3，Audio：4
	InfoType    int8      `form:"info_type" json:"info_type" db:"info_type"`
}

//插入日志信息
func (info *Info) Insert() (int64, error) {
	sqlstr := `
			INSERT INTO br_info(
				space_id,is_top,
				title,summary,body,view_num,like_num,comment_num,add_user_id,
				add_user_name,add_user_ip,add_utc_time,update_utc_time,status,info_type
			)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
			`
	res, err := db.Exec(sqlstr, info.SpaceID,
		info.IsTop, info.Title, info.Summary, info.Body, info.ViewNum, info.LikeNum,
		info.CommentNum, info.AddUserID, info.AddUserName, info.AddUserIP, info.AddUTCTime,
		info.UpdateUTCTime, info.Status, info.InfoType)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// 根据用户ID查询信息
func ListSpaceInfoByUserID(userID int64, sapceID int64)([]Info,error){
	infoArray := []Info{}
	sqlstr := 	"SELECT * FROM br_info WHERE space_id = ? AND add_user_id = ? ORDER BY info_id DESC"
	
	err := db.Select(&infoArray, sqlstr, sapceID,userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return infoArray, err
}