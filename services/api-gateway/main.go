package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kadraman/InsecureGoMonorepo/pkg/config"
	"github.com/kadraman/InsecureGoMonorepo/pkg/logging"
)

var (
	logger *logging.Logger
	cfg    *config.Config
)

// Service URLs
var (
	usersServiceURL    = getEnv("USERS_SERVICE_URL", "http://localhost:8081")
	productsServiceURL = getEnv("PRODUCTS_SERVICE_URL", "http://localhost:8082")
	ordersServiceURL   = getEnv("ORDERS_SERVICE_URL", "http://localhost:8083")
)

func main() {
	logger = logging.NewLogger("api-gateway")
	cfg, _ = config.LoadConfig("")

	router := gin.Default()

	// VULNERABILITY: No authentication middleware
	router.Use(insecureAuthMiddleware())

	// User routes
	router.POST("/api/users", proxyRequest(usersServiceURL, "/users"))
	router.GET("/api/users/:username", proxyRequest(usersServiceURL, "/users/:username"))
	router.POST("/api/login", proxyRequest(usersServiceURL, "/login"))
	router.GET("/api/users/search", proxyRequest(usersServiceURL, "/search"))

	// Product routes
	router.POST("/api/products", proxyRequest(productsServiceURL, "/products"))
	router.GET("/api/products/:id", proxyRequest(productsServiceURL, "/products/:id"))
	router.PUT("/api/products/:id", proxyRequest(productsServiceURL, "/products/:id"))
	router.DELETE("/api/products/:id", proxyRequest(productsServiceURL, "/products/:id"))
	router.GET("/api/products", proxyRequest(productsServiceURL, "/products"))

	// Order routes
	router.POST("/api/orders", proxyRequest(ordersServiceURL, "/orders"))
	router.GET("/api/orders/:id", proxyRequest(ordersServiceURL, "/orders/:id"))
	router.GET("/api/orders", proxyRequest(ordersServiceURL, "/orders"))
	router.PUT("/api/orders/:id/status", proxyRequest(ordersServiceURL, "/orders/:id/status"))

	// VULNERABILITY: Debug endpoint exposing system info
	router.GET("/api/debug", debugInfo)

	// VULNERABILITY: Open redirect
	router.GET("/api/redirect", openRedirect)

	port := getEnv("PORT", "8080")
	logger.Info(fmt.Sprintf("API Gateway starting on port %s", port))
	router.Run(":" + port)
}

// VULNERABILITY: Weak authentication - accepts hardcoded API key
func insecureAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// VULNERABILITY: Hardcoded API key check
		if apiKey != "" && apiKey != config.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// VULNERABILITY: No token validation for authenticated routes
		c.Next()
	}
}

func proxyRequest(serviceURL, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build target URL
		targetURL := serviceURL + strings.Replace(path, ":", "", -1)
		for _, param := range c.Params {
			targetURL = strings.Replace(targetURL, param.Key, param.Value, 1)
		}

		// Add query parameters
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// Create proxy request
		proxyReq, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
			return
		}

		// VULNERABILITY: Copy all headers including potentially sensitive ones
		for key, values := range c.Request.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Execute request
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			logger.Error(fmt.Sprintf("Proxy request failed: %v", err))
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Read response
		respBody, _ := io.ReadAll(resp.Body)

		// VULNERABILITY: Copy all response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	}
}

// VULNERABILITY: Debug endpoint exposing sensitive information
func debugInfo(c *gin.Context) {
	info := gin.H{
		"services": gin.H{
			"users":    usersServiceURL,
			"products": productsServiceURL,
			"orders":   ordersServiceURL,
		},
		"config": gin.H{
			"database_host": cfg.DatabaseHost,
			"database_user": cfg.DatabaseUser,
			"database_pass": cfg.DatabasePassword, // VULNERABILITY: Exposing password
			"api_key":       cfg.APIKey,           // VULNERABILITY: Exposing API key
			"jwt_secret":    cfg.JWTSecret,        // VULNERABILITY: Exposing JWT secret
		},
		"environment": os.Environ(), // VULNERABILITY: Exposing all env vars
	}

	c.JSON(http.StatusOK, info)
}

// VULNERABILITY: Open redirect
func openRedirect(c *gin.Context) {
	redirectURL := c.Query("url")

	if redirectURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter required"})
		return
	}

	// VULNERABILITY: No validation of redirect URL
	c.Redirect(http.StatusFound, redirectURL)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
