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

import (
	"github.com/biuren/hebe/internal/app/render/localization"
	"github.com/biuren/hebe/internal/app/render/types"
)

// NewPageData 构造函数，统一初始化流程
func NewPageData(title, description, keywords string, data map[string]interface{}, breadcrumb []types.BreadcrumbItem, context *types.RenderContext, localizationManager localization.ManagerInterface) *PageData {
	return &PageData{
		Title:               title,
		Description:         description,
		Keywords:            keywords,
		Data:                data,
		Breadcrumb:          breadcrumb,
		Context:             context,
		localizationManager: localizationManager,
	}
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
	Breadcrumb []types.BreadcrumbItem `json:"breadcrumb,omitempty"`

	// Context 渲染上下文，包含页面渲染所需的全局状态信息
	// 包括: 用户登录状态、用户信息、语言设置等全局状态
	// 所有模板都可以访问这些通用状态
	Context *types.RenderContext `json:"context,omitempty"`

	// localizationManager 本地化管理器引用，不序列化到JSON响应中
	// 提供模板内的翻译功能: .T("key") 和 .TWithParams("key", params)
	// 自动根据用户语言偏好、Cookie或浏览器设置选择合适的语言
	localizationManager localization.ManagerInterface `json:"-"`
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
func (p *PageData) GetAuthUser() types.RenderUser {
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
func (p *PageData) GetRenderContext() *types.RenderContext {
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

// 注意：JSONResponse 和 ErrorResponse 已移动到 response 包中
// 这里只保留 PageData 相关定义
