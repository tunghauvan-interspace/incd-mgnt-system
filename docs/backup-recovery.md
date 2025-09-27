# Database Backup and Recovery Procedures

This document describes the backup and recovery procedures for the Incident Management System's PostgreSQL database.

## Overview

The system provides automated scripts for database backup and restoration, supporting both compressed and uncompressed SQL dumps.

## Backup Procedures

### Manual Backup

To create a manual backup of the database:

```bash
# Create a backup with automatic timestamp
./scripts/backup-database.sh

# Create a backup with custom name
./scripts/backup-database.sh my_backup_name
```

### Backup Configuration

The backup script supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKUP_DIR` | `./backups` | Directory to store backups |
| `POSTGRES_HOST` | `localhost` | PostgreSQL server host |
| `POSTGRES_PORT` | `5432` | PostgreSQL server port |
| `POSTGRES_USER` | `user` | Database username |
| `POSTGRES_DB` | `incidentdb` | Database name |
| `POSTGRES_PASSWORD` | `password` | Database password |
| `BACKUP_RETENTION_DAYS` | `7` | Number of days to keep backups |

### Automated Backup

For automated backups, add a cron job:

```bash
# Example: Daily backup at 2:00 AM
0 2 * * * /path/to/incident-management/scripts/backup-database.sh daily_backup
```

## Recovery Procedures

### Full Database Restore

⚠️ **WARNING**: Database restore will completely replace the existing database and all current data will be lost.

```bash
# Restore from compressed backup
./scripts/restore-database.sh ./backups/backup_20231127_143022.sql.gz

# Restore with confirmation bypass (use with caution)
./scripts/restore-database.sh ./backups/backup_20231127_143022.sql.gz --force
```

### Pre-Restore Safety

The restore script automatically creates a pre-restore backup before proceeding with the restoration, providing an additional safety net.

## Docker Environment

When using Docker Compose, you can run backup/restore operations in several ways:

### Option 1: Host-based (Recommended)

Install PostgreSQL client tools on the host and use the scripts directly.

### Option 2: Container-based

```bash
# Create backup from inside container
docker exec incd-mgnt-system-postgres-1 pg_dump -U user -d incidentdb > backup.sql

# Restore backup to container
docker exec -i incd-mgnt-system-postgres-1 psql -U user -d postgres < backup.sql
```

## Database Health Monitoring

The system provides several endpoints to monitor database health:

- **Health Check**: `GET /health` - Basic application health
- **Readiness Check**: `GET /ready` - Detailed health including database connectivity
- **Database Stats**: `GET /db/stats` - Connection pool statistics and metrics

### Database Statistics

The `/db/stats` endpoint provides detailed information about the database connection pool:

```json
{
  "timestamp": "2025-09-27T09:51:58Z",
  "stats": {
    "database": {
      "health": "healthy",
      "type": "postgresql",
      "max_open_connections": 25,
      "open_connections": 1,
      "in_use": 1,
      "idle": 0,
      "wait_count": 0,
      "wait_duration": "0s",
      "max_idle_closed": 0,
      "max_idle_time_closed": 0,
      "max_lifetime_closed": 0
    }
  }
}
```

## Connection Pool Configuration

The database connection pool is configured through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_MAX_OPEN_CONNS` | `25` | Maximum number of open connections |
| `DB_MAX_IDLE_CONNS` | `5` | Maximum number of idle connections |
| `DB_CONN_MAX_LIFETIME` | `5m` | Maximum connection lifetime |

## Migration Management

Database migrations are handled automatically when the application starts:

1. **Up Migrations**: Applied automatically on startup
2. **Down Migrations**: Must be applied manually if needed
3. **Migration Files**: Located in `migrations/` directory

### Migration Files

- `001_initial_schema.up.sql` - Creates initial database schema
- `001_initial_schema.down.sql` - Drops database schema

## Troubleshooting

### Common Issues

#### Backup Fails with "Connection Refused"
- Verify PostgreSQL is running
- Check connection parameters (host, port, credentials)
- Ensure `pg_dump` is installed and accessible

#### Restore Fails with "Database Does Not Exist"
- The restore script handles database creation automatically
- Ensure the PostgreSQL user has sufficient privileges

#### Application Won't Start After Restore
- Check database connectivity with `/ready` endpoint
- Verify migration state in `schema_migrations` table
- Check application logs for detailed error messages

### Support Commands

```bash
# Check database connection
docker exec incd-mgnt-system-postgres-1 pg_isready -U user -d incidentdb

# View database tables
docker exec incd-mgnt-system-postgres-1 psql -U user -d incidentdb -c "\dt"

# Check migration status
docker exec incd-mgnt-system-postgres-1 psql -U user -d incidentdb -c "SELECT * FROM schema_migrations;"
```

## Recovery Testing

Regularly test backup and recovery procedures to ensure they work correctly:

```bash
# 1. Create test backup
./scripts/backup-database.sh test_recovery

# 2. Perform restore
./scripts/restore-database.sh ./backups/test_recovery.sql.gz --force

# 3. Verify application functionality
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

## Security Considerations

- Store backups in secure locations
- Encrypt sensitive backup files
- Limit access to backup scripts and files
- Use strong database passwords
- Regular security audits of backup procedures

## Contact Information

For questions or issues with backup and recovery procedures, please refer to the system documentation or contact the development team.