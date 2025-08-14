package user

import (
	"strings"

	"github.com/biuren/hebe/internal/domain/user"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/gin-gonic/gin"
)

// UpdateProfile 更新用户基本资料
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	// 解析更新资料请求
	req := &user.UpdateProfileRequest{}
	if err := c.Bind(req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 设置用户ID,只能修改登录用户 profile
	req.UserId = authUser.UserId

	// 验证全名格式
	if req.FullName != "" && !ctrl.validateFullNameFormat(req.FullName) {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("全名格式不正确"))
		return
	}

	// 调用service更新资料
	if err := ctrl.service.UpdateUserProfile(req.UserId, req.FullName, req.Remark); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 设置新的cookie
	authUser.FullName = req.FullName
	ctrl.auth.SetLoginCookie(c, authUser)

	ctrl.renderService.JSONSuccess(c, nil)
}

// UpdateEmail 更新用户邮箱
func (ctrl *UserController) UpdateEmail(c *gin.Context) {
	// 获取用户登录信息
	authUser := ctrl.auth.MustGetCurrentUser(c)

	req := &user.UpdateEmailRequest{}

	req.UserId = authUser.UserId

	// 获取表单数据
	req.Email = c.PostForm("email")

	// 调用service更新邮箱
	if err := ctrl.service.UpdateUserEmail(req.UserId, req.Email); err != nil {
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
	req := &user.ChangePasswordRequest{}
	if err := c.Bind(req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 验证输入字段不能为空
	if req.CurrentPwd == "" || req.NewPwd == "" || req.ConfirmPwd == "" {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("当前密码、新密码和确认密码不能为空"))
		return
	}

	// 判断输入密码是否正确
	if req.NewPwd != req.ConfirmPwd {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("两次输入密码不一致"))
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