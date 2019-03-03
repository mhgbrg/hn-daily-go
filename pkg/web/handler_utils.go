package web

import (
	"net/http"
)

type CustomHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Wrap(handler CustomHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			HandleError(err, w, r)
		}
	}
}
