COMPOSE := docker-compose
ENV ?= DATABASE_URL=postgres://quantatomai:quantatomai@localhost:5432/quantatomai?sslmode=disable REDIS_URL=localhost:6379

.PHONY: up down logs migrate test run lint fmt

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

migrate:
	$(ENV) go run ./src -migrate-only

test:
	$(ENV) go test -v ./...

run:
	$(ENV) go run ./src

seed:
	PGPASSWORD=quantatomai psql -h localhost -U quantatomai -d quantatomai -f sql/seed/seed_compat.sql

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...
