package pkg

import (
	"fmt"
	"log"

	"github.com/biuren/hebe/pkg"
)

// ExampleMarkdownUsage 演示如何使用 Markdown 处理器
func ExampleMarkdownUsage() {
	// 示例 Markdown 文本
	markdownText := `# Markdown 功能演示

这是一个 **Markdown** 处理器的演示文档。

## 支持的语法

### 文本格式
- *斜体文本* 或 _斜体文本_
- **粗体文本** 或 __粗体文本__
- **粗体和 _斜体_ 组合**
- ~~删除线文本~~

### 列表

#### 有序列表
1. 第一项
2. 第二项
3. 第三项

#### 无序列表
- 项目 A
- 项目 B
- 项目 C

### 代码

内联代码：` + "`fmt.Println(\"Hello\")`" + `

代码块：
` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

### 引用

> 这是一个引用块
> 可以包含多行内容

### 链接

[访问 GitHub](https://github.com)

### 分隔线

---

## 结束

感谢使用我们的 Markdown 处理器！`

	// 方法 1: 使用默认处理器
	fmt.Println("=== 使用默认处理器 ===")
	html, err := pkg.RenderMarkdown(markdownText)
	if err != nil {
		log.Printf("渲染失败: %v", err)
		return
	}
	fmt.Println(html)

	// 方法 2: 创建自定义处理器
	fmt.Println("\n=== 使用自定义处理器 ===")
	processor := pkg.NewMarkdownProcessor()
	html2, err := processor.RenderToHTML(markdownText)
	if err != nil {
		log.Printf("渲染失败: %v", err)
		return
	}
	fmt.Println(html2)

	// 方法 3: 渲染为 HTML 并转换为字节数组
	fmt.Println("\n=== 渲染为字节数组 ===")
	htmlResult, err := processor.RenderToHTML(markdownText)
	if err != nil {
		log.Printf("渲染失败: %v", err)
		return
	}
	bytes := []byte(htmlResult)
	fmt.Printf("渲染结果长度: %d 字节\n", len(bytes))
}

// ExampleSimpleMarkdown 简单的 Markdown 渲染示例
func ExampleSimpleMarkdown() {
	// 简单的 Markdown 文本
	simpleText := "**Hello** _World_!"

	// 渲染为 HTML
	html, err := pkg.RenderMarkdown(simpleText)
	if err != nil {
		log.Printf("渲染失败: %v", err)
		return
	}

	fmt.Printf("输入: %s\n", simpleText)
	fmt.Printf("输出: %s\n", html)
}

// ExampleMarkdownFeatures 演示各种 Markdown 功能
func ExampleMarkdownFeatures() {
	features := map[string]string{
		"斜体":   "*italic* 和 _italic_",
		"粗体":   "**bold** 和 __bold__",
		"组合":   "**bold and _italic_**",
		"删除线":  "~~strikethrough~~",
		"代码":   "`inline code`",
		"链接":   "[GitHub](https://github.com)",
		"标题":   "# H1\n## H2\n### H3",
		"列表":   "- Item 1\n- Item 2",
		"引用":   "> This is a quote",
		"分隔线":  "---",
	}

	processor := pkg.NewMarkdownProcessor()

	for name, markdown := range features {
		fmt.Printf("\n=== %s ===\n", name)
		fmt.Printf("Markdown: %s\n", markdown)
		
		html, err := processor.RenderToHTML(markdown)
		if err != nil {
			log.Printf("渲染 %s 失败: %v", name, err)
			continue
		}
		
		fmt.Printf("HTML: %s\n", html)
	}
}