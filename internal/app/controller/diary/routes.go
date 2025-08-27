package diary

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册所有日记本相关路由
// 将所有diary相关路由聚合在一个地方，实现高内聚
// 包括公开访问路由、需要鉴权的表单路由和数据操作路由
//
// 参数:
//   - router: Gin路由引擎实例
//
// 路由分组:
//   - /diary: 日记相关公开路由
//   - /u/:uid: 用户个人主页路由
//   - /form/diary: 需要鉴权的日记本表单路由
//   - /diary, /daybook: 需要鉴权的日记本数据操作路由
func (ctrl *DiaryController) RegisterRoutes(router *gin.Engine) {
	// 日记相关路由组（无需鉴权）
	diaryGroup := router.Group("/diary")
	{
		diaryGroup.GET("/", ctrl.ShowDiary)
	}

	// 用户个人主页路由（公开访问，显示用户的日记本列表）
	router.GET("/u/:uid", ctrl.ShowUserProfile)
	router.GET("/u/:uid/:pageindex", ctrl.ShowUserProfile)

	// 需要鉴权的日记本表单路由
	authForm := router.Group("/form")
	authForm.Use(ctrl.auth.RequireAuth())
	{
		authForm.GET("/diary", ctrl.ShowDiaryForm)
	}

	// 需要鉴权的日记本数据操作路由（保持原路径）
	authAPI := router.Group("/")
	authAPI.Use(ctrl.auth.RequireAuth())
	{
		authAPI.POST("/diary", ctrl.ProcessDiary)
		authAPI.POST("/daybook", ctrl.ProcessDiary)
	}
}