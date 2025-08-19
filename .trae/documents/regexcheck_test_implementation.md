# RegexCheck 模块测试实现方案

## 1. 模块概述

`regexcheck.go` 模块提供了一套完整的正则表达式验证功能，包含多种常用的数据格式验证方法。该模块是数据验证的核心组件，广泛用于用户输入验证、数据格式检查等场景。

### 1.1 核心功能

#### 结构体

* `RegexCheck` - 正则表达式验证器结构体

#### 验证方法

* `IsID()` - ID格式验证

* `IsUserName()` - 用户名格式验证

* `IsPwd()` - 密码格式验证

* `IsInteger()` - 整数格式验证

* `IsDecimals()` - 小数格式验证

* `IsCellphone()` - 手机号格式验证

* `IsTelephone()` - 电话号码格式验证

* `IsEngishLowCase()` - 英文小写验证

* `IsEnglishCap()` - 英文大写验证

* `IsEnglish()` - 英文字符验证

* `IsEmail()` - 邮箱格式验证

### 1.2 RegexCheck 结构体

```go
type RegexCheck struct {
    // 可能包含配置选项或缓存的正则表达式
}
```

## 2. 测试文件结构

### 2.1 文件位置

```
tests/unit/pkg/regexcheck_test.go
```

### 2.2 测试依赖

```go
package pkg

import (
    "strings"
    "testing"
    "unicode/utf8"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

## 3. RegexCheck 结构体测试

### 3.1 TestRegexCheck - 结构体初始化测试

```go
func TestRegexCheck(t *testing.T) {
    t.Run("创建RegexCheck实例", func(t *testing.T) {
        checker := &RegexCheck{}
        assert.NotNil(t, checker)
    })
    
    t.Run("多个实例独立性", func(t *testing.T) {
        checker1 := &RegexCheck{}
        checker2 := &RegexCheck{}
        
        assert.NotNil(t, checker1)
        assert.NotNil(t, checker2)
        // 验证实例独立性（如果有状态的话）
    })
}
```

## 4. IsID 方法测试

### 4.1 TestRegexCheck\_IsID - ID格式验证测试

```go
func TestRegexCheck_IsID(t *testing.T) {
    checker := &RegexCheck{}
    
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        // 有效ID格式
        {"纯数字ID", "123", true},
        {"长数字ID", "123456789", true},
        {"单个数字", "1", true},
        {"零", "0", true},
        
        // 数字字母组合（根据实际实现调整）
        {"字母数字组合", "abc123", true}, // 需要根据实际实现调整
        {"大写字母数字", "ABC123", true},
        {"混合大小写", "AbC123", true},
        
        // 无效ID格式
        {"空字符串", "", false},
        {"纯字母", "abc", false}, // 需要根据实际实现调整
        {"包含特殊字符", "123@abc", false},
        {"包含空格", "123 456", false},
        {"包含中文", "123中文", false},
        {"包含符号", "123-456", false},
        {"包含点号", "123.456", false},
        {"负数", "-123", false},
        {"浮点数", "123.45", false},
        
        // 边界情况
        {"很长的数字", "12345678901234567890", true},
        {"前导零", "0123", true},
        {"Unicode数字", "１２３", false}, // 全角数字
        {"十六进制", "0x123", false},
        {"八进制", "0123", true}, // 可能被当作普通数字
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := checker.IsID(tt.input)
            assert.Equal(t, tt.expected, result, "输入: %s", tt.input)
        })
    }
}
```

## 5. IsUserName 方法测试

### 5.1 TestRegexCheck\_IsUserName - 用户名格式验证测试

```go
func TestRegexCheck_IsUserName(t *testing.T) {
    checker := &RegexCheck{}
    
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        // 有效用户名
        {"纯字母", "username", true},
        {"字母数字", "user123", true},
        {"大写字母", "USERNAME", true},
        {"混合大小写", "UserName", true},
        {"字母数字混合", "User123Name", true},
        {"下划线", "user_name", true},
        {"数字开头", "123user", true}, // 根据实际规则调整
        {"单个字符", "a", true},
        {"最小长度", "ab", true}, // 根据实际最小长度调整
        
        // 无效用户名
        {"空字符串", "", false},
        {"包含空格", "user name", false},
        {"包含特殊字符@", "user@name", false},
        {"包含特殊字符#", "user#name", false},
        {"包含特殊字符$", "user$name", false},
        {"包含特殊字符%", "user%name", false},
        {"包含中文", "用户名", false},
        {"包含连字符", "user-name", false},
        {"包含点号", "user.name", false},
        {"包含加号", "user+name", false},
        {"包含等号", "user=name", false},
        
        // 边界情况
        {"很长的用户名", strings.Repeat("a", 100), true}, // 根据最大长度调整
        {"超长用户名", strings.Repeat("a", 1000), false}, // 根据最大长度调整
        {"纯下划线", "____", true}, // 根据实际规则调整
        {"下划线开头", "_username", true},
        {"下划线结尾", "username_", true},
        {"多个下划线", "user__name", true},
        
        // Unicode字符
        {"Unicode字母", "usernameé", false},
        {"表情符号", "user😀name", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := checker.IsUserName(tt.input)
            assert.Equal(t, tt.expected, result, "输入: %s", tt.input)
        })
    }
}
```

## 6. IsPwd 方法测试

### 6.1 TestRegexCheck\_IsPwd - 密码格式验证测试

```go
func TestRegexCheck_IsPwd(t *testing.T) {
    checker := &RegexCheck{}
    
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        // 有效密码（根据实际密码规则调整）
        {"简单密码", "password", true},
        {"字母数字", "password123", true},
        {"大小写数字", "Password123", true},
        {"包含特殊字符", "Password123!", true},
        {"复杂密码", "P@ssw0rd123", true},
        {"最小长度", "Pass123", true}, // 根据最小长度调整
        {"纯数字", "12345678", true}, // 根据实际规则调整
        {"纯字母", "password", true},
        {"纯大写", "PASSWORD", true},
        
        // 无效密码
        {"空字符串", "", false},
        {"太短", "123", false}, // 根据最小长度调整
        {"包含空格", "pass word", false},
        {"包含中文", "密码123", false},
        {"包含制表符", "pass\tword", false},
        {"包含换行符", "pass\nword", false},
        
        // 边界情况
        {"最大长度", strings.Repeat("a", 128), true}, // 根据最大长度调整
        {"超长密码", strings.Repeat("a", 1000), false},
        {"特殊字符组合", "!@#$%^&*()", true}, // 根据实际规则调整
        {"数字符号组合", "123!@#", true},
        
        // 常见密码模式
        {"键盘序列", "qwerty123", true},
        {"重复字符", "aaaa1111", true}, // 根据实际规则调整
        {"递增数字", "12345678", true},
        
        // Unicode字符
        {"Unicode字符", "pássword", false}, // 根据实际规则调整
        {"表情符号", "pass😀word", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := checker.IsPwd(tt.input)
            assert.Equal(t, tt.expected, result, "输入: %s", tt.input)
        })
    }
}
```

## 7. IsInteger 方法测试

### 7.1 TestRegexCheck\_IsInteger - 整数格式验证测试

```go
func TestRegexCheck_IsInteger(t *testing.T) {
    checker := &RegexCheck{}
    
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        // 有效整数
        {"正整数", "123", true},
        {"负整数", "-123", true},
        {"零", "0", true},
        {"负零", "-0", true},
        {"单个数字", "5", true},
        {"大整数", "123456789", true},
        {"前导零", "0123", true}, // 根据实际规则调整
        
        // 无效整数
        {"空字符串", "", false},
        {"浮点数", "123.45", false},
        {"包含字母", "123a", false},
        {"包含空格", "123 456", false},
        {"包含特殊字符", "123@", false},
        {"纯字母", "abc", false},
        {"多个负号", "--123", false},
        {"正号", "+123", false}, // 根据实际规则调整
        {"科学计数法", "1e5", false},
        {"十六进制", "0x123", false},
        {"八进制前缀", "0o123", false},
        {"二进制前缀", "0b101", false},
        
        // 边界情况
        {"很大的数", "999999999999999999", true},
        {"负的很大的数", "-999999999999999999", true},
        {"只有负号", "-", false},
        {"负号在末尾", "123-", false},
        {"中间有负号", "12-3", false},
        
        // Unicode数字
        {"全角数字", "１２３", false},
        {"阿拉伯数字", "٠١٢", false},
        {"中文数字", "一二三", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := checker.IsInteger(tt.input)
            assert.Equal(t, tt.expected, result, "输入: %s", tt.input)
        })
    }
}
```

## 8. IsDecimals 方法测试

### 8.1 TestRegexCheck\_IsDecimals - 小数格式验证测试

```go
func TestRegexCheck_IsDecimals(t *testing.T) {
    checker := &RegexCheck{}
    
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        // 有效小数
        {"正小数", "123.45", true},
        {"负小数", "-123.45", true},
        {"零点几", "0.5", true},
        {"负零点几", "-0.5", true},
        {"整数形式", "123", true}, // 根据实际规则调整
        {"负整数形式", "-123", true},
        {"小数点开头", ".5", true}, // 根据实际规则调整
        {"负小数点开头", "-.5", true},
        {"小数点结尾", "123.", true}, // 根据实际规则调整
        {"很多小数位", "123.123456789", true},
        {"前导零", "0123.45", true}, // 根据实际规则调整
        {"小数前导零", "123.0450", true},
        
        // 无效小数
        {"空字符串", "", false},
        {"纯字母", "abc", false},
        {"包含字
```

