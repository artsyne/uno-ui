package render

// PageData 页面数据结构
type PageData struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data"`
	Breadcrumb  []BreadcrumbItem       `json:"breadcrumb,omitempty"`
	IsLogin   bool                   `json:"is_login,omitempty"` // 是否登录
	IsOwner   bool                   `json:"is_owner,omitempty"` // 是否为当前用户
	TimeZone     string `json:"time_zone,omitempty"` //时区
	Lang	string `json:"lang,omitempty"` //语言
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
