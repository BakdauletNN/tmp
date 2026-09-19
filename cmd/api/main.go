package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"tmp/internal/config"
	"tmp/internal/database"
	"tmp/internal/handler/http"
	"tmp/internal/middleware"
	"tmp/internal/repository"
	"tmp/internal/service"
	_ "tmp/docs"
)

// @title           CoworkGo API
// @version         1.0
// @description     REST API для бронирования коворкингов
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configs", err)
	}

	pool, err := database.Connect(cfg.DataBaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWT_SECRET)
	authHandler := http.NewAuthHandler(authService)

	roomRepo := repository.NewRoomRepo(pool)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := http.NewRoomHanlder(roomService)

	officeRepo := repository.NewOfficeRepository(pool)
	officeService := service.NewOfficeService(officeRepo)
	officeHandler := http.NewOfficeHandler(officeService)

	bookingRepo := repository.NewBookingRepository(pool)
	bookingService := service.NewBookingService(bookingRepo, roomRepo)
	bookingHandler := http.NewBookingHandler(bookingService)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Start succes code 200",
			"status_db": "connected",
		})
	})

	authHandler.RegisterRoutes(router)
	roomHandler.RegisterRoutes(router)
	officeHandler.RegisterRoutes(router)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWT_SECRET))
	bookingHandler.RegisterRoutes(protected)

	router.Run(":" + cfg.Port)
}