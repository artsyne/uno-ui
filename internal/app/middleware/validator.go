package middleware

import (
	"fmt"
	"net/http"

	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validate 创建一个泛型验证中间件
// T: 请求参数绑定的目标结构体类型
func Validate[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		// 1. 尝试绑定 URI 参数 (如果有)
		// 忽略错误，因为不是所有请求都有 URI 参数，且 ShouldBindUri 在没有匹配时可能报错
		_ = c.ShouldBindUri(&req)

		// 2. 绑定 Body/Query 参数
		if err := c.ShouldBind(&req); err != nil {
			// 解析验证错误
			details := formatValidationError(err)

			// 构造标准错误响应
			errResponse := errors.NewAPIErrorResponse(
				errors.InvalidParams.Code(),
				errors.InvalidParams.Key(),
				details,
				getRequestID(c),
			)

			c.AbortWithStatusJSON(http.StatusBadRequest, errResponse)
			return
		}

		// 2. 将验证后的请求对象存入上下文
		c.Set("validatedRequest", &req)
		c.Next()
	}
}

// formatValidationError 格式化验证错误信息
func formatValidationError(err error) interface{} {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		result := make(map[string]string)
		for _, fieldError := range validationErrors {
			// 简单格式: FieldName: ErrorTag
			// 实际项目中可以结合 i18n 翻译错误信息
			fieldName := fieldError.Field()
			message := fmt.Sprintf("failed on tag '%s'", fieldError.Tag())
			result[fieldName] = message
		}
		return result
	}
	// 非验证错误（如 JSON 解析失败）
	return map[string]string{"error": err.Error()}
}

// GetValidatedRequest 从上下文中获取已验证的请求对象
func GetValidatedRequest[T any](c *gin.Context) *T {
	if v, exists := c.Get("validatedRequest"); exists {
		if req, ok := v.(*T); ok {
			return req
		}
	}
	return nil
}

// MustGetValidatedRequest 从上下文中获取已验证的请求对象，如果不存在则 panic
func MustGetValidatedRequest[T any](c *gin.Context) *T {
	req := GetValidatedRequest[T](c)
	if req == nil {
		panic("validatedRequest not found in context")
	}
	return req
}
