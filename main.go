package main

import (
	"database/sql"
	"fmt"
	"log"
	"url-shortening-service/config"
	"url-shortening-service/generate"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Connect to the database
	dbConn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.SSLMode),
	)
	defer dbConn.Close()
	if err != nil {
		log.Fatalf("Cannot connect database %s", err)
	}
	router := gin.Default()
	generateService := generate.NewGenerateService(dbConn)
	gernerateHandler := generate.NewGeneateHandler(generateService)
	router.POST("/api/urls", gernerateHandler.Generate)
	router.Run(":8080")
	fmt.Println("Server is running on port 8080", cfg)
}
