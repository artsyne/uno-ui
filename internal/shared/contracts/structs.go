// Package types 共享结构体定义
//
// 这个包定义了在多个模块之间共享的核心数据结构，
// 使用Go 1.24的现代特性，支持泛型和类型安全。
package types

// PageData 页面数据结构 - 核心页面渲染数据容器
//
// 这是整个应用的核心数据结构，用于在控制器和视图之间传递页面渲染所需的所有数据。
// 使用Go 1.24的泛型特性确保类型安全，支持任意类型的业务数据。
//
// 主要职责：
// - 封装页面元信息（标题、描述、关键词）
// - 提供渲染上下文（用户认证状态、语言设置等）
// - 承载业务数据（通过泛型T确保类型安全）
// - 提供便捷的访问方法（如IsAuthenticated、GetAuthUser）
//
// 使用示例：
//   data := &PageData[map[string]interface{}]{
//       Title: "首页",
//       Context: renderContext,
//       Data: map[string]interface{}{"posts": posts},
//   }
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
func (p *PageData[T]) IsAuthenticated() bool {
	if p.Context == nil {
		return false
	}
	return p.Context.IsLogin
}

// GetAuthUser 获取认证用户信息
func (p *PageData[T]) GetAuthUser() *AuthUser {
	if p.Context == nil {
		return nil
	}
	return p.Context.AuthUser
}

// RenderContext 渲染上下文
type RenderContext struct {
	IsLogin bool `json:"is_login"` // 用户是否登录

	// AuthUser 当前登录的用户信息
	AuthUser *AuthUser `json:"auth_user,omitempty"`

	// Lang 当前页面使用的语言代码
	Lang string `json:"lang"`

	// 请求上下文
	RequestID string `json:"request_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	IP        string `json:"ip,omitempty"`

	// 主题和布局
	Theme  string `json:"theme,omitempty"`
	Layout string `json:"layout,omitempty"`

	// 扩展属性
	Extensions map[string]any `json:"extensions,omitempty"`
}



// ResponseData 通用响应数据结构
type ResponseData[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
}