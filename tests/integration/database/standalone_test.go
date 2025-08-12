package database

import (
	"testing"

	framework "github.com/biuren/hebe/tests/framework"
)

// TestDatabaseOnly 仅使用数据库的测试示例
func TestDatabaseOnly(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 获取数据库连接
	db := ctx.GetDatabase().GetDB()

	// 测试数据库基本功能
	result := db.QueryRow("SELECT DATABASE() as current_db")
	var currentDB string
	err := result.Scan(&currentDB)
	if err != nil {
		t.Fatalf("查询当前数据库失败: %v", err)
	}

	// 验证使用的是测试数据库
	if currentDB == "" {
		t.Error("当前数据库名不应该为空")
	}

	t.Logf("当前使用的测试数据库: %s", currentDB)
}

// TestTableStructure 测试表结构
func TestTableStructure(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 获取数据库连接
	db := ctx.GetDatabase().GetDB()

	// 检查关键表是否存在
	tables := []string{"br_user", "br_diary", "br_comment", "br_follow"}

	for _, tableName := range tables {
		result := db.QueryRow(`
			SELECT COUNT(*) 
			FROM information_schema.tables 
			WHERE table_schema = DATABASE() AND table_name = ?
		`, tableName)

		var count int
		err := result.Scan(&count)
		if err != nil {
			t.Fatalf("检查表 %s 失败: %v", tableName, err)
		}

		if count != 1 {
			t.Errorf("表 %s 不存在", tableName)
		}
	}
}

// TestUserTableConstraints 测试用户表约束
func TestUserTableConstraints(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 获取数据库连接
	db := ctx.GetDatabase().GetDB()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("开始事务失败: %v", err)
	}
	defer tx.Rollback()

	// 插入正常用户
	_, err = tx.Exec(`
		INSERT INTO br_user (user_name, email, pwd, status, created_utc_time) 
		VALUES (?, ?, ?, ?, ?)
	`, "uniqueuser", "unique@example.com", "password", 1, 1642680000)
	if err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}

	// 测试用户名唯一约束
	_, err = tx.Exec(`
		INSERT INTO br_user (user_name, email, pwd, status, created_utc_time) 
		VALUES (?, ?, ?, ?, ?)
	`, "uniqueuser", "another@example.com", "password", 1, 1642680000)

	if err == nil {
		t.Error("期望用户名重复时报错，但没有报错")
	}

	// 测试邮箱唯一约束
	_, err = tx.Exec(`
		INSERT INTO br_user (user_name, email, pwd, status, created_utc_time) 
		VALUES (?, ?, ?, ?, ?)
	`, "anotheruser", "unique@example.com", "password", 1, 1642680000)

	if err == nil {
		t.Error("期望邮箱重复时报错，但没有报错")
	}
}

// TestTransactionIsolation 测试事务隔离
func TestTransactionIsolation(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 获取数据库连接
	db := ctx.GetDatabase().GetDB()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("开始事务失败: %v", err)
	}

	// 在事务中插入数据
	_, err = tx.Exec(`
		INSERT INTO br_user (user_name, email, pwd, status, created_utc_time) 
		VALUES (?, ?, ?, ?, ?)
	`, "txuser", "tx@example.com", "password", 1, 1642680000)
	if err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}

	// 在事务内部可以查询到数据
	row := tx.QueryRow("SELECT COUNT(*) FROM br_user WHERE user_name = ?", "txuser")
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Fatalf("查询用户数量失败: %v", err)
	}
	if count < 1 {
		t.Error("事务内部应该能查询到插入的数据")
	}

	// 回滚事务
	tx.Rollback()

	// 回滚后数据应该不存在
	row = db.QueryRow("SELECT COUNT(*) FROM br_user WHERE user_name = ?", "txuser")
	var userCount int
	err = row.Scan(&userCount)
	if err != nil {
		t.Fatalf("查询用户数量失败: %v", err)
	}

	if userCount != 0 {
		t.Errorf("回滚后用户应该不存在，但查询到 %d 个用户", userCount)
	}
}
