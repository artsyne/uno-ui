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
}