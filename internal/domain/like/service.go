package like

// Service 点赞业务逻辑层接口
type Service interface {
	// ToggleLike 切换点赞状态
	ToggleLike(userID, ownerID int64, likeType int) error

	// GetLikeStatus 获取点赞状态信息
	GetLikeStatus(userID, ownerID int64, likeType int) (*LikeStatus, error)

	// GetLikeCount 获取点赞总数
	GetLikeCount(ownerID int64, likeType int) (int, error)
}

// LikeStatus 点赞状态信息
type LikeStatus struct {
	IsLiked   bool `json:"is_liked"`   // 是否已点赞
	LikeCount int  `json:"like_count"` // 点赞总数
}

// LikeService 点赞服务实现
type LikeService struct {
	repo Repository
}

// NewLikeService 创建点赞服务实例
func NewLikeService(repo Repository) Service {
	return &LikeService{
		repo: repo,
	}
}

// ToggleLike 切换点赞状态
func (s *LikeService) ToggleLike(userID, ownerID int64, likeType int) error {
	return s.repo.ToggleLike(userID, ownerID, likeType)
}

// GetLikeStatus 获取点赞状态信息
func (s *LikeService) GetLikeStatus(userID, ownerID int64, likeType int) (*LikeStatus, error) {
	// 检查是否已点赞
	isLiked, err := s.repo.IsLiked(userID, ownerID, likeType)
	if err != nil {
		return nil, err
	}

	// 获取点赞总数
	likeCount, err := s.repo.GetLikeCount(ownerID, likeType)
	if err != nil {
		return nil, err
	}

	return &LikeStatus{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}

// GetLikeCount 获取点赞总数
func (s *LikeService) GetLikeCount(ownerID int64, likeType int) (int, error) {
	return s.repo.GetLikeCount(ownerID, likeType)
}