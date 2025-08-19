package view

import (
	"testing"

	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockI18nManager 模拟国际化管理器
type MockI18nManager struct{}

func (m *MockI18nManager) GetMessage(key, locale string) string {
	return key
}

func (m *MockI18nManager) GetMessageWithParams(key, locale string, params map[string]string) string {
	return key
}

func (m *MockI18nManager) GetSupportedLanguages() []config.SupportLang {
	return []config.SupportLang{
		{Code: "en", Label: "English", LocalText: "English"},
		{Code: "zh-cn", Label: "中文", LocalText: "中文"},
	}
}

func (m *MockI18nManager) GetSupportedLanguageCodes() []string {
	return []string{"en", "zh-cn"}
}

func (m *MockI18nManager) GetDefaultLanguage() string {
	return "en"
}

func (m *MockI18nManager) IsLanguageSupported(locale string) bool {
	return locale == "en" || locale == "zh-cn"
}

func (m *MockI18nManager) GetLanguageInfo(locale string) (config.SupportLang, bool) {
	for _, lang := range m.GetSupportedLanguages() {
		if lang.Code == locale {
			return lang, true
		}
	}
	return config.SupportLang{}, false
}

func (m *MockI18nManager) GetMessageCount() int {
	return 100
}

func (m *MockI18nManager) GetMessagesByLanguage(locale string) map[string]string {
	return map[string]string{"test.key": "test value"}
}

func (m *MockI18nManager) GetLocaleFromContext(c *gin.Context) string {
	return "en"
}

func (m *MockI18nManager) ReloadMessages() error {
	return nil
}

// TestNewViewService 测试创建新的视图服务
func TestNewViewService(t *testing.T) {
	cfg := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode: "test",
			Domain: "http://localhost",
			Port: "8080",
		},
	}
	mockI18n := &MockI18nManager{}

	service := view.NewViewService(cfg, mockI18n)
	assert.NotNil(t, service)
}

// TestNewViewService_NilDependencies 测试空依赖
func TestNewViewService_NilDependencies(t *testing.T) {
	tests := []struct {
		name string
		config *config.Config
		i18nManager i18n.I18nManager
		expectPanic bool
	}{
		{
			name: "valid config and i18n",
			config: &config.Config{
				Env: "test",
				App: config.AppConfig{GinMode: "test"},
			},
			i18nManager: &MockI18nManager{},
			expectPanic: false,
		},
		{
			name: "nil config",
			config: nil,
			i18nManager: &MockI18nManager{},
			expectPanic: false, // ViewService 可能允许 nil config
		},
		{
			name: "nil i18n manager",
			config: &config.Config{
				Env: "test",
				App: config.AppConfig{GinMode: "test"},
			},
			i18nManager: nil,
			expectPanic: false, // ViewService 可能允许 nil i18nManager
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					view.NewViewService(tt.config, tt.i18nManager)
				})
			} else {
				service := view.NewViewService(tt.config, tt.i18nManager)
				assert.NotNil(t, service)
			}
		})
	}
}

// TestViewService_CreatePageData 测试创建页面数据
func TestViewService_CreatePageData(t *testing.T) {
	cfg := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode: "test",
			Domain: "http://localhost",
			Port: "8080",
		},
	}
	mockI18n := &MockI18nManager{}

	service := view.NewViewService(cfg, mockI18n)
	assert.NotNil(t, service)

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 测试创建页面数据
	pageData := service.CreatePageData(c, "Test Title")
	assert.NotNil(t, pageData)
	assert.Equal(t, "Test Title - iDiary", pageData.Title)
	assert.NotNil(t, pageData.Context)
	assert.Equal(t, "en", pageData.Context.Lang)
	assert.NotNil(t, pageData.Data)
}

// TestViewService_CreatePageDataWithData 测试创建带数据的页面数据
func TestViewService_CreatePageDataWithData(t *testing.T) {
	cfg := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode: "test",
			Domain: "http://localhost",
			Port: "8080",
		},
	}
	mockI18n := &MockI18nManager{}

	service := view.NewViewService(cfg, mockI18n)
	assert.NotNil(t, service)

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 测试数据
	testData := map[string]interface{}{
		"message": "Hello World",
		"count": 42,
	}

	// 测试创建带数据的页面数据
	pageData := service.CreatePageDataWithData(c, "Test Title", testData)
	assert.NotNil(t, pageData)
	assert.Equal(t, "Test Title - iDiary", pageData.Title)
	assert.NotNil(t, pageData.Context)
	assert.Equal(t, "en", pageData.Context.Lang)
	assert.Equal(t, testData, pageData.Data)
}

// TestViewService_BuildRenderContext 测试构建渲染上下文
func TestViewService_BuildRenderContext(t *testing.T) {
	cfg := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode: "test",
			Domain: "http://localhost",
			Port: "8080",
		},
	}
	mockI18n := &MockI18nManager{}

	service := view.NewViewService(cfg, mockI18n)
	assert.NotNil(t, service)

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 测试构建渲染上下文
	renderContext := service.BuildRenderContext(c)
	assert.NotNil(t, renderContext)
	assert.Equal(t, "en", renderContext.Lang)
	assert.False(t, renderContext.IsLogin) // 默认未登录
}



// TestViewService_EmptyTitle 测试空标题
func TestViewService_EmptyTitle(t *testing.T) {
	cfg := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode: "test",
			Domain: "http://localhost",
			Port: "8080",
		},
	}
	mockI18n := &MockI18nManager{}

	service := view.NewViewService(cfg, mockI18n)
	assert.NotNil(t, service)

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 测试空标题
	pageData := service.CreatePageData(c, "")
	assert.NotNil(t, pageData)
	assert.Equal(t, " - iDiary", pageData.Title)
	assert.NotNil(t, pageData.Context)
}

// TestViewService_NilData 测试空数据
func TestViewService_NilData(t *testing.T) {
	cfg := &config.Config{
		Env: "test",
		App: config.AppConfig{
			GinMode: "test",
			Domain: "http://localhost",
			Port: "8080",
		},
	}
	mockI18n := &MockI18nManager{}

	service := view.NewViewService(cfg, mockI18n)
	assert.NotNil(t, service)

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 测试空数据
	pageData := service.CreatePageDataWithData(c, "Test Title", nil)
	assert.NotNil(t, pageData)
	assert.Equal(t, "Test Title - iDiary", pageData.Title)
	assert.Nil(t, pageData.Data)
	assert.NotNil(t, pageData.Context)
}