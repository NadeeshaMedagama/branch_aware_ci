.PHONY: build test clean install run help init docker-build docker-run

# Variables
BINARY_NAME=branch-aware-ci
VERSION=1.0.0
BUILD_DIR=bin
DOCKER_IMAGE=branch-aware-ci
GO=go

# Default target
all: build

## help: Show this help message
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install the binary"
	@echo "  run           - Build and run the binary"
	@echo "  init          - Initialize configuration file"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run in Docker container"
	@echo "  lint          - Run linters"
	@echo "  coverage      - Generate test coverage report"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-X 'main.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

## test-coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## install: Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install .
	@echo "Installed"

## run: Build and run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

## init: Initialize configuration file
init: build
	./$(BUILD_DIR)/$(BINARY_NAME) -init

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	@echo "Docker image built: $(DOCKER_IMAGE):$(VERSION)"

## docker-run: Run in Docker container
docker-run:
	@echo "Running in Docker..."
	docker run -v $(PWD):/repo $(DOCKER_IMAGE):latest -repo /repo

## lint: Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Formatted"

## tidy: Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy
	@echo "Dependencies tidied"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	@echo "Dependencies downloaded"

## version: Show version
version:
	@echo "$(VERSION)"

## build-services: Build all microservices
build-services:
	@echo "Building microservices..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/gateway ./services/gateway
	$(GO) build -o $(BUILD_DIR)/branch-detector ./services/branch-detector
	$(GO) build -o $(BUILD_DIR)/policy-engine ./services/policy-engine
	@echo "All services built"

## run-gateway: Run API Gateway
run-gateway:
	@echo "Starting API Gateway..."
	PORT=8080 ./$(BUILD_DIR)/gateway

## run-branch-detector: Run Branch Detector service
run-branch-detector:
	@echo "Starting Branch Detector..."
	HTTP_PORT=8081 GRPC_PORT=50051 ./$(BUILD_DIR)/branch-detector

## run-policy-engine: Run Policy Engine service
run-policy-engine:
	@echo "Starting Policy Engine..."
	HTTP_PORT=8082 GRPC_PORT=50052 ./$(BUILD_DIR)/policy-engine

## docker-compose-up: Start all services with Docker Compose
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

## docker-compose-down: Stop all services
docker-compose-down:
	@echo "Stopping services..."
	docker-compose down

## docker-compose-logs: Show service logs
docker-compose-logs:
	docker-compose logs -f

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	docker-compose up -d
	@sleep 10
	@echo "Testing Gateway health..."
	@curl -f http://localhost:8080/health || exit 1
	@echo "Testing Branch Detector health..."
	@curl -f http://localhost:8081/health || exit 1
	@echo "Testing Policy Engine health..."
	@curl -f http://localhost:8082/health || exit 1
	@echo "All integration tests passed!"
	docker-compose down

## docker-build-all: Build all Docker images
docker-build-all:
	@echo "Building all Docker images..."
	docker build -f services/gateway/Dockerfile -t branch-aware-gateway:$(VERSION) .
	docker build -f services/branch-detector/Dockerfile -t branch-aware-branch-detector:$(VERSION) .
	docker build -f services/policy-engine/Dockerfile -t branch-aware-policy-engine:$(VERSION) .
	@echo "All images built"

## proto-gen: Generate protobuf code
proto-gen:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/branchaware/v1/*.proto
	@echo "Protobuf code generated"


