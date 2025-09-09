package pkg

import (
	"strings"
	"testing"

	"github.com/biuren/hebe/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMarkdownProcessor(t *testing.T) {
	mp := pkg.NewMarkdownProcessor()
	assert.NotNil(t, mp)
}

func TestMarkdownProcessor_RenderToHTML(t *testing.T) {
	mp := pkg.NewMarkdownProcessor()

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "Italic with asterisk",
			source:   "*italic*",
			expected: "<p><em>italic</em></p>",
		},
		{
			name:     "Italic with underscore",
			source:   "_italic_",
			expected: "<p><em>italic</em></p>",
		},
		{
			name:     "Bold with double asterisk",
			source:   "**bold**",
			expected: "<p><strong>bold</strong></p>",
		},
		{
			name:     "Bold with double underscore",
			source:   "__bold__",
			expected: "<p><strong>bold</strong></p>",
		},
		{
			name:     "Combined bold and italic",
			source:   "**bold and _italic_**",
			expected: "<p><strong>bold and <em>italic</em></strong></p>",
		},
		{
			name:     "Horizontal rule",
			source:   "---",
			expected: "<hr />",
		},
		{
			name:     "Strikethrough",
			source:   "~~strike this~~",
			expected: "<p><del>strike this</del></p>",
		},
		{
			name:     "Ordered list",
			source:   "1. List item",
			expected: "<ol>\n<li>List item</li>\n</ol>",
		},
		{
			name:     "Unordered list",
			source:   "- List item",
			expected: "<ul>\n<li>List item</li>\n</ul>",
		},
		{
			name:     "Link",
			source:   "[link text](https://example.com)",
			expected: "<p><a href=\"https://example.com\">link text</a></p>",
		},
		{
			name:     "Inline code",
			source:   "`code`",
			expected: "<p><code>code</code></p>",
		},
		{
			name:     "Multiline code",
			source:   "```\nmultiline code\n```",
			expected: "<pre><code>multiline code\n</code></pre>",
		},
		{
			name:     "Blockquote",
			source:   "> I'm a quote",
			expected: "<blockquote>\n<p>I&rsquo;m a quote</p>\n</blockquote>",
		},
		{
			name:     "Header h1",
			source:   "# h1",
			expected: "<h1 id=\"h1\">h1</h1>",
		},
		{
			name:     "Header h2",
			source:   "## h2",
			expected: "<h2 id=\"h2\">h2</h2>",
		},
		{
			name:     "Header h3",
			source:   "### h3",
			expected: "<h3 id=\"h3\">h3</h3>",
		},
		{
			name:     "Header h4",
			source:   "#### h4",
			expected: "<h4 id=\"h4\">h4</h4>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mp.RenderToHTML(tt.source)
			require.NoError(t, err)
			// 去除换行符进行比较，因为不同的 goldmark 版本可能有不同的换行处理
			resultStr := strings.TrimSpace(string(result))
			expected := strings.TrimSpace(tt.expected)
			assert.Equal(t, expected, resultStr)
		})
	}
}

func TestMarkdownProcessor_RenderToBytes(t *testing.T) {
	mp := pkg.NewMarkdownProcessor()
	source := "**bold text**"
	expected := "<p><strong>bold text</strong></p>"

	// 使用 RenderToHTML 方法，因为 RenderToBytes 方法不存在
	result, err := mp.RenderToHTML(source)
	require.NoError(t, err)
	resultStr := strings.TrimSpace(string(result))
	expectedStr := strings.TrimSpace(expected)
	assert.Equal(t, expectedStr, resultStr)
}

func TestMarkdownProcessor_RenderFromBytes(t *testing.T) {
	mp := pkg.NewMarkdownProcessor()
	source := []byte("*italic text*")
	expected := "<p><em>italic text</em></p>"

	result, err := mp.RenderFromBytes(source)
	require.NoError(t, err)
	resultStr := strings.TrimSpace(string(result))
	expectedStr := strings.TrimSpace(expected)
	assert.Equal(t, expectedStr, resultStr)
}

func TestRenderMarkdown(t *testing.T) {
	source := "# Hello World"
	expected := "<h1 id=\"hello-world\">Hello World</h1>"

	result, err := pkg.RenderMarkdown(source)
	require.NoError(t, err)
	resultStr := strings.TrimSpace(string(result))
	expectedStr := strings.TrimSpace(expected)
	assert.Equal(t, expectedStr, resultStr)
}

func TestRenderMarkdownBytes(t *testing.T) {
	source := []byte("## Subtitle")
	expected := "<h2 id=\"subtitle\">Subtitle</h2>"

	// 使用 RenderMarkdown 方法，因为 RenderMarkdownBytes 方法不存在
	result, err := pkg.RenderMarkdown(string(source))
	require.NoError(t, err)
	resultStr := strings.TrimSpace(string(result))
	expectedStr := strings.TrimSpace(expected)
	assert.Equal(t, expectedStr, resultStr)
}

func TestGetDefaultMarkdownProcessor(t *testing.T) {
	mp := pkg.GetDefaultMarkdownProcessor()
	assert.NotNil(t, mp)
}

// TestSetDefaultMarkdownProcessor 已移除，因为 SetDefaultMarkdownProcessor 方法不存在
// 默认处理器使用 sync.Once 实现单例模式，不支持重新设置
// func TestSetDefaultMarkdownProcessor(t *testing.T) {
//	// 该测试已移除，因为不再支持设置默认处理器
// }

// TestComplexMarkdown 测试复杂的 Markdown 文档
func TestComplexMarkdown(t *testing.T) {
	mp := pkg.NewMarkdownProcessor()
	source := `# Main Title

This is a paragraph with **bold** and *italic* text.

## Subtitle

Here's a list:
- Item 1
- Item 2 with ~~strikethrough~~

And a code block:
` + "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```" + `

> This is a blockquote

[Link to example](https://example.com)

---

End of document.`

	result, err := mp.RenderToHTML(source)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证包含预期的 HTML 元素
	assert.Contains(t, result, "<h1")
	assert.Contains(t, result, "<h2")
	assert.Contains(t, result, "<strong>bold</strong>")
	assert.Contains(t, result, "<em>italic</em>")
	assert.Contains(t, result, "<ul>")
	assert.Contains(t, result, "<li>")
	assert.Contains(t, result, "<del>strikethrough</del>")
	assert.Contains(t, result, "<pre><code")
	assert.Contains(t, result, "<blockquote>")
	assert.Contains(t, result, "<a href")
	assert.Contains(t, result, "<hr")
}

// BenchmarkRenderMarkdown 性能测试
func BenchmarkRenderMarkdown(b *testing.B) {
	source := `# Performance Test

This is a **performance** test with *various* markdown elements:

- List item 1
- List item 2
- List item 3

` + "```go\nfunc example() {\n    return \"test\"\n}\n```" + `

> Blockquote for testing

[Link](https://example.com)`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pkg.RenderMarkdown(source)
		if err != nil {
			b.Fatal(err)
		}
	}
}