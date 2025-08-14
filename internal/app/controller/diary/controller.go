package diary

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// DiaryController 日记本控制器
// 负责处理所有与日记本相关的HTTP请求，包括渲染列表、表单、创建日记本、编辑日记本等操作
type DiaryController struct {
	config        *config.Config
	diaryService  *diary.Service
	userService   *user.Service
	auth          *user.AuthHelper
	renderService *view.ViewService
}

// NewDiaryController 创建新的日记本控制器实例
func NewDiaryController(config *config.Config, diaryService *diary.Service, userService *user.Service, auth *user.AuthHelper, renderService *view.ViewService) *DiaryController {
	return &DiaryController{
		config:        config,
		diaryService:  diaryService,
		userService:   userService,
		auth:          auth,
		renderService: renderService,
	}
}