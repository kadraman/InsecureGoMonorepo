package config

import (
	"encoding/json"
	"os"
)

// VULNERABILITY: Hardcoded Secrets
const (
	DefaultDBPassword = "admin123"            // Hardcoded password
	APIKey            = "sk-1234567890abcdef" // Hardcoded API key
	JWTSecret         = "my-secret-key"       // Hardcoded JWT secret
)

type Config struct {
	DatabaseHost     string `json:"database_host"`
	DatabasePort     int    `json:"database_port"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	ServerPort       int    `json:"server_port"`
	LogLevel         string `json:"log_level"`
	APIKey           string `json:"api_key"`
	JWTSecret        string `json:"jwt_secret"`
}

func LoadConfig(filepath string) (*Config, error) {
	// Default configuration with hardcoded credentials
	config := &Config{
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseUser:     "admin",
		DatabasePassword: DefaultDBPassword,
		ServerPort:       8080,
		LogLevel:         "info",
		APIKey:           APIKey,
		JWTSecret:        JWTSecret,
	}

	if filepath != "" {
		data, err := os.ReadFile(filepath)
		if err != nil {
			return config, nil // Return default config on error
		}
		json.Unmarshal(data, config) // Ignoring error intentionally
	}

	return config, nil
}

// VULNERABILITY: Insecure configuration
func (c *Config) GetConnectionString() string {
	// Returns connection string with password in plain text
	return c.DatabaseUser + ":" + c.DatabasePassword + "@" + c.DatabaseHost
}
