package tests

import (
	"testing"

	"github.com/biuren/hebe/pkg"
)

func TestNewTextProcessor(t *testing.T) {
	tp := pkg.NewTextProcessor()
	if tp == nil {
		t.Fatal("NewTextProcessor() returned nil")
	}
	if tp.DefaultSuffix != "..." {
		t.Errorf("Expected DefaultSuffix to be '...', got '%s'", tp.DefaultSuffix)
	}
	if tp.MaxSearchBack != 100 {
		t.Errorf("Expected MaxSearchBack to be 100, got %d", tp.MaxSearchBack)
	}
}

func TestTextProcessor_TruncateString(t *testing.T) {
	tp := pkg.NewTextProcessor()

	tests := []struct {
		name     string
		input    string
		length   int
		suffix   string
		expected string
	}{
		{
			name:     "空字符串",
			input:    "",
			length:   10,
			suffix:   "...",
			expected: "",
		},
		{
			name:     "长度为0",
			input:    "hello world",
			length:   0,
			suffix:   "...",
			expected: "...",
		},
		{
			name:     "不需要截断",
			input:    "hello",
			length:   10,
			suffix:   "...",
			expected: "hello",
		},
		{
			name:     "英文简单截断",
			input:    "hello world test",
			length:   10,
			suffix:   "...",
			expected: "hello...",
		},
		{
			name:     "中文截断",
			input:    "这是一个很长的中文测试字符串需要被截断",
			length:   10,
			suffix:   "...",
			expected: "这是一个很长的...",
		},
		{
			name:     "混合中英文截断",
			input:    "Hello 这是一个 mixed text 需要截断",
			length:   15,
			suffix:   "...",
			expected: "Hello 这是一个...",
		},
		{
			name:     "包含HTML标签",
			input:    "<p>Hello <strong>world</strong></p>",
			length:   10,
			suffix:   "...",
			expected: "Hello...",
		},
		{
			name:     "包含换行符",
			input:    "Hello\nworld\r\ntest",
			length:   10,
			suffix:   "...",
			expected: "Hello...",
		},
		{
			name:     "包含多个空格",
			input:    "Hello    world    test",
			length:   10,
			suffix:   "...",
			expected: "Hello...",
		},
		{
			name:     "句号截断",
			input:    "第一句话。第二句话。第三句话。",
			length:   8,
			suffix:   "...",
			expected: "第一句话。...",
		},
		{
			name:     "标点符号截断",
			input:    "Hello world! This is a test.",
			length:   15,
			suffix:   "...",
			expected: "Hello world!...",
		},
		{
			name:     "空格截断",
			input:    "Hello world this is a test",
			length:   15,
			suffix:   "...",
			expected: "Hello world...",
		},
		{
			name:     "日文截断",
			input:    "これは日本語のテストです。",
			length:   8,
			suffix:   "...",
			expected: "これは日本...",
		},
		{
			name:     "韩文截断",
			input:    "안녕하세요 이것은 한국어 테스트입니다",
			length:   10,
			suffix:   "...",
			expected: "안녕하세요...",
		},
		{
			name:     "特殊字符处理",
			input:    "Hello…world！测试？",
			length:   8,
			suffix:   "...",
			expected: "Hello...",
		},
		{
			name:     "后缀长度大于总长度",
			input:    "Hello",
			length:   3,
			suffix:   "...",
			expected: "...",
		},
		{
			name:     "复杂HTML内容",
			input:    `<div class="content"><p>这是<em>重要</em>内容</p><br/><span>更多内容</span></div>`,
			length:   10,
			suffix:   "...",
			expected: "这是重要内容更多内容",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tp.TruncateString(test.input, test.length, test.suffix)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestTextProcessor_TruncateWithDefault(t *testing.T) {
	tp := pkg.NewTextProcessor()
	result := tp.TruncateWithDefault("这是一个很长的测试字符串", 5)
	expected := "这是..."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestTextProcessor_CleanHTML(t *testing.T) {
	tp := pkg.NewTextProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "简单HTML标签",
			input:    "<p>Hello world</p>",
			expected: "Hello world",
		},
		{
			name:     "复杂HTML",
			input:    `<div class="test"><p>内容<span style="color:red">重要</span></p></div>`,
			expected: "内容重要",
		},
		{
			name:     "自闭合标签",
			input:    "Hello<br/>world<hr/>test",
			expected: "Helloworldtest",
		},
		{
			name:     "无HTML标签",
			input:    "纯文本内容",
			expected: "纯文本内容",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tp.CleanHTML(test.input)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestTextProcessor_NormalizeWhitespace(t *testing.T) {
	tp := pkg.NewTextProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "多个空格",
			input:    "Hello    world",
			expected: "Hello world",
		},
		{
			name:     "换行符",
			input:    "Hello\nworld\r\ntest",
			expected: "Hello world test",
		},
		{
			name:     "混合空白字符",
			input:    "  Hello\t\n  world  \r\n  ",
			expected: "Hello world",
		},
		{
			name:     "中文空白处理",
			input:    "你好\n\n世界\t\t测试",
			expected: "你好 世界 测试",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tp.NormalizeWhitespace(test.input)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestTextProcessorConfiguration(t *testing.T) {
	// 测试自定义配置
	tp := &pkg.TextProcessor{
		DefaultSuffix: ">>>",
		MaxSearchBack: 50,
	}

	result := tp.TruncateWithDefault("这是一个测试字符串，需要被截断", 5)
	expected := "这是>>>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestDefaultTextProcessor(t *testing.T) {
	// 测试默认全局实例
	defaultTP := pkg.GetDefaultTextProcessor()
	if defaultTP == nil {
		t.Fatal("GetDefaultTextProcessor() returned nil")
	}
	if defaultTP.DefaultSuffix != "..." {
		t.Errorf("Expected DefaultSuffix to be '...', got '%s'", defaultTP.DefaultSuffix)
	}

	// 测试向后兼容函数
	result := pkg.TruncateString("测试字符串", 5, "...")
	expected := "测试字符串"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	result = pkg.TruncateString("这是一个很长的测试字符串", 5, "...")
	expected = "这是..."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestSetDefaultTextProcessor(t *testing.T) {
	// 保存原始默认处理器
	originalTP := pkg.GetDefaultTextProcessor()

	// 设置新的默认处理器
	customTP := &pkg.TextProcessor{
		DefaultSuffix: "***",
		MaxSearchBack: 20,
	}
	pkg.SetDefaultTextProcessor(customTP)

	// 验证设置生效
	currentTP := pkg.GetDefaultTextProcessor()
	if currentTP.DefaultSuffix != "***" {
		t.Errorf("Expected DefaultSuffix to be '***', got '%s'", currentTP.DefaultSuffix)
	}
	if currentTP.MaxSearchBack != 20 {
		t.Errorf("Expected MaxSearchBack to be 20, got %d", currentTP.MaxSearchBack)
	}

	// 恢复原始默认处理器
	pkg.SetDefaultTextProcessor(originalTP)
}

func TestEdgeCases(t *testing.T) {
	tp := pkg.NewTextProcessor()

	tests := []struct {
		name     string
		input    string
		length   int
		suffix   string
		expected string
	}{
		{
			name:     "只有HTML标签",
			input:    "<p></p><div></div>",
			length:   10,
			suffix:   "...",
			expected: "...",
		},
		{
			name:     "只有空白字符",
			input:    "   \n\t  \r\n  ",
			length:   10,
			suffix:   "...",
			expected: "...",
		},
		{
			name:     "极长的后缀",
			input:    "Hello",
			length:   5,
			suffix:   "这是一个很长的后缀",
			expected: "Hello",
		},
		{
			name:     "单个字符",
			input:    "A",
			length:   1,
			suffix:   "",
			expected: "A",
		},
		{
			name:     "全角标点",
			input:    "你好！世界？这是测试。",
			length:   6,
			suffix:   "...",
			expected: "你好！...",
		},
		{
			name:     "长文本智能截断",
			input:    "第一句话。第二句话。第三句话。第四句话。第五句话。第六句话。第七句话。第八句话。第九句话。第十句话。",
			length:   50,
			suffix:   "...",
			expected: "第一句话。第二句话。第三句话。第四句话。第五句话。第六句话。第七句话。第八句话。第九句话。第十句话。",
		},
		{
			name:     "长文本强制截断",
			input:    "这是一个很长的没有标点符号的文本需要被强制截断处理",
			length:   20,
			suffix:   "...",
			expected: "这是一个很长的没有标点符号的文本需...",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := tp.TruncateString(test.input, test.length, test.suffix)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func BenchmarkTextProcessor_TruncateString(b *testing.B) {
	tp := pkg.NewTextProcessor()
	input := "这是一个用于性能测试的很长的字符串，包含中文和English mixed content，需要被智能截断处理。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tp.TruncateString(input, 20, "...")
	}
}

func BenchmarkTextProcessor_CleanHTML(b *testing.B) {
	tp := pkg.NewTextProcessor()
	input := `<div class="content"><p>这是<em>重要</em>内容</p><br/><span style="color:red">更多内容</span></div>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tp.CleanHTML(input)
	}
}
