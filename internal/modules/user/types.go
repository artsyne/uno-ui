package user

// Login 登录信息结构体
type Login struct {
	LoginName    string `form:"login_name" json:"login_name"`
	Pwd          string `form:"pwd" json:"pwd"`
	AccessUrl    string `form:"access_url" json:"access_url"`
	UserId       int64  `form:"user_id" json:"user_id"`
	UserName     string `form:"user_name" json:"user_name"`
	FullName     string `form:"full_name" json:"full_name"`
	UserGroup    int8   `form:"user_group" json:"user_group"`
	FaceSmallUrl string `form:"face_small_url" json:"face_small_url"`
	LoginUTCTime int64  `form:"login_utc_time" json:"login_utc_time"`
	TimeZone     string `form:"time_zone" json:"time_zone"`
	Lang         string `form:"lang" json:"lang"`
	Status       int8   `form:"status" json:"status"`
}

// User 用户信息结构体
type User struct {
	UserId           int64  `form:"user_id" json:"user_id" db:"user_id"`
	UserName         string `form:"user_name" json:"user_name" db:"user_name"`
	FullName         string `form:"full_name" json:"full_name" db:"full_name"`
	UserGroup        int8   `form:"user_group" json:"user_group" db:"user_group"`
	Pwd              string `form:"pwd" json:"pwd" db:"pwd"`
	PwdSalt          string `form:"pwd_salt" json:"pwd_salt" db:"pwd_salt"`
	ProfileUrl       string `form:"profile_url" json:"profile_url" db:"profile_url"`
	FaceUrl          string `form:"face_url" json:"face_url" db:"face_url"`
	Email            string `form:"email" json:"email" db:"email"`
	EmailVerified    int8   `form:"email_verified" json:"email_verified" db:"email_verified"`
	RegIp            string `form:"reg_ip" json:"reg_ip" db:"reg_ip"`
	CreatedUTCTime   int64  `form:"created_utc_time" json:"created_utc_time" db:"created_utc_time"`
	Remark           string `form:"remark" json:"remark" db:"remark"`
	LastIp           string `form:"last_ip" json:"last_ip" db:"last_ip"`
	LastLoginUTCTime int64  `form:"last_login_utc_time" json:"last_login_utc_time" db:"last_login_utc_time"`
	TimeZone         string `form:"time_zone" json:"time_zone" db:"time_zone"`
	Lang             string `form:"lang" json:"lang" db:"lang"`
	Status           int8   `form:"status" json:"status" db:"status"`
}

// ValidCode 验证码结构体
type ValidCode struct {
	ValidCodeId   int64  `form:"valid_code_id" json:"valid_code_id" db:"valid_code_id"`
	UserId        int64  `form:"user_id" json:"user_id" db:"user_id"`
	ValidCode     string `form:"valid_code" json:"valid_code" db:"valid_code"`
	ValidTypeId   int8   `form:"valid_type_id" json:"valid_type_id" db:"valid_type_id"`
	ExpireUTCTime int64  `form:"expire_utc_time" json:"expire_utc_time" db:"expire_utc_time"`
	CreateUTCTime int64  `form:"create_utc_time" json:"gmt_create" db:"create_utc_time"`
}

// JobQueue 作业队列结构体
type JobQueue struct {
	JobId       int64  `form:"job_id" json:"job_id" db:"job_id"`
	JobType     string `form:"job_type" json:"job_type" db:"job_type"`
	Body        string `form:"body" json:"body" db:"body"`
	AddUTCTime  int64  `form:"add_utc_time" json:"add_utc_time" db:"add_utc_time"`
	PlanUTCTime int64  `form:"plan_utc_time" json:"plan_utc_time" db:"plan_utc_time"`
	Result      int8   `form:"result" json:"result" db:"result"`
	TryNum      int8   `form:"try_num" json:"try_num" db:"try_num"`
}

// EmailTask 邮件任务结构体
type EmailTask struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// 请求结构体定义
type LoginRequest struct {
	LoginName string `form:"login_name" json:"login_name"`
	Pwd       string `form:"pwd" json:"pwd"`
	AccessUrl string `form:"access_url" json:"access_url"`
}

type RegisterRequest struct {
	UserName string `form:"user_name" json:"user_name"`
	FullName string `form:"full_name" json:"full_name"`
	Email    string `form:"email" json:"email"`
	Pwd      string `form:"pwd" json:"pwd"`
	TimeZone string `form:"time_zone" json:"time_zone"`
	Lang     string `form:"lang" json:"lang"`
	LastIp   string `form:"last_ip" json:"last_ip"`
}

type ResetPasswordRequest struct {
	UID    string `form:"uid" json:"uid"`
	Sign   string `form:"sign" json:"sign"`
	NewPwd string `form:"new_pwd" json:"new_pwd"`
}

type CheckUserNameRequest struct {
	UserName string `form:"username" json:"username"`
}

type CheckEmailRequest struct {
	Email string `form:"email" json:"email"`
}

type FindPasswordRequest struct {
	Email string `form:"email" json:"email"`
}

type ActivateUserParams struct {
	UID  string `form:"uid" json:"uid"`
	Sign string `form:"sign" json:"sign"`
}
