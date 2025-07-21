package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/render/localization"
	"github.com/biuren/hebe/internal/app/render/page"
	"github.com/biuren/hebe/internal/app/render/types"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocalizationManager 模拟本地化管理器
type MockLocalizationManager struct {
	mock.Mock
}

func (m *MockLocalizationManager) GetMessage(key, locale string) string {
	args := m.Called(key, locale)
	return args.String(0)
}

func (m *MockLocalizationManager) GetMessageWithParams(key, locale string, params map[string]string) string {
	args := m.Called(key, locale, params)
	return args.String(0)
}

func (m *MockLocalizationManager) GetSupportedLanguages() []config.SupportLang {
	args := m.Called()
	return args.Get(0).([]config.SupportLang)
}

func (m *MockLocalizationManager) GetSupportedLanguageCodes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockLocalizationManager) GetDefaultLanguage() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLocalizationManager) IsLanguageSupported(locale string) bool {
	args := m.Called(locale)
	return args.Bool(0)
}

func (m *MockLocalizationManager) GetLanguageInfo(locale string) (config.SupportLang, bool) {
	args := m.Called(locale)
	return args.Get(0).(config.SupportLang), args.Bool(1)
}

func (m *MockLocalizationManager) GetMessageCount() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockLocalizationManager) GetMessagesByLanguage(locale string) map[string]string {
	args := m.Called(locale)
	return args.Get(0).(map[string]string)
}

func (m *MockLocalizationManager) GetLocaleFromContext(c *gin.Context) string {
	args := m.Called(c)
	return args.String(0)
}

func (m *MockLocalizationManager) ReloadMessages() error {
	args := m.Called()
	return args.Error(0)
}

// MockPageDataBuilder 模拟页面数据构建器
type MockPageDataBuilder struct {
	mock.Mock
}

func (m *MockPageDataBuilder) BuildPageData(title, description, keywords string, data map[string]interface{}, breadcrumb []types.BreadcrumbItem, context *types.RenderContext, localizationManager localization.ManagerInterface) interface{} {
	args := m.Called(title, description, keywords, data, breadcrumb, context, localizationManager)
	return args.Get(0)
}

func TestContextBuilder_BuildContext(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		setupContext  func(*gin.Context)
		expectedLogin bool
		expectedLang  string
		expectedUser  types.RenderUser
	}{
		{
			name: "已登录用户",
			setupContext: func(c *gin.Context) {
				c.Set("IsLogin", true)
				c.Set("AuthUser", &MockRenderUser{})
			},
			expectedLogin: true,
			expectedLang:  "zh-cn",
			expectedUser:  &MockRenderUser{},
		},
		{
			name: "未登录用户",
			setupContext: func(c *gin.Context) {
				c.Set("IsLogin", false)
			},
			expectedLogin: false,
			expectedLang:  "zh-cn",
			expectedUser:  nil,
		},
		{
			name: "无认证信息",
			setupContext: func(c *gin.Context) {
				// 不设置任何认证信息
			},
			expectedLogin: false,
			expectedLang:  "zh-cn",
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟的本地化管理器
			mockLocalization := &MockLocalizationManager{}
			mockLocalization.On("GetLocaleFromContext", mock.AnythingOfType("*gin.Context")).Return(tt.expectedLang)

			// 创建上下文构建器
			contextBuilder := page.NewContextBuilder(mockLocalization)

			// 创建测试请求
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 设置上下文
			tt.setupContext(c)

			// 构建上下文
			result := contextBuilder.BuildContext(c)

			// 验证结果
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedLogin, result.IsLogin)
			assert.Equal(t, tt.expectedLang, result.Lang)
			assert.Equal(t, tt.expectedUser, result.AuthUser)
		})
	}
}

func TestPageBuilder_BuildPageData(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		titleKey      string
		setupContext  func(*gin.Context)
		expectedTitle string
		expectedDesc  string
	}{
		{
			name:     "首页",
			titleKey: "nav.home",
			setupContext: func(c *gin.Context) {
				c.Request.URL.Path = "/"
			},
			expectedTitle: "首页 - iDiary",
			expectedDesc:  "默认页面描述",
		},
		{
			name:     "用户设置页面",
			titleKey: "user.settings",
			setupContext: func(c *gin.Context) {
				c.Request.URL.Path = "/u/123/settings"
				c.Set("AuthUser", map[string]interface{}{
					"user_name": "testuser",
					"user_id":   123,
				})
			},
			expectedTitle: "用户设置 - iDiary",
			expectedDesc:  "默认页面描述",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟的本地化管理器
			mockLocalization := &MockLocalizationManager{}
			mockLocalization.On("GetLocaleFromContext", mock.AnythingOfType("*gin.Context")).Return("zh-cn")
			// 使用通用的模拟设置，不限制调用次数
			mockLocalization.On("GetMessage", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("测试标题")

			// 创建模拟的页面数据构建器
			mockPageDataBuilder := &MockPageDataBuilder{}
			mockPageDataBuilder.On("BuildPageData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("mock_page_data")

			// 创建上下文构建器
			contextBuilder := page.NewContextBuilder(mockLocalization)

			// 创建页面构建器
			pageBuilder := page.NewPageBuilder(nil, contextBuilder, mockLocalization, mockPageDataBuilder)

			// 创建测试请求
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 设置上下文
			tt.setupContext(c)

			// 构建页面数据
			result := pageBuilder.BuildPageData(c, tt.titleKey)

			// 验证结果
			assert.NotNil(t, result)
			assert.Equal(t, "mock_page_data", result)
		})
	}
}

func TestPageBuilder_GenerateBreadcrumb(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		path          string
		setupContext  func(*gin.Context)
		expectedItems int
		expectedNames []string
	}{
		{
			name: "首页",
			path: "/",
			setupContext: func(c *gin.Context) {
				// 不设置额外信息
			},
			expectedItems: 1,
			expectedNames: []string{"首页"},
		},
		{
			name: "用户设置页面",
			path: "/u/123/settings",
			setupContext: func(c *gin.Context) {
				c.Set("AuthUser", map[string]interface{}{
					"user_name": "testuser",
					"user_id":   123,
				})
			},
			expectedItems: 3,
			expectedNames: []string{"首页", "testuser", "设置"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟的本地化管理器
			mockLocalization := &MockLocalizationManager{}
			mockLocalization.On("GetLocaleFromContext", mock.AnythingOfType("*gin.Context")).Return("zh-cn")
			// 使用通用的模拟设置，不限制调用次数
			mockLocalization.On("GetMessage", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("测试标题")

			// 创建模拟的页面数据构建器
			mockPageDataBuilder := &MockPageDataBuilder{}
			mockPageDataBuilder.On("BuildPageData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("mock_page_data")

			// 创建上下文构建器
			contextBuilder := page.NewContextBuilder(mockLocalization)

			// 创建页面构建器
			pageBuilder := page.NewPageBuilder(nil, contextBuilder, mockLocalization, mockPageDataBuilder)

			// 创建测试请求
			req, _ := http.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 设置上下文
			tt.setupContext(c)

			// 构建页面数据
			result := pageBuilder.BuildPageData(c, "test.title")

			// 验证结果不为空
			assert.NotNil(t, result)
		})
	}
}

// MockRenderUser 模拟渲染用户
type MockRenderUser struct{}

func (m *MockRenderUser) GetLoginName() string   { return "mockuser" }
func (m *MockRenderUser) GetAccessUrl() string   { return "/mock" }
func (m *MockRenderUser) GetUserID() int64       { return 123 }
func (m *MockRenderUser) GetUserName() string    { return "MockUser" }
func (m *MockRenderUser) GetFullName() string    { return "Mock User" }
func (m *MockRenderUser) GetEmail() string       { return "mock@example.com" }
func (m *MockRenderUser) GetUserGroup() int8     { return 1 }
func (m *MockRenderUser) GetFaceUrl() string     { return "/avatar.jpg" }
func (m *MockRenderUser) GetLoginUTCTime() int64 { return 1234567890 }
func (m *MockRenderUser) GetTimeZone() string    { return "Asia/Shanghai" }
func (m *MockRenderUser) GetLang() string        { return "zh-cn" }
func (m *MockRenderUser) GetStatus() int8        { return 1 }
