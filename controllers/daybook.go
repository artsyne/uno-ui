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

// RenderDaybookIndex 渲染日记本首页
func RenderDaybookIndex(c *gin.Context) {
	h := utils.DefaultH(c)

	// 页面加载是否有异常
	hasError := false
	errMsg := ""

	h["Title"] = "日记本"
	h["Active"] = "users"

	// 获取 URL 中的 daybookID
	daybookID, err := strconv.ParseInt(strings.TrimSpace(c.Param("spaceid")), 10, 64)
	if err != nil {
		hasError = true
		h["SpaceID"] = daybookID
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

	daybook, _ := models.GetDaybookByDaybookID(daybookID)
	if daybook == nil {
		hasError = true
		h["SpaceID"] = daybookID
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
	h["Diaries"], diariesNum, err = models.ListInfoByUserID(daybook.AddUserID, daybookID, pageIndex, pageSize)
	if err != nil {
		hasError = true
	}

	// 初始化分页信息
	listOptions := utils.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/daybook/%d/%%d", daybookID),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := utils.GeneratePagination(listOptions)

	// 分页导航链接
	h["Space"] = daybook
	h["PageLinks"] = pageLinks
	h["SpaceID"] = daybookID
	h["hasError"] = hasError
	h["errMsg"] = errMsg
	h["DiariesNum"] = diariesNum
	h["PageSize"] = pageSize

	c.HTML(http.StatusOK, "daybook/index", h)
}

// 增加日记本表单渲染
func RenderDaybookForm(c *gin.Context) {
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
func HandleSpace(c *gin.Context) {
	// 序列化用户输入对象
	daybook := &models.Daybook{}
	err := c.Bind(daybook)
	if err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	var nowTime = utils.GetNowTimeUnixMilli()
	//其他设置默认值
	daybook.AccessPwd = ""
	daybook.AccessURL = ""
	daybook.AddUTCTime = nowTime
	daybook.AddUserID = userData.UserId
	daybook.AddUserName = userData.UserName
	daybook.CoverMidURL = ""
	daybook.CoverURL = ""
	daybook.LikeNum = 0
	daybook.MemberNum = 0
	daybook.Status = 1
	daybook.ViewNum = 0
	daybook.InfoNum = 0
	daybook.LastActionUTCTime = nowTime
	daybook.LastActionUserID = userData.UserId
	daybook.LastActionUserName = userData.UserName

	//保存日记本
	if action, _ := c.GetPostForm("action"); action == "My.insertDaybook" {
		daybookID, err := services.CreateNewDaybook(*daybook)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		daybook.SpaceID = daybookID
	}

	//创建返回对象
	h := utils.ResponseH(c)
	h["Data"] = daybook
	c.JSON(http.StatusOK, h)
}
