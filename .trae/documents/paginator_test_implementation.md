# Paginator 模块测试实现方案

## 1. 模块概述

`Paginator.go` 是一个分页导航生成器，提供了完整的分页逻辑实现。该模块包含以下核心功能：

* 分页链接生成

* 分页区间计算

* 页码边界处理

* 省略号显示逻辑

## 2. 测试文件结构

### 2.1 文件位置

```
tests/unit/pkg/paginator_test.go
```

### 2.2 测试用例规划

```go
package pkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGeneratePagination_NormalCase 测试正常分页情况
func TestGeneratePagination_NormalCase(t *testing.T) {
	// 测试中间页码的分页生成
}

// TestGeneratePagination_EdgeCases 测试边界情况
func TestGeneratePagination_EdgeCases(t *testing.T) {
	// 测试页码 < 1 和页码 > 总页数的情况
}

// TestGeneratePagination_SinglePage 测试单页情况
func TestGeneratePagination_SinglePage(t *testing.T) {
	// 测试只有一页时的分页逻辑
}

// TestGeneratePagination_LargePageCount 测试大页数情况
func TestGeneratePagination_LargePageCount(t *testing.T) {
	// 测试页数很多时的省略号逻辑
}

// TestGetInterval 测试分页区间计算
func TestGetInterval(t *testing.T) {
	// 测试不同页码和总页数下的区间计算
}

// TestAppendItem 测试分页项生成
func TestAppendItem(t *testing.T) {
	// 测试分页项的URL和文本生成
}
```

## 3. 详细测试用例设计

### 3.1 TestGeneratePagination\_NormalCase

**测试目标**: 验证正常情况下的分页生成逻辑

**测试场景**:

```go
func TestGeneratePagination_NormalCase(t *testing.T) {
	tests := []struct {
		name     string
		options  ListOptions
		expected []PageLinks
	}{
		{
			name: "中间页码分页",
			options: ListOptions{
				PageIndex: 5,
				PageSize:  10,
				Total:     100,
				URL:       "/page/%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expected: []PageLinks{
				{URL: "/page/4", Text: "上一页", PageIndex: 4},
				{URL: "/page/1", Text: "1", PageIndex: 1},
				{URL: "/page/2", Text: "2", PageIndex: 2},
				{Text: "..."},
				{URL: "/page/2", Text: "2", PageIndex: 2},
				{URL: "/page/3", Text: "3", PageIndex: 3},
				{URL: "/page/4", Text: "4", PageIndex: 4},
				{URL: "", Text: "5", PageIndex: 5, IsCurrentPage: true},
				{URL: "/page/6", Text: "6", PageIndex: 6},
				{URL: "/page/7", Text: "7", PageIndex: 7},
				{URL: "/page/8", Text: "8", PageIndex: 8},
				{Text: "..."},
				{URL: "/page/9", Text: "9", PageIndex: 9},
				{URL: "/page/10", Text: "10", PageIndex: 10},
				{URL: "/page/6", Text: "下一页", PageIndex: 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeneratePagination(tt.options)
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

### 3.2 TestGeneratePagination\_EdgeCases

**测试目标**: 验证边界条件处理

**测试场景**:

```go
func TestGeneratePagination_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		options     ListOptions
		expectedLen int
		description string
	}{
		{
			name: "页码小于1",
			options: ListOptions{
				PageIndex: -1,
				PageSize:  10,
				Total:     100,
				URL:       "/page/%d",
			},
			description: "应该自动调整为第1页",
		},
		{
			name: "页码大于总页数",
			options: ListOptions{
				PageIndex: 20,
				PageSize:  10,
				Total:     100,
				URL:       "/page/%d",
			},
			description: "应该自动调整为最后一页",
		},
		{
			name: "总数为0",
			options: ListOptions{
				PageIndex: 1,
				PageSize:  10,
				Total:     0,
				URL:       "/page/%d",
			},
			description: "应该返回空分页",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeneratePagination(tt.options)
			// 验证结果的合理性
			for _, link := range result {
				if link.IsCurrentPage {
					assert.GreaterOrEqual(t, link.PageIndex, 1)
					assert.LessOrEqual(t, link.PageIndex, (tt.options.Total+tt.options.PageSize-1)/tt.options.PageSize)
				}
			}
		})
	}
}
```

### 3.3 TestGeneratePagination\_SinglePage

**测试目标**: 验证单页情况

```go
func TestGeneratePagination_SinglePage(t *testing.T) {
	options := ListOptions{
		PageIndex: 1,
		PageSize:  10,
		Total:     5,
		URL:       "/page/%d",
		PagePrev:  "上一页",
		PageNext:  "下一页",
	}

	result := GeneratePagination(options)

	// 单页情况下不应该有上一页和下一页链接
	for _, link := range result {
		assert.NotEqual(t, "上一页", link.Text)
		assert.NotEqual(t, "下一页", link.Text)
	}

	// 应该只有一个当前页
	currentPageCount := 0
	for _, link := range result {
		if link.IsCurrentPage {
			currentPageCount++
			assert.Equal(t, 1, link.PageIndex)
		}
	}
	assert.Equal(t, 1, currentPageCount)
}
```

### 3.4 TestGeneratePagination\_LargePageCount

**测试目标**: 验证大页数时的省略号逻辑

```go
func TestGeneratePagination_LargePageCount(t *testing.T) {
	options := ListOptions{
		PageIndex: 50,
		PageSize:  10,
		Total:     1000, // 100页
		URL:       "/page/%d",
		PagePrev:  "上一页",
		PageNext:  "下一页",
	}

	result := GeneratePagination(options)

	// 应该包含省略号
	hasEllipsis := false
	for _, link := range result {
		if link.Text == "..." {
			hasEllipsis = true
			break
		}
	}
	assert.True(t, hasEllipsis, "大页数情况下应该包含省略号")

	// 应该包含第1页和最后一页
	hasFirstPage := false
	hasLastPage := false
	for _, link := range result {
		if link.PageIndex == 1 {
			hasFirstPage = true
		}
		if link.PageIndex == 100 {
			hasLastPage = true
		}
	}
	assert.True(t, hasFirstPage, "应该包含第1页")
	assert.True(t, hasLastPage, "应该包含最后一页")
}
```

### 3.5 TestGetInterval

**测试目标**: 验证分页区间计算逻辑

```go
func TestGetInterval(t *testing.T) {
	tests := []struct {
		name       string
		pageIndex  int
		totalPages int
		expected   []int
	}{
		{
			name:       "开始页码",
			pageIndex:  1,
			totalPages: 20,
			expected:   []int{1, 7},
		},
		{
			name:       "中间页码",
			pageIndex:  10,
			totalPages: 20,
			expected:   []int{7, 13},
		},
		{
			name:       "结束页码",
			pageIndex:  20,
			totalPages: 20,
			expected:   []int{14, 20},
		},
		{
			name:       "总页数小于显示页数",
			pageIndex:  3,
			totalPages: 5,
			expected:   []int{1, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getInterval(tt.pageIndex, tt.totalPages)
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

### 3.6 TestAppendItem

**测试目标**: 验证分页项生成逻辑

```go
func TestAppendItem(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		pageIndex   int
		currentPage int
		text        string
		expected    PageLinks
	}{
		{
			name:        "普通页码",
			url:         "/page/%d",
			pageIndex:   5,
			currentPage: 3,
			text:        "",
			expected: PageLinks{
				URL:           "/page/5",
				Text:          "5",
				PageIndex:     5,
				IsCurrentPage: false,
			},
		},
		{
			name:        "当前页码",
			url:         "/page/%d",
			pageIndex:   3,
			currentPage: 3,
			text:        "",
			expected: PageLinks{
				URL:           "",
				Text:          "3",
				PageIndex:     3,
				IsCurrentPage: true,
			},
		},
		{
			name:        "自定义文本",
			url:         "/page/%d",
			pageIndex:   2,
			currentPage: 3,
			text:        "上一页",
			expected: PageLinks{
				URL:           "/page/2",
				Text:          "上一页",
				PageIndex:     2,
				IsCurrentPage: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendItem(tt.url, tt.pageIndex, tt.currentPage, tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

## 4. 性能测试

### 4.1 基准测试

```go
func BenchmarkGeneratePagination(b *testing.B) {
	options := ListOptions{
		PageIndex: 50,
		PageSize:  10,
		Total:     10000,
		URL:       "/page/%d",
		PagePrev:  "上一页",
		PageNext:  "下一页",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GeneratePagination(options)
	}
}

func BenchmarkGetInterval(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getInterval(50, 100)
	}
}
```

## 5. 测试数据和辅助函数

### 5.1 测试数据生成器

```go
// generateTestOptions 生成测试用的分页选项
func generateTestOptions(pageIndex, pageSize, total int) ListOptions {
	return ListOptions{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     total,
		URL:       "/page/%d",
		PagePrev:  "上一页",
		PageNext:  "下一页",
	}
}

// assertValidPagination 验证分页结果的有效性
func assertValidPagination(t *testing.T, result []PageLinks, options ListOptions) {
	totalPages := (options.Total + options.PageSize - 1) / options.PageSize
	
	// 验证当前页标记
	currentPageCount := 0
	for _, link := range result {
		if link.IsCurrentPage {
			currentPageCount++
			assert.GreaterOrEqual(t, link.PageIndex, 1)
			assert.LessOrEqual(t, link.PageIndex, totalPages)
		}
	}
	assert.LessOrEqual(t, currentPageCount, 1, "应该最多有一个当前页")
	
	// 验证页码的有效性
	for _, link := range result {
		if link.PageIndex > 0 {
			assert.GreaterOrEqual(t, link.PageIndex, 1)
			assert.LessOrEqual(t, link.PageIndex, totalPages)
		}
	}
}
```

## 6. 测试覆盖率目标

* **函数覆盖率**: 100%

* **行覆盖率**: ≥ 95%

* **分支覆盖率**: ≥ 90%

## 7. 实施步骤

1. **创建测试文件**: `tests/unit/pkg/paginator_test.go`
2. **实现基础测试用例**: 正常情况和边界情况
3. **添加性能测试**: 基准测试
4. **验证测试覆盖率**: 使用 `go test -cover`
5. **优化测试用例**: 根据覆盖率报告补充遗漏的测试
6. **文档更新**: 更新测试文档和使用说明

## 8. 注意事项

* 测试用例应该覆盖所有可能的输入组合

* 特别关注边界条件和异常情况

* 确保测试的独立性，避免测试间的相互依赖

* 使用有意义的测试名称和描述

* 保持测试代码的简洁和可读性

