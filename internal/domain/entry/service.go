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
		diaryService: diary.NewService(config, diaryRepo),
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

	// 验证日记本是否存在且属于该用户
	diary, err := s.diaryService.GetDiaryByID(entry.DiaryID)
	if err != nil {
		return 0, errors.DiaryNotFound.WithDetails(map[string]interface{}{"reason": "diary_not_found"})
	}
	if diary == nil {
		return 0, errors.DiaryNotFound.WithDetails(map[string]interface{}{"reason": "diary_not_found_nil"})
	}

	// 检查日记本是否属于该用户
	if diary.AddUserID != entry.AddUserID {
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
	return s.repo.GetByIDForUser(entryID, userID)
}

// UpdateEntry 更新 entry
// 参数:
//   - entry: 包含更新信息的日记对象
//
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdateEntry(entry *Entry) error {
	// 验证内容不能为空
	if entry.Body == "" {
		return errors.ValidationError.WithDetails(map[string]interface{}{"field": "body", "reason": "required"})
	}

	// 更新日记本
	err := s.repo.Update(entry)
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
// 返回:
//   - []Entry: 日记列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListAllEntriesByUserID(userID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error) {
	return s.repo.ListAllByUserID(userID, status, pageIndex, pageSize)
}

// UpdateEntryByID 根据日记ID更新日记信息
// 参数:
//   - entry: 包含更新信息的日记对象
//
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdateEntryByID(entry *Entry) error {
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
	entry, err := s.repo.GetByIDForUser(entryID, userID)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}
	
	// 验证所有权
	if entry.AddUserID != userID {
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
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) HardDeleteEntryByID(entryID int64, userID int64) error {
	// 先获取日记信息
	entry, err := s.repo.GetByIDForUser(entryID, userID)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.EntryNotFound.WithDetails(map[string]interface{}{"entry_id": entryID})
	}

	// 验证所有权
	if entry.AddUserID != userID {
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
// 返回:
//   - error: 恢复失败时返回错误
func (s *Service) RecoverEntryByID(entryID int64, userID int64) error {
	// 获取日记信息
	entryData, err := s.repo.GetByIDForUser(entryID, userID)
	if err != nil {
		return err
	}
	if entryData == nil {
		return errors.EntryNotFound
	}

	// 验证所有权
	if entryData.AddUserID != userID {
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
