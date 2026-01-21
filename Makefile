.PHONY: all build test clean run-users run-products run-orders run-gateway docker-build docker-compose-up docker-compose-down

# Build all services
all: build

# Build all services
build:
	@echo "Building all services..."
	@cd services/users-service && go build -o ../../bin/users-service
	@cd services/products-service && go build -o ../../bin/products-service
	@cd services/orders-service && go build -o ../../bin/orders-service
	@cd services/api-gateway && go build -o ../../bin/api-gateway
	@echo "Build complete!"

# Run tests for all packages and services
test:
	@echo "Running tests..."
	@go test ./pkg/...
	@go test ./services/users-service/...
	@go test ./services/products-service/...
	@go test ./services/orders-service/...
	@go test ./services/api-gateway/...
	@echo "Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Run individual services
run-users:
	@cd services/users-service && go run main.go

run-products:
	@cd services/products-service && go run main.go

run-orders:
	@cd services/orders-service && go run main.go

run-gateway:
	@cd services/api-gateway && go run main.go

# Docker commands
docker-build:
	@echo "Building Docker images..."
	@docker build -t users-service:latest -f services/users-service/Dockerfile .
	@docker build -t products-service:latest -f services/products-service/Dockerfile .
	@docker build -t orders-service:latest -f services/orders-service/Dockerfile .
	@docker build -t api-gateway:latest -f services/api-gateway/Dockerfile .
	@echo "Docker images built!"

# Docker Compose commands
docker-compose-up:
	@docker-compose up -d

docker-compose-down:
	@docker-compose down

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin"

# Create bin directory
bin:
	@mkdir -p bin

# Help
help:
	@echo "Available targets:"
	@echo "  make build           - Build all services"
	@echo "  make test            - Run all tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make run-users       - Run users-service"
	@echo "  make run-products    - Run products-service"
	@echo "  make run-orders      - Run orders-service"
	@echo "  make run-gateway     - Run api-gateway"
	@echo "  make docker-build    - Build Docker images"
	@echo "  make deps            - Install dependencies"
	@echo "  make fmt             - Format code"
	@echo "  make lint            - Lint code"
