package repo

import (
	"database/sql"
)

type DbConn interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

type scannable interface {
	Scan(dest ...interface{}) error
}
