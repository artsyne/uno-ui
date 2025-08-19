package pkg

import (
	"regexp"
)

type RegexCheck struct{}

/************************* 自定义类型 ************************/
//数字+字母  不限制大小写 6~30位
func (ic *RegexCheck) IsID(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[0-9a-zA-Z]{6,30}$", s)
		if !b {
			return false
		}
	}
	return true
}

//检查用户名是否合法，数字+字母  不限制大小写 6~30位
func (ic *RegexCheck) IsUserName(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[0-9a-zA-Z]{6,30}$", s)
		if !b {
			return false
		}
	}
	return true
}

//数字+字母+符号 6~30位
func (ic *RegexCheck) IsPwd(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[0-9a-zA-Z@.]{6,30}$", s)
		if !b {
			return false
		}
	}
	return true
}


/************************* 数字类型 ************************/
//纯整数
func (ic *RegexCheck) IsInteger(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[0-9]+$", s)
		if !b {
			return false
		}
	}
	return true
}

//纯小数
func (ic *RegexCheck) IsDecimals(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^\\d+\\.[0-9]+$", s)
		if !b {
			return false
		}
	}
	return true
}

//手提电话（不带前缀）最高11位
func (ic *RegexCheck) IsCellphone(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^1[0-9]{10}$", s)
		if !b {
			return false
		}
	}
	return true
}

//家用电话（不带前缀） 最高8位
func (ic *RegexCheck) IsTelephone(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[0-9]{8}$", s)
		if !b {
			return false
		}
	}
	return true
}

/************************* 英文类型 *************************/
//仅小写
func (ic *RegexCheck) IsEngishLowCase(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[a-z]+$", s)
		if !b {
			return false
		}
	}
	return true
}

//仅大写
func (ic *RegexCheck) IsEnglishCap(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[A-Z]+$", s)
		if !b {
			return false
		}
	}
	return true
}

//大小写混合
func (ic *RegexCheck) IsEnglish(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		b, _ := regexp.MatchString("^[A-Za-z]+$", s)
		if !b {
			return false
		}
	}
	return true
}

//邮箱 最高30位
func (ic *RegexCheck) IsEmail(str ...string) bool {
	// 如果没有参数，返回true
	if len(str) == 0 {
		return true
	}
	
	for _, s := range str {
		// 基本的邮箱格式检查
		b, _ := regexp.MatchString("^[a-zA-Z0-9_.-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}$", s)
		if !b {
			return false
		}
		// 额外检查：不允许连续的点
		if regexp.MustCompile("\\.\\.").MatchString(s) {
			return false
		}
		// 不允许以点开头或结尾的用户名部分
		parts := regexp.MustCompile("@").Split(s, 2)
		if len(parts) == 2 {
			userPart := parts[0]
			if len(userPart) > 0 && (userPart[0] == '.' || userPart[len(userPart)-1] == '.') {
				return false
			}
		}
	}
	return true
}
