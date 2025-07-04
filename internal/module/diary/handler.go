package diary

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/core/config"
	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/biuren/hebe/internal/module/user"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

// HandlerInterface 日记本模块HTTP处理器接口
// 定义所有与日记本相关的HTTP处理操作
type HandlerInterface interface {
	// RegisterRoutes 注册所有日记本相关路由
	RegisterRoutes(router *gin.Engine)

	// RenderDiaryForm 渲染日记本创建/编辑表单页面
	RenderDiaryForm(c *gin.Context)

	// RenderDiary 渲染日记详细信息页面
	RenderDiary(c *gin.Context)

	// HandleDiary 处理日记本的创建和编辑请求
	HandleDiary(c *gin.Context)

	// RenderUserProfile 渲染用户个人主页
	RenderUserProfile(c *gin.Context)

	// ShowDiaryForm 显示创建/编辑日记表单页面
	ShowDiaryForm(c *gin.Context)

	// ShowDiaryDetail 显示日记详情页面
	ShowDiaryDetail(c *gin.Context)
}

// handler 日记本模块HTTP处理器实现
// 负责处理所有与日记本相关的HTTP请求，包括渲染列表、表单、创建日记本、编辑日记本等操作
type handler struct {
	service       ServiceInterface
	auth          user.AuthHelperInterface
	config        *config.Config
	renderService render.RenderServiceInterface
}

// NewHandlerWithServices 通过DI容器创建日记本处理器实例
func NewHandlerWithServices(renderService render.RenderServiceInterface, _ render.RenderServiceInterface) HandlerInterface {
	return &handler{
		service:       NewService(),
		auth:          user.NewAuthHelper(),
		config:        config.GetConfig(),
		renderService: renderService,
	}
}

// RegisterRoutes 注册所有日记本相关路由
// 将所有diary相关路由聚合在一个地方，实现高内聚
func (h *handler) RegisterRoutes(router *gin.Engine) {
	// 日记相关路由组（无需鉴权）
	diary := router.Group("/diary")
	{
		diary.GET("/", h.RenderDiary)
	}

	// 用户个人主页路由（公开访问，显示用户的日记本列表）
	router.GET("/u/:uid", h.RenderUserProfile)
	router.GET("/u/:uid/:pageindex", h.RenderUserProfile)

	// 需要鉴权的日记本表单路由
	authForm := router.Group("/form")
	authForm.Use(user.AuthRequired())
	{
		authForm.GET("/diary", h.RenderDiaryForm)
	}

	// 需要鉴权的日记本数据操作路由（保持原路径）
	authAPI := router.Group("/")
	authAPI.Use(user.AuthRequired())
	{
		authAPI.POST("/diary", h.HandleDiary)
		authAPI.POST("/daybook", h.HandleDiary)
	}
}

// parseDiaryID 解析并验证日记本ID
// 从字符串中解析出有效的日记本ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的ID
//   - error: 解析失败时的错误信息
func (h *handler) parseDiaryID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParamsErr()
	}
	return id, nil
}

// RenderDiaryForm 渲染日记本创建/编辑表单页面
// 根据URL参数中的action和diary_id决定是创建新日记本还是编辑现有日记本
//
// URL参数:
//   - action: 操作类型，"diary.add"表示创建，"diary.edit"表示编辑
//   - diary_id: 当action为edit时需要提供的日记本ID
func (h *handler) RenderDiaryForm(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := h.auth.MustGetCurrentUser(c)

	data := h.renderService.PageData(c, "diary.form.title")

	// 根据action参数处理不同的表单渲染逻辑
	action := c.Query("action")
	if action == "diary.edit" {
		// 编辑模式：获取并验证日记本数据
		diary, err := h.getAndValidateDiary(c, loginUser.UserId)
		if err != nil {
			h.renderService.JSONError(c, err)
			return
		}
		data.Data["Diary"] = diary
	}
	// 创建模式：直接渲染空表单

	h.renderService.HTML(c, "form/diary", data)
}

// RenderDiary 渲染日记详细信息页面
// 显示单个日记的详细内容
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (h *handler) RenderDiary(c *gin.Context) {
	data := h.renderService.PageData(c, "diary.detail.title")
	data.Data["Active"] = "users"
	h.renderService.HTML(c, "diary/diary", data)
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
func (h *handler) HandleDiary(c *gin.Context) {
	// 获取当前登录用户信息
	loginUser := h.auth.MustGetCurrentUser(c)

	// 解析表单数据到Diary结构体
	diary := &Diary{}
	if err := c.Bind(diary); err != nil {
		h.renderService.JSONError(c, errors.InvalidParamsErr())
		return
	}

	// 根据操作类型分发处理
	action, _ := c.GetPostForm("action")
	switch action {
	case "diary.add":
		h.handleCreateDiary(c, diary, loginUser)
	case "diary.edit":
		h.handleEditDiary(c, diary, loginUser)
	default:
		h.renderService.JSONError(c, errors.InvalidParamsErr())
	}
}

// handleCreateDiary 处理创建新日记本的业务逻辑
// 设置日记本的所有必要字段，调用服务层创建日记本，并返回创建结果
//
// 参数:
//   - c: Gin上下文
//   - formData: 解析后的表单数据
//   - user: 当前登录用户信息
func (h *handler) handleCreateDiary(c *gin.Context, formData *Diary, user *user.AuthUser) {
	nowTime := pkg.GetNowTimeUnixMilli()

	// 使用结构体字面量创建完整的日记本对象
	diary := &Diary{
		// 用户输入的字段
		Title:         formData.Title,
		Remark:        formData.Remark,
		ViewSetting:   formData.ViewSetting,
		SocialSetting: formData.SocialSetting,
		ThemeID:       formData.ThemeID,

		// 系统设置的字段
		DiaryID:            0, // 将由数据库生成
		AccessPwd:          "",
		AccessURL:          "",
		CoverURL:           "",
		CoverMidURL:        "",
		LikeNum:            0,
		MemberNum:          0,
		ViewNum:            0,
		InfoNum:            0, // 日记数量初始为0
		Status:             1, // 默认状态为正常
		AddUTCTime:         nowTime,
		AddUserID:          user.UserId,
		AddUserName:        user.UserName,
		LastActionUTCTime:  nowTime,
		LastActionUserID:   user.UserId,
		LastActionUserName: user.UserName,
	}

	// 调用服务层创建日记本
	diaryID, err := h.service.AddDiary(diary)
	if err != nil {
		h.renderService.JSONError(c, err)
		return
	}
	diary.DiaryID = diaryID

	// 返回成功响应
	h.respondSuccess(c, diary)
}

// handleEditDiary 处理编辑现有日记本的业务逻辑
// 验证用户权限，保留原有字段值，更新可编辑字段，调用服务层更新日记本
//
// 参数:
//   - c: Gin上下文
//   - diary: 解析后的日记本数据
//   - user: 当前登录用户信息
func (h *handler) handleEditDiary(c *gin.Context, diary *Diary, user *user.AuthUser) {
	// 验证日记本所有权
	originalDiary, err := h.validateDiaryOwnership(diary.DiaryID, user.UserId)
	if err != nil {
		h.renderService.JSONError(c, err)
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
	if err := h.service.UpdateDiary(diary); err != nil {
		h.renderService.JSONError(c, err)
		return
	}

	// 返回成功响应
	h.respondSuccess(c, diary)
}

// getAndValidateDiary 获取并验证日记本权限（用于表单渲染）
// 从URL参数中获取diary_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *Diary: 日记本对象，如果验证失败返回nil
//   - error: 错误信息
func (h *handler) getAndValidateDiary(c *gin.Context, userID int64) (*Diary, error) {
	diaryID, err := h.parseDiaryID(c.Query("diary_id"))
	if err != nil {
		return nil, err
	}

	return h.validateDiaryOwnership(diaryID, userID)
}

// validateDiaryOwnership 验证日记本所有权
// 检查指定用户是否拥有指定日记本的操作权限
//
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回值:
//   - *Diary: 日记本对象，如果验证失败返回nil
//   - error: 错误信息
func (h *handler) validateDiaryOwnership(diaryID, userID int64) (*Diary, error) {
	diary, err := h.service.GetDiaryByID(diaryID)
	if err != nil {
		return nil, errors.DiaryNotFoundErr()
	}

	// 验证所有权
	if diary.AddUserID != userID {
		return nil, errors.PermissionErr()
	}

	return diary, nil
}

// respondSuccess 返回成功响应
// 统一的成功响应格式，便于后续维护和修改
//
// 参数:
//   - c: Gin上下文
//   - diary: 日记本数据
func (h *handler) respondSuccess(c *gin.Context, diary *Diary) {
	h.renderService.JSONSuccess(c, diary)
}

// renderErrorPage 渲染错误页面
// 用于显示友好的错误信息页面
//
// 参数:
//   - c: Gin上下文
//   - title: 错误标题
//   - message: 错误消息
func (h *handler) renderErrorPage(c *gin.Context, title, message string) {
	h.renderService.HTMLError(c, title, message)
}

// RenderUserProfile 渲染用户个人主页
// 显示指定用户的日记本列表，支持分页显示
//
// URL格式:
//   - /u/{uid} - 显示第一页
//   - /u/{uid}/{pageindex} - 显示指定页
//
// 参数:
//   - c: Gin上下文，包含请求信息和响应操作
func (h *handler) RenderUserProfile(c *gin.Context) {
	// 获取url中的uid
	uid, err := h.parseDiaryID(c.Param("uid")) // 复用parseDiaryID方法
	if err != nil {
		h.renderService.Redirect(c, "/")
		return
	}

	// 创建页面数据
	data := h.renderService.PageData(c, "user.profile")

	// 判断是否为页面所有者
	isOwner := false
	if data.IsAuthenticated() {
		isOwner = (data.GetAuthUser().GetUserID() == uid)
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
	pageSize := h.config.App.PageSize

	// 获取用户的日记本列表
	diaries, diariesNum, _ := h.service.ListDiariesByUser(uid, pageIndex, pageSize)

	// 初始化分页信息
	listOptions := pkg.ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     diariesNum,
		URL:       fmt.Sprintf("/u/%d/%%d", uid),
		PagePrev:  "上页",
		PageNext:  "下页",
	}
	pageLinks := pkg.GeneratePagination(listOptions)

	// 根据uid 获取用户信息
	pageUser, err := user.NewService().GetUserByUserId(uid)
	if err != nil || pageUser == nil {
		h.renderService.Redirect(c, "/")
		return
	}

	// 组织业务数据
	data.Data["IsOwner"] = isOwner
	data.Data["User"] = pageUser //被访问用户
	data.Data["TimeZone"] = pageUser.TimeZone // 被访问用户时区
	data.Data["Diaries"] = diaries
	data.Data["DiariesNum"] = diariesNum
	data.Data["PageSize"] = pageSize
	data.Data["PageLinks"] = pageLinks

	h.renderService.HTML(c, "member/index", data)
}

// ShowDiaryForm 显示创建/编辑日记表单页面
func (h *handler) ShowDiaryForm(c *gin.Context) {
	// 获取当前登录用户（中间件保证已登录）
	authUser := h.auth.MustGetCurrentUser(c)

	// 根据Form Action决定是创建还是编辑
	formAction := c.DefaultQuery("form_action", "diary.create")
	var diary *Diary

	if formAction == "diary.edit" {
		// 编辑模式：获取日记信息
		diaryIdStr := c.Query("diary_id")
		if diaryIdStr == "" {
			h.renderService.JSONError(c, errors.ValidationErr().WithDetails("missing diary_id"))
			return
		}

		// 解析日记ID
		diaryId, err := strconv.ParseInt(diaryIdStr, 10, 64)
		if err != nil {
			h.renderService.Redirect(c, "/")
			return
		}

		// 获取日记详情
		diary, err = h.service.GetDiaryByID(diaryId)
		if err != nil || diary == nil {
			h.renderService.JSONError(c, errors.InternalErr().WithDetails("diary not found"))
			return
		}

		// 检查权限
		if diary.AddUserID != authUser.UserId {
			h.renderService.JSONError(c, errors.PermissionErr())
			return
		}
	}

	// 构建表单数据
	data := h.renderService.PageData(c, "diary.form.title")
	data.Data["FormAction"] = formAction
	data.Data["Diary"] = diary
	data.Data["SupportLangs"] = config.GetSupportLangs()
	data.Data["TimeZones"] = config.GetTimeZone()
	// 注意：用户信息现在通过 data.Context.AuthUser 在模板中访问，无需手动设置

	h.renderService.HTML(c, "form/diary", data)
}

// ShowDiaryDetail 显示日记详情页面
func (h *handler) ShowDiaryDetail(c *gin.Context) {
	diaryIdStr := c.Param("id")
	if diaryIdStr == "" {
		h.renderService.Redirect(c, "/")
		return
	}

	// 解析日记ID
	diaryId, err := strconv.ParseInt(diaryIdStr, 10, 64)
	if err != nil {
		h.renderService.Redirect(c, "/")
		return
	}

	// 获取日记详情
	diary, err := h.service.GetDiaryByID(diaryId)
	if err != nil || diary == nil {
		h.renderService.Redirect(c, "/")
		return
	}

	// 获取当前用户信息（可能未登录）
	isAuth, authUser := h.auth.IsAuthenticated(c)
	isOwner := isAuth && (authUser.UserId == diary.AddUserID)

	// 构建页面数据
	data := h.renderService.PageData(c, "diary.detail.title")
	data.Data["Diary"] = diary
	data.Data["IsOwner"] = isOwner

	h.renderService.HTML(c, "diary/diary", data)
}
