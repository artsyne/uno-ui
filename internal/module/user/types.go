package user

import (
	"github.com/biuren/hebe/internal/app/render"
)

// AuthUser 已认证用户信息结构体
type AuthUser struct {
	LoginName    string `form:"login_name" json:"login_name"`
	Pwd          string `form:"pwd" json:"pwd"`
	AccessUrl    string `form:"access_url" json:"access_url"`
	UserId       int64  `form:"user_id" json:"user_id"`
	UserName     string `form:"user_name" json:"user_name"`
	FullName     string `form:"full_name" json:"full_name"`
	Email        string `form:"email" json:"email" db:"email"`
	UserGroup    int8   `form:"user_group" json:"user_group"`
	FaceUrl      string `form:"face_url" json:"face_url"`
	LoginUTCTime int64  `form:"login_utc_time" json:"login_utc_time"`
	TimeZone     string `form:"time_zone" json:"time_zone"`
	Lang         string `form:"lang" json:"lang"`
	Status       int8   `form:"status" json:"status"`
}

// 实现 render.RenderUser 接口
func (u *AuthUser) GetLoginName() string   { return u.LoginName }
func (u *AuthUser) GetAccessUrl() string   { return u.AccessUrl }
func (u *AuthUser) GetUserID() int64       { return u.UserId }
func (u *AuthUser) GetUserName() string    { return u.UserName }
func (u *AuthUser) GetFullName() string    { return u.FullName }
func (u *AuthUser) GetEmail() string       { return u.Email }
func (u *AuthUser) GetUserGroup() int8     { return u.UserGroup }
func (u *AuthUser) GetFaceUrl() string     { return u.FaceUrl }
func (u *AuthUser) GetLoginUTCTime() int64 { return u.LoginUTCTime }
func (u *AuthUser) GetTimeZone() string    { return u.TimeZone }
func (u *AuthUser) GetLang() string        { return u.Lang }
func (u *AuthUser) GetStatus() int8        { return u.Status }

// 编译时检查：确保 AuthUser 实现了 render.RenderUser 接口
var _ render.RenderUser = (*AuthUser)(nil)

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
