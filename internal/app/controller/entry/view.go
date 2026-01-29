package entry

import (
	"fmt"
	"log/slog"
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
			ctrl.renderService.JSONError(c, errors.ValidationError)
			return
		}
		//获取日记信息
		entryData, err = ctrl.entryService.GetEntryByID(req.EntryID)
		if err != nil {
			ctrl.renderService.JSONError(c, errors.EntryNotFound)
			return
		}
		// 检查日记是否属于当前用户，检查权限
		if entryData.AddUserID != loginUser.UserId {
			ctrl.renderService.JSONError(c, errors.PermissionDenied)
			return
		}
	}
	// 初始化页面数据
	data := ctrl.renderService.CreatePageData(c, "entry.form.title")
	data.Data["Entry"] = entryData
	data.Data["DiaryID"] = c.Query("diary_id")

	ctrl.renderService.HTML(c, "form/entry", data)
}

// ShowDiaryIndex 显示日记本列表页面
// 显示指定日记本下的日记列表，支持分页显示
//
// URL格式:
//   - /diary/{diaryid} - 显示第一页
//   - /diary/{diaryid}/{pageindex} - 显示指定页
func (ctrl *EntryController) ShowDiaryIndex(c *gin.Context) {
	// 解析日记本ID
	diaryIDStr := c.Param("diaryid")
	diaryID, err := validation.ParseID(diaryIDStr)
	if err != nil {
		ctrl.renderErrorPage(c, "日记本ID格式错误", "请检查URL中的日记本ID是否正确")
		return
	}

	// 获取日记本信息 - 调用diary模块
	diaryInfo, err := ctrl.diaryService.GetDiaryByID(diaryID)
	if err != nil || diaryInfo == nil {
		ctrl.renderErrorPage(c, "日记本不存在", "未找到指定的日记本")
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

	// 获取日记列表
	entries, entryNum, err := ctrl.entryService.ListEntriesByUserID(diaryInfo.AddUserID, diaryID, pageIndex, pageSize)
	if err != nil {
		ctrl.renderErrorPage(c, "日记信息错误", "获取日记列表失败")
		return
	}

	// 构建分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     entryNum,
		URL:       fmt.Sprintf("/diary/%d/%%d", diaryID),
		PagePrev:  "<<",
		PageNext:  ">>",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 判断当前用户是否为日记本所有者
	isAuth, loginUser := ctrl.auth.IsAuthenticated(c)
	isOwner := isAuth && loginUser.UserId == diaryInfo.AddUserID

	// 获取用户时区
	timeZone, err := ctrl.userService.GetUserTimeZone(diaryInfo.AddUserID)
	if err != nil {
		ctrl.renderErrorPage(c, "时区获取失败", "获取用户时区失败")
		return
	}

	// 构建页面数据
	data := ctrl.renderService.CreatePageData(c, "diary.list.title")
	data.Title = diaryInfo.Title + " - " + diaryInfo.AddUserDisplayName
	data.Data["Diary"] = diaryInfo
	data.Data["Entries"] = entries
	data.Data["PageLinks"] = pageLinks
	data.Data["EntryNum"] = entryNum
	data.Data["PageSize"] = pageSize
	data.Data["TimeZone"] = timeZone
	data.Data["IsOwner"] = isOwner

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
		ctrl.renderErrorPage(c, "参数错误", "请检查URL参数是否正确")
		return
	}
	// 绑定 query
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.renderErrorPage(c, "参数错误", "请检查URL参数是否正确")
		return
	}

	// 兼容 HTMX 请求头
	isAjax := c.GetHeader("HX-Request") == "true"

	// 获取日记信息 - 调用entry模块
	entryInfo, err := ctrl.entryService.GetEntryByID(req.EntryID)
	if err != nil || entryInfo == nil {
		ctrl.renderErrorPage(c, "日记不存在", "未找到指定的日记")
		return
	}

	// 获取日记本信息 - 调用 diary 模块
	diaryInfo, err := ctrl.diaryService.GetDiaryByID(entryInfo.DiaryID)
	if err != nil || diaryInfo == nil {
		ctrl.renderErrorPage(c, "日记本不存在", "未找到指定的日记本")
		return
	}

	// 获取用户点赞状态
	isLiked := false
	if loginUser, exists := ctrl.auth.GetCurrentUser(c); exists && loginUser != nil {
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
		ctrl.renderErrorPage(c, "时区获取失败", "获取用户时区失败")
		return
	}

	// 获取上一篇日记ID
	previousEntryID, err := ctrl.entryService.GetPreviousEntryByDiaryID(entryInfo.EntryID, entryInfo.DiaryID)
	if err != nil {
		// 记录错误但不中断页面渲染
		previousEntryID = 0
	}

	// 获取下一篇日记ID
	nextEntryID, err := ctrl.entryService.GetNextEntryByDiaryID(entryInfo.EntryID, entryInfo.DiaryID)
	if err != nil {
		// 记录错误但不中断页面渲染
		nextEntryID = 0
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
	// 渲染 Markdown 内容为 HTML
	data.Data["BodyHTML"] = pkg.RenderMarkdown(entryInfo.Body)

	// 构建页面数据
	viewPath := "entry/view"
	if isAjax {
		viewPath = "entry/detail"
	}
	ctrl.renderService.HTML(c, viewPath, data)
}

// renderErrorPage 渲染错误页面
// 当发生错误时显示用户友好的错误信息，统一错误页面的展示格式
//
// 参数:
//   - c: Gin上下文对象
//   - title: 错误标题
//   - message: 错误详细信息
//
// 功能:
//   - 创建错误页面数据结构
//   - 设置错误标题和消息
//   - 渲染统一的错误页面模板
func (ctrl *EntryController) renderErrorPage(c *gin.Context, title, message string) {
	data := ctrl.renderService.CreatePageData(c, "error")
	data.Data["title"] = title
	data.Data["message"] = message
	ctrl.renderService.HTML(c, "error", data)
}
