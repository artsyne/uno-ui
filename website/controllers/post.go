package controllers

import (
	"net/http"

	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
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

	h := pkg.DefaultH(c)
	h["Title"] = "写日记"
	h["Active"] = "users"
	h["UserData"] = userData
	h["Diaries"] = diaries
	c.HTML(http.StatusOK, "post/write", h)
}

// 发布日记
func HandlePost(c *gin.Context) {
	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	// 序列化请求对象
	psotForm := &models.PostForm{}
	err := c.Bind(psotForm)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}
	var nowTime = pkg.GetNowTimeUnixMilli()

	//赋值post对象
	post := &models.Post{}

	// todo 需判断空间是否有权限发布日记
	post.DiaryID = psotForm.DiaryID
	post.IsTop = 0
	post.Title = psotForm.Title
	post.Body = psotForm.Body
	post.InfoType = 1 //默认为文本
	post.AddUserID = userData.UserId
	post.AddUserName = userData.UserName
	post.AddUTCTime = nowTime
	post.AddUserIP = c.ClientIP()
	post.UpdateUTCTime = nowTime
	post.Status = 1

	//插入数据库
	postID, err := services.AddPost(post)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	//创建返回对象
	h := pkg.ResponseH(c)
	h["Data"] = post
	h["Form"] = psotForm
	h["infoId"] = postID
	c.JSON(http.StatusOK, h)
}
