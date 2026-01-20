package entry

import (
	"fmt"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// Service 日记业务逻辑层实现
// 负责处理所有与日记相关的业务逻辑和规则验证
type Service struct {
	config       *config.Config        // 配置管理器
	repo         RepositoryInterface   // 日记数据访问层接口
	userService  *user.Service         // 用户业务逻辑层
	diaryService *diary.Service        // 日记本业务逻辑层
}

// NewService 创建新的文章业务逻辑层实例
// 参数:
//   - config: 配置管理器
//   - userRepo: 用户数据访问层接口
//   - taskService: 任务业务逻辑层
//   - diaryRepo: 日记本数据访问层接口
//   - entryRepo: 日记数据访问层接口
//   - i18nService: 国际化服务管理器
// 返回:
//   - *Service: 日记业务逻辑层实例
func NewService(config *config.Config, userRepo user.RepositoryInterface, taskService *task.Service, diaryRepo diary.RepositoryInterface, entryRepo RepositoryInterface, i18nService i18n.I18nManager) *Service {
	return &Service{
		config:       config, 
		repo:         entryRepo,
		userService:  user.NewService(config, userRepo, taskService, i18nService),
		diaryService: diary.NewService(config, diaryRepo),
	}
}

// ===================== 接口实现方法 =====================

// AddEntry 增加新的日记（entry）
// 返回新增日记ID
// 注意：此方法不再处理日记本数量更新，需要在调用方处理
// 参数:
//   - entry: 要添加的日记对象
// 返回:
//   - int64: 新增日记的ID
//   - error: 添加失败时返回错误
func (s *Service) AddEntry(entry *Entry) (int64, error) {
	// 验证用户ID
	if entry.AddUserID <= 0 {
		return 0, errors.InvalidParams.WithDetails("无效的用户ID")
	}
	
	// 验证用户是否存在
	_, err := s.userService.GetUserByUserId(entry.AddUserID)
	if err != nil {
		return 0, errors.UserNotFound.WithDetails("用户不存在")
	}
	
	// 验证日记本是否存在且属于该用户
	diary, err := s.diaryService.GetDiaryByID(entry.DiaryID)
	if err != nil {
		return 0, errors.DiaryNotFound.WithDetails("日记本不存在")
	}
	
	// 检查日记本是否属于该用户
	if diary.AddUserID != entry.AddUserID {
		return 0, errors.PermissionDenied.WithDetails("无权限访问该日记本")
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
// 返回:
//   - *Entry: 日记对象，未找到时返回nil
//   - error: 查询失败时返回错误
func (s *Service) GetEntryByID(entryID int64) (*Entry, error) {
	return s.repo.GetByID(entryID)
}

// UpdateEntry 更新 entry
// 参数:
//   - entry: 包含更新信息的日记对象
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdateEntry(entry *Entry) error {
	// 验证内容不能为空
	if entry.Body == "" {
		return errors.ValidationError.WithDetails("文章内容不能为空")
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
//   - pageIndex: 页码（从0开始）
//   - pageSize: 每页大小
// 返回:
//   - []Entry: 日记列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListEntriesByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Entry, int, error) {
	return s.repo.ListByUserID(userID, diaryID, pageIndex, pageSize)
}

// GetEntryCountByUserID 根据用户ID和日记本ID统计信息数量
// 参数:
//   - userID: 用户ID
//   - diaryID: 日记本ID
// 返回:
//   - int: 日记数量
//   - error: 查询失败时返回错误
func (s *Service) GetEntryCountByUserID(userID int64, diaryID int64) (int, error) {
	return s.repo.GetCountByUserID(userID, diaryID)
}

// UpdateEntryByID 根据日记ID更新日记信息
// 参数:
//   - entry: 包含更新信息的日记对象
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
//   - userName: 操作用户名
//   - actionTime: 操作时间（UTC时间戳）
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeleteEntryByID(entryID int64, userID int64, userName string, actionTime int64) error {
	// 先获取日记信息，用于后续减少日记本统计数
	entry, err := s.repo.GetByID(entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.EntryNotFound.WithDetails(fmt.Sprintf("entry_id: %d", entryID))
	}

	// 执行软删除
	err = s.repo.DeleteByID(entryID)
	if err != nil {
		return err
	}

	// 减少日记本的日记统计数
	err = s.diaryService.DecrementEntryCount(entry.DiaryID, userID, userName, actionTime)
	if err != nil {
		// 如果更新日记本数量失败，记录错误但不影响删除操作的成功
		// 这里可以添加日志记录
	}
	
	return nil
}

// GetPreviousEntryByDiaryID 根据当前日记ID和日记本ID获取上一篇日记ID
// 参数:
//   - currentEntryID: 当前日记ID
//   - diaryID: 日记本ID
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
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) IncrementLikeCount(entryID int64) error {
	return s.repo.IncrementLikeNum(entryID)
}

// DecrementLikeCount 减少日记的点赞数量
// 参数:
//   - entryID: 日记ID
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) DecrementLikeCount(entryID int64) error {
	return s.repo.DecrementLikeNum(entryID)
}