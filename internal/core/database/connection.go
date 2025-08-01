package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/biuren/hebe/internal/utils/config"
)

// 全局数据库连接
var db *sqlx.DB

// Connect 建立数据库连接
func Connect() error {
	config := config.GetConfig()
	connection := fmt.Sprintf(config.Database.Conn,
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Name)

	var err error
	db, err = sqlx.Connect("mysql", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

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
