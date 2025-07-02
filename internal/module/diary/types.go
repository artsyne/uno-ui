package diary

// Diary 日记本结构体
// 包含日记本的基本信息、权限设置、统计数据等
type Diary struct {
	DiaryID            int64  `form:"diary_id" json:"diary_id" db:"diary_id"`
	Title              string `form:"title" json:"title" db:"title"`
	Remark             string `form:"remark" json:"remark" db:"remark"`
	AccessURL          string `form:"access_url" json:"access_url" db:"access_url"`
	CoverURL           string `form:"cover_url" json:"cover_url" db:"cover_url"`
	CoverMidURL        string `form:"cover_mid_url" json:"cover_mid_url" db:"cover_mid_url"`
	LikeNum            int32  `form:"like_num" json:"like_num" db:"like_num"`
	MemberNum          int32  `form:"member_num" json:"member_num" db:"member_num"`
	ViewNum            int32  `form:"view_num" json:"view_num" db:"view_num"`
	ViewSetting        int8   `form:"view_setting" json:"view_setting" db:"view_setting"`
	AccessPwd          string `form:"access_pwd" json:"access_pwd" db:"access_pwd"`
	SocialSetting      int8   `form:"social_setting" json:"social_setting" db:"social_setting"`
	AddUTCTime         int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	AddUserID          int64  `form:"add_user_id" json:"add_user_id" db:"add_user_id"`
	AddUserName        string `form:"add_user_name" json:"add_user_name" db:"add_user_name"`
	Status             int8   `form:"status" json:"status" db:"status"`
	InfoNum            int32  `form:"info_num" json:"info_num" db:"info_num"`
	LastActionUTCTime  int64  `form:"last_action_utc_time" json:"last_action_utc_time" db:"last_action_utc_time"`
	LastActionUserID   int64  `form:"last_action_user_id" json:"last_action_user_id" db:"last_action_user_id"`
	LastActionUserName string `form:"last_action_user_name" json:"last_action_user_name" db:"last_action_user_name"`
	ThemeID            int8   `form:"theme_id" json:"theme_id" db:"theme_id"`
}
