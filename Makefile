# Makefile for openclaw-top
# Variables
BINARY_NAME=openclaw-top
ALIAS_NAME=octop
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/ateverychance/openclaw-top/internal/version.Version=$(VERSION) -X github.com/ateverychance/openclaw-top/internal/version.Commit=$(COMMIT) -X github.com/ateverychance/openclaw-top/internal/version.Date=$(DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: all build install test clean version help

all: build

# Build with version info
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/openclaw-top

# Build with alias name
alias:
	$(GOBUILD) $(LDFLAGS) -o $(ALIAS_NAME) ./cmd/openclaw-top

# Install to GOPATH/bin with version info
install:
	$(GOCMD) install $(LDFLAGS) ./cmd/openclaw-top

# Install with alias name
install-alias:
	$(GOCMD) install $(LDFLAGS) -o $(ALIAS_NAME) ./cmd/openclaw-top

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(ALIAS_NAME)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Show version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

# Run the binary
run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/openclaw-top && ./$(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  build        Build the binary with version info"
	@echo "  alias        Build binary as 'octop'"
	@echo "  install       Install to GOPATH/bin"
	@echo "  install-alias Install as 'octop' to GOPATH/bin"
	@echo "  test          Run tests"
	@echo "  clean         Clean build artifacts"
	@echo "  deps          Download dependencies"
	@echo "  version       Show version info"
	@echo "  run           Build and run"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION      (default: git tag or 'dev')"
	@echo "  COMMIT       (default: git short hash)"
	@echo "  DATE         (default: current UTC time)"