package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/GunarsK-portfolio/portfolio-common/health"
	"github.com/GunarsK-portfolio/portfolio-common/jwt"
	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/cache"
	"github.com/GunarsK-rpg/public-api/internal/config"
	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
	"github.com/GunarsK-rpg/public-api/internal/middleware"
)

// Setup configures all routes for the service.
func Setup(router *gin.Engine, handler *handlers.Handler, cfg *config.Config, healthAgg *health.Aggregator, pool *pgxpool.Pool, appCache *cache.Cache) error {
	// Public endpoints (no auth)
	router.GET("/health", healthAgg.Handler())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup JWT validator and auth middleware
	jwtService, err := jwt.NewValidatorOnly(cfg.JWTSecret)
	if err != nil {
		return fmt.Errorf("failed to create JWT validator: %w", err)
	}
	authMiddleware := commonMiddleware.NewAuthMiddleware(jwtService)

	// API v1 group - all routes require authentication
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware.ValidateToken())
	v1.Use(authMiddleware.AddTTLHeader())
	v1.Use(middleware.UserSync(pool, appCache))
	v1.Use(middleware.BodyLimit(cfg.MaxBodySize))

	// Classifiers routes (read-only, single batch endpoint)
	classifiers := v1.Group("/classifiers")
	classifiers.Use(commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelRead))
	{
		classifiers.GET("", handler.GetAllClassifiers)
		classifiers.GET("/source-books", handler.GetSourceBooks)
	}

	// Heroes routes
	heroes := v1.Group("/heroes")
	heroes.Use(commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelRead))
	{
		// Core hero CRUD
		heroes.GET("", handler.GetHeroes)
		heroes.GET("/:id", handler.GetHero)
		heroes.GET("/:id/sheet", handler.GetHeroSheet)
		heroes.POST("", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.CreateHero)
		heroes.PUT("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.UpdateHero)
		heroes.DELETE("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), handler.DeleteHero)

		// Sub-resource routes
		registerHeroSubResource(heroes, "attributes", handler.GetHeroAttributes, handler.UpsertHeroAttribute, handler.DeleteHeroAttribute)
		registerHeroSubResource(heroes, "defenses", handler.GetHeroDefenses, handler.UpsertHeroDefense, handler.DeleteHeroDefense)
		registerHeroSubResource(heroes, "derived-stats", handler.GetHeroDerivedStats, handler.UpsertHeroDerivedStat, handler.DeleteHeroDerivedStat)
		registerHeroSubResource(heroes, "skills", handler.GetHeroSkills, handler.UpsertHeroSkill, handler.DeleteHeroSkill)
		registerHeroSubResource(heroes, "expertises", handler.GetHeroExpertises, handler.UpsertHeroExpertise, handler.DeleteHeroExpertise)
		registerHeroSubResource(heroes, "talents", handler.GetHeroTalents, handler.UpsertHeroTalent, handler.DeleteHeroTalent)
		registerHeroSubResource(heroes, "equipment", handler.GetHeroEquipment, handler.UpsertHeroEquipment, handler.DeleteHeroEquipment)
		registerHeroSubResource(heroes, "conditions", handler.GetHeroConditions, handler.UpsertHeroCondition, handler.DeleteHeroCondition)
		registerHeroSubResource(heroes, "injuries", handler.GetHeroInjuries, handler.UpsertHeroInjury, handler.DeleteHeroInjury)
		registerHeroSubResource(heroes, "goals", handler.GetHeroGoals, handler.UpsertHeroGoal, handler.DeleteHeroGoal)
		registerHeroSubResource(heroes, "connections", handler.GetHeroConnections, handler.UpsertHeroConnection, handler.DeleteHeroConnection)
		registerHeroSubResource(heroes, "notes", handler.GetHeroNotes, handler.UpsertHeroNote, handler.DeleteHeroNote)
		registerHeroSubResource(heroes, "cultures", handler.GetHeroCultures, handler.UpsertHeroCulture, handler.DeleteHeroCulture)

		// Companion queries (read-only, instances managed via /npc-instances)
		heroes.GET("/:id/companions", handler.GetHeroCompanions)
		heroes.GET("/:id/companion-npcs", handler.GetCompanionNpcOptions)

		// Equipment modification management
		heroes.POST("/:id/equipment/:subId/modifications", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.AddEquipmentModification)
		heroes.DELETE("/:id/equipment/:subId/modifications/:modId", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), handler.RemoveEquipmentModification)

		// Favorite action management
		heroes.POST("/:id/favorites", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.AddFavoriteAction)
		heroes.DELETE("/:id/favorites/:subId", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), handler.RemoveFavoriteAction)

		// Avatar
		heroes.POST("/:id/avatar", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.SetHeroAvatar)
		heroes.DELETE("/:id/avatar", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), handler.DeleteHeroAvatar)

		// Resource patch routes
		heroes.PATCH("/:id/health", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.PatchHeroHealth)
		heroes.PATCH("/:id/focus", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.PatchHeroFocus)
		heroes.PATCH("/:id/investiture", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.PatchHeroInvestiture)
		heroes.PATCH("/:id/currency", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.PatchHeroCurrency)
	}

	// Campaigns routes
	campaigns := v1.Group("/campaigns")
	campaigns.Use(commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelRead))
	{
		campaigns.GET("", handler.GetCampaigns)
		campaigns.GET("/join/:code", handler.GetCampaignByCode)
		campaigns.GET("/:id", handler.GetCampaign)
		campaigns.POST("", commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelEdit), handler.CreateCampaign)
		campaigns.PUT("/:id", commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelEdit), handler.UpdateCampaign)
		campaigns.DELETE("/:id", commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelDelete), handler.DeleteCampaign)
		campaigns.DELETE("/:id/heroes/:hid", commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelDelete), handler.RemoveHeroFromCampaign)

		// NPC template library
		combatEdit := commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelEdit)
		combatDelete := commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelDelete)

		campaigns.GET("/:id/npcs", handler.GetNpcOptions)
		campaigns.GET("/:id/npcs/library", handler.GetNpcLibrary)
		campaigns.GET("/:id/npcs/:nid", handler.GetNpc)
		campaigns.POST("/:id/npcs", combatEdit, handler.CreateNpc)
		campaigns.PUT("/:id/npcs/:nid", combatEdit, handler.UpdateNpc)
		campaigns.DELETE("/:id/npcs/:nid", combatDelete, handler.DeleteNpc)
		campaigns.POST("/:id/npcs/:nid/avatar", combatEdit, handler.SetNpcAvatar)
		campaigns.DELETE("/:id/npcs/:nid/avatar", combatDelete, handler.DeleteNpcAvatar)

		// Combat encounters
		campaigns.GET("/:id/combats", handler.GetCombats)
		campaigns.GET("/:id/combats/:cid", handler.GetCombat)
		campaigns.POST("/:id/combats", combatEdit, handler.CreateCombat)
		campaigns.PUT("/:id/combats/:cid", combatEdit, handler.UpdateCombat)
		campaigns.DELETE("/:id/combats/:cid", combatDelete, handler.DeleteCombat)
		campaigns.POST("/:id/combats/:cid/end-round", combatEdit, handler.EndCombatRound)
	}

	// NPC templates (direct access — auth handled at DB level)
	npcs := v1.Group("/npcs")
	{
		npcs.GET("/:id", commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelRead), handler.GetNpcByID)
	}

	// NPC instances (combat + companion — auth handled at DB level)
	instances := v1.Group("/npc-instances")
	{
		instances.GET("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelRead), handler.GetNpcInstance)
		instances.POST("", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.CreateNpcInstance)
		instances.PATCH("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), handler.PatchNpcInstance)
		instances.DELETE("/:id", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), handler.DeleteNpcInstance)
	}

	return nil
}

// registerHeroSubResource registers GET/POST/DELETE routes for a hero sub-resource.
func registerHeroSubResource(group *gin.RouterGroup, name string, getHandler, upsertHandler, deleteHandler gin.HandlerFunc) {
	group.GET("/:id/"+name, getHandler)
	group.POST("/:id/"+name, commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), upsertHandler)
	group.DELETE("/:id/"+name+"/:subId", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), deleteHandler)
}
