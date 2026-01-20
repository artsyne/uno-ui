package comment

// CommentRequest 评论表单对象
// 用于接收前端表单提交的数据
type CommentRequest struct {
	CommentID int64  `form:"comment_id" json:"comment_id"` // 评论ID
	EntryID    int64  `form:"entry_id" json:"entry_id"`       // 日记ID
	Content   string `form:"content" json:"content"`       // 评论内容
}

// Comment 评论数据库对象
// 包含评论的完整信息，对应数据库br_comment表
type Comment struct {
	CommentID   int64  `form:"comment_id" json:"comment_id" db:"comment_id"`       // 评论ID
	EntryID     int64  `form:"entry_id" json:"entry_id" db:"entry_id"`             // 日记ID
	Content     string `form:"content" json:"content" db:"content"`                // 评论内容
	AddUtcTime  int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"` // 创建时间（UTC 毫秒）
	AddUserID   int64  `form:"add_user_id" json:"add_user_id" db:"add_user_id"`       // 创建用户ID
	AddUserName string `form:"add_user_name" json:"add_user_name" db:"add_user_name"` // 创建用户名
	AddUserIP   string `form:"add_user_ip" json:"add_user_ip" db:"add_user_ip"`       // 创建用户IP
}

// CommentListResponse 评论列表响应对象
// 用于返回分页的评论列表
type CommentListResponse struct {
	Comments   []Comment `json:"comments"`    // 评论列表
	Total      int       `json:"total"`       // 总数量
	PageIndex  int       `json:"page_index"`  // 当前页码
	PageSize   int       `json:"page_size"`   // 每页大小
	TotalPages int       `json:"total_pages"` // 总页数
	TimeZone   string    `json:"time_zone"`   // 用户时区（用于客户端格式化）
}

// CommentDeleteRequest 删除评论请求对象
type CommentDeleteRequest struct {
	CommentID int64 `form:"comment_id" json:"comment_id"` // 评论ID
}
