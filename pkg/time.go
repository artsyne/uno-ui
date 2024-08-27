package pkg

import "time"

// GetNowTime 获取当前的 UTC 时间，并返回一个 time.Time 对象。
// 这个函数不接受任何参数。
func GetNowTime() time.Time {
    return time.Now().UTC()
}

// GetNowTimeUnixMilli 获取当前的 UTC 时间，并将其转换为毫秒级的 Unix 时间戳。
// 这个函数不接受任何参数。
func GetNowTimeUnixMilli() int64{
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