APP_NAME=musorok
APP_PORT?=8080

.PHONY: dev run build test migrate seed docker-up docker-down fmt

dev: ## Run server locally (requires Postgres/Redis running)
	go run ./cmd/server

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP_NAME) ./cmd/server

test:
	go test ./... -v

fmt:
	gofmt -w .

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down -v

migrate:
	# Uses migrate Docker image to run SQL migrations against localhost
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate:4 -path=/migrations -database "postgres://postgres:postgres@localhost:5432/musorok?sslmode=disable" up

seed:
	go run ./seed
