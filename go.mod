module github.com/mhgbrg/hndaily

// +heroku goVersion go1.11
// +heroku install ./cmd/...

require (
	github.com/golang-migrate/migrate/v4 v4.2.5
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/sessions v1.1.3
	github.com/lib/pq v1.0.0
	github.com/pkg/errors v0.8.1
)
