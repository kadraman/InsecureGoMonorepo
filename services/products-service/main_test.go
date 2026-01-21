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
	logger = logging.NewLogger("products-service-test")
	cfg, _ = config.LoadConfig("")
	db, _ = database.NewDatabase(cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabasePort)

	router := gin.Default()
	router.POST("/products", createProduct)
	router.GET("/products/:id", getProduct)
	router.PUT("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProduct)
	router.GET("/products", listProducts)
	router.GET("/images/:filename", getImage)

	return router
}

func TestCreateProduct(t *testing.T) {
	router := setupTestRouter()

	product := Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Category:    "Electronics",
	}

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestGetProduct(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestUpdateProduct(t *testing.T) {
	router := setupTestRouter()

	product := Product{
		Name:        "Updated Product",
		Description: "An updated product",
		Price:       149.99,
		Category:    "Electronics",
	}

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteProduct(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListProducts(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/products?category=Electronics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
