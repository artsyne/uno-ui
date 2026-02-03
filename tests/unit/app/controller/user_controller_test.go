package controller_test

import (
	"testing"

	"github.com/biuren/hebe/internal/app/controller/user"
	"github.com/biuren/hebe/internal/app/view"
	domainUser "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/stretchr/testify/assert"
)

// TestNewUserController 测试控制器创建
func TestNewUserController(t *testing.T) {
	cfg := &config.Config{}

	// 创建真实的服务对象用于测试
	mockUserSvc := &domainUser.Service{}
	mockAuth := &domainUser.AuthHelper{}
	mockViewSvc := &view.ViewService{}

	// 创建控制器
	controller := user.NewUserController(
		cfg,
		mockUserSvc,
		mockViewSvc,
		mockAuth,
	)

	// 验证控制器不为空
	assert.NotNil(t, controller)
}

// TestUserController_BasicFunctionality 测试基本功能
func TestUserController_BasicFunctionality(t *testing.T) {
	cfg := &config.Config{}

	// 创建真实的服务对象
	userSvc := &domainUser.Service{}
	auth := &domainUser.AuthHelper{}
	viewSvc := &view.ViewService{}

	// 创建控制器
	controller := user.NewUserController(cfg, userSvc, viewSvc, auth)

	// 验证控制器不为空
	assert.NotNil(t, controller)
}

// TestUserController_ControllerCreation 测试控制器创建的不同场景
func TestUserController_ControllerCreation(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		userService *domainUser.Service
		auth        *domainUser.AuthHelper
		viewService *view.ViewService
		expectNil   bool
	}{
		{
			name:        "valid controller creation",
			config:      &config.Config{},
			userService: &domainUser.Service{},
			auth:        &domainUser.AuthHelper{},
			viewService: &view.ViewService{},
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := user.NewUserController(
				tt.config,
				tt.userService,
				tt.viewService,
				tt.auth,
			)

			if tt.expectNil {
				assert.Nil(t, controller)
			} else {
				assert.NotNil(t, controller)
			}
		})
	}
}