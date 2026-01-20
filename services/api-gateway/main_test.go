package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kadraman/InsecureGoMonorepo/pkg/config"
	"github.com/kadraman/InsecureGoMonorepo/pkg/logging"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger = logging.NewLogger("api-gateway-test")
	cfg, _ = config.LoadConfig("")

	router := gin.Default()
	router.Use(insecureAuthMiddleware())
	router.GET("/api/debug", debugInfo)
	router.GET("/api/redirect", openRedirect)

	return router
}

func TestInsecureAuthMiddleware(t *testing.T) {
	router := setupTestRouter()

	// Test without API key
	req, _ := http.NewRequest("GET", "/api/debug", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d without API key, got %d", http.StatusOK, w.Code)
	}

	// Test with correct API key
	req, _ = http.NewRequest("GET", "/api/debug", nil)
	req.Header.Set("X-API-Key", config.APIKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d with correct API key, got %d", http.StatusOK, w.Code)
	}

	// Test with incorrect API key
	req, _ = http.NewRequest("GET", "/api/debug", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d with wrong API key, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDebugInfo(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/debug", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check if response contains sensitive info
	body := w.Body.String()
	if body == "" {
		t.Error("Expected debug info in response body")
	}
}

func TestOpenRedirect(t *testing.T) {
	router := setupTestRouter()

	// Test with redirect URL
	req, _ := http.NewRequest("GET", "/api/redirect?url=http://example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status %d, got %d", http.StatusFound, w.Code)
	}

	// Test without redirect URL
	req, _ = http.NewRequest("GET", "/api/redirect", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetEnv(t *testing.T) {
	result := getEnv("NONEXISTENT_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}
