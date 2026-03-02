package pkg_test

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
		name        string
		options     pkg.ListOptions
		expectedLen int // 期望的长度
	}{
		{
			name: "第一页有更多数据",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  10,
				HasMore:   true,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expectedLen: 1,
		},
		{
			name: "第一页没有更多数据",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  10,
				HasMore:   false,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expectedLen: 0,
		},
		{
			name: "中间页有更多数据",
			options: pkg.ListOptions{
				PageIndex: 5,
				PageSize:  10,
				HasMore:   true,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expectedLen: 2,
		},
		{
			name: "最后一页没有更多数据",
			options: pkg.ListOptions{
				PageIndex: 10,
				PageSize:  10,
				HasMore:   false,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			expectedLen: 1,
		},
		{
			name: "页面大小为0",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  0,
				HasMore:   true,
				URL:       "/test?page=%d",
			},
			expectedLen: 0,
		},
		{
			name: "页面大小为负数",
			options: pkg.ListOptions{
				PageIndex: 1,
				PageSize:  -10,
				HasMore:   true,
				URL:       "/test?page=%d",
			},
			expectedLen: 0,
		},
		{
			name: "页码过小",
			options: pkg.ListOptions{
				PageIndex: -1,
				PageSize:  10,
				HasMore:   true,
				URL:       "/test?page=%d",
				PagePrev:  "上一页",
				PageNext:  "下一页",
			},
			// pageIndex < 1 会被修正为 1
			expectedLen: 1,
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

			// 验证长度
			assert.Equal(t, tt.expectedLen, len(result), "分页链接数量不匹配")

			// 验证具体内容
			for _, link := range result {
				if link.PageIndex > 0 {
					// 验证 URL 格式正确
					assert.Contains(t, link.URL, "/test?page=", "URL格式不正确")
				}
			}
		})
	}
}

func TestGeneratePaginationEdgeCases(t *testing.T) {
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
				HasMore:   true,
				URL:       "/test?page=%d",
			},
		},
		{
			name: "大数据量测试",
			options: pkg.ListOptions{
				PageIndex: 500,
				PageSize:  20,
				HasMore:   true,
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