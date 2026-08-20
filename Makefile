.PHONY: all build build-all clean test vet fmt run install help
.PHONY: package package-deb package-rpm

BINARY  := upftp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)
GO      := go

## build: 编译当前平台的二进制到 ./bin/$(BINARY)
build:
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

## build-all: 交叉编译 linux/darwin/windows × amd64/arm64 到 dist/
build-all:
	@mkdir -p dist
	@for os in linux darwin windows; do \
	  for arch in amd64 arm64; do \
	    ext=""; [ $$os = windows ] && ext=".exe"; \
	    echo "  → $$os/$$arch"; \
	    GOOS=$$os GOARCH=$$arch $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)_$${os}_$${arch}$$ext .; \
	  done; \
	done

## test: 运行测试
test:
	$(GO) test ./... -v

## vet: 静态检查
vet:
	$(GO) vet ./...

## fmt: 格式化
fmt:
	gofmt -w internal/ *.go

## clean: 清理构建产物
clean:
	rm -rf bin dist

## run: 直接运行(共享当前目录,默认端口)
run:
	$(GO) run . -d .

## install: 安装到 $$GOBIN
install: build
	$(GO) install -ldflags "$(LDFLAGS)" .

## package: Build .deb and .rpm for amd64+arm64
package:
	@./scripts/package.sh --version $(VERSION) --config upftp.example.yaml --service packaging/systemd/upftp.service --format all --arch amd64,arm64

## package-deb: Build .deb only
package-deb:
	@./scripts/package.sh --version $(VERSION) --config upftp.example.yaml --service packaging/systemd/upftp.service --format deb --arch amd64,arm64

## package-rpm: Build .rpm only
package-rpm:
	@./scripts/package.sh --version $(VERSION) --config upftp.example.yaml --service packaging/systemd/upftp.service --format rpm --arch amd64,arm64

## help: Show this help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'
