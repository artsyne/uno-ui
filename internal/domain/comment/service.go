package comment

import (
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
)

// Service 评论业务逻辑层实现
// 负责处理所有与评论相关的业务逻辑和规则验证
type Service struct {
	config       *config.Config      // 配置管理器
	repo         RepositoryInterface // 评论数据访问层接口
	userService  *user.Service       // 用户业务逻辑层
	entryService *entry.Service      // 日记业务逻辑层
	diaryService *diary.Service      // 日记本业务逻辑层
}

// NewService 创建新的评论业务逻辑层实例
// 参数:
//   - config: 配置管理器
//   - commentRepo: 评论数据访问层接口
//   - userService: 用户业务逻辑层
//   - entryService: 日记业务逻辑层
//   - diaryService: 日记本业务逻辑层
//
// 返回:
//   - *Service: 评论业务逻辑层实例
func NewService(config *config.Config, commentRepo RepositoryInterface, userService *user.Service, entryService *entry.Service, diaryService *diary.Service) *Service {
	return &Service{
		config:       config,
		repo:         commentRepo,
		userService:  userService,
		entryService: entryService,
		diaryService: diaryService,
	}
}

// ===================== 接口实现方法 =====================

// AddComment 增加新的评论
// 参数:
//   - comment: 要添加的评论对象
//
// 返回:
//   - int64: 新增评论的ID
//   - error: 添加失败时返回错误
func (s *Service) AddComment(comment *Comment) (int64, error) {
	// 验证用户ID
	if comment.AddUserID <= 0 {
		return 0, errors.InvalidParams.WithDetails(map[string]interface{}{"reason": "invalid_user_id"})
	}

	// 验证用户是否存在
	_, err := s.userService.GetUserByUserId(comment.AddUserID)
	if err != nil {
		return 0, errors.UserNotFound.WithDetails(map[string]interface{}{"reason": "user_not_found"})
	}

	// 验证日记是否存在
	// 优先检查用户是否有权访问（支持软删除/草稿），如果是日记作者或所有者
	entryData, err := s.entryService.GetEntryByIDForUser(comment.EntryID, comment.AddUserID)
	if err != nil {
		// 如果不是作者，尝试作为普通公开日记获取（仅Active）
		entryData, err = s.entryService.GetEntryByID(comment.EntryID)
		if err != nil {
			return 0, errors.EntryNotFound.WithDetails(map[string]interface{}{"reason": "entry_not_found"})
		}
	}

	// 验证日记本是否存在
	_, err = s.diaryService.GetDiaryByID(entryData.DiaryID)
	if err != nil {
		return 0, errors.DiaryNotFound.WithDetails(map[string]interface{}{"reason": "diary_not_found"})
	}

	// 设置创建时间为 UTC 毫秒时间戳
	comment.AddUtcTime = pkg.GetNowTimeUnixMilli()

	// 插入数据库
	commentID, err := s.repo.Insert(comment)
	if err != nil {
		return 0, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()})
	}

	// 增加日记的评论数量
	err = s.repo.IncrementEntryCommentNum(comment.EntryID)
	if err != nil {
		return 0, errors.DatabaseError.WithDetails(map[string]interface{}{"error": err.Error()})
	}

	return commentID, nil
}

// DeleteComment 删除评论
// 删除逻辑：日记作者或者日记本所有者或者评论作者可以删除
// 参数:
//   - commentID: 评论ID
//   - userID: 操作用户ID
//
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeleteComment(commentID int64, userID int64) error {
	// 获取评论信息
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return errors.CommentNotFound.WithDetails(map[string]interface{}{"reason": "comment_not_found"})
	}

	// 验证日记是否存在
	// 优先检查用户是否有权访问（支持软删除/草稿），如果是日记作者或所有者
	// 尝试作为所有者获取（支持所有状态）
	entryData, err := s.entryService.GetEntryByIDForUser(comment.EntryID, userID)
	if err != nil {
		// 如果不是作者，尝试作为普通公开日记获取（仅Active）
		entryData, err = s.entryService.GetEntryByID(comment.EntryID)
		if err != nil {
			return errors.EntryNotFound.WithDetails(map[string]interface{}{"reason": "entry_not_found"})
		}
	}

	// 获取日记本信息
	diaryData, err := s.diaryService.GetDiaryByID(entryData.DiaryID)
	if err != nil {
		return errors.DiaryNotFound.WithDetails(map[string]interface{}{"reason": "diary_not_found"})
	}

	// 检查删除权限：评论作者、日记作者或日记本所有者
	if comment.AddUserID != userID && entryData.AddUserID != userID && diaryData.AddUserID != userID {
		return errors.PermissionDenied.WithDetails(map[string]interface{}{"reason": "comment_permission_denied"})
	}

	// 删除评论
	err = s.repo.Delete(commentID)
	if err != nil {
		return &errors.DatabaseError
	}

	// 减少日记的评论数量
	err = s.repo.DecrementEntryCommentNum(comment.EntryID)
	if err != nil {
		return &errors.DatabaseError
	}

	return nil
}

// DeleteCommentsByEntryID 根据日记ID删除该日记下的所有评论
// 参数:
//   - entryID: 日记ID
//
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeleteCommentsByEntryID(entryID int64) error {
	return s.repo.DeleteByEntryID(entryID)
}

// GetCommentByID 根据评论ID获取评论信息
// 参数:
//   - commentID: 评论ID
//
// 返回:
//   - *Comment: 评论对象
//   - error: 获取失败时返回错误
func (s *Service) GetCommentByID(commentID int64) (*Comment, error) {
	return s.repo.GetByID(commentID)
}



// ListCommentsByEntryID 根据日记ID分页查询评论列表
// 参数:
//   - entryID: 日记ID
//   - pageIndex: 页码（从1开始）
//   - pageSize: 每页大小
//
// 返回:
//   - []Comment: 评论列表
//   - bool: 是否有更多数据
//   - error: 查询失败时返回错误
func (s *Service) ListCommentsByEntryID(entryID int64, pageIndex, pageSize int) ([]Comment, bool, error) {
	// 不再在此处验证日记是否存在或状态，因为调用方（如Controller）可能已经验证过
	// 或者调用方需要获取非Active状态日记的评论（如回收站中的日记）
	// 安全性由调用方保证

	return s.repo.ListByEntryID(entryID, pageIndex, pageSize)
}

// GetCommentCountByEntryID 根据日记ID统计评论数量
// 参数:
//   - entryID: 日记ID
//
// 返回:
//   - int: 评论数量
//   - error: 查询失败时返回错误
func (s *Service) GetCommentCountByEntryID(entryID int64) (int, error) {
	return s.repo.GetCountByEntryID(entryID)
}
