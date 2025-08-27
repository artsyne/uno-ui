package controller

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/config"
)

// IndexController 基础控制器，包含公共服务
// 提供网站首页渲染功能，管理配置和视图渲染服务
//
// 字段:
//   - config: 应用程序配置实例
//   - renderService: 视图渲染服务实例
type IndexController struct {
	config        *config.Config
	renderService *view.ViewService
}

// NewIndexController 创建控制器实例
// 初始化基础控制器，注入必要的依赖服务
//
// 参数:
//   - config: 应用程序配置实例
//   - renderService: 视图渲染服务实例
//
// 返回:
//   - *IndexController: 初始化完成的控制器实例
func NewIndexController(config *config.Config, renderService *view.ViewService) *IndexController {
	return &IndexController{
		config:        config,
		renderService: renderService,
	}
}

// RenderIndex 渲染网站首页
// 处理首页请求，创建页面数据并渲染首页模板
// 设置当前激活的导航项为"users"
//
// 参数:
//   - c: Gin上下文对象，包含请求和响应信息
func (ctrl *IndexController) RenderIndex(c *gin.Context) {
	data := ctrl.renderService.CreatePageData(c, "home.title")
	data.Data["Active"] = "users"
	ctrl.renderService.HTML(c, "index", data)
}
