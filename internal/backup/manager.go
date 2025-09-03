package backup

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

// Config represents backup configuration
type Config struct {
	DBType           string
	Host             string
	Port             string
	Username         string
	Password         string
	Database         string
	ConnectionString string
	BackupType       string
	Compress         bool
	StorageType      string
	StoragePath      string
	CloudProvider    string
	CloudBucket      string
	CloudRegion      string
	SelectiveTables  string
}

// Manager handles backup operations
type Manager struct {
	config         Config
	dbManager      *database.Manager
	storageManager *storage.Manager
}

// NewManager creates a new backup manager
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

	// Create storage manager
	storageConfig := storage.Config{
		Type:          config.StorageType,
		Path:          config.StoragePath,
		CloudProvider: config.CloudProvider,
		Bucket:        config.CloudBucket,
		Region:        config.CloudRegion,
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

// CreateBackup creates a backup of the database
func (m *Manager) CreateBackup() (string, error) {
	startTime := time.Now()
	logger.Info("Starting backup creation", "type", m.config.BackupType, "database", m.config.Database)

	// Generate backup filename
	filename := m.generateBackupFilename()
	backupPath := filepath.Join(m.config.StoragePath, filename)

	// Create backup based on database type
	var err error
	switch strings.ToLower(m.config.DBType) {
	case "mysql":
		err = m.createMySQLBackup(backupPath)
	case "postgres", "postgresql":
		err = m.createPostgreSQLBackup(backupPath)
	case "mongodb":
		err = m.createMongoDBBackup(backupPath)
	case "sqlite":
		err = m.createSQLiteBackup(backupPath)
	default:
		return "", fmt.Errorf("unsupported database type: %s", m.config.DBType)
	}

	if err != nil {
		logger.Error("Backup creation failed", "error", err)
		return "", err
	}

	// Compress backup if requested
	if m.config.Compress {
		compressedPath, err := m.compressBackup(backupPath)
		if err != nil {
			logger.Error("Backup compression failed", "error", err)
			return "", err
		}
		backupPath = compressedPath
	}

	// Upload to cloud storage if configured
	if m.config.StorageType == "cloud" {
		if err := m.storageManager.Upload(backupPath, filename); err != nil {
			logger.Error("Cloud upload failed", "error", err)
			return "", err
		}
	}

	duration := time.Since(startTime)
	logger.Info("Backup completed successfully",
		"path", backupPath,
		"duration", duration.String(),
		"size", m.getFileSize(backupPath))

	return backupPath, nil
}

// generateBackupFilename generates a unique backup filename
func (m *Manager) generateBackupFilename() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	dbType := strings.ToLower(m.config.DBType)
	database := m.config.Database
	backupType := m.config.BackupType

	// For SQLite, extract just the filename from the path
	if strings.ToLower(m.config.DBType) == "sqlite" {
		database = filepath.Base(database)
		// Remove the .db extension for the filename
		database = strings.TrimSuffix(database, ".db")
	}

	extension := ".sql"
	if strings.ToLower(m.config.DBType) == "mongodb" {
		extension = ".bson"
	}
	if strings.ToLower(m.config.DBType) == "sqlite" {
		extension = ".db"
	}

	filename := fmt.Sprintf("%s_%s_%s_%s%s", dbType, database, backupType, timestamp, extension)

	if m.config.Compress {
		filename += ".gz"
	}

	return filename
}

// createMySQLBackup creates a MySQL backup
func (m *Manager) createMySQLBackup(backupPath string) error {
	logger.Info("Creating MySQL backup")

	// For MySQL, we would typically use mysqldump command
	// This is a simplified implementation
	db := m.dbManager.GetDB()
	if db == nil {
		return fmt.Errorf("database connection not available")
	}

	// Get list of tables
	tables, err := m.getMySQLTables()
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Create backup file
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Write backup header
	file.WriteString("-- MySQL Backup\n")
	file.WriteString(fmt.Sprintf("-- Database: %s\n", m.config.Database))
	file.WriteString(fmt.Sprintf("-- Created: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString("-- \n\n")

	// Backup each table
	for _, table := range tables {
		if err := m.backupMySQLTable(file, table); err != nil {
			logger.Warn("Failed to backup table", "table", table, "error", err)
		}
	}

	return nil
}

// createPostgreSQLBackup creates a PostgreSQL backup
func (m *Manager) createPostgreSQLBackup(backupPath string) error {
	logger.Info("Creating PostgreSQL backup")

	// For PostgreSQL, we would typically use pg_dump command
	// This is a simplified implementation
	db := m.dbManager.GetDB()
	if db == nil {
		return fmt.Errorf("database connection not available")
	}

	// Get list of tables
	tables, err := m.getPostgreSQLTables()
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Create backup file
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Write backup header
	file.WriteString("-- PostgreSQL Backup\n")
	file.WriteString(fmt.Sprintf("-- Database: %s\n", m.config.Database))
	file.WriteString(fmt.Sprintf("-- Created: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString("-- \n\n")

	// Backup each table
	for _, table := range tables {
		if err := m.backupPostgreSQLTable(file, table); err != nil {
			logger.Warn("Failed to backup table", "table", table, "error", err)
		}
	}

	return nil
}

// createMongoDBBackup creates a MongoDB backup
func (m *Manager) createMongoDBBackup(backupPath string) error {
	logger.Info("Creating MongoDB backup")

	// For MongoDB, we would typically use mongodump command
	// This is a simplified implementation
	client := m.dbManager.GetMongoClient()
	if client == nil {
		return fmt.Errorf("MongoDB connection not available")
	}

	// Get list of collections
	collections, err := m.getMongoDBCollections()
	if err != nil {
		return fmt.Errorf("failed to get collections: %w", err)
	}

	// Create backup file
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Write backup header
	file.WriteString("-- MongoDB Backup\n")
	file.WriteString(fmt.Sprintf("-- Database: %s\n", m.config.Database))
	file.WriteString(fmt.Sprintf("-- Created: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString("-- \n\n")

	// Backup each collection
	for _, collection := range collections {
		if err := m.backupMongoDBCollection(file, collection); err != nil {
			logger.Warn("Failed to backup collection", "collection", collection, "error", err)
		}
	}

	return nil
}

// createSQLiteBackup creates a SQLite backup
func (m *Manager) createSQLiteBackup(backupPath string) error {
	logger.Info("Creating SQLite backup")

	// For SQLite, we can simply copy the database file
	sourceFile := m.config.Database
	if m.config.ConnectionString != "" {
		sourceFile = m.config.ConnectionString
	}

	// Copy the database file
	source, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to copy database file: %w", err)
	}

	return nil
}

// compressBackup compresses the backup file
func (m *Manager) compressBackup(backupPath string) (string, error) {
	logger.Info("Compressing backup file", "path", backupPath)

	compressedPath := backupPath + ".gz"

	// Open the original file
	file, err := os.Open(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	// Create the compressed file
	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create compressed file: %w", err)
	}
	defer compressedFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(compressedFile)
	defer gzWriter.Close()

	// Copy the file content to the gzip writer
	_, err = io.Copy(gzWriter, file)
	if err != nil {
		return "", fmt.Errorf("failed to compress file: %w", err)
	}

	// Remove the original file
	if err := os.Remove(backupPath); err != nil {
		logger.Warn("Failed to remove original backup file", "error", err)
	}

	logger.Info("Backup compressed successfully", "compressed_path", compressedPath)
	return compressedPath, nil
}

// getFileSize returns the size of a file in bytes
func (m *Manager) getFileSize(filepath string) int64 {
	file, err := os.Stat(filepath)
	if err != nil {
		return 0
	}
	return file.Size()
}

// Helper methods for different database types
func (m *Manager) getMySQLTables() ([]string, error) {
	// Implementation would query information_schema.tables
	return []string{"users", "orders", "products"}, nil
}

func (m *Manager) getPostgreSQLTables() ([]string, error) {
	// Implementation would query information_schema.tables
	return []string{"users", "orders", "products"}, nil
}

func (m *Manager) getMongoDBCollections() ([]string, error) {
	// Implementation would list collections from MongoDB
	return []string{"users", "orders", "products"}, nil
}

func (m *Manager) backupMySQLTable(file *os.File, table string) error {
	// Implementation would dump table data
	file.WriteString(fmt.Sprintf("-- Table: %s\n", table))
	file.WriteString(fmt.Sprintf("SELECT * FROM %s;\n\n", table))
	return nil
}

func (m *Manager) backupPostgreSQLTable(file *os.File, table string) error {
	// Implementation would dump table data
	file.WriteString(fmt.Sprintf("-- Table: %s\n", table))
	file.WriteString(fmt.Sprintf("SELECT * FROM %s;\n\n", table))
	return nil
}

func (m *Manager) backupMongoDBCollection(file *os.File, collection string) error {
	// Implementation would dump collection data
	file.WriteString(fmt.Sprintf("-- Collection: %s\n", collection))
	file.WriteString(fmt.Sprintf("db.%s.find();\n\n", collection))
	return nil
}
