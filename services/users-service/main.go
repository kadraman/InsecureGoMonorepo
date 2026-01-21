package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

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

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	logger = logging.NewLogger("users-service")
	// Set per-service defaults for DB auto-seeding if not already set
	if os.Getenv("DB_AUTO_SEED") == "" {
		os.Setenv("DB_AUTO_SEED", "1")
	}
	if os.Getenv("DB_SEED_FILE") == "" {
		os.Setenv("DB_SEED_FILE", "services/users-service/seed.sql")
	}

	cfg, _ = config.LoadConfig("")
	db, _ = database.NewDatabase(cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabasePort)

	router := gin.Default()

	// User endpoints
	router.POST("/users", createUser)
	router.GET("/users/:username", getUser)
	router.GET("/users/id/:id", getUserByID)
	router.POST("/login", login)
	router.GET("/search", searchUsers)
	router.GET("/export", exportUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger.Info(fmt.Sprintf("Users service starting on port %s", port))
	router.Run(":" + port)
}

// VULNERABILITY: SQL Injection
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VULNERABILITY: Weak password hashing (MD5)
	hashedPassword := hashPasswordMD5(user.Password)

	err := db.CreateUser(user.Username, user.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	logger.Info(fmt.Sprintf("User created: %s", user.Username))
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// VULNERABILITY: SQL Injection
func getUser(c *gin.Context) {
	username := c.Param("username")

	user, err := db.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// getUserByID retrieves a user by numeric id (used by other services for snapshotting)
func getUserByID(c *gin.Context) {
	id := c.Param("id")

	query := "SELECT id, username, email FROM users WHERE id = " + id
	results, err := db.ExecuteQuery(query)
	if err != nil || results == nil || len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the first row
	c.JSON(http.StatusOK, results[0])
}

// VULNERABILITY: SQL Injection via search parameter
func searchUsers(c *gin.Context) {
	searchTerm := c.Query("q")

	// Intentionally vulnerable: Direct string concatenation
	query := "SELECT * FROM users WHERE username LIKE '%" + searchTerm + "%'"

	results, _ := db.ExecuteQuery(query)
	c.JSON(http.StatusOK, results)
}

// VULNERABILITY: Weak authentication
func login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := db.GetUserByUsername(loginReq.Username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// VULNERABILITY: Weak password comparison
	hashedInput := hashPasswordMD5(loginReq.Password)
	storedPassword := user["password"].(string)

	if hashedInput != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// VULNERABILITY: Insecure token generation using hardcoded secret
	token := generateInsecureToken(loginReq.Username)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// VULNERABILITY: Command Injection via export functionality
func exportUsers(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		filename = "users.txt"
	}

	// VULNERABILITY: Command injection through logger
	logger.LogToFile(filename, "User export requested")

	c.JSON(http.StatusOK, gin.H{"message": "Export completed"})
}

// VULNERABILITY: Weak hashing algorithm (MD5)
func hashPasswordMD5(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VULNERABILITY: Insecure token generation
func generateInsecureToken(username string) string {
	// Using hardcoded secret from config
	data := username + config.JWTSecret
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
