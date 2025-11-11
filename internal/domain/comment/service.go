package comment

import (
	"time"

	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
)

// Service 评论业务逻辑层实现
// 负责处理所有与评论相关的业务逻辑和规则验证
type Service struct {
	config       *config.Config        // 配置管理器
	repo         RepositoryInterface   // 评论数据访问层接口
	userService  *user.Service         // 用户业务逻辑层
	postService  *post.Service         // 日记业务逻辑层
	diaryService *diary.Service        // 日记本业务逻辑层
}

// NewService 创建新的评论业务逻辑层实例
// 参数:
//   - config: 配置管理器
//   - commentRepo: 评论数据访问层接口
//   - userService: 用户业务逻辑层
//   - postService: 日记业务逻辑层
//   - diaryService: 日记本业务逻辑层
// 返回:
//   - *Service: 评论业务逻辑层实例
func NewService(config *config.Config, commentRepo RepositoryInterface, userService *user.Service, postService *post.Service, diaryService *diary.Service) *Service {
	return &Service{
		config:       config,
		repo:         commentRepo,
		userService:  userService,
		postService:  postService,
		diaryService: diaryService,
	}
}

// ===================== 接口实现方法 =====================

// AddComment 增加新的评论
// 参数:
//   - comment: 要添加的评论对象
// 返回:
//   - int64: 新增评论的ID
//   - error: 添加失败时返回错误
func (s *Service) AddComment(comment *Comment) (int64, error) {
	// 验证用户ID
	if comment.AddUserID <= 0 {
		return 0, errors.InvalidParams.WithDetails("无效的用户ID")
	}

	// 验证用户是否存在
	_, err := s.userService.GetUserByUserId(comment.AddUserID)
	if err != nil {
		return 0, errors.UserNotFound.WithDetails("用户不存在")
	}

	// 验证日记是否存在
	postData, err := s.postService.GetPostByID(comment.PostID)
	if err != nil {
		return 0, errors.PostNotFound.WithDetails("日记不存在")
	}

	// 验证日记本是否存在
	_, err = s.diaryService.GetDiaryByID(postData.DiaryID)
	if err != nil {
		return 0, errors.DiaryNotFound.WithDetails("日记本不存在")
	}

	// 设置创建时间
	comment.GmtCreate = time.Now().Format("2006-01-02 15:04:05")

	// 插入数据库
	commentID, err := s.repo.Insert(comment)
	if err != nil {
		return 0, &errors.DatabaseError
	}

	// 增加日记的评论数量
	err = s.repo.IncrementPostCommentNum(comment.PostID)
	if err != nil {
		return 0, &errors.DatabaseError
	}

	return commentID, nil
}

// DeleteComment 删除评论
// 删除逻辑：日记作者或者日记本所有者或者评论作者可以删除
// 参数:
//   - commentID: 评论ID
//   - userID: 操作用户ID
// 返回:
//   - error: 删除失败时返回错误
func (s *Service) DeleteComment(commentID int64, userID int64) error {
	// 获取评论信息
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return errors.CommentNotFound.WithDetails("评论不存在")
	}

	// 获取日记信息
	postData, err := s.postService.GetPostByID(comment.PostID)
	if err != nil {
		return errors.PostNotFound.WithDetails("日记不存在")
	}

	// 获取日记本信息
	diaryData, err := s.diaryService.GetDiaryByID(postData.DiaryID)
	if err != nil {
		return errors.DiaryNotFound.WithDetails("日记本不存在")
	}

	// 检查删除权限：评论作者、日记作者或日记本所有者
	if comment.AddUserID != userID && postData.AddUserID != userID && diaryData.AddUserID != userID {
		return errors.PermissionDenied.WithDetails("无权限删除该评论")
	}

	// 删除评论
	err = s.repo.Delete(commentID)
	if err != nil {
		return &errors.DatabaseError
	}

	// 减少日记的评论数量
	err = s.repo.DecrementPostCommentNum(comment.PostID)
	if err != nil {
		return &errors.DatabaseError
	}

	return nil
}

// GetCommentByID 根据评论ID获取评论信息
// 参数:
//   - commentID: 评论ID
// 返回:
//   - *Comment: 评论对象
//   - error: 获取失败时返回错误
func (s *Service) GetCommentByID(commentID int64) (*Comment, error) {
	return s.repo.GetByID(commentID)
}

// ListCommentsByPostID 根据日记ID分页查询评论列表
// 参数:
//   - postID: 日记ID
//   - pageIndex: 页码（从1开始）
//   - pageSize: 每页大小
// 返回:
//   - []Comment: 评论列表
//   - int: 总记录数
//   - error: 查询失败时返回错误
func (s *Service) ListCommentsByPostID(postID int64, pageIndex, pageSize int) ([]Comment, int, error) {
	// 验证日记是否存在
	_, err := s.postService.GetPostByID(postID)
	if err != nil {
		return nil, 0, errors.PostNotFound.WithDetails("日记不存在")
	}

	return s.repo.ListByPostID(postID, pageIndex, pageSize)
}

// GetCommentCountByPostID 根据日记ID统计评论数量
// 参数:
//   - postID: 日记ID
// 返回:
//   - int: 评论数量
//   - error: 查询失败时返回错误
func (s *Service) GetCommentCountByPostID(postID int64) (int, error) {
	return s.repo.GetCountByPostID(postID)
}