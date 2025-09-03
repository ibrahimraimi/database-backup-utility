# Database Backup Utility

A comprehensive command-line interface (CLI) utility for backing up and restoring any type of database. The utility supports various database management systems (DBMS) such as MySQL, PostgreSQL, MongoDB, SQLite, and others. This tool features automatic backup scheduling, compression of backup files, storage options (local and cloud), and logging of backup activities.

## Features

✅ **Multi-Database Support**: MySQL, PostgreSQL, MongoDB, SQLite  
✅ **Backup Types**: Full, incremental, and differential backups  
✅ **Compression**: Built-in gzip compression to reduce storage space  
✅ **Storage Options**: Local storage and cloud storage (AWS S3, Google Cloud Storage, Azure Blob Storage)  
✅ **Restore Operations**: Full and selective restore capabilities  
✅ **Logging**: Comprehensive logging with configurable levels and formats  
✅ **Notifications**: Slack and Discord notifications for backup status  
✅ **Cross-Platform**: Works on Windows, Linux, and macOS

## Installation

### Download Pre-built Binaries (Recommended)

Download the latest release for your platform:

**Linux (AMD64):**

```bash
curl -L -o dbu https://github.com/ibrahimraimi/database-backup-utility/releases/latest/download/dbu-linux-amd64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L -o dbu https://github.com/ibrahimraimi/database-backup-utility/releases/latest/download/dbu-darwin-amd64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

**macOS (Apple Silicon):**

```bash
curl -L -o dbu https://github.com/ibrahimraimi/database-backup-utility/releases/latest/download/dbu-darwin-arm64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

**Windows (AMD64):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/ibrahimraimi/database-backup-utility/releases/latest/download/dbu-windows-amd64.exe" -OutFile "dbu.exe"
```

### From Source

```bash
# Clone the repository
git clone https://github.com/ibrahimraimi/database-backup-utility.git
cd database-backup-utility

# Build the binary
make build

# Or build for all platforms
make build-all
```

### Using Go

```bash
go install github.com/ibrahimraimi/database-backup-utility@latest
```

> **Note:** For detailed download instructions and release information, see the [Releases Documentation](docs/releases.md).

## Quick Start

### Test Database Connection

```bash
# Test MySQL connection
./dbu test --db-type mysql --host localhost --port 3306 --username root --password mypassword --database mydb

# Test PostgreSQL connection
./dbu test --db-type postgres --host localhost --port 5432 --username postgres --password mypassword --database mydb

# Test MongoDB connection
./dbu test --db-type mongodb --host localhost --port 27017 --username admin --password mypassword --database mydb

# Test SQLite connection
./dbu test --db-type sqlite --database /path/to/database.db
```

### Create a Backup

```bash
# MySQL backup
./dbu backup --db-type mysql --host localhost --username root --password mypassword --database mydb --type full --compress

# PostgreSQL backup
./dbu backup --db-type postgres --host localhost --username postgres --password mypassword --database mydb --type full --compress

# MongoDB backup
./dbu backup --db-type mongodb --host localhost --username admin --password mypassword --database mydb --type full --compress

# SQLite backup
./dbu backup --db-type sqlite --database /path/to/database.db --type full --compress

# Selective backup (specific tables)
./dbu backup --db-type mysql --host localhost --username root --password mypassword --database mydb --tables "users,orders,products" --compress

# Cloud storage backup
./dbu backup --db-type mysql --host localhost --username root --password mypassword --database mydb --storage cloud --cloud-provider aws --bucket my-backup-bucket --region us-east-1
```

### Restore a Backup

```bash
# Restore from local backup
./dbu restore --db-type mysql --host localhost --username root --password mypassword --database mydb --file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# Restore from cloud backup
./dbu restore --db-type mysql --host localhost --username root --password mypassword --database mydb --file s3://my-backup-bucket/mysql_mydb_full_2024-01-15_10-30-00.sql.gz

# Selective restore (specific tables)
./dbu restore --db-type mysql --host localhost --username root --password mypassword --database mydb --file ./backups/mysql_mydb_full_2024-01-15_10-30-00.sql.gz --tables "users,orders"
```

## Configuration

Create a configuration file at `~/.dbu.yaml`:

```yaml
# Logging configuration
log:
  level: "info" # debug, info, warn, error
  format: "json" # json, text

# Storage configuration
storage:
  type: "local" # local, cloud
  path: "./backups"

# Cloud storage configuration
cloud:
  provider: "aws" # aws, gcp, azure
  bucket: "my-backup-bucket"
  region: "us-east-1"

# Notification configuration
notify:
  enabled: true
  type: "slack" # slack, discord
  webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  channel: "#backups"
```

## Environment Variables

```bash
# AWS S3
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1

# Google Cloud Storage
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

# Azure Blob Storage
export AZURE_STORAGE_ACCOUNT=your_storage_account
export AZURE_STORAGE_KEY=your_storage_key
```

## Command Reference

### Global Flags

- `--config`: Path to configuration file (default: ~/.dbu.yaml)
- `--log-level`: Log level (debug, info, warn, error)
- `--log-format`: Log format (json, text)

### Database Connection Flags

- `--db-type`: Database type (mysql, postgres, mongodb, sqlite)
- `--host`: Database host
- `--port`: Database port
- `--username`: Database username
- `--password`: Database password
- `--database`: Database name
- `--connection-string`: Full database connection string

### Backup Flags

- `--type`: Backup type (full, incremental, differential)
- `--compress`: Compress backup files (default: true)
- `--tables`: Comma-separated list of tables for selective backup
- `--storage`: Storage type (local, cloud)
- `--path`: Local storage path
- `--cloud-provider`: Cloud provider (aws, gcp, azure)
- `--bucket`: Cloud storage bucket name
- `--region`: Cloud storage region

### Restore Flags

- `--file`: Path to backup file to restore
- `--tables`: Comma-separated list of tables for selective restore
- `--drop-existing`: Drop existing tables before restore

## Examples

### Automated Backup Script

```bash
#!/bin/bash
# daily-backup.sh

DB_HOST="localhost"
DB_USER="root"
DB_PASS="mypassword"
DB_NAME="mydb"
BACKUP_DIR="/var/backups/db"

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

# Upload to cloud
./dbu backup \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME \
  --type full \
  --compress \
  --storage cloud \
  --cloud-provider aws \
  --bucket my-backup-bucket \
  --region us-east-1
```

### Cron Job Setup

```bash
# Add to crontab for daily backups at 2 AM
0 2 * * * /path/to/daily-backup.sh
```

## Development

### Building from Source

```bash
# Install dependencies
make deps

# Build the binary
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint
```

### Project Structure

```
database-backup-utility/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── backup.go          # Backup command
│   ├── restore.go         # Restore command
│   └── test.go            # Test command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── database/          # Database connection management
│   ├── backup/            # Backup operations
│   ├── restore/           # Restore operations
│   ├── storage/           # Storage management
│   ├── notification/      # Notification system
│   └── logger/            # Logging utilities
├── main.go                # Application entry point
├── go.mod                 # Go module file
├── Makefile              # Build automation
└── README.md             # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Documentation

Comprehensive documentation is available in the `/docs` directory:

- [Getting Started](docs/getting-started.md) - Quick start guide and installation
- [Releases](docs/releases.md) - Download pre-built binaries and release information
- [MySQL Guide](docs/mysql.md) - Complete MySQL backup and restore guide
- [PostgreSQL Guide](docs/postgresql.md) - Complete PostgreSQL backup and restore guide
- [MongoDB Guide](docs/mongodb.md) - Complete MongoDB backup and restore guide
- [SQLite Guide](docs/sqlite.md) - Complete SQLite backup and restore guide
- [Configuration](docs/configuration.md) - Configuration file setup and options
- [Cloud Storage](docs/cloud-storage.md) - Cloud storage integration guide
- [Notifications](docs/notifications.md) - Slack and Discord notification setup
- [Troubleshooting](docs/troubleshooting.md) - Common issues and solutions
- [Examples](docs/examples.md) - Real-world usage examples and scripts

## Support

For support and questions:

- Create an issue on GitHub
- Check the [documentation](docs/README.md)
- Review the [examples](docs/examples.md)
- Check the [troubleshooting guide](docs/troubleshooting.md)
