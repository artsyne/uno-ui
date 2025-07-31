package diary

import (
	"github.com/biuren/hebe/internal/core/config"
	errors "github.com/biuren/hebe/internal/core/error"
)

// ServiceInterface 日记本业务逻辑层接口
// 定义所有与日记本相关的业务逻辑操作
type ServiceInterface interface {
	// AddDiary 创建一个新的日记本
	AddDiary(diary *Diary) (int64, error)

	// UpdateDiary 更新日记本信息
	UpdateDiary(diary *Diary) error

	// GetDiaryByID 根据ID获取日记本
	GetDiaryByID(diaryID int64) (*Diary, error)

	// ListDiariesForWrite 获取用户的日记本列表（用于写日记界面）
	ListDiariesForWrite(userID int64) ([]Diary, error)

	// ListDiariesByUser 获取用户的分页日记本列表
	ListDiariesByUser(userID int64, pageIndex, pageSize int) ([]Diary, int, error)

	// IncrementPostCount 增加日记本的日记数量
	IncrementPostCount(diaryID int64, userID int64, userName string, actionTime int64) error

	// DecrementPostCount 减少日记本的日记数量
	DecrementPostCount(diaryID int64, userID int64, userName string, actionTime int64) error

	// CheckDiaryLimit 检查用户日记本数量是否超过限制
	CheckDiaryLimit(userID int64) error
}

// Service 日记本业务逻辑层实现
// 负责处理所有与日记本相关的业务逻辑和规则验证
type Service struct {
	repo   RepositoryInterface
	config *config.Config // 缓存配置，避免重复调用全局变量
}

// NewService 创建新的日记本业务逻辑层实例
func NewService() ServiceInterface {
	config := config.GetConfig()
	if config == nil {
		// 在初始化时就检查配置，如果配置未正确加载则panic
		// 这是合理的，因为没有配置应用无法正常运行
		panic("Failed to get config during diary service initialization - config is nil")
	}

	return &Service{
		repo:   NewRepository(),
		config: config, // 缓存配置
	}
}

// ===================== 接口实现方法 =====================

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

// CheckDiaryLimit 检查用户日记本数量是否超过限制
// 参数: userID - 用户ID
// 返回: 可能的错误
func (s *Service) CheckDiaryLimit(userID int64) error {
	daybookNum, err := s.repo.GetCountByUserID(userID)
	if err != nil {
		return errors.InvalidParams.Err()
	}

	if daybookNum > s.config.App.MaxDaybookLimit {
		return errors.DiaryLimitExceeded.Err()
	}

	return nil
}

// ===================== 私有辅助方法 =====================

// validateDiary 校验日记本参数
// 参数: diary - 要校验的日记本对象
// 返回: 可能的错误
func (s *Service) validateDiary(diary *Diary) error {
	if len(diary.Title) < 1 {
		return errors.InvalidParams.Err()
	}
	if diary.ViewSetting != 1 && diary.ViewSetting != 2 {
		return errors.InvalidParams.Err()
	}
	if diary.SocialSetting != 1 && diary.SocialSetting != 2 {
		return errors.InvalidParams.Err()
	}
	if diary.ThemeID < 1 || diary.ThemeID > 16 {
		return errors.InvalidParams.Err()
	}
	return nil
}