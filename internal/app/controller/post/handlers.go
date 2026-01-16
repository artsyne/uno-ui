package post

import (
	"log/slog"
	"strconv"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// CreatePost 处理创建新日记的请求
// 绑定JSON或表单数据到CreatePostRequest，并调用服务层创建日记
//
// 路由: POST /post
func (ctrl *PostController) CreatePost(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[CreatePostRequest](c)

	// 校验日记内容长度
	if len(req.Body) > ctrl.config.App.MaxPostLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	nowTime := pkg.GetNowTimeUnixMilli()

	// 构建新日记对象，设置所有必要字段
	postData := &post.Post{
		DiaryID:       req.DiaryID,
		Title:         req.Title,
		Body:          req.Body,
		Summary:       pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
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
	ctrl.renderService.JSONSuccess(c, postData)
}

// UpdatePost 处理更新日记的请求
// 绑定URI参数(id)和JSON/表单数据到UpdatePostRequest，并调用服务层更新
//
// 路由: PUT /post/:id
func (ctrl *PostController) UpdatePost(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[UpdatePostRequest](c)

	// 校验日记内容长度
	if len(req.Body) > ctrl.config.App.MaxPostLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	postId := req.PostID

	// 验证日记所有权
	// 获取日记信息
	originalPost, err := ctrl.postService.GetPostByID(postId)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.PostNotFound)
		return
	}
	// 判断是否有权限修改
	if originalPost.AddUserID != loginUser.UserId {
		ctrl.renderService.JSONError(c, errors.PermissionDenied)
		return
	}

	// 使用结构体字面量构建更新后的日记对象
	postData := &post.Post{
		PostID:        postId,
		DiaryID:       originalPost.DiaryID, // 保持原有日记本ID
		Title:         req.Title,
		Body:          req.Body,
		Summary:       pkg.TruncateString(req.Body, ctrl.config.App.SummaryMax, ctrl.config.App.Suffix),
		UpdateUTCTime: pkg.GetNowTimeUnixMilli(),
		// 保留原有字段值
		ViewNum:    originalPost.ViewNum,
		LikeNum:    originalPost.LikeNum,
		CommentNum: originalPost.CommentNum,
		IsTop:      originalPost.IsTop,
		InfoType:   originalPost.InfoType,
		Status:     originalPost.Status,
	}

	// 调用服务层更新日记
	if err := ctrl.postService.UpdatePost(postData); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, postData)
}

// DeletePost 处理删除日记请求
// 执行软删除操作，将日记状态设置为0，并减少日记本的统计数
//
// 路由: DELETE /post/:id
func (ctrl *PostController) DeletePost(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[DeletePostRequest](c)

	// 获取当前用户信息
	userData, isAuthenticated := ctrl.auth.GetCurrentUser(c)
	if !isAuthenticated {
		ctrl.renderService.JSONError(c, &errors.LoginRequired)
		return
	}

	// 验证日记所有权
	postData, postErr := ctrl.postService.GetPostByID(req.PostID)
	if postErr != nil {
		ctrl.renderService.JSONError(c, &errors.PostNotFound)
		return
	}
	if postData.AddUserID != userData.UserId {
		ctrl.renderService.JSONError(c, &errors.PermissionDenied)
		return
	}

	// 执行删除操作
	currentTime := pkg.GetNowTimeUnixMilli()
	deleteErr := ctrl.postService.DeletePostByID(req.PostID, userData.UserId, userData.UserName, currentTime)
	if deleteErr != nil {
		ctrl.renderService.JSONError(c, deleteErr)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// GetPostsStats 获取日记统计信息
// 根据POST请求体中的postids获取日记的统计信息
//
// 路由: POST /api/post/stats
// 参数:
//   - postids: 日记ID列表，逗号分隔 (JSON请求体)
//
// 返回:
//   - 包含统计信息的日记列表
func (ctrl *PostController) GetPostsStats(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetPostsStatsRequest](c)

	// 解析postids字符串为int64切片
	postIDs, err := req.ParsePostIDs()
	if err != nil {
		ctrl.renderService.JSONError(c, err) // ParsePostIDs已经返回了合适的错误类型
		return
	}

	// 调用服务层获取统计信息
	posts, err := ctrl.postService.GetPostsStats(postIDs)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 构建响应数据
	statsResponse := make([]PostStatsResponse, 0, len(posts))
	for _, p := range posts {
		statsResponse = append(statsResponse, PostStatsResponse{
			DiaryID:    p.DiaryID,
			PostID:     p.PostID,
			ViewNum:    p.ViewNum,
			LikeNum:    p.LikeNum,
			CommentNum: p.CommentNum,
		})
	}

	// 返回成功响应
	response := GetPostsStatsResponse{
		Stats: statsResponse,
	}
	ctrl.renderService.JSONSuccess(c, response)
}

// LikePost 日记点赞
// 处理日记点赞/取消点赞操作，需要用户登录
//
// 参数:
//   - c: Gin上下文对象
func (ctrl *PostController) LikePost(c *gin.Context) {
	// 获取登录用户信息
	loginUser, exists := ctrl.auth.GetCurrentUser(c)
	if !exists {
		ctrl.renderService.JSONError(c, &errors.Unauthorized)
		return
	}

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[LikePostRequest](c)
	postID := req.PostID

	// 获取点赞前的状态，用于判断是点赞还是取消点赞
	beforeStatus, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, postID, like.LikeTypePost)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "post_id", postID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	// 切换点赞状态
	err = ctrl.likeService.ToggleLike(loginUser.UserId, postID, like.LikeTypePost)
	if err != nil {
		slog.Error("点赞操作失败", "error", err, "user_id", loginUser.UserId, "post_id", postID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	// 根据之前的状态更新日记的like_num字段
	if beforeStatus.IsLiked {
		// 之前已点赞，现在是取消点赞，减少like_num
		err = ctrl.postService.DecrementLikeCount(postID)
	} else {
		// 之前未点赞，现在是点赞，增加like_num
		err = ctrl.postService.IncrementLikeCount(postID)
	}

	if err != nil {
		slog.Error("更新日记点赞数失败", "error", err, "post_id", postID)
		// 这里不返回错误，因为点赞状态已经更新成功，只是统计数字更新失败
		// 可以考虑记录日志或者异步重试
	}

	// 获取最新的点赞状态
	status, err := ctrl.likeService.GetLikeStatus(loginUser.UserId, postID, like.LikeTypePost)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", loginUser.UserId, "post_id", postID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	ctrl.renderService.JSONSuccess(c, status)
}

// GetPostsLikeStatus 获取多个日记的点赞状态
// 根据POST请求体中的postids获取用户对这些日记的点赞状态
//
// 路由: POST /api/post/like/status
// 参数:
//   - postids: 日记ID列表，逗号分隔 (JSON请求体)
//
// 返回:
//   - 包含点赞状态信息的日记列表
func (ctrl *PostController) GetPostsLikeStatus(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetPostLikeStatusRequest](c)

	// 解析postids字符串为int64切片
	postIDs, err := req.ParsePostIDs()
	if err != nil {
		ctrl.renderService.JSONError(c, err) // ParsePostIDs已经返回了合适的错误类型
		return
	}

	// 获取登录用户信息（可能为空，访客也可以查看点赞状态）
	loginUser, _ := ctrl.auth.GetCurrentUser(c)
	var userID int64 = 0
	if loginUser != nil {
		userID = loginUser.UserId
	}

	// 构建响应数据
	statusResponse := make([]PostLikeStatusResponse, 0, len(postIDs))
	for _, postID := range postIDs {
		// 获取点赞状态
		status, err := ctrl.likeService.GetLikeStatus(userID, postID, like.LikeTypePost)
		if err != nil {
			slog.Error("获取点赞状态失败", "error", err, "user_id", userID, "post_id", postID)
			// 对于单个日记的错误，我们继续处理其他日记，但记录错误
			continue
		}

		statusResponse = append(statusResponse, PostLikeStatusResponse{
			PostID: postID,
			Liked:  status.IsLiked,
		})
	}

	// 返回成功响应
	response := GetPostLikeStatusResponse{
		Status: statusResponse,
	}
	ctrl.renderService.JSONSuccess(c, response)
}

// GetPostLikeStatus 获取日记点赞状态
// 获取指定日记的点赞状态信息，包括是否已点赞和点赞总数
//
// 参数:
//   - c: Gin上下文对象
func (ctrl *PostController) GetPostLikeStatus(c *gin.Context) {
	// 从路径参数获取post_id
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
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
	status, err := ctrl.likeService.GetLikeStatus(userID, postID, like.LikeTypePost)
	if err != nil {
		slog.Error("获取点赞状态失败", "error", err, "user_id", userID, "post_id", postID)
		ctrl.renderService.JSONError(c, &errors.InternalError)
		return
	}

	ctrl.renderService.JSONSuccess(c, status)
}
