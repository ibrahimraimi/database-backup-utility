package cmd

import (
	"database-backup-utility/internal/logger"
	"database-backup-utility/internal/restore"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	restoreBackupFile string
	restoreTables     string
	dropExisting      bool
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a database from a backup file",
	Long: `Restore a database from a backup file. Supports selective restoration of 
specific tables or collections if supported by the DBMS. The restore operation 
can handle compressed backup files and restore from both local and cloud storage.`,
	RunE: runRestore,
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Database connection flags
	restoreCmd.Flags().String("db-type", "", "Database type (mysql, postgres, mongodb, sqlite)")
	restoreCmd.Flags().String("host", "", "Database host")
	restoreCmd.Flags().String("port", "", "Database port")
	restoreCmd.Flags().String("username", "", "Database username")
	restoreCmd.Flags().String("password", "", "Database password")
	restoreCmd.Flags().String("database", "", "Database name")
	restoreCmd.Flags().String("connection-string", "", "Full database connection string")

	// Restore configuration flags
	restoreCmd.Flags().StringVar(&restoreBackupFile, "file", "", "Path to backup file to restore")
	restoreCmd.Flags().StringVar(&restoreTables, "tables", "", "Comma-separated list of tables to restore (selective restore)")
	restoreCmd.Flags().BoolVar(&dropExisting, "drop-existing", false, "Drop existing tables before restore")

	// Mark required flags
	restoreCmd.MarkFlagRequired("db-type")
	restoreCmd.MarkFlagRequired("database")
	restoreCmd.MarkFlagRequired("file")
}

func runRestore(cmd *cobra.Command, args []string) error {
	logger.Info("Starting restore operation", "timestamp", time.Now())

	// Get database connection parameters
	dbType, _ := cmd.Flags().GetString("db-type")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	database, _ := cmd.Flags().GetString("database")
	connectionString, _ := cmd.Flags().GetString("connection-string")

	// Create restore configuration
	config := restore.Config{
		DBType:           dbType,
		Host:             host,
		Port:             port,
		Username:         username,
		Password:         password,
		Database:         database,
		ConnectionString: connectionString,
		BackupFile:       restoreBackupFile,
		SelectiveTables:  restoreTables,
		DropExisting:     dropExisting,
	}

	// Create restore manager
	manager, err := restore.NewManager(config)
	if err != nil {
		return fmt.Errorf("failed to create restore manager: %w", err)
	}

	// Test database connection
	logger.Info("Testing database connection")
	if err := manager.TestConnection(); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}
	logger.Info("Database connection successful")

	// Perform restore
	logger.Info("Starting restore process", "file", restoreBackupFile, "database", database)
	if err := manager.RestoreBackup(); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	logger.Info("Restore completed successfully")
	fmt.Printf("Database restored successfully from: %s\n", restoreBackupFile)

	return nil
}
