package database

import (
	"os"
	"testing"
)

// TestSeedFromFile verifies that SeedFromFile inserts rows into the DB.
func TestSeedFromFile(t *testing.T) {
	// Ensure auto-seed is off for this test
	prevAuto := os.Getenv("DB_AUTO_SEED")
	os.Setenv("DB_AUTO_SEED", "0")
	defer os.Setenv("DB_AUTO_SEED", prevAuto)

	db, err := NewDatabase("", "", "", 0)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	// Create a temp seed file
	f, err := os.CreateTemp("", "seed-*.sql")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	sql := "INSERT INTO users (username, email, password) VALUES ('seed_test','seed@test','md5:abc');"
	if _, err := f.WriteString(sql); err != nil {
		t.Fatalf("failed to write seed file: %v", err)
	}
	f.Close()

	if err := db.SeedFromFile(f.Name()); err != nil {
		t.Fatalf("SeedFromFile failed: %v", err)
	}

	rows, err := db.ExecuteQuery("SELECT username FROM users WHERE username='seed_test'")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected seeded user to exist")
	}
}

// TestAutoSeedFromEnv verifies that NewDatabase auto-seeds when DB_AUTO_SEED is set.
func TestAutoSeedFromEnv(t *testing.T) {
	// Create a temp seed file
	f, err := os.CreateTemp("", "seed-*.sql")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	sql := "INSERT INTO products (name, description, price, category) VALUES ('auto_seed_prod','auto','1.23','auto');"
	if _, err := f.WriteString(sql); err != nil {
		t.Fatalf("failed to write seed file: %v", err)
	}
	f.Close()

	prevAuto := os.Getenv("DB_AUTO_SEED")
	prevFile := os.Getenv("DB_SEED_FILE")
	os.Setenv("DB_AUTO_SEED", "1")
	os.Setenv("DB_SEED_FILE", f.Name())
	defer os.Setenv("DB_AUTO_SEED", prevAuto)
	defer os.Setenv("DB_SEED_FILE", prevFile)

	db, err := NewDatabase("", "", "", 0)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	rows, err := db.ExecuteQuery("SELECT name FROM products WHERE name='auto_seed_prod'")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected auto-seeded product to exist")
	}
}
