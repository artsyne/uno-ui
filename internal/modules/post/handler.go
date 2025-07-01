package post

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/modules/diary"
	"github.com/biuren/hebe/internal/modules/user"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// Handler 日记模块HTTP处理器
// 负责处理所有与日记相关的HTTP请求，包括渲染表单、创建日记、编辑日记等操作
type Handler struct {
	service *Service
	auth    *user.AuthHelper
}

// NewHandler 创建新的日记处理器实例
func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
		auth:    user.NewAuthHelper(),
	}
}

// RegisterRoutes 注册所有日记相关路由
// 将所有post相关路由聚合在一个地方，实现高内聚
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// 日记本页面路由组（无需鉴权，支持公开访问）
	// 显示日记本下的日记列表
	daybook := router.Group("/d")
	{
		daybook.GET("/:spaceid", h.RenderDiaryIndex)
		daybook.GET("/:spaceid/:pageindex", h.RenderDiaryIndex)
	}

	// 需要鉴权的日记表单路由
	authForm := router.Group("/form")
	authForm.Use(user.AuthRequired())
	{
		authForm.GET("/post", h.RenderPostForm)
	}

	// 需要鉴权的日记数据操作路由（保持原路径）
	authAPI := router.Group("/")
	authAPI.Use(user.AuthRequired())
	{
		authAPI.POST("/post", h.HandlePost)
		authAPI.POST("/info", h.HandlePost)
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
func (h *Handler) parsePostID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.ParseInputParamErr
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
func (h *Handler) parseDiaryID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.ParseInputParamErr
	}
	return id, nil
}

// RenderPostForm 渲染日记编辑表单页面
// 根据action参数决定是渲染创建表单还是编辑表单
// 对于编辑操作，会验证用户权限并加载现有日记数据
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 支持的URL参数:
//   - action: "post.create" 或 "post.edit"
//   - post_id: 当action为edit时需要提供的日记ID
func (h *Handler) RenderPostForm(c *gin.Context) {
	// 获取当前登录用户信息，如果未登录会panic（由中间件保证已登录）
	loginUser := h.auth.MustGetCurrentUser(c)

	// 初始化模板数据
	data := pkg.DefaultH(c)
	data["Active"] = "users"
	data["UserData"] = loginUser

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	if action == "post.edit" {
		// 编辑模式：获取并验证日记数据
		post, err := h.getAndValidatePost(c, loginUser.UserId)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		data["Post"] = post
	}
	// 创建模式：直接渲染空表单

	c.HTML(http.StatusOK, "form/post", data)
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
func (h *Handler) HandlePost(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := h.auth.MustGetCurrentUser(c)

	// 解析表单数据到PostForm结构体
	postForm := &PostForm{}
	if err := c.Bind(postForm); err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	// 根据操作类型分发处理
	action, _ := c.GetPostForm("action")
	switch action {
	case "post.create":
		h.handleCreatePost(c, postForm, loginUser)
	case "post.edit":
		h.handleEditPost(c, postForm, loginUser)
	default:
		errors.ResponseErr(c, errors.ParseInputParamErr)
	}
}

// handleCreatePost 处理创建新日记的业务逻辑
// 设置日记的所有必要字段，调用服务层创建日记，并返回创建结果
//
// 参数:
//   - c: Gin上下文
//   - postForm: 解析后的表单数据
//   - user: 当前登录用户信息
func (h *Handler) handleCreatePost(c *gin.Context, postForm *PostForm, user *user.Login) {
	config := configs.GetConfig()
	nowTime := pkg.GetNowTimeUnixMilli()

	// 构建新日记对象，设置所有必要字段
	post := &Post{
		DiaryID:       postForm.DiaryID,
		Title:         postForm.Title,
		Body:          postForm.Body,
		Summary:       pkg.TruncateString(postForm.Body, config.App.SummaryMax, config.App.Suffix),
		IsTop:         0,
		InfoType:      1, // 默认为文本类型
		AddUTCTime:    nowTime,
		AddUserID:     user.UserId,
		AddUserName:   user.UserName,
		AddUserIP:     c.ClientIP(),
		UpdateUTCTime: nowTime,
		Status:        1, // 默认状态为正常
		ViewNum:       0,
		LikeNum:       0,
		CommentNum:    0,
	}

	// 调用服务层创建日记
	postID, err := h.service.AddPost(post)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}
	post.PostID = postID

	// 日记本信息数量加一 - 在handler层处理跨模块调用
	diaryService := diary.NewService()
	err = diaryService.IncrementPostCount(post.DiaryID, post.AddUserID, post.AddUserName, post.AddUTCTime)
	if err != nil {
		// 如果更新日记本数量失败，记录错误但不影响日记创建的成功响应
		// 这里可以添加日志记录
	}

	// 返回成功响应
	h.respondSuccess(c, post)
}

// handleEditPost 处理编辑现有日记的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记
//
// 参数:
//   - c: Gin上下文
//   - postForm: 解析后的表单数据
//   - user: 当前登录用户信息
func (h *Handler) handleEditPost(c *gin.Context, postForm *PostForm, user *user.Login) {
	// 获取要编辑的日记ID
	postIdStr, _ := c.GetPostForm("post_id")
	postId, err := h.parsePostID(postIdStr)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 验证日记所有权
	originalPost, err := h.validatePostOwnership(postId, user.UserId)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 使用结构体字面量构建更新后的日记对象
	config := configs.GetConfig()
	post := &Post{
		PostID:        postId,
		DiaryID:       originalPost.DiaryID, // 保持原有日记本ID
		Title:         postForm.Title,
		Body:          postForm.Body,
		Summary:       pkg.TruncateString(postForm.Body, config.App.SummaryMax, config.App.Suffix),
		UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
		// 保留原有字段值
		ViewNum:    originalPost.ViewNum,
		LikeNum:    originalPost.LikeNum,
		CommentNum: originalPost.CommentNum,
		IsTop:      originalPost.IsTop,
		InfoType:   originalPost.InfoType,
	}

	// 调用服务层更新日记
	if err := h.service.UpdatePost(post); err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 返回成功响应
	h.respondSuccess(c, post)
}

// getAndValidatePost 获取并验证日记权限（用于表单渲染）
// 从URL参数中获取post_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *Post: 日记对象，如果验证失败返回nil
//   - error: 错误信息
func (h *Handler) getAndValidatePost(c *gin.Context, userID int64) (*Post, error) {
	postID, err := h.parsePostID(c.Query("post_id"))
	if err != nil {
		return nil, err
	}

	return h.validatePostOwnership(postID, userID)
}

// validatePostOwnership 验证日记所有权
// 检查指定日记是否存在，以及当前用户是否有权限操作该日记
//
// 参数:
//   - postID: 日记ID
//   - userID: 当前用户ID
//
// 返回值:
//   - *Post: 日记对象
//   - error: 验证失败时的错误信息
func (h *Handler) validatePostOwnership(postID, userID int64) (*Post, error) {
	post, err := h.service.GetPostByID(postID)
	if err != nil || post == nil {
		return nil, errors.ParseInputParamErr
	}

	// 检查是否是当前用户的日记
	if post.AddUserID != userID {
		return nil, errors.InputParamErr
	}

	return post, nil
}

// respondSuccess 统一的成功响应格式
// 严格保持原有的响应格式，确保前后端配合不出错
//
// 参数:
//   - c: Gin上下文
//   - post: 操作成功的日记对象
func (h *Handler) respondSuccess(c *gin.Context, post *Post) {
	data := pkg.ResponseH(c)
	data["Data"] = post
	c.JSON(http.StatusOK, data)
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
func (h *Handler) RenderDiaryIndex(c *gin.Context) {
	// 解析日记本ID
	diaryIDStr := c.Param("spaceid")
	diaryID, err := h.parseDiaryID(diaryIDStr)
	if err != nil {
		h.renderErrorPage(c, "日记本ID格式错误", "请检查URL中的日记本ID是否正确")
		return
	}

	// 获取日记本信息 - 调用diary模块
	diaryService := diary.NewService()
	diaryInfo, err := diaryService.GetDiaryByID(diaryID)
	if err != nil || diaryInfo == nil {
		h.renderErrorPage(c, "日记本不存在", "未找到指定的日记本")
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
	config := configs.GetConfig()
	pageSize := config.App.PageSize

	// 获取日记列表
	posts, diariesNum, err := h.service.ListPostsByUserID(diaryInfo.AddUserID, diaryID, pageIndex, pageSize)
	var hasError bool
	var errMsg string

	if err != nil {
		hasError = true
		errMsg = "获取日记列表失败"
		posts = []Post{} // 空列表
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
	isAuth, loginUser := h.auth.IsAuthenticated(c)
	isOwner := isAuth && loginUser.UserId == diaryInfo.AddUserID

	// 设置模板数据
	data := pkg.DefaultH(c)
	data["Diary"] = diaryInfo
	data["Posts"] = posts
	data["PageLinks"] = pageLinks
	data["hasError"] = hasError
	data["errMsg"] = errMsg
	data["DiariesNum"] = diariesNum
	data["PageSize"] = pageSize
	data["IsOwner"] = isOwner

	c.HTML(http.StatusOK, "diary/index", data)
}

// renderErrorPage 渲染错误页面
// 当发生错误时显示用户友好的错误信息
//
// 参数:
//   - c: Gin上下文
//   - title: 错误标题
//   - message: 错误详细信息
func (h *Handler) renderErrorPage(c *gin.Context, title, message string) {
	data := pkg.DefaultH(c)
	data["title"] = title
	data["message"] = message
	c.HTML(http.StatusBadRequest, "error", data)
}
