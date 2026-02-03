package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/GunarsK-portfolio/portfolio-common/health"

	"github.com/GunarsK-rpg/public-api/internal/config"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// Setup configures all routes for the service.
func Setup(router *gin.Engine, handler *handlers.Handler, cfg *config.Config, healthAgg *health.Aggregator, pool *pgxpool.Pool) {
	// Public endpoints (no auth)
	router.GET("/health", healthAgg.Handler())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
