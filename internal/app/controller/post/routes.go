package post

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册所有日记相关路由
// 将所有post相关路由聚合在一个地方，实现高内聚
func (ctrl *PostController) RegisterRoutes(router *gin.Engine) {
	// 日记本页面路由组（无需鉴权，支持公开访问）
	// 显示日记本下的日记列表
	daybook := router.Group("/d")
	{
		daybook.GET("/:spaceid", ctrl.ShowDiaryIndex)
		daybook.GET("/:spaceid/:pageindex", ctrl.ShowDiaryIndex)
	}

	// 需要鉴权的日记表单路由
	authForm := router.Group("/form")
	authForm.Use(ctrl.auth.RequireAuth())
	{
		authForm.GET("/post", ctrl.ShowPostForm)
	}

	// 需要鉴权的日记数据操作路由（保持原路径）
	authAPI := router.Group("/")
	authAPI.Use(ctrl.auth.RequireAuth())
	{
		authAPI.POST("/post", ctrl.ProcessPost)
		authAPI.POST("/info", ctrl.ProcessPost)
	}
}