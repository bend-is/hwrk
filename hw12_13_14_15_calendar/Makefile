BIN := "./bin/calendar"
DOCKER_IMG = "calendar:develop"
DC_FILE = "./deployments/docker-compose.yaml"
DB_DSN = "postgres://forge:secret@postgres:5432/app?sslmode=disable"

INPUT = $(filter-out $@, $(MAKECMDGOALS))
GIT_HASH := $(shell git log --format="%h" -n 1)
DC_RUN_ARGS = --rm --no-deps --user "$(shell id -u):$(shell id -g)"
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up: ## Up all containers
	docker-compose -f $(DC_FILE) up -d

migrate-create: ## Create migration
	docker-compose -f $(DC_FILE) run $(DC_RUN_ARGS) migrator create -ext sql -dir ./migrations/postgres -seq $(INPUT)

migrate: ## Migrate database
	docker-compose -f $(DC_FILE) run $(DC_RUN_ARGS) migrator -path ./migrations/postgres -database $(DB_DSN) up

migrate-down: ## Rollback n migrations or all
	docker-compose -f $(DC_FILE) run $(DC_RUN_ARGS) migrator -path ./migrations/postgres -database $(DB_DSN) down $(INPUT)

.PHONY: build run build-img run-img version test lint up migrate migrate-create migrate-down
