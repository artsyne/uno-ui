package comment

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
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
	if loginUser == nil {
		ctrl.renderService.JSONError(c, errors.Unauthorized)
		return
	}

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[CreateCommentRequest](c)

	// 校验评论内容不能超过1000字符
	if len(req.Content) > 1000 {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	// 构建新评论对象，设置所有必要字段
	commentData := &comment.Comment{
		EntryID:            req.EntryID,
		Content:            req.Content,
		AddUserID:          loginUser.UserId,
		AddUserHandle:      loginUser.Handle,
		AddUserDisplayName: loginUser.DisplayName,
		AddUserIP:          c.ClientIP(),
	}

	// 验证 Entry 是否存在
	// 尝试获取用户自己的 Entry（包含草稿、回收站等状态）
	entryData, err := ctrl.entryService.GetEntryByIDForUser(req.EntryID, loginUser.UserId)
	if err != nil {
		// 如果不是自己的 Entry，或者是自己的但未找到（理论上 GetEntryByIDForUser 已覆盖），
		// 则尝试作为普通 Active Entry 获取（用于查看他人的公开日记）
		entryData, err = ctrl.entryService.GetEntryByID(req.EntryID)
	}

	if err != nil {
		ctrl.renderService.JSONError(c, errors.EntryNotFound)
		return
	}

	// 验证 Diary 社交设置
	_, err = ctrl.diaryService.CheckSocialEnabled(entryData.DiaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 调用服务层创建评论
	commentID, err := ctrl.commentService.AddComment(commentData)
	if err != nil {
		slog.Error("Failed to add comment", "error", err, "entry_id", req.EntryID, "user_id", loginUser.UserId)
		ctrl.renderService.JSONError(c, err)
		return
	}
	commentData.CommentID = commentID

	// 读取完整评论对象
	latestComment, err := ctrl.commentService.GetCommentByID(commentID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 不返回用户提交ip
	latestComment.AddUserIP = ""

	// 检查是否为HTMX请求
	if c.GetHeader("HX-Request") == "true" {
		// 返回 200 OK 并触发 commentCreated 事件
		c.Header("HX-Trigger", "commentCreated")
		c.Status(200)
		return
	}

	// 返回最新评论对象和用户时区，供前端格式化显示
	ctrl.renderService.JSONSuccess(c, gin.H{
		"comment":   latestComment,
		"time_zone": loginUser.TimeZone,
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
	if loginUser == nil {
		ctrl.renderService.JSONError(c, errors.Unauthorized)
		return
	}

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[DeleteCommentRequest](c)

	// 调用服务层删除评论（服务层会验证删除权限）
	err := ctrl.commentService.DeleteComment(req.CommentID, loginUser.UserId)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 检查是否为HTMX请求
	if c.GetHeader("HX-Request") == "true" {
		// 返回空内容，配合前端 hx-swap="outerHTML" 实现删除效果
		c.Status(200)
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
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetCommentsRequest](c)

	// 验证 Entry 是否存在
	// 尝试获取当前登录用户
	loginUser, isAuth := ctrl.auth.GetCurrentUser(c)
	var entryData *entry.Entry
	var err error

	if isAuth {
		// 如果已登录，尝试获取用户自己的 Entry（包含草稿、回收站等状态）
		entryData, err = ctrl.entryService.GetEntryByIDForUser(req.EntryID, loginUser.UserId)
		if err != nil {
			// 如果不是自己的 Entry，或者是自己的但未找到（理论上 GetEntryByIDForUser 已覆盖），
			// 则尝试作为普通 Active Entry 获取（用于查看他人的公开日记）
			entryData, err = ctrl.entryService.GetEntryByID(req.EntryID)
		}
	} else {
		// 未登录，只获取 Active Entry
		entryData, err = ctrl.entryService.GetEntryByID(req.EntryID)
	}

	if err != nil {
		ctrl.renderService.JSONError(c, &errors.EntryNotFound)
		return
	}

	// 验证 Diary 社交设置
	_, err = ctrl.diaryService.CheckSocialEnabled(entryData.DiaryID)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 设置默认分页参数
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	// 始终使用系统配置的页面大小，不接受外部limit/page_size
	req.PageSize = ctrl.config.App.PageSize

	// 调用服务层获取评论列表
	comments, total, err := ctrl.commentService.ListCommentsByEntryID(req.EntryID, req.PageIndex, req.PageSize)
	if err != nil {
		slog.Error("ListCommentsByEntryID failed", "error", err, "entry_id", req.EntryID)
		if c.GetHeader("HX-Request") == "true" || c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.String(200, `<div class="text-center text-red-500 py-8">Failed to load comments</div>`)
			return
		}
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 计算总页数
	totalPages := (total + req.PageSize - 1) / req.PageSize

	// 获取用户时区（未登录则为UTC）
	var currentUserID int64
	tz := "UTC"
	if user, ok := ctrl.auth.GetCurrentUser(c); ok && user != nil {
		currentUserID = user.UserId
		if user.TimeZone != "" {
			tz = user.TimeZone
		}
	}

	// 获取日记作者ID，用于前端判断删除权限
	var entryOwnerID int64
	if entry, err := ctrl.entryService.GetEntryByID(req.EntryID); err == nil && entry != nil {
		entryOwnerID = entry.AddUserID
	}

	// 检查是否为HTMX请求
	if c.GetHeader("HX-Request") == "true" {
		lang := ctrl.renderService.GetLocaleFromContext(c)
		pageLinks := pkg.GeneratePagination(pkg.ListOptions{
			PageIndex: req.PageIndex,
			PageSize:  req.PageSize,
			Total:     total,
			URL:       fmt.Sprintf("/api/comment/list?entry_id=%d&page_index=%%d", req.EntryID),
			PagePrev:  ctrl.renderService.GetMessage("common.previous", lang),
			PageNext:  ctrl.renderService.GetMessage("common.next", lang),
		})

		// Calculate pagination for manual template rendering
		var prevPage, nextPage int
		if req.PageIndex > 1 {
			prevPage = req.PageIndex - 1
		}
		if req.PageIndex < totalPages {
			nextPage = req.PageIndex + 1
		}

		data := gin.H{
			"comments":      comments,
			"PageLinks":     pageLinks,
			"HxTarget":      "#comments-list-container",
			"time_zone":     tz,
			"CurrentUserID": currentUserID,
			"EntryID":       req.EntryID,
			"EntryOwnerID":  entryOwnerID,
			"TotalCount":    total,
			"pagination": gin.H{
				"total_pages":  totalPages,
				"current_page": req.PageIndex,
				"prev_page":    prevPage,
				"next_page":    nextPage,
			},
		}

		// 设置总数 Header，供前端 JS 更新 UI
		c.Header("X-Total-Count", strconv.Itoa(total))

		pageData := ctrl.renderService.CreatePageDataWithData(c, "comment.title", data)
		c.HTML(200, "comment/list", pageData)
		return
	}

	// 返回评论列表响应
	ctrl.renderService.JSONSuccess(c, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"total":       total,
			"page_index":  req.PageIndex,
			"page_size":   req.PageSize,
			"total_pages": totalPages,
		},
		"time_zone": tz,
	})
}

// GetCommentsByEntryIDFromPath 从URL路径参数获取日记ID并返回评论列表
// 这是一个便捷方法，用于处理RESTful风格的URL
//
// 参数:
//   - c: Gin上下文对象
//
// URL格式: /api/entry/:entry_id/comments
func (ctrl *CommentController) GetCommentsByEntryIDFromPath(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[GetCommentsByPathRequest](c)

	// 固定页面大小为配置中的 PageSize
	pageSize := ctrl.config.App.PageSize
	// 设置默认页码
	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}

	// 验证 Entry 是否存在且为 Active 状态
	// 由于此接口可能公开访问，默认只允许访问 Active 状态的日记
	// 如果需要访问非 Active 状态（如回收站），应使用 /api/comment/list 接口（带鉴权）
	_, err := ctrl.entryService.GetEntryByID(req.EntryID)
	if err != nil {
		ctrl.renderService.JSONError(c, &errors.EntryNotFound)
		return
	}

	// 调用服务层获取评论列表
	comments, total, err := ctrl.commentService.ListCommentsByEntryID(req.EntryID, pageIndex, pageSize)
	if err != nil {
		slog.Error("ListCommentsByEntryID failed", "error", err, "entry_id", req.EntryID)
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 计算总页数
	totalPages := (total + pageSize - 1) / pageSize
	nextPage := 0
	if pageIndex < totalPages {
		nextPage = pageIndex + 1
	}
	prevPage := 0
	if pageIndex > 1 {
		prevPage = pageIndex - 1
	}

	// 获取当前用户ID和时区
	var currentUserID int64
	tz := "UTC"
	if user, ok := ctrl.auth.GetCurrentUser(c); ok && user != nil {
		currentUserID = user.UserId
		if user.TimeZone != "" {
			tz = user.TimeZone
		}
	}

	// 检查是否为HTMX请求
	if c.GetHeader("HX-Request") == "true" {
		data := ctrl.renderService.CreatePageData(c, "")
		data.Data["comments"] = comments
		data.Data["pagination"] = gin.H{
			"total_pages":  totalPages,
			"next_page":    nextPage,
			"prev_page":    prevPage,
			"current_page": pageIndex,
		}
		data.Data["time_zone"] = tz
		data.Data["CurrentUserID"] = currentUserID
		data.Data["EntryID"] = req.EntryID
		data.Data["TotalCount"] = total

		ctrl.renderService.HTML(c, "comment/list", data)
		return
	}

	// 返回成功响应（包含分页和时区）
	ctrl.renderService.JSONSuccess(c, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"total_pages": totalPages,
			"next_page":   nextPage,
		},
		"time_zone": tz,
	})
}
