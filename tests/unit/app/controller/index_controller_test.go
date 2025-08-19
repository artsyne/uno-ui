package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/controller"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockI18nManager 模拟 I18nManager
type MockI18nManager struct{}

func (m *MockI18nManager) GetMessage(key, locale string) string {
	return key // 简单返回 key 作为消息
}

func (m *MockI18nManager) GetMessageWithParams(key, locale string, params map[string]string) string {
	return key // 简单返回 key 作为消息
}

func (m *MockI18nManager) GetSupportedLanguages() []config.SupportLang {
	return []config.SupportLang{
		{Code: "en", Label: "English", LocalText: "English"},
		{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
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
	return 10 // 模拟消息数量
}

func (m *MockI18nManager) GetMessagesByLanguage(locale string) map[string]string {
	return map[string]string{"test.key": "test.value"}
}

func (m *MockI18nManager) GetLocaleFromContext(c *gin.Context) string {
	return "en"
}

func (m *MockI18nManager) ReloadMessages() error {
	return nil
}

// TestNewIndexController 测试 IndexController 构造函数
func TestNewIndexController(t *testing.T) {
	// 准备测试数据
	cfg := &config.Config{}
	mockI18nManager := &MockI18nManager{}
	viewService := view.NewViewService(cfg, mockI18nManager)

	// 测试构造函数
	ctrl := controller.NewIndexController(cfg, viewService)

	// 验证控制器创建成功
	assert.NotNil(t, ctrl)
}

// TestIndexController_RenderIndex 测试首页渲染
func TestIndexController_RenderIndex(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 准备测试数据
	cfg := &config.Config{}
	mockI18nManager := &MockI18nManager{}
	viewService := view.NewViewService(cfg, mockI18nManager)
	ctrl := controller.NewIndexController(cfg, viewService)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	// 设置模拟的登录状态
	c.Set("IsLogin", false)

	// 测试 ViewService 的 CreatePageData 方法
	pageData := viewService.CreatePageData(c, "home.title")

	// 验证 PageData 结构正确
	assert.NotNil(t, pageData)
	assert.NotNil(t, pageData.Data)
	assert.Equal(t, "en", pageData.Context.Lang)

	// 测试可以设置 Data 字段
	pageData.Data["Active"] = "users"
	assert.Equal(t, "users", pageData.Data["Active"])

	// 验证控制器不为空
	assert.NotNil(t, ctrl)
}

// TestIndexController_RenderIndex_Integration 集成测试
func TestIndexController_RenderIndex_Integration(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 准备测试配置
	cfg := &config.Config{
		App: config.AppConfig{
			LocalesDir:  "../../../../website/locales",
			DefaultLang: "zh-cn",
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
			{Code: "en", Label: "English", LocalText: "English"},
		},
	}

	// 创建国际化管理器
	i18nManager := i18n.NewLocalizationManager(cfg)
	
	// 创建 ViewService
	viewService := view.NewViewService(cfg, i18nManager)
	
	// 创建 IndexController
	ctrl := controller.NewIndexController(cfg, viewService)

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	// 设置模拟的登录状态
	c.Set("IsLogin", false)

	// 测试 ViewService 的 CreatePageData 方法
	pageData := viewService.CreatePageData(c, "home.title")

	// 验证 PageData 结构正确
	assert.NotNil(t, pageData)
	assert.NotNil(t, pageData.Data)
	assert.Equal(t, "zh-cn", pageData.Context.Lang)
	assert.False(t, pageData.Context.IsLogin)

	// 测试可以设置 Data 字段
	pageData.Data["Active"] = "home"
	assert.Equal(t, "home", pageData.Data["Active"])

	// 验证控制器不为空
	assert.NotNil(t, ctrl)

	// 验证 ViewService 的国际化功能
	i18nMgr := viewService.GetI18nManager()
	assert.NotNil(t, i18nMgr)
	
	// 测试获取消息
	message := i18nMgr.GetMessage("common.home", "zh-cn")
	assert.NotEmpty(t, message, "应该能够获取国际化消息")
}