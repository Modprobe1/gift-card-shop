package main

import (
	"log"
	"os"

	"auth-backend/config"
	"auth-backend/controllers"
	"auth-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to MongoDB
	config.ConnectDB()

	// Initialize Gin router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Public routes
	api := router.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", controllers.GetProfile)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}