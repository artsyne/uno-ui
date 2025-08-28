package repository

import (
	"database/sql"
	"log/slog"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/pkg"
	"github.com/jmoiron/sqlx"
)

// taskRepository 任务数据访问层实现
// 负责处理所有与任务相关的数据库操作
type taskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository 创建新的任务数据访问层实例
func NewTaskRepository() task.RepositoryInterface {
	db := database.GetDB()
	if db == nil {
		panic("database connection is nil")
	}
	return &taskRepository{db: db}
}

// NewTaskRepositoryWithDB 创建新的任务数据访问层实例（用于测试）
func NewTaskRepositoryWithDB(db *sqlx.DB) task.RepositoryInterface {
	if db == nil {
		panic("database connection is nil")
	}
	return &taskRepository{db: db}
}

// ===================== JobHistory 相关方法 =====================

// InsertJobHistory 插入任务历史记录
func (r *taskRepository) InsertJobHistory(jobHistory *task.JobHistory) (int64, error) {
	sql := `
		INSERT INTO br_job_history (
			job_id,
			job_type,
			body,
			plan_utc_time,
			done_utc_time,
			result,
			try_num
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(sql,
		jobHistory.JobId,
		jobHistory.JobType,
		jobHistory.Body,
		jobHistory.PlanUTCTime,
		jobHistory.DoneUTCTime,
		jobHistory.Result,
		jobHistory.TryNum,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// ===================== JobQueue 相关方法 =====================

// InsertJobQueue 插入作业队列
func (r *taskRepository) InsertJobQueue(jobQueue *task.JobQueue) (int64, error) {
	sql := `
		INSERT INTO br_job_queue (
			job_type,
			body,
			add_utc_time,
			plan_utc_time,
			result,
			try_num
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(sql,
		jobQueue.JobType,
		jobQueue.Body,
		jobQueue.AddUTCTime,
		jobQueue.PlanUTCTime,
		jobQueue.Result,
		jobQueue.TryNum,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// GetJobTasks 获取作业任务
func (r *taskRepository) GetJobTasks(num int, jobType string) ([]task.JobQueue, error) {
	result := []task.JobQueue{}
	sqlstr := `SELECT * FROM br_job_queue WHERE job_type = ? and try_num < 4
		ORDER BY job_id desc limit ?`
	err := r.db.Select(&result, sqlstr, jobType, num)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// DeleteJobTask 删除作业任务
func (r *taskRepository) DeleteJobTask(jobId int64) (int64, error) {
	sql := "DELETE FROM br_job_queue WHERE job_id = ?"
	res, err := r.db.Exec(sql, jobId)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateJobTaskTryNum 更新任务失败次数
func (r *taskRepository) UpdateJobTaskTryNum(jobId int64) (int64, error) {
	sql := "UPDATE br_job_queue set try_num = try_num + 1 where job_id = ?"
	res, err := r.db.Exec(sql, jobId)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

// RemoveJobToHistory 将任务移到历史表中
func (r *taskRepository) RemoveJobToHistory(jobQueue *task.JobQueue) (int64, error) {
	// 删除任务
	rows, err := r.DeleteJobTask(jobQueue.JobId)
	if rows != 1 || err != nil {
		slog.Error("删除任务失败", "error", err, "job_id", jobQueue.JobId)
		// 任务重试次数加1
		_, _ = r.UpdateJobTaskTryNum(jobQueue.JobId)
		return 0, err
	}

	// 将任务写入历史表
	jobHistory := &task.JobHistory{
		JobId:       jobQueue.JobId,
		JobType:     jobQueue.JobType,
		Body:        jobQueue.Body,
		PlanUTCTime: jobQueue.PlanUTCTime,
		DoneUTCTime: pkg.GetNowTimeUnixMilli(),
		Result:      2,
		TryNum:      jobQueue.TryNum,
	}

	return r.InsertJobHistory(jobHistory)
}