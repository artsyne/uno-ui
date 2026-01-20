package entry

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/utils/errors"
)

// CreateEntryRequest 创建日记请求
// 用于创建新日记的请求数据结构
//
// 字段:
//   - DiaryID: 所属日记本ID，必需且大于0
//   - Title: 日记标题
//   - InfoType: 信息类型，1-5之间的值
//   - Body: 日记内容，必需且不能为空
//   - IsTop: 是否置顶，0=否，1=是
type CreateEntryRequest struct {
	DiaryID  int64  `form:"diary_id" json:"diary_id" binding:"required,min=1"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type" binding:"required,oneof=1 2 3 4 5"`
	Body     string `form:"info_body" json:"body" binding:"required,min=1"`
	IsTop    int8   `form:"is_top" json:"is_top"`
}

// UpdateEntryRequest 更新日记请求
// 用于更新现有日记的请求数据结构
//
// 字段:
//   - EntryID: 要更新的日记ID，必需且大于0
//   - Title: 日记标题
//   - InfoType: 信息类型，1-5之间的值
//   - Body: 日记内容，必需且不能为空
//   - IsTop: 是否置顶，0=否，1=是
type UpdateEntryRequest struct {
	EntryID  int64  `uri:"id" form:"entry_id" json:"entry_id" binding:"required,min=1"`
	Title    string `form:"info_title" json:"title"`
	InfoType int8   `form:"info_type" json:"info_type" binding:"required,oneof=1 2 3 4 5"`
	Body     string `form:"info_body" json:"body" binding:"required,min=1"`
	IsTop    int8   `form:"is_top" json:"is_top"`
}

// DeleteEntryRequest 删除日记请求
// 用于删除日记的请求数据结构
//
// 字段:
//   - EntryID: 要删除的日记ID，必需且大于0
type DeleteEntryRequest struct {
	EntryID int64 `uri:"id" form:"entry_id" json:"entry_id" binding:"required,min=1"`
	DiaryID int64 `form:"diary_id" json:"diary_id"`
}

// ShowEntryFormRequest 显示修改、增加entry表单
// 用于显示修改、增加entry表单的请求数据结构
//
// 字段:
//   - EntryID: 要操作的日记ID，必需且大于0
type ShowEntryFormRequest struct {
	EntryID int64 `form:"entry_id" json:"entry_id" binding:"required,min=1"`
}

// ListEntriesRequest 获取日记列表请求
// 用于获取指定日记本下的日记列表的请求数据结构
//
// 字段:
//   - DiaryID: 日记本ID，必需且大于0
//   - Page: 页码，用于分页
//   - Limit: 每页数量限制
type ListEntriesRequest struct {
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required,min=1"`
	Page    int   `form:"page" json:"page"`
	Limit   int   `form:"limit" json:"limit"`
}

// ShowEntryDetailRequest 获取日记详情请求
// 用于获取指定日记详情的请求数据结构
//
// 字段:
//   - EntryID: 日记ID，必需且大于0
//
// 路由:
//   - /entry/:entryid
type ShowEntryDetailRequest struct {
	EntryID int64 `uri:"entryid" form:"entry_id" json:"entry_id" binding:"required,min=1"`
}

// GetEntriesStatsRequest 获取日记统计信息请求
// 用于获取多个日记的统计信息
//
// 字段:
//   - EntryIDs: 日记ID列表，格式为"33,334,55,56"，必需且不能为空
//
// 路由:
//   - POST /api/entry/stats
type GetEntriesStatsRequest struct {
	EntryIDs string `json:"entryids" binding:"required,min=1"`
}

// EntryStatsResponse 单个日记统计信息响应
// 用于返回单个日记的统计数据
//
// 字段:
//   - DiaryID: 日记本ID
//   - EntryID: 日记ID
//   - ViewNum: 浏览次数
//   - LikeNum: 点赞次数
//   - CommentNum: 评论次数
type EntryStatsResponse struct {
	DiaryID    int64 `json:"diary_id"`
	EntryID    int64 `json:"entry_id"`
	ViewNum    int32 `json:"view_num"`
	LikeNum    int32 `json:"like_num"`
	CommentNum int32 `json:"comment_num"`
}

// GetEntriesStatsResponse 获取日记统计信息响应
// 用于返回多个日记的统计信息列表
//
// 字段:
//   - Stats: 日记统计信息列表
type GetEntriesStatsResponse struct {
	Stats []EntryStatsResponse `json:"stats"`
}

// ParseEntryIDs 解析entryids字符串为int64切片
// 切片id 不能超过 50，这是hardcode，确保数据库安全
func (r *GetEntriesStatsRequest) ParseEntryIDs() ([]int64, error) {
	if r.EntryIDs == "" {
		return []int64{}, nil
	}

	parts := strings.Split(r.EntryIDs, ",")
	entryIDs := make([]int64, 0, len(parts))
	if len(parts) > 50 {
		return nil, errors.ValidationError.WithDetails("entryids 最多 50 个")
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, errors.InvalidFormat.WithDetails(fmt.Sprintf("invalid entry id: %s", part))
		}
		entryIDs = append(entryIDs, id)
	}

	return entryIDs, nil
}

// GetEntryLikeStatusRequest 获取日记点赞状态请求
// 用于获取多个日记的点赞状态信息
//
// 字段:
//   - EntryIDs: 日记ID列表，格式为"33,334,55,56"，必需且不能为空
//
// 路由:
//   - POST /api/entry/like/status
type GetEntryLikeStatusRequest struct {
	EntryIDs string `json:"entryids" binding:"required,min=1"`
}

// ParseEntryIDs 解析entryids字符串为int64切片
// 切片id 不能超过 50，这是hardcode，确保数据库安全
func (r *GetEntryLikeStatusRequest) ParseEntryIDs() ([]int64, error) {
	if r.EntryIDs == "" {
		return []int64{}, nil
	}

	parts := strings.Split(r.EntryIDs, ",")
	entryIDs := make([]int64, 0, len(parts))
	if len(parts) > 50 {
		return nil, errors.ValidationError.WithDetails("entryids 最多 50 个")
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, errors.InvalidFormat.WithDetails(fmt.Sprintf("invalid entry id: %s", part))
		}
		entryIDs = append(entryIDs, id)
	}

	return entryIDs, nil
}

// EntryLikeStatusResponse 单个日记点赞状态响应
// 用于返回单个日记的点赞状态
//
// 字段:
//   - EntryID: 日记ID
//   - Liked: 是否已点赞
type EntryLikeStatusResponse struct {
	EntryID int64 `json:"entry_id"`
	Liked   bool  `json:"liked"`
}

// GetEntryLikeStatusResponse 获取日记点赞状态响应
// 用于返回多个日记的点赞状态列表
//
// 字段:
//   - Status: 日记点赞状态列表
type GetEntryLikeStatusResponse struct {
	Status []EntryLikeStatusResponse `json:"status"`
}

// LikeEntryRequest 日记点赞请求
// 用于日记点赞/取消点赞操作的请求数据结构
//
// 字段:
//   - EntryID: 要点赞的日记ID，必需且大于0
type LikeEntryRequest struct {
	EntryID int64 `json:"entry_id" binding:"required,min=1"`
}
