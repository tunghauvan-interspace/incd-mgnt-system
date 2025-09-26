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
   docker-compose up -d
   ```

4. **Access the dashboard:**
   - Incident Management System: http://localhost:8080
   - Prometheus: http://localhost:9090
   - Alertmanager: http://localhost:9093

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

The system is configured via environment variables:

### Core Settings
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Log level (default: info)
- `ALERTMANAGER_URL`: Alertmanager URL (default: http://localhost:9093)

### Notification Settings
- `SLACK_TOKEN`: Slack bot token
- `SLACK_CHANNEL`: Slack channel for notifications
- `EMAIL_SMTP_HOST`: SMTP server host
- `EMAIL_SMTP_PORT`: SMTP server port
- `EMAIL_USERNAME`: SMTP username
- `EMAIL_PASSWORD`: SMTP password
- `TELEGRAM_BOT_TOKEN`: Telegram bot token
- `TELEGRAM_CHAT_ID`: Telegram chat ID

### Example .env file:
```env
PORT=8080
LOG_LEVEL=info
ALERTMANAGER_URL=http://localhost:9093

# Slack Configuration
SLACK_TOKEN=xoxb-your-slack-token
SLACK_CHANNEL=#alerts

# Email Configuration
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-password

# Telegram Configuration
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
TELEGRAM_CHAT_ID=your-chat-id
```

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
