package cmd

import (
	"database-backup-utility/internal/backup"
	"database-backup-utility/internal/logger"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	backupType      string
	compress        bool
	storageType     string
	storagePath     string
	cloudProvider   string
	cloudBucket     string
	cloudRegion     string
	selectiveTables string
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of the specified database",
	Long: `Create a backup of the specified database. Supports multiple database types
including MySQL, PostgreSQL, MongoDB, and SQLite. The backup can be stored locally
or uploaded to cloud storage providers like AWS S3, Google Cloud Storage, or Azure Blob Storage.`,
	RunE: runBackup,
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Database connection flags
	backupCmd.Flags().String("db-type", "", "Database type (mysql, postgres, mongodb, sqlite)")
	backupCmd.Flags().String("host", "", "Database host")
	backupCmd.Flags().String("port", "", "Database port")
	backupCmd.Flags().String("username", "", "Database username")
	backupCmd.Flags().String("password", "", "Database password")
	backupCmd.Flags().String("database", "", "Database name")
	backupCmd.Flags().String("connection-string", "", "Full database connection string")

	// Backup configuration flags
	backupCmd.Flags().StringVar(&backupType, "type", "full", "Backup type (full, incremental, differential)")
	backupCmd.Flags().BoolVar(&compress, "compress", true, "Compress backup files")
	backupCmd.Flags().StringVar(&selectiveTables, "tables", "", "Comma-separated list of tables to backup (selective backup)")

	// Storage configuration flags
	backupCmd.Flags().StringVar(&storageType, "storage", "local", "Storage type (local, cloud)")
	backupCmd.Flags().StringVar(&storagePath, "path", "./backups", "Local storage path")
	backupCmd.Flags().StringVar(&cloudProvider, "cloud-provider", "", "Cloud provider (aws, gcp, azure)")
	backupCmd.Flags().StringVar(&cloudBucket, "bucket", "", "Cloud storage bucket name")
	backupCmd.Flags().StringVar(&cloudRegion, "region", "", "Cloud storage region")

	// Mark required flags
	_ = backupCmd.MarkFlagRequired("db-type")
	_ = backupCmd.MarkFlagRequired("database")
}

func runBackup(cmd *cobra.Command, args []string) error {
	logger.Info("Starting backup operation", "timestamp", time.Now())

	// Get database connection parameters
	dbType, _ := cmd.Flags().GetString("db-type")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	database, _ := cmd.Flags().GetString("database")
	connectionString, _ := cmd.Flags().GetString("connection-string")

	// Create backup configuration
	config := backup.Config{
		DBType:           dbType,
		Host:             host,
		Port:             port,
		Username:         username,
		Password:         password,
		Database:         database,
		ConnectionString: connectionString,
		BackupType:       backupType,
		Compress:         compress,
		StorageType:      storageType,
		StoragePath:      storagePath,
		CloudProvider:    cloudProvider,
		CloudBucket:      cloudBucket,
		CloudRegion:      cloudRegion,
		SelectiveTables:  selectiveTables,
	}

	// Create backup manager
	manager, err := backup.NewManager(config)
	if err != nil {
		return fmt.Errorf("failed to create backup manager: %w", err)
	}

	// Test database connection
	logger.Info("Testing database connection")
	if err := manager.TestConnection(); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}
	logger.Info("Database connection successful")

	// Perform backup
	logger.Info("Starting backup process", "type", backupType, "database", database)
	backupPath, err := manager.CreateBackup()
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	logger.Info("Backup completed successfully", "path", backupPath)
	fmt.Printf("Backup completed successfully. File saved to: %s\n", backupPath)

	return nil
}
