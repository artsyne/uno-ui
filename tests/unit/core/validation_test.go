package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/utils/validation"
	testhelpers "github.com/biuren/hebe/tests/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestValidateStruct 测试结构体验证
func TestValidateStruct(t *testing.T) {
	// 使用配置管理器确保配置正确初始化
	testhelpers.SetupTestConfigSafe(t)
	// 测试结构体
	type TestUser struct {
		Username string `validate:"required,min=3,max=20"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"min=18,max=100"`
	}

	t.Run("有效数据验证通过", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "test@example.com",
			Age:      25,
		}

		err := validation.ValidateStruct(user)
		assert.NoError(t, err)
	})

	t.Run("用户名太短", func(t *testing.T) {
		user := TestUser{
			Username: "ab", // 少于3个字符
			Email:    "test@example.com",
			Age:      25,
		}

		err := validation.ValidateStruct(user)
		assert.Error(t, err)
	})

	t.Run("邮箱格式无效", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "invalid-email", // 无效邮箱格式
			Age:      25,
		}

		err := validation.ValidateStruct(user)
		assert.Error(t, err)
	})

	t.Run("年龄超出范围", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "test@example.com",
			Age:      17, // 小于18岁
		}

		err := validation.ValidateStruct(user)
		assert.Error(t, err)
	})

	t.Run("必填字段为空", func(t *testing.T) {
		user := TestUser{
			Username: "", // 空用户名
			Email:    "test@example.com",
			Age:      25,
		}

		err := validation.ValidateStruct(user)
		assert.Error(t, err)
	})
}

// TestValidateVar 测试单个变量验证
func TestValidateVar(t *testing.T) {
	t.Run("邮箱验证", func(t *testing.T) {
		validEmails := []string{
			"user@example.com",
			"test.user@domain.co.uk",
			"user+tag@example.org",
		}

		for _, email := range validEmails {
			t.Run("有效邮箱: "+email, func(t *testing.T) {
				err := validation.ValidateVar(email, "email")
				assert.NoError(t, err, "邮箱 %s 应该是有效的", email)
			})
		}

		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"user@",
			"",
		}

		for _, email := range invalidEmails {
			t.Run("无效邮箱: "+email, func(t *testing.T) {
				err := validation.ValidateVar(email, "email")
				assert.Error(t, err, "邮箱 %s 应该是无效的", email)
			})
		}
	})

	t.Run("最小长度验证", func(t *testing.T) {
		// 测试最小长度为3
		err := validation.ValidateVar("ab", "min=3")
		assert.Error(t, err)

		err = validation.ValidateVar("abc", "min=3")
		assert.NoError(t, err)

		err = validation.ValidateVar("abcd", "min=3")
		assert.NoError(t, err)
	})

	t.Run("必填字段验证", func(t *testing.T) {
		err := validation.ValidateVar("", "required")
		assert.Error(t, err)

		err = validation.ValidateVar("value", "required")
		assert.NoError(t, err)
	})

	t.Run("数字范围验证", func(t *testing.T) {
		err := validation.ValidateVar(10, "min=18")
		assert.Error(t, err)

		err = validation.ValidateVar(25, "min=18,max=100")
		assert.NoError(t, err)

		err = validation.ValidateVar(150, "max=100")
		assert.Error(t, err)
	})
}

// TestGetValidator 测试全局验证器
func TestGetValidator(t *testing.T) {
	t.Run("获取验证器实例", func(t *testing.T) {
		validator1 := validation.GetValidator()
		validator2 := validation.GetValidator()

		assert.NotNil(t, validator1)
		assert.NotNil(t, validator2)

		// 验证单例模式：两次获取应该是同一个实例
		assert.Equal(t, validator1, validator2)
	})
}

// TestValidationWithComplexStruct 测试复杂结构体验证
func TestValidationWithComplexStruct(t *testing.T) {
	// 模拟用户注册数据结构
	type UserProfile struct {
		FirstName string `validate:"required,min=1"`
		LastName  string `validate:"required,min=1"`
		Phone     string `validate:"omitempty,e164"` // 可选，但如果提供必须是有效的E164格式
	}

	type UserRegistration struct {
		Username string      `validate:"required,min=3,max=20,alphanum"`
		Email    string      `validate:"required,email"`
		Password string      `validate:"required,min=8"`
		Age      int         `validate:"min=13,max=120"`
		Profile  UserProfile `validate:"required"`
	}

	t.Run("完整有效数据", func(t *testing.T) {
		registration := UserRegistration{
			Username: "testuser123",
			Email:    "test@example.com",
			Password: "securepassword",
			Age:      25,
			Profile: UserProfile{
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+1234567890",
			},
		}

		err := validation.ValidateStruct(registration)
		assert.NoError(t, err)
	})

	t.Run("嵌套结构体验证错误", func(t *testing.T) {
		registration := UserRegistration{
			Username: "testuser123",
			Email:    "test@example.com",
			Password: "securepassword",
			Age:      25,
			Profile: UserProfile{
				FirstName: "", // 空的FirstName应该导致验证错误
				LastName:  "Doe",
			},
		}

		err := validation.ValidateStruct(registration)
		assert.Error(t, err)
	})
}

// TestValidationRequestHandler 测试验证请求处理器
func TestValidationRequestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 使用配置管理器确保配置正确初始化
	testhelpers.SetupTestConfigSafe(t)

	// 测试结构体
	type TestRequest struct {
		Name  string `validate:"required,min=3"`
		Email string `validate:"required,email"`
	}

	t.Run("验证成功的处理", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := TestRequest{
			Name:  "testuser",
			Email: "test@example.com",
		}

		err := validation.ValidateStruct(req)
		assert.NoError(t, err)

		// 模拟成功响应
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": req})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("验证失败的处理", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := TestRequest{
			Name:  "ab",            // 太短，验证失败
			Email: "invalid-email", // 无效邮箱
		}

		err := validation.ValidateStruct(req)
		assert.Error(t, err)

		// 模拟错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestValidationIntegration 测试验证器与其他模块的集成
func TestValidationIntegration(t *testing.T) {
	// 使用配置管理器确保配置正确初始化
	cfg := testhelpers.SetupTestConfigSafe(t)

	// 模拟用户数据验证
	type UserData struct {
		Username string `validate:"required,min=3,max=20"`
		Email    string `validate:"required,email"`
		TimeZone string `validate:"required"`
		Lang     string `validate:"required"`
	}

	t.Run("使用真实配置的验证", func(t *testing.T) {
		userData := UserData{
			Username: "testuser",
			Email:    "test@example.com",
			TimeZone: "Asia/Shanghai",
			Lang:     cfg.App.DefaultLang, // 使用配置中的默认语言
		}

		err := validation.ValidateStruct(userData)
		assert.NoError(t, err, "使用真实配置的验证应该通过")
	})

	t.Run("验证支持的语言", func(t *testing.T) {

		for _, lang := range cfg.SupportLangs {
			userData := UserData{
				Username: "testuser",
				Email:    "test@example.com",
				TimeZone: "UTC",
				Lang:     lang.Code,
			}

			err := validation.ValidateStruct(userData)
			assert.NoError(t, err, "支持的语言 %s 应该通过验证", lang.Code)
		}
	})
}

// TestValidationPerformance 测试验证器性能
func TestValidationPerformance(t *testing.T) {
	type TestStruct struct {
		Field1 string `validate:"required,min=3,max=50"`
		Field2 string `validate:"required,email"`
		Field3 int    `validate:"min=1,max=100"`
	}

	data := TestStruct{
		Field1: "validstring",
		Field2: "test@example.com",
		Field3: 50,
	}

	// 基准测试：验证1000次
	t.Run("验证性能测试", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			err := validation.ValidateStruct(data)
			assert.NoError(t, err)
		}
	})
}

// TestValidationErrorMessages 测试验证错误消息
func TestValidationErrorMessages(t *testing.T) {
	type TestStruct struct {
		Username string `validate:"required,min=3,max=20"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"min=18,max=100"`
	}

	t.Run("获取详细错误信息", func(t *testing.T) {
		invalidData := TestStruct{
			Username: "ab",            // 太短
			Email:    "invalid-email", // 无效格式
			Age:      17,              // 太小
		}

		err := validation.ValidateStruct(invalidData)
		assert.Error(t, err)

		// 验证错误消息包含详细信息
		errMsg := err.Error()
		assert.NotEmpty(t, errMsg)
		// 错误消息应该包含字段名和验证规则信息
		t.Logf("验证错误消息: %s", errMsg)
	})
}
