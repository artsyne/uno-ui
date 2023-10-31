package utils

import (
	"fmt"
	"path"
	"strings"

	"github.com/biuren/hebe/settings"
)

// uploadImage 是一个用于处理图片上传的函数。
// 它接收一个 base64 编码的图片数据作为参数，将图片保存到服务器的指定位置，
// 并返回保存的文件的路径和文件名。如果在处理过程中出现错误，它将返回错误码。
// 参数：
// imgBase64Data: base64 编码的图片数据。
// 返回值：
// filePath: 保存的文件的路径。
// fileName: 保存的文件的名字。
// errCode: 错误码，如果没有错误发生，该值为 0。
// 请注意，这个函数可能需要根据你的实际需求进行修改或扩展。
func UploadImage(imgBase64Data string) (filePath, fileName string, errCode int) {

	//判断图片是否为空
	fileSize := len(imgBase64Data)
	if fileSize < 100 {
		return "", "", UPLOAD_FILE_EMPTY
	}

	// 获取上传文件配置信息
	config := settings.GetConfig()

	//判断上传文件类型是否正确
	filetype := imgBase64Data[:strings.IndexByte(imgBase64Data, ';')]
	if !strings.Contains(config.Upload.ImgType, filetype) {
		return "","", UPLOAD_FILE_TYPE_ERR
	}

	//判断文件大小
	if fileSize > config.Upload.ImgMaxSize {
		return "","", UPLOAD_FILE_TOO_LARGE
	}

	//路径格式为：config配置路径 + /年/月/日/uid.jpg
	filePath, fileName = CreateUploadPathFileName()
	filePath = path.Join(config.Upload.ImgSavePath, filePath)
	fileName = fmt.Sprintf("%s/%s.jpeg", filePath, fileName)

	//保存图片
	err := SaveBase64File(imgBase64Data, filePath, fileName)
	if err != nil {
		return "","", UPLOAD_FILE_ERR
	}

	return filePath, fileName, 0
}
