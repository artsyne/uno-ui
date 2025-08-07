package diary

// RepositoryInterface 日记本数据访问层接口
// 定义所有与日记本相关的数据库操作
type RepositoryInterface interface {
	// Insert 新增日记本到数据库
	Insert(diary *Diary) (int64, error)

	// Update 更新日记本信息
	Update(diary *Diary) error

	// IncrementPostNum 增加日记本的日记数量
	IncrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error

	// DecrementPostNum 减少日记本的日记数量
	DecrementPostNum(diaryID int64, userID int64, userName string, actionUTCTime int64) error

	// GetByDiaryID 根据ID获取日记本对象
	GetByDiaryID(diaryID int64) (*Diary, error)

	// ListForWriteByUserID 根据用户ID获取日记本列表，依据最后更新时间排序
	ListForWriteByUserID(userID int64) ([]Diary, error)

	// ListByUserID 根据特定用户ID检索分页的日记本列表
	ListByUserID(userID int64, pageIndex, pageSize int) ([]Diary, int, error)

	// GetCountByUserID 获取用户的日记本数量
	GetCountByUserID(userID int64) (int, error)
}