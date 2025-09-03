package database

import (
	"database/sql"
	"fmt"
	"strings"

	"database-backup-utility/internal/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config represents database connection configuration
type Config struct {
	Type             string
	Host             string
	Port             string
	Username         string
	Password         string
	Database         string
	ConnectionString string
}

// Manager handles database connections and operations
type Manager struct {
	config Config
	db     *sql.DB
	mongo  *mongo.Client
}

// NewManager creates a new database manager
func NewManager(config Config) (*Manager, error) {
	manager := &Manager{
		config: config,
	}

	// Validate configuration
	if err := manager.validateConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return manager, nil
}

// validateConfig validates the database configuration
func (m *Manager) validateConfig() error {
	if m.config.Type == "" {
		return fmt.Errorf("database type is required")
	}

	supportedTypes := []string{"mysql", "postgres", "postgresql", "mongodb", "sqlite"}
	if !contains(supportedTypes, strings.ToLower(m.config.Type)) {
		return fmt.Errorf("unsupported database type: %s", m.config.Type)
	}

	// For SQLite, only database name is required
	if strings.ToLower(m.config.Type) == "sqlite" {
		if m.config.Database == "" {
			return fmt.Errorf("database name is required for SQLite")
		}
		return nil
	}

	// For other databases, validate connection parameters
	if m.config.ConnectionString == "" {
		if m.config.Host == "" {
			return fmt.Errorf("host is required when connection string is not provided")
		}
		if m.config.Username == "" {
			return fmt.Errorf("username is required when connection string is not provided")
		}
		if m.config.Database == "" {
			return fmt.Errorf("database name is required when connection string is not provided")
		}
	}

	return nil
}

// TestConnection tests the database connection
func (m *Manager) TestConnection() error {
	switch strings.ToLower(m.config.Type) {
	case "mysql", "postgres", "postgresql", "sqlite":
		return m.testSQLConnection()
	case "mongodb":
		return m.testMongoConnection()
	default:
		return fmt.Errorf("unsupported database type: %s", m.config.Type)
	}
}

// testSQLConnection tests SQL database connection
func (m *Manager) testSQLConnection() error {
	connectionString := m.getSQLConnectionString()

	db, err := sql.Open(m.getSQLDriver(), connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	m.db = db
	logger.Info("SQL database connection successful", "type", m.config.Type)
	return nil
}

// testMongoConnection tests MongoDB connection
func (m *Manager) testMongoConnection() error {
	connectionString := m.getMongoConnectionString()

	client, err := mongo.Connect(nil, options.Client().ApplyURI(connectionString))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	if err := client.Ping(nil, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.mongo = client
	logger.Info("MongoDB connection successful")
	return nil
}

// getSQLConnectionString returns the SQL connection string
func (m *Manager) getSQLConnectionString() string {
	if m.config.ConnectionString != "" {
		return m.config.ConnectionString
	}

	switch strings.ToLower(m.config.Type) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			m.config.Username, m.config.Password, m.config.Host, m.config.Port, m.config.Database)
	case "postgres", "postgresql":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			m.config.Host, m.config.Port, m.config.Username, m.config.Password, m.config.Database)
	case "sqlite":
		return m.config.Database
	default:
		return ""
	}
}

// getMongoConnectionString returns the MongoDB connection string
func (m *Manager) getMongoConnectionString() string {
	if m.config.ConnectionString != "" {
		return m.config.ConnectionString
	}

	// Default port for MongoDB
	port := m.config.Port
	if port == "" {
		port = "27017"
	}

	if m.config.Username != "" && m.config.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
			m.config.Username, m.config.Password, m.config.Host, port, m.config.Database)
	}

	return fmt.Sprintf("mongodb://%s:%s/%s", m.config.Host, port, m.config.Database)
}

// getSQLDriver returns the SQL driver name
func (m *Manager) getSQLDriver() string {
	switch strings.ToLower(m.config.Type) {
	case "mysql":
		return "mysql"
	case "postgres", "postgresql":
		return "postgres"
	case "sqlite":
		return "sqlite3"
	default:
		return ""
	}
}

// GetDB returns the SQL database connection
func (m *Manager) GetDB() *sql.DB {
	return m.db
}

// GetMongoClient returns the MongoDB client
func (m *Manager) GetMongoClient() *mongo.Client {
	return m.mongo
}

// GetDatabaseName returns the database name
func (m *Manager) GetDatabaseName() string {
	return m.config.Database
}

// GetDatabaseType returns the database type
func (m *Manager) GetDatabaseType() string {
	return m.config.Type
}

// Close closes the database connection
func (m *Manager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	if m.mongo != nil {
		return m.mongo.Disconnect(nil)
	}
	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
