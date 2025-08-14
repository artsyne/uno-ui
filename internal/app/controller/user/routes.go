package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有用户相关路由
// 保持与原系统完全一致的路由路径
func (ctrl *UserController) RegisterRoutes(router *gin.Engine) {
	// 用户认证相关路由（无需鉴权）
	router.GET("/register", ctrl.ShowRegisterPage)
	router.POST("/register", ctrl.Register)
	router.POST("/check-username", ctrl.ValidateUsername)
	router.POST("/check-email", ctrl.ValidateEmail)

	router.GET("/login", ctrl.ShowLoginPage)
	router.POST("/login", ctrl.Login)

	router.POST("/loginout", ctrl.Logout)

	router.GET("/findpwd", ctrl.ShowPasswordResetPage)
	router.POST("/findpwd", ctrl.RequestPasswordReset)

	router.GET("/resetpwd", ctrl.ShowPasswordResetForm)
	router.POST("/resetpwd", ctrl.ResetPassword)

	router.GET("/activateuser", ctrl.ShowUserActivationPage)

	// 用户设置相关路由（需要鉴权）
	authGroup := router.Group("/u")
	authGroup.Use(ctrl.auth.RequireAuth())
	{
		authGroup.GET("/settings", ctrl.ShowUserSettingsPage)
		authGroup.POST("/settings/profile", ctrl.UpdateProfile)
		authGroup.POST("/settings/email", ctrl.UpdateEmail)
		authGroup.POST("/settings/password", ctrl.UpdatePassword)
	}
}