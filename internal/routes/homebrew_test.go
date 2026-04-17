package routes

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

func TestRegisterHomebrewRoutes_MountsAllEndpoints(t *testing.T) {
	router := gin.New()
	registerHomebrewRoutes(router.Group("/homebrew"), (*handlers.Handler)(nil))

	assertRoutesMounted(t, router, []string{
		"GET /homebrew/source-books",
		"GET /homebrew/source-books/:code",
		"POST /homebrew/source-books",
		"PUT /homebrew/source-books/:code",
		"POST /homebrew/source-books/:code/restore",
		"DELETE /homebrew/source-books/:code",
		"POST /homebrew/source-books/:code/:type",
		"PUT /homebrew/source-books/:code/:type/:cid",
		"POST /homebrew/source-books/:code/:type/:cid/restore",
		"DELETE /homebrew/source-books/:code/:type/:cid",
	})
}
