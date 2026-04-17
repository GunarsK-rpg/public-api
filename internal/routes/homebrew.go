package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// registerHomebrewRoutes mounts the homebrew CRUD surface on the given group.
// Factored out of Setup for direct testability (routes_test.go).
//
// Hero-scoped routes are registered behind a 501 stub while Phase 4.1 is
// unwritten: the handler methods are wired on *Handler for Phase 4.1 to
// consume, but the routes themselves serve a clear "not implemented" response
// so no caller accidentally relies on untested behaviour.
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

	// Hero-scoped classifiers — gated behind 501 until Phase 4.1.
	heroStub := func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "hero-scoped homebrew CRUD is gated on Phase 4.1"})
	}
	homebrew.POST("/heroes/:id/:type", hbEdit, heroStub)
	homebrew.PUT("/heroes/:id/:type/:cid", hbEdit, heroStub)
	homebrew.POST("/heroes/:id/:type/:cid/restore", hbEdit, heroStub)
	homebrew.DELETE("/heroes/:id/:type/:cid", hbDelete, heroStub)
}
