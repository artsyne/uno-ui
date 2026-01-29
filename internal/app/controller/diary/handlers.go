package diary

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// CreateDiary 处理创建新日记本的请求
// 绑定JSON或表单数据到CreateDiaryRequest，并调用服务层创建日记本
//
// 路由: POST /diary
func (ctrl *DiaryController) CreateDiary(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[CreateDiaryRequest](c)

	// 校验描述长度
	if len(req.Remark) > ctrl.config.App.MaxDescriptionLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	nowTime := pkg.GetNowTimeUnixMilli()

	// 使用结构体字面量创建完整的日记本对象
	diaryData := &diary.Diary{
		// 用户输入的字段
		Title:         req.Title,
		Remark:        req.Remark,
		ViewSetting:   int8(req.ViewSetting),
		SocialSetting: int8(req.SocialSetting),
		ThemeID:       int8(req.ThemeID),

		// 系统设置的字段
		DiaryID:            0, // 将由数据库生成
		AccessPwd:          req.AccessPwd,
		AccessURL:          "",
		CoverURL:           "",
		CoverMidURL:        "",
		LikeNum:            0,
		MemberNum:          0,
		ViewNum:            0,
		InfoNum:            0, // 日记数量初始为0
		Status:             1, // 默认状态为正常
		AddUTCTime:         nowTime,
		AddUserID:          loginUser.UserId,
		AddUserHandle:      loginUser.Handle,
	}

	// 调用服务层创建日记本
	diaryID, err := ctrl.diaryService.AddDiary(diaryData)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	diaryData.DiaryID = diaryID

	// 返回成功响应
	if c.GetHeader("HX-Request") == "true" {
		// 设置 HX-Trigger 响应头，触发前端事件
		// showToast: 显示成功提示
		// closeDiaryModal: 关闭日记本弹窗 (需要前端监听此事件)
		c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "common.save_success"}, "closeDiaryModal": true}`)

		data := map[string]interface{}{
			"Diary":    diaryData,
			"Lang":     loginUser.Lang,
			"TimeZone": loginUser.TimeZone,
		}
		c.HTML(http.StatusOK, "user/diary_card", data)
		return
	}

	ctrl.renderService.JSONSuccess(c, diaryData)
}

// UpdateDiary 处理更新日记本的请求
// 绑定URI参数(id)和JSON/表单数据到UpdateDiaryRequest，并调用服务层更新
//
// 路由: PUT /diary/:id
func (ctrl *DiaryController) UpdateDiary(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[UpdateDiaryRequest](c)

	// 校验描述长度
	if len(req.Remark) > ctrl.config.App.MaxDescriptionLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	// 验证日记本所有权
	originalDiary, err := ctrl.diaryService.GetDiaryByID(req.DiaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, &errors.PermissionDenied)
		return
	}
	if originalDiary.AddUserID != loginUser.UserId {
		ctrl.renderService.JSONError(c, &errors.PermissionDenied)
		return
	}

	// 创建更新的日记本数据
	formData := &diary.Diary{
		DiaryID:       req.DiaryID,
		Title:         req.Title,
		Remark:        req.Remark,
		ViewSetting:   int8(req.ViewSetting),
		SocialSetting: int8(req.SocialSetting),
		ThemeID:       int8(req.ThemeID),
	}

	// 保留原有字段值，只更新可编辑字段
	formData.AddUTCTime = originalDiary.AddUTCTime
	formData.AddUserID = originalDiary.AddUserID
	formData.AddUserHandle = originalDiary.AddUserHandle
	formData.AccessURL = originalDiary.AccessURL
	formData.CoverURL = originalDiary.CoverURL
	formData.CoverMidURL = originalDiary.CoverMidURL
	formData.LikeNum = originalDiary.LikeNum
	formData.MemberNum = originalDiary.MemberNum
	formData.ViewNum = originalDiary.ViewNum
	formData.InfoNum = originalDiary.InfoNum
	formData.Status = originalDiary.Status
	formData.AccessPwd = originalDiary.AccessPwd
	if req.AccessPwd != "" {
		formData.AccessPwd = req.AccessPwd
	}

	// 调用服务层更新日记本
	if err := ctrl.diaryService.UpdateDiary(formData); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	if c.GetHeader("HX-Request") == "true" {
		// 构造 HX-Trigger 数据
		triggerData := map[string]interface{}{
			"showToast":       map[string]string{"type": "success", "message": "common.save_success"},
			"closeDiaryModal": true,
			"diaryUpdated":    formData,
		}
		if triggerJSON, err := json.Marshal(triggerData); err == nil {
			c.Header("HX-Trigger", string(triggerJSON))
		}

		// UpdateDiary 通常在详情页调用，我们不需要返回HTML片段，而是让前端根据 diaryUpdated 事件刷新数据
		// 或者返回空内容，因为 HX-Swap="none"
		c.Status(http.StatusOK)
		return
	}

	ctrl.renderService.JSONSuccess(c, formData)
}

// DeleteDiary 删除日记本
// 处理删除日记本的请求，需要用户登录且为日记本所有者
//
// 参数:
//   - c: Gin上下文对象
//
// 路由: DELETE /diary/:id
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[DeleteDiaryRequest](c)

	// 调用服务层删除日记本
	err := ctrl.diaryService.DeleteDiary(req.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, nil)
}

// LikeDiary 日记本点赞
// 处理日记本点赞/取消点赞操作，需要用户登录
//
// 参数:
//   - c: Gin上下文对象
func (ctrl *DiaryController) LikeDiary(c *gin.Context) {
	// 获取登录用户信息
	loginUser, exists := ctrl.auth.GetCurrentUser(c)
	if !exists {
		ctrl.renderService.JSONError(c, &errors.Unauthorized)
		return
	}

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[LikeDiaryRequest](c)
	diaryID := req.DiaryID

	// 检查社交设置
	_, err := ctrl.diaryService.CheckSocialEnabled(diaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 获取点赞前的状态，用于判断是点赞还是取消点赞
	beforeStatus, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, diaryID, like.LikeTypeDiary)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "diary_id", diaryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	// 切换点赞状态
	err = ctrl.likeService.ToggleLike(loginUser.UserId, diaryID, like.LikeTypeDiary)
	if err != nil {
		slog.Error("点赞操作失败", "error", err, "user_id", loginUser.UserId, "diary_id", diaryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	// 根据之前的状态更新日记本的like_num字段
	if beforeStatus.IsLiked {
		// 之前已点赞，现在是取消点赞，减少like_num
		err = ctrl.diaryService.DecrementLikeCount(diaryID)
	} else {
		// 之前未点赞，现在是点赞，增加like_num
		err = ctrl.diaryService.IncrementLikeCount(diaryID)
	}

	if err != nil {
		slog.Error("更新日记本点赞数失败", "error", err, "diary_id", diaryID)
		// 这里不返回错误，因为点赞状态已经更新成功，只是统计数字更新失败
		// 可以考虑记录日志或者异步重试
	}

	// 获取最新的点赞状态
	status, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, diaryID, like.LikeTypeDiary)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "diary_id", diaryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	ctrl.renderService.JSONSuccess(c, status)
}

// GetDiaryLikeStatus 获取日记本点赞状态
// 获取指定日记本的点赞状态信息，包括是否已点赞和点赞总数
//
// 参数:
//   - c: Gin上下文对象
func (ctrl *DiaryController) GetDiaryLikeStatus(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetLikeStatusRequest](c)
	diaryID := req.DiaryID

	// 获取登录用户信息（可能为空，访客也可以查看点赞状态）
	loginUser, _ := ctrl.auth.GetCurrentUser(c)
	var userID int64 = 0
	if loginUser != nil {
		userID = loginUser.UserId
	}

	// 获取点赞状态
	status, err := ctrl.likeService.GetLikeStatus(userID, diaryID, like.LikeTypeDiary)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", userID, "diary_id", diaryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	ctrl.renderService.JSONSuccess(c, status)
}
