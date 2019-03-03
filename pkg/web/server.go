package web

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
)

func StartServer(addr string, db *sql.DB) error {
	templates, err := LoadTemplates()
	if err != nil {
		return errors.WithMessage(err, "failed to load templates")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", Wrap(GetLatestDigest(templates, db)))
	mux.HandleFunc("/digest/", Wrap(GetDigest(templates, db)))
	mux.HandleFunc("/story/", Wrap(ReadStory(db)))
	mux.HandleFunc("/archive/", Wrap(Archive(templates, db)))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	return http.ListenAndServe(
		addr,
		handlers.LoggingHandler(os.Stdout, mux),
	)
}
