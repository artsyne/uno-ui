package testhelpers

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/biuren/hebe/internal/core/config"
	"github.com/biuren/hebe/internal/core/database"
)

// DatabaseTestHelper 数据库测试辅助工具
type DatabaseTestHelper struct {
	db *sql.DB
	tx *sql.Tx
}

// setupTestEnvironment 设置测试环境变量
func setupTestEnvironment() {
	// 设置必需的JSON格式环境变量
	os.Setenv("SUPPORT_LANGS", `[{"code":"zh-cn","label":"中文","local-text":"中文"},{"code":"en","label":"English","local-text":"English"}]`)
	os.Setenv("TIMEZONES", `["Asia/Shanghai","UTC","America/New_York"]`)
	os.Setenv("PARAMS", `[{"key":"test_key","val":"test_value"}]`)

	// 设置其他必需的配置项
	os.Setenv("ENV", "test")
	os.Setenv("GIN_MODE", "debug")
	os.Setenv("DOMAIN", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("SUMMARY_MAX", "200")
	os.Setenv("SUFFIX", "...")
	os.Setenv("PAGE_SIZE", "10")
	os.Setenv("MAX_DAYBOOK_LIMIT", "100")
	os.Setenv("HASH_SALT", "test_salt_12345")
	os.Setenv("STATIC_DIR", "website/static")

	// 数据库配置
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_CONN", "mysql")

	// Redis配置
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")

	// 邮件配置
	os.Setenv("SMTP_SERVER", "smtp.test.com")
	os.Setenv("SMTP_USER", "test@test.com")
	os.Setenv("SMTP_PASSWORD", "test_password")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_FROM", "test@test.com")
	os.Setenv("EMAIL_MAX_RETRIES", "3")
	os.Setenv("EMAIL_BATCH_SIZE", "10")

	// Cookie配置
	os.Setenv("COOKIE_DOMAIN", "localhost")
	os.Setenv("COOKIE_SECRET", "test_secret_key_12345")
	os.Setenv("COOKIE_EXPIRE_DAYS", "30")
	os.Setenv("COOKIE_HTTP_ONLY", "true")

	// 上传配置
	os.Setenv("UPLOAD_PATH", "uploads")
	os.Setenv("UPLOAD_TYPES", "jpg,jpeg,png,gif")
	os.Setenv("UPLOAD_MAX_SIZE", "10485760")

	// 日志配置
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOG_FORMAT", "json")
	os.Setenv("LOG_PATH", "logs")
	os.Setenv("LOG_FILENAME", "test.log")
	os.Setenv("LOG_MAX_SIZE", "100")
	os.Setenv("LOG_MAX_AGE", "7")
	os.Setenv("LOG_MAX_BACKUPS", "5")
	os.Setenv("LOG_COMPRESS", "true")
	os.Setenv("LOG_LOCAL_TIME", "true")
}

// SetupTestDatabase 初始化测试数据库连接
// 在测试开始前调用，确保数据库连接正常
func SetupTestDatabase(t *testing.T) *DatabaseTestHelper {
	t.Helper()

	// 设置测试环境变量
	setupTestEnvironment()

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
		db: db,
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
}
