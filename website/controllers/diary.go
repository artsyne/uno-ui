package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"

	"github.com/gin-gonic/gin"
)

// RenderDiaryIndex 渲染日记本首页
func RenderDiaryIndex(c *gin.Context) {
	h := pkg.DefaultH(c)

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

	// 日记总数
	diariesNum := 0
	pageSize := configs.GetConfig().App.PageSize

	// 分页显示日记列表
	h["Posts"], diariesNum, err = models.ListPostsByUserID(diary.AddUserID, diaryID, pageIndex, pageSize)
	if err != nil {
		hasError = true
	}

	// 初始化分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/d/%d/%%d", diaryID),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 分页导航链接
	h["Diary"] = diary
	h["PageLinks"] = pageLinks
	h["DiaryID"] = diaryID
	h["hasError"] = hasError
	h["errMsg"] = errMsg
	h["DiariesNum"] = diariesNum
	h["PageSize"] = pageSize

	c.HTML(http.StatusOK, "diary/index", h)
}

// 增加日记本表单渲染
func RenderDiaryForm(c *gin.Context) {
	h := pkg.DefaultH(c)
	// 获取action 和 diary_id
	action := c.Query("action")
	// 如果是编辑日记本，则需要获取 diary_id
	if action == "diary.edit" {
		diaryID, err := strconv.ParseInt(c.Query("diary_id"), 10, 64)
		if err != nil || diaryID <= 0 {
			errors.ResponseErr(c, errors.ParseInputParamErr)
			return
		}
		//获取用户登录信息
		userSession, _ := c.Get("UserData")
		userData := userSession.(models.Login)

		// 获取日记本信息
		diary, err := models.GetDiaryByDiaryID(diaryID)
		if err != nil || diary == nil {
			errors.ResponseErr(c, errors.ParseInputParamErr)
			return
		}
		// 如果不是当前用户的日记本，则返回错误
		if diary.AddUserID != userData.UserId {
			errors.ResponseErr(c, errors.InputParamErr)
			return
		}
		// 返回日记本信息到模板
		h["Diary"] = diary

	}
	c.HTML(http.StatusOK, "form/diary", h)
}

// RenderDiary : 渲染日记详细信息
func RenderDiary(c *gin.Context) {
	h := pkg.DefaultH(c)
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

	var nowTime = pkg.GetNowTimeUnixMilli()
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
	action, _ := c.GetPostForm("action")
	if action == "diary.add" {
		diaryID, err := services.AddDiary(diary)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		diary.DiaryID = diaryID
	} else if action == "diary.edit" {
		// 获取原始日记本信息进行权限校验
		originalDiary, err := models.GetDiaryByDiaryID(diary.DiaryID)
		if err != nil || originalDiary == nil {
			errors.ResponseErr(c, errors.ParseInputParamErr)
			return
		}

		// 检查是否是当前用户的日记本
		if originalDiary.AddUserID != userData.UserId {
			errors.ResponseErr(c, errors.InputParamErr)
			return
		}

		// 保留原有的一些字段值
		diary.AddUTCTime = originalDiary.AddUTCTime
		diary.AddUserID = originalDiary.AddUserID
		diary.AddUserName = originalDiary.AddUserName
		diary.AccessURL = originalDiary.AccessURL
		diary.CoverURL = originalDiary.CoverURL
		diary.CoverMidURL = originalDiary.CoverMidURL
		diary.LikeNum = originalDiary.LikeNum
		diary.MemberNum = originalDiary.MemberNum
		diary.ViewNum = originalDiary.ViewNum
		diary.InfoNum = originalDiary.InfoNum
		diary.Status = originalDiary.Status
		diary.AccessPwd = originalDiary.AccessPwd

		// 更新日记本
		err = services.UpdateDiary(diary)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
	}

	//创建返回对象
	h := pkg.ResponseH(c)
	h["Data"] = diary
	c.JSON(http.StatusOK, h)
}
