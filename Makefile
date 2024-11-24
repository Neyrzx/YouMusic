# Configuration variables
GOBIN_PROJECT_RELATIVE_PATH=bin/
GOBIN_PROJECT_ABSOLUTE_PATH=$(CURDIR)/${GOBIN_PROJECT_RELATIVE_PATH}

SWAG_VERSION=v1.4.1
SWAG_DIRECTORIES=cmd/app
SWAG_OUTPUT=$(CURDIR)/docs/swagger

GOLANGCI_LINT_VERSION=v1.62.0
GOLANGCI_LINT_DIRECTORIES=./...

GOIMPORTS_VERSION=latest


# Установка всех нужных проекту тулов
deps: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
deps:
	@echo Installing Tools...

	@echo Installing Swagger...
	@go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}

	@echo Installing GolangCI-Lint...
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT_VERSION}

	@echo Installing GoImports...
	@go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION}
.PHONY: deps


# Генерация swagger спецификации
swagger: export GOBIN=${GOBIN_PROJECT_ABSOLUTE_PATH}
swagger:
	@echo Generating Swagger Docs...
	@swag init --dir ${SWAG_DIRECTORIES} 
.PHONY: swagger


# Основные рецепты
init: swagger
init: deps
	@go mod tidy
.PHONY: init


lint:
	@${GOBIN_PROJECT_ABSOLUTE_PATH}/golangci-lint run ${GOLANGCI_LINT_DIRECTORIES}
.PHONY: lint
