package routes

import (
	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// registerClassifiersRoutes mounts the read-only classifier endpoints.
func registerClassifiersRoutes(classifiers *gin.RouterGroup, handler *handlers.Handler) {
	classifiers.Use(commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelRead))
	classifiers.GET("", handler.GetAllClassifiers)
	classifiers.GET("/source-books", handler.GetSourceBooks)
}
