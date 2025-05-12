package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/pkg"
)

// RenderMemberIndex : 渲染会员首页
func RenderMemberIndex(c *gin.Context) {
	h := pkg.DefaultH(c)

	//获取url 中的 uid
	uid, _ := strconv.ParseInt(c.Param("uid"), 10, 64)

	h["Title"] = "个人主页"
	h["Active"] = "users"

	//读取用户认证信息
	userData, ok := h["UserData"].(models.Login)
	if !ok {
		// 用户未登录或登录信息无效，重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//返回判断是否是日记owner
	isOnwer := false
	if userData.UserId == uid {
		isOnwer = true
	}

	//获取日记本列表
	// 获取 URL 中的 pageIndex
	pageIndex, err := strconv.Atoi(strings.TrimSpace(c.Param("pageindex")))
	if err != nil || pageIndex <= 0 {
		pageIndex = 1
	}

	// 日记总数
	diariesNum := 0
	pageSize := configs.GetConfig().App.PageSize

	// 分页显示日记列表
	h["Diaries"], diariesNum, _ = models.ListDiariesByUserID(uid, pageIndex, pageSize)

	// 初始化分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/u/%d/%%d", uid),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	h["IsOwner"] = isOnwer
	h["UserData"] = userData
	h["PageLinks"] = pageLinks
	h["DiariesNum"] = diariesNum
	h["PageSize"] = pageSize

	c.HTML(http.StatusOK, "member/index", h)
}

// RenderMemberSetting: 渲染设置页面
func RenderMemberSettings(c *gin.Context) {
	h := pkg.DefaultH(c)

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
	h["UserData"] = userData

	//h["UserData"] = c
	c.HTML(http.StatusOK, "member/settings", h)
}
