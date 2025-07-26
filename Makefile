.PHONY: build run dev test clean docker-up docker-down migrate

# Build the application
build:
	go build -o bin/main cmd/main.go

# Run the application
run: build
	./bin/main

# Run in development mode with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start Docker containers
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# Start Docker containers and wait for database
docker-setup: docker-up
	@echo "Waiting for database to be ready..."
	@sleep 5

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
dev-deps:
	go install github.com/cosmtrek/air@latest

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run the application with database
start: docker-setup run

# Full setup for development
setup: deps dev-deps docker-setup
	@echo "Setup complete! Run 'make dev' to start development server"

