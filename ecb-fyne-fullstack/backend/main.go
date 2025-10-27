package main

import (
	"backend/database"
	"backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present (useful for Windows demo installations)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	// Ensure default host/port helpful for local Windows demo usage
	ensureDefaults()

	// Connect to MySQL and run migrations
	database.Connect()
	database.Migrate()

	// Setup Gin router with basic CORS for both desktop (Windows) and Pi clients
	router := gin.Default()
	router.Use(corsMiddleware)

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	addr := ":" + os.Getenv("API_PORT")
	log.Printf("Server running on %s (DB host: %s)\n", addr, os.Getenv("DB_HOST"))
	if err := router.Run(addr); err != nil {
		log.Fatal("Server terminated:", err)
	}
}

func ensureDefaults() {
	setIfEmpty("API_PORT", "8080")
	setIfEmpty("DB_HOST", "127.0.0.1")
	setIfEmpty("DB_PORT", "3306")
	setIfEmpty("DB_USER", "root")
	setIfEmpty("DB_NAME", "ecb_test")
}

func setIfEmpty(key, fallback string) {
	if os.Getenv(key) == "" {
		if err := os.Setenv(key, fallback); err != nil {
			log.Printf("Warning: unable to set default for %s: %v\n", key, err)
		}
	}
}

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
