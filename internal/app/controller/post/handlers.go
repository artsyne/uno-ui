package post

import (
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/post"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ProcessPost 处理日记的创建和编辑请求
// 根据表单中的action字段分发到对应的处理方法
//
// 支持的表单action值:
//   - "post.create": 创建新日记
//   - "post.edit": 编辑现有日记
func (ctrl *PostController) ProcessPost(c *gin.Context) {
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
		ctrl.CreatePost(c, postForm, loginUser)
	case "post.edit":
		ctrl.UpdatePost(c, postForm, loginUser)
	default:
		ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
	}
}

// CreatePost 处理创建新日记的业务逻辑
// 设置日记的所有必要字段，调用服务层创建日记，并返回创建结果
func (ctrl *PostController) CreatePost(c *gin.Context, postForm *post.PostRequest, loginUser *types.AuthUser) {
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

// UpdatePost 处理编辑现有日记的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记
func (ctrl *PostController) UpdatePost(c *gin.Context, postForm *post.PostRequest, loginUser *types.AuthUser) {
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

// DeletePost 处理删除日记请求（未来扩展功能）
func (ctrl *PostController) DeletePost(c *gin.Context) {
	// 预留的删除功能实现
	ctrl.renderService.JSONError(c, errors.InvalidParams.Err())
}

// respondSuccess 统一的成功响应格式
// 严格保持原有的响应格式，确保前后端配合不出错
func (ctrl *PostController) respondSuccess(c *gin.Context, postData *post.Post) {
	ctrl.renderService.JSONSuccess(c, postData)
}