package main

import (
	"fmt"
	"log"
	// "github.com/gin-gonic/gin"
	"url-shortening-service/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("Connected to PostgreSQL at:", cfg.Database.Host)
	// Connect to the database
	// db := sqlc.NewDB(cfg.DatabaseURL)
	//
	// // Initialize services
	// urlService := services.NewUrlService(db)
	//
	// // Initialize handlers
	// urlHandler := handlers.NewUrlHandler(urlService)
	//
	// // Set up the router
	// router := gin.Default()
	// routes.SetupRoutes(router, urlHandler)
	//
	// // Start the server
	// router.Run(":" + cfg.Port)
}
