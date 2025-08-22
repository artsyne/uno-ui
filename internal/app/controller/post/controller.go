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
type PostController struct {
	postService   *post.Service
	diaryService  *diary.Service
	auth          *user.AuthHelper
	config        *config.Config
	renderService *view.ViewService
	userService *user.Service
}

// NewPostController 创建新的日记控制器实例
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