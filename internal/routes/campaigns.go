package routes

import (
	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// registerCampaignsRoutes mounts campaign CRUD, NPC template library, and combat endpoints.
func registerCampaignsRoutes(campaigns *gin.RouterGroup, handler *handlers.Handler) {
	campaigns.Use(commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelRead))

	campaignEdit := commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelEdit)
	campaignDelete := commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelDelete)

	// Core campaign CRUD
	campaigns.GET("", handler.GetCampaigns)
	campaigns.GET("/join/:code", handler.GetCampaignByCode)
	campaigns.GET("/:id", handler.GetCampaign)
	campaigns.POST("", campaignEdit, handler.CreateCampaign)
	campaigns.PUT("/:id", campaignEdit, handler.UpdateCampaign)
	campaigns.DELETE("/:id", campaignDelete, handler.DeleteCampaign)
	campaigns.DELETE("/:id/heroes/:hid", campaignDelete, handler.RemoveHeroFromCampaign)

	// NPC template library
	campaigns.GET("/:id/npcs", handler.GetNpcOptions)
	campaigns.GET("/:id/npcs/library", handler.GetNpcLibrary)
	campaigns.GET("/:id/npcs/:nid", handler.GetNpc)
	campaigns.POST("/:id/npcs", campaignEdit, handler.CreateNpc)
	campaigns.PUT("/:id/npcs/:nid", campaignEdit, handler.UpdateNpc)
	campaigns.DELETE("/:id/npcs/:nid", campaignDelete, handler.DeleteNpc)
	campaigns.POST("/:id/npcs/:nid/avatar", campaignEdit, handler.SetNpcAvatar)
	campaigns.DELETE("/:id/npcs/:nid/avatar", campaignDelete, handler.DeleteNpcAvatar)

	// Combat encounters
	campaigns.GET("/:id/combats", handler.GetCombats)
	campaigns.GET("/:id/combats/:cid", handler.GetCombat)
	campaigns.POST("/:id/combats", campaignEdit, handler.CreateCombat)
	campaigns.PUT("/:id/combats/:cid", campaignEdit, handler.UpdateCombat)
	campaigns.DELETE("/:id/combats/:cid", campaignDelete, handler.DeleteCombat)
	campaigns.POST("/:id/combats/:cid/end-round", campaignEdit, handler.EndCombatRound)
}
