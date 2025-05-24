# === .env Loader (Подхватывает локальные переменные или CI Secrets) ===
include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

# === Settings ===
SQLC_BIN := ./bin/sqlc
MIGRATE_BIN := ./bin/migrate

DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

SQLC_VERSION := 1.29.0
MIGRATE_VERSION := v4.18.3

SQLC_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
SQLC_ARCH := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/arm64/arm64/')

# === Installers ===
bin/sqlc:
	@echo "⬇ Installing sqlc $(SQLC_VERSION) for $(SQLC_OS)-$(SQLC_ARCH)..."
	@mkdir -p bin
	@curl -sSL https://github.com/sqlc-dev/sqlc/releases/download/v$(SQLC_VERSION)/sqlc_$(SQLC_VERSION)_$(SQLC_OS)_$(SQLC_ARCH).tar.gz | tar -xz -C ./bin --strip-components=0
	@chmod +x ./bin/sqlc

bin/migrate:
	@echo "⬇ Installing migrate $(MIGRATE_VERSION)..."
	@mkdir -p ./bin
	@GO111MODULE=on GOBIN=$(shell pwd)/bin go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	@chmod +x ./bin/migrate

.PHONY: tools-install
tools-install: bin/sqlc bin/migrate
	@echo "✅ Tools installed: sqlc $(SQLC_VERSION) and migrate $(MIGRATE_VERSION)"

# === SQLC / Migrate targets ===
sqlc-generate: bin/sqlc
	$(SQLC_BIN) generate -f internal/storages/sqlc.yaml

migrate-up: bin/migrate
	$(MIGRATE_BIN) -path ./schema -database "$(DB_URL)" up

migrate-down: bin/migrate
	$(MIGRATE_BIN) -path ./schema -database "$(DB_URL)" down 1

migrate-status: bin/migrate
	$(MIGRATE_BIN) -path ./schema -database "$(DB_URL)" version

migrate-create: bin/migrate
	@if [ "$(name)" = "" ]; then echo "Usage: make migrate-create name=your_migration_name"; exit 1; fi
	$(MIGRATE_BIN) create -ext sql -dir ./schema $(name)