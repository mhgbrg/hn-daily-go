package web

import (
	"fmt"
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
