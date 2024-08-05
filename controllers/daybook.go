package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/errors"
	"github.com/biuren/hebe/models"
	"github.com/biuren/hebe/utils"
	"github.com/gin-gonic/gin"
)

// RenderDaybookIndex 渲染日记本首页
func RenderDaybookIndex(c *gin.Context) {
	h := utils.DefaultH(c)

	//页面加载是否有异常
	hasError := false
	errMsg := ""

	h["Title"] = "日记本"
	h["Active"] = "users"

	//获取url中spaceId
	spaceId, err := strconv.ParseInt(strings.TrimSpace(c.Param("spaceid")), 10, 64)
	if err == nil {
		space, _ := models.GetSpaceBySpaceID(spaceId)
		if space != nil {
			h["Space"] = space
			//日记总数
			diariesNum := 0
			pageIndex := 1
			pageSize := 5
			//分页显示日记列表
			h["Diaries"], diariesNum, err = models.ListInfoByUserID(space.AddUserID, spaceId, pageIndex, pageSize)
			if err != nil {
				hasError = true
			}
			//初始化分页信息

			listOptions := utils.ListOptions{
				PageIndex: pageIndex,
				PageSize:  pageSize,
				Total:     diariesNum,
				URL:       fmt.Sprintf("/daybook/%d/%%d",spaceId),
				PagePrev:  "上页",
				PageNext:  "下页",
			}
			pageLinks := utils.GeneratePagination(listOptions)
			//分页导航链接
			h["PageLinks"] = pageLinks

		} else {
			hasError = true
		}
	} else {
		hasError = true
	}

	h["SpaceID"] = spaceId
	h["hasError"] = hasError
	h["errMsg"] = errMsg

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
	space := &models.Space{}
	err := c.Bind(space)

	// 校验参数，标题必须填写，设置类型必须为1或者2，主题id必须在1-16之间
	if err != nil || len(space.Title) < 1 || (space.ViewSetting != 1 &&
		space.ViewSetting != 2 || space.ThemeID < 1 || space.ThemeID > 16) {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	//获取用户登录信息
	userSession, _ := c.Get("UserData")
	userData := userSession.(models.Login)

	var nowTime = utils.GetNowTimeUnixMilli()
	//其他设置默认值
	space.AccessPwd = ""
	space.AccessURL = ""
	space.AddUTCTime = nowTime
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
	space.LastActionUTCTime = nowTime
	space.LastActionUserID = userData.UserId
	space.LastActionUserName = userData.UserName

	//保存日记本
	if action, _ := c.GetPostForm("action"); action == "My.insertDaybook" {
		spaceID, err := space.Insert() // 插入对象
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		space.SpaceID = spaceID
	}

	//创建返回对象
	h := utils.ResponseH(c)
	h["Data"] = space
	c.JSON(http.StatusOK, h)
}
