package post

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册所有日记相关路由
// 将所有post相关路由聚合在一个地方，实现高内聚
// 包括公开访问的日记展示路由和需要鉴权的表单及数据操作路由
//
// 参数:
//   - router: Gin路由引擎实例
//
// 路由分组:
//   - /d/:diaryid: 日记本页面路由（公开访问）
//   - /p/:postid: 单个日记详情页面（公开访问）
//   - /form/post: 需要鉴权的日记表单路由
//   - /post, /info: 需要鉴权的日记数据操作路由
func (ctrl *PostController) RegisterRoutes(router *gin.Engine) {
	// 日记本页面路由组（无需鉴权，支持公开访问）
	// 显示日记本下的日记列表
	daybook := router.Group("/d")
	{
		daybook.GET("/:diaryid", ctrl.ShowDiaryIndex)
		daybook.GET("/:diaryid/:pageindex", ctrl.ShowDiaryIndex)
	}

	// 显示单个日记详情页面
	router.GET("/p/:postid", ctrl.ShowPost)

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

	// 需要鉴权的日记删除路由
	authDelete := router.Group("/post")
	authDelete.Use(ctrl.auth.RequireAuth())
	{
		authDelete.POST("/delete", ctrl.DeletePost)
	}
}