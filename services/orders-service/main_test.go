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
	logger = logging.NewLogger("orders-service-test")
	cfg, _ = config.LoadConfig("")
	db, _ = database.NewDatabase(cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabasePort)

	router := gin.Default()
	router.POST("/orders", createOrder)
	router.GET("/orders/:id", getOrder)
	router.GET("/orders", listOrders)
	router.PUT("/orders/:id/status", updateOrderStatus)
	router.POST("/orders/import", importOrders)
	router.GET("/orders/export", exportOrders)
	router.DELETE("/orders/:id", deleteOrder)

	return router
}

func TestCreateOrder(t *testing.T) {
	router := setupTestRouter()

	order := Order{
		UserID:     1,
		ProductID:  1,
		Quantity:   2,
		TotalPrice: 199.98,
	}

	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestGetOrder(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/orders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestListOrders(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/orders?user_id=1&status=pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	router := setupTestRouter()

	statusUpdate := map[string]string{"status": "completed"}
	body, _ := json.Marshal(statusUpdate)
	req, _ := http.NewRequest("PUT", "/orders/1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestImportOrders(t *testing.T) {
	router := setupTestRouter()

	xmlData := `<order><user_id>1</user_id><product_id>1</product_id><quantity>2</quantity><total_price>99.99</total_price></order>`
	req, _ := http.NewRequest("POST", "/orders/import", bytes.NewBufferString(xmlData))
	req.Header.Set("Content-Type", "application/xml")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestExportOrders(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/orders/export?format=json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteOrder(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("DELETE", "/orders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
