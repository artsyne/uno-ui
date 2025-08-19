package pkg

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/xid"
)

// UploadConfig 上传配置结构
type UploadConfig struct {
	Path    string // 上传路径
	Types   string // 允许的文件类型
	MaxSize int    // 最大文件大小
}

// UploadError 上传错误类型
type UploadError struct {
	Code    int
	Message string
	Details string
}

func (e *UploadError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// 预定义的错误
var (
	ErrUploadFailed          = &UploadError{Code: 1001, Message: "上传失败"}
	ErrFileFormatUnsupported = &UploadError{Code: 1002, Message: "文件格式不支持"}
	ErrFileTooLarge          = &UploadError{Code: 1003, Message: "文件太大"}
)

// WithDetails 添加错误详情
func (e *UploadError) WithDetails(details string) *UploadError {
	return &UploadError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

// UploadImage 是一个处理图片上传的函数。
// 它接受一个 base64 编码的图像数据作为参数，将图像保存到服务器的指定位置，
// 并返回保存的文件的路径和文件名。如果在过程中出现错误，它会返回一个错误代码。
// 参数：
// config: 上传配置
// imgBase64Data: base64 编码的图像数据。
// 返回值：
// filePath: 保存的文件的路径。
// fileName: 保存的文件的名称。
// error: 错误信息，如果没有错误发生，此值为nil。
func UploadImage(config *UploadConfig, imgBase64Data string) (string, string, error) {

	//判断图片是否为空
	if len(imgBase64Data) < 30 { // 至少需要包含 "data:image/xxx;base64," 这样的头部
		return "", "", ErrUploadFailed
	}

	//判断上传文件类型是否正确
	// 检查数据格式是否正确
	if !strings.HasPrefix(imgBase64Data, "data:") || !strings.Contains(imgBase64Data, ";base64,") {
		return "", "", ErrUploadFailed.WithDetails("invalid base64 format")
	}
	
	// 提取 MIME 类型，例如从 "data:image/jpeg;base64," 中提取 "image/jpeg"
	mimeType := imgBase64Data[5:strings.IndexByte(imgBase64Data, ';')] // 去掉 "data:" 前缀
	// 从 MIME 类型中提取文件扩展名
	var fileExt string
	switch mimeType {
	case "image/jpeg":
		fileExt = "jpg"
	case "image/png":
		fileExt = "png"
	case "image/gif":
		fileExt = "gif"
	default:
		return "", "", ErrFileFormatUnsupported
	}
	// 检查配置中是否支持该文件类型
	if !strings.Contains(config.Types, fileExt) {
		return "", "", ErrFileFormatUnsupported
	}

	// 解码base64数据以检查实际文件大小
	b64data := imgBase64Data[strings.IndexByte(imgBase64Data, ',')+1:]
	imgByte, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", "", ErrUploadFailed.WithDetails("invalid base64 data")
	}
	
	//判断文件大小
	if len(imgByte) > config.MaxSize {
		return "", "", ErrFileTooLarge.WithDetails("file too large")
	}

	//路径格式为：config配置路径 + /年/月/日/uid.jpg
	relativePath, baseFileName := CreateUploadPathFileName()
	filePath := path.Join(config.Path, relativePath)
	fileName := fmt.Sprintf("%s.%s", baseFileName, fileExt)
	fullFileName := path.Join(filePath, fileName)

	//保存图片
	err = SaveImageFile(imgByte, filePath, fullFileName)
	if err != nil {
		return "", "", ErrUploadFailed.WithDetails(err.Error())
	}

	// 返回相对路径和完整文件路径（测试期望fileName是完整路径）
	return relativePath, fullFileName, nil
}

// SaveImageFile 保存已解码的图片数据到文件
func SaveImageFile(imgData []byte, filePath, fileName string) error {
	//创建文件路径
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果路径不存在，则创建目录
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//文件保存
	err := os.WriteFile(fileName, imgData, 0666)
	if err != nil {
		return err
	}

	return nil
}

// SaveBase64File 是一个用于处理 base64 编码的文件数据的函数。
// 它接收一个 base64 编码的文件数据作为参数，将文件保存到服务器的指定位置，
// 并返回保存的文件的路径和文件名。如果在处理过程中出现错误，它将返回错误信息。
// 参数：
// base64FileData: base64 编码的文件数据。
// 返回值：
// filePath: 保存的文件的路径。
// fileName: 保存的文件的名字。
// err: 如果在处理过程中出现错误，err 将包含错误信息。
func SaveBase64File(base64Data, filePath, fileName string) error {
	// 验证输入参数
	if base64Data == "" {
		return fmt.Errorf("base64 data cannot be empty")
	}
	
	// 检查是否包含逗号分隔符
	commaIndex := strings.IndexByte(base64Data, ',')
	if commaIndex == -1 {
		return fmt.Errorf("invalid base64 data format: missing comma separator")
	}
	
	//去掉,号前内容，"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYA"
	b64data := base64Data[commaIndex+1:]
	
	// 验证base64数据不为空
	if b64data == "" {
		return fmt.Errorf("base64 data content cannot be empty")
	}

	//图片解码
	imgByte, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return err
	}

	//创建文件路径
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果路径不存在，则创建目录
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//文件保存
	err = os.WriteFile(fileName, imgByte, 0666)
	if err != nil {
		return err
	}

	return nil
}

// CreateUploadPathFileName 是一个用于创建上传文件的路径和文件名的函数。
// 它根据当前的日期和时间生成一个唯一的路径和文件名，以避免文件名冲突。
// 这个函数不接受任何参数。
// 返回值：
// uploadPath: 生成的文件路径。
// fileName: 生成的文件名。
// err: 如果在处理过程中出现错误，err 将包含错误信息。
func CreateUploadPathFileName() (string, string) {
	var nowTime time.Time = GetNowTime()
	guid := xid.New()

	filePath := fmt.Sprintf("%d/%02d/%02d/",
		nowTime.Year(), nowTime.Month(), nowTime.Day())

	return filePath, guid.String()
}
