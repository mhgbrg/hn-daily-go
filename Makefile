all: build

# --- building ---

build_digest:
	go build ./cmd/digest

build_server:
	go build ./cmd/server

clean:
	rm -f ./digest ./server

# --- digesting ---
digest: build_digest
	./digest ${date} ${start_date} ${end_date}

# --- server ---
serve: build_server
	./server

watch_server:
	ag -l -u | entr -r make serve

# --- db management ---
create_migration:
	migrate create -ext sql -dir db/migrations ${name}

apply_migrations:
	migrate -source file://db/migrations -database ${DATABASE_URL} up

rollback_migration:
	migrate -source file://db/migrations -database ${DATABASE_URL} down 1

drop_db:
	migrate -source file://db/migrations -database ${DATABASE_URL} drop
