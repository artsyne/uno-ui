package controllers

import (
	"net/http"

	"chinadiary.com/website/models"

	gobalerr "chinadiary.com/website/config"
	"chinadiary.com/website/helpers"
	"github.com/gin-gonic/gin"
)

// RenderDaybookIndex 渲染日记本首页
func RenderDaybookIndex(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "日记本"
	h["Active"] = "users"
	h["SpaceID"] = c.Param("spaceid")

	

	c.HTML(http.StatusOK, "daybook/index", h)
}

//增加日记本表单渲染
func RenderDaybookForm(c *gin.Context){
	h:= helpers.DefaultH(c)

	c.HTML(http.StatusOK,"form/daybook",h)
}

// RenderDiary : 渲染日记详细信息
func RenderDiary(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "我的日记本"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "daybook/diary", h)
}

// 保存日记本信息
func DoSpace(c *gin.Context) {
	// 序列化用户输入对象
	space := &models.Space{}
	err := c.Bind(space)

	// 校验参数，标题必须填写，设置类型必须为1或者2，主题id必须在1-16之间
	if err != nil || len(space.Title) < 1 || (space.ViewSetting != 1 &&
		space.ViewSetting != 2 || space.ThemeID < 1 || space.ThemeID > 16) {
		gobalerr.ResponeErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
		return
	}

	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	var nowTime = helpers.GetNowTime()
	//其他设置默认值
	space.AccessPwd = ""
	space.AccessURL = ""
	space.AddTime = nowTime
	space.AddUserID = userData.UserId
	space.AddUserName = userData.UserName
	space.CoverMidURL = ""
	space.CoverURL = ""
	space.IsShowComment = 1
	space.IsShowLike = 1
	space.IsShowView = 1
	space.LikeNum = 0
	space.MemberNum = 0
	space.Status = 1
	space.ViewNum = 0
	space.InfoNum = 0
	space.LastAddTime = nowTime
	space.LastAddUserID = userData.UserId
	space.LastAddUserName = userData.UserName

	//保存日记本
	if action, _ := c.GetPostForm("action"); action == "My.insertDaybook" {
		spaceID, err := space.Insert() // 插入对象
		if err != nil {
			gobalerr.ResponeErr(c, gobalerr.PARSE_INPUT_PARAM_ERR)
			return
		}
		space.SpaceID = spaceID
	}

	//创建返回对象
	h := helpers.ResponseH(c)
	h["Data"] = space
	c.JSON(http.StatusOK, h)
}
