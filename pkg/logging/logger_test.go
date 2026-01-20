package logging

import (
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test")
	if logger == nil {
		t.Fatal("Expected logger to be created")
	}
	if logger.prefix != "test" {
		t.Errorf("Expected prefix 'test', got '%s'", logger.prefix)
	}
}

func TestLoggerInfo(t *testing.T) {
	logger := NewLogger("test")
	// This just ensures the method doesn't panic
	logger.Info("test message")
}

func TestLoggerError(t *testing.T) {
	logger := NewLogger("test")
	logger.Error("test error")
}

func TestLoggerWarn(t *testing.T) {
	logger := NewLogger("test")
	logger.Warn("test warning")
}

func TestReadLogFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-log-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := "test log content"
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	logger := NewLogger("test")
	content, err := logger.ReadLogFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if content != testContent {
		t.Errorf("Expected '%s', got '%s'", testContent, content)
	}
}
