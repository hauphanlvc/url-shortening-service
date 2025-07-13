package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "log"
	"net/http"
	"url-shortening-service/config"
	"url-shortening-service/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"url-shortening-service/internal/repository/postgres"
	generate "url-shortening-service/internal/urls/generate/service"
	generatehttp "url-shortening-service/internal/urls/generate/transport/http"
	retrieve "url-shortening-service/internal/urls/retrieve/service"
	retrievehttp "url-shortening-service/internal/urls/retrieve/transport/http"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

 // Todo: Improve RateLimiter with client IP
func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 4)
	return func(c *gin.Context) {

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Limite exceed",
			})
		} else {
			c.Next()
		}

	}
}
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
	router.Use(RateLimiter())
	nanoIdGenerator := generate.NewNannoIdGenerator()
	postgresStore := postgres.NewPostgresStore(dbConn)
	generateService := generate.NewGenerateService(postgresStore, nanoIdGenerator, dragonFlyCache)
	gernerateHandler := generatehttp.NewGeneateHandler(generateService)
	router.POST("/shorten", gernerateHandler.Generate)

	retrieveService := retrieve.NewRetrieveService(postgresStore, dragonFlyCache)
	retrieveHandler := retrievehttp.NewRetrieveHandler(retrieveService)
	router.GET("/:shortUrl", retrieveHandler.Retrieve)

	if err := router.Run(":8080"); err != nil {
		log.Fatal().Msgf("failed to run server: %v", err)
	}

	router.GET("/urls/:shortUrl", func(c *gin.Context) {

	})

	router.DELETE("/urls/:shortUrl", func(c *gin.Context) {

	})
	fmt.Println("Server is running on port 8080", cfg)
}
