# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the Database Backup Utility.

## Common Issues

### Build and Installation Issues

#### Go Module Issues

**Problem**: `go: module database-backup-utility: cannot find module providing package`

**Solution**:

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Rebuild
go build -o build/dbu .
```

#### Missing Dependencies

**Problem**: `missing go.sum entry for module`

**Solution**:

```bash
# Remove go.sum and regenerate
rm go.sum
go mod tidy
go build -o build/dbu .
```

#### CGO Issues (SQLite)

**Problem**: `CGO_ENABLED=0` errors with SQLite

**Solution**:

```bash
# Install SQLite development libraries
# Ubuntu/Debian
sudo apt-get install libsqlite3-dev

# CentOS/RHEL
sudo yum install sqlite-devel

# macOS
brew install sqlite

# Rebuild with CGO
CGO_ENABLED=1 go build -o build/dbu .
```

### Database Connection Issues

#### MySQL Connection Problems

**Problem**: `dial tcp: connection refused`

**Solutions**:

```bash
# Check if MySQL is running
systemctl status mysql
# or
systemctl status mysqld

# Check if port is accessible
telnet localhost 3306

# Check MySQL configuration
mysql --version
mysql -h localhost -u root -p -e "SELECT 1"

# Check firewall
sudo ufw status
sudo iptables -L
```

**Problem**: `Access denied for user`

**Solutions**:

```bash
# Check user permissions
mysql -u root -p -e "SELECT User, Host FROM mysql.user WHERE User='backup_user';"

# Grant necessary permissions
mysql -u root -p -e "GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON mydb.* TO 'backup_user'@'%';"

# Test connection manually
mysql -h localhost -u backup_user -p mydb
```

#### PostgreSQL Connection Problems

**Problem**: `dial tcp: connection refused`

**Solutions**:

```bash
# Check if PostgreSQL is running
systemctl status postgresql

# Check if port is accessible
telnet localhost 5432

# Check PostgreSQL configuration
psql --version
psql -h localhost -U postgres -d postgres -c "SELECT 1"

# Check pg_hba.conf
sudo cat /etc/postgresql/*/main/pg_hba.conf
```

**Problem**: `FATAL: password authentication failed`

**Solutions**:

```bash
# Check user permissions
psql -h localhost -U postgres -d postgres -c "\du backup_user"

# Grant necessary permissions
psql -h localhost -U postgres -d postgres -c "GRANT CONNECT ON DATABASE mydb TO backup_user;"
psql -h localhost -U postgres -d postgres -c "GRANT USAGE ON SCHEMA public TO backup_user;"
psql -h localhost -U postgres -d postgres -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO backup_user;"

# Test connection manually
psql -h localhost -U backup_user -d mydb
```

#### MongoDB Connection Problems

**Problem**: `dial tcp: connection refused`

**Solutions**:

```bash
# Check if MongoDB is running
systemctl status mongod

# Check if port is accessible
telnet localhost 27017

# Check MongoDB configuration
mongosh --version
mongosh --host localhost --port 27017

# Check MongoDB logs
tail -f /var/log/mongodb/mongod.log
```

**Problem**: `Authentication failed`

**Solutions**:

```bash
# Check user permissions
mongosh --host localhost --port 27017 -u admin -p admin_password --authenticationDatabase admin
use mydb
db.getUsers()

# Create backup user
use admin
db.createUser({
  user: "backup_user",
  pwd: "secure_backup_password",
  roles: [
    { role: "read", db: "mydb" }
  ]
})

# Test connection manually
mongosh --host localhost --port 27017 -u backup_user -p secure_backup_password --authenticationDatabase admin mydb
```

#### SQLite Connection Problems

**Problem**: `no such file or directory`

**Solutions**:

```bash
# Check if database file exists
ls -la /path/to/database.db

# Check file permissions
ls -la /path/to/database.db
chmod 644 /path/to/database.db

# Check directory permissions
ls -la /path/to/
chmod 755 /path/to/
```

**Problem**: `database is locked`

**Solutions**:

```bash
# Check for active connections
lsof /path/to/database.db

# Kill processes using the database
kill -9 <process_id>

# Check for WAL files
ls -la /path/to/database.db*

# Remove WAL files if safe
rm /path/to/database.db-wal
rm /path/to/database.db-shm
```

### Backup Issues

#### Backup File Creation Problems

**Problem**: `permission denied`

**Solutions**:

```bash
# Check backup directory permissions
ls -la /var/backups/
chmod 755 /var/backups/
chown backup_user:backup_group /var/backups/

# Check disk space
df -h /var/backups/

# Create backup directory
mkdir -p /var/backups/
```

**Problem**: `no space left on device`

**Solutions**:

```bash
# Check disk space
df -h

# Clean up old backups
find /var/backups/ -name "*.sql.gz" -mtime +7 -delete

# Compress existing backups
gzip /var/backups/*.sql

# Move backups to different location
mv /var/backups/* /mnt/backup-storage/
```

#### Backup Corruption Issues

**Problem**: `backup file is corrupted`

**Solutions**:

```bash
# Check file integrity
file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# Test decompression
gunzip -t ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# View backup content
zcat ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz | head -20

# Recreate backup
./dbu backup --db-type mysql --database mydb --compress
```

### Restore Issues

#### Restore Permission Problems

**Problem**: `permission denied during restore`

**Solutions**:

```bash
# Check database user permissions
mysql -u root -p -e "SHOW GRANTS FOR 'restore_user'@'%';"

# Grant necessary permissions
mysql -u root -p -e "GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, INDEX, LOCK TABLES ON mydb.* TO 'restore_user'@'%';"

# Test restore with different user
./dbu restore --db-type mysql --username root --password root_password --database mydb --file backup.sql.gz
```

#### Restore Data Issues

**Problem**: `data not restored correctly`

**Solutions**:

```bash
# Check backup file content
zcat backup.sql.gz | grep "INSERT INTO"

# Verify database state before restore
mysql -u root -p -e "SELECT COUNT(*) FROM mydb.users;"

# Restore with verbose logging
./dbu --log-level debug restore --db-type mysql --database mydb --file backup.sql.gz

# Check database state after restore
mysql -u root -p -e "SELECT COUNT(*) FROM mydb.users;"
```

### Cloud Storage Issues

#### AWS S3 Issues

**Problem**: `Access Denied`

**Solutions**:

```bash
# Check AWS credentials
aws sts get-caller-identity

# Check S3 permissions
aws s3 ls s3://your-backup-bucket

# Verify IAM policy
aws iam get-policy --policy-arn arn:aws:iam::account:policy/BackupPolicy

# Test S3 access
aws s3 cp test.txt s3://your-backup-bucket/
```

**Problem**: `NoSuchBucket`

**Solutions**:

```bash
# Check if bucket exists
aws s3 ls s3://your-backup-bucket

# Create bucket if it doesn't exist
aws s3 mb s3://your-backup-bucket

# Check bucket region
aws s3api get-bucket-location --bucket your-backup-bucket
```

#### Google Cloud Storage Issues

**Problem**: `Authentication failed`

**Solutions**:

```bash
# Check service account
gcloud auth list

# Check service account key
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
gcloud auth activate-service-account --key-file=$GOOGLE_APPLICATION_CREDENTIALS

# Test GCS access
gsutil ls gs://your-backup-bucket
```

#### Azure Blob Storage Issues

**Problem**: `Authentication failed`

**Solutions**:

```bash
# Check Azure credentials
az account show

# Check storage account
az storage account show --name your_storage_account

# Test blob access
az storage blob list --container-name your-backup-container
```

### Notification Issues

#### Slack Notification Problems

**Problem**: `webhook not working`

**Solutions**:

```bash
# Test webhook manually
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Test notification"}' \
  $SLACK_WEBHOOK_URL

# Check webhook URL format
echo $SLACK_WEBHOOK_URL

# Verify channel permissions
# Ensure webhook has permission to post to the channel
```

#### Discord Notification Problems

**Problem**: `webhook not working`

**Solutions**:

```bash
# Test webhook manually
curl -X POST -H 'Content-type: application/json' \
  --data '{"content":"Test notification"}' \
  $DISCORD_WEBHOOK_URL

# Check webhook URL format
echo $DISCORD_WEBHOOK_URL

# Verify channel permissions
# Ensure webhook has permission to post to the channel
```

## Debugging Techniques

### Enable Debug Logging

```bash
# Enable debug logging
./dbu --log-level debug backup --db-type mysql --database mydb

# Enable debug logging with text format
./dbu --log-level debug --log-format text backup --db-type mysql --database mydb
```

### Verbose Output

```bash
# Test connection with verbose output
./dbu --log-level debug test --db-type mysql --host localhost --username root --password mypass --database mydb

# Backup with verbose output
./dbu --log-level debug backup --db-type mysql --database mydb --compress
```

### Check Configuration

```bash
# Check configuration file
cat ~/.dbu.yaml

# Validate configuration
./dbu --config ~/.dbu.yaml test --db-type mysql --database mydb
```

### Network Diagnostics

```bash
# Test database connectivity
telnet localhost 3306  # MySQL
telnet localhost 5432  # PostgreSQL
telnet localhost 27017 # MongoDB

# Check DNS resolution
nslookup localhost
nslookup your-database-host

# Check routing
traceroute your-database-host
```

### File System Diagnostics

```bash
# Check disk space
df -h

# Check inode usage
df -i

# Check file permissions
ls -la /path/to/database.db
ls -la /var/backups/

# Check file system errors
dmesg | grep -i error
```

## Performance Issues

### Slow Backup Performance

**Solutions**:

```bash
# Use incremental backups
./dbu backup --db-type mysql --type incremental --database mydb

# Backup specific tables only
./dbu backup --db-type mysql --tables "important_table1,important_table2" --database mydb

# Use compression
./dbu backup --db-type mysql --compress --database mydb

# Optimize database before backup
mysql -u root -p -e "OPTIMIZE TABLE mydb.important_table;"
```

### Memory Issues

**Solutions**:

```bash
# Check memory usage
free -h
top -p $(pgrep dbu)

# Increase swap space
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Optimize database settings
mysql -u root -p -e "SET GLOBAL innodb_buffer_pool_size = 1G;"
```

### Network Issues

**Solutions**:

```bash
# Use local backup first, then upload to cloud
./dbu backup --db-type mysql --storage local --database mydb
./dbu backup --db-type mysql --storage cloud --database mydb

# Use compression to reduce network traffic
./dbu backup --db-type mysql --compress --database mydb

# Use parallel processing if available
./dbu backup --db-type mysql --parallel 4 --database mydb
```

## Recovery Procedures

### Database Recovery

#### MySQL Recovery

```bash
# Stop MySQL service
sudo systemctl stop mysql

# Restore from backup
./dbu restore --db-type mysql --database mydb --file backup.sql.gz

# Start MySQL service
sudo systemctl start mysql

# Verify data
mysql -u root -p -e "SELECT COUNT(*) FROM mydb.users;"
```

#### PostgreSQL Recovery

```bash
# Stop PostgreSQL service
sudo systemctl stop postgresql

# Restore from backup
./dbu restore --db-type postgres --database mydb --file backup.sql.gz

# Start PostgreSQL service
sudo systemctl start postgresql

# Verify data
psql -h localhost -U postgres -d mydb -c "SELECT COUNT(*) FROM users;"
```

#### MongoDB Recovery

```bash
# Stop MongoDB service
sudo systemctl stop mongod

# Restore from backup
./dbu restore --db-type mongodb --database mydb --file backup.bson.gz

# Start MongoDB service
sudo systemctl start mongod

# Verify data
mongosh --host localhost --port 27017 -u admin -p admin_password --authenticationDatabase admin mydb
db.users.countDocuments()
```

#### SQLite Recovery

```bash
# Stop application using the database
sudo systemctl stop your-app

# Restore from backup
./dbu restore --db-type sqlite --database /path/to/database.db --file backup.db.gz

# Start application
sudo systemctl start your-app

# Verify data
sqlite3 /path/to/database.db "SELECT COUNT(*) FROM users;"
```

### Backup Recovery

#### Corrupted Backup Recovery

```bash
# Try to repair corrupted backup
gunzip -c backup.sql.gz | head -100 > backup_partial.sql

# Extract specific tables
gunzip -c backup.sql.gz | grep -A 1000 "CREATE TABLE users" > users_table.sql

# Use database-specific repair tools
mysql -u root -p -e "REPAIR TABLE mydb.users;"
```

#### Missing Backup Recovery

```bash
# Check for alternative backups
find /var/backups/ -name "*mydb*" -type f

# Check cloud storage
aws s3 ls s3://your-backup-bucket/ --recursive | grep mydb

# Use database replication if available
mysql -u root -p -e "SHOW SLAVE STATUS\G"
```

## Getting Help

### Log Analysis

```bash
# Check application logs
tail -f /var/log/dbu.log

# Check system logs
journalctl -u dbu -f

# Check database logs
tail -f /var/log/mysql/error.log
tail -f /var/log/postgresql/postgresql-*.log
tail -f /var/log/mongodb/mongod.log
```

### Community Support

1. **Check the documentation** in the `/docs` directory
2. **Review the examples** in the `/examples` directory
3. **Check GitHub issues** for similar problems
4. **Create a new issue** with detailed information

### Reporting Issues

When reporting issues, include:

1. **Database type and version**
2. **Operating system and version**
3. **Error messages and logs**
4. **Configuration files (sanitized)**
5. **Steps to reproduce**
6. **Expected vs actual behavior**

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [MySQL Guide](mysql.md)
- [PostgreSQL Guide](postgresql.md)
- [MongoDB Guide](mongodb.md)
- [SQLite Guide](sqlite.md)
- [Cloud Storage](cloud-storage.md)
- [Notifications](notifications.md)
