# 多作者日记本功能设计方案

## 1. 概述 (Overview)

本设计旨在实现日记本的多作者（共创）功能。该功能允许日记本所有者邀请其他用户（作者）加入并贡献日记内容。受邀用户可以接受邀请，并在加入的日记本中管理自己创作的日记。

## 2. 数据库结构 (Database Schema)

基于现有的 `br_diary_author` 表结构（已确认 `sql/schema.sql`）：

```sql
CREATE TABLE `br_diary_author` (
  `diary_id` bigint unsigned NOT NULL COMMENT '关联日记本ID',
  `user_id` bigint unsigned NOT NULL COMMENT '关联用户ID',
  `is_admin` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否管理员(0:否 1:是)',
  `join_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '加入时间(UTC时间戳)',
  `post_num` int unsigned NOT NULL DEFAULT '0' COMMENT '发表文章数',
  `status` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '状态(1:活跃 2:待确认)',
  PRIMARY KEY (`diary_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日记本作者表';
```

### 2.1 字段说明

| 字段            | 类型             | 说明                                       |
| --------------- | ---------------- | ------------------------------------------ |
| `diary_id`      | bigint unsigned  | 日记本ID，外键关联 `br_diary.diary_id`     |
| `user_id`       | bigint unsigned  | 用户ID，外键关联 `br_user.user_id`         |
| `is_admin`      | tinyint unsigned | 是否为管理员：0=普通作者，1=管理员         |
| `join_utc_time` | bigint unsigned  | 加入时间（UTC时间戳，毫秒）                |
| `post_num`      | int unsigned     | 该作者在该日记本中发表的日记数量           |
| `status`        | tinyint unsigned | 状态：1=活跃（已加入），2=待确认（已邀请） |

### 2.2 索引说明

- 主键：`(diary_id, user_id)` - 联合主键，确保一个用户在一个日记本中只能有一条记录
- 建议添加的索引（可根据实际查询需求添加）：
  - `idx_user_id_status` - 用于查询用户参与的日记本列表
  - `idx_status` - 用于查询待确认的邀请

## 3. 领域逻辑更新 (Domain Logic Updates)

### 3.1 日记本领域 (`internal/domain/diary`)

#### 3.1.1 类型定义 (`types.go`)

```go
// Author 日记本作者结构体
// 对应数据库 br_diary_author 表
type Author struct {
    DiaryID     int64  `json:"diary_id" db:"diary_id"`           // 日记本ID
    UserID      int64  `json:"user_id" db:"user_id"`             // 用户ID
    IsAdmin     int8   `json:"is_admin" db:"is_admin"`           // 是否管理员 (0:否, 1:是)
    JoinUTCTime int64  `json:"join_utc_time" db:"join_utc_time"` // 加入时间
    PostNum     int    `json:"post_num" db:"post_num"`           // 发表文章数
    Status      int8   `json:"status" db:"status"`               // 状态 (1:Active, 2:Pending)

    // Joined fields from br_user
    Handle      string `json:"handle" db:"handle"`               // 用户Handle (Join查询)
    DisplayName string `json:"display_name" db:"display_name"`   // 用户显示名称 (Join查询)
    Avatar      string `json:"avatar" db:"avatar"`               // 用户头像 (Join查询)
}

// AuthorStatus 作者状态常量
const (
    AuthorStatusActive  int8 = 1 // 活跃/已加入
    AuthorStatusPending int8 = 2 // 待确认/已邀请
)

// AuthorAdminStatus 管理员状态常量
const (
    AuthorIsAdminNo  int8 = 0 // 普通作者
    AuthorIsAdminYes int8 = 1 // 管理员
)
```

#### 3.1.2 Service 结构体定义

```go
// Service 日记本业务逻辑层实现
type Service struct {
    config       *config.Config
    repo         RepositoryInterface
    userService  *user.Service
    taskService  *task.Service
    i18nService  i18n.I18nManager
}

// NewService 创建新的日记本业务逻辑层实例
func NewService(config *config.Config, repo RepositoryInterface, userService *user.Service, taskService *task.Service, i18nService i18n.I18nManager) *Service {
    return &Service{
        config:      config,
        repo:        repo,
        userService: userService,
        taskService: taskService,
        i18nService: i18nService,
    }
}
```

#### 3.1.3 Repository 接口 (`repository.go`)

```go
// RepositoryInterface 新增作者相关方法
type RepositoryInterface interface {
    // ... 原有方法 ...

    // ===================== Author 相关方法 =====================

    // InsertAuthor 新增日记本作者
    InsertAuthor(author *Author) error

    // UpdateAuthorStatus 更新作者状态
    UpdateAuthorStatus(diaryID, userID int64, status int8) error

    // UpdateAuthorAdmin 更新作者管理员权限
    UpdateAuthorAdmin(diaryID, userID int64, isAdmin int8) error

    // DeleteAuthor 删除作者
    DeleteAuthor(diaryID, userID int64) error

    // ListAuthorsByDiaryID 获取日记本的所有作者
    ListAuthorsByDiaryID(diaryID int64) ([]Author, error)

    // ListJoinedDiariesByUserID 获取用户参与的日记本列表（状态为Active）
    ListJoinedDiariesByUserID(userID int64) ([]Diary, error)

    // ListJoinedDiariesByUserIDPaged 获取用户参与的分页日记本列表
    ListJoinedDiariesByUserIDPaged(userID int64, pageIndex, pageSize int) ([]Diary, int, error)

    // GetJoinedDiariesCountByUserID 获取用户参与的日记本数量
    GetJoinedDiariesCountByUserID(userID int64) (int, error)

    // GetDiaryAuthor 获取指定作者信息
    GetDiaryAuthor(diaryID, userID int64) (*Author, error)

    // IsDiaryAuthor 检查用户是否为活跃作者
    IsDiaryAuthor(diaryID, userID int64) (bool, error)

    // ListPendingInvitationsByUserID 获取用户的待确认邀请列表
    ListPendingInvitationsByUserID(userID int64) ([]Diary, error)

    // IncrementAuthorPostNum 增加作者的发表数量
    IncrementAuthorPostNum(diaryID, userID int64) error

    // DecrementAuthorPostNum 减少作者的发表数量
    DecrementAuthorPostNum(diaryID, userID int64) error
}
```

#### 3.1.4 Service 层实现 (`service.go`)

```go
// InviteAuthor 邀请作者加入日记本
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
    var inviteeID int64
    var user *user.User
    if pkg.IsEmail(inviteeIdentity) {
        user, err = s.userService.GetUserByEmail(inviteeIdentity)
    } else {
        // 尝试解析为 userID（数字）或 Handle
        if id, err := strconv.ParseInt(inviteeIdentity, 10, 64); err == nil && id > 0 {
            user, err = s.userService.GetUserByUserId(id)
        } else {
            user, err = s.userService.GetUserByHandle(inviteeIdentity)
        }
    }

    if err != nil {
        return errors.UserNotFound.WithDetails(map[string]interface{}{"identity": inviteeIdentity})
    }
    inviteeID = user.UserId

    // 3. 检查是否已经是作者或已邀请
    existingAuthor, err := s.repo.GetDiaryAuthor(diaryID, inviteeID)
    if err != nil && err != sql.ErrNoRows {
        return err
    }
    if existingAuthor != nil {
        if existingAuthor.Status == AuthorStatusActive {
            return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "already_author"})
        }
        return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "already_invited"})
    }

    // 4. 获取日记本信息（用于邮件）
    diary, err := s.repo.GetByDiaryID(diaryID)
    if err != nil {
        return err
    }

    // 5. 创建邀请记录
    author := &Author{
        DiaryID:     diaryID,
        UserID:      inviteeID,
        IsAdmin:     AuthorIsAdminNo,
        JoinUTCTime: pkg.GetNowTimeUnixMilli(),
        PostNum:     0,
        Status:      AuthorStatusPending,
    }

    err = s.repo.InsertAuthor(author)
    if err != nil {
        return err
    }

    // 6. 发送邀请邮件
    if err := s.sendInvitationEmail(diary, user, ownerLang); err != nil {
        slog.Error("发送邀请邮件失败", "error", err, "diary_id", diaryID, "user_id", inviteeID)
        // 邮件发送失败不影响邀请记录的创建
    }

    return nil
}

// sendInvitationEmail 发送邀请邮件
// 使用被邀请用户的语言设置发送邮件
func (s *Service) sendInvitationEmail(diary *Diary, invitee *user.User, ownerLang string) error {
    // 获取被邀请用户的语言设置，如果未设置则使用拥有者的语言
    userLang := invitee.Lang
    if userLang == "" {
        userLang = ownerLang
    }

    // 规范化语言代码
    normalizedLang := "en"
    if userLang == "zh-cn" || userLang == "zh" {
        normalizedLang = "zh-cn"
    }

    // 使用被邀请用户的语言获取邮件内容
    subject := s.i18nService.GetMessage("email.diary_invitation.subject", normalizedLang)

    // 构建邮件正文（支持 HTML）
    params := map[string]string{
        "diary_title": diary.Title,
        "owner_name":   diary.AddUserDisplayName,
        "owner_handle": diary.AddUserHandle,
        "invitee_name": invitee.DisplayName,
        "accept_url":   fmt.Sprintf("%s/diary/%d/join", s.config.App.Domain, diary.DiaryID),
        "domain":       s.config.App.Domain,
    }

    body := s.i18nService.GetMessageWithParams("email.diary_invitation.body", normalizedLang, params)

    // 发送邮件任务
    return s.taskService.CreateEmailTask(invitee.Email, subject, body)
}

// normalizeLanguage 规范化语言代码
func (s *Service) normalizeLanguage(lang string) string {
    lang = strings.ToLower(lang)
    switch lang {
    case "zh", "zh-hans", "zh_hans", "zh-cn", "zh_cn":
        return "zh-cn"
    case "zh-hant", "zh_hant", "zh-tw", "zh_tw", "zh-hk", "zh_hk":
        return "zh-tw"
    case "en", "en-us", "en_us":
        return "en"
    default:
        // 检查是否为支持的语言
        if s.i18nService.IsLanguageSupported(lang) {
            return lang
        }
        // 默认使用中文
        return "zh-cn"
    }
}

// SetAuthorAdmin 设置/取消作者管理员权限
func (s *Service) SetAuthorAdmin(diaryID, ownerID, targetUserID int64, isAdmin bool) error {
    // 1. 验证操作用户是否为日记本拥有者
    isOwner, err := s.repo.IsDiaryOwner(diaryID, ownerID)
    if err != nil {
        return err
    }
    if !isOwner {
        return errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_diary_owner"})
    }

    // 2. 不能设置拥有者为管理员
    if targetUserID == ownerID {
        return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "cannot_set_owner"})
    }

    // 3. 检查目标用户是否为作者
    author, err := s.repo.GetDiaryAuthor(diaryID, targetUserID)
    if err != nil {
        return errors.AuthorNotFound
    }

    // 4. 更新管理员状态
    var adminStatus int8 = AuthorIsAdminNo
    if isAdmin {
        adminStatus = AuthorIsAdminYes
    }

    return s.repo.UpdateAuthorAdmin(diaryID, targetUserID, adminStatus)
}

// AcceptInvitation 接受日记本邀请
func (s *Service) AcceptInvitation(diaryID, userID int64) error {
    // 1. 检查邀请是否存在
    author, err := s.repo.GetDiaryAuthor(diaryID, userID)
    if err != nil {
        return errors.InvitationNotFound
    }

    // 2. 检查邀请状态
    if author.Status != AuthorStatusPending {
        return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "not_pending_invitation"})
    }

    // 3. 更新状态为活跃
    return s.repo.UpdateAuthorStatus(diaryID, userID, AuthorStatusActive)
}

// DeclineInvitation 拒绝日记本邀请
func (s *Service) DeclineInvitation(diaryID, userID int64) error {
    // 1. 检查邀请是否存在
    author, err := s.repo.GetDiaryAuthor(diaryID, userID)
    if err != nil {
        return errors.InvitationNotFound
    }

    // 2. 检查邀请状态
    if author.Status != AuthorStatusPending {
        return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "not_pending_invitation"})
    }

    // 3. 删除邀请记录
    return s.repo.DeleteAuthor(diaryID, userID)
}

// QuitDiary 退出日记本
func (s *Service) QuitDiary(diaryID, userID int64) error {
    // 1. 检查是否为作者
    author, err := s.repo.GetDiaryAuthor(diaryID, userID)
    if err != nil {
        return errors.AuthorNotFound
    }

    // 2. 拥有者不能退出自己的日记本
    diary, err := s.repo.GetByDiaryID(diaryID)
    if err != nil {
        return err
    }
    if diary.AddUserID == userID {
        return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "owner_cannot_quit"})
    }

    // 3. 删除作者记录
    return s.repo.DeleteAuthor(diaryID, userID)
}

// RemoveAuthor 移除作者（仅限拥有者）
func (s *Service) RemoveAuthor(diaryID, ownerID, targetUserID int64) error {
    // 1. 验证操作用户是否为日记本拥有者
    isOwner, err := s.repo.IsDiaryOwner(diaryID, ownerID)
    if err != nil {
        return err
    }
    if !isOwner {
        return errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_diary_owner"})
    }

    // 2. 不能移除拥有者
    if targetUserID == ownerID {
        return errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "cannot_remove_owner"})
    }

    // 3. 检查目标用户是否为作者
    _, err = s.repo.GetDiaryAuthor(diaryID, targetUserID)
    if err != nil {
        return errors.AuthorNotFound
    }

    // 4. 删除作者记录
    return s.repo.DeleteAuthor(diaryID, targetUserID)
}

// ListAuthors 获取日记本的所有作者
func (s *Service) ListAuthors(diaryID, requestUserID int64) ([]Author, *Diary, bool, error) {
    // 1. 获取日记本信息
    diary, err := s.repo.GetByDiaryID(diaryID)
    if err != nil {
        return nil, nil, false, err
    }

    // 2. 验证查看权限
    isOwner := (diary.AddUserID == requestUserID)
    isAuthor, _ := s.repo.IsDiaryAuthor(diaryID, requestUserID)

    if diary.ViewSetting == 2 { // 私有日记本
        if !isOwner && !isAuthor {
            return nil, nil, false, errors.PermissionDenied
        }
    }

    // 3. 获取作者列表
    authors, err := s.repo.ListAuthorsByDiaryID(diaryID)
    if err != nil {
        return nil, nil, false, err
    }

    return authors, diary, isOwner, nil
}

// ListJoinedDiaries 获取用户参与的日记本列表
func (s *Service) ListJoinedDiaries(userID int64, pageIndex, pageSize int) ([]Diary, int, error) {
    return s.repo.ListJoinedDiariesByUserIDPaged(userID, pageIndex, pageSize)
}

// ListPendingInvitations 获取用户的待确认邀请列表
func (s *Service) ListPendingInvitations(userID int64) ([]Diary, error) {
    return s.repo.ListPendingInvitationsByUserID(userID)
}

// ListDiariesForWrite 获取用户可以写日记的日记本列表
// 包含：拥有的日记本 + 加入的日记本
func (s *Service) ListDiariesForWrite(userID int64) ([]Diary, error) {
    // 1. 获取拥有的日记本
    ownedDiaries, err := s.repo.ListForWriteByUserID(userID)
    if err != nil {
        return nil, err
    }

    // 2. 获取加入的日记本
    joinedDiaries, err := s.repo.ListJoinedDiariesByUserID(userID)
    if err != nil {
        return nil, err
    }

    // 3. 合并列表
    allDiaries := append(ownedDiaries, joinedDiaries...)

    // 4. 按 diary_id 去重（如果有重叠）
    uniqueDiaries := make([]Diary, 0)
    seen := make(map[int64]bool)
    for _, d := range allDiaries {
        if !seen[d.DiaryID] {
            seen[d.DiaryID] = true
            uniqueDiaries = append(uniqueDiaries, d)
        }
    }

    return uniqueDiaries, nil
}
```

### 3.2 日记领域 (`internal/domain/entry`)

#### 3.2.1 Service 层更新

```go
// checkWritePermission 验证用户是否有权限在日记本中写日记
func (s *Service) checkWritePermission(diaryID, userID int64) (bool, error) {
    // 1. 获取日记本信息
    diary, err := s.diaryService.GetDiaryByID(diaryID)
    if err != nil {
        return false, err
    }

    // 2. 拥有者有权限
    if diary.AddUserID == userID {
        return true, nil
    }

    // 3. 活跃作者有权限
    isAuthor, err := s.diaryService.IsDiaryAuthor(diaryID, userID)
    if err != nil {
        return false, err
    }

    return isAuthor, nil
}

// AddEntry 增加新的日记（更新权限验证）
func (s *Service) AddEntry(entry *Entry) (int64, error) {
    // ... 原有验证逻辑 ...

    // 验证用户是否有写权限
    hasPermission, err := s.checkWritePermission(entry.DiaryID, entry.AddUserID)
    if err != nil {
        return 0, err
    }
    if !hasPermission {
        return 0, errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "no_write_permission"})
    }

    // ... 插入逻辑 ...
}

// checkDeletePermission 验证用户是否有权限删除指定日记
func (s *Service) checkDeletePermission(entryID, userID int64) (bool, error) {
    // 1. 获取日记信息
    entryData, err := s.repo.GetByIDForUser(entryID, userID)
    if err != nil {
        return false, err
    }
    if entryData == nil {
        return false, errors.EntryNotFound
    }

    // 2. 获取日记本信息
    diary, err := s.diaryService.GetDiaryByID(entryData.DiaryID)
    if err != nil {
        return false, err
    }

    // 3. 日记本拥有者可以删除所有日记
    if diary.AddUserID == userID {
        return true, nil
    }

    // 4. 确保拥有者可以删除所有日记（包括管理员的日记）
    author, err := s.diaryService.GetDiaryAuthor(entryData.DiaryID, userID)
    if err != nil {
        return false, err
    }
    if author != nil && author.IsAdmin == AuthorIsAdminYes {
        // 管理员可以删除所有日记
        return true, nil
    }

    // 5. 普通作者只能删除自己的日记
    if entryData.AddUserID == userID {
        return true, nil
    }

    return false, errors.PermissionDenied
}

// DeleteEntryByID 根据日记ID删除日记信息（软删除，更新权限验证）
func (s *Service) DeleteEntryByID(entryID int64, userID int64, actionTime int64) error {
    // 1. 验证删除权限
    hasPermission, err := s.checkDeletePermission(entryID, userID)
    if err != nil {
        return err
    }
    if !hasPermission {
        return errors.PermissionDenied
    }

    // 2. 获取日记信息（用于减少统计数）
    entry, err := s.repo.GetByIDForUser(entryID, userID)
    if err != nil {
        return err
    }
    if entry == nil {
        return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
    }

    // 3. 执行软删除
    err = s.repo.DeleteByID(entryID)
    if err != nil {
        return err
    }

    // 4. 减少日记统计数（如果原状态是Active）
    if entry.Status == StatusActive {
        err = s.diaryService.DecrementEntryCount(entry.DiaryID)
        if err != nil {
            slog.Error("减少日记统计失败", "error", err, "entry_id", entryID, "diary_id", entry.DiaryID)
        }

        // 5. 减少作者的发表数量
        err = s.diaryService.DecrementAuthorPostNum(entry.DiaryID, entry.AddUserID)
        if err != nil {
            slog.Error("减少作者发表数量失败", "error", err, "entry_id", entryID, "author_id", entry.AddUserID)
        }
    }

    return nil
}
```

### 3.3 错误定义 (`internal/utils/errors`)

```go
// 新增多作者相关错误

var (
    // AuthorNotFound 作者未找到
    AuthorNotFound = NewAppError(10001, "AUTHOR_NOT_FOUND", "Author not found", http.StatusNotFound)

    // InvitationNotFound 邀请未找到
    InvitationNotFound = NewAppError(10002, "INVITATION_NOT_FOUND", "Invitation not found", http.StatusNotFound)

    // AlreadyAuthor 已经是作者
    AlreadyAuthor = NewAppError(10003, "ALREADY_AUTHOR", "User is already an author", http.StatusBadRequest)

    // AlreadyInvited 已经被邀请
    AlreadyInvited = NewAppError(10004, "ALREADY_INVITED", "User is already invited", http.StatusBadRequest)
)
)
```

### 3.4 权限矩阵 (Permissions)

| 角色                  | 定义                                   | is_admin        |
| --------------------- | -------------------------------------- | --------------- |
| **拥有者 (Owner)**    | 日记本的创建者，拥有最高权限           | N/A (AddUserID) |
| **管理员 (Admin)**    | 拥有者授权的管理者，可管理日记内容     | 1               |
| **普通作者 (Author)** | 被邀请的参与者，可创建和管理自己的日记 | 0               |
| **访客 (Guest)**      | 未登录或未加入的用户                   | N/A             |

**权限说明：**

1. **访问权限**：
   - 公开日记本（view_setting=1）：所有用户可访问
   - 私有日记本（view_setting=2）：拥有者、管理员、普通作者都可以访问
   - 密码保护日记本（view_setting=3）：拥有者、管理员、普通作者都可以访问；访客需要密码

2. **写权限**：
   - 拥有者、管理员、普通作者都可以在日记本中创建日记

3. **管理权限**（仅拥有者）：
   - 邀请作者
   - 设置/取消管理员权限
   - 移除作者
   - 编辑日记本设置

4. **删除权限**：
   - 拥有者可以删除任何日记
   - 管理员可以删除任何日记
   - 普通作者只能删除自己创建的日记

## 4. 控制器层设计 (`internal/app/controller/diary`)

### 4.1 请求定义 (`requests.go`)

```go
// InviteAuthorRequest 邀请作者请求
type InviteAuthorRequest struct {
    DiaryID         int64  `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
    InviteeIdentity string `form:"invitee_identity" json:"invitee_identity" binding:"required"` // Handle 或 Email
}

// SetAuthorAdminRequest 设置作者管理员请求
type SetAuthorAdminRequest struct {
    DiaryID int64 `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
    UserID  int64 `form:"user_id" json:"user_id" uri:"user_id" binding:"required"`
    IsAdmin bool  `form:"is_admin" json:"is_admin" binding:"required"`
}

// AcceptInvitationRequest 接受邀请请求
type AcceptInvitationRequest struct {
    DiaryID int64 `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
}

// DeclineInvitationRequest 拒绝邀请请求
type DeclineInvitationRequest struct {
    DiaryID int64 `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
}

// QuitDiaryRequest 退出日记本请求
type QuitDiaryRequest struct {
    DiaryID int64 `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
}

// RemoveAuthorRequest 移除作者请求
type RemoveAuthorRequest struct {
    DiaryID int64 `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
    UserID  int64 `form:"user_id" json:"user_id" uri:"user_id" binding:"required"`
}

// ListAuthorsRequest 获取作者列表请求
type ListAuthorsRequest struct {
    DiaryID int64 `form:"diary_id" json:"diary_id" uri:"id" binding:"required"`
}
```

### 4.2 处理器方法 (`handlers.go`)

```go
// ListAuthors 获取日记本作者列表
func (ctrl *DiaryController) ListAuthors(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[ListAuthorsRequest](c)

    authors, diary, isOwner, err := ctrl.diaryService.ListAuthors(req.DiaryID, loginUser.UserId)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        data := map[string]interface{}{
            "Diary":    diary,
            "Authors":  authors,
            "IsOwner":  isOwner,
            "UserID":   loginUser.UserId,
            "Lang":     loginUser.Lang,
            "TimeZone": loginUser.TimeZone,
        }
        c.HTML(http.StatusOK, "diary/authors_list", data)
        return
    }

    ctrl.renderService.JSONSuccess(c, map[string]interface{}{
        "authors":  authors,
        "diary":    diary,
        "is_owner": isOwner,
    })
}

// InviteAuthor 邀请作者
func (ctrl *DiaryController) InviteAuthor(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[InviteAuthorRequest](c)

    // 获取用户语言设置（用于邮件）
    userLang := loginUser.Lang
    if userLang == "" {
        userLang = "en"
    }

    err := ctrl.diaryService.InviteAuthor(req.DiaryID, loginUser.UserId, req.InviteeIdentity, userLang)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "diary.invitation_sent"}}`)
        ctrl.ListAuthors(c)
        return
    }

    ctrl.renderService.JSONSuccess(c, nil)
}

// SetAuthorAdmin 设置作者管理员权限
func (ctrl *DiaryController) SetAuthorAdmin(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[SetAuthorAdminRequest](c)

    err := ctrl.diaryService.SetAuthorAdmin(req.DiaryID, loginUser.UserId, req.UserID, req.IsAdmin)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "diary.admin_updated"}}`)
        ctrl.ListAuthors(c)
        return
    }

    ctrl.renderService.JSONSuccess(c, nil)
}

// AcceptInvitation 接受邀请
func (ctrl *DiaryController) AcceptInvitation(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[AcceptInvitationRequest](c)

    err := ctrl.diaryService.AcceptInvitation(req.DiaryID, loginUser.UserId)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "diary.joined_successfully"}, "refreshInvitations": true}`)
        c.Status(http.StatusOK)
        return
    }

    ctrl.renderService.JSONSuccess(c, nil)
}

// DeclineInvitation 拒绝邀请
func (ctrl *DiaryController) DeclineInvitation(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[DeclineInvitationRequest](c)

    err := ctrl.diaryService.DeclineInvitation(req.DiaryID, loginUser.UserId)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "diary.invitation_declined"}, "refreshInvitations": true}`)
        c.Status(http.StatusOK)
        return
    }

    ctrl.renderService.JSONSuccess(c, nil)
}

// QuitDiary 退出日记本
func (ctrl *DiaryController) QuitDiary(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[QuitDiaryRequest](c)

    err := ctrl.diaryService.QuitDiary(req.DiaryID, loginUser.UserId)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "diary.quit_successfully"}, "refreshJoinedDiaries": true}`)
        c.Status(http.StatusOK)
        return
    }

    ctrl.renderService.JSONSuccess(c, nil)
}

// RemoveAuthor 移除作者
func (ctrl *DiaryController) RemoveAuthor(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)
    req := middleware.MustGetValidatedRequest[RemoveAuthorRequest](c)

    err := ctrl.diaryService.RemoveAuthor(req.DiaryID, loginUser.UserId, req.UserID)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        c.Header("HX-Trigger", `{"showToast": {"type": "success", "message": "diary.author_removed"}}`)
        ctrl.ListAuthors(c)
        return
    }

    ctrl.renderService.JSONSuccess(c, nil)
}

// ListJoinedDiaries 获取用户参与的日记本列表
func (ctrl *DiaryController) ListJoinedDiaries(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)

    pageIndex := 0
    pageSize := 20
    if pageStr := c.Query("page"); pageStr != "" {
        if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
            pageIndex = page - 1
        }
    }

    diaries, total, err := ctrl.diaryService.ListJoinedDiaries(loginUser.UserId, pageIndex, pageSize)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        data := map[string]interface{}{
            "Diaries":    diaries,
            "Total":      total,
            "PageIndex":  pageIndex + 1,
            "PageSize":   pageSize,
            "TotalPages": (total + pageSize - 1) / pageSize,
            "Lang":       loginUser.Lang,
            "TimeZone":   loginUser.TimeZone,
            "UserID":     loginUser.UserId,
        }
        c.HTML(http.StatusOK, "user/joined_diaries_list", data)
        return
    }

    ctrl.renderService.JSONSuccess(c, map[string]interface{}{
        "diaries":  diaries,
        "total":    total,
        "page":     pageIndex + 1,
        "per_page": pageSize,
    })
}

// ListPendingInvitations 获取用户的待确认邀请列表
func (ctrl *DiaryController) ListPendingInvitations(c *gin.Context) {
    loginUser := ctrl.auth.MustGetCurrentUser(c)

    invitations, err := ctrl.diaryService.ListPendingInvitations(loginUser.UserId)
    if err != nil {
        ctrl.renderService.JSONError(c, err)
        return
    }

    if c.GetHeader("HX-Request") == "true" {
        data := map[string]interface{}{
            "Invitations": invitations,
            "Lang":        loginUser.Lang,
            "TimeZone":    loginUser.TimeZone,
            "UserID":      loginUser.UserId,
        }
        c.HTML(http.StatusOK, "user/invitations_list", data)
        return
    }

    ctrl.renderService.JSONSuccess(c, map[string]interface{}{
        "invitations": invitations,
    })
}
```

### 4.3 路由定义 (`routes.go`)

```go
// RegisterRoutes 注册日记本相关路由
func (ctrl *DiaryController) RegisterRoutes(router *gin.Engine) {
    // ... 原有路由 ...

    // ===================== 作者管理路由 =====================

    // 作者列表
    authAPI.GET("/diary/:id/authors", middleware.Validate[ListAuthorsRequest](), ctrl.ListAuthors)
    authAPI.POST("/diary/:id/authors", middleware.Validate[InviteAuthorRequest](), ctrl.InviteAuthor)
    authAPI.PUT("/diary/:id/authors/:user_id/admin", middleware.Validate[SetAuthorAdminRequest](), ctrl.SetAuthorAdmin)
    authAPI.DELETE("/diary/:id/authors/:user_id", middleware.Validate[RemoveAuthorRequest](), ctrl.RemoveAuthor)
    authAPI.POST("/diary/:id/authors/quit", middleware.Validate[QuitDiaryRequest](), ctrl.QuitDiary)

    // 邀请操作
    authAPI.POST("/diary/:id/join", middleware.Validate[AcceptInvitationRequest](), ctrl.AcceptInvitation)
    authAPI.POST("/diary/:id/decline", middleware.Validate[DeclineInvitationRequest](), ctrl.DeclineInvitation)

    // ===================== 用户主页相关路由 =====================

    // 用户个人主页（公开访问）
    router.GET("/u/:uid", ctrl.ShowUserProfile)
    router.GET("/u/:uid/:pageindex", ctrl.ShowUserProfile)

    // 加入的日记本列表
    homeAPI := router.Group("/home").Use(ctrl.auth.RequireAuth())
    {
        homeAPI.GET("/joined-diaries", ctrl.ListJoinedDiaries)
        homeAPI.GET("/invitations", ctrl.ListPendingInvitations)
    }
}
```

## 5. 前端设计 (Frontend Design)

### 5.1 安全机制

#### 5.1.1 项目统一安全方案

hebe 项目使用统一的安全机制，包括 CSRF 保护、表单验证和请求拦截。

**CSRF Token 机制**（来自 `website/static/src/js/core/init.js`）：

```javascript
// CSRF Token 获取
window.getCsrfToken = () => {
  return document
    .querySelector('meta[name="csrf-token"]')
    ?.getAttribute('content');
};

// Fetch 拦截器自动注入 CSRF Token
const originalFetch = window.fetch;
window.fetch = function (input, init) {
  const csrfToken = window.getCsrfToken();
  if (csrfToken) {
    init = init || {};
    init.headers = init.headers || {};
    if (init.headers instanceof Headers) {
      init.headers.append('X-CSRF-Token', csrfToken);
    } else {
      init.headers['X-CSRF-Token'] = csrfToken;
    }
  }
  // ... 进度条处理 ...
  return originalFetch(input, init).finally(() => {
    // ... 进度条结束 ...
  });
};

// HTMX 自动注入 CSRF Token
document.body.addEventListener('htmx:configRequest', (evt) => {
  const csrfToken = window.getCsrfToken();
  if (csrfToken) {
    evt.detail.headers['X-CSRF-Token'] = csrfToken;
  }
});
```

**表单验证机制**（使用 `x-uno-valid` 指令）：

```html
<!-- 示例：用户设置表单 -->
<form
  id="profile-form"
  class="relative"
  x-uno-valid
  hx-post="/u/settings/profile"
  hx-swap="none"
>
  <input
    type="text"
    id="display_name"
    name="display_name"
    x-uno-valid="[{ type: 'length', max: 30 }]"
  />
</form>
```

**后端验证中间件**（来自 `internal/app/middleware/validator.go`）：

```go
// 泛型验证中间件
func Validate[T any]() gin.HandlerFunc {
  return func(c *gin.Context) {
    var req T
    _ = c.ShouldBindUri(&req)      // 绑定 URI 参数
    if err := c.ShouldBind(&req); err != nil {  // 绑定 Body/Query 参数
      details := formatValidationError(err)
      errResponse := errors.NewAPIErrorResponse(
        errors.InvalidParams.Code(),
        errors.InvalidParams.Key(),
        details,
        getRequestID(c),
      )
      c.AbortWithStatusJSON(http.StatusBadRequest, errResponse)
      return
    }
    c.Set("validatedRequest", &req)
    c.Next()
  }
}

// 控制器中使用
func (ctrl *DiaryController) InviteAuthor(c *gin.Context) {
  loginUser := ctrl.auth.MustGetCurrentUser(c)
  req := middleware.MustGetValidatedRequest[InviteAuthorRequest](c)
  // ... 业务逻辑 ...
}
```

#### 5.1.2 前端实现（使用项目统一方案）

**JavaScript 模块**（`website/static/src/js/modules/authors.js`）：

```javascript
// 作者管理模块
export default function Authors() {
  return {
    // 打开邀请作者弹窗（使用 uno box）
    openInviteAuthorBox(diaryId) {
      const uno = window.getUno();
      if (!uno) return;

      let closeHandler;
      const boxId =
        'box-invite-author-' + Math.random().toString(36).slice(2, 10);

      uno.box.dialog({
        boxId: boxId,
        title: this.$store.i18n.t('diary.authors.invite_title'),
        widthClass: '!w-[480px]',
        showOkButton: true,
        showCancelButton: true,
        okText: this.$store.i18n.t('diary.authors.invite_send'),
        cancelText: this.$store.i18n.t('common.cancel'),
        containerClass: '!border',
        okButtonClass: 'btn',
        cancelButtonClass: 'btn',
        buttonAlign: 'right',
        onContentLoad: (el) => {
          if (!el) el = document.getElementById(boxId);
          if (window.htmx) window.htmx.process(el);
          if (window.Alpine) window.Alpine.initTree(el);
        },
        onClose: () => {
          if (closeHandler) {
            document.body.removeEventListener(
              'closeInviteAuthorBox',
              closeHandler,
              true,
            );
          }
        },
      });

      closeHandler = (evt) => {
        const eventBoxId = evt?.target?.dataset?.boxId;
        if (eventBoxId && eventBoxId !== boxId) {
          return;
        }
        uno.box.close(boxId);
      };

      document.body.addEventListener(
        'closeInviteAuthorBox',
        closeHandler,
        true,
      );
    },

    // 刷新作者列表
    refreshAuthorsList(diaryId) {
      htmx.ajax('GET', `/diary/${diaryId}/authors`, {
        target: '#authorsList',
        swap: 'innerHTML',
      });
    },

    // 切换管理员权限
    toggleAdmin(diaryId, userId, isAdmin) {
      const uno = window.getUno();
      const newIsAdmin = !isAdmin;
      const confirmMsg = newIsAdmin
        ? this.$store.i18n.t('diary.authors.confirm_grant_admin')
        : this.$store.i18n.t('diary.authors.confirm_revoke_admin');

      uno.box
        .confirm(confirmMsg, {
          title: this.$store.i18n.t('diary.authors.admin_toggle_title'),
          okText: this.$store.i18n.t('common.confirm'),
          cancelText: this.$store.i18n.t('common.cancel'),
          okButtonClass: 'btn',
          cancelButtonClass: 'btn',
        })
        .then((confirmed) => {
          if (confirmed) {
            // CSRF Token 自动通过 fetch 拦截器注入
            fetch(`/diary/${diaryId}/authors/${userId}/admin`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ is_admin: newIsAdmin }),
            }).then(() => {
              this.refreshAuthorsList(diaryId);
            });
          }
        });
    },

    // 移除作者
    removeAuthor(diaryId, userId) {
      const uno = window.getUno();

      uno.box
        .confirm(this.$store.i18n.t('diary.authors.confirm_remove'), {
          title: this.$store.i18n.t('diary.authors.remove_title'),
          okText: this.$store.i18n.t('common.confirm'),
          cancelText: this.$store.i18n.t('common.cancel'),
          okButtonClass: 'btn btn-danger',
          cancelButtonClass: 'btn',
        })
        .then((confirmed) => {
          if (confirmed) {
            // CSRF Token 自动通过 fetch 拦截器注入
            fetch(`/diary/${diaryId}/authors/${userId}`, {
              method: 'DELETE',
            }).then(() => {
              this.refreshAuthorsList(diaryId);
            });
          }
        });
    },

    // 退出日记本
    quitDiary(diaryId) {
      const uno = window.getUno();

      uno.box
        .confirm(this.$store.i18n.t('diary.authors.confirm_quit'), {
          title: this.$store.i18n.t('diary.authors.quit_title'),
          okText: this.$store.i18n.t('common.confirm'),
          cancelText: this.$store.i18n.t('common.cancel'),
          okButtonClass: 'btn',
          cancelButtonClass: 'btn',
        })
        .then((confirmed) => {
          if (confirmed) {
            // CSRF Token 自动通过 fetch 拦截器注入
            fetch(`/diary/${diaryId}/authors/quit`, {
              method: 'POST',
            }).then(() => {
              window.location.href = '/home';
            });
          }
        });
    },
  };
}
```

### 5.2 邀请作者表单 (`website/views/diary/box_invite_author.html`)

```html
{{define "diary/box_invite_author"}}
<form
  id="inviteAuthorForm"
  class="relative"
  x-uno-valid
  hx-post="/diary/{{.DiaryID}}/authors"
  hx-target="#authorsList"
  hx-swap="outerHTML"
  @htmx:after-request="uno?.box?.close(document.querySelector('[data-box-id]')?.dataset.boxId)"
>
  <div class="form-mask">
    <div class="loading-spinner"></div>
  </div>

  <div>
    <div class="mb-4">
      <label class="form-title" for="invitee_identity">
        {{getMessage "diary.authors.invite_label" .Lang}}
      </label>
      <input
        type="text"
        id="invitee_identity"
        name="invitee_identity"
        class="form-control"
        placeholder='{{getMessage "diary.authors.invite_placeholder" .Lang}}'
        x-uno-valid="['required']"
        autofocus
      />
      <p class="form-info">{{getMessage "diary.authors.invite_help" .Lang}}</p>
    </div>
  </div>
</form>
{{end}}
```

### 5.3 作者列表组件 (`website/views/diary/authors_list.html`)

```html
{{define "diary/authors_list"}}
<div class="authors-container mb-6">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-xl font-semibold">
      <i class="mdi mdi-account-group mr-2"></i>
      {{getMessage "diary.label_info.authors_list" .Lang}}
    </h2>

    {{if .Data.IsOwner}}
    <button
      class="btn"
      hx-get="/form/diary?action=diary.invite_author&diary_id={{.Data.Diary.DiaryID}}"
      hx-target="#invited-author-box"
      hx-swap="none"
    >
      <i class="mdi mdi-plus mr-1"></i>
      {{getMessage "diary.authors.invite" .Lang}}
    </button>
    {{end}}
  </div>

  <!-- Box 容器 -->
  <div id="invited-author-box" class="hidden"></div>

  {{if not .Data.Authors}}
  <div class="text-center py-8 text-gray-500">
    <i class="mdi mdi-account-off text-4xl mb-2"></i>
    <p>{{getMessage "diary.authors.empty" .Lang}}</p>
  </div>
  {{else}}
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {{range .Data.Authors}}
    <!-- 拥有者卡片 -->
    {{if eq .UserID .Data.Diary.AddUserID}}
    <div class="card author-card" data-user-id="{{.UserID}}">
      <div class="relative">
        <!-- 右上角拥有者标识 -->
        <div class="absolute top-2 right-2">
          <span
            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
          >
            <i class="mdi mdi-crown mr-1"></i>
            {{getMessage "diary.authors.owner" $.Lang}}
          </span>
        </div>

        <div class="text-center py-6">
          <div
            class="avatar-lg inline-flex items-center justify-center relative mx-auto"
          >
            {{getFirstCharOfName .}}
          </div>
          <h3 class="font-medium mt-3">{{or .DisplayName .Handle}}</h3>
          <p class="text-sm text-gray-500">
            {{.PostNum}} {{getMessage "diary.authors.posts" $.Lang}}
          </p>
        </div>
      </div>
    </div>
    {{else}}
    <!-- 普通作者/管理员卡片 -->
    <div class="card author-card" data-user-id="{{.UserID}}">
      <div class="relative">
        <!-- 右上角功能菜单按钮 -->
        <div class="absolute top-2 right-2">
          <a
            href="javascript:void(0)"
            x-uno-tooltip="{ content: '#author-menu-{{.UserID}}', trigger: 'click', placement: 'bottom' }"
          >
            <i class="mdi mdi-dots-horizontal"></i>
          </a>
        </div>

        <div class="text-center py-6">
          <div
            class="avatar-lg inline-flex items-center justify-center relative mx-auto"
          >
            {{getFirstCharOfName .}}
          </div>
          <h3 class="font-medium mt-3">
            {{or .DisplayName .Handle}} {{if eq .IsAdmin 1}}
            <span
              class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 ml-2"
            >
              <i class="mdi mdi-shield-account mr-1"></i>
              {{getMessage "diary.authors.admin" $.Lang}}
            </span>
            {{end}}
          </h3>
          <p class="text-sm text-gray-500">
            {{.PostNum}} {{getMessage "diary.authors.posts" $.Lang}}
          </p>
        </div>
      </div>
    </div>
    {{end}} {{end}}
  </div>

  <!-- 分页组件（只有多页时显示） -->
  {{if gt .Data.TotalPages 1}}
  <div class="flex justify-center mt-6">{{template "pagination" .}}</div>
  {{end}} {{end}}
</div>

<!-- 功能菜单模板（隐藏） -->
{{range .Data.Authors}} {{if ne .UserID .Data.Diary.AddUserID}}
<div id="author-menu-{{.UserID}}" class="hidden">
  <ul class="text-left">
    {{if .Data.IsOwner}}
    <li class="py-1">
      <button
        class="block w-full text-left hover:text-blue-500"
        @click="authorsMixin.toggleAdmin({{$.Data.Diary.DiaryID}}, {{.UserID}}, {{eq .IsAdmin 1}}); $store.uno?.tooltip?.hide()"
      >
        <i
          class="mdi mdi-shield{{if eq .IsAdmin 1}}{{else}}-outline{{end}} mr-2"
        ></i>
        {{if eq .IsAdmin 1}}{{getMessage "diary.authors.revoke_admin"
        $.Lang}}{{else}}{{getMessage "diary.authors.grant_admin" $.Lang}}{{end}}
      </button>
    </li>
    <li class="py-1">
      <button
        class="block w-full text-left hover:text-red-500"
        @click="authorsMixin.removeAuthor({{$.Data.Diary.DiaryID}}, {{.UserID}}); $store.uno?.tooltip?.hide()"
      >
        <i class="mdi mdi-delete mr-2"></i>
        {{getMessage "diary.authors.remove" $.Lang}}
      </button>
    </li>
    {{else if eq .UserID .Data.UserID}}
    <li class="py-1">
      <button
        class="block w-full text-left hover:text-red-500"
        @click="authorsMixin.quitDiary({{$.Data.Diary.DiaryID}}); $store.uno?.tooltip?.hide()"
      >
        <i class="mdi mdi-logout mr-2"></i>
        {{getMessage "diary.authors.quit" $.Lang}}
      </button>
    </li>
    {{end}}
  </ul>
</div>
{{end}} {{end}} {{end}}
```

### 5.4 日记本首页更新 (`website/views/diary/index.html`)

```html
{{define "diary/index"}} {{template "header" .}} {{template "navbar" .}}

<main
  class="main dt{{.Data.Diary.ThemeID}} transition-colors duration-200"
  data-social-setting="{{.Data.Diary.SocialSetting}}"
  data-is-owner="{{.Data.IsOwner}}"
  data-diary-id="{{.Data.Diary.DiaryID}}"
  x-data="diaryDetail"
>
  <section class="flex justify-between mb-10">
    <div id="page-nav">
      <ul class="flex text-lg">
        <li>
          <a
            href="/u/{{.Data.Diary.AddUserID}}"
            class="hover:text-red-500 dark:hover:text-red-400"
            >{{ or .Data.Diary.AddUserDisplayName .Data.Diary.AddUserHandle
            }}</a
          >
          <i class="mdi mdi-slash-forward mx-3"></i>
        </li>
        <li>
          <a
            href="/diary/{{.Data.Diary.DiaryID}}"
            id="diary_title"
            class="hover:text-red-500 dark:hover:text-red-400"
            >{{.Data.Diary.Title}}{{if eq .Data.Diary.ViewSetting 2}}<i
              class="mdi mdi-lock-outline mx-3"
            ></i
            >{{end}}
          </a>
        </li>
      </ul>
    </div>
    <div>
      {{if .Data.IsOwner}}
      <button
        class="mr-3 btn"
        @click="openDiarySettings({{.Data.Diary.DiaryID}})"
      >
        <i class="mdi mdi-cogs"></i> {{getMessage "diary.settings"
        .Context.Lang}}
      </button>
      {{end}}
      <!-- 日记本点赞按钮容器 -->
      {{if eq .Data.Diary.SocialSetting 1}}
      <span id="diary-like-wrapper" data-diary-id="{{.Data.Diary.DiaryID}}">
        <button class="btn" @click="toggleDiaryLike({{.Data.Diary.DiaryID}})">
          <i
            class="mdi mr-1"
            :class="diaryLike.liked ? 'mdi-thumb-up' : 'mdi-thumb-up-outline'"
          ></i>
          <span x-text="diaryLike.count"></span>
        </button>
      </span>
      {{end}}
    </div>
  </section>

  <section class="mb-10">
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <div>
        <h1 class="card-title">
          {{getMessage "diary.label_info.introduction" .Context.Lang}}
        </h1>
        <div class="text-sm">
          <p id="diary_remark">{{.Data.Diary.Remark}}</p>
        </div>
      </div>
      <div>
        <h1 class="card-title">{{getMessage "common.menu" .Context.Lang}}</h1>
        <div>
          <ul>
            <li class="py-1">
              <a
                href="/diary/{{.Data.Diary.DiaryID}}/authors"
                class="hover:text-red-500 dark:hover:text-red-400"
              >
                <i class="mdi mdi-account-group"></i>
                <span
                  >{{getMessage "diary.label_info.authors_list"
                  .Context.Lang}}</span
                >
              </a>
            </li>
            <li class="py-1">
              <a href="#" class="hover:text-red-500 dark:hover:text-red-400">
                <i class="mdi mdi-delete-outline"></i>
                <span
                  >{{getMessage "diary.label_info.liked_users_list"
                  .Context.Lang}}</span
                >
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </section>

  <section>{{ template "diary/list" . }}</section>
</main>

{{template "footer" .}} {{end}}
```

### 5.5 作者列表页面 (`website/views/diary/authors.html`)

```html
{{define "diary/authors"}} {{template "header" .}} {{template "navbar" .}}

<main class="main dt{{.Data.Diary.ThemeID}} transition-colors duration-200">
  <section class="mb-10">
    <div id="page-nav">
      <ul class="flex text-lg">
        <li>
          <a
            href="/u/{{.Data.Diary.AddUserID}}"
            class="hover:text-red-500 dark:hover:text-red-400"
            >{{ or .Data.Diary.AddUserDisplayName .Data.Diary.AddUserHandle
            }}</a
          >
          <i class="mdi mdi-slash-forward mx-3"></i>
        </li>
        <li>
          <a
            href="/diary/{{.Data.Diary.DiaryID}}"
            id="diary_title"
            class="hover:text-red-500 dark:hover:text-red-400"
            >{{.Data.Diary.Title}}</a
          >
          <i class="mdi mdi-slash-forward mx-3"></i>
        </li>
        <li>
          <a href="/diary/{{.Data.Diary.DiaryID}}/authors">
            {{getMessage "diary.label_info.authors_list" .Lang}}
          </a>
        </li>
      </ul>
    </div>
  </section>

  <div id="authorsList">{{template "diary/authors_list" .}}</div>
</main>

{{template "footer" .}} {{end}}
```

### 5.6 用户主页更新 (`website/views/user/index.html`)

```html
{{ define "user/index" }} {{ template "header" . }} {{ template "navbar" . }}

<main class="main" x-data="user">
  <section class="mb-10 flex justify-between">
    <div id="page-nav">
      <ul class="flex text-lg">
        <li>
          <a href="/u/{{ .Data.User.UserId }}"
            >{{ or .Data.User.DisplayName .Data.User.Handle }}</a
          >
        </li>
      </ul>
    </div>
  </section>

  <section class="mb-10">
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <div>
        <h1 class="card-title">
          {{getMessage "user.profile.introduction" .Context.Lang}}
        </h1>
        <div>{{ .Data.User.Remark }}</div>
      </div>
      <div>
        <h1 class="card-title">
          {{getMessage "user.profile.menu" .Context.Lang}}
        </h1>
        <div>
          <ul class="text-gray-900 dark:text-gray-100">
            <li class="py-1">
              <a href="#">
                <i class="mdi mdi-book-outline"></i>
                <span
                  >{{getMessage "user.profile.latest_public_diary"
                  .Context.Lang}}</span
                >
              </a>
            </li>
            <li class="py-1">
              <a href="#joined-diaries">
                <i class="mdi mdi-book-multiple"></i>
                <span
                  >{{getMessage "user.profile.joined_diaries"
                  .Context.Lang}}</span
                >
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </section>

  <section class="flex justify-between mb-10">
    <div class="text-lg">{{getMessage "diary.list.title" .Context.Lang}}</div>
    <div>
      {{if .Data.IsOwner}}
      <button class="btn" @click="openCreateDiary()">
        <i class="mdi mdi-plus"></i> {{getMessage "diary.create" .Context.Lang}}
      </button>
      {{end}}
    </div>
  </section>

  <section>{{ template "user/list" . }}</section>

  <!-- 加入的日记本区域（包含待确认邀请） -->
  <section class="mt-10" id="joined-diaries">
    <div id="joinedDiariesList">{{template "user/joined_diaries_list" .}}</div>
  </section>
</main>
{{ template "footer" . }} {{ end }}
```

### 5.7 已加入日记本列表组件（融合邀请和已加入）(`website/views/user/joined_diaries_list.html`)

```html
{{define "user/joined_diaries_list"}}
{{if .Data.Diaries}}
<div class="joined-diaries-section">
    <h2 class="text-xl font-semibold mb-4">
        <i class="mdi mdi-book-multiple mr-2"></i>
        {{getMessage "user.joined_diaries.title" .Lang}}
    </h2>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {{range .Data.Diaries}}
        <div class="diary-card p-4 bg-white dark:bg-gray-800 rounded-lg shadow-sm">
            {{template "user/diary_card" dict "Diary" . "UserID" $.Data.UserID "Lang" $.Data.Lang "TimeZone" $.Data.TimeZone}}

            <div class="flex space-x-2 mt-4">
                <a href="/diary/{{.DiaryID}}" class="btn btn-sm flex-1">
                    {{getMessage "diary.view" .Lang}}
                </a>
                {{if eq .Status 2}}
                <!-- 待确认邀请 -->
                <div class="flex space-x-2 flex-1">
                    <button
                        class="btn btn-sm flex-1"
                        hx-post="/diary/{{.DiaryID}}/join"
                        hx-target="#joinedDiariesList"
                        hx-swap="outerHTML">
                        {{getMessage "common.accept" .Lang}}
                    </button>
                    <button
                        class="btn btn-sm btn-secondary flex-1"
                        hx-post="/diary/{{.DiaryID}}/decline"
                        hx-target="#joinedDiariesList"
                        hx-swap="outerHTML">
                        {{getMessage "common.decline" .Lang}}
                    </button>
                </div>
                {{else}}
                <!-- 已加入的日记本 -->
                <button
                    class="btn btn-sm btn-secondary flex-1"
                    hx-post="/diary/{{.DiaryID}}/authors/quit"
                    hx-vals='{"diary_id": {{.DiaryID}}}'
                    hx-confirm="{{getMessage "diary.authors.confirm_quit" .Lang}}"
                    hx-target="#joinedDiariesList"
                    hx-swap="outerHTML">
                    {{getMessage "diary.authors.quit" .Lang}}
                </button>
                {{end}}
            </div>
        </div>
        {{end}}
    </div>

    <!-- 更多链接 -->
    {{if .Data.ShowMoreLink}}
    <div class="text-center mt-4">
        <a href="/home/joined-diaries" class="text-blue-500 hover:text-blue-600">
            {{getMessage "user.joined_diaries.view_more" .Lang}}
        </a>
    </div>
    {{end}}
</div>
{{else}}
<div class="text-center py-8 text-gray-500">
    <i class="mdi mdi-book-open-variant-outline text-4xl mb-2"></i>
    <p>{{getMessage "user.joined_diaries.empty" .Lang}}</p>
</div>
{{end}}
{{end}}
```

### 5.8 已加入日记本列表完整页面（分页）(`website/views/user/joined_diaries_full.html`)

```html
{{define "user/joined_diaries_full"}} {{template "header" .}} {{template
"navbar" .}}

<main class="main">
  <section class="mb-10">
    <div id="page-nav">
      <ul class="flex text-lg">
        <li>
          <a href="/u/{{ .Data.User.UserId }}"
            >{{ or .Data.User.DisplayName .Data.User.Handle }}</a
          >
        </li>
        <li class="mx-2">
          <span class="mdi mdi-slash-forward"></span>
        </li>
        <li>
          <a href="#"
            >{{getMessage "user.joined_diaries.title" .Context.Lang}}</a
          >
        </li>
      </ul>
    </div>
  </section>

  <div id="joinedDiariesList">{{template "user/joined_diaries_list" .}}</div>

  <!-- 分页组件（只有多页时显示） -->
  {{if gt .Data.TotalPages 1}}
  <div class="flex justify-center mt-6">{{template "pagination" .}}</div>
  {{end}}
</main>

{{template "footer" .}} {{end}}
```

### 5.9 环境配置项

需要在 `.env.local` 中增加配置项：

```bash
# 用户主页显示的加入日记本数量（包含待确认邀请）
IDIARY_USER_JOINED_DIARIES_LIMIT=5
```

### 5.10 服务端逻辑更新

#### 5.9.1 Service 层更新

```go
// ListJoinedDiariesWithInvitations 获取用户的加入日记本列表（包含待确认邀请）
// 排序规则：status desc（待确认优先），join_utc_time desc（最新优先）
// limit: 限制返回数量（用于用户主页展示），<= 0 表示返回全部（用于完整列表页）
func (s *Service) ListJoinedDiariesWithInvitations(userID int64, limit int) ([]Diary, bool, error) {
    diaries, err := s.repo.ListJoinedDiariesWithInvitations(userID, limit)
    if err != nil {
        return nil, false, err
    }

    // 计算是否有更多（用于显示"更多"链接）
    total, err := s.repo.GetJoinedDiariesCountByUserID(userID)
    if err != nil {
        return nil, false, err
    }

    showMore := (limit > 0 && total > limit)

    return diaries, showMore, nil
}

// ListJoinedDiariesWithInvitationsPaged 获取用户的加入日记本分页列表
func (s *Service) ListJoinedDiariesWithInvitationsPaged(userID int64, pageIndex, pageSize int) ([]Diary, int, int, error) {
    diaries, total, err := s.repo.ListJoinedDiariesWithInvitationsPaged(userID, pageIndex, pageSize)
    if err != nil {
        return nil, 0, 0, err
    }

    totalPages := (total + pageSize - 1) / pageSize

    return diaries, total, totalPages, nil
}
```

#### 5.9.2 Repository 层更新

```go
// ListJoinedDiariesWithInvitations 获取用户的加入日记本列表（包含待确认邀请）
// 排序规则：status desc（待确认优先），join_utc_time desc（最新优先）
func (r *Repository) ListJoinedDiariesWithInvitations(userID int64, limit int) ([]Diary, error) {
    query := `
        SELECT d.*, da.status, da.join_utc_time
        FROM br_diary d
        INNER JOIN br_diary_author da ON d.diary_id = da.diary_id
        WHERE da.user_id = ? AND da.status IN (1, 2)
        ORDER BY da.status DESC, da.join_utc_time DESC
    `

    if limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", limit)
    }

    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var diaries []Diary
    for rows.Next() {
        var d Diary
        // ... 扫描字段 ...
        diaries = append(diaries, d)
    }

    return diaries, nil
}

// ListJoinedDiariesWithInvitationsPaged 获取用户的加入日记本分页列表
func (r *Repository) ListJoinedDiariesWithInvitationsPaged(userID int64, pageIndex, pageSize int) ([]Diary, int, error) {
    offset := pageIndex * pageSize

    // 查询分页数据
    query := `
        SELECT d.*, da.status, da.join_utc_time
        FROM br_diary d
        INNER JOIN br_diary_author da ON d.diary_id = da.diary_id
        WHERE da.user_id = ? AND da.status IN (1, 2)
        ORDER BY da.status DESC, da.join_utc_time DESC
        LIMIT ? OFFSET ?
    `

    rows, err := r.db.Query(query, userID, pageSize, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var diaries []Diary
    for rows.Next() {
        var d Diary
        // ... 扫描字段 ...
        diaries = append(diaries, d)
    }

    // 查询总数
    countQuery := `
        SELECT COUNT(*)
        FROM br_diary d
        INNER JOIN br_diary_author da ON d.diary_id = da.diary_id
        WHERE da.user_id = ? AND da.status IN (1, 2)
    `

    var total int
    err = r.db.QueryRow(countQuery, userID).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    return diaries, total, nil
}

// GetJoinedDiariesCountByUserID 获取用户参与的日记本数量
func (r *Repository) GetJoinedDiariesCountByUserID(userID int64) (int, error) {
    query := `
        SELECT COUNT(*)
        FROM br_diary d
        INNER JOIN br_diary_author da ON d.diary_id = da.diary_id
        WHERE da.user_id = ? AND da.status IN (1, 2)
    `

    var count int
    err := r.db.QueryRow(query, userID).Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}
```

### 5.11 多语言文件更新

**原则说明：**

- **模板多语言**：用于模型显示、静态文本（如按钮标签、表单标签、标题等）
- **JavaScript 多语言**：用于前端交互（如确认对话框、操作提示、错误消息等）

#### 模板多语言 - `website/locales/zh-cn.json`

```json
{
  "diary": {
    "label_info": {
      "authors_list": "作者列表"
    },
    "authors": {
      "invite": "邀请作者",
      "invite_title": "邀请新作者",
      "invite_label": "用户标识",
      "invite_placeholder": "@handle 或 email@example.com",
      "invite_help": "输入用户的 @handle 或邮箱地址",
      "invite_send": "发送邀请",
      "owner": "拥有者",
      "admin": "管理员",
      "posts": "篇日记",
      "empty": "暂无作者"
    },
    "view": "查看"
  },
  "user": {
    "profile": {
      "joined_diaries": "加入的日记本"
    },
    "joined_diaries": {
      "title": "加入的日记本",
      "view_more": "查看更多",
      "empty": "您还没有加入任何日记本"
    }
  },
  "common": {
    "accept": "接受",
    "decline": "拒绝"
  },
  "email": {
    "diary_invitation": {
      "subject": "邀请您加入日记本「{{diary_title}}」",
      "body": "<!DOCTYPE html>\n<html lang=\"zh-CN\">\n<head>\n<meta charset=\"UTF-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n<style>body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;line-height:1.6;color:#333;max-width:600px;margin:0 auto;padding:20px;background:#f5f5f5}.email-container{background:#fff;border-radius:8px;padding:30px;box-shadow:0 2px 8px rgba(0,0,0,0.1)}.email-header{border-bottom:1px solid #eee;padding-bottom:20px;margin-bottom:20px}.email-header h1{margin:0;color:#e53e3e;font-size:24px}.email-body{margin-bottom:30px}.email-body p{margin:10px 0}.email-footer{border-top:1px solid #eee;padding-top:20px;margin-top:20px;color:#999;font-size:12px}.btn{display:inline-block;padding:10px 20px;background:#e53e3e;color:#fff;text-decoration:none;border-radius:4px;margin:10px 0}.btn:hover{background:#c53030}</style>\n</head>\n<body>\n<div class=\"email-container\">\n  <div class=\"email-header\">\n    <h1>{{diary_title}}</h1>\n  </div>\n  <div class=\"email-body\">\n    <p>亲爱的 <strong>{{invitee_name}}</strong>：</p>\n    <p><strong>{{owner_name}}</strong> (@{{owner_handle}}) 邀请您加入日记本「{{diary_title}}」。</p>\n    <p>这是一个多人协作的日记本，加入后您可以：</p>\n    <ul>\n      <li>与作者们一起记录生活点滴</li>\n      <li>查看和分享日记内容</li>\n      <li>参与日记本的建设和管理</li>\n    </ul>\n    <p style=\"text-align:center;\">\n      <a href=\"{{accept_url}}\" class=\"btn\">接受邀请</a>\n    </p>\n    <p>如果您不认识 {{owner_name}}，或者不想加入此日记本，请忽略此邮件。</p>\n  </div>\n  <div class=\"email-footer\">\n    <p>此邮件由 {{domain}} 自动发送，请勿直接回复。</p>\n  </div>\n</div>\n</body>\n</html>"
    }
  }
}
```

#### 模板多语言 - `website/locales/en.json`

```json
{
  "diary": {
    "label_info": {
      "authors_list": "Authors List"
    },
    "authors": {
      "invite": "Invite Author",
      "invite_title": "Invite New Author",
      "invite_label": "User Identity",
      "invite_placeholder": "@handle or email@example.com",
      "invite_help": "Enter user's @handle or email address",
      "invite_send": "Send Invitation",
      "owner": "Owner",
      "admin": "Admin",
      "posts": "posts",
      "empty": "No authors yet"
    },
    "view": "View"
  },
  "user": {
    "profile": {
      "joined_diaries": "Joined Diaries"
    },
    "joined_diaries": {
      "title": "Joined Diaries",
      "view_more": "View More",
      "empty": "You haven't joined any diaries yet"
    }
  },
  "common": {
    "accept": "Accept",
    "decline": "Decline"
  },
  "email": {
    "diary_invitation": {
      "subject": "You're invited to join \"{{diary_title}}\"",
      "body": "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"UTF-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n<style>body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;line-height:1.6;color:#333;max-width:600px;margin:0 auto;padding:20px;background:#f5f5f5}.email-container{background:#fff;border-radius:8px;padding:30px;box-shadow:0 2px 8px rgba(0,0,0,0.1)}.email-header{border-bottom:1px solid #eee;padding-bottom:20px;margin-bottom:20px}.email-header h1{margin:0;color:#e53e3e;font-size:24px}.email-body{margin-bottom:30px}.email-body p{margin:10px 0}.email-footer{border-top:1px solid #eee;padding-top:20px;margin-top:20px;color:#999;font-size:12px}.btn{display:inline-block;padding:10px 20px;background:#e53e3e;color:#fff;text-decoration:none;border-radius:4px;margin:10px 0}.btn:hover{background:#c53030}</style>\n</head>\n<body>\n<div class=\"email-container\">\n  <div class=\"email-header\">\n    <h1>{{diary_title}}</h1>\n  </div>\n  <div class=\"email-body\">\n    <p>Dear <strong>{{invitee_name}}</strong>,</p>\n    <p><strong>{{owner_name}}</strong> (@{{owner_handle}}) has invited you to join the diary \"{{diary_title}}\".</p>\n    <p>This is a collaborative diary where you can:</p>\n    <ul>\n      <li>Record life moments together with other authors</li>\n      <li>View and share diary entries</li>\n      <li>Participate in diary management</li>\n    </ul>\n    <p style=\"text-align:center;\">\n      <a href=\"{{accept_url}}\" class=\"btn\">Accept Invitation</a>\n    </p>\n    <p>If you don't know {{owner_name}} or don't want to join this diary, please ignore this email.</p>\n  </div>\n  <div class=\"email-footer\">\n    <p>This email was automatically sent by {{domain}}. Please do not reply directly.</p>\n  </div>\n</div>\n</body>\n</html>"
    }
  }
}
```

#### JavaScript 多语言 - `website/static/src/js/locales/zh-cn.js`

```javascript
export default {
  // Common
  'common.save': '保存',
  'common.edit': '编辑',
  'common.cancel': '取消',
  'common.delete': '删除',
  'common.op_fail': '操作失败',
  'common.op_success': '操作成功',
  'common.load_fail': '加载失败',
  'common.save_success': '保存成功',
  'common.save_fail': '保存失败',
  'common.net_error': '网络错误',
  'common.views': '浏览',
  'common.comments': '评论',
  'common.confirm': '确认',

  // Diary
  'diary.create': '新建日记本',
  'diary.settings': '日记本设置',

  // Diary Authors (前端交互相关)
  'diary.authors.invite_title': '邀请新作者',
  'diary.authors.invite_send': '发送邀请',
  'diary.authors.grant_admin': '设为管理员',
  'diary.authors.revoke_admin': '取消管理员',
  'diary.authors.remove_title': '移除作者',
  'diary.authors.confirm_grant_admin': '确定要授予该用户管理员权限吗？',
  'diary.authors.confirm_revoke_admin': '确定要撤销该用户的管理员权限吗？',
  'diary.authors.admin_toggle_title': '权限设置',
  'diary.authors.confirm_remove': '确定要移除该作者吗？',
  'diary.authors.quit_title': '退出日记本',
  'diary.authors.confirm_quit': '确定要退出此日记本吗？',
  'diary.authors.invitation_sent': '邀请已发送',
  'diary.authors.admin_updated': '权限已更新',
  'diary.authors.author_removed': '作者已移除',
  'diary.authors.joined_successfully': '已成功加入日记本',
  'diary.authors.invitation_declined': '已拒绝邀请',
  'diary.authors.quit_successfully': '已成功退出日记本',

  // Backend Error Keys (Synced with internal/utils/errors/errors.go)
  // ... 其他错误键 ...
};
```

#### JavaScript 多语言 - `website/static/src/js/locales/en.js`

```javascript
export default {
  // Common
  'common.save': 'Save',
  'common.edit': 'Edit',
  'common.cancel': 'Cancel',
  'common.delete': 'Delete',
  'common.op_fail': 'Operation failed',
  'common.op_success': 'Operation successful',
  'common.load_fail': 'Load failed',
  'common.save_success': 'Saved successfully',
  'common.save_fail': 'Save failed',
  'common.net_error': 'Network error',
  'common.views': 'Views',
  'common.comments': 'Comments',
  'common.confirm': 'Confirm',

  // Diary
  'diary.create': 'New Diary',
  'diary.settings': 'Diary Settings',

  // Diary Authors (前端交互相关)
  'diary.authors.invite_title': 'Invite New Author',
  'diary.authors.invite_send': 'Send Invitation',
  'diary.authors.grant_admin': 'Grant Admin',
  'diary.authors.revoke_admin': 'Revoke Admin',
  'diary.authors.remove_title': 'Remove Author',
  'diary.authors.confirm_grant_admin':
    'Are you sure you want to grant admin permission to this user?',
  'diary.authors.confirm_revoke_admin':
    'Are you sure you want to revoke admin permission from this user?',
  'diary.authors.admin_toggle_title': 'Permission Settings',
  'diary.authors.confirm_remove':
    'Are you sure you want to remove this author?',
  'diary.authors.quit_title': 'Quit Diary',
  'diary.authors.confirm_quit': 'Are you sure you want to quit this diary?',
  'diary.authors.invitation_sent': 'Invitation sent',
  'diary.authors.admin_updated': 'Permission updated',
  'diary.authors.author_removed': 'Author removed',
  'diary.authors.joined_successfully': 'Successfully joined the diary',
  'diary.authors.invitation_declined': 'Invitation declined',
  'diary.authors.quit_successfully': 'Successfully quit the diary',

  // Backend Error Keys (Synced with internal/utils/errors/errors.go)
  // ... 其他错误键 ...
};
```

### 5.12 日记本模块更新 (`website/static/src/js/modules/diary-detail.js`)

```javascript
import { SocialMixin } from './social.js';
import { EntryMixin } from './entry.js';
import Authors from './authors.js';

export default function DiaryDetail() {
  return {
    ...SocialMixin(),
    ...EntryMixin(),
    ...Authors(),

    init() {
      // Initialize social features
      this.initSocial();

      // Read isOwner from DOM
      const mainEl = document.querySelector('main');
      this.isOwner = mainEl?.dataset.isOwner === 'true';
      this.diaryId = mainEl?.dataset.diaryId;

      // Expose to global window
      window.openCreateEntry = (diaryId) => this.openCreateEntry(diaryId);
      window.openEditEntry = (diaryId, entryId) =>
        this.openEditEntry(diaryId, entryId);
      window.deleteEntry = (diaryId, entryId) =>
        this.deleteEntry(diaryId, entryId);
      window.openEntryModal = (url) => this.openEntryModal(url);

      window.openDiarySettings = (diaryId) => this.openDiarySettings(diaryId);
      window.toggleDiaryLike = (diaryId) => this.toggleDiaryLike(diaryId);

      // Expose authors methods
      window.authorsMixin = this;

      // Listen for diary updates globally
      this._bindDiaryUpdateListener();
    },

    _bindDiaryUpdateListener() {
      const updateHandler = (evt) => {
        const diary = evt.detail.diaryUpdated || evt.detail.value || evt.detail;
        if (diary && typeof diary === 'object' && diary.diary_id) {
          this._updateDiaryUI(diary);
        }
      };
      document.body.addEventListener('diaryUpdated', updateHandler);
    },

    openDiarySettings(diaryId) {
      const uno = window.getUno();
      if (!uno) return;
      const t = (key) => this.$store.i18n.t(key);

      let closeHandler;
      const boxId =
        'box-diary-settings-' + Math.random().toString(36).slice(2, 10);

      uno.box.dialog({
        boxId: boxId,
        title: t('diary.settings'),
        hxUrl: `/form/diary?action=diary.edit&diary_id=${diaryId}`,
        widthClass: '!w-[880px]',
        heightClass: '!h-[570px]',
        showOkButton: true,
        okText: t('common.save'),
        containerClass: '!border',
        okButtonClass: 'btn',
        buttonAlign: 'right',
        onContentLoad: (el) => {
          if (!el) el = document.getElementById(boxId);
          if (window.Alpine) window.Alpine.initTree(el);
          if (window.htmx) window.htmx.process(el);

          const form = el.querySelector('form');
          if (form) {
            form.dataset.boxId = boxId;
          }
        },
        onOk: async () => {
          const form = document.getElementById('diary_editor_form');
          if (form && (await uno.valid.validateForm(form))) {
            form.dispatchEvent(new CustomEvent('validated', { bubbles: true }));
            return false;
          }
          return false;
        },
        onClose: () => {
          if (closeHandler)
            document.body.removeEventListener(
              'closeDiaryModal',
              closeHandler,
              true,
            );
        },
      });

      closeHandler = (evt) => {
        const eventBoxId = evt?.target?.dataset?.boxId;
        if (eventBoxId && eventBoxId !== boxId) {
          return;
        }
        uno.box.close(boxId);
      };

      document.body.addEventListener('closeDiaryModal', closeHandler, true);
    },

    _updateDiaryUI(diary) {
      const titleEl = document.getElementById('diary_title');
      const remarkEl = document.getElementById('diary_remark');

      if (titleEl) {
        let iconHtml = '';
        if (diary.view_setting === 2) {
          iconHtml = '<i class="mdi mdi-lock-outline mx-3"></i>';
        }
        titleEl.innerHTML = diary.title + iconHtml;
      }
      if (remarkEl) {
        remarkEl.innerText = diary.remark || '';
      }

      const mainEl = document.querySelector('main');
      if (mainEl) {
        if (diary.theme_id) {
          mainEl.classList.forEach((cls) => {
            if (cls.startsWith('dt') && !isNaN(cls.substring(2))) {
              mainEl.classList.remove(cls);
            }
          });
          mainEl.classList.add(`dt${diary.theme_id}`);
        }
        mainEl.dataset.socialSetting = diary.social_setting;
      }

      this.socialSetting = diary.social_setting;

      const socialWrapper = document.getElementById('diary-like-wrapper');
      if (diary.social_setting === 1) {
        if (!socialWrapper) {
          const pageNav = document.getElementById('page-nav');
          if (pageNav && pageNav.nextElementSibling) {
            const btnContainer = pageNav.nextElementSibling;
            const span = document.createElement('span');
            span.id = 'diary-like-wrapper';
            span.dataset.diaryId = diary.diary_id;
            span.innerHTML = `
                <button class="btn" @click="toggleDiaryLike(${diary.diary_id})" :disabled="diaryLike.loading">
                    <i class="mdi mr-1"
                        :class="diaryLike.liked ? 'mdi-thumb-up' : 'mdi-thumb-up-outline'"></i>
                    <span x-text="diaryLike.count"></span>
                </button>`;
            btnContainer.appendChild(span);
            if (window.Alpine) {
              window.Alpine.initTree(span);
            }
          }
        }
        this.loadStats();
        this.loadDiaryLike();
      } else {
        if (socialWrapper) socialWrapper.remove();
        this.stats = {};
      }
    },
  };
}
```

## 6. 邮件多语言机制设计

### 6.1 多语言架构说明

hebe 项目中存在三套多语言系统：

1. **模板多语言** (`website/locales/*.json`)
   - 用于服务端模板渲染（Go template）
   - 通过 `getMessage()` 函数在模板中使用
   - 键名使用点号分隔，如 `diary.authors.invite_title`

2. **JavaScript 多语言** (`website/static/src/js/locales/*.js`)
   - 用于前端 JavaScript 代码
   - 通过 `this.$store.i18n.t(key)` 函数使用
   - 键名使用点号分隔，如 `diary.authors.invite_title`

3. **邮件多语言**（新增）
   - 邮件内容使用模板多语言系统
   - 通过 `i18nService.GetMessage()` 获取邮件主题和正文
   - 根据被邀请用户的语言设置发送对应语言的邮件

### 6.2 邮件发送流程

```
邀请作者完整流程：
1. 用户在前端输入 @handle 或 email
2. 前端防重复提交检查（isSubmitting 标志，禁用提交按钮）
3. 后端验证用户权限和身份
4. 后端检查是否已邀请/已加入（防止重复）
5. 创建邀请记录（status = 2: Pending）
6. 获取被邀请用户的语言设置（优先）或使用邀请人的语言（备用）
7. 规范化语言代码（zh, en, zh-cn, zh-tw 等）
8. 使用被邀请用户的语言获取邮件主题和正文
9. 构建邮件正文参数（日记本标题、邀请人姓名、接受链接等）
10. 通过任务队列异步发送邮件（避免阻塞主流程）
11. 返回成功响应
12. 前端刷新作者列表
```

### 6.3 邮件模板设计原则

1. **HTML 格式**: 邮件正文必须是完整的 HTML 文档
2. **响应式设计**: 邮件模板需要支持移动端显示
3. **主题统一**: 使用 iDiary 品牌色（红色 #e53e3e）
4. **清晰可读**: 避免使用复杂嵌套，保持简洁明了
5. **安全链接**: 所有链接必须是完整的 URL（包含域名）
6. **纯文本降级**: 邮件客户端不支持 HTML 时，仍然可读
7. **避免脚本**: 不在邮件中使用 JavaScript

### 6.4 邮件参数说明

| 参数           | 说明             | 示例                             |
| -------------- | ---------------- | -------------------------------- |
| `diary_title`  | 日记本标题       | "我的旅行日记"                   |
| `owner_name`   | 邀请人显示名称   | "张三"                           |
| `owner_handle` | 邀请人 handle    | @zhangsan                        |
| `invitee_name` | 被邀请人显示名称 | "李四"                           |
| `accept_url`   | 接受邀请的链接   | https://idiary.ai/diary/123/join |
| `domain`       | 应用域名         | idiary.ai                        |

## 7. 实施计划

### 7.1 阶段划分

#### 阶段 1: 后端基础

- [x] 实现 `diary/types.go` 中的 `Author` 结构体
- [x] 实现 `diary/repository.go` 中的作者相关方法
- [x] 实现 `diary/service.go` 中的作者管理业务逻辑
- [x] 实现邮件发送方法 `sendInvitationEmail`
- [x] 实现语言规范化方法 `normalizeLanguage`
- [x] 编写单元测试（types、repository、service）

#### 阶段 2: 邮件多语言

- [x] 添加邮件多语言配置到 `website/locales/zh-cn.json`
- [x] 添加邮件多语言配置到 `website/locales/en.json`
- [x] 测试邮件发送功能（中英文）

#### 阶段 3: 权限集成

- [x] 更新 `entry/service.go` 中的权限验证逻辑（确保私有日记本的所有作者都能访问）
- [x] 编写集成测试
- [x] 编写单元测试

#### 阶段 4: 控制器和路由

- [x] 实现 `diary/requests.go`
- [x] 实现 `diary/handlers.go`（使用 `middleware.Validate` 中间件）
- [x] 更新 `diary/routes.go`
- [x] 编写单元测试

#### 阶段 5: 前端开发

- [x] 创建 `website/static/src/js/modules/authors.js`（使用项目统一的 CSRF 安全机制）
- [x] 更新 `website/static/src/js/locales/zh-cn.js`（添加作者相关键值）
- [x] 更新 `website/static/src/js/locales/en.js`（添加作者相关键值）
- [x] 创建 `website/views/diary/box_invite_author.html`（使用 `x-uno-valid` 验证）
- [x] 创建 `website/views/diary/authors.html`（作者列表独立页面）
- [x] 更新 `website/views/diary/index.html`（Authors 链接指向独立页面）
- [x] 创建 `website/views/user/joined_diaries_list.html`（融合邀请和已加入）
- [x] 创建 `website/views/user/joined_diaries_full.html`（完整列表页）
- [x] 更新 `website/views/user/index.html`
- [ ] 更新 `website/static/src/js/modules/diary-detail.js`

#### 阶段 6: 模板多语言

- [ ] 更新 `website/locales/zh-cn.json`
- [ ] 更新 `website/locales/en.json`

#### 阶段 7: 配置更新

- [ ] 在 `.env.local` 中添加 `IDIARY_USER_JOINED_DIARIES_LIMIT` 配置项
- [ ] 更新 `internal/core/config` 配置结构体以支持新配置

#### 阶段 8: 测试和联调

- [ ] 前后端联调
- [ ] 测试 CSRF Token 注入（fetch 和 HTMX）
- [ ] 测试邮件发送功能（中英文）
- [ ] 测试私有日记本权限
- [ ] 测试作者列表分页功能
- [ ] 测试加入日记本列表排序和分页

### 7.2 优先级

| 优先级 | 功能                 | 说明                                     |
| ------ | -------------------- | ---------------------------------------- |
| P0     | 邀请作者             | 核心功能，必须实现                       |
| P0     | 邮件通知             | 核心功能，必须实现                       |
| P0     | 接受/拒绝邀请        | 核心功能，必须实现                       |
| P0     | 作者写日记           | 核心功能，必须实现                       |
| P0     | 权限验证             | 核心功能，必须实现（包括私有日记本权限） |
| P0     | CSRF 安全            | 核心功能，必须实现（使用项目统一方案）   |
| P1     | 设置管理员           | 重要功能，应该实现                       |
| P1     | 移除作者             | 重要功能，应该实现                       |
| P1     | 退出日记本           | 重要功能，应该实现                       |
| P2     | 作者列表展示（Card） | 改进功能，可以延后                       |
| P2     | 加入的日记本列表     | 改进功能，可以延后                       |
| P3     | 通知系统             | 未来功能，可以延后                       |

## 8. 边界情况和错误处理

### 8.1 边界情况

1. **邀请自己**: 防止用户邀请自己
2. **邀请拥有者**: 防止拥有者再次被邀请
3. **重复邀请**: 检查是否已邀请或已加入
4. **邀请不存在用户**: Handle 或 Email 对应的用户不存在
5. **退出日记本**: 拥有者不能退出自己的日记本
6. **移除作者**: 不能移除拥有者
7. **设置管理员**: 不能设置拥有者为管理员
8. **删除日记**: 确保作者的发表数量正确更新
9. **并发操作**: 处理同时邀请、接受、拒绝的并发情况
10. **日记本删除**: 删除日记本时级联删除作者记录
11. **邮件发送失败**: 邮件发送失败不影响邀请记录的创建
12. **语言设置缺失**: 用户没有设置语言时使用默认语言
13. **私有日记本权限**: 确保私有日记本的所有作者（拥有者、管理员、普通作者）都能访问

### 8.2 错误处理

| 错误场景         | 错误代码             | HTTP 状态码 | 用户提示                             |
| ---------------- | -------------------- | ----------- | ------------------------------------ |
| 用户不存在       | USER_NOT_FOUND       | 404         | "用户不存在"                         |
| 不是日记本拥有者 | PERMISSION_DENIED    | 403         | "只有拥有者可以执行此操作"           |
| 已经是作者       | ALREADY_AUTHOR       | 400         | "该用户已经是作者"                   |
| 已经被邀请       | ALREADY_INVITED      | 400         | "该用户已被邀请"                     |
| 邀请不存在       | INVITATION_NOT_FOUND | 404         | "邀请不存在"                         |
| 不是待确认状态   | INVALID_PARAMS       | 400         | "不是待确认的邀请"                   |
| 拥有者不能退出   | INVALID_PARAMS       | 400         | "拥有者不能退出自己的日记本"         |
| 不能移除拥有者   | INVALID_PARAMS       | 400         | "不能移除拥有者"                     |
| 不能设置拥有者   | INVALID_PARAMS       | 400         | "不能设置拥有者为管理员"             |
| 邮件发送失败     | EMAIL_SEND_FAILED    | 500         | "邀请已发送成功，但邮件通知发送失败" |
| CSRF Token 无效  | CSRF_TOKEN_INVALID   | 403         | "CSRF Token 无效"                    |

## 9. 性能优化建议

1. **缓存作者列表**: 对于活跃的日记本，可以缓存作者列表
2. **批量查询**: 在展示加入的日记本时，批量查询日记本信息
3. **索引优化**: 为 `br_diary_author` 表添加合适的索引
4. **分页加载**: 作者列表和日记本列表使用分页加载
5. **懒加载**: 使用 HTMX 的懒加载功能
6. **邮件队列**: 使用任务队列异步发送邮件，避免阻塞主流程
7. **邮件模板缓存**: 缓存编译后的邮件模板，减少每次生成开销

## 10. 未来扩展

1. **通知系统**: 邀请发送、接受、拒绝时发送通知
2. **活动日志**: 记录作者管理操作日志
3. **角色扩展**: 添加更多角色类型（如编辑、审核员等）
4. **权限粒度**: 更细粒度的权限控制
5. **团队功能**: 支持团队邀请
6. **数据分析**: 作者贡献统计、活跃度分析
7. **邮件模板**: 支持自定义邮件模板
8. **批量邀请**: 支持批量邀请作者
9. **邀请链接过期**: 邀请链接设置有效期
10. 邀请撤销: 支持撤销已发送的邀请
11. 邀请重发: 支持重新发送失败的邀请
12. 邮件追踪: 记录邮件打开和点击情况
13. 移动端优化: 邮件在移动端的显示优化
