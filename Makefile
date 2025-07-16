VERSION = $(shell git describe --always --tags)
BUILD = $(shell date +%F)
COMMIT_SHA=$(shell git rev-parse --shor# Release packages using GoReleaser
.PHONY: release
release:
	@echo "Creating release with GoReleaser..."
	goreleaser release --clean

# Test release locally
.PHONY: release-test
release-test:
	@echo "Testing release with GoReleaser..."
	goreleaser release --snapshot --clean

# Build packages manually (for testing)
.PHONY: build-packages
build-packages: build-all
	@echo "Building packages manually..."
	mkdir -p dist/packages
	
	# Create DEB package structure
	for arch in amd64 arm64; do \
		mkdir -p dist/packages/upftp_$(VERSION)_$$arch/usr/bin; \
		mkdir -p dist/packages/upftp_$(VERSION)_$$arch/DEBIAN; \
		cp dist/$(BINARY_NAME)_linux_$$arch dist/packages/upftp_$(VERSION)_$$arch/usr/bin/upftp; \
		sed "s/{{.Version}}/$(VERSION)/g; s/{{.Arch}}/$$arch/g" packaging/debian/control > dist/packages/upftp_$(VERSION)_$$arch/DEBIAN/control; \
		cp packaging/debian/postinst dist/packages/upftp_$(VERSION)_$$arch/DEBIAN/; \
		cp packaging/debian/prerm dist/packages/upftp_$(VERSION)_$$arch/DEBIAN/; \
		cp packaging/debian/postrm dist/packages/upftp_$(VERSION)_$$arch/DEBIAN/; \
		chmod 755 dist/packages/upftp_$(VERSION)_$$arch/DEBIAN/*; \
		dpkg-deb --build dist/packages/upftp_$(VERSION)_$$arch dist/packages/; \
	done

# Install from source
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete!"

# Uninstall
.PHONY: uninstall
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstall complete!"

# Setup APT repository (requires proper setup)
.PHONY: setup-apt-repo
setup-apt-repo:
	@echo "Setting up APT repository..."
	@if [ ! -d "/var/www/apt" ]; then \
		echo "Creating APT repository structure..."; \
		sudo mkdir -p /var/www/apt/{pool/main,dists/stable/main/binary-{amd64,arm64}}; \
		sudo chown -R $(USER):$(USER) /var/www/apt; \
	fi
	@echo "APT repository structure ready at /var/www/apt"

# Publish to APT repository
.PHONY: publish-apt
publish-apt: build-packages setup-apt-repo
	@echo "Publishing packages to APT repository..."
	./scripts/publish-apt.sh

# Generate Homebrew formula
.PHONY: generate-brew-formula
generate-brew-formula:
	@echo "Generating Homebrew formula..."
	@VERSION_NUM=$(VERSION); \
	SHA256=$$(curl -sL "https://github.com/zy84338719/upftp/archive/$$VERSION_NUM.tar.gz" | sha256sum | cut -d' ' -f1); \
	sed "s/{{.Version}}/$$VERSION_NUM/g; s/{{.SHA256}}/$$SHA256/g" packaging/brew/upftp.rb.template > upftp.rb
	@echo "Formula generated: upftp.rb"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, download deps, and build"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-linux  - Build for Linux (amd64, arm64)"
	@echo "  build-windows- Build for Windows (amd64, 386)"
	@echo "  build-darwin - Build for macOS (amd64, arm64)"
	@echo "  package      - Create release packages"
	@echo "  release      - Create release with GoReleaser"
	@echo "  release-test - Test release locally with GoReleaser"
	@echo "  build-packages - Build DEB packages manually"
	@echo "  install      - Install from source to /usr/local/bin"
	@echo "  uninstall    - Remove from /usr/local/bin"
	@echo "  publish-apt  - Publish packages to APT repository"
	@echo "  generate-brew-formula - Generate Homebrew formula"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run"
	@echo "  run-ftp      - Build and run with FTP enabled"
	@echo "  dev          - Download deps and run"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  deps         - Download dependencies"
	@echo "  debugInfo    - Show build information"
	@echo "  help         - Show this help"= upftp

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
