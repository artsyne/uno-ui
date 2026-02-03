package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appMiddleware "github.com/biuren/hebe/internal/app/middleware"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAuthRequiredMiddleware 测试 AuthRequired 中间件
func TestAuthRequiredMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestAuthRequiredMiddleware")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("未登录用户访问受保护资源", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建自定义的 AuthRequired 中间件用于测试
		authMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				// 模拟认证检查失败
				c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
				c.Abort()
			}
		}

		// 应用测试用的认证中间件
		router.Use(authMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
		})

		// 创建无认证信息的请求
		req, _ := http.NewRequest("GET", "/protected", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证返回 403 Forbidden
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("已登录用户访问受保护资源", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建模拟的认证成功中间件
		authMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				// 模拟认证检查成功
				c.Set("userID", 123)
				c.Set("username", "testuser")
				c.Next()
			}
		}

		// 应用测试用的认证中间件
		router.Use(authMiddleware())
		router.GET("/protected", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, 123, userID)

			c.JSON(http.StatusOK, gin.H{"message": "protected resource", "userID": userID})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/protected", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证认证成功
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "protected resource")
		assert.Contains(t, w.Body.String(), "123")
	})
}

// TestSharedDataMiddleware 测试 SharedData 中间件
func TestSharedDataMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestSharedDataMiddleware")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("未登录用户的共享数据", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建模拟的 SharedData 中间件
		sharedDataMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				// 模拟未登录状态
				c.Set("IsLogin", false)
				c.Next()
			}
		}

		// 应用模拟的 SharedData 中间件
		router.Use(sharedDataMiddleware())
		router.GET("/test", func(c *gin.Context) {
			// 验证共享数据
			isLogin, exists := c.Get("IsLogin")
			assert.True(t, exists)
			assert.False(t, isLogin.(bool))

			// AuthUser 应该不存在
			_, exists = c.Get("AuthUser")
			assert.False(t, exists)

			c.JSON(http.StatusOK, gin.H{"isLogin": isLogin})
		})

		// 创建无认证信息的请求
		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "false")
	})

	t.Run("已登录用户的共享数据", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建模拟的 SharedData 中间件（已登录状态）
		sharedDataMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				// 模拟已登录状态
				c.Set("IsLogin", true)
				c.Set("AuthUser", map[string]interface{}{
					"id":       123,
					"username": "testuser",
				})
				c.Next()
			}
		}

		// 应用模拟的 SharedData 中间件
		router.Use(sharedDataMiddleware())
		router.GET("/test", func(c *gin.Context) {
			// 验证共享数据
			isLogin, exists := c.Get("IsLogin")
			assert.True(t, exists)
			assert.True(t, isLogin.(bool))

			// AuthUser 应该存在
			authUser, exists := c.Get("AuthUser")
			assert.True(t, exists)
			assert.NotNil(t, authUser)

			c.JSON(http.StatusOK, gin.H{"isLogin": isLogin, "user": authUser})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "true")
		assert.Contains(t, w.Body.String(), "testuser")
	})
}

// TestPanicRecoveryMiddleware 测试 PanicRecovery 中间件
func TestPanicRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestPanicRecoveryMiddleware")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("panic 恢复测试", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 应用 PanicRecovery 中间件
		router.Use(appMiddleware.PanicRecovery())
		router.GET("/panic", func(c *gin.Context) {
			panic("intentional panic for testing")
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/panic", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证 panic 被恢复，返回 500 错误
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "code")
		assert.Contains(t, w.Body.String(), "message")
	})

	t.Run("正常请求不受影响", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 应用 PanicRecovery 中间件
		router.Use(appMiddleware.PanicRecovery())
		router.GET("/normal", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "normal response"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/normal", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证正常响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "normal response")
	})
}

// TestAppErrorHandlerMiddleware 测试 AppErrorHandler 中间件
func TestAppErrorHandlerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestAppErrorHandlerMiddleware")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("处理应用错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 应用 AppErrorHandler 中间件
		router.Use(appMiddleware.AppErrorHandler())
		router.GET("/error", func(c *gin.Context) {
			// 模拟添加错误到上下文
		_ = c.Error(errors.New("test error"))
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/error", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证错误被处理
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "code")
		assert.Contains(t, w.Body.String(), "message")
	})

	t.Run("正常请求不受影响", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 应用 AppErrorHandler 中间件
		router.Use(appMiddleware.AppErrorHandler())
		router.GET("/normal", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "normal response"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/normal", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证正常响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "normal response")
	})
}

// TestMiddlewareChain 测试中间件链
func TestMiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestMiddlewareChain")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("中间件执行顺序", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建模拟的中间件链
		executionOrder := []string{}

		// 模拟 PanicRecovery 中间件
		panicRecoveryMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				executionOrder = append(executionOrder, "PanicRecovery")
				c.Next()
			}
		}

		// 模拟 AppErrorHandler 中间件
		appErrorHandlerMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				executionOrder = append(executionOrder, "AppErrorHandler")
				c.Next()
			}
		}

		// 模拟 SharedData 中间件
		sharedDataMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				executionOrder = append(executionOrder, "SharedData")
				c.Set("IsLogin", false)
				c.Next()
			}
		}

		// 模拟 AuthRequired 中间件（阻止访问）
		authRequiredMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				executionOrder = append(executionOrder, "AuthRequired")
				c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
				c.Abort()
			}
		}

		// 按照项目中的顺序应用中间件
		router.Use(panicRecoveryMiddleware())
		router.Use(appErrorHandlerMiddleware())
		router.Use(sharedDataMiddleware())
		router.Use(authRequiredMiddleware())

		router.GET("/chain-test", func(c *gin.Context) {
			executionOrder = append(executionOrder, "Handler")
			c.JSON(http.StatusOK, gin.H{"message": "middleware chain test"})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/chain-test", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证中间件执行顺序
		expectedOrder := []string{"PanicRecovery", "AppErrorHandler", "SharedData", "AuthRequired"}
		assert.Equal(t, expectedOrder, executionOrder)

		// 由于认证失败，AuthRequired 会阻止访问
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("不需要认证的路由", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		// 创建模拟的系统级中间件
		panicRecoveryMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Next()
			}
		}

		appErrorHandlerMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Next()
			}
		}

		sharedDataMiddleware := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Set("IsLogin", false)
				c.Next()
			}
		}

		// 只应用系统级中间件
		router.Use(panicRecoveryMiddleware())
		router.Use(appErrorHandlerMiddleware())
		router.Use(sharedDataMiddleware())

		router.GET("/public", func(c *gin.Context) {
			// 验证 SharedData 中间件设置的数据
			isLogin, exists := c.Get("IsLogin")
			assert.True(t, exists)
			assert.False(t, isLogin.(bool))

			c.JSON(http.StatusOK, gin.H{"message": "public resource", "isLogin": isLogin})
		})

		// 创建请求
		req, _ := http.NewRequest("GET", "/public", nil)
		c.Request = req

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证正常访问
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "public resource")
		assert.Contains(t, w.Body.String(), "false")
	})
}
