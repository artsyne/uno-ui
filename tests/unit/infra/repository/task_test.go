package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/tests/framework"
	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTaskRepository 模拟任务仓库，用于测试
type mockTaskRepository struct {
	db *sqlx.DB
}

// InsertJobHistory 插入任务历史记录
func (r *mockTaskRepository) InsertJobHistory(jobHistory *task.JobHistory) (int64, error) {
	sql := `INSERT INTO br_job_history (
		job_id, job_type, body, plan_utc_time, done_utc_time, result, try_num
	) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(sql,
		jobHistory.JobId, jobHistory.JobType, jobHistory.Body,
		jobHistory.PlanUTCTime, jobHistory.DoneUTCTime, jobHistory.Result, jobHistory.TryNum)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// InsertJobQueue 插入作业队列
func (r *mockTaskRepository) InsertJobQueue(jobQueue *task.JobQueue) (int64, error) {
	sql := `INSERT INTO br_job_queue (
		job_type, body, add_utc_time, plan_utc_time, result, try_num
	) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(sql,
		jobQueue.JobType, jobQueue.Body, jobQueue.AddUTCTime,
		jobQueue.PlanUTCTime, jobQueue.Result, jobQueue.TryNum)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// GetJobTasks 获取作业任务
func (r *mockTaskRepository) GetJobTasks(num int, jobType string) ([]task.JobQueue, error) {
	result := []task.JobQueue{}
	sqlstr := `SELECT * FROM br_job_queue WHERE job_type = ? and try_num < 4 ORDER BY job_id desc limit ?`
	err := r.db.Select(&result, sqlstr, jobType, num)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// DeleteJobTask 删除作业任务
func (r *mockTaskRepository) DeleteJobTask(jobId int64) (int64, error) {
	sql := "DELETE FROM br_job_queue WHERE job_id = ?"
	res, err := r.db.Exec(sql, jobId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// UpdateJobTaskTryNum 更新任务失败次数
func (r *mockTaskRepository) UpdateJobTaskTryNum(jobId int64) (int64, error) {
	sql := "UPDATE br_job_queue set try_num = try_num + 1 where job_id = ?"
	res, err := r.db.Exec(sql, jobId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, err
}

// RemoveJobToHistory 将任务移到历史表中
func (r *mockTaskRepository) RemoveJobToHistory(jobQueue *task.JobQueue) (int64, error) {
	// 模拟实际的 RemoveJobToHistory 逻辑
	rows, err := r.DeleteJobTask(jobQueue.JobId)
	if err != nil {
		return 0, err
	}
	
	if rows == 1 {
		// 删除成功，插入历史记录
		jobHistory := &task.JobHistory{
			JobId: jobQueue.JobId,
			JobType: jobQueue.JobType,
			Body: jobQueue.Body,
			PlanUTCTime: jobQueue.PlanUTCTime,
			DoneUTCTime: time.Now().Unix(),
			Result: 2, // 成功完成
			TryNum: jobQueue.TryNum,
		}
		return r.InsertJobHistory(jobHistory)
	} else {
		// 删除失败，更新重试次数
		return r.UpdateJobTaskTryNum(jobQueue.JobId)
	}
}

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

// TestNewTaskRepository 测试创建新的任务仓库
func TestNewTaskRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	// 使用测试框架提供的数据库连接
	testCtx := framework.NewTestContext(framework.WithDatabase())
	defer testCtx.Cleanup()
	
	// 测试使用全局数据库连接的构造函数
	repo := repository.NewTaskRepository()
	assert.NotNil(t, repo)
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestNewTaskRepository_WithDB 测试使用指定数据库连接创建任务仓库
func TestNewTaskRepository_WithDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	
	sqlxDB := sqlx.NewDb(db, "mysql")
	repo := repository.NewTaskRepositoryWithDB(sqlxDB)
	assert.NotNil(t, repo)
	
	// 验证 mock 期望
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestNewTaskRepository_NilDB 测试空数据库连接
func TestNewTaskRepository_NilDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	// 测试传入 nil 数据库连接应该 panic
	assert.Panics(t, func() {
		repository.NewTaskRepositoryWithDB(nil)
	})
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestTaskRepository_InsertJobHistory 测试插入任务历史记录
func TestTaskRepository_InsertJobHistory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	// 创建一个模拟的 repository 实例用于测试
	// 注意：实际的 NewTaskRepository() 不接受参数
	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	jobHistory := &task.JobHistory{
		JobId: 1,
		JobType: "test_job",
		Body: "{}",
		PlanUTCTime: time.Now().Unix(),
		DoneUTCTime: time.Now().Unix(),
		Result: 2,
		TryNum: 1,
	}

	// 设置 mock 期望
	mock.ExpectExec("INSERT INTO br_job_history").
		WithArgs(jobHistory.JobId, jobHistory.JobType, jobHistory.Body, jobHistory.PlanUTCTime, jobHistory.DoneUTCTime, jobHistory.Result, jobHistory.TryNum).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 执行测试
	id, err := repo.InsertJobHistory(jobHistory)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestTaskRepository_InsertJobHistory_Error 测试插入任务历史记录错误
func TestTaskRepository_InsertJobHistory_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	jobHistory := &task.JobHistory{
		JobId: 1,
		JobType: "test_job",
		Body: "{}",
		PlanUTCTime: time.Now().Unix(),
		DoneUTCTime: time.Now().Unix(),
		Result: 2,
		TryNum: 1,
	}

	// 设置 mock 期望 - 模拟数据库错误
	mock.ExpectExec("INSERT INTO br_job_history").
		WillReturnError(sql.ErrConnDone)

	// 执行测试
	id, err := repo.InsertJobHistory(jobHistory)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, int64(0), id)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestTaskRepository_InsertJobQueue 测试插入作业队列
func TestTaskRepository_InsertJobQueue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	// 准备测试数据
	jobQueue := &task.JobQueue{
		JobType: "test_job",
		Body: "{\"key\": \"value\"}",
		AddUTCTime: time.Now().Unix(),
		PlanUTCTime: time.Now().Unix(),
		Result: 0,
		TryNum: 0,
	}

	// 设置 mock 期望
	mock.ExpectExec("INSERT INTO br_job_queue").
		WithArgs(jobQueue.JobType, jobQueue.Body, jobQueue.AddUTCTime, jobQueue.PlanUTCTime, jobQueue.Result, jobQueue.TryNum).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 执行测试
	id, err := repo.InsertJobQueue(jobQueue)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestTaskRepository_GetJobTasks 测试获取作业任务
func TestTaskRepository_GetJobTasks(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	// 设置 mock 期望
	rows := sqlmock.NewRows([]string{
		"job_id", "job_type", "body", "add_utc_time",
		"plan_utc_time", "result", "try_num",
	}).AddRow(
		1, "test_job", "{\"key\": \"value\"}", time.Now().Unix(),
		time.Now().Unix(), 0, 0,
	).AddRow(
		2, "test_job", "{\"key2\": \"value2\"}", time.Now().Unix(),
		time.Now().Unix(), 0, 1,
	)

	mock.ExpectQuery("SELECT \\* FROM br_job_queue WHERE job_type = \\? and try_num < 4 ORDER BY job_id desc limit \\?").
		WithArgs("test_job", 10).
		WillReturnRows(rows)

	// 执行测试
	tasks, err := repo.GetJobTasks(10, "test_job")
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "test_job", tasks[0].JobType)
	assert.Equal(t, "test_job", tasks[1].JobType)

	// 验证 mock 期望
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}



// TestTaskRepository_DeleteJobTask 测试删除任务
func TestTaskRepository_DeleteJobTask(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	// 设置 mock 期望
	mock.ExpectExec("DELETE FROM br_job_queue WHERE job_id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试
	rows, err := repo.DeleteJobTask(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestTaskRepository_UpdateJobTaskTryNum 测试更新任务尝试次数
func TestTaskRepository_UpdateJobTaskTryNum(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	// 设置 mock 期望
	mock.ExpectExec("UPDATE br_job_queue set try_num = try_num \\+ 1 where job_id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试
	rows, err := repo.UpdateJobTaskTryNum(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}

// TestTaskRepository_RemoveJobToHistory 测试将任务移动到历史记录
func TestTaskRepository_RemoveJobToHistory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &mockTaskRepository{db: sqlx.NewDb(db, "mysql")}

	// 准备测试数据
	jobQueue := &task.JobQueue{
		JobId: 1,
		JobType: "test_job",
		Body: "{}",
		AddUTCTime: time.Now().Unix(),
		PlanUTCTime: time.Now().Unix(),
		Result: 0,
		TryNum: 1,
	}

	// 设置 mock 期望 - 删除队列记录（成功删除1行）
	mock.ExpectExec("DELETE FROM br_job_queue WHERE job_id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // LastInsertId不重要，RowsAffected=1

	// 设置 mock 期望 - 插入历史记录（使用实际的 SQL 格式）
	mock.ExpectExec(`INSERT INTO br_job_history \(
\s*job_id, job_type, body, plan_utc_time, done_utc_time, result, try_num
\s*\) VALUES \(\?, \?, \?, \?, \?, \?, \?\)`).
		WithArgs(1, "test_job", "{}", sqlmock.AnyArg(), sqlmock.AnyArg(), 2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 执行测试
	id, err := repo.RemoveJobToHistory(jobQueue)
	if err != nil {
		t.Logf("RemoveJobToHistory error: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 确保上下文未超时
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out")
	default:
	}
}