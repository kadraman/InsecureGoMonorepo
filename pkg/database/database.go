package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

var schemaEmbedded = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT,
	password TEXT
);

CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	description TEXT,
	price REAL,
	category TEXT
);

CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	product_id INTEGER,
	quantity INTEGER,
	total_price REAL,
	status TEXT,
	created_at TEXT,
	user_snapshot TEXT,
	product_snapshot TEXT
);
`

// NewDatabase opens a real in-memory SQLite database for demonstration.
// The parameters are accepted to keep the API compatible but are ignored.
func NewDatabase(host, user, password string, port int) (*Database, error) {
	conn, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	db := &Database{conn: conn}

	// Prefer loading schema from external file; fall back to embedded schema
	schemaPath := os.Getenv("DB_SCHEMA_FILE")
	if schemaPath == "" {
		schemaPath = "pkg/database/schema.sql"
	}

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		// fallback: embedded schema (same as pkg/database/schema.sql)
		embedded := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT,
			password TEXT
		);

		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			description TEXT,
			price REAL,
			category TEXT
		);

		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			product_id INTEGER,
			quantity INTEGER,
			total_price REAL,
			status TEXT,
			created_at TEXT,
			user_snapshot TEXT,
			product_snapshot TEXT
		);
		`
		schemaBytes = []byte(embedded)
	}

	if _, err := conn.Exec(string(schemaBytes)); err != nil {
		conn.Close()
		return nil, err
	}

	// Auto-seed database if requested via environment
	if err := db.AutoSeedFromEnv(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// Close closes the underlying DB connection.
func (db *Database) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

// ExecuteQuery executes a raw SQL query and returns rows as a slice of maps.
// NOTE: Intentionally accepts raw SQL (vulnerable to injection) for demo purposes.
func (db *Database) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	fmt.Printf("Executing query: %s\n", query)
	// Use Exec for non-SELECT statements to support INSERT/UPDATE/DELETE
	isSelect := strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "SELECT")
	if !isSelect {
		_, err := db.conn.Exec(query)
		return nil, err
	}

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		colsData := make([]interface{}, len(cols))
		colsDataPtr := make([]interface{}, len(cols))
		for i := range colsData {
			colsDataPtr[i] = &colsData[i]
		}

		if err := rows.Scan(colsDataPtr...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{}, len(cols))
		for i, colName := range cols {
			rowMap[colName] = colsData[i]
		}
		results = append(results, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// GetUserByUsername retrieves a user by building a raw, concatenated SQL query.
// This is intentionally vulnerable to demonstrate SQL injection.
func (db *Database) GetUserByUsername(username string) (map[string]interface{}, error) {
	query := "SELECT id, username, email, password FROM users WHERE username = '" + username + "'"
	fmt.Printf("Executing query: %s\n", query)

	row := db.conn.QueryRow(query)
	var id int
	var uname, email, password string
	if err := row.Scan(&id, &uname, &email, &password); err != nil {
		if err == sql.ErrNoRows {
			// Preserve previous mock behavior for tests: return a fake user
			hashedPassword := "482c811da5d5b4bc6d497ffa98491e38" // MD5 of "password123"
			return map[string]interface{}{
				"id":       1,
				"username": username,
				"email":    "user@example.com",
				"password": hashedPassword,
			}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"id":       id,
		"username": uname,
		"email":    email,
		"password": password,
	}, nil
}

// CreateUser inserts a user using direct string formatting (vulnerable to injection).
func (db *Database) CreateUser(username, email, password string) error {
	query := fmt.Sprintf("INSERT INTO users (username, email, password) VALUES ('%s', '%s', '%s')",
		username, email, password)
	fmt.Printf("Executing query: %s\n", query)

	_, err := db.conn.Exec(query)
	return err
}

// HashPassword uses a weak MD5-based hash (INTENTIONALLY weak for demo).
func HashPassword(password string) string {
	h := md5.Sum([]byte(password))
	return "md5:" + hex.EncodeToString(h[:])
}

// SeedUser is a small helper to insert a user using HashPassword.
func (db *Database) SeedUser(username, email, plainPassword string) error {
	hashed := HashPassword(plainPassword)
	if err := db.CreateUser(username, email, hashed); err != nil {
		// Unique constraint may error - report for demo purposes
		log.Printf("Seed user error: %v", err)
		return err
	}
	return nil
}
