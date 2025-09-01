package post

// CreatePostRequest 创建日记请求
// 用于创建新日记的请求数据结构
//
// 字段:
//   - DiaryID: 所属日记本ID，必需且大于0
//   - Title: 日记标题
//   - InfoType: 信息类型，1-5之间的值
//   - Body: 日记内容，必需且不能为空
//   - IsTop: 是否置顶，0=否，1=是
type CreatePostRequest struct {
	DiaryID  int64  `form:"diary_id" json:"diary_id" binding:"required,min=1"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type" binding:"required,oneof=1 2 3 4 5"`
	Body     string `form:"info_body" json:"body" binding:"required,min=1"`
	IsTop    int8   `form:"is_top" json:"is_top"`
}

// UpdatePostRequest 更新日记请求
// 用于更新现有日记的请求数据结构
//
// 字段:
//   - PostID: 要更新的日记ID，必需且大于0
//   - Title: 日记标题
//   - InfoType: 信息类型，1-5之间的值
//   - Body: 日记内容，必需且不能为空
//   - IsTop: 是否置顶，0=否，1=是
type UpdatePostRequest struct {
	PostID   int64  `form:"post_id" json:"post_id" binding:"required,min=1"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type" binding:"required,oneof=1 2 3 4 5"`
	Body     string `form:"info_body" json:"body" binding:"required,min=1"`
	IsTop    int8   `form:"is_top" json:"is_top"`
}

// DeletePostRequest 删除日记请求
// 用于删除日记的请求数据结构
//
// 字段:
//   - PostID: 要删除的日记ID，必需且大于0
type DeletePostRequest struct {
	PostID int64 `form:"post_id" json:"post_id" binding:"required,min=1"`
}

// ShowPostFormRequest 显示修改、增加post表单
// 用于显示修改、增加post表单的请求数据结构
//
// 字段:
//   - PostID: 要操作的日记ID，必需且大于0
type ShowPostFormRequest struct {
	PostID int64 `form:"post_id" json:"post_id" binding:"required,min=1"`
}

// ListPostsRequest 获取日记列表请求
// 用于获取指定日记本下的日记列表的请求数据结构
//
// 字段:
//   - DiaryID: 日记本ID，必需且大于0
//   - Page: 页码，用于分页
//   - Limit: 每页数量限制
type ListPostsRequest struct {
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required,min=1"`
	Page    int   `form:"page" json:"page"`
	Limit   int   `form:"limit" json:"limit"`
}

// ShowPostDetailRequest 获取日记详情请求
// 用于获取指定日记详情的请求数据结构
//
// 字段:
//   - PostID: 日记ID，必需且大于0
//   - IsAjax: 是否为异步请求，true表示是异步请求
//
// 路由:
//   - /p/:postid
type ShowPostDetailRequest struct {
	PostID int64 `uri:"postid" form:"post_id" json:"post_id" binding:"required,min=1"`
	IsAjax bool `form:"isajax" json:"isajax"`
}