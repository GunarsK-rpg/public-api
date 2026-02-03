package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/GunarsK-rpg/public-api/internal/models/requests"
)

// Simple getters (no parameters)

func (h *Handler) GetAttributeTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetAttributeTypes)
}

func (h *Handler) GetAttributes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetAttributes)
}

func (h *Handler) GetDerivedStats(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetDerivedStats)
}

func (h *Handler) GetDerivedStatValues(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetDerivedStatValues)
}

func (h *Handler) GetSkills(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetSkills)
}

func (h *Handler) GetExpertiseTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetExpertiseTypes)
}

func (h *Handler) GetPaths(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetPaths)
}

func (h *Handler) GetSurges(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetSurges)
}

func (h *Handler) GetRadiantOrders(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetRadiantOrders)
}

func (h *Handler) GetAncestries(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetAncestries)
}

func (h *Handler) GetActivationTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetActivationTypes)
}

func (h *Handler) GetActionTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetActionTypes)
}

func (h *Handler) GetDamageTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetDamageTypes)
}

func (h *Handler) GetUnits(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetUnits)
}

func (h *Handler) GetEquipmentTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetEquipmentTypes)
}

func (h *Handler) GetEquipmentAttributes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetEquipmentAttributes)
}

func (h *Handler) GetConditions(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetConditions)
}

func (h *Handler) GetInjuries(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetInjuries)
}

func (h *Handler) GetGoalStatus(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetGoalStatus)
}

func (h *Handler) GetConnectionTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetConnectionTypes)
}

func (h *Handler) GetCompanionTypes(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetCompanionTypes)
}

func (h *Handler) GetCultures(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetCultures)
}

func (h *Handler) GetTiers(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetTiers)
}

func (h *Handler) GetLevels(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetLevels)
}

func (h *Handler) GetStartingKits(c *gin.Context) {
	handleJSONResponse(c, h.repo.GetStartingKits)
}

// Filtered getters (query params → JSONB filter)

func (h *Handler) GetExpertises(c *gin.Context) {
	handleFilteredResponse[requests.GetExpertisesQuery](c, h.repo.GetExpertises)
}

func (h *Handler) GetSpecialties(c *gin.Context) {
	handleFilteredResponse[requests.GetSpecialtiesQuery](c, h.repo.GetSpecialties)
}

func (h *Handler) GetSingerForms(c *gin.Context) {
	handleFilteredResponse[requests.GetSingerFormsQuery](c, h.repo.GetSingerForms)
}

func (h *Handler) GetTalents(c *gin.Context) {
	handleFilteredResponse[requests.GetTalentsQuery](c, h.repo.GetTalents)
}

func (h *Handler) GetActions(c *gin.Context) {
	handleFilteredResponse[requests.GetActionsQuery](c, h.repo.GetActions)
}

func (h *Handler) GetActionLinks(c *gin.Context) {
	handleFilteredResponse[requests.GetActionLinksQuery](c, h.repo.GetActionLinks)
}

func (h *Handler) GetEquipments(c *gin.Context) {
	handleFilteredResponse[requests.GetEquipmentsQuery](c, h.repo.GetEquipments)
}
