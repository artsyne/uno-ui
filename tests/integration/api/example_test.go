package api

import (
	"net/http"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	testhelpers "github.com/biuren/hebe/tests/helpers"
	"github.com/gin-gonic/gin"
)

// TestHealthCheck 测试健康检查端点
func TestHealthCheck(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		// 设置HTTP测试环境
		httpHelper := testhelpers.SetupHTTPTest(suite)
		router := httpHelper.GetRouter()

		// 设置健康检查路由
		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"config": suite.Config.GetTestConfig().Env,
			})
		})

		// 发送GET请求
		w := httpHelper.GET(t, "/health", nil)

		// 验证响应
		httpHelper.AssertStatus(t, w, http.StatusOK)
		data := httpHelper.AssertJSON(t, w)
		httpHelper.AssertJSONField(t, data, "status", "ok")
		httpHelper.AssertJSONField(t, data, "config", "test")
	})
}

// TestAPIWithDI 测试API和DI容器的集成
func TestAPIWithDI(t *testing.T) {
	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		// 设置HTTP测试环境
		httpHelper := testhelpers.SetupHTTPTest(suite)
		router := httpHelper.GetRouter()

		// 在DI容器中注册一个服务
		container := suite.Container
		err := di.Bind(container, func() string {
			return "service-from-di"
		})
		if err != nil {
			t.Fatalf("注册服务到DI容器失败: %v", err)
		}

		// 设置使用DI服务的路由
		router.GET("/service", func(c *gin.Context) {
			// 从DI容器获取服务
			service, err := di.TryGet[string](container)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "服务未找到",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": service,
				"env":     suite.Config.GetTestConfig().Env,
			})
		})

		// 发送GET请求
		w := httpHelper.GET(t, "/service", nil)

		// 验证响应
		httpHelper.AssertStatus(t, w, http.StatusOK)
		data := httpHelper.AssertJSON(t, w)
		httpHelper.AssertJSONField(t, data, "message", "service-from-di")
		httpHelper.AssertJSONField(t, data, "env", "test")
	})
}
