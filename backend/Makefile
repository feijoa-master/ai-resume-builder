.PHONY: help install run build test clean docker-up docker-down migrate

# Variables
BINARY_NAME=resume-builder-api
MAIN_PATH=./cmd/api

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	go mod download
	go mod tidy

run: ## Run the application
	go run $(MAIN_PATH)/main.go

build: ## Build the application
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@psql $(DATABASE_URL) -f migrations/001_initial_schema.sql

migrate-down: ## Run database migrations down (rollback)
	@echo "Rolling back migrations..."
	# Add rollback SQL when needed

db-shell: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d resume_builder

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

dev: docker-up ## Start development environment
	@sleep 3
	@make run