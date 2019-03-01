package repo

import (
	"database/sql"

	_ "github.com/lib/pq" // Load PSQL driver.

	"github.com/pkg/errors"
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

func ConnectToDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open connection to database at %s", url)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ping database at %s", url)
	}

	return db, nil
}
