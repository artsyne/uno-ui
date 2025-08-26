package user

import (
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/gin-gonic/gin"
)

// ValidateUsername 验证用户名是否可用
func (ctrl *UserController) ValidateUsername(c *gin.Context) {
	// 解析请求数据
	var req CheckUserNameRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 检查用户名是否已存在
	isExistUserName, err := ctrl.service.IsExistUserName(req.UserName)
	if err != nil {
		ctrl.renderService.JSONError(c, err, map[string]interface{}{})
		return
	}
	if isExistUserName {
		ctrl.renderService.JSONError(c, errors.UserExists.WithDetails("用户名: "+req.UserName), map[string]interface{}{})
		return
	}

	ctrl.renderService.JSONSuccess(c, map[string]interface{}{})
}

// ValidateEmail 验证邮箱是否可用
func (ctrl *UserController) ValidateEmail(c *gin.Context) {
	// 解析请求数据
	var req CheckEmailRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 检查邮箱是否已存在
	isExistEmail, err := ctrl.service.IsExistEmail(req.Email)
	if err != nil {
		ctrl.renderService.JSONError(c, err, map[string]interface{}{})
		return
	}
	if isExistEmail {
		ctrl.renderService.JSONError(c, errors.EmailExists.WithDetails("邮箱: "+req.Email), map[string]interface{}{})
		return
	}

	ctrl.renderService.JSONSuccess(c, map[string]interface{}{})
}



// parseUserID 解析用户ID
func (ctrl *UserController) parseUserID(idStr string) (int64, error) {
	return ctrl.service.ParseUserID(idStr)
}