package user

import (
	errors "github.com/biuren/hebe/internal/utils/error"
	pkg "github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ValidateUsername 验证用户名是否可用
func (ctrl *UserController) ValidateUsername(c *gin.Context) {
	userName, _ := c.GetPostForm("username")

	// 验证用户名格式
	if !ctrl.validateUserNameFormat(userName) {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("用户名格式不正确"), map[string]interface{}{})
		return
	}

	// 检查用户名是否已存在
	isExistUserName, err := ctrl.service.IsExistUserName(userName)
	if err != nil {
		ctrl.renderService.JSONError(c, err, map[string]interface{}{})
		return
	}
	if isExistUserName {
		ctrl.renderService.JSONError(c, errors.UserExists.WithDetails("用户名: "+userName), map[string]interface{}{})
		return
	}

	ctrl.renderService.JSONSuccess(c, map[string]interface{}{})
}

// ValidateEmail 验证邮箱是否可用
func (ctrl *UserController) ValidateEmail(c *gin.Context) {
	email, _ := c.GetPostForm("email")

	// 验证邮箱格式
	if !ctrl.validateEmailFormat(email) {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("邮箱格式不正确"), map[string]interface{}{})
		return
	}

	// 检查邮箱是否已存在
	isExistEmail, err := ctrl.service.IsExistEmail(email)
	if err != nil {
		ctrl.renderService.JSONError(c, err, map[string]interface{}{})
		return
	}
	if isExistEmail {
		ctrl.renderService.JSONError(c, errors.EmailExists.WithDetails("邮箱: "+email), map[string]interface{}{})
		return
	}

	ctrl.renderService.JSONSuccess(c, map[string]interface{}{})
}

// validateUserNameFormat 验证用户名格式
func (ctrl *UserController) validateUserNameFormat(userName string) bool {
	return pkg.IsValidUserName(userName)
}

// validateEmailFormat 验证邮箱格式
func (ctrl *UserController) validateEmailFormat(email string) bool {
	regexCheck := &pkg.RegexCheck{}
	return regexCheck.IsEmail(email)
}

// parseUserID 解析用户ID
func (ctrl *UserController) parseUserID(idStr string) (int64, error) {
	return ctrl.service.ParseUserID(idStr)
}