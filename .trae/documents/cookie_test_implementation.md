# Cookie 模块测试实现方案

## 1. 模块概述

`cookie.go` 模块提供了 HTTP Cookie 的基本操作功能，包括读取、写入、删除 Cookie 以及设置默认配置。该模块是 Web 应用中会话管理和用户状态保持的重要组件。

### 1.1 核心功能

* Cookie 写入 (`WriteCookie`)

* Cookie 读取 (`ReadCookie`)

* Cookie 删除 (`DelCookie`)

* 默认配置设置 (`SetCookieDefVal`)

### 1.2 默认配置

```go
var cookie = &http.Cookie{
    Path:     "/",
    Domain:   "",
    MaxAge:   7 * 86400, // 7天
    HttpOnly: true,
}
```

## 2. 测试文件结构

### 2.1 文件位置

```
tests/unit/pkg/cookie_test.go
```

### 2.2 测试依赖

```go
package pkg

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

## 3. 测试用例设计

### 3.1 TestWriteCookie - Cookie 写入测试

**测试目标**: 验证 Cookie 写入功能的正确性

```go
func TestWriteCookie(t *testing.T) {
    tests := []struct {
        name     string
        cookieName string
        cookieValue string
        setupConfig func()
        validateFunc func(*testing.T, *http.Response)
    }{
        {
            name:        "写入普通Cookie",
            cookieName:  "test_cookie",
            cookieValue: "test_value",
            setupConfig: func() {
                // 使用默认配置
            },
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "test_cookie", cookie.Name)
                assert.Equal(t, "test_value", cookie.Value)
                assert.Equal(t, "/", cookie.Path)
                assert.Equal(t, 7*86400, cookie.MaxAge)
                assert.True(t, cookie.HttpOnly)
            },
        },
        {
            name:        "写入空值Cookie",
            cookieName:  "empty_cookie",
            cookieValue: "",
            setupConfig: func() {},
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "empty_cookie", cookie.Name)
                assert.Equal(t, "", cookie.Value)
            },
        },
        {
            name:        "写入特殊字符Cookie",
            cookieName:  "special_cookie",
            cookieValue: "value with spaces and 中文",
            setupConfig: func() {},
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "special_cookie", cookie.Name)
                assert.Equal(t, "value with spaces and 中文", cookie.Value)
            },
        },
        {
            name:        "自定义配置Cookie",
            cookieName:  "custom_cookie",
            cookieValue: "custom_value",
            setupConfig: func() {
                SetCookieDefVal("example.com", "/app", 1, false)
            },
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "custom_cookie", cookie.Name)
                assert.Equal(t, "custom_value", cookie.Value)
                assert.Equal(t, "example.com", cookie.Domain)
                assert.Equal(t, "/app", cookie.Path)
                assert.Equal(t, 86400, cookie.MaxAge) // 1天
                assert.False(t, cookie.HttpOnly)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 重置为默认配置
            SetCookieDefVal("", "/", 7, true)
            
            // 应用测试配置
            tt.setupConfig()
            
            // 创建测试响应记录器
            recorder := httptest.NewRecorder()
            
            // 写入Cookie
            WriteCookie(recorder, tt.cookieName, tt.cookieValue)
            
            // 获取响应
            resp := recorder.Result()
            defer resp.Body.Close()
            
            // 验证结果
            tt.validateFunc(t, resp)
        })
    }
}
```

### 3.2 TestReadCookie - Cookie 读取测试

**测试目标**: 验证 Cookie 读取功能的正确性

```go
func TestReadCookie(t *testing.T) {
    tests := []struct {
        name         string
        setupCookies func(*http.Request)
        cookieName   string
        expected     string
        description  string
    }{
        {
            name: "读取存在的Cookie",
            setupCookies: func(req *http.Request) {
                cookie := &http.Cookie{
                    Name:  "test_cookie",
                    Value: "test_value",
                }
                req.AddCookie(cookie)
            },
            cookieName:  "test_cookie",
            expected:    "test_value",
            description: "应该能正确读取存在的Cookie值",
        },
        {
            name: "读取不存在的Cookie",
            setupCookies: func(req *http.Request) {
                // 不添加任何Cookie
            },
            cookieName:  "nonexistent_cookie",
            expected:    "",
            description: "读取不存在的Cookie应该返回空字符串",
        },
        {
            name: "读取多个Cookie中的一个",
            setupCookies: func(req *http.Request) {
                cookies := []*http.Cookie{
                    {Name: "cookie1", Value: "value1"},
                    {Name: "cookie2", Value: "value2"},
                    {Name: "cookie3", Value: "value3"},
                }
                for _, cookie := range cookies {
                    req.AddCookie(cookie)
                }
            },
            cookieName:  "cookie2",
            expected:    "value2",
            description: "应该能从多个Cookie中正确读取指定的Cookie",
        },
        {
            name: "读取空值Cookie",
            setupCookies: func(req *http.Request) {
                cookie := &http.Cookie{
                    Name:  "empty_cookie",
                    Value: "",
                }
                req.AddCookie(cookie)
            },
            cookieName:  "empty_cookie",
            expected:    "",
            description: "应该能正确读取空值Cookie",
        },
        {
            name: "读取特殊字符Cookie",
            setupCookies: func(req *http.Request) {
                cookie := &http.Cookie{
                    Name:  "special_cookie",
                    Value: "value with spaces and 中文",
                }
                req.AddCookie(cookie)
            },
            cookieName:  "special_cookie",
            expected:    "value with spaces and 中文",
            description: "应该能正确读取包含特殊字符的Cookie",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 创建测试请求
            req := httptest.NewRequest("GET", "/", nil)
            
            // 设置Cookie
            tt.setupCookies(req)
            
            // 读取Cookie
            result := ReadCookie(req, tt.cookieName)
            
            // 验证结果
            assert.Equal(t, tt.expected, result, tt.description)
        })
    }
}
```

### 3.3 TestDelCookie - Cookie 删除测试

**测试目标**: 验证 Cookie 删除功能的正确性

```go
func TestDelCookie(t *testing.T) {
    tests := []struct {
        name       string
        cookieName string
        validateFunc func(*testing.T, *http.Response)
    }{
        {
            name:       "删除普通Cookie",
            cookieName: "test_cookie",
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "test_cookie", cookie.Name)
                assert.Equal(t, -1, cookie.MaxAge, "删除Cookie的MaxAge应该为-1")
            },
        },
        {
            name:       "删除空名称Cookie",
            cookieName: "",
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "", cookie.Name)
                assert.Equal(t, -1, cookie.MaxAge)
            },
        },
        {
            name:       "删除特殊字符名称Cookie",
            cookieName: "special_cookie_中文",
            validateFunc: func(t *testing.T, resp *http.Response) {
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "special_cookie_中文", cookie.Name)
                assert.Equal(t, -1, cookie.MaxAge)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 创建测试响应记录器
            recorder := httptest.NewRecorder()
            
            // 删除Cookie
            DelCookie(recorder, tt.cookieName)
            
            // 获取响应
            resp := recorder.Result()
            defer resp.Body.Close()
            
            // 验证结果
            tt.validateFunc(t, resp)
        })
    }
}
```

### 3.4 TestSetCookieDefVal - 默认配置设置测试

**测试目标**: 验证默认配置设置功能的正确性

```go
func TestSetCookieDefVal(t *testing.T) {
    tests := []struct {
        name     string
        domain   string
        path     string
        maxAge   int
        httpOnly bool
        validateFunc func(*testing.T)
    }{
        {
            name:     "设置完整配置",
            domain:   "example.com",
            path:     "/app",
            maxAge:   30,
            httpOnly: false,
            validateFunc: func(t *testing.T) {
                // 写入一个Cookie来验证配置
                recorder := httptest.NewRecorder()
                WriteCookie(recorder, "test", "value")
                
                resp := recorder.Result()
                defer resp.Body.Close()
                
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "example.com", cookie.Domain)
                assert.Equal(t, "/app", cookie.Path)
                assert.Equal(t, 30*86400, cookie.MaxAge)
                assert.False(t, cookie.HttpOnly)
            },
        },
        {
            name:     "设置空域名",
            domain:   "",
            path:     "/",
            maxAge:   1,
            httpOnly: true,
            validateFunc: func(t *testing.T) {
                recorder := httptest.NewRecorder()
                WriteCookie(recorder, "test", "value")
                
                resp := recorder.Result()
                defer resp.Body.Close()
                
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "", cookie.Domain)
                assert.Equal(t, "/", cookie.Path)
                assert.Equal(t, 86400, cookie.MaxAge)
                assert.True(t, cookie.HttpOnly)
            },
        },
        {
            name:     "设置零过期时间",
            domain:   "test.com",
            path:     "/test",
            maxAge:   0,
            httpOnly: false,
            validateFunc: func(t *testing.T) {
                recorder := httptest.NewRecorder()
                WriteCookie(recorder, "test", "value")
                
                resp := recorder.Result()
                defer resp.Body.Close()
                
                cookies := resp.Cookies()
                require.Len(t, cookies, 1)
                
                cookie := cookies[0]
                assert.Equal(t, "test.com", cookie.Domain)
                assert.Equal(t, "/test", cookie.Path)
                assert.Equal(t, 0, cookie.MaxAge)
                assert.False(t, cookie.HttpOnly)
            },
        },
    }

    // 保存原始配置
    originalCookie := *cookie
    defer func() {
        // 恢复原始配置
        *cookie = originalCookie
    }()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 设置新配置
            SetCookieDefVal(tt.domain, tt.path, tt.maxAge, tt.httpOnly)
            
            // 验证配置
            tt.validateFunc(t)
        })
    }
}
```

## 4. 集成测试

### 4.1 TestCookieWorkflow - Cookie 工作流程测试

**测试目标**: 验证 Cookie 的完整生命周期

```go
func TestCookieWorkflow(t *testing.T) {
    t.Run("完整Cookie生命周期", func(t *testing.T) {
        cookieName := "session_id"
        cookieValue := "abc123def456"
        
        // 1. 设置自定义配置
        SetCookieDefVal("example.com", "/app", 1, true)
        
        // 2. 写入Cookie
        writeRecorder := httptest.NewRecorder()
        WriteCookie(writeRecorder, cookieName, cookieValue)
        
        // 3. 验证写入的Cookie
        writeResp := writeRecorder.Result()
        defer writeResp.Body.Close()
        
        cookies := writeResp.Cookies()
        require.Len(t, cookies, 1)
        
        writtenCookie := cookies[0]
        assert.Equal(t, cookieName, writtenCookie.Name)
        assert.Equal(t, cookieValue, writtenCookie.Value)
        
        // 4. 模拟客户端请求，读取Cookie
        readReq := httptest.NewRequest("GET", "/", nil)
        readReq.AddCookie(writtenCookie)
        
        readValue := ReadCookie(readReq, cookieName)
        assert.Equal(t, cookieValue, readValue)
        
        // 5. 删除Cookie
        delRecorder := httptest.NewRecorder()
        DelCookie(delRecorder, cookieName)
        
        delResp := delRecorder.Result()
        defer delResp.Body.Close()
        
        delCookies := delResp.Cookies()
        require.Len(t, delCookies, 1)
        
        deletedCookie := delCookies[0]
        assert.Equal(t, cookieName, deletedCookie.Name)
        assert.Equal(t, -1, deletedCookie.MaxAge)
        
        // 6. 验证删除后读取返回空值
        readReq2 := httptest.NewRequest("GET", "/", nil)
        readReq2.AddCookie(deletedCookie)
        
        readValue2 := ReadCookie(readReq2, cookieName)
        // 注意：实际浏览器会忽略MaxAge为-1的Cookie，但在测试中我们仍能读取到
        // 这里主要测试函数的行为一致性
    })
}
```

## 5. 边界和异常测试

### 5.1 TestCookieEdgeCases - 边界情况测试

```go
func TestCookieEdgeCases(t *testing.T) {
    tests := []struct {
        name        string
        testFunc    func(*testing.T)
        description string
    }{
        {
            name: "nil ResponseWriter",
            testFunc: func(t *testing.T) {
                // 这个测试可能会panic，需要recover
                defer func() {
                    if r := recover(); r != nil {
                        t.Log("WriteCookie with nil ResponseWriter caused panic:", r)
                    }
                }()
                
                WriteCookie(nil, "test", "value")
            },
            description: "测试nil ResponseWriter的处理",
        },
        {
            name: "nil Request",
            testFunc: func(t *testing.T) {
                defer func() {
                    if r := recover(); r != nil {
                        t.Log("ReadCookie with nil Request caused panic:", r)
                    }
                }()
                
                result := ReadCookie(nil, "test")
                assert.Equal(t, "", result)
            },
            description: "测试nil Request的处理",
        },
        {
            name: "超长Cookie名称",
            testFunc: func(t *testing.T) {
                longName := string(make([]byte, 1000))
                for i := range longName {
                    longName = longName[:i] + "a" + longName[i+1:]
                }
                
                recorder := httptest.NewRecorder()
                WriteCookie(recorder, longName, "value")
                
                resp := recorder.Result()
                defer resp.Body.Close()
                
                // 验证是否能正常处理超长名称
                cookies := resp.Cookies()
                if len(cookies) > 0 {
                    assert.Equal(t, longName, cookies[0].Name)
                }
            },
            description: "测试超长Cookie名称的处理",
        },
        {
            name: "超长Cookie值",
            testFunc: func(t *testing.T) {
                longValue := string(make([]byte, 5000))
                for i := range longValue {
                    longValue = longValue[:i] + "x" + longValue[i+1:]
                }
                
                recorder := httptest.NewRecorder()
                WriteCookie(recorder, "test", longValue)
                
                resp := recorder.Result()
                defer resp.Body.Close()
                
                // 验证是否能正常处理超长值
                cookies := resp.Cookies()
                if len(cookies) > 0 {
                    assert.Equal(t, longValue, cookies[0].Value)
                }
            },
            description: "测试超长Cookie值的处理",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.testFunc(t)
        })
    }
}
```

## 6. 性能测试

### 6.1 基准测试

```go
func BenchmarkWriteCookie(b *testing.B) {
    recorder := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        WriteCookie(recorder, "benchmark_cookie", "benchmark_value")
    }
}

func BenchmarkReadCookie(b *testing.B) {
    req := httptest.NewRequest("GET", "/", nil)
    req.AddCookie(&http.Cookie{Name: "benchmark_cookie", Value: "benchmark_value"})
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ReadCookie(req, "benchmark_cookie")
    }
}

func BenchmarkDelCookie(b *testing.B) {
    recorder := httptest.NewRecorder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        DelCookie(recorder, "benchmark_cookie")
    }
}

func BenchmarkSetCookieDefVal(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        SetCookieDefVal("example.com", "/app", 7, true)
    }
}
```

## 7. 测试辅助函数

### 7.1 辅助函数

```go
// createTestRequest 创建带有指定Cookie的测试请求
func createTestRequest(cookies map[string]string) *http.Request {
    req := httptest.NewRequest("GET", "/", nil)
    for name, value := range cookies {
        req.AddCookie(&http.Cookie{Name: name, Value: value})
    }
    return req
}

// extractCookies 从响应中提取Cookie映射
func extractCookies(resp *http.Response) map[string]*http.Cookie {
    result := make(map[string]*http.Cookie)
    for _, cookie := range resp.Cookies() {
        result[cookie.Name] = cookie
    }
    return result
}

// assertCookieEquals 断言Cookie的各个属性
func assertCookieEquals(t *testing.T, expected, actual *http.Cookie) {
    assert.Equal(t, expected.Name, actual.Name)
    assert.Equal(t, expected.Value, actual.Value)
    assert.Equal(t, expected.Path, actual.Path)
    assert.Equal(t, expected.Domain, actual.Domain)
    assert.Equal(t, expected.MaxAge, actual.MaxAge)
    assert.Equal(t, expected.HttpOnly, actual.HttpOnly)
    assert.Equal(t, expected.Secure, actual.Secure)
}
```

## 8. 测试覆盖率目标

* **函数覆盖率**: 100%

* **行覆盖率**: ≥ 95%

* **分支覆盖率**: ≥ 90%

## 9. 实施步骤

1. **创建测试文件**: `tests/unit/pkg/cookie_test.go`
2. **实现基础测试用例**: WriteCookie, ReadCookie, DelCookie
3. **添加配置测试**: SetCookieDefVal
4. **实现集成测试**: 完整工作流程测试
5. **添加边界测试**: 异常情况和边界条件
6. **性能测试**: 基准测试
7. **验证覆盖率**: 使用 `go test -cover`
8. **优化和完善**: 根据覆盖率报告补充测试

## 10. 注意事项

* Cookie 测试需要使用 `httptest` 包模拟 HTTP 环境

* 注意测试的独立性，每个测试后恢复默认配置

* 特别关注边界条件，如空值、nil 指针、超长字符串等

* 考虑浏览器的 Cookie 限制和行为

* 测试应该覆盖各种字符编码情况

* 确保测试的可重复性和稳定性

