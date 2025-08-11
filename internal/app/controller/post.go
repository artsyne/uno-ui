package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/utils/config"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)



// postController 日记模块HTTP处理器实现
// 负责处理所有与日记相关的HTTP请求，包括渲染表单、创建日记、编辑日记等操作
type PostController struct {
	postService   *post.Service
	diaryService  *diary.Service
	auth          *user.AuthHelper
	config        *config.Config
	renderService *view.ViewService
}

// NewPostController 创建新的日记控制器实例
func NewPostController(postService *post.Service, diaryService *diary.Service, auth *user.AuthHelper, config *config.Config, renderService *view.ViewService) *PostController {
	return &PostController{
		postService:   postService,
		diaryService:  diaryService,
		auth:          auth,
		config:        config,
		renderService: renderService,
	}
}

// RegisterRoutes 注册所有日记相关路由
// 将所有post相关路由聚合在一个地方，实现高内聚
func (ctrl *PostController) RegisterRoutes(router *gin.Engine) {
	// 日记本页面路由组（无需鉴权，支持公开访问）
	// 显示日记本下的日记列表
	daybook := router.Group("/d")
	{
		daybook.GET("/:spaceid", ctrl.RenderDiaryIndex)
		daybook.GET("/:spaceid/:pageindex", ctrl.RenderDiaryIndex)
	}

	// 需要鉴权的日记表单路由
	authForm := router.Group("/form")
	authForm.Use(ctrl.auth.RequireAuth())
	{
		authForm.GET("/post", ctrl.RenderPostForm)
	}

	// 需要鉴权的日记数据操作路由（保持原路径）
	authAPI := router.Group("/")
	authAPI.Use(ctrl.auth.RequireAuth())
	{
		authAPI.POST("/post", ctrl.HandlePost)
		authAPI.POST("/info", ctrl.HandlePost)
	}
}

// parsePostID 解析并验证日记ID
// 从字符串中解析出有效的日记ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的ID
//   - error: 解析失败时的错误信息
func (ctrl *PostController) parsePostID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParams.Err()
	}
	return id, nil
}

// parseDiaryID 解析并验证日记本ID
// 从字符串中解析出有效的日记本ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的ID
//   - error: 解析失败时的错误信息
func (ctrl *PostController) parseDiaryID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParams.Err()
	}
	return id, nil
}

// RenderPostForm 渲染日记编辑表单页面
// 支持创建新日记和编辑现有日记两种模式
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 支持的URL参数:
//   - action: "post.create" 或 "post.edit"
//   - post_id: 当action为edit时需要提供的日记ID
func (ctrl *PostController) RenderPostForm(c *gin.Context) {
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

// HandlePost 处理日记的创建和编辑请求
// 根据表单中的action字段分发到对应的处理方法
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 支持的表单action值:
//   - "post.create": 创建新日记
//   - "post.edit": 编辑现有日记
func (ctrl *PostController) HandlePost(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析表单数据到PostRequest结构体
	postForm := &post.PostRequest{}
	if err := c.Bind(postForm); err != nil {
		ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
		return
	}

	// 根据操作类型分发处理
	action, _ := c.GetPostForm("action")
	switch action {
	case "post.create":
		ctrl.handleCreatePost(c, postForm, loginUser)
	case "post.edit":
		ctrl.handleEditPost(c, postForm, loginUser)
	default:
		ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
	}
}

// handleCreatePost 处理创建新日记的业务逻辑
// 设置日记的所有必要字段，调用服务层创建日记，并返回创建结果
//
// 参数:
//   - c: Gin上下文
//   - postForm: 解析后的表单数据
//   - user: 当前登录用户信息
func (ctrl *PostController) handleCreatePost(c *gin.Context, postForm *post.PostRequest, loginUser *types.AuthUser) {
	nowTime := pkg.GetNowTimeUnixMilli()

	// 构建新日记对象，设置所有必要字段
	postData := &post.Post{
		DiaryID:       postForm.DiaryID,
		Title:         postForm.Title,
		Body:          postForm.Body,
		Summary:       pkg.TruncateString(postForm.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
		IsTop:         0,
		InfoType:      1, // 默认为文本类型
		AddUTCTime:    nowTime,
		AddUserID:     loginUser.UserId,
		AddUserName:   loginUser.UserName,
		AddUserIP:     c.ClientIP(),
		UpdateUTCTime: nowTime,
		Status:        1, // 默认状态为正常
		ViewNum:       0,
		LikeNum:       0,
		CommentNum:    0,
	}

	// 调用服务层创建日记
	postID, err := ctrl.postService.AddPost(postData)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	postData.PostID = postID

	// 日记本信息数量加一 - 在controller层处理跨模块调用
	err = ctrl.diaryService.IncrementPostCount(postData.DiaryID, postData.AddUserID, postData.AddUserName, postData.AddUTCTime)
	if err != nil {
		// 如果更新日记本数量失败，记录错误但不影响日记创建的成功响应
		// 这里可以添加日志记录
	}

	// 返回成功响应
	ctrl.respondSuccess(c, postData)
}

// handleEditPost 处理编辑现有日记的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记
//
// 参数:
//   - c: Gin上下文
//   - postForm: 解析后的表单数据
//   - user: 当前登录用户信息
func (ctrl *PostController) handleEditPost(c *gin.Context, postForm *post.PostRequest, loginUser *types.AuthUser) {
	// 获取要编辑的日记ID
	postIdStr, _ := c.GetPostForm("post_id")
	postId, err := ctrl.parsePostID(postIdStr)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 验证日记所有权
	originalPost, err := ctrl.validatePostOwnership(postId, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 使用结构体字面量构建更新后的日记对象
	postData := &post.Post{
		PostID:        postId,
		DiaryID:       originalPost.DiaryID, // 保持原有日记本ID
		Title:         postForm.Title,
		Body:          postForm.Body,
		Summary:       pkg.TruncateString(postForm.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
		UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
		// 保留原有字段值
		ViewNum:    originalPost.ViewNum,
		LikeNum:    originalPost.LikeNum,
		CommentNum: originalPost.CommentNum,
		IsTop:      originalPost.IsTop,
		InfoType:   originalPost.InfoType,
	}

	// 调用服务层更新日记
	if err := ctrl.postService.UpdatePost(postData); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	ctrl.respondSuccess(c, postData)
}

// getAndValidatePost 获取并验证日记权限（用于表单渲染）
// 从URL参数中获取post_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *post.Post: 日记对象，如果验证失败返回nil
//   - error: 错误信息
func (ctrl *PostController) getAndValidatePost(c *gin.Context, userID int64) (*post.Post, error) {
	postID, err := ctrl.parsePostID(c.Query("post_id"))
	if err != nil {
		return nil, err
	}

	return ctrl.validatePostOwnership(postID, userID)
}

// validatePostOwnership 验证日记所有权
// 检查指定日记是否存在，以及当前用户是否有权限操作该日记
//
// 参数:
//   - postID: 日记ID
//   - userID: 当前用户ID
//
// 返回值:
//   - *post.Post: 日记对象
//   - error: 验证失败时的错误信息
func (ctrl *PostController) validatePostOwnership(postID, userID int64) (*post.Post, error) {
	postData, err := ctrl.postService.GetPostByID(postID)
	if err != nil || postData == nil {
		return nil, errors.InvalidParams.Err()
	}

	// 检查是否是当前用户的日记
	if postData.AddUserID != userID {
		return nil, errors.PermissionDenied.Err()
	}

	return postData, nil
}

// respondSuccess 统一的成功响应格式
// 严格保持原有的响应格式，确保前后端配合不出错
//
// 参数:
//   - c: Gin上下文
//   - postData: 操作成功的日记对象
func (ctrl *PostController) respondSuccess(c *gin.Context, postData *post.Post) {
	ctrl.renderService.JSONSuccess(c, postData)
}

// RenderDiaryIndex 渲染日记本列表页面
// 显示指定日记本下的日记列表，支持分页显示
//
// URL格式:
//   - /d/{diaryid} - 显示第一页
//   - /d/{diaryid}/{pageindex} - 显示指定页
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (ctrl *PostController) RenderDiaryIndex(c *gin.Context) {
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
//
// 参数:
//   - c: Gin上下文
//   - title: 错误标题
//   - message: 错误详细信息
func (ctrl *PostController) renderErrorPage(c *gin.Context, title, message string) {
	data := ctrl.renderService.CreatePageData(c, "error")
	data.Data["title"] = title
	data.Data["message"] = message
	ctrl.renderService.HTML(c, "error", data)
}