package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite(dataSourceName string) (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}
