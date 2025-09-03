# MongoDB Database Guide

This guide covers using the Database Backup Utility with MongoDB databases.

## Prerequisites

- MongoDB server running and accessible
- User with appropriate permissions for backup/restore operations
- MongoDB client tools installed (optional, for advanced operations)

## Connection Testing

### Basic Connection Test

```bash
./dbu test \
  --db-type mongodb \
  --host localhost \
  --port 27017 \
  --username admin \
  --password mypassword \
  --database mydb
```

### Connection with Custom Port

```bash
./dbu test \
  --db-type mongodb \
  --host 192.168.1.100 \
  --port 27018 \
  --username backup_user \
  --password secure_password \
  --database production_db
```

### Connection String Method

```bash
./dbu test \
  --db-type mongodb \
  --connection-string "mongodb://admin:mypassword@localhost:27017/mydb"
```

### Connection with Authentication Database

```bash
./dbu test \
  --db-type mongodb \
  --connection-string "mongodb://admin:mypassword@localhost:27017/mydb?authSource=admin"
```

## Backup Operations

### Full Database Backup

```bash
./dbu backup \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --type full \
  --compress
```

### Selective Collection Backup

```bash
./dbu backup \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --tables "users,orders,products" \
  --compress
```

### Backup to Custom Location

```bash
./dbu backup \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --storage local \
  --path /var/backups/mongodb \
  --compress
```

### Incremental Backup

```bash
./dbu backup \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --type incremental \
  --compress
```

## Restore Operations

### Full Database Restore

```bash
./dbu restore \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --file ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz
```

### Selective Collection Restore

```bash
./dbu restore \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --file ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz \
  --tables "users,orders"
```

### Restore with Drop Existing Collections

```bash
./dbu restore \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --file ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz \
  --drop-existing
```

## User Permissions

### Required MongoDB Permissions

The MongoDB user needs the following permissions for backup operations:

```javascript
// For backup operations
use mydb
db.createUser({
  user: "backup_user",
  pwd: "secure_backup_password",
  roles: [
    { role: "read", db: "mydb" }
  ]
})

// For restore operations
use mydb
db.createUser({
  user: "restore_user",
  pwd: "secure_restore_password",
  roles: [
    { role: "readWrite", db: "mydb" }
  ]
})
```

### Create Dedicated Backup User

```javascript
// Connect to MongoDB as admin
use admin

// Create backup user with read permissions
db.createUser({
  user: "backup_user",
  pwd: "secure_backup_password",
  roles: [
    { role: "read", db: "mydb" },
    { role: "read", db: "admin" }  // For authentication
  ]
})

// Create restore user with readWrite permissions
db.createUser({
  user: "restore_user",
  pwd: "secure_restore_password",
  roles: [
    { role: "readWrite", db: "mydb" },
    { role: "read", db: "admin" }  // For authentication
  ]
})
```

## Configuration Examples

### MongoDB-Specific Configuration

```yaml
# ~/.dbu.yaml
log:
  level: "info"
  format: "json"

storage:
  type: "local"
  path: "/var/backups/mongodb"

# MongoDB-specific settings
database:
  mongodb:
    default_port: 27017
    auth_source: "admin"
    connection_timeout: 30
    socket_timeout: 30
    server_selection_timeout: 30
```

### Environment Variables

```bash
# MongoDB connection
export MONGODB_HOST=localhost
export MONGODB_PORT=27017
export MONGODB_USERNAME=backup_user
export MONGODB_PASSWORD=secure_backup_password
export MONGODB_DATABASE=production_db
export MONGODB_AUTH_SOURCE=admin
```

## Advanced Usage

### Automated Backup Script

```bash
#!/bin/bash
# mongodbu.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
BACKUP_DIR="/var/backups/mongodb"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
./dbu backup \
  --db-type mongodb \
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
    echo "MongoDB backup completed successfully at $(date)"

    # Optional: Clean up old backups (keep last 7 days)
    find $BACKUP_DIR -name "mongodb_${DB_NAME}_*.bson.gz" -mtime +7 -delete

    # Optional: Send notification
    # curl -X POST -H 'Content-type: application/json' \
    #   --data '{"text":"MongoDB backup completed successfully"}' \
    #   $SLACK_WEBHOOK_URL
else
    echo "MongoDB backup failed at $(date)"
    exit 1
fi
```

### Cron Job Setup

```bash
# Add to crontab for daily backups at 2 AM
0 2 * * * /path/to/mongodbu.sh

# Add to crontab for hourly incremental backups
0 * * * * /path/to/mongodb-incremental-backup.sh
```

### Docker Usage

```bash
# Backup MongoDB running in Docker
./dbu backup \
  --db-type mongodb \
  --host localhost \
  --port 27017 \
  --username admin \
  --password mypassword \
  --database mydb \
  --compress

# Using Docker Compose
docker-compose exec dbu ./dbu backup \
  --db-type mongodb \
  --host mongodb \
  --username admin \
  --password mypassword \
  --database mydb \
  --compress
```

## MongoDB Replica Set Support

### Connection to Replica Set

```bash
# Test connection to replica set
./dbu test \
  --db-type mongodb \
  --connection-string "mongodb://admin:mypassword@mongodb1:27017,mongodb2:27017,mongodb3:27017/mydb?replicaSet=rs0"

# Backup from replica set
./dbu backup \
  --db-type mongodb \
  --connection-string "mongodb://admin:mypassword@mongodb1:27017,mongodb2:27017,mongodb3:27017/mydb?replicaSet=rs0" \
  --compress
```

### Sharded Cluster Support

```bash
# Test connection to sharded cluster
./dbu test \
  --db-type mongodb \
  --connection-string "mongodb://admin:mypassword@mongos1:27017,mongos2:27017/mydb"

# Backup from sharded cluster
./dbu backup \
  --db-type mongodb \
  --connection-string "mongodb://admin:mypassword@mongos1:27017,mongos2:27017/mydb" \
  --compress
```

## Troubleshooting

### Common Issues

#### Connection Refused

```bash
# Check if MongoDB is running
systemctl status mongod

# Check if port is accessible
telnet localhost 27017

# Check MongoDB logs
tail -f /var/log/mongodb/mongod.log
```

#### Authentication Failed

```bash
# Test authentication manually
mongo --host localhost --port 27017 -u backup_user -p secure_backup_password --authenticationDatabase admin

# Check user permissions
mongo --host localhost --port 27017 -u admin -p admin_password --authenticationDatabase admin
use mydb
db.getUsers()
```

#### Backup File Issues

```bash
# Check backup file integrity
file ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz

# Test decompression
gunzip -t ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz

# View backup content (first few lines)
zcat ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz | head -20
```

### Performance Optimization

#### Large Database Backups

```bash
# For large databases, consider:
# 1. Use incremental backups
./dbu backup --db-type mongodb --type incremental

# 2. Backup specific collections only
./dbu backup --db-type mongodb --tables "important_collection1,important_collection2"

# 3. Use MongoDB's native mongodump for very large databases
# 4. Consider using MongoDB's oplog for incremental backups
```

#### Network Optimization

```bash
# For remote MongoDB servers, consider:
# 1. Using connection pooling
# 2. Compressing network traffic
# 3. Using dedicated backup network
# 4. Using read preferences for replica sets
```

## Best Practices

1. **Regular Backups**: Set up automated daily backups
2. **Test Restores**: Regularly test restore procedures
3. **Monitor Space**: Ensure sufficient disk space for backups
4. **Secure Credentials**: Use dedicated backup users with minimal permissions
5. **Backup Verification**: Verify backup integrity after creation
6. **Retention Policy**: Implement backup retention and cleanup policies
7. **Replica Set Backups**: Use secondary members for backups to reduce primary load
8. **Sharded Cluster**: Consider backing up each shard separately for large clusters
9. **Documentation**: Document your backup and restore procedures
10. **Monitoring**: Set up monitoring and alerting for backup failures
11. **Oplog Backups**: For replica sets, consider backing up the oplog for point-in-time recovery

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [Cloud Storage](cloud-storage.md)
- [Troubleshooting](troubleshooting.md)
- [Examples](examples.md)
