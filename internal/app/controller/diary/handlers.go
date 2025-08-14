package diary

import (
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/diary"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ProcessDiary 处理日记本的创建和编辑请求
// 根据表单中的action字段分发到对应的处理方法
//
// 支持的表单action值:
//   - "diary.add": 创建新日记本
//   - "diary.edit": 编辑现有日记本
func (ctrl *DiaryController) ProcessDiary(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析表单数据到Diary结构体
	diaryData := &diary.Diary{}
	if err := c.Bind(diaryData); err != nil {
		ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
		return
	}

	// 根据操作类型分发处理
	action, _ := c.GetPostForm("action")
	switch action {
	case "diary.add":
		ctrl.CreateDiary(c, diaryData, loginUser)
	case "diary.edit":
		ctrl.UpdateDiary(c, diaryData, loginUser)
	default:
		ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
	}
}

// CreateDiary 处理创建新日记本的业务逻辑
// 设置日记本的所有必要字段，调用服务层创建日记本，并返回创建结果
func (ctrl *DiaryController) CreateDiary(c *gin.Context, formData *diary.Diary, user *types.AuthUser) {
	nowTime := pkg.GetNowTimeUnixMilli()

	// 使用结构体字面量创建完整的日记本对象
	diaryData := &diary.Diary{
		// 用户输入的字段
		Title:         formData.Title,
		Remark:        formData.Remark,
		ViewSetting:   formData.ViewSetting,
		SocialSetting: formData.SocialSetting,
		ThemeID:       formData.ThemeID,

		// 系统设置的字段
		DiaryID:            0, // 将由数据库生成
		AccessPwd:          "",
		AccessURL:          "",
		CoverURL:           "",
		CoverMidURL:        "",
		LikeNum:            0,
		MemberNum:          0,
		ViewNum:            0,
		InfoNum:            0, // 日记数量初始为0
		Status:             1, // 默认状态为正常
		AddUTCTime:         nowTime,
		AddUserID:          user.UserId,
		AddUserName:        user.UserName,
		LastActionUTCTime:  nowTime,
		LastActionUserID:   user.UserId,
		LastActionUserName: user.UserName,
	}

	// 调用服务层创建日记本
	diaryID, err := ctrl.diaryService.AddDiary(diaryData)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	diaryData.DiaryID = diaryID

	// 返回成功响应
	ctrl.respondSuccess(c, diaryData)
}

// UpdateDiary 处理编辑现有日记本的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记本
func (ctrl *DiaryController) UpdateDiary(c *gin.Context, formData *diary.Diary, loginUser *types.AuthUser) {
	// 验证日记本所有权
	originalDiary, err := ctrl.validateDiaryOwnership(formData.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 保留原有字段值，只更新可编辑字段
	formData.AddUTCTime = originalDiary.AddUTCTime
	formData.AddUserID = originalDiary.AddUserID
	formData.AddUserName = originalDiary.AddUserName
	formData.AccessURL = originalDiary.AccessURL
	formData.CoverURL = originalDiary.CoverURL
	formData.CoverMidURL = originalDiary.CoverMidURL
	formData.LikeNum = originalDiary.LikeNum
	formData.MemberNum = originalDiary.MemberNum
	formData.ViewNum = originalDiary.ViewNum
	formData.InfoNum = originalDiary.InfoNum
	formData.Status = originalDiary.Status
	formData.AccessPwd = originalDiary.AccessPwd

	// 更新最后操作时间和操作用户
	formData.LastActionUTCTime = pkg.GetNowTimeUnixMilli()
	formData.LastActionUserID = loginUser.UserId
	formData.LastActionUserName = loginUser.UserName

	// 调用服务层更新日记本
	if err := ctrl.diaryService.UpdateDiary(formData); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	ctrl.respondSuccess(c, formData)
}

// DeleteDiary 处理删除日记本请求（未来扩展功能）
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	// 预留的删除功能实现
	ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
}

// respondSuccess 返回成功响应
// 统一的成功响应格式，便于后续维护和修改
func (ctrl *DiaryController) respondSuccess(c *gin.Context, diaryData *diary.Diary) {
	ctrl.renderService.JSONSuccess(c, diaryData)
}