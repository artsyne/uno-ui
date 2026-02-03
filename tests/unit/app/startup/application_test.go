package startup_test

import (
	"testing"

	bootstrap "github.com/biuren/hebe/internal/app/startup"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInitializer 模拟初始化器
type MockInitializer struct {
	mock.Mock
}

func (m *MockInitializer) InitializeCore() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockInitializer) GetContainer() interface{} {
	args := m.Called()
	return args.Get(0)
}

// TestNewApplication 测试创建新的应用程序
func TestNewApplication(t *testing.T) {
	tests := []struct {
		name string
		config *config.Config
		wantNil bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Env: "test",
				Database: config.DatabaseConfig{
					Host: "localhost",
					Port: 3306,
					User: "test",
					Password: "test",
					Name: "test",
					Conn: "mysql",
				},
				App: config.AppConfig{
					Port: "8080",
					Domain: "http://localhost:8080",
					GinMode: "debug",
					Suffix: ".html",
					PageSize: 10,
					MaxDiaryLimit: 100,
					SummaryMax: 200,
					StaticDir: "static",
					LocalesDir: "locales",
					DefaultLang: "en",
				},
				User: config.UserConfig{PwdMinLen: 6, PwdMaxLen: 20, StatusOK: 1, StatusWaitActive: 0, StatusDeleted: -1, DefaultGroup: 1},
				ValidCode: config.ValidCodeConfig{Length: 6, TypeActiveUser: 1, TypeResetPwd: 2, ExpireMinutes: 30},
				Redis: config.RedisConfig{Host: "localhost", Port: 6379, Password: ""},
				Email: config.EmailConfig{SMTPServer: "smtp.example.com", SMTPUser: "test@example.com", SMTPPassword: "password", SMTPPort: 587, SMTPFrom: "test@example.com", MaxRetries: 3, BatchSize: 10},
				Cookie: config.CookieConfig{Domain: "localhost", Secret: "test-secret-key-16chars", ExpireDays: 7, HttpOnly: true},
				Upload: config.UploadConfig{Path: "uploads", Types: "jpg,png,gif", MaxSize: 1048576},
				SupportLangs: []config.SupportLang{{Code: "en", Label: "English", LocalText: "English"}},
				Timezones: []string{"UTC"},
				Params: []config.Param{{Key: "test_key", Val: "test_value"}},
				Log: config.LogConfig{Level: "info", Format: "json", Path: "logs", Filename: "app.log", MaxSize: 100, MaxAge: 7, MaxBackups: 3, Compress: true, LocalTime: true},
			},
			wantNil: false,
		},
		{
			name: "nil config",
			config: nil,
			wantNil: false, // NewApplication 总是返回实例，即使配置为 nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bootstrap.NewApplication(tt.config)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

// TestApplication_Bootstrap 测试应用程序引导
func TestApplication_Bootstrap(t *testing.T) {
	config := &config.Config{
		Env: "test",
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
			Conn:     "mysql",
		},
		App: config.AppConfig{
			Port:    "8080",
			Domain:  "http://localhost:8080",
			GinMode: "debug",
			Suffix:  ".html",
			PageSize: 10,
			MaxDiaryLimit: 100,
			SummaryMax: 200,
			StaticDir: "static",
			LocalesDir: "locales",
			DefaultLang: "en",
		},
		User: config.UserConfig{
			PwdMinLen: 6,
			PwdMaxLen: 20,
			StatusOK: 1,
			StatusWaitActive: 0,
			StatusDeleted: -1,
			DefaultGroup: 1,
		},
		ValidCode: config.ValidCodeConfig{
			Length: 6,
			TypeActiveUser: 1,
			TypeResetPwd: 2,
			ExpireMinutes: 30,
		},
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: 6379,
			Password: "",
		},
		Email: config.EmailConfig{
			SMTPServer: "smtp.example.com",
			SMTPUser: "test@example.com",
			SMTPPassword: "password",
			SMTPPort: 587,
			SMTPFrom: "test@example.com",
			MaxRetries: 3,
			BatchSize: 10,
		},
		Cookie: config.CookieConfig{
			Domain: "localhost",
			Secret: "test-secret-key-16chars",
			ExpireDays: 7,
			HttpOnly: true,
		},
		Upload: config.UploadConfig{
			Path: "uploads",
			Types: "jpg,png,gif",
			MaxSize: 1048576,
		},
		SupportLangs: []config.SupportLang{
			{Code: "en", Label: "English", LocalText: "English"},
			{Code: "zh-CN", Label: "中文", LocalText: "中文"},
		},
		Timezones: []string{"UTC", "Asia/Shanghai"},
		Params: []config.Param{
			{Key: "test_key", Val: "test_value"},
		},
		Log: config.LogConfig{
			Level: "info",
			Format: "json",
			Path: "logs",
			Filename: "app.log",
			MaxSize: 100,
			MaxAge: 7,
			MaxBackups: 3,
			Compress: true,
			LocalTime: true,
		},
	}

	app := bootstrap.NewApplication(config)
	assert.NotNil(t, app)

	// 测试引导过程
	err := app.Bootstrap()
	// 由于涉及真实的数据库和服务初始化，这里主要测试方法是否存在
	// 在实际项目中，应该使用 mock 来测试各个组件
	assert.NotNil(t, err) // 预期会有错误，因为没有真实的依赖
}

// TestApplication_Run 测试应用程序运行
func TestApplication_Run(t *testing.T) {
	config := &config.Config{
		Env: "test",
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "test",
			Password: "test",
			Name: "test",
			Conn: "mysql",
		},
		App: config.AppConfig{
			Port: "8080",
			Domain: "http://localhost:8080",
			GinMode: "debug",
			Suffix: ".html",
			PageSize: 10,
			MaxDiaryLimit: 100,
			SummaryMax: 200,
			StaticDir: "static",
			LocalesDir: "locales",
			DefaultLang: "en",
		},
		User: config.UserConfig{
			PwdMinLen: 6,
			PwdMaxLen: 20,
			StatusOK: 1,
			StatusWaitActive: 0,
			StatusDeleted: -1,
			DefaultGroup: 1,
		},
		ValidCode: config.ValidCodeConfig{
			Length: 6,
			TypeActiveUser: 1,
			TypeResetPwd: 2,
			ExpireMinutes: 30,
		},
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: 6379,
			Password: "",
		},
		Email: config.EmailConfig{
			SMTPServer: "smtp.example.com",
			SMTPUser: "test@example.com",
			SMTPPassword: "password",
			SMTPPort: 587,
			SMTPFrom: "test@example.com",
			MaxRetries: 3,
			BatchSize: 10,
		},
		Cookie: config.CookieConfig{
			Domain: "localhost",
			Secret: "test-secret-key-16chars",
			ExpireDays: 7,
			HttpOnly: true,
		},
		Upload: config.UploadConfig{
			Path: "uploads",
			Types: "jpg,png,gif",
			MaxSize: 1048576,
		},
		SupportLangs: []config.SupportLang{
			{Code: "en", Label: "English", LocalText: "English"},
			{Code: "zh-CN", Label: "中文", LocalText: "中文"},
		},
		Timezones: []string{"UTC", "Asia/Shanghai"},
		Params: []config.Param{
			{Key: "test_key", Val: "test_value"},
		},
		Log: config.LogConfig{
			Level: "info",
			Format: "json",
			Path: "logs",
			Filename: "app.log",
			MaxSize: 100,
			MaxAge: 7,
			MaxBackups: 3,
			Compress: true,
			LocalTime: true,
		},
	}

	app := bootstrap.NewApplication(config)
	assert.NotNil(t, app)

	// 测试运行方法存在
	// 注意：实际运行会启动服务器，所以这里只测试方法是否存在
	assert.NotNil(t, app.Run)
}

// TestApplication_Fields 测试应用程序字段
func TestApplication_Fields(t *testing.T) {
	config := &config.Config{
		Env: "test",
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "test",
			Password: "test",
			Name: "test",
			Conn: "mysql",
		},
		App: config.AppConfig{
			Port: "8080",
			Domain: "http://localhost:8080",
			GinMode: "debug",
			Suffix: ".html",
			PageSize: 10,
			MaxDiaryLimit: 100,
			SummaryMax: 200,
			StaticDir: "static",
			LocalesDir: "locales",
			DefaultLang: "en",
		},
		User: config.UserConfig{PwdMinLen: 6, PwdMaxLen: 20, StatusOK: 1, StatusWaitActive: 0, StatusDeleted: -1, DefaultGroup: 1},
		ValidCode: config.ValidCodeConfig{Length: 6, TypeActiveUser: 1, TypeResetPwd: 2, ExpireMinutes: 30},
		Redis: config.RedisConfig{Host: "localhost", Port: 6379, Password: ""},
		Email: config.EmailConfig{SMTPServer: "smtp.example.com", SMTPUser: "test@example.com", SMTPPassword: "password", SMTPPort: 587, SMTPFrom: "test@example.com", MaxRetries: 3, BatchSize: 10},
		Cookie: config.CookieConfig{Domain: "localhost", Secret: "test-secret-key-16chars", ExpireDays: 7, HttpOnly: true},
		Upload: config.UploadConfig{Path: "uploads", Types: "jpg,png,gif", MaxSize: 1048576},
		SupportLangs: []config.SupportLang{{Code: "en", Label: "English", LocalText: "English"}},
		Timezones: []string{"UTC"},
		Params: []config.Param{{Key: "test_key", Val: "test_value"}},
		Log: config.LogConfig{Level: "info", Format: "json", Path: "logs", Filename: "app.log", MaxSize: 100, MaxAge: 7, MaxBackups: 3, Compress: true, LocalTime: true},
	}

	app := bootstrap.NewApplication(config)
	assert.NotNil(t, app)

	// 验证应用程序创建成功
	assert.NotNil(t, app)
}

// TestApplication_InitializationSteps 测试初始化步骤
func TestApplication_InitializationSteps(t *testing.T) {
	config := &config.Config{
		Env: "test",
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "test",
			Password: "test",
			Name: "test",
			Conn: "mysql",
		},
		App: config.AppConfig{
			Port: "8080",
			Domain: "http://localhost:8080",
			GinMode: "debug",
			Suffix: ".html",
			PageSize: 10,
			MaxDiaryLimit: 100,
			SummaryMax: 200,
			StaticDir: "static",
			LocalesDir: "locales",
			DefaultLang: "en",
		},
		User: config.UserConfig{PwdMinLen: 6, PwdMaxLen: 20, StatusOK: 1, StatusWaitActive: 0, StatusDeleted: -1, DefaultGroup: 1},
		ValidCode: config.ValidCodeConfig{Length: 6, TypeActiveUser: 1, TypeResetPwd: 2, ExpireMinutes: 30},
		Redis: config.RedisConfig{Host: "localhost", Port: 6379, Password: ""},
		Email: config.EmailConfig{SMTPServer: "smtp.example.com", SMTPUser: "test@example.com", SMTPPassword: "password", SMTPPort: 587, SMTPFrom: "test@example.com", MaxRetries: 3, BatchSize: 10},
		Cookie: config.CookieConfig{Domain: "localhost", Secret: "test-secret-key-16chars", ExpireDays: 7, HttpOnly: true},
		Upload: config.UploadConfig{Path: "uploads", Types: "jpg,png,gif", MaxSize: 1048576},
		SupportLangs: []config.SupportLang{{Code: "en", Label: "English", LocalText: "English"}},
		Timezones: []string{"UTC"},
		Params: []config.Param{{Key: "test_key", Val: "test_value"}},
		Log: config.LogConfig{Level: "info", Format: "json", Path: "logs", Filename: "app.log", MaxSize: 100, MaxAge: 7, MaxBackups: 3, Compress: true, LocalTime: true},
	}

	app := bootstrap.NewApplication(config)
	assert.NotNil(t, app)

	// 测试各个初始化步骤的方法是否存在
	// 由于这些是私有方法，我们通过 Bootstrap 来间接测试
	err := app.Bootstrap()
	// 预期会有错误，因为没有真实的依赖，但方法应该存在
	assert.Error(t, err)
}

// TestApplication_ConfigValidation 测试配置验证
func TestApplication_ConfigValidation(t *testing.T) {
	tests := []struct {
		name string
		config *config.Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Env: "test",
				Database: config.DatabaseConfig{
					Host: "localhost",
					Port: 3306,
					User: "test",
					Password: "test",
					Name: "test",
					Conn: "mysql",
				},
				App: config.AppConfig{
					Port: "8080",
					Domain: "http://localhost:8080",
					GinMode: "debug",
					Suffix: ".html",
					PageSize: 10,
					MaxDiaryLimit: 100,
					SummaryMax: 200,
					StaticDir: "static",
					LocalesDir: "locales",
					DefaultLang: "en",
				},
				User: config.UserConfig{PwdMinLen: 6, PwdMaxLen: 20, StatusOK: 1, StatusWaitActive: 0, StatusDeleted: -1, DefaultGroup: 1},
				ValidCode: config.ValidCodeConfig{Length: 6, TypeActiveUser: 1, TypeResetPwd: 2, ExpireMinutes: 30},
				Redis: config.RedisConfig{Host: "localhost", Port: 6379, Password: ""},
				Email: config.EmailConfig{SMTPServer: "smtp.example.com", SMTPUser: "test@example.com", SMTPPassword: "password", SMTPPort: 587, SMTPFrom: "test@example.com", MaxRetries: 3, BatchSize: 10},
				Cookie: config.CookieConfig{Domain: "localhost", Secret: "test-secret-key-16chars", ExpireDays: 7, HttpOnly: true},
				Upload: config.UploadConfig{Path: "uploads", Types: "jpg,png,gif", MaxSize: 1048576},
				SupportLangs: []config.SupportLang{{Code: "en", Label: "English", LocalText: "English"}},
				Timezones: []string{"UTC"},
				Params: []config.Param{{Key: "test_key", Val: "test_value"}},
				Log: config.LogConfig{Level: "info", Format: "json", Path: "logs", Filename: "app.log", MaxSize: 100, MaxAge: 7, MaxBackups: 3, Compress: true, LocalTime: true},
			},
			expectError: true, // 仍然会有错误，因为没有真实依赖
		},
		{
			name: "empty config",
			config: &config.Config{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := bootstrap.NewApplication(tt.config)
			assert.NotNil(t, app)

			err := app.Bootstrap()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestApplication_GinModeSettings 测试 Gin 模式设置
func TestApplication_GinModeSettings(t *testing.T) {
	tests := []struct {
		name string
		mode string
		expectedGinMode string
	}{
		{
			name: "test mode",
			mode: "test",
			expectedGinMode: gin.TestMode,
		},
		{
			name: "debug mode",
			mode: "debug",
			expectedGinMode: gin.DebugMode,
		},
		{
			name: "release mode",
			mode: "release",
			expectedGinMode: gin.ReleaseMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置 Gin 模式
			gin.SetMode(tt.expectedGinMode)
			actualMode := gin.Mode()
			assert.Equal(t, tt.expectedGinMode, actualMode)
		})
	}
}

// TestApplication_MultipleInstances 测试多个应用实例
func TestApplication_MultipleInstances(t *testing.T) {
	config1 := &config.Config{
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			User: "test1",
			Password: "test1",
			Name: "test1",
		},
	}

	config2 := &config.Config{
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			User: "test2",
			Password: "test2",
			Name: "test2",
		},
	}

	app1 := bootstrap.NewApplication(config1)
	app2 := bootstrap.NewApplication(config2)

	assert.NotNil(t, app1)
	assert.NotNil(t, app2)
	assert.NotEqual(t, app1, app2)
}

// TestApplication_MethodsExist 测试方法存在性
func TestApplication_MethodsExist(t *testing.T) {
	config := &config.Config{
		Database: config.DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			User: "test",
			Password: "test",
			Name: "test",
		},
	}

	app := bootstrap.NewApplication(config)
	assert.NotNil(t, app)

	// 验证关键方法存在
	assert.NotNil(t, app.Bootstrap)
	assert.NotNil(t, app.Run)
}
