package helpers

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/rs/xid"
)

// 获取当前时间
func GetNowTime() time.Time {
	return time.Now()
}

const (
	RAND_KIND_NUM   = 0 // 纯数字
	RAND_KIND_LOWER = 1 // 小写字母
	RAND_KIND_UPPER = 2 // 大写字母
	RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func Rand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

// 发送电子邮件，mailtype可以为html
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

// 根据base64编码保存用户上传图片
// base64Data：base64编码的图片数据
// filePath：保存路径，路径格式为：aaa/bbb/ccc/
// fileName：保存文件名，文件名需带路径,格式为：aaa/bbb/ccc/ddd.xxx
func SaveBase64Img(base64Data, filePath, fileName string) error {
	//去掉,号前内容，"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYA"
	b64data := base64Data[strings.IndexByte(base64Data, ',')+1:]

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

// 创建上传文件图片路径和文件名
// 路径格式为：年/月/日/，文件名为guid
func CreateUploadPathFileName() (string, string) {
	nowTime := GetNowTime()
	guid := xid.New()

	filePath := fmt.Sprintf("%d/%d/%d/",
		nowTime.Year(), nowTime.Month(), nowTime.Day())

	return filePath, guid.String()
}
