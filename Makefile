# Makefile for Incident Management System

# Variables
BINARY_NAME=incident-management
GO_VERSION=1.21

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build test test-unit test-integration test-config test-all clean lint fmt vet deps benchmark

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the application
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@go build -o $(BINARY_NAME) ./cmd/server
	@echo "$(GREEN)Build complete: $(BINARY_NAME)$(NC)"

build-docker: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(BINARY_NAME) .
	@echo "$(GREEN)Docker image built: $(BINARY_NAME)$(NC)"

# Test targets
test: test-unit ## Run unit tests (default)

test-unit: ## Run unit tests only
	@echo "$(GREEN)Running unit tests...$(NC)"
	@go test ./... -v -short

test-config: ## Run configuration tests only
	@echo "$(GREEN)Running configuration tests...$(NC)"
	@go test ./internal/config -v

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@./scripts/run-integration-tests.sh

test-all: ## Run all tests (unit + integration)
	@echo "$(GREEN)Running all tests...$(NC)"
	@go test ./... -v -short
	@./scripts/run-integration-tests.sh

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

benchmark: ## Run benchmark tests
	@echo "$(GREEN)Running benchmark tests...$(NC)"
	@go test ./internal/config -bench=. -benchmem

# Database targets
db-start: ## Start PostgreSQL database using Docker
	@echo "$(GREEN)Starting PostgreSQL database...$(NC)"
	@docker compose up -d postgres
	@echo "$(GREEN)Database started$(NC)"

db-stop: ## Stop PostgreSQL database
	@echo "$(GREEN)Stopping PostgreSQL database...$(NC)"
	@docker compose down postgres
	@echo "$(GREEN)Database stopped$(NC)"

db-reset: ## Reset database (stop, remove, start)
	@echo "$(GREEN)Resetting database...$(NC)"
	@docker compose down postgres
	@docker volume rm $$(docker volume ls -q | grep postgres) 2>/dev/null || true
	@docker compose up -d postgres
	@echo "$(GREEN)Database reset complete$(NC)"

db-verify: ## Verify database connectivity and schema
	@echo "$(GREEN)Verifying database...$(NC)"
	@./scripts/verify-database.sh

# Development targets
dev: ## Start development environment
	@echo "$(GREEN)Starting development environment...$(NC)"
	@docker compose up -d postgres prometheus alertmanager
	@echo "$(GREEN)Development environment started$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"
	@echo "$(YELLOW)Prometheus: http://localhost:9090$(NC)"
	@echo "$(YELLOW)Alertmanager: http://localhost:9093$(NC)"

dev-stop: ## Stop development environment
	@echo "$(GREEN)Stopping development environment...$(NC)"
	@docker compose down
	@echo "$(GREEN)Development environment stopped$(NC)"

run: build ## Build and run the application
	@echo "$(GREEN)Starting $(BINARY_NAME)...$(NC)"
	@./$(BINARY_NAME)

run-debug: build ## Build and run with debug mode
	@echo "$(GREEN)Starting $(BINARY_NAME) in debug mode...$(NC)"
	@DEBUG_MODE=true LOG_LEVEL=debug ./$(BINARY_NAME)

# Code quality targets
fmt: ## Format Go code
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "$(GREEN)Running linter...$(NC)"
	@golangci-lint run || echo "$(YELLOW)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"

deps: ## Download and tidy dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

deps-update: ## Update dependencies
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

# Configuration targets
config-validate: ## Validate current configuration
	@echo "$(GREEN)Validating configuration...$(NC)"
	@go run ./cmd/server --validate-config 2>/dev/null || echo "$(YELLOW)Note: --validate-config flag not implemented yet$(NC)"

config-example: ## Show example configuration
	@echo "$(GREEN)Example configuration:$(NC)"
	@cat .env.example

# Clean targets
clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "$(GREEN)Clean complete$(NC)"

clean-all: clean ## Clean everything including Docker containers and volumes
	@echo "$(GREEN)Cleaning all Docker resources...$(NC)"
	@docker compose down --volumes --remove-orphans 2>/dev/null || true
	@docker rmi $(BINARY_NAME) 2>/dev/null || true
	@echo "$(GREEN)Clean all complete$(NC)"

# Production targets
prod-build: ## Build for production
	@echo "$(GREEN)Building for production...$(NC)"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o $(BINARY_NAME) ./cmd/server
	@echo "$(GREEN)Production build complete$(NC)"

prod-test: ## Run production readiness tests
	@echo "$(GREEN)Running production readiness tests...$(NC)"
	@$(MAKE) test-all
	@$(MAKE) benchmark
	@echo "$(GREEN)Production tests complete$(NC)"

# Docker Compose targets
up: ## Start all services using Docker Compose
	@echo "$(GREEN)Starting all services...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)All services started$(NC)"

down: ## Stop all services
	@echo "$(GREEN)Stopping all services...$(NC)"
	@docker compose down
	@echo "$(GREEN)All services stopped$(NC)"

logs: ## Show logs from all services
	@docker compose logs -f

logs-app: ## Show logs from application only
	@docker compose logs -f incident-management

# Health check targets
health: ## Check health of running services
	@echo "$(GREEN)Checking service health...$(NC)"
	@curl -s http://localhost:8080/health || echo "$(YELLOW)Application not responding$(NC)"
	@curl -s http://localhost:9090/-/healthy || echo "$(YELLOW)Prometheus not responding$(NC)"
	@curl -s http://localhost:9093/-/healthy || echo "$(YELLOW)Alertmanager not responding$(NC)"

# Documentation targets
docs: ## Generate documentation
	@echo "$(GREEN)Generating documentation...$(NC)"
	@godoc -http=:6060 &
	@echo "$(GREEN)Documentation server started at http://localhost:6060$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"

# Install targets
install: build ## Install the application binary
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	@cp $(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)$(BINARY_NAME) installed to /usr/local/bin/$(NC)"

uninstall: ## Uninstall the application binary
	@echo "$(GREEN)Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)$(BINARY_NAME) uninstalled$(NC)"

# Additional targets from main branch
test-postgres: ## Run PostgreSQL tests
	@echo "$(GREEN)Running PostgreSQL tests...$(NC)"
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		echo "$(YELLOW)Setting up test database...$(NC)"; \
		docker compose up -d postgres; \
		sleep 10; \
		export TEST_DATABASE_URL="postgres://user:password@localhost:5432/incidentdb?sslmode=disable"; \
	fi
	TEST_DATABASE_URL="postgres://user:password@localhost:5432/incidentdb?sslmode=disable" $(GOTEST) -v ./internal/storage
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		docker compose down postgres; \
	fi

docker-run: build-docker ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -p 8080:8080 $(BINARY_NAME)

profile-cpu: ## CPU profiling (requires running app)
	@echo "$(GREEN)Starting CPU profiling (30s)...$(NC)"
	@curl -s "http://localhost:8080/debug/pprof/profile?seconds=30" > cpu.prof || echo "$(YELLOW)Application not running with pprof enabled$(NC)"

profile-mem: ## Memory profiling (requires running app)
	@echo "$(GREEN)Getting memory profile...$(NC)"
	@curl -s "http://localhost:8080/debug/pprof/heap" > mem.prof || echo "$(YELLOW)Application not running with pprof enabled$(NC)"

watch: ## Watch for changes and auto-rebuild
	@echo "$(GREEN)Watching for changes...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(YELLOW)Falling back to basic watch...$(NC)"; \
		while true; do \
			$(GOTEST) -v ./...; \
			sleep 5; \
		done; \
	fi

release: clean deps test-all build-linux ## Build release version
	@echo "$(GREEN)Release build complete!$(NC)"
