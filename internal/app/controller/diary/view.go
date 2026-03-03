package diary

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/biuren/hebe/internal/app/controller/validation"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ShowDiary 显示日记详细信息页面
// 显示单个日记的详细内容，渲染日记详情页面
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 创建页面数据对象
//   - 设置活动菜单为"users"
//   - 渲染diary/diary模板
func (ctrl *DiaryController) ShowDiary(c *gin.Context) {
	data := ctrl.renderService.CreatePageData(c, "diary.detail.title")
	data.Data["Active"] = "users"
	ctrl.renderService.HTML(c, "diary/diary", data)
}

// ShowDiaryForm 显示日记本创建/编辑表单页面
// 根据URL参数中的action和diary_id决定是创建新日记本还是编辑现有日记本
//
// URL参数:
//   - action: 操作类型，"diary.add"表示创建，"diary.edit"表示编辑
//   - diary_id: 当action为edit时需要提供的日记本ID
func (ctrl *DiaryController) ShowDiaryForm(c *gin.Context) {
	// 创建页面对象
	data := ctrl.renderService.CreatePageData(c, "diary.form.title")

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	if action == "diary.edit" {
		// 获取当前登录用户信息
		loginUser := ctrl.auth.MustGetCurrentUser(c)
		// 获取请求diary_id
		req := &ShowDiaryRequest{}
		if err := c.ShouldBindQuery(req); err != nil {
			ctrl.renderService.RenderErrorPage(c, errors.InvalidParams)
			return
		}

		// 编辑模式：获取并验证日记本数据
		diaryData, err := ctrl.diaryService.GetDiaryByID(req.DiaryID)
		if err != nil {
			ctrl.renderService.RenderErrorPage(c, errors.DiaryNotFound)
			return
		}

		// 检查日记是否属于当前用户
		if diaryData.AddUserID != loginUser.UserId {
			ctrl.renderService.RenderErrorPage(c, errors.PermissionDenied)
			return
		}

		data.Data["Diary"] = diaryData
	}
	ctrl.renderService.HTML(c, "form/diary", data)
}

// ShowUserProfile 显示用户个人主页
// 显示指定用户的日记本列表，支持分页显示
// 根据用户ID获取用户信息和日记本列表，并生成分页链接
//
// 参数:
//   - c: Gin上下文对象
//
// URL格式:
//   - /u/{uid} - 显示第一页
//   - /u/{uid}/{pageindex} - 显示指定页
//
// 功能:
//   - 解析URL中的用户ID和页码
//   - 验证用户是否为页面所有者
//   - 获取用户的日记本列表（分页）
//   - 生成分页导航链接
//   - 渲染用户主页模板
func (ctrl *DiaryController) ShowUserProfile(c *gin.Context) {
	// 获取url中的uid
	uid, err := validation.ParseID(c.Param("uid")) // 使用统一的ParseID方法
	if err != nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	// 创建页面数据
	data := ctrl.renderService.CreatePageData(c, "user.profile.title")

	// 判断是否为页面所有者
	isOwner := false
	if data.IsAuthenticated() {
		isOwner = (data.GetAuthUser().UserId == uid)
	}

	// 获取加入的日记本（包括待确认和已加入，按状态和时间排序）
	// 性能优化：只有 Owner 才会显示加入的日记本，避免多余查询
	if isOwner {
		// 使用 limit + 1 的方式来判断是否还有更多，避免 COUNT 查询
		limit := ctrl.config.App.UserJoinedDiariesLimit
		// 如果是Owner，包含私有日记本和Pending状态；否则只包含公开日记本且Active状态
		joinedDiaries, _ := ctrl.diaryService.ListJoinedDiariesByUserIDWithLimit(uid, isOwner, limit+1)

		hasMore := false
		if len(joinedDiaries) > limit {
			hasMore = true
			joinedDiaries = joinedDiaries[:limit]
		}

		if len(joinedDiaries) > 0 {
			data.Data["JoinedDiaries"] = joinedDiaries
			data.Data["HasMoreJoinedDiaries"] = hasMore
			data.Data["IsJoinedList"] = false
			data.Data["HxTarget"] = "#joined-diaries-section"
		}
	}

	// 获取URL中的pageIndex
	pageIndexStr := c.Param("pageindex")
	pageIndex := 1
	if pageIndexStr != "" {
		if idx, _ := strconv.Atoi(strings.TrimSpace(pageIndexStr)); idx > 0 {
			pageIndex = idx
		}
	}

	// 获取分页配置
	pageSize := ctrl.config.App.PageSize

	// 性能优化：先获取用户信息，如果用户不存在则直接返回，避免后续 diaries 查询
	pageUser, err := ctrl.userService.GetUserByUserId(uid)
	if err != nil || pageUser == nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	// 获取用户的日记本列表
	diaries, hasMore, _ := ctrl.diaryService.ListDiariesByUser(uid, isOwner, pageIndex, pageSize)

	// 初始化分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		HasMore:   hasMore,
		URL:       fmt.Sprintf("/u/%d/p/%%d", uid),
		PagePrev:  "<<",
		PageNext:  ">>",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 用户配置的时区
	if pageUser.TimeZone != "" {
		loc, err := time.LoadLocation(pageUser.TimeZone)
		if err == nil {
			// 在模板中显示时间时需要使用
			data.Data["Location"] = loc
		}
	}

	data.Data["IsOwner"] = isOwner
	data.Data["User"] = pageUser              //被访问用户
	data.Data["TimeZone"] = pageUser.TimeZone // 被访问用户时区
	data.Data["Diaries"] = diaries
	data.Data["DiariesNum"] = len(diaries) // 使用返回的条目数作为参考
	data.Data["PageSize"] = pageSize
	data.Data["PageIndex"] = pageIndex
	data.Data["PageLinks"] = pageLinks
	data.Data["HxTarget"] = "#diary-container"

	// 检查是否为HTMX请求
	if c.GetHeader("HX-Request") == "true" {
		ctrl.renderService.HTML(c, "user/list", data)
		return
	}

	ctrl.renderService.HTML(c, "user/index", data)
}

// ShowLatestPublicEntries 显示用户最新公开日记页面
//
// URL格式:
//   - /u/{uid}/latest
func (ctrl *DiaryController) ShowLatestPublicEntries(c *gin.Context) {
	uid, err := validation.ParseID(c.Param("uid"))
	if err != nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.profile.latest_public_diary")

	// 根据uid 获取用户信息
	pageUser, err := ctrl.userService.GetUserByUserId(uid)
	if err != nil || pageUser == nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}
	data.Data["User"] = pageUser
	data.Data["TimeZone"] = pageUser.TimeZone // 被访问用户时区

	// Calculate IsOwner
	isOwner := false
	if data.IsAuthenticated() {
		authUser := data.GetAuthUser()
		data.Data["CurrentUserID"] = authUser.UserId
		if authUser.UserId == pageUser.UserId {
			isOwner = true
		}
	}
	data.Data["IsOwner"] = isOwner

	// 获取最新公开日记
	latestEntriesLimit := ctrl.config.App.UserLatestPublicEntriesLimit
	latestEntries, err := ctrl.entryService.GetLatestPublicEntries(uid, latestEntriesLimit)
	if err != nil {
		latestEntries = []entry.Entry{}
	}
	data.Data["LatestPublicEntries"] = latestEntries

	ctrl.renderService.HTML(c, "user/latest_public_entries_page", data)
}

// ShowJoinedDiaries 显示用户加入的日记本列表（完整版）
// 显示用户参与的所有日记本，支持分页
//
// 参数:
//   - c: Gin上下文对象
//
// URL格式:
//   - /u/{uid}/joined - 显示第一页
//   - /u/{uid}/joined/{pageindex} - 显示指定页
func (ctrl *DiaryController) ShowJoinedDiaries(c *gin.Context) {
	uid, err := validation.ParseID(c.Param("uid"))
	if err != nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	data := ctrl.renderService.CreatePageData(c, "user.joined_diaries.title")

	// 判断是否为页面所有者
	isOwner := false
	if data.IsAuthenticated() {
		isOwner = (data.GetAuthUser().UserId == uid)
	}

	pageIndexStr := c.Param("pageindex")
	pageIndex := 1
	if pageIndexStr != "" {
		if idx, _ := strconv.Atoi(strings.TrimSpace(pageIndexStr)); idx > 0 {
			pageIndex = idx
		}
	}

	pageSize := ctrl.config.App.PageSize
	diaries, hasMore, _ := ctrl.diaryService.ListJoinedDiariesByUserIDPaged(uid, isOwner, pageIndex, pageSize)

	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		HasMore:   hasMore,
		URL:       fmt.Sprintf("/u/%d/joined/p/%%d", uid),
		PagePrev:  "<<",
		PageNext:  ">>",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	data.Data["Diaries"] = diaries
	data.Data["DiariesNum"] = len(diaries) // 使用返回的条目数作为参考
	data.Data["PageLinks"] = pageLinks
	data.Data["IsJoinedList"] = true
	data.Data["HxTarget"] = "#joined-diaries-content"
	data.Data["IsJoinedDiariesPage"] = true

	// Get User info for header/navigation if not present
	if data.Data["User"] == nil {
		user, err := ctrl.userService.GetUserByUserId(uid)
		if err == nil {
			data.Data["User"] = user
		}
	}

	if c.GetHeader("HX-Request") == "true" {
		ctrl.renderService.HTML(c, "user/joined_diaries_content", data)
		return
	}

	ctrl.renderService.HTML(c, "user/joined_diaries_full", data)
}
