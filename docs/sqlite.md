# SQLite Database Guide

This guide covers using the Database Backup Utility with SQLite databases.

## Prerequisites

- SQLite database file accessible
- Read/write permissions to the database file and backup directory
- SQLite3 command-line tools installed (optional, for advanced operations)

## Connection Testing

### Basic Connection Test

```bash
./db-backup test \
  --db-type sqlite \
  --database /path/to/database.db
```

### Connection with Relative Path

```bash
./db-backup test \
  --db-type sqlite \
  --database ./data/app.db
```

### Connection String Method

```bash
./db-backup test \
  --db-type sqlite \
  --connection-string "/path/to/database.db"
```

## Backup Operations

### Full Database Backup

```bash
./db-backup backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --type full \
  --compress
```

### Backup to Custom Location

```bash
./db-backup backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --storage local \
  --path /var/backups/sqlite \
  --compress
```

### Incremental Backup

```bash
./db-backup backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --type incremental \
  --compress
```

## Restore Operations

### Full Database Restore

```bash
./db-backup restore \
  --db-type sqlite \
  --database /path/to/database.db \
  --file ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz
```

### Restore to Different Location

```bash
./db-backup restore \
  --db-type sqlite \
  --database /path/to/new_database.db \
  --file ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz
```

## Configuration Examples

### SQLite-Specific Configuration

```yaml
# ~/.db-backup.yaml
log:
  level: "info"
  format: "json"

storage:
  type: "local"
  path: "/var/backups/sqlite"

# SQLite-specific settings
database:
  sqlite:
    journal_mode: "WAL"
    synchronous: "NORMAL"
    cache_size: 1000
    temp_store: "MEMORY"
```

### Environment Variables

```bash
# SQLite database path
export SQLITE_DATABASE="/path/to/database.db"
export SQLITE_BACKUP_DIR="/var/backups/sqlite"
```

## Advanced Usage

### Automated Backup Script

```bash
#!/bin/bash
# sqlite-backup.sh

# Configuration
DB_PATH="/path/to/database.db"
BACKUP_DIR="/var/backups/sqlite"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
./db-backup backup \
  --db-type sqlite \
  --database $DB_PATH \
  --type full \
  --compress \
  --storage local \
  --path $BACKUP_DIR

# Check if backup was successful
if [ $? -eq 0 ]; then
    echo "SQLite backup completed successfully at $(date)"

    # Optional: Clean up old backups (keep last 7 days)
    find $BACKUP_DIR -name "sqlite_*.db.gz" -mtime +7 -delete

    # Optional: Send notification
    # curl -X POST -H 'Content-type: application/json' \
    #   --data '{"text":"SQLite backup completed successfully"}' \
    #   $SLACK_WEBHOOK_URL
else
    echo "SQLite backup failed at $(date)"
    exit 1
fi
```

### Multiple Database Backup Script

```bash
#!/bin/bash
# sqlite-multi-backup.sh

# Configuration
BACKUP_DIR="/var/backups/sqlite"
DATABASES=(
    "/path/to/database1.db"
    "/path/to/database2.db"
    "/path/to/database3.db"
)

# Backup each database
for db in "${DATABASES[@]}"; do
    echo "Backing up: $db"

    ./db-backup backup \
        --db-type sqlite \
        --database "$db" \
        --type full \
        --compress \
        --storage local \
        --path $BACKUP_DIR

    if [ $? -eq 0 ]; then
        echo "Successfully backed up: $db"
    else
        echo "Failed to backup: $db"
    fi
done
```

### Cron Job Setup

```bash
# Add to crontab for daily backups at 2 AM
0 2 * * * /path/to/sqlite-backup.sh

# Add to crontab for hourly incremental backups
0 * * * * /path/to/sqlite-incremental-backup.sh
```

### Docker Usage

```bash
# Backup SQLite database in Docker container
docker run --rm \
  -v /path/to/database:/data \
  -v /path/to/backups:/backups \
  db-backup:latest \
  backup \
  --db-type sqlite \
  --database /data/database.db \
  --storage local \
  --path /backups \
  --compress
```

## SQLite-Specific Features

### WAL Mode Support

```bash
# Enable WAL mode for better concurrency
sqlite3 /path/to/database.db "PRAGMA journal_mode=WAL;"

# Backup with WAL mode
./db-backup backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --compress
```

### Database Optimization

```bash
# Optimize database before backup
sqlite3 /path/to/database.db "VACUUM; ANALYZE;"

# Backup optimized database
./db-backup backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --compress
```

### Schema and Data Separation

```bash
# Backup schema only
sqlite3 /path/to/database.db ".schema" > schema.sql

# Backup data only
sqlite3 /path/to/database.db ".dump --data-only" > data.sql

# Full backup
./db-backup backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --compress
```

## Troubleshooting

### Common Issues

#### Permission Denied

```bash
# Check file permissions
ls -la /path/to/database.db

# Fix permissions if needed
chmod 644 /path/to/database.db
chmod 755 /path/to/backup/directory
```

#### Database Locked

```bash
# Check for active connections
lsof /path/to/database.db

# Kill processes using the database
kill -9 <process_id>

# Or wait for locks to be released
```

#### Backup File Issues

```bash
# Check backup file integrity
file ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz

# Test decompression
gunzip -t ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz

# View backup content (first few lines)
zcat ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz | head -20
```

### Performance Optimization

#### Large Database Backups

```bash
# For large databases, consider:
# 1. Use incremental backups
./db-backup backup --db-type sqlite --type incremental

# 2. Optimize database before backup
sqlite3 /path/to/database.db "VACUUM; ANALYZE;"

# 3. Use WAL mode for better concurrency
sqlite3 /path/to/database.db "PRAGMA journal_mode=WAL;"
```

#### Disk Space Management

```bash
# Check disk space before backup
df -h /path/to/backup/directory

# Monitor backup size
du -h ./backups/sqlite_database_full_*.db.gz

# Clean up old backups
find /var/backups/sqlite -name "*.db.gz" -mtime +7 -delete
```

## Best Practices

1. **Regular Backups**: Set up automated daily backups
2. **Test Restores**: Regularly test restore procedures
3. **Monitor Space**: Ensure sufficient disk space for backups
4. **File Permissions**: Set appropriate permissions for database and backup files
5. **Backup Verification**: Verify backup integrity after creation
6. **Retention Policy**: Implement backup retention and cleanup policies
7. **WAL Mode**: Use WAL mode for better concurrency during backups
8. **Database Optimization**: Regularly vacuum and analyze databases
9. **Documentation**: Document your backup and restore procedures
10. **Monitoring**: Set up monitoring and alerting for backup failures
11. **Multiple Copies**: Keep backups in multiple locations
12. **Version Control**: Consider versioning your database schema separately

## Security Considerations

### File Permissions

```bash
# Secure database file
chmod 600 /path/to/database.db
chown app_user:app_group /path/to/database.db

# Secure backup directory
chmod 700 /var/backups/sqlite
chown backup_user:backup_group /var/backups/sqlite
```

### Backup Encryption

```bash
# Encrypt backup files
gpg --symmetric --cipher-algo AES256 ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz

# Decrypt backup files
gpg --decrypt ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz.gpg > backup.db.gz
```

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [Cloud Storage](cloud-storage.md)
- [Troubleshooting](troubleshooting.md)
- [Examples](examples.md)
