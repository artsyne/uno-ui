package user

import (
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/errors"
	pkg "github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// RequestPasswordReset 处理找回密码请求
func (ctrl *UserController) RequestPasswordReset(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[FindPasswordRequest](c)

	// 转换为Domain层结构体
	domainReq := &user.FindPasswordRequest{
		Email: req.Email,
	}

	// 获取用户信息
	userInfo, err := ctrl.service.GetUserByEmail(domainReq.Email)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.UserNotFound.WithDetails(map[string]interface{}{"email": req.Email, "reason": "user_not_found_by_email"}))
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
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[ResetPasswordRequest](c)

	// 转换为Domain层结构体
	domainReq := &user.ResetPasswordRequest{
		UID:    req.UID,
		Sign:   req.Sign,
		NewPwd: req.NewPwd,
	}

	// 解析用户ID
	userID, err := ctrl.service.ParseUserID(domainReq.UID)
	if err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(map[string]interface{}{"reason": "invalid_user_id_format"}))
		return
	}

	// 检查验证码
	valid, validCodeID := ctrl.service.CheckValidCode(userID, domainReq.Sign, ctrl.config.ValidCode.TypeResetPwd)
	if !valid {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(map[string]interface{}{"reason": "invalid_or_expired_code"}))
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