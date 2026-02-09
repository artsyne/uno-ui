package entry

// RepositoryInterface 日记数据访问层接口
// 定义所有与日记相关的数据库操作
type RepositoryInterface interface {
	// Insert 插入日记信息到数据库
	// 参数:
	//   - entry: 要插入的日记对象
	// 返回:
	//   - int64: 新插入记录的ID
	//   - error: 插入失败时返回错误
	Insert(entry *Entry) (int64, error)

	// Update 更新日记信息
	// 参数:
	//   - entry: 包含更新信息的日记对象
	// 返回:
	//   - error: 更新失败时返回错误
	Update(entry *Entry) error

	// ListByUserID 根据用户ID和日记本ID分页查询日记信息
	// 参数:
	//   - userID: 用户ID
	//   - diaryID: 日记本ID
	//   - status: 日记状态
	//   - pageIndex: 页码（从0开始）
	//   - pageSize: 每页大小
	// 返回:
	//   - []Entry: 日记列表
	//   - int: 总记录数
	//   - error: 查询失败时返回错误
	ListByUserID(userID int64, diaryID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error)

	// ListByDiaryID 根据日记本ID分页查询日记信息（所有作者）
	// 参数:
	//   - diaryID: 日记本ID
	//   - status: 日记状态
	//   - pageIndex: 页码（从0开始）
	//   - pageSize: 每页大小
	// 返回:
	//   - []Entry: 日记列表
	//   - int: 总记录数
	//   - error: 查询失败时返回错误
	ListByDiaryID(diaryID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error)

	// GetCountByUserID 根据用户ID和日记本ID统计信息数量
	// 参数:
	//   - userID: 用户ID
	//   - diaryID: 日记本ID
	//   - status: 日记状态
	// 返回:
	//   - int: 日记数量
	//   - error: 查询失败时返回错误
	GetCountByUserID(userID int64, diaryID int64, status int8) (int, error)

	// GetCountByDiaryID 根据日记本ID统计信息数量
	// 参数:
	//   - diaryID: 日记本ID
	//   - status: 日记状态
	// 返回:
	//   - int: 日记数量
	//   - error: 查询失败时返回错误
	GetCountByDiaryID(diaryID int64, status int8) (int, error)

	// ListAllByUserID 根据用户ID分页查询日记信息（跨日记本）
	// 参数:
	//   - userID: 用户ID
	//   - status: 日记状态
	//   - pageIndex: 页码（从0开始）
	//   - pageSize: 每页大小
	// 返回:
	//   - []Entry: 日记列表
	//   - int: 总记录数
	//   - error: 查询失败时返回错误
	ListAllByUserID(userID int64, status int8, pageIndex, pageSize int) ([]Entry, int, error)

	// GetByID 根据日记ID获取日记信息
	// 参数:
	//   - entryID: 日记ID
	// 返回:
	//   - *Entry: 日记对象，未找到时返回nil
	//   - error: 查询失败时返回错误
	GetByID(entryID int64) (*Entry, error)

	// GetByIDForUser 根据日记ID和用户ID获取日记信息（包含草稿/回收站状态）
	// 用于自动保存草稿、发布草稿等需要访问非正常状态数据的场景
	// 参数:
	//   - entryID: 日记ID
	//   - userID: 用户ID（用于权限校验）
	// 返回:
	//   - *Entry: 日记对象
	//   - error: 查询失败时返回错误
	GetByIDForUser(entryID int64, userID int64) (*Entry, error)

	// GetByIDRaw 根据日记ID获取日记信息（忽略状态和用户限制）
	// 用于权限检查等需要获取原始日记数据的场景
	// 参数:
	//   - entryID: 日记ID
	// 返回:
	//   - *Entry: 日记对象，未找到时返回nil
	//   - error: 查询失败时返回错误
	GetByIDRaw(entryID int64) (*Entry, error)

	// UpdateByID 根据日记ID更新日记信息
	// 参数:
	//   - entry: 包含更新信息的日记对象
	// 返回:
	//   - error: 更新失败时返回错误
	UpdateByID(entry *Entry) error

	// DeleteByID 根据日记ID删除日记信息（软删除）
	// 参数:
	//   - entryID: 日记ID
	// 返回:
	//   - error: 删除失败时返回错误
	DeleteByID(entryID int64) error

	// HardDeleteByID 根据日记ID物理删除日记信息
	// 参数:
	//   - entryID: 日记ID
	// 返回:
	//   - error: 删除失败时返回错误
	HardDeleteByID(entryID int64) error

	// GetByIDForUser 根据日记ID和用户ID获取日记信息（忽略状态）
	// 参数:
	//   - currentEntryID: 当前日记ID
	//   - diaryID: 日记本ID
	// 返回:
	//   - int64: 上一篇日记ID，未找到时返回0
	//   - error: 查询失败时返回错误
	GetPreviousEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error)

	// GetNextEntryByDiaryID 根据当前日记ID和日记本ID获取下一篇日记ID
	// 参数:
	//   - currentEntryID: 当前日记ID
	//   - diaryID: 日记本ID
	// 返回:
	//   - int64: 下一篇日记ID，未找到时返回0
	//   - error: 查询失败时返回错误
	GetNextEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error)

	// GetStatsByIDs 根据日记ID列表获取统计信息
	// 参数:
	//   - entryIDs: 日记ID列表
	// 返回:
	//   - []Entry: 包含统计信息的日记列表
	//   - error: 查询失败时返回错误
	GetStatsByIDs(entryIDs []int64) ([]Entry, error)

	// IncrementLikeNum 增加日记的点赞数量
	// 参数:
	//   - entryID: 日记ID
	// 返回:
	//   - error: 操作失败时返回错误
	IncrementLikeNum(entryID int64) error

	// DecrementLikeNum 减少日记的点赞数量
	// 参数:
	//   - entryID: 日记ID
	// 返回:
	//   - error: 操作失败时返回错误
	DecrementLikeNum(entryID int64) error
}
