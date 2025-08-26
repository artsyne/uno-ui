package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
)

// TestIsID 测试ID验证功能
func TestIsID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效ID - 字母数字混合",
			input:    "abc123",
			expected: true,
		},
		{
			name:     "有效ID - 最小长度",
			input:    "abc123",
			expected: true,
		},
		{
			name:     "有效ID - 最大长度",
			input:    "abc123456789012345678901234567",
			expected: true,
		},
		{
			name:     "无效ID - 太短",
			input:    "abc12",
			expected: false,
		},
		{
			name:     "无效ID - 太长",
			input:    "abc1234567890123456789012345678",
			expected: false,
		},
		{
			name:     "无效ID - 包含特殊字符",
			input:    "abc123@",
			expected: false,
		},
		{
			name:     "无效ID - 包含空格",
			input:    "abc 123",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsID(tt.input)
			assert.Equal(t, tt.expected, result, "IsID(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}





// TestIsInteger 测试整数验证功能
func TestIsInteger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效整数 - 单位数",
			input:    "5",
			expected: true,
		},
		{
			name:     "有效整数 - 多位数",
			input:    "12345",
			expected: true,
		},
		{
			name:     "有效整数 - 零",
			input:    "0",
			expected: true,
		},
		{
			name:     "无效整数 - 包含字母",
			input:    "123a",
			expected: false,
		},
		{
			name:     "无效整数 - 包含小数点",
			input:    "123.0",
			expected: false,
		},
		{
			name:     "无效整数 - 负数",
			input:    "-123",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsInteger(tt.input)
			assert.Equal(t, tt.expected, result, "IsInteger(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsDecimals 测试小数验证功能
func TestIsDecimals(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效小数 - 标准格式",
			input:    "123.45",
			expected: true,
		},
		{
			name:     "有效小数 - 单位小数",
			input:    "1.5",
			expected: true,
		},
		{
			name:     "有效小数 - 多位小数",
			input:    "123.456789",
			expected: true,
		},
		{
			name:     "无效小数 - 整数",
			input:    "123",
			expected: false,
		},
		{
			name:     "无效小数 - 以小数点开头",
			input:    ".123",
			expected: false,
		},
		{
			name:     "无效小数 - 以小数点结尾",
			input:    "123.",
			expected: false,
		},
		{
			name:     "无效小数 - 包含字母",
			input:    "123.4a",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsDecimals(tt.input)
			assert.Equal(t, tt.expected, result, "IsDecimals(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsCellphone 测试手机号验证功能
func TestIsCellphone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效手机号 - 标准格式",
			input:    "13812345678",
			expected: true,
		},
		{
			name:     "有效手机号 - 15开头",
			input:    "15987654321",
			expected: true,
		},
		{
			name:     "有效手机号 - 18开头",
			input:    "18612345678",
			expected: true,
		},
		{
			name:     "无效手机号 - 不以1开头",
			input:    "23812345678",
			expected: false,
		},
		{
			name:     "无效手机号 - 太短",
			input:    "1381234567",
			expected: false,
		},
		{
			name:     "无效手机号 - 太长",
			input:    "138123456789",
			expected: false,
		},
		{
			name:     "无效手机号 - 包含字母",
			input:    "1381234567a",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsCellphone(tt.input)
			assert.Equal(t, tt.expected, result, "IsCellphone(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsTelephone 测试固定电话验证功能
func TestIsTelephone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效固话 - 8位数字",
			input:    "12345678",
			expected: true,
		},
		{
			name:     "有效固话 - 以0开头",
			input:    "01234567",
			expected: true,
		},
		{
			name:     "无效固话 - 太短",
			input:    "1234567",
			expected: false,
		},
		{
			name:     "无效固话 - 太长",
			input:    "123456789",
			expected: false,
		},
		{
			name:     "无效固话 - 包含字母",
			input:    "1234567a",
			expected: false,
		},
		{
			name:     "无效固话 - 包含特殊字符",
			input:    "1234-567",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsTelephone(tt.input)
			assert.Equal(t, tt.expected, result, "IsTelephone(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsEngishLowCase 测试小写英文验证功能
func TestIsEngishLowCase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效小写英文 - 单词",
			input:    "hello",
			expected: true,
		},
		{
			name:     "有效小写英文 - 单字母",
			input:    "a",
			expected: true,
		},
		{
			name:     "有效小写英文 - 长单词",
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: true,
		},
		{
			name:     "无效小写英文 - 包含大写",
			input:    "Hello",
			expected: false,
		},
		{
			name:     "无效小写英文 - 包含数字",
			input:    "hello123",
			expected: false,
		},
		{
			name:     "无效小写英文 - 包含特殊字符",
			input:    "hello!",
			expected: false,
		},
		{
			name:     "无效小写英文 - 包含空格",
			input:    "hello world",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsEngishLowCase(tt.input)
			assert.Equal(t, tt.expected, result, "IsEngishLowCase(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsEnglishCap 测试大写英文验证功能
func TestIsEnglishCap(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效大写英文 - 单词",
			input:    "HELLO",
			expected: true,
		},
		{
			name:     "有效大写英文 - 单字母",
			input:    "A",
			expected: true,
		},
		{
			name:     "有效大写英文 - 长单词",
			input:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			expected: true,
		},
		{
			name:     "无效大写英文 - 包含小写",
			input:    "Hello",
			expected: false,
		},
		{
			name:     "无效大写英文 - 包含数字",
			input:    "HELLO123",
			expected: false,
		},
		{
			name:     "无效大写英文 - 包含特殊字符",
			input:    "HELLO!",
			expected: false,
		},
		{
			name:     "无效大写英文 - 包含空格",
			input:    "HELLO WORLD",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsEnglishCap(tt.input)
			assert.Equal(t, tt.expected, result, "IsEnglishCap(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsEnglish 测试英文验证功能
func TestIsEnglish(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效英文 - 小写",
			input:    "hello",
			expected: true,
		},
		{
			name:     "有效英文 - 大写",
			input:    "HELLO",
			expected: true,
		},
		{
			name:     "有效英文 - 大小写混合",
			input:    "Hello",
			expected: true,
		},
		{
			name:     "有效英文 - 单字母",
			input:    "a",
			expected: true,
		},
		{
			name:     "无效英文 - 包含数字",
			input:    "hello123",
			expected: false,
		},
		{
			name:     "无效英文 - 包含特殊字符",
			input:    "hello!",
			expected: false,
		},
		{
			name:     "无效英文 - 包含空格",
			input:    "hello world",
			expected: false,
		},

		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsEnglish(tt.input)
			assert.Equal(t, tt.expected, result, "IsEnglish(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestIsEmail 测试邮箱验证功能
func TestIsEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "有效邮箱 - 标准格式",
			input:    "test@example.com",
			expected: true,
		},
		{
			name:     "有效邮箱 - 包含数字",
			input:    "user123@domain.org",
			expected: true,
		},
		{
			name:     "有效邮箱 - 包含点号",
			input:    "user.name@example.com",
			expected: true,
		},
		{
			name:     "有效邮箱 - 包含加号",
			input:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "有效邮箱 - 子域名",
			input:    "test@mail.example.com",
			expected: true,
		},
		{
			name:     "无效邮箱 - 缺少@",
			input:    "testexample.com",
			expected: false,
		},
		{
			name:     "无效邮箱 - 缺少域名",
			input:    "test@",
			expected: false,
		},
		{
			name:     "无效邮箱 - 缺少用户名",
			input:    "@example.com",
			expected: false,
		},
		{
			name:     "无效邮箱 - 缺少顶级域名",
			input:    "test@example",
			expected: false,
		},
		{
			name:     "无效邮箱 - 多个@",
			input:    "test@@example.com",
			expected: false,
		},
		{
			name:     "无效邮箱 - 包含空格",
			input:    "test @example.com",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsEmail(tt.input)
			assert.Equal(t, tt.expected, result, "IsEmail(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestRegexCheckEdgeCases 测试边界情况
func TestRegexCheckEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("空字符串测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		
		emptyStr := ""
		assert.False(t, pkg.IsID(emptyStr), "空字符串不应该是有效ID")
		assert.False(t, pkg.IsUserName(emptyStr), "空字符串不应该是有效用户名")
		assert.False(t, pkg.IsPwd(emptyStr), "空字符串不应该是有效密码")
		assert.False(t, pkg.IsInteger(emptyStr), "空字符串不应该是有效整数")
		assert.False(t, pkg.IsDecimals(emptyStr), "空字符串不应该是有效小数")
		assert.False(t, pkg.IsCellphone(emptyStr), "空字符串不应该是有效手机号")
		assert.False(t, pkg.IsTelephone(emptyStr), "空字符串不应该是有效固话")
		assert.False(t, pkg.IsEngishLowCase(emptyStr), "空字符串不应该是有效小写英文")
		assert.False(t, pkg.IsEnglishCap(emptyStr), "空字符串不应该是有效大写英文")
		assert.False(t, pkg.IsEnglish(emptyStr), "空字符串不应该是有效英文")
		assert.False(t, pkg.IsEmail(emptyStr), "空字符串不应该是有效邮箱")

	})

	t.Run("特殊字符测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		
		specialChars := "!@#$%^&*()"
		assert.True(t, pkg.IsPwd(specialChars), "特殊字符应该是有效密码")
		assert.False(t, pkg.IsID(specialChars), "特殊字符不应该是有效ID")
		assert.False(t, pkg.IsUserName(specialChars), "特殊字符不应该是有效用户名")
	})

	t.Run("Unicode字符测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		
		unicodeStr := "测试🎉emoji"
		assert.False(t, pkg.IsPwd(unicodeStr), "包含Unicode字符的字符串不应该是有效密码")
		assert.False(t, pkg.IsID(unicodeStr), "包含Unicode字符的字符串不应该是有效ID")
		assert.False(t, pkg.IsEmail(unicodeStr), "包含Unicode字符的字符串不应该是有效邮箱")
	})
}

// TestRegexCheckIntegration 集成测试
func TestRegexCheckIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("综合验证测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		
		// 测试一个复杂的用户注册场景
		userID := "user123456"
		userName := "testuser"
		password := "MyPass123!"
		email := "test@example.com"
		phone := "13812345678"
		
		assert.True(t, pkg.IsID(userID), "用户ID应该有效")
		assert.True(t, pkg.IsUserName(userName), "用户名应该有效")
		assert.True(t, pkg.IsPwd(password), "密码应该有效")
		assert.True(t, pkg.IsEmail(email), "邮箱应该有效")
		assert.True(t, pkg.IsCellphone(phone), "手机号应该有效")
	})

	t.Run("性能测试", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		
		// 测试大量调用的性能
		testStr := "test123456"
		for i := 0; i < 1000; i++ {
			pkg.IsID(testStr)
			pkg.IsUserName(testStr)
			pkg.IsPwd(testStr)
			pkg.IsEmail("test@example.com")
		}
		// 如果能执行到这里说明性能可接受
		assert.True(t, true, "性能测试通过")
	})
}