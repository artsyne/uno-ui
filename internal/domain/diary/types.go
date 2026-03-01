package diary

// Diary 日记本结构体
// 包含日记本的基本信息、权限设置、统计数据等
type Diary struct {
	DiaryID       int64  `form:"diary_id" json:"diary_id" db:"diary_id"`                   // 日记本唯一标识ID
	Title         string `form:"title" json:"title" db:"title"`                            // 日记本标题
	Remark        string `form:"remark" json:"remark" db:"remark"`                         // 日记本描述/备注
	AccessURL     string `form:"access_url" json:"access_url" db:"access_url"`             // 访问URL地址
	CoverURL      string `form:"cover_url" json:"cover_url" db:"cover_url"`                // 封面图片URL
	CoverMidURL   string `form:"cover_mid_url" json:"cover_mid_url" db:"cover_mid_url"`    // 中等尺寸封面图片URL
	LikeNum       int32  `form:"like_num" json:"like_num" db:"like_num"`                   // 点赞数量
	MemberNum     int32  `form:"member_num" json:"member_num" db:"member_num"`             // 成员数量
	ViewNum       int32  `form:"view_num" json:"view_num" db:"view_num"`                   // 浏览次数
	ViewSetting   int8   `form:"view_setting" json:"view_setting" db:"view_setting"`       // 查看权限设置 (1:公开 2:私密)
	AccessPwd     string `form:"access_pwd" json:"access_pwd" db:"access_pwd"`             // 访问密码
	SocialSetting int8   `form:"social_setting" json:"social_setting" db:"social_setting"` // 社交功能设置
	JoinUTCTime   int64  `form:"join_utc_time" json:"join_utc_time" db:"join_utc_time"`    // 加入时间 (仅在查询参与日记本时填充)
	AddUTCTime    int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`       // 创建时间（UTC时间戳）
	AddUserID     int64  `form:"add_user_id" json:"add_user_id" db:"add_user_id"`          // 创建者用户ID
	Status        int8   `form:"status" json:"status" db:"status"`                         // 日记本状态（0:删除 1:正常）
	InfoNum       int32  `form:"info_num" json:"info_num" db:"info_num"`                   // 信息数量
	ThemeID       int8   `form:"theme_id" json:"theme_id" db:"theme_id"`                   // 主题ID

	// Joined fields (aliased in query)
	AddUserHandle      string `json:"add_user_handle" db:"handle"`             // 创建者Handle
	AddUserDisplayName string `json:"add_user_display_name" db:"display_name"` // 创建者显示名称
	AuthorStatus       int8   `json:"author_status" db:"author_status"`        // 作者状态 (1:Active, 2:Pending)
}

// ViewSetting constants
const (
	ViewSettingPublic  int8 = 1 // 公开
	ViewSettingPrivate int8 = 2 // 私密
)

// Author 日记本作者结构体
// 对应数据库 br_diary_author 表
type Author struct {
	DiaryID     int64 `json:"diary_id" db:"diary_id"`           // 日记本ID
	UserID      int64 `json:"user_id" db:"user_id"`             // 用户ID
	IsAdmin     int8  `json:"is_admin" db:"is_admin"`           // 是否管理员 (0:否, 1:是)
	JoinUTCTime int64 `json:"join_utc_time" db:"join_utc_time"` // 加入时间
	PostNum     int   `json:"post_num" db:"post_num"`           // 发表文章数
	Status      int8  `json:"status" db:"status"`               // 状态 (1:Active, 2:Pending)

	// Joined fields from br_user
	Handle      string `json:"handle" db:"handle"`             // 用户Handle (Join查询)
	DisplayName string `json:"display_name" db:"display_name"` // 用户显示名称 (Join查询)
	Avatar      string `json:"avatar" db:"face_url"`           // 用户头像 (Join查询)
}

// AuthorStatus 作者状态常量
const (
	AuthorStatusActive  int8 = 1 // 活跃/已加入
	AuthorStatusPending int8 = 2 // 待确认/已邀请
)

// AuthorAdminStatus 管理员状态常量
const (
	AuthorIsAdminNo  int8 = 0 // 普通作者
	AuthorIsAdminYes int8 = 1 // 管理员
)

const (
	StatusDeleted int8 = 0 // 删除
	StatusActive  int8 = 1 // 正常
)

// SocialSetting 社交功能设置常量
const (
	SocialSettingEnabled  int8 = 1 // 开启社交功能
	SocialSettingDisabled int8 = 2 // 关闭社交功能
)
