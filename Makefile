.PHONY: build run docker-up migrate test

build:
	go build -o bin/auth cmd/server

run:
	go run ./cmd/server

docker-up:
	docker-compose up --build

migrate:
	@echo "Run migrations manually: psql -U postgres -d authdb < internal/db/migrations.sql"

test:
	go test ./... -v