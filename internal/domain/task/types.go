package task

// JobHistory 任务历史记录
// 对应数据库br_job_history表
type JobHistory struct {
	JobId       int64  `form:"job_id" json:"job_id" db:"job_id"`                   // 任务ID
	JobType     string `form:"job_type" json:"job_type" db:"job_type"`             // 任务类型
	Body        string `form:"body" json:"body" db:"body"`                         // 任务内容（JSON格式）
	PlanUTCTime int64  `form:"plan_utc_time" json:"plan_utc_time" db:"plan_utc_time"` // 计划执行时间（UTC时间戳）
	DoneUTCTime int64  `form:"done_utc_time" json:"done_utc_time" db:"done_utc_time"` // 实际完成时间（UTC时间戳）
	Result      int8   `form:"result" json:"result" db:"result"`                   // 执行结果（0:失败, 1:成功）
	TryNum      int8   `form:"try_num" json:"try_num" db:"try_num"`                 // 尝试次数
}

// JobQueue 任务队列
// 对应数据库br_job_queue表（来自user模块）
type JobQueue struct {
	JobId       int64  `form:"job_id" json:"job_id" db:"job_id"`                   // 任务ID
	JobType     string `form:"job_type" json:"job_type" db:"job_type"`             // 任务类型
	Body        string `form:"body" json:"body" db:"body"`                         // 任务内容（JSON格式）
	AddUTCTime  int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`   // 添加时间（UTC时间戳）
	PlanUTCTime int64  `form:"plan_utc_time" json:"plan_utc_time" db:"plan_utc_time"` // 计划执行时间（UTC时间戳）
	Result      int8   `form:"result" json:"result" db:"result"`                   // 执行结果（0:失败, 1:成功）
	TryNum      int8   `form:"try_num" json:"try_num" db:"try_num"`                 // 尝试次数
}

// EmailTask 邮件任务
// 用于解析任务队列中的邮件任务内容
type EmailTask struct {
	To      string `json:"to"`      // 收件人邮箱地址
	Subject string `json:"subject"` // 邮件主题
	Body    string `json:"body"`    // 邮件正文内容
}