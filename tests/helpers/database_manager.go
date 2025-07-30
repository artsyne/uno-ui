package testhelpers

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/jmoiron/sqlx"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	db     *sqlx.DB
	mu     sync.RWMutex
	inited bool
}

var (
	dbManager     *DatabaseManager
	dbManagerOnce sync.Once
)

// GetDatabaseManager 获取数据库管理器单例
func GetDatabaseManager() *DatabaseManager {
	dbManagerOnce.Do(func() {
		dbManager = &DatabaseManager{}
	})
	return dbManager
}

// InitializeDatabase 初始化数据库连接
func (dm *DatabaseManager) InitializeDatabase(t *testing.T) *sqlx.DB {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.inited && dm.db != nil {
		return dm.db
	}

	// 确保配置已初始化
	SetupTestConfigSafe(t)

	// 初始化数据库
	err := database.Connect()
	if err != nil {
		if t != nil {
			t.Fatalf("数据库初始化失败: %v", err)
		}
		panic("数据库初始化失败: " + err.Error())
	}

	dm.db = database.GetDB()
	dm.inited = true
	return dm.db
}

// GetDatabase 获取数据库连接
func (dm *DatabaseManager) GetDatabase() *sqlx.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.db
}

// CleanupDatabase 清理数据库
func (dm *DatabaseManager) CleanupDatabase() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.db != nil {
		// 清理数据库中的测试数据
		if _, err := dm.db.Exec("DELETE FROM br_user WHERE username LIKE 'test_%'"); err != nil {
			return fmt.Errorf("清理用户数据失败: %w", err)
		}
		if _, err := dm.db.Exec("DELETE FROM br_diary WHERE title LIKE 'test_%'"); err != nil {
			return fmt.Errorf("清理日记数据失败: %w", err)
		}
		if _, err := dm.db.Exec("DELETE FROM br_comment WHERE content LIKE 'test_%'"); err != nil {
			return fmt.Errorf("清理评论数据失败: %w", err)
		}
	}

	return nil
}

// TransactionTestHelper 事务测试辅助器
type TransactionTestHelper struct {
	db *sqlx.DB
	tx *sql.Tx
}

// BeginTransaction 开始事务
func (tth *TransactionTestHelper) BeginTransaction(t *testing.T) {
	t.Helper()
	var err error
	tth.tx, err = tth.db.Begin()
	if err != nil {
		t.Fatalf("开始事务失败: %v", err)
	}
}

// RollbackTransaction 回滚事务
func (tth *TransactionTestHelper) RollbackTransaction(t *testing.T) {
	t.Helper()
	if tth.tx != nil {
		err := tth.tx.Rollback()
		if err != nil {
			t.Logf("回滚事务失败: %v", err)
		}
		tth.tx = nil
	}
}

// CommitTransaction 提交事务
func (tth *TransactionTestHelper) CommitTransaction(t *testing.T) {
	t.Helper()
	if tth.tx != nil {
		err := tth.tx.Commit()
		if err != nil {
			t.Fatalf("提交事务失败: %v", err)
		}
		tth.tx = nil
	}
}

// WithDatabaseTestImproved 改进的数据库测试辅助函数
// 提供更好的错误处理和资源管理
func WithDatabaseTestImproved(t *testing.T, testFunc func(*sqlx.DB)) {
	t.Helper()

	// 获取数据库连接
	db := GetDatabaseManager().GetDatabase()
	if db == nil {
		t.Fatal("无法获取数据库连接")
	}

	// 执行测试前清理
	if err := GetDatabaseManager().CleanupDatabase(); err != nil {
		t.Fatalf("测试前清理数据库失败: %v", err)
	}

	// 执行测试函数
	testFunc(db)

	// 执行测试后清理
	if err := GetDatabaseManager().CleanupDatabase(); err != nil {
		t.Errorf("测试后清理数据库失败: %v", err)
	}
}

// WithDatabaseTestSafe 安全的数据库测试辅助函数（向后兼容）
func WithDatabaseTestSafe(t *testing.T, testFunc func(*sqlx.DB)) {
	t.Helper()

	// 初始化数据库
	dbManager := GetDatabaseManager()
	db := dbManager.InitializeDatabase(t)

	// 在测试结束时清理数据库
	defer func() {
		if err := dbManager.CleanupDatabase(); err != nil {
			t.Errorf("清理数据库失败: %v", err)
		}
	}()

	// 执行测试函数
	testFunc(db)
}

// ResetDatabaseManager 重置数据库管理器（主要用于测试）
func ResetDatabaseManager() {
	dbManagerOnce = sync.Once{}
	dbManager = nil
}