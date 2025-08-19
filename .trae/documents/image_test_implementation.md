# Image 模块测试实现方案

## 1. 模块概述

`image.go` 模块提供了图片上传和处理功能，包括图片上传配置、Base64图片保存、上传路径生成等核心功能。该模块是Web应用中图片处理的重要组件。

### 1.1 核心功能

#### 结构体和类型

* `UploadConfig` - 上传配置结构体

* `UploadError` - 自定义错误类型

#### 核心函数

* `UploadImage()` - 图片上传处理

* `SaveBase64File()` - Base64图片保存

* `CreateUploadPathFileName()` - 生成上传路径和文件名

#### 错误定义

* `ErrInvalidFileType` - 无效文件类型错误

* `ErrFileSizeExceeded` - 文件大小超限错误

* `ErrUploadFailed` - 上传失败错误

### 1.2 UploadConfig 配置项

```go
type UploadConfig struct {
    MaxSize      int64    // 最大文件大小
    AllowedTypes []string // 允许的文件类型
    UploadPath   string   // 上传路径
    URLPrefix    string   // URL前缀
}
```

## 2. 测试文件结构

### 2.1 文件位置

```
tests/unit/pkg/image_test.go
```

### 2.2 测试依赖

```go
package pkg

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

## 3. 测试辅助函数

### 3.1 图片生成辅助函数

```go
// createTestImage 创建测试用的图片数据
func createTestImage(format string, width, height int) ([]byte, error) {
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    
    // 填充一些颜色
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 128, 255})
        }
    }
    
    var buf bytes.Buffer
    switch format {
    case "jpeg", "jpg":
        err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
        return buf.Bytes(), err
    case "png":
        err := png.Encode(&buf, img)
        return buf.Bytes(), err
    default:
        return nil, fmt.Errorf("unsupported format: %s", format)
    }
}

// createTestBase64Image 创建Base64编码的测试图片
func createTestBase64Image(format string, width, height int) (string, error) {
    imageData, err := createTestImage(format, width, height)
    if err != nil {
        return "", err
    }
    
    mimeType := "image/" + format
    if format == "jpg" {
        mimeType = "image/jpeg"
    }
    
    base64Data := base64.StdEncoding.EncodeToString(imageData)
    return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), nil
}

// createMultipartFile 创建multipart文件上传数据
func createMultipartFile(fieldName, fileName string, data []byte) (*bytes.Buffer, string, error) {
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)
    
    part, err := writer.CreateFormFile(fieldName, fileName)
    if err != nil {
        return nil, "", err
    }
    
    _, err = part.Write(data)
    if err != nil {
        return nil, "", err
    }
    
    err = writer.Close()
    if err != nil {
        return nil, "", err
    }
    
    return &buf, writer.FormDataContentType(), nil
}
```

## 4. UploadConfig 测试

### 4.1 TestUploadConfig - 配置结构体测试

```go
func TestUploadConfig(t *testing.T) {
    t.Run("默认配置创建", func(t *testing.T) {
        config := &UploadConfig{
            MaxSize:      10 * 1024 * 1024, // 10MB
            AllowedTypes: []string{"jpg", "jpeg", "png", "gif"},
            UploadPath:   "/uploads",
            URLPrefix:    "https://example.com",
        }
        
        assert.Equal(t, int64(10*1024*1024), config.MaxSize)
        assert.Equal(t, []string{"jpg", "jpeg", "png", "gif"}, config.AllowedTypes)
        assert.Equal(t, "/uploads", config.UploadPath)
        assert.Equal(t, "https://example.com", config.URLPrefix)
    })
    
    t.Run("配置验证", func(t *testing.T) {
        tests := []struct {
            name   string
            config *UploadConfig
            valid  bool
        }{
            {
                "有效配置",
                &UploadConfig{
                    MaxSize:      1024 * 1024,
                    AllowedTypes: []string{"jpg", "png"},
                    UploadPath:   "/uploads",
                    URLPrefix:    "https://example.com",
                },
                true,
            },
            {
                "零大小限制",
                &UploadConfig{
                    MaxSize:      0,
                    AllowedTypes: []string{"jpg"},
                    UploadPath:   "/uploads",
                    URLPrefix:    "https://example.com",
                },
                false,
            },
            {
                "空允许类型",
                &UploadConfig{
                    MaxSize:      1024 * 1024,
                    AllowedTypes: []string{},
                    UploadPath:   "/uploads",
                    URLPrefix:    "https://example.com",
                },
                false,
            },
            {
                "空上传路径",
                &UploadConfig{
                    MaxSize:      1024 * 1024,
                    AllowedTypes: []string{"jpg"},
                    UploadPath:   "",
                    URLPrefix:    "https://example.com",
                },
                false,
            },
        }
        
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                // 这里可以添加配置验证逻辑
                valid := tt.config.MaxSize > 0 && 
                        len(tt.config.AllowedTypes) > 0 && 
                        tt.config.UploadPath != ""
                assert.Equal(t, tt.valid, valid)
            })
        }
    })
}
```

## 5. UploadImage 函数测试

### 5.1 TestUploadImage - 图片上传测试

```go
func TestUploadImage(t *testing.T) {
    // 创建临时目录用于测试
    tempDir := t.TempDir()
    
    config := &UploadConfig{
        MaxSize:      2 * 1024 * 1024, // 2MB
        AllowedTypes: []string{"jpg", "jpeg", "png"},
        UploadPath:   tempDir,
        URLPrefix:    "https://example.com",
    }
    
    t.Run("成功上传JPEG图片", func(t *testing.T) {
        // 创建测试图片
        imageData, err := createTestImage("jpeg", 100, 100)
        require.NoError(t, err)
        
        // 创建multipart请求
        body, contentType, err := createMultipartFile("image", "test.jpg", imageData)
        require.NoError(t, err)
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        // 执行上传
        result, err := UploadImage(req, "image", config)
        
        assert.NoError(t, err)
        assert.NotEmpty(t, result)
        
        // 验证文件是否已保存
        if result != "" {
            // 从结果URL中提取文件路径
            // 这里需要根据实际的UploadImage实现来验证
        }
    })
    
    t.Run("成功上传PNG图片", func(t *testing.T) {
        imageData, err := createTestImage("png", 50, 50)
        require.NoError(t, err)
        
        body, contentType, err := createMultipartFile("image", "test.png", imageData)
        require.NoError(t, err)
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        result, err := UploadImage(req, "image", config)
        
        assert.NoError(t, err)
        assert.NotEmpty(t, result)
    })
    
    t.Run("文件类型不允许", func(t *testing.T) {
        // 创建一个文本文件伪装成图片
        textData := []byte("This is not an image")
        
        body, contentType, err := createMultipartFile("image", "test.txt", textData)
        require.NoError(t, err)
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        result, err := UploadImage(req, "image", config)
        
        assert.Error(t, err)
        assert.Empty(t, result)
        assert.Contains(t, err.Error(), "invalid file type")
    })
    
    t.Run("文件大小超限", func(t *testing.T) {
        // 创建一个大图片
        imageData, err := createTestImage("jpeg", 2000, 2000)
        require.NoError(t, err)
        
        // 确保图片大小超过限制
        if len(imageData) <= int(config.MaxSize) {
            // 如果生成的图片还是太小，手动创建大数据
            largeData := make([]byte, config.MaxSize+1000)
            copy(largeData, imageData)
            imageData = largeData
        }
        
        body, contentType, err := createMultipartFile("image", "large.jpg", imageData)
        require.NoError(t, err)
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        result, err := UploadImage(req, "image", config)
        
        assert.Error(t, err)
        assert.Empty(t, result)
        assert.Contains(t, err.Error(), "file size exceeded")
    })
    
    t.Run("缺少文件字段", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/upload", strings.NewReader(""))
        req.Header.Set("Content-Type", "multipart/form-data")
        
        result, err := UploadImage(req, "nonexistent", config)
        
        assert.Error(t, err)
        assert.Empty(t, result)
    })
    
    t.Run("无效的multipart请求", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/upload", strings.NewReader("invalid data"))
        req.Header.Set("Content-Type", "application/json")
        
        result, err := UploadImage(req, "image", config)
        
        assert.Error(t, err)
        assert.Empty(t, result)
    })
}
```

## 6. SaveBase64File 函数测试

### 6.1 TestSaveBase64File - Base64文件保存测试

```go
func TestSaveBase64File(t *testing.T) {
    tempDir := t.TempDir()
    
    t.Run("保存JPEG Base64图片", func(t *testing.T) {
        base64Data, err := createTestBase64Image("jpeg", 100, 100)
        require.NoError(t, err)
        
        fileName, err := SaveBase64File(base64Data, tempDir)
        
        assert.NoError(t, err)
        assert.NotEmpty(t, fileName)
        assert.True(t, strings.HasSuffix(fileName, ".jpg"))
        
        // 验证文件是否存在
        fullPath := filepath.Join(tempDir, fileName)
        assert.FileExists(t, fullPath)
        
        // 验证文件内容
        savedData, err := os.ReadFile(fullPath)
        assert.NoError(t, err)
        assert.NotEmpty(t, savedData)
    })
    
    t.Run("保存PNG Base64图片", func(t *testing.T) {
        base64Data, err := createTestBase64Image("png", 50, 50)
        require.NoError(t, err)
        
        fileName, err := SaveBase64File(base64Data, tempDir)
        
        assert.NoError(t, err)
        assert.NotEmpty(t, fileName)
        assert.True(t, strings.HasSuffix(fileName, ".png"))
        
        fullPath := filepath.Join(tempDir, fileName)
        assert.FileExists(t, fullPath)
    })
    
    t.Run("无效的Base64数据", func(t *testing.T) {
        invalidBase64 := "data:image/jpeg;base64,invalid_base64_data!!!"
        
        fileName, err := SaveBase64File(invalidBase64, tempDir)
        
        assert.Error(t, err)
        assert.Empty(t, fileName)
    })
    
    t.Run("不支持的图片格式", func(t *testing.T) {
        // 创建一个不支持格式的Base64数据
        textData := base64.StdEncoding.EncodeToString([]byte("not an image"))
        base64Data := fmt.Sprintf("data:text/plain;base64,%s", textData)
        
        fileName, err := SaveBase64File(base64Data, tempDir)
        
        assert.Error(t, err)
        assert.Empty(t, fileName)
        assert.Contains(t, err.Error(), "unsupported format")
    })
    
    t.Run("缺少data URL前缀", func(t *testing.T) {
        // 只有Base64数据，没有data URL前缀
        base64Data := base64.StdEncoding.EncodeToString([]byte("test data"))
        
        fileName, err := SaveBase64File(base64Data, tempDir)
        
        assert.Error(t, err)
        assert.Empty(t, fileName)
    })
    
    t.Run("空Base64数据", func(t *testing.T) {
        fileName, err := SaveBase64File("", tempDir)
        
        assert.Error(t, err)
        assert.Empty(t, fileName)
    })
    
    t.Run("无效的保存路径", func(t *testing.T) {
        base64Data, err := createTestBase64Image("jpeg", 50, 50)
        require.NoError(t, err)
        
        // 使用不存在的路径
        invalidPath := "/nonexistent/path"
        
        fileName, err := SaveBase64File(base64Data, invalidPath)
        
        assert.Error(t, err)
        assert.Empty(t, fileName)
    })
}
```

## 7. CreateUploadPathFileName 函数测试

### 7.1 TestCreateUploadPathFileName - 路径文件名生成测试

```go
func TestCreateUploadPathFileName(t *testing.T) {
    t.Run("生成基本路径和文件名", func(t *testing.T) {
        basePath := "/uploads"
        extension := "jpg"
        
        path, fileName := CreateUploadPathFileName(basePath, extension)
        
        assert.NotEmpty(t, path)
        assert.NotEmpty(t, fileName)
        assert.True(t, strings.HasPrefix(path, basePath))
        assert.True(t, strings.HasSuffix(fileName, "."+extension))
    })
    
    t.Run("不同扩展名", func(t *testing.T) {
        basePath := "/uploads"
        extensions := []string{"jpg", "jpeg", "png", "gif", "webp"}
        
        for _, ext := range extensions {
            t.Run("扩展名_"+ext, func(t *testing.T) {
                path, fileName := CreateUploadPathFileName(basePath, ext)
                
                assert.True(t, strings.HasSuffix(fileName, "."+ext))
                assert.Contains(t, path, basePath)
            })
        }
    })
    
    t.Run("路径包含日期结构", func(t *testing.T) {
        basePath := "/uploads"
        extension := "jpg"
        
        path, _ := CreateUploadPathFileName(basePath, extension)
        
        // 验证路径包含年月日结构
        now := time.Now()
        expectedYear := fmt.Sprintf("%d", now.Year())
        expectedMonth := fmt.Sprintf("%02d", now.Month())
        expectedDay := fmt.Sprintf("%02d", now.Day())
        
        assert.Contains(t, path, expectedYear)
        assert.Contains(t, path, expectedMonth)
        assert.Contains(t, path, expectedDay)
    })
    
    t.Run("文件名唯一性", func(t *testing.T) {
        basePath := "/uploads"
        extension := "jpg"
        
        // 生成多个文件名，验证唯一性
        fileNames := make(map[string]bool)
        for i := 0; i < 100; i++ {
            _, fileName := CreateUploadPathFileName(basePath, extension)
            
            assert.False(t, fileNames[fileName], "文件名应该是唯一的: %s", fileName)
            fileNames[fileName] = true
        }
    })
    
    t.Run("空基础路径", func(t *testing.T) {
        basePath := ""
        extension := "jpg"
        
        path, fileName := CreateUploadPathFileName(basePath, extension)
        
        assert.NotEmpty(t, fileName)
        // 路径可能为空或包含默认值
        assert.True(t, strings.HasSuffix(fileName, "."+extension))
    })
    
    t.Run("空扩展名", func(t *testing.T) {
        basePath := "/uploads"
        extension := ""
        
        path, fileName := CreateUploadPathFileName(basePath, extension)
        
        assert.NotEmpty(t, path)
        assert.NotEmpty(t, fileName)
        // 文件名可能没有扩展名或有默认扩展名
    })
}
```

## 8. 错误类型测试

### 8.1 TestUploadError - 自定义错误测试

```go
func TestUploadError(t *testing.T) {
    t.Run("错误类型定义", func(t *testing.T) {
        // 测试预定义的错误
        assert.NotNil(t, ErrInvalidFileType)
        assert.NotNil(t, ErrFileSizeExceeded)
        assert.NotNil(t, ErrUploadFailed)
        
        // 验证错误消息
        assert.Contains(t, ErrInvalidFileType.Error(), "invalid file type")
        assert.Contains(t, ErrFileSizeExceeded.Error(), "file size exceeded")
        assert.Contains(t, ErrUploadFailed.Error(), "upload failed")
    })
    
    t.Run("错误类型判断", func(t *testing.T) {
        // 测试错误类型判断
        var err error
        
        err = ErrInvalidFileType
        assert.True(t, errors.Is(err, ErrInvalidFileType))
        assert.False(t, errors.Is(err, ErrFileSizeExceeded))
        
        err = ErrFileSizeExceeded
        assert.True(t, errors.Is(err, ErrFileSizeExceeded))
        assert.False(t, errors.Is(err, ErrInvalidFileType))
        
        err = ErrUploadFailed
        assert.True(t, errors.Is(err, ErrUploadFailed))
        assert.False(t, errors.Is(err, ErrInvalidFileType))
    })
    
    t.Run("自定义UploadError", func(t *testing.T) {
        customErr := UploadError("custom upload error")
        
        assert.Equal(t, "custom upload error", customErr.Error())
        assert.Equal(t, "custom upload error", string(customErr))
    })
}
```

## 9. 集成测试

### 9.1 TestImageUploadWorkflow - 完整上传流程测试

```go
func TestImageUploadWorkflow(t *testing.T) {
    tempDir := t.TempDir()
    
    config := &UploadConfig{
        MaxSize:      5 * 1024 * 1024, // 5MB
        AllowedTypes: []string{"jpg", "jpeg", "png", "gif"},
        UploadPath:   tempDir,
        URLPrefix:    "https://example.com/uploads",
    }
    
    t.Run("完整上传流程", func(t *testing.T) {
        // 1. 创建测试图片
        imageData, err := createTestImage("jpeg", 200, 200)
        require.NoError(t, err)
        
        // 2. 通过multipart上传
        body, contentType, err := createMultipartFile("image", "test.jpg", imageData)
        require.NoError(t, err)
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        uploadResult, err := UploadImage(req, "image", config)
        assert.NoError(t, err)
        assert.NotEmpty(t, uploadResult)
        
        // 3. 验证文件确实被保存
        // 这里需要根据UploadImage的实际返回值来验证文件存在
        
        // 4. 测试Base64上传相同图片
        base64Data, err := createTestBase64Image("jpeg", 200, 200)
        require.NoError(t, err)
        
        base64FileName, err := SaveBase64File(base64Data, tempDir)
        assert.NoError(t, err)
        assert.NotEmpty(t, base64FileName)
        
        // 5. 验证两种方式都能成功保存文件
        base64FilePath := filepath.Join(tempDir, base64FileName)
        assert.FileExists(t, base64FilePath)
    })
    
    t.Run("多文件上传", func(t *testing.T) {
        formats := []string{"jpeg", "png"}
        results := make([]string, 0, len(formats))
        
        for _, format := range formats {
            imageData, err := createTestImage(format, 100, 100)
            require.NoError(t, err)
            
            fileName := fmt.Sprintf("test.%s", format)
            if format == "jpeg" {
                fileName = "test.jpg"
            }
            
            body, contentType, err := createMultipartFile("image", fileName, imageData)
            require.NoError(t, err)
            
            req := httptest.NewRequest("POST", "/upload", body)
            req.Header.Set("Content-Type", contentType)
            
            result, err := UploadImage(req, "image", config)
            assert.NoError(t, err)
            assert.NotEmpty(t, result)
            
            results = append(results, result)
        }
        
        // 验证所有上传结果都不同
        assert.Len(t, results, len(formats))
        for i := 0; i < len(results); i++ {
            for j := i + 1; j < len(results); j++ {
                assert.NotEqual(t, results[i], results[j], "上传结果应该是唯一的")
            }
        }
    })
}
```

## 10. 性能测试

### 10.1 基准测试

```go
func BenchmarkUploadImage(b *testing.B) {
    tempDir := b.TempDir()
    config := &UploadConfig{
        MaxSize:      10 * 1024 * 1024,
        AllowedTypes: []string{"jpg", "jpeg", "png"},
        UploadPath:   tempDir,
        URLPrefix:    "https://example.com",
    }
    
    // 预创建测试数据
    imageData, err := createTestImage("jpeg", 100, 100)
    if err != nil {
        b.Fatal(err)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        body, contentType, err := createMultipartFile("image", "test.jpg", imageData)
        if err != nil {
            b.Fatal(err)
        }
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        _, err = UploadImage(req, "image", config)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkSaveBase64File(b *testing.B) {
    tempDir := b.TempDir()
    
    // 预创建Base64数据
    base64Data, err := createTestBase64Image("jpeg", 100, 100)
    if err != nil {
        b.Fatal(err)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := SaveBase64File(base64Data, tempDir)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkCreateUploadPathFileName(b *testing.B) {
    basePath := "/uploads"
    extension := "jpg"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        CreateUploadPathFileName(basePath, extension)
    }
}
```

## 11. 边界和异常测试

### 11.1 TestImageEdgeCases - 边界情况测试

```go
func TestImageEdgeCases(t *testing.T) {
    tempDir := t.TempDir()
    
    t.Run("极小图片", func(t *testing.T) {
        // 1x1像素图片
        imageData, err := createTestImage("png", 1, 1)
        require.NoError(t, err)
        
        config := &UploadConfig{
            MaxSize:      1024 * 1024,
            AllowedTypes: []string{"png"},
            UploadPath:   tempDir,
            URLPrefix:    "https://example.com",
        }
        
        body, contentType, err := createMultipartFile("image", "tiny.png", imageData)
        require.NoError(t, err)
        
        req := httptest.NewRequest("POST", "/upload", body)
        req.Header.Set("Content-Type", contentType)
        
        result, err := UploadImage(req, "image", config)
        assert.NoError(t, err)
        assert.NotEmpty(t, result)
    })
    
    t.Run("文件名包含特殊字符", func(t *testing.T) {
        imageData, err := createTestImage("jpeg", 50, 50)
        require.NoError(t, err)
        
        config := &UploadConfig{
            MaxSize:      1024 * 1024,
            AllowedTypes: []string{"jpg", "jpeg"},
            UploadPath:   tempDir,
            URLPrefix:    "https://example.com",
        }
        
        // 包含特殊字符的文件名
        specialNames := []string{
            "test file.jpg",
            "测试图片.jpg",
            "test@#$%.jpg",
            "test(1).jpg",
        }
        
        for _, fileName := range specialNames {
            t.Run("文件名_"+fileName, func(t *testing.T) {
                body, contentType, err := createMultipartFile("image", fileName, imageData)
                require.NoError(t, err)
                
                req := httptest.NewRequest("POST", "/upload", body)
                req.Header.Set("Content-Type", contentType)
                
                result, err := UploadImage(req, "image", config)
                // 根据实现，可能成功或失败
                if err != nil {
                    t.Logf("特殊文件名 %s 上传失败: %v", fileName, err)
                } else {
                    assert.NotEmpty(t, result)
                }
            })
        }
    })
    
    t.Run("并发上传", func(t *testing.T) {
        config := &UploadConfig{
            MaxSize:      1024 * 1024,
            AllowedTypes: []string{"jpg", "jpeg"},
            UploadPath:   tempDir,
            URLPrefix:    "https://example.com",
        }
        
        const numGoroutines = 10
        results := make(chan string, numGoroutines)
        errors := make(chan error, numGoroutines)
        
        for i := 0; i < numGoroutines; i++ {
            go func(id int) {
                imageData, err := createTestImage("jpeg", 50, 50)
                if err != nil {
                    errors <- err
                    return
                }
                
                fileName := fmt.Sprintf("concurrent_%d.jpg", id)
                body, contentType, err := createMultipartFile("image", fileName, imageData)
                if err != nil {
                    errors <- err
                    return
                }
                
                req := httptest.NewRequest("POST", "/upload", body)
                req.Header.Set("Content-Type", contentType)
                
                result, err := UploadImage(req, "image", config)
                if err != nil {
                    errors <- err
                    return
                }
                
                results <- result
            }(i)
        }
        
        // 收集结果
        successCount := 0
        errorCount := 0
        
        for i := 0; i < numGoroutines; i++ {
            select {
            case result := <-results:
                assert.NotEmpty(t, result)
                successCount++
            case err := <-errors:
                t.Logf("并发上传错误: %v", err)
                errorCount++
            }
        }
        
        t.Logf("并发上传结果: 成功 %d, 失败 %d", successCount, errorCount)
        assert.Greater(t, successCount, 0, "至少应该有一些成功的上传")
    })
}
```

## 12. 测试覆盖率目标

* **函数覆盖率**: 100%

* **行覆盖率**: ≥ 95%

* **分支覆盖率**: ≥ 90%

## 13. 实施步骤

1. **创建测试文件**: `tests/unit/pkg/image_test.go`
2. **实现辅助函数**: 图片生成、multipart创建等
3. **实现配置测试**: UploadConfig结构体测试
4. **实现上传测试**: UploadImage函数测试
5. **实现Base64测试**: SaveBase64File函数测试
6. **实现路径生成测试**: CreateUploadPathFileName函数测试
7. **实现错误类型测试**: 自定义错误测试
8. **实现集成测试**: 完整上传流程测试
9. **实现性能测试**: 基准测试
10. **实现边界测试**: 特殊情况和并发测试
11. **验证覆盖率**: 使用 `go test -cover`
12. **优化和完善**: 根据覆盖率报告补充测试

## 14. 注意事项

* 图片测试需要生成真实的图片数据，确保测试的真实性

* 使用临时目录进行文件操作测试，确保测试的独立性

* 注意测试不同图片格式的兼容性

* 文件上传测试需要模拟真实的HTTP multipart请求

* Base64测试需要验证编码解码的正确性

* 并发测试需要验证线程安全性

* 性能测试应该关注大文件和高并发场景

* 确保测试后清理临时文件，避免磁盘空间浪费

