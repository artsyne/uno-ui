package user

import (
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

// UpdateProfile 更新用户基本资料
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[UpdateProfileRequest](c)

	// 校验描述长度
	if len(req.Remark) > ctrl.config.App.MaxDescriptionLength {
		ctrl.renderService.JSONError(c, errors.ValidationError)
		return
	}

	// 调用service更新资料（包括时区和语言）
	if err := ctrl.service.UpdateUserProfile(authUser.UserId, req.DisplayName, req.Remark, req.TimeZone, req.Lang); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 设置新的cookie
	authUser.DisplayName = req.DisplayName
	authUser.TimeZone = req.TimeZone
	authUser.Lang = req.Lang
	
	ctrl.auth.SetLoginCookie(c, authUser)

	ctrl.renderService.JSONSuccessWithMsg(c, nil, "Profile updated successfully")
}

// UpdateHandle 更新用户Handle
func (ctrl *UserController) UpdateHandle(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[UpdateHandleRequest](c)

	// 调用service更新Handle
	if err := ctrl.service.UpdateUserHandle(authUser.UserId, req.Handle); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 更新cookie中的handle
	authUser.Handle = req.Handle
	ctrl.auth.SetLoginCookie(c, authUser)

	msg := ctrl.renderService.GetMessage("user.settings.username_updated", ctrl.renderService.GetLocaleFromContext(c))
	ctrl.renderService.JSONSuccessWithMsg(c, nil, msg)
}

// UpdateEmail 更新用户邮箱
func (ctrl *UserController) UpdateEmail(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[UpdateEmailRequest](c)

	// 调用service更新邮箱
	if err := ctrl.service.UpdateUserEmail(authUser.UserId, req.Email); err != nil {
		ctrl.renderService.JSONError(c, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()}))
		return
	}

	ctrl.renderService.JSONSuccessWithMsg(c, nil, "Email updated successfully")
}

// UpdatePassword 修改用户密码
func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	// 读取用户认证信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[ChangePasswordRequest](c)

	// 调用service更改密码
	err := ctrl.service.ChangeUserPassword(authUser.UserId, req.CurrentPwd, req.NewPwd)
	if err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 更新 cookie
	ctrl.auth.SetLoginCookie(c, authUser)

	msg := ctrl.renderService.GetMessage("user.settings.password_changed", ctrl.renderService.GetLocaleFromContext(c))
	ctrl.renderService.JSONSuccessWithMsg(c, nil, msg)
}
