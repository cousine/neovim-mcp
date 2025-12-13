.PHONY: help build clean test test-unit test-integration fmt lint vet run install install-deps dev

# Variables
BINARY_NAME=neovim-mcp
BUILD_DIR=dist
MAIN_PATH=./cmd/neovim-mcp
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags="-s -w"

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-dev: ## Build without optimizations for debugging
	@echo "Building $(BINARY_NAME) (dev mode)..."
	$(GO) build $(GOFLAGS) -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Dev build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS) $(LDFLAGS) $(MAIN_PATH)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

clean: ## Remove built binaries and test cache
	@echo "Cleaning..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@$(GO) clean -testcache
	@echo "Clean complete"

# Test targets
test: ## Run all tests (unit + integration)
	@echo "Running all tests..."
	@gotestsum --format testdox --format-icons hivis -- -tags=integration -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "Coverage Summary (lowest coverage first):"
	@go tool cover -func=coverage.out | grep -v "^mode:" | sort -t: -k2 -n | head -20
	@echo ""
	@go tool cover -func=coverage.out | tail -1

test-unit: ## Run unit tests only (fast)
	@echo "Running unit tests..."
	@gotestsum --format testdox --format-icons hivis -- -race -coverprofile=coverage.out ./internal/... ./cmd/...
	@echo ""
	@echo "Coverage Summary (lowest coverage first):"
	@go tool cover -func=coverage.out | grep -v "^mode:" | sort -t: -k2 -n | head -20
	@echo ""
	@go tool cover -func=coverage.out | tail -1

test-integration: ## Run integration tests with Docker containers (default)
	@echo "Running integration tests with containers..."
	@echo "Note: Requires Docker to be running"
	@echo "Tip: Set NEOVIM_TEST_VERBOSE=1 to show container logs"
	@if [ "$$NEOVIM_TEST_VERBOSE" = "1" ]; then \
		gotestsum --format standard-verbose --format-icons hivis -- -tags=integration -race -coverprofile=coverage.out ./test/integration/...; \
	else \
		gotestsum --format testdox --format-icons hivis -- -tags=integration -race -coverprofile=coverage.out ./test/integration/...; \
	fi
	@echo ""
	@echo "Coverage Summary (lowest coverage first):"
	@go tool cover -func=coverage.out | grep -v "^mode:" | sort -t: -k2 -n | head -20
	@echo ""
	@go tool cover -func=coverage.out | tail -1

test-integration-local: ## Run integration tests with local Neovim
	@echo "Running integration tests with local Neovim..."
	@echo "Note: Requires Neovim to be installed"
	@echo "Tip: Set NEOVIM_TEST_VERBOSE=1 to show detailed test output"
	@if [ "$$NEOVIM_TEST_VERBOSE" = "1" ]; then \
		NEOVIM_TEST_LOCAL=1 gotestsum --format standard-verbose --format-icons hivis -- -tags=integration -race -coverprofile=coverage.out ./test/integration/...; \
	else \
		NEOVIM_TEST_LOCAL=1 gotestsum --format testdox --format-icons hivis -- -tags=integration -race -coverprofile=coverage.out ./test/integration/...; \
	fi
	@echo ""
	@echo "Coverage Summary (lowest coverage first):"
	@go tool cover -func=coverage.out | grep -v "^mode:" | sort -t: -k2 -n | head -20
	@echo ""
	@go tool cover -func=coverage.out | tail -1

test-coverage: ## Run tests with coverage report (HTML)
	@echo "Running tests with coverage..."
	@gotestsum --format testdox --format-icons hivis -- -tags=integration -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "Detailed Coverage Summary (lowest coverage first):"
	@go tool cover -func=coverage.out | grep -v "^mode:" | sort -t: -k2 -n | head -20
	@echo ""
	@echo "Total Coverage:"
	@go tool cover -func=coverage.out | tail -1
	@echo ""
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Code quality targets
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "Vet complete"

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run ./...
	@echo "Lint complete"

check: fmt vet ## Run fmt and vet

# Dependency management
mod-download: ## Download Go modules
	@echo "Downloading dependencies..."
	$(GO) mod download
	@echo "Download complete"

mod-tidy: ## Tidy Go modules
	@echo "Tidying modules..."
	$(GO) mod tidy
	@echo "Tidy complete"

mod-verify: ## Verify Go modules
	@echo "Verifying modules..."
	$(GO) mod verify
	@echo "Verify complete"

install-deps: ## Install required development dependencies
	@echo "Installing development dependencies..."
	@echo ""
	@echo "Installing gotestsum..."
	@go install gotest.tools/gotestsum@latest
	@echo "  ✓ gotestsum installed"
	@echo ""
	@echo "Installing golangci-lint (if not present)..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "  ✓ golangci-lint installed"
	@echo ""
	@echo "All dependencies installed successfully!"
	@echo ""
	@echo "Installed versions:"
	@echo "  gotestsum:      $(shell gotestsum --version 2>/dev/null || echo 'installed')"
	@echo "  golangci-lint:  $(shell golangci-lint --version 2>/dev/null | head -1 || echo 'installed')"

# Development targets
dev: ## Run in development mode (with auto-reload would require additional tools)
	@echo "Starting development server..."
	@echo "Note: Make sure Neovim is running with: nvim --listen /tmp/nvim.sock"
	NVIM_LISTEN_ADDRESS=/tmp/nvim.sock $(GO) run $(MAIN_PATH)

run: ## Run the server (requires Neovim to be running)
	@echo "Starting $(BINARY_NAME)..."
	@echo "Note: Make sure Neovim is running with: nvim --listen /tmp/nvim.sock"
	@test -S /tmp/nvim.sock || (echo "Error: Neovim socket not found at /tmp/nvim.sock"; exit 1)
	./$(BINARY_NAME)

# Neovim helpers
start-nvim: ## Start Neovim with RPC socket (for testing)
	@echo "Starting Neovim with RPC socket..."
	nvim --listen /tmp/nvim.sock

start-nvim-headless: ## Start headless Neovim for integration tests
	@echo "Starting headless Neovim..."
	nvim --headless --listen /tmp/nvim-test.sock &
	@echo "Neovim started on /tmp/nvim-test.sock"

kill-nvim: ## Kill any running Neovim instances
	@echo "Killing Neovim instances..."
	@pkill -f "nvim.*listen" || echo "No Neovim instances found"
	@rm -f /tmp/nvim.sock /tmp/nvim-test.sock

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	$(GO) doc -all ./... > docs.txt
	@echo "Documentation generated: docs.txt"

# CI/CD targets
ci: mod-download check test-unit ## Run CI checks (no integration tests)
	@echo "CI checks complete"

ci-full: mod-download check test ## Run full CI checks (including integration)
	@echo "Full CI checks complete"

# Release targets
release: clean test build ## Build a release version
	@echo "Release build complete"

release-all: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Multi-platform builds complete"

# Info targets
version: ## Show Go version
	@$(GO) version

deps: ## List dependencies
	@echo "Direct dependencies:"
	@$(GO) list -m all | grep -v "neovim-mcp$$"

info: ## Show project information
	@echo "Project: $(BINARY_NAME)"
	@echo "Main: $(MAIN_PATH)"
	@echo "Go version: $(shell go version)"
	@echo "Build directory: $(BUILD_DIR)"
	@echo ""
	@echo "File statistics:"
	@echo "  Go files: $(shell find . -name '*.go' -not -path './vendor/*' | wc -l | tr -d ' ')"
	@echo "  Total lines: $(shell find . -name '*.go' -not -path './vendor/*' -exec wc -l {} + | tail -1 | awk '{print $$1}')"
	@echo ""
	@echo "Dependencies:"
	@$(GO) list -m all | wc -l | awk '{print "  Modules: " $$1}'

.DEFAULT_GOAL := help
