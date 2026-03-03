package like

// Service 点赞业务逻辑层接口
type Service interface {
	// ToggleLike 切换点赞状态
	ToggleLike(userID, ownerID int64, likeType int) error

	// IsLiked 检查用户是否已点赞
	IsLiked(userID, ownerID int64, likeType int) (bool, error)

	// GetLikeStatus 获取点赞状态（只查询 IsLiked，LikeCount 由调用方从 owner 对象获取）
	// 注意：推荐使用 IsLiked 方法 + owner 对象的 likeNum 字段
	GetLikeStatus(userID, ownerID int64, likeType int) (*LikeStatus, error)

	// GetLikeCount 获取点赞总数（从 owner 对象获取，不查询 br_like 表）
	// 注意：此方法已废弃，推荐通过 owner 对象的 likeNum 字段获取
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

// IsLiked 检查用户是否已点赞
func (s *LikeService) IsLiked(userID, ownerID int64, likeType int) (bool, error) {
	return s.repo.IsLiked(userID, ownerID, likeType)
}

// GetLikeCount 获取点赞总数（已废弃，推荐通过 owner 对象的 likeNum 字段获取）
// 注意：此方法保留仅为兼容性，不推荐使用
func (s *LikeService) GetLikeCount(ownerID int64, likeType int) (int, error) {
	return s.repo.GetLikeCount(ownerID, likeType)
}

// GetLikeStatus 获取点赞状态（只查询 IsLiked，LikeCount 由调用方从 owner 对象获取）
// 注意：推荐使用 IsLiked 方法 + owner 对象的 likeNum 字段
func (s *LikeService) GetLikeStatus(userID, ownerID int64, likeType int) (*LikeStatus, error) {
	isLiked, err := s.repo.IsLiked(userID, ownerID, likeType)
	if err != nil {
		return nil, err
	}
	// LikeCount 由调用方填充
	return &LikeStatus{
		IsLiked:   isLiked,
		LikeCount: 0,
	}, nil
}