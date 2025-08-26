package post

// CreatePostRequest 创建日记请求
type CreatePostRequest struct {
	DiaryID  int64  `form:"diary_id" json:"diary_id" binding:"required,min=1"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type" binding:"required,oneof=1 2 3 4 5"`
	Body     string `form:"info_body" json:"body" binding:"required,min=1"`
	IsTop    int8   `form:"is_top" json:"is_top"`
}

// UpdatePostRequest 更新日记请求
type UpdatePostRequest struct {
	PostID   int64  `form:"post_id" json:"post_id" binding:"required,min=1"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type" binding:"required,oneof=1 2 3 4 5"`
	Body     string `form:"info_body" json:"body" binding:"required,min=1"`
	IsTop    int8   `form:"is_top" json:"is_top"`
}

// DeletePostRequest 删除日记请求
type DeletePostRequest struct {
	PostID int64 `form:"post_id" json:"post_id" binding:"required,min=1"`
}

// GetPostRequest 获取日记请求
type GetPostRequest struct {
	PostID int64 `form:"post_id" json:"post_id" binding:"required,min=1"`
}

// ListPostsRequest 获取日记列表请求
type ListPostsRequest struct {
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required,min=1"`
	Page    int   `form:"page" json:"page"`
	Limit   int   `form:"limit" json:"limit"`
}

