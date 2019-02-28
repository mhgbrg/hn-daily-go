package web

import (
	"fmt"
	"net/http"
)

type CustomHandlerFunc func(w http.ResponseWriter, r *http.Request) (fmt.Stringer, error)

func Wrap(handler CustomHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := handler(w, r)
		if err != nil {
			HandleError(err, w, r)
		} else if res != nil {
			fmt.Fprintf(w, res.String())
		}
	}
}
