# Makefile for Incident Management System

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=incident-management
BINARY_UNIX=$(BINARY_NAME)_unix

# Build targets
.PHONY: all build clean test coverage deps help

all: test build

# Build the application
build:
	@echo "🔨 Building application..."
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Build for Linux
build-linux:
	@echo "🔨 Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f /tmp/integration_test_*

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run unit tests
test:
	@echo "🧪 Running unit tests..."
	$(GOTEST) -v ./internal/...

# Run tests with coverage
coverage:
	@echo "📊 Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./internal/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📈 Coverage report generated: coverage.html"

# Run integration tests
integration-test:
	@echo "🧪 Running integration tests..."
	./scripts/integration-test.sh

# Run all tests (unit + integration)
test-all: test integration-test

# Run PostgreSQL tests
test-postgres:
	@echo "🐘 Running PostgreSQL tests..."
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		echo "⚠️  Setting up test database..."; \
		docker compose up -d postgres; \
		sleep 10; \
		export TEST_DATABASE_URL="postgres://user:password@localhost:5432/incidentdb?sslmode=disable"; \
	fi
	TEST_DATABASE_URL="postgres://user:password@localhost:5432/incidentdb?sslmode=disable" $(GOTEST) -v ./internal/storage
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		docker compose down postgres; \
	fi

# Start development environment
dev:
	@echo "🚀 Starting development environment..."
	docker compose up --build

# Start only the application
run:
	@echo "🚀 Starting application..."
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

# Run with PostgreSQL
run-postgres:
	@echo "🚀 Starting with PostgreSQL..."
	docker compose up -d postgres
	@echo "⏳ Waiting for PostgreSQL..."
	@sleep 10
	DATABASE_URL="postgres://user:password@localhost:5432/incidentdb?sslmode=disable" ./$(BINARY_NAME)

# Database management
db-up:
	@echo "🐘 Starting PostgreSQL..."
	docker compose up -d postgres

db-down:
	@echo "🛑 Stopping PostgreSQL..."
	docker compose down postgres

db-clean: db-down
	@echo "🧹 Cleaning PostgreSQL data..."
	docker compose down -v postgres

# Linting and formatting
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...

vet:
	@echo "🔍 Vetting code..."
	$(GOCMD) vet ./...

# Security scanning
sec:
	@echo "🔒 Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Docker targets
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t incident-management .

docker-run: docker-build
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 incident-management

# Monitoring and metrics
metrics:
	@echo "📊 Fetching metrics..."
	@curl -s http://localhost:8080/metrics 2>/dev/null || echo "❌ Application not running on localhost:8080"

health:
	@echo "❤️  Checking health..."
	@curl -s http://localhost:8080/health 2>/dev/null | jq . || echo "❌ Application not running on localhost:8080"

ready:
	@echo "✅ Checking readiness..."
	@curl -s http://localhost:8080/ready 2>/dev/null | jq . || echo "❌ Application not running on localhost:8080"

# Documentation
docs:
	@echo "📚 Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "🌐 Starting documentation server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "⚠️  godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Benchmarking
bench:
	@echo "⚡ Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Performance profiling
profile-cpu:
	@echo "🔍 Starting CPU profiling (30s)..."
	@curl -s "http://localhost:8080/debug/pprof/profile?seconds=30" > cpu.prof || echo "❌ Application not running with pprof enabled"

profile-mem:
	@echo "🔍 Getting memory profile..."
	@curl -s "http://localhost:8080/debug/pprof/heap" > mem.prof || echo "❌ Application not running with pprof enabled"

# Development helpers
watch:
	@echo "👀 Watching for changes..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "ℹ️  Falling back to basic watch..."; \
		while true; do \
			$(GOTEST) -v ./...; \
			sleep 5; \
		done; \
	fi

# Release targets
release: clean deps test-all build-linux
	@echo "🎉 Release build complete!"

# Help target
help:
	@echo "🚀 Incident Management System - Make Commands"
	@echo ""
	@echo "📦 Building:"
	@echo "  build              Build the application"
	@echo "  build-linux        Build for Linux (production)"
	@echo "  clean              Clean build artifacts"
	@echo "  deps               Install/update dependencies"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test               Run unit tests"
	@echo "  test-all           Run unit + integration tests"
	@echo "  integration-test   Run comprehensive integration tests"
	@echo "  test-postgres      Run PostgreSQL integration tests"
	@echo "  coverage           Generate test coverage report"
	@echo ""
	@echo "🚀 Running:"
	@echo "  run                Run the application (memory store)"
	@echo "  run-postgres       Run with PostgreSQL"
	@echo "  dev                Start development environment (Docker)"
	@echo "  watch              Watch for changes and auto-rebuild"
	@echo ""
	@echo "🐘 Database:"
	@echo "  db-up              Start PostgreSQL"
	@echo "  db-down            Stop PostgreSQL"
	@echo "  db-clean           Clean PostgreSQL data"
	@echo ""
	@echo "📊 Monitoring:"
	@echo "  health             Check application health"
	@echo "  ready              Check application readiness"
	@echo "  metrics            Fetch Prometheus metrics"
	@echo ""
	@echo "🔧 Quality:"
	@echo "  fmt                Format code"
	@echo "  vet                Vet code"
	@echo "  sec                Security scan (requires gosec)"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-run         Run Docker container"
	@echo ""
	@echo "📚 Other:"
	@echo "  docs               Start documentation server"
	@echo "  bench              Run benchmarks"
	@echo "  profile-cpu        CPU profiling (requires running app)"
	@echo "  profile-mem        Memory profiling (requires running app)"
	@echo "  release            Build release version"
	@echo "  help               Show this help message"

# Default target
.DEFAULT_GOAL := help