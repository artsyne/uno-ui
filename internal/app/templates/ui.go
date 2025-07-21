package templates

import (
	"fmt"
	"html/template"
	"strings"
)

// BreadcrumbItem 面包屑导航项目
type BreadcrumbItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GetUIFuncs 获取UI相关的模板函数
func GetUIFuncs() template.FuncMap {
	return template.FuncMap{
		"renderBreadcrumb": renderBreadcrumb,
	}
}

// renderBreadcrumb 渲染面包屑导航
func renderBreadcrumb(breadcrumb []BreadcrumbItem) template.HTML {
	if len(breadcrumb) == 0 {
		return ""
	}

	var parts []string
	for i, item := range breadcrumb {
		if i == len(breadcrumb)-1 {
			// 最后一个项目不需要链接
			parts = append(parts, fmt.Sprintf(`<span>%s</span>`, item.Name))
		} else if item.URL != "" {
			parts = append(parts, fmt.Sprintf(`<a href="%s">%s</a>`, item.URL, item.Name))
		} else {
			parts = append(parts, fmt.Sprintf(`<span>%s</span>`, item.Name))
		}
	}

	html := fmt.Sprintf(`<nav class="breadcrumb">%s</nav>`, strings.Join(parts, ` <span class="separator">/</span> `))
	return template.HTML(html)
}
