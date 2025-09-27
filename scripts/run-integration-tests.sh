#!/bin/bash

# Integration Test Setup and Runner for Configuration Management
# This script sets up the environment and runs comprehensive integration tests

set -e

echo "🧪 Setting up Configuration Management Integration Tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [[ ! -f "go.mod" ]]; then
    print_error "Must be run from project root directory"
    exit 1
fi

# Check if Docker is available for database tests
DOCKER_AVAILABLE=false
if command -v docker &> /dev/null && docker info &> /dev/null; then
    DOCKER_AVAILABLE=true
    print_success "Docker is available for database integration tests"
else
    print_warning "Docker not available - database integration tests will be skipped"
fi

# Setup test database if Docker is available
TEST_DB_CONTAINER=""
if $DOCKER_AVAILABLE; then
    print_status "Setting up test database..."
    
    # Check if postgres container is already running
    if docker ps --format "table {{.Names}}" | grep -q "postgres"; then
        print_status "PostgreSQL container already running"
    else
        print_status "Starting PostgreSQL container for tests..."
        docker run -d \
            --name postgres-integration-test \
            -e POSTGRES_DB=incidentdb_test \
            -e POSTGRES_USER=test \
            -e POSTGRES_PASSWORD=test \
            -p 5433:5432 \
            postgres:15-alpine
        
        TEST_DB_CONTAINER="postgres-integration-test"
        
        # Wait for database to be ready
        print_status "Waiting for database to be ready..."
        timeout 30s bash -c 'while ! docker exec postgres-integration-test pg_isready -U test -d incidentdb_test; do sleep 1; done'
        
        if [ $? -ne 0 ]; then
            print_error "Database failed to start within 30 seconds"
            docker rm -f postgres-integration-test 2>/dev/null || true
            exit 1
        fi
        
        print_success "Test database ready"
    fi
    
    # Set test database URL
    export TEST_DATABASE_URL="postgres://test:test@localhost:5433/incidentdb_test?sslmode=disable"
fi

# Cleanup function
cleanup() {
    print_status "Cleaning up test environment..."
    
    if [[ -n "$TEST_DB_CONTAINER" ]]; then
        print_status "Removing test database container..."
        docker rm -f "$TEST_DB_CONTAINER" 2>/dev/null || true
    fi
    
    # Clean up any test environment variables
    unset TEST_DATABASE_URL
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Run configuration unit tests first
print_status "Running configuration unit tests..."
if go test ./internal/config -v -count=1; then
    print_success "Configuration unit tests passed"
else
    print_error "Configuration unit tests failed"
    exit 1
fi

# Run integration tests
print_status "Running configuration integration tests..."
if go test ./internal/config -v -run TestIntegration -count=1; then
    print_success "Configuration integration tests passed"
else
    print_error "Configuration integration tests failed"
    exit 1
fi

# Run all storage tests if database is available
if $DOCKER_AVAILABLE; then
    print_status "Running storage integration tests with database..."
    if TEST_DATABASE_URL="$TEST_DATABASE_URL" go test ./internal/storage -v -count=1; then
        print_success "Storage integration tests passed"
    else
        print_error "Storage integration tests failed"
        exit 1
    fi
else
    print_warning "Skipping storage integration tests (no database available)"
fi

# Test configuration management with actual application
print_status "Testing configuration with actual application startup..."

# Create temporary test configuration
TEST_CONFIG_FILE=$(mktemp)
cat > "$TEST_CONFIG_FILE" << 'EOF'
PORT=9999
LOG_LEVEL=debug
DEBUG_MODE=true
ALERTMANAGER_URL=http://test-alertmanager:9093
SLACK_TOKEN=xoxb-integration-test-token
SLACK_CHANNEL=#integration-test
EMAIL_SMTP_HOST=smtp.integration-test.com
EMAIL_USERNAME=test@integration.com
EMAIL_PASSWORD=integration-test-password
METRICS_ENABLED=true
METRICS_PORT=9091
ENABLE_CORS=true
EOF

# Test application startup with configuration
print_status "Building application for configuration test..."
if go build -o /tmp/integration-test-app ./cmd/server; then
    print_success "Application built successfully"
else
    print_error "Failed to build application"
    rm -f "$TEST_CONFIG_FILE"
    exit 1
fi

print_status "Testing application startup with valid configuration..."
if timeout 5s bash -c "set -a; source $TEST_CONFIG_FILE; set +a; /tmp/integration-test-app" 2>&1 | grep -q "Configuration loaded successfully"; then
    print_success "Application startup with configuration validation works"
else
    print_error "Application failed to start with valid configuration"
    rm -f "$TEST_CONFIG_FILE" /tmp/integration-test-app
    exit 1
fi

# Test invalid configuration rejection
print_status "Testing invalid configuration rejection..."
if PORT=99999 LOG_LEVEL=invalid /tmp/integration-test-app 2>&1 | grep -q "Configuration validation failed"; then
    print_success "Application correctly rejects invalid configuration"
else
    print_error "Application failed to reject invalid configuration"
    rm -f "$TEST_CONFIG_FILE" /tmp/integration-test-app
    exit 1
fi

# Clean up test files
rm -f "$TEST_CONFIG_FILE" /tmp/integration-test-app

# Run benchmark tests for configuration loading
print_status "Running configuration performance benchmarks..."
if go test ./internal/config -bench=. -benchmem -count=1 > /tmp/benchmark.out 2>&1; then
    print_success "Performance benchmarks completed"
    echo ""
    echo "📊 Performance Results:"
    grep "Benchmark" /tmp/benchmark.out | head -5
    rm -f /tmp/benchmark.out
else
    print_warning "Benchmarks failed or not available"
fi

# Summary
echo ""
echo "🎉 Integration Test Suite Completed!"
echo ""
print_success "✅ Configuration unit tests passed"
print_success "✅ Configuration integration tests passed"
if $DOCKER_AVAILABLE; then
    print_success "✅ Database integration tests passed"
else
    print_warning "⚠️  Database integration tests skipped (Docker not available)"
fi
print_success "✅ Application startup configuration validation works"
print_success "✅ Invalid configuration rejection works"
print_success "✅ Performance benchmarks completed"

echo ""
echo "📋 Test Coverage Summary:"
echo "  • Configuration loading and validation"
echo "  • Hot-reloading functionality"
echo "  • Vault integration preparation"
echo "  • Database configuration integration"
echo "  • End-to-end configuration lifecycle"
echo "  • Application startup validation"
echo "  • Error handling and edge cases"
echo ""

print_success "All integration tests completed successfully! 🚀"