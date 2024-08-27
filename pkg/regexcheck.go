package pkg

import (
	"regexp"
)

type RegexCheck struct{}

/************************* 自定义类型 ************************/
//数字+字母  不限制大小写 6~30位
func (ic *RegexCheck) IsID(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[0-9a-zA-Z]{6,30}$", s)
		if false == b {
			return b
		}
	}
	return b
}

//数字+字母+符号 6~30位
func (ic *RegexCheck) IsPwd(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[0-9a-zA-Z@.]{6,30}$", s)
		if false == b {
			return b
		}
	}
	return b
}

/************************* 数字类型 ************************/
//纯整数
func (ic *RegexCheck) IsInteger(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[0-9]+$", s)
		if false == b {
			return b
		}
	}
	return b
}

//纯小数
func (ic *RegexCheck) IsDecimals(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^\\d+\\.[0-9]+$", s)
		if false == b {
			return b
		}
	}
	return b
}

//手提电话（不带前缀）最高11位
func (ic *RegexCheck) IsCellphone(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^1[0-9]{10}$", s)
		if false == b {
			return b
		}
	}
	return b
}

//家用电话（不带前缀） 最高8位
func (ic *RegexCheck) IsTelephone(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[0-9]{8}$", s)
		if false == b {
			return b
		}
	}
	return b
}

/************************* 英文类型 *************************/
//仅小写
func (ic *RegexCheck) IsEngishLowCase(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[a-z]+$", s)
		if false == b {
			return b
		}
	}
	return b
}

//仅大写
func (ic *RegexCheck) IsEnglishCap(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[A-Z]+$", s)
		if false == b {
			return b
		}
	}
	return b
}

//大小写混合
func (ic *RegexCheck) IsEnglish(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[A-Za-z]+$", s)
		if false == b {
			return b
		}
	}
	return b
}

//邮箱 最高30位
func (ic *RegexCheck) IsEmail(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$", s)
		if false == b {
			return b
		}
	}
	return b
}
