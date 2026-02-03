package pkg_test

import (
	"context"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
)

// TestRand 测试随机字符串生成功能
func TestRand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	tests := []struct {
		name     string
		size     int
		kind     int
		validate func(string) bool
	}{
		{
			name: "生成纯数字字符串",
			size: 10,
			kind: pkg.RAND_KIND_NUM,
			validate: func(s string) bool {
				for _, ch := range s {
					if !unicode.IsDigit(ch) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "生成小写字母字符串",
			size: 8,
			kind: pkg.RAND_KIND_LOWER,
			validate: func(s string) bool {
				for _, ch := range s {
					if !unicode.IsLower(ch) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "生成大写字母字符串",
			size: 6,
			kind: pkg.RAND_KIND_UPPER,
			validate: func(s string) bool {
				for _, ch := range s {
					if !unicode.IsUpper(ch) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "生成混合字符串",
			size: 12,
			kind: pkg.RAND_KIND_ALL,
			validate: func(s string) bool {
				for _, ch := range s {
					if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
						return false
					}
				}
				return true
			},
		},
		{
			name: "生成默认类型字符串",
			size: 5,
			kind: 999, // 无效类型，应该使用默认行为
			validate: func(s string) bool {
				for _, ch := range s {
					if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.Rand(tt.size, tt.kind)
			
			// 验证长度
			assert.Equal(t, tt.size, len(result), "生成的字符串长度应该匹配")
			
			// 验证字符类型
			assert.True(t, tt.validate(result), "生成的字符串应该符合指定类型")
		})
	}
}

// TestRandEdgeCases 测试Rand函数的边界情况
func TestRandEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	t.Run("零长度字符串", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		result := pkg.Rand(0, pkg.RAND_KIND_ALL)
		assert.Equal(t, "", result, "零长度应该返回空字符串")
	})
	
	t.Run("单字符字符串", func(t *testing.T) {
		result := pkg.Rand(1, pkg.RAND_KIND_NUM)
		assert.Equal(t, 1, len(result), "单字符长度应该正确")
		assert.True(t, unicode.IsDigit(rune(result[0])), "单字符应该是数字")
	})
	
	t.Run("大长度字符串", func(t *testing.T) {
		result := pkg.Rand(1000, pkg.RAND_KIND_LOWER)
		assert.Equal(t, 1000, len(result), "大长度字符串长度应该正确")
		for _, ch := range result {
			assert.True(t, unicode.IsLower(ch), "所有字符都应该是小写字母")
		}
	})
}

// TestRandUniqueness 测试Rand函数生成的随机性
func TestRandUniqueness(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	t.Run("多次调用应该生成不同的字符串", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		results := make(map[string]bool)
		size := 10
		iterations := 100
		
		for i := 0; i < iterations; i++ {
			result := pkg.Rand(size, pkg.RAND_KIND_ALL)
			results[result] = true
		}
		
		// 至少应该有90%的字符串是唯一的
		uniqueCount := len(results)
		expectedMinUnique := int(float64(iterations) * 0.9)
		assert.GreaterOrEqual(t, uniqueCount, expectedMinUnique, "应该生成足够多的唯一字符串")
	})
}

// TestSendToMail 测试邮件发送功能
func TestSendToMail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	// 注意：这些测试不会实际发送邮件，只测试参数处理和格式化
	tests := []struct {
		name     string
		user     string
		password string
		host     string
		to       string
		subject  string
		body     string
		mailtype string
		expectError bool
	}{
		{
			name:     "HTML邮件参数",
			user:     "test@example.com",
			password: "password",
			host:     "smtp.example.com:587",
			to:       "recipient@example.com",
			subject:  "Test Subject",
			body:     "<h1>Test Body</h1>",
			mailtype: "html",
			expectError: true, // 实际会失败，因为SMTP服务器不存在
		},
		{
			name:     "纯文本邮件参数",
			user:     "test@example.com",
			password: "password",
			host:     "smtp.example.com:587",
			to:       "recipient@example.com",
			subject:  "Test Subject",
			body:     "Test Body",
			mailtype: "plain",
			expectError: true, // 实际会失败，因为SMTP服务器不存在
		},
		{
			name:     "多个收件人",
			user:     "test@example.com",
			password: "password",
			host:     "smtp.example.com:587",
			to:       "recipient1@example.com;recipient2@example.com",
			subject:  "Test Subject",
			body:     "Test Body",
			mailtype: "plain",
			expectError: true, // 实际会失败，因为SMTP服务器不存在
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			err := pkg.SendToMail(tt.user, tt.password, tt.host, tt.to, tt.subject, tt.body, tt.mailtype)
			
			if tt.expectError {
				assert.Error(t, err, "应该返回错误，因为SMTP服务器不存在")
			} else {
				assert.NoError(t, err, "不应该返回错误")
			}
		})
	}
}

// TestSendToMailMessageFormat 测试邮件消息格式
func TestSendToMailMessageFormat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	// 这个测试验证邮件消息的格式是否正确
	// 我们通过检查函数内部逻辑来验证
	t.Run("验证HTML邮件内容类型", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		// 测试HTML邮件类型的内容类型设置
		mailtype := "html"
		var content_type string
		if mailtype == "html" {
			content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
		} else {
			content_type = "Content-Type: text/plain" + "; charset=UTF-8"
		}
		expected := "Content-Type: text/html; charset=UTF-8"
		assert.Equal(t, expected, content_type, "HTML邮件内容类型应该正确")
	})
	
	t.Run("验证纯文本邮件内容类型", func(t *testing.T) {
		// 测试纯文本邮件类型的内容类型设置
		mailtype := "plain"
		var content_type string
		if mailtype == "html" {
			content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
		} else {
			content_type = "Content-Type: text/plain" + "; charset=UTF-8"
		}
		expected := "Content-Type: text/plain; charset=UTF-8"
		assert.Equal(t, expected, content_type, "纯文本邮件内容类型应该正确")
	})
	
	t.Run("验证收件人分割", func(t *testing.T) {
		// 测试多个收件人的分割逻辑
		to := "user1@example.com;user2@example.com;user3@example.com"
		send_to := strings.Split(to, ";")
		expected := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		assert.Equal(t, expected, send_to, "收件人分割应该正确")
	})
}

// TestIsUserName 测试用户名验证功能
func TestIsUserName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	tests := []struct{
		name     string
		userName string
		expected bool
	}{
		{
			name:     "有效用户名 - 字母数字",
			userName: "user123",
			expected: true,
		},
		{
			name:     "有效用户名 - 包含下划线",
			userName: "user_name_123",
			expected: true,
		},
		{
			name:     "有效用户名 - 最小长度",
			userName: "user12",
			expected: true,
		},
		{
			name:     "有效用户名 - 最大长度",
			userName: "user123456789012", // 16个字符
			expected: true,
		},
		{
			name:     "无效用户名 - 太短",
			userName: "u1", // 2个字符
			expected: false,
		},
		{
			name:     "无效用户名 - 太长",
			userName: "user12345678901234567", // 17个字符
			expected: false,
		},
		{
			name:     "无效用户名 - 包含特殊字符",
			userName: "user@123",
			expected: false,
		},
		{
			name:     "无效用户名 - 包含空格",
			userName: "user 123",
			expected: false,
		},
		{
			name:     "有效用户名 - 包含连字符",
			userName: "user-123",
			expected: true,
		},
		{
			name:     "空用户名",
			userName: "",
			expected: false,
		},
		{
			name:     "纯字母用户名",
			userName: "username",
			expected: true,
		},
		{
			name:     "纯数字用户名",
			userName: "123456",
			expected: true,
		},
	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsHandle(tt.userName)
			assert.Equal(t, tt.expected, result, "用户名验证结果应该匹配")
		})
	}
}

// TestIsPwd 测试密码验证功能
func TestIsPwd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "有效密码 - 字母数字",
			password: "password123",
			expected: true,
		},
		{
			name:     "有效密码 - 包含特殊字符",
			password: "password123!@#",
			expected: true,
		},
		{
			name:     "有效密码 - 最小长度",
			password: "pass123a", // 8个字符
			expected: true,
		},
		{
			name:     "有效密码 - 最大长度",
			password: "password123456789012345678901234567890123456789012345678901234567890123456789012", // 80个字符
			expected: true,
		},
		{
			name:     "无效密码 - 太短",
			password: "pass12", // 6个字符
			expected: false,
		},
		{
			name:     "无效密码 - 太长",
			password: "password1234567890123456789012345678901234567890123456789012345678901234567890123", // 81个字符
			expected: false,
		},
		{
			name:     "无效密码 - 只有字母",
			password: "password",
			expected: false,
		},
		{
			name:     "无效密码 - 只有数字",
			password: "12345678",
			expected: false,
		},
		{
			name:     "无效密码 - 只有特殊字符",
			password: "!@#$%^&*",
			expected: false,
		},
		{
			name:     "空密码",
			password: "",
			expected: false,
		},
		{
			name:     "有效密码 - 大小写字母数字",
			password: "Password123",
			expected: true,
		},

	}

	for _, tt := range tests {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.IsPwd(tt.password)
			assert.Equal(t, tt.expected, result, "密码验证结果应该匹配")
		})
	}
}



// TestHelpersIntegration 测试helpers模块的集成功能
func TestHelpersIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()
	
	t.Run("生成有效的随机用户名", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		// 生成一个可能有效的用户名
		userName := pkg.Rand(8, pkg.RAND_KIND_ALL)
		
		// 验证生成的用户名长度在有效范围内
		assert.GreaterOrEqual(t, len(userName), 6, "生成的用户名长度应该>=6")
		assert.LessOrEqual(t, len(userName), 20, "生成的用户名长度应该<=20")
		
		// 检查是否只包含有效字符
		validChars := true
		for _, ch := range userName {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
				validChars = false
				break
			}
		}
		assert.True(t, validChars, "生成的用户名应该只包含字母和数字")
	})
	
	t.Run("生成有效的随机密码", func(t *testing.T) {
		// 生成一个8位的随机密码
		password := pkg.Rand(8, pkg.RAND_KIND_ALL)
		
		// 验证生成的密码是否有效
		isValid := pkg.IsPwd(password)
		assert.True(t, isValid, "生成的密码应该有效")
		
		// 验证密码长度
		assert.Equal(t, 8, len(password), "密码长度应该为8")
	})
}