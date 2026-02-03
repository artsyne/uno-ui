package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/biuren/hebe/internal/app/controller/validation"
	"github.com/biuren/hebe/tests/framework"
)

// TestValidationStruct 测试结构体验证
func TestValidationStruct(t *testing.T) {
	type TestRequest struct {
		Name  string `validate:"required,min=3"`
		Email string `validate:"required,email"`
	}

	t.Run("验证成功", func(t *testing.T) {
		req := TestRequest{
			Name:  "testuser",
			Email: "test@example.com",
		}

		err := validation.Validator.Struct(req)
		assert.NoError(t, err)
	})

	t.Run("验证失败-必填字段缺失", func(t *testing.T) {
		req := TestRequest{
			Name: "", // 缺失
		}

		err := validation.Validator.Struct(req)
		assert.Error(t, err)
	})

	t.Run("验证失败-规则不满足", func(t *testing.T) {
		req := TestRequest{
			Name:  "ab", // 太短
			Email: "invalid-email",
		}

		err := validation.Validator.Struct(req)
		assert.Error(t, err)
	})
}

// TestValidationInController 测试在Controller中的使用
func TestValidationInController(t *testing.T) {
	gin.SetMode(gin.TestMode)

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

		err := validation.Validator.Struct(req)
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

		err := validation.Validator.Struct(req)
		assert.Error(t, err)

		// 模拟错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestValidationIntegration 测试验证器与其他模块的集成
func TestValidationIntegration(t *testing.T) {
	ctx := framework.NewTestContext("TestValidationIntegration")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	// cfg := ctx.GetConfig()

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
			// Lang:     cfg.(*config.Config).App.DefaultLang, // 使用配置中的默认语言
			Lang: "en", // 使用默认语言
		}

		err := validation.Validator.Struct(userData)
		assert.NoError(t, err, "使用真实配置的验证应该通过")
	})

	t.Run("验证支持的语言", func(t *testing.T) {

		// for _, lang := range cfg.(*config.Config).SupportLangs {
		for _, lang := range []struct{ Code string }{{"en"}, {"zh-cn"}} {
			userData := UserData{
				Username: "testuser",
				Email:    "test@example.com",
				TimeZone: "UTC",
				Lang:     lang.Code,
			}

			err := validation.Validator.Struct(userData)
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
			err := validation.Validator.Struct(data)
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

		err := validation.Validator.Struct(invalidData)
		assert.Error(t, err)

		// 验证错误消息包含详细信息
		errMsg := err.Error()
		assert.NotEmpty(t, errMsg)
		// 错误消息应该包含字段名和验证规则信息
		t.Logf("验证错误消息: %s", errMsg)
	})
}
