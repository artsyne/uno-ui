package controllers

import (
	"chinadiary.com/website/models"
	"net/http"

	"chinadiary.com/website/helpers"
	"github.com/gin-gonic/gin"
)

// RenderMyIndex 个人空间搜
func RenderMyIndex(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "管理"
	h["Active"] = "users"
	h["IsAdminPage"] = true
	
	//读取用户认证信息
	userData := h["UserData"].(models.Login)
	spaces,_ := models.GetSpaceByUserID(userData.UserId)

	h["Spaces"] = spaces

	c.HTML(http.StatusOK, "my/index", h)
}
