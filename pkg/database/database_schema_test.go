package database

import (
	"os"
	"testing"
)

// Test that NewDatabase loads schema from DB_SCHEMA_FILE when set.
func TestNewDatabaseLoadsSchemaFile(t *testing.T) {
	// Create a temporary SQL file defining a custom table.
	f, err := os.CreateTemp("", "schema-*.sql")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	custom := `CREATE TABLE IF NOT EXISTS custom_table (id INTEGER PRIMARY KEY);`
	if _, err := f.WriteString(custom); err != nil {
		t.Fatalf("failed to write temp schema: %v", err)
	}
	f.Close()

	// Ensure auto-seed disabled for test
	prevAuto := os.Getenv("DB_AUTO_SEED")
	os.Setenv("DB_AUTO_SEED", "0")
	defer os.Setenv("DB_AUTO_SEED", prevAuto)

	// Point NewDatabase at our temp schema file
	prev := os.Getenv("DB_SCHEMA_FILE")
	os.Setenv("DB_SCHEMA_FILE", f.Name())
	defer os.Setenv("DB_SCHEMA_FILE", prev)

	db, err := NewDatabase("", "", "", 0)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	// Query sqlite_master for custom_table
	rows, err := db.ExecuteQuery("SELECT name FROM sqlite_master WHERE type='table' AND name='custom_table'")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected custom_table to exist based on DB_SCHEMA_FILE")
	}
}
