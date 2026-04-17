package routes

import (
	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/GunarsK-portfolio/portfolio-common/middleware"

	"github.com/GunarsK-rpg/public-api/internal/constants"
	"github.com/GunarsK-rpg/public-api/internal/handlers"
)

// registerHeroesRoutes mounts hero CRUD, sub-resources, equipment/favorite
// management, avatar handling, and resource patch endpoints.
func registerHeroesRoutes(heroes *gin.RouterGroup, handler *handlers.Handler) {
	heroes.Use(commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelRead))

	heroEdit := commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit)
	heroDelete := commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete)

	// Core hero CRUD
	heroes.GET("", handler.GetHeroes)
	heroes.GET("/:id", handler.GetHero)
	heroes.GET("/:id/sheet", handler.GetHeroSheet)
	heroes.POST("", heroEdit, handler.CreateHero)
	heroes.PUT("/:id", heroEdit, handler.UpdateHero)
	heroes.DELETE("/:id", heroDelete, handler.DeleteHero)

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
	registerHeroSubResource(heroes, "paths", handler.GetHeroPaths, handler.UpsertHeroPath, handler.DeleteHeroPath)

	// Companion queries (read-only, instances managed via /npc-instances)
	heroes.GET("/:id/companions", handler.GetHeroCompanions)
	heroes.GET("/:id/companion-npcs", handler.GetCompanionNpcOptions)

	// Equipment modification management
	heroes.POST("/:id/equipment/:subId/modifications", heroEdit, handler.AddEquipmentModification)
	heroes.DELETE("/:id/equipment/:subId/modifications/:modId", heroDelete, handler.RemoveEquipmentModification)

	// Favorite action management
	heroes.POST("/:id/favorites", heroEdit, handler.AddFavoriteAction)
	heroes.DELETE("/:id/favorites/:subId", heroDelete, handler.RemoveFavoriteAction)

	// Avatar
	heroes.POST("/:id/avatar", heroEdit, handler.SetHeroAvatar)
	heroes.DELETE("/:id/avatar", heroDelete, handler.DeleteHeroAvatar)

	// Resource patch routes
	heroes.PATCH("/:id/health", heroEdit, handler.PatchHeroHealth)
	heroes.PATCH("/:id/focus", heroEdit, handler.PatchHeroFocus)
	heroes.PATCH("/:id/magic", heroEdit, handler.PatchHeroMagic)
	heroes.PATCH("/:id/currency", heroEdit, handler.PatchHeroCurrency)
}

// registerHeroSubResource registers GET/POST/DELETE routes for a hero sub-resource.
func registerHeroSubResource(group *gin.RouterGroup, name string, getHandler, upsertHandler, deleteHandler gin.HandlerFunc) {
	group.GET("/:id/"+name, getHandler)
	group.POST("/:id/"+name, commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelEdit), upsertHandler)
	group.DELETE("/:id/"+name+"/:subId", commonMiddleware.RequirePermission(constants.ResourceHeroes, commonMiddleware.LevelDelete), deleteHandler)
}
