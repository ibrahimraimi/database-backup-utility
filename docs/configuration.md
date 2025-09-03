# Configuration Guide

This guide covers all configuration options available in the Database Backup Utility.

## Configuration File

The utility uses YAML configuration files. The default location is `~/.db-backup.yaml`.

### Basic Configuration Structure

```yaml
# Logging configuration
log:
  level: "info" # debug, info, warn, error
  format: "json" # json, text

# Storage configuration
storage:
  type: "local" # local, cloud
  path: "./backups"

# Cloud storage configuration (when storage.type is "cloud")
cloud:
  provider: "aws" # aws, gcp, azure
  bucket: "my-backup-bucket"
  region: "us-east-1"
  access_key: "" # Set via environment variable
  secret_key: "" # Set via environment variable

# Notification configuration
notify:
  enabled: false
  type: "slack" # slack, discord
  webhook: "" # Webhook URL
  channel: "" # Channel name (for Slack)

# Database-specific configurations
database:
  mysql:
    default_port: 3306
    connection_timeout: 30
    read_timeout: 60
    write_timeout: 60

  postgres:
    default_port: 5432
    ssl_mode: "disable"
    connection_timeout: 30
    statement_timeout: 300

  mongodb:
    default_port: 27017
    auth_source: "admin"
    connection_timeout: 30
    socket_timeout: 30
    server_selection_timeout: 30

  sqlite:
    journal_mode: "WAL"
    synchronous: "NORMAL"
    cache_size: 1000
    temp_store: "MEMORY"
```

## Logging Configuration

### Log Levels

| Level   | Description                           |
| ------- | ------------------------------------- |
| `debug` | Detailed information for debugging    |
| `info`  | General information about operations  |
| `warn`  | Warning messages for potential issues |
| `error` | Error messages for failed operations  |

### Log Formats

| Format | Description                                              |
| ------ | -------------------------------------------------------- |
| `json` | Structured JSON format (recommended for production)      |
| `text` | Human-readable text format (recommended for development) |

### Example Logging Configuration

```yaml
log:
  level: "info"
  format: "json"
```

## Storage Configuration

### Local Storage

```yaml
storage:
  type: "local"
  path: "/var/backups/database-backup-utility"
```

### Cloud Storage

```yaml
storage:
  type: "cloud"
  path: "./temp" # Temporary local path for processing

cloud:
  provider: "aws" # aws, gcp, azure
  bucket: "my-backup-bucket"
  region: "us-east-1"
  access_key: "" # Set via AWS_ACCESS_KEY_ID
  secret_key: "" # Set via AWS_SECRET_ACCESS_KEY
```

## Cloud Storage Providers

### AWS S3

```yaml
cloud:
  provider: "aws"
  bucket: "my-backup-bucket"
  region: "us-east-1"
  access_key: "" # Set via AWS_ACCESS_KEY_ID
  secret_key: "" # Set via AWS_SECRET_ACCESS_KEY
```

### Google Cloud Storage

```yaml
cloud:
  provider: "gcp"
  bucket: "my-backup-bucket"
  region: "us-central1"
  access_key: "" # Set via GOOGLE_APPLICATION_CREDENTIALS
  secret_key: ""
```

### Azure Blob Storage

```yaml
cloud:
  provider: "azure"
  bucket: "my-backup-container"
  region: "eastus"
  access_key: "" # Set via AZURE_STORAGE_ACCOUNT
  secret_key: "" # Set via AZURE_STORAGE_KEY
```

## Notification Configuration

### Slack Notifications

```yaml
notify:
  enabled: true
  type: "slack"
  webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  channel: "#backups"
```

### Discord Notifications

```yaml
notify:
  enabled: true
  type: "discord"
  webhook: "https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"
  channel: "backups"
```

## Database-Specific Configuration

### MySQL Configuration

```yaml
database:
  mysql:
    default_port: 3306
    connection_timeout: 30
    read_timeout: 60
    write_timeout: 60
    charset: "utf8mb4"
    collation: "utf8mb4_unicode_ci"
    ssl_mode: "preferred"
```

### PostgreSQL Configuration

```yaml
database:
  postgres:
    default_port: 5432
    ssl_mode: "disable"
    connection_timeout: 30
    statement_timeout: 300
    application_name: "db-backup-utility"
    timezone: "UTC"
```

### MongoDB Configuration

```yaml
database:
  mongodb:
    default_port: 27017
    auth_source: "admin"
    connection_timeout: 30
    socket_timeout: 30
    server_selection_timeout: 30
    max_pool_size: 100
    min_pool_size: 0
```

### SQLite Configuration

```yaml
database:
  sqlite:
    journal_mode: "WAL"
    synchronous: "NORMAL"
    cache_size: 1000
    temp_store: "MEMORY"
    foreign_keys: true
    busy_timeout: 30000
```

## Environment Variables

### Database Connection

```bash
# MySQL
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USERNAME=backup_user
export MYSQL_PASSWORD=secure_password
export MYSQL_DATABASE=production_db

# PostgreSQL
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USERNAME=backup_user
export POSTGRES_PASSWORD=secure_password
export POSTGRES_DATABASE=production_db
export POSTGRES_SSLMODE=disable

# MongoDB
export MONGODB_HOST=localhost
export MONGODB_PORT=27017
export MONGODB_USERNAME=backup_user
export MONGODB_PASSWORD=secure_password
export MONGODB_DATABASE=production_db
export MONGODB_AUTH_SOURCE=admin

# SQLite
export SQLITE_DATABASE=/path/to/database.db
```

### Cloud Storage

```bash
# AWS S3
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1
export AWS_S3_BUCKET=my-backup-bucket

# Google Cloud Storage
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
export GCP_PROJECT_ID=my-project
export GCP_BUCKET=my-backup-bucket

# Azure Blob Storage
export AZURE_STORAGE_ACCOUNT=my_storage_account
export AZURE_STORAGE_KEY=my_storage_key
export AZURE_CONTAINER=my-backup-container
```

### Notifications

```bash
# Slack
export SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
export SLACK_CHANNEL=#backups

# Discord
export DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK
export DISCORD_CHANNEL=backups
```

## Command Line Overrides

Configuration values can be overridden using command line flags:

```bash
# Override log level
./db-backup --log-level debug backup --db-type mysql ...

# Override storage path
./db-backup backup --db-type mysql --path /custom/backup/path ...

# Override cloud provider
./db-backup backup --db-type mysql --storage cloud --cloud-provider aws --bucket my-bucket ...
```

## Configuration Validation

The utility validates configuration on startup:

```bash
# Test configuration
./db-backup --config /path/to/config.yaml test --db-type mysql ...
```

## Multiple Configuration Files

You can use different configuration files for different environments:

```bash
# Development configuration
./db-backup --config ~/.db-backup-dev.yaml backup --db-type mysql ...

# Production configuration
./db-backup --config ~/.db-backup-prod.yaml backup --db-type mysql ...

# Staging configuration
./db-backup --config ~/.db-backup-staging.yaml backup --db-type mysql ...
```

## Configuration Examples

### Development Environment

```yaml
# ~/.db-backup-dev.yaml
log:
  level: "debug"
  format: "text"

storage:
  type: "local"
  path: "./backups/dev"

notify:
  enabled: false

database:
  mysql:
    default_port: 3306
    connection_timeout: 10
```

### Production Environment

```yaml
# ~/.db-backup-prod.yaml
log:
  level: "info"
  format: "json"

storage:
  type: "cloud"
  path: "./temp"

cloud:
  provider: "aws"
  bucket: "prod-backup-bucket"
  region: "us-east-1"

notify:
  enabled: true
  type: "slack"
  webhook: "https://hooks.slack.com/services/PROD/SLACK/WEBHOOK"
  channel: "#prod-backups"

database:
  mysql:
    default_port: 3306
    connection_timeout: 30
    ssl_mode: "required"
```

### Staging Environment

```yaml
# ~/.db-backup-staging.yaml
log:
  level: "info"
  format: "json"

storage:
  type: "local"
  path: "/var/backups/staging"

notify:
  enabled: true
  type: "discord"
  webhook: "https://discord.com/api/webhooks/STAGING/DISCORD/WEBHOOK"
  channel: "staging-backups"

database:
  mysql:
    default_port: 3306
    connection_timeout: 20
```

## Security Best Practices

### Credential Management

1. **Never store credentials in configuration files**
2. **Use environment variables for sensitive data**
3. **Use dedicated service accounts with minimal permissions**
4. **Rotate credentials regularly**
5. **Use IAM roles when possible (AWS)**

### File Permissions

```bash
# Secure configuration file
chmod 600 ~/.db-backup.yaml
chown $USER:$USER ~/.db-backup.yaml

# Secure backup directory
chmod 700 /var/backups/database-backup-utility
chown backup_user:backup_group /var/backups/database-backup-utility
```

### Network Security

```yaml
# Use SSL/TLS for database connections
database:
  mysql:
    ssl_mode: "required"

  postgres:
    ssl_mode: "require"
```

## Troubleshooting Configuration

### Common Issues

#### Configuration File Not Found

```bash
# Check if configuration file exists
ls -la ~/.db-backup.yaml

# Create default configuration
cp config.example.yaml ~/.db-backup.yaml
```

#### Invalid YAML Syntax

```bash
# Validate YAML syntax
python -c "import yaml; yaml.safe_load(open('~/.db-backup.yaml'))"

# Or use online YAML validator
```

#### Environment Variables Not Set

```bash
# Check environment variables
env | grep -E "(MYSQL|POSTGRES|MONGODB|AWS|SLACK|DISCORD)"

# Set missing variables
export MYSQL_HOST=localhost
export MYSQL_USERNAME=backup_user
export MYSQL_PASSWORD=secure_password
```

## Related Documentation

- [Getting Started](getting-started.md)
- [MySQL Guide](mysql.md)
- [PostgreSQL Guide](postgresql.md)
- [MongoDB Guide](mongodb.md)
- [SQLite Guide](sqlite.md)
- [Cloud Storage](cloud-storage.md)
- [Notifications](notifications.md)
- [Troubleshooting](troubleshooting.md)
