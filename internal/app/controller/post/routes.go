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
	diary := router.Group("/d")
	{
		diary.GET("/:diaryid", ctrl.ShowDiaryIndex)
		diary.GET("/:diaryid/:pageindex", ctrl.ShowDiaryIndex)
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

	// 需要鉴权的日记点赞操作路由
	api := router.Group("/api/post")
	api.Use(ctrl.auth.RequireAuth())
	{
		// 根据postIds获取日记社交统计信息
		api.POST("/stats", ctrl.GetPostsStats)
		// 根据postIds获取日记点赞状态
		api.POST("/like/status", ctrl.GetPostsLikeStatus)
		
		// 日记点赞/取消点赞
		apiAuth := api.Group("/").Use(ctrl.auth.RequireAuth())
		{
			apiAuth.POST("/like", ctrl.LikePost)
		}
	}
}