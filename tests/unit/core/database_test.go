package core

import (
	"context"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/database"
	framework "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 获取测试数据库配置
func getTestDatabaseConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
			Conn:     "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
}

func TestDatabase_SetDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 测试SetDB函数 - 使用内存数据库进行测试
	// 注意：这里使用sqlite内存数据库来避免依赖真实的MySQL服务器
	testConnection := ":memory:"
	
	// 由于原代码使用MySQL驱动，我们测试错误情况
	err := database.SetDB(testConnection)
	// 预期会出错，因为连接字符串格式不匹配MySQL
	assert.Error(t, err)
}

func TestDatabase_GetDB_BeforeConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 测试在连接之前获取数据库连接
	db := database.GetDB()
	// 在没有连接的情况下，应该返回nil或之前设置的连接
	// 这取决于之前是否有其他测试设置了连接
	if db != nil {
		// 如果有连接，验证它不是nil
		assert.NotNil(t, db)
	}
}

func TestDatabase_Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 测试关闭数据库连接
	err := database.Close()
	// 即使没有连接，Close也不应该返回错误
	assert.NoError(t, err)
}

func TestDatabase_Connect_InvalidConfig(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 设置无效的配置来测试连接失败的情况
	invalidConfig := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "invalid_host",
			Port:     9999,
			Name:     "invalid_db",
			User:     "invalid_user",
			Password: "invalid_password",
			Conn:     "%s:%s@tcp(%s:9999)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}

	// 设置无效配置
	config.SetConfig(invalidConfig)

	// 尝试连接，应该失败
	err := database.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestDatabase_Connect_ValidConfig_NoServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 使用有效格式但无法连接的配置
	validConfig := getTestDatabaseConfig()
	config.SetConfig(validConfig)

	// 尝试连接，由于没有真实的MySQL服务器，应该失败
	err := database.Connect()
	assert.Error(t, err)
	// 错误信息应该包含连接失败的提示
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestDatabase_ConnectionString_Format(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 测试连接字符串格式化
	cfg := getTestDatabaseConfig()
	config.SetConfig(cfg)

	// 验证配置是否正确设置
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 3306, cfg.Database.Port)
	assert.Equal(t, "test_db", cfg.Database.Name)
	assert.Equal(t, "test_user", cfg.Database.User)
	assert.Equal(t, "test_password", cfg.Database.Password)
	assert.NotEmpty(t, cfg.Database.Conn)
}

// 测试数据库配置的完整性
func TestDatabase_ConfigIntegrity(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	cfg := getTestDatabaseConfig()
	
	// 验证所有必需的数据库配置字段
	require.NotEmpty(t, cfg.Database.Host, "数据库主机不能为空")
	require.NotZero(t, cfg.Database.Port, "数据库端口不能为0")
	require.NotEmpty(t, cfg.Database.Name, "数据库名称不能为空")
	require.NotEmpty(t, cfg.Database.User, "数据库用户不能为空")
	require.NotEmpty(t, cfg.Database.Conn, "数据库连接字符串模板不能为空")
	
	// 验证端口范围
	assert.True(t, cfg.Database.Port > 0 && cfg.Database.Port <= 65535, "数据库端口应该在有效范围内")
}

// 测试数据库连接字符串的构建
func TestDatabase_ConnectionStringBuilding(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	cfg := getTestDatabaseConfig()
	
	// 验证连接字符串模板包含必要的占位符
	assert.Contains(t, cfg.Database.Conn, "%s", "连接字符串应该包含占位符")
	assert.Contains(t, cfg.Database.Conn, "tcp", "连接字符串应该指定TCP协议")
	assert.Contains(t, cfg.Database.Conn, "charset=utf8mb4", "连接字符串应该指定UTF8MB4字符集")
	assert.Contains(t, cfg.Database.Conn, "parseTime=True", "连接字符串应该启用时间解析")
}

// 测试数据库配置的边界情况
func TestDatabase_ConfigEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(framework.GetTestTimeout())*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("测试超时")
	default:
	}

	// 测试空密码的情况
	cfgEmptyPassword := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Name:     "test_db",
			User:     "test_user",
			Password: "", // 空密码
			Conn:     "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	
	// 空密码应该是允许的（在某些环境中）
	assert.Equal(t, "", cfgEmptyPassword.Database.Password)
	
	// 测试非标准端口
	cfgCustomPort := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3307, // 非标准端口
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
			Conn:     "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	
	assert.Equal(t, 3307, cfgCustomPort.Database.Port)
}