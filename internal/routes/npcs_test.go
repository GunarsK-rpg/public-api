package routes

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

func TestRegisterNpcRoutes_MountsAllEndpoints(t *testing.T) {
	router := gin.New()
	registerNpcRoutes(router.Group("/npcs"), (*handlers.Handler)(nil))

	assertRoutesMounted(t, router, []string{
		"GET /npcs/:id",
	})
}

func TestRegisterNpcInstanceRoutes_MountsAllEndpoints(t *testing.T) {
	router := gin.New()
	registerNpcInstanceRoutes(router.Group("/npc-instances"), (*handlers.Handler)(nil))

	assertRoutesMounted(t, router, []string{
		"GET /npc-instances/:id",
		"POST /npc-instances",
		"PATCH /npc-instances/:id",
		"DELETE /npc-instances/:id",
	})
}
