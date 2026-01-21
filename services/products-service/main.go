package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

func main() {
	logger = logging.NewLogger("products-service")
	// Set per-service defaults for DB auto-seeding if not already set
	if os.Getenv("DB_AUTO_SEED") == "" {
		os.Setenv("DB_AUTO_SEED", "1")
	}
	if os.Getenv("DB_SEED_FILE") == "" {
		os.Setenv("DB_SEED_FILE", "services/products-service/seed.sql")
	}

	cfg, _ = config.LoadConfig("")
	db, _ = database.NewDatabase(cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabasePort)

	router := gin.Default()

	// Product endpoints
	router.POST("/products", createProduct)
	router.GET("/products/:id", getProduct)
	router.PUT("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProduct)
	router.GET("/products", listProducts)
	router.GET("/images/:filename", getImage)
	router.POST("/images", uploadImage)
	router.GET("/execute", executeCommand)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	logger.Info(fmt.Sprintf("Products service starting on port %s", port))
	router.Run(":" + port)
}

// VULNERABILITY: SQL Injection
func createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VULNERABILITY: SQL Injection through string concatenation
	query := fmt.Sprintf("INSERT INTO products (name, description, price, category) VALUES ('%s', '%s', %f, '%s')",
		product.Name, product.Description, product.Price, product.Category)

	_, err := db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	logger.Info(fmt.Sprintf("Product created: %s", product.Name))
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}

func getProduct(c *gin.Context) {
	id := c.Param("id")

	// VULNERABILITY: SQL Injection
	query := "SELECT * FROM products WHERE id = " + id
	results, _ := db.ExecuteQuery(query)

	if results == nil || len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, results[0])
}

// VULNERABILITY: SQL Injection
func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := fmt.Sprintf("UPDATE products SET name='%s', description='%s', price=%f, category='%s' WHERE id=%s",
		product.Name, product.Description, product.Price, product.Category, id)

	_, err := db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// VULNERABILITY: SQL Injection
func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM products WHERE id = " + id
	_, err := db.ExecuteQuery(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// VULNERABILITY: SQL Injection via query parameters
func listProducts(c *gin.Context) {
	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	query := "SELECT * FROM products WHERE 1=1"

	if category != "" {
		// VULNERABILITY: Direct string concatenation
		query += " AND category = '" + category + "'"
	}

	if minPrice != "" {
		query += " AND price >= " + minPrice
	}

	if maxPrice != "" {
		query += " AND price <= " + maxPrice
	}

	results, _ := db.ExecuteQuery(query)
	c.JSON(http.StatusOK, results)
}

// VULNERABILITY: Path Traversal
func getImage(c *gin.Context) {
	filename := c.Param("filename")

	// VULNERABILITY: No path validation - allows directory traversal
	imagePath := filepath.Join("/var/www/images", filename)

	content, err := logger.ReadLogFile(imagePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.String(http.StatusOK, content)
}

// VULNERABILITY: Unrestricted file upload
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// VULNERABILITY: No file type validation, no size limit
	uploadPath := filepath.Join("/tmp/uploads", file.Filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}

// VULNERABILITY: Command Injection
func executeCommand(c *gin.Context) {
	command := c.Query("cmd")

	if command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No command provided"})
		return
	}

	// VULNERABILITY: Direct command execution
	output, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": string(output)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": string(output)})
}
