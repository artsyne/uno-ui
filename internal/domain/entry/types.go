package entry

// EntryRequest 日记表单对象
// 用于接收前端表单提交的数据
type EntryRequest struct {
	EntryID  int64  `form:"entry_id" json:"entry_id" `  // 日记ID
	DiaryID  int64  `form:"diary_id" json:"diary_id"`   // 日记本ID
	Title    string `form:"info_title" json:"title"`    // 日记标题
	InfoType InfoType `form:"info_type" json:"info_type"` // 信息类型（1:文本 2:图片 3:视频 4:音频 5:文件 6: URL）
	Body     string `form:"info_body" json:"body"`      // 日记内容
}

// Entry 日记数据库对象
// 包含日记的完整信息，对应数据库br_entry表
type Entry struct {
	EntryID       int64  `form:"entry_id" json:"entry_id" db:"entry_id"`                      // 日记ID
	DiaryID       int64  `form:"diary_id" json:"diary_id" db:"diary_id"`                      // 日记本ID
	IsTop         int8   `form:"is_top" json:"is_top" db:"is_top"`                            // 是否置顶（0:否 1:是）
	Title         string `form:"title" json:"title" db:"title"`                               // 日记标题
	Summary       string `form:"summary" json:"summary" db:"summary"`                         // 日记摘要
	Body          string `form:"body" json:"body" db:"body"`                                  // 日记内容
	ViewNum       int32  `form:"view_num" json:"view_num" db:"view_num"`                      // 浏览次数
	LikeNum       int32  `form:"like_num" json:"like_num" db:"like_num"`                      // 点赞次数
	CommentNum    int32  `form:"comment_num" json:"comment_num" db:"comment_num"`             // 评论次数
	AddUserID     int64  `form:"add_user_id" json:"add_user_id" db:"add_user_id"`             // 创建用户ID
	AddUserIP     string `form:"add_user_ip" json:"add_user_ip" db:"add_user_ip"`             // 创建用户IP
	AddUTCTime    int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`          // 创建时间（UTC时间戳）
	UpdateUTCTime int64  `form:"update_utc_time" json:"update_utc_time" db:"update_utc_time"` // 更新时间（UTC时间戳）
	Status        int8   `form:"status" json:"status" db:"status"`                            // 状态（0:删除 1:正常 2:草稿）
	InfoType      InfoType `form:"info_type" json:"info_type" db:"info_type"`                   // 信息类型（1:文本 2:图片 3:视频 4:音频 5:文件 6: URL）

	// Joined fields from br_user
	AddUserHandle      string `json:"add_user_handle" db:"handle"`             // 创建者Handle
	AddUserDisplayName string `json:"add_user_display_name" db:"display_name"` // 创建者显示名称

	// Joined fields from br_diary
	DiaryTitle    string `json:"diary_title" db:"diary_title"`       // 日记本标题
	DiaryOwnerID  int64  `json:"diary_owner_id" db:"diary_owner_id"` // 日记本所有者ID
	SocialSetting int8   `json:"social_setting" db:"social_setting"` // 社交设置 (1:开启 2:关闭)
}

// InfoType 信息类型
type InfoType int8

const (
	StatusDeleted int8 = 0 // 删除
	StatusActive  int8 = 1 // 正常
	StatusDraft   int8 = 2 // 草稿

	// 信息类型常量
	InfoTypeText  InfoType = 1 // 文本
	InfoTypeImage InfoType = 2 // 图片
	InfoTypeVideo InfoType = 3 // 视频
	InfoTypeAudio InfoType = 4 // 音频
	InfoTypeFile  InfoType = 5 // 文件
	InfoTypeURL   InfoType = 6 // URL
)
