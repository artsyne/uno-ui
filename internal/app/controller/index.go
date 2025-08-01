package controller

import (
	"github.com/biuren/hebe/internal/app/view"
	"github.com/gin-gonic/gin"
)

// Controller 基础控制器，包含公共服务
type Controller struct {
	renderService view.ViewServiceInterface
}

// NewController 创建控制器实例
func NewController(renderService view.ViewServiceInterface) *Controller {
	return &Controller{
		renderService: renderService,
	}
}

// RenderIndex 渲染网站首页
func (ctrl *Controller) RenderIndex(c *gin.Context) {
	data := ctrl.renderService.CreatePageData(c, "home.title")
	data.Data["Active"] = "users"
	ctrl.renderService.HTML(c, "index.html", data)
}
