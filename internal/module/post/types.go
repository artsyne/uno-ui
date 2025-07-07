package post

// PostRequest 日记表单对象
// 用于接收前端表单提交的数据
type PostRequest struct {
	PostID   int64  `form:"post_id" json:"post_id" `
	DiaryID  int64  `form:"diary_id" json:"diary_id"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type"`
	Body     string `form:"info_body" json:"body"`
}

// Post 日记数据库对象
// 包含日记的完整信息，对应数据库br_post表
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
