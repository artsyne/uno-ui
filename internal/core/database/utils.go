package database

import (
	"database/sql"
)

// TableColumns 数据库表列信息结构体
// 用于存储MySQL表列的详细信息，对应SHOW COLUMNS命令的返回结果
type TableColumns struct {
	Field   string         `json:"Field" db:"Field"`     // 列名
	Type    string         `json:"Type" db:"Type"`       // 数据类型
	Null    string         `json:"Null" db:"Null"`       // 是否允许NULL值
	Key     string         `json:"Key" db:"Key"`         // 键类型(PRI/UNI/MUL等)
	Default sql.NullString `json:"Default" db:"Default"` // 默认值
	Extra   string         `json:"Extra" db:"Extra"`     // 额外信息(auto_increment等)
}

// Tables 数据库表信息结构体
// 用于存储数据库中表的名称信息，对应SHOW TABLES命令的返回结果
type Tables struct {
	TableName string `json:"Tables_in_cd" db:"Tables_in_cd"` // 表名
}

// GetTables 获取所有数据库表名称
// 执行SHOW TABLES命令获取当前数据库中的所有表名
//
// 返回:
//   - []Tables: 表信息列表，如果没有表则返回nil
//   - error: 执行失败时返回错误，成功时返回nil
func GetTables() ([]Tables, error) {
	tables := []Tables{}
	sqlStr := "SHOW TABLES"
	err := db.Select(&tables, sqlStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tables, err
}

// GetTableColumns 根据表名获取表的列信息
// 执行SHOW COLUMNS FROM命令获取指定表的所有列详细信息
//
// 参数:
//   - tableName: 要查询的表名
//
// 返回:
//   - []TableColumns: 列信息列表，如果表不存在或无列则返回nil
//   - error: 执行失败时返回错误，成功时返回nil
func GetTableColumns(tableName string) ([]TableColumns, error) {
	columns := []TableColumns{}
	sqlStr := "SHOW COLUMNS FROM " + tableName
	err := db.Select(&columns, sqlStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return columns, err
}
