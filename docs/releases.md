# Releases

This document explains how to download and use pre-built releases of the database backup utility.

## Downloading Releases

### Latest Release

Download the latest release for your platform:

**Linux (AMD64):**

```bash
curl -L -o dbu https://github.com/your-username/database-backup-utility/releases/latest/download/dbu-linux-amd64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L -o dbu https://github.com/your-username/database-backup-utility/releases/latest/download/dbu-darwin-amd64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

**macOS (Apple Silicon):**

```bash
curl -L -o dbu https://github.com/your-username/database-backup-utility/releases/latest/download/dbu-darwin-arm64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

**Windows (AMD64):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/your-username/database-backup-utility/releases/latest/download/dbu-windows-amd64.exe" -OutFile "dbu.exe"
```

### Specific Version

To download a specific version, replace `latest` with the version tag:

```bash
# Example: Download v1.0.0 for Linux
curl -L -o dbu https://github.com/your-username/database-backup-utility/releases/download/v1.0.0/dbu-linux-amd64
chmod +x dbu
sudo mv dbu /usr/local/bin/
```

## Verification

After downloading, verify the installation:

```bash
dbu --version
dbu --help
```

## Release Process

### For Maintainers

To create a new release:

1. **Prepare the release:**

   ```bash
   # Make sure you're on the main branch
   git checkout main
   git pull origin main

   # Run tests
   make test
   ```

2. **Create the release:**

   ```bash
   # Use the release script
   ./scripts/release.sh v1.0.0
   ```

3. **Monitor the build:**
   - The script will create a Git tag
   - GitHub Actions will automatically build and create a release
   - Check the Actions tab for build progress

### Manual Release Process

If you prefer to create releases manually:

1. **Create and push a tag:**

   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically:**

   - Run tests
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload the binaries

3. **Edit the release:**
   - Go to the GitHub releases page
   - Edit the release notes
   - Mark as latest release if needed

## Release Assets

Each release includes the following assets:

- `dbu-linux-amd64` - Linux 64-bit
- `dbu-darwin-amd64` - macOS Intel
- `dbu-darwin-arm64` - macOS Apple Silicon
- `dbu-windows-amd64.exe` - Windows 64-bit

## Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Changelog

See the [CHANGELOG.md](../CHANGELOG.md) file for a detailed list of changes in each release.

## Support

If you encounter issues with a release:

1. Check the [troubleshooting guide](troubleshooting.md)
2. Open an issue on GitHub
3. Check if a newer version is available
