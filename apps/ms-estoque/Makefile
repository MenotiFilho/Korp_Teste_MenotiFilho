MIGRATE_IMAGE ?= migrate/migrate:v4.18.3
DB_URL ?= postgres://postgres:postgres@localhost:5433/estoque?sslmode=disable
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
MIGRATIONS_DIR := $(MAKEFILE_DIR)migrations

.PHONY: migrate-up migrate-down migrate-version migrate-force migrate-create

migrate-up:
	docker run --rm -v $(MIGRATIONS_DIR):/migrations --network host $(MIGRATE_IMAGE) -path=/migrations -database "$(DB_URL)" up

migrate-down:
	docker run --rm -v $(MIGRATIONS_DIR):/migrations --network host $(MIGRATE_IMAGE) -path=/migrations -database "$(DB_URL)" down 1

migrate-version:
	docker run --rm -v $(MIGRATIONS_DIR):/migrations --network host $(MIGRATE_IMAGE) -path=/migrations -database "$(DB_URL)" version

migrate-force:
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required. Example: make migrate-force VERSION=1"; exit 1; fi
	docker run --rm -v $(MIGRATIONS_DIR):/migrations --network host $(MIGRATE_IMAGE) -path=/migrations -database "$(DB_URL)" force $(VERSION)

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "NAME is required. Example: make migrate-create NAME=create_table_x"; exit 1; fi
	docker run --rm -v $(MIGRATIONS_DIR):/migrations $(MIGRATE_IMAGE) create -ext sql -dir /migrations -seq $(NAME)
