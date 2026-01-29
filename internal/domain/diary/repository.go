package diary

// RepositoryInterface 日记本数据访问层接口
// 定义所有与日记本相关的数据库操作
type RepositoryInterface interface {
	// Insert 新增日记本到数据库
	// 参数:
	//   - diary: 要插入的日记本对象
	// 返回:
	//   - int64: 新插入记录的ID
	//   - error: 插入失败时返回错误
	Insert(diary *Diary) (int64, error)

	// Update 更新日记本信息
	// 参数:
	//   - diary: 包含更新信息的日记本对象
	// 返回:
	//   - error: 更新失败时返回错误
	Update(diary *Diary) error

	// Delete 删除日记本
	// 参数:
	//   - diaryID: 日记本ID
	// 返回:
	//   - error: 删除失败时返回错误
	Delete(diaryID int64) error

	// IncrementEntryNum 增加日记本的日记数量
	// 参数:
	//   - diaryID: 日记本ID
	// 返回:
	//   - error: 操作失败时返回错误
	IncrementEntryNum(diaryID int64) error

	// DecrementEntryNum 减少日记本的日记数量
	// 参数:
	//   - diaryID: 日记本ID
	// 返回:
	//   - error: 操作失败时返回错误
	DecrementEntryNum(diaryID int64) error

	// GetByDiaryID 根据ID获取日记本对象
	// 参数:
	//   - diaryID: 日记本ID
	// 返回:
	//   - *Diary: 日记本对象，未找到时返回nil
	//   - error: 查询失败时返回错误
	GetByDiaryID(diaryID int64) (*Diary, error)

	// ListForWriteByUserID 根据用户ID获取日记本列表，依据最后更新时间排序
	// 参数:
	//   - userID: 用户ID
	// 返回:
	//   - []Diary: 日记本列表，按最后更新时间排序
	//   - error: 查询失败时返回错误
	ListForWriteByUserID(userID int64) ([]Diary, error)

	// ListByUserID 根据特定用户ID检索分页的日记本列表
	// 参数:
	//   - userID: 用户ID
	//   - pageIndex: 页码（从0开始）
	//   - pageSize: 每页大小
	// 返回:
	//   - []Diary: 日记本列表
	//   - int: 总记录数
	//   - error: 查询失败时返回错误
	ListByUserID(userID int64, pageIndex, pageSize int) ([]Diary, int, error)

	// GetCountByUserID 获取用户的日记本数量
	// 参数:
	//   - userID: 用户ID
	// 返回:
	//   - int: 日记本数量
	//   - error: 查询失败时返回错误
	GetCountByUserID(userID int64) (int, error)

	// IsDiaryOwner 验证用户是否是日记本的拥有者
	// 参数:
	//   - diaryID: 日记本ID
	//   - userID: 用户ID
	// 返回:
	//   - bool: 是否为拥有者
	//   - error: 验证失败时返回错误
	IsDiaryOwner(diaryID, userID int64) (bool, error)

	// IncrementLikeNum 增加日记本的点赞数量
	// 参数:
	//   - diaryID: 日记本ID
	// 返回:
	//   - error: 操作失败时返回错误
	IncrementLikeNum(diaryID int64) error

	// DecrementLikeNum 减少日记本的点赞数量
	// 参数:
	//   - diaryID: 日记本ID
	// 返回:
	//   - error: 操作失败时返回错误
	DecrementLikeNum(diaryID int64) error
}
