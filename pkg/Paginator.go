package pkg

import (
	"fmt"
)

// ListOptions 分页选项（limit+1 模式）
type ListOptions struct {
	PageIndex int    // 当前页码
	PageSize  int    // 每页显示的记录数
	HasMore   bool   // 是否还有更多数据（limit+1 模式）
	URL       string // 分页链接，其中应该包含 %d，用于替换为页码
	PagePrev  string // 上一页文字
	PageNext  string // 下一页文字
}

// PageLinks 分页链接结构
type PageLinks struct {
	URL       string // 链接地址
	Text      string // 链接文字
	PageIndex int    // 页码
}

// GeneratePagination 生成分页导航（HasMore 模式）
// 只显示上一页和下一页，适用于"加载更多"场景
func GeneratePagination(listOptions ListOptions) []PageLinks {
	pageSize := listOptions.PageSize
	pageIndex := listOptions.PageIndex

	// 如果页面大小小于等于0，直接返回空数组，避免除零错误
	if pageSize <= 0 {
		return []PageLinks{}
	}

	// 如果页面索引小于1，设为1
	if pageIndex < 1 {
		pageIndex = 1
	}

	aryLinks := make([]PageLinks, 0)

	// 如果当前页码大于 1，则添加上一页链接
	if pageIndex > 1 {
		aryLinks = append(aryLinks, PageLinks{
			URL:       fmt.Sprintf(listOptions.URL, pageIndex-1),
			Text:      listOptions.PagePrev,
			PageIndex: pageIndex - 1,
		})
	}

	// 如果还有更多数据，则添加下一页链接
	if listOptions.HasMore {
		aryLinks = append(aryLinks, PageLinks{
			URL:       fmt.Sprintf(listOptions.URL, pageIndex+1),
			Text:      listOptions.PageNext,
			PageIndex: pageIndex + 1,
		})
	}

	return aryLinks
}