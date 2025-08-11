package post

import (
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
)

// Service 日记业务逻辑层实现
// 负责处理所有与日记相关的业务逻辑和规则验证
type Service struct {
	config       *config.Config 
	repo         RepositoryInterface
	userService  *user.Service
	diaryService *diary.Service
}

// NewService 创建新的文章业务逻辑层实例
func NewService(config *config.Config, userRepo user.RepositoryInterface, taskService *task.Service, diaryRepo diary.RepositoryInterface, postRepo RepositoryInterface) *Service {
	return &Service{
		config:       config, 
		repo:         postRepo,
		userService:  user.NewService(config, userRepo, taskService),
		diaryService: diary.NewService(config, diaryRepo),
	}
}

// ===================== 接口实现方法 =====================

// AddPost 增加新的日记（post）
// 返回新增日记ID
// 注意：此方法不再处理日记本数量更新，需要在调用方处理
func (s *Service) AddPost(post *Post) (int64, error) {
	// 校验参数
	if err := s.validatePostForAdd(post); err != nil {
		return 0, err
	}

	//截取配置长度字符串，末尾增加 suffix
	post.Summary = pkg.TruncateString(post.Body, s.config.App.SummaryMax, s.config.App.Suffix)

	// todo：需判断空间是否有权限发布日记

	//插入数据库
	postID, err := s.repo.Insert(post)
	if err != nil {
		return 0, errors.DatabaseError.Err()
	}
	post.PostID = postID

	return post.PostID, nil
}

// GetPostByID 根据id获取post
func (s *Service) GetPostByID(postId int64) (*Post, error) {
	return s.repo.GetByID(postId)
}

// UpdatePost 更新 post
func (s *Service) UpdatePost(post *Post) error {
	// 校验参数
	if err := s.validatePost(post); err != nil {
		return err
	}

	// 更新日记本
	err := s.repo.Update(post)
	if err != nil {
		return err
	}

	return nil
}

// ListPostsByUserID 根据用户ID和日记本ID分页查询日记信息
func (s *Service) ListPostsByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error) {
	return s.repo.ListByUserID(userID, diaryID, pageIndex, pageSize)
}

// GetPostCountByUserID 根据用户ID和日记本ID统计信息数量
func (s *Service) GetPostCountByUserID(userID int64, diaryID int64) (int, error) {
	return s.repo.GetCountByUserID(userID, diaryID)
}

// UpdatePostByID 根据日记ID更新日记信息
func (s *Service) UpdatePostByID(post *Post) error {
	return s.repo.UpdateByID(post)
}

// DeletePostByID 根据日记ID删除日记信息
func (s *Service) DeletePostByID(postID int64) error {
	return s.repo.DeleteByID(postID)
}

// ===================== 私有辅助方法 =====================

// validatePost 验证文章参数
func (s *Service) validatePost(post *Post) error {
	// 验证标题
	if post.Title == "" {
		return errors.InvalidParams.WithDetails("标题不能为空")
	}

	// 验证内容
	if post.Body == "" {
		return errors.InvalidParams.WithDetails("内容不能为空")
	}

	// 验证日记本ID
	if post.DiaryID <= 0 {
		return errors.InvalidParams.WithDetails("日记本ID无效")
	}

	// 检查日记本是否存在
	_, err := s.diaryService.GetDiaryByID(post.DiaryID)
	if err != nil {
		return err // 直接返回日记本服务的错误
	}

	return nil
}

// validatePostForAdd 用于添加文章时的验证，包括用户ID验证
func (s *Service) validatePostForAdd(post *Post) error {
	// 先进行基本验证
	if err := s.validatePost(post); err != nil {
		return err
	}

	// 验证用户ID
	if post.AddUserID <= 0 {
		return errors.InvalidParams.WithDetails("用户ID无效")
	}

	// 检查用户是否存在
	_, err := s.userService.GetUserByUserId(post.AddUserID)
	if err != nil {
		return err // 直接返回用户服务的错误
	}

	return nil
}
