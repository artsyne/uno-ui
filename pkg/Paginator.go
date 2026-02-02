package pkg

import (
	"fmt"
	"math"
)

// ListOptions 分页选项
type ListOptions struct {
	PageIndex int    // 当前页码
	PageSize  int    // 每页显示的记录数
	Total     int    // 总记录数
	URL       string // 分页链接，其中应该包含 %d，用于替换为页码
	PagePrev  string // 上一页文字
	PageNext  string // 下一页文字
}

// PageLinks 分页链接结构
type PageLinks struct {
	URL           string // 链接地址
	Text          string // 链接文字
	PageIndex     int    // 页码
	IsCurrentPage bool   // 是否是当前页
}

// GeneratePagination 生成分页导航
func GeneratePagination(listOptions ListOptions) []PageLinks {
	// 如果总记录数小于等于0，直接返回空数组
	if listOptions.Total <= 0 {
		return []PageLinks{}
	}

	pageSize := listOptions.PageSize
	pageIndex := listOptions.PageIndex

	// 如果页面大小小于等于0，直接返回空数组，避免除零错误
	if pageSize <= 0 {
		return []PageLinks{}
	}

	// 根据总数和每页显示的记录数计算总页数
	totalPages := int(math.Ceil(float64(listOptions.Total) / float64(pageSize)))

	// 判断 pageIndex 参数是否合法，如果大于 totalPages 取 totalPages，如果小于 1 取 1
	if pageIndex < 1 {
		pageIndex = 1
	} else if pageIndex > totalPages {
		pageIndex = totalPages
	}

	var aryLinks []PageLinks
	interval := getInterval(pageIndex, totalPages)

	// 如果当前页码大于 1，则添加上一页链接
	if pageIndex > 1 {
		aryLinks = append(aryLinks, appendItem(listOptions.URL, pageIndex-1, pageIndex, listOptions.PagePrev))
	}

	// 添加起始页码链接
	if interval[0] > 1 {
		for i := 1; i <= 2 && i < interval[0]; i++ {
			aryLinks = append(aryLinks, appendItem(listOptions.URL, i, pageIndex, ""))
		}
		if interval[0] > 3 {
			aryLinks = append(aryLinks, PageLinks{Text: "..."})
		}
	}

	// 添加中间页码链接
	for i := interval[0]; i <= interval[1]; i++ {
		aryLinks = append(aryLinks, appendItem(listOptions.URL, i, pageIndex, ""))
	}

	// 添加结束页码链接
	if interval[1] < totalPages {
		if interval[1] < totalPages-2 {
			aryLinks = append(aryLinks, PageLinks{Text: "dots"})
		}
		for i := totalPages - 1; i <= totalPages; i++ {
			if i > interval[1] {
				aryLinks = append(aryLinks, appendItem(listOptions.URL, i, pageIndex, ""))
			}
		}
	}

	// 如果当前页码小于总页数，则显示下一页
	if pageIndex < totalPages {
		aryLinks = append(aryLinks, appendItem(listOptions.URL, pageIndex+1, pageIndex, listOptions.PageNext))
	}

	return aryLinks
}

// getInterval 计算分页区间
func getInterval(pageIndex, totalPages int) []int {
	const numPagesToShow = 7 // 当前页左右各显示 3 个页码，加上当前页共 7 个页码
	halfPagesToShow := numPagesToShow / 2

	start := pageIndex - halfPagesToShow
	if start < 1 {
		start = 1
	}

	end := pageIndex + halfPagesToShow
	if end > totalPages {
		end = totalPages
	}

	if end-start+1 < numPagesToShow {
		if start == 1 {
			end = int(math.Min(float64(start+numPagesToShow-1), float64(totalPages)))
		} else if end == totalPages {
			start = int(math.Max(float64(end-numPagesToShow+1), 1))
		}
	}

	return []int{start, end}
}

// appendItem 生成分页项
func appendItem(url string, pageIndex, currentPage int, text string) PageLinks {
	// 如果翻页链接文字为空，则显示页码
	if text == "" {
		text = fmt.Sprintf("%d", pageIndex)
	}

	// 如果当前页码等于当前页，则标记为当前页
	if pageIndex == currentPage {
		return PageLinks{
			URL:           "",
			Text:          text,
			PageIndex:     pageIndex,
			IsCurrentPage: true,
		}
	}
	return PageLinks{
		URL:       fmt.Sprintf(url, pageIndex),
		Text:      text,
		PageIndex: pageIndex,
	}
}
