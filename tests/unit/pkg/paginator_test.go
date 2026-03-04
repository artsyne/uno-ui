package pkg_test

import (
	"testing"

	"github.com/biuren/hebe/pkg"
)

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name        string
		pageIndex   int
		pageSize    int
		hasMore     bool
		urlTemplate string
		want        *pkg.Pagination
	}{
		{
			name:        "第一页有更多数据",
			pageIndex:   1,
			pageSize:    20,
			hasMore:     true,
			urlTemplate: "/list/p/%d",
			want: &pkg.Pagination{
				PageIndex:  1,
				HasMore:    true,
				HasPrev:    false,
				PrevURL:    "",
				NextURL:    "/list/p/2",
				NoPushURL:  false,
			},
		},
		{
			name:        "第一页没有更多数据",
			pageIndex:   1,
			pageSize:    20,
			hasMore:     false,
			urlTemplate: "/list/p/%d",
			want: &pkg.Pagination{
				PageIndex:  1,
				HasMore:    false,
				HasPrev:    false,
				PrevURL:    "",
				NextURL:    "",
				NoPushURL:  false,
			},
		},
		{
			name:        "中间页有更多数据",
			pageIndex:   3,
			pageSize:    20,
			hasMore:     true,
			urlTemplate: "/list/p/%d",
			want: &pkg.Pagination{
				PageIndex:  3,
				HasMore:    true,
				HasPrev:    true,
				PrevURL:    "/list/p/2",
				NextURL:    "/list/p/4",
				NoPushURL:  false,
			},
		},
		{
			name:        "最后一页没有更多数据",
			pageIndex:   5,
			pageSize:    20,
			hasMore:     false,
			urlTemplate: "/list/p/%d",
			want: &pkg.Pagination{
				PageIndex:  5,
				HasMore:    false,
				HasPrev:    true,
				PrevURL:    "/list/p/4",
				NextURL:    "",
				NoPushURL:  false,
			},
		},
		{
			name:        "页码小于1自动修正",
			pageIndex:   0,
			pageSize:    20,
			hasMore:     true,
			urlTemplate: "/list/p/%d",
			want: &pkg.Pagination{
				PageIndex:  1,
				HasMore:    true,
				HasPrev:    false,
				PrevURL:    "",
				NextURL:    "/list/p/2",
				NoPushURL:  false,
			},
		},
		{
			name:        "负数页码自动修正",
			pageIndex:   -5,
			pageSize:    20,
			hasMore:     true,
			urlTemplate: "/list/p/%d",
			want: &pkg.Pagination{
				PageIndex:  1,
				HasMore:    true,
				HasPrev:    false,
				PrevURL:    "",
				NextURL:    "/list/p/2",
				NoPushURL:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pkg.NewPagination(tt.pageIndex, tt.pageSize, tt.hasMore, tt.urlTemplate)
			if got.PageIndex != tt.want.PageIndex {
				t.Errorf("PageIndex = %v, want %v", got.PageIndex, tt.want.PageIndex)
			}
			if got.HasMore != tt.want.HasMore {
				t.Errorf("HasMore = %v, want %v", got.HasMore, tt.want.HasMore)
			}
			if got.HasPrev != tt.want.HasPrev {
				t.Errorf("HasPrev = %v, want %v", got.HasPrev, tt.want.HasPrev)
			}
			if got.PrevURL != tt.want.PrevURL {
				t.Errorf("PrevURL = %v, want %v", got.PrevURL, tt.want.PrevURL)
			}
			if got.NextURL != tt.want.NextURL {
				t.Errorf("NextURL = %v, want %v", got.NextURL, tt.want.NextURL)
			}
			if got.NoPushURL != tt.want.NoPushURL {
				t.Errorf("NoPushURL = %v, want %v", got.NoPushURL, tt.want.NoPushURL)
			}
		})
	}
}

func TestPaginationEdgeCases(t *testing.T) {
	t.Run("空URL模板", func(t *testing.T) {
		p := pkg.NewPagination(2, 20, true, "")
		if p.PrevURL != "%!d()" && p.PrevURL != "" {
			t.Logf("PrevURL with empty template: %s", p.PrevURL)
		}
	})

	t.Run("复杂URL模板", func(t *testing.T) {
		p := pkg.NewPagination(3, 20, true, "/api/entries?status=1&page=%d&sort=desc")
		expectedPrev := "/api/entries?status=1&page=2&sort=desc"
		expectedNext := "/api/entries?status=1&page=4&sort=desc"
		if p.PrevURL != expectedPrev {
			t.Errorf("PrevURL = %v, want %v", p.PrevURL, expectedPrev)
		}
		if p.NextURL != expectedNext {
			t.Errorf("NextURL = %v, want %v", p.NextURL, expectedNext)
		}
	})
}

func TestSetNoPushURL(t *testing.T) {
	t.Run("默认值为false", func(t *testing.T) {
		p := pkg.NewPagination(1, 20, true, "/list/p/%d")
		if p.NoPushURL != false {
			t.Errorf("NoPushURL default = %v, want false", p.NoPushURL)
		}
	})

	t.Run("设置为true", func(t *testing.T) {
		p := pkg.NewPagination(1, 20, true, "/list/p/%d").SetNoPushURL(true)
		if p.NoPushURL != true {
			t.Errorf("NoPushURL = %v, want true", p.NoPushURL)
		}
	})

	t.Run("链式调用", func(t *testing.T) {
		p := pkg.NewPagination(1, 20, true, "/list/p/%d").SetNoPushURL(true)
		if p.NextURL != "/list/p/2" {
			t.Errorf("链式调用后 NextURL = %v, want /list/p/2", p.NextURL)
		}
	})
}