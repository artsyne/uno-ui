package post

import (
	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/pkg"
)

// Service 日记业务逻辑层
// 负责处理所有与日记相关的业务逻辑和规则验证
type Service struct {
	repo   *Repository
	config *configs.Config // 缓存配置，避免重复调用全局变量
}

// NewService 创建新的日记业务逻辑层实例
func NewService() *Service {
	config := configs.GetConfig()
	if config == nil {
		// 在初始化时就检查配置，如果配置未正确加载则panic
		// 这是合理的，因为没有配置应用无法正常运行
		panic("Failed to get config during post service initialization - config is nil")
	}

	return &Service{
		repo:   NewRepository(),
		config: config, // 缓存配置
	}
}

// AddPost 增加新的日记（post）
// 返回新增日记ID
// 注意：此方法不再处理日记本数量更新，需要在调用方处理
func (s *Service) AddPost(post *Post) (int64, error) {
	//截取配置长度字符串，末尾增加 suffix
	post.Summary = pkg.TruncateString(post.Body, s.config.App.SummaryMax, s.config.App.Suffix)

	// todo：需判断空间是否有权限发布日记

	//插入数据库
	postID, err := s.repo.Insert(post)
	if err != nil {
		return 0, errors.DbOperationErr
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
	//校验参数：todo

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
