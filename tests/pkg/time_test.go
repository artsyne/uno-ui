package tests

import (
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
)

// TestGetNowTime 测试 GetNowTime 函数
func TestGetNowTime(t *testing.T) {
	now := pkg.GetNowTime()
	if time.Since(now) > time.Second {
		t.Errorf("GetNowTime is not returning current time")
	}
}

// TestGetNowTimeUnixMilli 测试 GetNowTimeUnixMilli 函数
func TestGetNowTimeUnixMilli(t *testing.T) {
	now := pkg.GetNowTimeUnixMilli()
	if time.Since(time.Unix(0, now*int64(time.Millisecond))) > time.Second {
		t.Errorf("GetNowTimeUnixMilli is not returning current time in milliseconds")
	}
}

// TestGetTimeFromUnixMilli 测试 GetTimeFromUnixMilli 函数
func TestGetTimeFromUnixMilli(t *testing.T) {
	now := time.Now().UnixMilli()
	timeFromUnixMilli := pkg.GetTimeFromUnixMilli(now)
	if timeFromUnixMilli.UnixMilli() != now {
		t.Errorf("GetTimeFromUnixMilli is not returning correct time")
	}
}

// TestGetLocalTimeFromUTC 测试 GetLocalTimeFromUTC 函数
func TestGetLocalTimeFromUTC(t *testing.T) {
	now := time.Now().UnixMilli()
	localTime, err := pkg.GetLocalTimeFromUTC(now, "America/New_York")
	if err != nil {
		t.Errorf("GetLocalTimeFromUTC returned an error: %v", err)
	}
	if localTime.UnixMilli() != now {
		t.Errorf("GetLocalTimeFromUTC is not returning correct local time")
	}
}