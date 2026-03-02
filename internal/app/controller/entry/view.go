package entry

import (
	stdErrors "errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/app/controller/validation"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ShowEntryForm 显示日记编辑表单页面
// 支持创建新日记和编辑现有日记两种模式
//
// 支持的URL参数:
//   - action: "entry.create" 或 "entry.edit"
//   - entry_id: 当action为edit时需要提供的日记ID
func (ctrl *EntryController) ShowEntryForm(c *gin.Context) {
	// 获取当前登录用户信息，如果未登录会panic（由中间件保证已登录）
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	var entryData *entry.Entry

	if action == "entry.edit" {
		// 编辑模式：获取并验证日记数据
		var req ShowEntryFormRequest
		var err error
		if err = c.ShouldBindQuery(&req); err != nil {
			ctrl.renderService.RenderErrorPage(c, errors.ValidationError)
			return
		}
		//获取日记信息
		entryData, err = ctrl.entryService.GetEntryByIDForUser(req.EntryID, loginUser.UserId)
		if err != nil {
			ctrl.renderService.RenderErrorPage(c, errors.EntryNotFound)
			return
		}
		// 权限检查已在 GetEntryByIDForUser 中处理 (Owner/Admin/Author 均可访问)

	}
	// 初始化页面数据
	data := ctrl.renderService.CreatePageData(c, "entry.form.title")
	data.Data["Entry"] = entryData
	data.Data["DiaryID"] = c.Query("diary_id")

	ctrl.renderService.HTML(c, "form/entry", data)
}

// ShowDrafts 显示草稿箱页面
// 路由: GET /drafts
func (ctrl *EntryController) ShowDrafts(c *gin.Context) {
	loginUser := ctrl.auth.MustGetCurrentUser(c)
	if loginUser == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 获取分页参数
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if pageIndex < 1 {
		pageIndex = 1
	}
	pageSize := ctrl.config.App.PageSize

	// 获取状态参数，默认为草稿
	statusStr := c.DefaultQuery("status", "2") // 2=Draft
	statusInt, _ := strconv.Atoi(statusStr)
	status := int8(statusInt)

	// 获取列表 - 使用 hasMore 模式
	entries, hasMore, err := ctrl.entryService.ListAllEntriesByUserID(loginUser.UserId, status, pageIndex, pageSize)
	if err != nil {
		ctrl.renderService.RenderErrorPage(c, errors.InternalError)
		return
	}

	// 构建分页信息 - 使用 hasMore 模式
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		HasMore:   hasMore,
		URL:       fmt.Sprintf("/drafts?status=%d&page=%%d", status),
		PagePrev:  "<<",
		PageNext:  ">>",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 准备渲染数据
	data := ctrl.renderService.CreatePageData(c, "entry.drafts.title")
	data.Data["Entries"] = entries
	data.Data["Total"] = len(entries) // hasMore 模式下使用返回的条目数作为参考
	data.Data["CurrentPage"] = pageIndex
	data.Data["PageSize"] = pageSize
	data.Data["Status"] = status
	data.Data["PageLinks"] = pageLinks
	data.Data["Tab"] = "drafts" // 用于前端高亮
	data.Data["HxTarget"] = "#drafts-list"

	if status == entry.StatusDeleted {
		data.Data["Tab"] = "trash"
	}

	// HTMX 请求：如果是请求列表部分，只渲染列表片段
	if c.GetHeader("HX-Request") == "true" && c.GetHeader("HX-Target") == "drafts-list" {
		ctrl.renderService.HTML(c, "diary/drafts_list", data)
		return
	}

	// 完整页面渲染
	ctrl.renderService.HTML(c, "diary/drafts", data)
}

// ShowDiaryIndex 显示日记本列表页面
// 显示指定日记本下的日记列表，支持分页显示
//
// URL格式:
//   - /diary/{id} - 显示第一页
//   - /diary/{id}/{pageindex} - 显示指定页
func (ctrl *EntryController) ShowDiaryIndex(c *gin.Context) {
	// 解析日记本ID
	diaryIDStr := c.Param("id")
	diaryID, err := validation.ParseID(diaryIDStr)
	if err != nil {
		ctrl.renderService.RenderErrorPage(c, errors.InvalidParams)
		return
	}

	// 获取页码参数
	pageIndexStr := c.Param("pageindex")
	pageIndex := 1
	if pageIndexStr != "" {
		if idx, err := strconv.Atoi(strings.TrimSpace(pageIndexStr)); err == nil && idx > 0 {
			pageIndex = idx
		}
	}

	// 获取分页配置
	pageSize := ctrl.config.App.PageSize

	// 获取当前登录用户ID（用于权限检查）
	var currentUserID int64
	isAuth, loginUser := ctrl.auth.IsAuthenticated(c)
	if isAuth {
		currentUserID = loginUser.UserId
	}

	// 获取日记列表（使用ListEntriesByDiaryID以支持多人日记）
	// 同时返回日记本信息
	diaryInfo, entries, hasMore, err := ctrl.entryService.ListEntriesByDiaryID(diaryID, currentUserID, entry.StatusActive, pageIndex, pageSize)
	if err != nil {
		if stdErrors.Is(err, errors.DiaryNotFound) {
			ctrl.renderService.RenderErrorPage(c, errors.DiaryNotFound)
			return
		}
		// 如果是权限错误，返回403
		if stdErrors.Is(err, errors.PermissionDenied) {
			ctrl.renderService.RenderErrorPage(c, errors.PermissionDenied)
			return
		}
		ctrl.renderService.RenderErrorPage(c, errors.InternalError)
		return
	}

	// 构建分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		HasMore:   hasMore,
		URL:       fmt.Sprintf("/diary/%d/p/%%d", diaryID),
		PagePrev:  "<<",
		PageNext:  ">>",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 判断当前用户是否为日记本所有者
	isOwner := isAuth && loginUser.UserId == diaryInfo.AddUserID

	// 判断是否有写入权限（Owner 或 Active Author）
	canWrite := isOwner
	if isAuth {
		if !canWrite {
			// 检查是否为活跃作者
			isAuthor, _ := ctrl.diaryService.IsActiveAuthor(diaryID, loginUser.UserId)
			if isAuthor {
				canWrite = true
			}
		}
	}

	// 获取用户时区
	timeZone, err := ctrl.userService.GetUserTimeZone(diaryInfo.AddUserID)
	if err != nil {
		ctrl.renderService.RenderErrorPage(c, errors.InternalError)
		return
	}

	// 构建页面数据
	data := ctrl.renderService.CreatePageData(c, "diary.list.title")
	data.Title = diaryInfo.Title + " - " + diaryInfo.AddUserDisplayName
	data.Data["Diary"] = diaryInfo
	data.Data["Entries"] = entries
	data.Data["PageLinks"] = pageLinks
	data.Data["EntryNum"] = len(entries) // 使用返回的条目数作为参考
	data.Data["PageSize"] = pageSize
	data.Data["TimeZone"] = timeZone
	data.Data["IsOwner"] = isOwner
	data.Data["CanWrite"] = canWrite
	data.Data["CurrentUserID"] = currentUserID
	data.Data["HxTarget"] = "#diary-content"

	// 检查是否为HTMX请求
	if c.GetHeader("HX-Request") == "true" {
		ctrl.renderService.HTML(c, "diary/list", data)
		return
	}

	ctrl.renderService.HTML(c, "diary/index", data)
}

// ShowEntry 显示单个日记详情页面
// 根据URL参数中的entryid获取并显示日记详情，包括日记内容、所属日记本信息等
//
// URL格式:
//   - /entry/{entryid}: 显示指定ID的日记详情
//
// 参数:
//   - c: Gin上下文对象，包含请求参数和响应方法
//
// 功能:
//   - 解析并验证日记ID
//   - 获取日记详情和所属日记本信息
//   - 获取用户时区设置
//   - 渲染日记详情页面模板
func (ctrl *EntryController) ShowEntry(c *gin.Context) {
	var req ShowEntryDetailRequest
	// 绑定 uri
	if err := c.ShouldBindUri(&req); err != nil {
		ctrl.renderService.RenderErrorPage(c, errors.InvalidParams)
		return
	}
	// 绑定 query
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.renderService.RenderErrorPage(c, errors.InvalidParams)
		return
	}

	// 兼容 HTMX 请求头
	isAjax := c.GetHeader("HX-Request") == "true"

	// 获取日记信息 - 调用entry模块
	loginUser, isAuth := ctrl.auth.GetCurrentUser(c)
	var err error
	var entryInfo *entry.Entry
	if req.WithDel {
		if !isAuth || loginUser == nil {
			ctrl.renderService.RenderErrorPage(c, errors.PermissionDenied)
			return
		}
		entryInfo, err = ctrl.entryService.GetEntryByIDForUser(req.EntryID, loginUser.UserId)
		if err != nil || entryInfo == nil {
			ctrl.renderService.RenderErrorPage(c, errors.EntryNotFound)
			return
		}
	} else {
		entryInfo, err = ctrl.entryService.GetEntryByID(req.EntryID)
		if err != nil || entryInfo == nil {
			ctrl.renderService.RenderErrorPage(c, errors.EntryNotFound)
			return
		}
	}

	// 获取日记本信息 - 调用 diary 模块
	diaryInfo, err := ctrl.diaryService.GetDiaryByID(entryInfo.DiaryID)
	if err != nil || diaryInfo == nil {
		ctrl.renderService.RenderErrorPage(c, errors.DiaryNotFound)
		return
	}

	// 获取用户点赞状态
	isLiked := false
	if loginUser != nil {
		status, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, entryInfo.EntryID, like.LikeTypeEntry)
		if err == nil {
			isLiked = status.IsLiked
		} else {
			// 记录错误但不中断页面渲染，默认为未点赞
			slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "entry_id", entryInfo.EntryID)
		}
	}

	// todo: 未来判断访客的浏览权限

	// 获取用户时区
	timeZone, err := ctrl.userService.GetUserTimeZone(diaryInfo.AddUserID)
	if err != nil {
		ctrl.renderService.RenderErrorPage(c, errors.InternalError)
		return
	}

	var previousEntryID int64
	var nextEntryID int64

	// 如果包含已删除内容（回收站/草稿模式），不显示上一篇/下一篇
	if !req.WithDel {
		// 获取上一篇日记ID
		previousEntryID, err = ctrl.entryService.GetPreviousEntryByDiaryID(entryInfo.EntryID, entryInfo.DiaryID)
		if err != nil {
			// 记录错误但不中断页面渲染
			previousEntryID = 0
		}

		// 获取下一篇日记ID
		nextEntryID, err = ctrl.entryService.GetNextEntryByDiaryID(entryInfo.EntryID, entryInfo.DiaryID)
		if err != nil {
			// 记录错误但不中断页面渲染
			nextEntryID = 0
		}
	}

	// 构建页面数据
	data := ctrl.renderService.CreatePageData(c, "entry.detail.title")
	data.Title = entryInfo.Title + " - " + entryInfo.AddUserDisplayName
	data.Data["Entry"] = entryInfo
	data.Data["Diary"] = diaryInfo
	data.Data["TimeZone"] = timeZone
	data.Data["IsLiked"] = isLiked
	data.Data["PreviousEntryID"] = previousEntryID
	data.Data["NextEntryID"] = nextEntryID
	data.Data["HideNav"] = req.WithDel // 传递 HideNav 标志
	// 添加 CurrentUserID 用于判断是否显示编辑按钮
	if loginUser != nil {
		data.Data["CurrentUserID"] = loginUser.UserId
	}
	// 渲染 Markdown 内容为 HTML
	data.Data["BodyHTML"] = pkg.RenderMarkdown(entryInfo.Body)

	// 构建页面数据
	viewPath := "entry/view"
	if isAjax {
		viewPath = "entry/detail"
	}
	ctrl.renderService.HTML(c, viewPath, data)
}
