# Complete Setup Guide - Incident Management System

This guide provides comprehensive instructions for setting up the Incident Management System with its modern Vue.js frontend and Go backend.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Development Setup](#development-setup) 
3. [Production Deployment](#production-deployment)
4. [Configuration](#configuration)
5. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- **Docker & Docker Compose**: Latest version
- **Node.js**: Version 20+ (for frontend development)
- **Go**: Version 1.21+ (for backend development)
- **Git**: Latest version

### System Requirements

- **Memory**: 4GB RAM minimum (8GB recommended)
- **Storage**: 10GB free space minimum
- **CPU**: 2 cores minimum (4 cores recommended)
- **Network**: Internet access for downloading dependencies

## Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/tunghauvan-interspace/incd-mgnt-system.git
cd incd-mgnt-system
```

### 2. Environment Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit configuration (see Configuration section below)
nano .env
```

### 3. Frontend Development Setup

```bash
# Navigate to frontend directory
cd web/frontend

# Install Node.js dependencies
npm install

# Start development server
npm run dev
```

The Vue.js development server will run on: http://localhost:5173

### 4. Backend Development Setup

#### Option A: Using Docker Compose (Recommended)

```bash
# Start all services including database and monitoring
docker-compose --profile development up -d

# Or start only specific services
docker-compose up -d postgres prometheus alertmanager
```

#### Option B: Manual Go Development

```bash
# Install Go dependencies
go mod download

# Run database migrations (if using local PostgreSQL)
go run cmd/migrate/main.go up

# Start the Go backend
go run cmd/server/main.go
```

### 5. Accessing Services

- **Frontend (Vue.js)**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093
- **PostgreSQL**: localhost:5432

## Production Deployment

### 1. Prepare Environment

```bash
# Clone repository on production server
git clone https://github.com/tunghauvan-interspace/incd-mgnt-system.git
cd incd-mgnt-system

# Setup production environment
cp .env.example .env
# Configure production values in .env
```

### 2. Build Frontend Assets

```bash
# Build Vue.js application for production
cd web/frontend
npm install
npm run build

# Assets will be created in web/frontend/dist/
```

### 3. Deploy with Docker

#### Option A: Full Stack Deployment

```bash
# Deploy all services
docker-compose --profile production up -d
```

#### Option B: Custom Deployment

```bash
# Build and deploy backend only
docker-compose up -d postgres prometheus alertmanager backend
```

### 4. Serve Frontend Assets

**Option A: Use a reverse proxy (Nginx/Apache)**

```nginx
# Example Nginx configuration
server {
    listen 80;
    server_name your-domain.com;
    
    # Serve Vue.js static files
    location / {
        root /path/to/web/frontend/dist;
        try_files $uri $uri/ /index.html;
    }
    
    # Proxy API requests to Go backend
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Option B: Copy to Go backend static directory**

```bash
# Copy built assets to Go backend (if serving from Go)
cp -r web/frontend/dist/* web/static/
```

## Configuration

### Environment Variables

#### Core Settings

```bash
# Server Configuration
PORT=8080
LOG_LEVEL=info
DEBUG_MODE=false

# Database Configuration
DATABASE_URL=postgres://user:password@localhost:5432/incidentdb?sslmode=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

#### Monitoring Integration

```bash
# Prometheus/Alertmanager
ALERTMANAGER_URL=http://alertmanager:9093
METRICS_ENABLED=true
METRICS_PORT=9090
```

#### Notification Channels

```bash
# Slack Integration
SLACK_TOKEN=xoxb-your-slack-token
SLACK_CHANNEL=#incidents

# Email/SMTP
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your-email@domain.com
EMAIL_PASSWORD=your-app-password

# Telegram
TELEGRAM_BOT_TOKEN=your-bot-token
TELEGRAM_CHAT_ID=your-chat-id
```

### Frontend Configuration

The Vue.js frontend automatically proxies API requests to the backend during development. For production, ensure your web server properly routes API requests.

## Project Architecture

### Frontend (Vue.js + TypeScript)

```
web/frontend/
├── src/
│   ├── components/     # Reusable UI components
│   ├── views/          # Page components (Dashboard, Incidents, Alerts)
│   ├── services/       # API integration layer
│   ├── types/          # TypeScript type definitions
│   ├── utils/          # Helper functions
│   └── assets/         # Static assets and global styles
├── dist/               # Built assets (production)
├── package.json        # Dependencies and scripts
└── vite.config.ts      # Build configuration
```

### Backend (Go)

```
internal/
├── handlers/           # HTTP request handlers
├── services/           # Business logic layer
├── storage/            # Database layer
├── models/             # Data structures
└── middleware/         # HTTP middleware

cmd/
├── server/             # Main application entry point
└── migrate/            # Database migration tool
```

## Available Commands

### Frontend Commands

```bash
cd web/frontend

# Development
npm run dev          # Start development server
npm run build        # Build for production
npm run preview      # Preview production build

# Code Quality
npm run lint         # Lint code
npm run format       # Format code
npm run type-check   # TypeScript type checking
```

### Backend Commands

```bash
# Development
go run cmd/server/main.go              # Start server
go run cmd/migrate/main.go up          # Run database migrations
go run cmd/migrate/main.go down        # Rollback migrations

# Production
go build -o bin/server cmd/server/main.go  # Build binary
./bin/server                               # Run binary

# Testing
go test ./...                          # Run all tests
go test -cover ./...                   # Run tests with coverage
```

### Docker Commands

```bash
# Development
docker-compose --profile development up -d     # Start dev environment
docker-compose --profile development down      # Stop dev environment

# Production
docker-compose --profile production up -d      # Start production
docker-compose --profile production down       # Stop production

# All services (default)
docker-compose up -d                           # Start all services
docker-compose down                            # Stop all services

# Individual services
docker-compose up -d postgres                  # Start only database
docker-compose up -d frontend backend          # Start app services
```

## Troubleshooting

### Common Issues

#### Frontend Issues

**1. Node modules not found**
```bash
cd web/frontend
rm -rf node_modules package-lock.json
npm install
```

**2. Build fails due to TypeScript errors**
```bash
npm run type-check  # Check for type errors
npm run lint        # Check for linting issues
```

**3. API requests failing (CORS/Proxy)**
- Ensure backend is running on port 8080
- Check Vite proxy configuration in `vite.config.ts`
- Verify API endpoints return proper CORS headers

#### Backend Issues

**1. Database connection failed**
```bash
# Check if PostgreSQL is running
docker-compose up -d postgres

# Verify database credentials in .env
# Test connection manually
psql -h localhost -U user -d incidentdb
```

**2. Port already in use**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

**3. Missing environment variables**
```bash
# Verify .env file exists and has required variables
cat .env | grep -E "(DATABASE_URL|PORT)"
```

#### Docker Issues

**1. Container fails to start**
```bash
# Check container logs
docker-compose logs <service-name>

# Common solutions
docker-compose down
docker-compose pull
docker-compose up -d
```

**2. Build context issues**
```bash
# Rebuild images
docker-compose build --no-cache

# Clean up Docker resources
docker system prune -a
```

### Health Checks

#### Service Status

```bash
# Check if all services are running
docker-compose ps

# Check service logs
docker-compose logs -f backend
docker-compose logs -f frontend
```

#### API Health

```bash
# Test backend API
curl http://localhost:8080/health

# Test metrics endpoint
curl http://localhost:8080/metrics
```

#### Frontend Health

```bash
# Check if frontend is accessible
curl http://localhost:5173
```

## Performance Optimization

### Frontend

- **Code Splitting**: Automatically handled by Vite
- **Asset Optimization**: Images and bundles are optimized during build
- **Caching**: Set appropriate cache headers for static assets

### Backend  

- **Database**: Configure connection pooling appropriately
- **Monitoring**: Use Prometheus metrics to identify bottlenecks
- **Caching**: Implement Redis caching for frequently accessed data

### Docker

- **Multi-stage builds**: Dockerfile uses multi-stage builds for smaller images
- **Volume optimization**: Use named volumes for persistent data
- **Resource limits**: Set memory and CPU limits for production

## Security Considerations

### Production Deployment

1. **Use HTTPS**: Configure SSL/TLS certificates
2. **Environment Variables**: Never commit secrets to git
3. **Database Security**: Use strong passwords and network isolation
4. **CORS Configuration**: Restrict origins in production
5. **Rate Limiting**: Implement rate limiting for API endpoints

### Development

1. **Local Environment**: Keep development isolated from production
2. **Dependency Updates**: Regularly update dependencies for security patches
3. **Code Review**: Review all code changes before deployment

## Support

For additional support:

1. **Documentation**: Check `web/FRONTEND_INTEGRATION.md` for detailed frontend integration
2. **Issues**: Report bugs and issues on the GitHub repository
3. **Logs**: Always check service logs when troubleshooting issues

## Next Steps

After successful setup:

1. **Configure Monitoring**: Set up Prometheus alerts and Grafana dashboards
2. **Customize Notifications**: Configure Slack, Email, or Telegram notifications
3. **Integration**: Connect with your existing Alertmanager setup
4. **Scaling**: Consider horizontal scaling for high-availability deployments