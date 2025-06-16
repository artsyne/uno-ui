package controllers

import (
	"net/http"
	"strconv"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/helpers"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// RenderPostForm: 渲染日本编辑界面
func RenderPostForm(c *gin.Context) {
	//获取用户登录信息
	loginUser := helpers.Auth.MustGetCurrentUser(c)
	h := pkg.DefaultH(c)

	// 获取action 和 post_id
	action := c.Query("action")
	if action == "post.edit" {
		postId, err := strconv.ParseInt(c.Query("post_id"), 10, 64)
		if err != nil || postId <= 0 {
			errors.ResponseErr(c, errors.ParseInputParamErr)
			return
		}

		// 获取日记信息
		post, err := services.GetPostByID(postId)

		//判断是否有权限修改
		if loginUser.UserId != post.AddUserID || err != nil {
			errors.ResponseErr(c, errors.InputParamErr)
			return
		}
		h["Post"] = post
	}

	h["Active"] = "users"
	h["UserData"] = loginUser
	c.HTML(http.StatusOK, "form/post", h)
}

// HandlePost 处理日记的创建和编辑
func HandlePost(c *gin.Context) {
	// 序列化用户输入对象
	postForm := &models.PostForm{}
	err := c.Bind(postForm)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}
	//获取配置信息
	config := configs.GetConfig()

	// 获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	//初始化修改数据
	var nowTime = pkg.GetNowTimeUnixMilli()
	post := &models.Post{}
	post.DiaryID = postForm.DiaryID
	post.Title = postForm.Title
	post.Body = postForm.Body
	// 摘要长度从配置中获取
	post.Summary = pkg.TruncateString(postForm.Body, config.App.SummaryMax,config.App.Suffix) 

	// 保存日记
	action, _ := c.GetPostForm("action")
	if action == "post.create" {
		// 其他设置默认值
		post.IsTop = 0
		post.InfoType = 1 // 默认为文本
		post.AddUTCTime = nowTime
		post.AddUserID = userData.UserId
		post.AddUserName = userData.UserName
		post.AddUserIP = c.ClientIP()
		post.UpdateUTCTime = nowTime
		post.Status = 1
		post.ViewNum = 0
		post.LikeNum = 0
		post.CommentNum = 0

		postID, err := services.AddPost(post)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		post.PostID = postID
	} else if action == "post.edit" {
		// 获取 post_id
		postIdStr, _ := c.GetPostForm("post_id")
		postId, err := strconv.ParseInt(postIdStr, 10, 64)
		if err != nil || postId <= 0 {
			errors.ResponseErr(c, errors.ParseInputParamErr)
			return
		}
		post.PostID = postId
		// 获取原始日记信息进行权限校验
		originalPost, err := services.GetPostByID(post.PostID)
		if err != nil || originalPost == nil {
			errors.ResponseErr(c, errors.ParseInputParamErr)
			return
		}

		// 检查是否是当前用户的日记
		if originalPost.AddUserID != userData.UserId {
			errors.ResponseErr(c, errors.InputParamErr)
			return
		}

		post.UpdateUTCTime = nowTime
		// 保留原有的一些字段值
		post.DiaryID = originalPost.DiaryID
		post.ViewNum = originalPost.ViewNum
		post.LikeNum = originalPost.LikeNum
		post.CommentNum = originalPost.CommentNum
		post.IsTop = originalPost.IsTop
		post.InfoType = originalPost.InfoType

		// 更新日记
		err = services.UpdatePost(post)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
	}

	// 创建返回对象
	h := pkg.ResponseH(c)
	h["Data"] = post
	c.JSON(http.StatusOK, h)
}
