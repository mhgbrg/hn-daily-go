package cmdutils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Load PSQL driver.

	"github.com/pkg/errors"
)

func ConnectToDB() *sql.DB {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal(errors.New("failed to read database connection string from environment variable DB_URL"))
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to open connection to database at %s", connectionString))
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to ping database at %s", connectionString))
	}

	return db
}
