package restore

import (
	"compress/gzip"
	"database-backup-utility/internal/database"
	"database-backup-utility/internal/logger"
	"database-backup-utility/internal/storage"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config represents restore configuration
type Config struct {
	DBType           string
	Host             string
	Port             string
	Username         string
	Password         string
	Database         string
	ConnectionString string
	BackupFile       string
	SelectiveTables  string
	DropExisting     bool
}

// Manager handles restore operations
type Manager struct {
	config         Config
	dbManager      *database.Manager
	storageManager *storage.Manager
}

// NewManager creates a new restore manager
func NewManager(config Config) (*Manager, error) {
	// Create database manager
	dbConfig := database.Config{
		Type:             config.DBType,
		Host:             config.Host,
		Port:             config.Port,
		Username:         config.Username,
		Password:         config.Password,
		Database:         config.Database,
		ConnectionString: config.ConnectionString,
	}

	dbManager, err := database.NewManager(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	// Create storage manager (for cloud storage downloads)
	storageConfig := storage.Config{
		Type: "local", // We'll download to local first
		Path: "./temp",
	}

	storageManager, err := storage.NewManager(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage manager: %w", err)
	}

	return &Manager{
		config:         config,
		dbManager:      dbManager,
		storageManager: storageManager,
	}, nil
}

// TestConnection tests the database connection
func (m *Manager) TestConnection() error {
	return m.dbManager.TestConnection()
}

// RestoreBackup restores a database from a backup file
func (m *Manager) RestoreBackup() error {
	startTime := time.Now()
	logger.Info("Starting restore operation", "file", m.config.BackupFile, "database", m.config.Database)

	// Determine if backup file is local or remote
	backupPath := m.config.BackupFile
	if m.isRemoteFile(m.config.BackupFile) {
		// Download from cloud storage
		tempPath, err := m.downloadBackupFile()
		if err != nil {
			return fmt.Errorf("failed to download backup file: %w", err)
		}
		backupPath = tempPath
		defer os.Remove(tempPath) // Clean up temp file
	}

	// Check if backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	// Determine if backup is compressed
	var restorePath string
	if m.isCompressed(backupPath) {
		// Decompress backup
		decompressedPath, err := m.decompressBackup(backupPath)
		if err != nil {
			return fmt.Errorf("failed to decompress backup: %w", err)
		}
		restorePath = decompressedPath
		defer os.Remove(decompressedPath) // Clean up temp file
	} else {
		restorePath = backupPath
	}

	// Restore based on database type
	var err error
	switch strings.ToLower(m.config.DBType) {
	case "mysql":
		err = m.restoreMySQL(restorePath)
	case "postgres", "postgresql":
		err = m.restorePostgreSQL(restorePath)
	case "mongodb":
		err = m.restoreMongoDB(restorePath)
	case "sqlite":
		err = m.restoreSQLite(restorePath)
	default:
		return fmt.Errorf("unsupported database type: %s", m.config.DBType)
	}

	if err != nil {
		logger.Error("Restore operation failed", "error", err)
		return err
	}

	duration := time.Since(startTime)
	logger.Info("Restore completed successfully", "duration", duration.String())

	return nil
}

// isRemoteFile checks if the backup file is a remote URL
func (m *Manager) isRemoteFile(filepath string) bool {
	return strings.HasPrefix(filepath, "s3://") ||
		strings.HasPrefix(filepath, "gs://") ||
		strings.HasPrefix(filepath, "azure://") ||
		strings.HasPrefix(filepath, "https://")
}

// isCompressed checks if the backup file is compressed
func (m *Manager) isCompressed(filepath string) bool {
	return strings.HasSuffix(filepath, ".gz") ||
		strings.HasSuffix(filepath, ".tar.gz") ||
		strings.HasSuffix(filepath, ".zip")
}

// downloadBackupFile downloads a backup file from cloud storage
func (m *Manager) downloadBackupFile() (string, error) {
	logger.Info("Downloading backup file from cloud storage")

	// Extract provider and path from URL
	provider, remotePath, err := m.parseCloudURL(m.config.BackupFile)
	if err != nil {
		return "", err
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", "backup_*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Configure storage manager for the specific cloud provider
	storageConfig := storage.Config{
		Type:          "cloud",
		CloudProvider: provider,
		Bucket:        m.extractBucketFromPath(remotePath),
		Region:        "us-east-1", // Default region
	}

	storageManager, err := storage.NewManager(storageConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create cloud storage manager: %w", err)
	}

	// Download file
	if err := storageManager.Download(remotePath, tempPath); err != nil {
		return "", fmt.Errorf("failed to download backup file: %w", err)
	}

	return tempPath, nil
}

// parseCloudURL parses a cloud storage URL
func (m *Manager) parseCloudURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "s3://") {
		path := strings.TrimPrefix(url, "s3://")
		return "aws", path, nil
	} else if strings.HasPrefix(url, "gs://") {
		path := strings.TrimPrefix(url, "gs://")
		return "gcp", path, nil
	} else if strings.HasPrefix(url, "azure://") {
		path := strings.TrimPrefix(url, "azure://")
		return "azure", path, nil
	}
	return "", "", fmt.Errorf("unsupported cloud URL format: %s", url)
}

// extractBucketFromPath extracts bucket name from cloud path
func (m *Manager) extractBucketFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// decompressBackup decompresses a backup file
func (m *Manager) decompressBackup(backupPath string) (string, error) {
	logger.Info("Decompressing backup file", "path", backupPath)

	// Create temp file for decompressed content
	tempFile, err := os.CreateTemp("", "decompressed_*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Open compressed file
	compressedFile, err := os.Open(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to open compressed file: %w", err)
	}
	defer compressedFile.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(compressedFile)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create decompressed file
	decompressedFile, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to create decompressed file: %w", err)
	}
	defer decompressedFile.Close()

	// Copy decompressed content
	_, err = io.Copy(decompressedFile, gzReader)
	if err != nil {
		return "", fmt.Errorf("failed to decompress file: %w", err)
	}

	logger.Info("Backup file decompressed successfully", "path", tempPath)
	return tempPath, nil
}

// restoreMySQL restores a MySQL database
func (m *Manager) restoreMySQL(backupPath string) error {
	logger.Info("Restoring MySQL database")

	db := m.dbManager.GetDB()
	if db == nil {
		return fmt.Errorf("database connection not available")
	}

	// Read backup file
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Split SQL statements
	statements := strings.Split(string(backupContent), ";")

	// Execute each statement
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		if _, err := db.Exec(statement); err != nil {
			logger.Warn("Failed to execute SQL statement", "statement", i+1, "error", err)
		}
	}

	logger.Info("MySQL database restored successfully")
	return nil
}

// restorePostgreSQL restores a PostgreSQL database
func (m *Manager) restorePostgreSQL(backupPath string) error {
	logger.Info("Restoring PostgreSQL database")

	db := m.dbManager.GetDB()
	if db == nil {
		return fmt.Errorf("database connection not available")
	}

	// Read backup file
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Split SQL statements
	statements := strings.Split(string(backupContent), ";")

	// Execute each statement
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		if _, err := db.Exec(statement); err != nil {
			logger.Warn("Failed to execute SQL statement", "statement", i+1, "error", err)
		}
	}

	logger.Info("PostgreSQL database restored successfully")
	return nil
}

// restoreMongoDB restores a MongoDB database
func (m *Manager) restoreMongoDB(backupPath string) error {
	logger.Info("Restoring MongoDB database")

	client := m.dbManager.GetMongoClient()
	if client == nil {
		return fmt.Errorf("MongoDB connection not available")
	}

	// For MongoDB, we would typically use mongorestore command
	// This is a simplified implementation
	logger.Info("MongoDB restore would be implemented using mongorestore command")

	// TODO: Implement actual MongoDB restore logic
	return fmt.Errorf("MongoDB restore not fully implemented yet")
}

// restoreSQLite restores a SQLite database
func (m *Manager) restoreSQLite(backupPath string) error {
	logger.Info("Restoring SQLite database")

	// For SQLite, we can simply copy the backup file to the database location
	targetPath := m.config.Database
	if m.config.ConnectionString != "" {
		targetPath = m.config.ConnectionString
	}

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Copy backup file to target location
	source, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer source.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target database file: %w", err)
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	if err != nil {
		return fmt.Errorf("failed to copy database file: %w", err)
	}

	logger.Info("SQLite database restored successfully")
	return nil
}
