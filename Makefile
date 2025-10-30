.PHONY: build install test clean lint help

BINARY_NAME=kcavo
BUILD_DIR=bin
INSTALL_PATH=$(HOME)/.local/bin

# Build the plugin
build:
	@echo "🔨 Building kubectl-cost..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install to local bin (adds to PATH if ~/.local/bin is in PATH)
install: build
	@echo "📦 Installing kubectl-cost..."
	@mkdir -p $(INSTALL_PATH)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "   Make sure $(INSTALL_PATH) is in your PATH"
	@echo "   Usage: kubectl cost <command>"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Lint code
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@echo "   This may take a minute on first run..."
	@go mod download
	@go mod tidy
	@go mod verify
	@echo "✅ Dependencies installed and verified"
	@echo ""
	@echo "📦 Installed packages:"
	@go list -m all | head -20

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ Multi-platform build complete"

# Quick test with real cluster
try:
	@echo "🚀 Testing with your cluster..."
	go run . analyze

# Display help
help:
	@echo "Kcavo - Kubernetes Cost Analysis Plugin"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build the plugin binary"
	@echo "  install     - Build and install to ~/.local/bin"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format code"
	@echo "  clean       - Remove build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  try         - Quick test with your cluster"
	@echo "  help        - Show this help message"