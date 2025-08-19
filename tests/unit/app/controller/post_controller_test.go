package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	diaryDomain "github.com/biuren/hebe/internal/domain/diary"
	postDomain "github.com/biuren/hebe/internal/domain/post"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/app/view"
	post "github.com/biuren/hebe/internal/app/controller/post"
)



// TestNewPostController 测试控制器创建
func TestNewPostController(t *testing.T) {
	cfg := &config.Config{}
	mockPostSvc := &postDomain.Service{}
	mockDiarySvc := &diaryDomain.Service{}
	mockAuth := &userDomain.AuthHelper{}
	mockViewSvc := &view.ViewService{}

	// 创建控制器
	controller := post.NewPostController(
		cfg,
		mockPostSvc,
		mockDiarySvc,
		mockAuth,
		mockViewSvc,
	)

	// 验证控制器不为空
	assert.NotNil(t, controller)
}