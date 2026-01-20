package logging

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(message string) {
	log.Printf("[%s] INFO: %s\n", l.prefix, message)
}

func (l *Logger) Error(message string) {
	log.Printf("[%s] ERROR: %s\n", l.prefix, message)
}

func (l *Logger) Warn(message string) {
	log.Printf("[%s] WARN: %s\n", l.prefix, message)
}

// VULNERABILITY: Command Injection
// LogToFile executes a shell command to write logs - vulnerable to command injection
func (l *Logger) LogToFile(filename, message string) error {
	// Intentionally vulnerable: Directly using user input in shell command
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' >> %s", message, filename))
	return cmd.Run()
}

// VULNERABILITY: Path Traversal
// ReadLogFile reads a log file - vulnerable to path traversal
func (l *Logger) ReadLogFile(filename string) (string, error) {
	// Intentionally vulnerable: No path validation
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
