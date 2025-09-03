package cmd

import (
	"database-backup-utility/internal/database"
	"database-backup-utility/internal/logger"
	"fmt"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test database connection",
	Long: `Test the database connection with the provided parameters.
This command validates credentials and connectivity before proceeding 
with backup or restore operations.`,
	RunE: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Database connection flags
	testCmd.Flags().String("db-type", "", "Database type (mysql, postgres, mongodb, sqlite)")
	testCmd.Flags().String("host", "", "Database host")
	testCmd.Flags().String("port", "", "Database port")
	testCmd.Flags().String("username", "", "Database username")
	testCmd.Flags().String("password", "", "Database password")
	testCmd.Flags().String("database", "", "Database name")
	testCmd.Flags().String("connection-string", "", "Full database connection string")

	// Mark required flags
	_ = testCmd.MarkFlagRequired("db-type")
}

func runTest(cmd *cobra.Command, args []string) error {
	// Get database connection parameters
	dbType, _ := cmd.Flags().GetString("db-type")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	dbName, _ := cmd.Flags().GetString("database")
	connectionString, _ := cmd.Flags().GetString("connection-string")

	// Create database connection configuration
	config := database.Config{
		Type:             dbType,
		Host:             host,
		Port:             port,
		Username:         username,
		Password:         password,
		Database:         dbName,
		ConnectionString: connectionString,
	}

	// Create database manager
	manager, err := database.NewManager(config)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}

	// Test connection
	logger.Info("Testing database connection", "type", dbType, "host", host, "database", dbName)
	if err := manager.TestConnection(); err != nil {
		logger.Error("Database connection test failed", "error", err)
		return fmt.Errorf("database connection test failed: %w", err)
	}

	logger.Info("Database connection test successful")
	fmt.Println("✅ Database connection test successful!")

	return nil
}
