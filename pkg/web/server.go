package web

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/pkg/repo"
)

func StartServer(config Config) error {
	gob.Register(&Flash{})

	db, err := repo.ConnectToDB(config.DatabaseURL)
	defer db.Close()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to connect to database"))
	}

	sessionStorage := CreateSessionStorage(db, config.CryptoKeys)

	templates, err := LoadTemplates()
	if err != nil {
		return errors.WithMessage(err, "failed to load templates")
	}

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", Wrap(GetLatestDigest(templates, db, sessionStorage)))
	router.HandleFunc("/digest/{date}", Wrap(GetDigest(templates, db, sessionStorage)))
	router.HandleFunc("/set-device-id", Wrap(SetDeviceID(db, sessionStorage)))
	router.HandleFunc("/story/{id}", Wrap(ReadStory(db, sessionStorage)))
	router.HandleFunc("/story/{id}/mark-as-read", Wrap(MarkStoryAsRead(db, sessionStorage)))
	router.HandleFunc("/archive/{yearMonth}", Wrap(Archive(templates, db)))
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		handlers.LoggingHandler(os.Stdout, router),
	)
}
