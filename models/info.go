package models

import (
	"database/sql"
	"time"
)

// info 表单对象
type InfoForm struct {
	InfoID       int64  `form:"info_id" json:"info_id" `
	SpaceID      int64  `form:"space_id" json:"space_id"`
	CoverImgData string `form:"cover_data" json:"cover_data" `
	CoverStyle   string `form:"cover_style" json:"cover_style"`
	Title        string `form:"info_title" json:"title"`
	Body         string `form:"info_body" json:"body"`
}

//info 数据库对象
type Info struct {
	InfoID      int64     `form:"info_id" json:"info_id" db:"info_id"`
	SpaceID     int64     `form:"space_id" json:"space_id" db:"space_id"`
	AccessURL   string    `form:"access_url" json:"access_url" db:"access_url"`
	CoverURL    string    `form:"cover_url" json:"cover_url" db:"cover_url"`
	CoverMidURL string    `form:"cover_mid_url" json:"cover_mid_url" db:"cover_mid_url"`
	CoverStyle  string    `form:"cover_style" json:"cover_style" db:"cover_style"`
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
	AddTime     time.Time `form:"add_time" json:"add_time" db:"add_time"`
	UpdateTime  time.Time `form:"update_time" json:"update_time" db:"update_time"`
	Status      int8      `form:"status" json:"status" db:"status"`
}

//插入日志信息
func (info *Info) Insert() (int64, error) {
	sqlstr := `
			INSERT INTO br_info(
				space_id,access_url,cover_url,cover_mid_url,cover_style,is_top,
				title,summary,body,view_num,like_num,comment_num,add_user_id,
				add_user_name,add_user_ip,add_time,update_time,status
			)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
			`
	res, err := db.Exec(sqlstr, info.SpaceID, info.AccessURL, info.CoverURL, info.CoverMidURL, info.CoverStyle,
		info.IsTop, info.Title, info.Summary, info.Body, info.ViewNum, info.LikeNum,
		info.CommentNum, info.AddUserID, info.AddUserName, info.AddUserIP, info.AddTime,
		info.UpdateTime, info.Status)

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