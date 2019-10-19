module github.com/mhgbrg/hndaily

// +heroku goVersion go1.11
// +heroku install ./cmd/...

go 1.11

require (
	github.com/golang-migrate/migrate/v4 v4.2.5
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/sessions v1.1.3
	github.com/lib/pq v1.0.0
	github.com/pkg/errors v0.8.1
	golang.org/x/crypto v0.0.0-20190228161510-8dd112bcdc25
)
