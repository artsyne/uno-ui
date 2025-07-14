// Package render 渲染层类型定义
//
// 这个包定义了应用渲染层的核心数据结构和接口，主要包括:
//
// 1. RenderUser: 渲染用户接口，抽象认证用户信息，避免循环依赖
// 2. RenderContext: 渲染上下文，包含页面通用状态(登录、语言等)
// 3. PageData: 页面数据结构，承载完整的页面渲染数据
// 4. 响应结构: JSON和错误响应的标准化格式
//
// 设计原则:
// - 接口隔离: 通过RenderUser接口避免渲染层对具体用户实现的依赖
// - 数据聚合: PageData集中管理页面所需的所有数据
// - 国际化支持: 内置多语言翻译能力
// - 用户体验: 自动处理SEO元数据、面包屑导航等
package render

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

// PageData 页面数据结构 - 用于渲染HTML页面的完整数据集合
//
// 业务逻辑说明:
// 1. 页面元数据管理: Title、Description、Keywords用于SEO优化和页面标识
// 2. 动态内容载体: Data字段承载页面特定的业务数据(如文章列表、用户信息等)
// 3. 导航体验: Breadcrumb提供用户路径导航，提升用户体验
// 4. 用户上下文: Context包含登录状态、用户信息、语言设置等全局状态
// 5. 国际化支持: localizationManager提供模板内的多语言翻译能力
//
// 使用场景:
// - 控制器构建页面数据传递给模板
// - 模板中使用.T()方法进行多语言翻译
// - 模板中访问用户状态(.IsAuthenticated(), .GetAuthUser())
// - SEO优化和社交媒体分享的元数据设置
type PageData struct {
	// Title 页面标题，用于HTML <title>标签和浏览器标签显示
	// 通过localization自动翻译，格式通常为: "页面名称 - 网站名称"
	Title string `json:"title"`

	// Description 页面描述，用于HTML <meta name="description">标签
	// 对SEO非常重要，搜索引擎会在搜索结果中显示此内容
	Description string `json:"description,omitempty"`

	// Keywords 页面关键词，用于HTML <meta name="keywords">标签
	// 虽然现代SEO中重要性降低，但仍有助于内容分类
	Keywords string `json:"keywords,omitempty"`

	// Data 页面特定的业务数据容器，灵活存储各种页面内容
	// 例如: 文章列表、用户资料、日记条目、分页信息等
	// 控制器根据业务需求向此字段填充相应数据
	Data map[string]interface{} `json:"data"`

	// Breadcrumb 面包屑导航，帮助用户了解当前页面在网站中的位置
	// 自动根据URL路径和用户状态生成，提升导航体验
	Breadcrumb []BreadcrumbItem `json:"breadcrumb,omitempty"`

	// Context 渲染上下文，包含页面渲染所需的全局状态信息
	// 包括: 用户登录状态、用户信息、当前语言等
	// 所有模板都可以访问这些通用状态
	Context *RenderContext `json:"context,omitempty"`

	// localizationManager 本地化管理器引用，不序列化到JSON响应中
	// 提供模板内的翻译功能: .T("key") 和 .TWithParams("key", params)
	// 自动根据用户语言偏好、Cookie或浏览器设置选择合适的语言
	localizationManager *LocalizationManager `json:"-"`
}

// ===================== PageData 便捷访问方法 =====================
// 以下方法提供模板中常用的状态检查和数据访问功能

// IsAuthenticated 检查当前用户是否已登录
// 业务用途: 模板中根据登录状态显示不同内容
// 使用示例: {{if .IsAuthenticated}} 显示用户菜单 {{else}} 显示登录按钮 {{end}}
func (p *PageData) IsAuthenticated() bool {
	return p.Context != nil && p.Context.IsLogin
}

// GetAuthUser 获取当前已认证的用户信息
// 业务用途: 在模板中访问用户头像、用户名、权限等信息
// 使用示例: {{.GetAuthUser.GetUserName}} 或 {{.GetAuthUser.GetFaceUrl}}
// 注意: 使用前应先检查 .IsAuthenticated 确保用户已登录
func (p *PageData) GetAuthUser() RenderUser {
	if p.Context != nil {
		return p.Context.AuthUser
	}
	return nil
}

// GetLang 获取当前页面使用的语言代码
// 业务用途: 模板中设置HTML lang属性、条件显示特定语言内容
// 使用示例: <html lang="{{.GetLang}}"> 或 {{if eq .GetLang "zh-cn"}} 中文特有内容 {{end}}
// 语言来源: 用户设置 > Cookie > 浏览器偏好 > 系统默认
func (p *PageData) GetLang() string {
	if p.Context != nil && p.Context.Lang != "" {
		return p.Context.Lang
	}
	return "en" // 默认语言
}

// GetRenderContext 获取完整的渲染上下文
// 业务用途: 访问完整的页面状态信息，通常在复杂模板逻辑中使用
// 包含: 登录状态、用户信息、语言设置等全局状态
func (p *PageData) GetRenderContext() *RenderContext {
	return p.Context
}

// ===================== 国际化翻译方法 =====================
// 这些方法支持模板中的多语言功能，自动根据用户语言偏好返回对应翻译

// T 基础翻译方法，根据键名获取对应语言的文本
// 业务用途: 模板中显示多语言文本，如导航菜单、按钮、标签等
// 使用示例: {{.T "nav.home"}} → "首页"(中文) 或 "Home"(英文)
// 回退机制: 如果翻译不存在，返回原始键名作为显示内容
func (p *PageData) T(key string) string {
	if p.localizationManager == nil || p.Context == nil || p.Context.Lang == "" {
		return key
	}
	return p.localizationManager.GetMessage(key, p.Context.Lang)
}

// TWithParams 带参数的翻译方法，支持在翻译文本中插入动态内容
// 业务用途: 显示包含用户名、数量等动态信息的多语言文本
// 使用示例:
//
//	{{$params := dict "username" .GetAuthUser.GetUserName}}
//	{{.TWithParams "greeting.welcome_user" $params}} → "欢迎回来，张三！"
//
// 参数格式: 翻译文本中使用 {{参数名}} 语法，如 "欢迎回来，{{username}}！"
func (p *PageData) TWithParams(key string, params map[string]string) string {
	if p.localizationManager == nil || p.Context == nil || p.Context.Lang == "" {
		return key
	}
	return p.localizationManager.GetMessageWithParams(key, p.Context.Lang, params)
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

// JSONResponse 标准JSON响应结构 - 用于API接口统一响应格式
// 业务逻辑: 保证前后端交互的数据格式一致性，便于前端统一处理
// 使用场景: AJAX请求、移动端API、第三方集成等
type JSONResponse struct {
	// Code 响应状态码，通常遵循HTTP状态码规范
	// 200: 成功, 400: 客户端错误, 500: 服务器错误等
	Code int `json:"code"`

	// Message 响应消息，给用户的提示信息
	// 支持多语言，根据用户语言偏好返回对应文本
	Message string `json:"message"`

	// Data 响应数据载体，根据不同接口返回相应的业务数据
	// 可以是对象、数组、字符串等任意类型
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应结构 - 专门用于错误情况的响应
// 业务逻辑: 统一错误信息格式，便于错误处理和用户提示
// 使用场景: 表单验证失败、权限不足、资源不存在等错误情况
type ErrorResponse struct {
	// Code 错误状态码，用于程序判断错误类型
	Code int `json:"code"`

	// Message 错误描述信息，给用户的错误提示
	// 应该是用户友好的语言，避免技术术语
	Message string `json:"message"`

	// Data 错误详细信息，可选字段
	// 例如: 表单验证错误的具体字段信息
	Data interface{} `json:"data,omitempty"`
}
