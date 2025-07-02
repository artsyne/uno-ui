package database

import (
	"database/sql"
)

// TableColumns 数据库表列信息结构体
type TableColumns struct {
	Field   string         `json:"Field" db:"Field"`
	Type    string         `json:"Type" db:"Type"`
	Null    string         `json:"Null" db:"Null"`
	Key     string         `json:"Key" db:"Key"`
	Default sql.NullString `json:"Default" db:"Default"`
	Extra   string         `json:"Extra" db:"Extra"`
}

// Tables 数据库表信息结构体
type Tables struct {
	TableName string `json:"Tables_in_cd" db:"Tables_in_cd"`
}

// GetTables 获取所有数据库表名称
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
func GetTableColumns(tableName string) ([]TableColumns, error) {
	columns := []TableColumns{}
	sqlStr := "SHOW COLUMNS FROM " + tableName
	err := db.Select(&columns, sqlStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return columns, err
}
