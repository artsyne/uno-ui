# Helpers 模块测试实现方案

## 1. 模块概述

`helpers.go` 模块提供了一系列通用的辅助函数，包括字符串处理、数组操作、类型转换、随机数生成等功能。这些函数是项目中广泛使用的基础工具函数。

### 1.1 核心功能分类

#### 字符串处理

* `Substr()` - 字符串截取

* `InArray()` - 检查元素是否在数组中

* `Implode()` - 数组连接为字符串

* `Explode()` - 字符串分割为数组

* `Ucfirst()` - 首字母大写

* `Lcfirst()` - 首字母小写

* `Addslashes()` - 添加转义字符

* `Stripslashes()` - 移除转义字符

* `Quotemeta()` - 转义元字符

* `Htmlspecialchars()` - HTML特殊字符转义

* `Htmlspecialchars_decode()` - HTML特殊字符解码

* `Nl2br()` - 换行符转换为HTML换行

* `StrReplace()` - 字符串替换

* `Trim()` - 去除首尾空白

* `Ltrim()` - 去除左侧空白

* `Rtrim()` - 去除右侧空白

* `Strrev()` - 字符串反转

* `NumberFormat()` - 数字格式化

* `StrShuffle()` - 字符串随机打乱

* `StrWordCount()` - 单词计数

* `Strlen()` - 字符串长度

* `MbStrlen()` - 多字节字符串长度

* `StrRepeat()` - 字符串重复

* `StrPad()` - 字符串填充

#### 类型转换

* `ToString()` - 转换为字符串

* `ToInt()` - 转换为整数

* `ToInt64()` - 转换为64位整数

* `ToFloat64()` - 转换为64位浮点数

* `ToBool()` - 转换为布尔值

#### 随机数和编码

* `Rand()` - 生成随机数

* `RandString()` - 生成随机字符串

* `Base64Encode()` - Base64编码

* `Base64Decode()` - Base64解码

* `UrlEncode()` - URL编码

* `UrlDecode()` - URL解码

#### 数组操作

* `ArrayReverse()` - 数组反转

* `ArrayUnique()` - 数组去重

* `ArrayMerge()` - 数组合并

* `ArrayKeys()` - 获取数组键

* `ArrayValues()` - 获取数组值

* `ArrayChunk()` - 数组分块

* `ArraySlice()` - 数组切片

* `ArrayPush()` - 数组追加

* `ArrayPop()` - 数组弹出

* `ArrayShift()` - 数组移除首元素

* `ArrayUnshift()` - 数组前置元素

#### 其他工具

* `Round()` - 数字四舍五入

* `Ceil()` - 向上取整

* `Floor()` - 向下取整

* `Abs()` - 绝对值

* `Max()` - 最大值

* `Min()` - 最小值

* `Md5()` - MD5哈希

* `Sha1()` - SHA1哈希

* `FileExists()` - 文件存在检查

* `IsDir()` - 目录检查

* `IsFile()` - 文件检查

* `Unlink()` - 删除文件

* `Mkdir()` - 创建目录

## 2. 测试文件结构

### 2.1 文件位置

```
tests/unit/pkg/helpers_test.go
```

### 2.2 测试依赖

```go
package pkg

import (
    "math"
    "os"
    "path/filepath"
    "reflect"
    "strings"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

## 3. 字符串处理函数测试

### 3.1 TestSubstr - 字符串截取测试

```go
func TestSubstr(t *testing.T) {
    tests := []struct {
        name     string
        str      string
        start    int
        length   int
        expected string
    }{
        {"正常截取", "hello world", 0, 5, "hello"},
        {"从中间截取", "hello world", 6, 5, "world"},
        {"负数起始位置", "hello world", -5, 5, "world"},
        {"超出长度", "hello", 0, 10, "hello"},
        {"零长度", "hello", 0, 0, ""},
        {"空字符串", "", 0, 5, ""},
        {"中文字符串", "你好世界", 0, 2, "你好"},
        {"起始位置超出", "hello", 10, 5, ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Substr(tt.str, tt.start, tt.length)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3.2 TestInArray - 数组包含检查测试

```go
func TestInArray(t *testing.T) {
    tests := []struct {
        name     string
        needle   interface{}
        haystack []interface{}
        expected bool
    }{
        {"字符串存在", "apple", []interface{}{"apple", "banana", "orange"}, true},
        {"字符串不存在", "grape", []interface{}{"apple", "banana", "orange"}, false},
        {"数字存在", 42, []interface{}{1, 42, 100}, true},
        {"数字不存在", 99, []interface{}{1, 42, 100}, false},
        {"空数组", "test", []interface{}{}, false},
        {"nil值", nil, []interface{}{nil, "test"}, true},
        {"混合类型", "42", []interface{}{42, "42", true}, true},
        {"布尔值", true, []interface{}{false, true, 1}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := InArray(tt.needle, tt.haystack)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3.3 TestImplodeExplode - 数组连接和分割测试

```go
func TestImplode(t *testing.T) {
    tests := []struct {
        name      string
        separator string
        pieces    []string
        expected  string
    }{
        {"逗号分隔", ",", []string{"a", "b", "c"}, "a,b,c"},
        {"空格分隔", " ", []string{"hello", "world"}, "hello world"},
        {"空分隔符", "", []string{"a", "b", "c"}, "abc"},
        {"单个元素", ",", []string{"single"}, "single"},
        {"空数组", ",", []string{}, ""},
        {"包含空字符串", ",", []string{"a", "", "c"}, "a,,c"},
        {"中文字符", "|", []string{"你好", "世界"}, "你好|世界"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Implode(tt.separator, tt.pieces)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestExplode(t *testing.T) {
    tests := []struct {
        name      string
        separator string
        str       string
        expected  []string
    }{
        {"逗号分割", ",", "a,b,c", []string{"a", "b", "c"}},
        {"空格分割", " ", "hello world", []string{"hello", "world"}},
        {"不存在分隔符", ",", "hello", []string{"hello"}},
        {"空字符串", ",", "", []string{""}},
        {"连续分隔符", ",", "a,,c", []string{"a", "", "c"}},
        {"中文分割", "|", "你好|世界", []string{"你好", "世界"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Explode(tt.separator, tt.str)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3.4 TestStringCase - 字符串大小写测试

```go
func TestUcfirst(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"小写首字母", "hello", "Hello"},
        {"已大写首字母", "Hello", "Hello"},
        {"空字符串", "", ""},
        {"单字符", "a", "A"},
        {"数字开头", "1hello", "1hello"},
        {"中文字符", "你好", "你好"},
        {"特殊字符开头", "@hello", "@hello"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Ucfirst(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestLcfirst(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"大写首字母", "Hello", "hello"},
        {"已小写首字母", "hello", "hello"},
        {"空字符串", "", ""},
        {"单字符", "A", "a"},
        {"数字开头", "1Hello", "1Hello"},
        {"中文字符", "你好", "你好"},
        {"特殊字符开头", "@Hello", "@Hello"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Lcfirst(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 4. 类型转换函数测试

### 4.1 TestTypeConversions - 类型转换测试

```go
func TestToString(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected string
    }{
        {"整数", 42, "42"},
        {"浮点数", 3.14, "3.14"},
        {"布尔值true", true, "true"},
        {"布尔值false", false, "false"},
        {"字符串", "hello", "hello"},
        {"nil值", nil, ""},
        {"零值", 0, "0"},
        {"负数", -42, "-42"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ToString(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestToInt(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected int
    }{
        {"字符串数字", "42", 42},
        {"整数", 42, 42},
        {"浮点数", 3.14, 3},
        {"布尔值true", true, 1},
        {"布尔值false", false, 0},
        {"无效字符串", "hello", 0},
        {"空字符串", "", 0},
        {"nil值", nil, 0},
        {"负数字符串", "-42", -42},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ToInt(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestToFloat64(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected float64
    }{
        {"字符串浮点数", "3.14", 3.14},
        {"整数", 42, 42.0},
        {"浮点数", 3.14, 3.14},
        {"布尔值true", true, 1.0},
        {"布尔值false", false, 0.0},
        {"无效字符串", "hello", 0.0},
        {"空字符串", "", 0.0},
        {"nil值", nil, 0.0},
        {"科学计数法", "1.23e2", 123.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ToFloat64(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestToBool(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected bool
    }{
        {"布尔值true", true, true},
        {"布尔值false", false, false},
        {"字符串true", "true", true},
        {"字符串false", "false", false},
        {"数字1", 1, true},
        {"数字0", 0, false},
        {"非零数字", 42, true},
        {"空字符串", "", false},
        {"非空字符串", "hello", true},
        {"nil值", nil, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ToBool(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 5. 编码和随机函数测试

### 5.1 TestEncodingFunctions - 编码函数测试

```go
func TestBase64EncodeDecode(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"普通字符串", "hello world"},
        {"空字符串", ""},
        {"中文字符串", "你好世界"},
        {"特殊字符", "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./ "},
        {"数字", "1234567890"},
        {"长字符串", strings.Repeat("a", 1000)},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            encoded := Base64Encode(tt.input)
            decoded := Base64Decode(encoded)
            assert.Equal(t, tt.input, decoded)
        })
    }
}

func TestUrlEncodeDecode(t *testing.T) {
    tests := []struct {
        name  string
        input string
    }{
        {"普通字符串", "hello world"},
        {"空字符串", ""},
        {"中文字符串", "你好世界"},
        {"特殊字符", "hello world!@#$%^&*()_+-="},
        {"URL路径", "/path/to/resource?param=value&other=test"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            encoded := UrlEncode(tt.input)
            decoded := UrlDecode(encoded)
            assert.Equal(t, tt.input, decoded)
        })
    }
}

func TestHashFunctions(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        md5Hash  string
        sha1Hash string
    }{
        {
            "空字符串",
            "",
            "d41d8cd98f00b204e9800998ecf8427e",
            "da39a3ee5e6b4b0d3255bfef95601890afd80709",
        },
        {
            "hello",
            "hello",
            "5d41402abc4b2a76b9719d911017c592",
            "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            md5Result := Md5(tt.input)
            sha1Result := Sha1(tt.input)
            
            assert.Equal(t, tt.md5Hash, md5Result)
            assert.Equal(t, tt.sha1Hash, sha1Result)
        })
    }
}
```

### 5.2 TestRandomFunctions - 随机函数测试

```go
func TestRand(t *testing.T) {
    tests := []struct {
        name string
        min  int
        max  int
    }{
        {"正数范围", 1, 10},
        {"包含零", 0, 5},
        {"负数范围", -10, -1},
        {"跨零范围", -5, 5},
        {"单点范围", 5, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            for i := 0; i < 100; i++ {
                result := Rand(tt.min, tt.max)
                assert.GreaterOrEqual(t, result, tt.min)
                assert.LessOrEqual(t, result, tt.max)
            }
        })
    }
}

func TestRandString(t *testing.T) {
    tests := []struct {
        name   string
        length int
    }{
        {"长度0", 0},
        {"长度1", 1},
        {"长度10", 10},
        {"长度100", 100},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := RandString(tt.length)
            assert.Equal(t, tt.length, len(result))
            
            // 验证字符串只包含预期字符
            for _, char := range result {
                assert.True(t, 
                    (char >= 'a' && char <= 'z') ||
                    (char >= 'A' && char <= 'Z') ||
                    (char >= '0' && char <= '9'),
                    "字符串应该只包含字母和数字")
            }
        })
    }
    
    // 测试随机性
    t.Run("随机性测试", func(t *testing.T) {
        results := make(map[string]bool)
        for i := 0; i < 100; i++ {
            result := RandString(10)
            results[result] = true
        }
        // 100次生成应该有很高的唯一性
        assert.Greater(t, len(results), 90, "随机字符串应该有足够的唯一性")
    })
}
```

## 6. 数组操作函数测试

### 6.1 TestArrayFunctions - 数组操作测试

```go
func TestArrayReverse(t *testing.T) {
    tests := []struct {
        name     string
        input    []interface{}
        expected []interface{}
    }{
        {"数字数组", []interface{}{1, 2, 3, 4}, []interface{}{4, 3, 2, 1}},
        {"字符串数组", []interface{}{"a", "b", "c"}, []interface{}{"c", "b", "a"}},
        {"单元素数组", []interface{}{"single"}, []interface{}{"single"}},
        {"空数组", []interface{}{}, []interface{}{}},
        {"混合类型", []interface{}{1, "hello", true}, []interface{}{true, "hello", 1}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ArrayReverse(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestArrayUnique(t *testing.T) {
    tests := []struct {
        name     string
        input    []interface{}
        expected []interface{}
    }{
        {"有重复元素", []interface{}{1, 2, 2, 3, 1}, []interface{}{1, 2, 3}},
        {"无重复元素", []interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
        {"全部相同", []interface{}{"a", "a", "a"}, []interface{}{"a"}},
        {"空数组", []interface{}{}, []interface{}{}},
        {"混合类型重复", []interface{}{1, "1", 1, "1"}, []interface{}{1, "1"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ArrayUnique(tt.input)
            assert.ElementsMatch(t, tt.expected, result)
        })
    }
}

func TestArrayMerge(t *testing.T) {
    tests := []struct {
        name     string
        arrays   [][]interface{}
        expected []interface{}
    }{
        {
            "两个数组",
            [][]interface{}{{1, 2}, {3, 4}},
            []interface{}{1, 2, 3, 4},
        },
        {
            "三个数组",
            [][]interface{}{{"a"}, {"b", "c"}, {"d"}},
            []interface{}{"a", "b", "c", "d"},
        },
        {
            "包含空数组",
            [][]interface{}{{1, 2}, {}, {3, 4}},
            []interface{}{1, 2, 3, 4},
        },
        {
            "单个数组",
            [][]interface{}{{1, 2, 3}},
            []interface{}{1, 2, 3},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ArrayMerge(tt.arrays...)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 7. 数学函数测试

### 7.1 TestMathFunctions - 数学函数测试

```go
func TestRound(t *testing.T) {
    tests := []struct {
        name      string
        value     float64
        precision int
        expected  float64
    }{
        {"正数四舍五入", 3.14159, 2, 3.14},
        {"正数五入", 3.145, 2, 3.15},
        {"负数四舍五入", -3.14159, 2, -3.14},
        {"整数", 5.0, 2, 5.0},
        {"零精度", 3.7, 0, 4.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Round(tt.value, tt.precision)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestCeilFloor(t *testing.T) {
    tests := []struct {
        name     string
        value    float64
        ceil     float64
        floor    float64
    }{
        {"正小数", 3.2, 4.0, 3.0},
        {"负小数", -3.2, -3.0, -4.0},
        {"整数", 5.0, 5.0, 5.0},
        {"接近整数", 4.999, 5.0, 4.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ceilResult := Ceil(tt.value)
            floorResult := Floor(tt.value)
            
            assert.Equal(t, tt.ceil, ceilResult)
            assert.Equal(t, tt.floor, floorResult)
        })
    }
}

func TestAbs(t *testing.T) {
    tests := []struct {
        name     string
        value    float64
        expected float64
    }{
        {"正数", 5.5, 5.5},
        {"负数", -5.5, 5.5},
        {"零", 0.0, 0.0},
        {"负零", -0.0, 0.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Abs(tt.value)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestMaxMin(t *testing.T) {
    tests := []struct {
        name     string
        values   []float64
        max      float64
        min      float64
    }{
        {"正数序列", []float64{1.0, 5.0, 3.0, 2.0}, 5.0, 1.0},
        {"包含负数", []float64{-1.0, 5.0, -3.0, 2.0}, 5.0, -3.0},
        {"单个值", []float64{42.0}, 42.0, 42.0},
        {"相同值", []float64{3.0, 3.0, 3.0}, 3.0, 3.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            maxResult := Max(tt.values...)
            minResult := Min(tt.values...)
            
            assert.Equal(t, tt.max, maxResult)
            assert.Equal(t, tt.min, minResult)
        })
    }
}
```

## 8. 文件操作函数测试

### 8.1 TestFileOperations - 文件操作测试

```go
func TestFileOperations(t *testing.T) {
    // 创建临时目录用于测试
    tempDir := t.TempDir()
    
    t.Run("文件存在检查", func(t *testing.T) {
        // 创建测试文件
        testFile := filepath.Join(tempDir, "test.txt")
        err := os.WriteFile(testFile, []byte("test content"), 0644)
        require.NoError(t, err)
        
        // 测试存在的文件
        assert.True(t, FileExists(testFile))
        assert.True(t, IsFile(testFile))
        assert.False(t, IsDir(testFile))
        
        // 测试不存在的文件
        nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
        assert.False(t, FileExists(nonExistentFile))
        assert.False(t, IsFile(nonExistentFile))
    })
    
    t.Run("目录检查", func(t *testing.T) {
        // 测试存在的目录
        assert.True(t, FileExists(tempDir))
        assert.True(t, IsDir(tempDir))
        assert.False(t, IsFile(tempDir))
        
        // 测试不存在的目录
        nonExistentDir := filepath.Join(tempDir, "nonexistent")
        assert.False(t, FileExists(nonExistentDir))
        assert.False(t, IsDir(nonExistentDir))
    })
    
    t.Run("创建目录", func(t *testing.T) {
        newDir := filepath.Join(tempDir, "newdir")
        
        // 目录不存在
        assert.False(t, IsDir(newDir))
        
        // 创建目录
        err := Mkdir(newDir, 0755)
        assert.NoError(t, err)
        
        // 验证目录已创建
        assert.True(t, IsDir(newDir))
        
        // 重复创建应该不报错
        err = Mkdir(newDir, 0755)
        assert.NoError(t, err)
    })
    
    t.Run("删除文件", func(t *testing.T) {
        // 创建测试文件
        testFile := filepath.Join(tempDir, "delete_test.txt")
        err := os.WriteFile(testFile, []byte("test"), 0644)
        require.NoError(t, err)
        
        // 文件存在
        assert.True(t, FileExists(testFile))
        
        // 删除文件
        err = Unlink(testFile)
        assert.NoError(t, err)
        
        // 验证文件已删除
        assert.False(t, FileExists(testFile))
        
        // 删除不存在的文件应该返回错误
        err = Unlink(testFile)
        assert.Error(t, err)
    })
}
```

## 9. HTML和字符串处理测试

### 9.1 TestHTMLFunctions - HTML处理测试

```go
func TestHtmlspecialchars(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"基本HTML标签", "<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
        {"引号转义", `"hello" & 'world'", "&quot;hello&quot; &amp; &#39;world&#39;"},
        {"空字符串", "", ""},
        {"普通文本", "hello world", "hello world"},
        {"中文字符", "你好<世界>", "你好&lt;世界&gt;"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Htmlspecialchars(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestHtmlspecialchars_decode(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"HTML实体解码", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;", "<script>alert('xss')</script>"},
        {"引号解码", "&quot;hello&quot; &amp; &#39;world&#39;", `"hello" & 'world'`},
        {"空字符串", "", ""},
        {"普通文本", "hello world", "hello world"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Htmlspecialchars_decode(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestNl2br(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"单行换行", "hello\nworld", "hello<br>world"},
        {"多行换行", "line1\nline2\nline3", "line1<br>line2<br>line3"},
        {"Windows换行", "hello\r\nworld", "hello<br>world"},
        {"无换行", "hello world", "hello world"},
        {"空字符串", "", ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Nl2br(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## 10. 性能测试

### 10.1 基准测试

```go
func BenchmarkSubstr(b *testing.B) {
    str := "hello world this is a test string"
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Substr(str, 0, 10)
    }
}

func BenchmarkInArray(b *testing.B) {
    haystack := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
    needle := "e"
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        InArray(needle, haystack)
    }
}

func BenchmarkToString(b *testing.B) {
    value := 12345
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ToString(value)
    }
}

func BenchmarkRandString(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        RandString(10)
    }
}

func BenchmarkBase64Encode(b *testing.B) {
    str := "hello world this is a test string for base64 encoding"
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Base64Encode(str)
    }
}
```

## 11. 测试覆盖率目标

* **函数覆盖率**: 100%

* **行覆盖率**: ≥ 95%

* **分支覆盖率**: ≥ 90%

## 12. 实施步骤

1. **创建测试文件**: `tests/unit/pkg/helpers_test.go`
2. **实现字符串处理测试**: Substr, InArray, Implode, Explode等
3. **实现类型转换测试**: ToString, ToInt, ToFloat64, ToBool
4. **实现编码函数测试**: Base64, URL编码, 哈希函数
5. **实现随机函数测试**: Rand, RandString
6. **实现数组操作测试**: ArrayReverse, ArrayUnique, ArrayMerge等
7. **实现数学函数测试**: Round, Ceil, Floor, Abs, Max, Min
8. **实现文件操作测试**: FileExists, IsDir, IsFile, Mkdir, Unlink
9. **实现HTML处理测试**: Htmlspecialchars, Nl2br等
10. **添加性能测试**: 关键函数的基准测试
11. **验证覆盖率**: 使用 `go test -cover`
12. **优化和完善**: 根据覆盖率报告补充测试

## 13. 注意事项

* 由于helpers.go包含大量函数，建议分批实现测试

* 特别注意边界条件和异常输入的测试

* 文件操作测试需要使用临时目录，确保测试的独立性

* 随机函数测试需要验证随机性和范围正确性

* 类型转换测试需要覆盖各种输入类型

* 确保测试的可重复性和稳定性

* 性能测试应该关注关键路径上的函数

