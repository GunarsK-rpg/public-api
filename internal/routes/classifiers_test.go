package routes

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

func TestRegisterClassifiersRoutes_MountsAllEndpoints(t *testing.T) {
	router := gin.New()
	registerClassifiersRoutes(router.Group("/classifiers"), (*handlers.Handler)(nil))

	assertRoutesMounted(t, router, []string{
		"GET /classifiers",
		"GET /classifiers/source-books",
	})
}
