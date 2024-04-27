.PHONY: migrate_force migrate_version migrate_down migrate_up new_migration local run build test

# ==============================================================================
# Main

run:
	go run ./cmd/shortify/main.go

build:
	go build ./cmd/shortify/main.go