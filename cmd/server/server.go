package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq" // load PSQL driver

	"github.com/mhgbrg/hndaily/pkg/handlers"
)

func main() {
	// TODO: Read database info from environment.
	db, err := sql.Open("postgres", "user=hndaily dbname=hndaily sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/digest/", handlers.GetDigest(db))
	http.HandleFunc("/story/", handlers.ReadStory(db))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
