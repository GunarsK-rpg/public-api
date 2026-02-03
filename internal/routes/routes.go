package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/GunarsK-portfolio/portfolio-common/health"
	"github.com/GunarsK-portfolio/portfolio-common/jwt"
	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/config"
	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
	"github.com/GunarsK-rpg/public-api/internal/middleware"
)

// Setup configures all routes for the service.
func Setup(router *gin.Engine, handler *handlers.Handler, cfg *config.Config, healthAgg *health.Aggregator, pool *pgxpool.Pool) {
	// Public endpoints (no auth)
	router.GET("/health", healthAgg.Handler())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup JWT validator and auth middleware
	jwtService, err := jwt.NewValidatorOnly(cfg.JWTSecret)
	if err != nil {
		log.Fatalf("failed to create JWT validator: %v", err)
	}
	authMiddleware := commonMiddleware.NewAuthMiddleware(jwtService)

	// API v1 group - all routes require authentication
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware.ValidateToken())
	v1.Use(authMiddleware.AddTTLHeader())
	v1.Use(middleware.UserSync(pool))

	// Classifiers routes (read-only)
	classifiers := v1.Group("/classifiers")
	classifiers.Use(commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelRead))

	// Heroes routes
	heroes := v1.Group("/heroes")

	// Campaigns routes
	campaigns := v1.Group("/campaigns")

	// Suppress unused variable warnings until endpoints are added
	_ = classifiers
	_ = heroes
	_ = campaigns
}
