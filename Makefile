all: build

build:
	go build ./cmd/digest

clean:
	rm -f ./digest

digest: build
	./digest ${date}

create-migration:
	migrate create -ext sql -dir db/migrations ${name}

apply-migrations:
	migrate -source file://db/migrations -database "postgresql://hndaily@localhost/hndaily?sslmode=disable" up

rollback-migration:
	migrate -source file://db/migrations -database "postgresql://hndaily@localhost/hndaily?sslmode=disable" down 1

drop-db:
	migrate -source file://db/migrations -database "postgresql://hndaily@localhost/hndaily?sslmode=disable" drop
