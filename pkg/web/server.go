package web

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func StartServer(db *sql.DB) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Wrap(GetLatestDigest(db)))
	mux.HandleFunc("/digest/", Wrap(GetDigest(db)))
	mux.HandleFunc("/story/", Wrap(ReadStory(db)))
	mux.HandleFunc("/archive/", Wrap(Archive(db)))
	return http.ListenAndServe(
		"localhost:8080",
		handlers.LoggingHandler(os.Stdout, mux),
	)
}
