package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/helpers"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/internal/services"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// DiaryController 日记本控制器
// 负责处理日记本相关的所有HTTP请求，包括日记本列表展示、表单渲染、创建和编辑等操作
type DiaryController struct{}

// NewDiaryController 创建新的日记本控制器实例
// 返回值: *DiaryController 日记本控制器指针
func NewDiaryController() *DiaryController {
	return &DiaryController{}
}

// RenderDiaryIndex 渲染日记本首页
// 显示指定日记本的详细信息和其包含的日记列表，支持分页显示
// 同时判断当前用户是否为日记本所有者，以决定是否显示管理功能
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// URL参数:
//   - diaryid: 日记本ID
//   - pageindex: 页码（可选，默认为1）
func (dc *DiaryController) RenderDiaryIndex(c *gin.Context) {
	h := pkg.DefaultH(c)
	h["Title"] = "日记本"
	h["Active"] = "users"

	// 页面加载错误标志
	hasError := false
	errMsg := ""

	// 解析URL中的日记本ID
	diaryID, err := strconv.ParseInt(strings.TrimSpace(c.Param("diaryid")), 10, 64)
	if err != nil {
		hasError = true
		h["SpaceID"] = diaryID
		h["hasError"] = hasError
		h["errMsg"] = errMsg
		c.HTML(http.StatusOK, "diary/index", h)
		return
	}

	// 解析页码参数，默认为第1页
	pageIndex, err := strconv.Atoi(strings.TrimSpace(c.Param("pageindex")))
	if err != nil || pageIndex <= 0 {
		pageIndex = 1
	}

	// 获取日记本信息
	diary, _ := models.GetDiaryByDiaryID(diaryID)
	if diary == nil {
		hasError = true
		h["DiaryID"] = diaryID
		h["hasError"] = hasError
		h["errMsg"] = errMsg
		c.HTML(http.StatusOK, "diary/index", h)
		return
	}

	// 获取分页配置
	pageSize := configs.GetConfig().App.PageSize

	// 获取分页日记列表
	posts, diariesNum, err := models.ListPostsByUserID(diary.AddUserID, diaryID, pageIndex, pageSize)
	if err != nil {
		hasError = true
	}
	h["Posts"] = posts

	// 生成分页导航
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/d/%d/%%d", diaryID),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 判断当前用户是否为日记本所有者
	isAuth, loginUser := helpers.Auth.IsAuthenticated(c)
	isOwner := isAuth && loginUser.UserId == diary.AddUserID
	h["IsOwner"] = isOwner

	// 设置模板数据
	h["Diary"] = diary
	h["PageLinks"] = pageLinks
	h["hasError"] = hasError
	h["errMsg"] = errMsg
	h["DiariesNum"] = diariesNum
	h["PageSize"] = pageSize

	c.HTML(http.StatusOK, "diary/index", h)
}

// RenderDiaryForm 渲染日记本创建/编辑表单页面
// 根据action参数决定是渲染创建表单还是编辑表单
// 对于编辑操作，会验证用户权限并加载现有日记本数据
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 支持的URL参数:
//   - action: "diary.add" 或 "diary.edit"
//   - diary_id: 当action为edit时需要提供的日记本ID
func (dc *DiaryController) RenderDiaryForm(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := helpers.Auth.MustGetCurrentUser(c)

	h := pkg.DefaultH(c)

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	if action == "diary.edit" {
		// 编辑模式：获取并验证日记本数据
		diary, err := dc.getAndValidateDiary(c, loginUser.UserId)
		if err != nil {
			errors.ResponseErr(c, err)
			return
		}
		h["Diary"] = diary
	}
	// 创建模式：直接渲染空表单

	c.HTML(http.StatusOK, "form/diary", h)
}

// RenderDiary 渲染日记详细信息页面
// 显示单个日记的详细内容
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (dc *DiaryController) RenderDiary(c *gin.Context) {
	h := pkg.DefaultH(c)
	h["Title"] = "我的日记本"
	h["Active"] = "users"
	c.HTML(http.StatusOK, "diary/diary", h)
}

// HandleDiary 处理日记本的创建和编辑请求
// 根据表单中的action字段分发到对应的处理方法
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
//
// 支持的表单action值:
//   - "diary.add": 创建新日记本
//   - "diary.edit": 编辑现有日记本
func (dc *DiaryController) HandleDiary(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := helpers.Auth.MustGetCurrentUser(c)

	// 解析表单数据到Diary结构体
	diary := &models.Diary{}
	if err := c.Bind(diary); err != nil {
		errors.ResponseErr(c, errors.ParseInputParamErr)
		return
	}

	// 根据操作类型分发处理
	action, _ := c.GetPostForm("action")
	switch action {
	case "diary.add":
		dc.handleCreateDiary(c, diary, loginUser)
	case "diary.edit":
		dc.handleEditDiary(c, diary, loginUser)
	default:
		errors.ResponseErr(c, errors.ParseInputParamErr)
	}
}

// handleCreateDiary 处理创建新日记本的业务逻辑
// 设置日记本的所有必要字段，调用服务层创建日记本，并返回创建结果
//
// 参数:
//   - c: Gin上下文
//   - diary: 解析后的日记本数据
//   - user: 当前登录用户信息
func (dc *DiaryController) handleCreateDiary(c *gin.Context, diary *models.Diary, user *models.Login) {
	nowTime := pkg.GetNowTimeUnixMilli()

	// 设置新日记本的默认值和用户信息
	diary.AccessPwd = ""
	diary.AccessURL = ""
	diary.AddUTCTime = nowTime
	diary.AddUserID = user.UserId
	diary.AddUserName = user.UserName
	diary.CoverMidURL = ""
	diary.CoverURL = ""
	diary.LikeNum = 0
	diary.MemberNum = 0
	diary.Status = 1 // 默认状态为正常
	diary.ViewNum = 0
	diary.InfoNum = 0 // 日记数量初始为0
	diary.LastActionUTCTime = nowTime
	diary.LastActionUserID = user.UserId
	diary.LastActionUserName = user.UserName

	// 调用服务层创建日记本
	diaryID, err := services.AddDiary(diary)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}
	diary.DiaryID = diaryID

	// 返回成功响应
	dc.respondSuccess(c, diary)
}

// handleEditDiary 处理编辑现有日记本的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记本
//
// 参数:
//   - c: Gin上下文
//   - diary: 解析后的日记本数据
//   - user: 当前登录用户信息
func (dc *DiaryController) handleEditDiary(c *gin.Context, diary *models.Diary, user *models.Login) {
	// 验证日记本所有权
	originalDiary, err := dc.validateDiaryOwnership(diary.DiaryID, user.UserId)
	if err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 保留原有字段值，只更新可编辑字段
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

	// 更新最后操作时间和操作用户
	diary.LastActionUTCTime = pkg.GetNowTimeUnixMilli()
	diary.LastActionUserID = user.UserId
	diary.LastActionUserName = user.UserName

	// 调用服务层更新日记本
	if err := services.UpdateDiary(diary); err != nil {
		errors.ResponseErr(c, err)
		return
	}

	// 返回成功响应
	dc.respondSuccess(c, diary)
}

// getAndValidateDiary 获取并验证日记本权限（用于表单渲染）
// 从URL参数中获取diary_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *models.Diary: 日记本对象，如果验证失败返回nil
//   - error: 错误信息
func (dc *DiaryController) getAndValidateDiary(c *gin.Context, userID int64) (*models.Diary, error) {
	diaryID, err := strconv.ParseInt(c.Query("diary_id"), 10, 64)
	if err != nil || diaryID <= 0 {
		return nil, errors.ParseInputParamErr
	}

	return dc.validateDiaryOwnership(diaryID, userID)
}

// validateDiaryOwnership 验证日记本所有权
// 检查指定日记本是否存在，以及当前用户是否有权限操作该日记本
//
// 参数:
//   - diaryID: 日记本ID
//   - userID: 当前用户ID
//
// 返回值:
//   - *models.Diary: 日记本对象
//   - error: 验证失败时的错误信息
func (dc *DiaryController) validateDiaryOwnership(diaryID, userID int64) (*models.Diary, error) {
	diary, err := models.GetDiaryByDiaryID(diaryID)
	if err != nil || diary == nil {
		return nil, errors.ParseInputParamErr
	}

	// 检查是否是当前用户的日记本
	if diary.AddUserID != userID {
		return nil, errors.InputParamErr
	}

	return diary, nil
}

// respondSuccess 统一的成功响应格式
// 封装标准的JSON响应格式，包含操作成功的日记本数据
//
// 参数:
//   - c: Gin上下文
//   - diary: 操作成功的日记本对象
func (dc *DiaryController) respondSuccess(c *gin.Context, diary *models.Diary) {
	h := pkg.ResponseH(c)
	h["Data"] = diary
	c.JSON(http.StatusOK, h)
}
