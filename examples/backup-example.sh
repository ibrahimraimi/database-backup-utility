#!/bin/bash

# Database Backup Utility - Example Usage Script
# This script demonstrates how to use the database backup utility

set -e

# Configuration
DB_HOST="localhost"
DB_USER="root"
DB_PASS="password"
DB_NAME="example_db"
BACKUP_DIR="./backups"

echo "🚀 Database Backup Utility Example"
echo "=================================="

# Create backup directory
mkdir -p $BACKUP_DIR

echo "📋 Testing database connection..."
./db-backup test \
  --db-type mysql \
  --host $DB_HOST \
  --username $DB_USER \
  --password $DB_PASS \
  --database $DB_NAME

echo "💾 Creating full backup..."
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

echo "📊 Listing backup files..."
ls -la $BACKUP_DIR/

echo "✅ Example completed successfully!"
echo ""
echo "To restore a backup, use:"
echo "./db-backup restore --db-type mysql --host $DB_HOST --username $DB_USER --password $DB_PASS --database $DB_NAME --file <backup-file>"
