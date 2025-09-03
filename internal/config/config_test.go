package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.LogLevel)
	}

	if cfg.LogFormat != "json" {
		t.Errorf("Expected log format 'json', got '%s'", cfg.LogFormat)
	}

	if cfg.Storage.Type != "local" {
		t.Errorf("Expected storage type 'local', got '%s'", cfg.Storage.Type)
	}
}
