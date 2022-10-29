package models

import (
	"time"
)

type JobHistory struct {
	JobId    int64     `form:"job_id" json:"job_id" db:"job_id"`
	JobType  string    `form:"job_type" json:"job_type" db:"job_type"`
	Body     string    `form:"body" json:"body" db:"body"`
	PlanTime time.Time `form:"plan_time" json:"plan_time" db:"plan_time"`
	DoneTime time.Time `form:"done_time" json:"done_time" db:"done_time"`
	Result   int8      `form:"result" json:"result" db:"result"`
	TryNum   int8      `form:"try_num" json:"try_num" db:"try_num"`
}

//插入数据
func (jobHistory *JobHistory) Insert() (int64, error) {
	sql := `
                insert into br_job_history  (
                    job_id  ,
                    job_type  ,
                    body  ,
                    plan_time  ,
                    done_time  ,
                    result ,
                    try_num
                )
                values (
                    ?  ,
                    ?  ,
                    ?  ,
                    ?  ,
                    ?  ,
                    ? ,
                    ?
                )
                `
	res, err := db.Exec(sql,
		jobHistory.JobId,
		jobHistory.JobType,
		jobHistory.Body,
		jobHistory.PlanTime,
		jobHistory.DoneTime,
		jobHistory.Result,
		jobHistory.TryNum,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

//删除数据
