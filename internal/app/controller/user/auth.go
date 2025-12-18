package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/biuren/hebe/internal/app/middleware"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/errors"
	pkg "github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// Login 处理用户登录请求
func (ctrl *UserController) Login(c *gin.Context) {
	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[LoginRequest](c)

	// 转换为Domain层结构体
	domainReq := &user.LoginRequest{
		LoginName: strings.ToLower(req.LoginName),
		Pwd:       req.Pwd,
		AccessUrl: req.AccessUrl,
	}

	// 验证用户名和密码
	user, err := ctrl.service.LoginWithUsernameOrEmail(domainReq.LoginName, domainReq.Pwd)
	if err != nil {
		if ctrl.isHtmxRequest(c) {
			ctrl.htmxAlert(c, http.StatusUnauthorized, err.Error())
			return
		}
		ctrl.renderService.JSONError(c, errors.PasswordIncorrect.WithDetails(err.Error()))
		return
	}

	// 创建登录会话并设置Cookie
	authUser := &types.AuthUser{
		UserId:   user.UserId,
		UserName: user.UserName,
		FullName: user.FullName,
		Email:    user.Email,
		Lang:     user.Lang,
		TimeZone: user.TimeZone,
	}

	// 设置登录Cookie
	ctrl.auth.SetLoginCookie(c, authUser)

	if ctrl.isHtmxRequest(c) {
		c.Header("HX-Redirect", fmt.Sprintf("/u/%d", user.UserId))
		c.Status(http.StatusOK)
		return
	}

	// 返回登录成功响应
	responseData := map[string]interface{}{
		"IsLogin": true,
		"UserId":  user.UserId,
	}
	ctrl.renderService.JSONSuccess(c, responseData)
}

// Register 处理用户注册请求
func (ctrl *UserController) Register(c *gin.Context) {
	// 如果已登录则跳转首页
	if isAuth, _ := ctrl.auth.IsAuthenticated(c); isAuth {
		// 检查是否为API或AJAX请求
		if ctrl.isHtmxRequest(c) ||
			c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
			strings.Contains(c.GetHeader("Accept"), "application/json") {
			ctrl.renderService.JSONError(c, errors.OperationNotAllowed.WithDetails("User already logged in"))
			return
		}
		ctrl.renderService.Redirect(c, "/")
		return
	}

	// 获取已验证的请求数据
	req := middleware.MustGetValidatedRequest[RegisterRequest](c)

	// 转换为Domain层结构体
	domainUser := &user.User{
		UserName: req.UserName,
		FullName: req.FullName,
		Email:    req.Email,
		Pwd:      req.Pwd,
		TimeZone: req.TimeZone,
		Lang:     req.Lang,
		LastIp:   c.ClientIP(),
	}

	// 调用服务层注册用户
	if regErr := ctrl.service.RegisterUser(domainUser); regErr != nil {
		if ctrl.isHtmxRequest(c) {
			ctrl.htmxAlert(c, http.StatusBadRequest, regErr.Error())
			return
		}
		ctrl.renderService.JSONError(c, regErr)
		return
	}

	if ctrl.isHtmxRequest(c) {
		c.Header("HX-Redirect", "/login")
		c.Status(http.StatusOK)
		return
	}

	// 返回注册成功响应
	responseData := map[string]interface{}{
		"UserId": domainUser.UserId,
	}
	ctrl.renderService.JSONSuccess(c, responseData)
}

// Logout 处理登出请求
func (ctrl *UserController) Logout(c *gin.Context) {
	// 清除登录cookie
	pkg.DelCookie(c.Writer, "ud")
	c.Redirect(http.StatusFound, "/index")
}
