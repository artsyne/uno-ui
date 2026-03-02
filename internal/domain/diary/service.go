package diary

import (
	"bytes"
	"database/sql"
	"embed"
	std_errors "errors"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/pkg"
)

//go:embed templates/invitation.html
var invitationTemplateFS embed.FS

// Service 日记本业务逻辑层实现
// 负责处理所有与日记本相关的业务逻辑和规则验证
type Service struct {
	config         *config.Config      // 应用配置信息
	repo           RepositoryInterface // 数据访问层接口
	userService    *user.Service       // 用户服务
	taskService    *task.Service       // 任务服务
	i18nService    i18n.I18nManager    // 国际化服务
	invitationTmpl *template.Template  // 邀请邮件模板
}

// NewService 创建新的日记本业务逻辑层实例
// 参数:
//   - config: 应用配置对象
//   - repo: 数据访问层接口实现
//   - userService: 用户服务实例
//   - taskService: 任务服务实例
//   - i18nService: 国际化服务实例
//
// 返回:
//   - *Service: 日记本业务逻辑层实例
func NewService(config *config.Config, repo RepositoryInterface, userService *user.Service, taskService *task.Service, i18nService i18n.I18nManager) *Service {
	tmpl, err := template.ParseFS(invitationTemplateFS, "templates/invitation.html")
	if err != nil {
		panic(fmt.Sprintf("failed to parse invitation template: %v", err))
	}

	return &Service{
		config:         config,
		repo:           repo,
		userService:    userService,
		taskService:    taskService,
		i18nService:    i18nService,
		invitationTmpl: tmpl,
	}
}

// ===================== 接口实现方法 =====================

// AddDiary 创建一个新的日记本
// 参数:
//   - diary: 要创建的日记本对象
//
// 返回:
//   - int64: 新创建的日记本ID
//   - error: 创建失败时返回错误
func (s *Service) AddDiary(diary *Diary) (int64, error) {
	// 验证用户ID有效性
	if diary.AddUserID <= 0 {
		return 0, &errors.InvalidParams
	}

	// 验证标题不能为空
	if len(diary.Title) < 1 {
		return 0, &errors.InvalidParams
	}

	// 判断用户增加日记本数量是否超过限制
	if err := s.CheckDiaryLimit(diary.AddUserID); err != nil {
		return 0, err
	}

	// 设置默认状态
	diary.Status = StatusActive
	if diary.ThemeID == 0 {
		diary.ThemeID = 1
	}

	// 新增日记本
	diaryID, err := s.repo.Insert(diary)
	if err != nil {
		return 0, err
	}
	diary.DiaryID = diaryID

	// 将拥有者添加为作者（管理员）
	ownerAuthor := &Author{
		DiaryID:     diaryID,
		UserID:      diary.AddUserID,
		IsAdmin:     AuthorIsAdminYes,
		JoinUTCTime: time.Now().Unix(),
		PostNum:     0,
		Status:      AuthorStatusActive,
	}
	if err := s.repo.InsertAuthor(ownerAuthor); err != nil {
		// 如果添加作者失败，尝试回滚（删除日记本）
		_ = s.repo.Delete(diaryID)
		return 0, err
	}

	return diaryID, nil
}

// UpdateDiary 更新日记本信息
// 参数:
//   - diary: 要更新的日记本对象
//
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdateDiary(diary *Diary) error {
	// 校验参数
	if err := s.validateDiary(diary); err != nil {
		return err
	}

	// 更新日记本
	err := s.repo.Update(diary)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDiary 删除日记本
// 参数:
//   - diaryID: 日记本ID
//   - userID: 操作用户ID
//
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeleteDiary(diaryID, userID int64) error {
	// 验证所有权
	isOwner, err := s.repo.IsDiaryOwner(diaryID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return &errors.PermissionDenied
	}

	// 检查是否有日记
	diary, err := s.repo.GetByDiaryID(diaryID)
	if err != nil {
		return err
	}
	if diary.InfoNum > 0 {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "diary_not_empty"})
	}

	// 删除日记本
	return s.repo.Delete(diaryID)
}

// GetDiaryByID 根据ID获取日记本
// 参数:
//   - diaryID: 日记本ID
//
// 返回:
//   - *Diary: 日记本对象，未找到时返回nil
//   - error: 查询失败时返回错误
func (s *Service) GetDiaryByID(diaryID int64) (*Diary, error) {
	return s.repo.GetByDiaryID(diaryID)
}

// CheckSocialEnabled 检查日记本社交功能是否开启
// 参数:
//   - diaryID: 日记本ID
//
// 返回:
//   - *Diary: 日记本对象
//   - error: 如果日记本不存在或社交功能关闭，返回相应错误
func (s *Service) CheckSocialEnabled(diaryID int64) (*Diary, error) {
	diary, err := s.GetDiaryByID(diaryID)
	if err != nil {
		return nil, &errors.DiaryNotFound
	}
	if diary.SocialSetting == SocialSettingDisabled {
		return nil, &errors.PermissionDenied
	}
	return diary, nil
}

// ListDiariesForWrite 获取用户的日记本列表（用于写日记界面）
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - []Diary: 日记本列表，按最后更新时间排序
//   - error: 查询失败时返回错误
func (s *Service) ListDiariesForWrite(userID int64) ([]Diary, error) {
	return s.repo.ListForWriteByUserID(userID)
}

// ListDiariesByUser 获取用户的分页日记本列表
// 参数:
//   - userID: 用户ID
//   - includePrivate: 是否包含私有日记本
//   - pageIndex: 页索引（从0开始）
//   - pageSize: 每页大小
//
// 返回:
//   - []Diary: 日记本列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListDiariesByUser(userID int64, includePrivate bool, pageIndex, pageSize int) ([]Diary, int, error) {
	diaries, _, err := s.repo.ListByUserID(userID, includePrivate, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if len(diaries) == 0 {
		return []Diary{}, 0, nil
	}

	count, err := s.repo.GetCountByUserID(userID, includePrivate)
	if err != nil {
		return nil, 0, err
	}

	return diaries, count, nil
}

// IncrementEntryCount 增加日记本的日记数量
// 参数:
//   - diaryID: 日记本ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) IncrementEntryCount(diaryID int64) error {
	return s.repo.IncrementEntryNum(diaryID)
}

// DecrementEntryCount 减少日记本的日记数量
// 参数:
//   - diaryID: 日记本ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) DecrementEntryCount(diaryID int64) error {
	return s.repo.DecrementEntryNum(diaryID)
}

// ===================== 多作者相关方法 =====================

// InviteAuthor 邀请作者加入日记本
// 参数:
//   - diaryID: 日记本ID
//   - ownerID: 操作用户ID（必须是日记本拥有者）
//   - inviteeIdentity: 被邀请人标识（邮箱或Handle）
//   - ownerLang: 拥有者语言（用于邮件通知）
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) InviteAuthor(diaryID, ownerID int64, inviteeIdentity string, ownerLang string) error {
	// 1. 验证操作用户是否为日记本拥有者
	isOwner, err := s.repo.IsDiaryOwner(diaryID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_diary_owner"})
	}

	// 2. 解析被邀请用户（Handle 或 Email）
	var user *user.User
	if pkg.IsEmail(inviteeIdentity) {
		user, err = s.userService.GetUserByEmail(inviteeIdentity)
	} else {
		// 尝试解析为 userID（数字）或 Handle
		if id, parseErr := strconv.ParseInt(inviteeIdentity, 10, 64); parseErr == nil && id > 0 {
			user, err = s.userService.GetUserByUserId(id)
		} else {
			user, err = s.userService.GetUserByHandle(inviteeIdentity)
		}
	}

	if err != nil {
		// 用户未找到
		return errors.UserNotFound.WithDetails(map[string]interface{}{"identity": inviteeIdentity})
	}
	inviteeID := user.UserId

	// 不能邀请自己
	if inviteeID == ownerID {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "cannot_invite_self"})
	}

	// 3. 检查是否已经是作者或已邀请
	existingAuthor, err := s.repo.GetDiaryAuthor(diaryID, inviteeID)
	if err != nil && err != sql.ErrNoRows && !std_errors.Is(err, errors.AuthorNotFound) {
		return err
	}
	if existingAuthor != nil {
		if existingAuthor.Status == AuthorStatusActive {
			return errors.AlreadyAuthor
		}
		return errors.AlreadyInvited
	}

	// 4. 获取日记本信息（用于邮件）
	diary, err := s.repo.GetByDiaryID(diaryID)
	if err != nil {
		return err
	}

	// 5. 创建邀请记录 (Pending状态)
	author := &Author{
		DiaryID:     diaryID,
		UserID:      inviteeID,
		IsAdmin:     AuthorIsAdminNo,
		JoinUTCTime: time.Now().Unix(),
		PostNum:     0,
		Status:      AuthorStatusPending,
	}

	err = s.repo.InsertAuthor(author)
	if err != nil {
		return err
	}

	// 6. 发送邀请邮件
	// 这里我们异步发送邮件，不阻塞主流程
	go func() {
		if err := s.sendInvitationEmail(diary, user, ownerLang); err != nil {
			// 记录日志，但不返回错误给用户，因为邀请已经创建成功
			fmt.Printf("Failed to send invitation email: %v\n", err)
		}
	}()

	return nil
}

// RemoveAuthor 移除作者（或取消邀请）
// 参数:
//   - diaryID: 日记本ID
//   - ownerID: 操作用户ID（必须是日记本拥有者）
//   - targetUserID: 目标作者ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) RemoveAuthor(diaryID, ownerID, targetUserID int64) error {
	// 1. 验证操作用户是否为日记本拥有者
	isOwner, err := s.repo.IsDiaryOwner(diaryID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_diary_owner"})
	}

	// 2. 不能移除自己（拥有者）
	if ownerID == targetUserID {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "cannot_remove_owner"})
	}

	// 3. 检查目标用户是否存在于该日记本
	_, err = s.repo.GetDiaryAuthor(diaryID, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.AuthorNotFound
		}
		return err
	}

	// 4. 删除作者
	return s.repo.DeleteAuthor(diaryID, targetUserID)
}

// ListAuthors 获取日记本的所有作者
// 参数:
//   - diaryID: 日记本ID
//   - onlyActive: 是否只获取活跃作者
//
// 返回:
//   - []Author: 作者列表
//   - error: 查询失败时返回错误
func (s *Service) ListAuthors(diaryID int64, onlyActive bool) ([]Author, error) {
	authors, _, err := s.repo.ListAuthorsByDiaryID(diaryID, onlyActive, 0, 0)
	return authors, err
}

// GetAuthor 获取指定作者信息
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回:
//   - *Author: 作者信息
//   - error: 查询失败时返回错误
func (s *Service) GetAuthor(diaryID, userID int64) (*Author, error) {
	return s.repo.GetDiaryAuthor(diaryID, userID)
}

// UpdateAuthorAdmin 更新作者管理员权限
// 参数:
//   - diaryID: 日记本ID
//   - ownerID: 操作用户ID（必须是日记本拥有者）
//   - targetUserID: 目标作者ID
//   - isAdmin: 是否管理员 (0:否, 1:是)
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) UpdateAuthorAdmin(diaryID, ownerID, targetUserID int64, isAdmin int8) error {
	// 1. 验证操作用户是否为日记本拥有者
	isOwner, err := s.repo.IsDiaryOwner(diaryID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_diary_owner"})
	}

	// 2. 检查目标用户是否为活跃作者
	author, err := s.repo.GetDiaryAuthor(diaryID, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.AuthorNotFound
		}
		return err
	}
	if author.Status != AuthorStatusActive {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "author_not_active"})
	}

	// 3. 更新权限
	return s.repo.UpdateAuthorAdmin(diaryID, targetUserID, isAdmin)
}

// AcceptInvitation 接受邀请
// 参数:
//   - diaryID: 日记本ID
//   - inviteeID: 被邀请人ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) AcceptInvitation(diaryID, inviteeID int64) error {
	// 1. 获取邀请记录
	author, err := s.repo.GetDiaryAuthor(diaryID, inviteeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.InvitationNotFound
		}
		return err
	}

	// 2. 验证状态
	if author.Status == AuthorStatusActive {
		return errors.AlreadyAuthor
	}
	if author.Status != AuthorStatusPending {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_status"})
	}

	// 3. 检查加入数量限制
	if err := s.CheckJoinedDiariesLimit(inviteeID); err != nil {
		return err
	}

	// 4. 更新状态
	return s.repo.UpdateAuthorStatus(diaryID, inviteeID, AuthorStatusActive)
}

// RejectInvitation 拒绝邀请
// 参数:
//   - diaryID: 日记本ID
//   - inviteeID: 被邀请人ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) RejectInvitation(diaryID, inviteeID int64) error {
	// 1. 获取邀请记录
	author, err := s.repo.GetDiaryAuthor(diaryID, inviteeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.InvitationNotFound
		}
		return err
	}

	// 2. 验证状态
	if author.Status != AuthorStatusPending {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_status"})
	}

	// 3. 删除记录
	return s.repo.DeleteAuthor(diaryID, inviteeID)
}

// GetDiaryAuthors 获取日记本的所有作者(支持分页)
// 参数:
//   - diaryID: 日记本ID
//   - onlyActive: 是否只获取活跃作者
//   - pageIndex, pageSize: 分页参数
//
// 返回:
//   - []Author: 作者列表
//   - int64: 总数
//   - error: 查询失败时返回错误
func (s *Service) GetDiaryAuthors(diaryID int64, onlyActive bool, pageIndex, pageSize int) ([]Author, int64, error) {
	authors, count, err := s.repo.ListAuthorsByDiaryID(diaryID, onlyActive, pageIndex, pageSize)
	return authors, int64(count), err
}

// IsAuthor 检查用户是否为活跃作者
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回:
//   - bool: 是否为活跃作者
//   - error: 查询失败时返回错误
func (s *Service) IsAuthor(diaryID, userID int64) (bool, error) {
	return s.repo.IsDiaryAuthor(diaryID, userID)
}

// GetAuthorRole 获取作者在日记本中的角色
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回:
//   - int8: 角色 (1:Owner, 2:Admin, 3:Member)
//   - error: 查询失败时返回错误（如非作者）
func (s *Service) GetAuthorRole(diaryID, userID int64) (int8, error) {
	// 1. 检查是否为拥有者
	isOwner, err := s.repo.IsDiaryOwner(diaryID, userID)
	if err != nil {
		return 0, err
	}
	if isOwner {
		return 1, nil // Owner
	}

	// 2. 获取作者信息
	author, err := s.repo.GetDiaryAuthor(diaryID, userID)
	if err != nil {
		return 0, err
	}

	// 3. 判断是否为管理员
	if author.IsAdmin == AuthorIsAdminYes {
		return 2, nil // Admin
	}

	return 3, nil // Member
}

// QuitDiary 退出日记本
// 参数:
//   - diaryID: 日记本ID
//   - userID: 操作用户ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) QuitDiary(diaryID, userID int64) error {
	// 1. 检查是否为拥有者（拥有者不能退出）
	isOwner, err := s.repo.IsDiaryOwner(diaryID, userID)
	if err != nil {
		return err
	}
	if isOwner {
		return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "owner_cannot_quit"})
	}

	// 2. 检查是否为作者
	_, err = s.repo.GetDiaryAuthor(diaryID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.AuthorNotFound
		}
		return err
	}

	// 3. 删除记录
	return s.repo.DeleteAuthor(diaryID, userID)
}

// ListJoinedDiariesByUserID 获取用户参与的日记本列表（状态为Active）
func (s *Service) ListJoinedDiariesByUserID(userID int64) ([]Diary, error) {
	return s.repo.ListJoinedDiariesByUserID(userID)
}

// ListJoinedDiariesByUserIDWithLimit 获取用户参与的日记本列表（带限制）
func (s *Service) ListJoinedDiariesByUserIDWithLimit(userID int64, includePrivate bool, limit int) ([]Diary, error) {
	return s.repo.ListJoinedDiariesByUserIDWithLimit(userID, includePrivate, limit)
}

// ListJoinedDiariesByUserIDPaged 获取用户参与的分页日记本列表
func (s *Service) ListJoinedDiariesByUserIDPaged(userID int64, includePrivate bool, pageIndex, pageSize int) ([]Diary, int, error) {
	diaries, _, err := s.repo.ListJoinedDiariesByUserIDPaged(userID, includePrivate, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if len(diaries) == 0 {
		return []Diary{}, 0, nil
	}

	count, err := s.repo.GetJoinedDiariesCountByUserID(userID, includePrivate)
	if err != nil {
		return nil, 0, err
	}

	return diaries, count, nil
}

// ListPendingInvitationsByUserID 获取用户的待确认邀请列表
func (s *Service) ListPendingInvitationsByUserID(userID int64) ([]Diary, error) {
	return s.repo.ListPendingInvitationsByUserID(userID)
}

// sendInvitationEmail 发送邀请邮件
func (s *Service) sendInvitationEmail(diary *Diary, invitee *user.User, lang string) error {
	// 规范化语言代码
	lang = s.normalizeLanguage(lang)

	// 构建接受邀请的链接 (假设路由是 /diary/accept_invitation?diary_id=xxx)
	// 实际项目中应该生成一个带签名的 token
	acceptLink := fmt.Sprintf("http://%s/diary/invitation?id=%d", s.config.App.Domain, diary.DiaryID)

	// 准备翻译参数
	params := map[string]string{
		"diary_title": diary.Title,
		"name":        invitee.DisplayName,
	}

	// 获取翻译内容
	subject := s.i18nService.GetMessageWithParams("email.invitation.subject", lang, params)
	greeting := s.i18nService.GetMessageWithParams("email.invitation.greeting", lang, params)
	bodyText := s.i18nService.GetMessageWithParams("email.invitation.body", lang, params)
	roleDesc := s.i18nService.GetMessage("email.invitation.role_desc", lang)
	actionBtn := s.i18nService.GetMessage("email.invitation.action_btn", lang)
	altLinkDesc := s.i18nService.GetMessage("email.invitation.alt_link_desc", lang)
	ignoreDesc := s.i18nService.GetMessage("email.invitation.ignore_desc", lang)

	data := struct {
		InviteeName string
		DiaryTitle  string
		AcceptLink  string
		// Translated fields
		Subject     string
		Greeting    string
		Body        string
		RoleDesc    string
		ActionBtn   string
		AltLinkDesc string
		IgnoreDesc  string
	}{
		InviteeName: invitee.DisplayName,
		DiaryTitle:  diary.Title,
		AcceptLink:  acceptLink,
		Subject:     subject,
		Greeting:    greeting,
		Body:        bodyText,
		RoleDesc:    roleDesc,
		ActionBtn:   actionBtn,
		AltLinkDesc: altLinkDesc,
		IgnoreDesc:  ignoreDesc,
	}

	var body bytes.Buffer
	if err := s.invitationTmpl.Execute(&body, data); err != nil {
		return err
	}

	return s.taskService.CreateEmailTask(invitee.Email, subject, body.String())
}

// normalizeLanguage 规范化语言代码
func (s *Service) normalizeLanguage(lang string) string {
	if s.i18nService.IsLanguageSupported(lang) {
		return lang
	}
	return s.i18nService.GetDefaultLanguage()
}

// CheckDiaryLimit 检查用户日记本数量是否超过限制
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - error: 超过限制时返回错误
func (s *Service) CheckDiaryLimit(userID int64) error {
	daybookNum, err := s.repo.GetCountByUserID(userID, true)
	if err != nil {
		return &errors.InvalidParams
	}

	if daybookNum > s.config.App.MaxDiaryLimit {
		return &errors.DiaryLimitExceeded
	}

	return nil
}

// CheckJoinedDiariesLimit 检查用户加入日记本数量是否超过限制
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - error: 超过限制时返回错误
func (s *Service) CheckJoinedDiariesLimit(userID int64) error {
	count, err := s.repo.GetActiveJoinedDiariesCountByUserID(userID)
	if err != nil {
		return &errors.InvalidParams
	}

	// DEBUG: Print count and limit
	// fmt.Printf("DEBUG: CheckJoinedDiariesLimit userID=%d count=%d limit=%d\n", userID, count, s.config.App.UserJoinedDiariesLimit)

	if count >= s.config.App.UserJoinedDiariesLimit {
		return &errors.JoinedDiariesLimitExceeded
	}

	return nil
}

// IsDiaryOwner 验证用户是否是日记本的所有者
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回:
//   - bool: 是否为拥有者
//   - error: 验证失败时返回错误
func (s *Service) IsDiaryOwner(diaryID, userID int64) (bool, error) {
	IsDiaryOwner, err := s.repo.IsDiaryOwner(diaryID, userID)
	if err != nil {
		return false, err
	}

	return IsDiaryOwner, nil
}

// IncrementLikeCount 增加日记本的点赞数量
// 参数:
//   - diaryID: 日记本ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) IncrementLikeCount(diaryID int64) error {
	return s.repo.IncrementLikeNum(diaryID)
}

// DecrementLikeCount 减少日记本的点赞数量
// 参数:
//   - diaryID: 日记本ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) DecrementLikeCount(diaryID int64) error {
	return s.repo.DecrementLikeNum(diaryID)
}

// IsActiveAuthor 检查用户是否为日记本的活跃作者
// 参数:
//   - diaryID: 日记本ID
//   - userID: 用户ID
//
// 返回:
//   - bool: 是否为活跃作者
//   - error: 查询失败时返回错误
func (s *Service) IsActiveAuthor(diaryID, userID int64) (bool, error) {
	author, err := s.repo.GetDiaryAuthor(diaryID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		// Check for AuthorNotFound error
		if appErr, ok := err.(errors.AppError); ok {
			if appErr.Code() == errors.AuthorNotFound.Code() {
				return false, nil
			}
		}
		return false, err
	}
	return author.Status == AuthorStatusActive, nil
}

// ===================== 私有辅助方法 =====================

// validateDiary 校验日记本参数
// 参数:
//   - diary: 要校验的日记本对象
//
// 返回:
//   - error: 校验失败时返回错误
func (s *Service) validateDiary(diary *Diary) error {
	if len(diary.Title) < 1 {
		return &errors.InvalidParams
	}
	if diary.ViewSetting != ViewSettingPublic && diary.ViewSetting != ViewSettingPrivate {
		return &errors.InvalidParams
	}
	if diary.SocialSetting != 1 && diary.SocialSetting != 2 {
		return &errors.InvalidParams
	}
	if diary.ThemeID < 1 || diary.ThemeID > 16 {
		return &errors.InvalidParams
	}
	return nil
}
