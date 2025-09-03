package main

import (
	"database-backup-utility/cmd"
	"database-backup-utility/internal/config"
	"database-backup-utility/internal/logger"
	"os"
)

func main() {
	// Initialize logger with default values first
	logger.Init("info", "text")

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		// If config loading fails, use defaults
		logger.Warn("Failed to load configuration, using defaults", "error", err)
	} else {
		// Re-initialize logger with config values
		logger.Init(cfg.LogLevel, cfg.LogFormat)
	}

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}
