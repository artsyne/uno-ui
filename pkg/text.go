package pkg

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TextProcessor 文本处理器，提供各种文本处理功能
type TextProcessor struct {
	// 可以添加配置字段，比如默认后缀、截断策略等
	DefaultSuffix string
	MaxSearchBack int // 向前搜索的最大字符数
}

// NewTextProcessor 创建一个新的文本处理器实例
func NewTextProcessor() *TextProcessor {
	return &TextProcessor{
		DefaultSuffix: "...",
		MaxSearchBack: 100, // 增加搜索范围，提供更好的截断点选择
	}
}

// TruncateString 智能截断字符串生成摘要
// 支持多语言，正确处理换行符、HTML标签、特殊字符等
// s: 要截断的字符串
// length: 截断的长度（以字符为单位）
// suffix: 截断后添加的字符串（如"..."）
func (tp *TextProcessor) TruncateString(s string, length int, suffix string) string {
	// 0. 处理空字符串
	if s == "" {
		return ""
	}

	if length <= 0 {
		return suffix
	}

	// 1. 清理和标准化文本
	cleaned := tp.cleanTextForSummary(s)

	// 2. 处理清理后为空的情况
	if cleaned == "" {
		return suffix
	}

	// 3. 检查是否需要截断
	if utf8.RuneCountInString(cleaned) <= length {
		return cleaned
	}

	// 4. 转换为rune数组进行Unicode安全处理
	runes := []rune(cleaned)

	// 5. 智能截断：寻找最佳截断点
	suffixLen := utf8.RuneCountInString(suffix)
	maxContentLength := length - suffixLen

	if maxContentLength <= 0 {
		return suffix
	}

	truncatePos := tp.findBestTruncatePosition(runes, maxContentLength)

	// 6. 生成最终结果
	result := string(runes[:truncatePos])

	// 7. 清理末尾空白并添加后缀
	result = strings.TrimRight(result, " \t\n\r")
	if result == "" {
		// 如果清理后为空，使用强制截断
		minLength := tp.min(maxContentLength, len(runes))
		if minLength > 0 {
			result = string(runes[:minLength])
		}
	}

	return result + suffix
}

// TruncateWithDefault 使用默认后缀截断字符串
func (tp *TextProcessor) TruncateWithDefault(s string, length int) string {
	return tp.TruncateString(s, length, tp.DefaultSuffix)
}

// CleanHTML 清理HTML标签
func (tp *TextProcessor) CleanHTML(s string) string {
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	return htmlRegex.ReplaceAllString(s, "")
}

// NormalizeWhitespace 规范化空白字符
func (tp *TextProcessor) NormalizeWhitespace(s string) string {
	// 规范化换行符
	s = regexp.MustCompile(`\r\n|\r|\n`).ReplaceAllString(s, " ")
	// 规范化空白字符
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	// 清理首尾空白
	return strings.TrimSpace(s)
}

// cleanTextForSummary 清理文本用于生成摘要
func (tp *TextProcessor) cleanTextForSummary(s string) string {
	// 1. 移除HTML标签
	s = tp.CleanHTML(s)
	// 2. 规范化空白字符
	return tp.NormalizeWhitespace(s)
}

// findBestTruncatePosition 寻找最佳截断位置
// 优先级：句号 > 其他标点 > 空格 > 强制截断
func (tp *TextProcessor) findBestTruncatePosition(runes []rune, maxLength int) int {
	if len(runes) <= maxLength {
		return len(runes)
	}

	// 在有效范围内寻找最佳截断点
	searchEnd := tp.min(maxLength, len(runes))

	// 1. 优先在句号处截断（中英文句号）
	for i := searchEnd - 1; i >= tp.max(0, searchEnd-tp.MaxSearchBack); i-- {
		if runes[i] == '.' || runes[i] == '。' {
			return i + 1
		}
	}

	// 2. 次选其他标点符号
	for i := searchEnd - 1; i >= tp.max(0, searchEnd-tp.MaxSearchBack); i-- {
		if tp.isPunctuation(runes[i]) {
			return i + 1
		}
	}

	// 3. 再选空格（主要用于英文等有空格分词的语言）
	for i := searchEnd - 1; i >= tp.max(0, searchEnd-tp.MaxSearchBack); i-- {
		if unicode.IsSpace(runes[i]) {
			return i
		}
	}

	// 4. 对于中文等无空格语言，在字符边界截断
	// 避免截断到符号或特殊字符中间
	for i := searchEnd - 1; i >= tp.max(0, searchEnd-10); i-- {
		if tp.isSafeToBreak(runes[i]) {
			return i + 1
		}
	}

	// 5. 强制截断
	return searchEnd
}

// isPunctuation 检查字符是否为标点符号
func (tp *TextProcessor) isPunctuation(r rune) bool {
	// 常见的结束标点符号
	endPuncts := []rune{'!', '?', ';', '！', '？', '；', '…'}
	for _, p := range endPuncts {
		if r == p {
			return true
		}
	}

	// 其他标点符号
	return unicode.IsPunct(r)
}

// isSafeToBreak 检查在该字符后是否可以安全截断
// 主要用于中文、日文等语言
func (tp *TextProcessor) isSafeToBreak(r rune) bool {
	// 中文、日文、韩文字符通常可以安全截断
	if unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r) {
		return true
	}

	// 英文字母数字后通常可以截断（如果后面跟着中文）
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// min 返回两个整数中的较小值
func (tp *TextProcessor) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func (tp *TextProcessor) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 全局默认实例，用于保持向后兼容性
var defaultTextProcessor = NewTextProcessor()

// 向后兼容的包级函数
// TruncateString 智能截断字符串生成摘要（向后兼容函数）
func TruncateString(s string, length int, suffix string) string {
	return defaultTextProcessor.TruncateString(s, length, suffix)
}

// GetDefaultTextProcessor 获取默认的文本处理器实例
func GetDefaultTextProcessor() *TextProcessor {
	return defaultTextProcessor
}

// SetDefaultTextProcessor 设置默认的文本处理器实例
func SetDefaultTextProcessor(tp *TextProcessor) {
	if tp != nil {
		defaultTextProcessor = tp
	}
}
