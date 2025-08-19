package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
)

// TestGetNowTime 测试 GetNowTime 函数
func TestGetNowTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	now := pkg.GetNowTime()
	if time.Since(now) > time.Second {
		t.Errorf("GetNowTime is not returning current time")
	}
}

// TestGetNowTimeUnixMilli 测试 GetNowTimeUnixMilli 函数
func TestGetNowTimeUnixMilli(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	// 获取当前时间的毫秒时间戳
	timestamp := pkg.GetNowTimeUnixMilli()

	// 验证时间戳是否合理（应该是一个正数且接近当前时间）
	if timestamp <= 0 {
		t.Error("时间戳应该是正数")
	}

	// 验证时间戳是否接近当前时间
	now := time.Now().UnixMilli()
	diff := now - timestamp
	if diff < 0 {
		diff = -diff
	}

	// 允许1秒的误差
	if diff > 1000 {
		t.Errorf("时间戳差异过大: %d ms", diff)
	}
}

// TestGetTimeFromUnixMilli 测试 GetTimeFromUnixMilli 函数
func TestGetTimeFromUnixMilli(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	now := time.Now().UnixMilli()
	timeFromUnixMilli := pkg.GetTimeFromUnixMilli(now)
	if timeFromUnixMilli.UnixMilli() != now {
		t.Errorf("GetTimeFromUnixMilli is not returning correct time")
	}
}

// TestGetLocalTimeFromUTC 测试 GetLocalTimeFromUTC 函数
func TestGetLocalTimeFromUTC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}
	
	now := time.Now().UnixMilli()
	localTime, err := pkg.GetLocalTimeFromUTC(now, "America/New_York")
	if err != nil {
		t.Errorf("GetLocalTimeFromUTC returned an error: %v", err)
	}
	if localTime.UnixMilli() != now {
		t.Errorf("GetLocalTimeFromUTC is not returning correct local time")
	}
}

func TestIsTimeZone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	tests := []struct {
		name     string
		timezone string
		want     bool
	}{
		{
			name:     "有效时区 - 亚洲/上海",
			timezone: "Asia/Shanghai",
			want:     true,
		},
		{
			name:     "有效时区 - UTC",
			timezone: "UTC",
			want:     true,
		},
		{
			name:     "有效时区 - 美国/纽约",
			timezone: "America/New_York",
			want:     true,
		},
		{
			name:     "无效时区",
			timezone: "Invalid/TimeZone",
			want:     false,
		},
		{
			name:     "空字符串",
			timezone: "",
			want:     true, // IsTimeZone实现中空字符串被认为是有效的
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			got := pkg.IsTimeZone(tt.timezone)
			if got != tt.want {
				t.Errorf("IsTimeZone(%q) = %v, want %v", tt.timezone, got, tt.want)
			}
		})
	}
}
