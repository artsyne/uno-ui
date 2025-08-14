package user

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// UserController 用户HTTP处理层实现
type UserController struct {
	config        *config.Config
	service       *user.Service
	auth          *user.AuthHelper
	renderService *view.ViewService
}

// NewUserController 创建用户HTTP处理层实例
func NewUserController(config *config.Config, service *user.Service, renderService *view.ViewService, auth *user.AuthHelper) *UserController {
	return &UserController{
		config:        config,
		service:       service,
		renderService: renderService,
		auth:          auth,
	}
}