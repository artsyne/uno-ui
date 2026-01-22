package pkg

import (
	"regexp"
)

/************************* 自定义类型 ************************/
//数字+字母  不限制大小写 6~30位
func IsID(str string) bool {
	b, _ := regexp.MatchString("^[0-9a-zA-Z]{6,30}$", str)
	return b
}

// 检查用户名是否合法，数字+字母  不限制大小写 3~16位
func IsUserName(str string) bool {
	b, _ := regexp.MatchString("^[0-9a-zA-Z_-]{3,16}$", str)
	return b
}

// IsPwd 检查密码是否合法
// 要求：密码需 8-80 位，必须包含字母和数字，可含_@#$%^&*()+-= 特殊字符
func IsPwd(str string) bool {
	// Go regex (RE2) doesn't support lookarounds, so we check conditions separately

	// 1. Check length and allowed characters
	// Note: escaped - as \- to be safe, though strict position works too
	if matched, _ := regexp.MatchString(`^[A-Za-z\d!_@#$%^&*()+\-=]{8,80}$`, str); !matched {
		return false
	}

	// 2. Check for at least one letter
	if matched, _ := regexp.MatchString(`[A-Za-z]`, str); !matched {
		return false
	}

	// 3. Check for at least one digit
	if matched, _ := regexp.MatchString(`\d`, str); !matched {
		return false
	}

	return true
}

/************************* 数字类型 ************************/
//纯整数
func IsInteger(str string) bool {
	b, _ := regexp.MatchString("^[0-9]+$", str)
	return b
}

// 纯小数
func IsDecimals(str string) bool {
	b, _ := regexp.MatchString("^\\d+\\.[0-9]+$", str)
	return b
}

// 手提电话（不带前缀）最高11位
func IsCellphone(str string) bool {
	b, _ := regexp.MatchString("^1[0-9]{10}$", str)
	return b
}

// 家用电话（不带前缀） 最高8位
func IsTelephone(str string) bool {
	b, _ := regexp.MatchString("^[0-9]{8}$", str)
	return b
}

/************************* 英文类型 *************************/
//仅小写
func IsEngishLowCase(str string) bool {
	b, _ := regexp.MatchString("^[a-z]+$", str)
	return b
}

// 仅大写
func IsEnglishCap(str string) bool {
	b, _ := regexp.MatchString("^[A-Z]+$", str)
	return b
}

// 大小写混合
func IsEnglish(str string) bool {
	b, _ := regexp.MatchString("^[A-Za-z]+$", str)
	return b
}

// 邮箱
func IsEmail(str string) bool {
	emailRegex := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)+$"
	b, _ := regexp.MatchString(emailRegex, str)
	return b
}
