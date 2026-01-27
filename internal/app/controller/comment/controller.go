package comment

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// CommentController 评论控制器
// 负责处理所有与评论相关的HTTP请求，包括创建评论、删除评论、获取评论列表等操作
//
// 字段:
//   - commentService: 评论业务逻辑服务
//   - diaryService: 日记本业务逻辑服务
//   - entryService: 日记业务逻辑服务
//   - auth: 用户认证辅助工具
//   - config: 应用配置
//   - renderService: 视图渲染服务
type CommentController struct {
	commentService *comment.Service
	diaryService   *diary.Service
	entryService   *entry.Service
	auth           *user.AuthHelper
	config         *config.Config
	renderService  *view.ViewService
}

// NewCommentController 创建新的评论控制器实例
// 通过依赖注入的方式初始化CommentController，确保所有必需的服务都已正确注入
//
// 参数:
//   - config: 应用配置实例
//   - commentService: 评论业务逻辑服务
//   - diaryService: 日记本业务逻辑服务
//   - entryService: 日记业务逻辑服务
//   - auth: 用户认证辅助工具
//   - renderService: 视图渲染服务
//
// 返回:
//   - *CommentController: 初始化完成的评论控制器实例
func NewCommentController(config *config.Config, commentService *comment.Service, diaryService *diary.Service, entryService *entry.Service, auth *user.AuthHelper, renderService *view.ViewService) *CommentController {
	return &CommentController{
		config:         config,
		commentService: commentService,
		diaryService:   diaryService,
		entryService:   entryService,
		auth:           auth,
		renderService:  renderService,
	}
}