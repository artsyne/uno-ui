package diary

import (
	"encoding/json"
	"fmt"
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
		DiaryID:       0, // 将由数据库生成
		AccessPwd:     req.AccessPwd,
		AccessURL:     "",
		CoverURL:      "",
		CoverMidURL:   "",
		LikeNum:       0,
		MemberNum:     0,
		ViewNum:       0,
		InfoNum:       0,                  // 日记数量初始为0
		Status:        diary.StatusActive, // 默认状态为正常
		AddUTCTime:    nowTime,
		AddUserID:     loginUser.UserId,
		AddUserHandle: loginUser.Handle,
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

	// 调用服务层更新日记本
	// 只有拥有者可以更新
	isOwner, err := ctrl.diaryService.IsDiaryOwner(req.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	if !isOwner {
		ctrl.renderService.JSONError(c, &errors.PermissionDenied)
		return
	}

	// 构建更新对象
	// 注意：这里我们重新从数据库获取一次，或者只更新允许的字段
	// 简单起见，我们构造一个Diary对象传给UpdateDiary，服务层会处理
	// 但服务层UpdateDiary需要完整对象吗？
	// 检查服务层实现：UpdateDiary 调用 repo.Update
	// repo.Update 是全量更新。
	// 所以最好先获取旧对象，再覆盖字段。

	oldDiary, err := ctrl.diaryService.GetDiaryByID(req.DiaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	oldDiary.Title = req.Title
	oldDiary.Remark = req.Remark
	oldDiary.ViewSetting = int8(req.ViewSetting)
	oldDiary.SocialSetting = int8(req.SocialSetting)
	oldDiary.ThemeID = int8(req.ThemeID)
	oldDiary.AccessPwd = req.AccessPwd

	err = ctrl.diaryService.UpdateDiary(oldDiary)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 如果是HTMX请求，返回触发事件头
	if c.GetHeader("HX-Request") == "true" {
		// 构造触发事件的数据
		// diaryUpdated: 更新前端日记本UI展示 (对应 diary-detail.js 中的 diaryUpdated 事件监听)
		// closeDiaryModal: 关闭编辑弹窗
		// showToast: 显示成功提示
		triggerData := map[string]interface{}{
			"diaryUpdated":    oldDiary,
			"closeDiaryModal": true,
			"showToast": map[string]string{
				"type":    "success",
				"message": "common.save_success",
			},
		}

		// 使用 MarshalASCII 确保 JSON 中的非 ASCII 字符（如中文）被转义为 \uXXXX
		// 这样可以安全地放入 HTTP 响应头中，避免乱码问题
		triggerJSON, err := pkg.MarshalASCII(triggerData)
		if err == nil {
			c.Header("HX-Trigger", string(triggerJSON))
		}
		// 直接返回 204 No Content，既能避免 ViewService.JSONSuccess 添加默认的 "success" 消息导致前端显示重复 Toast，
		// 又能避开 auth.js 中 htmx:beforeSwap 对 JSON 响应的自动 Toast 处理
		c.Status(http.StatusNoContent)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// DeleteDiary 删除日记本
// 路由: DELETE /diary/:id
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[DeleteDiaryRequest](c)

	err := ctrl.diaryService.DeleteDiary(req.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// LikeDiary 点赞/取消点赞日记本
// 路由: POST /api/diary/like
func (ctrl *DiaryController) LikeDiary(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[LikeDiaryRequest](c)

	// 调用点赞服务
	// 这里的 likeService 是通用的点赞服务
	// 参数: objectType (diary=1), objectID, userID
	// 返回: isLiked (true=点赞成功, false=取消点赞成功), likeCount
	err := ctrl.likeService.ToggleLike(loginUser.UserId, req.DiaryID, like.LikeTypeDiary)
	if err != nil {
		slog.Error("点赞失败", "error", err, "user_id", loginUser.UserId, "diary_id", req.DiaryID)
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 获取最新状态
	status, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, req.DiaryID, like.LikeTypeDiary)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 如果点赞成功，更新日记本的点赞数（缓存或数据库）
	// 注意：likeService 已经更新了数据库中的点赞记录，但是日记本表中的 like_num 也需要更新
	if status.IsLiked {
		_ = ctrl.diaryService.IncrementLikeCount(req.DiaryID)
	} else {
		_ = ctrl.diaryService.DecrementLikeCount(req.DiaryID)
	}

	// 返回结果
	ctrl.renderService.JSONSuccess(c, gin.H{
		"is_liked":   status.IsLiked,
		"like_count": status.LikeCount,
	})
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

// AuthorViewModel 作者视图模型
type AuthorViewModel struct {
	diary.Author
	Role int // 1: Owner, 2: Admin, 3: Member
}

// ShowAuthors 显示作者管理页面
// 路由: GET /diary/:id/authors
func (ctrl *DiaryController) ShowAuthors(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[GetAuthorsRequest](c)

	diaryObj, err := ctrl.diaryService.GetDiaryByID(req.DiaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 检查权限：
	// 1. 如果是公开日记 (ViewSetting == diary.ViewSettingPublic)，所有人可见
	// 2. 如果是私有日记，只有拥有者或活跃作者可见
	isOwner := diaryObj.AddUserID == loginUser.UserId
	canView := false
	if diaryObj.ViewSetting == diary.ViewSettingPublic {
		canView = true
	} else {
		if isOwner {
			canView = true
		} else {
			isAuthor, err := ctrl.diaryService.IsAuthor(req.DiaryID, loginUser.UserId)
			if err != nil {
				ctrl.renderService.JSONError(c, err)
				return
			}
			if isAuthor {
				canView = true
			}
		}
	}

	if !canView {
		ctrl.renderService.RenderErrorPage(c, errors.PermissionDenied)
		return
	}

	// Pagination Logic
	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	// pageSize strictly from config, ignore frontend input
	pageSize := ctrl.config.App.PageSize

	allAuthors, totalCount, err := ctrl.diaryService.GetDiaryAuthors(req.DiaryID, !isOwner, pageIndex, pageSize)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	canQuit := !isOwner

	// 转换为视图模型并计算角色
	authorViewModels := make([]AuthorViewModel, len(allAuthors))
	currentUserRole := 0
	if isOwner {
		currentUserRole = 1 // 默认为Owner
	}

	for i, a := range allAuthors {
		role := 3 // Default Member
		if a.UserID == diaryObj.AddUserID {
			role = 1 // Owner
		} else if a.IsAdmin == 1 {
			role = 2 // Admin
		}

		authorViewModels[i] = AuthorViewModel{
			Author: a,
			Role:   role,
		}

		if a.UserID == loginUser.UserId {
			currentUserRole = role
		}
	}

	canManage := currentUserRole == 1 || currentUserRole == 2

	pageData := ctrl.renderService.CreatePageData(c, "diary.author.list_title")
	pageData.Data["Diary"] = diaryObj
	pageData.Data["Authors"] = authorViewModels
	pageData.Data["IsOwner"] = isOwner
	pageData.Data["CanQuit"] = canQuit
	pageData.Data["CanManage"] = canManage
	pageData.Data["UserRole"] = currentUserRole
	pageData.Data["Lang"] = loginUser.Lang
	pageData.Data["DiaryID"] = diaryObj.DiaryID

	// Pagination Data
	pageData.Data["TotalPages"] = totalPages
	pageData.Data["HxTarget"] = "#authors-content"
	if totalPages > 1 {
		pageData.Data["PageLinks"] = pkg.GeneratePagination(pkg.ListOptions{
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Total:     int(totalCount),
			URL:       fmt.Sprintf("/diary/%d/authors?page=%%d", req.DiaryID),
			PagePrev:  "<",
			PageNext:  ">",
		})
	}

	if c.GetHeader("HX-Request") == "true" {
		ctrl.renderService.HTML(c, "diary/authors_content", pageData)
		return
	}

	ctrl.renderService.HTML(c, "diary/authors", pageData)
}

// ShowInviteForm 显示邀请作者表单（用于弹窗）
// 路由: GET /diary/:id/authors/invite
func (ctrl *DiaryController) ShowInviteForm(c *gin.Context) {
	// 鉴权在路由层处理，确保是作者且有权限邀请（Admin/Owner）
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[GetAuthorsRequest](c) // 复用GetAuthorsRequest因为只需要DiaryID

	// 再次检查权限，确保是 Admin 或 Owner
	// 这里可以调用 Service 检查
	role, err := ctrl.diaryService.GetAuthorRole(req.DiaryID, loginUser.UserId)
	if err != nil {
		lang := ctrl.renderService.GetLocaleFromContext(c)
		msg := ctrl.renderService.GetMessage("error.message", lang)
		c.String(http.StatusOK, fmt.Sprintf(`<div class="text-center text-red-500 py-8">%s</div>`, msg))
		return
	}
	if role != 1 && role != 2 { // 1: Owner, 2: Admin
		lang := ctrl.renderService.GetLocaleFromContext(c)
		msg := ctrl.renderService.GetMessage("error.permission_denied", lang)
		c.String(http.StatusOK, fmt.Sprintf(`<div class="text-center text-red-500 py-8">%s</div>`, msg))
		return
	}

	data := map[string]interface{}{
		"DiaryID": req.DiaryID,
		"Lang":    loginUser.Lang,
	}

	// 只渲染表单部分
	ctrl.renderService.HTML(c, "diary/box_invite_author", data)
}

// InviteAuthor 邀请作者
// 路由: POST /diary/:id/authors
func (ctrl *DiaryController) InviteAuthor(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[InviteAuthorRequest](c)

	err := ctrl.diaryService.InviteAuthor(req.DiaryID, loginUser.UserId, req.InviteeIdentity, loginUser.Lang)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		triggerData := map[string]interface{}{
			"closeAuthorInviteModal": true,
			"reloadAuthorList":       true,
			"showToast": map[string]string{
				"type":    "success",
				"message": "diary.author.invite_success",
			},
		}

		if triggerJSON, err := json.Marshal(triggerData); err == nil {
			c.Header("HX-Trigger", string(triggerJSON))
		}
		// Return 204 No Content to avoid JSON processing in auth.js
		// and strictly rely on HX-Trigger for events
		c.Status(http.StatusNoContent)
		return
	}

	ctrl.renderService.JSONSuccessWithMsg(c, nil, "diary.author.invite_success")
}

// RemoveAuthor 移除作者
// 路由: DELETE /diary/:id/authors/:user_id
func (ctrl *DiaryController) RemoveAuthor(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[RemoveAuthorRequest](c)

	err := ctrl.diaryService.RemoveAuthor(req.DiaryID, loginUser.UserId, req.TargetUserID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// UpdateAuthor 更新作者权限
// 路由: PUT /diary/:id/authors/:user_id
func (ctrl *DiaryController) UpdateAuthor(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[UpdateAuthorRequest](c)

	err := ctrl.diaryService.UpdateAuthorAdmin(req.DiaryID, loginUser.UserId, req.TargetUserID, req.IsAdmin)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// QuitDiary 退出日记本
// 路由: DELETE /diary/:id/authors/self
func (ctrl *DiaryController) QuitDiary(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[QuitDiaryRequest](c)

	err := ctrl.diaryService.QuitDiary(req.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.Status(http.StatusOK)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// ShowInvitation 显示邀请页面
// 路由: GET /diary/invitation
func (ctrl *DiaryController) ShowInvitation(c *gin.Context) {
	req := middleware.MustGetValidatedRequest[GetInvitationRequest](c)

	// 这里我们只是简单的显示页面，页面上应该有 accept/reject 按钮
	// 实际逻辑中可能需要校验 invitee 是否是当前登录用户，或者 token 是否有效
	// 简单起见，我们假设用户点击链接后，如果未登录会跳转登录，登录后看到此页面

	loginUser := ctrl.auth.MustGetCurrentUser(c)

	diary, err := ctrl.diaryService.GetDiaryByID(req.DiaryID)
	if err != nil {
		// 处理日记本不存在的情况
		ctrl.renderService.RenderErrorPage(c, errors.DiaryNotFound)
		return
	}

	pageData := ctrl.renderService.CreatePageDataWithData(c, "diary.invitation.title", map[string]interface{}{
		"Diary": diary,
		"User":  loginUser,
	})
	// 渲染邀请页面模板
	ctrl.renderService.HTML(c, "diary/invitation", pageData)
}

// AcceptInvitation 接受邀请
// 路由: POST /diary/invitation/accept
func (ctrl *DiaryController) AcceptInvitation(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[HandleInvitationRequest](c)

	err := ctrl.diaryService.AcceptInvitation(req.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		// 获取日记本信息以便渲染卡片
		diaryData, err := ctrl.diaryService.GetDiaryByID(req.DiaryID)
		if err != nil {
			ctrl.renderService.JSONError(c, err)
			return
		}
		// 手动设置状态为 Active (1)，因为 GetDiaryByID 不包含 AuthorStatus，且我们刚接受了邀请
		diaryData.AuthorStatus = diary.AuthorStatusActive
		// 设置 JoinUTCTime 为当前时间，虽然不是精确的加入时间，但对于显示顺序（如果按时间排）可能有用
		// 不过这里主要是为了显示卡片，卡片上不显示加入时间
		// 注意：JoinedDiaryCard 模板需要的 User 字段在 Diary.AddUser* 中

		data := map[string]interface{}{
			"Diary": diaryData,
			"Lang":  loginUser.Lang,
		}
		c.HTML(http.StatusOK, "user/joined_diary_card", data)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// RejectInvitation 拒绝邀请
// 路由: POST /diary/invitation/reject
func (ctrl *DiaryController) RejectInvitation(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	req := middleware.MustGetValidatedRequest[HandleInvitationRequest](c)

	err := ctrl.diaryService.RejectInvitation(req.DiaryID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		// 拒绝邀请后，卡片应该从列表中移除
		c.Status(http.StatusOK)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}
