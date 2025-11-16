.PHONY: build test run lint clean help reindex deps format vet test-coverage docker-build docker-run docker-down docker-logs install-linter clean-all

# Variables
BINARY_NAME=quran-api
BINARY_PATH=tmp/$(BINARY_NAME)
GO_VERSION := $(shell go version | awk '{print $$3}')

# Default target
.DEFAULT_GOAL := help

# Build the application
build:
	@echo "Building the application..."
	@mkdir -p tmp
	@go build -o $(BINARY_PATH) cmd/main.go
	@echo "Binary built at $(BINARY_PATH)"

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Run the application
run:
	@echo "Starting the application..."
	@go run cmd/main.go

# Run the unit tests
test:
	@echo "Running unit tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Re-index Quran data
reindex:
	@echo "Re-indexing Quran data..."
	@go run cmd/main.go -reindex

# Format code
format:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run the linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/"; \
		exit 1; \
	fi

# Install golangci-lint (if not installed)
install-linter:
	@echo "Installing golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint already installed"; \
	else \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
	fi

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t quran-api:latest .
	@echo "Docker image built: quran-api:latest"

# Build Docker image with pre-built search index
# This requires quran.bleve to exist in the project root
docker-build-with-index:
	@echo "Building Docker image with search index..."
	@if [ ! -d "quran.bleve" ]; then \
		echo "Error: quran.bleve not found. Building index first..."; \
		go run cmd/main.go -reindex || (echo "Failed to build index. Please run 'make reindex' first." && exit 1); \
	fi
	@docker build -f Dockerfile.with-index -t quran-api:with-index .
	@echo "Docker image with index built: quran-api:with-index"

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started. Use 'make docker-logs' to view logs."

# Stop Docker Compose services
docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down
	@echo "Services stopped"

# View Docker Compose logs
docker-logs:
	@docker-compose logs -f

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_PATH)
	@rm -f coverage.out coverage.html
	@rm -rf tmp/
	@echo "Cleanup complete"

# Full clean (including test cache)
clean-all: clean
	@echo "Cleaning test cache..."
	@go clean -testcache
	@echo "Full cleanup complete"

# Display help
help:
	@echo "Quran API - Makefile Commands"
	@echo "=============================="
	@echo ""
	@echo "Build & Run:"
	@echo "  make build          - Build the application binary"
	@echo "  make run            - Run the application"
	@echo "  make reindex        - Re-index Quran data for search"
	@echo ""
	@echo "Development:"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make format         - Format code with gofmt"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint (requires installation)"
	@echo "  make install-linter - Install golangci-lint"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run unit tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build          - Build Docker image (without index)"
	@echo "  make docker-build-with-index - Build Docker image with search index"
	@echo "  make docker-run            - Start services with Docker Compose"
	@echo "  make docker-down           - Stop Docker Compose services"
	@echo "  make docker-logs           - View Docker Compose logs"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make clean-all      - Remove build artifacts and test cache"
	@echo ""
	@echo "Info:"
	@echo "  make help           - Display this help message"
	@echo ""
	@echo "Go Version: $(GO_VERSION)"
