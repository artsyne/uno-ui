package comment

import (
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
//   - GET /api/post/:post_id/comments - 根据日记ID获取评论列表（RESTful风格）
func RegisterRoutes(router *gin.RouterGroup, ctrl *CommentController) {
	// 评论相关API路由组
	commentAPI := router.Group("/comment")
	{
		// 创建评论
		commentAPI.POST("/create", ctrl.CreateComment)
		
		// 删除评论
		commentAPI.POST("/delete", ctrl.DeleteComment)
		
		// 获取评论列表
		commentAPI.GET("/list", ctrl.GetComments)
	}

	// RESTful风格的路由
	// 根据日记ID获取评论列表
	router.GET("/post/:post_id/comments", ctrl.GetCommentsByPostIDFromPath)
}