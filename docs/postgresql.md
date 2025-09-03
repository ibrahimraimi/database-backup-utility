# PostgreSQL Database Guide

This guide covers using the Database Backup Utility with PostgreSQL databases.

## Prerequisites

- PostgreSQL server running and accessible
- User with appropriate permissions for backup/restore operations
- PostgreSQL client tools installed (optional, for advanced operations)

## Connection Testing

### Basic Connection Test

```bash
./db-backup test \
  --db-type postgres \
  --host localhost \
  --port 5432 \
  --username postgres \
  --password mypassword \
  --database mydb
```

### Connection with Custom Port

```bash
./db-backup test \
  --db-type postgres \
  --host 192.168.1.100 \
  --port 5433 \
  --username backup_user \
  --password secure_password \
  --database production_db
```

### Connection String Method

```bash
./db-backup test \
  --db-type postgres \
  --connection-string "host=localhost port=5432 user=postgres password=mypassword dbname=mydb sslmode=disable"
```

## Backup Operations

### Full Database Backup

```bash
./db-backup backup \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --type full \
  --compress
```

### Selective Table Backup

```bash
./db-backup backup \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --tables "users,orders,products" \
  --compress
```

### Backup to Custom Location

```bash
./db-backup backup \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --storage local \
  --path /var/backups/postgresql \
  --compress
```

### Incremental Backup

```bash
./db-backup backup \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --type incremental \
  --compress
```

## Restore Operations

### Full Database Restore

```bash
./db-backup restore \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --file ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz
```

### Selective Table Restore

```bash
./db-backup restore \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --file ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz \
  --tables "users,orders"
```

### Restore with Drop Existing Tables

```bash
./db-backup restore \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --file ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz \
  --drop-existing
```

## User Permissions

### Required PostgreSQL Permissions

The PostgreSQL user needs the following permissions for backup operations:

```sql
-- For backup operations
GRANT CONNECT ON DATABASE mydb TO backup_user;
GRANT USAGE ON SCHEMA public TO backup_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO backup_user;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO backup_user;

-- For restore operations
GRANT CONNECT ON DATABASE mydb TO backup_user;
GRANT USAGE ON SCHEMA public TO backup_user;
GRANT CREATE ON SCHEMA public TO backup_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO backup_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO backup_user;
```

### Create Dedicated Backup User

```sql
-- Create backup user
CREATE USER backup_user WITH PASSWORD 'secure_backup_password';

-- Grant database permissions
GRANT CONNECT ON DATABASE mydb TO backup_user;
GRANT USAGE ON SCHEMA public TO backup_user;

-- Grant table permissions
GRANT SELECT ON ALL TABLES IN SCHEMA public TO backup_user;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO backup_user;

-- Grant default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO backup_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON SEQUENCES TO backup_user;
```

## Configuration Examples

### PostgreSQL-Specific Configuration

```yaml
# ~/.db-backup.yaml
log:
  level: "info"
  format: "json"

storage:
  type: "local"
  path: "/var/backups/postgresql"

# PostgreSQL-specific settings
database:
  postgres:
    default_port: 5432
    ssl_mode: "disable"
    connection_timeout: 30
    statement_timeout: 300
```

### Environment Variables

```bash
# PostgreSQL connection
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USERNAME=backup_user
export POSTGRES_PASSWORD=secure_backup_password
export POSTGRES_DATABASE=production_db
export POSTGRES_SSLMODE=disable
```

## Advanced Usage

### Automated Backup Script

```bash
#!/bin/bash
# postgres-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
BACKUP_DIR="/var/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
./db-backup backup \
  --db-type postgres \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress \
  --storage local \
  --path $BACKUP_DIR

# Check if backup was successful
if [ $? -eq 0 ]; then
    echo "PostgreSQL backup completed successfully at $(date)"

    # Optional: Clean up old backups (keep last 7 days)
    find $BACKUP_DIR -name "postgres_${DB_NAME}_*.sql.gz" -mtime +7 -delete

    # Optional: Send notification
    # curl -X POST -H 'Content-type: application/json' \
    #   --data '{"text":"PostgreSQL backup completed successfully"}' \
    #   $SLACK_WEBHOOK_URL
else
    echo "PostgreSQL backup failed at $(date)"
    exit 1
fi
```

### Cron Job Setup

```bash
# Add to crontab for daily backups at 2 AM
0 2 * * * /path/to/postgres-backup.sh

# Add to crontab for hourly incremental backups
0 * * * * /path/to/postgres-incremental-backup.sh
```

### Docker Usage

```bash
# Backup PostgreSQL running in Docker
./db-backup backup \
  --db-type postgres \
  --host localhost \
  --port 5432 \
  --username postgres \
  --password mypassword \
  --database mydb \
  --compress

# Using Docker Compose
docker-compose exec db-backup ./db-backup backup \
  --db-type postgres \
  --host postgres \
  --username postgres \
  --password mypassword \
  --database mydb \
  --compress
```

## SSL/TLS Configuration

### SSL Connection

```bash
# Test SSL connection
./db-backup test \
  --db-type postgres \
  --connection-string "host=localhost port=5432 user=postgres password=mypassword dbname=mydb sslmode=require"

# Backup with SSL
./db-backup backup \
  --db-type postgres \
  --connection-string "host=localhost port=5432 user=postgres password=mypassword dbname=mydb sslmode=require" \
  --compress
```

### SSL Modes

| Mode          | Description                                   |
| ------------- | --------------------------------------------- |
| `disable`     | No SSL connection                             |
| `allow`       | Try non-SSL first, then SSL                   |
| `prefer`      | Try SSL first, then non-SSL                   |
| `require`     | SSL required, but don't verify certificate    |
| `verify-ca`   | SSL required, verify certificate authority    |
| `verify-full` | SSL required, verify certificate and hostname |

## Troubleshooting

### Common Issues

#### Connection Refused

```bash
# Check if PostgreSQL is running
systemctl status postgresql

# Check if port is accessible
telnet localhost 5432

# Check PostgreSQL logs
tail -f /var/log/postgresql/postgresql-*.log
```

#### Permission Denied

```bash
# Verify user permissions
psql -h localhost -U backup_user -d mydb -c "\du backup_user"

# Test connection manually
psql -h localhost -U backup_user -d mydb
```

#### Backup File Issues

```bash
# Check backup file integrity
file ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz

# Test decompression
gunzip -t ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz

# View backup content (first few lines)
zcat ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz | head -20
```

### Performance Optimization

#### Large Database Backups

```bash
# For large databases, consider:
# 1. Use incremental backups
./db-backup backup --db-type postgres --type incremental

# 2. Backup specific tables only
./db-backup backup --db-type postgres --tables "important_table1,important_table2"

# 3. Use parallel processing (if supported)
# 4. Consider using PostgreSQL's native pg_dump for very large databases
```

#### Network Optimization

```bash
# For remote PostgreSQL servers, consider:
# 1. Using connection pooling (pgBouncer)
# 2. Compressing network traffic
# 3. Using dedicated backup network
```

## Best Practices

1. **Regular Backups**: Set up automated daily backups
2. **Test Restores**: Regularly test restore procedures
3. **Monitor Space**: Ensure sufficient disk space for backups
4. **Secure Credentials**: Use dedicated backup users with minimal permissions
5. **Backup Verification**: Verify backup integrity after creation
6. **Retention Policy**: Implement backup retention and cleanup policies
7. **SSL Security**: Use SSL connections for remote backups
8. **Documentation**: Document your backup and restore procedures
9. **Monitoring**: Set up monitoring and alerting for backup failures
10. **Schema Backups**: Consider separate schema and data backups for large databases

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [Cloud Storage](cloud-storage.md)
- [Troubleshooting](troubleshooting.md)
- [Examples](examples.md)
