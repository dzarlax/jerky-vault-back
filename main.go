package main

import (
	"log"
	"mobile-backend-go/database"
	_ "mobile-backend-go/docs" // Import for Swagger documentation
	"mobile-backend-go/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// loadEnvVar loads an environment variable, if not set, tries to load from .env
func loadEnvVar(key string) string {
	// First try to get the variable value from environment
	value := os.Getenv(key)
	// If variable is not in environment, try to load from .env
	if value == "" {
		// Load .env only if the variable value is empty
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
		// Try again to get the value after loading from .env
		value = os.Getenv(key)
	}
	return value
}

// @title Jerky-vault Backend API
// @version 1.0
// @description Jerky-vault project API
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Define required environment variables
	//requiredEnvVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "FRONT_URL"}
	requiredEnvVars := []string{"DATABASE_URL", "FRONT_URL"}

	// Check for all required environment variables
	for _, envVar := range requiredEnvVars {
		value := loadEnvVar(envVar)
		if value == "" {
			log.Fatalf("Environment variable %s is not set", envVar)
		}
	}

	// Connect to database and run migrations
	database.ConnectDatabase()

	// Create Gin router instance
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONT_URL"), "http://localhost:3000"}, // Next.js frontend domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Route for Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup API routes
	routes.SetupRoutes(r)

	// Start server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
