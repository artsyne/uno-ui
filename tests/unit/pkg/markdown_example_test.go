package pkg_test

import (
	"fmt"
	"log"

	"github.com/biuren/hebe/pkg"
)

// ExampleRenderMarkdown_usage 演示如何使用 Markdown 处理器
func ExampleRenderMarkdown_usage() {
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
	html := pkg.RenderMarkdown(markdownText)
	if html == "" { // 检查是否为空，表示渲染可能失败或内容为空
		log.Printf("渲染失败或内容为空")
		return
	}
	fmt.Println(html)

	// 方法 2: 创建自定义处理器
	fmt.Println("=== 使用自定义处理器 ===")
	processor := pkg.NewMarkdownProcessor()
	html2, err := processor.RenderToHTML(markdownText)
	if err != nil {
		log.Printf("渲染失败: %v", err)
		return
	}
	fmt.Println(html2)

	// 方法 3: 渲染为 HTML 并转换为字节数组
	fmt.Println("=== 渲染为字节数组 ===")
	htmlResult, err := processor.RenderToHTML(markdownText)
	if err != nil {
		log.Printf("渲染失败: %v", err)
		return
	}
	bytes := []byte(htmlResult)
	fmt.Printf("渲染结果长度: %d 字节\n", len(bytes))

	// Output:
	// === 使用默认处理器 ===
	// <h1 id="markdown-">Markdown 功能演示</h1>
	// <p>这是一个 <strong>Markdown</strong> 处理器的演示文档。</p>
	// <h2 id="heading">支持的语法</h2>
	// <h3 id="heading-1">文本格式</h3>
	// <ul>
	// <li><em>斜体文本</em> 或 <em>斜体文本</em></li>
	// <li><strong>粗体文本</strong> 或 <strong>粗体文本</strong></li>
	// <li><strong>粗体和 <em>斜体</em> 组合</strong></li>
	// <li><del>删除线文本</del></li>
	// </ul>
	// <h3 id="heading-2">列表</h3>
	// <h4 id="heading-3">有序列表</h4>
	// <ol>
	// <li>第一项</li>
	// <li>第二项</li>
	// <li>第三项</li>
	// </ol>
	// <h4 id="heading-4">无序列表</h4>
	// <ul>
	// <li>项目 A</li>
	// <li>项目 B</li>
	// <li>项目 C</li>
	// </ul>
	// <h3 id="heading-5">代码</h3>
	// <p>内联代码：<code>fmt.Println(&quot;Hello&quot;)</code></p>
	// <p>代码块：</p>
	// <pre><code class="language-go">func main() {
	//     fmt.Println(&quot;Hello, World!&quot;)
	// }
	// </code></pre>
	// <h3 id="heading-6">引用</h3>
	// <blockquote>
	// <p>这是一个引用块<br />
	// 可以包含多行内容</p>
	// </blockquote>
	// <h3 id="heading-7">链接</h3>
	// <p><a href="https://github.com">访问 GitHub</a></p>
	// <h3 id="heading-8">分隔线</h3>
	// <hr />
	// <h2 id="heading-9">结束</h2>
	// <p>感谢使用我们的 Markdown 处理器！</p>
	//
	// === 使用自定义处理器 ===
	// <h1 id="markdown-">Markdown 功能演示</h1>
	// <p>这是一个 <strong>Markdown</strong> 处理器的演示文档。</p>
	// <h2 id="heading">支持的语法</h2>
	// <h3 id="heading-1">文本格式</h3>
	// <ul>
	// <li><em>斜体文本</em> 或 <em>斜体文本</em></li>
	// <li><strong>粗体文本</strong> 或 <strong>粗体文本</strong></li>
	// <li><strong>粗体和 <em>斜体</em> 组合</strong></li>
	// <li><del>删除线文本</del></li>
	// </ul>
	// <h3 id="heading-2">列表</h3>
	// <h4 id="heading-3">有序列表</h4>
	// <ol>
	// <li>第一项</li>
	// <li>第二项</li>
	// <li>第三项</li>
	// </ol>
	// <h4 id="heading-4">无序列表</h4>
	// <ul>
	// <li>项目 A</li>
	// <li>项目 B</li>
	// <li>项目 C</li>
	// </ul>
	// <h3 id="heading-5">代码</h3>
	// <p>内联代码：<code>fmt.Println(&quot;Hello&quot;)</code></p>
	// <p>代码块：</p>
	// <pre><code class="language-go">func main() {
	//     fmt.Println(&quot;Hello, World!&quot;)
	// }
	// </code></pre>
	// <h3 id="heading-6">引用</h3>
	// <blockquote>
	// <p>这是一个引用块<br />
	// 可以包含多行内容</p>
	// </blockquote>
	// <h3 id="heading-7">链接</h3>
	// <p><a href="https://github.com">访问 GitHub</a></p>
	// <h3 id="heading-8">分隔线</h3>
	// <hr />
	// <h2 id="heading-9">结束</h2>
	// <p>感谢使用我们的 Markdown 处理器！</p>
	//
	// === 渲染为字节数组 ===
	// 渲染结果长度: 1216 字节
}

// ExampleRenderMarkdown_simple 简单的 Markdown 渲染示例
func ExampleRenderMarkdown_simple() {
	// 简单的 Markdown 文本
	simpleText := "**Hello** _World_!"

	// 渲染为 HTML
	html := pkg.RenderMarkdown(simpleText)
	if html == "" { // 检查是否为空，表示渲染可能失败或内容为空
		log.Printf("渲染失败或内容为空")
		return
	}

	fmt.Printf("输入: %s\n", simpleText)
	fmt.Printf("输出: %s\n", html)

	// Output:
	// 输入: **Hello** _World_!
	// 输出: <p><strong>Hello</strong> <em>World</em>!</p>
}

// ExampleRenderMarkdown_features 演示 Markdown 支持的特性
func ExampleRenderMarkdown_features() {
	features := []struct {
		name     string
		markdown string
	}{
		{"斜体", "*italic* 和 _italic_"},
		{"粗体", "**bold** 和 __bold__"},
		{"组合", "**bold and _italic_**"},
		{"删除线", "~~strikethrough~~"},
		{"代码", "`inline code`"},
		{"链接", "[GitHub](https://github.com)"},
		{"标题", "# H1\n## H2\n### H3"},
		{"列表", "- Item 1\n- Item 2"},
		{"引用", "> This is a quote"},
		{"分隔线", "---"},
	}

	processor := pkg.NewMarkdownProcessor()

	for _, feature := range features {
		fmt.Printf("=== %s ===\n", feature.name)
		fmt.Printf("Markdown: %s\n", feature.markdown)

		html, err := processor.RenderToHTML(feature.markdown)
		if err != nil {
			log.Printf("渲染 %s 失败: %v", feature.name, err)
			continue
		}

		fmt.Printf("HTML: %s\n", html)
	}

	// Output:
	// === 斜体 ===
	// Markdown: *italic* 和 _italic_
	// HTML: <p><em>italic</em> 和 <em>italic</em></p>
	//
	// === 粗体 ===
	// Markdown: **bold** 和 __bold__
	// HTML: <p><strong>bold</strong> 和 <strong>bold</strong></p>
	//
	// === 组合 ===
	// Markdown: **bold and _italic_**
	// HTML: <p><strong>bold and <em>italic</em></strong></p>
	//
	// === 删除线 ===
	// Markdown: ~~strikethrough~~
	// HTML: <p><del>strikethrough</del></p>
	//
	// === 代码 ===
	// Markdown: `inline code`
	// HTML: <p><code>inline code</code></p>
	//
	// === 链接 ===
	// Markdown: [GitHub](https://github.com)
	// HTML: <p><a href="https://github.com">GitHub</a></p>
	//
	// === 标题 ===
	// Markdown: # H1
	// ## H2
	// ### H3
	// HTML: <h1 id="h1">H1</h1>
	// <h2 id="h2">H2</h2>
	// <h3 id="h3">H3</h3>
	//
	// === 列表 ===
	// Markdown: - Item 1
	// - Item 2
	// HTML: <ul>
	// <li>Item 1</li>
	// <li>Item 2</li>
	// </ul>
	//
	// === 引用 ===
	// Markdown: > This is a quote
	// HTML: <blockquote>
	// <p>This is a quote</p>
	// </blockquote>
	//
	// === 分隔线 ===
	// Markdown: ---
	// HTML: <hr />
}
