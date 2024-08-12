package controllers

import (
	"net/http"

	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
	"github.com/gin-gonic/gin"
)

// renderDiaryWrite ： 渲染日记编辑界面
func RenderDiaryWrite(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	//读取用户所有日记本
	spaces, err := models.ListDaybooksForDiaryByUserID(userData.UserId)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	h := utils.DefaultH(c)
	h["Title"] = "写日记"
	h["Active"] = "users"
	h["UserData"] = userData
	h["Spaces"] = spaces
	c.HTML(http.StatusOK, "diary/write", h)
}

// 发布日记
func HandleInfo(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	//获取配置信息
	config := settings.GetConfig()

	// 序列化请求对象
	infoForm := &models.InfoForm{}
	err := c.Bind(infoForm)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	var nowTime = utils.GetNowTimeUnixMilli()

	//赋值info对象
	info := &models.Info{}

	// todo 需判断空间是否有权限发布日记

	info.SpaceID = infoForm.SpaceID
	info.IsTop = 0
	info.Title = infoForm.Title
	//截取配置长度字符串，末尾增加 suffix
	info.Summary = utils.TruncateString(infoForm.Body, config.SummaryMax, config.Suffix)
	info.Body = infoForm.Body
	info.InfoType = 1 //默认为文本
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
	err = models.IncrementInfoNum(info.SpaceID, userData.UserId, userData.UserName, nowTime)
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
