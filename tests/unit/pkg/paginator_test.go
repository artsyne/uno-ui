package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
)

func TestGeneratePagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		options  pkg.ListOptions
		expected []pkg.PageLinks
	}{
		{
			name: "基本分页测试 - 第1页",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  10,
				Total:     100,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expected: []pkg.PageLinks{
				{URL: "", Text: "1", PageIndex: 1, IsCurrentPage: true},
				{URL: "/test?page=2", Text: "2", PageIndex: 2, IsCurrentPage: false},
				{URL: "/test?page=3", Text: "3", PageIndex: 3, IsCurrentPage: false},
				{URL: "/test?page=4", Text: "4", PageIndex: 4, IsCurrentPage: false},
				{URL: "/test?page=5", Text: "5", PageIndex: 5, IsCurrentPage: false},
				{URL: "/test?page=6", Text: "6", PageIndex: 6, IsCurrentPage: false},
				{URL: "/test?page=7", Text: "7", PageIndex: 7, IsCurrentPage: false},
				{URL: "", Text: "...", PageIndex: 0, IsCurrentPage: false},
				{URL: "/test?page=9", Text: "9", PageIndex: 9, IsCurrentPage: false},
				{URL: "/test?page=10", Text: "10", PageIndex: 10, IsCurrentPage: false},
				{URL: "/test?page=2", Text: "下一页", PageIndex: 2, IsCurrentPage: false},
			},
		},
		{
			name: "中间页测试 - 第5页",
			options: pkg.ListOptions{
				PageIndex: 5,
				PageSize:  10,
				Total:     100,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
		},
		{
			name: "最后一页测试",
			options: pkg.ListOptions{
				PageIndex: 10,
				PageSize:  10,
				Total:     100,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
		},
		{
			name: "单页测试",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  10,
				Total:     5,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expected: []pkg.PageLinks{
				{URL: "", Text: "1", PageIndex: 1, IsCurrentPage: true},
			},
		},
		{
			name: "页码超出范围测试 - 页码过大",
			options: pkg.ListOptions{
				PageIndex: 20,
				PageSize:  10,
				Total:     100,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
		},
		{
			name: "页码超出范围测试 - 页码过小",
			options: pkg.ListOptions{
				PageIndex: -1,
				PageSize:  10,
				Total:     100,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
		},
		{
			name: "空数据测试",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  10,
				Total:     0,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expected: []pkg.PageLinks{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			result := pkg.GeneratePagination(tt.options)
			
			// 验证结果不为nil
			assert.NotNil(t, result)
			
			// 对于有明确期望的测试用例，验证具体内容
			if tt.expected != nil {
				assert.Equal(t, len(tt.expected), len(result), "分页链接数量不匹配")
				for i, expected := range tt.expected {
					if i < len(result) {
						assert.Equal(t, expected.Text, result[i].Text, "链接文字不匹配")
						assert.Equal(t, expected.PageIndex, result[i].PageIndex, "页码不匹配")
						assert.Equal(t, expected.IsCurrentPage, result[i].IsCurrentPage, "当前页标记不匹配")
					}
				}
			}
			
			// 验证当前页标记
			currentPageFound := false
			for _, link := range result {
				if link.IsCurrentPage {
					currentPageFound = true
					break
				}
			}
			
			// 只有在有分页链接时才应该找到当前页
			if len(result) > 0 && tt.options.Total > 0 {
				assert.True(t, currentPageFound, "应该有一个当前页标记")
			}
		})
	}
}

func TestGetInterval(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	// 由于getInterval是私有函数，我们通过GeneratePagination间接测试
	tests := []struct {
		name       string
		pageIndex  int
		totalPages int
		pageSize   int
		total      int
	}{
		{
			name:       "开始区间测试",
			pageIndex:  1,
			totalPages: 10,
			pageSize:   10,
			total:      100,
		},
		{
			name:       "中间区间测试",
			pageIndex:  5,
			totalPages: 10,
			pageSize:   10,
			total:      100,
		},
		{
			name:       "结束区间测试",
			pageIndex:  10,
			totalPages: 10,
			pageSize:   10,
			total:      100,
		},
		{
			name:       "小范围测试",
			pageIndex:  3,
			totalPages: 5,
			pageSize:   10,
			total:      50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			options := pkg.ListOptions{
				PageIndex: tt.pageIndex,
				PageSize:  tt.pageSize,
				Total:     tt.total,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			}
			
			result := pkg.GeneratePagination(options)
			
			// 验证结果合理性
			assert.NotNil(t, result)
			
			// 验证页码范围合理
			for _, link := range result {
				if link.PageIndex > 0 { // 排除省略号等特殊项
					assert.True(t, link.PageIndex >= 1, "页码应该大于等于1")
					assert.True(t, link.PageIndex <= tt.totalPages, "页码应该小于等于总页数")
				}
			}
		})
	}
}

func TestAppendItem(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	// 由于appendItem是私有函数，我们通过GeneratePagination间接测试
	tests := []struct {
		name        string
		pageIndex   int
		currentPage int
		text        string
		url         string
	}{
		{
			name:        "普通页码项",
			pageIndex:   2,
			currentPage: 1,
			text:        "",
			url:         "/test?page=%d",
		},
		{
			name:        "当前页项",
			pageIndex:   1,
			currentPage: 1,
			text:        "",
			url:         "/test?page=%d",
		},
		{
			name:        "下一页按钮项",
			pageIndex:   0, // 不检查特定页码，而是检查文字
			currentPage: 1,
			text:        "下一页",
			url:         "/test?page=%d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			options := pkg.ListOptions{
				PageIndex: tt.currentPage,
				PageSize:  10,
				Total:     50,
				URL:       tt.url,
				PagePrev:  "上一页",
				PageNext:  "下一页",
			}
			
			result := pkg.GeneratePagination(options)
			
			// 验证结果包含预期的页码或文字
			found := false
			for _, link := range result {
				// 如果指定了特定页码，则按页码查找
				if tt.pageIndex > 0 && link.PageIndex == tt.pageIndex {
					found = true
					
					// 验证当前页标记
					if tt.pageIndex == tt.currentPage {
						assert.True(t, link.IsCurrentPage, "应该标记为当前页")
						assert.Empty(t, link.URL, "当前页不应该有URL")
					} else {
						assert.False(t, link.IsCurrentPage, "不应该标记为当前页")
						assert.NotEmpty(t, link.URL, "非当前页应该有URL")
					}
					
					// 验证文字内容
					if tt.text != "" {
						assert.Equal(t, tt.text, link.Text, "文字内容不匹配")
					}
					
					break
				}
				// 如果没有指定页码，则按文字查找
				if tt.pageIndex == 0 && tt.text != "" && link.Text == tt.text {
					found = true
					assert.False(t, link.IsCurrentPage, "导航按钮不应该标记为当前页")
					assert.NotEmpty(t, link.URL, "导航按钮应该有URL")
					break
				}
			}
			
			if tt.pageIndex > 0 && tt.pageIndex <= 5 { // 在合理范围内的页码应该能找到
				assert.True(t, found, "应该找到指定页码")
			} else if tt.pageIndex == 0 && tt.text != "" { // 按文字查找的情况
				assert.True(t, found, "应该找到指定文字的链接")
			}
		})
	}
}

func TestPaginationEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name    string
		options pkg.ListOptions
	}{
		{
			name: "页面大小为0",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  0,
				Total:     100,
				URL:       "/test?page=%d",
			},
		},
		{
			name: "总数为负数",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  10,
				Total:     -10,
				URL:       "/test?page=%d",
			},
		},
		{
			name: "页面大小为负数",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  -10,
				Total:     100,
				URL:       "/test?page=%d",
			},
		},
		{
			name: "大数据量测试",
			options: pkg.ListOptions{
				PageIndex: 500,
				PageSize:  20,
				Total:     10000,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			// 确保函数不会panic
			assert.NotPanics(t, func() {
				result := pkg.GeneratePagination(tt.options)
				assert.NotNil(t, result)
			})
		})
	}
}

func TestPaginationStructure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 测试分页结构的完整性
	options := pkg.ListOptions{
		PageIndex: 5,
		PageSize:  10,
		Total:     200,
		URL:       "/articles?page=%d",
		PagePrev:  "Previous",
		PageNext:  "Next",
	}
	
	result := pkg.GeneratePagination(options)
	
	// 验证基本结构
	assert.NotEmpty(t, result, "分页结果不应为空")
	
	// 验证包含上一页
	hasPrev := false
	for _, link := range result {
		if link.Text == "Previous" {
			hasPrev = true
			assert.Equal(t, 4, link.PageIndex, "上一页页码应该是4")
			break
		}
	}
	assert.True(t, hasPrev, "应该包含上一页链接")
	
	// 验证包含下一页
	hasNext := false
	for _, link := range result {
		if link.Text == "Next" {
			hasNext = true
			assert.Equal(t, 6, link.PageIndex, "下一页页码应该是6")
			break
		}
	}
	assert.True(t, hasNext, "应该包含下一页链接")
	
	// 验证当前页
	hasCurrent := false
	for _, link := range result {
		if link.IsCurrentPage {
			hasCurrent = true
			assert.Equal(t, 5, link.PageIndex, "当前页页码应该是5")
			assert.Equal(t, "5", link.Text, "当前页文字应该是5")
			assert.Empty(t, link.URL, "当前页不应该有URL")
			break
		}
	}
	assert.True(t, hasCurrent, "应该包含当前页标记")
}