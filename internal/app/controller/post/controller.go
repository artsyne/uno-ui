package post

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// PostController 日记控制器
// 负责处理所有与日记相关的HTTP请求，包括渲染表单、创建日记、编辑日记等操作
//
// 字段:
//   - postService: 日记业务逻辑服务
//   - diaryService: 日记本业务逻辑服务
//   - auth: 用户认证辅助工具
//   - config: 应用配置
//   - renderService: 视图渲染服务
//   - userService: 用户业务逻辑服务
type PostController struct {
	postService   *post.Service
	diaryService  *diary.Service
	auth          *user.AuthHelper
	config        *config.Config
	renderService *view.ViewService
	userService *user.Service
}

// NewPostController 创建新的日记控制器实例
// 通过依赖注入的方式初始化PostController，确保所有必需的服务都已正确注入
//
// 参数:
//   - config: 应用配置实例
//   - postService: 日记业务逻辑服务
//   - diaryService: 日记本业务逻辑服务
//   - auth: 用户认证辅助工具
//   - renderService: 视图渲染服务
//   - userService: 用户业务逻辑服务
//
// 返回:
//   - *PostController: 初始化完成的日记控制器实例
func NewPostController(config *config.Config, postService *post.Service, diaryService *diary.Service, auth *user.AuthHelper, renderService *view.ViewService, userService *user.Service) *PostController {
	return &PostController{
		config:        config,
		postService:   postService,
		diaryService:  diaryService,
		auth:          auth,
		renderService: renderService,
		userService:   userService,
	}
}