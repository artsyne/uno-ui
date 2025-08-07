package post

// RepositoryInterface 日记数据访问层接口
// 定义所有与日记相关的数据库操作
type RepositoryInterface interface {
	// Insert 插入日记信息到数据库
	Insert(post *Post) (int64, error)

	// Update 更新日记信息
	Update(post *Post) error

	// ListByUserID 根据用户ID和日记本ID分页查询日记信息
	ListByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error)

	// GetCountByUserID 根据用户ID和日记本ID统计信息数量
	GetCountByUserID(userID int64, diaryID int64) (int, error)

	// GetByID 根据日记ID获取日记信息
	GetByID(postID int64) (*Post, error)

	// UpdateByID 根据日记ID更新日记信息
	UpdateByID(post *Post) error

	// DeleteByID 根据日记ID删除日记信息
	DeleteByID(postID int64) error
}