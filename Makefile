VERSION = $(shell git describe --always --tags)
BUILD = $(shell date +%F)
COMMIT_SHA=$(shell git rev-parse --short HEAD)
BINARY_NAME = upftp

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Build flags
LDFLAGS = -X "main.Version=$(VERSION)" -X "main.LastCommit=$(COMMIT_SHA)"
BUILD_FLAGS = -a -ldflags "$(LDFLAGS)"

# Default target
.PHONY: all
all: clean deps build

# Show debug information
.PHONY: debugInfo
debugInfo:
	@echo "VERSION:"    $(VERSION)
	@echo "COMMIT_SHA:" $(COMMIT_SHA)
	@echo "BUILD:"      $(BUILD)

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for current platform
.PHONY: build
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) ./

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

# Build for Linux (amd64 and arm64)
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_linux_amd64 ./
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_linux_arm64 ./

# Build for Windows (amd64 and 386)
.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_windows_amd64.exe ./
	GOOS=windows GOARCH=386 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_windows_386.exe ./

# Build for macOS (amd64 and arm64)
.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_darwin_amd64 ./
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_darwin_arm64 ./

# Create release packages
.PHONY: package
package: build-all
	@mkdir -p releases
	# Package Linux builds
	tar -czf releases/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz -C dist $(BINARY_NAME)_linux_amd64
	tar -czf releases/$(BINARY_NAME)_$(VERSION)_linux_arm64.tar.gz -C dist $(BINARY_NAME)_linux_arm64
	# Package Windows builds
	cd dist && zip ../releases/$(BINARY_NAME)_$(VERSION)_windows_amd64.zip $(BINARY_NAME)_windows_amd64.exe
	cd dist && zip ../releases/$(BINARY_NAME)_$(VERSION)_windows_386.zip $(BINARY_NAME)_windows_386.exe
	# Package macOS builds
	tar -czf releases/$(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz -C dist $(BINARY_NAME)_darwin_amd64
	tar -czf releases/$(BINARY_NAME)_$(VERSION)_darwin_arm64.tar.gz -C dist $(BINARY_NAME)_darwin_arm64

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -rf releases/

# Development run
.PHONY: run
run: build
	./$(BINARY_NAME)

# Development run with FTP enabled
.PHONY: run-ftp
run-ftp: build
	./$(BINARY_NAME) -enable-ftp

# Install dependencies and run
.PHONY: dev
dev: deps run

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, download deps, and build"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms (Linux, Windows, macOS)"
	@echo "  build-linux  - Build for Linux (amd64, arm64)"
	@echo "  build-windows- Build for Windows (amd64, 386)"
	@echo "  build-darwin - Build for macOS (amd64, arm64)"
	@echo "  package      - Create release packages"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run"
	@echo "  run-ftp      - Build and run with FTP enabled"
	@echo "  dev          - Download deps and run"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  deps         - Download dependencies"
	@echo "  debugInfo    - Show build information"
	@echo "  help         - Show this help"
