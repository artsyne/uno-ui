package controllers

import (
	"net/http"
	"strconv"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/helpers"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// PostController 日记控制器
// 负责处理日记相关的所有HTTP请求，包括渲染表单、创建日记、编辑日记等操作
type PostController struct{}

// NewPostController 创建新的日记控制器实例
// 返回值: *PostController 日记控制器指针
func NewPostController() *PostController {
	return &PostController{}
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
func (pc *PostController) RenderPostForm(c *gin.Context) {
	// 获取当前登录用户信息，如果未登录会panic（由中间件保证已登录）
	loginUser := helpers.Auth.MustGetCurrentUser(c)

	// 初始化模板数据
	h := pkg.DefaultH(c)
	h["Active"] = "users"
	h["UserData"] = loginUser

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	if action == "post.edit" {
		// 编辑模式：获取并验证日记数据
		post, err := pc.getAndValidatePost(c, loginUser.UserId)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		h["Post"] = post
	}
	// 创建模式：直接渲染空表单

	c.HTML(http.StatusOK, "form/post", h)
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
func (pc *PostController) HandlePost(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := helpers.Auth.MustGetCurrentUser(c)

	// 解析表单数据到PostForm结构体
	postForm := &models.PostForm{}
	if err := c.Bind(postForm); err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	// 根据操作类型分发处理
	action, _ := c.GetPostForm("action")
	switch action {
	case "post.create":
		pc.handleCreatePost(c, postForm, loginUser)
	case "post.edit":
		pc.handleEditPost(c, postForm, loginUser)
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
func (pc *PostController) handleCreatePost(c *gin.Context, postForm *models.PostForm, user *models.Login) {
	config := configs.GetConfig()
	nowTime := pkg.GetNowTimeUnixMilli()

	// 构建新日记对象，设置所有必要字段
	post := &models.Post{
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
	postID, err := services.AddPost(post)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}
	post.PostID = postID

	// 返回成功响应
	pc.respondSuccess(c, post)
}

// handleEditPost 处理编辑现有日记的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记
//
// 参数:
//   - c: Gin上下文
//   - postForm: 解析后的表单数据
//   - user: 当前登录用户信息
func (pc *PostController) handleEditPost(c *gin.Context, postForm *models.PostForm, user *models.Login) {
	// 获取要编辑的日记ID
	postIdStr, _ := c.GetPostForm("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil || postId <= 0 {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	// 验证日记所有权
	originalPost, err := pc.validatePostOwnership(postId, user.UserId)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 构建更新后的日记对象
	config := configs.GetConfig()
	post := &models.Post{
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
	if err := services.UpdatePost(post); err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 返回成功响应
	pc.respondSuccess(c, post)
}

// getAndValidatePost 获取并验证日记权限（用于表单渲染）
// 从URL参数中获取post_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *models.Post: 日记对象，如果验证失败返回nil
//   - error: 错误信息
func (pc *PostController) getAndValidatePost(c *gin.Context, userID int64) (*models.Post, error) {
	postId, err := strconv.ParseInt(c.Query("post_id"), 10, 64)
	if err != nil || postId <= 0 {
		return nil, errors.ParseInputParamErr
	}

	return pc.validatePostOwnership(postId, userID)
}

// validatePostOwnership 验证日记所有权
// 检查指定日记是否存在，以及当前用户是否有权限操作该日记
//
// 参数:
//   - postID: 日记ID
//   - userID: 当前用户ID
//
// 返回值:
//   - *models.Post: 日记对象
//   - error: 验证失败时的错误信息
func (pc *PostController) validatePostOwnership(postID, userID int64) (*models.Post, error) {
	post, err := services.GetPostByID(postID)
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
// 封装标准的JSON响应格式，包含操作成功的日记数据
//
// 参数:
//   - c: Gin上下文
//   - post: 操作成功的日记对象
func (pc *PostController) respondSuccess(c *gin.Context, post *models.Post) {
	h := pkg.ResponseH(c)
	h["Data"] = post
	c.JSON(http.StatusOK, h)
}