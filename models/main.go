package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

//SetDB establishes connection to database and saves its handler into db *sqlx.DB
func SetDB(connection string) {
	db = sqlx.MustConnect("mysql", connection)
}

//GetDB returns database handler
func GetDB() *sqlx.DB {
	return db
}
