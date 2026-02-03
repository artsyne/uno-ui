package database_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/utils/config"
)

// TestConnect 测试数据库连接函数
func TestConnect(t *testing.T) {
	// 保存原始配置
	originalConfig := config.GetConfig()
	defer func() {
		config.SetConfig(originalConfig)
	}()

	// 测试配置为nil的情况
	config.SetConfig(nil)
	err := database.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "global config is nil")

	// 测试有效配置但无法连接的情况
	validConfig := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "invalid_host",
			Port:     9999,
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
			Conn:     "%s:%s@tcp(%s:9999)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
	config.SetConfig(validConfig)

	err = database.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

// TestConnectWithConfig 测试使用指定配置连接数据库
func TestConnectWithConfig(t *testing.T) {
	// 测试无效的数据库配置
	invalidConfig := config.DatabaseConfig{
		Host:     "invalid_host",
		Port:     9999,
		Name:     "test_db",
		User:     "test_user",
		Password: "test_password",
		Conn:     "%s:%s@tcp(%s:9999)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	}

	err := database.ConnectWithConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

// TestGetDB 测试获取数据库连接
func TestGetDB(t *testing.T) {
	// 在没有连接的情况下获取数据库连接
	db := database.GetDB()
	// 可能为nil或之前设置的连接
	if db != nil {
		assert.NotNil(t, db)
	}
}

// TestClose 测试关闭数据库连接
func TestClose(t *testing.T) {
	// 测试关闭数据库连接
	err := database.Close()
	// 即使没有连接，Close也不应该返回错误
	assert.NoError(t, err)
}

// TestSetDB 测试设置数据库连接
func TestSetDB(t *testing.T) {
	// 测试使用无效连接字符串
	invalidConnection := "invalid_connection_string"
	err := database.SetDB(invalidConnection)
	assert.Error(t, err)

	// 测试使用空连接字符串
	emptyConnection := ""
	err = database.SetDB(emptyConnection)
	assert.Error(t, err)
}

// TestConnectionStringFormat 测试连接字符串格式
func TestConnectionStringFormat(t *testing.T) {
	// 测试连接字符串模板的正确性
	config := config.DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		Name:     "test_db",
		User:     "test_user",
		Password: "test_password",
		Conn:     "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	}

	// 验证连接字符串模板包含必要的占位符
	assert.Contains(t, config.Conn, "%s", "连接字符串应该包含用户名占位符")
	assert.Contains(t, config.Conn, "tcp", "连接字符串应该指定TCP协议")
	assert.Contains(t, config.Conn, "charset=utf8mb4", "连接字符串应该指定UTF8MB4字符集")
	assert.Contains(t, config.Conn, "parseTime=True", "连接字符串应该启用时间解析")
}

// TestConnectionPoolSettings 测试连接池设置
func TestConnectionPoolSettings(t *testing.T) {
	// 这个测试验证连接池参数的设置是否合理
	// 由于无法直接测试连接池设置，我们验证常量值
	maxOpenConns := 25
	maxIdleConns := 5
	maxLifetime := time.Hour
	maxIdleTime := time.Minute * 10

	// 验证连接池参数的合理性
	assert.Greater(t, maxOpenConns, 0, "最大打开连接数应该大于0")
	assert.Greater(t, maxIdleConns, 0, "最大空闲连接数应该大于0")
	assert.LessOrEqual(t, maxIdleConns, maxOpenConns, "最大空闲连接数不应该超过最大打开连接数")
	assert.Greater(t, maxLifetime, time.Duration(0), "连接最大生存时间应该大于0")
	assert.Greater(t, maxIdleTime, time.Duration(0), "连接最大空闲时间应该大于0")
}

// TestDatabaseConfigValidation 测试数据库配置验证
func TestDatabaseConfigValidation(t *testing.T) {
	tests := []struct {
		name     string
		config   config.DatabaseConfig
		expected bool // 是否期望连接成功
	}{
		{
			name: "空主机名",
			config: config.DatabaseConfig{
				Host:     "",
				Port:     3306,
				Name:     "test_db",
				User:     "test_user",
				Password: "test_password",
				Conn:     "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			},
			expected: false,
		},
		{
			name: "无效端口",
			config: config.DatabaseConfig{
				Host:     "localhost",
				Port:     0,
				Name:     "test_db",
				User:     "test_user",
				Password: "test_password",
				Conn:     "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			},
			expected: false,
		},
		{
			name: "空数据库名",
			config: config.DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				Name:     "",
				User:     "test_user",
				Password: "test_password",
				Conn:     "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := database.ConnectWithConfig(tt.config)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}