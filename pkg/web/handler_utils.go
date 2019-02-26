package web

import (
	"fmt"
	"log"
	"net/http"
)

type CustomHandlerFunc func(w http.ResponseWriter, r *http.Request) (fmt.Stringer, error)

func Wrap(handler CustomHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := handler(w, r)
		if err != nil {
			log.Printf("%+v", err)
			switch err := err.(type) {
			case HTTPError:
				switch err.Code {
				case 404:
					http.NotFound(w, r)
				default:
					http.Error(w, err.Error(), err.Code)
				}
			default:
				http.Error(w, err.Error(), 500)
			}
		} else if res != nil {
			fmt.Fprintf(w, res.String())
		}
	}
}
