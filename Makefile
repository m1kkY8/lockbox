# Lockbox Makefile

# Variables
BINARY_NAME=lockbox
GO_MAIN=./main.go
BUILD_DIR=./bin
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-w -s -X main.version=$(VERSION)"

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
.PHONY: help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(GO_MAIN)

## test: Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

## test-coverage: Run tests and show coverage
.PHONY: test-coverage
test-coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

## run: Build and run the application with test server
.PHONY: run
run: build
	@$(BUILD_DIR)/$(BINARY_NAME) -h sjdoo.zapto.org

## dev: Build and run with local development server
.PHONY: dev
dev: build
	@$(BUILD_DIR)/$(BINARY_NAME) -h "localhost:1337" -u "test"

## install: Build and install the application
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	@go build $(LDFLAGS)
	@go install

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html

## deps: Download and verify dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## security: Run security analysis
.PHONY: security
security:
	@echo "Running security analysis..."
	@gosec ./...

## build-all: Build for all platforms
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(GO_MAIN)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(GO_MAIN)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(GO_MAIN)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(GO_MAIN)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(GO_MAIN)

## ci: Run CI pipeline locally (test, lint, build)
.PHONY: ci
ci: deps fmt vet lint test build
	@echo "CI pipeline completed successfully!"
