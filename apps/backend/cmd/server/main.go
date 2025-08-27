package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/config"
	"github.com/zeusnotfound04/Tranza/models"
	"github.com/zeusnotfound04/Tranza/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()
	defer config.CloseDB(db)

	// Auto-migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.EmailVerification{},
		&models.Transaction{},
		&models.APIKey{},
		&models.APIUsageLog{}, // Added API usage logging table
		&models.LinkedCard{},
		&models.Wallet{},
		&models.Address{}, // Added address management table
		&models.AIPaymentRequest{},
		&models.AISpendingLimit{},
		&models.AISpendingTracker{},
		&models.ExternalTransfer{}, // Added missing external transfers table
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS for HttpOnly cookies
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		os.Getenv("FRONTEND_URL"), // e.g., "http://localhost:3000"
		"http://localhost:3000",   // fallback for development
		"http://localhost:3001",   // docs app
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-API-Key",
		"X-Requested-With",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	router.Use(cors.New(corsConfig))

	// Add custom debug middleware to log all requests
	router.Use(func(c *gin.Context) {
		fmt.Printf("🔍 DEBUG: %s %s | Headers: %v | Cookies: %v\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Header.Get("Authorization"),
			c.Request.Header.Get("Cookie"))
		c.Next()
	})

	// Add request logging middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup all routes using the comprehensive routes system
	routes.SetupRoutes(router, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Tranza API Server starting on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/ping", port)
	log.Printf("🔐 Auth endpoints: http://localhost:%s/auth/*", port)
	log.Printf("📊 API endpoints: http://localhost:%s/api/v1/*", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
