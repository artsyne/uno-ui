package pkg_test

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



// TestFormatLocalTime 测试 FormatLocalTime 函数
func TestFormatLocalTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	// 假设一个固定的时间戳，以便测试结果可预测
	// 2023年10月27日 07:30:00 UTC
	fixedUnixMilli := int64(1698391800000)

	tests := []struct {
		name     string
		unixMilli int64
		timeZone string
		expected string
	}{
		{
			name:     "上海时区 - 简体中文",
			unixMilli: fixedUnixMilli,
			timeZone: "Asia/Shanghai",
			expected: "2023年10月27日 15:30", // UTC 07:30 +8 = 15:30
		},
		{
			name:     "东京时区 - 日语",
			unixMilli: fixedUnixMilli,
			timeZone: "Asia/Tokyo",
			expected: "2023年10月27日 16時30分", // UTC 07:30 +9 = 16:30
		},
		{
			name:     "纽约时区 - 英语",
			unixMilli: fixedUnixMilli,
			timeZone: "America/New_York",
			expected: "Oct 27, 2023 03:30", // UTC 07:30 -4 = 03:30 (EDT)
		},
		{
			name:     "伦敦时区 - 英语",
			unixMilli: fixedUnixMilli,
			timeZone: "Europe/London",
			expected: "Oct 27, 2023 08:30", // UTC 07:30 +1 = 08:30 (BST)
		},
		{
			name:     "圣保罗时区 - 巴西葡萄牙语",
			unixMilli: fixedUnixMilli,
			timeZone: "America/Sao_Paulo",
			expected: "27/10/2023 04:30", // UTC 07:30 -3 = 04:30 (BRT)
		},
		{
			name:     "无效时区 - 返回默认格式",
			unixMilli: fixedUnixMilli,
			timeZone: "Invalid/TimeZone",
			expected: "Oct 27, 2023 07:30", // UTC时间，默认英语格式
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			got := pkg.FormatLocalTime(tt.unixMilli, tt.timeZone)
			if got != tt.expected {
				t.Errorf("FormatLocalTime(%d, %q) = %q, want %q", tt.unixMilli, tt.timeZone, got, tt.expected)
				t.Logf("Actual: %q, Expected: %q", got, tt.expected)
			}
		})
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
	
	tests := []struct {
		name     string
		timeZone string
		wantErr  bool
	}{
		{
			name:     "Valid timezone - New York",
			timeZone: "America/New_York",
			wantErr:  false,
		},
		{
			name:     "Valid timezone - Shanghai",
			timeZone: "Asia/Shanghai",
			wantErr:  false,
		},
		{
			name:     "Invalid timezone",
			timeZone: "Invalid/TimeZone",
			wantErr:  true,
		},
		{
			name:     "Empty timezone",
			timeZone: "",
			wantErr:  false, // 空字符串会使用UTC，不会报错
		},
	}
	
	now := time.Now().UnixMilli()
	
	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			localTime, err := pkg.GetLocalTimeFromUTC(now, tt.timeZone)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetLocalTimeFromUTC() expected error for timezone %s, got nil", tt.timeZone)
				}
			} else {
				if err != nil {
					t.Errorf("GetLocalTimeFromUTC() returned unexpected error: %v", err)
				}
				if localTime.UnixMilli() != now {
					t.Errorf("GetLocalTimeFromUTC() time mismatch, got %d, want %d", localTime.UnixMilli(), now)
				}
			}
		})
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
