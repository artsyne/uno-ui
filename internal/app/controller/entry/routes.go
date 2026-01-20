package entry

import (
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有日记相关路由
// 将所有entry相关路由聚合在一个地方，实现高内聚
// 包括公开访问的日记展示路由和需要鉴权的表单及数据操作路由
//
// 参数:
//   - router: Gin路由引擎实例
//
// 路由分组:
//   - /diary/:diaryid: 日记本页面路由（公开访问）
//   - /entry/:entryid: 单个日记详情页面（公开访问）
//   - /form/entry: 需要鉴权的日记表单路由
//   - /entry: 需要鉴权的日记数据操作路由
func (ctrl *EntryController) RegisterRoutes(router *gin.Engine) {

	// 日记本页面路由组（无需鉴权，支持公开访问）
	// 显示日记本下的日记列表
	diary := router.Group("/diary")
	{
		diary.GET("/:diaryid", ctrl.ShowDiaryIndex)
		diary.GET("/:diaryid/:pageindex", ctrl.ShowDiaryIndex)
	}

	// 显示单个日记详情页面
	router.GET("/entry/:entryid", ctrl.ShowEntry)

	// 需要鉴权的日记表单路由
	authForm := router.Group("/form")
	authForm.Use(ctrl.auth.RequireAuth())
	{
		authForm.GET("/entry", ctrl.ShowEntryForm)
	}

	// 需要鉴权的日记数据操作路由
	authAPI := router.Group("/")
	authAPI.Use(ctrl.auth.RequireAuth())
	{
		authAPI.POST("/entry", middleware.Validate[CreateEntryRequest](), ctrl.CreateEntry)
		authAPI.PUT("/entry/:id", middleware.Validate[UpdateEntryRequest](), ctrl.UpdateEntry)
	}

	// 需要鉴权的日记删除路由
	authDelete := router.Group("/entry")
	authDelete.Use(ctrl.auth.RequireAuth())
	{
		authDelete.DELETE("/:id", middleware.Validate[DeleteEntryRequest](), ctrl.DeleteEntry)
	}

	// 日记点赞相关路由
	api := router.Group("/api/entry")
	{
		// 根据entryIds获取日记社交统计信息（无需鉴权）
		api.POST("/stats", middleware.Validate[GetEntriesStatsRequest](), ctrl.GetEntriesStats)
		// 根据entryIds获取日记点赞状态（无需鉴权，内部处理匿名用户）
		api.POST("/like/status", middleware.Validate[GetEntryLikeStatusRequest](), ctrl.GetEntriesLikeStatus)

		// 日记点赞/取消点赞（需要鉴权）
		apiAuth := api.Group("/")
		apiAuth.Use(ctrl.auth.RequireAuth())
		{
			apiAuth.POST("/like", middleware.Validate[LikeEntryRequest](), ctrl.LikeEntry)
		}
	}
}
