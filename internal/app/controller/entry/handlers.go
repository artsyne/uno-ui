package entry

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/biuren/hebe/internal/app/middleware"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// CreateEntry 已被 PublishEntry 替代，此函数应被删除或保留作为兼容。
// 目前已从路由中移除。

// AutoSaveEntry 自动保存草稿（创建/更新）
// 路由: POST /entry/auto-save
func (ctrl *EntryController) AutoSaveEntry(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[AutoSaveEntryRequest](c)

	if len(req.Body) > ctrl.config.App.MaxEntryLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	nowTime := pkg.GetNowTimeUnixMilli()
	if req.EntryID <= 0 {
		entryData := &entry.Entry{
			DiaryID:            req.DiaryID,
			Title:              req.Title,
			Body:               req.Body,
			Summary:            pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
			IsTop:              0,
			InfoType:           entry.InfoTypeText,
			AddUTCTime:         nowTime,
			AddUserID:          loginUser.UserId,
			AddUserHandle:      loginUser.Handle,
			AddUserDisplayName: loginUser.DisplayName,
			AddUserIP:          c.ClientIP(),
			UpdateUTCTime:      nowTime,
			Status:             entry.StatusDraft,
		}

		entryID, err := ctrl.entryService.AddEntry(entryData)
		if err != nil {
			ctrl.renderService.JSONError(c, err)
			return
		}

		ctrl.renderService.JSONSuccess(c, gin.H{"entry_id": entryID, "saved_at": nowTime})
		return
	}

	originalEntry, err := ctrl.entryService.GetEntryByIDForUser(req.EntryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	if originalEntry.DiaryID != req.DiaryID {
		ctrl.renderService.JSONError(c, errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "diary_id_mismatch"}))
		return
	}
	// 允许已发布的日记进行自动保存（实时更新）
	if originalEntry.Status != entry.StatusDraft && originalEntry.Status != entry.StatusActive {
		ctrl.renderService.JSONError(c, errors.OperationNotAllowed.WithDetails(map[string]interface{}{"reason": "only_draft_or_active_can_auto_save"}))
		return
	}

	// double check ownership (though GetEntryByIDForUser implies it)
	if originalEntry.AddUserID != loginUser.UserId {
		ctrl.renderService.JSONError(c, errors.PermissionDenied)
		return
	}

	// 使用结构体字面量构建更新后的日记对象
	entryData := &entry.Entry{
		EntryID:       originalEntry.EntryID,
		DiaryID:       originalEntry.DiaryID,
		Title:         req.Title,
		Body:          req.Body,
		Summary:       pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
		UpdateUTCTime: nowTime,
		ViewNum:       originalEntry.ViewNum,
		LikeNum:       originalEntry.LikeNum,
		CommentNum:    originalEntry.CommentNum,
		IsTop:         originalEntry.IsTop,
		InfoType:      originalEntry.InfoType,
		Status:        originalEntry.Status, // 保持原有状态
	}
	if err := ctrl.entryService.UpdateEntry(entryData, loginUser.UserId); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	ctrl.renderService.JSONSuccess(c, gin.H{"entry_id": originalEntry.EntryID, "saved_at": nowTime})
}

// PublishEntry 发布日记（新建/发布草稿）
// 路由: POST /entry/publish
func (ctrl *EntryController) PublishEntry(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[PublishEntryRequest](c)

	if len(req.Body) > ctrl.config.App.MaxEntryLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	nowTime := pkg.GetNowTimeUnixMilli()
	if req.EntryID > 0 {
		originalEntry, err := ctrl.entryService.GetEntryByIDForUser(req.EntryID, loginUser.UserId)
		if err != nil {
			ctrl.renderService.JSONError(c, err)
			return
		}
		if originalEntry.DiaryID != req.DiaryID {
			ctrl.renderService.JSONError(c, errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "diary_id_mismatch"}))
			return
		}
		if originalEntry.Status != entry.StatusDraft {
			ctrl.renderService.JSONError(c, errors.OperationNotAllowed.WithDetails(map[string]interface{}{"reason": "only_draft_can_publish"}))
			return
		}

		entryData := &entry.Entry{
			EntryID:            originalEntry.EntryID,
			DiaryID:            originalEntry.DiaryID,
			Title:              req.Title,
			Body:               req.Body,
			Summary:            pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
			UpdateUTCTime:      nowTime,
			ViewNum:            originalEntry.ViewNum,
			LikeNum:            originalEntry.LikeNum,
			CommentNum:         originalEntry.CommentNum,
			IsTop:              originalEntry.IsTop,
			InfoType:           originalEntry.InfoType,
			Status:             entry.StatusActive,
			AddUserID:          originalEntry.AddUserID,
			AddUserHandle:      originalEntry.AddUserHandle,
			AddUserDisplayName: originalEntry.AddUserDisplayName,
		}
		if err := ctrl.entryService.UpdateEntry(entryData, loginUser.UserId); err != nil {
			ctrl.renderService.JSONError(c, err)
			return
		}

		if err := ctrl.diaryService.IncrementEntryCount(entryData.DiaryID); err != nil {
			slog.Error("Failed to increment entry count", "error", err, "diary_id", entryData.DiaryID)
		}

		ctrl.renderPublishedEntryResponse(c, loginUser, entryData)
		return
	}

	entryData := &entry.Entry{
		DiaryID:            req.DiaryID,
		Title:              req.Title,
		Body:               req.Body,
		Summary:            pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
		IsTop:              0,
		InfoType:           entry.InfoTypeText,
		AddUTCTime:         nowTime,
		AddUserID:          loginUser.UserId,
		AddUserHandle:      loginUser.Handle,
		AddUserDisplayName: loginUser.DisplayName,
		AddUserIP:          c.ClientIP(),
		UpdateUTCTime:      nowTime,
		Status:             entry.StatusActive,
		ViewNum:            0,
		LikeNum:            0,
		CommentNum:         0,
	}
	entryID, err := ctrl.entryService.AddEntry(entryData)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	entryData.EntryID = entryID

	if err := ctrl.diaryService.IncrementEntryCount(entryData.DiaryID); err != nil {
		slog.Error("Failed to increment entry count", "error", err, "diary_id", entryData.DiaryID)
	}

	ctrl.renderPublishedEntryResponse(c, loginUser, entryData)
}

func (ctrl *EntryController) renderPublishedEntryResponse(c *gin.Context, loginUser *types.AuthUser, entryData *entry.Entry) {
	if c.GetHeader("HX-Request") == "true" {
		trigger := map[string]interface{}{
			"showToast":       map[string]string{"type": "success", "message": "common.save_success"},
			"closeEntryModal": true,
		}
		if c.Query("source") == "drafts" {
			trigger["refresh-drafts"] = true
		}
		if triggerJSON, err := json.Marshal(trigger); err == nil {
			c.Header("HX-Trigger", string(triggerJSON))
		}

		diaryData, _ := ctrl.diaryService.GetDiaryByID(entryData.DiaryID)
		socialSetting := 1
		if diaryData != nil {
			socialSetting = int(diaryData.SocialSetting)
		}

		data := map[string]interface{}{
			"Entry":         entryData,
			"IsOwner":       true,
			"DiaryID":       entryData.DiaryID,
			"Lang":          loginUser.Lang,
			"TimeZone":      loginUser.TimeZone,
			"IsLiked":       false,
			"SocialSetting": socialSetting,
		}
		c.HTML(http.StatusOK, "diary/entry_card", data)
		return
	}

	ctrl.renderService.JSONSuccess(c, entryData)
}

// UpdateEntry 处理更新日记的请求
// 绑定URI参数(id)和JSON/表单数据到UpdateEntryRequest，并调用服务层更新
//
// 路由: PUT /entry/:id
func (ctrl *EntryController) UpdateEntry(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[UpdateEntryRequest](c)

	// 校验日记内容长度
	if len(req.Body) > ctrl.config.App.MaxEntryLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	entryID := req.EntryID

	// 验证日记所有权
	// 使用 GetEntryByIDForUser 以支持修改草稿状态的日记
	// GetEntryByIDForUser 内部已经包含了基于角色的权限检查 (Owner/Admin/Author)
	originalEntry, err := ctrl.entryService.GetEntryByIDForUser(entryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.EntryNotFound)
		return
	}
	// 移除 Controller 层的硬编码所有权检查，完全依赖 Service 层


	// 使用结构体字面量构建更新后的日记对象
	entryData := &entry.Entry{
		EntryID:       entryID,
		DiaryID:       originalEntry.DiaryID, // 保持原有日记本ID
		Title:         req.Title,
		Body:          req.Body,
		Summary:       pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
		UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
		// 保留原有字段值
		ViewNum:    originalEntry.ViewNum,
		LikeNum:    originalEntry.LikeNum,
		CommentNum: originalEntry.CommentNum,
		IsTop:      originalEntry.IsTop,
		InfoType:   originalEntry.InfoType,
		Status:     entry.StatusActive, // 强制更新为Active状态（即发布）
	}

	// 调用服务层更新
	if err := ctrl.entryService.UpdateEntry(entryData, loginUser.UserId); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 如果是从草稿状态发布，需要增加日记本的日记计数
	if originalEntry.Status == entry.StatusDraft {
		if err := ctrl.diaryService.IncrementEntryCount(entryData.DiaryID); err != nil {
			slog.Error("Failed to increment entry count", "error", err, "diary_id", entryData.DiaryID)
		}
	}

	// HTMX Support
	if c.GetHeader("HX-Request") == "true" {
		trigger := map[string]interface{}{
			"showToast":       map[string]string{"type": "success", "message": "common.save_success"},
			"closeEntryModal": true,
		}
		if c.Query("source") == "drafts" {
			trigger["refreshDrafts"] = true
		}
		if triggerJSON, err := json.Marshal(trigger); err == nil {
			c.Header("HX-Trigger", string(triggerJSON))
		}

		// Get Like Status
		isLiked := false
		likeStatus, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, entryID, like.LikeTypeEntry)
		if err == nil {
			isLiked = likeStatus.IsLiked
		}

		// Get Diary Info for SocialSetting
		diaryData, _ := ctrl.diaryService.GetDiaryByID(entryData.DiaryID)
		socialSetting := 1
		if diaryData != nil {
			socialSetting = int(diaryData.SocialSetting)
		}

		data := map[string]interface{}{
			"Entry":         entryData,
			"IsOwner":       true, // User must be owner to update
			"DiaryID":       entryData.DiaryID,
			"Lang":          loginUser.Lang,
			"TimeZone":      loginUser.TimeZone,
			"IsLiked":       isLiked,
			"SocialSetting": socialSetting,
			"CurrentUserID": loginUser.UserId,
		}
		c.HTML(http.StatusOK, "diary/entry_card", data)
		return
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, entryData)
}

// DeleteEntry 处理删除日记请求
// 执行软删除操作，将日记状态设置为0，并减少日记本的统计数
//
// 路由: DELETE /entry/:id
func (ctrl *EntryController) DeleteEntry(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[DeleteEntryRequest](c)

	// 获取当前用户信息
	userData, isAuthenticated := ctrl.auth.GetCurrentUser(c)
	if !isAuthenticated {
		ctrl.renderService.JSONError(c, &errors.LoginRequired)
		return
	}

	// 验证日记所有权 - 交由 Service 层处理
	// DeleteEntryByID 内部会调用 CanUpdateEntry 检查权限 (Owner/Admin/Author)
	
	// 执行删除操作
	currentTime := pkg.GetNowTimeUnixMilli()
	deleteErr := ctrl.entryService.DeleteEntryByID(req.EntryID, userData.UserId, currentTime)
	if deleteErr != nil {
		ctrl.renderService.JSONError(c, deleteErr)
		return
	}

	// HTMX 请求支持
	if c.GetHeader("HX-Request") == "true" {
		c.Status(http.StatusOK)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// RecoverEntry 恢复已删除的日记
// 路由: PUT /entry/:id/recover
func (ctrl *EntryController) RecoverEntry(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取参数
	entryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.InvalidParams)
		return
	}

	// 调用服务层恢复日记
	if err := ctrl.entryService.RecoverEntryByID(entryID, loginUser.UserId); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// HTMX Support
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "common.recover_success"}}`)
		// 如果在列表中，移除该行
		c.Status(http.StatusOK)
		return
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, nil)
}

// HardDeleteEntry 物理删除日记
// 路由: DELETE /entry/:id/hard
func (ctrl *EntryController) HardDeleteEntry(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取参数
	entryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.InvalidParams)
		return
	}

	// 调用服务层物理删除日记
	if err := ctrl.entryService.HardDeleteEntryByID(entryID, loginUser.UserId); err != nil {
		if c.GetHeader("HX-Request") == "true" {
			c.Header("HX-Trigger", `{"showToast": {"type": "error", "message": "common.operation_failed"}}`)
			c.Status(http.StatusOK)
			return
		}
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 物理删除日记相关的评论
	if err := ctrl.commentService.DeleteCommentsByEntryID(entryID); err != nil {
		// 记录错误但中断流程，因为日记已经删除
		slog.Error("Failed to delete comments for hard deleted entry", "error", err, "entry_id", entryID)
	}

	// HTMX Support
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "common.delete_success"}}`)
		// 移除该行
		c.Status(http.StatusOK)
		return
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, nil)
}

// GetEntriesStats 获取日记统计信息
// 根据POST请求体中的entryids获取日记的统计信息
//
// 路由: POST /api/entry/stats
// 参数:
//   - entryids: 日记ID列表，逗号分隔 (JSON请求体)
//
// 返回:
//   - 包含统计信息的日记列表
func (ctrl *EntryController) GetEntriesStats(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetEntriesStatsRequest](c)

	// 解析entryids字符串为int64切片
	entryIDs, err := req.ParseEntryIDs()
	if err != nil {
		ctrl.renderService.JSONError(c, err) // ParseEntryIDs已经返回了合适的错误类型
		return
	}

	// 调用服务层获取统计信息
	entries, err := ctrl.entryService.GetEntriesStats(entryIDs)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 构建响应数据
	statsResponse := make([]EntryStatsResponse, 0, len(entries))
	for _, e := range entries {
		statsResponse = append(statsResponse, EntryStatsResponse{
			DiaryID:    e.DiaryID,
			EntryID:    e.EntryID,
			ViewNum:    e.ViewNum,
			LikeNum:    e.LikeNum,
			CommentNum: e.CommentNum,
		})
	}

	// 返回成功响应
	response := GetEntriesStatsResponse{
		Stats: statsResponse,
	}
	ctrl.renderService.JSONSuccess(c, response)
}

// LikeEntry 日记点赞
// 处理日记点赞/取消点赞操作，需要用户登录
//
// 参数:
//   - c: Gin上下文对象
func (ctrl *EntryController) LikeEntry(c *gin.Context) {
	// 获取登录用户信息
	loginUser, exists := ctrl.auth.GetCurrentUser(c)
	if !exists {
		ctrl.renderService.JSONError(c, &errors.Unauthorized)
		return
	}

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[LikeEntryRequest](c)
	entryID := req.EntryID

	// 验证 Entry 是否存在
	entryData, err := ctrl.entryService.GetEntryByID(entryID)
	if err != nil {
		ctrl.renderService.JSONError(c, &errors.EntryNotFound)
		return
	}

	// 验证 Diary 社交设置
	_, err = ctrl.diaryService.CheckSocialEnabled(entryData.DiaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 获取点赞前的状态，用于判断是点赞还是取消点赞
	beforeStatus, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, entryID, like.LikeTypeEntry)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "entry_id", entryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	// 切换点赞状态
	err = ctrl.likeService.ToggleLike(loginUser.UserId, entryID, like.LikeTypeEntry)
	if err != nil {
		slog.Error("点赞操作失败", "error", err, "user_id", loginUser.UserId, "entry_id", entryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	// 根据之前的状态更新日记的like_num字段
	if beforeStatus.IsLiked {
		// 之前已点赞，现在是取消点赞，减少like_num
		err = ctrl.entryService.DecrementLikeCount(entryID)
	} else {
		// 之前未点赞，现在是点赞，增加like_num
		err = ctrl.entryService.IncrementLikeCount(entryID)
	}

	if err != nil {
		slog.Error("更新日记点赞数失败", "error", err, "entry_id", entryID)
		// 这里不返回错误，因为点赞状态已经更新成功，只是统计数字更新失败
		// 可以考虑记录日志或者异步重试
	}

	// 获取最新的点赞状态
	status, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, entryID, like.LikeTypeEntry)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "entry_id", entryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	ctrl.renderService.JSONSuccess(c, status)
}

// GetEntriesLikeStatus 获取多个日记的点赞状态
// 根据POST请求体中的entryids获取用户对这些日记的点赞状态
//
// 路由: POST /api/entry/like/status
// 参数:
//   - entryids: 日记ID列表，逗号分隔 (JSON请求体)
//
// 返回:
//   - 包含点赞状态信息的日记列表
func (ctrl *EntryController) GetEntriesLikeStatus(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetEntryLikeStatusRequest](c)

	// 解析entryids字符串为int64切片
	entryIDs, err := req.ParseEntryIDs()
	if err != nil {
		ctrl.renderService.JSONError(c, err) // ParseEntryIDs已经返回了合适的错误类型
		return
	}

	// 获取登录用户信息（可能为空，访客也可以查看点赞状态）
	loginUser, _ := ctrl.auth.GetCurrentUser(c)
	var userID int64 = 0
	if loginUser != nil {
		userID = loginUser.UserId
	}

	// 构建响应数据
	statusResponse := make([]EntryLikeStatusResponse, 0, len(entryIDs))
	for _, entryID := range entryIDs {
		// 获取点赞状态
		status, err := ctrl.likeService.GetLikeStatus(userID, entryID, like.LikeTypeEntry)
		if err != nil {
			slog.Error("获取点赞状态失败", "error", err, "user_id", userID, "entry_id", entryID)
			// 对于单个日记的错误，我们继续处理其他日记，但记录错误
			continue
		}

		statusResponse = append(statusResponse, EntryLikeStatusResponse{
			EntryID: entryID,
			Liked:   status.IsLiked,
		})
	}

	// 返回成功响应
	response := GetEntryLikeStatusResponse{
		Status: statusResponse,
	}
	ctrl.renderService.JSONSuccess(c, response)
}

// GetEntryLikeStatus 获取日记点赞状态
// 获取指定日记的点赞状态信息，包括是否已点赞和点赞总数
//
// 参数:
//   - c: Gin上下文对象
func (ctrl *EntryController) GetEntryLikeStatus(c *gin.Context) {
	// 从路径参数获取entry_id
	entryIDStr := c.Param("entry_id")
	entryID, err := strconv.ParseInt(entryIDStr, 10, 64)
	if err != nil {
		ctrl.renderService.JSONError(c, &errors.InvalidParams)
		return
	}

	// 获取登录用户信息（可能为空，访客也可以查看点赞状态）
	loginUser, _ := ctrl.auth.GetCurrentUser(c)
	var userID int64 = 0
	if loginUser != nil {
		userID = loginUser.UserId
	}

	// 获取点赞状态
	status, err := ctrl.likeService.GetLikeStatus(userID, entryID, like.LikeTypeEntry)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", userID, "entry_id", entryID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	ctrl.renderService.JSONSuccess(c, status)
}
