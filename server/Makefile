# Makefile for Prime Trade Microservices

.PHONY: all build run stop clean test docker-build docker-up docker-down logs

# Default target
all: docker-up

# Build all services locally
build:
	@echo "Building API Gateway..."
	cd api-gateway && go build -o bin/api-gateway ./cmd/server
	@echo "Building Auth Service..."
	cd auth-service && go build -o bin/auth-service ./cmd/server
	@echo "Building Task Service..."
	cd task-service && go build -o bin/task-service ./cmd/server

# Run all services locally (requires databases to be running)
run:
	@echo "Starting services..."
	cd api-gateway && go run ./cmd/server &
	cd auth-service && go run ./cmd/server &
	cd task-service && go run ./cmd/server &

# Stop all running services
stop:
	@pkill -f "api-gateway" || true
	@pkill -f "auth-service" || true
	@pkill -f "task-service" || true

# Clean build artifacts
clean:
	rm -rf api-gateway/bin
	rm -rf auth-service/bin
	rm -rf task-service/bin

# Run tests for all services
test:
	@echo "Testing API Gateway..."
	cd api-gateway && go test ./...
	@echo "Testing Auth Service..."
	cd auth-service && go test ./...
	@echo "Testing Task Service..."
	cd task-service && go test ./...

# Build Docker images
docker-build:
	docker-compose build

# Start all services with Docker Compose
docker-up:
	docker-compose up -d --build

# Stop all services
docker-down:
	docker-compose down

# Stop all services and remove volumes
docker-clean:
	docker-compose down -v --remove-orphans

# View logs
logs:
	docker-compose logs -f

# View logs for specific service
logs-gateway:
	docker-compose logs -f api-gateway

logs-auth:
	docker-compose logs -f auth-service

logs-task:
	docker-compose logs -f task-service

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger for Auth Service..."
	cd auth-service && swag init -g cmd/server/main.go -o docs
	@echo "Generating Swagger for Task Service..."
	cd task-service && swag init -g cmd/server/main.go -o docs

# Download dependencies for all services
deps:
	cd api-gateway && go mod download
	cd auth-service && go mod download
	cd task-service && go mod download

# Format code
fmt:
	cd api-gateway && go fmt ./...
	cd auth-service && go fmt ./...
	cd task-service && go fmt ./...

# Lint code (requires golangci-lint)
lint:
	cd api-gateway && golangci-lint run
	cd auth-service && golangci-lint run
	cd task-service && golangci-lint run

# Health check all services
health:
	@echo "Checking API Gateway..."
	@curl -s http://localhost:8080/health | jq . || echo "API Gateway not responding"
	@echo "\nChecking Auth Service..."
	@curl -s http://localhost:8081/health | jq . || echo "Auth Service not responding"
	@echo "\nChecking Task Service..."
	@curl -s http://localhost:8082/health | jq . || echo "Task Service not responding"
