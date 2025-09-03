package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	LogLevel  string `mapstructure:"log_level"`
	LogFormat string `mapstructure:"log_format"`
	Database  DatabaseConfig
	Storage   StorageConfig
	Cloud     CloudConfig
	Notify    NotificationConfig
}

// DatabaseConfig represents database connection configuration
type DatabaseConfig struct {
	Type             string `mapstructure:"type"`
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
	Database         string `mapstructure:"database"`
	ConnectionString string `mapstructure:"connection_string"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

// CloudConfig represents cloud storage configuration
type CloudConfig struct {
	Provider  string `mapstructure:"provider"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// NotificationConfig represents notification configuration
type NotificationConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Type    string `mapstructure:"type"` // slack, discord
	Webhook string `mapstructure:"webhook"`
	Channel string `mapstructure:"channel"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_format", "json")
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.path", "./backups")

	// Set config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, "dbu.yaml")
	viper.SetConfigFile(configPath)

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults
	}

	// Bind environment variables
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		LogLevel:  "info",
		LogFormat: "json",
		Storage: StorageConfig{
			Type: "local",
			Path: "./backups",
		},
	}
}

// Save saves configuration to file
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, "dbu.yaml")

	// Set values in viper
	viper.Set("log_level", c.LogLevel)
	viper.Set("log_format", c.LogFormat)
	viper.Set("storage.type", c.Storage.Type)
	viper.Set("storage.path", c.Storage.Path)
	viper.Set("cloud.provider", c.Cloud.Provider)
	viper.Set("cloud.bucket", c.Cloud.Bucket)
	viper.Set("cloud.region", c.Cloud.Region)
	viper.Set("notify.enabled", c.Notify.Enabled)
	viper.Set("notify.type", c.Notify.Type)
	viper.Set("notify.webhook", c.Notify.Webhook)
	viper.Set("notify.channel", c.Notify.Channel)

	return viper.WriteConfigAs(configPath)
}
