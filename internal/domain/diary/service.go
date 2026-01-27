package diary

import (
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
)

// Service 日记本业务逻辑层实现
// 负责处理所有与日记本相关的业务逻辑和规则验证
type Service struct {
	config *config.Config      // 应用配置信息
	repo   RepositoryInterface // 数据访问层接口
}

// NewService 创建新的日记本业务逻辑层实例
// 参数:
//   - config: 应用配置对象
//   - repo: 数据访问层接口实现
//
// 返回:
//   - *Service: 日记本业务逻辑层实例
func NewService(config *config.Config, repo RepositoryInterface) *Service {
	return &Service{
		config: config,
		repo:   repo,
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

	// 新增日记本
	diaryID, err := s.repo.Insert(diary)
	if err != nil {
		return 0, err
	}
	diary.DiaryID = diaryID

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
		return errors.InvalidParams.WithDetails("日记本不为空，无法删除")
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
	if diary.SocialSetting == 2 {
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
//   - pageIndex: 页索引（从0开始）
//   - pageSize: 每页大小
//
// 返回:
//   - []Diary: 日记本列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListDiariesByUser(userID int64, pageIndex, pageSize int) ([]Diary, int, error) {
	diaries, _, err := s.repo.ListByUserID(userID, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if len(diaries) == 0 {
		return []Diary{}, 0, nil
	}

	count, err := s.repo.GetCountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	return diaries, count, nil
}

// IncrementEntryCount 增加日记本的日记数量
// 参数:
//   - diaryID: 日记本ID
//   - userID: 操作用户ID
//   - userName: 操作用户名
//   - actionTime: 操作时间（UTC时间戳）
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) IncrementEntryCount(diaryID int64, userID int64, userName string, actionTime int64) error {
	return s.repo.IncrementEntryNum(diaryID, userID, userName, actionTime)
}

// DecrementEntryCount 减少日记本的日记数量
// 参数:
//   - diaryID: 日记本ID
//   - userID: 操作用户ID
//   - userName: 操作用户名
//   - actionTime: 操作时间（UTC时间戳）
//
// 返回:
//   - error: 操作失败时返回错误
func (s *Service) DecrementEntryCount(diaryID int64, userID int64, userName string, actionTime int64) error {
	return s.repo.DecrementEntryNum(diaryID, userID, userName, actionTime)
}

// CheckDiaryLimit 检查用户日记本数量是否超过限制
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - error: 超过限制时返回错误
func (s *Service) CheckDiaryLimit(userID int64) error {
	daybookNum, err := s.repo.GetCountByUserID(userID)
	if err != nil {
		return &errors.InvalidParams
	}

	if daybookNum > s.config.App.MaxDiaryLimit {
		return &errors.DiaryLimitExceeded
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
	if diary.ViewSetting != 1 && diary.ViewSetting != 2 {
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
