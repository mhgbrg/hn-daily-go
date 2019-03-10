package web

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/pkg/repo"
)

func StartServer(config Config) error {
	gob.Register(&Flash{})

	db, err := repo.ConnectToDB(config.DatabaseURL)
	defer db.Close()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to connet to database"))
	}

	sessionStorage := CreateSessionStorage(db, config.CryptoKeys)

	templates, err := LoadTemplates()
	if err != nil {
		return errors.WithMessage(err, "failed to load templates")
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", Wrap(GetLatestDigest(templates, db, sessionStorage)))
	mux.HandleFunc("/digest/", Wrap(GetDigest(templates, db, sessionStorage)))
	mux.HandleFunc("/set-device-id", Wrap(SetDeviceID(templates, db, sessionStorage)))
	mux.HandleFunc("/story/", Wrap(ReadStory(db, sessionStorage)))
	mux.HandleFunc("/archive/", Wrap(Archive(templates, db)))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		handlers.LoggingHandler(os.Stdout, mux),
	)
}
