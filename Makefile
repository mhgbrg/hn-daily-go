.PHONY: all
all: build

# --- building ---
.PHONY: build_digest
build_digest:
	go build ./cmd/digest

.PHONY: build_server
build_server:
	go build ./cmd/server

.PHONY: clean
clean:
	rm -f ./digest ./server

# --- digesting ---
.PHONY: digest
digest: build_digest
	./digest ${date} ${start_date} ${end_date}

# --- server ---
.PHONY: serve
serve: build_server
	./server

.PHONY: watch_server
watch_server:
	ag -l -u | entr -r make serve

# --- db management ---
.PHONY: create_migration
create_migration:
	migrate create -ext sql -dir db/migrations ${name}

.PHONY: apply_migrations
apply_migrations:
	migrate -source file://db/migrations -database ${DATABASE_URL} up

.PHONY: rollback_migration
rollback_migration:
	migrate -source file://db/migrations -database ${DATABASE_URL} down 1

.PHONY: drop_db
drop_db:
	migrate -source file://db/migrations -database ${DATABASE_URL} drop
