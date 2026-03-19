# SPDX-License-Identifier: Apache-2.0
# Copyright 2026 Milos Vasic

.PHONY: all build test test-race test-cover vet fmt tidy clean help

# Default target
all: tidy vet test build

# Build the binary
build:
	go build ./...

# Run tests
test:
	go test ./... -count=1

# Run tests with race detection
test-race:
	go test ./... -race -count=1

# Run tests with coverage
test-cover:
	go test ./... -race -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run go vet
vet:
	go vet ./...

# Format code
fmt:
	gofmt -s -w .

# Tidy modules
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	go clean ./...

# Show help
help:
	@echo "DocProcessor - Documentation processing and feature map extraction"
	@echo ""
	@echo "Targets:"
	@echo "  all         - tidy, vet, test, build (default)"
	@echo "  build       - Build all packages"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race detection"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  vet         - Run go vet"
	@echo "  fmt         - Format code"
	@echo "  tidy        - Run go mod tidy"
	@echo "  clean       - Remove build artifacts"
	@echo "  help        - Show this help"
