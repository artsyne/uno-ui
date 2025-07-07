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
}

// PageData 页面数据结构
type PageData struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Keywords	string                 `json:"keywords,omitempty"` 
	Data        map[string]interface{} `json:"data"`
	Breadcrumb  []BreadcrumbItem       `json:"breadcrumb,omitempty"`
	Context     *RenderContext         `json:"context,omitempty"`
}

// 便捷访问方法
func (p *PageData) IsAuthenticated() bool {
	return p.Context != nil && p.Context.IsLogin
}

func (p *PageData) GetAuthUser() RenderUser {
	if p.Context != nil {
		return p.Context.AuthUser
	}
	return nil
}

// BreadcrumbItem 面包屑导航项
type BreadcrumbItem struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// JSONResponse 标准JSON响应结构
type JSONResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
