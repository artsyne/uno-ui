package task

// JobHistory 任务历史记录
// 对应数据库br_job_history表
type JobHistory struct {
	JobId       int64  `form:"job_id" json:"job_id" db:"job_id"`
	JobType     string `form:"job_type" json:"job_type" db:"job_type"`
	Body        string `form:"body" json:"body" db:"body"`
	PlanUTCTime int64  `form:"plan_utc_time" json:"plan_utc_time" db:"plan_utc_time"`
	DoneUTCTime int64  `form:"done_utc_time" json:"done_utc_time" db:"done_utc_time"`
	Result      int8   `form:"result" json:"result" db:"result"`
	TryNum      int8   `form:"try_num" json:"try_num" db:"try_num"`
}

// JobQueue 任务队列
// 对应数据库br_job_queue表（来自user模块）
type JobQueue struct {
	JobId       int64  `form:"job_id" json:"job_id" db:"job_id"`
	JobType     string `form:"job_type" json:"job_type" db:"job_type"`
	Body        string `form:"body" json:"body" db:"body"`
	AddUTCTime  int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	PlanUTCTime int64  `form:"plan_utc_time" json:"plan_utc_time" db:"plan_utc_time"`
	Result      string `form:"result" json:"result" db:"result"`
	TryNum      int8   `form:"try_num" json:"try_num" db:"try_num"`
}

// EmailTask 邮件任务
// 用于解析任务队列中的邮件任务内容
type EmailTask struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
