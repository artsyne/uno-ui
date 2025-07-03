package app_test

import (
	"os"
	"testing"

	"github.com/biuren/hebe/internal/app/render"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/core/di"
)

// initTestConfig 初始化测试配置（不连接数据库）
func initTestConfig(t *testing.T) {
	t.Helper()

	// 直接设置测试环境变量
	os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-cn","label":"中文","local-text":"中文"},{"code":"en","label":"English","local-text":"English"}]`)
	os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("IDIARY_PARAMS", `[{"key":"test_key","val":"test_value"}]`)
	os.Setenv("IDIARY_ENV", "test")
	os.Setenv("IDIARY_GIN_MODE", "debug")
	os.Setenv("IDIARY_DOMAIN", "http://localhost:8080")
	os.Setenv("IDIARY_PORT", "8080")
	os.Setenv("IDIARY_SUMMARY_MAX", "200")
	os.Setenv("IDIARY_SUFFIX", "...")
	os.Setenv("IDIARY_PAGE_SIZE", "10")
	os.Setenv("IDIARY_MAX_DAYBOOK_LIMIT", "100")
	os.Setenv("IDIARY_HASH_SALT", "test_salt_12345")
	os.Setenv("IDIARY_STATIC_DIR", "website/static")
	os.Setenv("IDIARY_DB_HOST", "localhost")
	os.Setenv("IDIARY_DB_PORT", "3306")
	os.Setenv("IDIARY_DB_NAME", "test_db")
	os.Setenv("IDIARY_DB_USER", "test_user")
	os.Setenv("IDIARY_DB_PASSWORD", "test_password")
	os.Setenv("IDIARY_DB_CONN", "mysql")
	os.Setenv("IDIARY_REDIS_HOST", "localhost")
	os.Setenv("IDIARY_REDIS_PORT", "6379")
	os.Setenv("IDIARY_REDIS_PASSWORD", "test_redis_password")
	os.Setenv("IDIARY_SMTP_SERVER", "smtp.test.com")
	os.Setenv("IDIARY_SMTP_USER", "test@test.com")
	os.Setenv("IDIARY_SMTP_PASSWORD", "test_password")
	os.Setenv("IDIARY_SMTP_PORT", "587")
	os.Setenv("IDIARY_SMTP_FROM", "test@test.com")
	os.Setenv("IDIARY_EMAIL_MAX_RETRIES", "3")
	os.Setenv("IDIARY_EMAIL_BATCH_SIZE", "10")
	os.Setenv("IDIARY_COOKIE_DOMAIN", "localhost")
	os.Setenv("IDIARY_COOKIE_SECRET", "test_secret_key_12345")
	os.Setenv("IDIARY_COOKIE_EXPIRE_DAYS", "30")
	os.Setenv("IDIARY_COOKIE_HTTP_ONLY", "true")
	os.Setenv("IDIARY_UPLOAD_PATH", "uploads")
	os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")
	os.Setenv("IDIARY_LOG_LEVEL", "info")
	os.Setenv("IDIARY_LOG_FORMAT", "json")
	os.Setenv("IDIARY_LOG_PATH", "logs")
	os.Setenv("IDIARY_LOG_FILENAME", "test.log")
	os.Setenv("IDIARY_LOG_MAX_SIZE", "100")
	os.Setenv("IDIARY_LOG_MAX_AGE", "7")
	os.Setenv("IDIARY_LOG_MAX_BACKUPS", "5")
	os.Setenv("IDIARY_LOG_COMPRESS", "true")
	os.Setenv("IDIARY_LOG_LOCAL_TIME", "true")

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}
}

// TestDIContainerIntegration 测试DI容器和App层服务的集成
func TestDIContainerIntegration(t *testing.T) {
	// 初始化测试配置
	initTestConfig(t)

	// 创建DI容器
	container := di.NewContainer()

	// 注册App层服务
	if err := di.RegisterAppServices(container); err != nil {
		t.Fatalf("注册App层服务失败: %v", err)
	}

	t.Run("测试RenderService解析", func(t *testing.T) {
		service, err := di.Resolve[render.RenderServiceInterface](container)
		if err != nil {
			t.Fatalf("解析RenderService失败: %v", err)
		}

		if service == nil {
			t.Fatal("RenderService不应该为nil")
		}
	})
}

// TestHandlerCreation 测试通过DI容器创建Handler
func TestHandlerCreation(t *testing.T) {
	// 初始化测试配置
	initTestConfig(t)

	// 创建DI容器
	container := di.NewContainer()

	// 注册App层服务
	if err := di.RegisterAppServices(container); err != nil {
		t.Fatalf("注册App层服务失败: %v", err)
	}

	// 测试获取服务用于创建handlers
	responseService, err := di.Resolve[response.ResponseServiceInterface](container)
	if err != nil {
		t.Fatalf("解析ResponseService失败: %v", err)
	}

	viewService, err := di.Resolve[view.ViewServiceInterface](container)
	if err != nil {
		t.Fatalf("解析ViewService失败: %v", err)
	}

	// 验证能够成功创建各模块的Handler
	t.Run("测试创建UserHandler", func(t *testing.T) {
		// 这里只是验证能够获取到服务，实际的Handler创建会在路由测试中完成
		if responseService == nil || viewService == nil {
			t.Fatal("创建UserHandler所需的服务不完整")
		}
	})

	t.Run("测试创建DiaryHandler", func(t *testing.T) {
		if responseService == nil || viewService == nil {
			t.Fatal("创建DiaryHandler所需的服务不完整")
		}
	})

	t.Run("测试创建PostHandler", func(t *testing.T) {
		if responseService == nil || viewService == nil {
			t.Fatal("创建PostHandler所需的服务不完整")
		}
	})
}

// TestConfigDependency 测试配置依赖
func TestConfigDependency(t *testing.T) {
	// 初始化测试配置
	initTestConfig(t)

	// 验证配置正确加载
	cfg := config.GetConfig()
	if cfg == nil {
		t.Fatal("配置未正确加载")
	}

	// 创建ViewService验证配置依赖
	viewService := view.NewViewService(cfg)
	if viewService == nil {
		t.Fatal("使用配置创建ViewService失败")
	}

	// 创建ResponseService验证依赖关系
	responseService := response.NewResponseService(viewService)
	if responseService == nil {
		t.Fatal("创建ResponseService失败")
	}
}
