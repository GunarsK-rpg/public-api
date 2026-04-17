package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/GunarsK-portfolio/portfolio-common/health"
	"github.com/GunarsK-portfolio/portfolio-common/jwt"
	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/cache"
	"github.com/GunarsK-rpg/public-api/internal/config"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
	"github.com/GunarsK-rpg/public-api/internal/middleware"
)

// Setup configures all routes for the service.
func Setup(router *gin.Engine, handler *handlers.Handler, cfg *config.Config, healthAgg *health.Aggregator, pool *pgxpool.Pool, appCache *cache.Cache) error {
	// Public endpoints (no auth)
	router.GET("/health", healthAgg.Handler())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup JWT validator and auth middleware
	jwtService, err := jwt.NewValidatorOnly(cfg.JWTSecret)
	if err != nil {
		return fmt.Errorf("failed to create JWT validator: %w", err)
	}
	authMiddleware := commonMiddleware.NewAuthMiddleware(jwtService)

	// API v1 group - all routes require authentication
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware.ValidateToken())
	v1.Use(authMiddleware.AddTTLHeader())
	v1.Use(middleware.UserSync(pool, appCache))
	v1.Use(middleware.BodyLimit(cfg.MaxBodySize))

	registerClassifiersRoutes(v1.Group("/classifiers"), handler)
	registerHomebrewRoutes(v1.Group("/homebrew"), handler)
	registerHeroesRoutes(v1.Group("/heroes"), handler)
	registerCampaignsRoutes(v1.Group("/campaigns"), handler)
	registerNpcRoutes(v1.Group("/npcs"), handler)
	registerNpcInstanceRoutes(v1.Group("/npc-instances"), handler)

	return nil
}
