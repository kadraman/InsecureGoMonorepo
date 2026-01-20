package database

import (
	"database/sql"
	"fmt"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(host, user, password string, port int) (*Database, error) {
	// Note: This is a mock implementation for demonstration
	// In a real scenario, this would connect to an actual database
	return &Database{}, nil
}

// VULNERABILITY: SQL Injection
// ExecuteQuery executes a raw SQL query - vulnerable to SQL injection
func (db *Database) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	// Intentionally vulnerable: Direct string concatenation in SQL
	// In real implementation, this would execute against actual database
	fmt.Printf("Executing query: %s\n", query)
	return nil, nil
}

// VULNERABILITY: SQL Injection
// GetUserByUsername retrieves user by username - vulnerable to SQL injection
func (db *Database) GetUserByUsername(username string) (map[string]interface{}, error) {
	// Intentionally vulnerable: String concatenation in SQL query
	query := "SELECT * FROM users WHERE username = '" + username + "'"
	fmt.Printf("Executing query: %s\n", query)
	// Mock password hash for testing
	hashedPassword := "482c811da5d5b4bc6d497ffa98491e38" // MD5 hash of "password123"
	return map[string]interface{}{
		"id":       1,
		"username": username,
		"email":    "user@example.com",
		"password": hashedPassword,
	}, nil
}

// VULNERABILITY: SQL Injection
// CreateUser creates a new user - vulnerable to SQL injection
func (db *Database) CreateUser(username, email, password string) error {
	// Intentionally vulnerable: String concatenation in SQL query
	query := fmt.Sprintf("INSERT INTO users (username, email, password) VALUES ('%s', '%s', '%s')",
		username, email, password)
	fmt.Printf("Executing query: %s\n", query)
	return nil
}

// VULNERABILITY: Weak Password Hashing
// HashPassword uses weak hashing algorithm
func HashPassword(password string) string {
	// Intentionally vulnerable: Using weak hashing
	// Should use bcrypt or similar
	return fmt.Sprintf("md5:%s", password) // Simulating weak hash
}
