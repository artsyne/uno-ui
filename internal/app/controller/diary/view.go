package diary

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// ShowDiary 显示日记详细信息页面
// 显示单个日记的详细内容
func (ctrl *DiaryController) ShowDiary(c *gin.Context) {
	data := ctrl.renderService.CreatePageData(c, "diary.detail.title")
	data.Data["Active"] = "users"
	ctrl.renderService.HTML(c, "diary/diary", data)
}

// ShowDiaryForm 显示日记本创建/编辑表单页面
// 根据URL参数中的action和diary_id决定是创建新日记本还是编辑现有日记本
//
// URL参数:
//   - action: 操作类型，"diary.add"表示创建，"diary.edit"表示编辑
//   - diary_id: 当action为edit时需要提供的日记本ID
func (ctrl *DiaryController) ShowDiaryForm(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := ctrl.auth.MustGetCurrentUser(c)

	data := ctrl.renderService.CreatePageData(c, "diary.form.title")

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	if action == "diary.edit" {
		// 编辑模式：获取并验证日记本数据
		diaryData, err := ctrl.getAndValidateDiary(c, loginUser.UserId)
		if err != nil {
			ctrl.renderService.JSONError(c, err)
			return
		}
		data.Data["Diary"] = diaryData
	}
	// 创建模式：直接渲染空表单

	ctrl.renderService.HTML(c, "form/diary", data)
}

// ShowUserProfile 显示用户个人主页
// 显示指定用户的日记本列表，支持分页显示
//
// URL格式:
//   - /u/{uid} - 显示第一页
//   - /u/{uid}/{pageindex} - 显示指定页
func (ctrl *DiaryController) ShowUserProfile(c *gin.Context) {
	// 获取url中的uid
	uid, err := ctrl.parseDiaryID(c.Param("uid")) // 复用parseDiaryID方法
	if err != nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	// 创建页面数据
	data := ctrl.renderService.CreatePageData(c, "user.profile.title")

	// 判断是否为页面所有者
	isOwner := false
	if data.IsAuthenticated() {
		isOwner = (data.GetAuthUser().UserId == uid)
	}

	// 获取URL中的pageIndex
	pageIndexStr := c.Param("pageindex")
	pageIndex := 1
	if pageIndexStr != "" {
		if idx, err := strconv.Atoi(strings.TrimSpace(pageIndexStr)); err == nil && idx > 0 {
			pageIndex = idx
		}
	}

	// 获取分页配置
	pageSize := ctrl.config.App.PageSize

	// 获取用户的日记本列表
	diaries, diariesNum, _ := ctrl.diaryService.ListDiariesByUser(uid, pageIndex, pageSize)

	// 初始化分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/u/%d/%%d", uid),
		PagePrev:  "<<",
		PageNext:  ">>",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 根据uid 获取用户信息
	pageUser, err := ctrl.userService.GetUserByUserId(uid)
	if err != nil || pageUser == nil {
		ctrl.renderService.Redirect(c, "/")
		return
	}

	// 组织业务数据
	data.Data["IsOwner"] = isOwner
	data.Data["User"] = pageUser              //被访问用户
	data.Data["TimeZone"] = pageUser.TimeZone // 被访问用户时区
	data.Data["Diaries"] = diaries
	data.Data["DiariesNum"] = diariesNum
	data.Data["PageSize"] = pageSize
	data.Data["PageIndex"] = pageIndex
	data.Data["PageLinks"] = pageLinks

	ctrl.renderService.HTML(c, "member/index", data)
}