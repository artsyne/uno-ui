package controllers

import (
	"net/http"

	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/utils"
	"github.com/gin-gonic/gin"
)

// renderDiaryWrite ： 渲染日记编辑
func RenderDiaryWrite(c *gin.Context) {
	h := utils.DefaultH(c)
	h["Title"] = "写日记"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "diary/write_form", h)
}

// 发布日记
func HandleInfo(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	// 序列化请求对象
	infoForm := &models.InfoForm{}
	err := c.Bind(infoForm)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	info := &models.Info{}
	var nowTime = utils.GetNowTimeUnixMilli()
	//赋值info对象
	info.SpaceID = infoForm.SpaceID
	info.IsTop = 0
	info.Title = infoForm.Title
	info.Summary = ""
	info.Body = infoForm.Body
	info.InfoType = 1
	info.AddUserID = userData.UserId
	info.AddUserName = userData.UserName
	info.AddUTCTime = nowTime
	info.AddUserIP = c.ClientIP()
	info.UpdateUTCTime = nowTime
	info.Status = 1

	//插入数据库
	infoID, err := info.Insert()
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	//日记本信息数量加一
	err = models.IncrementInfoNum(info.SpaceID)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	//创建返回对象
	h := utils.ResponseH(c)
	h["Data"] = info
	h["Form"] = infoForm
	h["infoId"] = infoID
	c.JSON(http.StatusOK, h)
}
