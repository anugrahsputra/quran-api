.PHONY: build test run lint clean help

# Variables
BINARY_NAME=main
BINARY_PATH=tmp/$(BINARY_NAME)

# Default target
all: build

# Build the application
build: 
	@echo "Building the application..."
	@go build -o $(BINARY_PATH) main.go

# Run the unit tests
test:
	@echo "Running unit tests..."
	@go test ./...

# Run the application
run:
	@echo "Starting the application..."
	@go run main.go

# Run the linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_PATH)

# Display help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run the unit tests"
	@echo "  make run      - Run the application"
	@echo "  make lint     - Run the linter"
	@echo "  make clean    - Clean up build artifacts"
	@echo "  make help     - Display this help message"
