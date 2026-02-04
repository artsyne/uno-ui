package core

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	. "github.com/biuren/hebe/tests/framework/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DatabaseManager 数据库管理器
// 提供统一的数据库连接管理，替代当前混乱的数据库连接机制
type DatabaseManager struct {
	db       *sqlx.DB
	rawDB    *sql.DB
	config   *TestConfig
	dsn      string
	mu       sync.RWMutex
	healthy  bool
	initOnce sync.Once
}

var (
	globalDBManager   *DatabaseManager
	dbManagerInitOnce sync.Once
)

// GetDatabaseManager 获取数据库管理器单例
func GetDatabaseManager() *DatabaseManager {
	dbManagerInitOnce.Do(func() {
		config := GetTestConfig()
		globalDBManager = &DatabaseManager{
			config: config,
		}
	})
	return globalDBManager
}

// GetDatabaseName 获取当前数据库的名称
func (dm *DatabaseManager) GetDatabaseName() string {
	return dm.config.Database.Name
}

// Initialize 初始化数据库连接
func (dm *DatabaseManager) Initialize() error {
	var initErr error
	dm.initOnce.Do(func() {
		initErr = dm.doInitialize()
	})
	return initErr
}

// doInitialize 执行实际的初始化
func (dm *DatabaseManager) doInitialize() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 构建DSN
	dm.dsn = dm.buildDSN(dm.config.Database.Name)

	// 确保测试数据库存在
	if err := dm.ensureTestDatabase(); err != nil {
		return fmt.Errorf("确保测试数据库存在失败: %w", err)
	}

	// 创建数据库连接
	db, err := sqlx.Connect("mysql", dm.dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(dm.config.Database.MaxConns)
	db.SetMaxIdleConns(dm.config.Database.MaxIdle)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	dm.db = db
	dm.rawDB = db.DB
	dm.healthy = true

	slog.Debug("数据库管理器初始化成功",
		"database", dm.config.Database.Name,
		"max_conns", dm.config.Database.MaxConns,
		"max_idle", dm.config.Database.MaxIdle)

	return nil
}

// buildDSN 构建数据库连接字符串
func (dm *DatabaseManager) buildDSN(dbName string) string {
	cfg := dm.config.Database
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&collation=%s_general_ci",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, dbName, cfg.Charset, cfg.Charset)
}

// ensureTestDatabase 确保测试数据库存在并加载schema
func (dm *DatabaseManager) ensureTestDatabase() error {
	slog.Debug("开始检查测试数据库", "database", dm.config.Database.Name)
	
	// 连接到MySQL服务器（不指定数据库）
	serverDSN := dm.buildDSN("")
	slog.Debug("连接MySQL服务器", "dsn", dm.maskPassword(serverDSN))
	
	serverDB, err := sql.Open("mysql", serverDSN)
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %w", err)
	}
	defer serverDB.Close()

	// 测试服务器连接
	if err := serverDB.Ping(); err != nil {
		return fmt.Errorf("MySQL服务器连接测试失败: %w", err)
	}
	slog.Debug("MySQL服务器连接成功")

	// 检查数据库是否存在
	var dbExists int
	err = serverDB.QueryRow("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?",
		dm.config.Database.Name).Scan(&dbExists)
	if err != nil {
		return fmt.Errorf("检查数据库是否存在失败: %w", err)
	}

	// 创建测试数据库（如果不存在）
	if dbExists == 0 {
		slog.Info("测试数据库不存在，正在创建", "database", dm.config.Database.Name)
		createDBSQL := fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET %s COLLATE %s_unicode_ci",
			dm.config.Database.Name, dm.config.Database.Charset, dm.config.Database.Charset)

		if _, err := serverDB.Exec(createDBSQL); err != nil {
			return fmt.Errorf("创建测试数据库失败: %w", err)
		}
		slog.Info("测试数据库创建成功", "database", dm.config.Database.Name)
	} else {
		slog.Debug("测试数据库已存在", "database", dm.config.Database.Name)
	}

	// 连接到测试数据库检查是否需要加载schema
	testDSN := dm.buildDSN(dm.config.Database.Name)
	testDB, err := sql.Open("mysql", testDSN)
	if err != nil {
		return fmt.Errorf("连接测试数据库失败: %w", err)
	}
	defer testDB.Close()

	// 测试数据库连接
	if err := testDB.Ping(); err != nil {
		return fmt.Errorf("测试数据库连接测试失败: %w", err)
	}
	slog.Debug("测试数据库连接成功")

	// 检查是否存在关键表
	var tableCount int
	err = testDB.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'br_user'",
		dm.config.Database.Name).Scan(&tableCount)
	if err != nil {
		return fmt.Errorf("检查表是否存在失败: %w", err)
	}

	// 如果表不存在，加载schema
	if tableCount == 0 {
		slog.Info("数据库表不存在，正在加载schema", "database", dm.config.Database.Name)
		if err := dm.loadDatabaseSchema(testDB); err != nil {
			return fmt.Errorf("加载数据库schema失败: %w", err)
		}
		slog.Info("数据库schema加载完成", "database", dm.config.Database.Name)
	} else {
		slog.Debug("数据库表已存在，跳过schema加载", "database", dm.config.Database.Name)
	}

	return nil
}

// maskPassword 隐藏DSN中的密码信息用于日志记录
func (dm *DatabaseManager) maskPassword(dsn string) string {
	// 查找密码部分并替换为***
	if strings.Contains(dsn, ":") && strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		if len(parts) >= 2 {
			userPart := parts[0]
			if colonIndex := strings.LastIndex(userPart, ":"); colonIndex != -1 {
				maskedUserPart := userPart[:colonIndex+1] + "***"
				return maskedUserPart + "@" + strings.Join(parts[1:], "@")
			}
		}
	}
	return dsn
}

// loadDatabaseSchema 加载数据库schema
func (dm *DatabaseManager) loadDatabaseSchema(db *sql.DB) error {
	projectRoot := GetProjectRoot()
	schemaPath := filepath.Join(projectRoot, "sql", "schema.sql")
	slog.Debug("准备加载数据库schema", "schema_path", schemaPath)

	// 检查schema文件是否存在
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema文件不存在: %s", schemaPath)
	}

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("读取schema文件失败: %w", err)
	}
	slog.Debug("schema文件读取成功", "size", len(schemaBytes))

	schemaSQL := string(schemaBytes)

	// 分割SQL语句
	sqlStatements := dm.splitSQL(schemaSQL)
	slog.Debug("SQL语句分割完成", "statement_count", len(sqlStatements))

	successCount := 0
	skippedCount := 0
	for i, stmt := range sqlStatements {
		// 处理语句中的注释行
		lines := strings.Split(stmt, "\n")
		var validLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "--") || (strings.HasPrefix(line, "/*") && !strings.HasPrefix(line, "/*!")) {
				continue
			}
			validLines = append(validLines, line)
		}
		stmt = strings.Join(validLines, " ")

		if stmt == "" {
			continue
		}

		// 跳过MySQL特定的语句
		if strings.Contains(stmt, "/*!") || strings.Contains(stmt, "SET @") ||
			strings.Contains(stmt, "SET NAMES") || strings.Contains(stmt, "LOCK TABLES") ||
			strings.Contains(stmt, "UNLOCK TABLES") {
			skippedCount++
			slog.Debug("跳过MySQL特定语句", "statement_index", i+1)
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			// 忽略一些常见的错误
			if strings.Contains(err.Error(), "already exists") ||
				strings.Contains(err.Error(), "Unknown table") ||
				strings.Contains(err.Error(), "Duplicate") {
				skippedCount++
				slog.Debug("跳过已存在的对象", "statement_index", i+1, "error", err.Error())
				continue
			}
			sqlPreview := stmt
			if len(stmt) > 100 {
				sqlPreview = stmt[:100] + "..."
			}
			slog.Error("执行SQL语句失败", "statement_index", i+1, "error", err.Error(), "sql", sqlPreview)
			return fmt.Errorf("执行SQL语句失败 (第%d条): %w\nSQL: %s", i+1, err, stmt)
		}
		successCount++
		slog.Debug("SQL语句执行成功", "statement_index", i+1)
	}

	slog.Info("数据库schema加载统计", "success_count", successCount, "skipped_count", skippedCount, "total_count", len(sqlStatements))
	return nil
}

// splitSQL 分割SQL语句
func (dm *DatabaseManager) splitSQL(sqlContent string) []string {
	// 简单的SQL分割，按分号分割
	statements := strings.Split(sqlContent, ";")
	var result []string

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}

// GetDB 获取sqlx数据库连接
func (dm *DatabaseManager) GetDB() *sqlx.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.db
}

// GetRawDB 获取原始数据库连接
func (dm *DatabaseManager) GetRawDB() *sql.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.rawDB
}

// IsHealthy 检查数据库是否健康
func (dm *DatabaseManager) IsHealthy() bool {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if !dm.healthy || dm.db == nil {
		return false
	}

	// 尝试ping数据库
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := dm.db.PingContext(ctx); err != nil {
		dm.healthy = false
		return false
	}

	return true
}

// Close 关闭数据库连接
func (dm *DatabaseManager) Close() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.db != nil {
		err := dm.db.Close()
		dm.db = nil
		dm.rawDB = nil
		dm.healthy = false
		return err
	}

	return nil
}

// Execute 执行SQL语句
func (dm *DatabaseManager) Execute(query string, args ...interface{}) (sql.Result, error) {
	if !dm.IsHealthy() {
		return nil, fmt.Errorf("数据库连接不健康")
	}

	return dm.db.Exec(query, args...)
}

// Query 查询数据
func (dm *DatabaseManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if !dm.IsHealthy() {
		return nil, fmt.Errorf("数据库连接不健康")
	}

	return dm.db.Query(query, args...)
}

// QueryRow 查询单行数据
func (dm *DatabaseManager) QueryRow(query string, args ...interface{}) *sql.Row {
	return dm.db.QueryRow(query, args...)
}

// WithTransaction 在事务中执行函数
// 注意：新框架中我们尽量避免使用事务，但保留此方法用于特殊需求
func (dm *DatabaseManager) WithTransaction(fn func(*sql.Tx) error) error {
	if !dm.IsHealthy() {
		return fmt.Errorf("数据库连接不健康")
	}

	tx, err := dm.db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("事务回滚失败: %w", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

// GetStats 获取数据库连接池统计信息
func (dm *DatabaseManager) GetStats() sql.DBStats {
	if dm.rawDB == nil {
		return sql.DBStats{}
	}
	return dm.rawDB.Stats()
}
