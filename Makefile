VERSION = $(shell git describe --always --tags)
BUILD = $(shell date +%F)
COMMIT_SHA=$(shell git rev-parse --short HEAD)

debugInfo:
	@echo "VERSION:"    $(VERSION)
	@echo "COMMIT_SHA:" $(COMMIT_SHA)
	@echo "BUILD:"      $(BUILD)


build:
	go build -a -ldflags " -X \"main.Version=$(VERSION)\" -X \"main.LastCommit=$(COMMIT_SHA)\" " -o upftp ./
