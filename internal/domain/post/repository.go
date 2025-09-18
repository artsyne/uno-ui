package post

// RepositoryInterface 日记数据访问层接口
// 定义所有与日记相关的数据库操作
type RepositoryInterface interface {
	// Insert 插入日记信息到数据库
	// 参数:
	//   - post: 要插入的日记对象
	// 返回:
	//   - int64: 新插入记录的ID
	//   - error: 插入失败时返回错误
	Insert(post *Post) (int64, error)

	// Update 更新日记信息
	// 参数:
	//   - post: 包含更新信息的日记对象
	// 返回:
	//   - error: 更新失败时返回错误
	Update(post *Post) error

	// ListByUserID 根据用户ID和日记本ID分页查询日记信息
	// 参数:
	//   - userID: 用户ID
	//   - diaryID: 日记本ID
	//   - pageIndex: 页码（从0开始）
	//   - pageSize: 每页大小
	// 返回:
	//   - []Post: 日记列表
	//   - int: 总记录数
	//   - error: 查询失败时返回错误
	ListByUserID(userID int64, diaryID int64, pageIndex, pageSize int) ([]Post, int, error)

	// GetCountByUserID 根据用户ID和日记本ID统计信息数量
	// 参数:
	//   - userID: 用户ID
	//   - diaryID: 日记本ID
	// 返回:
	//   - int: 日记数量
	//   - error: 查询失败时返回错误
	GetCountByUserID(userID int64, diaryID int64) (int, error)

	// GetByID 根据日记ID获取日记信息
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - *Post: 日记对象，未找到时返回nil
	//   - error: 查询失败时返回错误
	GetByID(postID int64) (*Post, error)

	// UpdateByID 根据日记ID更新日记信息
	// 参数:
	//   - post: 包含更新信息的日记对象
	// 返回:
	//   - error: 更新失败时返回错误
	UpdateByID(post *Post) error

	// DeleteByID 根据日记ID删除日记信息
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - error: 删除失败时返回错误
	DeleteByID(postID int64) error

	// GetPreviousPostByDiaryID 根据当前日记ID和日记本ID获取上一篇日记ID
	// 参数:
	//   - currentPostID: 当前日记ID
	//   - diaryID: 日记本ID
	// 返回:
	//   - int64: 上一篇日记ID，未找到时返回0
	//   - error: 查询失败时返回错误
	GetPreviousPostByDiaryID(currentPostID int64, diaryID int64) (int64, error)

	// GetNextPostByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
	// 参数:
	//   - currentPostID: 当前日记ID
	//   - diaryID: 日记本ID
	// 返回:
	//   - int64: 下一篇日记ID，未找到时返回0
	//   - error: 查询失败时返回错误
	GetNextPostByDiaryID(currentPostID int64, diaryID int64) (int64, error)

	// GetStatsByIDs 根据日记ID列表获取统计信息
	// 参数:
	//   - postIDs: 日记ID列表
	// 返回:
	//   - []Post: 包含统计信息的日记列表
	//   - error: 查询失败时返回错误
	GetStatsByIDs(postIDs []int64) ([]Post, error)

	// IncrementLikeNum 增加日记的点赞数量
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - error: 操作失败时返回错误
	IncrementLikeNum(postID int64) error

	// DecrementLikeNum 减少日记的点赞数量
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - error: 操作失败时返回错误
	DecrementLikeNum(postID int64) error
}