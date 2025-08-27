package templates

import (
	"html/template"
	"time"

	"github.com/biuren/hebe/pkg"
)

// GetTimeFuncs 获取时间相关的模板函数
// 返回包含时间格式化相关函数的模板函数映射
//
// 返回:
//   - template.FuncMap: 包含以下时间相关函数的映射:
//     - formatTime: 格式化时间戳为指定格式和时区的时间字符串
//     - formatLocalTime: 格式化时间戳为本地时间字符串
func GetTimeFuncs() template.FuncMap {
	return template.FuncMap{
		"formatTime": formatTime,
		"formatLocalTime": formatLocalTime,
	}
}

// formatTime 格式化时间
// 将毫秒时间戳格式化为指定格式和时区的时间字符串
//
// 参数:
//   - timestamp: 毫秒级时间戳
//   - format: 时间格式字符串（Go时间格式）
//   - timezone: 可选的时区参数，如"Asia/Shanghai"、"UTC"等
//
// 返回:
//   - string: 格式化后的时间字符串
func formatTime(timestamp int64, format string, timezone ...string) string {
	t := time.Unix(timestamp/1000, (timestamp%1000)*1000000)

	// 如果指定了时区，进行转换
	if len(timezone) > 0 && timezone[0] != "" {
		if loc, err := time.LoadLocation(timezone[0]); err == nil {
			t = t.In(loc)
		}
	}

	return t.Format(format)
}

// formatLocalTime 格式化本地时间
// 使用pkg包的FormatLocalTime函数格式化时间戳为本地时间字符串
//
// 参数:
//   - timestamp: 毫秒级时间戳
//   - timezone: 时区字符串，如"Asia/Shanghai"、"UTC"等
//
// 返回:
//   - string: 格式化后的本地时间字符串
func formatLocalTime(timestamp int64, timezone string) string {
	return pkg.FormatLocalTime(timestamp, timezone)
}
