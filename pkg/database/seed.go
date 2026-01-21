package database

import (
	"io/ioutil"
	"os"
	"strings"
)

// SeedFromFile reads SQL from a file and executes it against the database.
// The file may contain multiple statements separated by semicolons.
func (db *Database) SeedFromFile(path string) error {
	if db == nil || db.conn == nil {
		return nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	sql := string(b)
	// Execute the SQL directly. Drivers typically accept multiple statements.
	_, err = db.conn.Exec(sql)
	return err
}

// AutoSeedFromEnv seeds the DB if the DB_AUTO_SEED env var is set to true/1.
// It reads DB_SEED_FILE to get the path to the seed SQL file (default: pkg/database/seed.sql).
func (db *Database) AutoSeedFromEnv() error {
	val := os.Getenv("DB_AUTO_SEED")
	if val == "1" || strings.ToLower(val) == "true" {
		path := os.Getenv("DB_SEED_FILE")
		if path == "" {
			path = "pkg/database/seed.sql"
		}
		return db.SeedFromFile(path)
	}
	return nil
}
