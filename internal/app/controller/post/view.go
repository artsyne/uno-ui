package post

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/app/controller/validation"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ShowPostForm 显示日记编辑表单页面
// 支持创建新日记和编辑现有日记两种模式
//
// 支持的URL参数:
//   - action: "post.create" 或 "post.edit"
//   - post_id: 当action为edit时需要提供的日记ID
func (ctrl *PostController) ShowPostForm(c *gin.Context) {
	// 获取当前登录用户信息，如果未登录会panic（由中间件保证已登录）
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	var postData *post.Post

	if action == "post.edit" {
		// 编辑模式：获取并验证日记数据
		var req ShowPostFormRequest
		var err error
		if err = c.ShouldBindQuery(&req); err != nil {
			ctrl.renderService.JSONError(c, errors.ValidationError)
			return
		}
		//获取日记信息
		postData, err = ctrl.postService.GetPostByID(req.PostID)
		if err != nil {
			ctrl.renderService.JSONError(c, errors.PostNotFound)
			return
		}
		// 检查日记是否属于当前用户，检查权限
		if postData.AddUserID != loginUser.UserId {
			ctrl.renderService.JSONError(c, errors.PermissionDenied)
			return
		}
	}
	// 初始化页面数据
	data := ctrl.renderService.CreatePageData(c, "post.form.title")
	data.Data["Post"] = postData

	ctrl.renderService.HTML(c, "form/post", data)
}

// ShowDiaryIndex 显示日记本列表页面
// 显示指定日记本下的日记列表，支持分页显示
//
// URL格式:
//   - /d/{diaryid} - 显示第一页
//   - /d/{diaryid}/{pageindex} - 显示指定页
func (ctrl *PostController) ShowDiaryIndex(c *gin.Context) {
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
	posts, postNum, err := ctrl.postService.ListPostsByUserID(diaryInfo.AddUserID, diaryID, pageIndex, pageSize)
	if err != nil {
		ctrl.renderErrorPage(c,"日记信息错误", "获取日记列表失败")
		return
	}

	// 构建分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     postNum,
		URL:       fmt.Sprintf("/d/%d/%%d", diaryID),
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
	data := ctrl.renderService.CreatePageData(c, "diary.index.title")
	data.Title = diaryInfo.Title + " - " + diaryInfo.AddUserName
	data.Data["Diary"] = diaryInfo
	data.Data["Posts"] = posts
	data.Data["PageLinks"] = pageLinks
	data.Data["PostNum"] = postNum
	data.Data["PageSize"] = pageSize
	data.Data["TimeZone"] = timeZone
	data.Data["IsOwner"] = isOwner

	ctrl.renderService.HTML(c, "diary/index", data)
}

// ShowPost 显示单个日记详情页面
// 根据URL参数中的postid获取并显示日记详情，包括日记内容、所属日记本信息等
//
// URL格式:
//   - /p/{postid}: 显示指定ID的日记详情
//
// 参数:
//   - c: Gin上下文对象，包含请求参数和响应方法
//
// 功能:
//   - 解析并验证日记ID
//   - 获取日记详情和所属日记本信息
//   - 获取用户时区设置
//   - 渲染日记详情页面模板
func (ctrl *PostController) ShowPost(c *gin.Context) {
	var req ShowPostDetailRequest
	// 绑定 uri
	if err:= c.ShouldBindUri(&req); err!= nil {
		ctrl.renderErrorPage(c, "参数错误", "请检查URL参数是否正确")
		return
	}
	// 绑定 query
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.renderErrorPage(c, "参数错误", "请检查URL参数是否正确")
		return
	}

	// 获取日记信息 - 调用post模块
	postInfo, err := ctrl.postService.GetPostByID(req.PostID)
	if err != nil || postInfo == nil {
		ctrl.renderErrorPage(c, "日记不存在", "未找到指定的日记")
		return
	}

	// 获取日记本信息 - 调用 diary 模块
	diaryInfo, err := ctrl.diaryService.GetDiaryByID(postInfo.DiaryID)
	if err != nil || diaryInfo == nil {
		ctrl.renderErrorPage(c, "日记本不存在", "未找到指定的日记本")
		return
	}

	// todo: 未来判断访客的浏览权限

	// 获取用户时区
	timeZone, err := ctrl.userService.GetUserTimeZone(diaryInfo.AddUserID)
	if err != nil {
		ctrl.renderErrorPage(c, "时区获取失败", "获取用户时区失败")
		return
	}

	// 获取上一篇日记ID
	previousPostID, err := ctrl.postService.GetPreviousPostByDiaryID(postInfo.PostID, postInfo.DiaryID)
	if err != nil {
		// 记录错误但不中断页面渲染
		previousPostID = 0
	}

	// 获取下一篇日记ID
	nextPostID, err := ctrl.postService.GetNextPostByDiaryID(postInfo.PostID, postInfo.DiaryID)
	if err != nil {
		// 记录错误但不中断页面渲染
		nextPostID = 0
	}

	// 构建页面数据
	data := ctrl.renderService.CreatePageData(c, "post.show.title")
	data.Title = postInfo.Title + " - " + postInfo.AddUserName
	data.Data["Post"] = postInfo
	data.Data["Diary"] = diaryInfo
	data.Data["TimeZone"] = timeZone
	data.Data["IsAjax"] = req.IsAjax
	data.Data["PreviousPostID"] = previousPostID
	data.Data["NextPostID"] = nextPostID

	// 构建页面数据
	viewPath:="post/view"
	if req.IsAjax {
		viewPath="post/detail"
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
func (ctrl *PostController) renderErrorPage(c *gin.Context, title, message string) {
	data := ctrl.renderService.CreatePageData(c, "error")
	data.Data["title"] = title
	data.Data["message"] = message
	ctrl.renderService.HTML(c, "error", data)
}