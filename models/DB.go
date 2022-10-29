package models

import (
	"database/sql"
)

// 创建Struct code snap
type TableColumns struct {
	Field   string `json:"Field" db:"Field"`
	Type    string `json:"Type" db:"Type"`
	Null    string `json:"Null" db:"Null"`
	Key     string `json:"Key" db:"Key"`
	Default sql.NullString `json:"Default" db:"Default"`
	Extra   string `json:"Extra" db:"Extra"`
}

type Tables struct{
	TableName string 		`json:"Tables_in_cd" db:"Tables_in_cd"`
}

//获取所有数据库名称
func GetTables() ([]Tables,error) {
	tables := []Tables{}
	sqlStr := "show tables"
	err:= db.Select(&tables,sqlStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tables, err
}

//根据表名返回列
func GetTableColumns(tableName string) ([]TableColumns, error) {
	columns := []TableColumns{}
	sqlstr := "show columns from " + tableName
	err := db.Select(&columns, sqlstr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return columns, err
}
