package routes

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

func TestRegisterHeroesRoutes_MountsAllEndpoints(t *testing.T) {
	router := gin.New()
	registerHeroesRoutes(router.Group("/heroes"), (*handlers.Handler)(nil))

	assertRoutesMounted(t, router, []string{
		// Core hero CRUD
		"GET /heroes",
		"GET /heroes/:id",
		"GET /heroes/:id/sheet",
		"POST /heroes",
		"PUT /heroes/:id",
		"DELETE /heroes/:id",
		// A representative slice of sub-resources (full coverage is in TestRegisterHeroSubResource_*)
		"GET /heroes/:id/attributes",
		"POST /heroes/:id/attributes",
		"DELETE /heroes/:id/attributes/:subId",
		"GET /heroes/:id/paths",
		"POST /heroes/:id/paths",
		"DELETE /heroes/:id/paths/:subId",
		// Non-CRUD endpoints
		"GET /heroes/:id/companions",
		"GET /heroes/:id/companion-npcs",
		"POST /heroes/:id/equipment/:subId/modifications",
		"DELETE /heroes/:id/equipment/:subId/modifications/:modId",
		"POST /heroes/:id/favorites",
		"DELETE /heroes/:id/favorites/:subId",
		"POST /heroes/:id/avatar",
		"DELETE /heroes/:id/avatar",
		"PATCH /heroes/:id/health",
		"PATCH /heroes/:id/focus",
		"PATCH /heroes/:id/magic",
		"PATCH /heroes/:id/currency",
	})
}

func TestRegisterHeroSubResource_RegistersCorrectRoutes(t *testing.T) {
	router := gin.New()
	group := router.Group("/heroes")

	noop := func(c *gin.Context) { c.Status(200) }
	registerHeroSubResource(group, "attributes", noop, noop, noop)

	assertRoutesMounted(t, router, []string{
		"GET /heroes/:id/attributes",
		"POST /heroes/:id/attributes",
		"DELETE /heroes/:id/attributes/:subId",
	})
}

func TestRegisterHeroSubResource_MultipleResources(t *testing.T) {
	router := gin.New()
	group := router.Group("/heroes")

	noop := func(c *gin.Context) { c.Status(200) }
	registerHeroSubResource(group, "skills", noop, noop, noop)
	registerHeroSubResource(group, "talents", noop, noop, noop)

	assertRoutesMounted(t, router, []string{
		"GET /heroes/:id/skills",
		"POST /heroes/:id/skills",
		"DELETE /heroes/:id/skills/:subId",
		"GET /heroes/:id/talents",
		"POST /heroes/:id/talents",
		"DELETE /heroes/:id/talents/:subId",
	})
}
