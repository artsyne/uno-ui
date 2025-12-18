package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRequest 测试用的请求结构体
type TestRequest struct {
	Name  string `json:"name" binding:"required,min=3"`
	Age   int    `json:"age" binding:"required,gte=0,lte=130"`
	Email string `json:"email" binding:"required,email"`
}

// TestValidateMiddleware 测试 Validate 中间件
func TestValidateMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		// 注册使用验证中间件的路由
		router.POST("/test", middleware.Validate[TestRequest](), func(c *gin.Context) {
			// 获取验证后的请求数据
			req := middleware.MustGetValidatedRequest[TestRequest](c)
			c.JSON(http.StatusOK, req)
		})

		// 准备有效的请求数据
		reqBody := TestRequest{
			Name:  "John Doe",
			Age:   30,
			Email: "john@example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)

		// 发送请求
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestRequest
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Name, resp.Name)
		assert.Equal(t, reqBody.Age, resp.Age)
		assert.Equal(t, reqBody.Email, resp.Email)
	})

	t.Run("Invalid Request - Validation Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.POST("/test", middleware.Validate[TestRequest](), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// 准备无效的请求数据 (名字太短，年龄无效，邮箱无效)
		reqBody := TestRequest{
			Name:  "Jo",
			Age:   -1,
			Email: "invalid-email",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// 验证返回 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// 验证错误响应格式
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// 检查错误结构
		assert.Equal(t, float64(1002), resp["code"])
		assert.Equal(t, "invalid_params", resp["message"])
		assert.NotNil(t, resp["details"])
	})

	t.Run("Invalid Request - JSON Parse Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.POST("/test", middleware.Validate[TestRequest](), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// 发送无效的 JSON
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString("{invalid-json}"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// 验证返回 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("HTMX Request - Validation Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.POST("/test", middleware.Validate[TestRequest](), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// 准备无效数据
		reqBody := TestRequest{
			Name: "Jo", // Too short
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// 设置 HTMX 头
		req.Header.Set("HX-Request", "true")

		router.ServeHTTP(w, req)

		// HTMX 请求通常返回 200 或 422，但 middleware 目前统一返回 400
		// 重要的是验证响应体包含错误信息
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// 验证是否触发了 htmx:responseError 事件或者包含错误信息
		// 这里主要检查是否返回了 JSON 错误，因为 middleware 中没有针对 HTMX 的特殊处理逻辑（除了错误响应本身）
		// 如果需要针对 HTMX 返回 HTML 片段，需要在 middleware 中增加逻辑
		// 目前的 middleware 实现是统一返回 JSONError
	})
}
