package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biuren/hebe/internal/utils/config"
	framework "github.com/biuren/hebe/tests/framework"
)

// TestConfigIntegration 测试配置系统的集成功能
func TestConfigIntegration(t *testing.T) {
	ctx := framework.NewTestContext()
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("测试框架配置", func(t *testing.T) {
		// 获取全局配置（通过测试框架）
		cfg := framework.GetGlobalConfig()
		assert.NotNil(t, cfg, "测试框架配置不应该为空")

		// 验证数据库配置
		assert.NotEmpty(t, cfg.Database.Host, "数据库主机不应该为空")
		assert.Greater(t, cfg.Database.Port, 0, "数据库端口应该大于0")
		assert.NotEmpty(t, cfg.Database.Name, "数据库名称不应该为空")
		assert.NotEmpty(t, cfg.Database.User, "数据库用户不应该为空")
		assert.NotEmpty(t, cfg.Database.Conn, "数据库连接字符串模板不应该为空")

		// 验证应用配置
		assert.NotEmpty(t, cfg.App.Port, "应用端口不应该为空")
		assert.NotEmpty(t, cfg.App.GinMode, "Gin模式不应该为空")
		assert.NotEmpty(t, cfg.App.Domain, "应用域名不应该为空")
		assert.NotEmpty(t, cfg.Env, "应用环境不应该为空")
	})

	t.Run("配置验证", func(t *testing.T) {
		// 获取配置并验证基本字段
		cfg := framework.GetGlobalConfig()
		assert.NotNil(t, cfg, "配置不应该为空")
		
		// 验证必需的配置字段
		assert.NotEmpty(t, cfg.Env, "环境配置不应该为空")
		assert.NotEmpty(t, cfg.App.GinMode, "Gin模式不应该为空")
		assert.NotEmpty(t, cfg.App.Domain, "应用域名不应该为空")
		assert.NotEmpty(t, cfg.App.Port, "应用端口不应该为空")
		
		// 验证数据库配置
		assert.NotEmpty(t, cfg.Database.Host, "数据库主机不应该为空")
		assert.Greater(t, cfg.Database.Port, 0, "数据库端口应该大于0")
		assert.NotEmpty(t, cfg.Database.Name, "数据库名称不应该为空")
		assert.NotEmpty(t, cfg.Database.User, "数据库用户不应该为空")
	})

	t.Run("配置函数可调用", func(t *testing.T) {
		// 测试配置函数不会panic
		assert.NotPanics(t, func() {
			config.GetSupportLangs()
		}, "GetSupportLangs不应该panic")

		assert.NotPanics(t, func() {
			config.GetTimeZone()
		}, "GetTimeZone不应该panic")

		assert.NotPanics(t, func() {
			config.GetParams()
		}, "GetParams不应该panic")

		assert.NotPanics(t, func() {
			config.UploadsPath()
		}, "UploadsPath不应该panic")

		assert.NotPanics(t, func() {
			config.GetViewsPath()
		}, "GetViewsPath不应该panic")
	})
}

// TestConfigViper 测试Viper配置管理器
func TestConfigViper(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	t.Run("Viper实例", func(t *testing.T) {
		// 测试框架已经初始化了配置，直接获取Viper实例
		v := config.GetViper()
		assert.NotNil(t, v, "Viper实例不应该为空")
		
		// 验证Viper可以读取配置
		assert.True(t, v.IsSet("ENV"), "ENV配置应该存在")
	})
}

// TestConfigComponents 测试配置各个组件
func TestConfigComponents(t *testing.T) {
	ctx := framework.NewTestContext()
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	cfg := framework.GetGlobalConfig()
	assert.NotNil(t, cfg, "配置不应该为空")

	t.Run("配置结构完整性", func(t *testing.T) {
		// 验证配置结构体字段存在
		assert.NotNil(t, &cfg.User, "用户配置结构应该存在")
		assert.NotNil(t, &cfg.ValidCode, "验证码配置结构应该存在")
		assert.NotNil(t, &cfg.Redis, "Redis配置结构应该存在")
		assert.NotNil(t, &cfg.Email, "邮件配置结构应该存在")
		assert.NotNil(t, &cfg.Cookie, "Cookie配置结构应该存在")
		assert.NotNil(t, &cfg.Upload, "上传配置结构应该存在")
		assert.NotNil(t, &cfg.Log, "日志配置结构应该存在")
	})

	t.Run("配置字段类型", func(t *testing.T) {
		// 验证配置字段类型正确
		assert.IsType(t, "", cfg.Env, "环境配置应该是字符串类型")
		assert.IsType(t, []config.SupportLang{}, cfg.SupportLangs, "支持语言应该是切片类型")
		assert.IsType(t, []string{}, cfg.Timezones, "时区应该是字符串切片类型")
		assert.IsType(t, []config.Param{}, cfg.Params, "参数应该是参数切片类型")
	})
}
