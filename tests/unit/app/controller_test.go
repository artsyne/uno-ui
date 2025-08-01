package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/app/controller"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 设置测试环境变量
func setupTestEnv() {
	os.Setenv("IDIARY_ENV", "test")
	os.Setenv("IDIARY_GIN_MODE", "debug")
	os.Setenv("IDIARY_DOMAIN", "http://localhost")
	os.Setenv("IDIARY_PORT", "8080")
	os.Setenv("IDIARY_SUMMARY_MAX", "200")
	os.Setenv("IDIARY_SUFFIX", ".test")
	os.Setenv("IDIARY_PAGE_SIZE", "20")
	os.Setenv("IDIARY_MAX_DAYBOOK_LIMIT", "100")
	os.Setenv("IDIARY_HASH_SALT", "test_salt_12345")
	os.Setenv("IDIARY_STATIC_DIR", "static")
	os.Setenv("IDIARY_LOCALES_DIR", "locales")
	os.Setenv("IDIARY_DEFAULT_LANG", "zh-CN")

	// Database配置
	os.Setenv("IDIARY_DB_HOST", "localhost")
	os.Setenv("IDIARY_DB_PORT", "3306")
	os.Setenv("IDIARY_DB_NAME", "test_db")
	os.Setenv("IDIARY_DB_USER", "test_user")
	os.Setenv("IDIARY_DB_PASSWORD", "test_password")
	os.Setenv("IDIARY_DB_CONN", "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local")

	// Redis配置
	os.Setenv("IDIARY_REDIS_HOST", "localhost")
	os.Setenv("IDIARY_REDIS_PORT", "6379")
	os.Setenv("IDIARY_REDIS_PASSWORD", "")

	// Email配置
	os.Setenv("IDIARY_SMTP_SERVER", "smtp.test.com")
	os.Setenv("IDIARY_SMTP_USER", "test@test.com")
	os.Setenv("IDIARY_SMTP_PASSWORD", "test_password")
	os.Setenv("IDIARY_SMTP_PORT", "587")
	os.Setenv("IDIARY_SMTP_FROM", "test@test.com")
	os.Setenv("IDIARY_EMAIL_MAX_RETRIES", "3")
	os.Setenv("IDIARY_EMAIL_BATCH_SIZE", "10")

	// Cookie配置
	os.Setenv("IDIARY_COOKIE_DOMAIN", "localhost")
	os.Setenv("IDIARY_COOKIE_SECRET", "test_cookie_secret_12345")
	os.Setenv("IDIARY_COOKIE_EXPIRE_DAYS", "7")
	os.Setenv("IDIARY_COOKIE_HTTP_ONLY", "true")

	// Upload配置
	os.Setenv("IDIARY_UPLOAD_PATH", "uploads")
	os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")

	// Log配置
	os.Setenv("IDIARY_LOG_LEVEL", "debug")
	os.Setenv("IDIARY_LOG_FORMAT", "json")
	os.Setenv("IDIARY_LOG_PATH", "logs")
	os.Setenv("IDIARY_LOG_FILENAME", "app.log")
	os.Setenv("IDIARY_LOG_MAX_SIZE", "100")
	os.Setenv("IDIARY_LOG_MAX_AGE", "30")
	os.Setenv("IDIARY_LOG_MAX_BACKUPS", "10")
	os.Setenv("IDIARY_LOG_COMPRESS", "true")
	os.Setenv("IDIARY_LOG_LOCAL_TIME", "true")

	// JSON配置
	os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-CN","label":"简体中文","local-text":"简体中文"},{"code":"en-US","label":"English","local-text":"English"}]`)
	os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("IDIARY_PARAMS", `[{"key":"test_param","val":"test_value"}]`)
}

// 清理测试环境变量
func cleanupTestEnv() {
	envVars := []string{
		"IDIARY_ENV", "IDIARY_GIN_MODE", "IDIARY_DOMAIN", "IDIARY_PORT", "IDIARY_SUMMARY_MAX", "IDIARY_SUFFIX", "IDIARY_PAGE_SIZE",
		"IDIARY_MAX_DAYBOOK_LIMIT", "IDIARY_HASH_SALT", "IDIARY_STATIC_DIR", "IDIARY_LOCALES_DIR", "IDIARY_DEFAULT_LANG",
		"IDIARY_DB_HOST", "IDIARY_DB_PORT", "IDIARY_DB_NAME", "IDIARY_DB_USER", "IDIARY_DB_PASSWORD", "IDIARY_DB_CONN",
		"IDIARY_REDIS_HOST", "IDIARY_REDIS_PORT", "IDIARY_REDIS_PASSWORD",
		"IDIARY_SMTP_SERVER", "IDIARY_SMTP_USER", "IDIARY_SMTP_PASSWORD", "IDIARY_SMTP_PORT", "IDIARY_SMTP_FROM",
		"IDIARY_EMAIL_MAX_RETRIES", "IDIARY_EMAIL_BATCH_SIZE",
		"IDIARY_COOKIE_DOMAIN", "IDIARY_COOKIE_SECRET", "IDIARY_COOKIE_EXPIRE_DAYS", "IDIARY_COOKIE_HTTP_ONLY",
		"IDIARY_UPLOAD_PATH", "IDIARY_UPLOAD_TYPES", "IDIARY_UPLOAD_MAX_SIZE",
		"IDIARY_LOG_LEVEL", "IDIARY_LOG_FORMAT", "IDIARY_LOG_PATH", "IDIARY_LOG_FILENAME", "IDIARY_LOG_MAX_SIZE",
		"IDIARY_LOG_MAX_AGE", "IDIARY_LOG_MAX_BACKUPS", "IDIARY_LOG_COMPRESS", "IDIARY_LOG_LOCAL_TIME",
		"IDIARY_SUPPORT_LANGS", "IDIARY_TIMEZONES", "IDIARY_PARAMS",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

func TestUserControllerCreation(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()
	assert.NotNil(t, userController)
}

func TestShowLoginPage(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ShowLoginPage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShowRegisterPage(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/register", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ShowRegisterPage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShowPasswordResetPage(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/password-reset", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ShowPasswordResetPage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShowUserSettingsPage(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/settings", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ShowUserSettingsPage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginPOST(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("login_name", "testuser")
	form.Add("pwd", "testpassword")

	req, _ := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.Login(c)

	// 由于没有真实的数据库连接，期望返回错误状态
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestLoginPOSTEmptyData(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/login", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.Login(c)

	// 空数据应该返回错误
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestRegisterPOST(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("user_name", "newuser")
	form.Add("full_name", "New User")
	form.Add("email", "newuser@example.com")
	form.Add("pwd", "newpassword")
	form.Add("time_zone", "UTC")
	form.Add("lang", "en")

	req, _ := http.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.Register(c)

	// 由于没有真实的数据库连接，期望返回错误状态
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestLogout(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.Logout(c)

	// 退出登录应该成功
	assert.Equal(t, http.StatusFound, w.Code)
}

func TestValidateUsername(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("username", "testuser")

	req, _ := http.NewRequest("POST", "/check-username", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ValidateUsername(c)

	// 验证用户名应该返回JSON响应
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestValidateUsernameEmpty(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/check-username", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ValidateUsername(c)

	// 空用户名应该返回错误
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestValidateEmail(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("email", "test@example.com")

	req, _ := http.NewRequest("POST", "/check-email", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ValidateEmail(c)

	// 验证邮箱应该返回JSON响应
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestValidateEmailEmpty(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/check-email", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ValidateEmail(c)

	// 空邮箱应该返回错误
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestUpdatePasswordEmpty(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/change-password", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.UpdatePassword(c)

	// 空表单数据应该返回错误
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUpdateProfileEmpty(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/update-profile", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.UpdateProfile(c)

	// 空表单数据应该返回错误
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestRouteRegistration(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 注册路由
	userController.RegisterRoutes(router)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/login"},
		{"POST", "/login"},
		{"GET", "/register"},
		{"POST", "/register"},
		{"POST", "/check-username"},
		{"POST", "/check-email"},
		{"GET", "/findpwd"},
		{"POST", "/findpwd"},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(test.method, test.path, nil)
		router.ServeHTTP(w, req)

		// 路由应该存在（不返回404）
		assert.NotEqual(t, http.StatusNotFound, w.Code, "Route %s %s should exist", test.method, test.path)
	}
}

func TestContentTypeHandling(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 注册路由
	router.POST("/check-username", userController.ValidateUsername)
	router.POST("/check-email", userController.ValidateEmail)

	tests := []struct {
		name        string
		path        string
		contentType string
		body        string
	}{
		{"ValidateUsername JSON", "/check-username", "application/json", `{"username":"test"}`},
		{"ValidateUsername Form", "/check-username", "application/x-www-form-urlencoded", "username=test"},
		{"ValidateEmail JSON", "/check-email", "application/json", `{"email":"test@example.com"}`},
		{"ValidateEmail Form", "/check-email", "application/x-www-form-urlencoded", "email=test@example.com"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", test.path, strings.NewReader(test.body))
			req.Header.Set("Content-Type", test.contentType)
			router.ServeHTTP(w, req)

			// 应该能处理不同的内容类型
			assert.NotEqual(t, http.StatusUnsupportedMediaType, w.Code)
		})
	}
}

func TestJSONResponseFormat(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	config.InitConfig()
	// 设置测试数据库连接
	err := database.SetDB("test_user:test_password@tcp(localhost)/test_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Skip("Skipping test due to database connection error:", err)
	}
	
	userController := controller.NewUserController()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("username", "testuser")

	req, _ := http.NewRequest("POST", "/check-username", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	userController.ValidateUsername(c)

	// 验证响应是有效的JSON
	var response map[string]interface{}
	jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, jsonErr, "Response should be valid JSON")

	// 验证响应包含必要字段
	assert.Contains(t, response, "success", "Response should contain success field")
}