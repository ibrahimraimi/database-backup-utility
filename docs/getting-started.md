# Getting Started

This guide will help you get up and running with the Database Backup Utility quickly.

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/ibrahimraimi/database-backup-utility.git
cd database-backup-utility

# Build the binary
make build

# Or build for all platforms
make build-all
```

### Option 2: Using Go

```bash
go install github.com/ibrahimraimi/database-backup-utility@latest
```

### Option 3: Docker

```bash
# Pull the image
docker pull ibrahimraimi/database-backup-utility:latest

# Or build locally
docker build -t dbu .
```

## Quick Start

### 1. Test Your Database Connection

Before creating backups, always test your database connection:

```bash
# MySQL
./dbu test --db-type mysql --host localhost --username root --password mypassword --database mydb

# PostgreSQL
./dbu test --db-type postgres --host localhost --username postgres --password mypassword --database mydb

# MongoDB
./dbu test --db-type mongodb --host localhost --username admin --password mypassword --database mydb

# SQLite
./dbu test --db-type sqlite --database /path/to/database.db
```

### 2. Create Your First Backup

```bash
# MySQL backup
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --type full \
  --compress

# PostgreSQL backup
./dbu backup \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --type full \
  --compress

# MongoDB backup
./dbu backup \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --type full \
  --compress

# SQLite backup
./dbu backup \
  --db-type sqlite \
  --database /path/to/database.db \
  --type full \
  --compress
```

### 3. Restore from Backup

```bash
# Restore MySQL backup
./dbu restore \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# Restore PostgreSQL backup
./dbu restore \
  --db-type postgres \
  --host localhost \
  --username postgres \
  --password mypassword \
  --database mydb \
  --file ./backups/postgres_mydb_full_2024-01-15_10-30-00.sql.gz

# Restore MongoDB backup
./dbu restore \
  --db-type mongodb \
  --host localhost \
  --username admin \
  --password mypassword \
  --database mydb \
  --file ./backups/mongodb_mydb_full_2024-01-15_10-30-00.bson.gz

# Restore SQLite backup
./dbu restore \
  --db-type sqlite \
  --database /path/to/database.db \
  --file ./backups/sqlite_database_full_2024-01-15_10-30-00.db.gz
```

## Configuration

### Basic Configuration

Create a configuration file at `~/dbu.yaml`:

```yaml
# Logging configuration
log:
  level: "info" # debug, info, warn, error
  format: "json" # json, text

# Storage configuration
storage:
  type: "local" # local, cloud
  path: "./backups"

# Notification configuration (optional)
notify:
  enabled: false
  type: "slack" # slack, discord
  webhook: "" # Webhook URL
  channel: "" # Channel name
```

### Environment Variables

Set up environment variables for sensitive data:

```bash
# Database credentials
export DB_HOST=localhost
export DB_USERNAME=root
export DB_PASSWORD=mypassword
export DB_NAME=mydb

# Cloud storage (if using)
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1
```

## Command Line Options

### Global Flags

| Flag           | Description                          | Default      |
| -------------- | ------------------------------------ | ------------ |
| `--config`     | Path to configuration file           | `~/dbu.yaml` |
| `--log-level`  | Log level (debug, info, warn, error) | `info`       |
| `--log-format` | Log format (json, text)              | `json`       |

### Database Connection Flags

| Flag                  | Description            | Required             |
| --------------------- | ---------------------- | -------------------- |
| `--db-type`           | Database type          | Yes                  |
| `--host`              | Database host          | Yes (except SQLite)  |
| `--port`              | Database port          | No                   |
| `--username`          | Database username      | Yes (except SQLite)  |
| `--password`          | Database password      | Yes (except SQLite)  |
| `--database`          | Database name          | Yes                  |
| `--connection-string` | Full connection string | Alternative to above |

### Backup Flags

| Flag         | Description                                         | Default     |
| ------------ | --------------------------------------------------- | ----------- |
| `--type`     | Backup type (full, incremental, differential)       | `full`      |
| `--compress` | Compress backup files                               | `true`      |
| `--tables`   | Comma-separated list of tables for selective backup | All tables  |
| `--storage`  | Storage type (local, cloud)                         | `local`     |
| `--path`     | Local storage path                                  | `./backups` |

### Restore Flags

| Flag              | Description                                          | Default    |
| ----------------- | ---------------------------------------------------- | ---------- |
| `--file`          | Path to backup file to restore                       | Required   |
| `--tables`        | Comma-separated list of tables for selective restore | All tables |
| `--drop-existing` | Drop existing tables before restore                  | `false`    |

## Next Steps

1. **Choose your database**: Follow the specific guide for your database type:

   - [MySQL Guide](mysql.md)
   - [PostgreSQL Guide](postgresql.md)
   - [MongoDB Guide](mongodb.md)
   - [SQLite Guide](sqlite.md)

2. **Set up advanced features**:

   - [Cloud Storage](cloud-storage.md) - Store backups in the cloud
   - [Notifications](notifications.md) - Get notified of backup status
   - [Configuration](configuration.md) - Advanced configuration options

3. **Explore examples**:
   - [Examples](examples.md) - Real-world usage scenarios

## Troubleshooting

If you encounter issues:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Verify your database connection with the `test` command
3. Check the logs for detailed error messages
4. Ensure you have the necessary permissions for database operations

## Getting Help

- Use `./dbu --help` for general help
- Use `./dbu <command> --help` for command-specific help
- Check the logs for detailed error information
- Review the troubleshooting guide for common issues
