package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadraman/InsecureGoMonorepo/pkg/config"
	"github.com/kadraman/InsecureGoMonorepo/pkg/database"
	"github.com/kadraman/InsecureGoMonorepo/pkg/logging"
)

var (
	db     *database.Database
	logger *logging.Logger
	cfg    *config.Config
)

type Order struct {
	ID         int       `json:"id" xml:"id"`
	UserID     int       `json:"user_id" xml:"user_id"`
	ProductID  int       `json:"product_id" xml:"product_id"`
	Quantity   int       `json:"quantity" xml:"quantity"`
	TotalPrice float64   `json:"total_price" xml:"total_price"`
	Status     string    `json:"status" xml:"status"`
	CreatedAt  time.Time `json:"created_at" xml:"created_at"`
}

type XMLOrder struct {
	XMLName    xml.Name `xml:"order"`
	ID         int      `xml:"id"`
	UserID     int      `xml:"user_id"`
	ProductID  int      `xml:"product_id"`
	Quantity   int      `xml:"quantity"`
	TotalPrice float64  `xml:"total_price"`
}

func main() {
	logger = logging.NewLogger("orders-service")
	cfg, _ = config.LoadConfig("")
	db, _ = database.NewDatabase(cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabasePort)

	router := gin.Default()

	// Order endpoints
	router.POST("/orders", createOrder)
	router.GET("/orders/:id", getOrder)
	router.GET("/orders", listOrders)
	router.PUT("/orders/:id/status", updateOrderStatus)
	router.POST("/orders/import", importOrders)
	router.GET("/orders/export", exportOrders)
	router.DELETE("/orders/:id", deleteOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	logger.Info(fmt.Sprintf("Orders service starting on port %s", port))
	router.Run(":" + port)
}

// VULNERABILITY: SQL Injection
func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.CreatedAt = time.Now()
	order.Status = "pending"

	// VULNERABILITY: SQL Injection
	query := fmt.Sprintf("INSERT INTO orders (user_id, product_id, quantity, total_price, status, created_at) VALUES (%d, %d, %d, %f, '%s', '%s')",
		order.UserID, order.ProductID, order.Quantity, order.TotalPrice, order.Status, order.CreatedAt.Format(time.RFC3339))

	_, err := db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	logger.Info(fmt.Sprintf("Order created for user %d", order.UserID))
	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order": order})
}

// VULNERABILITY: SQL Injection
func getOrder(c *gin.Context) {
	id := c.Param("id")

	query := "SELECT * FROM orders WHERE id = " + id
	results, _ := db.ExecuteQuery(query)

	if results == nil || len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, results[0])
}

// VULNERABILITY: SQL Injection via query parameters
func listOrders(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.Query("status")
	sortBy := c.Query("sort_by")

	query := "SELECT * FROM orders WHERE 1=1"

	if userID != "" {
		// VULNERABILITY: No input validation
		query += " AND user_id = " + userID
	}

	if status != "" {
		// VULNERABILITY: Direct string concatenation
		query += " AND status = '" + status + "'"
	}

	if sortBy != "" {
		// VULNERABILITY: Dangerous - allows ORDER BY injection
		query += " ORDER BY " + sortBy
	}

	results, _ := db.ExecuteQuery(query)
	c.JSON(http.StatusOK, results)
}

// VULNERABILITY: SQL Injection
func updateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VULNERABILITY: Direct string concatenation
	query := "UPDATE orders SET status = '" + request.Status + "' WHERE id = " + id

	_, err := db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

// VULNERABILITY: XXE (XML External Entity) Attack
func importOrders(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	var xmlOrder XMLOrder
	// VULNERABILITY: Unsafe XML parsing - vulnerable to XXE
	err = xml.Unmarshal(body, &xmlOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML: " + err.Error()})
		return
	}

	// Create order from XML data
	query := fmt.Sprintf("INSERT INTO orders (user_id, product_id, quantity, total_price, status) VALUES (%d, %d, %d, %f, 'pending')",
		xmlOrder.UserID, xmlOrder.ProductID, xmlOrder.Quantity, xmlOrder.TotalPrice)

	_, err = db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order imported successfully"})
}

// VULNERABILITY: Information disclosure through export
func exportOrders(c *gin.Context) {
	format := c.Query("format")
	userID := c.Query("user_id")

	query := "SELECT * FROM orders"
	if userID != "" {
		// VULNERABILITY: SQL Injection
		query += " WHERE user_id = " + userID
	}

	results, _ := db.ExecuteQuery(query)

	if format == "xml" {
		c.XML(http.StatusOK, results)
	} else {
		// VULNERABILITY: May expose sensitive data
		c.JSON(http.StatusOK, results)
	}
}

// VULNERABILITY: SQL Injection
func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	// VULNERABILITY: Direct string concatenation
	query := "DELETE FROM orders WHERE id = " + id

	_, err := db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
