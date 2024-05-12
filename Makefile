.PHONY: migrate_force migrate_version migrate_down migrate_up new_migration local run build test

# ==============================================================================
# Main

run:
	go run main.go serve

build:
	go build -o shortify ./main.go

test:
	go test -cover ./...

bench:
	go test -bench=. ./...