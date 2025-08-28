package user

import (
	"strings"

	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

// UpdateProfile 更新用户基本资料
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析更新资料请求
	var req UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 调用service更新资料（包括时区和语言）
	if err := ctrl.service.UpdateUserProfile(authUser.UserId, req.FullName, req.Remark, req.TimeZone, req.Lang); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 设置新的cookie
	authUser.FullName = req.FullName
	authUser.TimeZone = req.TimeZone
	authUser.Lang = req.Lang
	
	ctrl.auth.SetLoginCookie(c, authUser)

	ctrl.renderService.JSONSuccess(c, nil)
}

// UpdateEmail 更新用户邮箱
func (ctrl *UserController) UpdateEmail(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析更新邮箱请求
	var req UpdateEmailRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 调用service更新邮箱
	if err := ctrl.service.UpdateUserEmail(authUser.UserId, req.Email); err != nil {
		ctrl.renderService.JSONError(c, errors.DatabaseError.WithDetails(err.Error()))
		return
	}

	ctrl.renderService.JSONSuccess(c, gin.H{"message": "邮箱更新成功"})
}

// UpdatePassword 修改用户密码
func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	// 读取用户认证信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 绑定请求数据
	var req ChangePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 调用service更改密码
	err := ctrl.service.ChangeUserPassword(authUser.UserId, req.CurrentPwd, req.NewPwd)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 更新 cookie
	ctrl.auth.SetLoginCookie(c, authUser)

	ctrl.renderService.JSONSuccess(c, gin.H{"message": "密码修改成功"})
}

// validateFullNameFormat 验证全名格式
func (ctrl *UserController) validateFullNameFormat(fullName string) bool {
	return len(strings.TrimSpace(fullName)) > 0 && len(fullName) <= 50
}