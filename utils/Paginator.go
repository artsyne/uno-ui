package utils

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
	pageSize := listOptions.PageSize
	pageIndex := listOptions.PageIndex

	// 根据总数和每页显示的记录数计算总页数
	totalPages := int(math.Floor(float64(listOptions.Total-1)/float64(pageSize)) + 1)

	// 判断 pageIndex 参数是否合法，如果大于 totalPages 取 totalPages，如果小于 0 取 1
	if pageIndex < 0 {
		pageIndex = 1
	} else if pageIndex >= totalPages {
		pageIndex = totalPages
	}

	var aryLinks []PageLinks
	interval := getInterval(pageIndex, totalPages)

	// 如果当前页码大于 1，则添加上一页链接
	if pageIndex > 1 {
		aryLinks = append(aryLinks, appendItem(listOptions.URL, pageIndex-1, pageIndex, totalPages, listOptions.PagePrev))
	}

	// 添加起始页码链接
	if interval[0] > 1 {
		end := int(math.Min(2, float64(interval[0])))
		for i := 0; i <= end; i++ {
			aryLinks = append(aryLinks, appendItem(listOptions.URL, i+1, pageIndex, totalPages, ""))
		}
		if interval[0] > 2 {
			aryLinks = append(aryLinks, PageLinks{Text: "..."})
		}
	}

	// 添加中间页码链接
	for i := interval[0]; i <= interval[1]; i++ {
		aryLinks = append(aryLinks, appendItem(listOptions.URL, i, pageIndex, totalPages, ""))
	}

	// 添加结束页码链接
	if interval[1] < totalPages {
		if totalPages-2 > interval[1] {
			aryLinks = append(aryLinks, PageLinks{Text: "..."})
		}
		begin := int(math.Max(float64(totalPages-2+1), float64(interval[1])))
		for i := begin; i < totalPages; i++ {
			aryLinks = append(aryLinks, appendItem(listOptions.URL, i+1, pageIndex, totalPages, ""))
		}
	}

	// 如果当前页码小于总页数-1，则显示下一页
	if pageIndex < totalPages-1 {
		aryLinks = append(aryLinks, appendItem(listOptions.URL, pageIndex+1, pageIndex, totalPages, listOptions.PageNext))
	}

	return aryLinks
}

// getInterval 计算分页区间
func getInterval(pageIndex, pages int) []int {
	neHalf := int(math.Ceil(8.0 / 2.0))
	upperLimit := pages - 8

	var start int
	if pageIndex > neHalf {
		start = int(math.Max(math.Min(float64(pageIndex-neHalf), float64(upperLimit)), 0))
	} else {
		start = 1
	}

	var end int
	if pageIndex > neHalf {
		end = int(math.Min(float64(pageIndex+neHalf), float64(pages)))
	} else {
		end = int(math.Min(float64(8), float64(pages)))
	}

	return []int{start, end}
}

// appendItem 生成分页项
func appendItem(url string, pageIndex, currentPage, totalPages int, text string) PageLinks {

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