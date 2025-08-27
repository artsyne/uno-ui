package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/biuren/hebe/internal/utils/config"
)

// db 全局数据库连接实例
// 使用 sqlx.DB 提供增强的数据库操作功能
var db *sqlx.DB

// Connect 建立数据库连接
// 使用全局配置中的数据库配置建立数据库连接
//
// 返回:
//   - error: 连接失败时返回错误，成功时返回nil
func Connect() error {
	config := config.GetConfig()
	if config == nil {
		return fmt.Errorf("global config is nil")
	}
	return ConnectWithConfig(config.Database)
}

// ConnectWithConfig 使用指定配置建立数据库连接
// 根据提供的数据库配置创建MySQL连接，并设置连接池参数
//
// 参数:
//   - dbConfig: 数据库配置对象，包含连接字符串、用户名、密码等信息
//
// 返回:
//   - error: 连接失败时返回错误，成功时返回nil
func ConnectWithConfig(dbConfig config.DatabaseConfig) error {
	connection := fmt.Sprintf(dbConfig.Conn,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Name)

	var err error
	db, err = sqlx.Connect("mysql", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)                      // 最大打开连接数
	db.SetMaxIdleConns(5)                       // 最大空闲连接数
	db.SetConnMaxLifetime(time.Hour)            // 连接最大生存时间
	db.SetConnMaxIdleTime(time.Minute * 10)     // 连接最大空闲时间

	// 测试连接
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// GetDB 获取数据库连接
// 返回全局数据库连接实例，用于执行数据库操作
//
// 返回:
//   - *sqlx.DB: 数据库连接实例，如果未初始化则返回nil
func GetDB() *sqlx.DB {
	return db
}

// Close 关闭数据库连接
// 安全地关闭全局数据库连接，释放相关资源
//
// 返回:
//   - error: 关闭失败时返回错误，成功或连接为nil时返回nil
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// SetDB 设置数据库连接 (主要用于测试)
// 使用指定的连接字符串创建新的数据库连接，主要用于测试环境
//
// 参数:
//   - connection: MySQL连接字符串
//
// 返回:
//   - error: 连接失败时返回错误，成功时返回nil
func SetDB(connection string) error {
	var err error
	db, err = sqlx.Connect("mysql", connection)
	return err
}

// SetDBForTesting 设置数据库连接实例 (仅用于测试)
// 直接设置全局数据库连接实例，仅用于单元测试中的依赖注入
//
// 参数:
//   - database: 数据库连接实例
func SetDBForTesting(database *sqlx.DB) {
	db = database
}
