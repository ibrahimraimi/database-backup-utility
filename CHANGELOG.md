# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of database backup utility
- Support for MySQL, PostgreSQL, MongoDB, and SQLite databases
- Full, incremental, and differential backup types
- Local and cloud storage options (AWS S3, Google Cloud Storage, Azure Blob Storage)
- Compression support with gzip
- Slack and Discord notifications
- Comprehensive CLI interface with Cobra
- Configuration management with Viper
- Structured logging with Logrus
- Cross-platform builds (Linux, macOS, Windows)
- GitHub Actions CI/CD pipeline
- Comprehensive documentation

### Features

- **Database Support**: MySQL, PostgreSQL, MongoDB, SQLite
- **Backup Types**: Full, incremental, differential
- **Storage**: Local filesystem, AWS S3, Google Cloud Storage, Azure Blob Storage
- **Compression**: Built-in gzip compression
- **Notifications**: Slack and Discord webhook support
- **CLI**: Full command-line interface with subcommands
- **Configuration**: YAML configuration file support
- **Logging**: Structured logging with multiple levels and formats
- **Cross-Platform**: Native binaries for Linux, macOS, and Windows

### Technical Details

- Built with Go 1.21+
- Uses Cobra for CLI framework
- Uses Viper for configuration management
- Uses Logrus for structured logging
- Uses official database drivers for each supported database
- Docker support with multi-stage builds
- GitHub Actions for automated testing and releases

## [v0.1.0] - 2024-09-03

### Added

- Initial release
- Basic backup and restore functionality
- Database connection testing
- Configuration file support
- Basic logging
- Cross-platform builds

---

## Release Process

### Creating a New Release

1. Update this CHANGELOG.md with the new version
2. Update version in relevant files if needed
3. Create a Git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push the tag: `git push origin v1.0.0`
5. GitHub Actions will automatically create a release with binaries

### Version Format

- **v0.x.x** - Development releases
- **v1.x.x** - Stable releases
- **v1.0.0** - First stable release

### Breaking Changes

Breaking changes will be clearly marked in the changelog and will result in a major version bump.
