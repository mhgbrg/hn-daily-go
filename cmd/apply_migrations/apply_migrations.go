package main

import (
	"log"
	"os"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/pkg/errors"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	m, err := migrate.New(
		"file://db/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create migrate instance"))
	}
	err = m.Up()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to apply migrations"))
	}
	log.Print("applied migrations")
}
