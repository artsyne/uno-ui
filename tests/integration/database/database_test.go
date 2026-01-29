package database

import (
	"testing"

	framework "github.com/biuren/hebe/tests/framework"
)

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	// 检查数据库连接
	dbManager := ctx.GetDatabase()
	if dbManager == nil {
		t.Fatal("数据库连接不应该为nil")
	}

	// 测试简单查询
	row := dbManager.GetDB().QueryRow("SELECT 1 as test_value")
	var testValue int
	err := row.Scan(&testValue)
	if err != nil {
		t.Fatalf("查询测试失败: %v", err)
	}

	if testValue != 1 {
		t.Errorf("期望查询结果为1，实际得到: %d", testValue)
	}
}

// TestUserTable 测试用户表操作
func TestUserTable(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	dbManager := ctx.GetDatabase()
	db := dbManager.GetDB()

		// 清理可能存在的测试数据
	db.Exec("DELETE FROM br_user WHERE handle = ?", "testuser")

	// 插入测试用户
	result, err := db.Exec(`
		INSERT INTO br_user (handle, email, pwd, status, created_utc_time) 
		VALUES (?, ?, ?, ?, ?)
	`, "testuser", "test@example.com", "hashed_password", 1, 1642680000)
	if err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}

	// 检查插入结果
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("获取影响行数失败: %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("期望插入1行，实际插入: %d", rowsAffected)
	}

	// 查询插入的用户
	row := db.QueryRow("SELECT handle, email FROM br_user WHERE handle = ?", "testuser")
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
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM br_user").Scan(&count)
	if err != nil {
		t.Fatalf("统计用户数量失败: %v", err)
	}
	if count < 1 {
		t.Errorf("期望至少有1个用户，实际得到: %d", count)
	}

	// 清理测试数据
	_, err = db.Exec("DELETE FROM br_user WHERE handle = ?", "testuser")
	if err != nil {
		t.Fatalf("清理测试数据失败: %v", err)
	}
}

// TestDiaryTable 测试日记表操作
func TestDiaryTable(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	dbManager := ctx.GetDatabase()
	db := dbManager.GetDB()

	// 清理可能存在的测试数据
	db.Exec("DELETE FROM br_user WHERE handle = ?", "diaryuser")

	// 先插入用户（日记需要关联用户）
		_, err := db.Exec(`
			INSERT INTO br_user (handle, email, pwd, status, created_utc_time) 
			VALUES (?, ?, ?, ?, ?)
		`, "diaryuser", "diary@example.com", "hashed_password", 1, 1642680000)
		if err != nil {
			t.Fatalf("插入用户失败: %v", err)
		}

		// 获取用户ID
		row := db.QueryRow("SELECT user_id FROM br_user WHERE handle = ?", "diaryuser")
		var userID int64
		err = row.Scan(&userID)
		if err != nil {
			t.Fatalf("获取用户ID失败: %v", err)
		}

		// 插入日记
		result, err := db.Exec(`
			INSERT INTO br_diary (title, remark, add_utc_time, add_user_id, status) 
			VALUES (?, ?, ?, ?, ?)
		`, "测试日记", "这是一个测试日记", 1642680000, userID, 1)
		if err != nil {
			t.Fatalf("插入日记失败: %v", err)
		}

		// 检查插入结果
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			t.Fatalf("获取影响行数失败: %v", err)
		}
		if rowsAffected != 1 {
			t.Errorf("期望插入1行，实际插入: %d", rowsAffected)
		}

		// 查询日记
		diaryRow := db.QueryRow("SELECT title, remark, add_user_id FROM br_diary WHERE add_user_id = ?", userID)
		var title, remark string
		var dbUserID int64
		err = diaryRow.Scan(&title, &remark, &dbUserID)
		if err != nil {
			t.Fatalf("查询日记失败: %v", err)
		}

		if title != "测试日记" {
			t.Errorf("期望日记标题为 '测试日记'，实际得到: %s", title)
		}
		if dbUserID != userID {
			t.Errorf("期望创建者ID为 %d，实际得到: %d", userID, dbUserID)
		}

		// 统计日记数量
		var diaryCount int
		err = db.QueryRow("SELECT COUNT(*) FROM br_diary").Scan(&diaryCount)
		if err != nil {
			t.Fatalf("统计日记数量失败: %v", err)
		}
		if diaryCount < 1 {
			t.Errorf("期望至少有1个日记，实际得到: %d", diaryCount)
		}

		// 清理测试数据
		_, err = db.Exec("DELETE FROM br_diary WHERE add_user_id = ?", userID)
		if err != nil {
			t.Fatalf("清理日记数据失败: %v", err)
		}
		_, err = db.Exec("DELETE FROM br_user WHERE handle = ?", "diaryuser")
		if err != nil {
			t.Fatalf("清理用户数据失败: %v", err)
		}
}
