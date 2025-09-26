package pkg

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownProcessor Markdown 处理器，提供 Markdown 渲染功能
type MarkdownProcessor struct {
	markdown goldmark.Markdown
}

// NewMarkdownProcessor 创建一个新的 Markdown 处理器实例
// 支持以下 Markdown 语法：
// - Italic: *italic* or _italic_
// - Bold: **bold** or __bold__
// - Combined: **bold and _italic_**
// - Horizontal rule: ---
// - Strikethrough: ~~strike this~~
// - Ordered list: 1. List item
// - Unordered list: - List item
// - Links: [link text](link URL)
// - Inline code: `code`
// - Multiline code: ```multiline code```
// - Blockquote: > I'm a quote
// - Headers: # h1 ## h2 ### h3 #### h4
//
// 注意：该处理器不支持 ==highlighted== 语法，因为这需要 goldmark 开启不安全的 HTML 模式，
// 与此函数的“只支持特定 Markdown 语法”的目标相冲突。
func NewMarkdownProcessor() *MarkdownProcessor {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown (支持 strikethrough, tables 等)
			extension.Linkify,    // 自动链接
			extension.Typographer, // 排版增强
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 自动生成标题 ID
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 硬换行
			html.WithXHTML(),     // XHTML 兼容
		),
	)

	return &MarkdownProcessor{
		markdown: md,
	}
}

// RenderToHTML 将 Markdown 文本渲染为安全的 HTML 字符串
// source: Markdown 源文本
// 返回: 渲染后的 HTML 字符串和可能的错误
func (mp *MarkdownProcessor) RenderToHTML(source string) (template.HTML, error) {
	var buf bytes.Buffer
	err := mp.markdown.Convert([]byte(source), &buf)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// RenderFromBytes 从字节数组渲染 Markdown 为 HTML
// source: Markdown 源字节数组
// 返回: 渲染后的 HTML 字符串和可能的错误
func (mp *MarkdownProcessor) RenderFromBytes(source []byte) (template.HTML, error) {
	var buf bytes.Buffer
	err := mp.markdown.Convert([]byte(source), &buf)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// defaultMarkdownProcessor 是一个线程安全的单例 Markdown 处理器
var (
	defaultMarkdownProcessor *MarkdownProcessor
	once                     sync.Once
)

// GetDefaultMarkdownProcessor 获取默认的 Markdown 处理器
// 该函数使用 sync.Once 确保处理器只被初始化一次，是线程安全的。
func GetDefaultMarkdownProcessor() *MarkdownProcessor {
	once.Do(func() {
		defaultMarkdownProcessor = NewMarkdownProcessor()
	})
	return defaultMarkdownProcessor
}

// RenderMarkdown 使用默认处理器将 Markdown 渲染为安全的 HTML
// 这是一个便捷函数，用于快速渲染 Markdown
func RenderMarkdown(source string) template.HTML {
	renderedHTML, err := GetDefaultMarkdownProcessor().RenderToHTML(source)
	if err != nil {
		// If rendering fails, return the source string escaped as HTML
		return template.HTML(template.HTMLEscapeString(source))
	}
	return renderedHTML
}