package post

import (
	"fmt"
	
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/internal/utils/errors"
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
//   - postRepo: 日记数据访问层接口
//   - i18nService: 国际化服务管理器
// 返回:
//   - *Service: 日记业务逻辑层实例
func NewService(config *config.Config, userRepo user.RepositoryInterface, taskService *task.Service, diaryRepo diary.RepositoryInterface, postRepo RepositoryInterface, i18nService i18n.I18nManager) *Service {
	return &Service{
		config:       config, 
		repo:         postRepo,
		userService:  user.NewService(config, userRepo, taskService, i18nService),
		diaryService: diary.NewService(config, diaryRepo),
	}
}

// ===================== 接口实现方法 =====================

// AddPost 增加新的日记（post）
// 返回新增日记ID
// 注意：此方法不再处理日记本数量更新，需要在调用方处理
// 参数:
//   - post: 要添加的日记对象
// 返回:
//   - int64: 新增日记的ID
//   - error: 添加失败时返回错误
func (s *Service) AddPost(post *Post) (int64, error) {
	// 验证用户ID
	if post.AddUserID <= 0 {
		return 0, errors.InvalidParams.WithDetails("无效的用户ID")
	}
	
	// 验证用户是否存在
	_, err := s.userService.GetUserByUserId(post.AddUserID)
	if err != nil {
		return 0, errors.UserNotFound.WithDetails("用户不存在")
	}
	
	// 验证日记本是否存在且属于该用户
	diary, err := s.diaryService.GetDiaryByID(post.DiaryID)
	if err != nil {
		return 0, errors.DiaryNotFound.WithDetails("日记本不存在")
	}
	
	// 检查日记本是否属于该用户
	if diary.AddUserID != post.AddUserID {
		return 0, errors.PermissionDenied.WithDetails("无权限访问该日记本")
	}
	
	//插入数据库
	postID, err := s.repo.Insert(post)
	if err != nil {
		return 0, &errors.DatabaseError
	}
	post.PostID = postID

	return post.PostID, nil
}

// GetPostByID 根据id获取post
// 参数:
//   - postId: 日记ID
// 返回:
//   - *Post: 日记对象，未找到时返回nil
//   - error: 查询失败时返回错误
func (s *Service) GetPostByID(postId int64) (*Post, error) {
	return s.repo.GetByID(postId)
}

// UpdatePost 更新 post
// 参数:
//   - post: 包含更新信息的日记对象
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdatePost(post *Post) error {
	// 验证内容不能为空
	if post.Body == "" {
		return errors.ValidationError.WithDetails("文章内容不能为空")
	}
	
	// 更新日记本
	err := s.repo.Update(post)
	if err != nil {
		return err
	}

	return nil
}

// ListPostsByUserID 根据用户ID和日记本ID分页查询日记信息
// 参数:
//   - userID: 用户ID
//   - diaryID: 日记本ID
//   - pageIndex: 页码（从0开始）
//   - pageSize: 每页大小
// 返回:
//   - []Post: 日记列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListPostsByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error) {
	return s.repo.ListByUserID(userID, diaryID, pageIndex, pageSize)
}

// GetPostCountByUserID 根据用户ID和日记本ID统计信息数量
// 参数:
//   - userID: 用户ID
//   - diaryID: 日记本ID
// 返回:
//   - int: 日记数量
//   - error: 查询失败时返回错误
func (s *Service) GetPostCountByUserID(userID int64, diaryID int64) (int, error) {
	return s.repo.GetCountByUserID(userID, diaryID)
}

// UpdatePostByID 根据日记ID更新日记信息
// 参数:
//   - post: 包含更新信息的日记对象
// 返回:
//   - error: 更新失败时返回错误
func (s *Service) UpdatePostByID(post *Post) error {
	return s.repo.UpdateByID(post)
}

// DeletePostByID 根据日记ID删除日记信息（软删除）
// 删除日记并减少对应日记本的统计数
// 参数:
//   - postID: 日记ID
//   - userID: 操作用户ID
//   - userName: 操作用户名
//   - actionTime: 操作时间（UTC时间戳）
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeletePostByID(postID int64, userID int64, userName string, actionTime int64) error {
	// 先获取日记信息，用于后续减少日记本统计数
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.PostNotFound.WithDetails(fmt.Sprintf("post_id: %d", postID))
	}
	
	// 执行软删除
	err = s.repo.DeleteByID(postID)
	if err != nil {
		return err
	}
	
	// 减少日记本的日记统计数
	err = s.diaryService.DecrementPostCount(post.DiaryID, userID, userName, actionTime)
	if err != nil {
		// 如果更新日记本数量失败，记录错误但不影响删除操作的成功
		// 这里可以添加日志记录
	}
	
	return nil
}

// GetPreviousPostByDiaryID 根据当前日记ID和日记本ID获取上一篇日记ID
// 参数:
//   - currentPostID: 当前日记ID
//   - diaryID: 日记本ID
// 返回:
//   - int64: 上一篇日记ID，未找到时返回0
//   - error: 查询失败时返回错误
func (s *Service) GetPreviousPostByDiaryID(currentPostID int64, diaryID int64) (int64, error) {
	return s.repo.GetPreviousPostByDiaryID(currentPostID, diaryID)
}

// GetNextPostByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
// 参数:
//   - currentPostID: 当前日记ID
//   - diaryID: 日记本ID
// 返回:
//   - int64: 下一篇日记ID，未找到时返回0
//   - error: 查询失败时返回错误
func (s *Service) GetNextPostByDiaryID(currentPostID int64, diaryID int64) (int64, error) {
	return s.repo.GetNextPostByDiaryID(currentPostID, diaryID)
}

// GetPostsStats 根据日记ID列表获取统计信息
func (s *Service) GetPostsStats(postIDs []int64) ([]Post, error) {
	if len(postIDs) == 0 {
		return []Post{}, nil
	}

	// 调用repository层获取统计信息
	posts, err := s.repo.GetStatsByIDs(postIDs)
	if err != nil {
		return nil, err
	}

	return posts, nil
}