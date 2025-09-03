package database

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid MySQL config",
			config: Config{
				Type:     "mysql",
				Host:     "localhost",
				Port:     "3306",
				Username: "root",
				Password: "password",
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "valid SQLite config",
			config: Config{
				Type:     "sqlite",
				Database: "/path/to/database.db",
			},
			wantErr: false,
		},
		{
			name: "invalid database type",
			config: Config{
				Type:     "invalid",
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing database type",
			config: Config{
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing database name for non-SQLite",
			config: Config{
				Type:     "mysql",
				Host:     "localhost",
				Username: "root",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{config: tt.config}
			err := manager.validateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSQLDriver(t *testing.T) {
	tests := []struct {
		name     string
		dbType   string
		expected string
	}{
		{"MySQL", "mysql", "mysql"},
		{"PostgreSQL", "postgres", "postgres"},
		{"PostgreSQL alt", "postgresql", "postgres"},
		{"SQLite", "sqlite", "sqlite3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{config: Config{Type: tt.dbType}}
			result := manager.getSQLDriver()
			if result != tt.expected {
				t.Errorf("getSQLDriver() = %v, want %v", result, tt.expected)
			}
		})
	}
}
