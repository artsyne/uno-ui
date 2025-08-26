package user

// LoginRequest 登录请求
type LoginRequest struct {
	LoginName string `form:"login_name" json:"login_name" binding:"required,min=3,max=50"`
	Pwd       string `form:"pwd" json:"pwd" binding:"required,min=6"`
	AccessUrl string `form:"access_url" json:"access_url"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,username"`
	FullName string `form:"full_name" json:"full_name" binding:"required,min=2,max=50"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Pwd      string `form:"pwd" json:"pwd" binding:"required,password"`
	TimeZone string `form:"time_zone" json:"time_zone"`
	Lang     string `form:"lang" json:"lang"`
	LastIp   string `form:"last_ip" json:"last_ip"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	UID    string `form:"uid" json:"uid" binding:"required"`
	Sign   string `form:"sign" json:"sign" binding:"required"`
	NewPwd string `form:"new_pwd" json:"new_pwd" binding:"required,password"`
}

// CheckUserNameRequest 检查用户名请求
type CheckUserNameRequest struct {
	UserName string `form:"username" json:"username" binding:"required,username"`
}

// CheckEmailRequest 检查邮箱请求
type CheckEmailRequest struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}

// FindPasswordRequest 找回密码请求
type FindPasswordRequest struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	FullName string `form:"full_name" json:"full_name" binding:"required,min=2,max=50"`
	Remark   string `form:"remark" json:"remark"`
	TimeZone string `form:"time_zone" json:"time_zone" binding:"required"`
	Lang     string `form:"lang" json:"lang" binding:"required"`
}

// UpdateEmailRequest 更新邮箱请求
type UpdateEmailRequest struct {
	UserId int64  `form:"user_id" json:"user_id"`
	Email  string `form:"email" json:"email" binding:"required,email"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	CurrentPwd string `form:"current_pwd" json:"current_pwd" binding:"required,min=6"`
	NewPwd     string `form:"new_pwd" json:"new_pwd" binding:"required,password"`
	ConfirmPwd string `form:"confirm_pwd" json:"confirm_pwd" binding:"required,eqfield=NewPwd"`
}
// ActivateUserRequest 激活用户请求
type ActivateUserRequest struct {
	UID  string `form:"uid" json:"uid" binding:"required"`
	Sign string `form:"sign" json:"sign" binding:"required"`
}
