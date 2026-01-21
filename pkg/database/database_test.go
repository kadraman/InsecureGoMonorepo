package database

import (
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase("localhost", "user", "pass", 5432)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if db == nil {
		t.Fatal("Expected database to be created")
	}
}

func TestExecuteQuery(t *testing.T) {
	db, _ := NewDatabase("localhost", "user", "pass", 5432)
	_, err := db.ExecuteQuery("SELECT * FROM users")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, _ := NewDatabase("localhost", "user", "pass", 5432)
	user, err := db.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("Expected user to be returned")
	}
}

func TestCreateUser(t *testing.T) {
	db, _ := NewDatabase("localhost", "user", "pass", 5432)
	err := db.CreateUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHashPassword(t *testing.T) {
	hash := HashPassword("password123")
	if hash == "" {
		t.Fatal("Expected hash to be generated")
	}
	if hash == "password123" {
		t.Error("Expected password to be hashed")
	}
}
