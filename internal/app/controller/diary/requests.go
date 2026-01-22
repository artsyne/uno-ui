package diary

// CreateDiaryRequest 创建日记本请求
//
// 字段:
//   - Title: 日记本标题，必需
//   - Remark: 日记本备注
//   - ViewSetting: 查看权限设置，0=公开，1=私有，2=密码保护
//   - SocialSetting: 社交设置，0=不允许，1=允许
//   - ThemeID: 主题ID
//   - AccessPwd: 访问密码，当ViewSetting为2时使用
type CreateDiaryRequest struct {
	Title         string `form:"title" json:"title" binding:"required"`
	Remark        string `form:"remark" json:"remark"`
	ViewSetting   int8   `form:"view_setting" json:"view_setting" binding:"required,oneof=1 2"`
	SocialSetting int8   `form:"social_setting" json:"social_setting" binding:"required,oneof=1 2"`
	ThemeID       int8   `form:"theme_id" json:"theme_id"`
	AccessPwd     string `form:"access_pwd" json:"access_pwd"`
}

// UpdateDiaryRequest 更新日记本请求
//
// 字段:
//   - DiaryID: 日记本ID，必需
//   - Title: 日记本标题，必需
//   - Remark: 日记本备注
//   - ViewSetting: 查看权限设置，0=公开，1=私有，2=密码保护
//   - SocialSetting: 社交设置，0=不允许，1=允许
//   - ThemeID: 主题ID
//   - AccessPwd: 访问密码，当ViewSetting为2时使用
type UpdateDiaryRequest struct {
	DiaryID       int64  `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
	Title         string `form:"title" json:"title" binding:"required"`
	Remark        string `form:"remark" json:"remark"`
	ViewSetting   int8   `form:"view_setting" json:"view_setting" binding:"required,oneof=1 2"`
	SocialSetting int8   `form:"social_setting" json:"social_setting" binding:"required,oneof=1 2"`
	ThemeID       int8   `form:"theme_id" json:"theme_id"`
	AccessPwd     string `form:"access_pwd" json:"access_pwd"`
}

// ShowDiaryRequest 显示日记请求
// 用于获取特定日记本信息的请求数据结构
//
// 字段:
//   - DiaryID: 要显示的日记本ID，必需参数
type ShowDiaryRequest struct {
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required"`
}

// LikeDiaryRequest 日记本点赞请求
// 用于处理日记本点赞操作的请求数据结构
//
// 字段:
//   - DiaryID: 要点赞的日记本ID，必需参数
type LikeDiaryRequest struct {
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required"`
}

// GetLikeStatusRequest 获取点赞状态请求
// 用于获取日记本点赞状态的请求数据结构
//
// 字段:
//   - DiaryID: 要查询的日记本ID，必需参数
type GetLikeStatusRequest struct {
	DiaryID int64 `form:"diary_id" json:"diary_id" binding:"required"`
}

// DeleteDiaryRequest 删除日记本请求
//
// 字段:
//   - DiaryID: 要删除的日记本ID，必需参数
type DeleteDiaryRequest struct {
	DiaryID int64 `uri:"id" binding:"required"`
}
