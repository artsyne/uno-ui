package comment

// CreateCommentRequest 创建评论请求结构体
// 用于接收前端创建评论的表单数据
type CreateCommentRequest struct {
	EntryID int64  `form:"entry_id" json:"entry_id" binding:"required,min=1"`  // 日记ID，必填且大于0
	Content string `form:"content" json:"content" binding:"required,max=5000"` // 评论内容，必填且不超过5000字符
}

// DeleteCommentRequest 删除评论请求结构体
// 用于接收前端删除评论的请求数据
type DeleteCommentRequest struct {
	CommentID int64 `form:"comment_id" json:"comment_id" binding:"required,min=1"` // 评论ID，必填且大于0
}

// GetCommentsRequest 获取评论列表请求结构体
// 用于接收前端获取评论列表的请求参数
type GetCommentsRequest struct {
	EntryID   int64 `form:"entry_id" json:"entry_id" binding:"required,min=1"`     // 日记ID，必填且大于0
	PageIndex int   `form:"page_index" json:"page_index" binding:"min=1"`        // 页码，最小为1，默认为1
	PageSize  int   `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`  // 每页大小，最小为1，最大为100，默认为10
}

// GetCommentsByPathRequest 获取评论列表请求结构体（路径参数）
// 用于接收前端获取评论列表的请求参数（RESTful风格）
type GetCommentsByPathRequest struct {
	EntryID   int64 `uri:"entry_id" binding:"required,min=1"`                       // 日记ID，必填且大于0
	PageIndex int   `form:"page_index" json:"page_index" binding:"omitempty,min=1"` // 页码，最小为1，默认为1
}