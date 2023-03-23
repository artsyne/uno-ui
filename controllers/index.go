package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/helpers"
)

// RenderIndex : 渲染网站首页
func RenderIndex(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "List of users"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "index", h)
}
