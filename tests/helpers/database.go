package testhelpers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/core/database"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseTestHelper 数据库测试辅助工具
type DatabaseTestHelper struct {
	db         *sql.DB
	tx         *sql.Tx
	testDBName string
}

// setupTestEnvironment 设置测试环境变量
// 基于.env.local文件配置，但使用测试数据库
func setupTestEnvironment(testDBName string) {
	// 设置环境变量指向.env.local文件
	projectRoot := getProjectRoot()
	envFile := filepath.Join(projectRoot, ".env.local")
	os.Setenv("ENV_FILE", envFile)

	// 设置必需的JSON格式环境变量（使用IDIARY_前缀）
	os.Setenv("IDIARY_SUPPORT_LANGS", `[{"code":"zh-cn","label":"中文","local-text":"中文"},{"code":"en","label":"English","local-text":"English"}]`)
	os.Setenv("IDIARY_TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("IDIARY_PARAMS", `[{"key":"test_key","val":"test_value"}]`)

	// 设置其他必需的配置项
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
	os.Setenv("IDIARY_LOCALES_DIR", "website/locales")
	os.Setenv("IDIARY_DEFAULT_LANG", "zh-cn")

	// 数据库配置 - 使用测试数据库名称覆盖.env.local配置
	os.Setenv("IDIARY_DB_HOST", "127.0.0.1")
	os.Setenv("IDIARY_DB_PORT", "3306")
	os.Setenv("IDIARY_DB_NAME", testDBName)
	os.Setenv("IDIARY_DB_USER", "root")
	os.Setenv("IDIARY_DB_PASSWORD", "123123")
	os.Setenv("IDIARY_DB_CONN", "%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_general_ci")

	// Redis配置
	os.Setenv("IDIARY_REDIS_HOST", "127.0.0.1")
	os.Setenv("IDIARY_REDIS_PORT", "6379")
	os.Setenv("IDIARY_REDIS_PASSWORD", "woaibeijingtiananmeng")

	// 邮件配置
	os.Setenv("IDIARY_SMTP_SERVER", "smtp.test.com")
	os.Setenv("IDIARY_SMTP_USER", "test@test.com")
	os.Setenv("IDIARY_SMTP_PASSWORD", "test_password")
	os.Setenv("IDIARY_SMTP_PORT", "587")
	os.Setenv("IDIARY_SMTP_FROM", "test@test.com")
	os.Setenv("IDIARY_EMAIL_MAX_RETRIES", "3")
	os.Setenv("IDIARY_EMAIL_BATCH_SIZE", "10")

	// Cookie配置
	os.Setenv("IDIARY_COOKIE_DOMAIN", "localhost")
	os.Setenv("IDIARY_COOKIE_SECRET", "test_secret_key_12345678901234567890")
	os.Setenv("IDIARY_COOKIE_EXPIRE_DAYS", "30")
	os.Setenv("IDIARY_COOKIE_HTTP_ONLY", "true")

	// 上传配置
	os.Setenv("IDIARY_UPLOAD_PATH", "uploads")
	os.Setenv("IDIARY_UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("IDIARY_UPLOAD_MAX_SIZE", "10485760")

	// 日志配置
	os.Setenv("IDIARY_LOG_LEVEL", "info")
	os.Setenv("IDIARY_LOG_FORMAT", "json")
	os.Setenv("IDIARY_LOG_PATH", "logs")
	os.Setenv("IDIARY_LOG_FILENAME", "test.log")
	os.Setenv("IDIARY_LOG_MAX_SIZE", "100")
	os.Setenv("IDIARY_LOG_MAX_AGE", "7")
	os.Setenv("IDIARY_LOG_MAX_BACKUPS", "5")
	os.Setenv("IDIARY_LOG_COMPRESS", "true")
	os.Setenv("IDIARY_LOG_LOCAL_TIME", "true")
}

// createTestDatabase 创建测试数据库
func createTestDatabase(testDBName string) error {
	// 连接到 MySQL 服务器（不指定数据库）- 使用.env.local中的真实配置
	dsn := "root:123123@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接 MySQL 服务器失败: %v", err)
	}
	defer db.Close()

	// 创建测试数据库
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", testDBName))
	if err != nil {
		return fmt.Errorf("创建测试数据库失败: %v", err)
	}

	return nil
}

// loadDatabaseSchema 加载数据库 schema
func loadDatabaseSchema(db *sql.DB) error {
	// 获取项目根目录
	projectRoot := getProjectRoot()

	// 读取 schema 文件
	schemaPath := filepath.Join(projectRoot, "sql", "schema.sql")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("读取 schema 文件失败: %v", err)
	}

	schemaSQL := string(schemaBytes)

	// 分割 SQL 语句（简单分割，按分号）
	statements := strings.Split(schemaSQL, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") || strings.HasPrefix(stmt, "/*") {
			continue
		}

		// 跳过一些 MySQL 特定的语句
		if strings.Contains(stmt, "/*!") || strings.Contains(stmt, "SET @") || strings.Contains(stmt, "SET NAMES") {
			continue
		}

		_, err := db.Exec(stmt)
		if err != nil {
			// 忽略一些常见的错误（如表已存在）
			if !strings.Contains(err.Error(), "already exists") &&
				!strings.Contains(err.Error(), "Table") &&
				!strings.Contains(err.Error(), "Unknown table") {
				return fmt.Errorf("执行 SQL 语句失败 [%s]: %v", stmt, err)
			}
		}
	}

	return nil
}

// dropTestDatabase 删除测试数据库
func dropTestDatabase(testDBName string) error {
	// 连接到 MySQL 服务器 - 使用.env.local中的真实配置
	dsn := "root:123123@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接 MySQL 服务器失败: %v", err)
	}
	defer db.Close()

	// 删除测试数据库
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", testDBName))
	if err != nil {
		return fmt.Errorf("删除测试数据库失败: %v", err)
	}

	return nil
}

// SetupTestDatabase 初始化测试数据库连接
// 在测试开始前调用，确保数据库连接正常
func SetupTestDatabase(t *testing.T) *DatabaseTestHelper {
	t.Helper()

	// 使用固定的测试数据库名称
	testDBName := "hebe_test"

	// 确保测试数据库存在并已加载schema
	if err := ensureTestDatabase(testDBName); err != nil {
		t.Fatalf("确保测试数据库存在失败: %v", err)
	}

	// 设置测试环境变量
	setupTestEnvironment(testDBName)

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 连接数据库
	database.Connect()

	db := database.GetDB().DB
	if db == nil {
		t.Fatal("数据库连接为空")
	}

	return &DatabaseTestHelper{
		db:         db,
		testDBName: testDBName,
	}
}

// BeginTransaction 开始数据库事务
// 用于测试中的数据隔离，测试结束后可以回滚
func (h *DatabaseTestHelper) BeginTransaction(t *testing.T) {
	t.Helper()

	var err error
	h.tx, err = h.db.Begin()
	if err != nil {
		t.Fatalf("开始事务失败: %v", err)
	}
}

// RollbackTransaction 回滚数据库事务
// 在测试结束后调用，清理测试数据
func (h *DatabaseTestHelper) RollbackTransaction(t *testing.T) {
	t.Helper()

	if h.tx != nil {
		if err := h.tx.Rollback(); err != nil {
			t.Errorf("回滚事务失败: %v", err)
		}
		h.tx = nil
	}
}

// CommitTransaction 提交数据库事务
// 在需要持久化测试数据时使用
func (h *DatabaseTestHelper) CommitTransaction(t *testing.T) {
	t.Helper()

	if h.tx != nil {
		if err := h.tx.Commit(); err != nil {
			t.Fatalf("提交事务失败: %v", err)
		}
		h.tx = nil
	}
}

// ExecuteSQL 执行SQL语句
// 用于测试中的数据准备和验证
func (h *DatabaseTestHelper) ExecuteSQL(t *testing.T, query string, args ...interface{}) sql.Result {
	t.Helper()

	var result sql.Result
	var err error

	if h.tx != nil {
		result, err = h.tx.Exec(query, args...)
	} else {
		result, err = h.db.Exec(query, args...)
	}

	if err != nil {
		t.Fatalf("执行SQL失败 [%s]: %v", query, err)
	}

	return result
}

// QueryRow 查询单行数据
func (h *DatabaseTestHelper) QueryRow(t *testing.T, query string, args ...interface{}) *sql.Row {
	t.Helper()

	if h.tx != nil {
		return h.tx.QueryRow(query, args...)
	}
	return h.db.QueryRow(query, args...)
}

// Query 查询多行数据
func (h *DatabaseTestHelper) Query(t *testing.T, query string, args ...interface{}) *sql.Rows {
	t.Helper()

	var rows *sql.Rows
	var err error

	if h.tx != nil {
		rows, err = h.tx.Query(query, args...)
	} else {
		rows, err = h.db.Query(query, args...)
	}

	if err != nil {
		t.Fatalf("查询数据失败 [%s]: %v", query, err)
	}

	return rows
}

// TruncateTable 清空表数据
// 用于测试前的数据清理
func (h *DatabaseTestHelper) TruncateTable(t *testing.T, tableName string) {
	t.Helper()

	query := fmt.Sprintf("DELETE FROM %s", tableName)
	h.ExecuteSQL(t, query)
}

// CountRows 统计表行数
// 用于验证测试结果
func (h *DatabaseTestHelper) CountRows(t *testing.T, tableName string) int {
	t.Helper()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	row := h.QueryRow(t, query)

	var count int
	if err := row.Scan(&count); err != nil {
		t.Fatalf("统计行数失败 [%s]: %v", tableName, err)
	}

	return count
}

// AssertRowCount 断言表行数
func (h *DatabaseTestHelper) AssertRowCount(t *testing.T, tableName string, expected int) {
	t.Helper()

	actual := h.CountRows(t, tableName)
	if actual != expected {
		t.Errorf("表 %s 行数不匹配: 期望 %d, 实际 %d", tableName, expected, actual)
	}
}

// Cleanup 清理测试环境
// 应该在测试结束时调用
func (h *DatabaseTestHelper) Cleanup(t *testing.T) {
	t.Helper()

	// 回滚未提交的事务
	h.RollbackTransaction(t)

	// 关闭数据库连接
	if h.db != nil {
		h.db.Close()
	}
}

// WithDatabaseTest 提供独立的数据库测试环境
// 适用于只需要数据库连接的测试，无需完整的测试套件
func WithDatabaseTest(t *testing.T, testFunc func(*DatabaseTestHelper)) {
	t.Helper()

	db := SetupTestDatabase(t)
	defer db.Cleanup(t)

	testFunc(db)
}

// GetDSN 生成数据库连接字符串
func GetDSN(dbName string) string {
	return fmt.Sprintf("root:@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbName)
}

// GetDB 获取数据库连接（用于直接数据库操作）
func (h *DatabaseTestHelper) GetDB() *sql.DB {
	return h.db
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 从当前文件位置开始，向上查找包含go.mod的目录
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达根目录，回退到默认路径
			return "/Users/yasong.cy/go/src/idiary.ai/hebe"
		}
		dir = parent
	}
}

// ensureTestDatabase 确保测试数据库存在并已加载schema
func ensureTestDatabase(testDBName string) error {
	// 连接到 MySQL 服务器（不指定数据库）
	dsn := "root:123123@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接 MySQL 服务器失败: %v", err)
	}
	defer db.Close()

	// 创建测试数据库（如果不存在）
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", testDBName))
	if err != nil {
		return fmt.Errorf("创建测试数据库失败: %v", err)
	}

	// 连接到测试数据库
	testDSN := fmt.Sprintf("root:123123@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", testDBName)
	testDB, err := sql.Open("mysql", testDSN)
	if err != nil {
		return fmt.Errorf("连接测试数据库失败: %v", err)
	}
	defer testDB.Close()

	// 检查是否需要加载schema（检查是否存在某个关键表）
	var tableCount int
	err = testDB.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'br_user'", testDBName).Scan(&tableCount)
	if err != nil {
		return fmt.Errorf("检查表是否存在失败: %v", err)
	}

	// 如果表不存在，加载schema
	if tableCount == 0 {
		if err := loadDatabaseSchema(testDB); err != nil {
			return fmt.Errorf("加载数据库schema失败: %v", err)
		}
	}

	return nil
}

// CleanupAllTables 清理所有表数据
// 用于测试前的全面数据清理，支持可选的表清理
func (h *DatabaseTestHelper) CleanupAllTables(t *testing.T, tables ...string) {
	t.Helper()

	// 如果没有指定表，使用默认的表列表
	if len(tables) == 0 {
		tables = []string{
			"posts",
			"diaries",
			"tasks",
			"users",
		}
	}

	// 禁用外键检查
	h.ExecuteSQL(t, "SET FOREIGN_KEY_CHECKS = 0")

	// 清理指定的表
	for _, table := range tables {
		h.TruncateTable(t, table)
	}

	// 重新启用外键检查
	h.ExecuteSQL(t, "SET FOREIGN_KEY_CHECKS = 1")
}
