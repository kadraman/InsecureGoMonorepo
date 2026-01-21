package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadConfigDefault(t *testing.T) {
	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if config.DatabasePassword != DefaultDBPassword {
		t.Errorf("Expected default password, got %s", config.DatabasePassword)
	}
	if config.APIKey != APIKey {
		t.Errorf("Expected default API key, got %s", config.APIKey)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testConfig := &Config{
		DatabaseHost: "testhost",
		DatabasePort: 3306,
		ServerPort:   9090,
	}

	data, _ := json.Marshal(testConfig)
	tmpFile.Write(data)
	tmpFile.Close()

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if config.DatabaseHost != "testhost" {
		t.Errorf("Expected 'testhost', got %s", config.DatabaseHost)
	}
}

func TestGetConnectionString(t *testing.T) {
	config := &Config{
		DatabaseHost:     "localhost",
		DatabaseUser:     "user",
		DatabasePassword: "pass",
	}
	connStr := config.GetConnectionString()
	expected := "user:pass@localhost"
	if connStr != expected {
		t.Errorf("Expected '%s', got '%s'", expected, connStr)
	}
}
