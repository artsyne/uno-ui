package post

import (
	"fmt"
	"strconv"
	"strings"

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
		var err error
		postData, err = ctrl.getAndValidatePost(c, loginUser.UserId)
		if err != nil {
			ctrl.renderService.JSONError(c, err)
			return
		}
	}
	// 创建模式：直接渲染空表单

	// 初始化页面数据
	data := ctrl.renderService.CreatePageData(c, "post.form.title")
	data.Data["Active"] = "users"
	data.Data["Post"] = postData
	// 注意：用户信息现在通过 data.Context.AuthUser 在模板中访问，无需在 Data 中设置 UserData

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
	diaryIDStr := c.Param("spaceid")
	diaryID, err := ctrl.parseDiaryID(diaryIDStr)
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
	posts, diariesNum, err := ctrl.postService.ListPostsByUserID(diaryInfo.AddUserID, diaryID, pageIndex, pageSize)
	var hasError bool
	var errMsg string

	if err != nil {
		hasError = true
		errMsg = "获取日记列表失败"
		posts = []post.Post{} // 空列表
		diariesNum = 0
	}

	// 构建分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/d/%d/%%d", diaryID),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 判断当前用户是否为日记本所有者
	isAuth, loginUser := ctrl.auth.IsAuthenticated(c)
	isOwner := isAuth && loginUser.UserId == diaryInfo.AddUserID

	// 构建页面数据
	data := ctrl.renderService.CreatePageData(c, "diary.index.title")
	data.Title = diaryInfo.Title + " - " + diaryInfo.AddUserName
	data.Data["Diary"] = diaryInfo
	data.Data["Posts"] = posts
	data.Data["PageLinks"] = pageLinks
	data.Data["hasError"] = hasError
	data.Data["errMsg"] = errMsg
	data.Data["DiariesNum"] = diariesNum
	data.Data["PageSize"] = pageSize
	data.Data["IsOwner"] = isOwner

	ctrl.renderService.HTML(c, "diary/index", data)
}

// renderErrorPage 渲染错误页面
// 当发生错误时显示用户友好的错误信息
func (ctrl *PostController) renderErrorPage(c *gin.Context, title, message string) {
	data := ctrl.renderService.CreatePageData(c, "error")
	data.Data["title"] = title
	data.Data["message"] = message
	ctrl.renderService.HTML(c, "error", data)
}