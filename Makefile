.PHONY: help build test test-coverage lint clean docker-build run

# Default target
help:
	@echo "Drone OpenAI Plugin - Available targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  run            - Run the plugin locally (requires env vars)"

# Build the binary
build:
	@echo "Building drone-openai-plugin..."
	@go build -o drone-openai-plugin ./cmd/plugin
	@echo "Build complete: ./drone-openai-plugin"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@echo ""
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter (requires golangci-lint to be installed)
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f drone-openai-plugin
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t drone-openai-plugin:latest .
	@echo "Docker image built: drone-openai-plugin:latest"

# Run the plugin locally (for testing)
# Example: make run PROMPT="Say hello" API_KEY=sk-xxx
run: build
	@echo "Running plugin..."
	@PLUGIN_API_KEY=$(API_KEY) \
	 PLUGIN_PROMPT="$(PROMPT)" \
	 PLUGIN_MODEL=$(or $(MODEL),gpt-4o-mini) \
	 PLUGIN_FILE=$(FILE) \
	 PLUGIN_OUTPUT_FILE=$(OUTPUT_FILE) \
	 ./drone-openai-plugin

