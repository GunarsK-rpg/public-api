package routes

import (
	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// registerNpcRoutes mounts read-only NPC template lookup (auth enforced at DB level).
func registerNpcRoutes(npcs *gin.RouterGroup, handler *handlers.Handler) {
	npcs.GET("/:id", commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelRead), handler.GetNpcByID)
}

// registerNpcInstanceRoutes mounts the NPC instance CRUD surface used by
// both combat encounters and companion tracking. Auth is enforced at the
// DB layer (instances can belong to any campaign the caller can see).
func registerNpcInstanceRoutes(instances *gin.RouterGroup, handler *handlers.Handler) {
	instances.GET("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelRead), handler.GetNpcInstance)
	instances.POST("", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.CreateNpcInstance)
	instances.PATCH("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.PatchNpcInstance)
	instances.DELETE("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), handler.DeleteNpcInstance)
}
