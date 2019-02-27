package web

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	_ "github.com/lib/pq" // Load PostgreSQL driver.
)

func StartServer() {
	// TODO: Read database info from environment.
	db, err := sql.Open("postgres", "user=hndaily dbname=hndaily sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", Wrap(GetLatestDigest(db)))
	mux.HandleFunc("/digest/", Wrap(GetDigest(db)))
	mux.HandleFunc("/story/", Wrap(ReadStory(db)))
	log.Fatal(
		http.ListenAndServe(
			"localhost:8080",
			handlers.LoggingHandler(os.Stdout, mux),
		),
	)
}
