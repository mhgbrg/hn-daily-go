package web

import (
	"fmt"
	"log"
	"net/http"
)

type HTTPError struct {
	Err  error
	Code int
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("status %d: %+v", err.Code, err.Err)
}

func NotFoundError(err error) HTTPError {
	return HTTPError{err, http.StatusNotFound}
}

func InternalServerError(err error) HTTPError {
	return HTTPError{err, http.StatusInternalServerError}
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", err)
	switch err := err.(type) {
	case HTTPError:
		switch err.Code {
		case 404:
			http.NotFound(w, r)
		default:
			http.Error(w, "internal server error", 500)
		}
	default:
		http.Error(w, "internal server error", 500)
	}
}
