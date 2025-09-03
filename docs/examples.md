# Examples and Use Cases

This guide provides real-world examples and use cases for the Database Backup Utility.

## Basic Examples

### Simple MySQL Backup

```bash
#!/bin/bash
# simple-mysql-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="root"
DB_PASS="mypassword"
DB_NAME="mydb"

# Create backup
./db-backup backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress

echo "Backup completed!"
```

### Simple PostgreSQL Backup

```bash
#!/bin/bash
# simple-postgres-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="postgres"
DB_PASS="mypassword"
DB_NAME="mydb"

# Create backup
./db-backup backup \
  --db-type postgres \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress

echo "Backup completed!"
```

### Simple MongoDB Backup

```bash
#!/bin/bash
# simple-mongodb-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="admin"
DB_PASS="mypassword"
DB_NAME="mydb"

# Create backup
./db-backup backup \
  --db-type mongodb \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress

echo "Backup completed!"
```

### Simple SQLite Backup

```bash
#!/bin/bash
# simple-sqlite-backup.sh

# Configuration
DB_PATH="/path/to/database.db"

# Create backup
./db-backup backup \
  --db-type sqlite \
  --database $DB_PATH \
  --type full \
  --compress

echo "Backup completed!"
```

## Production Examples

### Production MySQL Backup with Cleanup

```bash
#!/bin/bash
# production-mysql-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
BACKUP_DIR="/var/backups/mysql"
RETENTION_DAYS=7
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create backup
echo "Starting MySQL backup for $DB_NAME at $(date)"
./db-backup backup \
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

    # Clean up old backups
    echo "Cleaning up backups older than $RETENTION_DAYS days"
    find $BACKUP_DIR -name "mysql_${DB_NAME}_*.sql.gz" -mtime +$RETENTION_DAYS -delete

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ MySQL backup completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "MySQL backup failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ MySQL backup failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### Production PostgreSQL Backup with Verification

```bash
#!/bin/bash
# production-postgres-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
BACKUP_DIR="/var/backups/postgresql"
RETENTION_DAYS=7
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create backup
echo "Starting PostgreSQL backup for $DB_NAME at $(date)"
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

    # Verify backup integrity
    BACKUP_FILE=$(find $BACKUP_DIR -name "postgres_${DB_NAME}_*.sql.gz" -newer $BACKUP_DIR -type f | head -1)
    if [ -f "$BACKUP_FILE" ]; then
        echo "Verifying backup integrity: $BACKUP_FILE"
        gunzip -t "$BACKUP_FILE"
        if [ $? -eq 0 ]; then
            echo "Backup integrity verified successfully"
        else
            echo "Backup integrity check failed"
            exit 1
        fi
    fi

    # Clean up old backups
    echo "Cleaning up backups older than $RETENTION_DAYS days"
    find $BACKUP_DIR -name "postgres_${DB_NAME}_*.sql.gz" -mtime +$RETENTION_DAYS -delete

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ PostgreSQL backup completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "PostgreSQL backup failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ PostgreSQL backup failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### Multi-Database Backup Script

```bash
#!/bin/bash
# multi-database-backup.sh

# Configuration
BACKUP_DIR="/var/backups"
RETENTION_DAYS=7
DATE=$(date +%Y%m%d_%H%M%S)

# Database configurations
declare -A DATABASES=(
    ["mysql_prod"]="mysql:localhost:backup_user:secure_password:production_db"
    ["postgres_prod"]="postgres:localhost:backup_user:secure_password:production_db"
    ["mongodb_prod"]="mongodb:localhost:backup_user:secure_password:production_db"
)

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup each database
for db_name in "${!DATABASES[@]}"; do
    IFS=':' read -r db_type host username password database <<< "${DATABASES[$db_name]}"

    echo "Starting backup for $db_name ($db_type) at $(date)"

    ./db-backup backup \
        --db-type $db_type \
        --host $host \
        --username $username \
        --password $password \
        --database $database \
        --type full \
        --compress \
        --storage local \
        --path $BACKUP_DIR

    if [ $? -eq 0 ]; then
        echo "Backup completed successfully for $db_name"
    else
        echo "Backup failed for $db_name"
        # Send failure notification
        curl -X POST -H 'Content-type: application/json' \
          --data "{\"text\":\"❌ Backup failed for $db_name\"}" \
          $SLACK_WEBHOOK_URL
    fi
done

# Clean up old backups
echo "Cleaning up backups older than $RETENTION_DAYS days"
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "*.bson.gz" -mtime +$RETENTION_DAYS -delete

echo "Multi-database backup completed at $(date)"
```

## Cloud Storage Examples

### AWS S3 Backup

```bash
#!/bin/bash
# aws-s3-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
S3_BUCKET="my-backup-bucket"
S3_REGION="us-east-1"

# Create backup and upload to S3
echo "Starting backup to S3 for $DB_NAME at $(date)"
./db-backup backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress \
  --storage cloud \
  --cloud-provider aws \
  --bucket $S3_BUCKET \
  --region $S3_REGION

if [ $? -eq 0 ]; then
    echo "Backup to S3 completed successfully at $(date)"

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ Backup to S3 completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "Backup to S3 failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ Backup to S3 failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### Google Cloud Storage Backup

```bash
#!/bin/bash
# gcp-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
GCP_BUCKET="my-backup-bucket"
GCP_REGION="us-central1"

# Set up Google Cloud credentials
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"

# Create backup and upload to GCS
echo "Starting backup to GCS for $DB_NAME at $(date)"
./db-backup backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress \
  --storage cloud \
  --cloud-provider gcp \
  --bucket $GCP_BUCKET \
  --region $GCP_REGION

if [ $? -eq 0 ]; then
    echo "Backup to GCS completed successfully at $(date)"

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ Backup to GCS completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "Backup to GCS failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ Backup to GCS failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### Azure Blob Storage Backup

```bash
#!/bin/bash
# azure-backup.sh

# Configuration
DB_HOST="localhost"
DB_USER="backup_user"
DB_PASS="secure_backup_password"
DB_NAME="production_db"
AZURE_CONTAINER="my-backup-container"
AZURE_REGION="eastus"

# Set up Azure credentials
export AZURE_STORAGE_ACCOUNT="my_storage_account"
export AZURE_STORAGE_KEY="my_storage_key"

# Create backup and upload to Azure
echo "Starting backup to Azure for $DB_NAME at $(date)"
./db-backup backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress \
  --storage cloud \
  --cloud-provider azure \
  --bucket $AZURE_CONTAINER \
  --region $AZURE_REGION

if [ $? -eq 0 ]; then
    echo "Backup to Azure completed successfully at $(date)"

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ Backup to Azure completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "Backup to Azure failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ Backup to Azure failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

## Restore Examples

### MySQL Restore

```bash
#!/bin/bash
# mysql-restore.sh

# Configuration
DB_HOST="localhost"
DB_USER="restore_user"
DB_PASS="secure_restore_password"
DB_NAME="production_db"
BACKUP_FILE="/var/backups/mysql/mysql_production_db_full_2024-01-15_10-30-00.sql.gz"

# Restore database
echo "Starting MySQL restore for $DB_NAME at $(date)"
./db-backup restore \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --file $BACKUP_FILE \
  --drop-existing

if [ $? -eq 0 ]; then
    echo "MySQL restore completed successfully at $(date)"

    # Verify restore
    mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "SELECT COUNT(*) FROM $DB_NAME.users;"

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ MySQL restore completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "MySQL restore failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ MySQL restore failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### PostgreSQL Restore

```bash
#!/bin/bash
# postgres-restore.sh

# Configuration
DB_HOST="localhost"
DB_USER="restore_user"
DB_PASS="secure_restore_password"
DB_NAME="production_db"
BACKUP_FILE="/var/backups/postgresql/postgres_production_db_full_2024-01-15_10-30-00.sql.gz"

# Restore database
echo "Starting PostgreSQL restore for $DB_NAME at $(date)"
./db-backup restore \
  --db-type postgres \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --file $BACKUP_FILE \
  --drop-existing

if [ $? -eq 0 ]; then
    echo "PostgreSQL restore completed successfully at $(date)"

    # Verify restore
    psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) FROM users;"

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ PostgreSQL restore completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "PostgreSQL restore failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ PostgreSQL restore failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### MongoDB Restore

```bash
#!/bin/bash
# mongodb-restore.sh

# Configuration
DB_HOST="localhost"
DB_USER="restore_user"
DB_PASS="secure_restore_password"
DB_NAME="production_db"
BACKUP_FILE="/var/backups/mongodb/mongodb_production_db_full_2024-01-15_10-30-00.bson.gz"

# Restore database
echo "Starting MongoDB restore for $DB_NAME at $(date)"
./db-backup restore \
  --db-type mongodb \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --file $BACKUP_FILE \
  --drop-existing

if [ $? -eq 0 ]; then
    echo "MongoDB restore completed successfully at $(date)"

    # Verify restore
    mongosh --host $DB_HOST --port 27017 -u $DB_USER -p $DB_PASS --authenticationDatabase admin $DB_NAME --eval "db.users.countDocuments()"

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ MongoDB restore completed successfully for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "MongoDB restore failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ MongoDB restore failed for $DB_NAME\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

### SQLite Restore

```bash
#!/bin/bash
# sqlite-restore.sh

# Configuration
DB_PATH="/path/to/database.db"
BACKUP_FILE="/var/backups/sqlite/sqlite_database_full_2024-01-15_10-30-00.db.gz"

# Stop application using the database
echo "Stopping application..."
sudo systemctl stop your-app

# Restore database
echo "Starting SQLite restore at $(date)"
./db-backup restore \
  --db-type sqlite \
  --database $DB_PATH \
  --file $BACKUP_FILE

if [ $? -eq 0 ]; then
    echo "SQLite restore completed successfully at $(date)"

    # Verify restore
    sqlite3 $DB_PATH "SELECT COUNT(*) FROM users;"

    # Start application
    echo "Starting application..."
    sudo systemctl start your-app

    # Send success notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"✅ SQLite restore completed successfully\"}" \
      $SLACK_WEBHOOK_URL
else
    echo "SQLite restore failed at $(date)"

    # Send failure notification
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ SQLite restore failed\"}" \
      $SLACK_WEBHOOK_URL

    exit 1
fi
```

## Cron Job Examples

### Daily Backup Cron Job

```bash
# Add to crontab for daily backups at 2 AM
0 2 * * * /path/to/production-mysql-backup.sh

# Add to crontab for daily PostgreSQL backups at 3 AM
0 3 * * * /path/to/production-postgres-backup.sh

# Add to crontab for daily MongoDB backups at 4 AM
0 4 * * * /path/to/production-mongodb-backup.sh
```

### Hourly Incremental Backup Cron Job

```bash
# Add to crontab for hourly incremental backups
0 * * * * /path/to/incremental-backup.sh
```

### Weekly Full Backup Cron Job

```bash
# Add to crontab for weekly full backups on Sunday at 1 AM
0 1 * * 0 /path/to/weekly-full-backup.sh
```

## Docker Examples

### Docker Compose Backup

```yaml
# docker-compose.yml
version: "3.8"

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: mydb
    volumes:
      - mysql_data:/var/lib/mysql

  db-backup:
    build: .
    depends_on:
      - mysql
    volumes:
      - ./backups:/backups
    environment:
      - DB_HOST=mysql
      - DB_USER=root
      - DB_PASS=password
      - DB_NAME=mydb
    command: >
      sh -c "
        ./db-backup backup \
          --db-type mysql \
          --host mysql \
          --username root \
          --password password \
          --database mydb \
          --compress
      "

volumes:
  mysql_data:
```

### Docker Backup Script

```bash
#!/bin/bash
# docker-backup.sh

# Configuration
CONTAINER_NAME="mysql-container"
DB_USER="root"
DB_PASS="password"
DB_NAME="mydb"
BACKUP_DIR="./backups"

# Create backup directory
mkdir -p $BACKUP_DIR

# Create backup using Docker
docker run --rm \
  --network host \
  -v $(pwd)/backups:/backups \
  db-backup:latest \
  backup \
  --db-type mysql \
  --host localhost \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --storage local \
  --path /backups \
  --compress

echo "Docker backup completed!"
```

## Monitoring Examples

### Backup Monitoring Script

```bash
#!/bin/bash
# backup-monitor.sh

# Configuration
BACKUP_DIR="/var/backups"
ALERT_DAYS=2

# Check if backups are being created
echo "Checking backup status..."

# Check MySQL backups
MYSQL_BACKUPS=$(find $BACKUP_DIR -name "mysql_*.sql.gz" -mtime -$ALERT_DAYS | wc -l)
if [ $MYSQL_BACKUPS -eq 0 ]; then
    echo "⚠️ No MySQL backups found in the last $ALERT_DAYS days"
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"⚠️ No MySQL backups found in the last $ALERT_DAYS days\"}" \
      $SLACK_WEBHOOK_URL
fi

# Check PostgreSQL backups
POSTGRES_BACKUPS=$(find $BACKUP_DIR -name "postgres_*.sql.gz" -mtime -$ALERT_DAYS | wc -l)
if [ $POSTGRES_BACKUPS -eq 0 ]; then
    echo "⚠️ No PostgreSQL backups found in the last $ALERT_DAYS days"
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"⚠️ No PostgreSQL backups found in the last $ALERT_DAYS days\"}" \
      $SLACK_WEBHOOK_URL
fi

# Check MongoDB backups
MONGO_BACKUPS=$(find $BACKUP_DIR -name "mongodb_*.bson.gz" -mtime -$ALERT_DAYS | wc -l)
if [ $MONGO_BACKUPS -eq 0 ]; then
    echo "⚠️ No MongoDB backups found in the last $ALERT_DAYS days"
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"⚠️ No MongoDB backups found in the last $ALERT_DAYS days\"}" \
      $SLACK_WEBHOOK_URL
fi

echo "Backup monitoring completed!"
```

### Health Check Script

```bash
#!/bin/bash
# health-check.sh

# Configuration
DB_HOST="localhost"
DB_USER="health_check_user"
DB_PASS="health_check_password"
DB_NAME="health_check_db"

# Test database connection
echo "Testing database connection..."
./db-backup test \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME

if [ $? -eq 0 ]; then
    echo "✅ Database connection test passed"
else
    echo "❌ Database connection test failed"
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"❌ Database connection test failed\"}" \
      $SLACK_WEBHOOK_URL
    exit 1
fi

echo "Health check completed!"
```

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [MySQL Guide](mysql.md)
- [PostgreSQL Guide](postgresql.md)
- [MongoDB Guide](mongodb.md)
- [SQLite Guide](sqlite.md)
- [Cloud Storage](cloud-storage.md)
- [Notifications](notifications.md)
- [Troubleshooting](troubleshooting.md)
