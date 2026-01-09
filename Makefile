.PHONY: build run test clean docker-build docker-build-api docker-up docker-down docker-db docker-redis docker-api docker-pgadmin docker-redis-ui docker-ui docker-logs migrate wire wire-gen docker-network docker-clean docker-clean-all

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application (builds and starts both Docker services)
run: docker-build docker-up
	@echo "Services are starting..."
	@echo "API: http://localhost:8080"
	@echo "pgAdmin: http://localhost:5050 (admin@admin.com / admin)"
	@echo "RedisInsight: http://localhost:8081"
	@echo "View logs with: make docker-logs"

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-network:
	@echo "Creating Docker network 'studoto'..."
	@./setup-network.sh

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

docker-down-volumes:
	docker-compose down -v
	@echo "⚠️  All volumes have been removed. Database will be re-initialized on next run."

docker-clean:
	@echo "Cleaning up containers and volumes..."
	docker-compose down -v --remove-orphans
	@echo "✅ Containers and volumes removed"

docker-clean-all:
	@echo "⚠️  Removing containers, volumes, and images..."
	docker-compose down -v --remove-orphans --rmi local
	@echo "✅ All containers, volumes, and images removed"

docker-db:
	docker-compose up -d postgres

docker-redis:
	docker-compose up -d redis

docker-pgadmin:
	docker-compose up -d pgadmin
	@echo "pgAdmin is available at http://localhost:5050"
	@echo "Email: admin@admin.com"
	@echo "Password: admin"
	@echo ""
	@echo "Waiting for pgAdmin to initialize (this may take 30-60 seconds)..."
	@sleep 30
	@echo "PostgreSQL server 'Studoto PostgreSQL' should appear automatically"

docker-redis-ui:
	docker-compose up -d redisinsight
	@echo "RedisInsight is available at http://localhost:8081"
	@echo "After opening, click 'Add Redis Database' and use:"
	@echo "  Host: 172.28.5.4"
	@echo "  Port: 6379"
	@echo "  Password: Da3ZqphucUo4zw9b"
	@echo ""
	@echo "Note: RedisInsight may take 30-60 seconds to fully start"

docker-ui:
	docker-compose up -d pgadmin redisinsight
	@echo "pgAdmin: http://localhost:5050 (admin@admin.com / admin)"
	@echo "  PostgreSQL server 'Studoto PostgreSQL' should appear automatically"
	@echo "RedisInsight: http://localhost:8081"
	@echo "  Click 'Add Redis Database' → Host: 172.28.5.4, Port: 6379, Password: Da3ZqphucUo4zw9b"
	@echo ""
	@echo "Note: Both services may take 30-60 seconds to fully initialize"

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

# Generate Wire code
wire:
	cd cmd/server && wire

# Generate Wire code and tidy dependencies
wire-gen: wire
	go mod tidy
