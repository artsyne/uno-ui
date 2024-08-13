package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/services"
	"github.com/biuren/hebe/settings"
	"github.com/biuren/hebe/utils"
	"github.com/gin-gonic/gin"
)

// RenderDiaryIndex 渲染日记本首页
func RenderDiaryIndex(c *gin.Context) {
	h := utils.DefaultH(c)

	// 页面加载是否有异常
	hasError := false
	errMsg := ""

	h["Title"] = "日记本"
	h["Active"] = "users"

	// 获取 URL 中的 diaryID
	diaryID, err := strconv.ParseInt(strings.TrimSpace(c.Param("spaceid")), 10, 64)
	if err != nil {
		hasError = true
		h["SpaceID"] = diaryID
		h["hasError"] = hasError
		h["errMsg"] = errMsg
		c.HTML(http.StatusOK, "daybook/index", h)
		return
	}

	// 获取 URL 中的 pageIndex
	pageIndex, err := strconv.Atoi(strings.TrimSpace(c.Param("pageindex")))
	if err != nil || pageIndex <= 0 {
		pageIndex = 1
	}

	diary, _ := models.GetDiaryByDiaryID(diaryID)
	if diary == nil {
		hasError = true
		h["DiaryID"] = diaryID
		h["hasError"] = hasError
		h["errMsg"] = errMsg
		c.HTML(http.StatusOK, "daybook/index", h)
		return
	}

	//读取配置
	config := settings.GetConfig()

	// 日记总数
	diariesNum := 0
	pageSize := config.PageSize

	// 分页显示日记列表
	h["Diaries"], diariesNum, err = models.ListInfoByUserID(diary.AddUserID, diaryID, pageIndex, pageSize)
	if err != nil {
		hasError = true
	}

	// 初始化分页信息
	listOptions := utils.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/daybook/%d/%%d", diaryID),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := utils.GeneratePagination(listOptions)

	// 分页导航链接
	h["Diary"] = diary
	h["PageLinks"] = pageLinks
	h["DiaryID"] = diaryID
	h["hasError"] = hasError
	h["errMsg"] = errMsg
	h["DiariesNum"] = diariesNum
	h["PageSize"] = pageSize

	c.HTML(http.StatusOK, "daybook/index", h)
}

// 增加日记本表单渲染
func RenderDiaryForm(c *gin.Context) {
	h := utils.DefaultH(c)

	c.HTML(http.StatusOK, "form/daybook", h)
}

// RenderDiary : 渲染日记详细信息
func RenderDiary(c *gin.Context) {
	h := utils.DefaultH(c)
	h["Title"] = "我的日记本"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "daybook/diary", h)
}

// 保存日记本信息
func HandleDiary(c *gin.Context) {
	// 序列化用户输入对象
	diary := &models.Diary{}
	err := c.Bind(diary)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	var nowTime = utils.GetNowTimeUnixMilli()
	//其他设置默认值
	diary.AccessPwd = ""
	diary.AccessURL = ""
	diary.AddUTCTime = nowTime
	diary.AddUserID = userData.UserId
	diary.AddUserName = userData.UserName
	diary.CoverMidURL = ""
	diary.CoverURL = ""
	diary.LikeNum = 0
	diary.MemberNum = 0
	diary.Status = 1
	diary.ViewNum = 0
	diary.InfoNum = 0
	diary.LastActionUTCTime = nowTime
	diary.LastActionUserID = userData.UserId
	diary.LastActionUserName = userData.UserName

	//保存日记本
	if action, _ := c.GetPostForm("action"); action == "My.insertDaybook" {
		diaryID, err := services.CreateNewDiary(*diary)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		diary.DiaryID = diaryID
	}

	//创建返回对象
	h := utils.ResponseH(c)
	h["Data"] = diary
	c.JSON(http.StatusOK, h)
}
