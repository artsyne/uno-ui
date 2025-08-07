package task

// RepositoryInterface 任务数据访问层接口
// 定义所有与任务相关的数据库操作
type RepositoryInterface interface {
	// InsertJobHistory 插入任务历史记录
	InsertJobHistory(jobHistory *JobHistory) (int64, error)

	// InsertJobQueue 插入作业队列
	InsertJobQueue(jobQueue *JobQueue) (int64, error)

	// GetJobTasks 获取作业任务
	GetJobTasks(num int, jobType string) ([]JobQueue, error)

	// DeleteJobTask 删除作业任务
	DeleteJobTask(jobId int64) (int64, error)

	// UpdateJobTaskTryNum 更新任务失败次数
	UpdateJobTaskTryNum(jobId int64) (int64, error)

	// RemoveJobToHistory 将任务移到历史表中
	RemoveJobToHistory(jobQueue *JobQueue) (int64, error)
}