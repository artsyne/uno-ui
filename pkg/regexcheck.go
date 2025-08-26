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

//检查用户名是否合法，数字+字母  不限制大小写 6~30位
func IsUserName(str string) bool {
	b, _ := regexp.MatchString("^[0-9a-zA-Z]{6,30}$", str)
	return b
}

// IsPwd 检查密码是否合法
// 要求：8-30位字母、数字或特殊字符
func IsPwd(str string) bool {
	b, _ := regexp.MatchString("^[\\x21-\\x7E]{8,30}$", str)
	return b
}

/************************* 数字类型 ************************/
//纯整数
func IsInteger(str string) bool {
	b, _ := regexp.MatchString("^[0-9]+$", str)
	return b
}

//纯小数
func IsDecimals(str string) bool {
	b, _ := regexp.MatchString("^\\d+\\.[0-9]+$", str)
	return b
}

//手提电话（不带前缀）最高11位
func IsCellphone(str string) bool {
	b, _ := regexp.MatchString("^1[0-9]{10}$", str)
	return b
}

//家用电话（不带前缀） 最高8位
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

//仅大写
func IsEnglishCap(str string) bool {
	b, _ := regexp.MatchString("^[A-Z]+$", str)
	return b
}

//大小写混合
func IsEnglish(str string) bool {
	b, _ := regexp.MatchString("^[A-Za-z]+$", str)
	return b
}

//邮箱
func IsEmail(str string) bool {
	b, _ := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", str)
	return b
}