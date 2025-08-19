package pkg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
)

// TestWriteCookie 测试写入cookie功能
func TestWriteCookie(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		cookieName string
		cookieValue string
		expectedPath string
		expectedMaxAge int
		expectedHttpOnly bool
	}{
		{
			name:        "写入普通cookie",
			cookieName:  "test_cookie",
			cookieValue: "test_value",
			expectedPath: "/",
			expectedMaxAge: 7 * 86400,
			expectedHttpOnly: true,
		},
		{
			name:        "写入空值cookie",
			cookieName:  "empty_cookie",
			cookieValue: "",
			expectedPath: "/",
			expectedMaxAge: 7 * 86400,
			expectedHttpOnly: true,
		},
		{
			name:        "写入特殊字符cookie",
			cookieName:  "special_cookie",
			cookieValue: "value with spaces",
			expectedPath: "/",
			expectedMaxAge: 7 * 86400,
			expectedHttpOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			// 创建响应记录器
			rw := httptest.NewRecorder()
			
			// 调用WriteCookie函数
			pkg.WriteCookie(rw, tt.cookieName, tt.cookieValue)
			
			// 获取设置的cookie
			result := rw.Result()
			cookies := result.Cookies()
			
			// 验证cookie是否正确设置
			assert.Len(t, cookies, 1, "应该设置一个cookie")
			
			cookie := cookies[0]
			assert.Equal(t, tt.cookieName, cookie.Name, "cookie名称应该匹配")
			assert.Equal(t, tt.cookieValue, cookie.Value, "cookie值应该匹配")
			assert.Equal(t, tt.expectedPath, cookie.Path, "cookie路径应该匹配")
			assert.Equal(t, tt.expectedMaxAge, cookie.MaxAge, "cookie最大年龄应该匹配")
			assert.Equal(t, tt.expectedHttpOnly, cookie.HttpOnly, "cookie HttpOnly属性应该匹配")
		})
	}
}

// TestDelCookie 测试删除cookie功能
func TestDelCookie(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name       string
		cookieName string
	}{
		{
			name:       "删除普通cookie",
			cookieName: "test_cookie",
		},
		{
			name:       "删除不存在的cookie",
			cookieName: "nonexistent_cookie",
		},
		{
			name:       "删除空名称cookie",
			cookieName: "",
		},
	}

	for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				select {
				case <-ctx.Done():
					t.Fatal("测试超时")
				default:
				}

				// 创建响应记录器
				rw := httptest.NewRecorder()
				
				// 调用DelCookie函数
				pkg.DelCookie(rw, tt.cookieName)
				
				// 获取设置的cookie
				result := rw.Result()
				cookies := result.Cookies()
				
				// 对于空名称的cookie，可能不会设置删除cookie
				if tt.cookieName == "" {
					// 空名称cookie可能不会被处理
					assert.Len(t, cookies, 0, "空名称cookie不应该设置删除cookie")
				} else {
					// 验证cookie是否正确设置为删除
					assert.Len(t, cookies, 1, "应该设置一个cookie用于删除")
					
					cookie := cookies[0]
					assert.Equal(t, tt.cookieName, cookie.Name, "cookie名称应该匹配")
					assert.Equal(t, -1, cookie.MaxAge, "cookie MaxAge应该为-1以删除cookie")
				}
			})
		}
}

// TestReadCookie 测试读取cookie功能
func TestReadCookie(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name         string
		cookieName   string
		cookieValue  string
		setCookie    bool
		expectedValue string
	}{
		{
			name:         "读取存在的cookie",
			cookieName:   "test_cookie",
			cookieValue:  "test_value",
			setCookie:    true,
			expectedValue: "test_value",
		},
		{
			name:         "读取不存在的cookie",
			cookieName:   "nonexistent_cookie",
			cookieValue:  "",
			setCookie:    false,
			expectedValue: "",
		},
		{
			name:         "读取空值cookie",
			cookieName:   "empty_cookie",
			cookieValue:  "",
			setCookie:    true,
			expectedValue: "",
		},
		{
			name:         "读取特殊字符cookie",
			cookieName:   "special_cookie",
			cookieValue:  "value with spaces",
			setCookie:    true,
			expectedValue: "value with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			// 创建HTTP请求
			req := httptest.NewRequest("GET", "/", nil)
			
			// 如果需要设置cookie，则添加到请求中
			if tt.setCookie {
				req.AddCookie(&http.Cookie{
					Name:  tt.cookieName,
					Value: tt.cookieValue,
				})
			}
			
			// 调用ReadCookie函数
			result := pkg.ReadCookie(req, tt.cookieName)
			
			// 验证结果
			assert.Equal(t, tt.expectedValue, result, "读取的cookie值应该匹配")
		})
	}
}

// TestSetCookieDefVal 测试设置cookie默认值功能
func TestSetCookieDefVal(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		domain   string
		path     string
		maxAge   int
		httpOnly bool
	}{
		{
			name:     "设置基本默认值",
			domain:   "example.com",
			path:     "/app",
			maxAge:   30,
			httpOnly: false,
		},
		{
			name:     "设置空域名",
			domain:   "",
			path:     "/",
			maxAge:   1,
			httpOnly: true,
		},
		{
			name:     "设置长期cookie",
			domain:   "long.example.com",
			path:     "/long",
			maxAge:   365,
			httpOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}

			// 调用SetCookieDefVal函数
			pkg.SetCookieDefVal(tt.domain, tt.path, tt.maxAge, tt.httpOnly)
			
			// 创建一个测试cookie来验证设置
			rw := httptest.NewRecorder()
			pkg.WriteCookie(rw, "test", "value")
			
			result := rw.Result()
			cookies := result.Cookies()
			assert.Len(t, cookies, 1, "应该设置一个cookie")
			
			cookie := cookies[0]
			assert.Equal(t, tt.path, cookie.Path, "路径应该匹配")
			assert.Equal(t, tt.maxAge*86400, cookie.MaxAge, "最大年龄应该匹配（转换为秒）")
			assert.Equal(t, tt.httpOnly, cookie.HttpOnly, "HttpOnly属性应该匹配")
		})
	}
}

// TestCookieIntegration 测试cookie功能的集成测试
func TestCookieIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("完整的cookie生命周期", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		// 创建响应记录器
		rw := httptest.NewRecorder()
		
		// 1. 写入cookie
		cookieName := "integration_test"
		cookieValue := "integration_value"
		pkg.WriteCookie(rw, cookieName, cookieValue)
		
		// 2. 模拟浏览器接收cookie并在下次请求中发送
		result := rw.Result()
		cookies := result.Cookies()
		assert.Len(t, cookies, 1, "应该设置一个cookie")
		
		// 将cookie添加到新请求中
		newReq := httptest.NewRequest("GET", "/", nil)
		newReq.AddCookie(cookies[0])
		
		// 3. 读取cookie
		readValue := pkg.ReadCookie(newReq, cookieName)
		assert.Equal(t, cookieValue, readValue, "读取的cookie值应该与写入的值匹配")
		
		// 4. 删除cookie
		delRw := httptest.NewRecorder()
		pkg.DelCookie(delRw, cookieName)
		
		delResult := delRw.Result()
		delCookies := delResult.Cookies()
		assert.Len(t, delCookies, 1, "应该设置一个删除cookie")
		assert.Equal(t, -1, delCookies[0].MaxAge, "删除cookie的MaxAge应该为-1")
	})
}

// TestCookieEdgeCases 测试cookie的边界情况
func TestCookieEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("处理特殊字符", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		rw := httptest.NewRecorder()
		
		// 测试包含特殊字符的cookie名称和值
		specialName := "test-cookie_123"
		specialValue := "value@#$%^&*()"
		
		pkg.WriteCookie(rw, specialName, specialValue)
		
		result := rw.Result()
		cookies := result.Cookies()
		assert.Len(t, cookies, 1, "应该设置一个cookie")
		assert.Equal(t, specialName, cookies[0].Name, "特殊字符cookie名称应该正确")
	})
	
	t.Run("处理长cookie值", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}

		rw := httptest.NewRecorder()
		
		// 创建一个长cookie值
		longValue := ""
		for i := 0; i < 1000; i++ {
			longValue += "a"
		}
		
		pkg.WriteCookie(rw, "long_cookie", longValue)
		
		result := rw.Result()
		cookies := result.Cookies()
		assert.Len(t, cookies, 1, "应该设置一个长cookie")
		assert.Equal(t, longValue, cookies[0].Value, "长cookie值应该正确")
	})
}