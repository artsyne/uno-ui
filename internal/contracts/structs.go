// Package types 共享结构体定义
//
// 这个包定义了在多个模块之间共享的核心数据结构，
// 使用Go 1.25的现代特性，支持泛型和类型安全。
package types

// AuthUser 已认证用户信息结构体
// 包含用户的完整认证信息和个人设置
type AuthUser struct {
	LoginName    string `form:"login_name" json:"login_name"`         // 登录名
	AccessUrl    string `form:"access_url" json:"access_url"`         // 访问URL
	UserId       int64  `form:"user_id" json:"user_id"`               // 用户唯一标识ID
	UserName     string `form:"user_name" json:"user_name"`           // 用户名
	FullName     string `form:"full_name" json:"full_name"`           // 用户全名
	Email        string `form:"email" json:"email" db:"email"`        // 邮箱地址
	UserGroup    int8   `form:"user_group" json:"user_group"`         // 用户组标识
	FaceUrl      string `form:"face_url" json:"face_url"`             // 头像URL
	LoginUTCTime int64  `form:"login_utc_time" json:"login_utc_time"` // 登录时间（UTC时间戳）
	TimeZone     string `form:"time_zone" json:"time_zone"`           // 用户时区设置
	Lang         string `form:"lang" json:"lang"`                     // 用户语言偏好
	Status       int8   `form:"status" json:"status"`                 // 用户账户状态
}

// PageData 页面数据结构 - 核心页面渲染数据容器
//
// 这是整个应用的核心数据结构，用于在控制器和视图之间传递页面渲染所需的所有数据。
// 使用Go 1.25的泛型特性确保类型安全，支持任意类型的业务数据。
//
// 主要职责：
// - 封装页面元信息（标题、描述、关键词）
// - 提供渲染上下文（用户认证状态、语言设置等）
// - 承载业务数据（通过泛型T确保类型安全）
// - 提供便捷的访问方法（如IsAuthenticated、GetAuthUser）
//
// 使用示例：
//
//	data := &PageData[map[string]interface{}]{
//	    Title: "首页",
//	    Context: renderContext,
//	    Data: map[string]interface{}{"entries": entries},
//	}
type PageData[T any] struct {
	// 页面元信息
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`

	// 渲染上下文
	Context *RenderContext `json:"context"`

	// 业务数据 - 使用泛型确保类型安全
	Data T `json:"data"`
}

// IsAuthenticated 检查用户是否已认证
// 返回:
//   - bool: true表示用户已认证，false表示未认证
func (p *PageData[T]) IsAuthenticated() bool {
	if p.Context == nil {
		return false
	}
	return p.Context.IsLogin
}

// GetAuthUser 获取认证用户信息
// 返回:
//   - *AuthUser: 认证用户信息，未认证时返回nil
func (p *PageData[T]) GetAuthUser() *AuthUser {
	if p.Context == nil {
		return nil
	}
	return p.Context.AuthUser
}

// RenderContext 渲染上下文
// 包含页面渲染所需的上下文信息
type RenderContext struct {
	IsLogin bool `json:"is_login"` // 用户是否登录

	// AuthUser 当前登录的用户信息
	AuthUser *AuthUser `json:"auth_user,omitempty"`

	// Lang 当前页面使用的语言代码
	Lang string `json:"lang"`

	// Theme 页面主题设置
	Theme string `json:"theme,omitempty"`
	// Layout 页面布局设置
	IsHXRequest bool `json:"is_hx_request"`

	// CsrfToken CSRF令牌
	CsrfToken string `json:"csrf_token"`
}
