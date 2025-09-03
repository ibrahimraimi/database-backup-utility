# Database Backup Utility Documentation

This directory contains comprehensive documentation for using the Database Backup Utility with different database systems.

## Documentation Index

- [Getting Started](getting-started.md) - Quick start guide and installation
- [MySQL Guide](mysql.md) - Complete MySQL backup and restore guide
- [PostgreSQL Guide](postgresql.md) - Complete PostgreSQL backup and restore guide
- [MongoDB Guide](mongodb.md) - Complete MongoDB backup and restore guide
- [SQLite Guide](sqlite.md) - Complete SQLite backup and restore guide
- [Configuration](configuration.md) - Configuration file setup and options
- [Cloud Storage](cloud-storage.md) - Cloud storage integration guide
- [Notifications](notifications.md) - Slack and Discord notification setup
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
- [Examples](examples.md) - Real-world usage examples and scripts

## Quick Reference

### Test Database Connection

```bash
./dbu test --db-type <database_type> --host <host> --username <user> --password <pass> --database <db_name>
```

### Create Backup

```bash
./dbu backup --db-type <database_type> --host <host> --username <user> --password <pass> --database <db_name> --compress
```

### Restore Backup

```bash
./dbu restore --db-type <database_type> --host <host> --username <user> --password <pass> --database <db_name> --file <backup_file>
```

## Supported Database Types

| Database   | Type                       | Default Port | Notes                       |
| ---------- | -------------------------- | ------------ | --------------------------- |
| MySQL      | `mysql`                    | 3306         | Full support with mysqldump |
| PostgreSQL | `postgres` or `postgresql` | 5432         | Full support with pg_dump   |
| MongoDB    | `mongodb`                  | 27017        | Full support with mongodump |
| SQLite     | `sqlite`                   | N/A          | File-based database         |

## Common Flags

| Flag                  | Description                                      | Required                        |
| --------------------- | ------------------------------------------------ | ------------------------------- |
| `--db-type`           | Database type (mysql, postgres, mongodb, sqlite) | Yes                             |
| `--host`              | Database host                                    | Yes (except SQLite)             |
| `--port`              | Database port                                    | No (uses defaults)              |
| `--username`          | Database username                                | Yes (except SQLite)             |
| `--password`          | Database password                                | Yes (except SQLite)             |
| `--database`          | Database name                                    | Yes                             |
| `--connection-string` | Full connection string                           | Alternative to individual flags |

## Next Steps

1. Start with [Getting Started](getting-started.md) for installation and basic usage
2. Choose your database type and follow the specific guide
3. Set up [Configuration](configuration.md) for advanced options
4. Explore [Examples](examples.md) for real-world scenarios
