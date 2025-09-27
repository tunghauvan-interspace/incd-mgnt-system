#!/bin/bash

# Database Restore Script for Incident Management System
# Usage: ./restore-database.sh <backup_file> [--force]

set -e

# Configuration
POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
POSTGRES_USER="${POSTGRES_USER:-user}"
POSTGRES_DB="${POSTGRES_DB:-incidentdb}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-password}"

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "❌ Error: Backup file is required"
    echo "Usage: $0 <backup_file> [--force]"
    echo ""
    echo "Examples:"
    echo "  $0 ./backups/backup_20231127_143022.sql.gz"
    echo "  $0 ./backups/backup_20231127_143022.sql --force"
    exit 1
fi

BACKUP_FILE="$1"
FORCE_RESTORE="${2:-}"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Determine if file is compressed
IS_COMPRESSED=false
TEMP_SQL_FILE=""

if [[ "$BACKUP_FILE" == *.gz ]]; then
    IS_COMPRESSED=true
    TEMP_SQL_FILE="/tmp/restore_$(date +%s).sql"
    echo "🗜️  Decompressing backup file..."
    if gunzip -c "$BACKUP_FILE" > "$TEMP_SQL_FILE"; then
        echo "✅ Backup decompressed to: $TEMP_SQL_FILE"
        SQL_FILE="$TEMP_SQL_FILE"
    else
        echo "❌ Failed to decompress backup file"
        exit 1
    fi
else
    SQL_FILE="$BACKUP_FILE"
fi

echo "🔄 Starting database restore..."
echo "📁 Backup file: $BACKUP_FILE"
echo "🎯 Target database: $POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB"
echo "👤 User: $POSTGRES_USER"

# Warning about data loss
if [ "$FORCE_RESTORE" != "--force" ]; then
    echo ""
    echo "⚠️  WARNING: This operation will COMPLETELY REPLACE the existing database!"
    echo "⚠️  All current data will be LOST and cannot be recovered!"
    echo ""
    read -p "Are you sure you want to continue? (type 'yes' to confirm): " confirm
    
    if [ "$confirm" != "yes" ]; then
        echo "❌ Restore cancelled by user"
        [ -n "$TEMP_SQL_FILE" ] && rm -f "$TEMP_SQL_FILE"
        exit 1
    fi
fi

# Set password for psql
export PGPASSWORD="$POSTGRES_PASSWORD"

# Test database connection
echo "🔌 Testing database connection..."
if ! psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d postgres -c "\q" 2>/dev/null; then
    echo "❌ Error: Cannot connect to database server"
    [ -n "$TEMP_SQL_FILE" ] && rm -f "$TEMP_SQL_FILE"
    exit 1
fi
echo "✅ Database connection successful"

# Create a pre-restore backup (if database exists and has data)
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
PRE_RESTORE_BACKUP="./backups/pre_restore_backup_${TIMESTAMP}.sql"

echo "💾 Creating pre-restore backup..."
mkdir -p "./backups"
if pg_dump -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
    --clean --if-exists --create > "$PRE_RESTORE_BACKUP" 2>/dev/null; then
    echo "✅ Pre-restore backup created: $PRE_RESTORE_BACKUP"
    gzip "$PRE_RESTORE_BACKUP" && echo "✅ Pre-restore backup compressed"
else
    echo "⚠️  Pre-restore backup failed (database may not exist yet)"
fi

# Restore the database
echo "📥 Restoring database from backup..."
if psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d postgres \
    --quiet --file="$SQL_FILE" 2>/dev/null; then
    echo "✅ Database restored successfully!"
    
    # Verify the restore
    echo "🔍 Verifying database restore..."
    TABLE_COUNT=$(psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
        -t -c "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';" 2>/dev/null | tr -d ' ')
    
    if [ "$TABLE_COUNT" -gt 0 ]; then
        echo "✅ Database verification successful: $TABLE_COUNT tables found"
        echo ""
        echo "📊 Database tables:"
        psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
            -c "\dt" 2>/dev/null | grep -E "^ public" || true
        echo ""
        echo "🎉 Database restore completed successfully!"
    else
        echo "⚠️  Warning: No tables found in database after restore"
    fi
else
    echo "❌ Database restore failed!"
    [ -n "$TEMP_SQL_FILE" ] && rm -f "$TEMP_SQL_FILE"
    exit 1
fi

# Clean up temporary files
if [ -n "$TEMP_SQL_FILE" ]; then
    rm -f "$TEMP_SQL_FILE"
    echo "🧹 Temporary files cleaned up"
fi

echo ""
echo "✅ Restore operation completed successfully!"
echo "📄 Backup restored from: $BACKUP_FILE"
echo "🕒 Completed at: $(date)"

# Remind user about application restart
echo ""
echo "⚠️  Important: You may need to restart the incident management application"
echo "   to ensure it picks up the restored database state."