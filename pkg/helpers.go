package pkg

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"net/smtp"
	"strings"
)

const (
	RAND_KIND_NUM   = iota // 纯数字
	RAND_KIND_LOWER        // 小写字母
	RAND_KIND_UPPER        // 大写字母
	RAND_KIND_ALL          // 数字、大小写字母
)

var (
	// 全局随机数生成器，在包初始化时设置种子
	globalRand *mathrand.Rand
)

func init() {
	// 使用加密安全的随机数作为种子
	seedBytes := make([]byte, 8)
	cryptorand.Read(seedBytes)
	seed := int64(0)
	for i, b := range seedBytes {
		seed |= int64(b) << (i * 8)
	}
	globalRand = mathrand.New(mathrand.NewSource(seed))
}

// Rand 函数生成指定长度和类型的随机字符串。
// size 参数指定了生成的随机字符串的长度。
// kind 参数指定了生成的随机字符串的类型，可以是 RAND_KIND_NUM、RAND_KIND_LOWER、RAND_KIND_UPPER 或 RAND_KIND_ALL。
// 函数返回生成的随机字符串。
func Rand(size int, kind int) string {
	if size <= 0 {
		return ""
	}
	
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		switch kind {
		case RAND_KIND_NUM:
			result[i] = '0' + uint8(globalRand.Intn(10))
		case RAND_KIND_LOWER:
			result[i] = 'a' + uint8(globalRand.Intn(26))
		case RAND_KIND_UPPER:
			result[i] = 'A' + uint8(globalRand.Intn(26))
		case RAND_KIND_ALL:
			fallthrough
		default:
			switch globalRand.Intn(3) {
			case 0:
				result[i] = '0' + uint8(globalRand.Intn(10))
			case 1:
				result[i] = 'a' + uint8(globalRand.Intn(26))
			case 2:
				result[i] = 'A' + uint8(globalRand.Intn(26))
			}
		}
	}
	return string(result)
}

// SendToMail 函数用于发送电子邮件。
// user 参数是发送邮件的用户的邮箱地址。
// password 参数是发送邮件的用户的邮箱密码。
// host 参数是发送邮件的 SMTP 服务器的主机名和端口号。
// to 参数是接收邮件的用户的邮箱地址。
// subject 参数是邮件的主题。
// body 参数是邮件的正文。
// mailtype 参数是邮件的类型，可以是 "html" 或者 "plain"。
// 函数返回一个 error，如果发送邮件过程中没有错误发生，error 为 nil。
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

// IsValidUserName 和 IsValidPassword 函数已移至 regexcheck.go
// 请使用 pkg.IsUserName 和 pkg.IsPwd 替代
