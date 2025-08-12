package core

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	// "github.com/biuren/hebe/internal/utils/config"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRecoveryMiddleware 测试恢复中间件
func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestRecoveryMiddleware")
	defer ctx.Cleanup()

	t.Run("panic恢复", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 添加恢复中间件
		router.Use(gin.Recovery())

		router.GET("/panic-test", func(c *gin.Context) {
			panic("intentional panic for testing")
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/panic-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证panic被恢复（应该返回500而不是崩溃程序）
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("正常请求不受影响", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 添加恢复中间件
		router.Use(gin.Recovery())

		router.GET("/normal-test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "normal"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/normal-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证正常响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "normal")
	})
}

// TestCustomErrorHandling 测试自定义错误处理
func TestCustomErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestCustomErrorHandling")
	defer ctx.Cleanup()

	// 自定义错误处理中间件
	errorHandlerMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Next()

			// 处理所有错误
			if len(c.Errors) > 0 {
				err := c.Errors.Last()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
					"code":  "INTERNAL_ERROR",
				})
			}
		}
	}

	t.Run("处理应用错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 应用错误处理中间件
		router.Use(errorHandlerMiddleware())

		// 设置会产生错误的处理器
		router.GET("/error", func(c *gin.Context) {
			// 模拟应用错误
			c.Error(errors.New("test application error"))
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/error", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证错误处理
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "test application error")
		assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
	})

	t.Run("正常请求不受影响", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 应用错误处理中间件
		router.Use(errorHandlerMiddleware())

		// 设置正常的处理器
		router.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/success", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证正常响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}

// TestMiddlewareChain 测试中间件链
func TestMiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestMiddlewareChain")
	defer ctx.Cleanup()

	t.Run("中间件执行顺序", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		var executionOrder []string

		// 添加多个中间件测试执行顺序
		router.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "middleware1")
			c.Next()
			executionOrder = append(executionOrder, "middleware1-after")
		})

		router.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "middleware2")
			c.Next()
			executionOrder = append(executionOrder, "middleware2-after")
		})

		router.GET("/order-test", func(c *gin.Context) {
			executionOrder = append(executionOrder, "handler")
			c.JSON(http.StatusOK, gin.H{"message": "order test"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/order-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证执行顺序
		assert.Equal(t, http.StatusOK, w.Code)
		expectedOrder := []string{
			"middleware1",
			"middleware2",
			"handler",
			"middleware2-after",
			"middleware1-after",
		}
		assert.Equal(t, expectedOrder, executionOrder)
	})

	t.Run("中间件中断请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 添加会中断请求的中间件
		router.Use(func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
		})

		router.GET("/abort-test", func(c *gin.Context) {
			// 这个处理器不应该被执行
			c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/abort-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证请求被中断
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
		assert.NotContains(t, w.Body.String(), "should not reach here")
	})
}

// TestAuthenticationMiddleware 测试认证中间件
func TestAuthenticationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestAuthenticationMiddleware")
	defer ctx.Cleanup()

	// 模拟认证中间件
	authMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
				c.Abort()
				return
			}

			if token != "Bearer valid-token" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}

			// 设置用户信息到上下文
			c.Set("userID", 123)
			c.Set("username", "testuser")
			c.Next()
		}
	}

	t.Run("缺少认证头", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(authMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
		})

		// 创建无认证头的请求
		req, _ := http.NewRequest("GET", "/protected", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证认证失败
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "missing authorization header")
	})

	t.Run("无效token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(authMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
		})

		// 创建无效token的请求
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证认证失败
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("有效token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(authMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, 123, userID)

			username, exists := c.Get("username")
			assert.True(t, exists)
			assert.Equal(t, "testuser", username)

			c.JSON(http.StatusOK, gin.H{"message": "protected resource", "user": username})
		})

		// 创建有效token的请求
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证认证成功
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "protected resource")
		assert.Contains(t, w.Body.String(), "testuser")
	})
}

// TestCORSMiddleware 测试CORS中间件
func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestCORSMiddleware")
	defer ctx.Cleanup()

	// 模拟CORS中间件
	corsMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			c.Next()
		}
	}

	t.Run("OPTIONS预检请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(corsMiddleware())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "cors test"})
		})

		// 创建OPTIONS请求
		req, _ := http.NewRequest("OPTIONS", "/api/test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证CORS头部
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("实际API请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(corsMiddleware())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "cors test"})
		})

		// 创建GET请求
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证响应和CORS头部
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Body.String(), "cors test")
	})
}

// TestMiddlewareWithRealServices 测试中间件与真实服务的集成
func TestMiddlewareWithRealServices(t *testing.T) {
	ctx := framework.NewTestContext("TestMiddlewareWithRealServices")
	defer ctx.Cleanup()

	// 创建一个简化的测试套件组件
	container := di.New()
	cfg := ctx.GetConfig()

	gin.SetMode(gin.TestMode)

	t.Run("中间件与DI容器集成", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建一个中间件，从DI容器获取服务
		router.Use(func(c *gin.Context) {
			// 将DI容器设置到上下文中
			c.Set("container", container)
			c.Next()
		})

		router.GET("/di-test", func(c *gin.Context) {
			container, exists := c.Get("container")
			assert.True(t, exists)
			assert.NotNil(t, container)

			c.JSON(http.StatusOK, gin.H{"message": "di integration success"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/di-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "di integration success")
	})

	t.Run("中间件与配置集成", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建使用配置的中间件
		router.Use(func(c *gin.Context) {
			c.Set("config", cfg)
			c.Set("env", "test")
			c.Next()
		})

		router.GET("/config-test", func(c *gin.Context) {
			env, exists := c.Get("env")
			assert.True(t, exists)
			assert.Equal(t, "test", env)

			c.JSON(http.StatusOK, gin.H{"env": env})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/config-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"env":"test"`)
	})
}
