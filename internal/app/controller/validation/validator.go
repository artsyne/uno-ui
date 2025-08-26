package validation

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/biuren/hebe/pkg"
)

// Validator 全局校验器实例
var Validator *validator.Validate

// init 初始化校验器
func init() {
	Validator = validator.New()
	
	// 注册自定义校验器
	registerCustomValidators(Validator)
}

// registerCustomValidators 注册自定义校验器
func registerCustomValidators(v *validator.Validate) {
	// 用户名格式校验
	v.RegisterValidation("username", validateUsername)
	
	// 密码强度校验
	v.RegisterValidation("password", validatePassword)
	

}

// validateUsername 验证用户名格式
// 使用pkg库的统一校验规则
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return pkg.IsUserName(username)
}

// validatePassword 验证密码格式
// 使用pkg库的统一校验规则
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return pkg.IsPwd(password)
}



// GetValidationError 获取校验错误的错误码
func GetValidationError(err error) *errors.AppError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 返回第一个校验错误的错误码
		for _, fieldError := range validationErrors {
			return getFieldErrorCode(fieldError)
		}
	}
	// 如果不是校验错误，返回通用校验错误
	return errors.ValidationError.Err()
}

// getFieldErrorCode 获取字段错误的错误码
func getFieldErrorCode(fe validator.FieldError) *errors.AppError {
	switch fe.Tag() {
	case "required":
		return errors.MissingRequiredField.Err()
	case "email":
		return errors.InvalidEmail.Err()
	case "min":
		return errors.FieldTooShort.Err()
	case "max":
		return errors.FieldTooLong.Err()
	case "len":
		return errors.ValueOutOfRange.Err()
	case "eqfield":
		return errors.FieldMismatch.Err()
	case "username":
		return errors.InvalidUsername.Err()
	case "password":
		return errors.InvalidPassword.Err()

	default:
		return errors.ValidationError.Err()
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
		return 0, errors.InvalidParams.Err()
	}
	return id, nil
}

// ValidateUsernameFormat 验证用户名格式
// 检查用户名是否符合系统要求的格式
//
// 参数:
//   - userName: 待验证的用户名
//
// 返回值:
//   - bool: 格式是否正确
func ValidateUsernameFormat(userName string) bool {
	return pkg.IsUserName(userName)
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