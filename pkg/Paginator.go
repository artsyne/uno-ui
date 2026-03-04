package pkg

import "fmt"

// Pagination 分页信息（limit+1 模式）
// 简化设计：直接包含上一页/下一页链接，无需中间结构
type Pagination struct {
	PageIndex int    // 当前页码（从1开始）
	HasMore   bool   // 是否有更多数据
	HasPrev   bool   // 是否有上一页
	PrevURL   string // 上一页链接
	NextURL   string // 下一页链接
	NoPushURL bool   // 是否禁止更新浏览器URL（评论分页需要禁止）
}

// NewPagination 创建分页信息
// urlTemplate: 分页链接模板，如 "/u/123/p/%d"，%d 会被替换为页码
func NewPagination(pageIndex, pageSize int, hasMore bool, urlTemplate string) *Pagination {
	// 参数校验
	if pageIndex < 1 {
		pageIndex = 1
	}

	p := &Pagination{
		PageIndex: pageIndex,
		HasMore:   hasMore,
		HasPrev:   pageIndex > 1,
	}

	// 生成上一页链接
	if p.HasPrev {
		p.PrevURL = fmt.Sprintf(urlTemplate, pageIndex-1)
	}

	// 生成下一页链接
	if hasMore {
		p.NextURL = fmt.Sprintf(urlTemplate, pageIndex+1)
	}

	return p
}

// SetNoPushURL 设置禁止更新浏览器URL
func (p *Pagination) SetNoPushURL(noPush bool) *Pagination {
	p.NoPushURL = noPush
	return p
}
