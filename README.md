# Incident Management System

An open-source, self-hosted incident management platform built for Prometheus and Alertmanager. It manages the full incident lifecycle (create, acknowledge, resolve), groups alerts, and provides real-time notifications via Slack, Email, and Telegram.

## Features

- **Full Incident Lifecycle Management**: Create, acknowledge, and resolve incidents
- **Alert Grouping**: Automatically groups related alerts into incidents
- **Multi-channel Notifications**: Slack, Email, and Telegram integration
- **Prometheus/Alertmanager Integration**: Seamless webhook integration
- **Modern Dashboard**: Real-time incident tracking with MTTA/MTTR metrics
- **Escalation Policies**: (Framework in place for future enhancement)
- **On-call Scheduling**: (Framework in place for future enhancement)
- **REST API**: Complete API for incident and alert management

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/tunghauvan-interspace/incd-mgnt-system.git
   cd incd-mgnt-system
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the services:**
   ```bash
   # For development (includes Vue.js frontend dev server)
   docker-compose --profile development up -d
   
   # For production
   docker-compose --profile production up -d
   ```

4. **Access the services:**
   - **Vue.js Frontend** (Development): http://localhost:5173  
   - **Go Backend API**: http://localhost:8080
   - **Prometheus**: http://localhost:9090
   - **Alertmanager**: http://localhost:9093

### Complete Setup Guide

For detailed setup instructions, troubleshooting, and production deployment guidance, see:

📖 **[Complete Setup Guide](SETUP_GUIDE.md)**

This guide covers:
- Development environment setup
- Production deployment strategies  
- Configuration options
- Frontend build process
- Docker service profiles
- Troubleshooting common issues

### Manual Installation

1. **Install Go 1.21+**

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the application:**
   ```bash
   go build -o incident-management ./cmd/server
   ```

4. **Run the application:**
   ```bash
   ./incident-management
   ```

## Configuration

The Incident Management System is highly configurable through environment variables. All configuration options are documented in the `.env.example` file with detailed explanations.

### Configuration Management Features

- **Comprehensive Validation**: All configuration is validated on startup with clear error messages
- **Hot Reloading**: Non-sensitive configuration can be reloaded without restart
- **Secure Credential Management**: Prepared for HashiCorp Vault integration
- **Environment Variable Support**: Full fallback to environment variables
- **Type Safety**: Automatic parsing and validation of different data types

### Quick Setup

1. **Copy the example configuration:**
   ```bash
   cp .env.example .env
   ```

2. **Configure required settings:**
   ```bash
   # Edit .env with your specific settings
   vim .env
   ```

3. **Validate configuration:**
   The application will validate your configuration on startup and provide clear error messages for any issues.

### Core Settings

#### Server Configuration
- `PORT` - HTTP server port (default: 8080)
- `LOG_LEVEL` - Logging level: debug, info, warn, error (default: info)  
- `ALERTMANAGER_URL` - Alertmanager webhook source (default: http://localhost:9093)
- `ALERTMANAGER_TIMEOUT` - Connection timeout in seconds (default: 30)

#### Database Configuration
- `DATABASE_URL` - PostgreSQL connection string (optional - uses in-memory if not set)
- `DB_MAX_OPEN_CONNS` - Maximum open connections (default: 25)
- `DB_MAX_IDLE_CONNS` - Maximum idle connections (default: 5)
- `DB_CONN_MAX_LIFETIME` - Connection lifetime (default: 5m)

### Notification Settings

All notification channels are optional. Configure any combination based on your needs.

#### Slack Integration
- `SLACK_TOKEN` - Bot OAuth token (xoxb-*)
- `SLACK_CHANNEL` - Channel for notifications (#channel-name)

Required Slack bot permissions:
- `chat:write` - Send messages
- `channels:read` - Access public channels
- `groups:read` - Access private channels

#### Email/SMTP Configuration  
- `EMAIL_SMTP_HOST` - SMTP server hostname
- `EMAIL_SMTP_PORT` - SMTP port (default: 587)
- `EMAIL_USERNAME` - SMTP authentication username
- `EMAIL_PASSWORD` - SMTP authentication password
- `EMAIL_FROM` - From address (optional, defaults to EMAIL_USERNAME)
- `EMAIL_TO` - Default recipient (optional)

Common SMTP configurations:
- **Gmail**: smtp.gmail.com:587 (use App Passwords)
- **Office365**: smtp.office365.com:587
- **SendGrid**: smtp.sendgrid.net:587

#### Telegram Integration
- `TELEGRAM_BOT_TOKEN` - Bot API token from @BotFather
- `TELEGRAM_CHAT_ID` - Chat ID for notifications

### Security Settings

#### TLS/HTTPS Configuration
- `TLS_CERT_FILE` - Path to TLS certificate file
- `TLS_KEY_FILE` - Path to TLS private key file

Both must be provided to enable HTTPS.

#### Server Timeouts
- `SERVER_READ_TIMEOUT` - Request read timeout (default: 30s)
- `SERVER_WRITE_TIMEOUT` - Response write timeout (default: 30s)
- `SERVER_IDLE_TIMEOUT` - Keep-alive timeout (default: 120s)

### Advanced Configuration

#### Metrics and Monitoring
- `METRICS_ENABLED` - Enable Prometheus metrics (default: true)
- `METRICS_PORT` - Metrics endpoint port (default: 9090)

#### Operational Settings
- `WEBHOOK_TIMEOUT` - Webhook processing timeout (default: 30s)
- `NOTIFICATION_TIMEOUT` - Notification delivery timeout (default: 15s)
- `MAX_INCIDENT_AGE` - Auto-resolve incidents after duration (default: 24h)

#### CORS Configuration
- `ENABLE_CORS` - Enable CORS headers (default: true)
- `CORS_ORIGIN` - Allowed origins (default: *)

#### Development Settings
- `DEBUG_MODE` - Enable debug features (default: false)

**⚠️ Never enable DEBUG_MODE in production!**

### HashiCorp Vault Integration (Future)

The system is prepared for HashiCorp Vault integration for secure credential management:

```bash
# Vault Configuration (Future Feature)
VAULT_ENABLED=true
VAULT_ADDR=https://vault.company.com:8200
VAULT_TOKEN=your-vault-token
VAULT_SECRET_PATH=secret/incident-management
```

### Configuration Validation

The system performs comprehensive validation on startup:

- **Port numbers**: Must be valid (1-65535)
- **Log levels**: Must be debug, info, warn, or error  
- **Database connections**: Idle connections cannot exceed open connections
- **Notification settings**: If one part is configured, required parts must also be set
- **TLS settings**: Both cert and key files required for HTTPS
- **Timeouts**: Must be positive values
- **Type validation**: Automatic parsing with error messages

Example validation error:
```
Configuration validation failed:
  - config validation failed for PORT: must be between 1 and 65535
  - config validation failed for LOG_LEVEL: must be one of: debug, info, warn, error
  - config validation failed for SLACK_CHANNEL: required when SLACK_TOKEN is provided
```

### Hot Configuration Reloading

Non-sensitive configuration can be reloaded without restarting:

**Reloadable Settings:**
- Log level
- Timeouts  
- CORS settings
- Debug mode
- Metrics settings

**Non-Reloadable Settings (require restart):**
- Server port
- Database connection settings
- TLS certificates
- Notification credentials

### Configuration Examples

#### Minimal Configuration (In-Memory)
```env
PORT=8080
LOG_LEVEL=info
```

#### Production with PostgreSQL and Slack
```env
# Server
PORT=8080
LOG_LEVEL=warn

# Database  
DATABASE_URL=postgres://incident_user:secure_pass@db:5432/incidentdb?sslmode=require
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10

# Notifications
SLACK_TOKEN=xoxb-your-production-token
SLACK_CHANNEL=#alerts-production

# Security
TLS_CERT_FILE=/etc/ssl/certs/server.crt
TLS_KEY_FILE=/etc/ssl/private/server.key
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9090
```

#### Development Configuration
```env
PORT=8080
LOG_LEVEL=debug
DEBUG_MODE=true
DATABASE_URL=postgres://dev:dev@localhost:5432/incident_dev
```

For complete configuration options and detailed explanations, see the [`.env.example`](.env.example) file.

## API Endpoints

### Incidents
- `GET /api/incidents` - List all incidents
- `GET /api/incidents/{id}` - Get incident details
- `POST /api/incidents/{id}/acknowledge` - Acknowledge an incident
- `POST /api/incidents/{id}/resolve` - Resolve an incident

### Alerts
- `GET /api/alerts` - List all alerts
- `POST /api/webhooks/alertmanager` - Alertmanager webhook endpoint

### Metrics
- `GET /api/metrics` - Get incident metrics (MTTA, MTTR, etc.)

### Health
- `GET /health` - Health check endpoint

## Dashboard

The web dashboard provides:

- **Real-time Metrics**: MTTA, MTTR, incident counts by status and severity
- **Incident Management**: View, acknowledge, and resolve incidents
- **Alert Monitoring**: View all alerts and their status
- **Charts**: Visual representation of incident trends

## Alertmanager Integration

Configure Alertmanager to send webhooks to the incident management system:

```yaml
# alertmanager.yml
route:
  receiver: 'incident-management'

receivers:
  - name: 'incident-management'
    webhook_configs:
      - url: 'http://incident-management:8080/api/webhooks/alertmanager'
        send_resolved: true
```

## Alert Grouping

Alerts are automatically grouped into incidents based on:
- Service label matching
- Instance label matching
- Alert name matching

The system creates new incidents for ungrouped alerts and adds related alerts to existing open incidents.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Prometheus    │───▶│   Alertmanager   │───▶│  Incident Mgmt  │
│                 │    │                  │    │     System      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │  Notifications  │
                                               │ Slack|Email|TG  │
                                               └─────────────────┘
```

## Development

1. **Run tests:**
   ```bash
   go test ./...
   ```

2. **Run with live reload:**
   ```bash
   # Install air for live reloading
   go install github.com/cosmtrek/air@latest
   air
   ```

3. **Build for production:**
   ```bash
   CGO_ENABLED=0 go build -ldflags="-w -s" -o incident-management ./cmd/server
   ```

## Docker Deployment

### Build and run with Docker:
```bash
docker build -t incident-management .
docker run -p 8080:8080 --env-file .env incident-management
```

### Using docker-compose with custom configuration:
```bash
# Edit docker-compose.yml as needed
docker-compose up -d
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Support

For issues, questions, or contributions, please visit our [GitHub repository](https://github.com/tunghauvan-interspace/incd-mgnt-system).

## Roadmap

- [ ] Database persistence (PostgreSQL, MySQL)
- [ ] User authentication and authorization
- [ ] Advanced escalation policies
- [ ] On-call scheduling implementation
- [ ] Incident templates
- [ ] SLA tracking
- [ ] Incident post-mortems
- [ ] Mobile app
- [ ] More notification channels (PagerDuty, Discord)
- [ ] Custom dashboards and reports
