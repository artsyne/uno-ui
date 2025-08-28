package user

import (
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/errors"
	pkg "github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// RequestPasswordReset 处理找回密码请求
func (ctrl *UserController) RequestPasswordReset(c *gin.Context) {
	// 解析请求数据
	var req FindPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 转换为Domain层结构体
	domainReq := &user.FindPasswordRequest{
		Email: req.Email,
	}

	// 获取用户信息
	userInfo, err := ctrl.service.GetUserByEmail(domainReq.Email)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.UserNotFound.WithDetails("邮箱: "+req.Email))
		return
	}

	// 创建重置密码验证码
	if err := ctrl.service.CreateResetPasswordCode(userInfo.UserId); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 发送重置密码邮件
	if err := ctrl.service.SendResetPasswordEmail(userInfo); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	ctrl.renderService.JSONSuccess(c, nil)
}

// ResetPassword 处理重置密码请求
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	// 解析重置密码请求
	var req ResetPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 转换为Domain层结构体
	domainReq := &user.ResetPasswordRequest{
		UID:    req.UID,
		Sign:   req.Sign,
		NewPwd: req.NewPwd,
	}

	// 解析用户ID
	userID, err := ctrl.service.ParseUserID(domainReq.UID)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("用户ID格式错误"))
		return
	}

	// 检查验证码
	valid, validCodeID := ctrl.service.CheckValidCode(userID, domainReq.Sign, ctrl.config.ValidCode.TypeResetPwd)
	if !valid {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("验证码无效或已过期"))
		return
	}

	// 重置密码
	if err := ctrl.service.ResetUserPassword(userID, domainReq.NewPwd); err != nil {
		ctrl.renderService.JSONError(c, err)
		return
	}

	// 删除验证码
	ctrl.service.DeleteValidCode(validCodeID)

	ctrl.renderService.JSONSuccess(c, nil)
}

// validatePasswordFormat 验证密码格式
func (ctrl *UserController) validatePasswordFormat(password string) bool {
	return pkg.IsPwd(password)
}