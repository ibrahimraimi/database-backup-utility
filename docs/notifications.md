# Notifications Guide

This guide covers setting up and configuring notifications for the Database Backup Utility.

## Supported Notification Providers

- **Slack** - Team communication platform
- **Discord** - Gaming and community platform

## Slack Integration

### Prerequisites

- Slack workspace with admin access
- Webhook URL for the target channel

### Creating a Slack Webhook

1. Go to your Slack workspace
2. Navigate to Apps → Incoming Webhooks
3. Click "Add to Slack"
4. Choose the channel for notifications
5. Copy the webhook URL

### Configuration

```yaml
# ~/dbu.yaml
notify:
  enabled: true
  type: "slack"
  webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  channel: "#backups"
```

### Environment Variables

```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
export SLACK_CHANNEL="#backups"
```

### Usage Examples

#### Basic Backup with Slack Notification

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --compress
```

#### Slack Notification Message Format

**Success Notification:**

```json
{
  "channel": "#backups",
  "attachments": [
    {
      "color": "good",
      "text": "✅ Database backup completed successfully!\nDatabase: mydb\nBackup file: mysql_mydb_full_2024-01-15_10-30-00.sql.gz\nDuration: 2m 30s\nSize: 1.2 GB",
      "timestamp": 1642248600
    }
  ]
}
```

**Failure Notification:**

```json
{
  "channel": "#backups",
  "attachments": [
    {
      "color": "danger",
      "text": "❌ Database backup failed!\nDatabase: mydb\nError: Connection timeout",
      "timestamp": 1642248600
    }
  ]
}
```

### Advanced Slack Configuration

#### Custom Slack App (Recommended for Production)

1. Create a Slack App at https://api.slack.com/apps
2. Enable Incoming Webhooks
3. Create a webhook for your channel
4. Use the webhook URL in your configuration

#### Slack Bot Token (Alternative)

```yaml
notify:
  enabled: true
  type: "slack"
  webhook: "" # Not needed for bot token
  channel: "#backups"
  bot_token: "xoxb-your-bot-token" # Set via SLACK_BOT_TOKEN
```

```bash
export SLACK_BOT_TOKEN="xoxb-your-bot-token"
```

## Discord Integration

### Prerequisites

- Discord server with admin access
- Webhook URL for the target channel

### Creating a Discord Webhook

1. Go to your Discord server
2. Right-click on the target channel
3. Select "Edit Channel"
4. Go to "Integrations" → "Webhooks"
5. Click "Create Webhook"
6. Copy the webhook URL

### Configuration

```yaml
# ~/dbu.yaml
notify:
  enabled: true
  type: "discord"
  webhook: "https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"
  channel: "backups"
```

### Environment Variables

```bash
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"
export DISCORD_CHANNEL="backups"
```

### Usage Examples

#### Basic Backup with Discord Notification

```bash
./dbu backup \
  --db-type mysql \
  --host localhost \
  --username root \
  --password mypassword \
  --database mydb \
  --compress
```

#### Discord Notification Message Format

**Success Notification:**

```json
{
  "embeds": [
    {
      "title": "Database Backup Utility",
      "description": "✅ Database backup completed successfully!\nDatabase: mydb\nBackup file: mysql_mydb_full_2024-01-15_10-30-00.sql.gz\nDuration: 2m 30s\nSize: 1.2 GB",
      "color": 65280,
      "timestamp": "2024-01-15T10:30:00.000Z"
    }
  ]
}
```

**Failure Notification:**

```json
{
  "embeds": [
    {
      "title": "Database Backup Utility",
      "description": "❌ Database backup failed!\nDatabase: mydb\nError: Connection timeout",
      "color": 16711680,
      "timestamp": "2024-01-15T10:30:00.000Z"
    }
  ]
}
```

## Notification Types

### Backup Notifications

#### Success Notifications

- Backup completion time
- Database name
- Backup file path
- Duration
- File size
- Compression ratio

#### Failure Notifications

- Error message
- Database name
- Timestamp
- Suggested actions

### Restore Notifications

#### Success Notifications

- Restore completion time
- Database name
- Backup file used
- Duration

#### Failure Notifications

- Error message
- Database name
- Backup file path
- Suggested actions

## Advanced Configuration

### Multiple Notification Channels

```yaml
# ~/dbu.yaml
notify:
  enabled: true
  type: "slack"
  webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  channel: "#backups"

  # Additional notification channels
  additional_channels:
    - type: "discord"
      webhook: "https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"
      channel: "backups"
    - type: "slack"
      webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK2"
      channel: "#alerts"
```

### Conditional Notifications

```yaml
# ~/dbu.yaml
notify:
  enabled: true
  type: "slack"
  webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  channel: "#backups"

  # Only notify on failures
  notify_on_success: false
  notify_on_failure: true

  # Only notify for specific databases
  databases:
    - "production_db"
    - "staging_db"
```

### Notification Templates

```yaml
# ~/dbu.yaml
notify:
  enabled: true
  type: "slack"
  webhook: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  channel: "#backups"

  # Custom message templates
  templates:
    success: "🎉 Backup completed for {{.Database}} in {{.Duration}}"
    failure: "🚨 Backup failed for {{.Database}}: {{.Error}}"
```

## Testing Notifications

### Test Slack Notification

```bash
# Test Slack webhook
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Test notification from Database Backup Utility"}' \
  $SLACK_WEBHOOK_URL
```

### Test Discord Notification

```bash
# Test Discord webhook
curl -X POST -H 'Content-type: application/json' \
  --data '{"content":"Test notification from Database Backup Utility"}' \
  $DISCORD_WEBHOOK_URL
```

### Test with Backup Command

```bash
# Test with a small backup
./dbu backup \
  --db-type sqlite \
  --database ./test.db \
  --compress
```

## Notification Scripts

### Custom Notification Script

```bash
#!/bin/bash
# custom-notification.sh

# Function to send custom notification
send_notification() {
    local status=$1
    local database=$2
    local message=$3

    if [ "$status" = "success" ]; then
        color="good"
        emoji="✅"
    else
        color="danger"
        emoji="❌"
    fi

    # Send to Slack
    curl -X POST -H 'Content-type: application/json' \
      --data "{
        \"channel\": \"#backups\",
        \"attachments\": [
          {
            \"color\": \"$color\",
            \"text\": \"$emoji $message\",
            \"timestamp\": $(date +%s)
          }
        ]
      }" \
      $SLACK_WEBHOOK_URL
}

# Usage in backup script
./dbu backup --db-type mysql --database mydb --compress

if [ $? -eq 0 ]; then
    send_notification "success" "mydb" "Backup completed successfully"
else
    send_notification "failure" "mydb" "Backup failed"
fi
```

### Multi-Provider Notification Script

```bash
#!/bin/bash
# multi-notification.sh

send_slack() {
    local message=$1
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"$message\"}" \
      $SLACK_WEBHOOK_URL
}

send_discord() {
    local message=$1
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"content\":\"$message\"}" \
      $DISCORD_WEBHOOK_URL
}

send_notification() {
    local message=$1
    send_slack "$message"
    send_discord "$message"
}

# Usage
send_notification "Database backup completed successfully"
```

## Monitoring and Alerting

### Health Checks

```bash
#!/bin/bash
# health-check.sh

# Check if backup utility is working
if ! ./dbu test --db-type mysql --host localhost --username root --password mypass --database mydb; then
    # Send alert
    curl -X POST -H 'Content-type: application/json' \
      --data '{"text":"🚨 Database Backup Utility health check failed"}' \
      $SLACK_WEBHOOK_URL
fi
```

### Backup Monitoring

```bash
#!/bin/bash
# backup-monitor.sh

# Check if backups are being created
BACKUP_DIR="/var/backups"
LAST_BACKUP=$(find $BACKUP_DIR -name "*.sql.gz" -mtime -1 | wc -l)

if [ $LAST_BACKUP -eq 0 ]; then
    # Send alert
    curl -X POST -H 'Content-type: application/json' \
      --data '{"text":"🚨 No backups created in the last 24 hours"}' \
      $SLACK_WEBHOOK_URL
fi
```

## Security Considerations

### Webhook Security

1. **Keep webhook URLs secret**
2. **Use environment variables**
3. **Rotate webhook URLs regularly**
4. **Monitor webhook usage**
5. **Use HTTPS only**

### Message Security

1. **Don't include sensitive data in notifications**
2. **Use generic error messages**
3. **Log detailed errors separately**
4. **Implement message filtering**

## Troubleshooting

### Common Issues

#### Webhook Not Working

```bash
# Test webhook connectivity
curl -I $SLACK_WEBHOOK_URL
curl -I $DISCORD_WEBHOOK_URL

# Check webhook URL format
echo $SLACK_WEBHOOK_URL
echo $DISCORD_WEBHOOK_URL
```

#### Permission Issues

```bash
# Check Slack channel permissions
# Ensure the webhook has permission to post to the channel

# Check Discord channel permissions
# Ensure the webhook has permission to post to the channel
```

#### Rate Limiting

```bash
# Slack rate limits: 1 message per second per webhook
# Discord rate limits: 30 requests per minute per webhook

# Implement rate limiting in your scripts
sleep 1  # Wait 1 second between notifications
```

### Debug Mode

```bash
# Enable debug logging to see notification details
./dbu --log-level debug backup --db-type mysql --database mydb
```

## Best Practices

1. **Use dedicated channels for backups**
2. **Implement notification filtering**
3. **Set up proper error handling**
4. **Monitor notification delivery**
5. **Use appropriate message formatting**
6. **Implement rate limiting**
7. **Keep webhook URLs secure**
8. **Test notifications regularly**
9. **Use different channels for different environments**
10. **Implement notification acknowledgments**

## Related Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [MySQL Guide](mysql.md)
- [PostgreSQL Guide](postgresql.md)
- [MongoDB Guide](mongodb.md)
- [SQLite Guide](sqlite.md)
- [Troubleshooting](troubleshooting.md)
