package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "log"
	"url-shortening-service/cache"
	"url-shortening-service/config"
	"url-shortening-service/generate"
	"url-shortening-service/internal/rest"
	"url-shortening-service/repository"
	"url-shortening-service/retrieve"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("Failed to load config: %v", err)
	}
	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port,
		cfg.User, cfg.Password,
		cfg.DBName, cfg.SSLMode)
	log.Debug().Msgf("connStr: %s", connStr)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal().Msgf("Cannot connect database %s", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal().Msgf("Failed to ping database")
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal().Msgf("Got the error when closing the DB connection")
		}
	}()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		log.Fatal().Err(err)
	}
	dragonFlyCache := cache.NewDrangonFlyCache(client)
	router := gin.Default()
	router.Use(gin.Recovery())

	nanoIdGenerator := generate.NewNannoIdGenerator()
	postgresStore := repository.NewPostgresStore(dbConn)
	generateService := generate.NewGenerateService(postgresStore, nanoIdGenerator, dragonFlyCache)
	gernerateHandler := rest.NewGeneateHandler(generateService)
	router.POST("/shorten", gernerateHandler.Generate)

	retrieveService := retrieve.NewRetrieveService(postgresStore, dragonFlyCache)
	retrieveHandler := rest.NewRetrieveHandler(retrieveService)
	router.GET("/:shortUrl", retrieveHandler.Retrieve)
	router.Run(":8080")
	fmt.Println("Server is running on port 8080", cfg)
}
