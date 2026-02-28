.PHONY: help build run test clean mocks infra-start infra-stop run-local smoke-test

GCP_PROJECT_ID ?= parish-local

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application (requires env vars)"
	@echo "  make test           - Run unit tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make mocks          - Generate/update mocks"
	@echo "  make infra-start    - Start local infrastructure (Datastore emulator + Redis)"
	@echo "  make infra-stop     - Stop local infrastructure"
	@echo "  make run-local      - Start infrastructure + run the API locally"
	@echo "  make smoke-test     - Run happy-path smoke tests against local API"

build:
	go build -o bin/parish-api ./cmd

run:
	go run ./cmd

test:
	go test -v ./...

clean:
	rm -rf bin/

mocks:
	@echo "Generating mocks..."
	go generate ./...

infra-start:
	@echo "Starting local infrastructure (Datastore emulator + Redis)..."
	docker compose up -d
	@echo "Waiting for Datastore emulator to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		curl -s http://localhost:8081 >/dev/null 2>&1 && break; \
		sleep 1; \
	done
	@echo "Infrastructure ready (Datastore :8081, Redis :6379)"

infra-stop:
	@echo "Stopping local infrastructure..."
	docker compose down
	@echo "Infrastructure stopped."

run-local: infra-start
	DATASTORE_EMULATOR_HOST=localhost:8081 \
	GCP_PROJECT_ID=$(GCP_PROJECT_ID) \
	COOKIE_SECURE=false \
	CORS_ORIGIN=http://localhost:3000 \
	PORT=8080 \
	REDIS_URL=localhost:6379 \
	ADMIN_EMAIL=admin@parish.local \
	ADMIN_PASSWORD=Admin@Str0ng!Pass \
	go run ./cmd

smoke-test:
	@./scripts/smoke_test.sh
