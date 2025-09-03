package notification

import (
	"bytes"
	"database-backup-utility/internal/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Config represents notification configuration
type Config struct {
	Enabled bool
	Type    string // slack, discord
	Webhook string
	Channel string
}

// Manager handles notification operations
type Manager struct {
	config Config
}

// NewManager creates a new notification manager
func NewManager(config Config) (*Manager, error) {
	return &Manager{
		config: config,
	}, nil
}

// SendBackupSuccess sends a success notification
func (m *Manager) SendBackupSuccess(database, backupPath string, duration time.Duration, size int64) error {
	if !m.config.Enabled {
		return nil
	}

	message := fmt.Sprintf("✅ Database backup completed successfully!\n"+
		"Database: %s\n"+
		"Backup file: %s\n"+
		"Duration: %s\n"+
		"Size: %d bytes", database, backupPath, duration.String(), size)

	return m.sendNotification(message, "success")
}

// SendBackupFailure sends a failure notification
func (m *Manager) SendBackupFailure(database string, errorMsg string) error {
	if !m.config.Enabled {
		return nil
	}

	message := fmt.Sprintf("❌ Database backup failed!\n"+
		"Database: %s\n"+
		"Error: %s", database, errorMsg)

	return m.sendNotification(message, "error")
}

// SendRestoreSuccess sends a restore success notification
func (m *Manager) SendRestoreSuccess(database, backupFile string, duration time.Duration) error {
	if !m.config.Enabled {
		return nil
	}

	message := fmt.Sprintf("✅ Database restore completed successfully!\n"+
		"Database: %s\n"+
		"Backup file: %s\n"+
		"Duration: %s", database, backupFile, duration.String())

	return m.sendNotification(message, "success")
}

// SendRestoreFailure sends a restore failure notification
func (m *Manager) SendRestoreFailure(database string, errorMsg string) error {
	if !m.config.Enabled {
		return nil
	}

	message := fmt.Sprintf("❌ Database restore failed!\n"+
		"Database: %s\n"+
		"Error: %s", database, errorMsg)

	return m.sendNotification(message, "error")
}

// sendNotification sends a notification based on the configured type
func (m *Manager) sendNotification(message, status string) error {
	switch strings.ToLower(m.config.Type) {
	case "slack":
		return m.sendSlackNotification(message, status)
	case "discord":
		return m.sendDiscordNotification(message, status)
	default:
		return fmt.Errorf("unsupported notification type: %s", m.config.Type)
	}
}

// sendSlackNotification sends a Slack notification
func (m *Manager) sendSlackNotification(message, status string) error {
	logger.Info("Sending Slack notification", "status", status)

	// Determine color based on status
	color := "good"
	if status == "error" {
		color = "danger"
	}

	// Create Slack payload
	payload := map[string]interface{}{
		"channel": m.config.Channel,
		"attachments": []map[string]interface{}{
			{
				"color":     color,
				"text":      message,
				"timestamp": time.Now().Unix(),
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	// Send HTTP request
	resp, err := http.Post(m.config.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack notification failed with status: %d", resp.StatusCode)
	}

	logger.Info("Slack notification sent successfully")
	return nil
}

// sendDiscordNotification sends a Discord notification
func (m *Manager) sendDiscordNotification(message, status string) error {
	logger.Info("Sending Discord notification", "status", status)

	// Determine color based on status
	color := 0x00ff00 // Green for success
	if status == "error" {
		color = 0xff0000 // Red for error
	}

	// Create Discord embed
	embed := map[string]interface{}{
		"title":       "Database Backup Utility",
		"description": message,
		"color":       color,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	// Create Discord payload
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	// Send HTTP request
	resp, err := http.Post(m.config.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Discord notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Discord notification failed with status: %d", resp.StatusCode)
	}

	logger.Info("Discord notification sent successfully")
	return nil
}
