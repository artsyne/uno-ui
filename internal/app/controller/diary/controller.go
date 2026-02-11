package diary

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/like"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// DiaryController 日记本控制器
// 负责处理所有与日记本相关的HTTP请求，包括渲染列表、表单、创建日记本、编辑日记本等操作
//
// 字段:
//   - config: 应用程序配置实例
//   - diaryService: 日记本业务逻辑服务
//   - entryService: 日记业务逻辑服务
//   - userService: 用户业务逻辑服务
//   - likeService: 点赞业务逻辑服务
//   - auth: 用户认证辅助工具
//   - renderService: 视图渲染服务实例
type DiaryController struct {
	config        *config.Config
	diaryService  *diary.Service
	entryService  *entry.Service
	userService   *user.Service
	likeService   like.Service
	auth          *user.AuthHelper
	renderService *view.ViewService
}

// NewDiaryController 创建新的日记本控制器实例
// 初始化日记本控制器，注入所有必要的依赖服务
//
// 参数:
//   - config: 应用程序配置实例
//   - diaryService: 日记本业务逻辑服务
//   - entryService: 日记业务逻辑服务
//   - userService: 用户业务逻辑服务
//   - likeService: 点赞业务逻辑服务
//   - auth: 用户认证辅助工具
//   - renderService: 视图渲染服务实例
//
// 返回:
//   - *DiaryController: 初始化完成的日记本控制器实例
func NewDiaryController(config *config.Config, diaryService *diary.Service, entryService *entry.Service, userService *user.Service, likeService like.Service, auth *user.AuthHelper, renderService *view.ViewService) *DiaryController {
	return &DiaryController{
		config:        config,
		diaryService:  diaryService,
		entryService:  entryService,
		userService:   userService,
		likeService:   likeService,
		auth:          auth,
		renderService: renderService,
	}
}