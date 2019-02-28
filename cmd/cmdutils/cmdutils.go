package cmdutils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Load PSQL driver.

	"github.com/pkg/errors"
)

func ConnectToDB() *sql.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal(errors.New("failed to read database connection string from environment variable DATABASE_URL"))
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to open connection to database at %s", databaseURL))
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to ping database at %s", databaseURL))
	}

	return db
}
