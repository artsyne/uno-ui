package user

import (
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有用户相关路由
// 保持与原系统完全一致的路由路径，包括用户认证和设置管理功能
//
// 参数:
//   - router: Gin路由引擎实例
//
// 路由分组:
//   - 认证路由（无需鉴权）: 注册、登录、密码重置、用户激活等
//   - 设置路由（需要鉴权）: 个人信息、邮箱、密码修改等
//
// 功能:
//   - 用户注册和登录流程
//   - 用户名和邮箱验证
//   - 密码重置和找回
//   - 用户激活
//   - 个人设置管理
func (ctrl *UserController) RegisterRoutes(router *gin.Engine) {
	// 认证限流中间件
	authLimit := middleware.AuthRateLimit(ctrl.config)

	// 用户认证相关路由（无需鉴权）
	router.GET("/register", ctrl.ShowRegisterPage)
	router.POST("/register", authLimit, middleware.Validate[RegisterRequest](), ctrl.Register)
	router.POST("/check-handle", authLimit, middleware.Validate[CheckHandleRequest](), ctrl.ValidateHandle)
	router.POST("/check-email", authLimit, middleware.Validate[CheckEmailRequest](), ctrl.ValidateEmail)

	router.GET("/login", ctrl.ShowLoginPage)
	router.POST("/login", authLimit, middleware.Validate[LoginRequest](), ctrl.Login)

	router.GET("/loginout", ctrl.Logout)

	router.GET("/findpwd", ctrl.ShowPasswordResetPage)
	router.POST("/findpwd", authLimit, middleware.Validate[FindPasswordRequest](), ctrl.RequestPasswordReset)

	router.GET("/resetpwd", ctrl.ShowPasswordResetForm)
	router.POST("/resetpwd", authLimit, middleware.Validate[ResetPasswordRequest](), ctrl.ResetPassword)

	router.GET("/activateuser", ctrl.ShowUserActivationPage)

	// 用户设置相关路由（需要鉴权）
	authGroup := router.Group("/u")
	authGroup.Use(ctrl.auth.RequireAuth())
	{
		authGroup.GET("/settings", ctrl.ShowUserSettingsPage)
		authGroup.POST("/settings/profile", middleware.Validate[UpdateProfileRequest](), ctrl.UpdateProfile)
		authGroup.POST("/settings/handle", middleware.Validate[UpdateHandleRequest](), ctrl.UpdateHandle)
		authGroup.POST("/settings/email", middleware.Validate[UpdateEmailRequest](), ctrl.UpdateEmail)
		authGroup.POST("/settings/password", middleware.Validate[ChangePasswordRequest](), ctrl.UpdatePassword)
	}
}
