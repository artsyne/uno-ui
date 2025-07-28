package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试配置
func getTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			DefaultLang: "zh-CN",
			LocalesDir:  "../../../locales",
			StaticDir:   "../../../static",
			PageSize:    20,
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-CN", Label: "简体中文", LocalText: "简体中文"},
			{Code: "en-US", Label: "English", LocalText: "English"},
		},
	}
}

func TestViewService_Creation(t *testing.T) {
	// 测试ViewService的创建
	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	require.NotNil(t, viewService)
	assert.Equal(t, cfg, viewService.GetConfig())
}

func TestViewService_PageData(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	// 测试PageData方法
	pageData := viewService.PageData(c, "test.title")
	require.NotNil(t, pageData)
	assert.NotEmpty(t, pageData.Title)
	assert.NotNil(t, pageData.Data)
}

func TestViewService_BuildPageData(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 测试数据
	testData := map[string]interface{}{
		"message": "hello world",
		"count":   42,
	}

	// 测试BuildPageData方法
	pageData := viewService.BuildPageData(c, "test.title", testData)
	require.NotNil(t, pageData)
	assert.NotEmpty(t, pageData.Title)
	assert.Equal(t, testData, pageData.Data)
}

func TestViewService_BuildRenderContext(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Accept-Language", "en-US")

	// 测试BuildRenderContext方法
	renderContext := viewService.BuildRenderContext(c)
	require.NotNil(t, renderContext)
	assert.NotEmpty(t, renderContext.Lang)
	assert.NotNil(t, renderContext.AuthUser) // 测试AuthUser字段
	assert.False(t, renderContext.IsLogin)   // 默认未登录
}

func TestViewService_JSONResponse(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试JSON响应
	testData := map[string]interface{}{"status": "ok"}
	viewService.JSON(c, http.StatusOK, testData)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestViewService_JSONSuccess(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试成功JSON响应
	testData := map[string]interface{}{"message": "success"}
	viewService.JSONSuccess(c, testData)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestViewService_JSONError(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试错误JSON响应
	testError := assert.AnError
	viewService.JSONError(c, testError)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestViewService_Redirect(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试重定向
	redirectURL := "/login"
	viewService.Redirect(c, redirectURL)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, redirectURL, w.Header().Get("Location"))
}

func TestViewService_I18n(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	cfg := getTestConfig()
	viewService := view.NewViewService(cfg)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Accept-Language", "zh-CN")

	// 测试获取语言设置
	lang := viewService.GetLocaleFromContext(c)
	assert.NotEmpty(t, lang)

	// 测试获取国际化消息（这里只测试方法调用，不测试具体内容）
	message := viewService.GetMessage("test.key", "zh-CN")
	assert.NotNil(t, message) // 即使没有找到消息，也应该返回非nil值
}