.PHONY: dev build test lint migrate-up migrate-down

dev:
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -o bin/api ./cmd/api

test:
	go test ./... -v

lint:
	go vet ./...

migrate-up:
	migrate -path ./migrations -database "postgres://asana:asana_password@localhost:5432/asana?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://asana:asana_password@localhost:5432/asana?sslmode=disable" down

tidy:
	go mod tidy
