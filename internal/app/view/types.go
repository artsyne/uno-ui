package view

// PageData 结构化页面数据
type PageData struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	User        interface{}            `json:"user,omitempty"`
	Data        map[string]interface{} `json:"data"`
	Navigation  interface{}            `json:"navigation,omitempty"`
	Lang        string                 `json:"lang"`
}

// NavigationData 导航数据
type NavigationData struct {
	Current string                 `json:"current"`
	Items   []NavigationItem       `json:"items"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// NavigationItem 导航项
type NavigationItem struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}
