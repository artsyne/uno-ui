// Package types 渲染层共享类型定义
//
// 这个包定义了渲染层中需要在多个包之间共享的核心类型和接口，
// 主要用于避免循环导入问题。
package types

// RenderUser 渲染用户接口 - 抽象 AuthUser 除 Pwd 外的所有属性
type RenderUser interface {
	GetLoginName() string
	GetAccessUrl() string
	GetUserID() int64
	GetUserName() string
	GetFullName() string
	GetEmail() string
	GetUserGroup() int8
	GetFaceUrl() string
	GetLoginUTCTime() int64
	GetTimeZone() string
	GetLang() string
	GetStatus() int8
}

// RenderContext 渲染上下文
// 包含页面渲染所需的普适性、非业务相关的状态信息。
type RenderContext struct {
	IsLogin bool `json:"is_login"` // 用户是否登录

	// AuthUser 抽象了当前登录的用户信息，通过接口避免循环依赖。
	// 允许渲染层访问认证用户的通用属性，而无需了解具体实现。
	AuthUser RenderUser `json:"auth_user,omitempty"`

	// Lang 当前页面使用的语言代码，通过智能检测获得
	Lang string `json:"lang"`
}

// BreadcrumbItem 面包屑导航项 - 用于构建页面导航路径
// 业务逻辑: 帮助用户了解当前页面在网站层级中的位置，提升导航体验
// 自动生成规则:
// - 首页 > 用户页面 > 具体功能页面
// - 首页 > 日记本 > 具体日记
// - 支持多级嵌套，最后一项通常不包含URL(当前页面)
type BreadcrumbItem struct {
	// Name 导航项显示名称，支持多语言翻译
	Name string `json:"name"`

	// URL 导航项链接地址，可选字段
	// 如果为空，表示当前页面，不可点击
	URL string `json:"url,omitempty"`
}
