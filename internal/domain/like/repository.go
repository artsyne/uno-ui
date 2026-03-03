package like

// Repository 点赞数据访问层接口
type Repository interface {
	// ToggleLike 切换点赞状态（如果已点赞则取消，未点赞则添加）
	ToggleLike(userID, ownerID int64, likeType int) error

	// IsLiked 检查用户是否已点赞
	IsLiked(userID, ownerID int64, likeType int) (bool, error)

	// GetLikeCount 获取点赞总数（从 owner 对象获取，不查询 br_like 表）
	GetLikeCount(ownerID int64, likeType int) (int, error)

}