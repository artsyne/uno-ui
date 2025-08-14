package user

import (
	"net/http"
	"strings"

	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/user"
	errors "github.com/biuren/hebe/internal/utils/error"
	pkg "github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// Login 处理用户登录请求
func (ctrl *UserController) Login(c *gin.Context) {
	// 解析登录表单数据，创建 LoginRequest 对象
	req := &user.LoginRequest{}

	if err := c.Bind(req); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 验证输入参数长度
	if len(req.LoginName) < 6 || len(req.Pwd) < 6 {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails("用户名和密码长度至少6位"))
		return
	}

	// 验证用户名和密码，登录名统一转为小写
	req.LoginName = strings.ToLower(req.LoginName)
	user, err := ctrl.service.LoginWithUsernameOrEmail(req.LoginName, req.Pwd)
	if err != nil {
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
		ctrl.renderService.Redirect(c, "/")
		return
	}

	// 解析注册表单数据
	user := &user.User{}
	if err := c.Bind(user); err != nil {
		ctrl.renderService.JSONError(c, errors.ValidationError.WithDetails(err.Error()))
		return
	}

	// 设置客户端IP
	user.LastIp = c.ClientIP()

	// 调用服务层注册用户
	if regErr := ctrl.service.RegisterUser(user); regErr != nil {
		ctrl.renderService.JSONError(c, regErr)
		return
	}

	// 返回注册成功响应
	responseData := map[string]interface{}{
		"UserId": user.UserId,
	}
	ctrl.renderService.JSONSuccess(c, responseData)
}

// Logout 处理登出请求
func (ctrl *UserController) Logout(c *gin.Context) {
	// 清除登录cookie
	pkg.DelCookie(c.Writer, "ud")
	c.Redirect(http.StatusFound, "/passport/login")
}