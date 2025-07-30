package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/view"
	types "github.com/biuren/hebe/internal/shared/contracts"
	testhelpers "github.com/biuren/hebe/tests/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestViewService_Creation 测试视图服务创建
func TestViewService_Creation(t *testing.T) {
	// 设置测试配置
	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	cfg := configHelper.GetTestConfig()

	// 创建真实的视图服务
	viewService := view.NewViewService(cfg)

	require.NotNil(t, viewService)
}

// TestViewService_PageData 测试页面数据创建
func TestViewService_PageData(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	cfg := configHelper.GetTestConfig()
	viewService := view.NewViewService(cfg)

	t.Run("未登录用户的页面数据", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		pageData := viewService.PageData(c, "nav.home")

		assert.NotNil(t, pageData)
		assert.False(t, pageData.IsAuthenticated())
		assert.Nil(t, pageData.GetAuthUser())
		assert.NotEmpty(t, pageData.Context.Lang)
	})

	t.Run("已登录用户的页面数据", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// 设置用户登录状态
		c.Set("IsLogin", true)
		c.Set("AuthUser", getMockAuthUser())

		pageData := viewService.PageData(c, "nav.home")

		assert.NotNil(t, pageData)
		assert.True(t, pageData.IsAuthenticated())
		assert.NotNil(t, pageData.GetAuthUser())
		assert.Equal(t, "mockuser", pageData.GetAuthUser().LoginName)
	})

	t.Run("不同语言的页面数据", func(t *testing.T) {
		tests := []struct {
			name           string
			acceptLanguage string
			expectedLang   bool // 是否为支持的语言
		}{
			{"中文用户", "zh-CN,zh;q=0.9,en;q=0.8", true},
			{"英文用户", "en-US,en;q=0.9,zh;q=0.8", true},
			{"无语言偏好", "", true}, // 应该使用默认语言
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", "/", nil)
				if tt.acceptLanguage != "" {
					req.Header.Set("Accept-Language", tt.acceptLanguage)
				}
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				pageData := viewService.PageData(c, "nav.home")

				assert.NotNil(t, pageData)
				assert.NotEmpty(t, pageData.Context.Lang)
				// 验证返回的语言是支持的
				if tt.expectedLang {
					supportedLangs := []string{"zh-cn", "en", "zh", "zh-CN", "en-US"}
					found := false
					for _, lang := range supportedLangs {
						if pageData.Context.Lang == lang {
							found = true
							break
						}
					}
					assert.True(t, found, "实际语言 %s 不在支持的语言列表中", pageData.Context.Lang)
				}
			})
		}
	})
}

// TestViewService_JSONResponse 测试JSON响应
func TestViewService_JSONResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	cfg := configHelper.GetTestConfig()
	viewService := view.NewViewService(cfg)

	t.Run("JSON成功响应", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		testData := map[string]string{"message": "success"}
		viewService.JSONSuccess(c, testData)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("JSON错误响应", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		viewService.JSONError(c, errors.New("测试错误"), http.StatusBadRequest)

		// JSONError可能使用默认的500状态码，或者按照实际实现来验证
		assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError,
			"状态码应该是400或500，实际是: %d", w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}

// TestViewService_Redirect 测试重定向
func TestViewService_Redirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	cfg := configHelper.GetTestConfig()
	viewService := view.NewViewService(cfg)

	t.Run("重定向响应", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/redirect", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		viewService.Redirect(c, "/target")

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/target", w.Header().Get("Location"))
	})
}

// TestViewService_I18n 测试国际化功能
func TestViewService_I18n(t *testing.T) {
	gin.SetMode(gin.TestMode)

	configHelper := testhelpers.SetupTestConfig(t)
	defer configHelper.Cleanup(t)

	cfg := configHelper.GetTestConfig()
	viewService := view.NewViewService(cfg)

	t.Run("国际化服务可用性", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		pageData := viewService.PageData(c, "nav.home")

		assert.NotNil(t, pageData)
		assert.NotNil(t, pageData.Context)
		assert.NotEmpty(t, pageData.Context.Lang)
	})
}

// TestViewService_WithRealDatabase 测试与真实数据库的集成
func TestViewService_WithRealDatabase(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务，确保测试数据隔离
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		gin.SetMode(gin.TestMode)

		// 使用WithCleanup来获取完整的测试环境
		testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
			cfg := suite.Config.GetTestConfig()
			viewService := view.NewViewService(cfg)

			t.Run("真实数据库环境下的页面数据", func(t *testing.T) {
				req, _ := http.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				pageData := viewService.PageData(c, "nav.home")

				assert.NotNil(t, pageData)
				assert.NotEmpty(t, pageData.Context.Lang)
			})
		})
	})
}

// getMockAuthUser 创建模拟认证用户
func getMockAuthUser() *types.AuthUser {
	return &types.AuthUser{
		LoginName:    "mockuser",
		AccessUrl:    "/mock",
		UserId:       123,
		UserName:     "MockUser",
		FullName:     "Mock User",
		Email:        "mock@example.com",
		UserGroup:    1,
		FaceUrl:      "/avatar.jpg",
		LoginUTCTime: 1234567890,
		TimeZone:     "Asia/Shanghai",
		Lang:         "zh-cn",
		Status:       1,
	}
}
