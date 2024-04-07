package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/utils"
	"github.com/biuren/hebe/models"
)

// RenderMemberIndex : 渲染会员首页
func RenderMemberIndex(c *gin.Context) {
	h := utils.DefaultH(c)

	//获取url 中的 uid
	uid, _ := strconv.ParseInt(c.Param("uid"), 10, 64)

	h["Title"] = "个人主页"
	h["Active"] = "users"

	//读取用户认证信息
	userData := h["UserData"].(models.Login)

	//返回判断是否是日记owner
	isOnwer := false
	if userData.UserId == uid {
		isOnwer = true
	}

	//获取日记本列表
	spaces,_ := models.GetSpaceByUserID(uid)

	h["IsOwner"] = isOnwer
	h["Spaces"] = spaces

	//h["UserData"] = c
	c.HTML(http.StatusOK, "member/index", h)
}

// RenderMemberSetting: 渲染设置页面
func RenderMemberSettings(c *gin.Context){
	h := utils.DefaultH(c)

	//获取url 中的 uid
	uid, _ := strconv.ParseInt(c.Param("uid"), 10, 64)

	h["Title"] = "设置"
	h["Active"] = "setting"

	//读取用户认证信息
	userData := h["UserData"].(models.Login)

	//返回判断是否是owner
	isOnwer := false
	if userData.UserId == uid {
		isOnwer = true

		//如果是用户本人，将获取用户信息
		userInfo, _ := models.GetUserByUserId(userData.UserId)
		
		//返回用户对象到模版
		h["UserInfo"] = userInfo
	}

	h["IsOwner"] = isOnwer

	//h["UserData"] = c
	c.HTML(http.StatusOK, "member/settings", h)
}