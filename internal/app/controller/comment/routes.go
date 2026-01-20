package comment

import (
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论相关的路由
// 定义所有与评论相关的API端点
//
// 参数:
//   - router: Gin路由组
//   - ctrl: 评论控制器实例
//
// 路由定义:
//   - POST /api/comment/create - 创建评论
//   - POST /api/comment/delete - 删除评论
//   - GET /api/comment/list - 获取评论列表
//   - GET /api/entry/:entry_id/comments - 根据日记ID获取评论列表（RESTful风格）
func RegisterRoutes(router *gin.RouterGroup, ctrl *CommentController) {
	// 评论相关API路由组
	commentAPI := router.Group("/comment")
	{
		// 获取评论列表 (无需鉴权)
		commentAPI.GET("/list", middleware.Validate[GetCommentsRequest](), ctrl.GetComments)

		// 需要鉴权的路由
		authGroup := commentAPI.Group("/")
		authGroup.Use(ctrl.auth.RequireAuth())
		{
			// 创建评论
			authGroup.POST("/create", middleware.Validate[CreateCommentRequest](), ctrl.CreateComment)

			// 删除评论
			authGroup.POST("/delete", middleware.Validate[DeleteCommentRequest](), ctrl.DeleteComment)
		}
	}

	// RESTful风格的路由
	// 根据日记ID获取评论列表
	router.GET("/entry/:entry_id/comments", middleware.Validate[GetCommentsByPathRequest](), ctrl.GetCommentsByEntryIDFromPath)
}
