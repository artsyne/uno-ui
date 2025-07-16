package database

import (
	"os"
	"testing"

	testhelpers "github.com/biuren/hebe/tests/helpers"
)

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	// 设置使用真实数据库
	os.Setenv("USE_REAL_DB", "true")
	defer os.Unsetenv("USE_REAL_DB")

	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		// 检查数据库连接
		if suite.Database == nil {
			t.Fatal("数据库连接不应该为nil")
		}

		// 测试简单查询
		result := suite.Database.QueryRow(t, "SELECT 1 as test_value")
		var testValue int
		err := result.Scan(&testValue)
		if err != nil {
			t.Fatalf("查询测试失败: %v", err)
		}

		if testValue != 1 {
			t.Errorf("期望查询结果为1，实际得到: %d", testValue)
		}
	})
}

// TestUserTable 测试用户表操作
func TestUserTable(t *testing.T) {
	// 设置使用真实数据库
	os.Setenv("USE_REAL_DB", "true")
	defer os.Unsetenv("USE_REAL_DB")

	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		db := suite.Database

		// 开始事务
		db.BeginTransaction(t)

		// 插入测试用户
		result := db.ExecuteSQL(t, `
			INSERT INTO br_user (user_name, email, pwd, pwd_salt, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "testuser", "test@example.com", "hashed_password", "salt123", 1, 1642680000)

		// 检查插入结果
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("获取影响行数失败: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("期望插入1行，实际插入: %d", rowsAffected)
		}

		// 查询插入的用户
		row := db.QueryRow(t, "SELECT user_name, email FROM br_user WHERE user_name = ?", "testuser")
		var userName, email string
		err = row.Scan(&userName, &email)
		if err != nil {
			t.Fatalf("查询用户失败: %v", err)
		}

		if userName != "testuser" {
			t.Errorf("期望用户名为 'testuser'，实际得到: %s", userName)
		}
		if email != "test@example.com" {
			t.Errorf("期望邮箱为 'test@example.com'，实际得到: %s", email)
		}

		// 测试行数统计
		count := db.CountRows(t, "br_user")
		if count < 1 {
			t.Errorf("期望至少有1个用户，实际得到: %d", count)
		}

		// 回滚事务（清理测试数据）
		db.RollbackTransaction(t)

		// 验证回滚后数据已清理
		countAfterRollback := db.CountRows(t, "br_user")
		if countAfterRollback >= count {
			t.Error("回滚后用户数量没有减少")
		}
	})
}

// TestDiaryTable 测试日记表操作
func TestDiaryTable(t *testing.T) {
	// 设置使用真实数据库
	os.Setenv("USE_REAL_DB", "true")
	defer os.Unsetenv("USE_REAL_DB")

	testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
		db := suite.Database

		// 开始事务
		db.BeginTransaction(t)

		// 先插入用户（日记需要关联用户）
		db.ExecuteSQL(t, `
			INSERT INTO br_user (user_name, email, pwd, pwd_salt, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "diaryuser", "diary@example.com", "hashed_password", "salt123", 1, 1642680000)

		// 获取用户ID
		row := db.QueryRow(t, "SELECT user_id FROM br_user WHERE user_name = ?", "diaryuser")
		var userID int64
		err := row.Scan(&userID)
		if err != nil {
			t.Fatalf("获取用户ID失败: %v", err)
		}

		// 插入日记
		result := db.ExecuteSQL(t, `
			INSERT INTO br_diary (title, remark, add_utc_time, add_user_id, add_user_name, status) 
			VALUES (?, ?, ?, ?, ?, ?)
		`, "测试日记", "这是一个测试日记", 1642680000, userID, "diaryuser", 1)

		// 检查插入结果
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("获取影响行数失败: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("期望插入1行，实际插入: %d", rowsAffected)
		}

		// 查询日记
		diaryRow := db.QueryRow(t, "SELECT title, remark, add_user_name FROM br_diary WHERE add_user_id = ?", userID)
		var title, remark, addUserName string
		err = diaryRow.Scan(&title, &remark, &addUserName)
		if err != nil {
			t.Fatalf("查询日记失败: %v", err)
		}

		if title != "测试日记" {
			t.Errorf("期望日记标题为 '测试日记'，实际得到: %s", title)
		}
		if addUserName != "diaryuser" {
			t.Errorf("期望创建者为 'diaryuser'，实际得到: %s", addUserName)
		}

		// 统计日记数量
		diaryCount := db.CountRows(t, "br_diary")
		if diaryCount < 1 {
			t.Errorf("期望至少有1个日记，实际得到: %d", diaryCount)
		}

		// 测试断言方法
		db.AssertRowCount(t, "br_diary", diaryCount)

		// 回滚事务
		db.RollbackTransaction(t)
	})
}
