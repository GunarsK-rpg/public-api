package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/GunarsK-portfolio/portfolio-common/health"
	"github.com/GunarsK-portfolio/portfolio-common/jwt"
	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/config"
	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
	"github.com/GunarsK-rpg/public-api/internal/middleware"
)

// Setup configures all routes for the service.
func Setup(router *gin.Engine, handler *handlers.Handler, cfg *config.Config, healthAgg *health.Aggregator, pool *pgxpool.Pool) {
	// Public endpoints (no auth)
	router.GET("/health", healthAgg.Handler())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup JWT validator and auth middleware
	jwtService, err := jwt.NewValidatorOnly(cfg.JWTSecret)
	if err != nil {
		log.Fatalf("failed to create JWT validator: %v", err)
	}
	authMiddleware := commonMiddleware.NewAuthMiddleware(jwtService)

	// API v1 group - all routes require authentication
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware.ValidateToken())
	v1.Use(authMiddleware.AddTTLHeader())
	v1.Use(middleware.UserSync(pool))

	// Classifiers routes (read-only)
	classifiers := v1.Group("/classifiers")
	classifiers.Use(commonMiddleware.RequirePermission(constants.ResourceClassifiers, commonMiddleware.LevelRead))
	{
		// Simple getters (no parameters)
		classifiers.GET("/attribute-types", handler.GetAttributeTypes)
		classifiers.GET("/attributes", handler.GetAttributes)
		classifiers.GET("/derived-stats", handler.GetDerivedStats)
		classifiers.GET("/derived-stat-values", handler.GetDerivedStatValues)
		classifiers.GET("/skills", handler.GetSkills)
		classifiers.GET("/expertise-types", handler.GetExpertiseTypes)
		classifiers.GET("/paths", handler.GetPaths)
		classifiers.GET("/surges", handler.GetSurges)
		classifiers.GET("/radiant-orders", handler.GetRadiantOrders)
		classifiers.GET("/ancestries", handler.GetAncestries)
		classifiers.GET("/activation-types", handler.GetActivationTypes)
		classifiers.GET("/action-types", handler.GetActionTypes)
		classifiers.GET("/damage-types", handler.GetDamageTypes)
		classifiers.GET("/units", handler.GetUnits)
		classifiers.GET("/equipment-types", handler.GetEquipmentTypes)
		classifiers.GET("/equipment-attributes", handler.GetEquipmentAttributes)
		classifiers.GET("/conditions", handler.GetConditions)
		classifiers.GET("/injuries", handler.GetInjuries)
		classifiers.GET("/goal-status", handler.GetGoalStatus)
		classifiers.GET("/connection-types", handler.GetConnectionTypes)
		classifiers.GET("/companion-types", handler.GetCompanionTypes)
		classifiers.GET("/cultures", handler.GetCultures)
		classifiers.GET("/tiers", handler.GetTiers)
		classifiers.GET("/levels", handler.GetLevels)
		classifiers.GET("/starting-kits", handler.GetStartingKits)

		// Getters with optional filters
		classifiers.GET("/expertises", handler.GetExpertises)    // ?type_code=
		classifiers.GET("/specialties", handler.GetSpecialties)  // ?path_code=
		classifiers.GET("/singer-forms", handler.GetSingerForms) // ?base_forms_only=
		classifiers.GET("/talents", handler.GetTalents)          // ?path_code=&specialty_code=&ancestry_code=&radiant_order_code=&surge_code=&is_key=
		classifiers.GET("/actions", handler.GetActions)          // ?action_type_code=&activation_type_code=&damage_type_code=
		classifiers.GET("/action-links", handler.GetActionLinks) // ?object_id= OR ?action_code= (one required)
		classifiers.GET("/equipments", handler.GetEquipments)    // ?type_code=&damage_type_code=
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
		registerHeroSubResource(heroes, "companions", handler.GetHeroCompanions, handler.UpsertHeroCompanion, handler.DeleteHeroCompanion)
		registerHeroSubResource(heroes, "cultures", handler.GetHeroCultures, handler.UpsertHeroCulture, handler.DeleteHeroCulture)
	}

	// Campaigns routes
	campaigns := v1.Group("/campaigns")
	campaigns.Use(commonMiddleware.RequirePermission(constants.ResourceCampaigns, commonMiddleware.LevelRead))

	// Suppress unused variable warnings until endpoints are added
	_ = campaigns
}

// registerHeroSubResource registers GET/POST/DELETE routes for a hero sub-resource.
func registerHeroSubResource(group *gin.RouterGroup, name string, getHandler, upsertHandler, deleteHandler gin.HandlerFunc) {
	group.GET("/:id/"+name, getHandler)
	group.POST("/:id/"+name, commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), upsertHandler)
	group.DELETE("/:id/"+name+"/:subId", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), deleteHandler)
}
