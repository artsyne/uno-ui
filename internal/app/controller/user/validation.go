package user

import (
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/utils/errors"
	pkg "github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ValidateUsername 验证用户名是否可用
func (ctrl *UserController) ValidateUsername(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[CheckUserNameRequest](c)

	// 检查用户名是否已存在
	isExistUserName, err := ctrl.service.IsExistUserName(req.UserName)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	if isExistUserName {
		ctrl.renderService.JSONError(c, errors.UserExists.WithDetails("用户名: "+req.UserName))
		return
	}

	ctrl.renderService.JSONSuccess(c, map[string]interface{}{"available": true})
}

// ValidateEmail 验证邮箱是否可用
func (ctrl *UserController) ValidateEmail(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[CheckEmailRequest](c)

	// 验证邮箱格式
	if !pkg.IsEmail(req.Email) {
		ctrl.renderService.JSONError(c, errors.InvalidEmail)
		return
	}

	// 检查邮箱是否已存在
	isExistEmail, err := ctrl.service.IsExistEmail(req.Email)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}
	if isExistEmail {
		ctrl.renderService.JSONError(c, errors.EmailExists.WithDetails("邮箱: "+req.Email))
		return
	}

	ctrl.renderService.JSONSuccess(c, map[string]interface{}{"available": true})
}
