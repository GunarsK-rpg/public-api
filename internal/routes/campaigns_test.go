package routes

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

func TestRegisterCampaignsRoutes_MountsAllEndpoints(t *testing.T) {
	router := gin.New()
	registerCampaignsRoutes(router.Group("/campaigns"), (*handlers.Handler)(nil))

	assertRoutesMounted(t, router, []string{
		// Core campaign CRUD
		"GET /campaigns",
		"GET /campaigns/join/:code",
		"GET /campaigns/:id",
		"POST /campaigns",
		"PUT /campaigns/:id",
		"DELETE /campaigns/:id",
		"DELETE /campaigns/:id/heroes/:hid",
		// NPC template library
		"GET /campaigns/:id/npcs",
		"GET /campaigns/:id/npcs/library",
		"GET /campaigns/:id/npcs/:nid",
		"POST /campaigns/:id/npcs",
		"PUT /campaigns/:id/npcs/:nid",
		"DELETE /campaigns/:id/npcs/:nid",
		"POST /campaigns/:id/npcs/:nid/avatar",
		"DELETE /campaigns/:id/npcs/:nid/avatar",
		// Combat encounters
		"GET /campaigns/:id/combats",
		"GET /campaigns/:id/combats/:cid",
		"POST /campaigns/:id/combats",
		"PUT /campaigns/:id/combats/:cid",
		"DELETE /campaigns/:id/combats/:cid",
		"POST /campaigns/:id/combats/:cid/end-round",
	})
}
