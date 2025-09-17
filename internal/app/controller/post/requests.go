package post

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/utils/errors"
)

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
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required,min=1"`
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

// GetPostsStatsRequest 获取日记统计信息请求
// 用于获取指定日记本下多个日记的统计信息
//
// 字段:
//   - DiaryID: 日记本ID，必需且大于0
//   - PostIDs: 日记ID列表，格式为"33,334,55,56"，必需且不能为空
//
// 路由:
//   - /d/:diaryid/stats?postids=33,334,55,56
type GetPostsStatsRequest struct {
	DiaryID int64  `uri:"diaryid" json:"diary_id" binding:"required,min=1"`
	PostIDs string `form:"postids" json:"postids" binding:"required,min=1"`
}

// PostStatsResponse 单个日记统计信息响应
// 用于返回单个日记的统计数据
//
// 字段:
//   - DiaryID: 日记本ID
//   - PostID: 日记ID
//   - ViewNum: 浏览次数
//   - LikeNum: 点赞次数
//   - CommentNum: 评论次数
type PostStatsResponse struct {
	DiaryID    int64 `json:"diary_id"`
	PostID     int64 `json:"post_id"`
	ViewNum    int32 `json:"view_num"`
	LikeNum    int32 `json:"like_num"`
	CommentNum int32 `json:"comment_num"`
}

// GetPostsStatsResponse 获取日记统计信息响应
// 用于返回多个日记的统计信息列表
//
// 字段:
//   - Stats: 日记统计信息列表
type GetPostsStatsResponse struct {
	Stats []PostStatsResponse `json:"stats"`
}

// ParsePostIDs 解析postids字符串为int64切片
// 切片id 不能超过 50，这是hardcode，确保数据库安全
func (r *GetPostsStatsRequest) ParsePostIDs() ([]int64, error) {
	if r.PostIDs == "" {
		return []int64{}, nil
	}

	parts := strings.Split(r.PostIDs, ",")
	postIDs := make([]int64, 0, len(parts))
	if len(parts) > 50 {
		return nil, errors.ValidationError.WithDetails("postids 最多 50 个")
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, errors.InvalidFormat.WithDetails(fmt.Sprintf("invalid post id: %s", part))
		}
		postIDs = append(postIDs, id)
	}

	return postIDs, nil
}