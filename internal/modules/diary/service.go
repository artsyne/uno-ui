package diary

import (
	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
)

// Service 日记本业务逻辑层
// 负责处理所有与日记本相关的业务逻辑和规则验证
type Service struct {
	repo *Repository
}

// NewService 创建新的日记本业务逻辑层实例
func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// AddDiary 创建一个新的日记本
// 参数: diary - 要创建的日记本对象
// 返回: 新创建的日记本ID和可能的错误
func (s *Service) AddDiary(diary *Diary) (int64, error) {
	// 校验参数
	if err := s.validateDiary(diary); err != nil {
		return 0, err
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
// 参数: diary - 要更新的日记本对象
// 返回: 可能的错误
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

// GetDiaryByID 根据ID获取日记本
// 参数: diaryID - 日记本ID
// 返回: 日记本对象和可能的错误
func (s *Service) GetDiaryByID(diaryID int64) (*Diary, error) {
	return s.repo.GetByDiaryID(diaryID)
}

// ListDiariesForWrite 获取用户的日记本列表（用于写日记界面）
// 参数: userID - 用户ID
// 返回: 日记本列表和可能的错误
func (s *Service) ListDiariesForWrite(userID int64) ([]Diary, error) {
	return s.repo.ListForWriteByUserID(userID)
}

// ListDiariesByUser 获取用户的分页日记本列表
// 参数: userID - 用户ID, pageIndex - 页索引, pageSize - 每页大小
// 返回: 日记本列表、总数和可能的错误
func (s *Service) ListDiariesByUser(userID int64, pageIndex, pageSize int) ([]Diary, int, error) {
	return s.repo.ListByUserID(userID, pageIndex, pageSize)
}

// IncrementPostCount 增加日记本的日记数量
// 参数: diaryID - 日记本ID, userID - 用户ID, userName - 用户名, actionTime - 操作时间
// 返回: 可能的错误
func (s *Service) IncrementPostCount(diaryID int64, userID int64, userName string, actionTime int64) error {
	return s.repo.IncrementPostNum(diaryID, userID, userName, actionTime)
}

// DecrementPostCount 减少日记本的日记数量
// 参数: diaryID - 日记本ID, userID - 用户ID, userName - 用户名, actionTime - 操作时间
// 返回: 可能的错误
func (s *Service) DecrementPostCount(diaryID int64, userID int64, userName string, actionTime int64) error {
	return s.repo.DecrementPostNum(diaryID, userID, userName, actionTime)
}

// validateDiary 校验日记本参数
// 参数: diary - 要校验的日记本对象
// 返回: 可能的错误
func (s *Service) validateDiary(diary *Diary) error {
	if len(diary.Title) < 1 {
		return errors.ParseInputParamErr
	}
	if diary.ViewSetting != 1 && diary.ViewSetting != 2 {
		return errors.ParseInputParamErr
	}
	if diary.SocialSetting != 1 && diary.SocialSetting != 2 {
		return errors.ParseInputParamErr
	}
	if diary.ThemeID < 1 || diary.ThemeID > 16 {
		return errors.ParseInputParamErr
	}
	return nil
}

// CheckDiaryLimit 检查用户日记本数量是否超过限制
// 参数: userID - 用户ID
// 返回: 可能的错误
func (s *Service) CheckDiaryLimit(userID int64) error {
	daybookNum, err := s.repo.GetCountByUserID(userID)
	if err != nil {
		return errors.ParseInputParamErr
	}

	if daybookNum > configs.GetConfig().App.MaxDaybookLimit {
		return errors.MaxDaybookLimitErr
	}

	return nil
}
