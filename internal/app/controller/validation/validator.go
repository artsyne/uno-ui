package validation

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/pkg"
)

// Validator 全局校验器实例 (保留用于非Gin上下文的校验)
var Validator *validator.Validate

// init 初始化校验器
// 在包加载时自动执行，创建全局校验器实例并注册自定义校验规则到 Gin
//
// 功能:
//   - 创建validator实例
//   - 注册自定义校验器到 Gin binding
func init() {
	Validator = validator.New()

	// 注册自定义校验器到本地实例
	registerCustomValidators(Validator)

	// 注册自定义校验器到 Gin 的 binding 引擎
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidators(v)
	}
}

// registerCustomValidators 注册自定义校验器
// 向校验器实例注册项目特定的自定义校验规则
//
// 参数:
//   - v: validator校验器实例
//
// 功能:
//   - 注册用户名格式校验器
//   - 注册密码强度校验器
func registerCustomValidators(v *validator.Validate) {
	// 用户名格式校验
	if err := v.RegisterValidation("username", validateUsername); err != nil {
		panic(err)
	}

	// 密码强度校验
	if err := v.RegisterValidation("password", validatePassword); err != nil {
		panic(err)
	}

}

// validateUsername 验证用户名格式
// 使用pkg库的统一校验规则，确保用户名符合系统要求
//
// 参数:
//   - fl: validator字段级别信息
//
// 返回:
//   - bool: 用户名格式是否有效
//
// 校验规则:
//   - 使用pkg.IsUserName进行格式验证
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return pkg.IsHandle(username)
}

// validatePassword 验证密码格式
// 使用pkg库的统一校验规则，确保密码符合安全要求
//
// 参数:
//   - fl: validator字段级别信息
//
// 返回:
//   - bool: 密码格式是否有效
//
// 校验规则:
//   - 使用pkg.IsPwd进行强度验证
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return pkg.IsPwd(password)
}

// GetValidationError 获取校验错误的错误码
// 将validator校验错误转换为应用程序统一的错误格式
//
// 参数:
//   - err: validator产生的错误
//
// 返回:
//   - errors.AppError: 转换后的应用错误
//
// 功能:
//   - 解析validator.ValidationErrors
//   - 返回第一个字段错误的错误码
//   - 非校验错误返回通用校验错误
func GetValidationError(err error) errors.AppError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 返回第一个校验错误的错误码
		for _, fieldError := range validationErrors {
			return getFieldErrorCode(fieldError)
		}
	}
	// 如果不是校验错误，返回通用校验错误
	return errors.ValidationError
}

// getFieldErrorCode 获取字段错误的错误码
// 根据validator字段错误类型映射到对应的应用错误码
//
// 参数:
//   - fe: validator字段错误信息
//
// 返回:
//   - errors.AppError: 对应的应用错误
//
// 支持的错误类型:
//   - required: 必填字段缺失
//   - email: 邮箱格式错误
//   - min/max: 长度限制错误
//   - len: 长度不匹配
//   - eqfield: 字段不匹配
//   - username: 用户名格式错误
//   - password: 密码格式错误
func getFieldErrorCode(fe validator.FieldError) errors.AppError {
	switch fe.Tag() {
	case "required":
		return errors.MissingRequiredField
	case "email":
		return errors.InvalidEmail
	case "min":
		return errors.FieldTooShort
	case "max":
		return errors.FieldTooLong
	case "len":
		return errors.ValueOutOfRange
	case "eqfield":
		return errors.FieldMismatch
	case "username":
		return errors.InvalidUsername
	case "password":
		return errors.InvalidPassword

	default:
		return errors.ValidationError
	}
}

// ParseID 通用ID解析方法
// 从字符串中解析出有效的ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的ID
//   - error: 解析失败时的错误信息
func ParseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParams
	}
	return id, nil
}

// ValidateHandleFormat 验证Handle格式
// 检查Handle是否符合系统要求的格式
//
// 参数:
//   - handle: 待验证的Handle
//
// 返回值:
//   - bool: 格式是否正确
func ValidateHandleFormat(handle string) bool {
	return pkg.IsHandle(handle)
}

// ValidateEmailFormat 验证邮箱格式
// 检查邮箱地址是否符合标准格式
//
// 参数:
//   - email: 待验证的邮箱地址
//
// 返回值:
//   - bool: 格式是否正确
func ValidateEmailFormat(email string) bool {
	return pkg.IsEmail(email)
}
