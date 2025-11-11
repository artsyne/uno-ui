package comment

import (
	"strconv"

	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

// CreateComment 处理创建新评论的业务逻辑
// 验证用户登录状态，解析请求参数，调用服务层创建评论
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 验证用户登录状态
//   - 解析创建评论请求参数
//   - 构建新评论对象并设置必要字段
//   - 调用服务层创建评论
//   - 返回创建成功的响应
func (ctrl *CommentController) CreateComment(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析创建评论请求
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 构建新评论对象，设置所有必要字段
	commentData := &comment.Comment{
		PostID:      req.PostID,
		Content:     req.Content,
		AddUserID:   loginUser.UserId,
		AddUserName: loginUser.UserName,
		AddUserIP:   c.ClientIP(),
	}

	// 调用服务层创建评论
	commentID, err := ctrl.commentService.AddComment(commentData)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	commentData.CommentID = commentID

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, gin.H{
		"comment_id": commentID,
		"message":    "评论创建成功",
	})
}

// DeleteComment 处理删除评论的业务逻辑
// 验证用户登录状态，解析请求参数，调用服务层删除评论
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 验证用户登录状态
//   - 解析删除评论请求参数
//   - 调用服务层删除评论（包含权限验证）
//   - 返回删除成功的响应
func (ctrl *CommentController) DeleteComment(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析删除评论请求
	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 调用服务层删除评论（服务层会验证删除权限）
	err := ctrl.commentService.DeleteComment(req.CommentID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, gin.H{
		"message": "评论删除成功",
	})
}

// GetComments 处理获取评论列表的业务逻辑
// 解析请求参数，调用服务层获取分页评论列表
//
// 参数:
//   - c: Gin上下文对象
//
// 功能:
//   - 解析获取评论列表请求参数
//   - 设置默认分页参数
//   - 调用服务层获取评论列表
//   - 构建分页响应数据
//   - 返回评论列表响应
func (ctrl *CommentController) GetComments(c *gin.Context) {
	// 解析获取评论列表请求
	var req GetCommentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 设置默认分页参数
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务层获取评论列表
	comments, total, err := ctrl.commentService.ListCommentsByPostID(req.PostID, req.PageIndex, req.PageSize)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 计算总页数
	totalPages := (total + req.PageSize - 1) / req.PageSize

	// 构建响应数据
	response := comment.CommentListResponse{
		Comments:   comments,
		Total:      total,
		PageIndex:  req.PageIndex,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, response)
}

// GetCommentsByPostIDFromPath 从URL路径参数获取日记ID并返回评论列表
// 这是一个便捷方法，用于处理RESTful风格的URL
//
// 参数:
//   - c: Gin上下文对象
//
// URL格式: /api/post/:post_id/comments
func (ctrl *CommentController) GetCommentsByPostIDFromPath(c *gin.Context) {
	// 从URL路径参数获取日记ID
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil || postID <= 0 {
		ctrl.renderService.JSONError(c, errors.InvalidParams.WithDetails("无效的日记ID"))
		return
	}

	// 获取分页参数
	pageIndexStr := c.DefaultQuery("page_index", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil || pageIndex <= 0 {
		pageIndex = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 调用服务层获取评论列表
	comments, total, err := ctrl.commentService.ListCommentsByPostID(postID, pageIndex, pageSize)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 计算总页数
	totalPages := (total + pageSize - 1) / pageSize

	// 构建响应数据
	response := comment.CommentListResponse{
		Comments:   comments,
		Total:      total,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// 返回成功响应
	ctrl.renderService.JSONSuccess(c, response)
}