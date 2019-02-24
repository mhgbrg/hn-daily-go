all: build

# building

build_digest:
	go build ./cmd/digest

build_server:
	go build ./cmd/server

clean:
	rm -f ./digest

# digesting

digest: build_digest
	./digest ${date} ${start_date} ${end_date}

# server

serve: build_server
	./server

# db management

create-migration:
	migrate create -ext sql -dir db/migrations ${name}

apply-migrations:
	migrate -source file://db/migrations -database "postgresql://hndaily@localhost/hndaily?sslmode=disable" up

rollback-migration:
	migrate -source file://db/migrations -database "postgresql://hndaily@localhost/hndaily?sslmode=disable" down 1

drop-db:
	migrate -source file://db/migrations -database "postgresql://hndaily@localhost/hndaily?sslmode=disable" drop
