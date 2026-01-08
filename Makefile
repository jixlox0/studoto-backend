.PHONY: build run test clean docker-build docker-build-api docker-up docker-down docker-db docker-api docker-logs migrate

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application (builds and starts both Docker services)
run: docker-build docker-up
	@echo "Services are starting..."
	@echo "API will be available at http://localhost:8080"
	@echo "View logs with: make docker-logs"

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-build:
	docker-compose build

docker-build-api:
	docker-compose build api

docker-up:
	docker-compose up -d

docker-api:
	docker-compose up -d api

docker-down:
	docker-compose down

docker-db:
	docker-compose up -d postgres

docker-logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Run all checks
check: fmt test

