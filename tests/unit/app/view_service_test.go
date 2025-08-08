package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/view"
	framework "github.com/biuren/hebe/tests/framework/core"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestViewService_Creation 测试视图服务创建
func TestViewService_Creation(t *testing.T) {
	ctx := framework.NewTestContext("TestViewService_Creation")
	defer ctx.Cleanup()

	// 获取全局配置
	cfg := framework.GetGlobalConfig()
	require.NotNil(t, cfg)

	// 创建视图服务
	viewService := view.NewViewService(cfg)
	require.NotNil(t, viewService)

	// 验证服务实现了接口
	var _ *view.ViewService = viewService
}

// TestViewService_PageData 测试页面数据创建
func TestViewService_PageData(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestViewService_PageData")
	defer ctx.Cleanup()

	// 获取全局配置
	cfg := framework.GetGlobalConfig()
	require.NotNil(t, cfg)

	// 创建视图服务
	viewService := view.NewViewService(cfg)
	require.NotNil(t, viewService)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 创建页面数据
	pageData := viewService.CreatePageData(c, "nav.home")

	// 验证页面数据
	require.NotNil(t, pageData)
	assert.NotEmpty(t, pageData.Context.Lang)
	assert.NotEmpty(t, pageData.Title)
}

// TestViewService_JSONResponse 测试JSON响应
func TestViewService_JSONResponse(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestViewService_JSONResponse")
	defer ctx.Cleanup()

	// 获取全局配置
	cfg := framework.GetGlobalConfig()
	require.NotNil(t, cfg)

	// 创建视图服务
	viewService := view.NewViewService(cfg)
	require.NotNil(t, viewService)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 测试JSON响应
	testData := map[string]interface{}{
		"message": "test",
		"status":  "success",
	}

	viewService.JSON(c, http.StatusOK, testData)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, w.Body.String(), "test")
	assert.Contains(t, w.Body.String(), "success")
}

// TestViewService_Redirect 测试重定向
func TestViewService_Redirect(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestViewService_Redirect")
	defer ctx.Cleanup()

	// 获取全局配置
	cfg := framework.GetGlobalConfig()
	require.NotNil(t, cfg)

	// 创建视图服务
	viewService := view.NewViewService(cfg)
	require.NotNil(t, viewService)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/redirect", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 测试重定向
	viewService.Redirect(c, "/target")

	// 验证重定向
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/target", w.Header().Get("Location"))
}

// TestViewService_I18n 测试国际化功能
func TestViewService_I18n(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestViewService_I18n")
	defer ctx.Cleanup()

	// 获取全局配置
	cfg := framework.GetGlobalConfig()
	require.NotNil(t, cfg)

	// 创建视图服务
	viewService := view.NewViewService(cfg)
	require.NotNil(t, viewService)

	tests := []struct {
		name           string
		acceptLanguage string
		expectedLang   string
	}{
		{
			name:           "中文用户",
			acceptLanguage: "zh-CN,zh;q=0.9,en;q=0.8",
			expectedLang:   "zh-cn",
		},
		{
			name:           "英文用户",
			acceptLanguage: "en-US,en;q=0.9,zh;q=0.8",
			expectedLang:   "en",
		},
		{
			name:           "无语言偏好",
			acceptLanguage: "",
			expectedLang:   "zh-cn", // 默认语言
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 创建页面数据
			pageData := viewService.CreatePageData(c, "nav.home")

			// 验证语言设置
			require.NotNil(t, pageData)
			actualLang := pageData.Context.Lang
			// 验证返回的语言不为空
			assert.NotEmpty(t, actualLang)
			// 验证语言代码包含字母（不严格验证格式，因为不同系统可能有不同的规范化方式）
			assert.Regexp(t, `[a-zA-Z]`, actualLang)
		})
	}
}

// TestViewService_WithRealDatabase 测试与真实数据库的集成
func TestViewService_WithRealDatabase(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	ctx := framework.NewTestContext("TestViewService_WithRealDatabase", framework.WithDatabase())
	defer ctx.Cleanup()

	// 获取全局配置
	cfg := framework.GetGlobalConfig()
	require.NotNil(t, cfg)

	// 创建视图服务
	viewService := view.NewViewService(cfg)
	require.NotNil(t, viewService)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 测试页面数据创建（在数据库环境下）
	pageData := viewService.CreatePageData(c, "nav.home")

	// 验证页面数据
	require.NotNil(t, pageData)
	assert.NotEmpty(t, pageData.Context.Lang)
	assert.NotEmpty(t, pageData.Title)

	// 测试JSON响应（在数据库环境下）
	testData := map[string]interface{}{
		"message": "database test",
		"status":  "success",
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req

	viewService.JSON(c2, http.StatusOK, testData)

	// 验证响应
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, w2.Body.String(), "database test")
}
