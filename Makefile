.PHONY: migrate_force migrate_version migrate_down migrate_up new_migration local run build test

# ==============================================================================
# Main

run:
	go run ./cmd/shortify/main.go serve

migrate:
	go run ./cmd/shortify/main.go migrate

build:
	go build -o shortify ./cmd/shortify/main.go

test:
	go test -cover ./...

bench:
	go test -bench=. ./...