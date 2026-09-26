package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "tmp/docs"
	"tmp/internal/admin"
	"tmp/internal/cache"
	"tmp/internal/config"
	"tmp/internal/database"
	handlerHTTP "tmp/internal/handler/http"
	"tmp/internal/logger"
	"tmp/internal/middleware"
	"tmp/cmd/worker/notifier"
	"tmp/internal/repository"
	"tmp/internal/service"
)

// @title           CoworkGo API
// @version         1.0
// @description     REST API Cowork
// @host            localhost:8080
// @BasePath        /
func main() {
	log := logger.SetupLogger(appEnv())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("starting api server", slog.String("env", appEnv()), slog.String("port", cfg.Port))

	var cacheClient *cache.Client
	if cfg.RedisURL != "" {
		cacheClient, err = cache.New(cfg.RedisURL)
		if err != nil {
			log.Warn("redis cache disabled: invalid configuration", slog.String("error", err.Error()))
		} else if err := cacheClient.Ping(context.Background()); err != nil {
			log.Warn("redis cache disabled: connection failed", slog.String("error", err.Error()))
			_ = cacheClient.Close()
			cacheClient = nil
		} else {
			defer func() {
				if err := cacheClient.Close(); err != nil {
					log.Error("redis close failed", slog.String("error", err.Error()))
				}
			}()
			log.Info("redis cache connected")
		}
	}

	pool, err := database.Connect(cfg.DataBaseURL)
	if err != nil {
		log.Error("database connection failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("database connected successfully")

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWT_SECRET)
	authService.ConfigurePasswordReset(
		notifier.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom),
		cfg.WebBaseURL,
	)
	authHandler := handlerHTTP.NewAuthHandler(authService)

	roomRepo := repository.NewRoomRepo(pool)
	roomService := service.NewRoomService(roomRepo, cacheClient)
	roomHandler := handlerHTTP.NewRoomHanlder(roomService)
	deskHandler := handlerHTTP.NewDeskHandler(repository.NewDeskRepo(pool))

	officeRepo := repository.NewOfficeRepository(pool)
	officeService := service.NewOfficeService(officeRepo, cacheClient)
	officeHandler := handlerHTTP.NewOfficeHandler(officeService)

	bookingRepo := repository.NewBookingRepository(pool)
	bookingService := service.NewBookingService(bookingRepo, roomRepo)
	bookingHandler := handlerHTTP.NewBookingHandler(bookingService)

	adminRepo := admin.NewAdminRepository(pool)
	adminHandler := admin.NewAdminHandler(adminRepo)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.GinLogger(log))
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler.RegisterRoutes(router)
	roomHandler.RegisterRoutes(router)
	deskHandler.RegisterRoutes(router)
	officeHandler.RegisterRoutes(router)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWT_SECRET))
	bookingHandler.RegisterRoutes(protected)
	adminHandler.RegisterRoutes(protected)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Info("listening for http requests", slog.String("addr", ":"+cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	log.Info("server stopped successfully")
}

func appEnv() string {
	if env := os.Getenv("ENV"); env != "" {
		return env
	}
	return "local"
}
