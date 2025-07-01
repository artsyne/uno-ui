package pkg

import (
	"math/rand"
	"net/smtp"
	"strings"
	"time"
	"unicode"
)

const (
	RAND_KIND_NUM   = iota // 纯数字
	RAND_KIND_LOWER        // 小写字母
	RAND_KIND_UPPER        // 大写字母
	RAND_KIND_ALL          // 数字、大小写字母
)

// Rand 函数生成指定长度和类型的随机字符串。
// size 参数指定了生成的随机字符串的长度。
// kind 参数指定了生成的随机字符串的类型，可以是 RAND_KIND_NUM、RAND_KIND_LOWER、RAND_KIND_UPPER 或 RAND_KIND_ALL。
// 函数返回生成的随机字符串。
func Rand(size int, kind int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		switch kind {
		case RAND_KIND_NUM:
			result[i] = '0' + uint8(r.Intn(10))
		case RAND_KIND_LOWER:
			result[i] = 'a' + uint8(r.Intn(26))
		case RAND_KIND_UPPER:
			result[i] = 'A' + uint8(r.Intn(26))
		case RAND_KIND_ALL:
			fallthrough
		default:
			switch r.Intn(3) {
			case 0:
				result[i] = '0' + uint8(r.Intn(10))
			case 1:
				result[i] = 'a' + uint8(r.Intn(26))
			case 2:
				result[i] = 'A' + uint8(r.Intn(26))
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

// IsValidUserName 验证用户名是否有效
// userName: 用户名
// 返回值：如果用户名有效，返回true，否则返回false
// 用户名必须至少包含6个字符，最多包含20个字符，只能包含字母、数字和下划线
func IsValidUserName(userName string) bool {
	if len(userName) < 6 || len(userName) > 20 {
		return false
	}
	for _, ch := range userName {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}
	return true
}

// IsValidPassword 验证密码是否有效
// password: 密码
// 返回值：如果密码有效，返回true，否则返回false
// 密码必须至少包含8个字符，最多包含30个字符，必须包含至少一个字母和一个数字
func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 30 {
		return false
	}
	hasLetter := false
	hasDigit := false
	for _, ch := range password {
		if unicode.IsLetter(ch) {
			hasLetter = true
		} else if unicode.IsDigit(ch) {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
