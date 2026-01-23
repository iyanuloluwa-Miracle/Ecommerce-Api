// @title E-Commerce API
// @version 1.0
// @description A RESTful e-commerce API built with Go, Gin framework, MongoDB, and Stripe for payment processing.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/iyanuloluwa-Miracle/ecommerce-api
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/config"
	_ "github.com/iyanuloluwa-Miracle/ecommerce-api/docs"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	config.ConnectDB()

	// Initialize Stripe
	config.InitStripe()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register routes
	routes.RegisterRoutes(r)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger docs available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
