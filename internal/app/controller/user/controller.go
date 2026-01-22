package user

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// UserController 用户HTTP处理层实现
// 负责处理所有与用户相关的HTTP请求，包括用户认证、注册、个人信息管理等
//
// 字段:
//   - config: 应用配置信息
//   - service: 用户业务逻辑服务
//   - auth: 用户认证辅助工具
//   - renderService: 视图渲染服务
type UserController struct {
	config        *config.Config
	service       *user.Service
	auth          *user.AuthHelper
	renderService *view.ViewService
}

// NewUserController 创建用户HTTP处理层实例
// 通过依赖注入的方式初始化用户控制器，确保所有必要的服务都已正确配置
//
// 参数:
//   - config: 应用配置信息
//   - service: 用户业务逻辑服务
//   - renderService: 视图渲染服务
//   - auth: 用户认证辅助工具
//
// 返回:
//   - *UserController: 初始化完成的用户控制器实例
func NewUserController(config *config.Config, service *user.Service, renderService *view.ViewService, auth *user.AuthHelper) *UserController {
	return &UserController{
		config:        config,
		service:       service,
		renderService: renderService,
		auth:          auth,
	}
}
