# MySQL Database Guide

This guide covers using the Database Backup Utility with MySQL databases.

## Prerequisites

- MySQL server running and accessible
- User with appropriate permissions for backup/restore operations
- MySQL client tools installed (optional, for advanced operations)

## Connection Testing

### Basic Connection Test

```bash
./dbu test \
  --db-type mysql \
  --host localhost \
  --port 3306 \
  --username root \
  --password mypassword \
  --database mydb
```

### Connection with Custom Port

```bash
./dbu test \
  --db-type mysql \
  --host 192.168.1.100 \
  --port 3307 \
  --username backup_user \
  --password secure_password \
  --database production_db
```

### Connection String Method

```bash
./dbu test \
  --db-type mysql \
  --connection-string "root:mypassword@tcp(localhost:3306)/mydb?parseTime=true"
```

## Backup Operations

### Full Database Backup

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --type full \
  --compress
```

### Selective Table Backup

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --tables "users,orders,products" \
  --compress
```

### Backup to Custom Location

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --storage local \
  --path /var/backups/mysql \
  --compress
```

### Incremental Backup

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --type incremental \
  --compress
```

## Restore Operations

### Full Database Restore

```bash
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz
```

### Selective Table Restore

```bash
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz \
  --tables "users,orders"
```

### Restore with Drop Existing Tables

```bash
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz \
  --drop-existing
```

## User Permissions

### Required MySQL Permissions

The MySQL user needs the following permissions for backup operations:

```sql
-- For backup operations
GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON mydb.* TO 'backup_user'@'%';

-- For restore operations
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, INDEX, LOCK TABLES ON mydb.* TO 'backup_user'@'%';

-- For global operations (if needed)
GRANT PROCESS, REPLICATION CLIENT ON *.* TO 'backup_user'@'%';

FLUSH PRIVILEGES;
```

### Create Dedicated Backup User

```sql
-- Create backup user
CREATE USER 'backup_user'@'%' IDENTIFIED BY 'secure_backup_password';

-- Grant necessary permissions
GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON *.* TO 'backup_user'@'%';
GRANT PROCESS, REPLICATION CLIENT ON *.* TO 'backup_user'@'%';

FLUSH PRIVILEGES;
```

## Configuration Examples

### MySQL-Specific Configuration

```yaml
# ~/.dbu.yaml
log:
  level: "info"
  format: "json"

storage:
  type: "local"
  path: "/var/backups/mysql"

# MySQL-specific settings
database:
  mysql:
    default_port: 3306
    connection_timeout: 30
    read_timeout: 60
    write_timeout: 60
```

### Environment Variables

```bash
# MySQL connection
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USERNAME=backup_user
export MYSQL_PASSWORD=secure_backup_password
export MYSQL_DATABASE=production_db
```

## Advanced Usage

### Automated Backup Script

```bash
#!/bin/bash
# mysql-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
BACKUP_DIR="/var/backups/mysql"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
./dbu backup \
  --db-type mysql \
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
    echo "MySQL backup completed successfully at $(date)"

    # Optional: Clean up old backups (keep last 7 days)
    find $BACKUP_DIR -name "mysql_${DB_NAME}_*.sql.gz" -mtime +7 -delete

    # Optional: Send notification
    # curl -X POST -H 'Content-type: application/json' \
    #   --data '{"text":"MySQL backup completed successfully"}' \
    #   $SLACK_WEBHOOK_URL
else
    echo "MySQL backup failed at $(date)"
    exit 1
fi
```

### Cron Job Setup

```bash
# Add to crontab for daily backups at 2 AM
0 2 * * * /path/to/mysql-backup.sh

# Add to crontab for hourly incremental backups
0 * * * * /path/to/mysql-incremental-backup.sh
```

### Docker Usage

```bash
# Backup MySQL running in Docker
./dbu backup \
  --db-type mysql \
  --host localhost \
  --port 3306 \
  --username root \
  --password mypassword \
  --database mydb \
  --compress

# Using Docker Compose
docker-compose exec dbu ./dbu backup \
  --db-type mysql \
  --host mysql \
  --username root \
  --password mypassword \
  --database mydb \
  --compress
```

## Troubleshooting

### Common Issues

#### Connection Refused

```bash
# Check if MySQL is running
systemctl status mysql

# Check if port is accessible
telnet localhost 3306

# Check MySQL error logs
tail -f /var/log/mysql/error.log
```

#### Permission Denied

```bash
# Verify user permissions
mysql -u backup_user -p -e "SHOW GRANTS FOR 'backup_user'@'%';"

# Test connection manually
mysql -h localhost -u backup_user -p mydb
```

#### Backup File Issues

```bash
# Check backup file integrity
file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# Test decompression
gunzip -t ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# View backup content (first few lines)
zcat ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz | head -20
```

### Performance Optimization

#### Large Database Backups

```bash
# For large databases, consider:
# 1. Use incremental backups
./dbu backup --db-type mysql --type incremental

# 2. Backup specific tables only
./dbu backup --db-type mysql --tables "important_table1,important_table2"

# 3. Use parallel processing (if supported)
# 4. Consider using MySQL's native tools for very large databases
```

#### Network Optimization

```bash
# For remote MySQL servers, consider:
# 1. Using connection pooling
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
7. **Documentation**: Document your backup and restore procedures
8. **Monitoring**: Set up monitoring and alerting for backup failures

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [Cloud Storage](cloud-storage.md)
- [Troubleshooting](troubleshooting.md)
- [Examples](examples.md)
