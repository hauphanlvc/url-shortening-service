package main

import (
	"database/sql"
	"fmt"
	"log"
	"url-shortening-service/config"
	"url-shortening-service/generate"
	"url-shortening-service/retrieve"

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
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.SSLMode)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect database %s", err)
	}
	defer dbConn.Close()

	router := gin.Default()
	generateService := generate.NewGenerateService(dbConn)
	gernerateHandler := generate.NewGeneateHandler(generateService)
	router.POST("/shorten", gernerateHandler.Generate)

	retrieveService := retrieve.NewRetrieveService(dbConn)
	retrieveHandler := retrieve.NewRetrieveHandler(retrieveService)
	router.GET("/:shortUrl", retrieveHandler.Retrieve)
	router.Run(":8080")
	fmt.Println("Server is running on port 8080", cfg)
}
