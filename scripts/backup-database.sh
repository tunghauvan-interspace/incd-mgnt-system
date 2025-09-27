#!/bin/bash

# Database Backup Script for Incident Management System
# Usage: ./backup-database.sh [backup_name]

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
POSTGRES_USER="${POSTGRES_USER:-user}"
POSTGRES_DB="${POSTGRES_DB:-incidentdb}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-password}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Generate backup name
if [ -n "$1" ]; then
    BACKUP_NAME="$1"
else
    BACKUP_NAME="backup_$(date +%Y%m%d_%H%M%S)"
fi

BACKUP_FILE="$BACKUP_DIR/${BACKUP_NAME}.sql"
BACKUP_COMPRESSED="$BACKUP_DIR/${BACKUP_NAME}.sql.gz"

echo "🔄 Starting database backup..."
echo "📁 Backup directory: $BACKUP_DIR"
echo "🏷️  Backup name: $BACKUP_NAME"

# Set password for pg_dump
export PGPASSWORD="$POSTGRES_PASSWORD"

# Create backup
echo "📥 Creating database dump..."
if pg_dump -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
    --verbose --clean --if-exists --create > "$BACKUP_FILE"; then
    echo "✅ Database dump created: $BACKUP_FILE"
    
    # Compress the backup
    echo "🗜️  Compressing backup..."
    if gzip "$BACKUP_FILE"; then
        echo "✅ Backup compressed: $BACKUP_COMPRESSED"
        FINAL_BACKUP_FILE="$BACKUP_COMPRESSED"
    else
        echo "⚠️  Compression failed, keeping uncompressed backup"
        FINAL_BACKUP_FILE="$BACKUP_FILE"
    fi
    
    # Show backup information
    echo "📊 Backup information:"
    echo "   📄 File: $FINAL_BACKUP_FILE"
    echo "   📏 Size: $(ls -lh "$FINAL_BACKUP_FILE" | awk '{print $5}')"
    echo "   🕒 Created: $(date)"
    
    # Clean up old backups (keep last 7 days by default)
    RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-7}"
    echo "🧹 Cleaning up backups older than $RETENTION_DAYS days..."
    find "$BACKUP_DIR" -name "backup_*.sql.gz" -mtime +$RETENTION_DAYS -delete 2>/dev/null || true
    find "$BACKUP_DIR" -name "backup_*.sql" -mtime +$RETENTION_DAYS -delete 2>/dev/null || true
    
    echo "🎉 Backup completed successfully!"
    exit 0
else
    echo "❌ Backup failed!"
    rm -f "$BACKUP_FILE" 2>/dev/null || true
    exit 1
fi