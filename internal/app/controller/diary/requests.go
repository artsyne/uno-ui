package diary

// ProcessDiaryRequest 处理日记本请求（用于统一的处理入口）
type ProcessDiaryRequest struct {
	Action        string `form:"action" json:"action" binding:"required,oneof=diary.add diary.edit"`
	DiaryID       int64  `form:"diary_id" json:"diary_id"`
	Title         string `form:"title" json:"title"`
	Remark        string `form:"remark" json:"remark"`
	ViewSetting   int8   `form:"view_setting" json:"view_setting" binding:"required,oneof=0 1 2"`
	SocialSetting int8   `form:"social_setting" json:"social_setting" binding:"required,oneof=0 1"`
	ThemeID       int8   `form:"theme_id" json:"theme_id"`
	AccessPwd     string `form:"access_pwd" json:"access_pwd"`
}