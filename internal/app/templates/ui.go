package templates

import (
	"fmt"
	"html/template"
	"strings"
)

// BreadcrumbItem 面包屑导航项目
// 表示面包屑导航中的单个项目，包含显示名称和链接地址
type BreadcrumbItem struct {
	Name string `json:"name"` // 显示名称
	URL  string `json:"url"`  // 链接地址，可为空表示不可点击
}

// GetUIFuncs 获取UI相关的模板函数
// 返回包含用户界面相关函数的模板函数映射
//
// 返回:
//   - template.FuncMap: 包含以下UI相关函数的映射:
//     - renderBreadcrumb: 渲染面包屑导航HTML
func GetUIFuncs() template.FuncMap {
	return template.FuncMap{
		"renderBreadcrumb": renderBreadcrumb,
	}
}

// renderBreadcrumb 渲染面包屑导航
// 将面包屑导航项目列表渲染为HTML格式的导航栏
//
// 参数:
//   - breadcrumb: 面包屑导航项目列表
//
// 返回:
//   - template.HTML: 渲染后的HTML内容，包含完整的面包屑导航结构
//
// 渲染规则:
//   - 空列表返回空字符串
//   - 最后一个项目显示为普通文本（不可点击）
//   - 其他项目根据URL是否为空决定是否渲染为链接
//   - 项目间使用"/"分隔符连接
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
