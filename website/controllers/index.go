package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/pkg"
)

// RenderIndex : 渲染网站首页
func RenderIndex(c *gin.Context) {
	h := pkg.DefaultH(c)
	h["Title"] = "List of users"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "index", h)
}
