package pkg

import (
	"strings"
	"time"
)

// GetNowTime 获取当前的 UTC 时间，并返回一个 time.Time 对象。
// 这个函数不接受任何参数。
func GetNowTime() time.Time {
	return time.Now().UTC()
}

// GetNowTimeUnixMilli 获取当前的 UTC 时间，并将其转换为毫秒级的 Unix 时间戳。
// 这个函数不接受任何参数。
func GetNowTimeUnixMilli() int64 {
	return GetNowTime().UnixMilli()
}

// GetTimeFromUnixMilli 接受一个毫秒级的 Unix 时间戳作为参数，
// 并将其转换为对应的 UTC 时间，返回一个 time.Time 对象。
func GetTimeFromUnixMilli(unixMilli int64) time.Time {
	return time.Unix(0, unixMilli*int64(time.Millisecond)).UTC()
}

// GetLocalTimeFromUTC 接受一个毫秒级的 Unix 时间戳和一个时区字符串作为参数，
// 将 Unix 时间戳转换为对应的本地时间，并返回一个 time.Time 对象。
// 如果时区字符串无效，这个函数会返回一个错误。
func GetLocalTimeFromUTC(utcTime int64, timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, utcTime*int64(time.Millisecond)).In(loc), nil
}

// IsTimeZone 函数接受一个时区字符串作为参数，
// 并尝试加载对应的时区。如果加载成功，函数返回 true；
// 如果加载失败（例如，时区字符串无效），函数返回 false。
func IsTimeZone(timeZone string) bool {
	_, err := time.LoadLocation(timeZone)
	return err == nil
}

// FormatLocalTime 接受一个毫秒级的 Unix 时间戳和时区字符串，
// 根据时区返回格式化的时间字符串。
// 支持主流时区的本地化格式，出错时返回默认英语格式。
func FormatLocalTime(unixMilli int64, timeZone string) string {
	// 尝试加载时区，失败时使用UTC
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		loc = time.UTC
	}

	t := time.Unix(0, unixMilli*int64(time.Millisecond)).In(loc)

	// 根据时区选择格式
	switch {
	case strings.HasPrefix(timeZone, "Asia/") && (strings.Contains(timeZone, "Shanghai") || strings.Contains(timeZone, "Beijing") || strings.Contains(timeZone, "Hong_Kong") || strings.Contains(timeZone, "Taipei")):
		return t.Format("2006年1月2日 15:04")
	case strings.HasPrefix(timeZone, "Asia/") && (strings.Contains(timeZone, "Tokyo") || strings.Contains(timeZone, "Seoul")):
		return t.Format("2006年1月2日 15時04分")
	case strings.HasPrefix(timeZone, "America/") && (strings.Contains(timeZone, "Sao_Paulo") || strings.Contains(timeZone, "Rio")):
		return t.Format("02/01/2006 15:04")
	default:
		// 默认英语格式
		return t.Format("Jan 2, 2006 15:04")
	}
}
