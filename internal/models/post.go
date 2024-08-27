package models

import (
	"database/sql"
)

// info 表单对象
type PostForm struct {
	PostID   int64  `form:"post_id" json:"post_id" `
	DiaryID  int64  `form:"diary_id" json:"diary_id"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type"`
	Body     string `form:"info_body" json:"body"`
}

// info 数据库对象
type Post struct {
	PostID        int64  `form:"post_id" json:"post_id" db:"post_id"`
	DiaryID       int64  `form:"diary_id" json:"diary_id" db:"diary_id"`
	IsTop         int8   `form:"is_top" json:"is_top" db:"is_top"`
	Title         string `form:"title" json:"title" db:"title"`
	Summary       string `form:"summary" json:"summary" db:"summary"`
	Body          string `form:"body" json:"body" db:"body"`
	ViewNum       int32  `form:"view_num" json:"view_num" db:"view_num"`
	LikeNum       int32  `form:"like_num" json:"like_num" db:"like_num"`
	CommentNum    int32  `form:"comment_num" json:"comment_num" db:"comment_num"`
	AddUserID     int64  `form:"add_user_id" json:"add_user_id" db:"add_user_id"`
	AddUserName   string `form:"add_user_name" json:"add_user_name" db:"add_user_name"`
	AddUserIP     string `form:"add_user_ip" json:"add_user_ip" db:"add_user_ip"`
	AddUTCTime    int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	UpdateUTCTime int64  `form:"update_utc_time" json:"update_utc_time" db:"update_utc_time"`
	// 用户状态：删除：0，邮箱待激活：1，锁定: 2，正常：3
	Status int8 `form:"status" json:"status" db:"status"`
	// 信息类型参数：Text：1，Image：2，Video: 3，Audio：4, File: 5
	InfoType int8 `form:"info_type" json:"info_type" db:"info_type"`
}

// 插入日志信息
func (post *Post) Insert() (int64, error) {
	sqlstr := `
			INSERT INTO br_post(
				diary_id,is_top,
				title,summary,body,view_num,like_num,comment_num,add_user_id,
				add_user_name,add_user_ip,add_utc_time,update_utc_time,status,info_type
			)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
			`
	res, err := db.Exec(sqlstr, post.DiaryID,
		post.IsTop, post.Title, post.Summary, post.Body, post.ViewNum, post.LikeNum,
		post.CommentNum, post.AddUserID, post.AddUserName, post.AddUserIP, post.AddUTCTime,
		post.UpdateUTCTime, post.Status, post.InfoType)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// ListPostsByUserID 根据用户ID和日记本ID分页查询日记信息
// 参数:
// - userID: 用户ID
// - spaceID: 日记本ID
// - pageIndex: 当前页码，从1开始
// - pageSize: 每页显示的记录数
// 返回值:
// - []Info: 查询到的日记信息数组
// - totalRecords: 总记录数
// - error: 错误信息
func ListPostsByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error) {
	infoArray := []Post{}
	totalRecords := 0

	// 计算偏移量，注意：页码通常从1开始，所以需要减1
	offset := (pageIndex - 1) * pageSize

	sqlstr := `SELECT * FROM br_post WHERE diary_id = ? AND add_user_id = ? 
				ORDER BY post_id DESC LIMIT ? offset ? `

	err := db.Select(&infoArray, sqlstr, diaryID, userID, pageSize, offset)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	//获取总记录数
	totalRecords, err = GetPostCountByUserID(userID, diaryID)

	return infoArray, totalRecords, err
}

// GetPostCountByUserID 根据用户ID和日记本ID统计信息数量
// 参数:
// - userID: 用户ID
// - spaceID: 日记本ID
// 返回值:
// - int: 总记录数
// - error: 错误信息
func GetPostCountByUserID(userID int64, diaryID int64) (int, error) {
	var count int
	sqlstr := "SELECT COUNT(*) FROM br_post WHERE diary_id = ? AND add_user_id = ?"

	err := db.Get(&count, sqlstr, diaryID, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}
