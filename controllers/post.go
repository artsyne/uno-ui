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
func RenderPostWrite(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	//读取用户所有日记本
	diaries, err := models.ListDiariesForWriteByUserID(userData.UserId)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	h := utils.DefaultH(c)
	h["Title"] = "写日记"
	h["Active"] = "users"
	h["UserData"] = userData
	h["Spaces"] = diaries
	c.HTML(http.StatusOK, "diary/write", h)
}

// 发布日记
func HandlePost(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	//获取配置信息
	config := settings.GetConfig()

	// 序列化请求对象
	psotForm := &models.PostForm{}
	err := c.Bind(psotForm)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	var nowTime = utils.GetNowTimeUnixMilli()

	//赋值post对象
	post := &models.Post{}

	// todo 需判断空间是否有权限发布日记

	post.SpaceID = psotForm.SpaceID
	post.IsTop = 0
	post.Title = psotForm.Title
	//截取配置长度字符串，末尾增加 suffix
	post.Summary = utils.TruncateString(psotForm.Body, config.SummaryMax, config.Suffix)
	post.Body = psotForm.Body
	post.InfoType = 1 //默认为文本
	post.AddUserID = userData.UserId
	post.AddUserName = userData.UserName
	post.AddUTCTime = nowTime
	post.AddUserIP = c.ClientIP()
	post.UpdateUTCTime = nowTime
	post.Status = 1

	//插入数据库
	postID, err := post.Insert()
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	//日记本信息数量加一
	err = models.IncrementPostNum(post.SpaceID, userData.UserId, userData.UserName, nowTime)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	//创建返回对象
	h := utils.ResponseH(c)
	h["Data"] = post
	h["Form"] = psotForm
	h["infoId"] = postID
	c.JSON(http.StatusOK, h)
}
