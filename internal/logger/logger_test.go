package logger

import (
	"testing"
)

func TestInit(t *testing.T) {
	// Test initialization with default values
	Init("info", "text")

	// Test initialization with different values
	Init("debug", "json")

	// If we get here without panicking, the test passes
}
