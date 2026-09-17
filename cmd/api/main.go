package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/config"
	"tmp/internal/database"
	"tmp/internal/handler/http"
	"tmp/internal/repository"
	"tmp/internal/service"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configs", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DataBaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo)
	authHandler := http.NewAuthHandler(authService)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Start succes code 200",
			"status_db": "connected",
		})
	})

	authHandler.RegisterRoutes(router)

	router.Run(":" + cfg.Port) // последняя строка!
}
