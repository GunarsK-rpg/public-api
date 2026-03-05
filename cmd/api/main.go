package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-portfolio/portfolio-common/health"
	"github.com/GunarsK-portfolio/portfolio-common/logger"
	"github.com/GunarsK-portfolio/portfolio-common/metrics"
	"github.com/GunarsK-portfolio/portfolio-common/middleware"
	"github.com/GunarsK-portfolio/portfolio-common/server"

	"github.com/GunarsK-rpg/public-api/internal/cache"
	"github.com/GunarsK-rpg/public-api/internal/config"
	"github.com/GunarsK-rpg/public-api/internal/database"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
	"github.com/GunarsK-rpg/public-api/internal/repository"
	"github.com/GunarsK-rpg/public-api/internal/routes"
)

func main() {
	cfg := config.Load()

	appLogger := logger.New(logger.Config{
		Level:       os.Getenv("LOG_LEVEL"),
		Format:      os.Getenv("LOG_FORMAT"),
		ServiceName: "public-api",
		AddSource:   os.Getenv("LOG_SOURCE") == "true",
	})

	metricsCollector := metrics.New(metrics.Config{
		ServiceName: "public-api",
		Namespace:   "rpg",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	redisClient := cache.NewRedisClient(cfg.Redis, cfg.Service.Environment)
	defer redisClient.Close()

	healthAgg := health.NewAggregator(3 * time.Second)
	healthAgg.Register(database.NewPgxChecker(pool))
	healthAgg.Register(health.NewRedisChecker(redisClient))

	appCache := cache.New(redisClient)
	repo := repository.New(pool)
	handler := handlers.New(repo, appCache)

	if cfg.Service.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(logger.Recovery(appLogger))
	router.Use(logger.RequestLogger(appLogger))
	router.Use(metricsCollector.Middleware())

	securityMiddleware := middleware.NewSecurityMiddleware(
		cfg.Service.AllowedOrigins,
		"GET,POST,PUT,DELETE,OPTIONS",
		"Content-Type,Authorization",
		true,
	)
	router.Use(securityMiddleware.Apply())

	routes.Setup(router, handler, cfg, healthAgg, pool)

	appLogger.Info("Public API ready", "port", cfg.Service.Port)
	serverCfg := server.DefaultConfig(strconv.Itoa(cfg.Service.Port))
	if err := server.Run(router, serverCfg, appLogger); err != nil {
		appLogger.Error("Server error", "error", err)
		os.Exit(1)
	}
}
