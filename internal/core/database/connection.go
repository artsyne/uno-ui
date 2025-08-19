package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/biuren/hebe/internal/utils/config"
)

// 全局数据库连接
var db *sqlx.DB

// Connect 建立数据库连接
func Connect() error {
	config := config.GetConfig()
	if config == nil {
		return fmt.Errorf("global config is nil")
	}
	return ConnectWithConfig(config.Database)
}

// ConnectWithConfig 使用指定配置建立数据库连接
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
func GetDB() *sqlx.DB {
	return db
}

// Close 关闭数据库连接
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// SetDB 设置数据库连接 (主要用于测试)
func SetDB(connection string) error {
	var err error
	db, err = sqlx.Connect("mysql", connection)
	return err
}

// SetDBForTesting 设置数据库连接实例 (仅用于测试)
func SetDBForTesting(database *sqlx.DB) {
	db = database
}
