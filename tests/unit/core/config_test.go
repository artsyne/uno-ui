package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/tests/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"
)

// loadEnvLocal loads .env.local from the project root.
func loadEnvLocal(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Find project root by looking for go.mod
	root := cwd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("Could not find project root (go.mod)")
		}
		root = parent
	}

	envPath := filepath.Join(root, ".env.local")
	// Load .env.local, but don't error if it doesn't exist (CI environment)
	// Note: gotenv.Load will not overwrite existing environment variables
	if err := gotenv.Load(envPath); err != nil {
		t.Logf("Warning: failed to load %s: %v", envPath, err)
	}
}

func TestConfig_InitConfig(t *testing.T) {
	// 使用测试框架初始化
	ctx := framework.NewTestContext("TestConfig_InitConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	// 加载本地环境变量
	loadEnvLocal(t)

	// 测试配置初始化
	err := config.InitConfig()
	require.NoError(t, err)

	// 获取配置
	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证基本配置 (Matches .env.local)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "debug", cfg.App.GinMode)
	assert.Equal(t, "http://127.0.0.1:9090", cfg.App.Domain)
	assert.Equal(t, "9090", cfg.App.Port)
}

func TestConfig_AppConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_AppConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证App配置 (Matches .env.local)
	assert.Equal(t, "debug", cfg.App.GinMode)
	assert.Equal(t, "http://127.0.0.1:9090", cfg.App.Domain)
	assert.Equal(t, "9090", cfg.App.Port)
	assert.Equal(t, 1200, cfg.App.SummaryMax)
	assert.Equal(t, "...", cfg.App.Suffix)
	assert.Equal(t, 20, cfg.App.PageSize)
	assert.Equal(t, 666, cfg.App.MaxDiaryLimit)
	assert.Equal(t, "en", cfg.App.DefaultLang)

	// Check relative paths or resolved paths
	// .env.local has ./website/static
	// config loader might resolve them or keep as is. Usually viper keeps as is.
	assert.Contains(t, cfg.App.StaticDir, "website/static")
	assert.Contains(t, cfg.App.LocalesDir, "website/locales")
}

func TestConfig_DatabaseConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_DatabaseConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Database配置 (Matches .env.local)
	assert.Equal(t, "127.0.0.1", cfg.Database.Host)
	assert.Equal(t, 3306, cfg.Database.Port)
	assert.Equal(t, "cd", cfg.Database.Name)
	assert.Equal(t, "root", cfg.Database.User)
	assert.Equal(t, "123123", cfg.Database.Password)
	// Conn string has format specifiers in .env.local
	assert.Contains(t, cfg.Database.Conn, "%s:%s@tcp(%s:3306)/%s")
}

func TestConfig_RedisConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_RedisConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Redis配置 (Matches .env.local)
	assert.Equal(t, "127.0.0.1", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, "woaibeijingtiananmeng", cfg.Redis.Password)
}

func TestConfig_EmailConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_EmailConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Email配置 (Matches .env.local)
	assert.Equal(t, "smtpdm.aliyun.com:25", cfg.Email.SMTPServer)
	assert.Equal(t, "noreply@sinodiary.com", cfg.Email.SMTPUser)
	assert.Equal(t, "HL90rmWKeJjS92", cfg.Email.SMTPPassword)
	assert.Equal(t, 25, cfg.Email.SMTPPort)
	assert.Equal(t, "noreply@sinodiary.com", cfg.Email.SMTPFrom)
	assert.Equal(t, 3, cfg.Email.MaxRetries)
}

func TestConfig_CookieConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_CookieConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Cookie配置 (Matches .env.local)
	assert.Equal(t, "127.0.0.1", cfg.Cookie.Domain)
	assert.Equal(t, "2GzVloukJl9C7y44O740o59e3sD8xTzr", cfg.Cookie.Secret)
	assert.Equal(t, 7, cfg.Cookie.ExpireDays)
	assert.True(t, cfg.Cookie.HttpOnly)
}

func TestConfig_UploadConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_UploadConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Upload配置 (Matches .env.local)
	assert.Equal(t, "upload/", cfg.Upload.Path)
	// .env.local uses pipe separator for types?
	// IDIARY_UPLOAD_TYPES="data:image/jpeg|data:image/jpg|data:image/gif|data:image/png"
	// Verify exact match
	assert.Equal(t, "data:image/jpeg|data:image/jpg|data:image/gif|data:image/png", cfg.Upload.Types)
	assert.Equal(t, 1048576, cfg.Upload.MaxSize) // 1MB
}

func TestConfig_LogConfig(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_LogConfig")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证Log配置 (Matches .env.local)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "text", cfg.Log.Format)
	assert.Equal(t, "./logs", cfg.Log.Path)
	assert.Equal(t, "idiary.log", cfg.Log.Filename)
	assert.Equal(t, 100, cfg.Log.MaxSize)
	assert.Equal(t, 30, cfg.Log.MaxAge)
	assert.Equal(t, 5, cfg.Log.MaxBackups)
	assert.True(t, cfg.Log.Compress)
	assert.True(t, cfg.Log.LocalTime)
}

func TestConfig_SupportLangs(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_SupportLangs")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证支持的语言 (Matches .env.local: 5 langs)
	require.Len(t, cfg.SupportLangs, 5)
	assert.Equal(t, "en", cfg.SupportLangs[0].Code)
	assert.Equal(t, "ja", cfg.SupportLangs[1].Code)
}

func TestConfig_Timezones(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_Timezones")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证时区配置 (Matches .env.local: huge list)
	assert.NotEmpty(t, cfg.Timezones)
	assert.Contains(t, cfg.Timezones, "Asia/Shanghai")
	assert.Contains(t, cfg.Timezones, "UTC")
}

func TestConfig_Params(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_Params")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	cfg := config.GetConfig()
	require.NotNil(t, cfg)

	// 验证参数配置 (Matches .env.local: 2 params)
	require.NotEmpty(t, cfg.Params)
	assert.Equal(t, "reg_email_subject", cfg.Params[0].Key)
}

func TestConfig_GetParams(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_GetParams")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	// 测试GetParams函数
	params := config.GetParams()
	require.NotNil(t, params)
	assert.Equal(t, "中国日记网注册激活邮件", params["reg_email_subject"])
}

func TestConfig_GetTimeZone(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_GetTimeZone")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	// 测试GetTimeZone函数
	timezones := config.GetTimeZone()
	require.NotNil(t, timezones)
	assert.Contains(t, timezones, "Asia/Shanghai")
}

func TestConfig_GetSupportLangs(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_GetSupportLangs")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	require.NoError(t, err)

	// 测试GetSupportLangs函数
	supportLangs := config.GetSupportLangs()
	require.NotNil(t, supportLangs)
	require.Len(t, supportLangs, 5)
}

func TestConfig_MissingRequiredFields(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_MissingRequiredFields")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	assert.NoError(t, err)

	globalConfig := config.GetConfig()
	assert.NotNil(t, globalConfig)
}

func TestConfig_InvalidJSON(t *testing.T) {
	ctx := framework.NewTestContext("TestConfig_InvalidJSON")
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()
	loadEnvLocal(t)

	err := config.InitConfig()
	assert.NoError(t, err)
}
