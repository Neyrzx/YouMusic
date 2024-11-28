include .env

# Configuration variables
GOBIN_PROJECT_RELATIVE_PATH=bin
GOBIN_PROJECT_ABSOLUTE_PATH=$(CURDIR)/${GOBIN_PROJECT_RELATIVE_PATH}

# Tools versions
VERSION_SWAG = v1.4.1
VERSION_GOLANGCI_LINT = v1.62.0
VERSION_MIGRATE = v4.18.1
VERSION_MOCKERY = v2.49.1

.DEFAULT_GOAL := init

# swag
# https://github.com/swaggo/swag
install-swag-cli: 
	@echo Installing Swag CLI...
	@go install github.com/swaggo/swag/cmd/swag@${VERSION_SWAG}
.PHONY: install-swag-cli

# golangci-lint
# https://github.com/golangci/golangci-lint
install-golangcilint-cli:
	@echo Installing Golangci-lint CLI...
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${VERSION_GOLANGCI_LINT}
.PHONY: install-golangcilint-cli

# golang-migrate
# https://github.com/golang-migrate/migrate
install-migrate-cli: export OS=$(shell uname -s | tr A-Z a-z)
install-migrate-cli: export ARCH=amd64
install-migrate-cli:
	@echo Installing Migrate CLI...
	@curl -L https://github.com/golang-migrate/migrate/releases/download/${VERSION_MIGRATE}/migrate.${OS}-${ARCH}.tar.gz | tar xvz -C ${GOBIN_PROJECT_ABSOLUTE_PATH}
.PHONY: install-migrate-cli

# mockery
# https://vektra.github.io/mockery/latest/
install-mockery-cli:
	@echo Installing Mockery CLI...
	@go install github.com/vektra/mockery/v2@${VERSION_MOCKERY}
.PHONY: install-mockery-cli

# Установка всех нужных проекту тулов
install: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
install: install-migrate-cli install-golangcilint-cli install-swag-cli install-mockery-cli
	@echo Tools installed successful.
.PHONY: install

# Генерация swagger спецификации
swagger: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
swagger:
	@echo Generating Swagger Docs...
	@swag init -g internal/delivery/rest/api.go
.PHONY: swagger

# Установка всех зависимостей проекта, генерация swagger документации.
init: swagger install
	@go mod tidy
.PHONY: init

# Mocks
mocks:
	@${GOBIN_PROJECT_ABSOLUTE_PATH}/mockery
.PHONY: mocks

# Linter
lint: export DIRECTORIES=./...
lint:
	@${GOBIN_PROJECT_ABSOLUTE_PATH}/golangci-lint run ${DIRECTORIES}
.PHONY: lint

# Tests
test:
	go test -timeout=30s ./... -coverprofile=coverage.out
.PHONY: test

# Tests
test-race:
	go test -race -timeout=30s ./... -coverprofile=coverage.out
.PHONY: test-race

# Tests coverage
test-coverage:
	@go tool cover -html=coverage.out
.PHONY: test-coverage

# Dev environment
compose-dev:
	@docker compose -f deployments/compose.dev.yaml --env-file=.env up -d
.PHONY: compose-dev

compose-dev-clean:
	@docker compose -f deployments/compose.dev.yaml --env-file=.env down -v
.PHONY: compose-dev-clean

# Migrations recipies.
MIGRATIONS_DIR=${MIGRATIONS_SOURCE}
MIGRATIONS_DB_URL="${DB_PROTO}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

migration-create:
	@read -p "Migration name: " name && \
    ${GOBIN_PROJECT_RELATIVE_PATH}/migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name
.PHONY: migration-create

# Apply all or N up migrations.
migration-up: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
migration-up:
	@read -p "Count to up migrations (empty to apply all): " count && \
	${GOBIN_PROJECT_RELATIVE_PATH}/migrate -database $(MIGRATIONS_DB_URL) -path $(MIGRATIONS_DIR) up $$count
.PHONY: migration-up

# Apply all or N down migrations.
migration-down: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
migration-down:
	@read -p "Count to rollback (empty to rollback all): " count && \
	${GOBIN_PROJECT_RELATIVE_PATH}/migrate -database $(MIGRATIONS_DB_URL) -path $(MIGRATIONS_DIR) down $$count && echo $$count
.PHONY: migration-down

migration: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
migration:
	@read -p "Command: " cmd && \
	${GOBIN_PROJECT_RELATIVE_PATH}/migrate -database $(MIGRATIONS_DB_URL) -path $(MIGRATIONS_DIR) $$cmd
.PHONY: migration
