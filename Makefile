# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)

.DEFAULT_GOAL := help

# Build the application binary
build:
	@echo "Building the application..."
	@mkdir -p bin
	@$(GOBUILD) -o $(BINARY_PATH) ./cmd/server/main.go

# Run tests
test: test-unit test-integration test-e2e
	@echo "Running all tests..."

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	@./scripts/test-unit.sh

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@./scripts/test-integration.sh

# Run end-to-end tests
test-e2e:
	@echo "Running e2e tests..."
	@./scripts/test-e2e.sh

# Lint the code
lint:
	@echo ">> running golangci-lint"
	@golangci-lint run ./...

# Run the application
run:
	@echo "Running the application..."
	@$(GORUN) ./cmd/server/main.go

# Clean the binary
clean:
	@echo "Cleaning up..."
	@rm -rf bin

# Show help
help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build              Build the application binary to bin/"
	@echo "  run                Build and run the application"
	@echo "  test               Run all tests (unit, integration, e2e)"
	@echo "  test-unit          Run unit tests"
	@echo "  test-integration   Run integration tests"
	@echo "  test-e2e           Run end-to-end tests"
	@echo "  lint               Lint the code (to be implemented)"
	@echo "  clean              Clean the generated binary"
	@echo "  help               Show this help message"
	@echo ""

.PHONY: build test test-unit test-integration test-e2e lint run clean help