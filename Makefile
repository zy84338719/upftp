VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BINARY_NAME = upftp

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod

LDFLAGS = -s -w \
	-X "main.Version=$(VERSION)" \
	-X "main.LastCommit=$(COMMIT_SHA)" \
	-X "main.BuildDate=$(BUILD_DATE)"

BUILD_FLAGS = -a -ldflags "$(LDFLAGS)"

.PHONY: all build build-all clean test deps fmt lint help
.PHONY: run run-ftp install uninstall
.PHONY: release release-snapshot release-local validate

all: clean deps build

deps:
	$(GOMOD) download
	$(GOMOD) tidy

MAIN_PATH = ./cmd/server

build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_linux_amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_linux_arm64 $(MAIN_PATH)

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_darwin_amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_darwin_arm64 $(MAIN_PATH)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_windows_amd64.exe $(MAIN_PATH)

test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -rf releases/
	rm -f coverage.out coverage.html

fmt:
	go fmt ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

run: build
	./$(BINARY_NAME)

run-ftp: build
	./$(BINARY_NAME) -enable-ftp -auto

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete!"

uninstall:
	@echo "Removing $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstall complete!"

validate:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser check; \
	else \
		echo "goreleaser not installed. Install with: go install github.com/goreleaser/goreleaser/v2@latest"; \
	fi

release-snapshot:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser build --snapshot --clean; \
	else \
		echo "goreleaser not installed"; exit 1; \
	fi

release-local:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed"; exit 1; \
	fi

release:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
	else \
		echo "goreleaser not installed"; exit 1; \
	fi

dev: deps run

debug-info:
	@echo "Version:     $(VERSION)"
	@echo "Commit:      $(COMMIT_SHA)"
	@echo "Build Date:  $(BUILD_DATE)"
	@echo "Binary:      $(BINARY_NAME)"

help:
	@echo "UPFTP Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Building:"
	@echo "  all           Clean, deps, and build (default)"
	@echo "  build         Build for current platform"
	@echo "  build-all     Build for all platforms"
	@echo "  build-linux   Build for Linux (amd64, arm64)"
	@echo "  build-darwin  Build for macOS (amd64, arm64)"
	@echo "  build-windows Build for Windows (amd64)"
	@echo ""
	@echo "Development:"
	@echo "  run           Build and run"
	@echo "  run-ftp       Build and run with FTP enabled"
	@echo "  dev           Download deps and run"
	@echo "  fmt           Format code"
	@echo "  lint          Lint code (requires golangci-lint)"
	@echo "  test          Run tests"
	@echo "  coverage      Run tests with coverage report"
	@echo ""
	@echo "Release:"
	@echo "  validate            Validate goreleaser config"
	@echo "  release-snapshot    Build snapshot (no tag needed)"
	@echo "  release-local       Local release test (no publish)"
	@echo "  release             Create and publish release"
	@echo ""
	@echo "Install/Uninstall:"
	@echo "  install        Install to /usr/local/bin"
	@echo "  uninstall      Remove from /usr/local/bin"
	@echo ""
	@echo "Other:"
	@echo "  clean          Clean build artifacts"
	@echo "  deps           Download dependencies"
	@echo "  debug-info     Show build information"
	@echo "  help           Show this help"
