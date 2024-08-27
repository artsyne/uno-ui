package models

import (
	"database/sql"

	"github.com/biuren/hebe/pkg"
	"github.com/sirupsen/logrus"
)

type EmailTask struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type JobQueue struct {
	JobId    int64     `form:"job_id" json:"job_id" db:"job_id"`
	JobType  string    `form:"job_type" json:"job_type" db:"job_type"`
	Body     string    `form:"body" json:"body" db:"body"`
	AddUTCTime  int64 `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	PlanUTCTime int64 `form:"plan_utc_time" json:"plan_utc_time" db:"plan_utc_time"`
	Result   int8      `form:"result" json:"result" db:"result"`
	TryNum   int8      `form:"try_num" json:"try_num" db:"try_num"`
}

//插入数据
func (jobQueue *JobQueue) Insert() (int64, error) {
	sql := `
                insert into br_job_queue  (
                    job_type  ,
                    body  ,
                    add_utc_time  ,
                    plan_utc_time  ,
                    result ,
                    try_num
                )
                values (
                    ?  ,
                    ?  ,
                    ?  ,
                    ?  ,
                    ? ,
                    ?
                )
                `
	res, err := db.Exec(sql,
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

//获取任务，获取失败小于3次的任务
func GetJobTasks(num int, jobType string) ([]JobQueue, error) {
	result := []JobQueue{}
	sqlstr := `SELECT *  FROM br_job_queue WHERE job_type = ? and try_num <4
		ORDER BY job_id desc limit ? `
	err := db.Select(&result, sqlstr, jobType, num)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

//删除一个任务
func DelJobTask(job_id int64) (int64, error) {
	sql := "DELETE FROM br_job_queue WHERE job_id = ?"
	res, err := db.Exec(sql, job_id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

//更新任务失败次数
func UpdateJobTaskTryNum(job_id int64) (int64, error) {
	sql := "UPDATE br_job_queue set try_num = try_num + 1  where job_id = ?"
	res, err := db.Exec(sql, job_id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	return rows, err
}

//将任务移到历史表中
func (jobQueue *JobQueue) RemoveToHistory() (int64, error) {
	//删除任务
	rows, err := DelJobTask(jobQueue.JobId)
	if rows != 1 || err != nil {
		logrus.Error(err)
		//任务重试次数加1
		_, _ = UpdateJobTaskTryNum(jobQueue.JobId)
		return 0, err
	}
	//将任务写入历史表
	jobHistory := &JobHistory{}
	jobHistory.JobId = jobQueue.JobId
	jobHistory.Body = jobQueue.Body
	jobHistory.DoneUTCTime = pkg.GetNowTimeUnixMilli()
	jobHistory.JobType = jobQueue.JobType
	jobHistory.PlanUTCTime = jobQueue.PlanUTCTime
	jobHistory.TryNum = jobQueue.TryNum
	jobHistory.Result = 2
	rows, err = jobHistory.Insert()
	return rows, err
}
