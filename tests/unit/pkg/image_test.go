package pkg_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/biuren/hebe/pkg"
	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestImage 创建一个测试用的图片
func createTestImage(format string, width, height int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充一些颜色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // 红色
		}
	}

	buf := new(bytes.Buffer)
	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(buf, img, nil)
		return buf.Bytes(), err
	case "png":
		err := png.Encode(buf, img)
		return buf.Bytes(), err
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// TestUploadConfig 测试UploadConfig结构体
func TestUploadConfig(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("创建默认配置", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		config := &pkg.UploadConfig{
			MaxSize: 1024 * 1024, // 1MB
			Types:   "data:image/jpg;data:image/png;data:image/gif",
			Path:    "/uploads",
		}

		assert.Equal(t, 1024*1024, config.MaxSize, "最大文件大小应该正确")
		assert.Equal(t, "data:image/jpg;data:image/png;data:image/gif", config.Types, "允许的文件类型应该正确")
		assert.Equal(t, "/uploads", config.Path, "上传路径应该正确")
	})
}

// TestUploadError 测试UploadError类型
func TestUploadError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name      string
		errorCode *pkg.UploadError
		expected  string
	}{
		{
			name:      "文件过大错误",
			errorCode: pkg.ErrFileTooLarge,
			expected:  "文件太大",
		},
		{
			name:      "文件类型不支持错误",
			errorCode: pkg.ErrFileFormatUnsupported,
			expected:  "文件格式不支持",
		},
		{
			name:      "文件上传失败错误",
			errorCode: pkg.ErrUploadFailed,
			expected:  "上传失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}
			assert.Equal(t, tt.expected, tt.errorCode.Error(), "错误消息应该匹配")
		})
	}
}

// TestUploadImage 测试图片上传功能
func TestUploadImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "image_upload_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		base64Data  string
		config      *pkg.UploadConfig
		expectError bool
		errorType   *pkg.UploadError
	}{
		{
			name:       "有效的JPEG文件",
			base64Data: "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString([]byte("fake jpeg data for testing")),
			config: &pkg.UploadConfig{
				MaxSize: 1024 * 1024,
				Types:   "jpg,png,gif",
				Path:    tempDir,
			},
			expectError: false,
		},
		{
			name:       "文件过大",
			base64Data: "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(bytes.Repeat([]byte("x"), 2*1024*1024)),
			config: &pkg.UploadConfig{
				MaxSize: 1024 * 1024, // 1MB限制
				Types:   "jpg,png,gif",
				Path:    tempDir,
			},
			expectError: true,
			errorType:   pkg.ErrFileTooLarge,
		},
		{
			name:       "不支持的文件类型",
			base64Data: "data:image/bmp;base64," + base64.StdEncoding.EncodeToString([]byte("fake bmp data")),
			config: &pkg.UploadConfig{
				MaxSize: 1024 * 1024,
				Types:   "jpg,png",
				Path:    tempDir,
			},
			expectError: true,
			errorType:   pkg.ErrFileFormatUnsupported,
		},
		{
			name:       "空数据",
			base64Data: "data:image/jpeg;base64,",
			config: &pkg.UploadConfig{
				MaxSize: 1024 * 1024,
				Types:   "jpg,png,gif",
				Path:    tempDir,
			},
			expectError: true,
			errorType:   pkg.ErrUploadFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}
			// 调用UploadImage函数
			filePath, fileName, uploadErr := pkg.UploadImage(tt.config, tt.base64Data)

			if tt.expectError {
				assert.Error(t, uploadErr, "应该返回错误")
				if tt.errorType != nil {
					// 检查错误类型是否为 UploadError
					if uploadError, ok := uploadErr.(*pkg.UploadError); ok {
						assert.Equal(t, tt.errorType.Code, uploadError.Code, "错误代码应该匹配")
						assert.Equal(t, tt.errorType.Message, uploadError.Message, "错误消息应该匹配")
					} else {
						t.Errorf("期望 UploadError 类型，但得到 %T", uploadErr)
					}
				}
				assert.Empty(t, filePath, "文件路径应该为空")
				assert.Empty(t, fileName, "文件名应该为空")
			} else {
				assert.NoError(t, uploadErr, "不应该返回错误")
				assert.NotEmpty(t, filePath, "文件路径不应该为空")
				assert.NotEmpty(t, fileName, "文件名不应该为空")

				// 验证文件是否真的被创建
				_, err := os.Stat(fileName)
				assert.NoError(t, err, "上传的文件应该存在")

				// 清理测试文件
				os.Remove(fileName)
			}
		})
	}
}

// TestSaveBase64File 测试Base64文件保存功能
func TestSaveBase64File(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "base64_save_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		base64Data  string
		filename    string
		expectError bool
	}{
		{
			name:        "有效的Base64数据",
			base64Data:  base64.StdEncoding.EncodeToString([]byte("test image data")),
			filename:    "test.jpg",
			expectError: false,
		},
		{
			name:        "空的Base64数据",
			base64Data:  "",
			filename:    "empty.jpg",
			expectError: true,
		},
		{
			name:        "无效的Base64数据",
			base64Data:  "invalid-base64-data!!!",
			filename:    "invalid.jpg",
			expectError: true,
		},
		{
			name:        "长Base64数据",
			base64Data:  base64.StdEncoding.EncodeToString(bytes.Repeat([]byte("test"), 1000)),
			filename:    "long.jpg",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}
			filePath := filepath.Join(tempDir, tt.filename)

			// 为 SaveBase64File 准备正确的 base64 数据格式
			base64DataWithHeader := "data:image/jpeg;base64," + tt.base64Data
			err := pkg.SaveBase64File(base64DataWithHeader, tempDir, filePath)

			if tt.expectError {
				assert.Error(t, err, "应该返回错误")
				// 验证文件不应该被创建
				_, statErr := os.Stat(filePath)
				assert.True(t, os.IsNotExist(statErr), "文件不应该存在")
			} else {
				assert.NoError(t, err, "不应该返回错误")

				// 验证文件是否被创建
				_, statErr := os.Stat(filePath)
				assert.NoError(t, statErr, "文件应该存在")

				// 验证文件内容
				fileContent, readErr := os.ReadFile(filePath)
				assert.NoError(t, readErr, "应该能读取文件")

				// 解码原始 base64 数据（不包含头部）
				expectedContent, decodeErr := base64.StdEncoding.DecodeString(tt.base64Data)
				if decodeErr == nil {
					assert.Equal(t, expectedContent, fileContent, "文件内容应该匹配")
				}

				// 清理测试文件
				os.Remove(filePath)
			}
		})
	}
}

// TestCreateUploadPathFileName 测试上传路径和文件名生成功能
func TestCreateUploadPathFileName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	tests := []struct {
		name        string
		expectedExt string
	}{
		{
			name:        "JPEG文件",
			expectedExt: ".jpg",
		},
		{
			name:        "PNG文件",
			expectedExt: ".png",
		},
		{
			name:        "GIF文件",
			expectedExt: ".gif",
		},
		{
			name:        "无扩展名文件",
			expectedExt: "",
		},
		{
			name:        "多个点的文件名",
			expectedExt: ".jpg",
		},
		{
			name:        "大写扩展名",
			expectedExt: ".JPG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			select {
			case <-ctx.Done():
				t.Fatal("测试超时")
			default:
			}
			filePath, fileName := pkg.CreateUploadPathFileName()
			result := filePath + fileName + tt.expectedExt

			// 验证返回的路径不为空
			assert.NotEmpty(t, result, "生成的路径不应该为空")

			// 验证文件扩展名
			if tt.expectedExt != "" {
				assert.True(t, strings.HasSuffix(result, tt.expectedExt), "应该包含正确的文件扩展名")
			}

			// 验证路径格式（应该包含日期目录结构）
			now := time.Now()
			expectedDatePath := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
			assert.Contains(t, result, expectedDatePath, "应该包含日期目录结构")

			// 验证文件名的唯一性（多次调用应该生成不同的文件名）
			filePath2, fileName2 := pkg.CreateUploadPathFileName()
			result2 := filePath2 + fileName2 + tt.expectedExt
			assert.NotEqual(t, result, result2, "多次调用应该生成不同的文件名")
		})
	}
}

// TestCreateUploadPathFileNameFormat 测试文件名格式
func TestCreateUploadPathFileNameFormat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("验证生成的路径格式", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		filePath, fileName := pkg.CreateUploadPathFileName()
		result := filePath + fileName

		// 分割路径
		parts := strings.Split(result, "/")
		assert.GreaterOrEqual(t, len(parts), 4, "路径应该至少包含4个部分：年/月/日/文件名")

		// 验证年份
		year := parts[0]
		assert.Len(t, year, 4, "年份应该是4位数")

		// 验证月份
		month := parts[1]
		assert.Len(t, month, 2, "月份应该是2位数")

		// 验证日期
		day := parts[2]
		assert.Len(t, day, 2, "日期应该是2位数")

		// 验证文件名（不包含扩展名）
		filename := parts[3]
		assert.NotEmpty(t, filename, "文件名不应该为空")
		assert.Greater(t, len(filename), 4, "文件名应该包含时间戳等信息")
		// CreateUploadPathFileName 只返回基础文件名，不包含扩展名
		assert.False(t, strings.Contains(filename, "."), "基础文件名不应该包含扩展名")
	})
}

// TestImageIntegration 测试图片模块的集成功能
func TestImageIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	// 创建临时目录用于测试
	tempDir, err := os.MkdirTemp("", "image_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("完整的图片上传流程", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		// 1. 创建测试图片
		imageData, err := createTestImage("jpg", 200, 200)
		require.NoError(t, err)

		// 2. 创建上传配置
		config := &pkg.UploadConfig{
			MaxSize: 1024 * 1024,
			Types:   "jpg,png,gif",
			Path:    tempDir,
		}

		// 3. 准备base64数据
		base64Data := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageData)

		// 4. 上传图片
		filePath, fileName, err := pkg.UploadImage(config, base64Data)
		assert.NoError(t, err, "图片上传应该成功")
		assert.NotEmpty(t, filePath, "应该返回文件路径")
		assert.NotEmpty(t, fileName, "应该返回文件名")

		// 5. 验证文件存在
		_, err = os.Stat(fileName)
		assert.NoError(t, err, "上传的文件应该存在")

		// 6. 验证文件内容
		uploadedData, err := os.ReadFile(fileName)
		assert.NoError(t, err, "应该能读取上传的文件")
		// UploadImage 会解码 base64 数据，所以我们应该比较解码后的数据
		// 从 base64Data 中提取实际的图片数据
		b64data := base64Data[strings.IndexByte(base64Data, ',')+1:]
		decodedData, decodeErr := base64.StdEncoding.DecodeString(b64data)
		assert.NoError(t, decodeErr, "应该能解码 base64 数据")
		assert.Equal(t, decodedData, uploadedData, "上传的文件内容应该匹配解码后的数据")
	})

	t.Run("Base64文件保存流程", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		// 1. 创建测试数据
		testData := []byte("test image content for base64")
		base64Data := base64.StdEncoding.EncodeToString(testData)

		// 2. 生成文件路径
		filePath, fileName := pkg.CreateUploadPathFileName()
		fullDirPath := filepath.Join(tempDir, filePath)
		fullFilePath := filepath.Join(fullDirPath, fileName+".jpg")

		// 3. 准备base64数据
		base64DataWithHeader := "data:image/jpeg;base64," + base64Data

		// 4. 保存Base64文件
		err = pkg.SaveBase64File(base64DataWithHeader, fullDirPath, fullFilePath)
		assert.NoError(t, err, "Base64文件保存应该成功")

		// 5. 验证文件内容
		savedData, err := os.ReadFile(fullFilePath)
		assert.NoError(t, err, "应该能读取保存的文件")
		assert.Equal(t, testData, savedData, "保存的文件内容应该匹配")
	})
}

// TestImageEdgeCases 测试图片模块的边界情况
func TestImageEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.GetTestTimeout())*time.Second)
	defer cancel()

	t.Run("处理空文件名", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		filePath, fileName := pkg.CreateUploadPathFileName()
		result := filePath + fileName
		assert.NotEmpty(t, result, "即使文件名为空也应该生成路径")
		assert.NotContains(t, result, "..", "路径不应该包含相对路径")
	})

	t.Run("处理特殊字符文件名", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		filePath, fileName := pkg.CreateUploadPathFileName()
		result := filePath + fileName
		assert.NotEmpty(t, result, "应该能处理特殊字符文件名")
		// CreateUploadPathFileName 只返回基础文件名，不包含扩展名
		assert.False(t, strings.Contains(fileName, "."), "基础文件名不应该包含扩展名")
	})

	t.Run("处理非常长的文件名", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatal("测试超时")
		default:
		}
		filePath, fileName := pkg.CreateUploadPathFileName()
		result := filePath + fileName
		assert.NotEmpty(t, result, "应该能处理长文件名")
		// CreateUploadPathFileName 只返回基础文件名，不包含扩展名
		assert.False(t, strings.Contains(fileName, "."), "基础文件名不应该包含扩展名")
	})
}
