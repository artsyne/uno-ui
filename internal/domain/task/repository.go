package task

// RepositoryInterface 任务数据访问层接口
// 定义所有与任务相关的数据库操作
type RepositoryInterface interface {
	// InsertJobHistory 插入任务历史记录
	// 参数:
	//   - jobHistory: 要插入的任务历史记录对象
	// 返回:
	//   - int64: 新插入记录的ID
	//   - error: 插入失败时返回错误
	InsertJobHistory(jobHistory *JobHistory) (int64, error)

	// InsertJobQueue 插入作业队列
	// 参数:
	//   - jobQueue: 要插入的作业队列对象
	// 返回:
	//   - int64: 新插入记录的ID
	//   - error: 插入失败时返回错误
	InsertJobQueue(jobQueue *JobQueue) (int64, error)

	// GetJobTasks 获取作业任务
	// 参数:
	//   - num: 获取任务的数量
	//   - jobType: 任务类型
	// 返回:
	//   - []JobQueue: 任务队列列表
	//   - error: 查询失败时返回错误
	GetJobTasks(num int, jobType string) ([]JobQueue, error)

	// DeleteJobTask 删除作业任务
	// 参数:
	//   - jobId: 任务ID
	// 返回:
	//   - int64: 受影响的行数
	//   - error: 删除失败时返回错误
	DeleteJobTask(jobId int64) (int64, error)

	// UpdateJobTaskTryNum 更新任务失败次数
	// 参数:
	//   - jobId: 任务ID
	// 返回:
	//   - int64: 受影响的行数
	//   - error: 更新失败时返回错误
	UpdateJobTaskTryNum(jobId int64) (int64, error)

	// RemoveJobToHistory 将任务移到历史表中
	// 参数:
	//   - jobQueue: 要移动的任务队列对象
	// 返回:
	//   - int64: 受影响的行数
	//   - error: 操作失败时返回错误
	RemoveJobToHistory(jobQueue *JobQueue) (int64, error)
}