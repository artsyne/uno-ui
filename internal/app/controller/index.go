package controller

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/utils/config"
)

// IndexController 基础控制器，包含公共服务
type IndexController struct {
	config        *config.Config
	renderService *view.ViewService
}

// NewIndexController 创建控制器实例
func NewIndexController(config *config.Config, renderService *view.ViewService) *IndexController {
	return &IndexController{
		config:        config,
		renderService: renderService,
	}
}

// RenderIndex 渲染网站首页
func (ctrl *IndexController) RenderIndex(c *gin.Context) {
	data := ctrl.renderService.CreatePageData(c, "home.title")
	data.Data["Active"] = "users"
	ctrl.renderService.HTML(c, "index.html", data)
}
