package templates

import (
	"html/template"
	"time"
)

// GetTimeFuncs 获取时间相关的模板函数
func GetTimeFuncs() template.FuncMap {
	return template.FuncMap{
		"formatTime": formatTime,
	}
}

// formatTime 格式化时间
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
