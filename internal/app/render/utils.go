package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// RenderUtils 渲染工具类
type RenderUtils struct {
	service *RenderService
}

// NewRenderUtils 创建渲染工具类
func NewRenderUtils(service *RenderService) *RenderUtils {
	return &RenderUtils{
		service: service,
	}
}

// SafeHTML 安全的HTML输出
func (ru *RenderUtils) SafeHTML(content string) template.HTML {
	return template.HTML(content)
}

// FormatFileSize 格式化文件大小
func (ru *RenderUtils) FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// TruncateString 截断字符串
func (ru *RenderUtils) TruncateString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength] + "..."
}

// FormatDate 格式化日期
func (ru *RenderUtils) FormatDate(timestamp int64, format string) string {
	t := time.Unix(timestamp, 0)
	return t.Format(format)
}

// FormatDateWithTimeZone 根据时区格式化日期
func (ru *RenderUtils) FormatDateWithTimeZone(timestamp int64, format string, timezone string) string {
	t := time.Unix(timestamp, 0)

	if timezone != "" {
		if loc, err := time.LoadLocation(timezone); err == nil {
			t = t.In(loc)
		}
	}

	return t.Format(format)
}

// URLJoin 连接URL路径
func (ru *RenderUtils) URLJoin(base string, parts ...string) string {
	result := strings.TrimRight(base, "/")
	for _, part := range parts {
		part = strings.Trim(part, "/")
		if part != "" {
			result += "/" + part
		}
	}
	return result
}

// JSONMarshal JSON序列化
func (ru *RenderUtils) JSONMarshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSONUnmarshal JSON反序列化
func (ru *RenderUtils) JSONUnmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// GetTemplatePath 获取模板路径
func (ru *RenderUtils) GetTemplatePath(templateName string) string {
	return filepath.Join("website", "views", templateName)
}

// RenderTemplate 渲染模板到字符串
func (ru *RenderUtils) RenderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(ru.GetTemplatePath(templateName))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetClientIP 获取客户端IP地址
func (ru *RenderUtils) GetClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For 头
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 检查 X-Real-IP 头
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// 使用 RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	return ip
}

// GetUserAgent 获取用户代理
func (ru *RenderUtils) GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

// IsAjaxRequest 检查是否为Ajax请求
func (ru *RenderUtils) IsAjaxRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

// GetReferer 获取引用页面
func (ru *RenderUtils) GetReferer(r *http.Request) string {
	return r.Header.Get("Referer")
}

// IsSecureRequest 检查是否为HTTPS请求
func (ru *RenderUtils) IsSecureRequest(r *http.Request) bool {
	return r.TLS != nil || strings.ToLower(r.Header.Get("X-Forwarded-Proto")) == "https"
}

// GetContentType 获取内容类型
func (ru *RenderUtils) GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// EscapeHTML HTML转义
func (ru *RenderUtils) EscapeHTML(str string) string {
	return template.HTMLEscapeString(str)
}

// UnescapeHTML HTML反转义
func (ru *RenderUtils) UnescapeHTML(str string) string {
	// Go的template包中没有HTMLUnescapeString，需要手动实现
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&quot;", "\"")
	str = strings.ReplaceAll(str, "&#39;", "'")
	return str
}

// GetCurrentTimestamp 获取当前时间戳
func (ru *RenderUtils) GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentTimestampMs 获取当前时间戳（毫秒）
func (ru *RenderUtils) GetCurrentTimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// ParseTimestamp 解析时间戳
func (ru *RenderUtils) ParseTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// ParseTimestampMs 解析时间戳（毫秒）
func (ru *RenderUtils) ParseTimestampMs(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000)
}

// GenerateRandomString 生成随机字符串
func (ru *RenderUtils) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// IsEmptyString 检查字符串是否为空
func (ru *RenderUtils) IsEmptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}

// ContainsString 检查字符串是否包含子字符串
func (ru *RenderUtils) ContainsString(str, substr string) bool {
	return strings.Contains(str, substr)
}

// ReplaceString 替换字符串
func (ru *RenderUtils) ReplaceString(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}
