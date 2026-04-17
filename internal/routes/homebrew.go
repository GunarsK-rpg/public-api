package routes

import (
	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// registerHomebrewRoutes mounts the homebrew CRUD surface on the given group.
// Factored out of Setup for direct testability (routes_test.go).
func registerHomebrewRoutes(homebrew *gin.RouterGroup, handler *handlers.Handler) {
	hbRead := commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelRead)
	hbEdit := commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelEdit)
	hbDelete := commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelDelete)

	// Source books
	homebrew.GET("/source-books", hbRead, handler.ListMyHomebrewSourceBooks)
	homebrew.GET("/source-books/:code", hbRead, handler.GetSourceBookByCode)
	homebrew.POST("/source-books", hbEdit, handler.CreateSourceBook)
	homebrew.PUT("/source-books/:code", hbEdit, handler.UpdateSourceBook)
	homebrew.POST("/source-books/:code/restore", hbEdit, handler.RestoreSourceBook)
	homebrew.DELETE("/source-books/:code", hbDelete, handler.DeleteSourceBook)

	// Book-scoped classifiers
	homebrew.POST("/source-books/:code/:type", hbEdit, handler.UpsertBookClassifier)
	homebrew.PUT("/source-books/:code/:type/:cid", hbEdit, handler.UpsertBookClassifier)
	homebrew.POST("/source-books/:code/:type/:cid/restore", hbEdit, handler.RestoreBookClassifier)
	homebrew.DELETE("/source-books/:code/:type/:cid", hbDelete, handler.DeleteBookClassifier)
}
