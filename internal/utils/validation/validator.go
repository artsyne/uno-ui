package validation

import (
	"strings"
	"sync"

	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/go-playground/validator/v10"
)

var (
	// 全局验证器实例
	once         sync.Once
	validatorV10 *validator.Validate
)

// GetValidator 获取全局单例验证器
// 使用 sync.Once 确保只初始化一次
func GetValidator() *validator.Validate {
	once.Do(func() {
		validatorV10 = validator.New()
		// registerCustomValidations(validatorV10) // 业务方应该在自己的模块中注册
	})
	return validatorV10
}

// ValidateStruct 验证结构体
// 对给定的结构体执行验证，并将验证错误转换为统一的错误格式
//
// 参数：
//
//	s: 要验证的结构体实例
//
// 返回值：
//
//	error: 验证错误，如果验证通过返回 nil
func ValidateStruct(s interface{}) error {
	if err := GetValidator().Struct(s); err != nil {
		return errors.ValidationError.WithDetails(err.Error())
	}
	return nil
}

// ValidateVar 验证单个变量
// 对单个变量执行验证规则
//
// 参数：
//
//	field: 要验证的变量
//	tag: 验证标签（如 "required,min=1,max=100"）
//
// 返回值：
//
//	error: 验证错误，如果验证通过返回 nil
func ValidateVar(field interface{}, tag string) error {
	if err := GetValidator().Var(field, tag); err != nil {
		return errors.ValidationError.WithDetails(err.Error())
	}
	return nil
}

// RegisterValidationRule 注册自定义验证规则
// 提供更灵活的验证规则注册方式
//
// 参数：
//
//	tag: 验证标签名
//	fn: 验证函数
//	callValidationEvenIfNull: 即使字段为空也执行验证
//
// 返回值：
//
//	error: 注册错误
func RegisterValidationRule(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return GetValidator().RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

// ValidateWithRules 使用指定规则验证变量
// 支持动态验证规则组合
//
// 参数：
//
//	field: 要验证的变量
//	rules: 验证规则数组
//
// 返回值：
//
//	error: 验证错误
func ValidateWithRules(field interface{}, rules []string) error {
	tag := strings.Join(rules, ",")
	return ValidateVar(field, tag)
}
