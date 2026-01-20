package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kadraman/InsecureGoMonorepo/pkg/config"
	"github.com/kadraman/InsecureGoMonorepo/pkg/database"
	"github.com/kadraman/InsecureGoMonorepo/pkg/logging"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger = logging.NewLogger("users-service-test")
	cfg, _ = config.LoadConfig("")
	db, _ = database.NewDatabase(cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabasePort)

	router := gin.Default()
	router.POST("/users", createUser)
	router.GET("/users/:username", getUser)
	router.POST("/login", login)
	router.GET("/search", searchUsers)
	router.GET("/export", exportUsers)

	return router
}

func TestCreateUser(t *testing.T) {
	router := setupTestRouter()

	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestGetUser(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/users/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestLogin(t *testing.T) {
	router := setupTestRouter()

	loginReq := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSearchUsers(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/search?q=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHashPasswordMD5(t *testing.T) {
	hash := hashPasswordMD5("password123")
	if hash == "" {
		t.Error("Expected hash to be generated")
	}
	if hash == "password123" {
		t.Error("Expected password to be hashed")
	}
}
