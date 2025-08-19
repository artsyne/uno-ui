package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/biuren/hebe/internal/core/database"
)

// MockDatabase 模拟数据库连接用于测试
type MockDatabase struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

// setupMockDB 设置模拟数据库
func setupMockDB(t *testing.T) *MockDatabase {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// 设置模拟数据库连接到 database 包
	sqlxDB := sqlx.NewDb(db, "mysql")
	database.SetDBForTesting(sqlxDB)

	return &MockDatabase{
		db:   db,
		mock: mock,
	}
}

// cleanup 清理模拟数据库
func (m *MockDatabase) cleanup() {
	m.db.Close()
}

// TestGetTables 测试获取所有数据库表名称
func TestGetTables(t *testing.T) {
	mockDB := setupMockDB(t)
	defer mockDB.cleanup()

	// 测试成功获取表列表
	t.Run("成功获取表列表", func(t *testing.T) {
		// 设置期望的查询和结果
		rows := sqlmock.NewRows([]string{"Tables_in_cd"}).
			AddRow("br_user").
			AddRow("br_diary").
			AddRow("br_post")

		mockDB.mock.ExpectQuery("SHOW TABLES").WillReturnRows(rows)

		// 执行测试
		tables, err := database.GetTables()

		// 验证结果
		assert.NoError(t, err)
		assert.Len(t, tables, 3)
		assert.Equal(t, "br_user", tables[0].TableName)
		assert.Equal(t, "br_diary", tables[1].TableName)
		assert.Equal(t, "br_post", tables[2].TableName)

		// 验证所有期望都被满足
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})

	// 测试没有表的情况
	t.Run("没有表的情况", func(t *testing.T) {
		mockDB.mock.ExpectQuery("SHOW TABLES").WillReturnError(sql.ErrNoRows)

		tables, err := database.GetTables()

		assert.NoError(t, err)
		assert.Nil(t, tables)
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})

	// 测试数据库错误
	t.Run("数据库错误", func(t *testing.T) {
		mockDB.mock.ExpectQuery("SHOW TABLES").WillReturnError(sql.ErrConnDone)

		tables, err := database.GetTables()

		assert.Error(t, err)
		assert.Empty(t, tables)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})
}

// TestGetTableColumns 测试根据表名获取表的列信息
func TestGetTableColumns(t *testing.T) {
	mockDB := setupMockDB(t)
	defer mockDB.cleanup()

	// 测试成功获取列信息
	t.Run("成功获取列信息", func(t *testing.T) {
		// 设置期望的查询和结果
		rows := sqlmock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
			AddRow("id", "int(11)", "NO", "PRI", sql.NullString{}, "auto_increment").
			AddRow("name", "varchar(255)", "NO", "", sql.NullString{String: "default_name", Valid: true}, "").
			AddRow("email", "varchar(255)", "YES", "UNI", sql.NullString{}, "")

		mockDB.mock.ExpectQuery("SHOW COLUMNS FROM br_user").WillReturnRows(rows)

		// 执行测试
		columns, err := database.GetTableColumns("br_user")

		// 验证结果
		assert.NoError(t, err)
		assert.Len(t, columns, 3)

		// 验证第一列
		assert.Equal(t, "id", columns[0].Field)
		assert.Equal(t, "int(11)", columns[0].Type)
		assert.Equal(t, "NO", columns[0].Null)
		assert.Equal(t, "PRI", columns[0].Key)
		assert.False(t, columns[0].Default.Valid)
		assert.Equal(t, "auto_increment", columns[0].Extra)

		// 验证第二列
		assert.Equal(t, "name", columns[1].Field)
		assert.Equal(t, "varchar(255)", columns[1].Type)
		assert.Equal(t, "NO", columns[1].Null)
		assert.Equal(t, "", columns[1].Key)
		assert.True(t, columns[1].Default.Valid)
		assert.Equal(t, "default_name", columns[1].Default.String)
		assert.Equal(t, "", columns[1].Extra)

		// 验证第三列
		assert.Equal(t, "email", columns[2].Field)
		assert.Equal(t, "varchar(255)", columns[2].Type)
		assert.Equal(t, "YES", columns[2].Null)
		assert.Equal(t, "UNI", columns[2].Key)
		assert.False(t, columns[2].Default.Valid)
		assert.Equal(t, "", columns[2].Extra)

		// 验证所有期望都被满足
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})

	// 测试表不存在的情况
	t.Run("表不存在的情况", func(t *testing.T) {
		mockDB.mock.ExpectQuery("SHOW COLUMNS FROM nonexistent_table").WillReturnError(sql.ErrNoRows)

		columns, err := database.GetTableColumns("nonexistent_table")

		assert.NoError(t, err)
		assert.Nil(t, columns)
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})

	// 测试数据库错误
	t.Run("数据库错误", func(t *testing.T) {
		mockDB.mock.ExpectQuery("SHOW COLUMNS FROM br_user").WillReturnError(sql.ErrConnDone)

		columns, err := database.GetTableColumns("br_user")

		assert.Error(t, err)
		assert.Empty(t, columns)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})

	// 测试空表名
	t.Run("空表名", func(t *testing.T) {
		mockDB.mock.ExpectQuery("SHOW COLUMNS FROM ").WillReturnError(sql.ErrNoRows)

		columns, err := database.GetTableColumns("")

		assert.NoError(t, err)
		assert.Nil(t, columns)
		assert.NoError(t, mockDB.mock.ExpectationsWereMet())
	})
}

// TestTableColumns 测试 TableColumns 结构体
func TestTableColumns(t *testing.T) {
	// 测试结构体字段
	column := database.TableColumns{
		Field:   "test_field",
		Type:    "varchar(255)",
		Null:    "NO",
		Key:     "PRI",
		Default: sql.NullString{String: "default_value", Valid: true},
		Extra:   "auto_increment",
	}

	assert.Equal(t, "test_field", column.Field)
	assert.Equal(t, "varchar(255)", column.Type)
	assert.Equal(t, "NO", column.Null)
	assert.Equal(t, "PRI", column.Key)
	assert.True(t, column.Default.Valid)
	assert.Equal(t, "default_value", column.Default.String)
	assert.Equal(t, "auto_increment", column.Extra)
}

// TestTables 测试 Tables 结构体
func TestTables(t *testing.T) {
	// 测试结构体字段
	table := database.Tables{
		TableName: "br_user",
	}

	assert.Equal(t, "br_user", table.TableName)
}

// TestDatabaseUtilsIntegration 测试数据库工具函数的集成
func TestDatabaseUtilsIntegration(t *testing.T) {
	mockDB := setupMockDB(t)
	defer mockDB.cleanup()

	// 模拟获取表列表
	tableRows := sqlmock.NewRows([]string{"Tables_in_cd"}).
		AddRow("br_user").
		AddRow("br_diary")

	mockDB.mock.ExpectQuery("SHOW TABLES").WillReturnRows(tableRows)

	// 模拟获取 br_user 表的列信息
	userColumnRows := sqlmock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
		AddRow("id", "int(11)", "NO", "PRI", sql.NullString{}, "auto_increment").
		AddRow("username", "varchar(255)", "NO", "UNI", sql.NullString{}, "")

	mockDB.mock.ExpectQuery("SHOW COLUMNS FROM br_user").WillReturnRows(userColumnRows)

	// 执行集成测试
	tables, err := database.GetTables()
	assert.NoError(t, err)
	assert.Len(t, tables, 2)

	columns, err := database.GetTableColumns("br_user")
	assert.NoError(t, err)
	assert.Len(t, columns, 2)
	assert.Equal(t, "id", columns[0].Field)
	assert.Equal(t, "username", columns[1].Field)

	// 验证所有期望都被满足
	assert.NoError(t, mockDB.mock.ExpectationsWereMet())
}