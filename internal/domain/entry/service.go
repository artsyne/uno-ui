package entry

import (
	"log/slog"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/pkg"
)

// Service 日记业务逻辑层实现
// 负责处理所有与日记相关的业务逻辑和规则验证
type Service struct {
	config       *config.Config      // 配置管理器
	repo         RepositoryInterface // 日记数据访问层接口
	userService  *user.Service       // 用户业务逻辑层
	diaryService *diary.Service      // 日记本业务逻辑层
}

// NewService 创建新的文章业务逻辑层实例
// 参数:
//   - config: 配置管理器
//   - userService: 用户业务逻辑层
//   - taskService: 任务业务逻辑层
//   - diaryRepo: 日记本数据访问层接口
//   - entryRepo: 日记数据访问层接口
//   - i18nService: 国际化服务管理器
//
// 返回:
//   - *Service: 日记业务逻辑层实例
func NewService(config *config.Config, userService *user.Service, taskService *task.Service, diaryRepo diary.RepositoryInterface, entryRepo RepositoryInterface, i18nService i18n.I18nManager) *Service {
	return &Service{
		config:       config,
		repo:         entryRepo,
		userService:  userService,
		diaryService: diary.NewService(config, diaryRepo, userService, taskService, i18nService),
	}
}

// ===================== 接口实现方法 =====================

// AddEntry 增加新的日记（entry）
// 返回新增日记ID
// 注意：此方法不再处理日记本数量更新，需要在调用方处理
// 参数:
//   - entry: 要添加的日记对象
//
// 返回:
//   - int64: 新增日记的ID
//   - error: 添加失败时返回错误
func (s *Service) AddEntry(entry *Entry) (int64, error) {
	// 验证用户ID
	if entry.AddUserID <= 0 {
		return 0, errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_user_id"})
	}

	// 验证用户是否存在
	_, err := s.userService.GetUserByUserId(entry.AddUserID)
	if err != nil {
		return 0, errors.UserNotFound.WithDetails(map[string]interface{}{"reason": "user_not_found"})
	}

	// 验证日记本是否存在
	diary, err := s.diaryService.GetDiaryByID(entry.DiaryID)
	if err != nil {
		return 0, errors.DiaryNotFound.WithDetails(map[string]interface{}{"reason": "diary_not_found"})
	}
	if diary == nil {
		return 0, errors.DiaryNotFound.WithDetails(map[string]interface{}{"reason": "diary_not_found_nil"})
	}

	// 检查是否有写权限 (拥有者或活跃作者)
	hasPermission, err := s.checkWritePermission(diary, entry.AddUserID)
	if err != nil {
		return 0, err
	}
	if !hasPermission {
		return 0, errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "diary_permission_denied"})
	}

	//插入数据库
	entryID, err := s.repo.Insert(entry)
	if err != nil {
		return 0, &errors.DatabaseError
	}
	entry.EntryID = entryID

	return entry.EntryID, nil
}

// checkWritePermission 检查是否有写权限
func (s *Service) checkWritePermission(diary *diary.Diary, userID int64) (bool, error) {
	// 1. 拥有者有权限
	if diary.AddUserID == userID {
		return true, nil
	}

	// 2. 活跃作者有权限
	isAuthor, err := s.diaryService.IsActiveAuthor(diary.DiaryID, userID)
	if err != nil {
		return false, err
	}

	return isAuthor, nil
}

// GetEntryByID 根据id获取entry
// 参数:
//   - entryID: 日记ID
//
// 返回:
//   - *Entry: 日记对象，未找到时返回nil
//   - error: 查询失败时返回错误
func (s *Service) GetEntryByID(entryID int64) (*Entry, error) {
	return s.repo.GetByID(entryID)
}

// GetEntryByIDForUser 根据id获取entry（包含草稿/回收站），并校验归属用户
// 关键逻辑：自动保存草稿/发布草稿需要访问非正常状态的 Entry，因此不能复用只查 status=active 的 GetEntryByID
func (s *Service) GetEntryByIDForUser(entryID int64, userID int64) (*Entry, error) {
	if userID <= 0 {
		return nil, errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_user_id"})
	}

	// 1. 获取原始日记数据（忽略用户限制）
	entry, err := s.repo.GetByIDRaw(entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}

	// 2. 获取日记本信息
	diaryData, err := s.diaryService.GetDiaryByID(entry.DiaryID)
	if err != nil {
		return nil, err
	}
	if diaryData == nil {
		return nil, errors.DiaryNotFound
	}

	// 3. 检查读权限
	// 3.1 如果是作者本人，始终有权限
	if entry.AddUserID == userID {
		return entry, nil
	}

	// 3.2 如果不是作者本人，检查是否为非Active状态
	// 草稿和回收站状态：作者、日记本拥有者、管理员可见
	if entry.Status != StatusActive {
		// 检查是否有管理权限 (拥有者或管理员)
		// 注意：这里重复了一部分 CanUpdateEntry 的逻辑，但为了性能（避免循环调用或过度复杂），直接在这里检查
		if diaryData.AddUserID == userID {
			// 拥有者可见
		} else {
			// 检查管理员权限
			author, err := s.diaryService.GetAuthor(entry.DiaryID, userID)
			if err != nil {
				return nil, errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_author_of_draft"})
			}
			if author.IsAdmin == diary.AuthorIsAdminYes && author.Status == diary.AuthorStatusActive {
				// 管理员可见
			} else {
				return nil, errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "not_author_of_draft"})
			}
		}
	}

	// 3.3 Active状态检查日记本权限
	// 私有日记本：需要是成员（拥有者/管理员/活跃作者）
	// 公开/密码日记本：成员始终可见
	if diaryData.ViewSetting == 2 { // Private
		isMember, err := s.diaryService.IsActiveAuthor(diaryData.DiaryID, userID)
		if err != nil {
			return nil, err
		}
		// 拥有者也是成员（IsActiveAuthor logic should cover Owner if added as author, but usually Owner is separate check）
		// IsActiveAuthor implementation in diary/service.go usually checks br_diary_author table.
		// Owner is usually implicitly an author or handled separately.
		// checkWritePermission handles Owner explicitly. Let's do same here.
		if diaryData.AddUserID == userID {
			isMember = true
		}

		if !isMember {
			return nil, errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "private_diary_access_denied"})
		}
	}

	// 对于公开(1)和密码(3)日记本，只要是Active状态，通常允许访问（GetByIDForUser 通常用于编辑/管理，前端展示用 GetEntryByID）
	// 但 GetEntryByIDForUser 也可能被用于 "View as User".
	// 如果是用于 "编辑"，那么需要写权限。这里是 Get，是读权限。
	// 假设此接口用于"获取详情以展示或编辑"。

	return entry, nil
}

// UpdateEntry 更新 entry
// 参数:
//   - entry: 包含更新信息的日记对象
//   - operatorID: 操作人ID
//
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdateEntry(entry *Entry, operatorID int64) error {
	// 验证内容不能为空
	if entry.Body == "" {
		return errors.ValidationError.WithDetails(map[string]interface{}{"field": "body", "reason": "required"})
	}

	// 权限检查
	canUpdate, err := s.CanUpdateEntry(entry.EntryID, operatorID)
	if err != nil {
		return err
	}
	if !canUpdate {
		return errors.PermissionDenied
	}

	// 更新日记本
	err = s.repo.Update(entry)
	if err != nil {
		return err
	}

	return nil
}

// ListEntriesByUserID 根据用户ID和日记本ID分页查询日记信息
// 参数:
//   - userID: 用户ID
//   - diaryID: 日记本ID
//   - status: 日记状态
//   - pageIndex: 页码（从0开始）
//   - pageSize: 每页大小
//
// 返回:
//   - []Entry: 日记列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListEntriesByUserID(userID int64, diaryID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error) {
	return s.repo.ListByUserID(userID, diaryID, status, pageIndex, pageSize)
}

// ListEntriesByDiaryID 根据日记本ID分页查询日记信息
// 权限校验：
// - 私有日记本：仅成员（拥有者、管理员、活跃作者）可见
// - 公开/密码日记本：通常可见（密码校验通常在Controller层处理或通过Session状态，此处暂不强制校验密码，仅校验成员资格用于私有日记本）
func (s *Service) ListEntriesByDiaryID(diaryID int64, userID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error) {
	// 1. 获取日记本信息
	diaryData, err := s.diaryService.GetDiaryByID(diaryID)
	if err != nil {
		return nil, 0, err
	}
	if diaryData == nil {
		return nil, 0, errors.DiaryNotFound
	}

	// 2. 权限检查
	// 如果是私有日记本，必须是成员
	if diaryData.ViewSetting == 2 { // Private
		// 检查是否为拥有者
		if diaryData.AddUserID == userID {
			// 拥有者有权限
		} else {
			// 检查是否为活跃作者
			isAuthor, err := s.diaryService.IsActiveAuthor(diaryID, userID)
			if err != nil {
				return nil, 0, err
			}
			if !isAuthor {
				return nil, 0, errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "private_diary_access_denied"})
			}
		}
	}

	// 3. 查询
	return s.repo.ListByDiaryID(diaryID, status, pageIndex, pageSize)
}

// GetEntryCountByUserID 根据用户ID和日记本ID统计信息数量
// 参数:
//   - userID: 用户ID
//   - diaryID: 日记本ID
//   - status: 日记状态
//
// 返回:
//   - int: 日记数量
//   - error: 查询失败时返回错误
func (s *Service) GetEntryCountByUserID(userID int64, diaryID int64, status int8) (int, error) {
	return s.repo.GetCountByUserID(userID, diaryID, status)
}

// ListAllEntriesByUserID 根据用户ID分页查询日记信息（跨日记本）
// 参数:
//   - userID: 用户ID
//   - status: 日记状态
//   - pageIndex: 页码（从0开始）
//   - pageSize: 每页大小
//
// 返回:
//   - []Entry: 日记列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListAllEntriesByUserID(userID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error) {
	return s.repo.ListAllByUserID(userID, status, pageIndex, pageSize)
}

// CanUpdateEntry 检查用户是否有权限更新指定日记
// 权限规则：
// 1. 日记作者本人可以更新
// 2. 日记本拥有者可以更新
// 3. 日记本管理员可以更新
func (s *Service) CanUpdateEntry(entryID, userID int64) (bool, error) {
	// 获取原始日记数据
	entry, err := s.repo.GetByIDRaw(entryID)
	if err != nil {
		return false, err
	}

	// 1. 检查是否为日记作者
	if entry.AddUserID == userID {
		return true, nil
	}

	// 获取日记本信息
	diaryData, err := s.diaryService.GetDiaryByID(entry.DiaryID)
	if err != nil {
		return false, err
	}
	if diaryData == nil {
		return false, errors.DiaryNotFound
	}

	// 2. 检查是否为日记本拥有者
	if diaryData.AddUserID == userID {
		return true, nil
	}

	// 3. 检查是否为管理员
	author, err := s.diaryService.GetAuthor(entry.DiaryID, userID)
	if err != nil {
		// 如果不是作者，或者查询出错，视为无权限
		return false, nil
	}

	if author.IsAdmin == diary.AuthorIsAdminYes && author.Status == diary.AuthorStatusActive {
		return true, nil
	}

	return false, nil
}

// UpdateEntryByID 根据日记ID更新日记信息
// 参数:
//   - entry: 包含更新信息的日记对象
//   - operatorID: 操作人ID
//
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdateEntryByID(entry *Entry, operatorID int64) error {
	// 权限检查
	canUpdate, err := s.CanUpdateEntry(entry.EntryID, operatorID)
	if err != nil {
		return err
	}
	if !canUpdate {
		return errors.PermissionDenied
	}
	return s.repo.UpdateByID(entry)
}

// DeleteEntryByID 根据日记ID删除日记信息（软删除）
// 删除日记并减少对应日记本的统计数
// 参数:
//   - entryID: 日记ID
//   - userID: 操作用户ID
//   - actionTime: 操作时间（UTC时间戳）
//
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeleteEntryByID(entryID int64, userID int64, actionTime int64) error {
	// 先获取日记信息，用于后续减少日记本统计数
	// 使用 GetByIDRaw 以支持管理员/拥有者删除他人日记
	entry, err := s.repo.GetByIDRaw(entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}

	// 验证删除权限
	// 权限规则与更新一致：作者、拥有者、管理员可删除
	canDelete, err := s.CanUpdateEntry(entryID, userID)
	if err != nil {
		return err
	}
	if !canDelete {
		return errors.PermissionDenied
	}

	// 执行软删除
	err = s.repo.DeleteByID(entryID)
	if err != nil {
		return err
	}

	// 减少日记本的日记统计数（如果原状态是Active）
	if entry.Status == StatusActive {
		err = s.diaryService.DecrementEntryCount(entry.DiaryID)
		if err != nil {
			slog.Error("减少日记统计失败", "error", err, "entry_id", entryID, "diary_id", entry.DiaryID)
		}
	}

	return nil
}

// HardDeleteEntryByID 根据日记ID物理删除日记信息
// 参数:
//   - entryID: 日记ID
//   - userID: 操作用户ID
//
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) HardDeleteEntryByID(entryID int64, userID int64) error {
	// 先获取日记信息
	entry, err := s.repo.GetByIDRaw(entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}

	// 验证删除权限
	canDelete, err := s.CanUpdateEntry(entryID, userID)
	if err != nil {
		return err
	}
	if !canDelete {
		return errors.PermissionDenied
	}

	// 执行物理删除
	err = s.repo.HardDeleteByID(entryID)
	if err != nil {
		return err
	}

	// 如果原状态是Active，减少统计数
	if entry.Status == StatusActive {
		_ = s.diaryService.DecrementEntryCount(entry.DiaryID)
	}

	return nil
}

// RecoverEntryByID 恢复已删除的日记
// 参数:
//   - entryID: 日记ID
//   - userID: 操作用户ID
//
// 返回:
//   - error: 恢复失败时返回错误
func (s *Service) RecoverEntryByID(entryID int64, userID int64) error {
	// 获取日记信息
	entryData, err := s.repo.GetByIDRaw(entryID)
	if err != nil {
		return err
	}
	if entryData == nil {
		return errors.EntryNotFound
	}

	// 验证权限 (使用 CanUpdateEntry 统一管理)
	canRecover, err := s.CanUpdateEntry(entryID, userID)
	if err != nil {
		return err
	}
	if !canRecover {
		return errors.PermissionDenied
	}

	// 如果已经是Active，直接返回
	if entryData.Status == StatusActive {
		return nil
	}

	// 更新状态为Active
	entryData.Status = StatusActive
	entryData.UpdateUTCTime = pkg.GetNowTimeUnixMilli()

	if err := s.repo.UpdateByID(entryData); err != nil {
		return err
	}

	// 增加统计数
	_ = s.diaryService.IncrementEntryCount(entryData.DiaryID)

	return nil
}

// GetPreviousEntryByDiaryID 根据当前日记ID和日记本ID获取上一篇日记ID
// 参数:
//   - currentEntryID: 当前日记ID
//   - diaryID: 日记本ID
//
// 返回:
//   - int64: 上一篇日记ID，未找到时返回0
//   - error: 查询失败时返回错误
func (s *Service) GetPreviousEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	return s.repo.GetPreviousEntryByDiaryID(currentEntryID, diaryID)
}

// GetNextEntryByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
// 参数:
//   - currentEntryID: 当前日记ID
//   - diaryID: 日记本ID
//
// 返回:
//   - int64: 下一篇日记ID，未找到时返回0
//   - error: 查询失败时返回错误
func (s *Service) GetNextEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	return s.repo.GetNextEntryByDiaryID(currentEntryID, diaryID)
}

// GetEntriesStats 根据日记ID列表获取统计信息
func (s *Service) GetEntriesStats(entryIDs []int64) ([]Entry, error) {
	if len(entryIDs) == 0 {
		return []Entry{}, nil
	}

	// 调用repository层获取统计信息
	entries, err := s.repo.GetStatsByIDs(entryIDs)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// IncrementLikeCount 增加日记的点赞数量
// 参数:
//   - entryID: 日记ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) IncrementLikeCount(entryID int64) error {
	return s.repo.IncrementLikeNum(entryID)
}

// DecrementLikeCount 减少日记的点赞数量
// 参数:
//   - entryID: 日记ID
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) DecrementLikeCount(entryID int64) error {
	return s.repo.DecrementLikeNum(entryID)
}
