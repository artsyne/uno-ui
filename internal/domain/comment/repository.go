package comment

// RepositoryInterface 评论数据访问层接口
// 定义所有与评论相关的数据库操作
type RepositoryInterface interface {
	// Insert 插入评论信息到数据库
	// 参数:
	//   - comment: 要插入的评论对象
	// 返回:
	//   - int64: 新插入记录的ID
	//   - error: 插入失败时返回错误
	Insert(comment *Comment) (int64, error)

	// Delete 删除评论
	// 参数:
	//   - commentID: 评论ID
	// 返回:
	//   - error: 删除失败时返回错误
	Delete(commentID int64) error

	// GetByID 根据评论ID获取评论信息
	// 参数:
	//   - commentID: 评论ID
	// 返回:
	//   - *Comment: 评论对象，未找到时返回nil
	//   - error: 查询失败时返回错误
	GetByID(commentID int64) (*Comment, error)

	// ListByPostID 根据日记ID分页查询评论信息
	// 参数:
	//   - postID: 日记ID
	//   - pageIndex: 页码（从1开始）
	//   - pageSize: 每页大小
	// 返回:
	//   - []Comment: 评论列表
	//   - int: 总记录数
	//   - error: 查询失败时返回错误
	ListByPostID(postID int64, pageIndex, pageSize int) ([]Comment, int, error)

	// GetCountByPostID 根据日记ID统计评论数量
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - int: 评论数量
	//   - error: 查询失败时返回错误
	GetCountByPostID(postID int64) (int, error)

	// IncrementPostCommentNum 增加日记的评论数量
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - error: 更新失败时返回错误
	IncrementPostCommentNum(postID int64) error

	// DecrementPostCommentNum 减少日记的评论数量
	// 参数:
	//   - postID: 日记ID
	// 返回:
	//   - error: 更新失败时返回错误
	DecrementPostCommentNum(postID int64) error
}