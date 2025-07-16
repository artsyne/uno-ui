package database

import (
	"testing"

	testhelpers "github.com/biuren/hebe/tests/helpers"
)

// TestDatabaseOnly 仅使用数据库的测试示例
func TestDatabaseOnly(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 测试数据库基本功能
		result := db.QueryRow(t, "SELECT DATABASE() as current_db")
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
	})
}

// TestTableStructure 测试表结构
func TestTableStructure(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 检查关键表是否存在
		tables := []string{"br_user", "br_diary", "br_comment", "br_follow"}

		for _, tableName := range tables {
			result := db.QueryRow(t, `
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
	})
}

// TestUserTableConstraints 测试用户表约束
func TestUserTableConstraints(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务
		db.BeginTransaction(t)
		defer db.RollbackTransaction(t)

		// 插入正常用户
		db.ExecuteSQL(t, `
			INSERT INTO br_user (user_name, email, pwd, pwd_salt, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "uniqueuser", "unique@example.com", "password", "salt", 1, 1642680000)

		// 测试用户名唯一约束
		_, err := db.GetDB().Exec(`
			INSERT INTO br_user (user_name, email, pwd, pwd_salt, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "uniqueuser", "another@example.com", "password", "salt", 1, 1642680000)

		if err == nil {
			t.Error("期望用户名重复时报错，但没有报错")
		}

		// 测试邮箱唯一约束
		_, err = db.GetDB().Exec(`
			INSERT INTO br_user (user_name, email, pwd, pwd_salt, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "anotheruser", "unique@example.com", "password", "salt", 1, 1642680000)

		if err == nil {
			t.Error("期望邮箱重复时报错，但没有报错")
		}
	})
}

// TestTransactionIsolation 测试事务隔离
func TestTransactionIsolation(t *testing.T) {
	testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
		// 开始事务1
		db.BeginTransaction(t)

		// 在事务中插入数据
		db.ExecuteSQL(t, `
			INSERT INTO br_user (user_name, email, pwd, pwd_salt, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "txuser", "tx@example.com", "password", "salt", 1, 1642680000)

		// 在事务内部可以查询到数据
		count := db.CountRows(t, "br_user")
		if count < 1 {
			t.Error("事务内部应该能查询到插入的数据")
		}

		// 回滚事务
		db.RollbackTransaction(t)

		// 回滚后数据应该不存在
		row := db.QueryRow(t, "SELECT COUNT(*) FROM br_user WHERE user_name = ?", "txuser")
		var userCount int
		err := row.Scan(&userCount)
		if err != nil {
			t.Fatalf("查询用户数量失败: %v", err)
		}

		if userCount != 0 {
			t.Errorf("回滚后用户应该不存在，但查询到 %d 个用户", userCount)
		}
	})
}
