package main

import (
	"incident-management/database"
	"incident-management/handlers"
	"incident-management/utils"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize validator
	utils.InitValidator()

	// Create Gin router
	r := gin.Default()

	// Create handler
	handler := handlers.NewIncidentHandler()
	// Allow everything (for development/testing only)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Define routes
	api := r.Group("/api/v1")
	{
		api.POST("/incidents", handler.CreateIncident)
		api.GET("/incidents", handler.GetAllIncidents)
	}

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	// Start server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
